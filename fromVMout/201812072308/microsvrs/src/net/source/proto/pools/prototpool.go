package pools

import (
	"sync"
	"net/source/utils/bytes"
)

var bytesChunkBufferSize = 1024 //规定指令内容体长度     (指令头(协议号+指令体长度)+指令体)  ；容易

var CtlBytesSlicePool = sync.Pool{
	New: func() interface{} {
		return make([]byte, bytesChunkBufferSize)
	},
}

var businessBytesChunkBufferSize = 512
var BusinessBytesSlicePool = sync.Pool{
	New: func() interface{} {
		return make([]byte, businessBytesChunkBufferSize)
	},
}

const headCmdSize =8
var HeadCmdBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, headCmdSize)
	},
}


const modbusHeadCmdSize =2
var ModBusHeadCmdBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, modbusHeadCmdSize)
	},
}


var ByteArrayPool = sync.Pool{
	New: func() interface{} {
		return  bytes.NewByteArray([]byte{})
	},
}
