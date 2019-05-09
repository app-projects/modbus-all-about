package main

import (
	"fmt"
)

type A interface {
	Run()
}

type B struct {
	A `json:"-"`
	Name string
}

func (this *B) Run() {

}

func main() {
	var a byte = 100

	s := fmt.Sprintf("%d", a)
	fmt.Println([]byte(s))
	fmt.Println(s)
}
