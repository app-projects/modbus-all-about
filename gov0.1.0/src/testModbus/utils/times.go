package utils

import "time"

func TransTime2MillSec(t time.Time) int64 {
	return t.UnixNano()/1e6
}

func SvrNowTimestamp() int64 {
	return TransTime2MillSec(time.Now())
}