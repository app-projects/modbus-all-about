package connection

import (
	"net/source/userapi"
	"net"
	"net/source/utils/bytes"
	"testModbus/echocontext"
	"fmt"
	"sync"
	"time"

	"log"
	"sync/atomic"
	"testModbus/utils"
)

type Connection struct {
	userapi.IClient
	Conn net.Conn
	Id   int64
	//*************解析包 过程中需要的数据 用于恢复作用*****************
	ToolBucket      []byte           //重用
	ToolTotalBytes  []byte           //重用
	ToolTotalCache  *bytes.ByteArray //操作缓冲    //cache 是一个环形的执行队列，r w ,是无止境的，是不需要重置的
	packGlobalCount int64
	NewPackFlag     bool

	isOpenFlag bool

	SvrCenter              userapi.IServiceIOCenter
	ModBusProtoBinPackChan chan userapi.IModBusProtoBinPack

	CommitCmdChan chan echocontext.NetworkCmd

	networkEchoCtx *echocontext.EchoContext
}

func (this *Connection) GetClientProtoBinChan() chan userapi.IModBusProtoBinPack {
	return this.ModBusProtoBinPackChan
}

//resendFlag 一个读 一个写 不需要加锁
func (this *Connection) PopModBusProtoBin() (userapi.IModBusProtoBinPack, error) {
	return <-this.ModBusProtoBinPackChan, nil
}

//什么是系统控制 有什么是外部发起控制；
//响应是系统内部控制，所以
//输入时外部控制，外部的环境是可能是并发的
//下面这个函数 是否线程安全的，因为外界有多个 线程对 同一个Connection 对象，故要加锁

var selectDelayMillSecTimeout time.Duration = time.Millisecond * 200

func (this *Connection) CommitReqPck(cmd echocontext.NetworkCmd) bool {
	this.CommitCmdChan <- cmd

	/*	var t = time.NewTimer(selectDelayMillSecTimeout)
		defer t.Stop()
		for {
			select {
			case this.CommitCmdChan <- cmd:
				break
			}
			select {
			case <-t.C:
				log.Println("case <-tick.C....")
			default:
				log.Println("no block")
			}
		}*/

	//	fmt.Println("CommitCmdChan 当前提交积压队列个数是：", len(this.CommitCmdChan))
	return true
}

func (this *Connection) SynRequestQueue() {
	var t = time.NewTimer(selectDelayMillSecTimeout)
	defer t.Stop()
	for {
		if !this.IsOpen() {
			//to do 关闭管道
			return //终端关闭状态
		}

		t.Reset(selectDelayMillSecTimeout)
		select {
		case cmd := <-this.CommitCmdChan:
			this.networkEchoCtx.CommitRequest(cmd)
		}
		select {
		case <-t.C:
			log.Println("SynRequestQueue time out tick...")
		default:
			log.Println("SynRequestQueue no block ..")
		}
		t.Stop()
	}
	log.Println("out of SynRequestQueue....")
}

func (this *Connection) AskRoutine() {
	this.networkEchoCtx.AskRoute()
}

func (this *Connection) ExecuteResponse() {
	this.networkEchoCtx.ExecuteResponse()
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
	c.isOpenFlag = false
	if c.SvrCenter != nil {
		c.SvrCenter.ClearClient(c)
	}
	c.networkEchoCtx.SetState(echocontext.ECHO_CTX_STATE_CLOSE)
	c.ResetRecvNewPack()
	subUid()
	c.Conn.Close()
	GetDevConnRouter().Unbind(utils.Int64_2Byte(c.GetId()))
}

func (c *Connection) PushProtoBin(bin userapi.IModBusProtoBinPack) {
	c.networkEchoCtx.PushNetworkData(bin)
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
func (c *Connection) SvrSend(bytes []byte) {
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
func (c *Connection) SendNetworkFn(cmd echocontext.NetworkCmd) {
	c.SvrSend(cmd.GetSndPack())
	//fmt.Println("服务器发送一个命令........xxxxx.........")
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

var TIMEOUT_MillSEC int64 = 10 * 1000

type Dev2ConnRouteTable struct {
	mac2ConnTable sync.Map
}

func (this *Dev2ConnRouteTable) GetConnByMac(mac byte) *Connection {

	var v, ok = this.mac2ConnTable.Load(mac)
	if ok {
		return v.(*Connection)
	}
	return nil
}

var dev2ConnRouteTable Dev2ConnRouteTable

func GetDevConnRouter() *Dev2ConnRouteTable {
	return &dev2ConnRouteTable
}

func (this *Dev2ConnRouteTable) Unbind(mac byte) {
	this.mac2ConnTable.Delete(mac)
}

func (this *Dev2ConnRouteTable) Bind(mac byte, conn *Connection) {
	var v, ok = this.mac2ConnTable.Load(mac)
	this.mac2ConnTable.Store(mac, conn)
	if ok {
		fmt.Println("设备的mac=", mac, " 和连接通道的关系有更新...,由 ", v.(*Connection).Id, "--->变为:", conn.Id)
		return
	}
	fmt.Println("设备的mac=", mac, " 和连接通道建立关系...---> 为:", conn.Id)
}

var clientUid int64

func getUId() int64 {
	return atomic.AddInt64(&clientUid, 1)
}

func subUid() {
	atomic.AddInt64(&clientUid, -1)
}

func NewConnection(conn net.Conn, defaultPackCacheSize int) *Connection {
	c := Connection{}
	c.NewPackFlag = true
	c.Conn = conn
	c.CommitCmdChan = make(chan echocontext.NetworkCmd, 1000000)
	c.isOpenFlag = true

	c.Id = getUId()
	c.ToolBucket = make([]byte, 256)
	c.ToolTotalBytes = make([]byte, defaultPackCacheSize*2)
	c.ModBusProtoBinPackChan = make(chan userapi.IModBusProtoBinPack)
	c.ToolTotalCache = bytes.NewByteArray(c.ToolTotalBytes)
	c.ToolTotalCache.Seek(0)
	fmt.Println("length valid", c.ToolTotalCache.Available())

	c.networkEchoCtx = echocontext.NewEchoContext()
	c.networkEchoCtx.SendNetworkFn = c.SendNetworkFn
	c.networkEchoCtx.SetTimeoutMillSec(TIMEOUT_MillSEC)

	if res := c.networkEchoCtx.Init(); res != 0 {
		fmt.Println("init networkEchoCtx failed ")
		return nil
	}

	return &c
}
