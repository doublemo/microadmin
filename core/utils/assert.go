package utils

import (
	"fmt"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-kit/kit/log"
)

// Assert 错误检查
// 如果发现错误那么直接panic
func Assert(err error) {
	if err != nil {
		panic(err)
	}
}

// RecoverStack 产生panic时的调用栈打印
func RecoverStack(log log.Logger, extras ...interface{}) {
	if r := recover(); r != nil {
		log.Log("panic", fmt.Sprint(r))
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			log.Log("frame", i, "func", runtime.FuncForPC(funcName).Name(), "file", file, "line", line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			log.Log("EXRAS", k, "DATA", fmt.Sprintf("%v", spew.Sdump(extras[k])))
		}
	}
}

// RecoverStackPanic 产生panic时的调用栈打印
func RecoverStackPanic(log log.Logger, extras ...interface{}) {
	if r := recover(); r != nil {
		log.Log("panic", fmt.Sprint(r))
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			log.Log("frame", i, "func", runtime.FuncForPC(funcName).Name(), "file", file, "line", line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			log.Log("EXRAS", k, "DATA", fmt.Sprintf("%v", spew.Sdump(extras[k])))
		}

		panic(r)
	}
}
