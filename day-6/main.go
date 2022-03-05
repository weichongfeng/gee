package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.Default()

	r.Get("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	r.Get("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})
	r.Run("127.0.0.1:8888")
}
