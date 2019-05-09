package init_reg_msg_proc

import (
	"net/source/app/modules/login"
	"net/source/msg/msgproc"
	_ "net/source/proto/initproto"

	"fmt"
	"net/source/app/toclientmsgbuilder"
)

const (
	C2S_LOGIN_REPORT_INFO = 21
)

func init() {
	msgproc.RepertoryPut(C2S_LOGIN_REPORT_INFO, login.CreateRepoMsgProc())

	caller := msgproc.CreateMsgCaller(C2S_LOGIN_REPORT_INFO, "hellworld", func(msg interface{}, tools *msgproc.UserApi) {
		s := msg.(*login.RepoMsg)
		if s.CtlMsg.ConnId%1000==0{
			fmt.Println("成功接收到了设备信息devID：", s.DevId)
			fmt.Println("成功接收到了设备信息devName", s.Name)
			bdmsg:=toclientmsgbuilder.TotalBuilder(s)
			tools.Commit(s.CtlMsg,bdmsg)
		}

	})
	msgproc.MsgHandlerAdd(C2S_LOGIN_REPORT_INFO, caller)
}
