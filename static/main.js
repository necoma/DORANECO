
// tab を作る
AddTab_Init("#GraphBody");
AddTab("#GraphBody", "SFlowWatcher", "sflow watcher");
AddTab("#GraphBody", "NetFlowWatcher", "netflow watcher");
AddTab("#GraphBody", "tab3", "TAB3");
$('#tab3').html("新規コンテンツ募集中");

InitSFlowGraphComponent("#SFlowWatcher", "SFlowWatcherComponent");
InitNetFlowGraphComponent("#NetFlowWatcher", "NetFlowWatcherComponent");

$("#SFlowWatcherComponent").css("background", "#ffffdd");
$("#NetFlowWatcherComponent").css("background", "#ffddff");


