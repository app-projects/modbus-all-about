package responseData2Svr

import (
	"net/source/userapi"
	"net/source/proto/constant/modbus"
	pack2 "testModbus/pack"
	"net/source/proto/binfiles"
	"testModbus/connection"
	"fmt"
)

type ClientBinHdl = func(c userapi.IClient, pack userapi.IModBusProtoBinPack)

var binHadlMap map[byte]ClientBinHdl

func init() {
	binHadlMap = make(map[byte]ClientBinHdl, 3)
	binHadlMap[modbus.FunCode03] = handl03Bin
	binHadlMap[modbus.FunCode06] = handl06Bin
	//binHadlMap[modbus.FunCode10] = handl10Bin
}

func ResponseBin(c *connection.Connection) {
	for {
		protoBin, err := c.PopModBusProtoBin()
		if (err == nil && protoBin != nil) {
			fn := binHadlMap[protoBin.GetFnCode()]
			fmt.Println("终端收到 功能请求 响应顺序——————————————————————————————：", protoBin.GetFnCode())
			if fn != nil {
				fn(c, protoBin)
			} else {
				fmt.Printf("提示:终端 数据包protoBin没有响应 处理函数:funcode:%d , clientid:%d , mac:%d \n", protoBin.GetFnCode(), protoBin.GetClientId(), protoBin.GetMac())
			}
		}
	}
}

func handl03Bin(c userapi.IClient, pack userapi.IModBusProtoBinPack) {
	modbus03Bin := pack.(*binfiles.Mod03_ProtoBinPack)
	var registerNum byte =modbus03Bin.ReadRegL
	response03H := pack2.C2s03App(pack.GetMac(), registerNum,modbus03Bin.StartRegH8,modbus03Bin.StartRegL8) //6个寄存器 暂定
	c.Send(response03H)
}

func handl06Bin(c userapi.IClient, pack userapi.IModBusProtoBinPack) {
	modbus06Bin := pack.(*binfiles.Mod06_ProtoBinPack)
	switch modbus06Bin.RegSetterH8 {
	case 0x20: //应用的空间
		response := pack2.ModifyAppReg(modbus06Bin.Mac, modbus06Bin.GetFnCode(), modbus06Bin.RegSetterH8, modbus06Bin.RegSetterL8, modbus06Bin.DataH8, modbus06Bin.DataL8)
		c.Send(response)
	default:
		response := pack2.ModifyAppReg(modbus06Bin.Mac, modbus06Bin.GetFnCode(), modbus06Bin.RegSetterH8, modbus06Bin.RegSetterL8, modbus06Bin.DataH8, modbus06Bin.DataL8)
		c.Send(response)
	}

}

func handl10Bin(c userapi.IClient, pack userapi.IModBusProtoBinPack) {
}
