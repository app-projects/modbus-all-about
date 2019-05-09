package mb_rtu_03h_decoder_resp

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
	"time"
	"net/source/proto/outputconfig"
)

type modbus_rtu_03_decoder_resp struct {
}

var selfPool = sync.Pool{
	New: func() interface{} {
		return &modbus_rtu_03_decoder_resp{}
	},
}

func CreateIntance() interfaces.Decoder {
	return selfPool.Get().(*modbus_rtu_03_decoder_resp)
}

func ReleaseInstance(ins interfaces.Decoder) {
	selfPool.Put(ins)
}

/***
Modbus  协议分段 分区格式：

<指令头>  mac,code    int16
<content>
<crc16>      int16

设备响应：
[设备地址] +[命令号03H] +       <指令头>

[返回的字节个数] +
[数据1] +
[数据2] +
...+ [数据n] +
 



03H协议请求特征是固定的：
1 / 总字节是， 8byte

*/

var pack_content_total_size = 1024 //06H包体  除了指令头 2byte

var allPackBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, pack_content_total_size)
	},
}

func (this *modbus_rtu_03_decoder_resp) Decode(c userapi.IClient) int {

	//copy data 通过上面获得 ver protolen 把数据 幺出来
	var availableLen = c.GetToolTotalCache().Available()
	if availableLen < 3 {
		return errcode.ERR_TRNAS_DECODE_TRY_AGAIN //整体的信息，但是没有整体的内容
	}

	var retCode = errcode.ERR_TRNAS_DECODE_TRY_AGAIN

	//预测处理
	availidtmp := c.GetToolTotalCache().BytesAvailable()
	var dataFieldLen byte
	dataFieldLen = availidtmp[2]
	var allPackLength = int(dataFieldLen) + 2 + 2 + 1 //headcmd 2 + crc2+dataField-length 1

	canRead := (availableLen >= allPackLength) //2 是crc所需要的占位
	if !canRead {
		return errcode.ERR_TRNAS_DECODE_TRY_AGAIN //整体的信息，但是没有整体的内容
	}

	var allPackBytesRf = allPackBytesPool.Get().([]byte)
	allPackBytes := allPackBytesRf[:allPackLength]

	_, e := c.GetToolTotalCache().Read(allPackBytes)
	var byteArray *bytes.ByteArray
	if e == nil {
		byteArray = pools.ByteArrayPool.Get().(*bytes.ByteArray) //借出
		byteArray.Reset()
		byteArray.WriteBytes(allPackBytes)

		var ok = utils.CheckContentCRC16OK(byteArray.Bytes())

		if !ok {

			goto ExitFlag
			return -1
		}

		mac, _ := byteArray.ReadByte()
		byteArray.ReadByte() //fnCode
		var modbusBin = binfiles.CreateModBusProtoBinResp(mac, modbus.FunCode03, c.GetId())
		modbusBin.SetMac(mac)

		// 取值

		var dataFieldLenth, _ = byteArray.ReadByte()

		modbus03BinResp := modbusBin.(*binfiles.Mod03_ProtoBinPackResp)
		modbus03BinResp.Mac = mac
		modbus03BinResp.DataFieldLength = dataFieldLenth

		byteArray.ReadBytes(modbus03BinResp.DataFields, int(dataFieldLenth), 0)
		//数据存档
		//c.DispathModBusProtoBin(modbusBin)

		fmt.Println("服务器收到终端查询响应数据如下")
		fmt.Println("mac地址：", modbus03BinResp.Mac)
		fmt.Println("响应功能号：", modbus03BinResp.FnCode)
		fmt.Println("返回字节数：", modbus03BinResp.DataFieldLength)
		/***
		设备响应：
		[设备地址] +[命令号03H] +       <指令头>

		[返回的字节个数] +
		[数据1] +
		[数据2] +
		...+ [数据n] +
		**/
		var j byte = 0
		for ; j < dataFieldLenth; j++ {
			fmt.Printf("    字段%d：%d\n", j, modbus03BinResp.DataFields[j])
		}
		//推送restful api 消息到数据仓库
		commit2Store(c, modbus03BinResp)
		fmt.Println("**********************************************")

		retCode = errcode.TRNAS_DECODE_COMPLETE
	}
ExitFlag:
	if byteArray != nil {
		pools.ByteArrayPool.Put(byteArray) //归还
	}
	allPackBytesPool.Put(allPackBytesRf)

	return retCode

}
//int
func commit2Store(c userapi.IClient, modbus03BinResp *binfiles.Mod03_ProtoBinPackResp) {
	var dataKv = make(map[string]interface{})
	dataKv["mac"] =modbus03BinResp.Mac
	dataKv["funcode"] = modbus03BinResp.FnCode
	dataKv["ip"] = c.GetConn().RemoteAddr().String()
	dataKv["templ"] = modbus03BinResp.DataFields[0]



	dataKv["timestamp"] = utils.TransTime2MillSec(time.Now())
	fmt.Println("push push ....mills:", utils.TransTime2MillSec(time.Now()))
	utils.HttpPost( outputconfig.RemoteStorePushQueryAddr, dataKv, nil)

}
