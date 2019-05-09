package login

import (
	"net/source/utils/bytes"
	"net/source/proto/defs"
	"net/source/msg/msgproc"
)

type RepoMsgProc struct {
	msgproc.MsgProc
}

/*
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
}
*/
//这里应该是垂直化的时候
func (this *RepoMsgProc) Proc(ctlMsg *defs.CtlMsg, businessBytes []byte) (interface{}, error) {
	reportMsg := RepoMsg{}
	reportMsg.CtlMsg = ctlMsg
	byteArray := bytes.NewByteArray(businessBytes)
	reportMsg.Ver = ctlMsg.Ver
	devId, err := byteArray.ReadInt64()
	if err == nil {
		reportMsg.DevId = devId
	}

	isCharge, err2 := byteArray.ReadByte()
	if err2 == nil {
		reportMsg.IsCharged = isCharge
	}

	name, err2 := byteArray.ReadUTF()
	if err2 == nil {
		reportMsg.Name = name
	}

	return &reportMsg, nil

}

func CreateRepoMsgProc() *RepoMsgProc {
	return &RepoMsgProc{}
}
