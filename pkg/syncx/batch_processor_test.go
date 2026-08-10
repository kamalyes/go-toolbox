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

// ============================================================================
// WithClone 测试
// ============================================================================

// TestBatchProcessor_WithClone verifies that WithClone clones items on Submit,
// protecting against caller modification after submission
func TestBatchProcessor_WithClone(t *testing.T) {
	type item struct{ value int }
	var mu sync.Mutex
	var flushed []item

	p := NewBatchProcessor(100, 100, 50*time.Millisecond, func(batch []item) {
		mu.Lock()
		defer mu.Unlock()
		flushed = append(flushed, batch...)
	},
		WithBatchProcessorClone(func(i item) item { return item{value: i.value} }),
	)

	original := item{value: 42}
	require.True(t, p.Submit(original))

	// Modify original after Submit — clone should protect the queued copy
	original.value = 999

	// Wait for interval flush
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if len(flushed) > 0 {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(5 * time.Millisecond)
	}

	p.Stop()

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, flushed, 1, "应该 flush 1 条")
	assert.Equal(t, 42, flushed[0].value, "Clone 应保护数据不被调用方修改影响")
}

// TestBatchProcessor_WithClone_SubmitBlocking verifies clone is applied in SubmitBlocking
func TestBatchProcessor_WithClone_SubmitBlocking(t *testing.T) {
	type item struct{ value int }
	var mu sync.Mutex
	var flushed []item

	cloneCount := 0
	p := NewBatchProcessor(100, 100, 50*time.Millisecond, func(batch []item) {
		mu.Lock()
		defer mu.Unlock()
		flushed = append(flushed, batch...)
	},
		WithBatchProcessorClone(func(i item) item {
			cloneCount++
			return item{value: i.value}
		}),
	)

	require.True(t, p.SubmitBlocking(context.Background(), item{value: 7}))

	// Wait for interval flush
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if len(flushed) > 0 {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(5 * time.Millisecond)
	}

	p.Stop()

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 1, cloneCount, "SubmitBlocking 也应触发 Clone")
	assert.Len(t, flushed, 1)
	assert.Equal(t, 7, flushed[0].value)
}

// ============================================================================
// WithPanicHandler / safeFlush 测试
// ============================================================================

// TestBatchProcessor_WithPanicHandler verifies that a flush panic does not crash
// the worker, the panic handler is called, and subsequent flushes succeed
func TestBatchProcessor_WithPanicHandler(t *testing.T) {
	var mu sync.Mutex
	panicCount := 0
	flushAfterPanic := false
	callCount := 0

	handler := func(r any) {
		mu.Lock()
		panicCount++
		mu.Unlock()
	}

	p := NewBatchProcessor(100, 1, 10*time.Second, func(batch []int) {
		mu.Lock()
		callCount++
		shouldPanic := callCount == 1
		mu.Unlock()

		if shouldPanic {
			panic("intentional flush panic")
		}
		// After panic, mark success
		mu.Lock()
		flushAfterPanic = true
		mu.Unlock()
	},
		WithBatchProcessorPanicHandler[int](handler),
	)
	defer p.Stop()

	// First Submit triggers flush (batchSize=1) → panic
	require.True(t, p.Submit(1))

	// Wait for panic recovery
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if panicCount > 0 {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(2 * time.Millisecond)
	}

	mu.Lock()
	require.Equal(t, 1, panicCount, "panic handler 应被调用 1 次")
	mu.Unlock()

	// Second Submit → flush should succeed (worker still alive)
	require.True(t, p.Submit(2))

	deadline = time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if flushAfterPanic {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(2 * time.Millisecond)
	}

	mu.Lock()
	assert.True(t, flushAfterPanic, "panic 后 worker 应继续处理后续 flush")
	mu.Unlock()
}

// TestBatchProcessor_PanicRecovery_Default verifies that without a panic handler,
// flush panic is silently recovered and the worker continues
func TestBatchProcessor_PanicRecovery_Default(t *testing.T) {
	var mu sync.Mutex
	callCount := 0
	survived := false

	p := NewBatchProcessor(100, 1, 10*time.Second, func(batch []int) {
		mu.Lock()
		callCount++
		shouldPanic := callCount == 1
		mu.Unlock()

		if shouldPanic {
			panic("no handler panic")
		}
		mu.Lock()
		survived = true
		mu.Unlock()
	})
	defer p.Stop()

	// First flush panics
	p.Submit(1)

	// Wait for panic recovery
	time.Sleep(50 * time.Millisecond)

	// Second flush should succeed
	p.Submit(2)

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		if survived {
			mu.Unlock()
			break
		}
		mu.Unlock()
		time.Sleep(2 * time.Millisecond)
	}

	mu.Lock()
	assert.True(t, survived, "无 handler 时也应静默恢复，worker 继续工作")
	mu.Unlock()
}

// ============================================================================
// DroppedCount 测试
// ============================================================================

// TestBatchProcessor_DroppedCount verifies that dropped items are counted
func TestBatchProcessor_DroppedCount(t *testing.T) {
	// Manually create without starting run goroutine
	p := &BatchProcessor[int]{
		queue:         make(chan int, 3),
		flushInterval: 10 * time.Second,
		batchSize:     100,
		flushFn:       func(batch []int) {},
		stopChan:      make(chan struct{}),
		done:          make(chan struct{}),
	}

	// Fill queue
	for i := 0; i < 3; i++ {
		require.True(t, p.Submit(i))
	}

	// Submit 5 more — all should be dropped
	for i := 0; i < 5; i++ {
		require.False(t, p.Submit(99))
	}

	assert.Equal(t, int64(5), p.DroppedCount(), "应累计 5 次丢弃")
}

// TestBatchProcessor_DroppedCount_SubmitBlocking verifies drop count on ctx timeout
func TestBatchProcessor_DroppedCount_SubmitBlocking(t *testing.T) {
	p := &BatchProcessor[int]{
		queue:         make(chan int, 1),
		flushInterval: 10 * time.Second,
		batchSize:     100,
		flushFn:       func(batch []int) {},
		stopChan:      make(chan struct{}),
		done:          make(chan struct{}),
	}

	require.True(t, p.Submit(1)) // fill queue

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	require.False(t, p.SubmitBlocking(ctx, 2)) // timeout → dropped

	assert.Equal(t, int64(1), p.DroppedCount(), "SubmitBlocking 超时也应计入 DroppedCount")
}

// ============================================================================
// WithName 测试
// ============================================================================

// TestBatchProcessor_WithName verifies that the name is set correctly
func TestBatchProcessor_WithName(t *testing.T) {
	p := NewBatchProcessor(10, 10, time.Second, func(batch []int) {},
		WithBatchProcessorName[int]("my-batcher"),
	)
	defer p.Stop()

	assert.Equal(t, "my-batcher", p.Name())
}

// TestBatchProcessor_Name_Empty verifies default name is empty
func TestBatchProcessor_Name_Empty(t *testing.T) {
	p := NewBatchProcessor(10, 10, time.Second, func(batch []int) {})
	defer p.Stop()

	assert.Equal(t, "", p.Name())
}

// ============================================================================
// 基准测试
// ============================================================================

// BenchmarkBatchProcessor_Submit measures Submit throughput (non-blocking, no clone)
func BenchmarkBatchProcessor_Submit(b *testing.B) {
	p := NewBatchProcessor(65536, 1000, 50*time.Millisecond, func(batch []int) {})
	defer p.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Submit(1)
		}
	})
}

// BenchmarkBatchProcessor_Submit_WithClone measures Submit with clone overhead
func BenchmarkBatchProcessor_Submit_WithClone(b *testing.B) {
	type item struct{ value int }
	p := NewBatchProcessor(65536, 1000, 50*time.Millisecond, func(batch []item) {},
		WithBatchProcessorClone(func(i item) item { return item{value: i.value} }),
	)
	defer p.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Submit(item{value: 1})
		}
	})
}

// BenchmarkBatchProcessor_SubmitBlocking measures SubmitBlocking throughput
func BenchmarkBatchProcessor_SubmitBlocking(b *testing.B) {
	p := NewBatchProcessor(65536, 1000, 50*time.Millisecond, func(batch []int) {})
	defer p.Stop()

	ctx := context.Background()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.SubmitBlocking(ctx, 1)
		}
	})
}
