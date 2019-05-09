package toclientmsgbuilder

import (
	"net/source/app/modules/login"
	"net/source/utils/bytes"
)

/*

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

/****
namesapce msg{
type BaseMsg struct {
MsgId int32
Ver       int32
Timestamp int64
ctlMsg *proto.CtlMsg

}
}


类型 c2s
编号 10
type RepoMsg struct {
	msg.BaseMsg  //extends
	DevId     int64
	IsCharged byte
	Name      string
}*/

var CRC32 int32 = 0xfff5

const SIZE_CRC32 = 4

func builderPreHead(m *login.RepoMsg) []byte {
	var bytess = make([]byte, 0)
	byteArray := bytes.NewByteArray(bytess)
	byteArray.WriteInt32(m.Ver)
	return byteArray.Bytes()
}

func builderCtls(msg *login.RepoMsg, businessBytes []byte) []byte {
	ctlMsgBytes := make([]byte, 0)
	ctlByteArray := bytes.NewByteArray(ctlMsgBytes)
	ctlByteArray.WriteInt32(msg.CtlMsg.RoomId)
	ctlByteArray.WriteInt32(msg.CtlMsg.SvrId)
	ctlByteArray.WriteInt64(msg.CtlMsg.Timstamp)
	ctlByteArray.WriteByte(byte(msg.CtlMsg.EncryptType))
	ctlByteArray.WriteByte(byte(msg.CtlMsg.CtlOptType))

	ctlByteArray.WriteInt16(msg.CtlMsg.BusinessDomainType)
	ctlByteArray.WriteInt32(msg.CtlMsg.BusinessMsgCode)

	msg.CtlMsg.BusinessMsgLength = int32(len(businessBytes))

	//下层域描述长度
	ctlByteArray.WriteInt32(msg.CtlMsg.BusinessMsgLength)

	//上面业务长度
	ctlLen := ctlByteArray.Length()
	//返回
	ret := bytes.NewByteArray(make([]byte, 0))
	ret.WriteInt32(int32(ctlLen))
	ret.WriteBytes(ctlByteArray.Bytes())

	return ret.Bytes()

}

func builderMsg(msg *login.RepoMsg) []byte {
	businessBytes := make([]byte, 0)
	businessByteArray := bytes.NewByteArray(businessBytes)
	businessByteArray.WriteInt64(msg.DevId)
	businessByteArray.WriteByte(msg.IsCharged)
	businessByteArray.WriteUTF(msg.Name)

	return businessByteArray.Bytes()
}

func TotalBuilder(msg *login.RepoMsg) []byte {

	businessBytes := builderMsg(msg)
	ctlBytes := builderCtls(msg, businessBytes)
	preheadBytes := builderPreHead(msg)

	//写入消息
	totalBytes := make([]byte, 0)
	totalBty := bytes.NewByteArray(totalBytes)

	totalBty.WriteBytes(preheadBytes)

	var contentLen = int32(len(ctlBytes) + len(businessBytes))
	//内容长度
	totalBty.WriteInt32(contentLen)

	//控制域所有字段写入 ，包含了业务域的长度
	totalBty.WriteBytes(ctlBytes)

	//业务域所有字段写入
	totalBty.WriteBytes(businessBytes)

	return totalBty.Bytes()
}
