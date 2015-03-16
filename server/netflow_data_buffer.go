package main

import (
	"time"
	"fmt"
	"container/list"
)

// NetFlowData $B$O(B flow $B$N>pJs$r4^$s$@(B JSON $B7A<0$NJ8;zNs$rJ];}$7$^$9(B
type NetFlowData struct {
	Time        jsonTime  `json:"time"`
	Flow        string    `json:"flow"`
}

// NetFlowBuffer $B$O(B NetFlowData $B$r(B MaxLength $B$^$GJ];}$7$^$9(B
type NetFlowBuffer struct {
	MaxLength    int
	NetFlowList *list.List
}

// MakeNetFlowBuffer $B$O(B NetFlowBuffer $B$r@8@.$7$^$9(B
func MakeNetFlowBuffer(maxLength int) *NetFlowBuffer{
	return &NetFlowBuffer{ MaxLength: maxLength,
		NetFlowList: list.New() }
}

// AddNetFlowData $B$O!"(BNetFlowBuffer $B$K(B NetFlowData $B$rDI2C$7$^$9(B
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

// AddFlowJSONString $B$O(B
// NetFlowBuffer $B$K(B JSON$B$@$H;W$o$l$kJ8;zNs(B $B$rDI2C$7$^$9!#8=:_$N;~4V$rIU2C>pJs$H$7$FKd$a9~$_$^$9(B
func (buf NetFlowBuffer) AddFlowJSONString(data string) error {
	if ( data == "" ) {
		return fmt.Errorf("nil input")
	}
	flowData := &NetFlowData{jsonTime{time.Now()}, data}

	return buf.AddNetFlowData(flowData)
}

// GetDataFromTime $B$O(B
// NetFlowBuffer $B$+$i;XDj$5$l$?;~4V0J9_$N%G!<%?$r<h$j=P$7$^$9!#(B
// $B%G!<%?$,L5$1$l$PD9$5(B 0 $B$N%9%i%$%9$rJV$7$^$9(B
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



