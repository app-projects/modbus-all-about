package utils

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

//kvData map[string]interface{}
//struct 对象
func HttpPost(url string, kvData interface{}, headFields map[string]string) {

	bytesData, err := json.Marshal(kvData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
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
	fmt.Println("HttpPost has send...")

}
