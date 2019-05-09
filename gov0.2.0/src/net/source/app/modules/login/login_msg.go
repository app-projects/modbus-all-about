package login

import (
	"net/source/msg/msgproc"
)

type RepoMsg struct {
	msgproc.BaseMsg
	DevId     int64
	IsCharged byte
	Name      string
}
