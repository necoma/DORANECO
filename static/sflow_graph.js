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
		{"IPv4Header":
			{"Version":4,"Len":20,"TOS":0,"TotalLen":1500,"ID":39884,"Flags":2,"FragOff":0,"TTL":62,"Protocol":6,"Checksum":18520,"Src":"157.82.117.15","Dst":"192.231.127.174","Options":null},
		"IPv6Header":null,
		"L4Header":
			{"Tcp":{"SourcePort":56474,"DestinationPort":24308},"Udp":null,"Icmp":null}
		},
		null
	];

// flowDataList を src IP で sort します
function PacketSort_SrcIP(flowDataList){
	flowDataList.sort(function(a, b){
		if('IPv4Header' in a && a.IPv4Header != null){
			if('IPv4Header' in b && b.IPv4Header != null){
				var aIP = a.IPv4Header.Src;
				var bIP = b.IPv4Header.Src;
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
		if('IPv6Header' in a && a.IPv6Header != null){
			if('IPv6Header' in b && b.IPv6Header != null){
				var aIP = a.IPv6Header.Src;
				var bIP = b.IPv6Header.Src;
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
function PacketUniq_SrcIP(flowDataList){
	var uniqList = [];
	var currentIP = "";
	var currentList = [];
	for(var i = 0; i < flowDataList.length; i++){
		var packet = flowDataList[i];
		var targetIP = "";
		if('IPv4Header' in packet && packet.IPv4Header != null) {
			targetIP = packet.IPv4Header.Src;
		}else if('IPv6Header' in packet && packet.IPv6Header != null){
			targetIP = packet.IPv6Header.Src;
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
function PacketSort_UniqedData(uniqedFlowDataList){
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

// logの部分にタブを一個追加してタブのセレクタに使える ID を返します
function AddLogTab(tabTitle, flowDataList, targetTabSelector){
	var originalData = $(targetTabSelector).data();
	var tabList = originalData.tabList;
	if(tabList == null || !Array.isArray(tabList) ){
		tabList = [];
	}
	var tabName = "log1";
	if(tabList.length > 0) {
		var lastTabData = tabList[tabList.length - 1];
		if("tabName" in lastTabData) {
			var lastTabName = lastTabData.tabName;
			var n = lastTabName.replace("log", "");
			var lastTabNumber = parseInt(n);
			lastTabNumber++;
			tabName = "log" + lastTabNumber;
		}
	}
	var data = {tabName: tabName, data: flowDataList};
	tabList.push(data);

	originalData.tabList = tabList;	
	$(targetTabSelector).data(originalData);
	AddTab(targetTabSelector, tabName, tabTitle);
	return tabName;
}

// srcIP, dstIP, SourcePort, DestinationPort をプロパティに持つ配列を表にして表示します
function PopupFlowDataList(flowDataList, targetTabSelector){
	var tabTitle = "RAW:" + flowDataList.port + " - " + FormatDate(flowDataList.dateTime, "hh:mm:ss");
	var tabName = AddLogTab(tabTitle, flowDataList, targetTabSelector);
	var template = $.templates("#PopupAddressPortTableTemplate");
	var html = template.render(flowDataList);
	$("#" + tabName).html(html);
	//$("#" + tabName + " table").tableSort();
	//$("#" + tabName + " .table-sort").tablesort();
	//$("table.table-sort").tablesort();
	//$("#LI_" + tabName).show();
}

function PacketListSortUniq_SrcIP(originalFlowDataList){
	var flowDataList = $.extend(true, {}, originalFlowDataList);
	var sortedFlowDataList = PacketSort_SrcIP(flowDataList.flowList);
	var uniqedFlowDataList = PacketUniq_SrcIP(sortedFlowDataList);
	var uniqSortedFlowDataList = PacketSort_UniqedData(uniqedFlowDataList);
	flowDataList.flowList = uniqSortedFlowDataList;
	flowDataList.uniqTarget = "src IP";
	return flowDataList;
}

// sort | uniq -c された srcIP, dstIP, SourcePort, DestinationPort をプロパティに持つ配列を表にして表示します
function PopupSortUniqedFlowDataList(flowDataList, targetTabSelector){
	var tabTitle = "UNIQ:" + flowDataList.port + " - " + FormatDate(flowDataList.dateTime, "hh:mm:ss");
	var tabName = AddLogTab(tabTitle, flowDataList, targetTabSelector);
	var template = $.templates("#PopupAddressPortTableTemplate_Uniq");
	var html = template.render(flowDataList);
	$("#" + tabName).html(html);
	//$("#" + tabName + " table").tableSort();
	//$("#" + tabName + " .table-sort").tablesort();
	//$("table.table-sort").tablesort();
	//$("#LI_" + tabName).show();
	$("#" + tabName + " .close").click(function(){ DelTab(targetTabSelector, tabName); });
}

// サーバから送られてきたデータを、TCP/UDP/ICMP に分けます
function splitL4Proto(rawEtherDataList){
	var tcp = [];
	var udp = [];
	var icmp = [];
	for(var i = 0; i < rawEtherDataList.length; i++){
		var data = rawEtherDataList[i];
		var l3Addr = {};
		if("IPv4Header" in data && data.IPv4Header != null) {
			var header = data.IPv4Header;
			l3Addr.srcIP = header.Src;
			l3Addr.dstIP = header.Dst;
		}
		if("IPv6Header" in data && data.IPv6Header != null) {
			var header = data.IPv6Header;
			l3Addr.srcIP = header.Src;
			l3Addr.dstIP = header.Dst;
		}
		if(!("L4Header" in data) || data.L4Header == null) {
			continue;
		}
		var l4Header = data.L4Header;
		if("Tcp" in l4Header && l4Header.Tcp != null){
			//tcp.push($.extend(true, l3Addr, l4Header.Tcp));
			tcp.push(data);
		}
		if("Udp" in l4Header && l4Header.Udp != null){
			//udp.push($.extend(true, l3Addr, l4Header.Udp));
			udp.push(data);
		}
		if("Icmp" in l4Header && l4Header.Icmp != null){
			//icmp.push($.extend(true, l3Addr, l4Header.Icmp));
			icmp.push(data);
		}
	}
	return {"tcp": tcp, "udp": udp, "icmp": icmp};
}

// TCP か UDP 形式のデータを表示用データに変換します。
// L4Header の source port, dest port を用いて、port の count とその port に纏わるpacketの束に変換します
function ProcessPacketDataToCountData(packetDataList){
	var portDictionary = {};
	for(var i = 0; i < packetDataList.length; i++){
		var packet = packetDataList[i];
		if ( !("L4Header" in packet) ) {
			continue;
		}
		var l4Header = packet.L4Header;
		var SourcePort = 0;
		var DestinationPort = 0;
		if ("Tcp" in l4Header && l4Header.Tcp != null){
			if( "SourcePort" in l4Header.Tcp && l4Header.Tcp.SourcePort != null){
				SourcePort = l4Header.Tcp.SourcePort;
			}
			if( "DestinationPort" in l4Header.Tcp && l4Header.Tcp.DestinationPort != null){
				DestinationPort = l4Header.Tcp.DestinationPort;
			}
		}else if("Udp" in l4Header && l4Header.Udp != null){
			if( "SourcePort" in l4Header.Udp && l4Header.Udp.SourcePort != null){
				SourcePort = l4Header.Udp.SourcePort;
			}
			if( "DestinationPort" in l4Header.Udp && l4Header.Udp.DestinationPort != null){
				DestinationPort = l4Header.Udp.DestinationPort;
			}
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

// サーバから取ってきたデータは source と destination に分かれているので、一つにまとめます
function AppendSrcDstData(srcDstCountData) {
	var result = {};
	if ('Source' in srcDstCountData ){
		var src = srcDstCountData.Source
		for ( var key in src ) {
			if (isNaN(result[key])) {
				result[key] = 0;
			}
			result[key] += src[key];
		}
	}
	if ('Destination' in srcDstCountData ){
		var dst = srcDstCountData.Destination
		for ( var key in dst ) {
			if (isNaN(result[key])) {
				result[key] = 0;
			}
			result[key] += dst[key];
		}
	}
	return result;
}

// {"port": count} の形式の map を count で sort してリストにして [{"port": "port", "count": count},,.]の返します
function ConvertCountMapToSortedList(countMap){
	var tmpList = [];
	for( var key in countMap ){
		var data = countMap[key];
		tmpList.push({"port": key, "count": data.count, "flowList": data.flowList});
	}
	tmpList.sort(function(a, b){
		if(a.count > b.count){
			return -1;
		}
		if(a.count < b.count){
			return 1;
		}
		return 0;
	});
	return tmpList;
}

// 反転した色を作ります
// http://yasu0120.blog130.fc2.com/blog-entry-21.html
function CreateInverseColor(color){
	baseColor = color.replace('#', "");
	if(baseColor.length != 6){
		return '#000000';
	}
	newColor = '';
	for(var x = 0; x < 3; x++){
		colorWK = 255 - parseInt(baseColor.substr((x*2),2),16);
		if(colorWK < 0){
			colorWK = 0;
		}else{
			colorWK = colorWK.toString(16);
		}
		if(colorWK.length < 2){
			colorWK = '0' + colorWK;
		}
		newColor += colorWK;
	}
	return "#" + newColor;
}

// sort されたリストからグラフ用のデータリストに変換します
function ConvertSortedListToGraphDataList(sortedList){
	color = d3.scale.category20();
	for (var i = 0; i < sortedList.length; i++ ){
		var data = sortedList[i];
		sortedList[i].title = data.port + " (" + data.count + ")";
		sortedList[i].color = color(data.port % 20);
		//sortedList[i].textColor = color((data.port+2) % 20);
		sortedList[i].textColor = CreateInverseColor(sortedList[i].color);
		sortedList[i]["dateTime"] = new Date();

		//console.log("color: ", sortedList[i].port, " to ", color(sortedList[i].port % 20), " 22 -> ", color(22 % 20));
	}
	return sortedList;
}

// 与えられたデータリストの count を書き換えて返します
function ModifyData(dataList){
	// 新しく生成して返すことにします。
	var resultList = [];
	var del = Math.random() * 4;
	for(var i = 0; i < (dataList.length - del); i++){
		var data = dataList[i];
		var newData = Math.random() * 100;
		resultList[i] = {};
		resultList[i]["color"] = data["color"];
		resultList[i]["title"] = data["title"];
		resultList[i]["textColor"] = data["textColor"];
		resultList[i]["count"] = newData;
		resultList[i]["dateTime"] = new Date();
	}
	return resultList;
}

var dataList = originalDataList.concat();

// targetSelector に svg を追加してその中に width/height の大きさの円グラフを描くための箱を作って返します
function CreateSVGArcGraphCanbas(targetSelector, width, height, clickEventHandler)
{
	var graphObj = {};

	// 指定されたセレクタの中に <svg>...</svg> を作る
	var svg = d3.select(targetSelector)
		.append("svg")
	// svg に width, height の attribute をつける
	svg
		.attr("width", width)
		.attr("height", height)

	// d3.layout.pie() で、sort せずに、値の ["percentage"] を値として使う関数を作る
	var pie = d3.layout.pie().sort(null).value(function(d){return d["count"];});
	// d3.svg.arc() で innerRadius, outerRadius を指定した関数を作る
	var arc = d3.svg.arc().innerRadius(innerRadius).outerRadius(outerRadius);

	// 円弧になる path の style
	var pathStyle = {
		fill: function(d){return d.data["color"];}
		, d: arc
		, stroke: "white"
	}

	// 文字を表示するための text の style
	var textStyle = {
		"transform": function(d) {
			d.innerRadius = innerRadius;
			d.outerRadius = outerRadius;
			return "translate(" + arc.centroid(d) + ")";
		}
		, "fill": function(d){return d.data['textColor'] }
		, "text-anchor": "middle"
	};

	// 作った svg等 を graphObj に格納
	graphObj.svg = svg;
	//graphObj.svg_g = svg_g;
	graphObj.pie = pie;
	graphObj.arc = arc;
	graphObj.pathStyle = pathStyle;
	graphObj.textStyle = textStyle;
	graphObj.clickEventHandler = clickEventHandler;
	graphObj.width = width;
	graphObj.height = height;

	return graphObj;
}
// CreateSVGArcGraphCanbas で作った graphObj に、newDataList で指示される円グラフを描きます
// duration が 0 より大きければ、指定された時間(ミリ秒)で新しい値まで遷移するようにします
function DrawArchGraph(graphObj, newDataList, duration){
	//console.log("draw", newDataList);

	// svg.g.arc を作る準備をする
	var dataInjectedArcGroup = graphObj.svg
		.selectAll("g.arc")
		.data(graphObj.pie(newDataList));

	// 既存の円弧を書き換える
	dataInjectedArcGroup
		.select("path")
		.transition()
		.attr(graphObj.pathStyle)
		.duration(duration)
		;
	dataInjectedArcGroup
		.select("text")
		.transition()
		.attr(graphObj.textStyle)
		.duration(duration)
		.text(function(d){return d.data["title"];});

	// 円弧一つづつ用の g を追加する
	var appendedArcGroup = dataInjectedArcGroup
		.enter()
		.append("g")
		.attr("class", "arc")
		.attr("transform", "translate(" + graphObj.width / 2 + "," + graphObj.height / 2 + ")")
		.on("click", function(d){if(graphObj.clickEventHandler){graphObj.clickEventHandler(d.data);}})
		;

	// 円弧を表示するための svg.g.path を作る
	appendedArcGroup
		.append("svg:path")
		.transition()
		.attr(graphObj.pathStyle)
		.duration(duration)
		;

	// svg.g.text を作る
	appendedArcGroup
		.append("svg:text")
		.transition()
		.attr(graphObj.textStyle)
		.duration(duration)
		.text(function(d){return d.data["title"];});  

	// svg.g.* のいらないものを削除
	dataInjectedArcGroup
		.exit()
		.transition()
		.attr("transform", "scale(0,0)")
		.duration(duration)
		.remove()
}

// データを読み込んで反映させます
function LoadNewData(tcpObj, udpObj){
	var duration = 10;
	var argDuration = Math.floor(location.href.split("?")[1]);
	if (argDuration > 0) {
		duration = argDuration;
	}
	GetJSON("/current_data.json?duration=" + duration, {}, function (sflowDataList){
		var l4ProtoCount = splitL4Proto(sflowDataList);
		var tcpCountMap = ProcessPacketDataToCountData(l4ProtoCount.tcp);
		var udpCountMap = ProcessPacketDataToCountData(l4ProtoCount.udp);
		var tcpSortedList = ConvertCountMapToSortedList(tcpCountMap);
		var udpSortedList = ConvertCountMapToSortedList(udpCountMap);
			
		if (tcpSortedList.length > 20) {
			tcpSortedList = tcpSortedList.slice(0, 20);
		}
		if (udpSortedList.length > 20) {
			udpSortedList = udpSortedList.slice(0, 20);
		}
		var tcpNewGraphData = ConvertSortedListToGraphDataList(tcpSortedList);
		DrawArchGraph(tcpObj, tcpNewGraphData, 300);

		var udpNewGraphData = ConvertSortedListToGraphDataList(udpSortedList);
		DrawArchGraph(udpObj, udpNewGraphData, 300);
	});
}

function LoopLoadNewData(tcpObj, udpObj, interval){
	LoadNewData(tcpObj, udpObj);
	setTimeout(function(){LoopLoadNewData(tcpObj, udpObj, interval);}, interval);
}

var alertDataList = [
		{ port: 1, time: "time", packetList: [] }
	];

function InjectPacketListIDFromAlertDataList(alertDataList){
	var newAlertDataList = [];
	for(var i = 0; i < alertDataList.alertList.length; i++){
		var data = {};
		var originalData = alertDataList.alertList[i];
		data.port = originalData.port;
		data.time = originalData.time;
		data.dateTime = new Date();
		data.flowList = originalData.packetList;
		data.packetListID = "PACKET_LIST_" + i;
		newAlertDataList.push(data);
	}
	return {alertList: newAlertDataList};
}
// alert の結果を popup します。
function AlertPopup(alertDataList, LogTabSelector){
	PopupSortUniqedFlowDataList(PacketListSortUniq_SrcIP(alertDataList), LogTabSelector);
	PopupFlowDataList(alertDataList, LogTabSelector);
};
function DrawAlertTable(alertDivSelector, alertDataList, LogTabSelector){
	var alertDataListWithPacketListID = InjectPacketListIDFromAlertDataList(alertDataList);
	var template = $.templates("#AlertTableTemplate");
	var html = template.render(alertDataListWithPacketListID);
	$(alertDivSelector).html(html);

	for(var i = 0; i < alertDataListWithPacketListID.alertList.length; i++){
		var data = alertDataListWithPacketListID.alertList[i];
		var popupFuncBinded = AlertPopup.bind(null, data).bind(null, LogTabSelector);
		$(alertDivSelector + " #" + data.packetListID).click(popupFuncBinded);
	}
}
function LoadAlertLog(alertDivSelector, LogTabSelector) {
	var duration = 600;
	GetJSON("/alert_data.json?duration=" + duration, {}, function (alertDataList){
		DrawAlertTable(alertDivSelector, {alertList: alertDataList.reverse()}, LogTabSelector);
	});
}

function LoopLoadAlertLog(alertDivSelector, LogTabSelector, interval){
	LoadAlertLog(alertDivSelector, LogTabSelector);
	setTimeout(function(){LoopLoadAlertLog(alertDivSelector, LogTabSelector, interval);}, interval);
}

function InitSFlowAlertLog(alertDivSelector, interval, LogTabSelector){
	//LoadAlertLog(alertDivSelector, LogTabSelector);
	LoopLoadAlertLog(alertDivSelector, LogTabSelector, interval);
	/*
	var prevRequestTime = new Date();
	setInterval(function(){
		// 長いこと表示されてないと setInterval が山ほど呼ばれるらしいので無視するようなのを入れます
		var now = new Date();
		var msec = now - prevRequestTime;
		if(msec < interval){
			return;
		}
		prevRequestTime = now;
		LoadAlertLog(alertDivSelector, LogTabSelector);
	}, interval);
	*/
}

// sflow のグラフを描くコンポーネントを作ります。
// コンポーネントを書き込む div 要素へのセレクタと、コンポーネントが使う id の名前空間(?)を指定します
// 例: InitSFlowGraphComponent("#tab1", "SflowGraph");
function InitSFlowGraphComponent(targetSelector, IDName) {
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
		PopupSortUniqedFlowDataList(PacketListSortUniq_SrcIP(d), LogTabSelector);
		PopupFlowDataList(d, LogTabSelector);
	});
	var udpObj = CreateSVGArcGraphCanbas(ID + "_UdpGraph", width, height, function(d){
		PopupSortUniqedFlowDataList(PacketListSortUniq_SrcIP(d), LogTabSelector);
		PopupFlowDataList(d, LogTabSelector);
	});

	// 右側の上は alert log 用に使います
	var AlertLogTabSelector = ID + "_RIGHT_UP";
	var AlertLogDivName = IDName + "_RIGHT_UP_ALERT_LOG";
	$(AlertLogTabSelector).html('<h3>alert list</h3><div id="' + AlertLogDivName + '"></div>');
	InitSFlowAlertLog("#" + AlertLogDivName, interval, LogTabSelector);

	// 初期状態を load します
	//LoadNewData(tcpObj, udpObj);
	// 定期更新を仕掛けます
	LoopLoadNewData(tcpObj, udpObj, interval);
	/*
	var prevRequestTime = new Date();
	setInterval(function(){
		// 長いこと表示されてないと setInterval が山ほど呼ばれるらしいので無視するようなのを入れます
		var now = new Date();
		var msec = now - prevRequestTime;
		if(msec < interval){
			return;
		}
		prevRequestTime = now;
		LoadNewData(tcpObj, udpObj);
	}, interval);
	*/
}


