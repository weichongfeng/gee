package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

type HandleFunc func(*Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
	htmlTemplates *template.Template
	funcMap template.FuncMap
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

func Default() *Engine {
	engine := New()
	engine.Use(Recovery())
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

func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file);err != nil{
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (g *RouterGroup) Static(relativePath string, root string)  {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	g.Get(urlPattern, handler)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap)  {
	e.funcMap = funcMap
}

func (e *Engine) loadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr,e)
}

