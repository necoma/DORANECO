package main

import (
	"time"
	"fmt"
	"container/list"
)

// NetFlowAlertData は、NetFlowのAlertを保持します
type NetFlowAlertData struct {
	Port        int           `json:"port"`
	Time        jsonTime      `json:"time"`
	Duration    time.Duration `json:"duration"`
	Flow        string        `json:"flow"`
}

// NetFlowAlertBuffer は NetFlowAlertData を MaxLength の数まで保持するための保管庫です
type NetFlowAlertBuffer struct {
	MaxLength    int
	AlertList    *list.List
}

// MakeNetFlowAlertBuffer は最大長を指定して NetFlowAlertBuffer を作成します
func MakeNetFlowAlertBuffer(maxLength int) *NetFlowAlertBuffer{
	return &NetFlowAlertBuffer{ MaxLength: maxLength,
		AlertList: list.New() }
}

// AddAlertData は NetFlowAlertData を追加します
func (buf NetFlowAlertBuffer) AddAlertData(data *NetFlowAlertData) error {
	if ( data == nil ) {
		return fmt.Errorf("nil input")
	}

	if buf.AlertList.Len() >= buf.MaxLength {
		buf.AlertList.Remove(buf.AlertList.Front())
	}

	buf.AlertList.PushBack(data)

	return nil
}

// AddAlert は NetFlowAlertBuffer に Flow を追加します。現在の時間を付加情報として埋め込みます
func (buf NetFlowAlertBuffer) AddAlert(port int, duration time.Duration, flow string) error {
	if ( flow == "" ) {
		return fmt.Errorf("nil input")
	}
	alertData := &NetFlowAlertData{port, jsonTime{time.Now()}, duration, flow}

	return buf.AddAlertData(alertData)
}


// GetDataFromTime は
// AlertBuffer から指定された時間以降のデータを取り出します。
// データが無ければ長さ 0 のスライスを返します
func (buf NetFlowAlertBuffer) GetDataFromTime(beforeTime time.Time) ([]*NetFlowAlertData, error) {
	ite := buf.AlertList.Front()
	for ; ite != nil; ite = ite.Next() {
		data := ite.Value.(*NetFlowAlertData)
		if ( data.Time.Sub(beforeTime) >= 0 ) {
			break
		}
	}
	retBuf := []*NetFlowAlertData{}
	for ; ite != nil; ite = ite.Next() {
		data := ite.Value.(*NetFlowAlertData)
		retBuf = append(retBuf, data)
	}
	return retBuf, nil
}



