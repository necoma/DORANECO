package main

import (
	"time"
	"fmt"
	"container/list"
)

// NetFlowData は flow の情報を含んだ JSON 形式の文字列を保持します
type NetFlowData struct {
	Time        jsonTime  `json:"time"`
	Flow        string    `json:"flow"`
}

// NetFlowBuffer は NetFlowData を MaxLength まで保持します
type NetFlowBuffer struct {
	MaxLength    int
	NetFlowList *list.List
}

// MakeNetFlowBuffer は NetFlowBuffer を生成します
func MakeNetFlowBuffer(maxLength int) *NetFlowBuffer{
	return &NetFlowBuffer{ MaxLength: maxLength,
		NetFlowList: list.New() }
}

// AddNetFlowData は、NetFlowBuffer に NetFlowData を追加します
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

// AddFlowJSONString は
// NetFlowBuffer に JSONだと思われる文字列 を追加します。現在の時間を付加情報として埋め込みます
func (buf NetFlowBuffer) AddFlowJSONString(data string) error {
	if ( data == "" ) {
		return fmt.Errorf("nil input")
	}
	flowData := &NetFlowData{jsonTime{time.Now()}, data}

	return buf.AddNetFlowData(flowData)
}

// GetDataFromTime は
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



