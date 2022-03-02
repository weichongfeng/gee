package gee

import (
	"fmt"
	"net/http"
)

type HandleFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandleFunc
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := request.Method + "-" + request.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(writer, request)
	}else {
		fmt.Fprintf(writer, "404 NOT FOUND: %s\n", request.URL)
	}
	
}

func New() *Engine {
	return &Engine{router: make(map[string]HandleFunc)}
}

func (e *Engine) AddRouter(method string, pattern string, handle HandleFunc)  {
	key := method + "-" + pattern
	e.router[key] = handle
}

func (e *Engine) Get(pattern string, handle HandleFunc) {
	e.AddRouter("GET", pattern, handle)
}

func (e *Engine) Post(pattern string, handle HandleFunc) {
	e.AddRouter("POST", pattern, handle)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr,e)
}

