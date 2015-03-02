
// tab を作る
AddTab_Init("#GraphBody");
AddTab("#GraphBody", "DORANECO", "DORANECO: Dashboard for Observing Realtime Attacks by NECOMA");
//AddTab("#GraphBody", "SFlowWatcher", "sflow watcher");
//AddTab("#GraphBody", "NetFlowWatcher", "netflow watcher");
AddTab("#GraphBody", "tab2", "TAB2");
$('#tab2').html("新規コンテンツ募集中");

SeparateUpDown("#DORANECO", "SFlowWatcher", "NetFlowWatcher");

InitSFlowGraphComponent("#SFlowWatcher", "SFlowWatcherComponent", 200, 170);
InitNetFlowGraphComponent("#NetFlowWatcher", "NetFlowWatcherComponent", 200, 170);

$("#SFlowWatcherComponent").css("background", "#ffffdd");
$("#NetFlowWatcherComponent").css("background", "#ffddff");


