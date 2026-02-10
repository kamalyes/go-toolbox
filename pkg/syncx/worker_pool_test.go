/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-10 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-10 00:00:00
 * @FilePath: \go-toolbox\pkg\syncx\worker_pool_test.go
 * @Description: Worker 池测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testTaskCount = 100

// TestWorkerPoolBasic 测试 Worker 池基本功能
func TestWorkerPoolBasic(t *testing.T) {
	pool := NewWorkerPool(5, testTaskCount)
	defer pool.Close()

	counter := 0
	var mu sync.Mutex

	for i := 0; i < testTaskCount; i++ {
		err := pool.Submit(context.Background(), func() {
			mu.Lock()
			counter++
			mu.Unlock()
		})
		assert.NoError(t, err)
	}

	// 等待所有任务完成
	pool.Wait()

	mu.Lock()
	assert.Equal(t, testTaskCount, counter)
	mu.Unlock()
}

// TestWorkerPoolQueueFull 测试 Worker 池队列满的情况
func TestWorkerPoolQueueFull(t *testing.T) {
	pool := NewWorkerPool(1, 2)
	defer pool.Close()

	// 提交会阻塞的任务
	blockChan := make(chan struct{})
	err := pool.Submit(context.Background(), func() {
		<-blockChan
	})
	assert.NoError(t, err)

	// 提交第二个任务（填满队列）
	err = pool.Submit(context.Background(), func() {})
	assert.NoError(t, err)

	// 非阻塞提交应该返回队列满错误
	err = pool.SubmitNonBlocking(func() {})
	assert.ErrorIs(t, err, ErrQueueFull)

	// 解除阻塞
	close(blockChan)
	time.Sleep(100 * time.Millisecond)
}

// TestWorkerPoolContextCancellation 测试 Worker 池的 context 取消
func TestWorkerPoolContextCancellation(t *testing.T) {
	pool := NewWorkerPool(2, 10)

	counter := 0
	var mu sync.Mutex

	// 提交一些任务
	for i := 0; i < 5; i++ {
		err := pool.Submit(context.Background(), func() {
			mu.Lock()
			counter++
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
		})
		assert.NoError(t, err)
	}

	// 等待任务完成
	pool.Wait()

	// 验证任务已完成
	mu.Lock()
	assert.Greater(t, counter, 0)
	mu.Unlock()

	// 关闭后提交应该返回错误
	pool.Close()
	err := pool.Submit(context.Background(), func() {})
	assert.ErrorIs(t, err, ErrClosed)
}

// TestWorkerPoolClosed 测试 Worker 池关闭后的行为
func TestWorkerPoolClosed(t *testing.T) {
	pool := NewWorkerPool(2, 10)
	pool.Close()

	// 关闭后提交应该返回错误
	err := pool.Submit(context.Background(), func() {})
	assert.ErrorIs(t, err, ErrClosed)

	// 非阻塞提交也应该返回错误
	err = pool.SubmitNonBlocking(func() {})
	assert.ErrorIs(t, err, ErrClosed)
}

// TestWorkerPoolGoroutineCount 测试 Worker 池的 goroutine 数量
func TestWorkerPoolGoroutineCount(t *testing.T) {
	initialGoroutines := runtime.NumGoroutine()

	pool := NewWorkerPool(10, 100)
	time.Sleep(100 * time.Millisecond)

	// 创建 pool 后应该增加 10 个 worker goroutine
	afterCreateGoroutines := runtime.NumGoroutine()
	assert.Greater(t, afterCreateGoroutines, initialGoroutines)

	pool.Close()
	time.Sleep(100 * time.Millisecond)

	// 关闭后 goroutine 数量应该恢复
	afterCloseGoroutines := runtime.NumGoroutine()
	assert.LessOrEqual(t, afterCloseGoroutines, initialGoroutines+2) // 允许小的偏差
}

// TestWorkerPoolStress 压力测试 Worker 池
func TestWorkerPoolStress(t *testing.T) {
	pool := NewWorkerPool(20, 1000)
	defer pool.Close()

	taskCount := 10000
	completedCount := 0
	var mu sync.Mutex

	for i := 0; i < taskCount; i++ {
		err := pool.Submit(context.Background(), func() {
			mu.Lock()
			completedCount++
			mu.Unlock()
		})
		assert.NoError(t, err)
	}

	pool.Wait()

	mu.Lock()
	assert.Equal(t, taskCount, completedCount)
	mu.Unlock()
}

// TestWorkerPoolQueueSize 测试 Worker 池队列大小
func TestWorkerPoolQueueSize(t *testing.T) {
	pool := NewWorkerPool(2, 10)
	defer pool.Close()

	// 初始队列应该为空
	assert.Equal(t, 0, pool.GetQueueSize())

	// 提交阻塞任务
	blockChan := make(chan struct{})
	pool.Submit(context.Background(), func() {
		<-blockChan
	})

	// 提交更多任务
	for i := 0; i < 5; i++ {
		pool.Submit(context.Background(), func() {})
	}

	// 队列大小应该反映待处理任务数
	queueSize := pool.GetQueueSize()
	assert.Greater(t, queueSize, 0)

	close(blockChan)
	time.Sleep(100 * time.Millisecond)
}

// TestWorkerPoolWorkerCount 测试 Worker 池的 worker 数量
func TestWorkerPoolWorkerCount(t *testing.T) {
	workerCount := 15
	pool := NewWorkerPool(workerCount, 100)
	defer pool.Close()

	assert.Equal(t, workerCount, pool.GetWorkerCount())
}

// TestWorkerPoolIsClosed 测试 Worker 池的关闭状态
func TestWorkerPoolIsClosed(t *testing.T) {
	pool := NewWorkerPool(5, 10)

	assert.False(t, pool.IsClosed())

	pool.Close()

	assert.True(t, pool.IsClosed())
}

// TestWorkerPoolNilTask 测试提交 nil 任务
func TestWorkerPoolNilTask(t *testing.T) {
	pool := NewWorkerPool(5, 10)
	defer pool.Close()

	// 提交 nil 任务应该返回 nil
	err := pool.Submit(context.Background(), nil)
	assert.NoError(t, err)

	// 非阻塞提交 nil 任务也应该返回 nil
	err = pool.SubmitNonBlocking(nil)
	assert.NoError(t, err)
}

// TestWorkerPoolContextTimeout 测试 context 超时
func TestWorkerPoolContextTimeout(t *testing.T) {
	pool := NewWorkerPool(1, 1)
	defer pool.Close()

	// 提交会阻塞的任务，填满队列
	blockChan := make(chan struct{})
	err := pool.Submit(context.Background(), func() {
		<-blockChan
	})
	assert.NoError(t, err)

	// 队列现在满了，再提交一个任务到队列
	err = pool.Submit(context.Background(), func() {})
	assert.NoError(t, err)

	// 创建超时 context，队列已满，应该超时
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// 提交任务应该超时
	err = pool.Submit(ctx, func() {})
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)

	close(blockChan)
}
