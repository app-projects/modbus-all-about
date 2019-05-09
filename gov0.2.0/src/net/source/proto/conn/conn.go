package conn

import (
	"net"
	"sync/atomic"
	"fmt"
	"net/source/userapi"
	"net/source/msg/msgproc"
	"net/source/utils/bytes"
)

var clienGenator int64 = 0

type Cli struct {
	userapi.IClient
	Conn     net.Conn
	Id       int64
	BufCache *bytes.SafeBuff
	//*************解析包 过程中需要的数据 用于恢复作用*****************
	ToolBucket     []byte           //重用
	ToolTotalBytes []byte           //重用
	ToolTotalCache *bytes.ByteArray //操作缓冲    //cache 是一个环形的执行队列，r w ,是无止境的，是不需要重置的

	NewPackFlag bool
	isOpenFlag bool

	//
	SvrCenter userapi.IServiceIOCenter

	ModBusProtoBinPackChan chan userapi.IModBusProtoBinPack
	/*	bucket := make([]byte, 256) //桶子
		totalBytes := make([]byte, 1024)
		totalBuf := bytes.NewBuffer(totalBytes)*/
}

func newClientConnection(defaultPackCacheSize int) *Cli {
	c := Cli{}
	c.Id = clienGenator
	c.ToolBucket = make([]byte, 256)
	c.ToolTotalBytes = make([]byte, defaultPackCacheSize*2)
	c.ModBusProtoBinPackChan = make(chan userapi.IModBusProtoBinPack, 1000)

	c.BufCache = bytes.NewBuffer()
	c.ToolTotalCache = bytes.NewByteArray(c.ToolTotalBytes)
	atomic.AddInt64(&clienGenator, 1)
	return &c
}

func init() {
	msgproc.GetAppTools().UserCreator = &UserConnCreator{}
}

type UserConnCreator struct {
}

func (*UserConnCreator) InitClient(conn net.Conn, size int) userapi.IClient {
	c := newClientConnection(size)
	//initConnection(conn)
	c.Conn = conn
	c.isOpenFlag = true
	return c
}
func (c *Cli) GetToolBucket() []byte {
	return c.GetToolBucket()
}
func (c *Cli) GetToolTotalCache() *bytes.ByteArray {
	return c.ToolTotalCache
}

//善后工作
func (c *Cli) Exit() {
	//clear others

	c.isOpenFlag = false

	if c.SvrCenter != nil {
		c.SvrCenter.ClearClient(c)
	}

	atomic.AddInt64(&clienGenator, ^int64(0))
	c.ResetRecvNewPack()
	c.Conn.Close()
	fmt.Printf("\n iot client id :%d 断开连接。。存活连接数：%d\n", c.Id, clienGenator)
}

func (c *Cli) GetModBusProtoBinChan() chan userapi.IModBusProtoBinPack {
	return c.ModBusProtoBinPackChan
}

func (c *Cli) GetId() int64 {
	return c.Id
}
func (c *Cli) GetConn() net.Conn {
	return c.Conn
}

func (c *Cli) Send(bytes []byte) {
	c.Conn.Write(bytes)
}
func (c *Cli) SetSvrCenter(svr userapi.IServiceIOCenter) {
	c.SvrCenter = svr
}

func (c *Cli) ResetRecvNewPack() {
	c.NewPackFlag = true //开启下一个数据包接收
}
func (c *Cli) GetNewPackFlag() bool {
	return c.NewPackFlag
}

func (c *Cli) SetNewPackFlag(b bool) {
	c.NewPackFlag = b
}


func (c *Cli) SetClientId(id int64){
	c.Id = id

}
func (c *Cli)  GetClientId()int64{
	return c.Id
}