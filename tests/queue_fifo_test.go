/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 22:49:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 23:58:47
 * @FilePath: \go-toolbox\tests\queue_fifo_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/queue"
	"github.com/stretchr/testify/assert"
)

// 公共数据定义
var (
	commonItems     = []int{1, 2, 3, 4, 5, 10, 20, 30}
	autoResizeIndex = 3                               // 自动扩容触发的索引
	enqueueItems    = commonItems[:autoResizeIndex+1] // 需要入队的元素
	dequeueOrder    = []int{1, 2}                     // 出队顺序验证
	contextDeadline = "context deadline exceeded"     // 超时错误信息
	queueEmptyErr   = "队列为空"                          // 空队列错误信息
)

// 公共错误消息
const (
	msgQueueAutoResize = "队列应已自动扩容并包含预期数量的元素"
)

// 测试 FIFO 队列的入队和出队功能
func TestFIFOQueueEnqueueDequeue(t *testing.T) {
	q := queue.NewFIFOQueue(3, true) // 初始容量3，启用自动扩容
	assert := assert.New(t)

	// 测试入队
	for _, item := range enqueueItems {
		err := q.Enqueue(context.Background(), item)
		assert.NoError(err, msgEnqueueNoError)
	}

	// 测试出队
	for _, expected := range dequeueOrder {
		item, err := q.Dequeue(context.Background())
		assert.NoError(err, msgDequeueNoError)
		assert.Equal(expected, item, "出队元素应符合预期")
	}

	assert.Equal(2, q.Size(), msgQueueAutoResize)
}

// 测试 FIFO 队列是否为空
func TestIsEmpty(t *testing.T) {
	assert := assert.New(t)
	q := queue.NewFIFOQueue(3, true)

	assert.True(q.IsEmpty(), "新创建的队列应为空")

	q.Enqueue(context.Background(), 1)
	assert.False(q.IsEmpty(), "入队后队列不应为空")

	q.Dequeue(context.Background())
	assert.True(q.IsEmpty(), "出队所有元素后队列应为空")
}

// 测试从空队列中出队
func TestDequeueFromEmptyQueue(t *testing.T) {
	assert := assert.New(t)
	q := queue.NewFIFOQueue(2, false)

	_, err := q.Dequeue(context.Background())
	assert.Error(err, "从空队列出队应返回错误")
	assert.Equal(t, "队列为空", err.Error(), "错误信息应指示队列为空")
}

// 测试上下文取消对队列操作的影响
func TestContextCancellation(t *testing.T) {
	assert := assert.New(t)
	q := queue.NewFIFOQueue(3, true)
	ctx, cancel := context.WithCancel(context.Background())
	// 在操作前取消上下文
	cancel()

	err := q.Enqueue(ctx, commonItems[autoResizeIndex+1])
	assert.Equal(t, contextDeadline, err.Error())

	_, err = q.Dequeue(ctx)
	assert.Equal(t, contextDeadline, err.Error())
}

// 测试队列大小
func TestQueueSize(t *testing.T) {
	assert := assert.New(t)
	q := queue.NewFIFOQueue(3, true)

	for _, item := range commonItems[5:] {
		err := q.Enqueue(context.Background(), item)
		assert.NoError(err)
	}

	assert.Equal(3, q.Size(), msgQueueSizeAfterEnqueue)
}

// 测试队列的出队顺序
func TestDequeueOrder(t *testing.T) {
	assert := assert.New(t)
	q := queue.NewFIFOQueue(3, true)

	for _, item := range enqueueItems {
		q.Enqueue(context.Background(), item)
	}

	expectedOrder := []int{1, 2, 3, 4}
	for _, expected := range expectedOrder {
		item, err := q.Dequeue(context.Background())
		assert.NoError(err, msgDequeueNoError)
		assert.Equal(expected, item, "出队元素应符合预期，遵循FIFO原则")
	}
}

// 测试 FIFO 队列的并发入队和出队操作
func TestFIFOQueueConcurrency(t *testing.T) {
	q := queue.NewFIFOQueue(5, true)
	assert := assert.New(t)

	items := make([]int, 100)
	for i := range items {
		items[i] = i
	}

	// 并发入队测试
	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go func(item int) {
			defer wg.Done()
			err := q.Enqueue(context.Background(), item)
			assert.NoError(err, msgEnqueueNoError)
		}(item)
	}
	wg.Wait()

	// 检查队列的大小
	assert.Equal(len(items), q.Size(), "并发入队后，队列大小应与入队元素数量相匹配")

	// 并发出队测试
	wg = sync.WaitGroup{}
	for range items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := q.Dequeue(context.Background())
			assert.NoError(err, msgDequeueNoError)
		}()
	}
	wg.Wait()

	assert.True(q.IsEmpty(), msgQueueShouldBeEmpty)
}

// 测试 FIFO 队列的自动扩容功能
func TestFIFOQueueAutoResize(t *testing.T) {
	q := queue.NewFIFOQueue(2, true) // 初始容量2，启用自动扩容
	assert := assert.New(t)

	// 入队超出初始容量，应该自动扩容
	for i := 0; i < 3; i++ {
		err := q.Enqueue(context.Background(), i)
		assert.NoError(err, "入队操作不应产生错误")
	}

	assert.Equal(3, q.Size(), "队列大小应为3，表示已自动扩容")
}

// 测试 FIFO 队列的上下文取消
func TestFIFOQueueContextCancellation(t *testing.T) {
	q := queue.NewFIFOQueue(5, true)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消上下文

	// 测试入队时，已取消的上下文
	err := q.Enqueue(ctx, 10)
	assert.Error(t, err, "上下文取消后，入队操作应失败")

	// 测试出队时，已取消的上下文
	_, err = q.Dequeue(ctx)
	assert.Error(t, err, "上下文取消后，出队操作应失败")
}

// 测试 FIFO 队列为空的出队操作
func TestFIFOQueueEmptyDequeue(t *testing.T) {
	q := queue.NewFIFOQueue(3, true)
	_, err := q.Dequeue(context.Background())
	assert.Error(t, err, "从空队列出队应产生错误")
	assert.Equal(t, queueEmptyErr, err.Error(), "错误信息应指示队列为空")
}
