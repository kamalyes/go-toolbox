/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-23 09:08:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 19:35:00
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

// ExecutionContext 执行上下文，包含执行过程中的全部信息
type ExecutionContext struct {
	Index      int                    // 执行索引
	Attempt    int                    // 尝试次数（从0开始）
	StartTime  time.Time              // 开始时间
	EndTime    time.Time              // 结束时间
	Duration   time.Duration          // 执行耗时
	Error      error                  // 执行错误
	Skipped    bool                   // 是否被跳过
	Delay      time.Duration          // 实际延迟时间
	Strategy   DelayStrategy          // 使用的延迟策略
	Concurrent bool                   // 是否并发执行
	Cancelled  bool                   // 是否被取消
	Retryable  bool                   // 是否可重试
	Metadata   map[string]interface{} // 自定义元数据
}

// ExecutionResult 执行结果（为了向后兼容保留）
type ExecutionResult struct {
	Index     int           // 执行索引
	StartTime time.Time     // 开始时间
	EndTime   time.Time     // 结束时间
	Duration  time.Duration // 执行耗时
	Error     error         // 执行错误
	Skipped   bool          // 是否被跳过
}

// ToExecutionContext 将 ExecutionResult 转换为 ExecutionContext
func (r *ExecutionResult) ToExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		Index:     r.Index,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Duration:  r.Duration,
		Error:     r.Error,
		Skipped:   r.Skipped,
		Metadata:  make(map[string]interface{}),
	}
}

// TaskFunc 泛型任务函数类型
type TaskFunc[T any] func(ctx *ExecutionContext) (T, error)

// SimpleTaskFunc 简单任务函数类型（无返回值）
type SimpleTaskFunc func(ctx *ExecutionContext) error

// CallbackFunc 泛型回调函数类型
type CallbackFunc[T any] func(ctx *ExecutionContext, result T)

// ErrorHandlerFunc 错误处理回调函数类型
type ErrorHandlerFunc func(ctx *ExecutionContext) (shouldContinue bool)

// TaskProgressFunc 任务进度回调函数类型
type TaskProgressFunc func(completed int64, total int64, percentage float64)

// ExecutionStats 执行统计
type ExecutionStats struct {
	StartTime      time.Time     // 开始时间
	EndTime        time.Time     // 结束时间
	TotalDuration  time.Duration // 总耗时
	SuccessCount   int64         // 成功次数
	ErrorCount     int64         // 错误次数
	SkippedCount   int64         // 跳过次数
	CancelledCount int64         // 取消次数
}

// Delayer 统一的泛型延迟执行器
type Delayer[T any] struct {
	// 基本配置
	delay           time.Duration // 基础延迟时间
	execCount       int           // 执行次数
	function        func() error  // 要执行的函数（支持返回错误）
	functionNoError func()        // 无错误返回的函数
	taskFunc        TaskFunc[T]   // 泛型任务函数

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
	onBeforeStart   func(ctx *ExecutionContext) // 开始执行前回调
	onAfterComplete func(ctx *ExecutionContext) // 完成执行后回调
	onError         ErrorHandlerFunc            // 错误处理回调
	onProgress      TaskProgressFunc            // 进度回调
	onSuccess       CallbackFunc[T]             // 泛型成功回调

	// 性能优化相关
	callbackPool     sync.Pool // ExecutionContext 对象池
	disableCallbacks bool      // 禁用回调以提升性能

	// 泛型结果相关
	resultChannel    chan T       // 结果通道
	genericResults   []T          // 泛型任务结果集合
	genericResultsMu sync.RWMutex // 泛型结果集合读写锁
	channelClosed    int64        // 通道是否已关闭 (0=开启, 1=关闭)

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

// NewDelayer 创建一个新的泛型 Delayer 实例
func NewDelayer[T any]() *Delayer[T] {
	ctx, cancel := context.WithCancel(context.Background())
	d := &Delayer[T]{
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
		resultChannel:  make(chan T, 100), // 泛型结果通道
		genericResults: make([]T, 0, 100), // 预分配容量以减少重新分配
	}

	// 初始化对象池
	d.callbackPool = sync.Pool{
		New: func() interface{} {
			return &ExecutionContext{
				Metadata: make(map[string]interface{}),
			}
		},
	}

	return d
}

// WithDelay 设置基础延迟时间
func (d *Delayer[T]) WithDelay(delay time.Duration) *Delayer[T] {
	d.delay = delay
	return d
}

// WithFunction 设置要执行的函数
func (d *Delayer[T]) WithFunction(f func() error) *Delayer[T] {
	d.function = f
	return d
}

// WithSimpleFunction 设置无返回值的函数
func (d *Delayer[T]) WithSimpleFunction(f func()) *Delayer[T] {
	d.functionNoError = f
	return d
}

// WithTaskFunc 设置泛型任务函数
func (d *Delayer[T]) WithTaskFunc(taskFunc TaskFunc[T]) *Delayer[T] {
	d.taskFunc = taskFunc
	return d
}

// WithSimpleTaskFunc 设置简单任务函数（无返回值）
func (d *Delayer[T]) WithSimpleTaskFunc(taskFunc SimpleTaskFunc) *Delayer[T] {
	d.taskFunc = func(ctx *ExecutionContext) (T, error) {
		var zero T
		err := taskFunc(ctx)
		return zero, err
	}
	return d
}

// WithTimes 设置执行次数
func (d *Delayer[T]) WithTimes(count int) *Delayer[T] {
	if count <= 0 {
		count = 1
	}
	d.execCount = count
	return d
}

// WithStrategy 设置延迟策略
func (d *Delayer[T]) WithStrategy(strategy DelayStrategy) *Delayer[T] {
	d.strategy = strategy
	return d
}

// WithCustomDelay 设置自定义延迟函数
func (d *Delayer[T]) WithCustomDelay(delayFunc DelayFunc) *Delayer[T] {
	d.customDelayFunc = delayFunc
	d.strategy = CustomDelayStrategy
	return d
}

// WithRandomBase 设置随机基数（用于随机延迟策略）
func (d *Delayer[T]) WithRandomBase(base float64) *Delayer[T] {
	if base <= 0 {
		base = 1.0
	}
	d.randomBase = base
	return d
}

// WithMaxDelay 设置最大延迟时间
func (d *Delayer[T]) WithMaxDelay(maxDelay time.Duration) *Delayer[T] {
	d.maxDelay = maxDelay
	return d
}

// WithMultiplier 设置指数策略的倍数
func (d *Delayer[T]) WithMultiplier(multiplier float64) *Delayer[T] {
	if multiplier <= 1.0 {
		multiplier = 2.0
	}
	d.multiplier = multiplier
	return d
}

// WithContext 设置上下文
func (d *Delayer[T]) WithContext(ctx context.Context) *Delayer[T] {
	if d.cancelFunc != nil {
		d.cancelFunc()
	}
	d.ctx, d.cancelFunc = context.WithCancel(ctx)
	return d
}

// WithConcurrent 设置是否并发执行
func (d *Delayer[T]) WithConcurrent(concurrent bool) *Delayer[T] {
	d.concurrent = concurrent
	return d
}

// WithMaxConcurrency 设置最大并发数
func (d *Delayer[T]) WithMaxConcurrency(maxConcurrency int) *Delayer[T] {
	if maxConcurrency <= 0 {
		maxConcurrency = 10
	}
	d.maxConcurrency = maxConcurrency
	return d
}

// WithStopOnError 设置遇到错误是否停止
func (d *Delayer[T]) WithStopOnError(stopOnError bool) *Delayer[T] {
	d.stopOnError = stopOnError
	return d
}

// WithOnBeforeStart 设置开始执行前回调
func (d *Delayer[T]) WithOnBeforeStart(callback func(ctx *ExecutionContext)) *Delayer[T] {
	d.onBeforeStart = callback
	return d
}

// WithOnStart 设置开始执行回调（向后兼容）
func (d *Delayer[T]) WithOnStart(callback func(index int)) *Delayer[T] {
	d.onBeforeStart = func(ctx *ExecutionContext) {
		callback(ctx.Index)
	}
	return d
}

// WithOnAfterComplete 设置完成执行后回调
func (d *Delayer[T]) WithOnAfterComplete(callback func(ctx *ExecutionContext)) *Delayer[T] {
	d.onAfterComplete = callback
	return d
}

// WithOnComplete 设置完成执行回调（向后兼容）
func (d *Delayer[T]) WithOnComplete(callback func(result *ExecutionResult)) *Delayer[T] {
	d.onAfterComplete = func(ctx *ExecutionContext) {
		result := &ExecutionResult{
			Index:     ctx.Index,
			StartTime: ctx.StartTime,
			EndTime:   ctx.EndTime,
			Duration:  ctx.Duration,
			Error:     ctx.Error,
			Skipped:   ctx.Skipped,
		}
		callback(result)
	}
	return d
}

// WithOnError 设置错误处理回调（向后兼容）
func (d *Delayer[T]) WithOnError(callback func(index int, err error) bool) *Delayer[T] {
	d.onError = func(ctx *ExecutionContext) bool {
		return callback(ctx.Index, ctx.Error)
	}
	return d
}

// WithOnErrorContext 设置错误处理回调
func (d *Delayer[T]) WithOnErrorContext(callback ErrorHandlerFunc) *Delayer[T] {
	d.onError = callback
	return d
}

// WithOnProgress 设置进度回调
func (d *Delayer[T]) WithOnProgress(callback TaskProgressFunc) *Delayer[T] {
	d.onProgress = callback
	return d
}

// WithOnSuccess 设置成功回调
func (d *Delayer[T]) WithOnSuccess(callback CallbackFunc[T]) *Delayer[T] {
	d.onSuccess = callback
	return d
}

// WithDisableCallbacks 设置是否禁用回调以提升性能
func (d *Delayer[T]) WithDisableCallbacks(disable bool) *Delayer[T] {
	d.disableCallbacks = disable
	return d
}

// GetResults 获取所有泛型任务结果
func (d *Delayer[T]) GetResults() []T {
	d.genericResultsMu.RLock()
	defer d.genericResultsMu.RUnlock()
	results := make([]T, len(d.genericResults))
	copy(results, d.genericResults)
	return results
}

// GetResultChannel 获取结果通道（只读）
func (d *Delayer[T]) GetResultChannel() <-chan T {
	return d.resultChannel
}

// calculateDelay 根据策略计算延迟时间
func (d *Delayer[T]) calculateDelay(attempt int) time.Duration {
	switch d.strategy {
	case FixedDelayStrategy:
		return d.delay

	case LinearDelayStrategy:
		delay := d.delay * time.Duration(attempt+1)
		if d.maxDelay > 0 && delay > d.maxDelay {
			return d.maxDelay
		}
		return delay

	case ExponentialDelayStrategy:
		delay := time.Duration(float64(d.delay) * math.Pow(d.multiplier, float64(attempt)))
		if d.maxDelay > 0 && delay > d.maxDelay {
			return d.maxDelay
		}
		return delay

	case RandomDelayStrategy:
		min := float64(d.delay) / d.randomBase
		max := float64(d.delay) * d.randomBase
		randomDelay := time.Duration(min + rand.Float64()*(max-min))
		if d.maxDelay > 0 && randomDelay > d.maxDelay {
			return d.maxDelay
		}
		return randomDelay

	case CustomDelayStrategy:
		if d.customDelayFunc != nil {
			delay := d.customDelayFunc(attempt, d.delay)
			if d.maxDelay > 0 && delay > d.maxDelay {
				return d.maxDelay
			}
			return delay
		}
		return d.delay

	default:
		return d.delay
	}
}

// getExecutionContext 从对象池获取 ExecutionContext
func (d *Delayer[T]) getExecutionContext() *ExecutionContext {
	if d.disableCallbacks {
		// 当回调被禁用时，返回一个简单的上下文
		return &ExecutionContext{
			Metadata: make(map[string]interface{}),
		}
	}

	// 从对象池获取
	ctx := d.callbackPool.Get().(*ExecutionContext)

	// 重置状态
	ctx.Error = nil
	ctx.Skipped = false
	ctx.Cancelled = false
	ctx.Retryable = true
	for k := range ctx.Metadata {
		delete(ctx.Metadata, k)
	}

	return ctx
}

// putExecutionContext 归还 ExecutionContext 到对象池
func (d *Delayer[T]) putExecutionContext(ctx *ExecutionContext) {
	if !d.disableCallbacks {
		d.callbackPool.Put(ctx)
	}
}

// Close 关闭结果通道
func (d *Delayer[T]) Close() {
	// 使用原子操作设置关闭状态
	if atomic.CompareAndSwapInt64(&d.channelClosed, 0, 1) {
		close(d.resultChannel)
	}
}

// safeChannelSend 安全地向通道发送结果（使用原子操作优化性能）
func (d *Delayer[T]) safeChannelSend(result T) {
	// 使用原子操作检查通道状态
	if atomic.LoadInt64(&d.channelClosed) == 0 {
		// 使用 select 防止阻塞，如果通道在发送过程中被关闭也不会 panic
		select {
		case d.resultChannel <- result:
			// 成功发送
		default:
			// 通道已满，丢弃结果
		}
	}
}

// Execute 执行延迟任务
func (d *Delayer[T]) Execute() error {
	if d.function == nil && d.functionNoError == nil && d.taskFunc == nil {
		return fmt.Errorf("no function to execute")
	}

	// 执行延迟任务
	if d.taskFunc != nil {
		// 如果有泛型任务函数，直接执行
		return d.executeSequentially()
	}

	d.stats.StartTime = time.Now()
	defer func() {
		d.stats.EndTime = time.Now()
		d.stats.TotalDuration = d.stats.EndTime.Sub(d.stats.StartTime)
	}()

	if d.concurrent {
		return d.executeConcurrently()
	}
	return d.executeSequentially()
}

// executeSequentially 顺序执行
func (d *Delayer[T]) executeSequentially() error {
	for i := 0; i < d.execCount; i++ {
		if atomic.LoadInt64(&d.stopped) == 1 {
			break
		}

		// 检查上下文是否被取消
		select {
		case <-d.ctx.Done():
			return d.ctx.Err()
		default:
		}

		// 延迟执行
		if d.delay > 0 || d.strategy != FixedDelayStrategy {
			delay := d.calculateDelay(i)
			if delay > 0 {
				timer := time.NewTimer(delay)
				d.timers = append(d.timers, timer)

				select {
				case <-timer.C:
				case <-d.ctx.Done():
					timer.Stop()
					return d.ctx.Err()
				}
			}
		}

		// 执行任务
		err := d.executeTask(i)
		if err != nil {
			atomic.AddInt64(&d.stats.ErrorCount, 1)
			if d.stopOnError {
				return err
			}
		} else {
			atomic.AddInt64(&d.stats.SuccessCount, 1)
		}

		// 更新进度
		if d.onProgress != nil && !d.disableCallbacks {
			completed := int64(i + 1)
			total := int64(d.execCount)
			percentage := float64(completed) / float64(total) * 100
			d.onProgress(completed, total, percentage)
		}
	}

	// 通知任务完成
	close(d.completionChan)
	return nil
}

// executeConcurrently 并发执行
func (d *Delayer[T]) executeConcurrently() error {
	semaphore := make(chan struct{}, d.maxConcurrency)
	errChan := make(chan error, d.execCount)
	var wg sync.WaitGroup

	atomic.StoreInt64(&d.pendingTasks, int64(d.execCount))

	for i := 0; i < d.execCount; i++ {
		if atomic.LoadInt64(&d.stopped) == 1 {
			break
		}

		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			defer atomic.AddInt64(&d.pendingTasks, -1)

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 增加运行计数
			atomic.AddInt64(&d.running, 1)
			defer atomic.AddInt64(&d.running, -1)

			// 检查上下文是否被取消
			select {
			case <-d.ctx.Done():
				errChan <- d.ctx.Err()
				return
			default:
			}

			// 延迟执行
			if d.delay > 0 || d.strategy != FixedDelayStrategy {
				delay := d.calculateDelay(index)
				if delay > 0 {
					timer := time.NewTimer(delay)
					select {
					case <-timer.C:
					case <-d.ctx.Done():
						timer.Stop()
						errChan <- d.ctx.Err()
						return
					}
				}
			}

			// 执行任务
			err := d.executeTask(index)
			if err != nil {
				atomic.AddInt64(&d.stats.ErrorCount, 1)
				if d.stopOnError {
					errChan <- err
					return
				}
			} else {
				atomic.AddInt64(&d.stats.SuccessCount, 1)
			}

			// 更新进度
			if d.onProgress != nil && !d.disableCallbacks {
				completed := atomic.LoadInt64(&d.stats.SuccessCount) + atomic.LoadInt64(&d.stats.ErrorCount)
				total := int64(d.execCount)
				percentage := float64(completed) / float64(total) * 100
				d.onProgress(completed, total, percentage)
			}
		}(i)
	}

	// 等待所有任务完成
	go func() {
		wg.Wait()
		close(errChan)
		close(d.completionChan)
	}()

	// 收集错误
	var lastErr error
	for err := range errChan {
		if err != nil {
			lastErr = err
			if d.stopOnError {
				atomic.StoreInt64(&d.stopped, 1)
			}
		}
	}

	return lastErr
}

// executeTask 执行单个任务
func (d *Delayer[T]) executeTask(index int) error {
	ctx := d.getExecutionContext()
	defer d.putExecutionContext(ctx)

	ctx.Index = index
	ctx.Attempt = 0
	ctx.StartTime = time.Now()
	ctx.Strategy = d.strategy
	ctx.Concurrent = d.concurrent

	// 前置回调
	if d.onBeforeStart != nil && !d.disableCallbacks {
		d.onBeforeStart(ctx)
	}

	var err error
	// 优先执行泛型任务函数
	if d.taskFunc != nil {
		result, taskErr := d.taskFunc(ctx)
		err = taskErr
		if err == nil {
			// 使用预分配的切片以减少内存分配
			d.genericResultsMu.Lock()
			d.genericResults = append(d.genericResults, result)
			d.genericResultsMu.Unlock()

			// 高性能的通道发送
			d.safeChannelSend(result)

			// 调用成功回调
			if d.onSuccess != nil && !d.disableCallbacks {
				d.onSuccess(ctx, result)
			}
		}
	} else if d.function != nil {
		err = d.function()
	} else if d.functionNoError != nil {
		d.functionNoError()
	}

	ctx.EndTime = time.Now()
	ctx.Duration = ctx.EndTime.Sub(ctx.StartTime)
	ctx.Error = err

	// 如果有错误，尝试调用错误处理回调
	if err != nil && d.onError != nil && !d.disableCallbacks {
		shouldContinue := d.onError(ctx)
		if !shouldContinue {
			return err
		}
	}

	// 后置回调
	if d.onAfterComplete != nil && !d.disableCallbacks {
		d.onAfterComplete(ctx)
	}

	// 保存执行结果
	result := &ExecutionResult{
		Index:     index,
		StartTime: ctx.StartTime,
		EndTime:   ctx.EndTime,
		Duration:  ctx.Duration,
		Error:     err,
		Skipped:   ctx.Skipped,
	}

	d.mu.Lock()
	d.results = append(d.results, result)
	d.mu.Unlock()

	return err
}

// Stop 停止执行
func (d *Delayer[T]) Stop() {
	atomic.StoreInt64(&d.stopped, 1)
	if d.cancelFunc != nil {
		d.cancelFunc()
	}

	// 停止所有计时器
	d.mu.Lock()
	for _, timer := range d.timers {
		timer.Stop()
	}
	d.timers = d.timers[:0]
	d.mu.Unlock()
}

// IsRunning 检查是否正在运行
func (d *Delayer[T]) IsRunning() bool {
	return atomic.LoadInt64(&d.running) > 0
}

// GetStats 获取执行统计
func (d *Delayer[T]) GetStats() *ExecutionStats {
	return d.stats
}

// Wait 等待上下文取消
func (d *Delayer[T]) Wait() error {
	<-d.ctx.Done()
	return d.ctx.Err()
}

// WaitForCompletion 等待任务完成
func (d *Delayer[T]) WaitForCompletion() {
	<-d.completionChan
}

// GetLegacyResults 获取所有执行结果（向后兼容）
func (d *Delayer[T]) GetLegacyResults() []*ExecutionResult {
	d.mu.RLock()
	defer d.mu.RUnlock()
	results := make([]*ExecutionResult, len(d.results))
	copy(results, d.results)
	return results
}
