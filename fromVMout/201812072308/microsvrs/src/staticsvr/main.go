package main

import (
	"net/http"
	"os"
	"fmt"
	"strconv"
)

// staticweb
// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

var statidir = "D:\\webstore\\modbusstate"

func main() {

	if len(os.Args) < 4 {
		fmt.Println("please input format : ./app.exe ip port  absoulute-web-home-dir")
		return
	}
	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])

	if err != nil {
		fmt.Println("baseport is a int num >0")
		return
	}
	statidir = os.Args[3]

	ex, e := PathExists(statidir)
	if ex {
		address := fmt.Sprintf("%s:%d", ip, port)
		fmt.Println("will listner :", address)
		fmt.Println("dir  :", statidir)
		http.Handle("/", http.FileServer(http.Dir(statidir)))
		http.ListenAndServe(address, nil)
	} else {
		fmt.Println("err:", e.Error())
	}

}
