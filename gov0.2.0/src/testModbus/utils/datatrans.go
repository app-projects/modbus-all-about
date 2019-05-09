package utils

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

func Uint642Byte(i uint64) byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b[7]
}

func Uint162Byte(i uint16) (bH byte,bL byte) {
	bL = byte(i)
	bH = byte((i>>8))
	return
}

func Byte2Uint64(b byte) uint64 {
	bts := make([]byte, 8)
	bts[7] = b
	return uint64(binary.BigEndian.Uint64(bts))
}

func Int64_2String(i int64) string {
	return fmt.Sprintf("%d", i)
}

func Byte_2HexString(i byte) string {
	return fmt.Sprintf("%x", i)
}

//取得低8位
func Int64_2Byte(i int64) byte {
	s := fmt.Sprintf("%d", i)
	u, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		fmt.Println("result:", u)
	}
	return Uint642Byte(u)
}

func Bytes2Uint16(bH byte, bL byte) uint16 {
	var h uint16 = 0
	h = uint16(bH)
	var l uint16 = 0
	l = uint16(bL)
	return (h << 8) | l
}
