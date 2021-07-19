package gee

import (
	"net/http"
)

//包具体实现

//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)
type Engine struct {
	//router map[string]HandlerFunc
	router *router
}

//初始化，创建Engine实例
func Default() *Engine {
	//	return &Engine{router: make(map[string]HandlerFunc)}
	return &Engine{router: newRouter()}
}

//增加路由
//将请求方式，路径，函数都添加到Engine结构体
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// key := method + "-" + pattern
	// engine.router[key] = handler
	engine.router.addRoute(method, pattern, handler)
}

//请求的方法
//GET
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

//POST
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

//启动
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// key := r.Method + "-" + r.URL.Path
	// handler, ok := engine.router[key]
	// if ok {
	// 	handler(w, r)
	// } else {
	// 	fmt.Printf("404 NOT FOUND: %s\n", r.URL)
	// }
	c := newContext(w, r)
	engine.router.handle(c)
}
