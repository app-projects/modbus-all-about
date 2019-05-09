package proto

import (
	"errors"
	"net/source/utils/bytes"
	"hash/crc32"
	"io"
)

func ParserProtoCmdHead(data []byte) (int32, int32, error) {

	if len(data) >= 8 {
		bytearray := bytes.NewByteArray(data)
		first, _ := bytearray.ReadInt32()
		sec, _ := bytearray.ReadInt32()
		return first, sec, nil
	}
	return -1, -1, errors.New("proto:ParserProtoCmdHead非法的字节：指令头必须是 >=8 bytes\n")
}

func ParserProtoCmdHead16(data []byte) (int8, int8, error) {
	if len(data) >= 16 {
		bytearray := bytes.NewByteArray(data)
		mac, _ := bytearray.ReadInt8()
		cmdCode, _ := bytearray.ReadInt8()
		return mac, cmdCode, nil
	}
	return -1, -1, errors.New("proto:ParserProtoCmdHead非法的字节：指令头必须是 >=8 bytes\n")
}

func CRC32GenId(uniqueStr string) uint32 {
	ieee := crc32.NewIEEE()
	io.WriteString(ieee, uniqueStr)
	return ieee.Sum32()
}
