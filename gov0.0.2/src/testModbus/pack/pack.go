package pack

import (
	"net/source/utils/bytes"
	"testModbus/utils"
	"net/source/proto/binfiles"
	"math/rand"
	"time"
)

/****
[设备地址]+[命令号03H] +       <指令头>    16bit  2byte
[起始寄存器地址高8位] +[低8位] +           16bit    2byte   //offset
[读取的寄存器数高8位] +[低8位] +           16bit    2byte    //length
[CRC校验的低8位] + [CRC校验的高8位]    16bit CRC   2byte

**/

type proto_req struct{
	Mac byte
	FunCode byte
}

type proto3_req struct{
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



func EncodeAskProto06H() []byte  {
	var proto06 binfiles.Mod06_ProtoBinPack
	proto06.Mac=0x01
	proto06.FnCode = 0x06
	proto06.RegSetterH8 = 0x14
	proto06.RegSetterL8=0x15
	proto06.DataH8 = 0x16
	proto06.DataL8 =0x17

	var byteArr =bytes.NewByteArray([]byte{})

	byteArr.WriteByte(proto06.Mac)
	byteArr.WriteByte(proto06.FnCode)

	byteArr.WriteByte(proto06.RegSetterH8)
	byteArr.WriteByte(proto06.RegSetterL8)

	byteArr.WriteByte(proto06.DataH8)
	byteArr.WriteByte(proto06.DataL8)

	//得出crc16
	upBytes:=byteArr.Bytes()
	crc16:=utils.GetCrc16Tool()
	crc16.PushBytes(upBytes)
	h,l:=crc16.Get()
	byteArr.WriteByte(l)
	byteArr.WriteByte(h)
	utils.ReleaseCrc16Tool(crc16)
	return byteArr.Bytes()
}

/*****

01
03
20
00
00
20
***/
func EncodeAskProto03H() []byte  {
	var proto03 proto3_req
	proto03.Mac=0x01
	proto03.FunCode = 0x03
	proto03.OffsetRegH = 0x20
	proto03.OffsetRegL=0x00
	proto03.ReadRegH = 0x00
	proto03.ReadRegL =0x20
	var byteArr =bytes.NewByteArray([]byte{})

	byteArr.WriteByte(proto03.Mac)
	byteArr.WriteByte(proto03.FunCode)

	byteArr.WriteByte(proto03.OffsetRegH)
	byteArr.WriteByte(proto03.OffsetRegL)

	byteArr.WriteByte(proto03.ReadRegH)
	byteArr.WriteByte(proto03.ReadRegL)

	//得出crc16
	upBytes:=byteArr.Bytes()
	crc16:=utils.GetCrc16Tool()
	crc16.PushBytes(upBytes)
	h,l:=crc16.Get()
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
func C2s03(mac byte,dataLen byte) []byte {
	infoByteArray:=bytes.NewByteArray([]byte{})
	infoByteArray.WriteByte(mac)
	infoByteArray.WriteByte(0x03)
	infoByteArray.WriteByte(dataLen)

	rand.Seed(time.Now().Unix())
	var i byte
	for i = 0; i<dataLen;i++ {
		infoByteArray.WriteByte(byte(rand.Intn(40)))
	}
	crcTool:=utils.GetCrc16Tool()

	crcTool.PushBytes(infoByteArray.Bytes())
	h,l :=crcTool.Get()
	infoByteArray.WriteByte(l)
	infoByteArray.WriteByte(h)
	utils.ReleaseCrc16Tool(crcTool)
	return infoByteArray.Bytes()
}

func DecodeProto03HBytes(packBytes []byte) (*proto3_req, int) {

	if len(packBytes) <=2{
		return nil, -1
	}

	var crc16Bytes =packBytes[len(packBytes)-2:]
	var crc16Tool =utils.GetCrc16Tool()
	crc16Tool.PushByte(crc16Bytes[1],crc16Bytes[0])
	crcValue:=crc16Tool.Value()

	crc16Tool.Reset()
	crc16Tool.PushBytes(packBytes[:len(packBytes)-2])

	realValue:= crc16Tool.Value()

	utils.ReleaseCrc16Tool(crc16Tool)

	if crcValue!=realValue{ //crc检查失败
		return nil, -2  //非法包
	}

	var byteArr =bytes.NewByteArray(packBytes)

	var proto03 proto3_req
	mac,err:=byteArr.ReadByte()
	if err==nil{
		proto03.Mac= mac
	}

	funCode,err:=byteArr.ReadByte()
	if err==nil{
		proto03.FunCode= funCode
	}


	offsetRegH,err:=byteArr.ReadByte()
	if err==nil{
		proto03.OffsetRegH= offsetRegH
	}


	offsetRegL,err:=byteArr.ReadByte()
	if err==nil{
		proto03.OffsetRegL= offsetRegL
	}



	readRegH,err:=byteArr.ReadByte()
	if err==nil{
		proto03.ReadRegH= readRegH
	}


	readRegL,err:=byteArr.ReadByte()
	if err==nil{
		proto03.ReadRegL= readRegL
	}

	return &proto03,0
}

