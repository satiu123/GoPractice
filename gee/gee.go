package gee

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup //存储所有分组
}
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc //中间件
	parent      *RouterGroup  //
	engine      *Engine       //所有分组共享同一个Engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: group.engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.addRoute("PUT", pattern, handler)
}
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
