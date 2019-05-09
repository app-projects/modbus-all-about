package connection

import (
	"net/source/userapi"
	"net"
	"net/source/utils/bytes"
)

type Connection struct {
	userapi.IClient
	Conn net.Conn
	Id   int64
	//*************解析包 过程中需要的数据 用于恢复作用*****************
	ToolBucket     []byte           //重用
	ToolTotalBytes []byte           //重用
	ToolTotalCache *bytes.ByteArray //操作缓冲    //cache 是一个环形的执行队列，r w ,是无止境的，是不需要重置的

	NewPackFlag bool

	isOpenFlag bool

	SvrCenter              userapi.IServiceIOCenter
	ModBusProtoBinPackChan chan userapi.IModBusProtoBinPack
}

func (this *Connection) IsOpen() bool {
	return this.isOpenFlag
}
func (c *Connection) GetToolBucket() []byte {
	return c.ToolBucket
}
func (c *Connection) GetToolTotalCache() *bytes.ByteArray {
	return c.ToolTotalCache
}

//善后工作
func (c *Connection) Exit() {
	//clear others

	c.isOpenFlag = false

	if c.SvrCenter != nil {
		c.SvrCenter.ClearClient(c)
	}

	//atomic.AddInt64(&clienGenator, ^int64(0))
	c.ResetRecvNewPack()
	c.Conn.Close()
	//fmt.Printf("\n iot client id :%d 断开连接。。存活连接数：%d\n", c.Id, clienGenator)
}

func (c *Connection) GetModBusProtoBinChan() chan userapi.IModBusProtoBinPack {
	return c.ModBusProtoBinPackChan
}

func (c *Connection) GetId() int64 {
	return c.Id
}
func (c *Connection) GetConn() net.Conn {
	return c.Conn
}

func (c *Connection) Send(bytes []byte) {
	c.Conn.Write(bytes)
}
func (c *Connection) SetSvrCenter(svr userapi.IServiceIOCenter) {
	c.SvrCenter = svr
}

func (c *Connection) ResetRecvNewPack() {
	c.NewPackFlag = true //开启下一个数据包接收
}
func (c *Connection) GetNewPackFlag() bool {
	return c.NewPackFlag
}

func (c *Connection) SetNewPackFlag(b bool) {
	c.NewPackFlag = b
}

func (c *Connection) SetClientId(id int64) {
	c.Id = id
}
func (c *Connection) GetClientId() int64 {
	return c.Id
}

func NewConnection(conn net.Conn, defaultPackCacheSize int) *Connection {
	c := Connection{}
	c.Conn = conn

	c.NewPackFlag = true
	c.isOpenFlag = true
	c.Id = 1
	c.ToolBucket = make([]byte, 256)
	c.ToolTotalBytes = make([]byte, defaultPackCacheSize*2)
	c.ModBusProtoBinPackChan = make(chan userapi.IModBusProtoBinPack, 1000)
	c.ToolTotalCache = bytes.NewByteArray(c.ToolTotalBytes)
	return &c
}
