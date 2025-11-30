/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-29 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-30 23:40:10
 * @FilePath: \engine-im-service\go-toolbox\pkg\syncx\periodic_task.go
 * @Description: 周期性任务管理器 - 用于管理多个定时执行的任务
 *
 * 功能特性：
 * - 支持多个周期性任务的并发执行
 * - 统一的错误处理和日志记录
 * - 优雅的启动和停止机制
 * - 支持任务立即执行选项
 * - 自动资源清理和上下文管理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package syncx

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PeriodicTask 表示一个周期性任务
type PeriodicTask struct {
	name             string                          // 任务名称
	interval         time.Duration                   // 执行间隔
	executeFunc      func(ctx context.Context) error // 执行函数
	immediateStart   bool                            // 是否立即执行首次任务
	preventOverlap   bool                            // 是否防止任务重叠执行
	onError          func(name string, err error)    // 错误处理回调
	onStart          func(name string)               // 启动回调
	onStop           func(name string)               // 停止回调
	onOverlapSkipped func(name string)               // 重叠跳过回调

	// 内部字段（重叠保护和取消控制）
	executeMutex sync.Mutex         // 执行保护锁
	isExecuting  bool               // 是否正在执行
	cancelFunc   context.CancelFunc // 任务取消函数
	taskCtx      context.Context    // 任务专用上下文
}

// PeriodicTaskManager 周期性任务管理器
type PeriodicTaskManager struct {
	tasks               []*PeriodicTask
	taskMap             map[string]*PeriodicTask // 任务名称到任务的映射
	ctx                 context.Context
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
	isRunning           bool
	mu                  sync.RWMutex
	defaultErrorHandler func(name string, err error)
	defaultOnStart      func(name string)
	defaultOnStop       func(name string)
}

// NewPeriodicTaskManager 创建新的周期性任务管理器
func NewPeriodicTaskManager() *PeriodicTaskManager {
	return &PeriodicTaskManager{
		tasks:   make([]*PeriodicTask, 0),
		taskMap: make(map[string]*PeriodicTask),
	}
}

// NewPeriodicTask 创建新的周期性任务
func NewPeriodicTask(name string, interval time.Duration, executeFunc func(ctx context.Context) error) *PeriodicTask {
	return &PeriodicTask{
		name:        name,
		interval:    interval,
		executeFunc: executeFunc,
	}
}

// SetImmediateStart 设置是否立即执行首次任务
func (t *PeriodicTask) SetImmediateStart(immediateStart bool) *PeriodicTask {
	t.immediateStart = immediateStart
	return t
}

// SetPreventOverlap 设置是否防止任务重叠执行
func (t *PeriodicTask) SetPreventOverlap(preventOverlap bool) *PeriodicTask {
	t.preventOverlap = preventOverlap
	return t
}

// SetOnError 设置错误处理回调
func (t *PeriodicTask) SetOnError(onError func(name string, err error)) *PeriodicTask {
	t.onError = onError
	return t
}

// SetOnStart 设置启动回调
func (t *PeriodicTask) SetOnStart(onStart func(name string)) *PeriodicTask {
	t.onStart = onStart
	return t
}

// SetOnStop 设置停止回调
func (t *PeriodicTask) SetOnStop(onStop func(name string)) *PeriodicTask {
	t.onStop = onStop
	return t
}

// SetOnOverlapSkipped 设置重叠跳过回调
func (t *PeriodicTask) SetOnOverlapSkipped(onOverlapSkipped func(name string)) *PeriodicTask {
	t.onOverlapSkipped = onOverlapSkipped
	return t
}

// AddTask 添加周期性任务
func (m *PeriodicTaskManager) AddTask(task *PeriodicTask) *PeriodicTaskManager {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 应用默认处理器
	if task.onError == nil && m.defaultErrorHandler != nil {
		task.onError = m.defaultErrorHandler
	}
	if task.onStart == nil && m.defaultOnStart != nil {
		task.onStart = m.defaultOnStart
	}
	if task.onStop == nil && m.defaultOnStop != nil {
		task.onStop = m.defaultOnStop
	}

	m.tasks = append(m.tasks, task)
	// 同时维护任务名称映射
	m.taskMap[task.name] = task
	return m
}

// RemoveTask 移除指定名称的任务
func (m *PeriodicTaskManager) RemoveTask(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查任务是否存在
	task, exists := m.taskMap[name]
	if !exists {
		return false
	}

	// 如果任务正在运行，先取消它
	if task.cancelFunc != nil {
		task.cancelFunc() // 取消任务上下文
	}

	// 等待正在执行的任务完成（带超时）
	if task.preventOverlap && task.IsExecuting() {
		// 异步等待任务完成，避免死锁
		go func() {
			timeout := time.NewTimer(10 * time.Second)
			defer timeout.Stop()

			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-timeout.C:
					// 超时，强制认为任务已完成
					return
				case <-ticker.C:
					if !task.IsExecuting() {
						return
					}
				}
			}
		}()
	}

	// 从map中删除
	delete(m.taskMap, name)

	// 从slice中删除
	for i, t := range m.tasks {
		if t.name == name {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			break
		}
	}
	return true
}

// RemoveTaskWithTimeout 移除指定名称的任务（带超时等待）
func (m *PeriodicTaskManager) RemoveTaskWithTimeout(name string, timeout time.Duration) bool {
	m.mu.Lock()

	// 检查任务是否存在
	task, exists := m.taskMap[name]
	if !exists {
		m.mu.Unlock()
		return false
	}

	// 如果任务正在运行，先取消它
	if task.cancelFunc != nil {
		task.cancelFunc()
	}

	m.mu.Unlock()

	// 如果需要等待正在执行的任务完成
	if task.preventOverlap && task.IsExecuting() {
		timeoutTimer := time.NewTimer(timeout)
		defer timeoutTimer.Stop()

		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-timeoutTimer.C:
				// 超时，继续移除操作
				goto removeTask
			case <-ticker.C:
				if !task.IsExecuting() {
					goto removeTask
				}
			}
		}
	}

removeTask:
	// 重新获取锁进行删除操作
	m.mu.Lock()
	defer m.mu.Unlock()

	// 再次检查任务是否还存在（防止并发删除）
	if _, exists := m.taskMap[name]; !exists {
		return false
	}

	// 从map中删除
	delete(m.taskMap, name)

	// 从slice中删除
	for i, t := range m.tasks {
		if t.name == name {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			break
		}
	}
	return true
}

// ClearAllTasks 清除所有任务
func (m *PeriodicTaskManager) ClearAllTasks() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tasks = m.tasks[:0]                      // 清空slice但保留容量
	m.taskMap = make(map[string]*PeriodicTask) // 重新创建map
}

// AddSimpleTask 添加简单的周期性任务
func (m *PeriodicTaskManager) AddSimpleTask(name string, interval time.Duration, executeFunc func(ctx context.Context) error) *PeriodicTaskManager {
	task := NewPeriodicTask(name, interval, executeFunc)
	return m.AddTask(task)
}

// AddTaskWithImmediateStart 添加立即执行的周期性任务
func (m *PeriodicTaskManager) AddTaskWithImmediateStart(name string, interval time.Duration, executeFunc func(ctx context.Context) error) *PeriodicTaskManager {
	task := NewPeriodicTask(name, interval, executeFunc).SetImmediateStart(true)
	return m.AddTask(task)
}

// AddTaskWithOverlapPrevention 添加防重叠执行的周期性任务
func (m *PeriodicTaskManager) AddTaskWithOverlapPrevention(name string, interval time.Duration, executeFunc func(ctx context.Context) error) *PeriodicTaskManager {
	task := NewPeriodicTask(name, interval, executeFunc).SetPreventOverlap(true)
	return m.AddTask(task)
}

// AddTaskWithOverlapPreventionAndCallback 添加防重叠执行的周期性任务（带重叠跳过回调）
func (m *PeriodicTaskManager) AddTaskWithOverlapPreventionAndCallback(
	name string,
	interval time.Duration,
	executeFunc func(ctx context.Context) error,
	onOverlapSkipped func(name string),
) *PeriodicTaskManager {
	task := NewPeriodicTask(name, interval, executeFunc).
		SetPreventOverlap(true).
		SetOnOverlapSkipped(onOverlapSkipped)
	return m.AddTask(task)
}

// AddTaskWithOverlapPreventionImmediateAndCallback 添加防重叠执行的周期性任务（带立即执行和重叠跳过回调）
func (m *PeriodicTaskManager) AddTaskWithOverlapPreventionImmediateAndCallback(
	name string,
	interval time.Duration,
	executeFunc func(ctx context.Context) error,
	onOverlapSkipped func(name string),
) *PeriodicTaskManager {
	task := NewPeriodicTask(name, interval, executeFunc).
		SetPreventOverlap(true).
		SetImmediateStart(true).
		SetOnOverlapSkipped(onOverlapSkipped)
	return m.AddTask(task)
}

// SetDefaultErrorHandler 设置默认错误处理器
func (m *PeriodicTaskManager) SetDefaultErrorHandler(handler func(name string, err error)) *PeriodicTaskManager {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.defaultErrorHandler = handler

	// 为已有任务设置默认处理器
	for _, task := range m.tasks {
		if task.onError == nil {
			task.onError = handler
		}
	}
	return m
}

// SetDefaultCallbacks 设置默认回调函数
func (m *PeriodicTaskManager) SetDefaultCallbacks(
	onStart func(name string),
	onStop func(name string),
) *PeriodicTaskManager {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.defaultOnStart = onStart
	m.defaultOnStop = onStop

	// 为已有任务设置默认回调
	for _, task := range m.tasks {
		if task.onStart == nil {
			task.onStart = onStart
		}
		if task.onStop == nil {
			task.onStop = onStop
		}
	}
	return m
}

// Start 启动所有周期性任务
func (m *PeriodicTaskManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("periodic task manager is already running")
	}

	// 创建上下文
	m.ctx, m.cancel = context.WithCancel(context.Background())

	// 启动每个任务
	for _, task := range m.tasks {
		m.wg.Add(1)
		go m.runTask(task)
	}

	m.isRunning = true
	return nil
}

// StartWithContext 使用指定上下文启动所有周期性任务
func (m *PeriodicTaskManager) StartWithContext(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("periodic task manager is already running")
	}

	// 使用传入的上下文创建子上下文
	m.ctx, m.cancel = context.WithCancel(ctx)

	// 启动每个任务
	for _, task := range m.tasks {
		m.wg.Add(1)
		go m.runTask(task)
	}

	m.isRunning = true
	return nil
}

// runTask 运行单个周期性任务
func (m *PeriodicTaskManager) runTask(task *PeriodicTask) {
	defer m.wg.Done()

	// 为任务创建独立的上下文，支持单独取消
	task.taskCtx, task.cancelFunc = context.WithCancel(m.ctx)

	// 调用启动回调
	if task.onStart != nil {
		task.onStart(task.name)
	}

	// 处理非正数间隔
	interval := task.interval
	if interval <= 0 {
		interval = time.Millisecond // 最小间隔为1毫秒
	}

	// 创建定时器
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 如果需要立即执行
	if task.immediateStart {
		m.executeTask(task)
	}

	// 主循环
	for {
		select {
		case <-m.ctx.Done():
			// 全局管理器停止
			if task.onStop != nil {
				task.onStop(task.name)
			}
			return
		case <-task.taskCtx.Done():
			// 单个任务被取消
			if task.onStop != nil {
				task.onStop(task.name)
			}
			return
		case <-ticker.C:
			// 每次 tick 都尝试执行任务，在 executeTask 中处理重叠保护
			go m.executeTask(task) // 使用 goroutine 避免阻塞 ticker
		}
	}
}

// executeTask 执行单个任务
func (m *PeriodicTaskManager) executeTask(task *PeriodicTask) {
	// 检查任务上下文是否已被取消
	if task.taskCtx != nil && task.taskCtx.Err() != nil {
		return // 任务已被取消，直接返回
	}

	// 🔒 重叠保护检查
	if task.preventOverlap {
		task.executeMutex.Lock()
		if task.isExecuting {
			task.executeMutex.Unlock()
			// 调用重叠跳过回调
			if task.onOverlapSkipped != nil {
				task.onOverlapSkipped(task.name)
			}
			return
		}
		task.isExecuting = true
		task.executeMutex.Unlock()

		// 🎯 确保执行完成后重置状态
		defer func() {
			task.executeMutex.Lock()
			task.isExecuting = false
			task.executeMutex.Unlock()
		}()
	}

	defer func() {
		if r := recover(); r != nil {
			// panic恢复：如果有错误处理器，将panic转换为错误
			if task.onError != nil {
				err := fmt.Errorf("task panic: %v", r)
				task.onError(task.name, err)
			}
		}
	}()

	// 使用任务专用的上下文执行任务
	ctx := m.ctx
	if task.taskCtx != nil {
		ctx = task.taskCtx
	}

	if err := task.executeFunc(ctx); err != nil {
		if task.onError != nil {
			task.onError(task.name, err)
		}
	}
}

// Stop 停止所有周期性任务
func (m *PeriodicTaskManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil
	}

	// 取消上下文
	if m.cancel != nil {
		m.cancel()
	}

	// 等待所有任务完成
	m.wg.Wait()

	m.isRunning = false
	return nil
}

// StopWithTimeout 在指定超时时间内停止所有周期性任务
func (m *PeriodicTaskManager) StopWithTimeout(timeout time.Duration) error {
	done := make(chan error, 1)

	go func() {
		done <- m.Stop()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("failed to stop periodic task manager within timeout %v", timeout)
	}
}

// IsRunning 检查任务管理器是否正在运行
func (m *PeriodicTaskManager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isRunning
}

// Wait 等待所有任务完成
func (m *PeriodicTaskManager) Wait() {
	m.wg.Wait()
}

// ======================== PeriodicTask Getter Methods ========================

// GetName 获取任务名称
func (t *PeriodicTask) GetName() string {
	return t.name
}

// GetInterval 获取执行间隔
func (t *PeriodicTask) GetInterval() time.Duration {
	return t.interval
}

// GetExecuteFunc 获取执行函数
func (t *PeriodicTask) GetExecuteFunc() func(ctx context.Context) error {
	return t.executeFunc
}

// GetImmediateStart 获取是否立即执行标志
func (t *PeriodicTask) GetImmediateStart() bool {
	return t.immediateStart
}

// GetPreventOverlap 获取是否防止重叠执行标志
func (t *PeriodicTask) GetPreventOverlap() bool {
	return t.preventOverlap
}

// GetOnError 获取错误处理回调
func (t *PeriodicTask) GetOnError() func(name string, err error) {
	return t.onError
}

// GetOnStart 获取启动回调
func (t *PeriodicTask) GetOnStart() func(name string) {
	return t.onStart
}

// GetOnStop 获取停止回调
func (t *PeriodicTask) GetOnStop() func(name string) {
	return t.onStop
}

// GetOnOverlapSkipped 获取重叠跳过回调
func (t *PeriodicTask) GetOnOverlapSkipped() func(name string) {
	return t.onOverlapSkipped
}

// IsExecuting 获取当前是否正在执行状态
func (t *PeriodicTask) IsExecuting() bool {
	t.executeMutex.Lock()
	defer t.executeMutex.Unlock()
	return t.isExecuting
}

// GetTaskCount 获取任务数量
func (m *PeriodicTaskManager) GetTaskCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.tasks)
}

// GetTaskNames 获取所有任务名称
func (m *PeriodicTaskManager) GetTaskNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, len(m.tasks))
	for i, task := range m.tasks {
		names[i] = task.name
	}
	return names
}

// TaskDetailInfo 任务详细信息
type TaskDetailInfo struct {
	Name           string        `json:"name"`
	Interval       time.Duration `json:"interval"`
	ImmediateStart bool          `json:"immediate_start"`
	PreventOverlap bool          `json:"prevent_overlap"`
	IsExecuting    bool          `json:"is_executing"`
}

// GetTaskDetails 获取任务详细信息
// 如果name为空，返回所有任务；如果指定name，返回匹配的任务
func (m *PeriodicTaskManager) GetTaskDetails(name ...string) []TaskDetailInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 辅助函数：构建任务详情
	buildTaskDetail := func(task *PeriodicTask) TaskDetailInfo {
		return TaskDetailInfo{
			Name:           task.name,
			Interval:       task.interval,
			ImmediateStart: task.immediateStart,
			PreventOverlap: task.preventOverlap,
			IsExecuting:    task.IsExecuting(),
		}
	}

	// 指定了name，直接从map查找 - O(1)查找
	if len(name) > 0 && name[0] != "" {
		if task, exists := m.taskMap[name[0]]; exists {
			return []TaskDetailInfo{buildTaskDetail(task)}
		}
		return []TaskDetailInfo{} // 没找到
	}

	// 没有指定name，返回所有任务
	details := make([]TaskDetailInfo, 0, len(m.tasks))
	for _, task := range m.tasks {
		details = append(details, buildTaskDetail(task))
	}
	return details
}
