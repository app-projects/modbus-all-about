package entity

import "datastore/web/store"

type iDevlog struct {
	Mac       string `json:"mac,omitempty"`
	FunCode   uint8  `json:"funcode,omitempty"`
	Ip        string `json:ip,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Templ     int    `json:"templ"`
}

type DevlogList struct {
	List []iDevlog `json:"list,omitempty"`
}

func NewLostList() *DevlogList {
	list := DevlogList{}
	list.List = make([]iDevlog, 2)
	return &list
}

func (this *iDevlog) GetKey() int64 {
	return this.Timestamp
}

func (this *iDevlog) SetMac(mac string) {
	this.Mac = mac
}
func (this *iDevlog) SetIp(ip string) {
	this.Ip = ip
}
func (this *iDevlog) SetValue(v interface{}) {
	this.Templ = v.(int)
}
func (this *iDevlog) SetTimeStmap(t int64) {
	this.Timestamp = t
}

func (this *iDevlog) GetMac() string {
	return this.Mac
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