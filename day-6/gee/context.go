package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request
	Method string
	Path string
	Params map[string]string
	StatusCode int
	handlers []HandleFunc
	index int
	engine *Engine
}

func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Writer: writer,
		Request: request,
		Method: request.Method,
		Path: request.URL.Path,
		index: -1,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string)  {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{})  {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{})  {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte)  {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) Fail(code int, err string)  {
	c.Status(code)
	c.Writer.Write([]byte(err))
}

func (c *Context) HTML(code int, name string, data interface{})  {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer,name,data); err != nil{
		c.Fail(500, err.Error())
	}
}

func (c *Context) Next()  {
	c.index++
	len := len(c.handlers)
	for ; c.index < len; c.index++ {
		c.handlers[c.index](c)
	}
}