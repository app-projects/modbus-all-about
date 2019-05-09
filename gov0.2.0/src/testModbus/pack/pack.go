package pack

import (
	"net/source/utils/bytes"
	"testModbus/utils"
	"fmt"
	"testModbus/data"
)

/****
[设备地址]+[命令号03H] +       <指令头>    16bit  2byte
[起始寄存器地址高8位] +[低8位] +           16bit    2byte   //offset
[读取的寄存器数高8位] +[低8位] +           16bit    2byte    //length
[CRC校验的低8位] + [CRC校验的高8位]    16bit CRC   2byte

**/

type proto_req struct {
	Mac     byte
	FunCode byte
}

type proto3_req struct {
	proto_req
	OffsetRegH byte
	OffsetRegL byte
	ReadRegH   byte
	ReadRegL   byte
	Crc16      int16
}

/*
*
[设备地址] + [命令号06H] +       <指令头>
[需下置的寄存器地址高8位] + [低8位] +
[下置的数据高8位] + [低8位] +
[CRC校验的低8位] + [CRC校验的高8位]

*/

func EncodeModifyOneProto(mac byte, fnCode byte, offsetH byte, offsetL, dataH byte, dataL byte) []byte {
	var byteArr = bytes.NewByteArray([]byte{})
	byteArr.WriteByte(mac)
	byteArr.WriteByte(fnCode)

	byteArr.WriteByte(offsetH)
	byteArr.WriteByte(offsetL)

	byteArr.WriteByte(dataH)
	byteArr.WriteByte(dataL)

	//得出crc16
	upBytes := byteArr.Bytes()
	crc16 := utils.GetCrc16Tool()
	crc16.PushBytes(upBytes)
	h, l := crc16.Get()
	byteArr.WriteByte(l)
	byteArr.WriteByte(h)
	utils.ReleaseCrc16Tool(crc16)
	return byteArr.Bytes()
}

func EncodeQueryProto(mac byte, fnCode byte, offsetH byte, offsetLow, dataCountH byte, dataCountL byte) []byte {
	var byteArr = bytes.NewByteArray([]byte{})
	byteArr.WriteByte(mac)
	byteArr.WriteByte(fnCode)

	byteArr.WriteByte(offsetH)
	byteArr.WriteByte(offsetLow)

	byteArr.WriteByte(dataCountH)
	byteArr.WriteByte(dataCountL)

	//得出crc16
	upBytes := byteArr.Bytes()
	crc16 := utils.GetCrc16Tool()
	crc16.PushBytes(upBytes)
	h, l := crc16.Get()
	byteArr.WriteByte(l)
	byteArr.WriteByte(h)
	utils.ReleaseCrc16Tool(crc16)
	return byteArr.Bytes()
}

/*
*设备响应：
[设备地址] +[命令号03H] + <指令头>
[返回的字节个数] +
[数据1] +
[数据2] +
...+ [数据n] +
[CRC校验的低8位] + [CRC校验的高8位]
*
*/

//readRegNum 是寄存器的个数
//-------------------------------------------终端数据------------------------------------------------

const (
	TMPL_ADDR_H = 0x20
	TMPL_ADDR_L = 0x07

	EL_ADDR_H = 0x20
	EL_ADDR_L = 0x08
)

//-------------------------------------------终端数据------------------------------------------------

func C2s03App(mac byte, readRegNum byte, oh byte, ol byte) []byte {
	devData := data.GetDevDataContext().GetDevData(mac)
	dataLen := readRegNum * 2
	infoByteArray := bytes.NewByteArray([]byte{})
	infoByteArray.WriteByte(mac)
	infoByteArray.WriteByte(0x03)
	infoByteArray.WriteByte(dataLen)

	if readRegNum <= 0 {
		goto outflag
	} else {
		var offsetKey = utils.Bytes2Uint16(oh, ol)
		var i uint16 = 0
		var count uint16 = uint16(readRegNum)

		for i = 0; i < count; i++ {
			dh, dl := devData.GetAppDataByKey16(offsetKey + i)
			infoByteArray.WriteByte(dh)
			infoByteArray.WriteByte(dl)
		}
	}

outflag:

	crcTool := utils.GetCrc16Tool()
	crcTool.PushBytes(infoByteArray.Bytes())
	h, l := crcTool.Get()
	infoByteArray.WriteByte(l)
	infoByteArray.WriteByte(h)
	utils.ReleaseCrc16Tool(crcTool)
	return infoByteArray.Bytes()
}

func ModifyAppReg(mac byte, fnCode byte, offsetH byte, offsetLow, dataH byte, dataL byte) []byte {
	var devData = data.GetDevDataContext().GetDevData(mac)
	devData.PutAppData(offsetH, offsetLow, dataH, dataL)
	fmt.Printf("终端寄存器=%x %x,被修改:0x%x,%x,\n", offsetH, offsetLow, dataH, dataL)
	return EncodeModifyOneProto(mac, fnCode, offsetH, offsetLow, dataH, dataL)
}
