package area

import (
	"fmt"
	"net"
	"sync"
	"net/source/userapi"

	"net/source/msg/msgproc"
	"sync/atomic"
	"net/source/proto/endpoint_poller"
)

type ServiceIOCenter struct {
}

var service ServiceIOCenter

func init() {
	msgproc.GetAppTools().SvrIO = &service
}

func (s *ServiceIOCenter) ClearClient(c userapi.IClient) error {
	Remove(c)
	return nil
}

var poller = endpoint_poller.CreateEndPointPoll()
var allUsers sync.Map

func Put(c userapi.IClient) {
	Id := c.GetId()
	_, ok := allUsers.Load(Id)
	if !ok {
		allUsers.LoadOrStore(Id, c)
	}
}

func Remove(c userapi.IClient) {
	Id := c.GetId()
	_, ok := allUsers.Load(Id)
	if !ok {
		allUsers.Delete(Id)
	}
}

func (s *ServiceIOCenter) FindUser(cId int64) (userapi.IClient) {
	v, ok := allUsers.Load(cId)
	if ok {
		return v.(userapi.IClient)
	}
	return nil
}

var totalConn int64

func addConn() {
	atomic.AddInt64(&totalConn, 1)
	if totalConn%1000 == 0 {
		fmt.Println("totalConn: %d\n", totalConn)
	}
}

func PollClient(cli userapi.IClient) {
	poller.Poll(cli)
}

func AcceptorServer(ip string, basePort int, maxPorts int) {
	//启动多个接待口
	for i := 0; i < maxPorts; i++ {

		go func(indx int) {
			address := fmt.Sprintf("%s:%d", ip, basePort+indx)
			//fmt.Printf("IP IS: %s \n",address)
			pollListener, err := net.Listen("tcp", address)
			if err != nil {
				fmt.Errorf("Listen err: %s\n", err.Error())
				return
			}
			fmt.Println("server listener:", address, " has started")
			for {
				clientConn, err := pollListener.Accept()
				addConn()
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				c := msgproc.GetAppTools().UserCreator.InitClient(clientConn, 1024) //recordClient 支持并发输入
				go PollClient(c)
			}

		}(i)

	}
}
