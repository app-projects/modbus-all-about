package main

import (
	"strings"
	"fmt"
)

func main() {
	var str="0 1 0 50 0 50 0 12 0 2 0 0 0 5 0 50 0 60 255 230 0 40 0 9 0 0 0 45 0 45 255 254 0 8 0 13 0 5 0 13 0 8 0 60 0 3 0 100 0 50 0 10 0 30 0 20 0 4 0 7 0 12 0 10"
	slstr:=strings.Split(str," ")
	fmt.Println(len(slstr))
}
