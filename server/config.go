package main

import (
	"encoding/json"
	"io/ioutil"
)

/* example
{
	"runningUserId": 1001,
	"runningGroupId": 1001,
	"sFlowListener":{
		"listen": "0.0.0.0:6343",
		"maxPacketCount": 100000,
		"maxAlertCount": 1024,
		"alertPortList": [
			123, 53, 161, 1900
		],
		"alertPortRangeList": [
			{"low": 1024, "high": 32768}
		],
		"alertTickTimeSecond": 60
	},
	"netFlowListener":{
		"listen": "0.0.0.0:2055",
		"maxPacketCount": 100000,
		"maxAlertCount": 1024,
		"alertTargetMap": { "443": 300 }
	},
	"httpServer":{
		"listen": "0.0.0.0:30000",
		"basicAuthUser": "hello",
		"basicAuthPassword": "world"
	},
	"necomatter":{
		"userName": "iimura",
		"apiKey": ""
	}
}
*/

// PortRange $B$O!"%]!<%HHV9f$NHO0O$r<($7$^$9(B
type PortRange struct {
	Low      int     `json:"low"`
	High     int     `json:"high"`
}

// SFlowListenerConfig $B$O(B SFlow $B$N%j%9%J$KI,MW$J@_Dj9`L\$rJ];}$7$^$9(B
type SFlowListenerConfig struct {
	Listen               string      `json:"listen"`
	MaxPacketCount       int         `json:"maxPacketCount"`
	MaxAlertCount        int         `json:"maxAlertCount"`
	AlertPortList        []int       `json:"alertPortList"`
	AlertTickTimeSecond  int         `json:"alertTickTimeSecond"`
	AlertPortRangeList   []PortRange `json:"alertPortRangeList"`
}

// NetFlowAlertTargetMap $B$O(B NetFlow $B$G(B alert $B$r>e$2$k;~$K;H$o$l$kCM$G!"(B
// key $B$G$"$k(B string $B$K$OJ8;zNs$G%]!<%HHV9f$r!"(B
// value $B$OC10L$,IC$N?tCM$rF~$l$^$9!#(B
// key $B$,J8;zNs$G$"$k$N$O!"(BJSON $B$N(B object $B$N(B key $B$,J8;zNs$7$+5v$5$l$F$$$J$+$C$?$?$a$G$9!#(B
type NetFlowAlertTargetMap map[string]float64

// NetFlowListenerConfig $B$O!"(BNetFlow $B$N%j%9%J$KI,MW$J@_Dj9`L\$rJ];}$7$^$9(B
type NetFlowListenerConfig struct {
	Listen               string                 `json:"listen"`
	MaxPacketCount       int                    `json:"maxPacketCount"`
	MaxAlertCount        int                    `json:"maxAlertCount"`
	AlertTargetMap       NetFlowAlertTargetMap  `json:"alertTargetMap"`
}

// HTTPServerConfig $B$O!"(BHTTPServer $B$KI,MW$J@_Dj9`L\$rJ];}$7$^$9(B
type HTTPServerConfig struct {
	Listen               string   `json:"listen"`
	BasicAuthUser        string   `json"basicAuthUser"`
	BasicAuthPassword    string   `json"basicAuthPassword"`
}

// WatcherConfig $B$O!"A4BN$N@_Dj$rJ];}$7$^$9(B
type WatcherConfig struct {
	SFlowListener        SFlowListenerConfig    `json:"sFlowListener"`
	NetFlowListener      NetFlowListenerConfig  `json:"netFlowListener"`
	HTTPServer           HTTPServerConfig       `json:"httpServer"`	
	NECOMAtter           MewUserData            `json:"necomatter"`
	RunningUserID        int                    `json:"runningUserId"`
	RunningGroupID       int                    `json:"runningGroupId"`
}

// ReadConfig $B$O(B WatcherConfig $B$r(B filename $B$+$iFI$_9~$_$^$9(B
func ReadConfig(filename string) (*WatcherConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config WatcherConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.HTTPServer.Listen == "" {
		config.HTTPServer.Listen = ":3000"
	}

	if config.SFlowListener.Listen == "" {
		config.SFlowListener.Listen = ":6343"
	}
	if config.SFlowListener.MaxPacketCount <= 0 {
		config.SFlowListener.MaxPacketCount = 102400
	}
	if config.SFlowListener.MaxAlertCount <= 0 {
		config.SFlowListener.MaxAlertCount = 1024
	}
	if config.SFlowListener.AlertTickTimeSecond <= 0 {
		config.SFlowListener.AlertTickTimeSecond = 60
	}

	if config.NetFlowListener.Listen == "" {
		config.NetFlowListener.Listen = ":2055"
	}
	if config.NetFlowListener.MaxPacketCount <= 0 {
		config.NetFlowListener.MaxPacketCount = 102400
	}
	if config.NetFlowListener.MaxAlertCount <= 0 {
		config.NetFlowListener.MaxAlertCount = 1024
	}

	return &config, nil
}

