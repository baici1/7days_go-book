package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node //建立一个前缀树路由 去映射handler
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc)}
}

//作用主要是分割路由地址（以/分割成各个部分）
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}

	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method] //可以添加自定义的请求方式
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0) //注册路由
	r.handlers[key] = handler
}

//获取路由以及匹配字段（param）值
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) //拿到请求地址每个部分
	params := make(map[string]string)
	root, ok := r.roots[method] //从请求方式作为根节点找
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0) //匹配 如果匹配成功就会返回子节点（最小的）
	if n != nil {
		parts := parsePattern(n.pattern) //拿到子节点路由的每个部分
		for index, part := range parts { //请求地址和路由地址进行匹配（：和*）
			if part[0] == ':' {
				params[part[1:]] = searchParts[index] //将请求地址的值与路由地址的params进行映射
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

//当有了中间件，就不再直接执行了而是通过添加到handlers里面了
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path) //获取请求地址和params值
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
		//	r.handlers[key](c) //映射对应的handler
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})

	}
	c.Next()
}
