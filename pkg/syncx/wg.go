/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 13:55:22
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-08 13:55:22
 * @FilePath: \go-toolbox\pkg\syncx\wg.go
 * @Description: 自定义的 WaitGroup 结构体，封装了 sync.WaitGroup，用于管理并发操作的等待。
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import "sync"

// WaitGroupWithMutex 结构体封装了 sync.WaitGroup
type WaitGroupWithMutex struct {
	wg sync.WaitGroup // 使用 sync.WaitGroup 进行并发控制
}

// Add 方法增加等待组的计数
func (w *WaitGroupWithMutex) Add(delta int) {
	w.wg.Add(delta) // 增加 delta，表示有新的 goroutine 开始
}

// Done 方法减少等待组的计数
func (w *WaitGroupWithMutex) Done() {
	w.wg.Done() // 表示一个 goroutine 完成
}

// Wait 方法阻塞直到等待组计数为零
func (w *WaitGroupWithMutex) Wait() {
	w.wg.Wait() // 等待所有 goroutine 完成
}
