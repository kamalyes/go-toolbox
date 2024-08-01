/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 17:53:10
 * @FilePath: \go-toolbox\retry\retry.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package retry

import (
	"time"
)

type Retry struct {
	attemptCount int
	interval     time.Duration
	errCallFun   ErrCallbackFunc
	funcName     string
}

type Option func(*Retry)

type DoFun func() error

// ErrCallbackFunc 是重试时的错误回调函数类型
// nowAttemptCount 表示当前尝试次数
// remainCount 表示剩余尝试次数
// err 是当前执行时的错误信息
type ErrCallbackFunc func(nowAttemptCount, remainCount int, err error, funcName ...string)

// NewRetry 创建一个重试器
func NewRetry(options ...Option) *Retry {
	r := &Retry{}
	for _, o := range options {
		o(r)
	}
	return r
}

// WithInterval 设置重试间隔时间的选项
func WithInterval(interval time.Duration) Option {
	return func(retry *Retry) {
		retry.interval = interval
	}
}

// WithAttemptCount 设置最大尝试次数的选项，0 表示不限次数
func WithAttemptCount(attemptCount int) Option {
	return func(retry *Retry) {
		retry.attemptCount = attemptCount
	}
}

// WithErrCallback 设置错误回调函数的选项，每次执行时有任何错误都会调用该回调函数
func WithErrCallback(errCallbackFunc ErrCallbackFunc) Option {
	return func(retry *Retry) {
		retry.errCallFun = errCallbackFunc
	}
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
		fName = funcName[0]
	}
	for {
		nowAttemptCount++

		err = f()
		if err == nil {
			return
		}
		if errCallFun != nil {
			errCallFun(nowAttemptCount, attemptCount-nowAttemptCount, err, fName)
		}

		if nowAttemptCount >= attemptCount {
			break
		}

		if interval > 0 {
			time.Sleep(interval)
		}
	}
	return
}
