// 基于准备好的dom，初始化echarts实例
var myChart = echarts.init(document.getElementById('templ'));
// 指定图表的配置项和数据

var opt1 = {
    tooltip: {
        formatter: "{a} <br/>{b} : {c}%"
    },
    toolbox: {
        feature: {
            restore: {},
            saveAsImage: {}
        }
    },
    series: [{
        name: '业务指标',
        type: 'gauge',

        title: {
            textStyle: {
                fontWeight: 'bolder',
                fontSize: 20,
                /*fontStyle: 'italic',*/
                color: "#25c36c"
            }
        },

        axisLabel: {            // 坐标轴小标记
            textStyle: {       // 属性lineStyle控制线条样式
                /*  color:"red",*/
                fontSize: 16,   //改变仪表盘内刻度数字的大小
                shadowColor: '#000', //默认透明
                fontWeight: 'bolder',
            }
        },

        detail: {formatter: '{value}%'},
        data: [{value: 0, name: '实时温度'}]
    }]
}

myChart.setOption(opt1);


//-----------------------------------------------------------


var myChart2 = echarts.init(document.getElementById('waterMark'));

// 指定图表的配置项和数据

var opt2 = {
    tooltip: {
        formatter: "{a} <br/>{b} : {c}%"
    },
    toolbox: {
        feature: {
            restore: {},
            saveAsImage: {}
        }
    },
    series: [{
        name: '业务指标',
        type: 'gauge',
        title: {
            textStyle: {
                fontWeight: 'bolder',
                fontSize: 20,
                /*fontStyle: 'italic',*/
                color: "#25c36c"
            }
        },
        axisLabel: {            // 坐标轴小标记
            textStyle: {       // 属性lineStyle控制线条样式
                /*  color:"red",*/
                fontSize: 16,   //改变仪表盘内刻度数字的大小
                shadowColor: '#000', //默认透明
                fontWeight: 'bolder',
            }
        },
        detail: {formatter: '{value}%'},
        data: [{value: 0, name: '实时水位'}]
    }]
}
myChart2.setOption(opt2);


//-----------------------------------------------------------


var myChart3 = echarts.init(document.getElementById('wet'));
// 指定图表的配置项和数据

var opt3 = {
    tooltip: {
        formatter: "{a} <br/>{b} : {c}%"
    },
    toolbox: {
        feature: {
            restore: {},
            saveAsImage: {}
        }
    },
    series: [{
        name: '业务指标',
        type: 'gauge',
        title: {
            textStyle: {
                fontWeight: 'bolder',
                fontSize: 20,
                /*fontStyle: 'italic',*/
                color: "#25c36c"
            }
        },
        axisLabel: {            // 坐标轴小标记
            textStyle: {       // 属性lineStyle控制线条样式
                /*  color:"red",*/
                fontSize: 16,   //改变仪表盘内刻度数字的大小
                shadowColor: '#000', //默认透明
                fontWeight: 'bolder',
            }
        },
        detail: {formatter: '{value}%'},
        data: [{value: 0, name: '实时湿度'}]
    }]
}
myChart3.setOption(opt3);


//elect
//-----------------------------------------------------------


var electC = echarts.init(document.getElementById('elect'));
// 指定图表的配置项和数据

var electOpt = {
    tooltip: {
        formatter: "{a} <br/>{b} : {c}%"
    },
    toolbox: {
        feature: {
            restore: {},
            saveAsImage: {}
        }
    },
    series: [{
        name: '业务指标',
        type: 'gauge',
        title: {
            textStyle: {
                fontWeight: 'bolder',
                fontSize: 20,
                /*fontStyle: 'italic',*/
                color: "#25c36c"
            }
        },
        axisLabel: {            // 坐标轴小标记
            textStyle: {       // 属性lineStyle控制线条样式
                /*  color:"red",*/
                fontSize: 16,   //改变仪表盘内刻度数字的大小
                shadowColor: '#000', //默认透明
                fontWeight: 'bolder',
            }
        },
        detail: {formatter: '{value}%'},
        data: [{value: 0, name: '实时电流'}]
    }]
}
electC.setOption(electOpt);


//light
//-----------------------------------------------------------


var lightC = echarts.init(document.getElementById('light'));
// 指定图表的配置项和数据

var lightOpt = {
    tooltip: {
        formatter: "{a} <br/>{b} : {c}%"
    },
    toolbox: {
        feature: {
            restore: {},
            saveAsImage: {}
        }
    },
    series: [{
        name: '业务指标',
        type: 'gauge',
        title: {
            textStyle: {
                fontWeight: 'bolder',
                fontSize: 20,
                /*fontStyle: 'italic',*/
                color: "#25c36c"
            }
        },
        axisLabel: {            // 坐标轴小标记
            textStyle: {       // 属性lineStyle控制线条样式
                /*  color:"red",*/
                fontSize: 16,   //改变仪表盘内刻度数字的大小
                shadowColor: '#000', //默认透明
                fontWeight: 'bolder',
            }
        },
        detail: {formatter: '{value}%'},
        data: [{value: 0, name: '实时亮度'}]
    }]
}
lightC.setOption(lightOpt);


//pressValue


var pressC = echarts.init(document.getElementById('press'));
// 指定图表的配置项和数据

var pressOpt = {
    tooltip: {
        formatter: "{a} <br/>{b} : {c}%"
    },
    toolbox: {
        feature: {
            restore: {},
            saveAsImage: {}
        }
    },
    series: [{
        name: '业务指标',
        type: 'gauge',
        title: {
            textStyle: {
                fontWeight: 'bolder',
                fontSize: 20,
                /*fontStyle: 'italic',*/
                color: "#25c36c"
            }
        },
        axisLabel: {            // 坐标轴小标记
            textStyle: {       // 属性lineStyle控制线条样式
                /*  color:"red",*/
                fontSize: 16,   //改变仪表盘内刻度数字的大小
                shadowColor: '#000', //默认透明
                fontWeight: 'bolder',
            }
        },
        detail: {formatter: '{value}%'},
        data: [{value: 0, name: '实时压强'}]
    }]
}
pressC.setOption(pressOpt);






