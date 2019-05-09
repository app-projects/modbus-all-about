package datalayer

import (
	"errors"
	"debug/dwarf"
)

type LayerNode struct {
	pre  *LayerNode
	next *LayerNode

	layerId  byte
}
//数据发往的方向：发向用户的方向
func (l *LayerNode)UpStream()  {


}

//数据发往的方向：硬件方向
func (l *LayerNode)DownStream()  {


}

type SwitchLayer struct {
	 LayerNode

}

type MsgDispatchLayer struct {
	  LayerNode

}

type DataLayer struct {
	 first     *LayerNode
	 last     *LayerNode
	 layersDict map[byte]*LayerNode
}

var dataLayer DataLayer

func init(){
   dataLayer = DataLayer{
	        first:nil,
	        last:nil,
   			layersDict:make(map[byte]*LayerNode,0),
   }

}

func (dataLayer * DataLayer) InjectLayer(layer *LayerNode) error {
	if layer==nil {
		return errors.New("inject layer nil")
	}
	if dataLayer.first==nil{
		dataLayer.first =layer
		dataLayer.last = layer
	}else {
		//增加一层
       if dataLayer.last!=nil{
       	  dataLayer.last.next = layer
       	  layer.pre = layer
       	  //保存搜索 索引
		   dataLayer.layersDict[layer.layerId] = layer
	   }

	}





}

