package main

import (
	"net/http"
	"fmt"
	"log"
)

func test1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is a test1 page server for you")
}
func main() {
	http.HandleFunc("/test1", test1)

	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal("ListenerAndServe:", err)
	}
}