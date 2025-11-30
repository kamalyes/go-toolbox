/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-20 10:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 10:00:00
 * @FilePath: \go-toolbox\pkg\queue\bounded_queue_test.go
 * @Description: BoundedQueue 单元测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestBoundedQueueBasic 基本功能测试
func TestBoundedQueueBasic(t *testing.T) {
	q := NewBoundedQueue[int](10, 1000)
	defer q.Close()

	ctx := context.Background()

	// 测试入队
	for i := 0; i < 5; i++ {
		err := q.Enqueue(ctx, i)
		if err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
	}

	if q.Size() != 5 {
		t.Errorf("Expected size 5, got %d", q.Size())
	}

	// 测试出队
	for i := 0; i < 5; i++ {
		val, err := q.Dequeue(ctx)
		if err != nil {
			t.Fatalf("Dequeue failed: %v", err)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	if !q.IsEmpty() {
		t.Error("Queue should be empty")
	}
}

// TestBoundedQueueAutoResize 测试自动扩缩容
func TestBoundedQueueAutoResize(t *testing.T) {
	q := NewBoundedQueue[string](10, 10000)
	defer q.Close()

	ctx := context.Background()

	// 填充超过初始容量,触发扩容
	for i := 0; i < 100; i++ {
		err := q.Enqueue(ctx, "item")
		if err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
	}

	stats := q.Stats()
	resizeCount := stats["resizeCount"].(int64)
	if resizeCount == 0 {
		t.Error("Expected auto-resize to occur")
	}

	initialCap := q.Cap()
	t.Logf("After expansion: size=%d, cap=%d, resizes=%d", q.Size(), initialCap, resizeCount)

	// 消费大部分元素,触发缩容
	for i := 0; i < 95; i++ {
		_, err := q.Dequeue(ctx)
		if err != nil {
			t.Fatalf("Dequeue failed: %v", err)
		}
	}

	stats = q.Stats()
	shrinkCount := stats["shrinkCount"].(int64)
	finalCap := q.Cap()

	t.Logf("After shrink: size=%d, cap=%d, shrinks=%d", q.Size(), finalCap, shrinkCount)

	if shrinkCount == 0 {
		t.Error("Expected auto-shrink to occur")
	}

	if finalCap >= initialCap {
		t.Errorf("Expected capacity to shrink, got %d >= %d", finalCap, initialCap)
	}
}

// TestBoundedQueueMaxCapacity 测试最大容量限制
func TestBoundedQueueMaxCapacity(t *testing.T) {
	maxCap := 50
	q := NewBoundedQueue[int](10, maxCap)
	defer q.Close()

	ctx := context.Background()

	// 禁用自动扩容后填充到最大容量
	q.SetAutoResize(false)

	// 填充到容量满
	for i := 0; i < 10; i++ {
		err := q.Enqueue(ctx, i)
		if err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
	}

	// 超过容量应该失败
	err := q.Enqueue(ctx, 999)
	if err != ErrQueueFull {
		t.Errorf("Expected ErrQueueFull, got %v", err)
	}

	// 启用自动扩容
	q.SetAutoResize(true)

	// 填充到最大容量
	for q.Size() < maxCap {
		err := q.Enqueue(ctx, 1)
		if err != nil {
			t.Fatalf("Enqueue with resize failed: %v", err)
		}
	}

	// 达到最大容量后应该失败
	err = q.Enqueue(ctx, 999)
	if err != ErrQueueFull {
		t.Errorf("Expected ErrQueueFull at max capacity, got %v", err)
	}
}

// TestBoundedQueueBlocking 测试阻塞操作
func TestBoundedQueueBlocking(t *testing.T) {
	q := NewBoundedQueue[int](10, 1000)
	defer q.Close()

	ctx := context.Background()

	// 启动消费者(延迟消费)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		val, err := q.Dequeue(ctx)
		if err != nil {
			t.Errorf("Dequeue failed: %v", err)
			return
		}
		if val != 42 {
			t.Errorf("Expected 42, got %d", val)
		}
	}()

	// 先入队一个元素
	err := q.Enqueue(ctx, 42)
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	wg.Wait()
}

// TestBoundedQueueTryDequeue 测试非阻塞取出
func TestBoundedQueueTryDequeue(t *testing.T) {
	q := NewBoundedQueue[string](10, 1000)
	defer q.Close()

	ctx := context.Background()

	// 空队列时 TryDequeue 应该立即返回
	_, ok := q.TryDequeue()
	if ok {
		t.Error("Expected TryDequeue to fail on empty queue")
	}

	// 入队后成功取出
	err := q.Enqueue(ctx, "test")
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	val, ok := q.TryDequeue()
	if !ok {
		t.Error("Expected TryDequeue to succeed")
	}
	if val != "test" {
		t.Errorf("Expected 'test', got '%s'", val)
	}
}

// TestBoundedQueueClose 测试关闭队列
func TestBoundedQueueClose(t *testing.T) {
	q := NewBoundedQueue[int](10, 1000)

	ctx := context.Background()

	// 入队一些元素
	for i := 0; i < 5; i++ {
		err := q.Enqueue(ctx, i)
		if err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
	}

	// 关闭队列
	q.Close()

	if !q.IsClosed() {
		t.Error("Queue should be closed")
	}

	// 关闭后入队应该失败
	err := q.Enqueue(ctx, 999)
	if err != ErrQueueClosed {
		t.Errorf("Expected ErrQueueClosed, got %v", err)
	}

	// 可以继续取出已有元素
	for i := 0; i < 5; i++ {
		val, err := q.Dequeue(ctx)
		if err != nil {
			t.Fatalf("Dequeue failed: %v", err)
		}
		if val != i {
			t.Errorf("Expected %d, got %d", i, val)
		}
	}

	// 队列空且关闭后, Dequeue应该返回错误
	_, err = q.Dequeue(ctx)
	if err != ErrQueueClosed {
		t.Errorf("Expected ErrQueueClosed on empty closed queue, got %v", err)
	}
}

// TestBoundedQueueContextCancel 测试上下文取消
func TestBoundedQueueContextCancel(t *testing.T) {
	q := NewBoundedQueue[int](10, 1000)
	defer q.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// 空队列上阻塞 Dequeue,应该在超时后返回
	_, err := q.Dequeue(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

// TestBoundedQueueConcurrent 并发测试
func TestBoundedQueueConcurrent(t *testing.T) {
	q := NewBoundedQueue[int](100, 100000)
	defer q.Close()

	ctx := context.Background()

	const producers = 10
	const consumers = 10
	const itemsPerProducer = 1000

	var wg sync.WaitGroup
	var enqueued, dequeued int64

	// 启动生产者
	for i := 0; i < producers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < itemsPerProducer; j++ {
				err := q.Enqueue(ctx, id*itemsPerProducer+j)
				if err != nil {
					t.Errorf("Producer %d enqueue failed: %v", id, err)
					return
				}
				atomic.AddInt64(&enqueued, 1)
			}
		}(i)
	}

	// 启动消费者
	for i := 0; i < consumers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				_, err := q.Dequeue(ctx)
				cancel()

				if err == context.DeadlineExceeded {
					// 超时,可能所有元素已消费完
					if atomic.LoadInt64(&enqueued) >= producers*itemsPerProducer {
						return
					}
					continue
				}

				if err != nil {
					t.Errorf("Consumer %d dequeue failed: %v", id, err)
					return
				}

				atomic.AddInt64(&dequeued, 1)
				if atomic.LoadInt64(&dequeued) >= producers*itemsPerProducer {
					return
				}
			}
		}(i)
	}

	wg.Wait()

	finalEnqueued := atomic.LoadInt64(&enqueued)
	finalDequeued := atomic.LoadInt64(&dequeued)

	t.Logf("Enqueued: %d, Dequeued: %d, Remaining: %d",
		finalEnqueued, finalDequeued, q.Size())

	expectedTotal := int64(producers * itemsPerProducer)
	if finalEnqueued != expectedTotal {
		t.Errorf("Expected %d enqueued, got %d", expectedTotal, finalEnqueued)
	}

	if finalDequeued+int64(q.Size()) != expectedTotal {
		t.Errorf("Lost items: enqueued=%d, dequeued=%d, remaining=%d",
			finalEnqueued, finalDequeued, q.Size())
	}
}

// TestBoundedQueueSmartGrowth 测试智能增长策略
func TestBoundedQueueSmartGrowth(t *testing.T) {
	q := NewBoundedQueue[int](10, 100000)
	defer q.Close()

	ctx := context.Background()

	// 逐步增加元素,观察容量变化
	steps := []int{10, 100, 500, 1000, 5000, 10000, 20000}

	t.Log("容量增长过程:")
	var prevResizeCount int64 = 0

	for _, targetCount := range steps {
		// 填充到目标数量
		for q.Size() < targetCount {
			err := q.Enqueue(ctx, 1)
			if err != nil {
				t.Fatalf("Enqueue failed: %v", err)
			}
		}

		stats := q.Stats()
		resizeCount := stats["resizeCount"].(int64)
		newResizes := resizeCount - prevResizeCount

		t.Logf("  元素数: %5d, 容量: %6d, 本阶段扩容: %d次, 累计扩容: %d次, 使用率: %.1f%%",
			targetCount, q.Cap(), newResizes, resizeCount, stats["utilization"])

		prevResizeCount = resizeCount
	}

	// 验证智能增长比2倍增长更节省内存
	finalCap := q.Cap()
	doubleGrowthCap := 10
	for doubleGrowthCap < 20000 {
		doubleGrowthCap *= 2
	}

	savings := float64(doubleGrowthCap-finalCap) / float64(doubleGrowthCap) * 100
	t.Logf("\n内存节省: 2倍增长需要 %d, 智能增长需要 %d, 节省 %.1f%%",
		doubleGrowthCap, finalCap, savings)

	if finalCap >= doubleGrowthCap*2 {
		t.Errorf("Smart growth should not exceed 2x of double growth: %d >= %d",
			finalCap, doubleGrowthCap*2)
	}
}

// BenchmarkBoundedQueueEnqueue 入队性能测试
func BenchmarkBoundedQueueEnqueue(b *testing.B) {
	q := NewBoundedQueue[int](1000, 1000000)
	defer q.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Enqueue(ctx, i)
	}
}

// BenchmarkBoundedQueueDequeue 出队性能测试
func BenchmarkBoundedQueueDequeue(b *testing.B) {
	q := NewBoundedQueue[int](b.N, b.N*2)
	defer q.Close()

	ctx := context.Background()

	// 预填充
	for i := 0; i < b.N; i++ {
		_ = q.Enqueue(ctx, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Dequeue(ctx)
	}
}

// BenchmarkBoundedQueueConcurrent 并发性能测试
func BenchmarkBoundedQueueConcurrent(b *testing.B) {
	q := NewBoundedQueue[int](10000, 1000000)
	defer q.Close()

	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				_ = q.Enqueue(ctx, i)
			} else {
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
				_, _ = q.Dequeue(ctx)
				cancel()
			}
			i++
		}
	})
}
