package simulator

import (
	"testModbus/pack"
	"fmt"
	"net/source/proto/binfiles"
	"testModbus/echocontext"
	"testModbus/msgcmd"
	"time"
	"net/source/userapi"
	"net/source/msg/msgpusher"
	"net/source/proto/constant/modbus"
	"testModbus/connection"
	"testModbus/utils"
)

//01h,03h,40h,1dh,00h,01h,01h,0cch，
//系统层面 下发 查询 1个寄存器

func commitPack(c *connection.Connection, pck []byte, hdl func(args ...interface{}) interface{}) {
	var cmd echocontext.NetworkCmd
	cmd = msgcmd.NewNetworkCmd()
	cmd.SetRespdHandler(hdl)
	cmd.SetPack(pck)
	c.CommitReqPck(cmd)
}

func CommitGetSysInfoAsk(c *connection.Connection) {
	mac := utils.Int64_2Byte(c.GetId())
	var sysBytes = pack.EncodeQueryProto(mac, 0x03, 0x40, 0x07, 0x00, 0x06)
	for {

		commitPack(c, sysBytes, func(args ...interface{}) interface{} {
			fmt.Println("this is sys info hdl query---------start------------------")
			var b = args[0].(userapi.IModBusProtoBinPack)
			printQueryBin(b)
			fmt.Println("this is sys info hdl query----------end-----------------")
			msgpusher.Commit2Store(c, b.(*binfiles.Mod03_ProtoBinPackResp))

			return nil
		})
		//fmt.Println("out of CommitGetSysInfoAsk")
		time.Sleep(time.Millisecond *50)
	}
}

func CommitReadRegAsk(handler func(args ...interface{}) interface{}, mac byte, offsetH byte, offsetLow, dataCountH byte, dataCountL byte) int {
	c := connection.GetDevConnRouter().GetConnByMac(mac)
	if c == nil { //存在终端
		return -1
	}
	if !c.IsOpen() {
		return -2 //终端关闭状态
	}

	var readRegBytes = pack.EncodeQueryProto(mac, 0x03, offsetH, offsetLow, dataCountH, dataCountL)
	commitPack(c, readRegBytes, handler)

	return 0
}

func CommitGetAppInfoAsk(c *connection.Connection) {
	var appBytes []byte
	mac := utils.Int64_2Byte(c.GetId())
	appBytes = pack.EncodeQueryProto(mac, 0x03, 0x20, 0x00, 0x00, 0x06)
	for {
		commitPack(c, appBytes, func(args ...interface{}) interface{} {
			fmt.Println("this is appHdl query---------start------------------")
			var b = args[0].(userapi.IModBusProtoBinPack)
			printQueryBin(b)
			fmt.Println("this is appHdl query----------end-----------------")
			msgpusher.Commit2Store(c, b.(*binfiles.Mod03_ProtoBinPackResp))

			return nil
		})
		time.Sleep(time.Millisecond * 10)
		fmt.Println("commit ask tick .......................")
	}
}

func CommitModifyAppInfoByMac(mac byte, offsetH byte, offsetL, dataH byte, dataL byte) int {
	c := connection.GetDevConnRouter().GetConnByMac(mac)
	if c == nil { //存在终端
		return -1
	}
	if !c.IsOpen() {
		return -2 //终端关闭状态
	}

	var appModifyBytes []byte
	appModifyBytes = pack.EncodeModifyOneProto(mac, modbus.FunCode06, offsetH, offsetL, dataH, dataL)

	commitPack(c, appModifyBytes, func(args ...interface{}) interface{} {
		fmt.Println("this is modifyHdl---------start------------------")
		var b = args[0].(userapi.IModBusProtoBinPack)
		printModifyBin(b)
		fmt.Println("this modifyHdl--------end------------------")
		//推送restful api 消息到数据仓库
		msgpusher.Commit2Store06(c, b.(*binfiles.Mod06_ProtoBinPack))
		return nil
	})
	return 0
}

func CommitModifyAppInfo(c *connection.Connection) {
	var appModifyBytes []byte
	mac := utils.Int64_2Byte(c.GetId())
	appModifyBytes = pack.EncodeModifyOneProto(mac, 0x06, 0x20, 0x10, 0x05, 0x06)
	for {
		commitPack(c, appModifyBytes, func(args ...interface{}) interface{} {
			fmt.Println("this is modifyHdl---------start------------------")
			var b = args[0].(userapi.IModBusProtoBinPack)
			printModifyBin(b)
			fmt.Println("this modifyHdl--------end------------------")
			//推送restful api 消息到数据仓库
			msgpusher.Commit2Store06(c, b.(*binfiles.Mod06_ProtoBinPack))
			return nil
		})

		time.Sleep(time.Millisecond *50)
		//fmt.Println("out of CommitGetAppInfoAsk")
	}
}

func printQueryBin(bin userapi.IModBusProtoBinPack) {
	//数据存档
	modbus03BinResp := bin.(*binfiles.Mod03_ProtoBinPackResp)
	fmt.Println("服务器收到终端查询响应数据如下")
	fmt.Println("mac地址：", modbus03BinResp.Mac)
	fmt.Println("响应功能号：", modbus03BinResp.FnCode)
	fmt.Println("返回字节数：", modbus03BinResp.DataFieldLength)
	var j byte = 0
	var regNum = modbus03BinResp.DataFieldLength / 2
	for ; j < regNum; j++ {
		fmt.Printf("    字段%d：0x%x,0x%x\n", j, modbus03BinResp.DataFields[2*j], modbus03BinResp.DataFields[2*j+1])
	}
}

func printModifyBin(b userapi.IModBusProtoBinPack) {
	bin := b.(*binfiles.Mod06_ProtoBinPack)
	fmt.Println("服务器端设置成功数据原路返回：格式如下")
	fmt.Println("mac地址：", bin.Mac)
	fmt.Println("查询功能号：", bin.FnCode)

	fmt.Printf("下置的寄存器地址高8位：0x%x\n", bin.RegSetterH8)
	fmt.Printf("下置的寄存器地址低8位：0x%x\n", bin.RegSetterL8)
	fmt.Printf("下置的数据高8位：0x%x\n", bin.DataH8)
	fmt.Printf("下置的数据低8位：0x%x\n", bin.DataL8)
	fmt.Println("**********************************************")

}
