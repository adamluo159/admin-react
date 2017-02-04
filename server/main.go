package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("server begin")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
	})
	http.ListenAndServe(":3030", nil)
}
