package msgproc

import (
	"sync"
	"net/source/proto"
	"fmt"
)

var instanceMsgProcRepertory *msgProcRepertory

type msgProcRepertory struct {
	//msg parse
	msgProcMap sync.Map

	//callbacks
	allCallbackMap sync.Map
}

type MsgCaller struct {
	Id          string
	uuid        uint32
	MsgTypeCode int32
	Handler     func(msg interface{}, tools *UserApi)
}

func CreateMsgCaller(msgTypeCode int32, uniqueid string, hdl func(msg interface{}, tools *UserApi)) *MsgCaller {
	caller := MsgCaller{}
	caller.Id = uniqueid
	caller.Handler = hdl
	caller.MsgTypeCode = msgTypeCode
	uuidForMsgCaller(&caller)
	return &caller
}

func uuidForMsgCaller(caller *MsgCaller) {
	uuidstr := fmt.Sprintf("caller.Id= %s+caller.MsgTypeCode=%d", caller.Id, caller.MsgTypeCode)
	caller.uuid = proto.CRC32GenId(uuidstr)
}

func init() {
	instanceMsgProcRepertory = &msgProcRepertory{}
}

//***************************Call back section******************************************

type IDispatcherMsg interface {
	CommitMsg(msg interface{}) error
}

type OneCodeCallbackMap struct {
	callbacksDict sync.Map //uuid ----> MsgCaller
}

func newOneCodeCallbackMap() *OneCodeCallbackMap {
	return &OneCodeCallbackMap{}
}

//commit msg and will publish  to all listener and for each msg callbacker
func (this *OneCodeCallbackMap) CommitMsg(msg interface{}) error {
	dict := this.callbacksDict

	var tools = GetAppTools()
	//dispatch  这里预留 处理 ，作为消息分发接口
	dict.Range(func(key, value interface{}) bool {
		caller := value.(*MsgCaller)
		caller.Handler(msg, tools)
		return true
	})
	return nil
}

func (this *OneCodeCallbackMap) add(caller *MsgCaller) {
	v, ok := this.callbacksDict.Load(caller.uuid)
	if !ok {
		this.callbacksDict.Store(caller.uuid, caller)
	} else {
		fmt.Errorf("OneCodeCallbackMap.push.caller has exist one type=%d, id=%s\n detail=%T\n", caller.MsgTypeCode, caller.Id, v)
		return
	}
}

func (this *OneCodeCallbackMap) remove(callerUuid uint32) {
	this.callbacksDict.Delete(callerUuid)
}

func (this *OneCodeCallbackMap) get(callerUuid uint32) *MsgCaller {
	v, ok := this.callbacksDict.Load(callerUuid)
	if ok {
		return v.(*MsgCaller)
	}
	return nil
}

func (p *msgProcRepertory) addMsgHandler(msgTypeCode int32, callback *MsgCaller) {
	codeCallbackMap, ok := p.allCallbackMap.Load(msgTypeCode)
	if ok {
		oldOneCodeMap := codeCallbackMap.(*OneCodeCallbackMap)
		oldOneCodeMap.add(callback)

	} else {
		newOneMap := newOneCodeCallbackMap()
		newOneMap.add(callback)
		p.allCallbackMap.Store(msgTypeCode, newOneMap)
	}
}

func (p *msgProcRepertory) getMsgHanlders(msgTypeCode int32) *OneCodeCallbackMap {
	codeCallbackMap, ok := p.allCallbackMap.Load(msgTypeCode)
	if ok {
		oldOneCodeMap := codeCallbackMap.(*OneCodeCallbackMap)
		return oldOneCodeMap
	}
	//else
	//不存在 就创建一个 dispather
	newOneMap := newOneCodeCallbackMap()
	p.allCallbackMap.Store(msgTypeCode, newOneMap)

	return newOneMap
}

func (p *msgProcRepertory) removeMsgHdl(msgTypeCode int32, callerUuid uint32) {
	codeCallbackMap, ok := p.allCallbackMap.Load(msgTypeCode)
	if ok {
		oldOneCodeMap := codeCallbackMap.(*OneCodeCallbackMap)
		oldOneCodeMap.remove(callerUuid)
	}
}

//export handler api

func MsgHandlerAdd(msgTypeCode int32, callback *MsgCaller) {
	instanceMsgProcRepertory.addMsgHandler(msgTypeCode, callback)
}

func MsgHandlerGet(msgTypeCode int32) IDispatcherMsg {
	return instanceMsgProcRepertory.getMsgHanlders(msgTypeCode)
}

func MsgHandlerRemove(msgTypeCode int32, callerUuid uint32) {
	instanceMsgProcRepertory.removeMsgHdl(msgTypeCode, callerUuid)
}

//************************PROC *********************************************

func (p *msgProcRepertory) put(msgTypeCode int32, processor MsgProc) {
	_, ok := p.msgProcMap.Load(msgTypeCode)
	if !ok {
		p.msgProcMap.Store(msgTypeCode, processor)
	}
}
func (p *msgProcRepertory) get(msgTypeCode int32) MsgProc {
	v, ok := instanceMsgProcRepertory.msgProcMap.Load(msgTypeCode)
	if ok {
		return v.(MsgProc)
	}
	return nil
}
func (p *msgProcRepertory) remove(msgTypeCode int32) {
	_, ok := p.msgProcMap.Load(msgTypeCode)
	if ok {
		p.msgProcMap.Delete(msgTypeCode)
	}
}

//----------------export MsgProc API

func RepertoryPut(msgTypeCode int32, processor MsgProc) {
	_, ok := instanceMsgProcRepertory.msgProcMap.Load(msgTypeCode)
	if !ok {
		instanceMsgProcRepertory.msgProcMap.Store(msgTypeCode, processor)
		fmt.Printf("注册了消息处理器：%d ， %T \n", msgTypeCode, processor)
	}
}
func RepertoryGet(msgTypeCode int32) MsgProc {
	v, ok := instanceMsgProcRepertory.msgProcMap.Load(msgTypeCode)
	if ok {
		return v.(MsgProc)
	}
	return nil
}
func RepertoryRemove(msgTypeCode int32) {
	_, ok := instanceMsgProcRepertory.msgProcMap.Load(msgTypeCode)
	if ok {
		instanceMsgProcRepertory.msgProcMap.Delete(msgTypeCode)
	}
}
