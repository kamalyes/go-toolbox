/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:55:15
 * @FilePath: \go-toolbox\tests\priority_queue_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"context"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/queue"
	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue(t *testing.T) {
	ctx := context.Background()
	pq := queue.NewPriorityQueue() // 创建一个新的优先队列

	// 测试队列是否为空
	assert.True(t, pq.IsEmpty(), "期望队列为空")

	// 测试入队
	err := pq.Enqueue(ctx, "task1", 1)
	assert.NoError(t, err, "入队任务 'task1' 时发生错误")

	err = pq.Enqueue(ctx, "task2", 2)
	assert.NoError(t, err, "入队任务 'task2' 时发生错误")

	// 测试队列大小
	assert.Equal(t, 2, pq.Size(), "期望队列大小为 2")

	// 测试出队
	item, err := pq.Dequeue(ctx)
	assert.NoError(t, err, "出队时发生错误")
	assert.Equal(t, "task2", item, "期望出队的任务为 'task2'")

	// 测试队列大小
	assert.Equal(t, 1, pq.Size(), "期望队列大小为 1")

	// 测试再次出队
	item, err = pq.Dequeue(ctx)
	assert.NoError(t, err, "出队时发生错误")
	assert.Equal(t, "task1", item, "期望出队的任务为 'task1'")

	// 测试队列是否为空
	assert.True(t, pq.IsEmpty(), "期望所有任务出队后队列为空")

	// 测试从空队列出队
	item, err = pq.Dequeue(ctx)
	assert.Error(t, err, "期望从空队列出队时返回错误")
	assert.Nil(t, item, "期望空队列出队时返回 nil")
}

func TestEnqueueWithCancel(t *testing.T) {
	pq := queue.NewPriorityQueue() // 创建一个新的优先队列
	ctx, cancel := context.WithCancel(context.Background())

	// 测试上下文取消
	cancel() // 立即取消上下文
	err := pq.Enqueue(ctx, "task1", 1)
	assert.Error(t, err, "期望在上下文被取消时入队返回错误")
}

func TestDequeueWithCancel(t *testing.T) {
	pq := queue.NewPriorityQueue() // 创建一个新的优先队列
	ctx, cancel := context.WithCancel(context.Background())
	pq.Enqueue(ctx, "task1", 1) // 入队一个任务

	// 测试上下文取消
	cancel() // 立即取消上下文
	item, err := pq.Dequeue(ctx)
	assert.Error(t, err, "期望在上下文被取消时出队返回错误")
	assert.Nil(t, item, "期望在上下文被取消时返回 nil")
}
