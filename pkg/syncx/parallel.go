/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 08:52:55
 * @FilePath: \go-toolbox\pkg\syncx\parallel.go
 * @Description: 并发执行工具函数
 *
 * 使用说明:
 *
 * 1. 回调风格 API (推荐):
 *    - NewParallelExecutor: 创建 map 并发执行器,支持链式调用设置回调
 *    - NewParallelSliceExecutor: 创建 slice 并发执行器,支持链式调用设置回调
 *
 *    可用回调:
 *      OnSuccess(fn)      - 每个任务成功时调用
 *      OnError(fn)        - 每个任务失败时调用
 *      OnEachComplete(fn) - 每个任务完成时调用(无论成败)
 *      OnComplete(fn)     - 所有任务完成后调用,返回结果和错误集合
 *
 *    示例:
 *      NewParallelExecutor(clientMap).
 *          OnSuccess(func(key, val, result) { log.Info("成功") }).
 *          OnError(func(key, val, err) { log.Error("失败", err) }).
 *          OnComplete(func(results, errors) { log.Info("全部完成") }).
 *          Execute(func(key, val) (result, error) { return process(val) })
 *
 * 2. 简化函数 (无返回值场景):
 *    - ParallelForEach: 并发遍历 map,无返回值
 *    - ParallelForEachSlice: 并发遍历 slice,无返回值
 *
 *    示例:
 *      ParallelForEach(clientMap, func(key, client) { client.Close() })
 *
 * 性能对比 (WaitGroup vs Channel):
 *   - 小数据(3元素):   WaitGroup 快 21%, 内存少 37.5%
 *   - 中等数据(100元素): WaitGroup 快 25%
 *   - 大数据(1000元素):  性能相近
 *   结论: 优先使用 WaitGroup 实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"fmt"
	"sync"
)

// ParallelExecuteFunc 并发执行函数类型 (Map)
type ParallelExecuteFunc[K comparable, V any, R any] func(K, V) (R, error)

// ParallelSliceExecuteFunc 并发执行函数类型 (Slice)
type ParallelSliceExecuteFunc[T any, R any] func(int, T) (R, error)

// ParallelSuccessCallback 成功回调函数类型 (Map)
type ParallelSuccessCallback[K comparable, V any, R any] func(K, V, R)

// ParallelSliceSuccessCallback 成功回调函数类型 (Slice)
type ParallelSliceSuccessCallback[T any, R any] func(int, T, R)

// ParallelErrorCallback 错误回调函数类型 (Map)
type ParallelErrorCallback[K comparable, V any] func(K, V, error)

// ParallelSliceErrorCallback 错误回调函数类型 (Slice)
type ParallelSliceErrorCallback[T any] func(int, T, error)

// ParallelCompleteCallback 完成回调函数类型 (Map)
type ParallelCompleteCallback[K comparable, R any] func(map[K]R, map[K]error)

// ParallelSliceCompleteCallback 完成回调函数类型 (Slice)
type ParallelSliceCompleteCallback[R any] func([]R, []error)

// ParallelEachCompleteCallback 每个任务完成回调函数类型 (Map)
type ParallelEachCompleteCallback[K comparable] func(K)

// ParallelSliceEachCompleteCallback 每个任务完成回调函数类型 (Slice)
type ParallelSliceEachCompleteCallback func(int)

// ParallelPanicCallback panic 回调函数类型 (Map)
type ParallelPanicCallback[K comparable, V any] func(K, V, any)

// ParallelSlicePanicCallback panic 回调函数类型 (Slice)
type ParallelSlicePanicCallback[T any] func(int, T, any)

// ParallelForEachFunc 遍历函数类型 (Map)
type ParallelForEachFunc[K comparable, V any] func(K, V)

// ParallelForEachSliceFunc 遍历函数类型 (Slice)
type ParallelForEachSliceFunc[T any] func(int, T)

// ParallelExecutor 并发执行器配置
//
// 泛型参数说明:
//   - K: Map 的键类型 (Key),必须是可比较类型
//   - V: Map 的值类型 (Value),任意类型
//   - R: 执行结果类型 (Result),任意类型
//
// 示例:
//
//	ParallelExecutor[string, *Client, bool]
//	- K=string: 用户ID作为键
//	- V=*Client: 客户端连接对象
//	- R=bool: 返回发送是否成功
type ParallelExecutor[K comparable, V any, R any] struct {
	data           map[K]V
	onSuccess      ParallelSuccessCallback[K, V, R]
	onError        ParallelErrorCallback[K, V]
	onComplete     ParallelCompleteCallback[K, R]
	onEachComplete ParallelEachCompleteCallback[K]
	onPanic        ParallelPanicCallback[K, V]
}

// ParallelSliceExecutor 并发 Slice 执行器配置
//
// 泛型参数说明:
//   - T: Slice 元素类型 (Type),任意类型
//   - R: 执行结果类型 (Result),任意类型
//
// 示例:
//
//	ParallelSliceExecutor[*Client, int]
//	- T=*Client: 客户端连接对象
//	- R=int: 返回发送的字节数
type ParallelSliceExecutor[T any, R any] struct {
	data           []T
	maxConcurrency int // 最大并发数，0表示不限制
	onSuccess      ParallelSliceSuccessCallback[T, R]
	onError        ParallelSliceErrorCallback[T]
	onComplete     ParallelSliceCompleteCallback[R]
	onEachComplete ParallelSliceEachCompleteCallback
	onPanic        ParallelSlicePanicCallback[T]
}

// NewParallelExecutor 创建并发执行器
func NewParallelExecutor[K comparable, V any, R any](m map[K]V) *ParallelExecutor[K, V, R] {
	return &ParallelExecutor[K, V, R]{
		data: m,
	}
}

// NewParallelSliceExecutor 创建并发 Slice 执行器
func NewParallelSliceExecutor[T any, R any](s []T) *ParallelSliceExecutor[T, R] {
	return &ParallelSliceExecutor[T, R]{
		data: s,
	}
}

// OnSuccess 设置成功回调
func (p *ParallelExecutor[K, V, R]) OnSuccess(fn ParallelSuccessCallback[K, V, R]) *ParallelExecutor[K, V, R] {
	p.onSuccess = fn
	return p
}

// OnError 设置错误回调
func (p *ParallelExecutor[K, V, R]) OnError(fn ParallelErrorCallback[K, V]) *ParallelExecutor[K, V, R] {
	p.onError = fn
	return p
}

// OnComplete 设置完成回调(所有任务完成后调用)
func (p *ParallelExecutor[K, V, R]) OnComplete(fn ParallelCompleteCallback[K, R]) *ParallelExecutor[K, V, R] {
	p.onComplete = fn
	return p
}

// OnEachComplete 设置每个任务完成时的回调
func (p *ParallelExecutor[K, V, R]) OnEachComplete(fn ParallelEachCompleteCallback[K]) *ParallelExecutor[K, V, R] {
	p.onEachComplete = fn
	return p
}

// OnPanic 设置 panic 回调
func (p *ParallelExecutor[K, V, R]) OnPanic(fn ParallelPanicCallback[K, V]) *ParallelExecutor[K, V, R] {
	p.onPanic = fn
	return p
}

// Execute 执行并发任务
func (p *ParallelExecutor[K, V, R]) Execute(fn ParallelExecuteFunc[K, V, R]) {
	if len(p.data) == 0 {
		if p.onComplete != nil {
			p.onComplete(make(map[K]R), make(map[K]error))
		}
		return
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results = make(map[K]R, len(p.data))
		errors  = make(map[K]error)
	)

	for k, v := range p.data {
		wg.Add(1)
		go func(key K, val V) {
			defer func() {
				if r := recover(); r != nil {
					mu.Lock()
					errors[key] = fmt.Errorf("panic: %v", r)
					mu.Unlock()

					if p.onPanic != nil {
						p.onPanic(key, val, r)
					}
				}
				wg.Done()
			}()

			result, err := fn(key, val)

			mu.Lock()
			if err != nil {
				errors[key] = err
				if p.onError != nil {
					p.onError(key, val, err)
				}
			} else {
				results[key] = result
				if p.onSuccess != nil {
					p.onSuccess(key, val, result)
				}
			}

			if p.onEachComplete != nil {
				p.onEachComplete(key)
			}
			mu.Unlock()
		}(k, v)
	}

	wg.Wait()

	if p.onComplete != nil {
		p.onComplete(results, errors)
	}
}

// OnSuccess 设置成功回调
func (p *ParallelSliceExecutor[T, R]) OnSuccess(fn ParallelSliceSuccessCallback[T, R]) *ParallelSliceExecutor[T, R] {
	p.onSuccess = fn
	return p
}

// OnError 设置错误回调
func (p *ParallelSliceExecutor[T, R]) OnError(fn ParallelSliceErrorCallback[T]) *ParallelSliceExecutor[T, R] {
	p.onError = fn
	return p
}

// OnComplete 设置完成回调(所有任务完成后调用)
func (p *ParallelSliceExecutor[T, R]) OnComplete(fn ParallelSliceCompleteCallback[R]) *ParallelSliceExecutor[T, R] {
	p.onComplete = fn
	return p
}

// OnEachComplete 设置每个任务完成时的回调
func (p *ParallelSliceExecutor[T, R]) OnEachComplete(fn ParallelSliceEachCompleteCallback) *ParallelSliceExecutor[T, R] {
	p.onEachComplete = fn
	return p
}

// OnPanic 设置 panic 回调
func (p *ParallelSliceExecutor[T, R]) OnPanic(fn ParallelSlicePanicCallback[T]) *ParallelSliceExecutor[T, R] {
	p.onPanic = fn
	return p
}

// WithConcurrency 设置最大并发数
//
// 参数:
//   - maxConcurrency: 最大并发 goroutine 数量
//   - 0: 不限制并发数，所有任务立即启动（默认行为）
//   - >0: 限制同时运行的 goroutine 数量，使用信号量控制
//
// 注意:
//   - 这只是控制并发数，所有任务最终都会被执行完
//   - 例如: 100个任务，maxConcurrency=10，表示同时最多10个任务在执行，
//     当一个任务完成后，会自动启动下一个任务，直到所有100个任务都完成
func (p *ParallelSliceExecutor[T, R]) WithConcurrency(maxConcurrency int) *ParallelSliceExecutor[T, R] {
	p.maxConcurrency = maxConcurrency
	return p
}

// Execute 执行并发任务
//
// 工作原理:
//   - 如果设置了 maxConcurrency，使用信号量控制同时运行的 goroutine 数量
//   - 所有任务都会被执行，只是控制了并发度
//   - 使用 WaitGroup 确保所有任务完成后才返回
//
// 示例:
//   - 100个任务，maxConcurrency=10:
//     同时最多10个 goroutine 在运行，当一个完成后，立即启动下一个
//     最终所有100个任务都会执行完成
func (p *ParallelSliceExecutor[T, R]) Execute(fn ParallelSliceExecuteFunc[T, R]) {
	if len(p.data) == 0 {
		if p.onComplete != nil {
			p.onComplete(nil, nil)
		}
		return
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results = make([]R, len(p.data))
		errors  = make([]error, len(p.data))
	)

	// 处理单个任务的公共逻辑
	processTask := func(idx int, val T) {
		defer func() {
			if r := recover(); r != nil {
				mu.Lock()
				errors[idx] = fmt.Errorf("panic: %v", r)
				mu.Unlock()

				if p.onPanic != nil {
					p.onPanic(idx, val, r)
				}
			}
			wg.Done()
		}()

		result, err := fn(idx, val)

		mu.Lock()
		defer mu.Unlock()

		if err != nil {
			errors[idx] = err
			if p.onError != nil {
				p.onError(idx, val, err)
			}
		} else {
			results[idx] = result
			if p.onSuccess != nil {
				p.onSuccess(idx, val, result)
			}
		}

		if p.onEachComplete != nil {
			p.onEachComplete(idx)
		}
	}

	// 如果设置了并发限制，使用信号量控制
	// 注意: 所有任务都会被执行，信号量只是控制同时运行的数量
	if p.maxConcurrency > 0 {
		semaphore := make(chan struct{}, p.maxConcurrency)

		for i, v := range p.data {
			wg.Add(1)
			semaphore <- struct{}{} // 获取信号量（阻塞直到有空位）

			go func(idx int, val T) {
				defer func() { <-semaphore }() // 释放信号量，让下一个任务可以启动
				processTask(idx, val)
			}(i, v)
		}
	} else {
		// 无并发限制，直接启动所有 goroutine
		for i, v := range p.data {
			wg.Add(1)
			go processTask(i, v)
		}
	}

	wg.Wait()

	if p.onComplete != nil {
		p.onComplete(results, errors)
	}
}

// ParallelForEach 并发遍历 map 并对每个元素执行操作
// 使用 WaitGroup 确保所有 goroutine 完成后才返回
//
// 参数:
//   - m: 要遍历的 map
//   - fn: 对每个元素执行的函数
//
// 示例:
//
//	syncx.ParallelForEach(clientMap, func(key string, client *Client) {
//	    client.Send(msg)
//	})
func ParallelForEach[K comparable, V any](m map[K]V, fn ParallelForEachFunc[K, V]) {
	if len(m) == 0 {
		return
	}

	var wg sync.WaitGroup
	for k, v := range m {
		wg.Add(1)
		go func(key K, val V) {
			defer wg.Done()
			fn(key, val)
		}(k, v)
	}
	wg.Wait()
}

// ParallelForEachSlice 并发遍历 slice 并对每个元素执行操作
// 使用 WaitGroup 确保所有 goroutine 完成后才返回
//
// 参数:
//   - s: 要遍历的 slice
//   - fn: 对每个元素执行的函数
//
// 示例:
//
//	syncx.ParallelForEachSlice(clients, func(i int, client *Client) {
//	    client.Send(msg)
//	})
func ParallelForEachSlice[T any](s []T, fn ParallelForEachSliceFunc[T]) {
	if len(s) == 0 {
		return
	}

	var wg sync.WaitGroup
	for i, v := range s {
		wg.Add(1)
		go func(idx int, val T) {
			defer wg.Done()
			fn(idx, val)
		}(i, v)
	}
	wg.Wait()
}
