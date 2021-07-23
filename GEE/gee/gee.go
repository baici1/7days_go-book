package gee

import (
	"net/http"
	"path"
	"strings"
	"text/template"
)

//包具体实现

//路由分组
//满足的条件
//1.前缀 ----分组的路径
//2.具有中间件
//3.可以进行嵌套
//4.提供分组的接口
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)
type Engine struct {
	//router map[string]HandlerFunc
	*RouterGroup
	router *router
	groups []*RouterGroup //存所有分组路由
	//模板
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

//初始化，创建Engine实例
func Default() *Engine {
	//	return &Engine{router: make(map[string]HandlerFunc)}
	//进行初始化
	//此时engine是最顶层的分组，它可以调用RouterGroup的所有接口
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//为分组创建一个engine
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//增加路由
//将请求方式，路径，函数都添加到Engine结构体
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	// key := method + "-" + pattern
	// engine.router[key] = handler

	//engine.router.addRoute(method, pattern, handler)

	pattern := group.prefix + comp
	//log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//请求的方法
//GET
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//POST
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
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
	//将路由组的中间件和函数 都放到了handlers
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		//匹配路由组
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//创建静态的映射关系
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath) //相对位置
	//找到静态文件的相对位置
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

//静态文件服务
func (group *RouterGroup) Static(relativePath string, root string) {
	hanler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, hanler)
}

//设置模板
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

//所有的自定义模板渲染函数
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
