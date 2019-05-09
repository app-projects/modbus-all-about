package endpoint_poller

import (
	"fmt"
	"net/source/userapi"
	"net"
	"syscall"
	"io"
	"net/source/proto/pools"
	"net/source/proto"
	"net/source/proto/trans"
)

type modBusRtuPoller struct {
}

var cmdheadSize = 2 //<指令头>  mac,code    int16

func (this *modBusRtuPoller) Poll(c userapi.IClient) int {

	if !c.IsOpen() {
		return -1
	}
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
			if c.GetNewPackFlag() && c.GetToolTotalCache().Available() >= cmdheadSize { // 就说明一个数据包完成

				cmdHead := pools.ModBusHeadCmdBytesPool.Get().([]byte)
				nread, er := c.GetToolTotalCache().Read(cmdHead)
				if nread > 0 && er == nil {
					//fmt.Println("after toolTotalCache.Len():", c.toolTotalCache.Len())
					mac, fnCode, err := proto.ParserProtoCmdHead16(cmdHead[:]) //optBucketBuf就是PROTO_HEAD_CMD_SIZE大小
					pools.HeadCmdBytesPool.Put(cmdHead)

					if err != nil {
						fmt.Errorf(err.Error())
						continue
					}

					decoder, err := trans.GetDecoderPluginByFnCode(int32(fnCode))
					if err == nil {
						c.SetNewPackFlag(false)
						res := decoder.Decode(c,mac)
						if res == 0 { //当前包解析完成
							c.ResetRecvNewPack()
							continue
						}
					}

				}

			}

		}
	}
endflag:

	defer c.Exit()
	return 0
}
