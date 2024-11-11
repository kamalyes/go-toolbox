/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 23:44:06
 * @FilePath: \go-toolbox\pkg\queue\lifo_queue.go
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

// LIFOQueue 实现了 LIFO 队列（栈）
type LIFOQueue struct {
	items []interface{} // 用切片存储队列元素
	mu    sync.RWMutex  // 读写锁，保证并发安全
}

// NewLIFOQueue 创建一个新的 LIFO 队列（栈）
func NewLIFOQueue() *LIFOQueue {
	return &LIFOQueue{
		items: []interface{}{}, // 初始化一个空切片
	}
}

// Enqueue 将元素添加到队列中（栈的压栈操作）
func (s *LIFOQueue) Enqueue(ctx context.Context, item interface{}) error {
	s.mu.Lock()         // 锁定操作，确保线程安全
	defer s.mu.Unlock() // 解锁操作

	select {
	case <-ctx.Done(): // 如果上下文被取消，返回取消错误
		return ctx.Err()
	default:
	}

	// 将元素添加到队列末尾（栈顶）
	s.items = append(s.items, item)
	return nil
}

// Dequeue 从队列中取出元素（栈的弹栈操作）
func (s *LIFOQueue) Dequeue(ctx context.Context) (interface{}, error) {
	s.mu.Lock()         // 锁定操作，确保线程安全
	defer s.mu.Unlock() // 解锁操作

	select {
	case <-ctx.Done(): // 如果上下文被取消，返回取消错误
		return nil, ctx.Err()
	default:
	}

	if len(s.items) == 0 { // 如果队列为空，返回错误
		return nil, errors.New("队列为空")
	}

	// 取出栈顶元素并删除
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

// IsEmpty 判断队列是否为空
func (s *LIFOQueue) IsEmpty() bool {
	s.mu.RLock()             // 锁定操作，确保线程安全
	defer s.mu.RUnlock()     // 解锁操作
	return len(s.items) == 0 // 如果队列为空，返回 true
}

// Size 返回队列的大小
func (s *LIFOQueue) Size() int {
	s.mu.RLock()         // 锁定操作，确保线程安全
	defer s.mu.RUnlock() // 解锁操作
	return len(s.items)  // 返回队列中元素的数量
}
