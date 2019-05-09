package utils

import "net/http"

func Fixcross(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
}

type handler = func(http.ResponseWriter, *http.Request)

func DecorateCrossMethod(fn handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		Fixcross(w)
		fn(w, req)
	}

}
