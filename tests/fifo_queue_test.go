/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-21 18:55:55
 * @FilePath: \go-toolbox\tests\fifo_queue_test.go
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
	// 创建一个新的 FIFOQueue，初始容量为 2，启用自动扩容
	q := queue.NewFIFOQueue(2, true)

	// 测试入队
	err := q.Enqueue(context.Background(), "item1")
	assert.NoError(t, err, "应该没有错误")

	err = q.Enqueue(context.Background(), "item2")
	assert.NoError(t, err, "应该没有错误")

	// 测试队列大小
	assert.Equal(t, 2, q.Size(), "队列大小应该为 2")

	// 测试扩容
	err = q.Enqueue(context.Background(), "item3") // 触发扩容
	assert.NoError(t, err, "应该没有错误")
	assert.Equal(t, 3, q.Size(), "队列大小应该为 3")
	assert.Equal(t, 4, q.Capacity(), "队列容量应该扩展到 4")

	// 测试出队
	item, err := q.Dequeue(context.Background())
	assert.NoError(t, err, "应该没有错误")
	assert.Equal(t, "item1", item, "出队的元素应该是 item1")

	item, err = q.Dequeue(context.Background())
	assert.NoError(t, err, "应该没有错误")
	assert.Equal(t, "item2", item, "出队的元素应该是 item2")

	// 测试队列是否为空
	assert.False(t, q.IsEmpty(), "队列不应该为空")

	// 再次出队
	item, err = q.Dequeue(context.Background())
	assert.NoError(t, err, "应该没有错误")
	assert.Equal(t, "item3", item, "出队的元素应该是 item3")

	// 测试队列是否为空
	assert.True(t, q.IsEmpty(), "队列应该为空")

	// 测试缩容
	q.Enqueue(context.Background(), "item4")
	q.Enqueue(context.Background(), "item5")
	q.Enqueue(context.Background(), "item6") // 触发扩容
	q.Dequeue(context.Background())          // 移除一个元素
	assert.Equal(t, 4, q.Capacity(), "队列容量应该为 4")
	q.Dequeue(context.Background()) // 移除一个元素
	assert.Equal(t, 2, q.Capacity(), "队列容量应该为 2")
	q.Dequeue(context.Background()) // 移除一个元素
	assert.Equal(t, 1, q.Capacity(), "队列容量应该为 1")

	// 进行缩容，检查最低缩容
	q.Dequeue(context.Background()) // 移除一个元素，触发缩容
	assert.Equal(t, 1, q.Capacity(), "队列容量应该缩容到 1")

	// 测试最小容量限制
	for i := 0; i < 10; i++ {
		q.Enqueue(context.Background(), i)
	}
	for i := 0; i < 9; i++ {
		q.Dequeue(context.Background())
	}
	assert.Equal(t, 2, q.Capacity(), "队列容量应该保持在最小容量限制")
}

func TestFIFOQueue_ContextCancel(t *testing.T) {
	// 创建一个新的 FIFOQueue
	q := queue.NewFIFOQueue(2, true)

	// 测试入队
	err := q.Enqueue(context.Background(), "item1")
	assert.NoError(t, err, "应该没有错误")

	// 测试出队时上下文取消
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 先入队一个元素
	err = q.Enqueue(ctx, "item2")
	assert.NoError(t, err, "应该没有错误")

	// 取消上下文
	cancel()

	// 测试入队时上下文取消
	err = q.Enqueue(ctx, "item3")
	assert.Error(t, err, "应该返回上下文取消错误")

	// 测试出队时上下文取消
	_, err = q.Dequeue(ctx)
	assert.Error(t, err, "应该返回上下文取消错误")
}

func TestFIFOQueue_Concurrent(t *testing.T) {
	q := queue.NewFIFOQueue(2, true)

	// 使用 WaitGroup 来等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 启动多个 goroutine 进行入队操作
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(item int) {
			defer wg.Done()
			err := q.Enqueue(context.Background(), item)
			assert.NoError(t, err, "应该没有错误")
		}(i)
	}

	wg.Wait() // 等待所有入队操作完成

	assert.Equal(t, 10, q.Size(), "队列大小应该为 10")

	// 启动多个 goroutine 进行出队操作
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := q.Dequeue(context.Background())
			assert.NoError(t, err, "应该没有错误")
		}()
	}

	wg.Wait() // 等待所有出队操作完成

	assert.True(t, q.IsEmpty(), "队列应该为空")
}

func TestFIFOQueue_SetGrowthFactor(t *testing.T) {
	q := queue.NewFIFOQueue(2, true)

	// 设置新的扩容因子
	q.SetGrowthFactor(1.5)

	// 确保扩容因子已更新
	assert.Equal(t, 1.5, q.GrowthFactor(), "扩容因子应该是 1.5")
}

func TestFIFOQueue_SetShrinkFactor(t *testing.T) {
	q := queue.NewFIFOQueue(2, true)

	// 设置新的缩容因子
	q.SetShrinkFactor(0.75)

	// 确保缩容因子已更新
	assert.Equal(t, 0.75, q.ShrinkFactor(), "缩容因子应该是 0.75")
}

func TestFIFOQueue_MinCapacity(t *testing.T) {
	q := queue.NewFIFOQueue(2, true)

	// 确保最小容量限制正常工作
	assert.Equal(t, 1, q.MinCapacity(), "最小容量应该是 1")

	// 尝试缩容到小于最小容量
	for i := 0; i < 10; i++ {
		q.Enqueue(context.Background(), i)
	}
	for i := 0; i < 9; i++ {
		q.Dequeue(context.Background())
	}
	assert.Equal(t, 2, q.Capacity(), "队列容量应该保持在最小容量限制")
}
