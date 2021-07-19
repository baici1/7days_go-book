# 序言

此项目是模仿的是[七天用Go从零实现系列](https://geektutu.com/post/gee.html) 

这个系列让我学习到了很多！！感谢博主

# Gee框架

语雀文档：https://www.yuque.com/docs/share/a3b7cc1a-12aa-4b42-bd6f-5e98ae8c080b?# 《🐹Gee框架》

## 序言

这个教程将使用 Go 语言实现一个简单的 Web 框架，起名叫做`Gee`，[`geektutu.com`](https://geektutu.com/)的前三个字母。我第一次接触的 Go 语言的 Web 框架是`Gin`，`Gin`的代码总共是14K，其中测试代码9K，也就是说实际代码量只有**5K**。`Gin`也是我非常喜欢的一个框架，与Python中的`Flask`很像，小而美。

`7天实现Gee框架`这个教程的很多设计，包括源码，参考了`Gin`，大家可以看到很多`Gin`的影子。

时间关系，同时为了尽可能地简洁明了，这个框架中的很多部分实现的功能都很简单，但是尽可能地体现一个框架核心的设计原则。例如`Router`的设计，虽然支持的动态路由规则有限，但为了性能考虑匹配算法是用`Trie树`实现的，`Router`最重要的指标之一便是性能。

## Day1 HTTP基础

go中是内置了 `http` 库的。最原生的写web应用其实就是用的是 `http` 库

```
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/hello", helloHandler)
    log.Fatal(http.ListenAndServe(":9999", nil))
}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
    for k, v := range req.Header {
        fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
    }
}
```

上面是最原生的写web应用的写法

`http.HandleFunc("/", indexHandler)` 这句话第一个参数是路由地址，第二个参数是绑定的逻辑函数。

`http.ListenAndServe(":9999", nil)` 是用来启动 Web 服务的，第一个参数是地址，`:9999`表示在 *9999* 端口监听。而第二个参数则代表处理所有的HTTP请求的实例，`nil` 代表使用标准库中的实例处理。第二个参数，则是我们基于`net/http`标准库实现Web框架的入口。

我们可以看看源码

![image.png](https://cdn.nlark.com/yuque/0/2021/png/21455688/1626602018996-f9e8c55d-cb8e-47b5-8f74-828fb67c5d87.png)

第二个参数就是 `Handler` 类型

![image.png](https://cdn.nlark.com/yuque/0/2021/png/21455688/1626602044737-c9f3d114-50a6-450b-a3e0-fa9a82a67425.png)

实现这样一个 `ServeHTTP` 接口



其实你也会发现你写的路由绑定的函数也是在变相的写 `ServeHTTP` 接口，因为参数都是一致的。

那么 `http`库根据web应用的原理可能是`http.ListenAndServe(":9999", nil)` 起的是web应用的入口，当你进行请求的时候，他会匹配路由，做对应对应的逻辑函数直到应用关闭,这一系列的工作都是在一个 `Handler` 实例完成，虽然此时传入的是 `nil` 但是你写的路由都挂载到了![image.png](https://cdn.nlark.com/yuque/0/2021/png/21455688/1626602791649-466c8e8c-bc76-4504-bc57-4b07f9b455d6.png)这个实例当中

根据上面的原理，我们就可以自己简单封装 `http` 库

```
package main

import (
    "fmt"
    "log"
    "net/http"
)

// Engine is the uni handler for all requests
type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    switch req.URL.Path {
    case "/":
        fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
    case "/hello":
        for k, v := range req.Header {
            fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
        }
    default:
        fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
    }
}

func main() {
    engine := new(Engine)
    log.Fatal(http.ListenAndServe(":9999", engine))
}
```

- 我们定义了一个结构体 `Engine` ，实现 `ServeHTTP` 的接口了，里面的逻辑就是在匹配路由，实现不同的逻辑工作。
- `http.ListenAndServe(":9999", nil)` 的第二个参数传入创建的 `Engine` 的实例。
- 当你利用 `curl` 工具进行测试时，你会发现，结果符合上诉猜测过程。



我们开始模仿gin框架，进行封装 `http` 库

### 1.导入本地包

地址：https://www.liwenzhou.com/posts/Go/import_local_package_in_go_module/

- 创建一个gee的文件夹，然后 `go mod init gee` 生成 `go.mod` 
- 在 `demo` 的 `go.mod` 中写入

```
module one

go 1.16

require gee v0.0.0

replace gee => ../gee
```

- 在主文件中直接引入 `gee` 
- 在gee文件内创建文件 `gee.go` 写入 `package gee` 就不会有报错了 

### 2.本地文件代码

```
package main

import (
    "fmt"
    "gee"
    "net/http"
)

// type Engine struct{}

// func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//  switch r.URL.Path {
//  case "/":
//      fmt.Printf("Path%q", r.URL.Path)
//  case "/hello":
//      for k, v := range r.Header {
//          fmt.Printf("Header[%q] = %q\n", k, v)
//      }
//  default:
//      fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
//  }
// }

func main() {

    // http.HandleFunc("/", indexHandler)
    // http.HandleFunc("/hello", helloHandler)
    // log.Fatal(http.ListenAndServe(":9999", nil))

    //创建实例
    r := gee.Default()
    //下面就是路由  参照gin框架
    r.GET("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Path=%q\n", r.URL.Path)
    })
    r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
        for k, v := range r.Header {
            fmt.Printf("Header[%q]=%q\n", k, v)
        }
    })
    //跑项目
    r.Run(":9999")
}

// // handler echoes r.URL.Path
// func indexHandler(w http.ResponseWriter, req *http.Request) {
//  fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
// }

// // handler echoes r.URL.Header
// func helloHandler(w http.ResponseWriter, req *http.Request) {
//  for k, v := range req.Header {
//      fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
//  }
// }
```

### 3.gee包的基本架构

```
package gee

import (
    "fmt"
    "net/http"
)

//包具体实现

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
    router map[string]HandlerFunc
}

//初始化，创建Engine实例
func Default() *Engine {
    return &Engine{router: make(map[string]HandlerFunc)}
}

//增加路由
//将请求方式，路径，函数都添加到Engine结构体
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
    key := method + "-" + pattern
    engine.router[key] = handler
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
    key := r.Method + "-" + r.URL.Path
    handler, ok := engine.router[key]
    if ok {
        handler(w, r)
    } else {
        fmt.Printf("404 NOT FOUND: %s\n", r.URL)
    }
}
```

整个`Gee`框架的原型已经出来了。实现了路由映射表，提供了用户注册静态路由的方法，包装了启动服务的函数。

以上就是Gee框架的雏形

完成时间：2021/7/18

## Days2 设计上下文Context

我们第一天完成的是一个静态路由的雏形，但是返回消息和请求的方式都比较简单，不能达到一个处理业务的能力，并且都不够完整，如果我们要构造一个完整的响应和请求，那么我们就得去设置请求头和主体，但是有些部分我们要不断地去设置例如状态码，消息类型等，这些重复工作我们可以进行封装，达到一个能快速，完整的构造HTTP响应的能力。

用返回 JSON 数据作比较，感受下封装前后的差距。

封装前

```
obj = map[string]interface{}{
    "name": "geektutu",
    "password": "1234",
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
encoder := json.NewEncoder(w)
if err := encoder.Encode(obj); err != nil {
    http.Error(w, err.Error(), 500)
}
```

封装后

```
c.JSON(http.StatusOK, gee.H{
    "username": c.PostForm("username"),
    "password": c.PostForm("password"),
})
```

当你利用gin框架写项目时候，你要处理不同的返回体，就会利用不同的接口，而且都比较简单，也需要去处理动态路由，以及一些中间件，这么多的工作下，函数的参数都没有发生改变都是 `*gin.Context` ,可以判断出上诉那些工作和信息都由Context去承载了。设计 Context 结构，扩展性和复杂性留在了内部，而对外简化了接口。路由的处理函数，以及将要实现的中间件，参数都统一使用 Context 实例， Context 就像一次会话的百宝箱，可以找到任何东西。



所以我们要对Context进行设计，封装。

- 返回体

- - 请求头
  - 主体
  - 。。。。

- 请求体

- - 动态路由
  - 。。。



```
package gee

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type H map[string]interface{}

type Context struct {
    //起始结构
    Writer http.ResponseWriter
    Req    *http.Request
    //请求信息
    Path   string
    Method string
    //返回信息
    StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
    return &Context{
        Writer: w,
        Req:    req,
        Path:   req.URL.Path,
        Method: req.Method,
    }
}

//获取post信息
func (c *Context) PostForm(key string) string {
    return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
    return c.Req.URL.Query().Get(key)
}

//设置状态码
func (c *Context) Status(code int) {
    c.StatusCode = code
    c.Writer.WriteHeader(code)
}

//设置头部信息
func (c *Context) SetHeader(key string, value string) {
    c.Writer.Header().Set(key, value)
}

//处理返回类型
//TEXT
func (c *Context) String(code int, format string, values ...interface{}) {
    c.SetHeader("Content-Type", "text/plain")
    c.Status(code)
    c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//JSON
func (c *Context) JSON(code int, obj interface{}) {
    c.SetHeader("Content-Type", "application/json")
    c.Status(code)
    encoder := json.NewEncoder(c.Writer)
    if err := encoder.Encode(obj); err != nil {
        http.Error(c.Writer, err.Error(), 500)
    }
}

//Data
func (c *Context) Data(code int, data []byte) {
    c.Status(code)
    c.Writer.Write(data)
}

//HTML
func (c *Context) HTML(code int, html string) {
    c.SetHeader("Content-Type", "text/html")
    c.Status(code)
    c.Writer.Write([]byte(html))
}
```

- 代码最开头，给`map[string]interface{}`起了一个别名`gee.H`，构建JSON数据时，显得更简洁。
- `Context`目前只包含了`http.ResponseWriter`和`*http.Request`，另外提供了对 Method 和 Path 这两个常用属性的直接访问。
- 提供了访问Query和PostForm参数的方法。
- 提供了快速构造String/Data/JSON/HTML响应的方法。



第一天的代码，路由仅仅完成了查找和绑定函数的在作用，还有其他的功能并未写上，所以为了解耦以及增强路由功能，简化代码，我们将路由方法和结构提取出来。方便后面对路由进行加强

```
type router struct {
    handlers map[string]HandlerFunc
}

func newRouter() *router {
    return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
    log.Printf("Route %4s - %s", method, pattern)
    key := method + "-" + pattern
    r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
    key := c.Method + "-" + c.Path
    if handler, ok := r.handlers[key]; ok {
        handler(c)
    } else {
        c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
    }
}
```

我们增加了上下文Context以及提取了路由，那么主文件我们也需要进行改变。

```
// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
    router *router
}

// New is the constructor of gee.Engine
func New() *Engine {
    return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
    engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
    engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
    engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
    return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    c := newContext(w, req)
    engine.router.handle(c)
}
```

将`router`相关的代码独立后，`gee.go`简单了不少。最重要的还是通过实现了 ServeHTTP 接口，接管了所有的 HTTP 请求。相比第一天的代码，这个方法也有细微的调整，在调用 router.handle 之前，构造了一个 Context 对象。这个对象目前还非常简单，仅仅是包装了原来的两个参数，之后我们会慢慢地给Context插上翅膀。

时间：2021/7/19
