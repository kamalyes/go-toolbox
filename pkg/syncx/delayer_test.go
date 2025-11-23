/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 19:03:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 19:03:15
 * @FilePath: \go-toolbox\syncx\delayer_test.go
 * @Description: Delayer 测试文件
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestDelayerBasicFunctionality(t *testing.T) {
	execCount := 0
	delayer := NewDelayer().
		WithDelay(100 * time.Millisecond).
		WithTimes(3).
		WithSimpleFunction(func() {
			execCount++
		})

	delayer.Build()
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 3, execCount, "Expected 3 executions")
}

func TestDelayerWithError(t *testing.T) {
	delayer := NewDelayer().
		WithDelay(50 * time.Millisecond).
		WithTimes(3).
		WithFunction(func() error {
			return errors.New("test error")
		})

	delayer.Build()
	time.Sleep(300 * time.Millisecond) // 增加等待时间

	stats := delayer.GetStats()
	assert.Greater(t, stats.Failed, int64(0), "Expected failed executions")
	assert.Equal(t, int64(3), stats.Total, "Expected 3 total executions")
}

func TestDelayerLinearStrategy(t *testing.T) {
	execCount := 0

	delayer := NewDelayer().
		WithDelay(50 * time.Millisecond).
		WithStrategy(LinearDelayStrategy).
		WithTimes(3).
		WithSimpleFunction(func() {
			execCount++
		})

	delayer.Build()
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 3, execCount, "Expected 3 executions")
}

func TestDelayerConcurrentExecution(t *testing.T) {
	execCount := 0

	delayer := NewDelayer().
		WithDelay(50 * time.Millisecond).
		WithTimes(5).
		WithConcurrent(true).
		WithMaxConcurrency(3).
		WithSimpleFunction(func() {
			execCount++
		})

	delayer.Build()
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 5, execCount, "Expected 5 executions")
}

func TestDelayerContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	execCount := 0
	var mu sync.Mutex

	delayer := NewDelayer().
		WithContext(ctx).
		WithDelay(200 * time.Millisecond). // 增加延迟时间
		WithTimes(10).
		WithSimpleFunction(func() {
			mu.Lock()
			execCount++
			mu.Unlock()
		})

	delayer.Build()

	// 在短时间后取消上下文（在第一个任务执行之前）
	time.AfterFunc(100*time.Millisecond, cancel)

	// 等待上下文取消
	delayer.Wait()

	// 额外等待一点时间确保所有计时器都有机会被检查
	time.Sleep(100 * time.Millisecond)

	// 读取最终的执行计数
	mu.Lock()
	finalCount := execCount
	mu.Unlock()

	// 由于我们在第一个任务执行前就取消了上下文，应该没有任务执行
	assert.Equal(t, 0, finalCount, "Expected 0 executions due to early cancellation, got %d", finalCount)
}

func TestDelayerWaitForCompletion(t *testing.T) {
	execCount := 0
	var mu sync.Mutex

	delayer := NewDelayer().
		WithDelay(50 * time.Millisecond).
		WithTimes(3).
		WithSimpleFunction(func() {
			mu.Lock()
			execCount++
			mu.Unlock()
			time.Sleep(100 * time.Millisecond) // 模拟耗时操作
		})

	start := time.Now()
	delayer.Build()

	// 等待所有任务完成
	delayer.WaitForCompletion()
	duration := time.Since(start)

	mu.Lock()
	finalCount := execCount
	mu.Unlock()

	assert.Equal(t, 3, finalCount, "Expected 3 executions")
	// 验证确实等待了任务完成（应该超过最后一个任务的延迟 + 执行时间）
	// 最后一个任务: 50ms延迟 + 100ms执行 = 150ms
	assert.GreaterOrEqual(t, duration, 140*time.Millisecond, "Should wait for tasks to complete, took %v", duration)
}

func TestDelayerWaitForCompletionWithTimeout(t *testing.T) {
	delayer := NewDelayer().
		WithDelay(50 * time.Millisecond).
		WithTimes(3).
		WithSimpleFunction(func() {
			time.Sleep(100 * time.Millisecond) // 模拟耗时操作
		})

	delayer.Build()

	// 使用较短的超时时间
	completed := delayer.WaitForCompletionWithTimeout(100 * time.Millisecond)
	assert.False(t, completed, "Should timeout before completion")

	// 使用足够长的超时时间
	completed = delayer.WaitForCompletionWithTimeout(1 * time.Second)
	assert.True(t, completed, "Should complete within timeout")
}
