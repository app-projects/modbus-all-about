package main

import (
	"runtime"
	_ "net/source/msg/init_reg_msg_proc"
	"net/source/proto/area"
)

var (
	IP       = "192.168.1.120"
	BASEPORT = 9000

	PACKCACHE = 1024
)

const SERVICE_AREA_WRITE_NUM = 1
const SERVICE_AREA_READ_NUM =  1

var SERVICE_LISTNER_PORT_MAX = 1

func main1() {

	IP = "192.168.1.120"
	BASEPORT = 9000
	SERVICE_LISTNER_PORT_MAX = 1

	runtime.GOMAXPROCS(runtime.NumCPU())
	area.AcceptorServer(IP, BASEPORT, SERVICE_LISTNER_PORT_MAX)
	select {}

}
