package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe("127.0.0.1:8888", nil))
}

func hello(w http.ResponseWriter, r *http.Request)  {
	path := r.URL.Path
	w.Write([]byte(path))
}

func index(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("hello world!"))
}
