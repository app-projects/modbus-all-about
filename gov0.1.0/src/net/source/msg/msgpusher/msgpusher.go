package msgpusher

import (
	"net/source/userapi"
	"net/source/proto/binfiles"
	"encoding/json"
	"fmt"
	"testModbus/utils"
	"time"
	"net/source/proto/outputconfig"
)

//int
func Commit2Store(c userapi.IClient, modbus03BinResp *binfiles.Mod03_ProtoBinPackResp) {


	var dataKv = make(map[string]interface{})
	dataKv["mac"] = modbus03BinResp.Mac
	dataKv["funcode"] = modbus03BinResp.FnCode
	dataKv["ip"] = c.GetConn().RemoteAddr().String()
	dataKv["templ"] = modbus03BinResp.DataFields[0]
	dataKv["datafields"] = modbus03BinResp.DataFields //[0:modbus03BinResp.DataFieldLength]

	var dataLen = modbus03BinResp.DataFieldLength
	dataKv["data_fields_len"] = dataLen

	dataFieldsMap := make(map[string]interface{})

	var i byte
	for i = 0; i < dataLen; i++ {
		dataFieldsMap[byte2String(i)] = modbus03BinResp.DataFields[i]
	}
	dataKv["data_fieldsmap"] = dataFieldsMap

	datafs, e := json.Marshal(modbus03BinResp.DataFields)
	fmt.Println(e, "----", string(datafs))

	dataKv["timestamp"] = utils.TransTime2MillSec(time.Now())
	fmt.Println("push push ....mills:", utils.TransTime2MillSec(time.Now()))
	go utils.HttpPost(outputconfig.RemoteStorePushQueryAddr, dataKv, nil)

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





func Commit2Store06(c userapi.IClient, bin *binfiles.Mod06_ProtoBinPack) {

	var out = getInnerBinPost(c, bin)
	go utils.HttpPost(outputconfig.RemoteStoreModifyAddr, out, nil)
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

