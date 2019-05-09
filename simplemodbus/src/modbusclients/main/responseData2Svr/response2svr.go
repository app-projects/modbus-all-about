package responseData2Svr

import (
	"net/source/userapi"
	"net/source/proto/constant/modbus"
	pack2 "testModbus/pack"
)

//模拟一个终端 反馈数据给 从机服务器  ，从机服务器是 一个虚拟化远端 控制台
func DispathModBusProtoBin(c userapi.IClient, pack userapi.IModBusProtoBinPack) {
	//pack.GetFnCode()
	//暂时 就是直接返回
	switch pack.GetFnCode() {
	case modbus.FunCode03:
		//封装一个数据 扔给服务器
		response03H := pack2.C2s03(pack.GetMac(), 10)
		c.Send(response03H)

	case modbus.FunCode06:



	case modbus.FunCode10:

	default:

	}

}
