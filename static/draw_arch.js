
// targetSelector に svg を追加してその中に width/height の大きさの円グラフを描くための箱を作って返します
function CreateSVGArcGraphCanbas(targetSelector, width, height, clickEventHandler)
{
	var outerRadius = Math.min(width, height) / 2 - 5;
	var innerRadius = outerRadius * 0.2;
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
	graphObj.scale = 1.0;
	graphObj.targetSelector = targetSelector;

	return graphObj;
}
// CreateSVGArcGraphCanbas で作った graphObj に、newDataList で指示される円グラフを描きます
// duration が 0 より大きければ、指定された時間(ミリ秒)で新しい値まで遷移するようにします
function DrawArchGraph(graphObj, newDataList, duration, scale){
	//console.log("draw", newDataList);
	if ( scale < 0.5 ){
		scale = 1.0;
	}
	// scale が前回と違っていたら全部作り直します。
	if (graphObj.scale != scale){
		graphObj.svg.selectAll().remove();
		//graphObj = CreateSVGArcGraphCanbas(graphObj.targetSelector
		//	, graphObj.width, graphObj.height, graphObj.clickEventHandler);
	}

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
		//.attr("transform", "translate(" + graphObj.width / 2 + "," + graphObj.height / 2 + ")")
		.attr("transform", "translate(" + graphObj.width / 2 + "," + graphObj.height / 2 + ")scale(" + scale + ")")
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


function Count2Scale(count, min, max){
	count -= min;
	scale = (count / (max - min)) * 1.5;
	if(scale < 0.5){
		return 0.5;
	}
	if(scale > 1.3){
		return 1.3;
	}
	return scale;
}

