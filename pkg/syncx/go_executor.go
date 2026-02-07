/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-08 15:25:26
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
	"sync"
	"sync/atomic"
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
	wg       *sync.WaitGroup // 可选：管理子 goroutine
	subTasks []subTask       // 子任务列表
}

// subTask 子任务定义
type subTask struct {
	fn      GoExecutorFunc
	fnError GoExecutorWaitFunc
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

// ChildRunner 子任务运行器，用于在父 goroutine 中启动子 goroutine
type ChildRunner struct {
	wg      sync.WaitGroup
	onPanic GoExecutorPanicHandler
	onError GoExecutorErrorHandler
}

// Go 启动一个子 goroutine
func (c *ChildRunner) Go(fn GoExecutorFunc) *ChildRunner {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer RecoverWithHandler(c.onPanic)
		fn()
	}()
	return c
}

// GoWithError 启动一个带错误返回的子 goroutine
func (c *ChildRunner) GoWithError(fn GoExecutorWaitFunc) *ChildRunner {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer RecoverWithHandler(c.onPanic)

		if err := fn(); err != nil && c.onError != nil {
			c.onError(err)
		}
	}()
	return c
}

// Wait 等待所有子 goroutine 完成
func (c *ChildRunner) Wait() {
	c.wg.Wait()
}

// Sub 添加子 goroutine 任务（链式调用）
//
// 示例:
//
//	Go().OnPanic(handler).
//	    Sub(func() { startHeartbeat() }).
//	    Sub(func() { subscribe1() }).
//	    Sub(func() { subscribe2() }).
//	    ExecSubs()
func (g *GoExecutor) Sub(fn GoExecutorFunc) *GoExecutor {
	g.subTasks = append(g.subTasks, subTask{fn: fn})
	return g
}

// SubWithError 添加带错误返回的子 goroutine 任务（链式调用）
func (g *GoExecutor) SubWithError(fn GoExecutorWaitFunc) *GoExecutor {
	g.subTasks = append(g.subTasks, subTask{fnError: fn})
	return g
}

// ExecSubs 执行所有子 goroutine 任务并等待完成
//
// 示例:
//
//	Go().OnPanic(handler).
//	    Sub(func() { task1() }).
//	    Sub(func() { task2() }).
//	    SubWithError(func() error { return task3() }).
//	    ExecSubs()
func (g *GoExecutor) ExecSubs() {
	go func() {
		defer RecoverWithHandler(g.onPanic)

		// 延迟执行
		if !g.waitWithDelay(g.ctx) {
			return
		}

		// 创建子任务运行器（继承父级的 panic/error handler）
		children := &ChildRunner{
			onPanic: g.onPanic,
			onError: g.onError,
		}

		// 启动所有子任务
		for _, task := range g.subTasks {
			if task.fn != nil {
				children.Go(task.fn)
			} else if task.fnError != nil {
				children.GoWithError(task.fnError)
			}
		}

		// 等待所有子任务完成
		children.Wait()
	}()
}

// ExecWithChildren 在父 goroutine 中执行，并管理多个子 goroutine
//
// 示例:
//
//	Go().OnPanic(handler).ExecWithChildren(func(children *ChildRunner) {
//	    children.Go(func() { startHeartbeat() })
//	    children.Go(func() { subscribe1() })
//	    children.Go(func() { subscribe2() })
//	})
func (g *GoExecutor) ExecWithChildren(fn func(*ChildRunner)) {
	go func() {
		defer RecoverWithHandler(g.onPanic)

		// 延迟执行
		if !g.waitWithDelay(g.ctx) {
			return
		}

		// 创建子任务运行器（继承父级的 panic/error handler）
		children := &ChildRunner{
			onPanic: g.onPanic,
			onError: g.onError,
		}

		// 执行用户函数（启动子任务）
		fn(children)

		// 等待所有子任务完成
		children.Wait()
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

// BatchExecutorMode 批量执行器模式
type BatchExecutorMode int

const (
	// FailFastMode 快速失败模式：遇到第一个错误立即停止提交新任务
	FailFastMode BatchExecutorMode = iota
	// ContinueOnErrorMode 继续执行模式：即使有错误也继续执行所有任务
	ContinueOnErrorMode
)

// BatchExecutor 批量并发执行器，支持并发限制和两种错误处理模式
// 所有方法都是并发安全的
type BatchExecutor struct {
	ctx       context.Context
	cancel    context.CancelFunc
	mode      BatchExecutorMode // 执行模式
	limit     int
	onPanic   GoExecutorPanicHandler
	onError   GoExecutorErrorHandler
	semaphore chan struct{}   // channel 天然并发安全
	wg        sync.WaitGroup  // WaitGroup 并发安全
	taskID    atomic.Int64    // 任务ID生成器
	errMu     sync.RWMutex    // 保护 errors 的读写
	errors    map[int64]error // 每个任务的错误映射 taskID -> error
	firstErr  error           // 第一个错误（快速失败模式用）
}

// NewBatchExecutor 创建批量并发执行器（默认快速失败模式）
//
// 示例:
//
//	exec := NewBatchExecutor(ctx).
//	    SetLimit(10).
//	    SetMode(ContinueOnErrorMode).  // 可选：设置继续执行模式
//	    OnPanic(func(r interface{}) { log.Error("panic", r) }).
//	    OnError(func(err error) { log.Error(err) })
//
//	for _, item := range items {
//	    exec.Go(func() error {
//	        return processItem(item)
//	    })
//	}
//
//	err := exec.Wait()  // 快速失败模式返回第一个错误
//	errors := exec.Errors()  // 获取所有错误映射
func NewBatchExecutor(ctx context.Context) *BatchExecutor {
	ctx, cancel := context.WithCancel(ctx)
	return &BatchExecutor{
		ctx:    ctx,
		cancel: cancel,
		mode:   FailFastMode, // 默认快速失败
		errors: make(map[int64]error),
	}
}

// SetLimit 设置最大并发数
func (b *BatchExecutor) SetLimit(n int) *BatchExecutor {
	b.limit = n
	if n > 0 {
		b.semaphore = make(chan struct{}, n)
	}
	return b
}

// SetMode 设置执行模式
func (b *BatchExecutor) SetMode(mode BatchExecutorMode) *BatchExecutor {
	b.mode = mode
	return b
}

// OnPanic 设置 panic 处理器
func (b *BatchExecutor) OnPanic(fn GoExecutorPanicHandler) *BatchExecutor {
	b.onPanic = fn
	return b
}

// OnError 设置错误处理器（每个错误都会调用）
func (b *BatchExecutor) OnError(fn GoExecutorErrorHandler) *BatchExecutor {
	b.onError = fn
	return b
}

// Go 提交一个任务（并发安全）
func (b *BatchExecutor) Go(fn func() error) {
	// 快速失败模式：检查是否已有错误
	if b.mode == FailFastMode {
		select {
		case <-b.ctx.Done():
			return
		default:
		}
	}

	// 生成任务ID
	taskID := b.taskID.Add(1)

	// 限流控制（阻塞式）
	if b.semaphore != nil {
		select {
		case b.semaphore <- struct{}{}:
		case <-b.ctx.Done():
			// 快速失败模式下，context 取消后不再执行
			if b.mode == FailFastMode {
				return
			}
		}
	}

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		if b.semaphore != nil {
			defer func() { <-b.semaphore }()
		}

		// panic 恢复和错误处理（局部变量，无并发问题）
		var err error
		defer RecoverAndHandle(&err, b.onPanic, func(e error) {
			// 统一处理错误（panic 转换的或正常返回的）
			if b.onError != nil {
				b.onError(e)
			}

			// 记录错误到 map
			b.errMu.Lock()
			b.errors[taskID] = e

			// 记录第一个错误
			if b.firstErr == nil {
				b.firstErr = e
				// 快速失败模式：取消 context
				if b.mode == FailFastMode {
					b.cancel()
				}
			}
			b.errMu.Unlock()
		})

		// 二次检查 context（可能在等待 semaphore 期间被取消）
		if b.mode == FailFastMode {
			select {
			case <-b.ctx.Done():
				return
			default:
			}
		}

		// 执行任务
		err = fn()
	}()
}

// Wait 等待所有任务完成，返回第一个错误（并发安全）
func (b *BatchExecutor) Wait() error {
	b.wg.Wait()
	b.cancel() // 确保 context 被取消

	b.errMu.RLock()
	err := b.firstErr
	b.errMu.RUnlock()
	return err
}

// Errors 获取所有错误映射（并发安全）
// 返回 map[taskID]error，其中 taskID 是任务提交的顺序编号（从1开始）
func (b *BatchExecutor) Errors() map[int64]error {
	b.errMu.RLock()
	defer b.errMu.RUnlock()

	// 返回副本，避免外部修改
	result := make(map[int64]error, len(b.errors))
	for k, v := range b.errors {
		result[k] = v
	}
	return result
}

// ErrorCount 获取错误总数（并发安全）
func (b *BatchExecutor) ErrorCount() int {
	b.errMu.RLock()
	defer b.errMu.RUnlock()
	return len(b.errors)
}

// HasErrors 检查是否有错误（并发安全）
func (b *BatchExecutor) HasErrors() bool {
	b.errMu.RLock()
	defer b.errMu.RUnlock()
	return len(b.errors) > 0
}

// Context 返回执行器的 context
func (b *BatchExecutor) Context() context.Context {
	return b.ctx
}
