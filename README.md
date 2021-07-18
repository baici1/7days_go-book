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

