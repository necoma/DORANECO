package main

import (
	"encoding/json"
	"io/ioutil"
)

/*
{
	"sFlowListener":{
		"listen": "0.0.0.0:6343",
		"maxPacketCount": 100000,
		"maxAlertCount": 1024,
		"alertPortList": [
			123, 53, 161, 1900
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
	}
	"necomatter":{
		"userName": "iimura",
		"apiKey": ""
	}
}
*/


type SFlowListenerConfig struct {
	Listen               string   `json:"listen"`
	MaxPacketCount       int      `json:"maxPacketCount"`
	MaxAlertCount        int      `json:"maxAlertCount"`
	AlertPortList        []int    `json:"alertPortList"`
	AlertTickTimeSecond  int      `json:"alertTickTimeSecond"`
}

type NetFlowAlertTargetMap map[string]float64

type NetFlowListenerConfig struct {
	Listen               string                 `json:"listen"`
	MaxPacketCount       int                    `json:"maxPacketCount"`
	MaxAlertCount        int                    `json:"maxAlertCount"`
	AlertTargetMap       NetFlowAlertTargetMap  `json:"alertTargetMap"`
}

type HttpServerConfig struct {
	Listen               string   `json:"listen"`
	BasicAuthUser        string   `json"basicAuthUser"`
	BasicAuthPassword    string   `json"basicAuthPassword"`
}

type WatcherConfig struct {
	SFlowListener        SFlowListenerConfig    `json:"sFlowListener"`
	NetFlowListener      NetFlowListenerConfig  `json:"netFlowListener"`
	HttpServer           HttpServerConfig       `json:"httpServer"`	
	NECOMAtter           MewUserData            `json:"necomatter"`
}

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

	if config.HttpServer.Listen == "" {
		config.HttpServer.Listen = ":3000"
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

