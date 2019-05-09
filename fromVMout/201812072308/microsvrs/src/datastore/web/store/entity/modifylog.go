package entity

import "datastore/web/store"

type ModifyLog06 struct {
	RegSetterH8 string `json:"regsetterh8"`
	RegSetterL8 string `json:"regsetterl8"`
	DataH8      string `json:"datah8"`
	DataL8      string `json:"datal8"`
	Value       int16 `json:"-"`

	clientId  string  `json:"clientid"`
	Mac       string `json:"mac"`
	FnCode    string   `json:"fncode"`
	Timestamp int64  `json:"timestamp"`
	Ip        string  `json:"ip"`
}

func (this *ModifyLog06) GetKey() int64 {
	return this.Timestamp
}

func (this *ModifyLog06) SetMac(mac string) {
	this.Mac = mac
}
func (this *ModifyLog06) SetIp(ip string) {
	this.Ip = ip
}
func (this *ModifyLog06) SetValue(v interface{}) {
	this.Value = v.(int16)
}
func (this *ModifyLog06) SetTimeStmap(t int64) {
	this.Timestamp = t
}

func (this *ModifyLog06) GetMac() string {
	return this.Mac
}
func (this *ModifyLog06) GetIp() string {
	return this.Ip
}
func (this *ModifyLog06) GetValue() interface{} {
	return this.Value
}
func (this *ModifyLog06) GetTimeStmap() int64 {
	return this.Timestamp
}

func NewDevModifyLog() store.Devlog {
	return &ModifyLog06{}
}
