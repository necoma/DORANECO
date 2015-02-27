package main

import (
	"time"
	"fmt"
	"container/list"
)

type NetFlowAlertData struct {
	Port        int           `json:"port"`
	Time        jsonTime     `json:"time"`
	Duration    time.Duration `json:"duration"`
	Flow        string        `json:"flow"`
}

type NetFlowAlertBuffer struct {
	MaxLength    int
	AlertList    *list.List
}

func MakeNetFlowAlertBuffer(maxLength int) *NetFlowAlertBuffer{
	return &NetFlowAlertBuffer{ MaxLength: maxLength,
		AlertList: list.New() }
}

// NetFlowAlertData $B$rDI2C$7$^$9(B
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

// NetFlowAlertBuffer $B$K(B Flow $B$rDI2C$7$^$9!#8=:_$N;~4V$rIU2C>pJs$H$7$FKd$a9~$_$^$9(B
func (buf NetFlowAlertBuffer) AddAlert(port int, duration time.Duration, flow string) error {
	if ( flow == "" ) {
		return fmt.Errorf("nil input")
	}
	alertData := &NetFlowAlertData{port, jsonTime{time.Now()}, duration, flow}

	return buf.AddAlertData(alertData)
}


// AlertBuffer $B$+$i;XDj$5$l$?;~4V0J9_$N%G!<%?$r<h$j=P$7$^$9!#(B
// $B%G!<%?$,L5$1$l$PD9$5(B 0 $B$N%9%i%$%9$rJV$7$^$9(B
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



