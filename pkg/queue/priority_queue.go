/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 11:26:09
 * @FilePath: \go-toolbox\pkg\queue\priority_queue.go
 * @Description: 优先队列实现
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package queue

import (
	"container/heap"
	"context"
	"errors"
	"sync"
)

// Item 定义了优先队列中的元素
type Item struct {
	value    interface{} // 队列元素的值
	priority int         // 优先级，数字越大优先级越高
}

// PriorityQueue 实现了优先队列，实际上是一个堆
type PriorityQueue struct {
	items []*Item
	mu    sync.RWMutex // 使用读写锁以支持并发读写
}

// Len 返回队列的长度
func (pq *PriorityQueue) Len() int {
	return len(pq.items) // 直接返回长度
}

// Less 判断队列中第 i 和第 j 个元素的优先级
func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.items[i].priority > pq.items[j].priority // 高优先级在前
}

// Swap 交换队列中第 i 和第 j 个元素的位置
func (pq *PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

// Push 将一个元素添加到队列中（堆）
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	pq.items = append(pq.items, item)
}

// Pop 从队列中移除并返回最优先的元素（堆顶元素）
func (pq *PriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.items = old[0 : n-1]
	return item
}

// NewPriorityQueue 创建并返回一个新的优先队列
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{}
	heap.Init(pq)
	return pq
}

// Enqueue 将一个元素添加到优先队列中，支持上下文取消
func (pq *PriorityQueue) Enqueue(ctx context.Context, item interface{}, priority int) error {
	select {
	case <-ctx.Done(): // 如果上下文被取消，返回取消错误
		return ctx.Err()
	default:
	}

	pq.mu.Lock()         // 锁定操作，确保线程安全
	defer pq.mu.Unlock() // 解锁操作

	heap.Push(pq, &Item{value: item, priority: priority})
	return nil
}

// Dequeue 从优先队列中取出最优先的元素，支持上下文取消
func (pq *PriorityQueue) Dequeue(ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done(): // 如果上下文被取消，返回取消错误
		return nil, ctx.Err()
	default:
	}

	pq.mu.Lock()         // 锁定操作，确保线程安全
	defer pq.mu.Unlock() // 解锁操作

	if pq.Len() == 0 {
		return nil, errors.New("队列为空")
	}
	return heap.Pop(pq).(*Item).value, nil
}

// IsEmpty 判断优先队列是否为空
func (pq *PriorityQueue) IsEmpty() bool {
	pq.mu.RLock()         // 读锁定操作
	defer pq.mu.RUnlock() // 解锁操作
	return pq.Len() == 0
}

// Size 返回优先队列的大小
func (pq *PriorityQueue) Size() int {
	pq.mu.RLock()         // 读锁定操作
	defer pq.mu.RUnlock() // 解锁操作
	return pq.Len()
}
