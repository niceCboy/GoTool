package tool

import (
   "runtime/debug"
)

type PanicLogger interface {
	Println(v ...interface{})
}

/*
用于截获程序的PANIC,并输出相应的错误及堆栈信息到日志中，最后调用错误处理函数call
*/
func CutPanic(logger PanicLogger, call func()) {
	if r := recover(); r != nil {
		logger.Println(r, string(debug.Stack()))
		if call != nil {
			call()
		}
	}
}
