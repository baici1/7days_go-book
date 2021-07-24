package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr
	//Callers 用来返回调用栈的程序计数器,
	//第 0 个 Caller 是 Callers 本身，
	//第 1 个是上一层 trace，
	//第 2 个是再上一层的 defer func。
	//因此，为了日志简洁一点，我们跳过了前 3 个 Caller。
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder //一个string类型的数值写入builder
	//Builder的底层实现其实就是一个string类型的切片
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	//当执行WriteString操作时，实际上就是append操作，最后利用String（）函数将他拼接成一个字符串
	return str.String()
}

//错误处理机制
//可以避免因为 panic 发生而导致整个程序终止，recover 函数只在 defer 中生效。
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message)) //获取触发 panic 的堆栈信息
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}
