package outputconfig

//linux
//数据仓库restful api

/* 内网地址*/
var RemoteStorePushQueryAddr string = "http://172.16.169.186:8501/dev/log"
var RemoteStoreModifyAddr string = "http://172.16.169.186:8501/dev/log/modify" //06号

/*var RemoteStorePushQueryAddr string = "http://47.110.78.124:8501/dev/log"
var RemoteStoreModifyAddr string = "http://47.110.78.124:8501/dev/log/modify" //06号*/

/* local 地址*/
/*var RemoteStorePushQueryAddr string = "http://192.168.1.102:8102/dev/log"
var RemoteStoreModifyAddr string = "http://192.168.1.102:8102/dev/log/modify"*/

//var pushAddrFmt = "http://%s:%d//dev/svr/msg/push"

//remote
var RemotePushSvrMsgIp = "172.17.0.2"
var RemotePushSvrMsgPort = 8520

//local
//http://192.168.1.102:8520/dev/svr/msg/push
//http://192.168.1.102:8520/nowtime


/*var RemotePushSvrMsgIp = "192.168.1.102"
var RemotePushSvrMsgPort = 8522*/
