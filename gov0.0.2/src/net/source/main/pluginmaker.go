package main

import (
	"net/source/proto/trans/interfaces"
	"net/source/proto/trans/decode_suit/mb_rtu_03h_decoder"
)


func CreateIntance() interfaces.Decoder {
	return  mb_rtu_03h_decoder.CreateIntance()
}

func ReleaseInstance(ins interfaces.Decoder)  {
	mb_rtu_03h_decoder.ReleaseInstance(ins)
}