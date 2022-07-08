package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(ctx *Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func NewEngine() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine:      engine,
		prefix:      "",
		middlewares: make([]HandleFunc, 0),
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func (g *RouterGroup) Use(handleFuncs ...HandleFunc) {
	g.middlewares = append(g.middlewares, handleFuncs...)
}
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req)
	handlers := make([]HandleFunc, 0)
	if req.RequestURI == "/" {
		handlers = append(handlers, e.RouterGroup.middlewares...)
	} else {
		for _, group := range e.groups {
			if strings.HasPrefix(req.RequestURI, group.prefix) {
				handlers = append(handlers, group.middlewares...)
			}
		}
	}
	ctx.handlers = handlers
	e.router.handle(ctx)
}
func (e *Engine) addRoute(method string, uriPath string, handleFunc HandleFunc) {
	e.router.addRoute(method, uriPath, handleFunc)
}

type RouterGroup struct {
	prefix      string
	engine      *Engine
	middlewares []HandleFunc
	parent      *RouterGroup
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		parent:      group,
		engine:      engine,
		middlewares: make([]HandleFunc, 0),
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRoute("GET", pattern, handler)
}
func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRoute("POST", pattern, handler)
}
