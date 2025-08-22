/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-21 18:57:15
 * @FilePath: \go-toolbox\tests\fifo_queue_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/queue"
)

func BenchmarkFIFOQueue_Enqueue(b *testing.B) {
	q := queue.NewFIFOQueue(1000, true)
	ctx := context.Background()

	b.ResetTimer() // 重置计时器，确保不包括设置时间
	for i := 0; i < b.N; i++ {
		if err := q.Enqueue(ctx, i); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkFIFOQueue_Dequeue(b *testing.B) {
	q := queue.NewFIFOQueue(1000, true)
	ctx := context.Background()

	// 先填充队列
	for i := 0; i < 1000; i++ {
		if err := q.Enqueue(ctx, i); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}

	b.ResetTimer() // 重置计时器，确保不包括设置时间
	for i := 0; i < b.N; i++ {
		if _, err := q.Dequeue(ctx); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkFIFOQueue_Concurrent(b *testing.B) {
	q := queue.NewFIFOQueue(1000, true)
	ctx := context.Background()
	var wg sync.WaitGroup

	// 启动多个 goroutine 来并发入队
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < b.N/10; j++ { // 每个 goroutine 入队一定数量的元素
				if err := q.Enqueue(ctx, time.Now().UnixNano()); err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		}()
	}

	// 等待所有入队操作完成
	wg.Wait()

	// 进行出队操作
	b.ResetTimer() // 重置计时器，确保不包括设置时间
	for i := 0; i < b.N; i++ {
		if _, err := q.Dequeue(ctx); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
