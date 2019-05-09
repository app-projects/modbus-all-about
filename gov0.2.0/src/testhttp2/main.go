package main

import (
	"net/http"
	"fmt"
	"log"
)

func test2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is a test2 page server  for you")
}
func main() {
	http.HandleFunc("/test2", test2)

	err := http.ListenAndServe(":8002", nil)
	if err != nil {
		log.Fatal("ListenerAndServe:", err)
	}
}