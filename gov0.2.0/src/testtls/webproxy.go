package main

/*******
 web proxy transfer server

*****/

import (
	"net/http"
	"mux"
	"render"
	"fmt"
	"log"
	"io"
	"modbusSvrRestful/utils"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"os"
	"strconv"
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

func HttpDo(url string, kvData interface{}, headFields map[string]string, method string) []byte {

	bytesData, err := json.Marshal(kvData)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if headFields != nil {
		for k, v := range headFields {
			request.Header.Set(k, v)
		}
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error(), b)
		return nil
	}
	//byte数组直接转成string，优化内存
	//str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("HttpPost has send...", string(b))
	return b

}

func forward(w http.ResponseWriter, req *http.Request, remote_addr string) {
	reqUrl := remote_addr + req.URL.Path
	log.Println("reqUrl:", reqUrl, " Method:", req.Method)
	b := HttpDo(reqUrl, nil, nil, req.Method)
	log.Println(string(b))


	io.WriteString(w, string(b))
}

//http://47.110.78.124:8520/dev/svr/msg/push
var proxy_msg = "http://47.110.78.124:8520"

func apiMsgSvr(w http.ResponseWriter, r *http.Request) {
	fmt.Println("apimsg prox uri:", r.RequestURI)
	utils.Fixcross(w)
	forward(w, r, proxy_msg)
}

//http://47.110.78.124:8501/dev/log/modifyquery"

var prox_data_store_view = "http://47.110.78.124:8501"

func dataview(w http.ResponseWriter, r *http.Request) {
	fmt.Println("dataview prox uri:", r.RequestURI)
	utils.Fixcross(w)
	forward(w, r, prox_data_store_view)
}

func initRestfulLogSvr() {

	//msg svr
	router.HandleFunc("/dev/svr/msg/push/{mac}/{msg_type}/{offset_h}/{offset_l}/{data_h}/{data_l}", apiMsgSvr).Methods("GET")
	router.HandleFunc("/nowtime", apiMsgSvr).Methods("GET")
	router.HandleFunc("/dev/msg/getreg/{mac}/{offset_h}/{offset_l}/{data_h}/{data_l}", apiMsgSvr).Methods("GET")

	//api proxy for datastore svr
	//http://47.110.78.124:8501/dev/log/modifyquery

	router.HandleFunc("/dev/log", dataview).Methods("POST")
	router.HandleFunc("/dev/log/{mac}", dataview).Methods("DELETE")
	router.HandleFunc("/dev/log/modify", dataview).Methods("POST")
	router.HandleFunc("/dev/log/modifyquery/{mac}", dataview).Methods("GET")
	//http://192.168.1.108:12345/dev/log/modifyquery/2
	router.HandleFunc("/dev/log/{mac}/{lastDataTime}", dataview).Methods("GET")
	router.HandleFunc("/dev/log/{mac}", dataview).Methods("GET")
	//http://192.168.1.108:12345/dev/log/2/0
	router.HandleFunc("/dev/log/stat/exist/{mac}", dataview).Methods("GET")
	router.HandleFunc("/nowtime", dataview).Methods("GET")

}

func startListeners(address string) {
	router.Use(mux.CORSMethodMiddleware(router))
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal("ListenerAndServe:", err)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("please input format : ./proxyapp.exe ip port")
		return
	}
	ip := os.Args[1]
	basePort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("baseport is a int num >0")
	}
	var address = fmt.Sprintf("%s:%d", ip, basePort)
	log.Println("will start web proxy server....", address)
	startListeners(address)
}
