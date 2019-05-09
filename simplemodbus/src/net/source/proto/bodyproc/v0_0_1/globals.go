package v0_0_1

import (
	"net/source/utils/bytes"
	"net/source/msg/msgproc"
	"net/source/proto/defs"
)

type GlobalBody struct {
}

//body 业务转发层
func (g *GlobalBody) Proc(ctlMsg *defs.CtlMsg, bodyBytes []byte) error {

	var byteArray = bytes.NewByteArray(bodyBytes)
	var businessMsgLen = ctlMsg.BusinessMsgLength
	var businessMsgCode = ctlMsg.BusinessMsgCode

	if businessMsgLen > 0 {
		businessBytes := make([]byte, businessMsgLen)
		//copy 数据数据出来
		byteArray.ReadBytes(businessBytes,int(businessMsgLen), byteArray.GetReadPos())
		proc := msgproc.RepertoryGet(businessMsgCode)
		msg, err := proc.Proc(ctlMsg, businessBytes)
		if err == nil {
			// find handler
			dispather:=msgproc.MsgHandlerGet(ctlMsg.BusinessMsgCode)
			if dispather!=nil{
				dispather.CommitMsg(msg)
			}
		}
	}

	return nil
}
