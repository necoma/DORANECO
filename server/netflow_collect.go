package main

import (
	"log"
	"fmt"
	"io"
	"bytes"
	"github.com/fln/nf9packet"
	"time"
	"strconv"
	//"encoding/json"
	//"encoding/binary"
)

// NetFlowTemplateMap は NetFlow の template を保持するものです。
// key となる値が何であるかの情報を常に更新し続けます。
type NetFlowTemplateMap map[int]nf9packet.TemplateRecord

// listenNetFlow は
// sflow を addrAndPort で待ち受けます。取得した datagram を ch に書き込みます
func listenNetFlow(addrAndPort string, packetChannel chan *nf9packet.Packet, closeChannel chan error) {
	mtu := 10000
	udpConn, err := listenUDP(addrAndPort)
	if err != nil {
		closeChannel <- nil
		return
	}
	defer udpConn.Close()
	for {
		select {
		case <- closeChannel:
			fmt.Println("close msg got. close UDP listen socket")
			break
		default:
		}

		buf := make([]byte, mtu)
		length, _, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				fmt.Println("udp socket EOF.")
				break
			}
			log.Println("ReadFromUDP error: ", err)
			continue
		}
		netFlowPacket, err := nf9packet.Decode(buf[:length])
		if err != nil {
			fmt.Println("nf9packet.Decode() error: ", err)
			continue
		}

		packetChannel <- netFlowPacket
	}
}

// UpdateTemplate は NetFlowTemplateMap の　template を更新します
func (templateMap *NetFlowTemplateMap) UpdateTemplate(flowSet *nf9packet.TemplateFlowSet) {
	for _, record := range flowSet.Records {
		(*templateMap)[int(record.TemplateId)] = record
	}
}

// ConvertFlowSetToJSON は
// flowSet を TemplateMap にある情報から JSON でそのまま使えるような文字列に変換します。
func (templateMap *NetFlowTemplateMap) ConvertFlowSetToJSON(flowSet *nf9packet.DataFlowSet) (string, error) {
	var retBuffer bytes.Buffer

	if templateMap == nil || flowSet == nil {
		return "", fmt.Errorf("nil input")
	}
	template, ok := (*templateMap)[int(flowSet.Id)]
	if !ok {
		return "", fmt.Errorf("template id %v not defined now", flowSet.Id)
	}
	buffer := bytes.NewBuffer(flowSet.Data)
	
	for _, field := range template.Fields {
		if buffer.Len() < int(field.Length) {
			break
		}
		currentBuffer := buffer.Next(int(field.Length))
		if retBuffer.Len() > 0 {
			retBuffer.WriteString(", ")
		}
		// 数字とかは数字のままにしたかったので怪しく source をかっぱらってカスタムした DB を作った
		fieldConverter, ok := fieldDb[field.Type]
		if ok {
			retBuffer.WriteString(fmt.Sprintf("\"%v\": %v",
				field.Name(),
				fieldConverter.String(currentBuffer)))
		}else{
			retBuffer.WriteString(fmt.Sprintf("\"%v\": \"%v\"",
				field.Name(),
				field.DataToString(currentBuffer)))
		}
	}
	return fmt.Sprintf("{%v}", retBuffer.String()), nil
}

// checkAlertData は
/// flowSet が alert を上げるべきデータであるかどうかを判定して、必要なら alertBuffer に追加します
func (templateMap *NetFlowTemplateMap) checkAlertData(flowSet *nf9packet.DataFlowSet, alertTargetMap *NetFlowAlertTargetMap, alertBuffer *NetFlowAlertBuffer, jsonString string) error {
	if templateMap == nil || flowSet == nil {
		return fmt.Errorf("nil input")
	}
	template, ok := (*templateMap)[int(flowSet.Id)]
	if !ok {
		return fmt.Errorf("template id %v not defined now", flowSet.Id)
	}
	buffer := bytes.NewBuffer(flowSet.Data)

	srcPort := ""
	dstPort := ""
	var firstSwitched time.Duration
	var lastSwitched time.Duration
	for _, field := range template.Fields {
		if buffer.Len() < int(field.Length) {
			break
		}
		currentBuffer := buffer.Next(int(field.Length))
		if field.Name() == "L4_SRC_PORT" {
			srcPort = field.DataToString(currentBuffer)
		}
		if field.Name() == "L4_DST_PORT" {
			dstPort = field.DataToString(currentBuffer)
		}
		if field.Name() == "FIRST_SWITCHED" {
			firstSwitched = time.Duration(fieldToUInteger(currentBuffer)) * time.Millisecond
		}
		if field.Name() == "LAST_SWITCHED" {
			lastSwitched = time.Duration(fieldToUInteger(currentBuffer)) * time.Millisecond
		}
	}

	duration := lastSwitched - firstSwitched
	if duration <= 0 {
		return nil
	}
	if srcPort != "" {
		targetDurationFloat, ok := (*alertTargetMap)[srcPort]
		targetDuration := time.Duration(targetDurationFloat) * time.Second
		if ok && targetDuration <= duration {
			intValue, ok := strconv.Atoi(srcPort)
			if ok == nil {
				alertBuffer.AddAlert(intValue, duration, jsonString)
			}
		}
	}
	if dstPort != "" {
		targetDurationFloat, ok := (*alertTargetMap)[dstPort]
		targetDuration := time.Duration(targetDurationFloat) * time.Second
		if ok && targetDuration <= duration {
			intValue, ok := strconv.Atoi(dstPort)
			if ok == nil {
				alertBuffer.AddAlert(intValue, duration, jsonString)
			}
		}
	}
	return nil
}

// startNetFlowCollector は、NetFlowCollector の動作を開始する点です
func startNetFlowCollector(addrAndPort string, netFlowBuffer *NetFlowBuffer, netFlowAlertBuffer *NetFlowAlertBuffer, netFlowAlertTargetMap *NetFlowAlertTargetMap) error {
	packetChannel := make(chan *nf9packet.Packet)
	listenStopChannel := make(chan error)
	templateMap := make(NetFlowTemplateMap)
	go listenNetFlow(addrAndPort, packetChannel, listenStopChannel)
	for {
		select {
		case closeError := <- listenStopChannel:
			fmt.Println("listen sflow stoped")
			return closeError
			break
		default:
		}
		netFlowPacket := <- packetChannel

		for _, flowSet := range netFlowPacket.FlowSets {
			switch set := flowSet.(type) {
			case nf9packet.TemplateFlowSet:
				templateMap.UpdateTemplate(&set)
				break
			case nf9packet.DataFlowSet:
				jsonString, err := templateMap.ConvertFlowSetToJSON(&set)
				if err != nil {
					break
				}
				netFlowBuffer.AddFlowJSONString(jsonString)
				templateMap.checkAlertData(&set, netFlowAlertTargetMap, netFlowAlertBuffer, jsonString)
				break
			case nf9packet.OptionsTemplateFlowSet:
				fmt.Println("WARN: OptionsTemplateFlowSet got.")
				break
			}
		}
	}
}


