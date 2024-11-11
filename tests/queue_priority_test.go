/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 23:38:55
 * @FilePath: \go-toolbox\tests\queue_priority_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"context"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/queue"
	"github.com/stretchr/testify/assert"
)

// TestPriorityQueue 测试优先队列的基本功能
func TestPriorityQueue(t *testing.T) {
	// 创建一个新的优先队列
	pq := queue.NewPriorityQueue()

	// 使用上下文设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Test 1: Enqueue操作
	err := pq.Enqueue(ctx, "Task 1", 1)
	assert.NoError(t, err, "Enqueue should succeed when adding Task 1 with priority 1")
	err = pq.Enqueue(ctx, "Task 2", 3)
	assert.NoError(t, err, "Enqueue should succeed when adding Task 2 with priority 3")
	err = pq.Enqueue(ctx, "Task 3", 2)
	assert.NoError(t, err, "Enqueue should succeed when adding Task 3 with priority 2")

	// Test 2: Size 操作
	assert.Equal(t, 3, pq.Size(), "Size should return the number of elements in the queue")

	// Test 3: Dequeue 操作（按优先级顺序）
	item, err := pq.Dequeue(ctx)
	assert.NoError(t, err, "Dequeue should succeed when the queue is not empty")
	assert.Equal(t, "Task 2", item, "Task 2 should be dequeued first due to its highest priority")

	item, err = pq.Dequeue(ctx)
	assert.NoError(t, err, "Dequeue should succeed when the queue is not empty")
	assert.Equal(t, "Task 3", item, "Task 3 should be dequeued second")

	item, err = pq.Dequeue(ctx)
	assert.NoError(t, err, "Dequeue should succeed when the queue is not empty")
	assert.Equal(t, "Task 1", item, "Task 1 should be dequeued last")

	// Test 4: Queue should be empty after all dequeue operations
	assert.True(t, pq.IsEmpty(), "Queue should be empty after all tasks are dequeued")
	assert.Equal(t, 0, pq.Size(), "Queue size should be 0 after all elements are dequeued")
}

// TestPriorityQueue_EmptyQueue 测试空队列的Dequeue操作
func TestPriorityQueue_EmptyQueue(t *testing.T) {
	pq := queue.NewPriorityQueue()

	// 使用上下文设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Test: Dequeue 操作 on empty queue
	item, err := pq.Dequeue(ctx)
	assert.Error(t, err, "Dequeue should return an error when the queue is empty")
	assert.Nil(t, item, "Dequeued item should be nil when the queue is empty")
}

// TestPriorityQueue_Timeout 测试上下文超时的情况
func TestPriorityQueue_Timeout(t *testing.T) {
	pq := queue.NewPriorityQueue()
	ctx, cancel := context.WithCancel(context.Background())

	// 在操作前取消上下文
	cancel()

	// 添加一个任务
	err := pq.Enqueue(ctx, "Task 1", 1)
	assert.Error(t, err, contextCancelErr)

	_, err = pq.Dequeue(ctx)
	assert.Error(t, err, contextCancelErr)
}
