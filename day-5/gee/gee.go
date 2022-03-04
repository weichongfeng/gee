package gee

import (
	"net/http"
	"strings"
)

type HandleFunc func(*Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix string
	middlewares []HandleFunc
	parent *RouterGroup
	engine *Engine
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := NewContext(writer, request)
	var middleware  []HandleFunc
	for _,routerGroup := range e.groups{
		if strings.HasPrefix(request.URL.Path,routerGroup.prefix) {
			middleware = routerGroup.middlewares
		}
	}
	c.handlers = middleware

	e.router.handle(c)
	
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (g *RouterGroup) AddRouter(method string, comp string, handler HandleFunc)  {
	pattern := g.prefix + comp
	g.engine.router.addRouter(method, pattern, handler)
}

func (g *RouterGroup) Get(pattern string, handler HandleFunc) {
	g.AddRouter("GET", pattern, handler)
}

func (g *RouterGroup) Post(pattern string, handler HandleFunc) {
	g.AddRouter("POST", pattern, handler)
}

func (g *RouterGroup) Use(handler ...HandleFunc)  {
	g.middlewares = append(g.middlewares, handler...)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr,e)
}

