package main

import (
	"bytes"
	"log"
	"net"
	"fmt"
	"io"
	"encoding/hex"
	"encoding/binary"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"

	"github.com/PreetamJinka/sflow"
	//"github.com/PreetamJinka/udpchan"
	//"github.com/limura/sflow"
)

// UDP で listen します。
// addrAndPort は 0.0.0.0:12345 といったような文字列を指定します。
func listenUDP(addrAndPort string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", addrAndPort)
	if err != nil {
		return nil, err
	}
	return net.ListenUDP("udp", addr)
}

// sflow を addrAndPort で待ち受けます。取得した datagram を ch に書き込みます
func listenSflow(addrAndPort string, datagramChannel chan *sflow.Datagram, closeChannel chan error) {
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
		decorder := sflow.NewDecoder(bytes.NewReader(buf[:length]))
		
		dgram, err := decorder.Decode()
		if err != nil {
			log.Println("decorder.Decode() error: ", err)
			continue
		}
		datagramChannel <- dgram
	}
}

// json にする時に残すものだけを頭文字を大文字にしています
type L4HeaderTCP struct {
	SourcePort       uint16
	DestinationPort  uint16
	sequenceNumber   uint32
	acknowledgementNumber uint32
        // unused        uint8
	flags            uint8
	window           uint16
	checksum         uint16
	urgentPointer    uint16
}

func (f L4HeaderTCP) String() string {
	return fmt.Sprint(
		"srcPort: ", f.SourcePort,
		", dstPort: ", f.DestinationPort,
		", seqNum: ", f.sequenceNumber,
		", ackNum: ", f.acknowledgementNumber,
		fmt.Sprintf(", flags: 0x%x", f.flags),
		", window: ", f.window,
		fmt.Sprintf(", checksum 0x%x", f.checksum),
		"")
}

// json にする時に残すものだけを頭文字を大文字にしています
type L4HeaderUDP struct {
	SourcePort      uint16
	DestinationPort uint16
	length          uint16
	checksum        uint16
}

func (f L4HeaderUDP) String() string {
	return fmt.Sprint(
		"srcPort: ", f.SourcePort,
		", dstPort: ", f.DestinationPort,
		", length: ", f.length,
		fmt.Sprintf(", checksum 0x%x", f.checksum),
		"")
}

type L4HeaderICMP struct {
	Type     uint8
	Code     uint8
	checksum uint16
}

func (f L4HeaderICMP) String() string {
	return fmt.Sprint(
		"type: ", f.Type,
		", code: ", f.Code,
		", checksum: ", f.checksum,
		"")
}

type L4Header struct {
	Tcp  *L4HeaderTCP
	Udp  *L4HeaderUDP
	Icmp *L4HeaderICMP
}

func (f L4Header) String() string {
	return fmt.Sprint(
		"TCP: ", f.Tcp,
		", UDP: ", f.Udp,
		", ICMP: ", f.Icmp,
		"")
}

type FlowEthernetISO8023 struct {
	destinationMacAddress [6]byte
	sourceMacAddress [6]byte
	vlanID uint16
	vlanPriority uint8
	IPv4Header *ipv4.Header
	IPv6Header *ipv6.Header
	L4Header *L4Header
}

func (f FlowEthernetISO8023) String() string {
	return fmt.Sprint(
		"dstMac: ", hex.EncodeToString(f.destinationMacAddress[:]),
		", srcMac: ", hex.EncodeToString(f.sourceMacAddress[:]),
		", vlanID: ", f.vlanID,
		", vlanPriority: ", f.vlanPriority,
		", ipv4: ", f.IPv4Header,
		", ipv6: ", f.IPv6Header,
		", l4header: ", f.L4Header,
		"")
}

func parseL4HeaderTCP(buf []byte) (*L4HeaderTCP, error) {
	reader := bytes.NewReader(buf)
	f := new(L4HeaderTCP)
	var err error

	err = binary.Read(reader, binary.BigEndian, &f.SourcePort)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.DestinationPort)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.sequenceNumber)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.acknowledgementNumber)
	if err != nil {
		return nil, err
	}

	// unused な領域なのでとりあえずサイズの同じ Flags にダミーで読み込みます
	err = binary.Read(reader, binary.BigEndian, &f.flags)
	if err != nil {
		return nil, err
	}

	// ということで、こちらで読んでいるのが正解
	err = binary.Read(reader, binary.BigEndian, &f.flags)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.window)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.checksum)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.urgentPointer)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func parseL4HeaderUDP(buf []byte) (*L4HeaderUDP, error) {
	reader := bytes.NewReader(buf)
	f := new(L4HeaderUDP)
	var err error

	err = binary.Read(reader, binary.BigEndian, &f.SourcePort)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.DestinationPort)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.length)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.checksum)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func parseL4HeaderICMP(buf []byte) (*L4HeaderICMP, error) {
	reader := bytes.NewReader(buf)
	f := new(L4HeaderICMP)
	var err error

	err = binary.Read(reader, binary.BigEndian, &f.Type)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.Code)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.checksum)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// protocol から L4 header を解析します。
func parseL4Header(buf []byte, protocol int) (*L4Header, error) {
	var err error
	f := new(L4Header)

	switch protocol {
	case 1: // ICMP
		f.Icmp, err = parseL4HeaderICMP(buf)
		if err != nil {
			return nil, err
		}
	case 6: // TCP
		f.Tcp, err = parseL4HeaderTCP(buf)
		if err != nil {
			return nil, err
		}
	case 17: // UDP
		f.Udp, err = parseL4HeaderUDP(buf)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

// sflow.RawPacketFlow.Header が EthernetISO8023 のつもりで decode します
// sflowtool.c の  decodeinkLayer() の丸パクリです
func processRawPacketFlowEthernetISO8023(rawPacketFlow sflow.RawPacketFlow) (FlowEthernetISO8023, error) {
	reader := bytes.NewReader(rawPacketFlow.Header)
	f := FlowEthernetISO8023{}
	var err error

	err = binary.Read(reader, binary.BigEndian, &f.destinationMacAddress)
	if err != nil {
		return f, err
	}

	err = binary.Read(reader, binary.BigEndian, &f.sourceMacAddress)
	if err != nil {
		return f, err
	}

	var type_len uint16
	err = binary.Read(reader, binary.BigEndian, &type_len)
	if err != nil {
		return f, err
	}

	if (type_len == 0x8100) {
		var vlanData uint16
		err := binary.Read(reader, binary.BigEndian, &vlanData)
		if err != nil {
			return f, err
		}
		f.vlanID = vlanData & 0x0fff;
		f.vlanPriority = (uint8)(vlanData >> 13);

		// vlan を読んだら type_len をもう一回読む必要があるっぽい
		err = binary.Read(reader, binary.BigEndian, &type_len)
		if err != nil {
			return f, err
		}
	}else{
		f.vlanID = 0
		f.vlanPriority = 0
	}

	// 802.3+802.2 かどうかを観ないと駄目っぽい
	if (type_len < 1500) { // 1500 = NFT_MAX_8023_LEN
		var tmpData [3]byte
		err = binary.Read(reader, binary.BigEndian, &tmpData)
		if err != nil {
			return f, err
		}
		if (tmpData[0] == 0xAA &&
		    tmpData[1] == 0xAA &&
		    tmpData[2] == 0x03 ) {
			err = binary.Read(reader, binary.BigEndian, &tmpData)
			if err != nil {
				return f, err
			}
			if (tmpData[0] != 0x00 &&
			    tmpData[1] != 0x00 &&
			    tmpData[2] != 0x00 ) {
				return f, fmt.Errorf("invalid header for 802.3+802.2")
			}
			// もう一回 type_len を読みます	
			err = binary.Read(reader, binary.BigEndian, &type_len)
			if err != nil {
				return f, err
			}
		} else {
			var tmpData [3]byte
			err := binary.Read(reader, binary.BigEndian, &tmpData)
			if err != nil {
				return f, err
			}
			if (tmpData[0] == 0x06 &&
		    	tmpData[1] == 0x06 &&
		    	(tmpData[2] & 0x01) != 0x0 ) {
				// IP over 8022
				type_len = 0x0800
			}else{
				return f, nil
			}
		}
	}

	if (type_len == 0x0800) {
		// IPv4
		buf := make([]byte, reader.Len())
		err := binary.Read(reader, binary.BigEndian, &buf)
		if err != nil {
			return f, err
		}
		f.IPv4Header, err = ipv4.ParseHeader(buf)
		if err != nil || f.IPv4Header == nil {
			return f, err
		}
		if (len(buf) > f.IPv4Header.Len) {
			l4Header := buf[f.IPv4Header.Len:]
			f.L4Header, err = parseL4Header(l4Header, f.IPv4Header.Protocol)
			if err != nil {
				f.L4Header = nil
			}
		}
	}
	if (type_len == 0x86DD) {
		// IPv6
		buf := make([]byte, reader.Len())
		err := binary.Read(reader, binary.BigEndian, &buf)
		if err != nil {
			return f, err
		}
		f.IPv6Header, err = ipv6.ParseHeader(buf)
		if err != nil || f.IPv6Header == nil {
			return f, err
		}
		// TODO: XXX nextHeader を確認せずに、いきなり parseL4Header に入れています。
		// 2つ目の ipv6 header とか知らない！
		if (len(buf) > 40) {
			l4Header := buf[40:]
			f.L4Header, err = parseL4Header(l4Header, f.IPv6Header.NextHeader)
			if err != nil {
				f.L4Header = nil
			}
		}
	}

	return f, nil
}

const (
  SFLHEADER_ETHERNET_ISO8023     = 1
  SFLHEADER_ISO88024_TOKENBUS    = 2
  SFLHEADER_ISO88025_TOKENRING   = 3
  SFLHEADER_FDDI                 = 4
  SFLHEADER_FRAME_RELAY          = 5
  SFLHEADER_X25                  = 6
  SFLHEADER_PPP                  = 7
  SFLHEADER_SMDS                 = 8
  SFLHEADER_AAL5                 = 9
  SFLHEADER_AAL5_IP              = 10 /* e.g. Cisco AAL5 mux */
  SFLHEADER_IPv4                 = 11
  SFLHEADER_IPv6                 = 12
  SFLHEADER_MPLS                 = 13
  SFLHEADER_POS                  = 14
  SFLHEADER_IEEE80211MAC         = 15
  SFLHEADER_IEEE80211_AMPDU      = 16
  SFLHEADER_IEEE80211_AMSDU_SUBFRAME = 17
)

func processRawPacketFlow(rawPacketFlow sflow.RawPacketFlow, headerBuffer *L3HeaderBuffer) {
	switch rawPacketFlow.Protocol {
	case SFLHEADER_ETHERNET_ISO8023:
		etherFlow, err := processRawPacketFlowEthernetISO8023(rawPacketFlow)
		if err != nil {
			log.Println("  Raw EthernetISO8023 decode error: ", err)
		}else{
			headerBuffer.AddL3Header(&etherFlow)
			//log.Println("  Raw EthernetISO8023: ", etherFlow)
		}
	default:
		log.Println("  unknwon protocol: ", rawPacketFlow.Protocol)
	}
}

func processFlowSample(sample sflow.Sample, headerBuffer *L3HeaderBuffer) {
	records := sample.GetRecords()
	for i := range records {
		record := records[i]
		switch record.RecordType() {
		case sflow.TypeRawPacketFlowRecord:
			rawPacketFlow, ok := record.(sflow.RawPacketFlow)
			if !ok {
				continue
			}
			/*
			log.Println("  raw packet flow. proto: ", rawPacketFlow.Protocol,
				", HeaderSize: ", rawPacketFlow.HeaderSize,
				", Header: ", hex.EncodeToString(rawPacketFlow.Header))
			*/
			processRawPacketFlow(rawPacketFlow, headerBuffer)
		case sflow.TypeEthernetFrameFlowRecord:
			etherPacketFlow, ok := record.(sflow.ExtendedSwitchFlow)
			if !ok {
				continue
			}
			log.Println("  ethernet frame packet flow. sourceVlan: ", etherPacketFlow.SourceVlan)
		default:
			//log.Println("  unknwon type packet flow.", record.RecordType())
		}
	}
}

func startSflowCollector(addrAndPort string, headerBuffer *L3HeaderBuffer) error {
	datagramChannel := make(chan *sflow.Datagram)
	listenStopChannel := make(chan error)
	go listenSflow(addrAndPort, datagramChannel, listenStopChannel)
	for {
		select {
		case closeError := <- listenStopChannel:
			fmt.Println("listen sflow stoped")
			return closeError
			break
		default:
		}
		datagram := <- datagramChannel
		//fmt.Println("sflow packet got: ip: ", datagram.IpAddress, " subAgentID ", datagram.SubAgentId);
		i := 0
		for _, sample := range datagram.Samples {
			i += 1
			switch sample.SampleType() {
			case sflow.TypeFlowSample:
				//log.Println("Flow Sample", i)
				processFlowSample(sample, headerBuffer)
				break
			case sflow.TypeCounterSample:
				//log.Println("Counter Sample", i)
				//processCounterSample(sample)
				break
			case sflow.TypeExpandedFlowSample:
				log.Println("Expanded Flow Sample", i)
				break
			case sflow.TypeExpandedCounterSample:
				//log.Println("Expanded Counter Sample", i)
				break
			default:
				log.Println("  Unknown")
				break
			}
		}
	}
}


