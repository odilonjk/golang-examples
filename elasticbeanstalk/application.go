package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello World")
	http.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	_ = http.ListenAndServe(":5000", nil)
}
