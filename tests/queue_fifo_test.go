/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:55:15
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

func TestFIFOQueue(t *testing.T) {
	ctx := context.Background()
	queue := queue.NewFIFOQueue(2, true) // 创建一个初始容量为2的队列，启用自动扩容

	// 测试入队
	assert.NoError(t, queue.Enqueue(ctx, 1))
	assert.NoError(t, queue.Enqueue(ctx, 2))

	// 测试自动扩容
	assert.NoError(t, queue.Enqueue(ctx, 3))
	assert.Equal(t, 3, queue.Size())

	// 测试出队
	item, err := queue.Dequeue(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, item)

	// 测试队列大小
	assert.Equal(t, 2, queue.Size())

	// 测试出队直到空
	item, err = queue.Dequeue(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 2, item)

	item, err = queue.Dequeue(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 3, item)

	// 测试队列为空
	assert.True(t, queue.IsEmpty())

	// 测试出队错误
	_, err = queue.Dequeue(ctx)
	assert.Error(t, err)

	// 测试上下文取消
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消上下文
	assert.Error(t, queue.Enqueue(ctx, 4))
}

func TestFIFOQueueConcurrency(t *testing.T) {
	queue := queue.NewFIFOQueue(5, true)
	ctx := context.Background()
	var wg sync.WaitGroup // 主要改动：使用 WaitGroup 等待所有 goroutine 完成

	// 启动多个 goroutine 测试并发入队
	for i := 0; i < 10; i++ {
		wg.Add(1) // 主要改动：每次启动 goroutine 前增加计数
		go func(i int) {
			defer wg.Done() // 主要改动：确保 goroutine 完成时减少计数
			assert.NoError(t, queue.Enqueue(ctx, i))
		}(i)
	}

	wg.Wait() // 主要改动：等待所有入队操作完成

	// 检查队列大小
	assert.GreaterOrEqual(t, queue.Size(), 10)

	// 测试并发出队
	for i := 0; i < 10; i++ {
		wg.Add(1) // 主要改动：每次启动 goroutine 前增加计数
		go func() {
			defer wg.Done() // 主要改动：确保 goroutine 完成时减少计数
			_, err := queue.Dequeue(ctx)
			assert.NoError(t, err)
		}()
	}

	wg.Wait() // 主要改动：等待所有出队操作完成

	// 最终检查队列是否为空
	assert.True(t, queue.IsEmpty())
}
