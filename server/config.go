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

// PortRange は、ポート番号の範囲を示します
type PortRange struct {
	Low      int     `json:"low"`
	High     int     `json:"high"`
}

// SFlowListenerConfig は SFlow のリスナに必要な設定項目を保持します
type SFlowListenerConfig struct {
	Listen               string      `json:"listen"`
	MaxPacketCount       int         `json:"maxPacketCount"`
	MaxAlertCount        int         `json:"maxAlertCount"`
	AlertPortList        []int       `json:"alertPortList"`
	AlertTickTimeSecond  int         `json:"alertTickTimeSecond"`
	AlertPortRangeList   []PortRange `json:"alertPortRangeList"`
}

// NetFlowAlertTargetMap は NetFlow で alert を上げる時に使われる値で、
// key である string には文字列でポート番号を、
// value は単位が秒の数値を入れます。
// key が文字列であるのは、JSON の object の key が文字列しか許されていなかったためです。
type NetFlowAlertTargetMap map[string]float64

// NetFlowListenerConfig は、NetFlow のリスナに必要な設定項目を保持します
type NetFlowListenerConfig struct {
	Listen               string                 `json:"listen"`
	MaxPacketCount       int                    `json:"maxPacketCount"`
	MaxAlertCount        int                    `json:"maxAlertCount"`
	AlertTargetMap       NetFlowAlertTargetMap  `json:"alertTargetMap"`
}

// HTTPServerConfig は、HTTPServer に必要な設定項目を保持します
type HTTPServerConfig struct {
	Listen               string   `json:"listen"`
	BasicAuthUser        string   `json"basicAuthUser"`
	BasicAuthPassword    string   `json"basicAuthPassword"`
}

// WatcherConfig は、全体の設定を保持します
type WatcherConfig struct {
	SFlowListener        SFlowListenerConfig    `json:"sFlowListener"`
	NetFlowListener      NetFlowListenerConfig  `json:"netFlowListener"`
	HTTPServer           HTTPServerConfig       `json:"httpServer"`	
	NECOMAtter           MewUserData            `json:"necomatter"`
	RunningUserID        int                    `json:"runningUserId"`
	RunningGroupID       int                    `json:"runningGroupId"`
}

// ReadConfig は WatcherConfig を filename から読み込みます
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

