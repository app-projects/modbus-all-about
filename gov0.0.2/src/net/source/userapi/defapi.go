package userapi

import (
	"net"
	"net/source/utils/bytes"
)

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
	GetModBusProtoBinChan() chan IModBusProtoBinPack

}


type IParse interface{
	Parse() error
}

type IProtoBinPack interface {
	IParse
	Len() int
	GetBytes()[]byte
	SetRLen(l int32)

}

type IModBusProtoBinPack interface {
	SetClientId(id int64)
	GetClientId()int64
	SetFnCode(code byte)
	GetFnCode()byte
	SetMac(mac byte)
	GetMac()byte
    Reset()
}


type UserCreator interface {
	InitClient(conn net.Conn, size int) IClient
}

type IServiceIOCenter interface {
	ClearClient(c IClient) error
	FindUser(uid int64) IClient
}
