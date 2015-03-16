package main

import (
	"time"
	"fmt"
	"log"
)

// UDP パケットだけを取り出します
func filterPacketUDP(packetList []*FlowEthernetISO8023) ([]*FlowEthernetISO8023) {
	result := []*FlowEthernetISO8023{}
	for i := 0; i < len(packetList); i++ {
		packet := packetList[i]
		if packet.L4Header == nil || packet.L4Header.UDP == nil {
			continue
		}
		result = append(result, packet)
	}
	return result
}

// 同じポート番号でパケットを集めます
type packetPortMap map[uint16][]*FlowEthernetISO8023

// 同じポート番号にパケットを集めて map にして返します
func countPort(packetArray []*FlowEthernetISO8023) packetPortMap {
	counter := packetPortMap{}
	for i := 0; i < len(packetArray); i++ {
		packet := packetArray[i]
		if packet.L4Header == nil {
			continue
		}
		if packet.L4Header.TCP != nil {
			counter[packet.L4Header.TCP.SourcePort] =
				append(counter[packet.L4Header.TCP.SourcePort], packet)
			counter[packet.L4Header.TCP.DestinationPort] =
				append(counter[packet.L4Header.TCP.DestinationPort], packet)
		}
		if packet.L4Header.UDP != nil {
			counter[packet.L4Header.UDP.SourcePort] =
				append(counter[packet.L4Header.UDP.SourcePort], packet)
			counter[packet.L4Header.UDP.DestinationPort] =
				append(counter[packet.L4Header.UDP.DestinationPort], packet)
		}
	}
	return counter
}

func calcMaxCountPort(counter packetPortMap) (uint16, []*FlowEthernetISO8023) {
	var maxPort uint16
	var maxPacketList []*FlowEthernetISO8023
	maxPort = 0
	for port := range counter {
		if len(maxPacketList) < len(counter[port]) {
			maxPacketList = counter[port]
			maxPort = port
		}
	}
	return maxPort, maxPacketList
}

// PortWatcher は sflow のフィードを受けてデータを溜め込みます
func PortWatcher(stopChannel chan bool,
	durationSecond int,
	watchPortList []int,
	watchPortRangeList []PortRange,
	l3HeaderBuffer *L3HeaderBuffer,
	alertBuffer *AlertBuffer,
	mewUserData *MewUserData) {
	for {
		select {
		case <- stopChannel:
			return
		case <- time.After(time.Duration(durationSecond) * time.Second):
		}

		packetList, err := l3HeaderBuffer.GetDataFromTime(time.Now().Add(time.Duration(durationSecond * -1) * time.Second))
		if err != nil {
			log.Println("l3HeaderBuffer.GetDataFromTime() failed.")
			continue
		}
		udpPacketArray := filterPacketUDP(packetList)
		if ( len(udpPacketArray) <= 10 ) { // パケットが10個以下なら何もしないことにします
			continue
		}

		portCount := countPort(udpPacketArray)
		maxPort, maxPacketList := calcMaxCountPort(portCount)

		hit := false
		maxPortInt := int(maxPort)
		for i := 0; i < len(watchPortList); i++ {
			if (maxPortInt == watchPortList[i]) {
				hit = true
				break
			}
		}
		for i := 0; i < len(watchPortRangeList) && hit != true; i++ {
			portRange := watchPortRangeList[i]
			if (maxPortInt >= portRange.Low && maxPortInt <= portRange.High) {
				hit = true
				break
			}
		}
		if hit {
			mewResult, err := Mew(fmt.Sprintf("port: %d is max count now. count: %d", maxPort, len(maxPacketList)), mewUserData)
			if err != nil {
				fmt.Println("mew error: ", err)
			}
			alertBuffer.AddPacketList(maxPort, maxPacketList, mewResult)
		}
	}
}

