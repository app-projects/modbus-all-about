package balance

import (
	"sync"
	"testModbus/connection"
)

type VirtualConnection struct{
	realConn *connection.Connection
	macDict  sync.Map
}

func (this *VirtualConnection)SetRConn(rcon *connection.Connection)  {
	this.realConn = rcon
}

type BalanceMac2Conn struct{
	connGroup []*VirtualConnection
	Size int
}



func newBalanceMac2Conn(connPoolSize int)  {
	ins:=BalanceMac2Conn{}
	ins.Size =connPoolSize
	ins.connGroup = make([]*VirtualConnection,connPoolSize)

}

/*
var dev2ConnRouteTable *Dev2ConnRouteTable
func newDev2ConnRouteTable() * Dev2ConnRouteTable{
	return &Dev2ConnRouteTable{}
}

func init()  {
	dev2ConnRouteTable = newDev2ConnRouteTable()

}

func GetDevConnRouter() * Dev2ConnRouteTable {
	return dev2ConnRouteTable
}
*/


