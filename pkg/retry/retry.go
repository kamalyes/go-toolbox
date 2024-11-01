/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 22:07:18
 * @FilePath: \go-toolbox\pkg\retry\retry.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package retry

import (
	"time"
)

// Retry 结构体用于实现重试机制
type Retry struct {
	attemptCount int             // 最大尝试次数
	interval     time.Duration   // 重试间隔时间
	errCallFun   ErrCallbackFunc // 错误回调函数
	funcName     string          // 函数名称（可选，用于回调）
}

// DoFun 定义执行函数的类型
type DoFun func() error

// ErrCallbackFunc 是重试时的错误回调函数类型
// nowAttemptCount 表示当前尝试次数
// remainCount 表示剩余尝试次数
// err 是当前执行时的错误信息
type ErrCallbackFunc func(nowAttemptCount, remainCount int, err error, funcName ...string)

// NewRetry 创建一个重试器，返回一个 Retry 实例
func NewRetry() *Retry {
	return &Retry{}
}

// SetAttemptCount 设置最大尝试次数，返回 Retry 实例以支持链式调用
func (r *Retry) SetAttemptCount(attemptCount int) *Retry {
	r.attemptCount = attemptCount
	return r
}

// SetInterval 设置重试间隔时间，返回 Retry 实例以支持链式调用
func (r *Retry) SetInterval(interval time.Duration) *Retry {
	r.interval = interval
	return r
}

// SetErrCallback 设置错误回调函数，返回 Retry 实例以支持链式调用
func (r *Retry) SetErrCallback(errCallbackFunc ErrCallbackFunc) *Retry {
	r.errCallFun = errCallbackFunc
	return r
}

// Do 为 Retry 结构体定义执行函数，执行指定函数 f
func (m *Retry) Do(f DoFun) (err error) {
	return DoRetry(m.attemptCount, m.interval, f, m.errCallFun, m.funcName)
}

// DoRetry 定义了重试操作，执行指定次数的尝试
func DoRetry(attemptCount int, interval time.Duration, f DoFun, errCallFun ErrCallbackFunc, funcName ...string) (err error) {
	nowAttemptCount := 0
	var fName string
	if len(funcName) > 0 {
		fName = funcName[0] // 获取函数名称
	}
	for {
		nowAttemptCount++

		err = f() // 执行传入的函数
		if err == nil {
			return // 如果没有错误，返回
		}
		if errCallFun != nil {
			errCallFun(nowAttemptCount, attemptCount-nowAttemptCount, err, fName) // 调用错误回调函数
		}

		if nowAttemptCount >= attemptCount {
			break // 达到最大尝试次数，退出循环
		}

		if interval > 0 {
			time.Sleep(interval) // 等待指定的间隔时间
		}
	}
	return
}
