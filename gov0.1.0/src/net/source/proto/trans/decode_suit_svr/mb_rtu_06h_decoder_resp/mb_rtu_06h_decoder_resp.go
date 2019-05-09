package mb_rtu_06h_decoder_resp

import (
	"net/source/userapi"
	"net/source/proto/trans/interfaces"
	"sync"
	"net/source/proto/binfiles"
	"net/source/proto/constant/modbus"
	"net/source/proto/pools"
	"net/source/utils/bytes"
	"testModbus/utils"
	"net/source/proto/trans/errcode"
	"fmt"
	"time"
)

type modbus_rtu_06_decoder_resp struct {
}

var selfPool = sync.Pool{
	New: func() interface{} {
		return &modbus_rtu_06_decoder_resp{}
	},
}

func CreateIntance() interfaces.Decoder {
	return selfPool.Get().(*modbus_rtu_06_decoder_resp)
}

func ReleaseInstance(ins interfaces.Decoder) {
	selfPool.Put(ins)
}

/***
Modbus  协议分段 分区格式：

<指令头>  mac,code    int16
<content>
<crc16>      int16

命令06H     （修改单个寄存器）   modify one  (就像http协议有put )

发送命令：

[设备地址] + [命令号06H] +       <指令头>   2byte

[需下置的寄存器地址高8位] + [低8位] +      2byte
[下置的数据高8位] + [低8位] +             2byte
[CRC校验的低8位] + [CRC校验的高8位]       2byte


设备响应：
如果成功把计算机发送的命令原样返回，
否则不响应

06H协议请求特征是固定的：
1 / 总字节是， 8byte

*/
var pack_all_total_size = 8 //03H包体  除了指令头 2byte

var allPackBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, pack_all_total_size)
	},
}

func (this *modbus_rtu_06_decoder_resp) Decode(c userapi.IClient) int {
	if c.GetToolTotalCache().Available() < pack_all_total_size {
		return errcode.ERR_TRNAS_DECODE_TRY_AGAIN
	}

	//copy data 通过上面获得 ver protolen 把数据 幺出来
	var availableLen = c.GetToolTotalCache().Available()
	if availableLen < 8 {
		return errcode.ERR_TRNAS_DECODE_TRY_AGAIN //整体的信息，但是没有整体的内容
	}

	var retCode = errcode.ERR_TRNAS_DECODE_TRY_AGAIN

	fmt.Println("enter:",c.GetToolTotalCache().Available())


	//预测处理

	var allPackLength = 8 //headcmd 2 + crc2+ data2 reg2

	canRead := (availableLen >= allPackLength) //2 是crc所需要的占位
	if !canRead {
		return errcode.ERR_TRNAS_DECODE_TRY_AGAIN //整体的信息，但是没有整体的内容
	}

/*	var allPackBytesRf = allPackBytesPool.Get().([]byte)
	allPackBytes := allPackBytesRf[:allPackLength]*/
     var allPackBytes = make([]byte,allPackLength)
	_, e := c.GetToolTotalCache().Read(allPackBytes)



	var byteArray *bytes.ByteArray
	if e == nil {
		byteArray = pools.ByteArrayPool.Get().(*bytes.ByteArray) //借出
		byteArray.Reset()
		byteArray.WriteBytes(allPackBytes)

		var res = utils.CheckContentCRC16OK(byteArray.Bytes())

		if !res {
			goto willExit
		}
		// 取值
		mac, _ := byteArray.ReadByte()
		byteArray.ReadByte() //fnCode

		var modbusBin = binfiles.CreateModBusProtoBin(mac, modbus.FunCode06, c.GetId())
		modbus06Bin := modbusBin.(*binfiles.Mod06_ProtoBinPack)
		modbus06Bin.EndTimestamp = time.Now().Unix()*1000
		h8, err := byteArray.ReadByte()
		if err == nil {
			modbus06Bin.RegSetterH8 = h8
		}

		l8, err := byteArray.ReadByte()
		if err == nil {
			modbus06Bin.RegSetterL8 = l8
		}

		dataH, err := byteArray.ReadByte()
		if err == nil {
			modbus06Bin.DataH8 = dataH
		}

		dataL, err := byteArray.ReadByte()
		if err == nil {
			modbus06Bin.DataL8 = dataL
		}
		c.PushProtoBin(modbus06Bin)
		//fmt.Println("服务器 接受 功能顺序-----modbus_rtu_06_decoder_resp-------------receive fncode : ",modbus06Bin.FnCode)
		//fmt.Println("after:",c.GetToolTotalCache().Available())
		retCode = errcode.TRNAS_DECODE_COMPLETE
	}

willExit:
	if byteArray != nil {
		pools.ByteArrayPool.Put(byteArray) //归还
	}
	//提前返回该回收的
	//allPackBytesPool.Put(allPackBytesRf)

	return retCode

}
