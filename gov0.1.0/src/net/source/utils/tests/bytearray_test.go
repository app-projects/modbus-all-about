package main

import (
	"fmt"
	bytes2 "net/source/utils/bytes"
	"testing"
)



func BenchmarkStringJoin1(t *testing.B) {
	var ba * bytes2.ByteArray= bytes2.NewByteArray([]byte{})

	ba.WriteBytes([]byte("abc"))
	ba.WriteByte('A')
	ba.WriteBool(true)
	ba.WriteBool(false)
	ba.WriteInt8(11)
	ba.WriteInt16(123)
	ba.WriteInt32(123)
	ba.WriteInt64(513)
	ba.WriteFloat32(123.456)
	ba.WriteFloat64(456.789)
	ba.WriteString("hello ")
	ba.WriteUTF("world!")

	bytes := make([]byte, 3)
	fmt.Println(ba.ReadBytes(bytes, 3, 0))
	fmt.Println(ba.ReadByte())
	fmt.Println(ba.ReadBool())
	fmt.Println(ba.ReadBool())
	fmt.Println(ba.ReadInt8())
	fmt.Println(ba.ReadInt16())
	fmt.Println(ba.ReadInt32())
	fmt.Println(ba.ReadInt64())
	fmt.Println(ba.ReadFloat32())
	fmt.Println(ba.ReadFloat64())
	fmt.Println(ba.ReadString(6))
	fmt.Println(ba.ReadUTF())

	byte,err := ba.ReadByte()
	if err == nil{
		fmt.Println(byte)
	}else{
		fmt.Println("end of file")
	}

	ba.Seek(3)    //back to 3
	fmt.Println(ba.ReadByte())
	ba.Seek(39)    //back to 39
	fmt.Printf("ba has %d bytes available!\n", ba.Available())
	fmt.Println(ba.ReadUTF())
}