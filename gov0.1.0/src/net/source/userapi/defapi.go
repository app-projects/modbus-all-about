package userapi

import (
	"net"
	"net/source/utils/bytes"
	"testModbus/echocontext"
)

/*
type PackBinHandler = func(modbusBinPack IModBusProtoBinPack)

type NetworkCmd interface {
	GetTimestamp() int64
	SetTimestamp(t int64)
	SetRspdArgs(args ...interface{})
	DoResponse()
	GetSndPack() []byte
}

*/


//一般不要返回系统的组织给应用层，因为它们并不关心，它们关心的是数据；而对于系统来说，更需要返回系统间的接口
type IClient interface {
	Exit()
	ResetRecvNewPack()
	GetId() int64
	GetConn() net.Conn
	GetToolBucket() []byte
	GetToolTotalCache() *bytes.ByteArray
	SetSvrCenter(svr IServiceIOCenter)
	Send([]byte)
	IsOpen() bool
	GetNewPackFlag() bool
	SetNewPackFlag(b bool)
    PushProtoBin(IModBusProtoBinPack)

    //client api
     PopModBusProtoBin() (IModBusProtoBinPack, error)
     GetClientProtoBinChan () chan IModBusProtoBinPack
}

type IParse interface {
	Parse() error
}

type IProtoBinPack interface {
	IParse
	Len() int
	GetBytes() []byte
	SetRLen(l int32)
}

type IModBusProtoBinPack interface {
	echocontext.NetworkData

	SetClientId(id int64)
	GetClientId() int64
	SetFnCode(code byte)
	GetFnCode() byte
	SetMac(mac byte)
	GetMac() byte
	Reset()
}

type UserCreator interface {
	InitClient(conn net.Conn, size int) IClient
}

type IServiceIOCenter interface {
	ClearClient(c IClient) error
	FindUser(uid int64) IClient
}
