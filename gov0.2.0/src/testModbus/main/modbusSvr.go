package main

import (
	"fmt"
	"net"
	"net/source/proto/trans"
	"syscall"
	"io"

	"sync"
	"testModbus/connection"
	"testModbus/utils"
	"net/source/proto/trans/errcode"
	"os"
	"strconv"
	"runtime"
	"modbusSvrRestful/web"
	"testModbus/outinterface"
	"log"
	"testModbus/simulator"
)

var cmdheadSize = 2

const modbusHeadCmdSize = 2

var inBusHeadCmdBytesPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, modbusHeadCmdSize)
	},
}

func poll(c *connection.Connection) int {

	for {
		if c.IsOpen() == false {
			return -1
		}

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
				ok := filterPackIgnoreCompleted(fnCode, c)

				if ok {
					decoder, err := trans.GetDecoderPluginByFnCodeResp(int32(fnCode))
					if err == nil && decoder != nil {
						res := decoder.Decode(c)
						if res == errcode.ERR_TRNAS_DECODE_TRY_AGAIN {
							continue
						} else if errcode.TRNAS_DECODE_COMPLETE == res {
						}

					} else {
						fmt.Println(err.Error())
						continue
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

var heartbitSize = 10
var heartBytes = make([]byte, 10)
//返回值表示，true 是 被拦截过
func filterPackIgnoreCompleted(fnCode byte, c *connection.Connection) bool {
	switch fnCode {
	case 119:
		if (c.GetToolTotalCache().Available() >= 10) {
			c.GetToolTotalCache().Read(heartBytes) //抛弃心跳包
			log.Println("丢弃新跳包 10个字节....")
			return true
		} else {
			return false
		}
	default:
		return true
	}
	return true
}

func LbDivideMac2Conn(mac byte, conn *connection.Connection) {
	connection.GetDevConnRouter().Bind(mac, conn)
}

func fixedClient(clientConn net.Conn) {
	//prepare ask
	conn := connection.NewConnection(clientConn, 1024)

	LbDivideMac2Conn(utils.Int64_2Byte(conn.GetId()), conn)
	log.Println("a connectoin id:", conn.GetId(), " coming")

	go conn.SynRequestQueue()
	//go simulator.CommitGetSysInfoAsk(conn)
	go simulator.CommitGetAppInfoAsk(conn)
	go outinterface.TickComingMsg()
	go simulator.CommitModifyAppInfo(conn)
	go poll(conn)
	go conn.ExecuteResponse()
	go conn.AskRoutine()

}

//var wg sync.WaitGroup

func acceptorServer(ip string, basePort int, maxPorts int) {
	//启动多个接待口
	for i := 0; i < maxPorts; i++ {
		go func(indx int) {
			address := fmt.Sprintf("%s:%d", ip, basePort+indx)
			pollListener, err := net.Listen("tcp", address)
			if err != nil {
				fmt.Errorf("Listen err: %s\n", err.Error())
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
	runtime.GOMAXPROCS(runtime.NumCPU())

	go acceptorServer(ip, basePort, portNum)
	go web.SvrWebRestfulMain()

	select {}

}
