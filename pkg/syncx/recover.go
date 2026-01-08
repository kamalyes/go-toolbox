/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-08 15:27:17
 * @FilePath: \go-toolbox\pkg\syncx\recover.go
 * @Description: 统一的 panic 恢复处理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

// RecoverFunc panic 恢复处理函数类型
type RecoverFunc func(interface{})

// SafeGo 安全的 goroutine 启动，带 panic 恢复
//
// 示例:
//
//	SafeGo(func() {
//	    // 可能会 panic 的代码
//	}, func(r interface{}) {
//	    log.Error("panic recovered", r)
//	})
func SafeGo(fn func(), onPanic RecoverFunc) {
	go func() {
		defer RecoverWithHandler(onPanic)
		fn()
	}()
}

// RecoverWithHandler 使用自定义处理器恢复 panic
//
// 示例:
//
//	defer RecoverWithHandler(func(r interface{}) {
//	    log.Error("panic", r)
//	})
func RecoverWithHandler(handler RecoverFunc) {
	if r := recover(); r != nil && handler != nil {
		handler(r)
	}
}

// Recover 简单的 panic 恢复（忽略错误）
//
// 示例:
//
//	defer Recover()
func Recover() {
	recover()
}

// MustRecover 必须恢复的 panic（会重新 panic）
//
// 示例:
//
//	defer MustRecover(func(r interface{}) {
//	    log.Error("critical panic", r)
//	})
func MustRecover(handler RecoverFunc) {
	if r := recover(); r != nil {
		if handler != nil {
			handler(r)
		}
		panic(r) // 重新抛出
	}
}

// RecoverToError 将 panic 转换为 error
//
// 示例:
//
//	func example() (err error) {
//	    defer RecoverToError(&err, nil)
//	    可能会 panic 的代码
//	    return nil
//	}
func RecoverToError(err *error, handler RecoverFunc) {
	if r := recover(); r != nil {
		if handler != nil {
			handler(r)
		}
		if err != nil {
			if e, ok := r.(error); ok {
				*err = e
			} else {
				*err = &panicError{value: r}
			}
		}
	}
}

// RecoverAndHandle 恢复 panic 并在 defer 中处理错误
// 用于需要在同一个 defer 中完成 panic 恢复和错误处理的场景
//
// 示例:
//
//	func example() error {
//	    var err error
//	    defer RecoverAndHandle(&err,
//	        func(r interface{}) { log.Error("panic", r) },  // panic handler
//	        func(e error) { log.Error("error", e) })        // error handler
//
//	    err = doSomething()
//	    return err
//	}
func RecoverAndHandle(err *error, panicHandler RecoverFunc, errorHandler func(error)) {
	// 先恢复 panic
	if r := recover(); r != nil {
		if panicHandler != nil {
			panicHandler(r)
		}
		if err != nil {
			if e, ok := r.(error); ok {
				*err = e
			} else {
				*err = &panicError{value: r}
			}
		}
	}

	// 再处理错误（无论是 panic 转换的还是正常返回的）
	if err != nil && *err != nil && errorHandler != nil {
		errorHandler(*err)
	}
}

// panicError panic 错误包装
type panicError struct {
	value interface{}
}

func (p *panicError) Error() string {
	return formatPanic(p.value)
}

// formatPanic 格式化 panic 值
func formatPanic(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return "panic occurred"
	}
}
