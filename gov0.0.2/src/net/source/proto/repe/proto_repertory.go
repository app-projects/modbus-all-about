package repe

import (
	"sync"
	"net/source/proto/defs"
	"fmt"
)

var instanceProcessRepe *processorRepertory

type processorRepertory struct {
	protoMap sync.Map
}

func init() {
	instanceProcessRepe = &processorRepertory{}
}

func (p *processorRepertory) put(ver int32, processor *defs.ProtoProcessor) {
	_, ok := p.protoMap.Load(ver)
	if !ok {
		p.protoMap.Store(ver, processor)
	}
}
func (p *processorRepertory) get(ver int32) *defs.ProtoProcessor {
	v, ok := instanceProcessRepe.protoMap.Load(ver)
	if ok {
		return v.(*defs.ProtoProcessor)
	}
	return nil
}
func (p *processorRepertory) remove(ver int32) {
	_, ok := p.protoMap.Load(ver)
	if ok {
		p.protoMap.Delete(ver)
	}
}

func RepertoryPut(ver int32, processor *defs.ProtoProcessor) {
	_, ok := instanceProcessRepe.protoMap.Load(ver)
	if !ok {
		instanceProcessRepe.protoMap.Store(ver, processor)
	}
}
func RepertoryGet(ver int32) *defs.ProtoProcessor {
	defer func() {
		  if err:=recover() ; err!=nil{
             fmt.Println(err)
		  }
	}()
	v, ok := instanceProcessRepe.protoMap.Load(ver)
	//fmt.Printf("RepertoryGet %T",v)
	if ok {
		return v.(*defs.ProtoProcessor)  //转型有问题
	}
	return nil
}
func RepertoryRemove(ver int32) {
	_, ok := instanceProcessRepe.protoMap.Load(ver)
	if ok {
		instanceProcessRepe.protoMap.Delete(ver)
	}
}
