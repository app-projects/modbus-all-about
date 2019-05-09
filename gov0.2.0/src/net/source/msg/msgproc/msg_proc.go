package msgproc

import (
	"net/source/proto/defs"
)

type MsgProc interface {
	Proc(ctlMsg *defs.CtlMsg, businessBytes []byte) (interface{}, error) //返回是 msg err
}
