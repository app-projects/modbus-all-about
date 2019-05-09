package main

import (
	"net"
	"fmt"
	"syscall"
	"io"
	"net/source/proto/trans"
	"sync"
	"errors"
	"net/source/utils/bytes"
	"testModbus/connection"
	"net/source/proto/trans/errcode"
	"os"
	"strconv"
)

/****
[设备地址]+[命令号03H] +       <指令头>    16bit  2byte
[起始寄存器地址高8位] +[低8位] +           16bit    2byte   //offset
[读取的寄存器数高8位] +[低8位] +           16bit    2byte    //length
[CRC校验的低8位] + [CRC校验的高8位]    16bit CRC   2byte

**/

type proto_req struct {
	Mac     int8
	FunCode int8
}

type proto3_req struct {
	proto_req
	OffsetRegH int8
	OffsetRegL int8
	ReadRegH   int8
	ReadRegL   int8
	Crc16      int16
}

func listenAsk(con net.Conn) {

}

var cmdheadSize = 2

const modbusHeadCmdSize = 2

var inBusHeadCmdBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, modbusHeadCmdSize)
	},
}

func parserProtoCmdHead16(data []byte) (byte, byte, error) {
	if len(data) >= 2 {
		bytearray := bytes.NewByteArray(data)
		mac, _ := bytearray.ReadByte()
		cmdCode, _ := bytearray.ReadByte()
		return mac, cmdCode, nil
	}
	return 0, 0, errors.New("proto:ParserProtoCmdHead非法的字节：指令头必须是 >=8 bytes\n")
}

func processConn(con net.Conn) {
	c := connection.NewConnection(con, 1024)
	go pollClient(c)

}

func pollClient(c *connection.Connection) int {
	/*if c.IsOpen() == false {
		return -1
	}*/
	for {
		nread, err := c.GetConn().Read(c.GetToolBucket())
		if err != nil { //桶子没有 获得数据
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() { //空闲等待时间，空档期时间
				continue
			}

			if err == syscall.EAGAIN || err == syscall.EINTR {
				fmt.Println(".conn.Read err:", err.Error())
				//c.GetConn().Write([]byte("hello please send data to server process"))
				continue
			} else if err == io.EOF || err == syscall.ECONNABORTED { //对端 关闭

				fmt.Println("退出善后  do something exit for client: id", c.GetId())
				goto endflag
			}
		} else {

			if nread > 0 { //盛到了水
				c.GetToolTotalCache().Write(c.GetToolBucket()[:nread]) //buffer 标准库 是对 []byte 切片的 控制操作封装
			}

			//分析
			if/* c.GetNewPackFlag() && */c.GetToolTotalCache().Available() >= cmdheadSize { // 就说明一个数据包完成

				//cmdHead := inBusHeadCmdBytesPool.Get().([]byte)
				cmdHead := c.GetToolTotalCache().BytesAvailable()[:cmdheadSize]

				//nread, er := c.GetToolTotalCache().Read(cmdHead)

				_, fnCode, err := parserProtoCmdHead16(cmdHead[:]) //optBucketBuf就是PROTO_HEAD_CMD_SIZE大小
				//inBusHeadCmdBytesPool.Put(cmdHead)
				if err != nil {
					fmt.Errorf(err.Error())
					continue
				}

				decoder, err := trans.GetDecoderPluginByFnCode(int32(fnCode))
				if err == nil {
					//c.SetNewPackFlag(false)
					res := decoder.Decode(c)
					if res==errcode.ERR_TRNAS_DECODE_TRY_AGAIN{
						continue
					}else if errcode.TRNAS_DECODE_COMPLETE==res{
						//c.ResetRecvNewPack()
					}

				}

			}

		}
	}
endflag:

	defer c.Exit()
	return 0

}

func main() {


	if len(os.Args) < 3 {
		fmt.Println("please input format : ./app.exe ip port ")
		return
	}
	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("baseport is a int num >0")
	}

    addr:=fmt.Sprintf("%s:%d",ip,port)
    fmt.Println(" connect target:",addr)
	con, err := net.Dial("tcp", addr)

	if err == nil {
		fmt.Println("连接成功：", con.RemoteAddr())

		processConn(con)

	} else {
		fmt.Println(err.Error())
	}

	select {}

}
