package trans

import (
	"net/source/proto/trans/interfaces"
	"plugin"
	"sync"
	"net/source/proto/constant/modbus"
	"net/source/proto/trans/decode_suit/mb_rtu_03h_decoder"
	"net/source/proto/trans/decode_suit/mb_rtu_06h_decoder"
	"net/source/proto/trans/decode_suit/mb_rtu_10h_decoder"
	"net/source/proto/trans/decode_suit_svr/mb_rtu_03h_decoder_resp"
	"net/source/proto/trans/decode_suit_svr/mb_rtu_06h_decoder_resp"
)

var createFnCache sync.Map // soname fun CreateIntance

func GetDecoderPluginByName(soName string) (interfaces.Decoder, error) {
	createFn, exist := createFnCache.Load(soName)
	if exist {
		return createFn.(func() interfaces.Decoder)(), nil
	}

	plg, err1 := plugin.Open(soName)
	if err1 == nil && plg != nil {
		fn, err1 := plg.Lookup("CreateIntance")
		if err1 == nil && fn != nil {
			ins := fn.(func() interfaces.Decoder)()
			//save
			createFnCache.Store(soName, fn)
			return ins, nil
		}
	}
	return nil, err1
}

func GetDecoderPluginByFnCode(funCode int32) (interfaces.Decoder, error) {
	return codeTestMap[byte(funCode)](), nil ///test

	soName, _ := decoderPluginsMap.Load(funCode)
	return GetDecoderPluginByName(soName.(string))
}

func GetDecoderPluginByFnCodeResp(funCode int32) (interfaces.Decoder, error) {
	return codeTestMapResp[byte(funCode)](), nil ///test

	soName, _ := decoderPluginsMap.Load(funCode)
	return GetDecoderPluginByName(soName.(string))
}



var decoderPluginsMap sync.Map

var codeTestMap map[byte]func() interfaces.Decoder
var codeTestMapResp map[byte]func() interfaces.Decoder
func init() {
	codeTestMap= make(map[byte]func() interfaces.Decoder,3)

	codeTestMap[modbus.FunCode03] = mb_rtu_03h_decoder.CreateIntance
	codeTestMap[modbus.FunCode06] = mb_rtu_06h_decoder.CreateIntance
	codeTestMap[modbus.FunCode10] = mb_rtu_10h_decoder.CreateIntance


	codeTestMapResp= make(map[byte]func() interfaces.Decoder,3)
	codeTestMapResp[modbus.FunCode03] = mb_rtu_03h_decoder_resp.CreateIntance
	codeTestMapResp[modbus.FunCode06] = mb_rtu_06h_decoder_resp.CreateIntance



	/*	PutDecoder(int32(modbus.FunCode03), "./decode_suit/mb_rtu_03h_decoder/mb_rtu_03h_decoder.so")
		PutDecoder(int32(modbus.FunCode06), "./decode_suit/mb_rtu_03h_decoder/mb_rtu_03h_decoder.so")
		PutDecoder(int32(modbus.FunCode10), "./decode_suit/mb_rtu_03h_decoder/mb_rtu_10h_decoder.so")*/
}

func PutDecoder(code int32, soPath string) {
	decoderPluginsMap.Store(code, soPath)
}
