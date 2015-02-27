package main

import (
	"time"
	"fmt"
	//"strconv"
	"container/list"
)

type SrcDstCountData struct {
	SourceCount      map[string]int  `json:"Source"`
	DestinationCount map[string]int  `json:"Destination"`
}

type L4CountData struct {
	Tcp   SrcDstCountData  `json:"tcp"`
	Udp   SrcDstCountData  `json:"udp"`
	Icmp  SrcDstCountData  `json:"icmp"`
}

type L3HeaderWithTime struct {
	time     time.Time
	l3Header *FlowEthernetISO8023
}

type CountDataBuffer struct {
	MaxLength     int
	CountDataList *list.List
}

type L3HeaderBuffer struct {
	MaxLength    int
	L3HeaderList *list.List
}

//
func MakeL3HeaderBuffer(maxLength int) *L3HeaderBuffer{
	return &L3HeaderBuffer{ MaxLength: maxLength,
		L3HeaderList: list.New() }
}

// L3HeaderWithTime を追加します
func (buf L3HeaderBuffer) AddL3HeaderWithTime(data *L3HeaderWithTime) error {
	if ( data == nil ) {
		return fmt.Errorf("nil input")
	}

	if buf.L3HeaderList.Len() >= buf.MaxLength {
		buf.L3HeaderList.Remove(buf.L3HeaderList.Front())
	}

	buf.L3HeaderList.PushBack(data)

	return nil
}

// L3HeaderBuffer に L3Header(FlowEthernetISO8023) を追加します。現在の時間を付加情報として埋め込みます
func (buf L3HeaderBuffer) AddL3Header(data *FlowEthernetISO8023) error {
	if ( data == nil ) {
		return fmt.Errorf("nil input")
	}
	l3HeaderWithTime := &L3HeaderWithTime{time.Now(), data}

	return buf.AddL3HeaderWithTime(l3HeaderWithTime)
}


// CountDataBuffer から指定された時間以降のデータを取り出します。
// データが無ければ長さ 0 のスライスを返します
func (buf L3HeaderBuffer) GetDataFromTime(beforeTime time.Time) ([]*FlowEthernetISO8023, error) {
	ite := buf.L3HeaderList.Front()
	for ; ite != nil; ite = ite.Next() {
		data := ite.Value.(*L3HeaderWithTime)
		if ( data.time.Sub(beforeTime) >= 0 ) {
			break
		}
	}
	retBuf := []*FlowEthernetISO8023{}
	for ; ite != nil; ite = ite.Next() {
		data := ite.Value.(*L3HeaderWithTime)
		retBuf = append(retBuf, data.l3Header)
	}
	return retBuf, nil
}

/*
// []*L4Header を受け取り、中にあるもののポート番号で sort uniq -c した値を返します
func CalcPortCount(l4HeaderArray []*L4Header) (*SrcDstCountData, error) {
	srcPortMap := map[uint16]int{}
	dstPortMap := map[uint16]int{}

	i := 0
	for i = 0; i < len(l4HeaderArray); i++ {
		l4Header := l4HeaderArray[i]
		if ( l4Header.Tcp != nil ) {
			srcPortMap[l4Header.Tcp.SourcePort]++
			dstPortMap[l4Header.Tcp.DestinationPort]++
		}else if ( l4Header.Udp != nil ) {
			srcPortMap[l4Header.Udp.SourcePort]++
			dstPortMap[l4Header.Udp.DestinationPort]++
		}
	}

	countData := &SrcDstCountData{}

	countData.SourceCount = map[string]int{}
	countData.DestinationCount = map[string]int{}

	for port, count := range srcPortMap {
		countData.SourceCount[strconv.Itoa(int(port))] += count
	}
	for port, count := range dstPortMap {
		countData.DestinationCount[strconv.Itoa(int(port))] += count
	}

	return countData, nil
}
*/


