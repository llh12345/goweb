package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//origin
	Writer http.ResponseWriter
	Req    *http.Request
	//
	Path   string
	Method string
	//
	StatusCode int

	Params map[string]string

	handlers []HandleFunc
	index    int
}

func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}
func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: writer,
		Req:    req,
		Path:   req.RequestURI,
		Method: req.Method,
		index:  -1,
	}
}
func (c *Context) Json(statusCode int, data interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(statusCode)
	err := json.NewEncoder(c.Writer).Encode(data)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) PostForm(key string) (value string) {
	return c.Req.FormValue(key)
}
func (c *Context) Query(key string) (value string) {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	fmt.Fprintf(c.Writer, format, values...)
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
