/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 00:00:00 09:00:00
 * @FilePath: \go-toolbox\pkg\syncx\go_executor.go
 * @Description: Goroutine 执行器 - 链式调用风格，集成 contextx
 *
 * 使用说明:
 *
 * 1. 基础 Goroutine 执行:
 *    Go().
 *        OnPanic(func(r interface{}) { log.Error("panic", r) }).
 *        Exec(func() { doSomething() })
 *
 * 2. 带超时的 Context 执行:
 *    Go(ctx).
 *        WithTimeout(2 * time.Second).
 *        OnError(func(err error) { log.Error(err) }).
 *        OnPanic(func(r interface{}) { log.Error("panic", r) }).
 *        ExecWithContext(func(ctx context.Context) error {
 *            return repo.Save(ctx, data)
 *        })
 *
 * 3. 带延迟的执行:
 *    Go(ctx).
 *        WithDelay(5 * time.Second).
 *        OnCancel(func() { log.Info("cancelled") }).
 *        ExecWithContext(func(ctx context.Context) error {
 *            return sendNotification(ctx)
 *        })
 *
 * 4. 完整示例:
 *    Go(ctx).
 *        WithTimeout(10 * time.Second).
 *        WithDelay(1 * time.Second).
 *        OnError(func(err error) { metrics.RecordError(err) }).
 *        OnPanic(func(r interface{}) { sentry.CaptureException(r) }).
 *        OnCancel(func() { log.Info("task cancelled") }).
 *        ExecWithContext(func(ctx context.Context) error {
 *            return processData(ctx)
 *        })
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"time"
)

// GoExecutorErrorHandler 错误处理函数类型
type GoExecutorErrorHandler func(error)

// GoExecutorPanicHandler panic 处理函数类型（复用 RecoverFunc）
type GoExecutorPanicHandler = RecoverFunc

// GoExecutorCancelHandler 取消处理函数类型
type GoExecutorCancelHandler func()

// GoExecutorFunc 无返回值的执行函数类型
type GoExecutorFunc func()

// GoExecutorContextFunc 带 Context 的执行函数类型
type GoExecutorContextFunc func(context.Context) error

// GoExecutorResultFunc 带返回值的执行函数类型
type GoExecutorResultFunc func() (interface{}, error)

// GoExecutorWaitFunc 同步等待的执行函数类型
type GoExecutorWaitFunc func() error

// waitWithDelay 等待延迟时间，支持取消
// 返回 true 表示正常完成延迟，false 表示被取消
func (g *GoExecutor) waitWithDelay(ctx context.Context) bool {
	if g.delay <= 0 {
		return true
	}

	timer := time.NewTimer(g.delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		if g.onCancel != nil {
			g.onCancel()
		}
		return false
	}
}

// GoExecutor Goroutine 执行器
type GoExecutor struct {
	timeout  time.Duration
	delay    time.Duration
	onError  GoExecutorErrorHandler
	onPanic  GoExecutorPanicHandler
	onCancel GoExecutorCancelHandler
	ctx      context.Context
}

// Go 创建一个新的 Goroutine 执行器
// 参数 ctx 可选，如果不传或传 nil，则使用 context.Background()
//
// 示例:
//
//	Go(ctx).OnPanic(handler).Exec(func() { ... })
//	Go().OnPanic(handler).Exec(func() { ... })  // 使用 background context
func Go(ctx ...context.Context) *GoExecutor {
	var parentCtx context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		parentCtx = ctx[0]
	} else {
		parentCtx = context.Background()
	}

	return &GoExecutor{
		ctx: parentCtx,
	}
}

// GoWithContext 创建带父 Context 的 Goroutine 执行器（已废弃，使用 Go(ctx) 代替）
//
// 示例:
//
//	Go(parentCtx).WithTimeout(5*time.Second).ExecWithContext(fn)
//
// Deprecated: 使用 Go(ctx) 代替
func GoWithContext(parent context.Context) *GoExecutor {
	return Go(parent)
}

// WithTimeout 设置超时时间
//
// 示例:
//
//	Go().WithTimeout(5*time.Second).ExecWithContext(fn)
func (g *GoExecutor) WithTimeout(timeout time.Duration) *GoExecutor {
	g.timeout = timeout
	return g
}

// WithDelay 设置延迟执行时间
//
// 示例:
//
//	Go().WithDelay(2*time.Second).Exec(fn)
func (g *GoExecutor) WithDelay(delay time.Duration) *GoExecutor {
	g.delay = delay
	return g
}

// OnError 设置错误回调
//
// 示例:
//
//	Go().OnError(func(err error) { log.Error(err) }).ExecWithContext(fn)
func (g *GoExecutor) OnError(fn GoExecutorErrorHandler) *GoExecutor {
	g.onError = fn
	return g
}

// OnPanic 设置 panic 回调
//
// 示例:
//
//	Go().OnPanic(func(r interface{}) { log.Error("panic", r) }).Exec(fn)
func (g *GoExecutor) OnPanic(fn GoExecutorPanicHandler) *GoExecutor {
	g.onPanic = fn
	return g
}

// OnCancel 设置取消回调
//
// 示例:
//
//	Go().OnCancel(func() { log.Info("cancelled") }).ExecWithContext(fn)
func (g *GoExecutor) OnCancel(fn GoExecutorCancelHandler) *GoExecutor {
	g.onCancel = fn
	return g
}

// Exec 执行无参数的函数
//
// 示例:
//
//	Go().OnPanic(handler).Exec(func() {
//	    doSomething()
//	})
func (g *GoExecutor) Exec(fn GoExecutorFunc) {
	go func() {
		defer RecoverWithHandler(g.onPanic)

		// 延迟执行
		if !g.waitWithDelay(g.ctx) {
			return
		}

		fn()
	}()
}

// ExecWithContext 执行带 Context 的函数
//
// 示例:
//
//	Go().WithTimeout(5*time.Second).OnError(handler).ExecWithContext(func(ctx context.Context) error {
//	    return repo.Save(ctx, data)
//	})
func (g *GoExecutor) ExecWithContext(fn GoExecutorContextFunc) {
	go func() {
		defer RecoverWithHandler(g.onPanic)

		// 创建 context
		ctx := g.ctx
		var cancel context.CancelFunc

		if g.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, g.timeout)
			defer cancel()
		}

		// 延迟执行
		if !g.waitWithDelay(ctx) {
			return
		}

		// 执行函数
		if err := fn(ctx); err != nil {
			if g.onError != nil {
				g.onError(err)
			}
		}
	}()
}

// ExecWithResult 执行带返回值的函数(异步获取结果)
//
// 示例:
//
//	resultChan := Go().ExecWithResult(func() (int, error) {
//	    return computeValue(), nil
//	})
//	select {
//	case result := <-resultChan:
//	    if result.Err != nil { ... }
//	    value := result.Value
//	}
func (g *GoExecutor) ExecWithResult(fn GoExecutorResultFunc) <-chan struct {
	Value interface{}
	Err   error
} {
	resultChan := make(chan struct {
		Value interface{}
		Err   error
	}, 1)

	go func() {
		defer close(resultChan)
		defer RecoverWithHandler(g.onPanic)

		// 延迟执行
		if !g.waitWithDelay(g.ctx) {
			return
		}

		value, err := fn()
		if err != nil && g.onError != nil {
			g.onError(err)
		}

		resultChan <- struct {
			Value interface{}
			Err   error
		}{value, err}
	}()

	return resultChan
}

// Wait 同步等待执行完成
//
// 示例:
//
//	Go().Wait(func() error {
//	    return doSomething()
//	})
func (g *GoExecutor) Wait(fn GoExecutorWaitFunc) error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		defer RecoverWithHandler(g.onPanic)

		// 延迟执行
		if !g.waitWithDelay(g.ctx) {
			errChan <- g.ctx.Err()
			return
		}

		if err := fn(); err != nil {
			if g.onError != nil {
				g.onError(err)
			}
			errChan <- err
			return
		}
		errChan <- nil
	}()

	return <-errChan
}
