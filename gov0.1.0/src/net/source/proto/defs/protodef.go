package defs

import (
	"net/source/utils/bytes"
	"fmt"
	"net/source/proto/pools"
)

// msgver   protover(unpack protoversion)  handler(parse msg handler)
//  1                 1                           handler1
//  3                 1                           handler3

type ProtoHead struct {
	//协议头：
	Ver     int32 //协议版本号                     2byte      (必有字段)
	BodyLen int32
}

//协议的服务器组织，采用builder模式实现
type ProtoProcessor struct {
	Ver           int32 //协议版本号
	CtlHeaderProc ProtoCtlHeadProc
	BodyProc      ProtoBodyProc
}

type ProtoCtlHeadProc interface {
	Proc(ctlbytes []byte) (*CtlMsg, error)
}
type ProtoBodyProc interface {
	Proc(ctlMsg *CtlMsg, bodyBytes []byte) error
}

//------------------------prehead---------------------------
//指令头
//内容长度
/*---------------------------控制域CtlMsg----------------------------------
控制域长度
机房id
服务器id
时间戳
加密类型编号
协议控制类型编号
业务领域类型
业务领域编号
业务数据长度
------------------------------------------------------------------------------------

协议体：(领域业务消息)
----------------------------------业务域数据-BusinessMsg------------------------------------*/

type CtlMsg struct {
	Ver         int32
	RoomId      int32
	SvrId       int32
	Timstamp    int64
	EncryptType int8
	CtlOptType  int8

	BusinessDomainType int16
	BusinessMsgCode    int32
	BusinessMsgLength  int32

	ConnId int64 //客户端连接
}

//ProtoProcessor 是 就像网卡出来的数据，决定送到那一层处理
//协议版本 分道扬镳的地方，分层的关键地方  跳转 转发 路由的地方，分发 调度的地方

//层间路由 总控制器

var  dstBytes = [1024]byte{0}
func (this *ProtoProcessor) DepartDomain(ver int32, cliId int64, contentCtlAndBusiness []byte) {
	var byteArray = bytes.NewByteArray(contentCtlAndBusiness)
	ctlLen, err1 := byteArray.ReadInt32()
	if err1 == nil {
		if ctlLen > 0 {
			ctlBytes := dstBytes[0:ctlLen]
			byteArray.Read(ctlBytes)

			// ctl progress  //协议控制域层
			ctlMsg, err := this.CtlHeaderProc.Proc(ctlBytes)
			ctlMsg.Ver = ver
			ctlMsg.ConnId = cliId
			if err == nil {

				// next progress2 bodyProc   //协议业务域层
				copyDst := pools.BusinessBytesSlicePool.Get().([]byte)
				bodyBytes := copyDst[0:ctlMsg.BusinessMsgLength] //截取 业务真实的长度
				byteArray.Read(bodyBytes)
				this.BodyProc.Proc(ctlMsg, bodyBytes)
				pools.BusinessBytesSlicePool.Put(copyDst)
			}
		}
	} else {
		fmt.Errorf("DepartDomain error ctlLen get:", err1)
	}

}
