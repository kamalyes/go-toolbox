/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-08-11 09:27:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-20 13:15:15
 * @FilePath: \go-toolbox\pkg\syncx\func_chain.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"fmt"
	"sort"
	"sync"
)

// ReturnFunc 是一个可返回错误的函数类型
type ReturnFunc[T any] func() (T, error)

// FuncItem 表示一个带有优先级的函数
type FuncItem[T any] struct {
	priority int           // 优先级
	fn       ReturnFunc[T] // 执行函数
	result   T             // 存储执行结果
	err      error         // 存储执行过程中产生的错误
}

// FuncChain 是一个支持可选返回值和错误的函数链
type FuncChain[T any] struct {
	funcs []FuncItem[T] // 使用切片存储函数和它们的优先级
	mu    sync.RWMutex  // 读写互斥锁
}

// NewFuncChain 创建一个新的 FuncChain 实例
func NewFuncChain[T any]() *FuncChain[T] {
	return &FuncChain[T]{
		funcs: []FuncItem[T]{},
	}
}

// NewFuncItem 创建一个新的 FuncItem 实例，并设置默认超时和优先级
func NewFuncItem[T any](f ReturnFunc[T]) *FuncItem[T] {
	return &FuncItem[T]{
		fn:       f,
		priority: -1, // 默认优先级设置为-1
	}
}

// WithPriority 设置 FuncItem 的优先级
func (fi *FuncItem[T]) WithPriority(priority int) *FuncItem[T] {
	fi.priority = priority
	return fi
}

// GetResult 获取 FuncItem Result
func (fi *FuncItem[T]) GetResult() T {
	return fi.result
}

// GetError 获取 FuncItem Error
func (fi *FuncItem[T]) GetError() error {
	return fi.err
}

// AddFuncItem 将 FuncItem 添加到 FuncChain 中
func (fc *FuncChain[T]) AddFuncItem(item *FuncItem[T]) *FuncChain[T] {
	return WithLockReturnValue(&fc.mu, func() *FuncChain[T] {
		fc.funcs = append(fc.funcs, *item)
		return fc
	})
}

// Clear 清空 FuncChain 中的所有 FuncItem
func (fc *FuncChain[T]) Clear() {
	WithLock(&fc.mu, func() {
		fc.funcs = []FuncItem[T]{}
	})
}

// GetFuncItems 返回 FuncChain 中的所有 FuncItem
func (fc *FuncChain[T]) GetFuncItems() []FuncItem[T] {
	return WithRLockReturnValue(&fc.mu, func() []FuncItem[T] {
		return fc.funcs
	})
}

// Execute 按顺序执行所有添加的函数，使用协程，并支持超时
func (fc *FuncChain[T]) Execute() {
	WithLock(&fc.mu, func() {
		// 对函数按优先级进行排序，优先级越小的排在前面
		sort.Slice(fc.funcs, func(i, j int) bool {
			return fc.funcs[i].priority < fc.funcs[j].priority // 优先级越小，排在前面
		})
	})

	var wg sync.WaitGroup

	// 为每个 FuncItem 启动一个 goroutine
	for i := range fc.funcs {
		item := &fc.funcs[i] // 获取指向当前 FuncItem 的指针
		wg.Add(1)

		go func(item *FuncItem[T]) {
			defer wg.Done()
			// 调用函数并捕获结果和错误
			defer func() {
				if r := recover(); r != nil {
					item.err = fmt.Errorf("panic: %v", r) // 处理恐慌
				}
			}()
			// 在超时上下文内执行函数
			item.result, item.err = item.fn() // 正常调用
		}(item)
	}

	wg.Wait() // 等待所有 goroutine 完成
}
