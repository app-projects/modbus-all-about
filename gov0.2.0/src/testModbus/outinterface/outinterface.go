package outinterface

import (
	"testModbus/simulator"
	"testModbus/connection"
	"time"
	"log"
	"sync"
)

type WebMsg struct {
	MsgType byte `json:"msg_type"`
	Mac     byte `json:"mac"`
	OffsetH byte `json:"offset_h"`
	OffsetL byte `json:"offset_l"`
	DataH   byte `json:"data_h"`
	DataL   byte `json:"data_l"`

	Handler func(args ...interface{}) interface{}

	IsReleaseNature  bool
}

const WEB_MSG_TYPE_MODIFY = 1
const WEB_MSG_TYPE_READ_REG = 2


var webMsgPool = sync.Pool{
	New: func() interface{} {
		return &WebMsg{IsReleaseNature:true}
	},
}

func NewWebMsg(natrueRelease bool) *WebMsg {
	m :=webMsgPool.Get().(*WebMsg)
	m.IsReleaseNature =natrueRelease
	 return m
}

func ReleaseWebMsg(m *WebMsg)  {
   if m!=nil{
	   webMsgPool.Put(m)
   }
}

var msgQueue chan *WebMsg

func init() {
	msgQueue = make(chan *WebMsg, 2000000)
}

func PushMsg(msg *WebMsg) int {
	c := connection.GetDevConnRouter().GetConnByMac(msg.Mac)
	if c == nil {
		return -2
	}
	if msg == nil {
		return -1
	}
	msgQueue <- msg
	return 0
}

var selectDelayMillSecTimeout time.Duration = time.Millisecond * 300

//如果来自页面的消息 很多，导致 积压队列非常长，那么可以考虑通过 msg.MsgType ,来分治 ，放入不同的线程
//目前先放入到同一个队列
func TickComingMsg() {
	var t = time.NewTimer(selectDelayMillSecTimeout)
	defer t.Stop()
	for {
		t.Reset(selectDelayMillSecTimeout)
		select {
		case msg := <-msgQueue:
			if (msg.MsgType == WEB_MSG_TYPE_READ_REG) {
				simulator.CommitReadRegAsk(msg.Handler, msg.Mac, msg.OffsetH, msg.OffsetL, msg.DataH, msg.DataL)
			} else if (msg.MsgType ==WEB_MSG_TYPE_MODIFY){
				simulator.CommitModifyAppInfoByMac(msg.Mac, msg.OffsetH, msg.OffsetL, msg.DataH, msg.DataL)
			}
			if msg.IsReleaseNature{
				ReleaseWebMsg(msg)
			}
		case <-t.C:
			 log.Println("TickComingMsg  tick ")
		}
		t.Stop()
	}
}
