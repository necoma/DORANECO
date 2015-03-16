package main

import (
	"time"
	"fmt"
	"container/list"
)

// NetFlowAlertData $B$O!"(BNetFlow$B$N(BAlert$B$rJ];}$7$^$9(B
type NetFlowAlertData struct {
	Port        int           `json:"port"`
	Time        jsonTime      `json:"time"`
	Duration    time.Duration `json:"duration"`
	Flow        string        `json:"flow"`
}

// NetFlowAlertBuffer $B$O(B NetFlowAlertData $B$r(B MaxLength $B$N?t$^$GJ];}$9$k$?$a$NJ]4I8K$G$9(B
type NetFlowAlertBuffer struct {
	MaxLength    int
	AlertList    *list.List
}

// MakeNetFlowAlertBuffer $B$O:GBgD9$r;XDj$7$F(B NetFlowAlertBuffer $B$r:n@.$7$^$9(B
func MakeNetFlowAlertBuffer(maxLength int) *NetFlowAlertBuffer{
	return &NetFlowAlertBuffer{ MaxLength: maxLength,
		AlertList: list.New() }
}

// AddAlertData $B$O(B NetFlowAlertData $B$rDI2C$7$^$9(B
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

// AddAlert $B$O(B NetFlowAlertBuffer $B$K(B Flow $B$rDI2C$7$^$9!#8=:_$N;~4V$rIU2C>pJs$H$7$FKd$a9~$_$^$9(B
func (buf NetFlowAlertBuffer) AddAlert(port int, duration time.Duration, flow string) error {
	if ( flow == "" ) {
		return fmt.Errorf("nil input")
	}
	alertData := &NetFlowAlertData{port, jsonTime{time.Now()}, duration, flow}

	return buf.AddAlertData(alertData)
}


// GetDataFromTime $B$O(B
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



