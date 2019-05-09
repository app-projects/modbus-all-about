package main

import (
	"fmt"
	"net"
	"net/source/proto/trans"
	"syscall"
	"io"

	"sync"
	"time"
	"testModbus/pack"
	"testModbus/connection"
	"testModbus/utils"
	"net/source/proto/trans/errcode"
	"os"
	"strconv"
	"github.com/gpmgo/gopm/modules/log"
)

var cmdheadSize = 2

const modbusHeadCmdSize = 2

var inBusHeadCmdBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, modbusHeadCmdSize)
	},
}

func poll(c *connection.Connection) int {
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
			} else if err == io.EOF { //对端 关闭

				fmt.Println("退出善后  do something exit for client: id", c.GetId())
				goto endflag
			}
		} else {

			if nread > 0 { //盛到了水
				c.GetToolTotalCache().Write(c.GetToolBucket()[:nread]) //buffer 标准库 是对 []byte 切片的 控制操作封装
			}

			//分析
			if /* c.GetNewPackFlag() && */ c.GetToolTotalCache().Available() >= cmdheadSize { // 就说明一个数据包完成
				cmdHead := c.GetToolTotalCache().BytesAvailable()[:cmdheadSize]
				_, fnCode, err := utils.ParserProtoCmdHead16(cmdHead[:]) //optBucketBuf就是PROTO_HEAD_CMD_SIZE大小
				if err != nil {
					fmt.Errorf(err.Error())
					continue
				}

				decoder, err := trans.GetDecoderPluginByFnCodeResp(int32(fnCode))
				if err == nil {
					res := decoder.Decode(c)
					if res == errcode.ERR_TRNAS_DECODE_TRY_AGAIN {
						continue
					} else if errcode.TRNAS_DECODE_COMPLETE == res {
					}

				}

			}

		}
	}
endflag:
	fmt.Println("call exit")
	defer c.Exit()
	return 0

}

func processModbusPack(con *connection.Connection) {

}

func AskRoutine(con *connection.Connection) {
	askBytes := pack.EncodeAskProto06H()

	for {

		con.Conn.Write(askBytes)
		//fmt.Println("查询指令下发成功")
		time.Sleep(time.Second * 5)
	}
}

func fixedClient(clientConn net.Conn) {
	conn := connection.NewConnection(clientConn, 1024)
	go AskRoutine(conn)
	go poll(conn)
	//go processModbusPack(conn)
}

func acceptorServer(ip string, basePort int, maxPorts int) {
	//启动多个接待口
	for i := 0; i < maxPorts; i++ {

		go func(indx int) {
			address := fmt.Sprintf("%s:%d", ip, basePort+indx)
			//fmt.Printf("IP IS: %s \n",address)
			pollListener, err := net.Listen("tcp", address)
			if err != nil {
				log.Fatal("Listen err: %s\n", err.Error())
				return
			}
			fmt.Println("server listener:", address, " has started")
			for {
				clientConn, err := pollListener.Accept()
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				fixedClient(clientConn)
			}

		}(i)

	}
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("please input format : ./app.exe ip baseport portNum")
		return
	}
	ip := os.Args[1]
	basePort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("baseport is a int num >0")
	}
	portNum, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("portNum is a int num >0")
	}
	fmt.Printf("fmt:ip=%s,basePort:%d,portNum=%d \n", ip, basePort, portNum)
	acceptorServer(ip, basePort, portNum)

	select {}

}
