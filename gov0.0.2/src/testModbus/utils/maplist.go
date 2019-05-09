package utils

import (
	"container/list"
)

type Keyer interface {
	GetKey() int64
}

type MapList struct {
	dataMap  map[int64]*list.Element
	dataList *list.List
}

func NewMapList() *MapList {
	return &MapList{
		dataMap:  make(map[int64]*list.Element),
		dataList: list.New(),
	}
}

func (mapList *MapList) Exists(data Keyer) bool {
	_, exists := mapList.dataMap[int64(data.GetKey())]
	return exists
}

func (mapList *MapList) Push(data Keyer) bool {
	if mapList.Exists(data) {
		return false
	}
	elem := mapList.dataList.PushBack(data)
	mapList.dataMap[data.GetKey()] = elem
	return true
}

func (mapList *MapList) Remove(data Keyer) {
	if !mapList.Exists(data) {
		return
	}
	mapList.dataList.Remove(mapList.dataMap[data.GetKey()])
	delete(mapList.dataMap, data.GetKey())
}

func (mapList *MapList) Size() int {
	return mapList.dataList.Len()
}
const STATUS_LOOP_EXIT =-1
const STATUS_LOOP_OK =0
//从前到后  正序
func (mapList *MapList) Walk(cb func(data Keyer)int) {
	for elem := mapList.dataList.Front(); elem != nil; elem = elem.Next() {
		res:=cb(elem.Value.(Keyer))
		if res==STATUS_LOOP_EXIT{
           break
		}
	}
}
//从后 到前 反序
func (mapList *MapList) WalkFromEnd(cb func(data Keyer)int) {
	for elem := mapList.dataList.Back(); elem != nil; elem = elem.Prev() {
		res:=cb(elem.Value.(Keyer))
		if res==STATUS_LOOP_EXIT{
			break
		}
	}
}


func (mapList *MapList) Shift( num int) {
	var size =mapList.Size()
	if num>size{
		num = size
	}
	for i:=0;i<num;i++{
		elem:=mapList.dataList.Front()
		mapList.Remove(elem.Value.(Keyer))
	}
}


type Entity struct {
	value int64
}

func (e Entity) GetKey() int64 {
	return e.value
}

/*
func main() {
	fmt.Println("Starting test...")
	ml := NewMapList()
	var a, b, c Keyer
	a = &Elements{"Alice"}
	b = &Elements{"Bob"}
	c = &Elements{"Conrad"}
	ml.Push(a)
	ml.Push(b)
	ml.Push(c)
	cb := func(data Keyer) {
		fmt.Println(ml.dataMap[data.GetKey()].Value.(*Elements).value)
	}
	fmt.Println("Print elements in the order of pushing:")
	ml.Walk(cb)
	fmt.Printf("Size of MapList: %d \n", ml.Size())
	ml.Remove(b)
	fmt.Println("After removing b:")
	ml.Walk(cb)
	fmt.Printf("Size of MapList: %d \n", ml.Size())
}
*/
