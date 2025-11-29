/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-29 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-29 13:56:14
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
	Name           string                          // 任务名称
	Interval       time.Duration                   // 执行间隔
	ExecuteFunc    func(ctx context.Context) error // 执行函数
	ImmediateStart bool                            // 是否立即执行首次任务
	OnError        func(name string, err error)    // 错误处理回调
	OnStart        func(name string)               // 启动回调
	OnStop         func(name string)               // 停止回调
}

// PeriodicTaskManager 周期性任务管理器
type PeriodicTaskManager struct {
	tasks               []*PeriodicTask
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
		tasks: make([]*PeriodicTask, 0),
	}
}

// AddTask 添加周期性任务
func (m *PeriodicTaskManager) AddTask(task *PeriodicTask) *PeriodicTaskManager {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 应用默认处理器
	if task.OnError == nil && m.defaultErrorHandler != nil {
		task.OnError = m.defaultErrorHandler
	}
	if task.OnStart == nil && m.defaultOnStart != nil {
		task.OnStart = m.defaultOnStart
	}
	if task.OnStop == nil && m.defaultOnStop != nil {
		task.OnStop = m.defaultOnStop
	}

	m.tasks = append(m.tasks, task)
	return m
}

// AddSimpleTask 添加简单的周期性任务
func (m *PeriodicTaskManager) AddSimpleTask(name string, interval time.Duration, executeFunc func(ctx context.Context) error) *PeriodicTaskManager {
	task := &PeriodicTask{
		Name:        name,
		Interval:    interval,
		ExecuteFunc: executeFunc,
	}
	return m.AddTask(task)
}

// AddTaskWithImmediateStart 添加立即执行的周期性任务
func (m *PeriodicTaskManager) AddTaskWithImmediateStart(name string, interval time.Duration, executeFunc func(ctx context.Context) error) *PeriodicTaskManager {
	task := &PeriodicTask{
		Name:           name,
		Interval:       interval,
		ExecuteFunc:    executeFunc,
		ImmediateStart: true,
	}
	return m.AddTask(task)
}

// SetDefaultErrorHandler 设置默认错误处理器
func (m *PeriodicTaskManager) SetDefaultErrorHandler(handler func(name string, err error)) *PeriodicTaskManager {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.defaultErrorHandler = handler

	// 为已有任务设置默认处理器
	for _, task := range m.tasks {
		if task.OnError == nil {
			task.OnError = handler
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
		if task.OnStart == nil {
			task.OnStart = onStart
		}
		if task.OnStop == nil {
			task.OnStop = onStop
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

	// 调用启动回调
	if task.OnStart != nil {
		task.OnStart(task.Name)
	}

	// 处理非正数间隔
	interval := task.Interval
	if interval <= 0 {
		interval = time.Millisecond // 最小间隔为1毫秒
	}

	// 创建定时器
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 如果需要立即执行
	if task.ImmediateStart {
		m.executeTask(task)
	}

	// 主循环
	for {
		select {
		case <-m.ctx.Done():
			// 调用停止回调
			if task.OnStop != nil {
				task.OnStop(task.Name)
			}
			return
		case <-ticker.C:
			m.executeTask(task)
		}
	}
} // executeTask 执行单个任务
func (m *PeriodicTaskManager) executeTask(task *PeriodicTask) {
	defer func() {
		if r := recover(); r != nil {
			// panic恢复：如果有错误处理器，将panic转换为错误
			if task.OnError != nil {
				err := fmt.Errorf("task panic: %v", r)
				task.OnError(task.Name, err)
			}
		}
	}()

	if err := task.ExecuteFunc(m.ctx); err != nil {
		if task.OnError != nil {
			task.OnError(task.Name, err)
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
		names[i] = task.Name
	}
	return names
}

// Wait 等待所有任务完成
func (m *PeriodicTaskManager) Wait() {
	m.wg.Wait()
}
