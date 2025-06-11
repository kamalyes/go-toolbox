/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-11 15:56:57
 * @FilePath: \go-toolbox\pkg\retry\retry.go
 * @Description: 重试机制
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// Retry 结构体用于实现重试机制
type Retry struct {
	ctx            context.Context     // 上下文
	lock           syncx.Locker        // 并发执行锁
	attemptCount   int                 // 最大尝试次数
	interval       time.Duration       // 重试间隔时间
	errCallFun     ErrCallbackFunc     // 错误回调函数
	successCallFun SuccessCallbackFunc // 成功回调函数
	conditionFunc  func(error) bool    // 重试条件函数
}

// DoFun 定义执行函数的类型
type DoFun func() error

// ErrCallbackFunc 是重试时的错误回调函数类型
// nowAttemptCount 表示当前尝试次数
// remainCount 表示剩余尝试次数
// err 是当前执行时的错误信息
type ErrCallbackFunc func(nowAttemptCount, remainCount int, err error, funcName ...string)

// SuccessCallbackFunc 是成功回调函数类型
type SuccessCallbackFunc func(funcName ...string)

// NewRetry 创建一个重试器，返回一个 Retry 实例
func NewRetry() *Retry {
	return NewRetryWithCtx(context.Background())
}

// NewRetryWithCtx 创建一个自定义上下文的重试器，返回一个 Retry实例
func NewRetryWithCtx(ctx context.Context) *Retry {
	return &Retry{
		ctx: ctx,
	}
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

// SetSuccessCallback 设置成功回调函数，返回 Retry 实例以支持链式调用
func (r *Retry) SetSuccessCallback(successCallbackFunc SuccessCallbackFunc) *Retry {
	r.successCallFun = successCallbackFunc
	return r
}

// SetLock 设置锁
func (r *Retry) SetLock(lock syncx.Locker) *Retry {
	r.lock = lock
	return r
}

// SetConditionFunc 设置重试条件函数
func (r *Retry) SetConditionFunc(f func(error) bool) *Retry {
	r.conditionFunc = f
	return r
}

// Do 为 Retry 结构体定义执行函数，执行指定函数 f
func (r *Retry) Do(f DoFun) (err error) {
	caller := osx.GetRuntimeCaller(3)
	defer caller.Release()
	exec := func() error {
		return doRetryWithCondition(r.ctx, r.attemptCount, r.interval, f, r.errCallFun, r.successCallFun, r.conditionFunc, caller.String())
	}
	return mathx.IfDo(
		r.lock != nil,
		func() error {
			return syncx.WithLockReturnValue(r.lock, exec)
		},
		exec())
}

// doRetryWithCondition 内部函数，定义了重试操作，执行指定次数的尝试
func doRetryWithCondition(ctx context.Context, attemptCount int, interval time.Duration, f DoFun, errCallFun ErrCallbackFunc, successCallFun SuccessCallbackFunc, conditionFunc func(error) bool, funcName ...string) (err error) {
	// 确保尝试次数为正数
	if attemptCount <= 0 {
		return fmt.Errorf("attemptCount must be greater than zero")
	}

	var (
		fName           = mathx.IF(len(funcName) > 0, funcName[0], "") // 获取函数名称
		nowAttemptCount int
	)

	for nowAttemptCount < attemptCount {
		nowAttemptCount++

		select {
		case <-ctx.Done():
			return ctx.Err() // 如果上下文被取消，返回错误
		default:
			// 捕获 panic
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic occurred: %v", r)
				}
			}()

			err = f() // 执行传入的函数
			if err == nil {
				if successCallFun != nil {
					successCallFun(fName) // 调用成功回调函数
				}
				return // 如果没有错误，返回
			}

			// 判断是否满足重试条件，默认全部重试
			if conditionFunc != nil && !conditionFunc(err) {
				return err // 不满足重试条件，直接返回错误
			}

			if errCallFun != nil {
				errCallFun(nowAttemptCount, attemptCount-nowAttemptCount, err, fName) // 调用错误回调函数
			}

			// 等待指定的间隔时间
			if interval > 0 {
				time.Sleep(interval)
			}
		}
	}

	return err // 返回最后的错误
}
