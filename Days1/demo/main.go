package main

import (
	"fmt"
	"gee"
	"net/http"
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
// 	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
// }

// // handler echoes r.URL.Header
// func helloHandler(w http.ResponseWriter, req *http.Request) {
// 	for k, v := range req.Header {
// 		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
// 	}
// }