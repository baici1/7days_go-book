package main

import (
	"fmt"
	"gee"
	"net/http"
	"time"
)

// type Engine struct{}

// func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		fmt.Printf("Path%q", r.URL.Path)
// 	case "/hello":
// 		for k, v := range r.Header {
// 			fmt.Printf("Header[%q] = %q\n", k, v)
// 		}
// 	default:
// 		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
// 	}
// }
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
func main() {

	// http.HandleFunc("/", indexHandler)
	// http.HandleFunc("/hello", helloHandler)
	// log.Fatal(http.ListenAndServe(":9999", nil))

	// //创建实例
	// r := gee.Default()
	// //下面就是路由  参照gin框架
	// r.GET("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Printf("Path=%q\n", r.URL.Path)
	// })
	// r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	for k, v := range r.Header {
	// 		fmt.Printf("Header[%q]=%q\n", k, v)
	// 	}
	// })
	// //跑项目
	// r.Run(":9999")
	//Days2
	//封装request与response
	//创建实例
	// r := gee.Default()
	// //下面就是路由  参照gin框架
	// r.GET("/", func(c *gee.Context) {
	// 	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	// })
	// r.GET("/hello", func(c *gee.Context) {
	// 	// expect /hello?name=geektutu
	// 	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	// })

	// r.POST("/login", func(c *gee.Context) {
	// 	c.JSON(http.StatusOK, gee.H{
	// 		"username": c.PostForm("username"),
	// 		"password": c.PostForm("password"),
	// 	})
	// })
	// //跑项目
	// r.Run(":9999")
	//Days3
	//创建实例
	// 	r := gee.Default()
	// 	//下面就是路由  参照gin框架
	// 	r.GET("/", func(c *gee.Context) {
	// 		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	// 	})

	// 	r.GET("/hello", func(c *gee.Context) {
	// 		// expect /hello?name=geektutu
	// 		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	// 	})

	// 	r.GET("/hello/:name", func(c *gee.Context) {
	// 		// expect /hello/geektutu
	// 		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	// 	})

	// 	r.GET("/assets/*filepath", func(c *gee.Context) {
	// 		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	// 	})
	// 	//跑项目
	// 	r.Run(":9999")
	// }
	//Days4
	r := gee.Default()
	//下面就是路由  参照gin框架
	r.Use(gee.Logger())

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	// v1 := r.Group("/v1")
	// {

	// 	v1.GET("/hello", func(c *gee.Context) {
	// 		// expect /hello?name=geektutu
	// 		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	// 	})
	// }
	// v2 := r.Group("/v2")
	// {
	// 	v2.GET("/hello/:name", func(c *gee.Context) {
	// 		// expect /hello/geektutu
	// 		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	// 	})
	// 	v2.POST("/login", func(c *gee.Context) {
	// 		c.JSON(http.StatusOK, gee.H{
	// 			"username": c.PostForm("username"),
	// 			"password": c.PostForm("password"),
	// 		})
	// 	})

	//跑项目
	r.Run(":9999")
}

// // handler echoes r.URL.Path
// func indexHandler(w http.ResponseWriter, req *http.Request) {
// 	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
// }

// // handler echoes r.URL.Header
// func helloHandler(w http.ResponseWriter, req *http.Request) {
// 	for k, v := range req.Header {
// 		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
// 	}
// }
