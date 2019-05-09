package web

import (
	"net/http"
	"mux"
	"testModbus/utils"
	"render"
	"fmt"
	utils2 "modbusSvrRestful/utils"
	"net/source/proto/outputconfig"
	"strconv"
	"testModbus/outinterface"
	"net/source/userapi"
	"net/source/proto/binfiles"
	"sync"
	"testModbus/data"
	"log"
)

var r *render.Render
var sm *http.ServeMux
var router *mux.Router
var now nowTime

type stateResult struct {
	Res int    `json:"res"`
	Tip string `json:"tip"`
}

//var sr stateResult

var stateResPool = sync.Pool{
	New: func() interface{} {
		return &stateResult{}
	},
}

func init() {
	r = render.New(render.Options{
		Directory: "templates", // Specify what path to load the templates from.
		Asset: func(name string) ([]byte, error) { // Load from an Asset function instead of file.
			return []byte("template content"), nil
		},
		AssetNames: func() []string { // Return a list of asset names for the Asset function
			return []string{"filename.tmpl"}
		},
		Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template or {{ partial "css" }} to render a partial from the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		//Funcs: []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		Delims:                    render.Delims{"{[{", "}]}"},                      // Sets delimiters to the specified strings.
		Charset:                   "UTF-8",                                          // Sets encoding for content-types. Default is "UTF-8".
		DisableCharset:            true,                                             // Prevents the charset from being appended to the content type header.
		IndentJSON:                true,                                             // Output human readable JSON.
		IndentXML:                 true,                                             // Output human readable XML.
		PrefixJSON:                []byte(")]}',\n"),                                // Prefixes JSON responses with the given bytes.
		PrefixXML:                 []byte("<?xml version='1.0' encoding='UTF-8'?>"), // Prefixes XML responses with the given bytes.
		HTMLContentType:           "application/xhtml+xml",                          // Output XHTML content type instead of default "text/html".
		IsDevelopment:             true,                                             // Render will now recompile the templates on every HTML response.
		UnEscapeHTML:              true,                                             // Replace ensure '&<>' are output correctly (JSON only).
		StreamingJSON:             true,                                             // Streams the JSON response via json.Encoder.
		RequirePartials:           true,                                             // Return an error if a template is missing a partial used in a layout.
		DisableHTTPErrorRendering: true,                                             // Disables automatic rendering of http.StatusInternalServerError when an error occurs.
	})

	sm = http.NewServeMux()
	router = mux.NewRouter()
	initRestfulLogSvr()
}

func PostWebMsg(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	//msg_type, err1 := strconv.ParseUint(params["msg_type"], 10, 8)

	mac, err2 := strconv.ParseUint(params["mac"], 10, 8)
	offset_h, err3 := strconv.ParseUint(params["offset_h"], 10, 8)
	offset_l, err4 := strconv.ParseUint(params["offset_l"], 10, 8)
	data_h, err5 := strconv.ParseUint(params["data_h"], 10, 8)
	data_l, err6 := strconv.ParseUint(params["data_l"], 10, 8)

	var sr = stateResPool.Get().(*stateResult)

	if  err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		sr.Tip = "所有必须是整数"
		sr.Res = -1
		r.JSONP(w, http.StatusOK, "success", sr)
		return
	}

	msg := outinterface.NewWebMsg(true)
	msg.MsgType = outinterface.WEB_MSG_TYPE_MODIFY;
	msg.Mac = utils.Uint642Byte(mac)
	msg.OffsetH = utils.Uint642Byte(offset_h)
	msg.OffsetL = utils.Uint642Byte(offset_l)
	msg.DataH = utils.Uint642Byte(data_h)
	msg.DataL = utils.Uint642Byte(data_l)

	res := outinterface.PushMsg(msg)
	sr.Res = res
	if res == 0 {
		sr.Tip = "修改请求已提交"
	} else if res == -2 {
		sr.Tip = fmt.Sprintf("终端mac=%d没有建立连接", msg.Mac)
	} else if res == -3 {
		sr.Tip = "try again "
	} else if (res == -1) {
		sr.Tip = "msg is nil "
	}
	r.JSONP(w, http.StatusOK, "success", sr)

	stateResPool.Put(sr)
}

type ReadReg struct {
	Res   int    `json:"res"`
	Reg   uint16 `json:"reg"`
	Value uint16 `json:"value"`
	Timestamp int64 `json:"timestamp"`
}

var ReadRegPool = sync.Pool{
	New: func() interface{} {
		return &ReadReg{}
	},
}

func GetRegWebMsg(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	mac, err2 := strconv.ParseUint(params["mac"], 10, 8)
	offset_h, err3 := strconv.ParseUint(params["offset_h"], 10, 8)
	offset_l, err4 := strconv.ParseUint(params["offset_l"], 10, 8)
	data_h, err5 := strconv.ParseUint(params["data_h"], 10, 8)
	data_l, err6 := strconv.ParseUint(params["data_l"], 10, 8)
	if data_l == 0 {
		data_l = 1
	}

	var sr = stateResPool.Get().(*stateResult)

	if err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		sr.Tip = "所有必须是整数"
		sr.Res = -1
		r.JSONP(w, http.StatusOK, "success", sr)
		stateResPool.Put(sr)
		return
	}
	msg := outinterface.NewWebMsg(false)
	//这里都是读取WEB_MSG_TYPE_READ_REG 类型， 由于这里是异步（要先查询 终端 ，而这里是同步 不堵塞，所以客户端要请求两次获得数据）
	//opt here

	msg.MsgType = outinterface.WEB_MSG_TYPE_READ_REG;
	msg.Mac = utils.Uint642Byte(mac)
	msg.OffsetH = utils.Uint642Byte(offset_h)
	msg.OffsetL = utils.Uint642Byte(offset_l)

	msg.DataH = utils.Uint642Byte(data_h) //无用
	msg.DataL = utils.Uint642Byte(data_l) // 无用

	msg.Handler = func(args ...interface{}) interface{} { //请求响应 更新操作
		//成功的回调函数 ;封装好数据发送到客户端
		var b = args[0].(userapi.IModBusProtoBinPack)
		modbus03BinResp := b.(*binfiles.Mod03_ProtoBinPackResp)

		var result = make(map[string]string)
		result["res"] = "0"
		result["data_h"] = utils.Byte_2HexString(modbus03BinResp.DataFields[0])
		result["data_l"] = utils.Byte_2HexString(modbus03BinResp.DataFields[1])
		devData := data.GetDevDataContext().GetDevData(modbus03BinResp.Mac)
		if msg.OffsetH == 0x40 {
			devData.PutSysData(msg.OffsetH, msg.OffsetL, modbus03BinResp.DataFields[0], modbus03BinResp.DataFields[1])
		} else if msg.OffsetH == 0x20 {
			devData.PutAppData(msg.OffsetH, msg.OffsetL, modbus03BinResp.DataFields[0], modbus03BinResp.DataFields[1])
		}
		fmt.Printf("终端返回addr:0x %x%x 查询值 :data=%s%s:", msg.OffsetH, msg.OffsetL, result["data_h"], result["data_l"])

		//手动释放
		outinterface.ReleaseWebMsg(msg)
		return nil
	}

	//请求提交
	res := outinterface.PushMsg(msg)
	if res == -2 {
		sr.Tip = fmt.Sprintf("终端mac=%d没有建立连接", msg.Mac)
		r.JSONP(w, http.StatusOK, "success", sr)
		goto endflag
	} else if (res == -1) {
		sr.Tip = "msg is nil "
		r.JSONP(w, http.StatusOK, "success", sr)
		goto endflag
	} else {
		regRespd := ReadRegPool.Get().(*ReadReg)
		regRespd.Res = 0
		var dh, dl byte
		devData := data.GetDevDataContext().GetDevData(utils.Uint642Byte(mac))
		if msg.OffsetH == 0x40 {
			dh, dl = devData.GetSysData(utils.Uint642Byte(offset_h), utils.Uint642Byte(offset_l))

		} else if msg.OffsetH == 0x20 {
			dh, dl = devData.GetAppData(utils.Uint642Byte(offset_h), utils.Uint642Byte(offset_l))
		}
		regRespd.Reg = utils.Bytes2Uint16(msg.OffsetH, msg.OffsetL)
		regRespd.Value = utils.Bytes2Uint16(dh, dl)
		regRespd.Timestamp = utils.SvrNowTimestamp()
		r.JSONP(w, http.StatusOK, "success", regRespd)
		ReadRegPool.Put(regRespd)
	}

endflag:
	stateResPool.Put(sr)
}

type nowTime struct {
	SvrTime int64 `json:"svr_time"`
}

func GetNowTime(w http.ResponseWriter, req *http.Request) {
	now.SvrTime = utils.SvrNowTimestamp()
	r.JSONP(w, http.StatusOK, "success", now)
}

//restful 只提供 数据的对外api,不做网页显示相关的内容， 显示展示服务的功能，通过 web 资源服务器来提供
func initRestfulLogSvr() {
	router.HandleFunc("/dev/svr/msg/push/{mac}/{msg_type}/{offset_h}/{offset_l}/{data_h}/{data_l}", utils2.DecorateCrossMethod(PostWebMsg)).Methods("GET")
	router.HandleFunc("/nowtime", utils2.DecorateCrossMethod(GetNowTime)).Methods("GET")
	router.HandleFunc("/dev/msg/getreg/{mac}/{offset_h}/{offset_l}/{data_h}/{data_l}", utils2.DecorateCrossMethod(GetRegWebMsg)).Methods("GET")
}

func startListeners(address string) {
	router.Use(mux.CORSMethodMiddleware(router))
	//http.ListenAndServe(address, router)

	err := http.ListenAndServeTLS(address, "1671679_www.cloudchip.net.pem", "1671679_www.cloudchip.net.key", router)
	if err != nil {
		log.Fatal("ListenerAndServe:", err)
	}
}

func SvrWebRestfulMain() {
	address := fmt.Sprintf("%s:%d", outputconfig.RemotePushSvrMsgIp, outputconfig.RemotePushSvrMsgPort)
	startListeners(address)
	fmt.Println("center -srr msg listner :", address)
	//group.Done()

}
