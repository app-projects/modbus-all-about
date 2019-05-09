package proto

import (
	"fmt"
	"net/source/proto/defs"
	"net/source/proto/repe"
)

func AssembleProtoProc(ver int32, headProc defs.ProtoCtlHeadProc, bodyProc defs.ProtoBodyProc) {
	p := defs.ProtoProcessor{}
	p.Ver = ver
	p.BodyProc = bodyProc
	p.CtlHeaderProc = headProc
	repe.RepertoryPut(ver, &p)
	fmt.Println("注册了 头处理器 体处理器", ver)
}
