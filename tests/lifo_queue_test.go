/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:55:15
 * @FilePath: \go-toolbox\tests\lifo_queue_test.go
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

// 测试数据和文案
var (
	testItems        = []int{1, 2, 3}
	testItemsSize    = len(testItems)
	contextCancelErr = "上下文取消后，操作应失败"
	emptyQueueErr    = "队列为空"
	testItemForSize  = 10
)

// 公共错误消息
const (
	msgNewLIFOQueueEmpty                  = "NewLIFOQueue 应该创建一个非空的 LIFOQueue 实例"
	msgQueueEmpty                         = "新创建的 LIFOQueue 应该是空的"
	msgEnqueueNoError                     = "入队操作不应产生错误"
	msgDequeueNoError                     = "出队操作不应产生错误"
	msgDequeueLastInFirstOut              = "出队应该返回最后入队的元素（LIFO顺序）"
	msgQueueShouldBeEmpty                 = "所有元素出队后，队列应该为空"
	msgEmptyDequeueErr                    = "从空队列出队应产生错误"
	msgQueueSizeAfterEnqueue              = "入队两个元素后，队列大小应为2"
	msgQueueSizeAfterDequeue              = "出队一个元素后，队列大小应为1"
	msgConcurrencySizeMatch               = "并发入队后，队列大小应与入队元素数量相匹配"
	msgQueueShouldBeEmptyAfterConcurrency = "并发出队后，队列应为空"
)

// 测试创建新的 LIFO 队列
func TestNewLIFOQueue(t *testing.T) {
	q := queue.NewLIFOQueue()
	assert.NotNil(t, q, msgNewLIFOQueueEmpty)
	assert.True(t, q.IsEmpty(), msgQueueEmpty)
}

// 测试 LIFO 队列的入队和出队功能
func TestLIFOQueueEnqueueDequeue(t *testing.T) {
	q := queue.NewLIFOQueue()
	assert := assert.New(t)

	// 测试入队
	for _, item := range testItems {
		err := q.Enqueue(context.Background(), item)
		assert.NoError(err, msgEnqueueNoError)
	}

	// 测试出队
	for i := len(testItems) - 1; i >= 0; i-- {
		result, err := q.Dequeue(context.Background())
		assert.NoError(err, msgDequeueNoError)
		assert.Equal(testItems[i], result, msgDequeueLastInFirstOut)
	}

	// 确认队列为空
	assert.True(q.IsEmpty(), msgQueueShouldBeEmpty)
}

// 测试从空队列中出队
func TestLIFOQueueEmptyDequeue(t *testing.T) {
	q := queue.NewLIFOQueue()
	_, err := q.Dequeue(context.Background())
	assert.Error(t, err, msgEmptyDequeueErr)
	assert.Equal(t, emptyQueueErr, err.Error(), "错误信息应指示队列为空")
}

// 测试上下文取消对队列操作的影响
func TestLIFOQueueContextCancellation(t *testing.T) {
	q := queue.NewLIFOQueue()
	ctx, cancel := context.WithCancel(context.Background())

	// 在操作前取消上下文
	cancel()

	err := q.Enqueue(ctx, testItemForSize)
	assert.Error(t, err, contextCancelErr)

	_, err = q.Dequeue(ctx)
	assert.Error(t, err, contextCancelErr)
}

// 测试队列大小
func TestLIFOQueueSize(t *testing.T) {
	q := queue.NewLIFOQueue()
	assert := assert.New(t)

	assert.Equal(0, q.Size(), "新队列的大小应为0")

	q.Enqueue(context.Background(), 1)
	q.Enqueue(context.Background(), 2)
	assert.Equal(2, q.Size(), msgQueueSizeAfterEnqueue)

	q.Dequeue(context.Background())
	assert.Equal(1, q.Size(), msgQueueSizeAfterDequeue)
}

// 测试队列的并发安全性
func TestLIFOQueueConcurrency(t *testing.T) {
	q := queue.NewLIFOQueue()
	assert := assert.New(t)

	// 并发入队测试
	var wg sync.WaitGroup
	for _, item := range testItems {
		wg.Add(1)
		go func(item int) {
			defer wg.Done()
			q.Enqueue(context.Background(), item)
		}(item)
	}
	wg.Wait() // 等待所有入队操作完成

	assert.Equal(testItemsSize, q.Size(), msgConcurrencySizeMatch)

	// 并发出队测试
	wg = sync.WaitGroup{} // 重置 WaitGroup
	for range testItems {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q.Dequeue(context.Background())
		}()
	}
	wg.Wait() // 等待所有出队操作完成

	assert.True(q.IsEmpty(), msgQueueShouldBeEmptyAfterConcurrency)
}
