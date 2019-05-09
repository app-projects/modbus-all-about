package entity

import (
	"datastore/web/store"
	"fmt"
)

type iDevlog struct {
	Mac           byte                   `json:"mac,omitempty"`
	FunCode       byte                   `json:"funcode,omitempty"`
	Ip            string                 `json:ip,omitempty"`
	Timestamp     int64                  `json:"timestamp,omitempty"`
	Templ         byte                   `json:"templ"`
	DataFieldsLen byte                   `json:"data_fields_len"`
	DataFieldsmap map[string]interface{} `json:"data_fieldsmap"`
	//DataFields []int `json:datafields` 数组本质是内存块 ，原生的数组作为序列号，不好跨语言；所以数组一定要给定长度 给定长度+map
	//所以通常 采用  给定长度+map
}

type DevlogList struct {
	List []iDevlog `json:"list,omitempty"`
}

func NewLostList() *DevlogList {
	list := DevlogList{}
	list.List = make([]iDevlog, 0)
	return &list
}

func (this *iDevlog) GetKey() int64 {
	return this.Timestamp
}

func (this *iDevlog) SetMac(mac string) {
	var l = len(mac)
	if l > 0 {
		this.Mac = []byte(mac)[0]
	} else {
		this.Mac = 0
	}
}
func (this *iDevlog) SetIp(ip string) {
	this.Ip = ip
}
func (this *iDevlog) SetValue(v interface{}) {
	this.Templ = v.(byte)
}
func (this *iDevlog) SetTimeStmap(t int64) {
	this.Timestamp = t
}

func (this *iDevlog) GetMac() string {
	return byte2String(this.Mac)
}
func (this *iDevlog) GetIp() string {
	return this.Ip
}
func (this *iDevlog) GetValue() interface{} {
	return this.Templ
}
func (this *iDevlog) GetTimeStmap() int64 {
	return this.Timestamp
}

func NewDevQueryLog() store.Devlog {
	return &iDevlog{}
}

func byte2String(b byte) string {
	var s = fmt.Sprintf("%d", b)
	fmt.Println(s)
	return s
}
