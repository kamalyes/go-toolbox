/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 11:16:12
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
type FIFOQueue struct {
	items      []interface{} // 队列的存储数组
	head       int           // 队列头部指针
	tail       int           // 队列尾部指针
	size       int           // 队列当前大小
	cap        int           // 队列的容量
	autoResize bool          // 是否自动扩容
	mu         sync.RWMutex  // 读写锁，保证并发安全
	cond       *sync.Cond    // 条件变量，用于等待和通知
}

// NewFIFOQueue 创建并返回一个新的 FIFO 队列
func NewFIFOQueue(capacity int, autoResize bool) *FIFOQueue {
	q := &FIFOQueue{
		items:      make([]interface{}, capacity),
		cap:        capacity,
		autoResize: autoResize,
	}
	q.cond = sync.NewCond(&q.mu) // 初始化条件变量
	return q
}

// Enqueue 向队列尾部添加一个元素
func (q *FIFOQueue) Enqueue(ctx context.Context, item interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 检查上下文是否已取消
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
			for i := 0; i < q.size; i++ {
				newItems[i] = q.items[(q.head+i)%q.cap]
			}
			q.items = newItems
			q.head = 0
			q.tail = q.size
			q.cap = newCap
		} else {
			return errors.New("队列已满，且未启用自动扩容")
		}
	}

	// 将元素添加到队列尾部
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.cap
	q.size++

	// 通知其他等待的 goroutine
	q.cond.Signal()
	return nil
}

// Dequeue 从队列头部移除并返回一个元素
func (q *FIFOQueue) Dequeue(ctx context.Context) (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 如果队列为空，等待直到有元素可用
	if q.size == 0 {
		return nil, errors.New("queue is empty") // 立即返回错误
	}

	// 获取队列头部的元素
	item := q.items[q.head]
	q.head = (q.head + 1) % q.cap
	q.size--

	return item, nil
}

// IsEmpty 检查队列是否为空
func (q *FIFOQueue) IsEmpty() bool {
	q.mu.RLock()         // 锁定操作，确保线程安全
	defer q.mu.RUnlock() // 解锁操作
	return q.size == 0
}

// Size 返回队列中的元素数量
func (q *FIFOQueue) Size() int {
	q.mu.RLock()         // 锁定操作，确保线程安全
	defer q.mu.RUnlock() // 解锁操作
	return q.size
}
