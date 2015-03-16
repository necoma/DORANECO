package main

import (
	"time"
	"fmt"
	//"strconv"
	"container/list"
)

// SrcDstCountData は、
// source と destination のポート番号でまとめたパケットの数を保持します
type SrcDstCountData struct {
	SourceCount      map[string]int  `json:"Source"`
	DestinationCount map[string]int  `json:"Destination"`
}

// L4CountData は、Layer4 のそれぞれのプロトコル毎に SrcDstCountData を保持します
type L4CountData struct {
	TCP   SrcDstCountData  `json:"tcp"`
	UDP   SrcDstCountData  `json:"udp"`
	ICMP  SrcDstCountData  `json:"icmp"`
}

// L3HeaderWithTime は Etherフレーム とそのフレームが観測された時間を保持します
type L3HeaderWithTime struct {
	time     time.Time
	l3Header *FlowEthernetISO8023
}

// CountDataBuffer は指定された最大長までの L4CountData を保持します
type CountDataBuffer struct {
	MaxLength     int
	CountDataList *list.List
}

// L3HeaderBuffer は指定された最大長までの L3HeaderWithTime を保持します
type L3HeaderBuffer struct {
	MaxLength    int
	L3HeaderList *list.List
}

// MakeL3HeaderBuffer は最大長の指定された L3HeaderBuffer を生成します
func MakeL3HeaderBuffer(maxLength int) *L3HeaderBuffer{
	return &L3HeaderBuffer{ MaxLength: maxLength,
		L3HeaderList: list.New() }
}

// AddL3HeaderWithTime は、L3HeaderBuffer に L3HeaderWithTime を追加します
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

// AddL3Header は L3HeaderBuffer に Etherフレーム を追加します
// L3HeaderBuffer に L3Header(FlowEthernetISO8023) を追加します。現在の時間を付加情報として埋め込みます
func (buf L3HeaderBuffer) AddL3Header(data *FlowEthernetISO8023) error {
	if ( data == nil ) {
		return fmt.Errorf("nil input")
	}
	l3HeaderWithTime := &L3HeaderWithTime{time.Now(), data}

	return buf.AddL3HeaderWithTime(l3HeaderWithTime)
}


// GetDataFromTime は
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

