/*
*
*  {
      "mac": 2,
      "funcode": 3,
      "Ip": "61.140.21.146:13397",
      "timestamp": 1543489996399,
      "templ": 12,
      "data_fields_len": 10,
      "data_fieldsmap": {
        "0": 12,             ---//温度
        "1": 20,              ---//水位
        "2": 0,                  ---//湿度
        "3": 21,                    ---//电流
        "4": 30,                   ---//亮度
        "5": 31,                 ---//压强
        "6": 17,
        "7": 0,
        "8": 8,
        "9": 24
      }
    }
*
*
* */


var targetItems = []

function pushItem(chart, opt) {
    var item = {}
    item.chart = chart
    item.opt = opt
    targetItems.push(item)
}

function initTargets() {
    pushItem(myChart, opt1)
    pushItem(myChart2, opt2)
    pushItem(myChart3, opt3)
    pushItem(electC, electOpt)
    pushItem(lightC, lightOpt)
    pushItem(pressC, pressOpt)
}
initTargets()




$(document).ready(function () {
    var baseUrl = $("#data-store-api").attr("targetUrl")

    function getQueryString(name) {
        var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
        var r = window.location.search.substr(1).match(reg);
        if (r != null) return unescape(r[2]);
        return null;
    }


    var macValue =  getQueryString("mac")

    //config
    var url =  baseUrl+"/dev/log/"+macValue+"/{lastDataTime}"
    var syntimeurl = baseUrl+"/nowtime"


    var existUrl = baseUrl+"/dev/log/stat/exist/"+macValue

    function checkMacInput() {
        if  (macValue==null){
            $('#mac').text("无")
            $('#statetip').text("错误提醒，url连接地址需要用'问号'(?)携带mac地址参数！格式：http://ip:port/rtdata.html?mac=2")
            $('#statetip').addClass("state_err")
            return false
        }
        return true
    }


    function isExistMac(successFun) {

        $.ajax({
            url: existUrl,
            type: "GET",
            dataType: "jsonp",  //指定服务器返回的数据类型
            jsonpCallback: "success",
            timeout : 3000,
            success: function (data) {
                if (data != null) {
                    if (data.res == 1) { //存在改mac终端的监控记录
                        successFun()
                        $('#statetip').text("该设备mac有效")
                        $('#statetip').addClass("state_ok")
                    }else{
                        $('#statetip').text("注意，设备mac没有历史采集记录")
                        $('#statetip').addClass("state_tip")
                    }
                }else{
                    $('#statetip').text("错误：server返回数据有不对")
                    $('#statetip').addClass("state_tip")
                }
            },
            complete : function(XMLHttpRequest,status){ //请求完成后最终执行参数
                if(status=='timeout'){//超时,status还有success,error等值的情况
                   // ajaxTimeoutTest.abort();
                   // alert("超时");
                    $('#statetip').text("数据仓库服务器不在工作状态")
                    $('#statetip').addClass("state_err")
                }
            }

        })
    }


    function refreshPanData(panItem, value) {
        panItem.opt.series[0].data[0].value = value
        panItem.chart.setOption(panItem.opt, true);
    }

    var currItem
    var tempData

    function updatePan(len, datamap) {
        currItem = null
        tempData = null

        if (len > 0 && datamap != null) {
            for (var i = 0; i < len; i++) {
                currItem = targetItems[i]
                tempData = datamap[i]
                if (currItem != null && !isNaN(Number(tempData))) {
                    refreshPanData(currItem, Number(tempData))
                }
            }

        }

    }
   /* function getLocalTime(nS) {
        return new Date(parseInt(nS)).toLocaleString().replace(/:\d{1,2}$/,' ');
    }*/

    function getLocalTime(nS) {
        return new Date(parseInt(nS)).toLocaleString().replace(/年|月/g, "-").replace(/日/g, " ");
    }
    // ------ui-------

    function renderData(data) {
        if (data != null) {
            $('#mac').text(data.mac | "null")
            $('#ip').text(data.Ip)
            $('#time').text(getLocalTime(data.timestamp))

            //更新仪表盘
            updatePan(data.data_fields_len, data.data_fieldsmap)
        }


    }


    var svrTimeStmap = 0
    var lastDataTime = 0
    var dynamicData
    var pollDataInterval = 2000

    function insertLog(datalist) {
        var len = datalist.length

        lastDataTime = datalist[len - 1].timestamp

        var data = datalist.pop()  //获得最新一条数据

        renderData(data)

    }


    function tickRequst() {
        var reqUrl = url.replace("{lastDataTime}", lastDataTime)

        $.ajax({
            url: reqUrl,
            type: "GET",
            dataType: "jsonp",  //指定服务器返回的数据类型
            jsonpCallback: "success",
            success: function (data) {
                if (data != null) {
                    dynamicData = data
                    if (data.svrtime) {
                        svrTimeStmap = data.svrtime
                    }

                    if (data.list != null && data.list.constructor == Array && data.list.length > 0) {
                        insertLog(data.list)
                    }

                }


            }
        });
    }

    //异步获取远程数据
    function pollDataFromRoute() {
        setInterval(function () {
            tickRequst()
        }, pollDataInterval)
    }


    function synSvrtime() {
        $.ajax({
            url: syntimeurl,
            type: "GET",
            dataType: "jsonp",  //指定服务器返回的数据类型
            jsonpCallback: "success",
            success: function (data) {
                if (data != null) {
                    svrTimeStmap = data.svr_time

                }

            }
        });
    }

    var res=checkMacInput()
    if (res)
    {
        isExistMac(function(){
            synSvrtime()
            pollDataFromRoute()
        })
    }

});

