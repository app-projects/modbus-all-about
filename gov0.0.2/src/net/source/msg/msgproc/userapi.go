package msgproc

import (
	"net/source/proto/defs"

	"net/source/userapi"
	"fmt"
)

//type CollectionBytesFn func()[]byte

type MsgPulisher interface {
	Commit(msg *defs.CtlMsg, sndBytes []byte)
	Pubish(msg *defs.CtlMsg, sndBytes []byte)
}

type UserApi struct {
	MsgPulisher
	SvrIO userapi.IServiceIOCenter
 	UserCreator userapi.UserCreator
}

type appMsgPuber struct {
}

func (this *appMsgPuber) Commit(msg *defs.CtlMsg, sndBytes []byte) {
	 cliconn:= GetAppTools().SvrIO.FindUser(msg.ConnId)
	 fmt.Println(msg.ConnId)
	 if cliconn!=nil{
	 	cliconn.Send(sndBytes)
	 }
}

func (this *appMsgPuber) Pubish(msg *defs.CtlMsg, sndBytes []byte) {

}

var tools *UserApi

func init() {
	tools = &UserApi{MsgPulisher: &appMsgPuber{}}
}

func GetAppTools() *UserApi {
	return tools
}
