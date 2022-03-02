package main

import (
	"log"
	"net/http"
)

type Engine struct {

}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/hello":
		writer.Write([]byte("访问/hello页面"))
		break
	default:
		writer.Write([]byte(request.URL.Path))
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe("127.0.0.1:8888", engine))
}
