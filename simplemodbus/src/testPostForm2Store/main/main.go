package main

import (
	"fmt"
	"testModbus/utils"
	"time"
	"math/rand"
)

func testCommitQueuryLog()  {
	rand.Seed(time.Now().Unix())
	for {
		time.Sleep(time.Second * 5)
		fmt.Println("push push ....mills:", utils.TransTime2MillSec(time.Now()))
		var toUrl = "http://192.168.1.108:12345/dev/log"

		var dataKv = make(map[string]interface{})
		dataKv["mac"] = "2"
		dataKv["funcode"] = 6
		dataKv["ip"] = "192.168.1.104"
		dataKv["templ"] = rand.Intn(40)
		dataKv["timestamp"] = utils.TransTime2MillSec(time.Now())
		utils.HttpPost(toUrl, dataKv, nil)

	}
}



/*
	RegSetterH8 byte `json:"reg_setter_h_8"`
	RegSetterL8 byte `json:"reg_setter_l_8"`
	DataH8      byte `json:"data_h_8"`
	DataL8      byte `json:"data_l_8"`


	ClientId  int64  `json:"client_id"`
	Mac       string `json:"mac"`
	FnCode    byte   `json:"fn_code"`
	Timestamp int64  `json:"timestamp"`
	Ip        string `json:"ip"`


*/
func testCommitModifyLog() {
	for {
		time.Sleep(time.Second * 5)
		rand.Seed(time.Now().Unix())
		fmt.Println("push push ....mills:", utils.TransTime2MillSec(time.Now()))
		var toUrl= "http://192.168.1.108:12345/dev/log/modify"

		var dataKv= make(map[string]interface{})
		dataKv["reg_setter_h_8"] = 0x10
		dataKv["reg_setter_l_8"] = 0x11
		dataKv["data_h_8"] = 0x12
		dataKv["data_l_8"] = 0x17

		dataKv["client_id"] = 1
		dataKv["mac"] = "2"
		dataKv["fn_code"] = 6
		dataKv["timestamp"] = utils.TransTime2MillSec(time.Now())
		dataKv["ip"] = "192.168.1.104"

		utils.HttpPost(toUrl, dataKv, nil)
	}

}

func testCrc16mobbus()  {
	var b =[]byte{0x02,0x03, 0x14 ,0x15 ,0x16, 0x17}
	t:=utils.GetCrc16Tool()
	t.PushBytes(b)
	fmt.Printf("%x",t.Value())
	utils.ReleaseCrc16Tool(t)

}

func main() {
	testCrc16mobbus()
	//testCommitModifyLog()
	//testCommitQueuryLog()
}
