// 与えられたセレクタの中に、縦2つに割った領域のタグを生成する。
// この操作は破壊的です。指定されたセレクタの中身が全部吹き飛びます
function SeparateUpDown(targetElementSelector, upID, downID){
	var template = $.templates("#ComponentUpDownTemplate");
	var data = {upID: upID, downID: downID};
	var html = template.render(data);
	var target = $(targetElementSelector);
	target.html(html);
	$(targetElementSelector + " div.split-pane").splitPane();
}

// 与えられたセレクタの中に、横2つに割った領域のタグを生成する。
// この操作は破壊的です。指定されたセレクタの中身が全部吹き飛びます
function SeparateLeftRight(targetElementSelector, leftID, rightID){
	var template = $.templates('#ComponentRightLeftTemplate');
	var data = {leftID: leftID, rightID: rightID};
	var html = template.render(data);
	var target = $(targetElementSelector);
	target.html(html);
	$(targetElementSelector + " div.split-pane").splitPane();
}

// 与えられたセレクタに、bootstrap のタブを追加します
// タブの構造は
/*
<ul class="nav nav-tabs">
	<li class="active"><a href="#TABID1" data-toggle="tab">TABTITLE1</a></li>
	<li><a href="#TABID2" data-toggole="tab">TABTITLE2</a></li>
</ul>
<div class="tab-content">
	<div class="tab-pane fade in active" id="TABID1"></div>
	<div class="tab-pane fade" id="TABID2"></div>
</div>
*/
// と複雑なので、ちょっと面倒くさい指定になります。

// ひな形となるタグ <ul class="nav nav-tabs"></ul><div class="tab-contnnt"></div> を書き込みます。
// この操作は破壊的です。指定されたセレクタの中身が全部吹き飛びます
function AddTab_Init(targetElementSelector){
	$(targetElementSelector).html('<ul class="nav nav-tabs"></ul><div class="tab-content"></div>');
	$(targetElementSelector + " .tab-content").height("100%");
}

// AddTab_Init() で書き込まれたタブの中に、タブを一つ追加します。
function AddTab(targetElementSelector, tabID, tabTitle){
	var ulElement = $(targetElementSelector + " ul.nav");
	var divElement = $(targetElementSelector + " div.tab-content");

	var closeButtonHTML = '&nbsp;<button class="close">&times;</button>';

	ulElement.append('<li id="LI_' + tabID + '"><a href="#' + tabID + '" data-toggle="tab">' + tabTitle + closeButtonHTML + '</a></li>');
	divElement.append('<div class="tab-pane fade" id="' + tabID + '"></div>');

	// 何もない初期状態でタグだけ追加しても show ってやってやらないと何も表示されないので
	if($(targetElementSelector + " ul.nav li").size() <= 1){
		$(targetElementSelector + ' a:first').tab('show')
	}
	$('#' + tabID).height('100%');
	$("#LI_" + tabID + " .close").click(function(){ DelTab(targetElementSelector, tabID); });
}

// AddTab() で追加したタブを消します。
function DelTab(targetElementSelector, tabID){
	var liElement = $("#LI_" + tabID);
	var divElement = $("#" + tabID);
	var isActive = liElement.hasClass("active");
	liElement.remove();
	divElement.remove();

	// 表示されているタブを消したので、別のタブを表示させます
	if(isActive){
		if($(targetElementSelector + " ul.nav li").size() >= 1){
			$(targetElementSelector + ' a:first').tab('show');
		}
	}
}

