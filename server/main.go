package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"time"
	"strconv"
	"os"
	"bytes"
	"syscall"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
)

// host order (usually little endian) -> network order (big endian)
func htons(n int) int {
	return int(int16(byte(n))<<8 | int16(byte(n>>8)))
}

func getCurrentSFlowDataJSON(res http.ResponseWriter, req *http.Request, l4HeaderBuffer *L3HeaderBuffer) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	if l4HeaderBuffer == nil {
		fmt.Fprintf(res, "%s", `{"result": "error", "description": "count data is not anavailable."}`)
		return
	}

	req.ParseForm()
	durationSecString := req.Form.Get("duration")
	durationSec, err := strconv.Atoi(durationSecString)
	if err != nil {
		durationSec = 5
	}
	// 指定された時間(秒)までのデータを取り出します
	currentDataList, err := l4HeaderBuffer.GetDataFromTime(time.Now().Add(time.Duration(durationSec * -1) * time.Second))
	if err != nil {
		fmt.Fprintf(res, "%s", `{"result": "error", "description": "can not collect count data."}`)
		return
	}
	js, err := json.Marshal(currentDataList)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(res, "%s", "{\"result\": \"error\", \"description\": \"convert json error\"}")
	}
	res.Write(js)
}

func getAlertDataJSON(res http.ResponseWriter, req *http.Request, alertBuffer *AlertBuffer) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	if alertBuffer == nil {
		fmt.Fprintf(res, "%s", `{"result": "error", "description": "alert data is not anavailable."}`)
		return
	}
	req.ParseForm()
	durationSecString := req.Form.Get("duration")
	durationSec, err := strconv.Atoi(durationSecString)
	if err != nil {
		durationSec = 5
	}
	// 指定された時間(秒)までのデータを取り出します
	currentDataList, err := alertBuffer.GetDataFromTime(time.Now().Add(time.Duration(durationSec * -1) * time.Second))
	if err != nil {
		fmt.Fprintf(res, "%s", `{"result": "error", "description": "can not collect alert data."}`)
		return
	}
	js, err := json.Marshal(currentDataList)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(res, "%s", "{\"result\": \"error\", \"description\": \"convert json error\"}")
	}
	res.Write(js)
}

func getCurrentFlowDataJSON(res http.ResponseWriter, req *http.Request, netFlowBuffer *NetFlowBuffer) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	if netFlowBuffer == nil {
		fmt.Fprintf(res, "%s", `{"result": "error", "description": "count data is not anavailable."}`)
		return
	}

	req.ParseForm()
	durationSecString := req.Form.Get("duration")
	durationSec, err := strconv.Atoi(durationSecString)
	if err != nil {
		durationSec = 5
	}
	// 指定された時間(秒)までのデータを取り出します
	currentDataList, err := netFlowBuffer.GetDataFromTime(time.Now().Add(time.Duration(durationSec * -1) * time.Second))
	if err != nil {
		fmt.Fprintf(res, "%s", `{"result": "error", "description": "can not collect count data."}`)
		return
	}
	var strBuffer bytes.Buffer
	for i := range currentDataList {
		if strBuffer.Len() > 0 {
			strBuffer.WriteString(",")
		}
		strBuffer.WriteString(currentDataList[i].Flow)
	}
	fmt.Fprintf(res, "[%v]", strBuffer.String())
}

func getNetFlowAlertDataJSON(res http.ResponseWriter, req *http.Request, alertBuffer *NetFlowAlertBuffer) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	if alertBuffer == nil {
		fmt.Fprintf(res, "%s", `{"result": "error", "description": "alert data is not anavailable."}`)
		return
	}
	req.ParseForm()
	durationSecString := req.Form.Get("duration")
	durationSec, err := strconv.Atoi(durationSecString)
	if err != nil {
		durationSec = 300
	}
	// 指定された時間(秒)までのデータを取り出します
	currentDataList, err := alertBuffer.GetDataFromTime(time.Now().Add(time.Duration(durationSec * -1) * time.Second))
	if err != nil {
		fmt.Fprintf(res, "%s", `{"result": "error", "description": "can not collect alert data."}`)
		return
	}
	var strBuffer bytes.Buffer
	for i := range currentDataList {
		currentData := currentDataList[i]
		if strBuffer.Len() > 0 {
			strBuffer.WriteString(",")
		}
		strBuffer.WriteString(fmt.Sprintf("{\"port\":%v, \"duration\":%v, \"time\":\"%v\", \"flow\":%v}",
			currentData.Port, currentData.Duration.Seconds(), currentData.Time.format(), currentData.Flow))
	}
	fmt.Fprintf(res, "[%v]", strBuffer.String())
}

func martiniMain(listen string, l4HeaderBuffer *L3HeaderBuffer, alertBuffer *AlertBuffer, netFlowBuffer *NetFlowBuffer, netFlowAlertBuffer *NetFlowAlertBuffer, authUser string, authPassword string) {
	m := martini.Classic()
	m.Use(martini.Static("static"))
	m.Get("/current_data.json", func(res http.ResponseWriter, req *http.Request) {
		getCurrentSFlowDataJSON(res, req, l4HeaderBuffer)
	})
	m.Get("/alert_data.json", func(res http.ResponseWriter, req *http.Request) {
		getAlertDataJSON(res, req, alertBuffer)
	})
	m.Get("/netflow_current_data.json", func(res http.ResponseWriter, req *http.Request) {
		getCurrentFlowDataJSON(res, req, netFlowBuffer)
	})
	m.Get("/netflow_alert_data.json", func(res http.ResponseWriter, req *http.Request) {
		getNetFlowAlertDataJSON(res, req, netFlowAlertBuffer)
	})
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Use(auth.Basic(authUser, authPassword))
	m.RunOnAddr(listen)
}

func switchUser(uid int, gid int){
	if syscall.Getuid() == 0 {
		syscall.Setgid(gid)
		syscall.Setuid(uid)
	}
}

func main(){
	configFile := "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	watcherConfig, err := ReadConfig(configFile)
	if err != nil {
		fmt.Println("config file parse error: ", err)
		fmt.Println("usage: ", os.Args[0], " config.json")
		return
	}
	l3HeaderBuffer := MakeL3HeaderBuffer(watcherConfig.SFlowListener.MaxPacketCount)
	alertBuffer := MakeAlertBuffer(watcherConfig.SFlowListener.MaxAlertCount)
	netFlowBuffer := MakeNetFlowBuffer(watcherConfig.NetFlowListener.MaxPacketCount)
	netFlowAlertBuffer := MakeNetFlowAlertBuffer(watcherConfig.NetFlowListener.MaxAlertCount)

	go martiniMain(watcherConfig.HTTPServer.Listen,
		l3HeaderBuffer,
		alertBuffer,
		netFlowBuffer,
		netFlowAlertBuffer,
		watcherConfig.HTTPServer.BasicAuthUser,
		watcherConfig.HTTPServer.BasicAuthPassword)
	switchUser(watcherConfig.RunningUserID, watcherConfig.RunningGroupID)
	
	stopChannel := make(chan bool)
	go PortWatcher(stopChannel,
		watcherConfig.SFlowListener.AlertTickTimeSecond,
		watcherConfig.SFlowListener.AlertPortList,
		watcherConfig.SFlowListener.AlertPortRangeList,
		l3HeaderBuffer,
		alertBuffer,
		&watcherConfig.NECOMAtter)

	go startNetFlowCollector(watcherConfig.NetFlowListener.Listen,
		netFlowBuffer,
		netFlowAlertBuffer,
		&watcherConfig.NetFlowListener.AlertTargetMap)

	fmt.Println("start sflow listener")
	err = startSflowCollector(watcherConfig.SFlowListener.Listen, l3HeaderBuffer)
	if err != nil {
		fmt.Println("startSflowCollector error: ", err)
	}
	stopChannel <- true
}
