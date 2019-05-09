package interfaces

import "net/source/userapi"

type Encoder interface{

	Encode() int

}


type Decoder interface{

	Decode(client userapi.IClient) int

}