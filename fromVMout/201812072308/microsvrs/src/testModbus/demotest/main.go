package main

import (
	"fmt"
)

func trans(s string) {
	/* var bt = []byte(s)
	 unsafe.Pointer(&bt)*/

}

type data struct {
	S []int
}

func main() {

	var s = ""
	var sl = []byte(s)
	fmt.Println(sl[0])
}
