package utils

import (
	"net/source/utils/bytes"
	"errors"
)

func ParserProtoCmdHead16(data []byte) (byte, byte, error) {
	if len(data) >= 2 {
		bytearray := bytes.NewByteArray(data)
		mac, _ := bytearray.ReadByte()
		cmdCode, _ := bytearray.ReadByte()
		return mac, cmdCode, nil
	}
	return 0, 0, errors.New("proto:ParserProtoCmdHead非法的字节：指令头必须是 >=8 bytes\n")
}

func CheckContentCRC16OK(wholePack []byte) bool {

	if (len(wholePack) < 3) {
		return false
	}
	var crc16Bytes = wholePack[len(wholePack)-2:]
	var targetCH = crc16Bytes[1]
	var targetCL =  crc16Bytes[0]

	var crc16Tool = GetCrc16Tool()
	crc16Tool.PushBytes(wholePack[:len(wholePack)-2])
	rH,rL := crc16Tool.Get()
	ReleaseCrc16Tool(crc16Tool)

	if rH==targetCH && rL==targetCL{
        return  true
	}
	return true
}
