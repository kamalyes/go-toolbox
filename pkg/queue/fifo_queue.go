/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 23:41:00
 * @FilePath: \go-toolbox\pkg\queue\fifo_queue.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package queue

import (
	"context"
	"errors"
	"sync"
)

// FIFOQueue 实现了先进先出（FIFO）的队列
// - 队列支持扩容（通过设置 autoResize 字段）
// - 具有同步控制，支持并发操作
// - 使用环形数组实现队列，避免数组元素的移动
type FIFOQueue struct {
	items      []interface{} // 队列的存储数组
	head       int           // 队列头部指针
	tail       int           // 队列尾部指针
	size       int           // 队列当前大小（元素数量）
	cap        int           // 队列的容量（最大存储数量）
	autoResize bool          // 是否自动扩容，当队列满时使用
	mu         sync.RWMutex  // 读写锁，保证并发安全
}

// NewFIFOQueue 创建并返回一个新的 FIFO 队列，
// 参数 `capacity` 是初始容量，`autoResize` 决定是否在队列满时自动扩容
func NewFIFOQueue(capacity int, autoResize bool) *FIFOQueue {
	return &FIFOQueue{
		items:      make([]interface{}, capacity), // 初始化队列的存储数组
		cap:        capacity,                      // 设置队列的容量
		autoResize: autoResize,                    // 设置是否启用自动扩容
	}
}

// Enqueue 向队列尾部添加一个元素，队列满时会根据 `autoResize` 进行扩容
// 如果上下文已取消，返回 `ctx.Err()` 错误
func (q *FIFOQueue) Enqueue(ctx context.Context, item interface{}) error {
	q.mu.Lock() // 锁定队列，防止并发修改
	defer q.mu.Unlock()

	// 如果上下文已取消，则返回错误
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 如果队列已满，则根据 autoResize 决定是否扩容
	if q.size == q.cap {
		if q.autoResize {
			// 扩容：将队列容量翻倍
			newCap := q.cap * 2
			newItems := make([]interface{}, newCap)
			// 重新排列元素，因为当前队列是环形的
			for i := 0; i < q.size; i++ {
				newItems[i] = q.items[(q.head+i)%q.cap]
			}
			// 更新队列为新的扩容数组，并重置 head 和 tail 指针
			q.items = newItems
			q.head = 0
			q.tail = q.size
			q.cap = newCap
		} else {
			// 如果未启用自动扩容，则返回错误
			return errors.New("队列已满，且未启用自动扩容")
		}
	}

	// 将元素添加到队列尾部
	q.items[q.tail] = item
	// 更新尾部指针，保持环形结构
	q.tail = (q.tail + 1) % q.cap
	// 增加队列大小
	q.size++
	return nil
}

// Dequeue 从队列头部移除并返回一个元素
// 如果队列为空，返回错误
// 如果上下文已取消，返回 `ctx.Err()` 错误
func (q *FIFOQueue) Dequeue(ctx context.Context) (interface{}, error) {
	q.mu.Lock() // 锁定队列，防止并发修改
	defer q.mu.Unlock()

	// 如果上下文已取消，则返回错误
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 如果队列为空，返回错误
	if q.size == 0 {
		return nil, errors.New("队列为空")
	}

	// 获取队列头部的元素
	item := q.items[q.head]
	// 更新队列头部指针，保持环形结构
	q.head = (q.head + 1) % q.cap
	// 减少队列大小
	q.size--
	return item, nil
}

// IsEmpty 检查队列是否为空
func (q *FIFOQueue) IsEmpty() bool {
	q.mu.RLock() // 只读锁定，避免修改队列时阻塞
	defer q.mu.RUnlock()
	return q.size == 0
}

// Size 返回队列中的元素数量
func (q *FIFOQueue) Size() int {
	q.mu.RLock() // 只读锁定
	defer q.mu.RUnlock()
	return q.size
}
