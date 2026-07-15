/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-07-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-11 22:08:59
 * @FilePath: \go-toolbox\pkg\syncx\batch_processor_test.go
 * @Description:
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package syncx

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBatchProcessor_BatchFlushAtBatchSize 满 batchSize 时立即 flush
func TestBatchProcessor_BatchFlushAtBatchSize(t *testing.T) {
	var mu sync.Mutex
	var flushCalls [][]int

	p := NewBatchProcessor(100, 10, 5*time.Second, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		flushCalls = append(flushCalls, append([]int(nil), batch...))
	})
	defer p.Stop()

	for i := 0; i < 10; i++ {
		require.True(t, p.Submit(i))
	}

	// 等待 flush
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if len(flushCalls) >= 1 {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(5 * time.Millisecond)
	}

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, flushCalls, 1, "应该 flush 1 次")
	assert.Len(t, flushCalls[0], 10, "应该包含 10 个元素")
}

// TestBatchProcessor_FlushOnInterval 未满 batchSize 但超时时 flush
func TestBatchProcessor_FlushOnInterval(t *testing.T) {
	var mu sync.Mutex
	var flushCalls [][]int

	p := NewBatchProcessor(100, 100, 50*time.Millisecond, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		flushCalls = append(flushCalls, append([]int(nil), batch...))
	})
	defer p.Stop()

	for i := 0; i < 3; i++ {
		p.Submit(i)
	}

	// 等待 ticker 触发
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if len(flushCalls) >= 1 {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(5 * time.Millisecond)
	}

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, flushCalls, 1, "应该 flush 1 次")
	assert.Len(t, flushCalls[0], 3, "应该包含 3 个元素")
}

// TestBatchProcessor_StopFlushes Stop 时 flush 剩余数据
func TestBatchProcessor_StopFlushes(t *testing.T) {
	var mu sync.Mutex
	var flushCalls [][]int

	p := NewBatchProcessor(100, 100, 10*time.Second, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		flushCalls = append(flushCalls, append([]int(nil), batch...))
	})

	for i := 0; i < 5; i++ {
		p.Submit(i)
	}

	p.Stop()

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, flushCalls, 1, "Stop 应该 flush 1 次")
	assert.Len(t, flushCalls[0], 5, "应该包含 5 个元素")
}

// TestBatchProcessor_QueueFull 队列满时 Submit 返回 false
func TestBatchProcessor_QueueFull(t *testing.T) {
	// 手动创建，不启动 run goroutine（避免后台消费 channel）
	p := &BatchProcessor[int]{
		queue:         make(chan int, 5),
		flushInterval: 10 * time.Second,
		batchSize:     100,
		flushFn:       func(batch []int) {},
		stopChan:      make(chan struct{}),
		done:          make(chan struct{}),
	}

	for i := 0; i < 5; i++ {
		require.True(t, p.Submit(i), "前 5 条应该成功")
	}

	ok := p.Submit(99)
	assert.False(t, ok, "队列满时 Submit 应该返回 false")
}

// TestBatchProcessor_MultipleBatches 多批 flush
func TestBatchProcessor_MultipleBatches(t *testing.T) {
	var mu sync.Mutex
	var totalFlushed int

	p := NewBatchProcessor(1000, 50, 5*time.Second, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		totalFlushed += len(batch)
	})
	defer p.Stop()

	for i := 0; i < 150; i++ {
		p.Submit(i)
	}

	// 等待 3 次 flush（50+50+50）
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if totalFlushed >= 150 {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(5 * time.Millisecond)
	}

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 150, totalFlushed, "3 批合计应该 150 个元素")
}

// TestBatchProcessor_SubmitBlocking_WaitsForSpace 队列满时 SubmitBlocking 阻塞，腾出空位后成功写入
func TestBatchProcessor_SubmitBlocking_WaitsForSpace(t *testing.T) {
	// 手动创建，不启动 run goroutine（避免后台消费 channel）
	p := &BatchProcessor[int]{
		queue:         make(chan int, 2),
		flushInterval: 10 * time.Second,
		batchSize:     100,
		flushFn:       func(batch []int) {},
		stopChan:      make(chan struct{}),
		done:          make(chan struct{}),
	}

	// 填满队列
	require.True(t, p.Submit(1))
	require.True(t, p.Submit(2))

	// SubmitBlocking 应该阻塞（队列已满）
	resultCh := make(chan bool, 1)
	go func() {
		resultCh <- p.SubmitBlocking(context.Background(), 3)
	}()

	// 确认阻塞：短时间内不应有返回
	select {
	case <-resultCh:
		t.Fatal("队列满时 SubmitBlocking 不应该立即返回")
	case <-time.After(50 * time.Millisecond):
		// 预期阻塞中
	}

	// 腾出一个空位（FIFO，读出最先入队的 1）
	first := <-p.queue
	assert.Equal(t, 1, first)

	// 现在 SubmitBlocking 应该成功返回 true
	select {
	case ok := <-resultCh:
		require.True(t, ok, "腾位后 SubmitBlocking 应返回 true")
	case <-time.After(time.Second):
		t.Fatal("腾位后 SubmitBlocking 应及时返回")
	}

	// 队列剩余：[2, 3]，依次读出校验
	assert.Equal(t, 2, <-p.queue)
	assert.Equal(t, 3, <-p.queue, "最后写入的应该是被阻塞的 item 3")
}

// TestBatchProcessor_SubmitBlocking_ContextTimeout 队列满时 SubmitBlocking 等到 ctx 超时返回 false
func TestBatchProcessor_SubmitBlocking_ContextTimeout(t *testing.T) {
	p := &BatchProcessor[int]{
		queue:         make(chan int, 1),
		flushInterval: 10 * time.Second,
		batchSize:     100,
		flushFn:       func(batch []int) {},
		stopChan:      make(chan struct{}),
		done:          make(chan struct{}),
	}

	// 填满队列
	require.True(t, p.Submit(1))

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	ok := p.SubmitBlocking(ctx, 2)
	elapsed := time.Since(start)

	require.False(t, ok, "ctx 超时后应返回 false")
	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond, "应该等待到 ctx 超时才返回")
}
