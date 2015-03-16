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

// NetFlowTemplateMap $B$O(B NetFlow $B$N(B template $B$rJ];}$9$k$b$N$G$9!#(B
// key $B$H$J$kCM$,2?$G$"$k$+$N>pJs$r>o$K99?7$7B3$1$^$9!#(B
type NetFlowTemplateMap map[int]nf9packet.TemplateRecord

// listenNetFlow $B$O(B
// sflow $B$r(B addrAndPort $B$GBT$A<u$1$^$9!#<hF@$7$?(B datagram $B$r(B ch $B$K=q$-9~$_$^$9(B
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

// UpdateTemplate $B$O(B NetFlowTemplateMap $B$N!!(Btemplate $B$r99?7$7$^$9(B
func (templateMap *NetFlowTemplateMap) UpdateTemplate(flowSet *nf9packet.TemplateFlowSet) {
	for _, record := range flowSet.Records {
		(*templateMap)[int(record.TemplateId)] = record
	}
}

// ConvertFlowSetToJSON $B$O(B
// flowSet $B$r(B TemplateMap $B$K$"$k>pJs$+$i(B JSON $B$G$=$N$^$^;H$($k$h$&$JJ8;zNs$KJQ49$7$^$9!#(B
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
		// $B?t;z$H$+$O?t;z$N$^$^$K$7$?$+$C$?$N$G2x$7$/(B source $B$r$+$C$Q$i$C$F%+%9%?%`$7$?(B DB $B$r:n$C$?(B
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

// checkAlertData $B$O(B
/// flowSet $B$,(B alert $B$r>e$2$k$Y$-%G!<%?$G$"$k$+$I$&$+$rH=Dj$7$F!"I,MW$J$i(B alertBuffer $B$KDI2C$7$^$9(B
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

// startNetFlowCollector $B$O!"(BNetFlowCollector $B$NF0:n$r3+;O$9$kE@$G$9(B
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


