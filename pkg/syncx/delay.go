/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-23 09:08:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 19:03:15
 * @FilePath: \go-toolbox\pkg\syncx\delay.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// DelayStrategy 定义延迟策略
type DelayStrategy int

const (
	// FixedDelayStrategy 固定延迟策略
	FixedDelayStrategy DelayStrategy = iota
	// LinearDelayStrategy 线性递增延迟策略
	LinearDelayStrategy
	// ExponentialDelayStrategy 指数延迟策略
	ExponentialDelayStrategy
	// RandomDelayStrategy 随机延迟策略
	RandomDelayStrategy
	// CustomDelayStrategy 自定义延迟策略
	CustomDelayStrategy
)

// DelayFunc 自定义延迟函数类型
type DelayFunc func(attempt int, baseDelay time.Duration) time.Duration

// ExecutionResult 执行结果
type ExecutionResult struct {
	Index     int           // 执行索引
	StartTime time.Time     // 开始时间
	EndTime   time.Time     // 结束时间
	Duration  time.Duration // 执行耗时
	Error     error         // 执行错误
	Skipped   bool          // 是否被跳过
}

// ExecutionStats 执行统计
type ExecutionStats struct {
	Total     int64         // 总执行次数
	Success   int64         // 成功次数
	Failed    int64         // 失败次数
	Skipped   int64         // 跳过次数
	TotalTime time.Duration // 总耗时
	AvgTime   time.Duration // 平均耗时
}

// Delayer 结构体，用于支持链式调用的延迟执行器
type Delayer struct {
	// 基本配置
	delay           time.Duration // 基础延迟时间
	execCount       int           // 执行次数
	function        func() error  // 要执行的函数（支持返回错误）
	functionNoError func()        // 无错误返回的函数

	// 延迟策略相关
	strategy        DelayStrategy // 延迟策略
	customDelayFunc DelayFunc     // 自定义延迟函数
	randomBase      float64       // 随机基数
	maxDelay        time.Duration // 最大延迟时间
	multiplier      float64       // 指数策略的倍数

	// 执行控制
	ctx            context.Context    // 上下文
	cancelFunc     context.CancelFunc // 取消函数
	concurrent     bool               // 是否并发执行
	maxConcurrency int                // 最大并发数
	stopOnError    bool               // 遇到错误是否停止

	// 回调和监控
	onStart    func(index int)                 // 开始执行回调
	onComplete func(result *ExecutionResult)   // 完成执行回调
	onError    func(index int, err error) bool // 错误处理回调（返回true表示继续执行）

	// 内部状态
	timers         []*time.Timer      // 计时器列表
	results        []*ExecutionResult // 执行结果
	stats          *ExecutionStats    // 统计信息
	mu             sync.RWMutex       // 读写锁
	running        int64              // 正在运行的任务数
	stopped        int64              // 是否已停止
	completionChan chan struct{}      // 等待任务完成的通道
	pendingTasks   int64              // 待执行任务数
}

// NewDelayer 创建一个新的 Delayer 实例
func NewDelayer() *Delayer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Delayer{
		delay:          0,
		execCount:      1,
		strategy:       FixedDelayStrategy,
		randomBase:     1.0,
		maxDelay:       time.Hour,
		multiplier:     2.0,
		ctx:            ctx,
		cancelFunc:     cancel,
		concurrent:     false,
		maxConcurrency: 10,
		stopOnError:    false,
		stats:          &ExecutionStats{},
		results:        make([]*ExecutionResult, 0),
		timers:         make([]*time.Timer, 0),
		completionChan: make(chan struct{}),
	}
}

// WithDelay 设置基础延迟时间
func (d *Delayer) WithDelay(delay time.Duration) *Delayer {
	d.delay = delay
	return d
}

// WithFunction 设置要执行的函数（支持错误返回）
func (d *Delayer) WithFunction(f func() error) *Delayer {
	d.function = f
	return d
}

// WithSimpleFunction 设置要执行的函数（无错误返回）
func (d *Delayer) WithSimpleFunction(f func()) *Delayer {
	d.functionNoError = f
	return d
}

// WithTimes 设置执行次数
func (d *Delayer) WithTimes(count int) *Delayer {
	if count > 0 {
		d.execCount = count
	}
	return d
}

// WithStrategy 设置延迟策略
func (d *Delayer) WithStrategy(strategy DelayStrategy) *Delayer {
	d.strategy = strategy
	return d
}

// WithCustomDelay 设置自定义延迟函数
func (d *Delayer) WithCustomDelay(delayFunc DelayFunc) *Delayer {
	d.strategy = CustomDelayStrategy
	d.customDelayFunc = delayFunc
	return d
}

// WithRandomBase 设置随机基数
func (d *Delayer) WithRandomBase(base float64) *Delayer {
	if base > 0 {
		d.randomBase = base
	}
	return d
}

// WithMaxDelay 设置最大延迟时间
func (d *Delayer) WithMaxDelay(maxDelay time.Duration) *Delayer {
	d.maxDelay = maxDelay
	return d
}

// WithMultiplier 设置指数策略的倍数
func (d *Delayer) WithMultiplier(multiplier float64) *Delayer {
	if multiplier > 0 {
		d.multiplier = multiplier
	}
	return d
}

// WithContext 设置上下文
func (d *Delayer) WithContext(ctx context.Context) *Delayer {
	d.cancelFunc() // 取消之前的上下文
	d.ctx, d.cancelFunc = context.WithCancel(ctx)
	return d
}

// WithConcurrent 设置是否并发执行
func (d *Delayer) WithConcurrent(concurrent bool) *Delayer {
	d.concurrent = concurrent
	return d
}

// WithMaxConcurrency 设置最大并发数
func (d *Delayer) WithMaxConcurrency(maxConcurrency int) *Delayer {
	if maxConcurrency > 0 {
		d.maxConcurrency = maxConcurrency
	}
	return d
}

// WithStopOnError 设置遇到错误是否停止
func (d *Delayer) WithStopOnError(stopOnError bool) *Delayer {
	d.stopOnError = stopOnError
	return d
}

// WithOnStart 设置开始执行回调
func (d *Delayer) WithOnStart(callback func(index int)) *Delayer {
	d.onStart = callback
	return d
}

// WithOnComplete 设置完成执行回调
func (d *Delayer) WithOnComplete(callback func(result *ExecutionResult)) *Delayer {
	d.onComplete = callback
	return d
}

// WithOnError 设置错误处理回调
func (d *Delayer) WithOnError(callback func(index int, err error) bool) *Delayer {
	d.onError = callback
	return d
}

// calculateDelay 计算延迟时间
func (d *Delayer) calculateDelay(attempt int) time.Duration {
	var delay time.Duration

	switch d.strategy {
	case FixedDelayStrategy:
		delay = d.delay
	case LinearDelayStrategy:
		delay = d.delay * time.Duration(attempt+1)
	case ExponentialDelayStrategy:
		delay = time.Duration(float64(d.delay) * math.Pow(d.multiplier, float64(attempt)))
	case RandomDelayStrategy:
		randomFactor := rand.Float64() * d.randomBase
		delay = time.Duration(float64(d.delay) * randomFactor)
	case CustomDelayStrategy:
		if d.customDelayFunc != nil {
			delay = d.customDelayFunc(attempt, d.delay)
		} else {
			delay = d.delay
		}
	default:
		delay = d.delay
	}

	// 应用最大延迟限制
	if delay > d.maxDelay {
		delay = d.maxDelay
	}

	return delay
}

// executeTask 执行单个任务
func (d *Delayer) executeTask(index int) {
	// 增加运行中的任务计数
	atomic.AddInt64(&d.running, 1)
	defer func() {
		// 减少运行中的任务计数
		running := atomic.AddInt64(&d.running, -1)
		pending := atomic.AddInt64(&d.pendingTasks, -1)

		// 如果所有任务都完成了，关闭完成通道
		if running == 0 && pending == 0 {
			select {
			case <-d.completionChan:
				// 通道已经关闭
			default:
				close(d.completionChan)
			}
		}
	}()

	if atomic.LoadInt64(&d.stopped) == 1 {
		return
	}

	result := &ExecutionResult{
		Index:     index,
		StartTime: time.Now(),
	}

	// 检查上下文是否已取消
	select {
	case <-d.ctx.Done():
		result.Error = d.ctx.Err()
		result.Skipped = true
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		d.recordResult(result)
		return
	default:
	}

	// 调用开始回调
	if d.onStart != nil {
		d.onStart(index)
	}

	// 执行函数
	var err error
	if d.function != nil {
		err = d.function()
	} else if d.functionNoError != nil {
		d.functionNoError()
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Error = err

	// 记录结果
	d.recordResult(result)

	// 处理错误
	if err != nil {
		shouldContinue := true
		if d.onError != nil {
			shouldContinue = d.onError(index, err)
		}

		if d.stopOnError || !shouldContinue {
			atomic.StoreInt64(&d.stopped, 1)
			d.Stop()
		}
	}

	// 调用完成回调
	if d.onComplete != nil {
		d.onComplete(result)
	}
}

// recordResult 记录执行结果
func (d *Delayer) recordResult(result *ExecutionResult) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.results = append(d.results, result)

	// 更新统计信息
	atomic.AddInt64(&d.stats.Total, 1)
	d.stats.TotalTime += result.Duration

	if result.Skipped {
		atomic.AddInt64(&d.stats.Skipped, 1)
	} else if result.Error != nil {
		atomic.AddInt64(&d.stats.Failed, 1)
	} else {
		atomic.AddInt64(&d.stats.Success, 1)
	}

	// 计算平均时间
	if d.stats.Total > 0 {
		d.stats.AvgTime = d.stats.TotalTime / time.Duration(d.stats.Total)
	}
}

// Build 开始延迟执行
func (d *Delayer) Build() *Delayer {
	if d.function == nil && d.functionNoError == nil {
		return d
	}

	atomic.StoreInt64(&d.stopped, 0)
	atomic.StoreInt64(&d.pendingTasks, int64(d.execCount))
	atomic.StoreInt64(&d.running, 0)

	// 重新创建完成通道
	d.mu.Lock()
	d.completionChan = make(chan struct{})
	d.mu.Unlock()

	if d.concurrent {
		d.buildConcurrent()
	} else {
		d.buildSequential()
	}

	return d
}

// buildSequential 顺序执行
func (d *Delayer) buildSequential() {
	for i := 0; i < d.execCount; i++ {
		if atomic.LoadInt64(&d.stopped) == 1 {
			break
		}

		delay := d.calculateDelay(i)

		timer := time.AfterFunc(delay, func(index int) func() {
			return func() {
				// 在执行前再次检查是否应该停止
				if atomic.LoadInt64(&d.stopped) == 1 {
					// 减少待执行任务计数
					pending := atomic.AddInt64(&d.pendingTasks, -1)
					running := atomic.LoadInt64(&d.running)
					if running == 0 && pending == 0 {
						select {
						case <-d.completionChan:
							// 通道已经关闭
						default:
							close(d.completionChan)
						}
					}
					return
				}
				// 检查上下文是否已取消
				select {
				case <-d.ctx.Done():
					// 减少待执行任务计数
					pending := atomic.AddInt64(&d.pendingTasks, -1)
					running := atomic.LoadInt64(&d.running)
					if running == 0 && pending == 0 {
						select {
						case <-d.completionChan:
							// 通道已经关闭
						default:
							close(d.completionChan)
						}
					}
					return
				default:
					d.executeTask(index)
				}
			}
		}(i))

		d.mu.Lock()
		d.timers = append(d.timers, timer)
		d.mu.Unlock()
	}
}

// buildConcurrent 并发执行
func (d *Delayer) buildConcurrent() {
	semaphore := make(chan struct{}, d.maxConcurrency)

	for i := 0; i < d.execCount; i++ {
		if atomic.LoadInt64(&d.stopped) == 1 {
			break
		}

		delay := d.calculateDelay(i)

		timer := time.AfterFunc(delay, func(index int) func() {
			return func() {
				// 在执行前再次检查是否应该停止
				if atomic.LoadInt64(&d.stopped) == 1 {
					// 减少待执行任务计数
					pending := atomic.AddInt64(&d.pendingTasks, -1)
					running := atomic.LoadInt64(&d.running)
					if running == 0 && pending == 0 {
						select {
						case <-d.completionChan:
							// 通道已经关闭
						default:
							close(d.completionChan)
						}
					}
					return
				}
				// 检查上下文是否已取消
				select {
				case <-d.ctx.Done():
					// 减少待执行任务计数
					pending := atomic.AddInt64(&d.pendingTasks, -1)
					running := atomic.LoadInt64(&d.running)
					if running == 0 && pending == 0 {
						select {
						case <-d.completionChan:
							// 通道已经关闭
						default:
							close(d.completionChan)
						}
					}
					return
				default:
					semaphore <- struct{}{} // 获取信号量
					go func() {
						defer func() { <-semaphore }() // 释放信号量
						d.executeTask(index)
					}()
				}
			}
		}(i))

		d.mu.Lock()
		d.timers = append(d.timers, timer)
		d.mu.Unlock()
	}
}

// Stop 停止所有待执行的任务
func (d *Delayer) Stop() {
	atomic.StoreInt64(&d.stopped, 1)
	d.cancelFunc()

	d.mu.Lock()
	defer d.mu.Unlock()

	// 停止所有计时器
	for _, timer := range d.timers {
		timer.Stop()
	}
}

// Wait 等待所有任务完成
// Wait 等待上下文取消
func (d *Delayer) Wait() {
	<-d.ctx.Done()
}

// WaitForCompletion 等待所有任务完成
func (d *Delayer) WaitForCompletion() {
	<-d.completionChan
}

// WaitForCompletionWithTimeout 等待所有任务完成（带超时）
func (d *Delayer) WaitForCompletionWithTimeout(timeout time.Duration) bool {
	select {
	case <-d.completionChan:
		return true
	case <-time.After(timeout):
		return false
	}
}

// GetResults 获取所有执行结果
func (d *Delayer) GetResults() []*ExecutionResult {
	d.mu.RLock()
	defer d.mu.RUnlock()

	results := make([]*ExecutionResult, len(d.results))
	copy(results, d.results)
	return results
}

// GetStats 获取执行统计信息
func (d *Delayer) GetStats() *ExecutionStats {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return &ExecutionStats{
		Total:     atomic.LoadInt64(&d.stats.Total),
		Success:   atomic.LoadInt64(&d.stats.Success),
		Failed:    atomic.LoadInt64(&d.stats.Failed),
		Skipped:   atomic.LoadInt64(&d.stats.Skipped),
		TotalTime: d.stats.TotalTime,
		AvgTime:   d.stats.AvgTime,
	}
}

// IsRunning 检查是否有任务正在运行
func (d *Delayer) IsRunning() bool {
	return atomic.LoadInt64(&d.running) > 0
}

// IsStopped 检查是否已停止
func (d *Delayer) IsStopped() bool {
	return atomic.LoadInt64(&d.stopped) == 1
}

// Reset 重置所有状态
func (d *Delayer) Reset() *Delayer {
	d.Stop()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.results = make([]*ExecutionResult, 0)
	d.timers = make([]*time.Timer, 0)
	d.stats = &ExecutionStats{}
	d.completionChan = make(chan struct{})
	atomic.StoreInt64(&d.stopped, 0)
	atomic.StoreInt64(&d.running, 0)
	atomic.StoreInt64(&d.pendingTasks, 0)

	// 重新创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	d.ctx = ctx
	d.cancelFunc = cancel

	return d
}

// String 返回当前配置的字符串表示
func (d *Delayer) String() string {
	return fmt.Sprintf("Delayer{delay=%v, count=%d, strategy=%d, concurrent=%v, maxConcurrency=%d}",
		d.delay, d.execCount, d.strategy, d.concurrent, d.maxConcurrency)
}
