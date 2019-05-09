package initproto

import (
	"fmt"
	"net/source/proto/bodyproc/v0_0_1"
	v0_0_12 "net/source/proto/ctlheadproc/v0_0_1"
	"net/source/proto"
)

const (
	VERSION_1 = 1 + iota
	VERSION_2
	VERSION_3
	VERSION_4
	VERSION_5
)

func init() {
	fmt.Println("模板启动了....")
	var gbProc = &v0_0_1.GlobalBody{}
	var ghProc = &v0_0_12.GlobalCtlHeader{}
	//add line processor
	proto.AssembleProtoProc(VERSION_1, ghProc, gbProc)
}
