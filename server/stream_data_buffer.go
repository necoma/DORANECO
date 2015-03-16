package main

import (
	"time"
	"fmt"
	//"strconv"
	"container/list"
)

// SrcDstCountData $B$O!"(B
// source $B$H(B destination $B$N%]!<%HHV9f$G$^$H$a$?%Q%1%C%H$N?t$rJ];}$7$^$9(B
type SrcDstCountData struct {
	SourceCount      map[string]int  `json:"Source"`
	DestinationCount map[string]int  `json:"Destination"`
}

// L4CountData $B$O!"(BLayer4 $B$N$=$l$>$l$N%W%m%H%3%kKh$K(B SrcDstCountData $B$rJ];}$7$^$9(B
type L4CountData struct {
	TCP   SrcDstCountData  `json:"tcp"`
	UDP   SrcDstCountData  `json:"udp"`
	ICMP  SrcDstCountData  `json:"icmp"`
}

// L3HeaderWithTime $B$O(B Ether$B%U%l!<%`(B $B$H$=$N%U%l!<%`$,4QB,$5$l$?;~4V$rJ];}$7$^$9(B
type L3HeaderWithTime struct {
	time     time.Time
	l3Header *FlowEthernetISO8023
}

// CountDataBuffer $B$O;XDj$5$l$?:GBgD9$^$G$N(B L4CountData $B$rJ];}$7$^$9(B
type CountDataBuffer struct {
	MaxLength     int
	CountDataList *list.List
}

// L3HeaderBuffer $B$O;XDj$5$l$?:GBgD9$^$G$N(B L3HeaderWithTime $B$rJ];}$7$^$9(B
type L3HeaderBuffer struct {
	MaxLength    int
	L3HeaderList *list.List
}

// MakeL3HeaderBuffer $B$O:GBgD9$N;XDj$5$l$?(B L3HeaderBuffer $B$r@8@.$7$^$9(B
func MakeL3HeaderBuffer(maxLength int) *L3HeaderBuffer{
	return &L3HeaderBuffer{ MaxLength: maxLength,
		L3HeaderList: list.New() }
}

// AddL3HeaderWithTime $B$O!"(BL3HeaderBuffer $B$K(B L3HeaderWithTime $B$rDI2C$7$^$9(B
// L3HeaderWithTime $B$rDI2C$7$^$9(B
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

// AddL3Header $B$O(B L3HeaderBuffer $B$K(B Ether$B%U%l!<%`(B $B$rDI2C$7$^$9(B
// L3HeaderBuffer $B$K(B L3Header(FlowEthernetISO8023) $B$rDI2C$7$^$9!#8=:_$N;~4V$rIU2C>pJs$H$7$FKd$a9~$_$^$9(B
func (buf L3HeaderBuffer) AddL3Header(data *FlowEthernetISO8023) error {
	if ( data == nil ) {
		return fmt.Errorf("nil input")
	}
	l3HeaderWithTime := &L3HeaderWithTime{time.Now(), data}

	return buf.AddL3HeaderWithTime(l3HeaderWithTime)
}


// GetDataFromTime $B$O(B
// CountDataBuffer $B$+$i;XDj$5$l$?;~4V0J9_$N%G!<%?$r<h$j=P$7$^$9!#(B
// $B%G!<%?$,L5$1$l$PD9$5(B 0 $B$N%9%i%$%9$rJV$7$^$9(B
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

