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
	"fmt"
	"net/source/proto/trans/errcode"
	"net/source/proto/outputconfig"
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

var allpackBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, pack_all_total_size)
	},
}

func (this *modbus_rtu_06_decoder_resp) Decode(c userapi.IClient) int {
	if c.GetToolTotalCache().Available() < pack_all_total_size {
		return errcode.ERR_TRNAS_DECODE_TRY_AGAIN
	}

	var retCode = errcode.ERR_TRNAS_DECODE_TRY_AGAIN
	//copy data 通过上面获得 ver protolen 把数据 幺出来

	var allPackBytes = allpackBytesPool.Get().([]byte)
	//allPackBytes := allPackBytesRf[0:pack_all_total_size]

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
		mac, err := byteArray.ReadByte()
		byteArray.ReadByte() //fnCode

		var modbusBin = binfiles.CreateModBusProtoBin(mac, modbus.FunCode06, c.GetId())
		modbus06Bin := modbusBin.(*binfiles.Mod06_ProtoBinPack)

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

		//数据copy 成功
		//responseData2Svr.DispathModBusProtoBin(c, modbusBin)
		//提交到远程仓库
		commit2Store(c, modbus06Bin)
		retCode = errcode.TRNAS_DECODE_COMPLETE
	}

willExit:
	if byteArray != nil {
		pools.ByteArrayPool.Put(byteArray) //归还
	}
	//提前返回该回收的
	allpackBytesPool.Put(allPackBytes)

	return retCode

}

func commit2Store(c userapi.IClient, bin *binfiles.Mod06_ProtoBinPack) {

	fmt.Println("服务器端设置成功数据原路返回：格式如下")
	fmt.Println("mac地址：", bin.Mac)
	fmt.Println("查询功能号：", bin.FnCode)
	/***
	[设备地址]
	[命令号03H] +       <指令头>    16bit  2byte
	[需下置的寄存器地址高8位] + [低8位] +      2byte
	[下置的数据高8位] + [低8位] +             2byte
	[CRC校验的低8位] + [CRC校验的高8位]       2byte
	**/
	fmt.Printf("下置的寄存器地址高8位：0x%x\n", bin.RegSetterH8)
	fmt.Printf("下置的寄存器地址低8位：0x%x\n", bin.RegSetterL8)
	fmt.Printf("下置的数据高8位：0x%x\n", bin.DataH8)
	fmt.Printf("下置的数据低8位：0x%x\n", bin.DataL8)
	fmt.Println("**********************************************")
	var out = getInnerBinPost(c, bin)
	/*	v, e := json.Marshal(out)
		if e == nil {
			fmt.Println("marshell:", string(v))
		}*/
	utils.HttpPost(outputconfig.RemoteStoreModifyAddr, out, nil)

}

func getInnerBinPost(c userapi.IClient, src *binfiles.Mod06_ProtoBinPack) map[string]interface{} {
	var dataKv = make(map[string]interface{})
	dataKv["regsetterh8"] = byte2String(src.RegSetterH8)
	dataKv["regsetterl8"] = byte2String(src.RegSetterL8)
	dataKv["datah8"] = byte2String(src.DataH8)
	dataKv["datal8"] = byte2String(src.DataL8)
	dataKv["clientid"] = int642String(c.GetId())
	dataKv["mac"] = byte2String(src.Mac)
	dataKv["fncode"] = byte2String(src.FnCode)
	dataKv["ip"] = c.GetConn().RemoteAddr().String()
	dataKv["timestamp"] = utils.TransTime2MillSec(time.Now())

	return dataKv
}

func byte2String(b byte) string {
	var s = fmt.Sprintf("%d", b)
	fmt.Println(s)
	return s
}

func int642String(b int64) string {
	var s = fmt.Sprintf("%d", b)
	fmt.Println(s)
	return s
}
