/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-10 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-10 00:00:00
 * @FilePath: \go-toolbox\pkg\syncx\worker_pool.go
 * @Description: Worker 池实现，用于限制并发 goroutine 数量，防止 OOM
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerTask 任务接口
type WorkerTask func()

// WorkerPool Worker 池，用于限制并发 goroutine 数量
// 防止高频操作导致 goroutine 无限增长导致 OOM
//
// 使用场景:
//   - PubSub 消息处理
//   - 高频任务分发
//   - 并发控制
//
// 示例:
//
//	pool := NewWorkerPool(20, 100)
//	defer pool.Close()
//
//	for i := 0; i < 1000; i++ {
//	    err := pool.Submit(ctx, func() {
//	        处理任务
//	    })
//	    if err != nil {
//	        log.Errorf("submit failed: %v", err)
//	    }
//	}
type WorkerPool struct {
	workers int             // worker 数量
	queue   chan WorkerTask // 任务队列
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	once    sync.Once
	closed  bool
	mu      sync.Mutex
	active  int32 // 活跃任务计数
}

var (
	// ErrClosed 表示 Worker 池已关闭
	ErrClosed = errors.New("worker pool is closed")

	// ErrQueueFull 表示 Worker 池队列已满
	ErrQueueFull = errors.New("worker pool queue is full")
)

// NewWorkerPool 创建 Worker 池
// workers: worker 数量，建议 10-50，根据 CPU 核心数调整
// queueSize: 任务队列大小，建议 100-1000
//
// 示例:
//
//	创建 20 个 worker，队列大小 100
//	pool := NewWorkerPool(20, 100)
//	defer pool.Close()
func NewWorkerPool(workers, queueSize int) *WorkerPool {
	if workers <= 0 {
		workers = 10 // 默认 10 个 worker
	}
	if queueSize <= 0 {
		queueSize = 100 // 默认队列大小 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workers: workers,
		queue:   make(chan WorkerTask, queueSize),
		ctx:     ctx,
		cancel:  cancel,
	}

	// 启动 worker goroutine
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	return pool
}

// worker 工作 goroutine，从队列中取任务执行
func (p *WorkerPool) worker() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			// context 已取消，优先检查，确保快速退出
			return
		case task, ok := <-p.queue:
			if !ok {
				// 队列已关闭，退出 worker
				return
			}
			if task != nil {
				atomic.AddInt32(&p.active, 1)
				task()
				atomic.AddInt32(&p.active, -1)
			}
		}
	}
}

// Submit 提交任务到队列
// 如果队列满，会阻塞直到有空位或 context 取消
//
// 参数:
//   - ctx: 上下文，用于超时和取消控制
//   - task: 要执行的任务
//
// 返回:
//   - nil: 任务成功提交
//   - ErrClosed: Worker 池已关闭
//   - context.Err(): context 被取消或超时
//
// 示例:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	err := pool.Submit(ctx, func() {
//	   处理任务
//	})
//	if err != nil {
//	    log.Errorf("submit failed: %v", err)
//	}
func (p *WorkerPool) Submit(ctx context.Context, task WorkerTask) error {
	if task == nil {
		return nil
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return ErrClosed
	}
	p.mu.Unlock()

	select {
	case p.queue <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return ErrClosed
	}
}

// SubmitNonBlocking 非阻塞提交任务
// 如果队列满，返回错误而不是阻塞
//
// 参数:
//   - task: 要执行的任务
//
// 返回:
//   - nil: 任务成功提交
//   - ErrClosed: Worker 池已关闭
//   - ErrQueueFull: 队列已满
//
// 示例:
//
//	err := pool.SubmitNonBlocking(func() {
//	   处理任务
//	})
//	if err == ErrQueueFull {
//	    log.Warn("queue is full, task dropped")
//	}
func (p *WorkerPool) SubmitNonBlocking(task WorkerTask) error {
	if task == nil {
		return nil
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return ErrClosed
	}
	p.mu.Unlock()

	select {
	case p.queue <- task:
		return nil
	default:
		return ErrQueueFull
	}
}

// Wait 等待所有已提交的任务完成
// 不关闭 pool，可以继续提交新任务
//
// 示例:
//
//	pool := NewWorkerPool(20, 100)
//	提交任务 ...
//	pool.Wait()  // 等待所有任务完成
//	继续提交任务 ...
func (p *WorkerPool) Wait() {
	for {
		queueLen := len(p.queue)
		activeCount := atomic.LoadInt32(&p.active)

		if queueLen == 0 && activeCount == 0 {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
}

// Close 关闭 Worker 池，等待所有任务完成
// 调用此方法后，所有新的提交都会返回 ErrClosed
//
// 示例:
//
//	pool := NewWorkerPool(20, 100)
//	提交任务 ...
//	pool.Close()  // 等待所有任务完成后返回
func (p *WorkerPool) Close() error {
	p.once.Do(func() {
		p.mu.Lock()
		p.closed = true
		p.mu.Unlock()

		// 关闭队列，停止接收新任务
		close(p.queue)

		// 取消 context，通知所有 worker 退出
		p.cancel()

		// 等待所有 worker 完成
		p.wg.Wait()
	})

	return nil
}

// GetQueueSize 获取队列中待处理任务数
func (p *WorkerPool) GetQueueSize() int {
	return len(p.queue)
}

// GetWorkerCount 获取 worker 数量
func (p *WorkerPool) GetWorkerCount() int {
	return p.workers
}

// IsClosed 检查 Worker 池是否已关闭
func (p *WorkerPool) IsClosed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.closed
}
