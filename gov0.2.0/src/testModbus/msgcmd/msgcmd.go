package msgcmd

import "testModbus/echocontext"

type iNetworkCmd struct {
	Handler   echocontext.RespHandler
	RespArgs  []interface{}
	ReqPack   []byte
	Timestamp int64
	state int
}

func (this *iNetworkCmd) SetRespdHandler(h echocontext.RespHandler) {
	this.Handler = h
}

func (this *iNetworkCmd) GetRespdHandler() echocontext.RespHandler {
	return this.Handler
}

func (this *iNetworkCmd) GetTimestamp() int64 {
	return this.Timestamp
}

func (this *iNetworkCmd) SetTimestamp(t int64) {
	this.Timestamp = t
}

func (this *iNetworkCmd) GetPack() []byte {
	return this.ReqPack
}

func (this *iNetworkCmd) SetPack(pck []byte) {
	this.ReqPack = pck
}


func (this *iNetworkCmd) SetRspdArgs(args ...interface{}) {
	this.RespArgs = args
}
func (this *iNetworkCmd) DoResponse() {
	if nil != this.Handler {
		this.Handler(this.RespArgs...)
	}
}
func (this *iNetworkCmd) GetState() int {
	return this.state
}
func (this *iNetworkCmd) SetState(state int){
	  this.state =state
}



func (this *iNetworkCmd) GetSndPack() []byte {
	return this.ReqPack
}

func NewNetworkCmd() echocontext.NetworkCmd {
	cmd:=iNetworkCmd{}
	return &cmd
}
