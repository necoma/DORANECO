// JSON をGETします。
function GetJSON(url, data, success_func, error_func){
        $.ajax({url: url
                , type: "GET"
                , data: data
                , dataType: 'json'
                , success: success_func
                , error: error_func
        });
}

