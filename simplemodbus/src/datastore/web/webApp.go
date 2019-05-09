package main

import (
	"net/http"
	"mux"
	"encoding/json"
	"testModbus/utils"
	"time"
	"datastore/web/store"
	"datastore/web/pages"
	"render"
	"math/rand"
	"fmt"
	"strconv"
	"datastore/web/store/entity"
	"os"
)

var r *render.Render
var sm *http.ServeMux
var router *mux.Router

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

type err_NotTarget struct {
	Msg string `json:"msg"`
}

var notFindmsg = err_NotTarget{"no mac dev "}

type outlogfmt struct {
	List    []store.Devlog `json:"list"`
	SvrTime int64          `json:"svrtime"`
}

var outlogF = outlogfmt{}
var recentMaxLogNum = 20

func GetDevLog(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	mac := params["mac"]
	var svrTime = utils.SvrNowTimestamp()
	//时间的合法检查
	var lastDataTimeStr string
	var lastDataTime int64 = 0
	lastDataTimeStr = params["lastDataTime"]
	if lastDataTimeStr == "" || lastDataTimeStr == "0" {
		lastDataTime = -1
	} else {
		_, ok := strconv.ParseInt(lastDataTimeStr, 10, 0)
		if ok != nil {
			lastDataTime = svrTime
		}
	}
	//时间的合法检查
	var loglist = make([]store.Devlog, 0)
	if !store.GetStoreInstance().QueryExist(mac) {
		//log.Fatal("not exist mac dev id:%s\n", mac)
		r.JSONP(w, http.StatusNotFound, "failed", notFindmsg)
		return
	}

	store.GetStoreInstance().LoadLogRecentByLowWaterTimeLine(mac, func(log store.Devlog) int {
		loglist = append(loglist, log)
		return utils.STATUS_LOOP_OK
	}, svrTime, recentMaxLogNum, lastDataTime)

	outlogF.List = loglist
	outlogF.SvrTime = svrTime
	r.JSONP(w, http.StatusOK, "success", outlogF)
}

func PostDevLog(w http.ResponseWriter, req *http.Request) {
	var log store.Devlog = entity.NewDevQueryLog()
	_ = json.NewDecoder(req.Body).Decode(log)
	store.GetStoreInstance().PushLog(log)
	json.NewEncoder(w).Encode(log)
}

func statFun(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	pages.RenderPage(res, "F:/gopaths/simplemodbus/src/datastore/web/pages/rtdata.html", params)
}

func GetDevModifyLog(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)
	mac := params["mac"]
	var svrTime = utils.SvrNowTimestamp()

	var loglist = make([]store.Devlog, 0)
	if !store.GetModifyStoreInstance().QueryExist(mac) {
		r.JSONP(w, http.StatusNotFound, "failed", notFindmsg)
		return
	}

	store.GetModifyStoreInstance().LoadLogRecentByLowWaterTimeLine(mac, func(log store.Devlog) int {
		loglist = append(loglist, log)
		return utils.STATUS_LOOP_OK
	}, -1, recentMaxLogNum, -1)

	outlogF.List = loglist
	outlogF.SvrTime = svrTime
	r.JSONP(w, http.StatusOK, "success", outlogF)

}

func PostDevModifyLog(w http.ResponseWriter, req *http.Request) {
	var log store.Devlog = entity.NewDevModifyLog()
	_ = json.NewDecoder(req.Body).Decode(log)
	store.GetModifyStoreInstance().PushLog(log)
	json.NewEncoder(w).Encode(log)
}

func initRestfulLogSvr() {

	router.HandleFunc("/dev/log", PostDevLog).Methods("POST")
	router.HandleFunc("/dev/log/{mac}", DeleteDevLog).Methods("DELETE")
	router.HandleFunc("/dev/log/modify", PostDevModifyLog).Methods("POST")

	router.HandleFunc("/dev/log/modifyquery/{mac}", GetDevModifyLog).Methods("GET")
	//http://192.168.1.108:12345/dev/log/modifyquery/2
	router.HandleFunc("/dev/log/{mac}/{lastDataTime}", GetDevLog).Methods("GET")
	//http://192.168.1.108:12345/dev/log/2/0

	router.HandleFunc("/dev/stat/{mac}", statFun)
}

func CommitLogSvr() {
	rand.Seed(time.Now().Unix())
	for {
		time.Sleep(time.Second * 2)
		fmt.Println("push push ....mills:", utils.TransTime2MillSec(time.Now()))

		/*		var log = entity.NewDevQueryLog()

				store.GetStoreInstance().PushLog(&store.Devlog{"1", 3, "192.16.1.123", utils.TransTime2MillSec(time.Now()), rand.Intn(30)})
				store.GetStoreInstance().PushLog(&store.Devlog{"2", 6, "192.16.1.125", utils.TransTime2MillSec(time.Now()), rand.Intn(40)})*/
	}
}

func startListeners(address string) {
	go http.ListenAndServe(address, router)
	//go CommitLogSvr()
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("please input format : ./app.exe ip port")
		return
	}
	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("baseport is a int num >0")
		return
	}
    address:=fmt.Sprintf("%s:%d",ip,port)
    fmt.Println("will listner :",address)
	startListeners(address)
	select {}
}

func DeleteDevLog(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)
	//mac:=params["mac"]
	//json.NewEncoder(w).Encode(people)
}

func initTmpltRender() {

	/*	sm.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("Welcome, visit sub pages now."))
		})

		sm.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
			r.Data(w, http.StatusOK, []byte("Some binary data here."))
		})

		sm.HandleFunc("/text", func(w http.ResponseWriter, req *http.Request) {
			r.Text(w, http.StatusOK, "Plain text here")
		})

		sm.HandleFunc("/json", func(w http.ResponseWriter, req *http.Request) {
			r.JSON(w, http.StatusOK, map[string]string{"hello": "json"})
		})

		sm.HandleFunc("/jsonp", func(w http.ResponseWriter, req *http.Request) {
			r.JSONP(w, http.StatusOK, "success", map[string]string{"hello": "jsonp"})
		})

		sm.HandleFunc("/jsonp", func(w http.ResponseWriter, req *http.Request) {
			r.JSONP(w, http.StatusOK, "success", map[string]string{"hello": "jsonp"})
		})

		sm.HandleFunc("/html", func(w http.ResponseWriter, req *http.Request) {
			// Assumes you have a template in ./templates called "example.tmpl"
			// $ mkdir -p templates && echo "<h1>Hello {{.}}.</h1>" > templates/example.tmpl
			r.HTML(w, http.StatusOK, "example", "World")
		})*/

}
