package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.Get("/", func(ctx *gee.Context) {
		fmt.Println("1111")
		ctx.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.Run("127.0.0.1:8888")
}
