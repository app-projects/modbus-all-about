package echocontext

import (
	"fmt"
	"time"
	"container/list"
	"log"
)

const STATE_READY = 0
const STATE_SENDED = 1
const STATE_RECVED = 2
const STATE_TIMEOUT = -1

type RespHandler = func(args ...interface{}) interface{}
type NetworkCmd interface {
	GetTimestamp() int64
	SetTimestamp(t int64)
	SetRspdArgs(args ...interface{})
	DoResponse()
	GetSndPack() []byte

	SetRespdHandler(h RespHandler)
	GetRespdHandler() RespHandler

	GetState() int
	SetState(state int)

	GetPack() []byte
	SetPack(pck []byte)
}

type NetworkData interface {
	GetTimestamp() int64
	GetData() interface{}
}

type EchoContext struct {
	Queue           *list.List
	AskGoHead       chan interface{}
	ResponseTasks   chan NetworkCmd
	networkDataChan chan NetworkData
	SendNetworkFn   func(NetworkCmd)

	timeoutMillSec int64
	timer          *time.Timer
	timeoutDur     time.Duration

	state int
}

const ECHO_CTX_STATE_CLOSE = 1
const ECHO_CTX_STATE_RUNNING = 0

func (ctx *EchoContext) SetState(s int) {
	ctx.state = s

}
func (ctx *EchoContext) SetTimeoutMillSec(timeout int64) {
	ctx.timeoutMillSec = timeout
	ctx.timeoutDur = time.Duration(time.Millisecond) * time.Duration(timeout)
	if ctx.timer == nil {
		ctx.timer = time.NewTimer(ctx.timeoutDur)
	} else {
		ctx.timer.Reset(ctx.timeoutDur)
	}
	ctx.timer.Stop()
}

func (this *EchoContext) printInfo() {
	fmt.Println("NetworkCmd    数据命令 队列饱和度:", this.Queue.Len())
	fmt.Println("ResponseTasks 网络包处理任务 ：队列饱和度:", len(this.ResponseTasks))
	fmt.Println("NetworkData   响应网络接收包 ：    队列饱和度:", len(this.networkDataChan))
}

func (ctx *EchoContext) pushRespondCmd(cmd NetworkCmd) {
	cmd.SetState(STATE_RECVED)
	ctx.ResponseTasks <- cmd
	/*
		var t = time.NewTimer(selectDelayMillSecTimeout)
		defer t.Stop()
		for {
			t.Reset(selectDelayMillSecTimeout)
			select {
			case ctx.ResponseTasks <- cmd:
				break
			}
			select {
			case <-t.C:
				log.Println("pushRespondCmd: tick ....")
			default:
				log.Println("pushRespondCmd: no block ....")
			}
		}*/
}
func (ctx *EchoContext) CommitRequest(t NetworkCmd) {
	t.SetState(STATE_READY)
	ctx.Queue.PushBack(t)
}

func sendNetwork(ctx *EchoContext, t NetworkCmd) {
	if ctx.SendNetworkFn != nil {
		t.SetTimestamp(time.Now().Unix() * 1000)
		ctx.SendNetworkFn(t)
		t.SetState(STATE_SENDED)
	}
	//fmt.Println("send to network timestapm------:", t.GetTimestamp())
}

//modbus 是一个echo模型 终端有反馈，再发才有意义
func (ctx *EchoContext) AskRoute() {
	var willReq *list.Element = nil
	var willReqCmd NetworkCmd = nil
	var timeOutReq *list.Element = nil
	var preRequest_PckTimestamp int64
	for {
		if (ctx.state == ECHO_CTX_STATE_CLOSE) {
			//存盘 持久化
			goto ExitFlag
		}
		//---------------------------------------------发送 逻辑 start------------------------------------------------------------------------------------------

		if timeOutReq == nil {
			var le = ctx.Queue.Len()
			if le > 0 && willReq == nil {
				willReq = ctx.Queue.Front()
			}
		} else {
			willReq = timeOutReq
			timeOutReq = nil //超时包完成交接工作 end
		}

		if willReq != nil { //发送 逻辑
			willReqCmd = willReq.Value.(NetworkCmd)
			if willReqCmd.GetState() == STATE_TIMEOUT || willReqCmd.GetState() == STATE_READY {
				sendNetwork(ctx, willReqCmd)
				preRequest_PckTimestamp = willReqCmd.GetTimestamp()
			}
		}
		//------------------------------------------------发送 逻辑 end--------------------------------------------------------------------------------------------
		//============================上下同步过程 ask 和接受回答======== 一问 一答 模式==================================

		//******************************************处理网络包逻辑 start**************************************************************************
		if willReqCmd != nil && (willReqCmd.GetState() == STATE_SENDED) { //有已经发送了数据包，才有一下操作
			ctx.timer.Reset(ctx.timeoutDur)
			select {
			case responseData := <-ctx.networkDataChan: //next request flag signal
				//如果存在 迷途包 则丢弃，当做不存在 ，因为下面本身具有  对响应超时的处理（补救重发） ，客户端 断点 或者 异常容易出现这个情况
				if (preRequest_PckTimestamp > 0 && (responseData.GetTimestamp()-preRequest_PckTimestamp > ctx.timeoutMillSec)) { // 对数据 超时校验 ，万一网络原因 不符合协议，来数据包当然 2MSL 到了该包自动在网络消失
					fmt.Println("drop a timeout pack...")
				} else {
					//process cmd
					willReqCmd.SetRspdArgs(responseData)
					ctx.pushRespondCmd(willReqCmd)
					ctx.Queue.Remove(willReq)
					willReq = nil  //willReq=nil 开启新包 ，故只能再次 被赋值 nil   ，如果放在下面 外面 赋值nil ，不能因为 丢弃包，而开启新包请求（在加断点情况会出现异常）
				}
				//willReq = nil 如果放在 外面 赋值nil ，那么   //willReq=nil 开启新包 ，故只能再次 被赋值 nil


				fmt.Println("responseDataresponseDataresponseDataresponseData-----len:", len(ctx.networkDataChan))

				//网络响应超时处理：补救重发
			case <-ctx.timer.C: // 对行为逻辑的过滤  时间阀值 对外行为
				//标记为超时包
				//重新起航
				if willReq != nil {
					fmt.Println("超时重发timeout retry send pack timestapm------")
					timeOutReq = willReq //获得超时包  开启上面超时包的交接工作
					willReqCmd.SetState(STATE_TIMEOUT)
				} else {
					log.Println("AskRoute tick ........")
				}
				//<-ctx.networkDataChan
				willReq = nil
			}
			ctx.timer.Stop()
		}

		//******************************************处理网络包逻辑 end**************************************************************************

		time.Sleep(time.Duration(time.Millisecond * 5))
		//ctx.printInfo()
	}
ExitFlag:
	log.Println("client echo server out of service ....")
}

//will delete
func processPack(ctx *EchoContext, responseData NetworkData) {
	el := ctx.Queue.Front()
	var cmd = el.Value.(NetworkCmd)
	cmd.SetRspdArgs(responseData)
	ctx.pushRespondCmd(cmd)
	ctx.Queue.Remove(el)
}

const Cache_SIZE = 10000000

func NewEchoContext() *EchoContext {
	var ctx = EchoContext{}
	ctx.Queue = list.New()
	ctx.state = ECHO_CTX_STATE_RUNNING
	ctx.AskGoHead = make(chan interface{})
	ctx.ResponseTasks = make(chan NetworkCmd, Cache_SIZE)
	ctx.networkDataChan = make(chan NetworkData, Cache_SIZE)
	return &ctx
}

var selectDelayMillSecTimeout time.Duration = time.Millisecond * 300

func (this *EchoContext) PushNetworkData(netData NetworkData) {
	this.networkDataChan <- netData
	/*
		var t = time.NewTimer(selectDelayMillSecTimeout)
		defer t.Stop()
		for {
			t.Reset(selectDelayMillSecTimeout)
			select {
			case this.networkDataChan <- netData:
				break
			}
			select {
			case <-t.C:
				log.Println("PushNetworkData：timeout tick ...")
			default:
				log.Println("PushNetworkData：no block ...")
			}

		}*/

}

func (this *EchoContext) Init() int {
	if this.SendNetworkFn == nil {
		fmt.Println("this.SendNetworkFn==nil")
		return -1
	}
	return 0
}

func (ctx *EchoContext) ExecuteResponse() {
	var t = time.NewTimer(selectDelayMillSecTimeout)
	defer t.Stop()
	for {
		if (ctx.state == ECHO_CTX_STATE_CLOSE) {
			//存盘 持久化
			goto ExitFlag
		}
		t.Reset(selectDelayMillSecTimeout)
		//log.Println("ExecuteResponse  task len:", len(ctx.ResponseTasks), "ctx.Queue", ctx.Queue.Len(), " ctx.ResponseTasks", len(ctx.ResponseTasks))
		select {
		case task := <-ctx.ResponseTasks:
			if task != nil {
				task.DoResponse()
			}
		case <-t.C:
			//log.Println("ExecuteResponse tick  task len:", len(ctx.ResponseTasks), "ctx.Queue", ctx.Queue.Len(), " ctx.ResponseTasks", len(ctx.ResponseTasks))
		}
		t.Stop()
	}
ExitFlag:
	log.Println("exit from routine ExecuteResponse....")
}

/*
模拟 接收到终端数据
var endpointRecvQue = list.New()

type ResponseBin struct {
	Timestamp int64
}

func (this *ResponseBin) GetTimestamp() int64 {
	return this.Timestamp
}

func (this *ResponseBin) GetData() interface{} {
	return this
}

func AnswerRoute(ctx *EchoContext) {
	rand.Seed(time.Now().Unix())
	for ; endpointRecvQue.Len() > 0; {
		req := endpointRecvQue.Front().Value.(NetworkCmd)
		if req != nil {
			fmt.Println("endpoint 接收到了 请求:", req.GetTimestamp())
		}
		//延迟 3秒响应服务器
		var n = rand.Intn(5)
		time.Sleep(time.Duration(time.Second * (time.Duration(n))))
		ctx.NetworkDataChan <- &ResponseBin{Timestamp: time.Now().Unix()}

	}
}
*/

/*
模拟发起 网络请求
type iNetworkCmd struct {
	RespHandler func(args ...interface{}) interface{}
	RespArgs    []interface{}
	ReqPack     []byte
	Timestamp   int64
}
func (this *iNetworkCmd) GetTimestamp() int64 {
	return this.Timestamp
}

func (this *iNetworkCmd) SetTimestamp(t int64) {
	this.Timestamp = t
}

func (this *iNetworkCmd) SetRspdArgs(args ...interface{}) {
	this.RespArgs = args
}
func (this *iNetworkCmd) DoResponse() {
	if nil != this.RespHandler {
		this.RespHandler(this.RespArgs...)
	}
}

func (this *iNetworkCmd) GetSndPack() []byte {
	return nil
}*/

/*func AddRequest(ctx *EchoContext) {

	var t NetworkCmd
	for {
		t = &iNetworkCmd{
			RespHandler: func(args ...interface{}) interface{} {

				fmt.Printf("task result is:%T \n", args[0])
				return nil
			},
			RespArgs: []interface{}{1, 2, 3, 4},
		}
		t.Timestamp = time.Now().Unix()
		ctx.CommitRequest(t)
		time.Sleep(time.Duration(time.Second))

	}

}*/

/*
初始化上下文框架

func main() {

	var ctx = NewEchoContext()
	res := ctx.Init()
	if res == 0 {
		go AddRequest(ctx)
		go AskRoute(ctx)
		go AnswerRoute(ctx)
		go ExecuteResponse(ctx)
	} else {
		fmt.Println("初始化上下文失败....")
	}

	select {}

}
*/
