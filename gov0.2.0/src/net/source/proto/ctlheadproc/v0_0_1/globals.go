package v0_0_1

import (
	"net/source/utils/bytes"
	"net/source/proto/defs"
)

type GlobalCtlHeader struct {
}

/*---------------------------控制域----------------------------------
机房id
服务器id
时间戳
加密类型编号
协议控制类型编号
业务领域类型
业务编号
业务协议体长度   <协议体可扩展>
------------------------------------------------------------------------------------
****/
//控制域信息


//控制域层
func (g *GlobalCtlHeader) Proc(ctlBytes []byte) (*defs.CtlMsg, error) {
	byteArray := bytes.NewByteArray(ctlBytes)
	msg := defs.CtlMsg{}

	roomId, err := byteArray.ReadInt32()
	if err == nil {
		msg.RoomId = roomId
	}

	svrId, err := byteArray.ReadInt32()
	if err == nil {
		msg.SvrId = svrId
	}

	timeStamp, err := byteArray.ReadInt64()
	if err == nil {
		msg.Timstamp = timeStamp
	}
	encryptType, err := byteArray.ReadByte()
	if err == nil {
		msg.EncryptType = int8(encryptType)
	}

	ctlOptType, err := byteArray.ReadByte()
	if err == nil {
		msg.CtlOptType = int8(ctlOptType)
	}

	domainType, err := byteArray.ReadInt16()
	if err == nil {
		msg.BusinessDomainType = domainType
	}

	businessMsgCode, err := byteArray.ReadInt32()
	if err == nil {
		msg.BusinessMsgCode = businessMsgCode
	}

	businessMsgLength, err := byteArray.ReadInt32()
	if err == nil {
		msg.BusinessMsgLength = businessMsgLength
	}

	return &msg, nil
}
