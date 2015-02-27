package main

import (
	"time"
	"fmt"
	"container/list"
)

type AlertData struct {
	Port        uint16                   `json:"port"`
	Time        jsonTime                 `json:"time"`
	PacketList  []*FlowEthernetISO8023   `json:"packetList"`
}

type AlertBuffer struct {
	MaxLength    int
	AlertList    *list.List
}

func MakeAlertBuffer(maxLength int) *AlertBuffer{
	return &AlertBuffer{ MaxLength: maxLength,
		AlertList: list.New() }
}

// AlertData を追加します
func (buf AlertBuffer) AddAlertData(data *AlertData) error {
	if ( data == nil ) {
		return fmt.Errorf("nil input")
	}

	if buf.AlertList.Len() >= buf.MaxLength {
		buf.AlertList.Remove(buf.AlertList.Front())
	}

	buf.AlertList.PushBack(data)

	return nil
}

// AlertBuffer に PacketList([]*FlowEthernetISO8023) を追加します。現在の時間を付加情報として埋め込みます
func (buf AlertBuffer) AddPacketList(port uint16, data []*FlowEthernetISO8023) error {
	if ( data == nil ) {
		return fmt.Errorf("nil input")
	}
	alertData := &AlertData{port, jsonTime{time.Now()}, data}

	return buf.AddAlertData(alertData)
}


// AlertBuffer から指定された時間以降のデータを取り出します。
// データが無ければ長さ 0 のスライスを返します
func (buf AlertBuffer) GetDataFromTime(beforeTime time.Time) ([]*AlertData, error) {
	ite := buf.AlertList.Front()
	for ; ite != nil; ite = ite.Next() {
		data := ite.Value.(*AlertData)
		if ( data.Time.Sub(beforeTime) >= 0 ) {
			break
		}
	}
	retBuf := []*AlertData{}
	for ; ite != nil; ite = ite.Next() {
		data := ite.Value.(*AlertData)
		retBuf = append(retBuf, data)
	}
	return retBuf, nil
}



