package data

import (
	"sync"
	"testModbus/utils"
)

type DataBlock struct {
	BaseAddr  uint16
	Uint16Map sync.Map
}

type DevData struct {
	Mac      byte
	AppBlock *DataBlock
	SysBlock *DataBlock
}

func newDevData(mac byte, appOffset uint16, sysOffset uint16) *DevData {
	d := DevData{}
	d.Mac = mac
	d.AppBlock = &DataBlock{BaseAddr: appOffset}
	d.SysBlock = &DataBlock{BaseAddr: sysOffset}

	return &d
}

func (this *DevData) PutAppData(offsetH byte, offsetL byte, dH byte, dL byte) {
	var key = utils.Bytes2Uint16(offsetH, offsetL)
	var value = utils.Bytes2Uint16(dH, dL)
	this.AppBlock.Uint16Map.Store(key, value)
}

func (this *DevData) PutAppDataByKey16(key uint16, dH byte, dL byte) {
	var value = utils.Bytes2Uint16(dH, dL)
	this.AppBlock.Uint16Map.Store(key, value)
}


func (this *DevData) PutSysData(offsetH byte, offsetL byte, dH byte, dL byte) {
	var key = utils.Bytes2Uint16(offsetH, offsetL)
	var value = utils.Bytes2Uint16(dH, dL)
	this.SysBlock.Uint16Map.Store(key, value)
}


func (this *DevData) GetAppDataByKey16(key uint16) (dH byte, dL byte) {
	v, ok := this.AppBlock.Uint16Map.Load(key)
	if ok {
		value, okk := v.(uint16)
		if okk {
			return utils.Uint162Byte(value)
		}
	}
	return 0, 0
}



func (this *DevData) GetAppData(offsetH byte, offsetL byte) (dH byte, dL byte) {
	var key = utils.Bytes2Uint16(offsetH, offsetL)
	v, ok := this.AppBlock.Uint16Map.Load(key)
	if ok {
		value, okk := v.(uint16)
		if okk {
			return utils.Uint162Byte(value)
		}
	}
	return 0, 0
}




func (this *DevData) GetSysData(offsetH byte, offsetL byte) (dH byte, dL byte) {
	var key = utils.Bytes2Uint16(offsetH, offsetL)
	v, ok := this.SysBlock.Uint16Map.Load(key)
	if ok {
		value, okk := v.(uint16)
		if okk {
			return utils.Uint162Byte(value)
		}
	}
	return 0, 0
}

func init() {
	devDataContext = newDevDataContext()
}

var devDataContext *DevDataContext

func GetDevDataContext() *DevDataContext {
	return devDataContext
}

type DevDataContext struct {
	devAllData sync.Map
}

func newDevDataContext() *DevDataContext {
	return &DevDataContext{}
}

const AppOffset = 0x2000
const SysOffset = 0x4000

func (this *DevDataContext) GetDevData(mac byte) *DevData {
	v, ok := this.devAllData.Load(mac)
	if ok {
		return v.(*DevData)
	} else {
		devData := newDevData(mac, AppOffset, SysOffset)
		this.devAllData.Store(mac, devData)
		return devData
	}
	return nil
}
