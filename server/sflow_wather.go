package main

import (
	"time"
	"fmt"
	"log"
)

// UDP パケットだけを取り出します
func FilterPacket_UDP(packetList []*FlowEthernetISO8023) ([]*FlowEthernetISO8023) {
	result := []*FlowEthernetISO8023{}
	for i := 0; i < len(packetList); i++ {
		packet := packetList[i]
		if packet.L4Header == nil || packet.L4Header.Udp == nil {
			continue
		}
		result = append(result, packet)
	}
	return result
}

// 同じポート番号でパケットを集めます
type PacketPortMap map[uint16][]*FlowEthernetISO8023

// 同じポート番号にパケットを集めて map にして返します
func CountPort(packetArray []*FlowEthernetISO8023) PacketPortMap {
	counter := PacketPortMap{}
	for i := 0; i < len(packetArray); i++ {
		packet := packetArray[i]
		if packet.L4Header == nil {
			continue
		}
		if packet.L4Header.Tcp != nil {
			counter[packet.L4Header.Tcp.SourcePort] =
				append(counter[packet.L4Header.Tcp.SourcePort], packet)
			counter[packet.L4Header.Tcp.DestinationPort] =
				append(counter[packet.L4Header.Tcp.DestinationPort], packet)
		}
		if packet.L4Header.Udp != nil {
			counter[packet.L4Header.Udp.SourcePort] =
				append(counter[packet.L4Header.Udp.SourcePort], packet)
			counter[packet.L4Header.Udp.DestinationPort] =
				append(counter[packet.L4Header.Udp.DestinationPort], packet)
		}
	}
	return counter
}

func CalcMaxCountPort(counter PacketPortMap) (uint16, []*FlowEthernetISO8023) {
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

func PortWatcher(stopChannel chan bool,
	durationSecond int,
	watchPortList []int,
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
		udpPacketArray := FilterPacket_UDP(packetList)
		if ( len(udpPacketArray) <= 10 ) { // パケットが10個以下なら何もしないことにします
			continue
		}

		portCount := CountPort(udpPacketArray)
		maxPort, maxPacketList := CalcMaxCountPort(portCount)

		hit := false
		for i := 0; i < len(watchPortList); i++ {
			if (int(maxPort) == watchPortList[i]) {
				hit = true
				break
			}
		}
		if hit {
			err := Mew(fmt.Sprintf("port: %d is max count now. count: %d", maxPort, len(maxPacketList)), mewUserData)
			if err != nil {
				fmt.Println("mew error: ", err)
			}
			alertBuffer.AddPacketList(maxPort, maxPacketList)
		}
	}
}

