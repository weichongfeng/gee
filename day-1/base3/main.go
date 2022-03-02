package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	gee := gee.New()
	gee.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
	})
	
	gee.Get("/hello", func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
		}
	})

	gee.Run("127.0.0.1:8888")
}
