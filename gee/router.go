package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandleFunc
	roots    map[string]*node
}

func newRouter() *router {
	roots := make(map[string]*node)
	roots["GET"] = &node{}
	roots["POST"] = &node{}
	return &router{
		roots: roots,
	}
}
func (r *router) addRoute(method string, pattern string, handleFunc HandleFunc) {
	root := r.roots[method]
	parts := r.parsePattern(pattern)
	root.insert(pattern, parts, 0, handleFunc)
	//key := method + "-" + pattern
	//r.handlers[key] = handleFunc
}

func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	root := r.roots[method]
	searchParts := r.parsePattern(pattern)
	n := root.search(pattern, searchParts, 0)
	if n != nil {
		routeParts := r.parsePattern(n.pattern)
		params := make(map[string]string)
		for index, part := range routeParts {
			if len(part) > 0 && part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
		}
		return n, params
	}
	return nil, nil
}
func (r *router) handle(c *Context) {
	log.Printf("request method:%s, uri: %s\n", c.Method, c.Path)
	n, params := r.getRoute(c.Method, c.Path)
	if n == nil || n.f == nil {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		return
	}
	c.Params = params
	c.handlers = append(c.handlers, n.f)
	c.Next()
}
func (r *router) parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")
	return parts[1:]
}
