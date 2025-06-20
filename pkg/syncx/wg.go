/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 13:55:22
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-09 18:25:16
 * @FilePath: \go-toolbox\pkg\syncx\wg.go
 * @Description: 自定义的 WaitGroup 结构体，封装了 sync.WaitGroup，用于管理并发操作的等待
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"fmt"
	"sync"
)

// WaitGroup 是一个自定义的等待组，支持并发操作并捕获错误
type WaitGroup struct {
	catchPanic bool           // 捕获异常
	err        error          // 用于存储第一个发生的错误
	ch         chan struct{}  // 用于限制并发数量的通道
	wg         sync.WaitGroup // 内置的等待组，用于等待 goroutine 完成
	mu         sync.RWMutex   // 读写锁，用于保护 err 字段
}

// NewWaitGroup 创建一个新的 WaitGroup，支持最大并发数量
//
// params:
//   - catchPanic: 是否捕获 goroutine 中的 panic
//   - maxChannel: 可选参数，最大并发数量
//
// return:
//   - *WaitGroup: 返回初始化后的 WaitGroup 实例
func NewWaitGroup(catchPanic bool, maxChannel ...uint) *WaitGroup {
	wg := &WaitGroup{
		catchPanic: catchPanic,
	}
	// 如果提供了最大并发数量，则初始化通道
	if len(maxChannel) > 0 {
		wg.ch = make(chan struct{}, maxChannel[0])
	}

	return wg
}

// SetError 设置错误
//
// params:
//   - err: 要设置的错误，如果为 nil 则不进行设置
func (wg *WaitGroup) SetError(err error) {
	WithLock(&wg.mu, func() {
		if err != nil {
			wg.err = err
		}
	})
}

// GetError 返回错误，如果没有错误则返回 nil
//
// return:
//   - error: 捕获的错误，如果没有则返回 nil
func (wg *WaitGroup) GetError() error {
	return WithRLockReturnValue[error](&wg.mu, func() error {
		return wg.err
	})
}

// GetChannelSize 返回当前通道的可用容量
//
// return:
//   - int: 当前通道的可用容量
func (wg *WaitGroup) GetChannelSize() int {
	return WithRLockReturnValue[int](&wg.mu, func() int {
		return len(wg.ch)
	})
}

// Add 增加等待计数
//
// params:
//   - delta: 增加的计数值，通常为 1
func (h *WaitGroup) Add(delta int) {
	h.wg.Add(delta)
}

// Done 减少等待计数
// 每当一个 goroutine 完成时调用此方法
func (h *WaitGroup) Done() {
	h.wg.Done()
}

// Go 启动一个新的 goroutine，并根据需要处理 panic
//
// params:
//   - f: 要执行的函数
func (h *WaitGroup) Go(f func()) {
	// 如果通道不为 nil，则向通道发送一个信号，表示一个新的 goroutine 启动
	if h.ch != nil {
		h.ch <- struct{}{}
	}
	h.Add(1) // 增加等待计数
	started := make(chan struct{})

	go func() {
		defer func() {
			// 在 goroutine 完成后，释放通道信号并减少等待计数
			if h.ch != nil {
				<-h.ch
			}
			h.Done()
		}()
		close(started) // 通知goroutine已启动

		if h.catchPanic {
			h.handlePanic(f) // 调用处理 panic 的方法
			return
		}
		f() // 执行传入的函数
	}()

	<-started // 等待goroutine启动
}

// handlePanic 捕获并处理 panic
//
// params:
//   - f: 要执行的函数
func (h *WaitGroup) handlePanic(f func()) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			var err error
			// 判断 recoverErr 的类型
			switch e := recoverErr.(type) {
			case error:
				err = e // 如果是 error 类型，直接赋值
			default:
				err = fmt.Errorf("发生了未知错误: %v", recoverErr) // 否则封装成未知错误
			}
			h.SetError(err)
		}
	}()
	f() // 执行传入的函数
}

// Wait 等待所有 goroutine 完成并返回任何捕获的错误
//
// return:
//   - error: 如果没有错误发生，则返回 nil
func (h *WaitGroup) Wait() error {
	h.wg.Wait()         // 等待所有 goroutine 完成
	return h.GetError() // 返回捕获的错误
}
