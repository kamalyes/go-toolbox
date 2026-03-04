/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-04 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-04 12:00:00
 * @FilePath: \go-toolbox\pkg\syncx\dynamic_ticker.go
 * @Description: 动态定时器，支持运行时调整频率，可与 EventLoop 配合使用
 *
 * 使用说明:
 *
 * 1. 基础用法:
 *    ticker := NewDynamicTicker(1 * time.Second)
 *    ticker.Start()
 *    defer ticker.Stop()
 *
 *    for t := range ticker.C {
 *        // 处理定时事件
 *        fmt.Println("tick at", t)
 *    }
 *
 * 2. 动态调整频率:
 *    ticker := NewDynamicTicker(1 * time.Second)
 *    ticker.Start()
 *
 *    // 运行时调整为 2 秒
 *    ticker.UpdateInterval(2 * time.Second)
 *
 * 3. 与 EventLoop 配合使用:
 *    ticker := NewDynamicTicker(1 * time.Second)
 *    ticker.Start()
 *    defer ticker.Stop()
 *
 *    NewEventLoop(ctx).
 *        OnChannel(ticker.C, func(t time.Time) {
 *            // 处理定时事件
 *            // 可以根据条件动态调整频率
 *            if needFaster {
 *                ticker.UpdateInterval(500 * time.Millisecond)
 *            }
 *        }).
 *        Run()
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"sync"
	"time"
)

// DynamicTicker 动态定时器，支持运行时调整频率
//
// 特性：
//   - 支持运行时动态调整定时器间隔
//   - 线程安全，可在多个 goroutine 中调用
//   - 非阻塞发送，不会因为接收方慢而阻塞定时器
//   - 可与 EventLoop 无缝集成
//
// 注意事项：
//   - 必须调用 Start() 才能开始工作
//   - 使用完毕后应调用 Stop() 释放资源
//   - C 通道有 1 个缓冲，如果接收方处理慢会丢弃部分 tick
type DynamicTicker struct {
	mu       sync.RWMutex   // 保护 interval、ticker、running 字段
	interval time.Duration  // 当前定时器间隔
	ticker   *time.Ticker   // 底层定时器
	C        chan time.Time // 对外暴露的通道，接收定时事件
	quit     chan struct{}  // 停止信号通道
	running  bool           // 定时器是否正在运行
}

// NewDynamicTicker 创建动态定时器
// 注意：
//   - 创建后需要调用 Start() 才能开始工作
//   - 定时器不会自动启动
func NewDynamicTicker(interval time.Duration) *DynamicTicker {
	dt := &DynamicTicker{
		interval: interval,
		C:        make(chan time.Time, 1), // 1 个缓冲，避免短暂阻塞
		quit:     make(chan struct{}),
	}
	return dt
}

// Start 启动定时器
//
// 功能：
//   - 创建底层 time.Ticker 并启动后台 goroutine
//   - 如果已经启动，重复调用不会有任何效果
//
// 注意：
//   - 线程安全，可以在多个 goroutine 中调用
//   - 启动后会立即开始按照 interval 发送 tick 事件
func (dt *DynamicTicker) Start() {
	dt.mu.Lock()
	if dt.running {
		dt.mu.Unlock()
		return
	}
	dt.running = true
	dt.ticker = time.NewTicker(dt.interval)
	dt.mu.Unlock()

	go dt.run()
}

// run 运行定时器循环（内部方法）
//
// 功能：
//   - 在独立的 goroutine 中运行
//   - 从底层 ticker 接收事件并转发到 C 通道
//   - 使用非阻塞发送，避免因接收方慢而阻塞定时器
//
// 退出条件：
//   - 收到 quit 信号时退出
func (dt *DynamicTicker) run() {
	for {
		select {
		case <-dt.quit:
			return
		case t := <-dt.ticker.C:
			// 非阻塞发送，避免阻塞定时器
			// 如果 C 通道已满（接收方处理慢），则丢弃本次 tick
			select {
			case dt.C <- t:
			default:
			}
		}
	}
}

// UpdateInterval 动态更新定时器间隔
// 功能：
//   - 如果新间隔与当前间隔相同，不做任何操作
//   - 如果定时器正在运行，会立即应用新的间隔
//   - 使用 ticker.Reset() 实现无缝切换，不会丢失正在进行的 tick
//
// 注意：
//   - 线程安全，可以在任何时候调用
//   - 可以在定时器启动前或运行中调用
func (dt *DynamicTicker) UpdateInterval(interval time.Duration) {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	if dt.interval == interval || interval <= 0 {
		return
	}

	dt.interval = interval

	if dt.ticker != nil {
		dt.ticker.Reset(interval)
	}
}

// GetInterval 获取当前间隔
func (dt *DynamicTicker) GetInterval() time.Duration {
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	return dt.interval
}

// Stop 停止定时器
func (dt *DynamicTicker) Stop() {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	if !dt.running {
		return
	}

	dt.running = false
	close(dt.quit)

	if dt.ticker != nil {
		dt.ticker.Stop()
	}
}
