package msgproc

import (
	"net/source/proto/defs"
)

/*---------------------------控制域----------------------------------
机房id
服务器id
时间戳
加密类型编号
协议控制类型编号
业务领域编号
业务协议体长度   <协议体可扩展>
------------------------------------------------------------------------------------
**/

//数据的上层体现,按照实际情
type BaseMsg struct {
	MsgId     int32
	Ver       int32
	Timestamp int64
	CtlMsg    *defs.CtlMsg
}
