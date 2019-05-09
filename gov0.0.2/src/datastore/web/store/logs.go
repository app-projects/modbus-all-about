package store

import (
	"sync"
	"sync/atomic"
	"testModbus/utils"
)

type Devlog interface{
	GetKey() int64

	SetMac(mac string)
	SetIp(ip string)
	SetValue(v interface{})
	SetTimeStmap(t int64)

	GetMac()string
	GetIp()string
	GetValue()interface{}
	GetTimeStmap()int64
}


const maxSize = 20000

//毫秒级的log

type devLogHistory struct {
	millSecMap *utils.MapList
	Size       int64
}

func newDevLogHistory() *devLogHistory {
	ins := devLogHistory{}
	ins.millSecMap = utils.NewMapList()
	return &ins
}

func (this *devLogHistory) put(log Devlog) {
	atomic.AddInt64(&this.Size, 1)
	this.millSecMap.Push(log)
	if this.Size > maxSize*2/3 {
		// 立刻做
		this.millSecMap.Shift(maxSize / 3)
	}
}

func (this *devLogHistory) LoadLog(recevicerFun func(log Devlog) int, fromTimeStamp int64, toTimestamp int64) {
	if recevicerFun == nil || this.Size <= 0 {
		return
	}
	this.millSecMap.Walk(func(data utils.Keyer) int {
		log := data.(Devlog)
		if log != nil && (log.GetTimeStmap() >= fromTimeStamp && log.GetTimeStmap() <= toTimestamp) {
			recevicerFun(log)
		}
		return utils.STATUS_LOOP_OK
	})
}

func (this *devLogHistory) LoadLogRecent(recevicerFun func(log Devlog) int, offsetTimestamp int64, recentBeforeNum int) {
	if recevicerFun == nil || this.Size <= 0 {
		return
	}

	this.millSecMap.WalkFromEnd(func(data utils.Keyer) int {
		log := data.(Devlog)
		if log != nil && (log.GetTimeStmap() <= offsetTimestamp) {
			recevicerFun(log)
			recentBeforeNum--
			if recentBeforeNum <= 0 {
				return utils.STATUS_LOOP_EXIT
			}
		}
		return utils.STATUS_LOOP_OK
	})
}

//lowWaterBaseTime 低水位基线  小于0 lowWaterBaseTime < 0 ignore
func (this *devLogHistory) LoadLogRecentByLowWaterTimeLine(recevicerFun func(log Devlog) int, offsetTimestamp int64, recentBeforeNum int,lowWaterBaseTime int64) {
	if recevicerFun == nil || this.Size <= 0 {
		return
	}
	var lowIgnore = false
	if lowWaterBaseTime<0{ //ignore this condition
		lowIgnore = true
	}

	var offsetTimeIgnore = false
	if offsetTimestamp<0{ //ignore this condition
		offsetTimeIgnore = true
	}



	this.millSecMap.WalkFromEnd(func(data utils.Keyer) int {
		log := data.(Devlog)
		if !lowIgnore{ //不能忽略 则从新赋予条件
			lowIgnore = log.GetTimeStmap()>=lowWaterBaseTime
		}
		if !offsetTimeIgnore{ //不能忽略 则从新赋予条件
			offsetTimeIgnore = log.GetTimeStmap() <= offsetTimestamp
		}

		if log != nil && (offsetTimeIgnore &&lowIgnore) {
			recevicerFun(log)
			recentBeforeNum--
			if recentBeforeNum <= 0 {
				return utils.STATUS_LOOP_EXIT
			}
		}
		return utils.STATUS_LOOP_OK
	})
}

type logStore struct {
	devLogMap sync.Map
}

func (store *logStore) LoadLogByTimeStamp(mac string, recevicerFun func(log Devlog) int, fromTimeStamp int64, toTimestamp int64) {
	devHistory, ok := store.devLogMap.Load(mac)
	if ok {
		devHistory.(*devLogHistory).LoadLog(recevicerFun, fromTimeStamp, toTimestamp)
	}
}

func (store *logStore) LoadLogByRecent(mac string, recevicerFun func(log Devlog) int, offsetTimestamp int64, recentBeforeNum int) {
	devHistory, ok := store.devLogMap.Load(mac)
	if ok {
		devHistory.(*devLogHistory).LoadLogRecent(recevicerFun, offsetTimestamp, recentBeforeNum)
	}
}

func (store *logStore) LoadLogRecentByLowWaterTimeLine(mac string,recevicerFun func(log Devlog) int, offsetTimestamp int64, recentBeforeNum int,lowWaterBaseTime int64){
	devHistory, ok := store.devLogMap.Load(mac)
	if ok {
		devHistory.(*devLogHistory).LoadLogRecentByLowWaterTimeLine(recevicerFun, offsetTimestamp, recentBeforeNum,lowWaterBaseTime)
	}
}


func (store *logStore) PushLog(devlog Devlog) {
	devLogH, ok := store.devLogMap.Load(devlog.GetMac())
	if !ok {
		devLogH = newDevLogHistory()
		store.devLogMap.Store(devlog.GetMac(), devLogH)
	}
	devLogH.(*devLogHistory).put(devlog)
}

func (store *logStore) QueryExist(mac string) bool {
	_, ok := store.devLogMap.Load(mac)
	return ok
}

//查询记录仓库
var store logStore

type LogStore interface {
	PushLog(devlog Devlog)
	LoadLogByTimeStamp(mac string, recevicerFun func(log Devlog) int, fromTimeStamp int64, toTimestamp int64)
	QueryExist(mac string) bool
	LoadLogByRecent(mac string, recevicerFun func(log Devlog) int, offsetTimestamp int64, recentBeforeNum int)
	//lowWaterBaseTime if <0 will ignore this condition
	LoadLogRecentByLowWaterTimeLine(mac string,recevicerFun func(log Devlog) int, offsetTimestamp int64, recentBeforeNum int,lowWaterBaseTime int64)
}

func GetStoreInstance() LogStore {
	return &store
}
//修改记录仓库
var storeModify logStore
func GetModifyStoreInstance() LogStore {
	return &storeModify
}

func StartLogStore() {

}
