package mb_rtu_03h_decoder

import (
	"net/source/userapi"
	"net/source/proto/trans/interfaces"
	"sync"
	"net/source/proto/binfiles"
	"net/source/proto/constant/modbus"
	"net/source/utils/bytes"
	"net/source/proto/pools"
	"testModbus/utils"
	"fmt"
	"net/source/proto/trans/errcode"
)

type modbus_rtu_03_decoder struct {
}

var selfPool = sync.Pool{
	New: func() interface{} {
		return &modbus_rtu_03_decoder{}
	},
}

func CreateIntance() interfaces.Decoder {
	return selfPool.Get().(*modbus_rtu_03_decoder)
}

func ReleaseInstance(ins interfaces.Decoder) {
	selfPool.Put(ins)
}

/***
Modbus  协议分段 分区格式：

<指令头>  mac,code    int16
<content>
<crc16>      int16


发送命令：
[设备地址]+[命令号03H] +       <指令头>    16bit  2byte
[起始寄存器地址高8位] +[低8位] +           16bit    2byte
[读取的寄存器数高8位] +[低8位] +           16bit    2byte
[CRC校验的低8位] + [CRC校验的高8位]    16bit CRC   2byte

设备响应：
[设备地址] +[命令号03H] +       <指令头>

[返回的字节个数] +
[数据1] +
[数据2] +
...+ [数据n] +
 



03H协议请求特征是固定的：
1 / 总字节是， 8byte

*/

var pack_all_total_size = 8

var allpackBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, pack_all_total_size)
	},
}

func (this *modbus_rtu_03_decoder) Decode(c userapi.IClient) int {

	if c.GetToolTotalCache().Available() < pack_all_total_size {
		return errcode.ERR_TRNAS_DECODE_TRY_AGAIN
	}

	var retCode = errcode.ERR_TRNAS_DECODE_SKIP

	//copy data 通过上面获得 ver protolen 把数据 幺出来

	var allPackBytes = allpackBytesPool.Get().([]byte)

	var byteArray *bytes.ByteArray
	_, e := c.GetToolTotalCache().Read(allPackBytes)
	if e == nil {

		byteArray = pools.ByteArrayPool.Get().(*bytes.ByteArray) //借出
		byteArray.Reset()
		byteArray.WriteBytes(allPackBytes)

		var ok = utils.CheckContentCRC16OK(byteArray.Bytes())

		if !ok {
			goto willExit
		}

		// 取值
		//cmd head
		mac, err := byteArray.ReadByte()
		byteArray.ReadByte() //fnCode

		var modbusBin = binfiles.CreateModBusProtoBin(mac, modbus.FunCode03, c.GetId())
		modbus03Bin := modbusBin.(*binfiles.Mod03_ProtoBinPack)
		sh8, err := byteArray.ReadByte()
		if err == nil {
			modbus03Bin.StartRegH8 = sh8
		}

		sl8, err := byteArray.ReadByte()
		if err == nil {
			modbus03Bin.StartRegL8 = sl8
		}

		rh, err := byteArray.ReadByte()
		if err == nil {
			modbus03Bin.ReadRegH = rh
		}

		rl, err := byteArray.ReadByte()
		if err == nil {
			modbus03Bin.ReadRegL = rl
		}

		c.GetClientProtoBinChan()<-modbusBin
		allpackBytesPool.Put(allPackBytes)

		fmt.Println("终端收到服务器的查询指令：格式如下")
		fmt.Println("mac地址：", modbus03Bin.Mac)
		fmt.Println("查询功能号：", modbus03Bin.FnCode)
		/***
		[设备地址]
		[命令号03H] +       <指令头>    16bit  2byte
		[起始寄存器地址高8位] +[低8位] +           16bit    2byte   //offset
		[读取的寄存器数高8位] +[低8位] +           16bit    2byte    //length
		[CRC校验的低8位] + [CRC校验的高8位]    16bit CRC   2byte
		**/
		fmt.Printf("起始寄存器地址高8位：0x%x\n", modbus03Bin.StartRegH8)
		fmt.Printf("起始寄存器地址低8位：0x%x\n", modbus03Bin.StartRegL8)
		fmt.Printf("读取的寄存器数高8位：0x%x\n", modbus03Bin.ReadRegH)
		fmt.Printf("读取的寄存器数低8位：0x%x\n", modbus03Bin.ReadRegL)
		fmt.Println("**********************************************")

		retCode = errcode.TRNAS_DECODE_COMPLETE
	}

willExit:

	if byteArray != nil {
		pools.ByteArrayPool.Put(byteArray) //归还
	}

	allpackBytesPool.Put(allPackBytes)
	return retCode
}
