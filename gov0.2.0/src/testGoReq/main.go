package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
)

func HttpPost(url string, kvData interface{}, headFields map[string]string) {

	bytesData, err := json.Marshal(kvData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("GET", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
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
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error(), b)
		return
	}
	//byte数组直接转成string，优化内存
	//str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("HttpPost has send...",string(b))

}

func main() {
	var url = "http://47.110.78.124:8501/dev/log/modifyquery/1"
	HttpPost(url, nil, nil)
}
