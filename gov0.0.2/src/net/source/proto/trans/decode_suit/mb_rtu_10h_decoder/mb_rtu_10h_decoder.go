package mb_rtu_10h_decoder

import (
	"net/source/userapi"
	"net/source/proto/trans/interfaces"
	"sync"
)

type modbus_rtu_10_decoder struct {
}

var selfPool = sync.Pool{
	New: func() interface{} {
		return &modbus_rtu_10_decoder{}
	},
}

func CreateIntance() interfaces.Decoder {
	return selfPool.Get().(*modbus_rtu_10_decoder)
}

func ReleaseInstance(ins interfaces.Decoder) {
	selfPool.Put(ins)
}

/***
Modbus  协议分段 分区格式：

<指令头>  mac,code    int16
<content>
<crc16>      int16

命令10H      （修改多个寄存器） modify multi(就像http协议有  put delete)

发送命令：
[设备地址] + [命令号10H] +    <指令头>

[起始寄存器地址高8位] + [低8位] +
[寄存器数高8位] + [低8位] +

[寄存器字节数] +

[数据1高8位] + [低8位]
+…. +
[数据N高8位] + [低8位] +

[CRC校验的低8位] + [CRC校验的高8位]


设备响应：
如果成功把计算机发送的命令原样返回，
否则不响应

03H协议请求特征是固定的：
1 / 总字节是， 8byte

*/
var pack_content_total_size = 6 //03H包体  除了指令头 2byte

var contentBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, pack_content_total_size)
	},
}

func (this *modbus_rtu_10_decoder) Decode(c userapi.IClient) int {

	return -1

}
