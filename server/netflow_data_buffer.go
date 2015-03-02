package main

import (
	"time"
	"fmt"
	"container/list"
)

type NetFlowData struct {
	Time        jsonTime  `json:"time"`
	Flow        string    `json:"flow"`
}

type NetFlowBuffer struct {
	MaxLength    int
	NetFlowList *list.List
}

func MakeNetFlowBuffer(maxLength int) *NetFlowBuffer{
	return &NetFlowBuffer{ MaxLength: maxLength,
		NetFlowList: list.New() }
}

// NetFlowData を追加します
func (buf NetFlowBuffer) AddNetFlowData(data *NetFlowData) error {
	if ( data == nil ) {
		return fmt.Errorf("nil input")
	}

	if buf.NetFlowList.Len() >= buf.MaxLength {
		buf.NetFlowList.Remove(buf.NetFlowList.Front())
	}

	buf.NetFlowList.PushBack(data)

	return nil
}

// NetFlowBuffer に JSONだと思われる文字列 を追加します。現在の時間を付加情報として埋め込みます
func (buf NetFlowBuffer) AddFlowJsonString(data string) error {
	if ( data == "" ) {
		return fmt.Errorf("nil input")
	}
	flowData := &NetFlowData{jsonTime{time.Now()}, data}

	return buf.AddNetFlowData(flowData)
}

// NetFlowBuffer から指定された時間以降のデータを取り出します。
// データが無ければ長さ 0 のスライスを返します
func (buf NetFlowBuffer) GetDataFromTime(beforeTime time.Time) ([]*NetFlowData, error) {
	ite := buf.NetFlowList.Front()
	for ; ite != nil; ite = ite.Next() {
		data := ite.Value.(*NetFlowData)
		if ( data.Time.Sub(beforeTime) >= 0 ) {
			break
		}
	}
	retBuf := []*NetFlowData{}
	for ; ite != nil; ite = ite.Next() {
		data := ite.Value.(*NetFlowData)
		retBuf = append(retBuf, data)
	}
	return retBuf, nil
}



