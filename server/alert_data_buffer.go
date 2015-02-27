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

// AlertData $B$rDI2C$7$^$9(B
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

// AlertBuffer $B$K(B PacketList([]*FlowEthernetISO8023) $B$rDI2C$7$^$9!#8=:_$N;~4V$rIU2C>pJs$H$7$FKd$a9~$_$^$9(B
func (buf AlertBuffer) AddPacketList(port uint16, data []*FlowEthernetISO8023) error {
	if ( data == nil ) {
		return fmt.Errorf("nil input")
	}
	alertData := &AlertData{port, jsonTime{time.Now()}, data}

	return buf.AddAlertData(alertData)
}


// AlertBuffer $B$+$i;XDj$5$l$?;~4V0J9_$N%G!<%?$r<h$j=P$7$^$9!#(B
// $B%G!<%?$,L5$1$l$PD9$5(B 0 $B$N%9%i%$%9$rJV$7$^$9(B
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



