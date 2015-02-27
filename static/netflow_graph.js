var width = 400;
var height = 340;
var outerRadius = Math.min(width, height) / 2 - 5;
var innerRadius = outerRadius * 0.2;
var duration = 2000;

var originalDataList = [
	{ "color": "#ff0000", "title": "ff0000", "count": 100, "textColor": "#000000" }
	, { "color": "#00ff00", "title": "00ff00", "count": 200, "textColor": "#000000" }
	, { "color": "#00ffff", "title": "00ffff", "count": 300, "textColor": "#000000" }
	, { "color": "#ffff00", "title": "ffff00", "count": 400, "textColor": "#000000" }
	, { "color": "#0000ff", "title": "0000ff", "count": 500, "textColor": "#ffffff" }
];

var etherDataList =
[
  {
    "SRC_AS": 0,
    "DST_AS": 0,
    "DIRECTION": "Ingress",
    "TCP_FLAGS": "   A    ",
    "SRC_MASK": 15,
    "IPV4_DST_ADDR": "202.249.39.226",
    "IPV4_SRC_ADDR": "184.105.139.118",
    "OUTPUT_SNMP": 81,
    "INPUT_SNMP": 93,
    "IN_PKTS": 1,
    "IN_BYTES": 40,
    "FIRST_SWITCHED": 4049479.828,
    "LAST_SWITCHED": 4049479.828,
    "PROTOCOL": "0x11",
    "SRC_TOS": "0x00",
    "L4_SRC_PORT": 58995,
    "L4_DST_PORT": 123,
    "FLOW_SAMPLER_ID": 0,
    "VENDOR_PROPRIETARY_50": "0x00",
    "IPV4_NEXT_HOP": "0.0.0.0",
    "DST_MASK": 0
  },
  {
    "TCP_FLAGS": "        ",
    "IPV6_SRC_MASK": 0,
    "IPV6_DST_MASK": 0,
    "IPV6_NEXT_HOP": "::",
    "L4_DST_PORT": 0,
    "L4_SRC_PORT": 128,
    "SRC_TOS": "0x00",
    "PROTOCOL": "0x3a",
    "LAST_SWITCHED": 4049621.511,
    "FIRST_SWITCHED": 4049619.527,
    "IN_BYTES": 288,
    "IN_PKTS": 3,
    "INPUT_SNMP": 105,
    "OUTPUT_SNMP": 85,
    "IPV6_SRC_ADDR": "2001:200:0:6002::a10:1a2",
    "IPV6_DST_ADDR": "2001:558:6000:4::2"
  },
  null
];

// flowDataList を src IP で sort します
function NetFlowPacketSort_SrcIP(flowDataList){
	flowDataList.sort(function(a, b){
		if('IPV4_SRC_ADDR' in a && a.IPV4_SRC_ADDR != null){
			if('IPV4_SRC_ADDR' in b && b.IPV4_SRC_ADDR != null){
				var aIP = a.IPV4_SRC_ADDR;
				var bIP = b.IPV4_SRC_ADDR;
				if(aIP == bIP){
					return 0;
				}
				if(aIP > bIP){
					return -1;
				}
				return 1;
			}
			return -1;
		}
		if('IPV6_SRC_ADDR' in a && a.IPV6_SRC_ADDR != null){
			if('IPV6_SRC_ADDR' in b && b.IPV6_SRC_ADDR != null){
				var aIP = a.IPV6_SRC_ADDR;
				var bIP = b.IPV6_SRC_ADDR;
				if(aIP == bIP){
					return 0;
				}
				if(aIP > bIP){
					return -1;
				}
				return 1;
			}
			return 1;
		}
		return -1;
	});
	return flowDataList;
}

// sort された flowDataList を src IP で uniq します。
// 返されるデータの型は
// [ [flowData, flowData, ...], [flowData, flowData, ...], ...]
// と、元の flowData を src IP 毎の flowData のリストにまとめたリストになります。
function NetFlowPacketUniq_SrcIP(flowDataList){
	var uniqList = [];
	var currentIP = "";
	var currentList = [];
	for(var i = 0; i < flowDataList.length; i++){
		var packet = flowDataList[i];
		var targetIP = "";
		if('IPV4_SRC_ADDR' in packet && packet.IPV4_SRC_ADDR != null) {
			targetIP = packet.IPV4_SRC_ADDR;
		}else if('IPV6_SRC_ADDR' in packet && packet.IPV6_SRC_ADDR != null){
			targetIP = packet.IPV6_SRC_ADDR;
		}
		if(targetIP == ""){
			continue;
		}
		if(currentIP != targetIP){
			if(currentList.length > 0){
				uniqList.push(currentList);
			}
			currentList = [];
		}
		currentList.push(packet);
		currentIP = targetIP;
	}
	if(currentList.length > 0){
		uniqList.push(currentList);
	}
	return uniqList;
}

// uniq された flowDataList を、uniqされた数で sort (数の多い方が先頭に集まる) します
function NetFlowPacketSort_UniqedData(uniqedFlowDataList){
	uniqedFlowDataList.sort(function(a, b){
		if(a.length > b.length){
			return -1;
		}else if(a.length < b.length){
			return 1;
		}
		return 0;
	});
	return uniqedFlowDataList;
}

// srcIP, dstIP, SourcePort, DestinationPort をプロパティに持つ配列を表にして表示します
function NetFlowPopupFlowDataList(flowDataList, targetTabSelector){
	var tabTitle = "RAW:" + flowDataList.port + " - " + FormatDate(flowDataList.dateTime, "hh:mm:ss");
	var tabName = AddLogTab(tabTitle, flowDataList, targetTabSelector);
	var template = $.templates("#NetFlowPopupAddressPortTableTemplate");
	var html = template.render(flowDataList);
	$("#" + tabName).html(html);
}

function NetFlowPacketListSortUniq_SrcIP(originalFlowDataList){
	var flowDataList = $.extend(true, {}, originalFlowDataList);
	var sortedFlowDataList = NetFlowPacketSort_SrcIP(flowDataList.flowList);
	var uniqedFlowDataList = NetFlowPacketUniq_SrcIP(sortedFlowDataList);
	var uniqSortedFlowDataList = NetFlowPacketSort_UniqedData(uniqedFlowDataList);
	flowDataList.flowList = uniqSortedFlowDataList;
	flowDataList.uniqTarget = "src IP";
	return flowDataList;
}

// sort | uniq -c された srcIP, dstIP, SourcePort, DestinationPort をプロパティに持つ配列を表にして表示します
function NetFlowPopupSortUniqedFlowDataList(flowDataList, targetTabSelector){
	var tabTitle = "UNIQ:" + flowDataList.port + " - " + FormatDate(flowDataList.dateTime, "hh:mm:ss");
	var tabName = AddLogTab(tabTitle, flowDataList, targetTabSelector);
	var template = $.templates("#NetFlowPopupAddressPortTableTemplate_Uniq");
	var html = template.render(flowDataList);
	$("#" + tabName).html(html);
}

// サーバから送られてきたデータを、TCP/UDP/ICMP に分けます
function NetFlowSplitL4Proto(rawEtherDataList){
	var tcp = [];
	var udp = [];
	var icmp = [];
	for(var i = 0; i < rawEtherDataList.length; i++){
		var data = rawEtherDataList[i];
		var l3Addr = {};
		if("IPV4_SRC_ADDR" in data && data.IPV4_SRC_ADDR != null) {
			l3Addr.srcIP = data.IPV4_SRC_ADDR;
		}
		if("IPV4_DST_ADDR" in data && data.IPV4_DST_ADDR != null) {
			l3Addr.dstIP = data.IPV4_DST_ADDR;
		}
		if("IPV6_SRC_ADDR" in data && data.IPV6_SRC_ADDR != null) {
			l3Addr.srcIP = data.IPV6_SRC_ADDR;
		}
		if("IPV6_DST_ADDR" in data && data.IPV6_DST_ADDR != null) {
			l3Addr.dstIP = data.IPV6_DST_ADDR;
		}
		if("PROTOCOL" in data && data.PROTOCOL != null){
			if(data.PROTOCOL == "TCP"){
				tcp.push(data);
			}
			if(data.PROTOCOL == "UDP"){
				udp.push(data);
			}
			if(data.PROTOCOL == "ICMP"){
				icmp.push(data);
			}
		}
	}
	return {"tcp": tcp, "udp": udp, "icmp": icmp};
}

// TCP か UDP 形式のデータを表示用データに変換します。
// L4Header の source port, dest port を用いて、port の count とその port に纏わるpacketの束に変換します
function NetFlowProcessPacketDataToCountData(packetDataList){
	var portDictionary = {};
	for(var i = 0; i < packetDataList.length; i++){
		var packet = packetDataList[i];
		var SourcePort = 0;
		var DestinationPort = 0;
		if("L4_SRC_PORT" in packet && packet.L4_SRC_PORT != null){
			SourcePort = packet.L4_SRC_PORT;
		}
		if("L4_DST_PORT" in packet && packet.L4_DST_PORT != null){
			DestinationPort = packet.L4_DST_PORT;
		}
		if( !(SourcePort in portDictionary) ) {
			portDictionary[SourcePort] = { "count": 0, "flowList": [] };
		}
		portDictionary[SourcePort].count++;
		portDictionary[SourcePort].flowList.push(packet);
		if( !(DestinationPort in portDictionary) ) {
			portDictionary[DestinationPort] = { "count": 0, "flowList": [] };
		}
		portDictionary[DestinationPort].count++;
		portDictionary[DestinationPort].flowList.push(packet);
	}
	return portDictionary;
}

// データを読み込んで反映させます
function NetFlowLoadNewData(tcpObj, udpObj){
	var duration = 10;
	var argDuration = Math.floor(location.href.split("?")[1]);
	if (argDuration > 0) {
		duration = argDuration;
	}
	GetJSON("/netflow_current_data.json?duration=" + duration, {}, function (sflowDataList){
		var l4ProtoCount = NetFlowSplitL4Proto(sflowDataList);
		var tcpCountMap = NetFlowProcessPacketDataToCountData(l4ProtoCount.tcp);
		var udpCountMap = NetFlowProcessPacketDataToCountData(l4ProtoCount.udp);
		var tcpSortedList = ConvertCountMapToSortedList(tcpCountMap);
		var udpSortedList = ConvertCountMapToSortedList(udpCountMap);
			
		if (tcpSortedList.length > 20) {
			var dropedPackets = ConcatCountMapList(tcpSortedList.slice(20));
			tcpSortedList = tcpSortedList.slice(0, 20);
			tcpSortedList.push(dropedPackets);
		}
		if (udpSortedList.length > 20) {
			var dropedPackets = ConcatCountMapList(udpSortedList.slice(20));
			udpSortedList = udpSortedList.slice(0, 20);
			udpSortedList.push(dropedPackets);
		}
		var tcpNewGraphData = ConvertSortedListToGraphDataList(tcpSortedList);
		DrawArchGraph(tcpObj, tcpNewGraphData, 300);

		var udpNewGraphData = ConvertSortedListToGraphDataList(udpSortedList);
		DrawArchGraph(udpObj, udpNewGraphData, 300);
	});
}

function NetFlowLoopLoadNewData(tcpObj, udpObj, interval){
	NetFlowLoadNewData(tcpObj, udpObj);
	setTimeout(function(){NetFlowLoopLoadNewData(tcpObj, udpObj, interval);}, interval);
}

var alertDataList = [
		{ port: 1, time: "time", packetList: [] }
	];

function NetFlowInjectPacketListIDFromAlertDataList(alertDataList){
	var newAlertDataList = [];
	for(var i = 0; i < alertDataList.alertList.length; i++){
		var data = {};
		var originalData = alertDataList.alertList[i];
		data.port = originalData.port;
		data.time = originalData.time;
		data.duration = originalData.duration;
		data.dateTime = new Date();
		data.flow = originalData.flow;
		data.packetListID = "NETFLOW_PACKET_LIST_" + i;
		newAlertDataList.push(data);
	}
	return {alertList: newAlertDataList};
}
// alert の結果を popup します。
function AlertPopup(alertDataList, LogTabSelector){
	NetFlowPopupSortUniqedFlowDataList(NetFlowPacketListSortUniq_SrcIP(alertDataList), LogTabSelector);
	NetFlowPopupFlowDataList(alertDataList, LogTabSelector);
};
function NetFlowDrawAlertTable(alertDivSelector, alertDataList, LogTabSelector){
	var alertDataListWithPacketListID = NetFlowInjectPacketListIDFromAlertDataList(alertDataList);
	var template = $.templates("#NetFlowAlertTableTemplate");
	var html = template.render(alertDataListWithPacketListID);
	$(alertDivSelector).html(html);

	for(var i = 0; i < alertDataListWithPacketListID.alertList.length; i++){
		var data = alertDataListWithPacketListID.alertList[i];
		var popupFuncBinded = AlertPopup.bind(null, data).bind(null, LogTabSelector);
		$(alertDivSelector + " #" + data.packetListID).click(popupFuncBinded);
	}
}
function NetFlowLoadAlertLog(alertDivSelector, LogTabSelector) {
	var duration = 60000;
	GetJSON("/netflow_alert_data.json?duration=" + duration, {}, function (alertDataList){
		NetFlowDrawAlertTable(alertDivSelector, {alertList: alertDataList.reverse()}, LogTabSelector);
	});
}

function NetFlowLoopLoadAlertLog(alertDivSelector, LogTabSelector, interval){
	NetFlowLoadAlertLog(alertDivSelector, LogTabSelector);
	setTimeout(function(){NetFlowLoopLoadAlertLog(alertDivSelector, LogTabSelector, interval);}, interval);
}

function InitNetFlowAlertLog(alertDivSelector, interval, LogTabSelector){
	NetFlowLoopLoadAlertLog(alertDivSelector, LogTabSelector, interval);
}

// netflow のグラフを描くコンポーネントを作ります。
// コンポーネントを書き込む div 要素へのセレクタと、コンポーネントが使う id の名前空間(?)を指定します
// 例: InitNetFlowGraphComponent("#tab1", "SflowGraph");
function InitNetFlowGraphComponent(targetSelector, IDName) {
	var interval = 2000;

	var ID = "#" + IDName;
	// 横に分割
	SeparateLeftRight(targetSelector, IDName + "_LEFT", IDName + "_RIGHT");
	// 左側を縦に分割
	SeparateUpDown(ID + "_LEFT", IDName + "_LEFT_UP", IDName + "_LEFT_DOWN");
	// tcpGraph, udpGraph の領域を作成
	$(ID + "_LEFT_UP").html("<h2>TCP graph</h2><div id=\"" + IDName + "_TcpGraph\"></div>");
	$(ID + "_LEFT_DOWN").html("<h2>UDP graph</h2><div id=\"" + IDName + "_UdpGraph\"></div>");
	// 右側も縦に分割
	SeparateUpDown(ID + "_RIGHT", IDName + "_RIGHT_UP", IDName + "_RIGHT_DOWN");

	// 右側の下は log 用に使います
	var LogTabSelector = ID + "_RIGHT_DOWN";
	AddTab_Init(LogTabSelector);
	var tcpObj = CreateSVGArcGraphCanbas(ID + "_TcpGraph", width, height, function(d){
		NetFlowPopupSortUniqedFlowDataList(NetFlowPacketListSortUniq_SrcIP(d), LogTabSelector);
		NetFlowPopupFlowDataList(d, LogTabSelector);
	});
	var udpObj = CreateSVGArcGraphCanbas(ID + "_UdpGraph", width, height, function(d){
		NetFlowPopupSortUniqedFlowDataList(NetFlowPacketListSortUniq_SrcIP(d), LogTabSelector);
		NetFlowPopupFlowDataList(d, LogTabSelector);
	});

	// 右側の上は alert log 用に使います
	var AlertLogTabSelector = ID + "_RIGHT_UP";
	var AlertLogDivName = IDName + "_RIGHT_UP_ALERT_LOG";
	$(AlertLogTabSelector).html('<h3>alert list</h3><div id="' + AlertLogDivName + '"></div>');
	InitNetFlowAlertLog("#" + AlertLogDivName, interval, LogTabSelector);

	// 定期更新を仕掛けます
	NetFlowLoopLoadNewData(tcpObj, udpObj, interval);
}


