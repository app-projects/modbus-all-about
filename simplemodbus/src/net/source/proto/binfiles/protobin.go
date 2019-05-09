package binfiles

import (
	"bytes"
	"net/source/proto/repe"
	"errors"
	"fmt"
	"net/source/proto/pools"
	"sync"
	bytes2 "net/source/utils/bytes"
	"net/source/userapi"
)

type ProtoBinPack struct {
	Ver          int32
	len          int32 //设想发送长度
	Rlen         int32 //实际获得长度
	Bytes        []byte
	refPoolBytes []byte //Bytes字段的 吃对象引用

	S_cliconnId int64 //服务器端 的客户连接id
	Buf         *bytes.Buffer
}

var ProtobinPool = sync.Pool{
	New: func() interface{} {
		p := &ProtoBinPack{
			Ver:         0,
			len:         0,
			S_cliconnId: -1,
			Rlen:        0,
		}
		p.Bytes = nil
		p.Buf = nil
		return p
	},
}

var ModBusProtobinPool_03 = sync.Pool{
	New: func() interface{} {
		p := &Mod03_ProtoBinPack{
		}
		return p
	},
}

var ModBusProtobinPool_06 = sync.Pool{
	New: func() interface{} {
		p := &Mod06_ProtoBinPack{
		}
		return p
	},
}

var ModBusProtobinPool_10 = sync.Pool{
	New: func() interface{} {
		p := &Mod10_ProtoBinPack{
		}
		return p
	},
}


func (proto *ProtoBinPack) reset() {
	proto.Ver = 0
	proto.len = 0

	proto.Rlen = 0
	proto.S_cliconnId = -1
	if proto.refPoolBytes != nil {
		pools.CtlBytesSlicePool.Put(proto.refPoolBytes)
		proto.refPoolBytes = nil
	}
	proto.Buf.Reset()
	proto.Bytes = nil
	proto.Buf = nil

}

func (this *ProtoBinPack) SetClientId(id int64) {
	this.S_cliconnId = id
}
func (this *ProtoBinPack) GetClientId() int64 {
	return this.S_cliconnId
}

func (proto *ProtoBinPack) Release() {
	pools.CtlBytesSlicePool.Put(proto.refPoolBytes)
	proto.reset()
}

func CreateProtoBin(ver, sizel int32, cliconnId int64) *ProtoBinPack {
	tmp := ProtobinPool.Get().(*ProtoBinPack)
	tmp.refPoolBytes = pools.CtlBytesSlicePool.Get().([]byte)

	tmp.Bytes = tmp.refPoolBytes[0:sizel]
	tmp.Buf = bytes.NewBuffer(tmp.Bytes)

	tmp.Ver = ver
	tmp.len = sizel
	tmp.S_cliconnId = cliconnId
	return tmp
}

type Mod03_ProtoBinPack struct {
	Modbase_ProtoBinPack
	StartRegH8 byte
	StartRegL8 byte
	ReadRegH   byte
	ReadRegL   byte

	Crc16 int16
}


/*
[设备地址] +[命令号03H] +       <指令头>
[返回的字节个数] +
[数据1] +
[数据2] +
...+ [数据n] +

[CRC校验的低8位] + [CRC校验的高8位]*/

type Mod03_ProtoBinPackResp struct {
	Modbase_ProtoBinPack
	DataFieldLength byte
	DataFields []byte
	Crc16 int16
}


var ModBusProtobinRespPool_03 = sync.Pool{
	New: func() interface{} {
		p := &Mod03_ProtoBinPackResp{
			DataFields:make([]byte,50),
			DataFieldLength:0,
		}
		return p
	},
}




type Mod06_ProtoBinPack struct {
	Modbase_ProtoBinPack    `json:"-"`
	RegSetterH8  byte		`json:"reg_setter_h_8"`
	RegSetterL8  byte		`json:"reg_setter_l_8"`
	DataH8 byte				`json:"data_h_8"`
	DataL8 byte				`json:"data_l_8"`
	Crc16  uint16			`json:"-"`
}

func (this *Mod06_ProtoBinPack) SetFnCode(code byte) {
	this.FnCode = code

}
func (this *Mod06_ProtoBinPack) SetMac(mac byte) {
	this.Mac = mac
}
func (this *Mod06_ProtoBinPack) Reset() {
	//this.Mac = mac
}

type Modbase_ProtoBinPack struct {
	userapi.IModBusProtoBinPack   `json:"-"`
	clientId int64    `json:"client_id"`
	Mac      byte	   `json:"mac"`
	FnCode   byte      `json:"fn_code"`
}

func (this *Modbase_ProtoBinPack) GetFnCode() byte {
	return this.FnCode
}

func (this *Modbase_ProtoBinPack) SetClientId(id int64) {
	this.clientId = id
}

func (this *Modbase_ProtoBinPack) GetClientId() int64 {
	return this.clientId
}

func (this *Modbase_ProtoBinPack) SetFnCode(code byte) {
	this.FnCode = code

}
func (this *Modbase_ProtoBinPack) SetMac(mac byte) {
	this.Mac = mac
}
func (this *Modbase_ProtoBinPack) GetMac() byte {
	return this.Mac
}



func (this *Modbase_ProtoBinPack) Reset() {
	//this.Mac = mac
}

type Mod10_ProtoBinPack struct {
	Modbase_ProtoBinPack
	StartRegH8 byte
	StartRegL8 byte

	RegNum       byte
	DataFields16 []int16
	Crc16        uint16
}

var creator map[byte]*sync.Pool
var code03 byte = 0x03
var code06 byte = 0x06
var code10 byte = 0x10


var creator_resp map[byte]*sync.Pool
func init() {

	creator = make(map[byte]*sync.Pool, 5)

	creator[code03] = &ModBusProtobinPool_03
	creator[code06] = &ModBusProtobinPool_06
	creator[code10] = &ModBusProtobinPool_10

	creator_resp = make(map[byte]*sync.Pool, 5)
	creator_resp[code03] = &ModBusProtobinRespPool_03
}

func CreateModBusProtoBin(mac, fnCode byte, cliconnId int64) userapi.IModBusProtoBinPack {
	creatorIns := creator[fnCode]
	if creatorIns == nil {
		return nil
	}
	tmp := creatorIns.Get().(userapi.IModBusProtoBinPack)
	tmp.Reset()
	tmp.SetFnCode(fnCode)
	tmp.SetMac(mac)
	tmp.SetClientId(cliconnId)
	return tmp
}

func CreateModBusProtoBinResp(mac, fnCode byte, cliconnId int64) userapi.IModBusProtoBinPack {
	creatorIns := creator_resp[fnCode]
	if creatorIns == nil {
		return nil
	}
	tmp := creatorIns.Get().(userapi.IModBusProtoBinPack)
	tmp.Reset()
	tmp.SetFnCode(fnCode)
	tmp.SetMac(mac)
	tmp.SetClientId(cliconnId)
	return tmp
}

// translate  proto to  app msg

var dstBytes = [1024]byte{0}

func (proto *ProtoBinPack) Len() int32 {

	return proto.len

}

func (proto *ProtoBinPack) GetBytes() []byte {
	return proto.Bytes
}

func (proto *ProtoBinPack) SetRLen(rl int32) {
	proto.Rlen = rl
}

func (proto *ProtoBinPack) Parse() error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("parse proto bin err:", err)
		}
	}()

	processor := repe.RepertoryGet(proto.Ver)
	if processor != nil {
		processor.Ver = proto.Ver
		if proto.Rlen > 0 {
			//bytesEntity_Ref := pools.CtlBytesSlicePool.Get().([]byte)
			useBytes := dstBytes[0:proto.Rlen]
			var bt = bytes2.NewByteArray(useBytes)
			bt.Write(proto.Bytes)
			processor.DepartDomain(processor.Ver, proto.S_cliconnId, bt.Bytes())
			//pools.CtlBytesSlicePool.Put(bytesEntity_Ref)
		}
		return nil
	}
	return errors.New(fmt.Sprintln("err: processor := repe.RepertoryGet(proto.Ver) :processor==nil not exist %d ", proto.Ver))
}
