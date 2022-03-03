package gee

import (
	"net/http"
)

type HandleFunc func(*Context)

type Engine struct {
	router *router
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := NewContext(writer, request)
	e.router.handle(c)
	
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (e *Engine) AddRouter(method string, pattern string, handler HandleFunc)  {
	e.router.addRouter(method, pattern, handler)
}

func (e *Engine) Get(pattern string, handler HandleFunc) {
	e.AddRouter("GET", pattern, handler)
}

func (e *Engine) Post(pattern string, handler HandleFunc) {
	e.AddRouter("POST", pattern, handler)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr,e)
}

