/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-09 19:15:01
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-07 17:25:07
 * @FilePath: \go-toolbox\pkg\moment\timer.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package moment

import (
	"fmt"
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// TimerInterface 定义计时器的接口
type TimerInterface interface {
	Run()                       // 启动计时器
	Pause()                     // 暂停计时器
	Resume()                    // 恢复计时器
	Finish()                    // 完成停止计时器
	GetTraceId() string         // 获取跟踪ID
	GetStartTime() time.Time    // 获取开始时间
	GetEndTime() time.Time      // 获取结束时间
	GetDuration() time.Duration // 获取持续时间
}

// Timer 结构体用于表示一个计时器
type Timer struct {
	traceId       string        // 计时器的唯一标识符
	startTime     time.Time     // 计时器开始时间
	endTime       time.Time     // 计时器结束时间
	duration      time.Duration // 持续时间
	paused        bool          // 计时器是否已暂停
	pauseStart    time.Time
	pauseDuration time.Duration // 暂停的总时长
	stopChan      chan struct{} // 用于停止计时器的通道
	mu            sync.RWMutex  // 保护共享数据的互斥锁(读锁和写锁)
}

// NewTimer 创建一个新的计时器
func NewTimer() *Timer {
	t := &Timer{
		stopChan: make(chan struct{}),           // 初始化停止通道
		traceId:  osx.HashUnixMicroCipherText(), // 初始化跟踪Id
	}
	return t
}

// NewTimerWithTraceId 创建一个带有自定义跟踪Id的新计时器
func NewTimerWithTraceId(traceId string) *Timer {
	return NewTimer().SetTraceId(traceId)
}

// Run 启动计时器，开始计时并打印时差
func (t *Timer) Run() {
	// 检查是否已暂停或已开始
	if t.paused || !t.startTime.IsZero() {
		return
	}
	syncx.WithLock(&t.mu, func() {
		t.startTime = time.Now() // 记录开始时间
		t.endTime = t.startTime  // 初始化结束时间
		go func() {
			// 只需在 goroutine 中执行一次打印
			select {
			case <-t.stopChan: // 如果接收到停止信号，退出并打印持续时间
				t.PrintLog()
				return
			}
		}()
	})
}

// Pause 暂停计时器
func (t *Timer) Pause() {
	syncx.WithLock(&t.mu, func() {
		// 检查是否已暂停或未开始
		if t.paused || t.startTime.IsZero() {
			return
		}
		t.pauseStart = time.Now() // 记录暂停开始时间
		t.paused = true
	})
}

// Resume 恢复计时器
func (t *Timer) Resume() {
	syncx.WithLock(&t.mu, func() {
		// 检查是否未暂停或未开始
		if !t.paused || t.startTime.IsZero() {
			return
		}
		t.pauseDuration += time.Since(t.pauseStart) // 更新总暂停时间
		t.paused = false
	})
}

// Finish 停止计时器
func (t *Timer) Finish() {
	syncx.WithLock(&t.mu, func() {
		// 检查是否未开始
		if t.startTime.IsZero() {
			return
		}
		t.endTime = time.Now()
		t.duration = t.endTime.Sub(t.startTime) - t.pauseDuration
		close(t.stopChan)
	})
}

// SetTraceId 设置计时器的跟踪Id
func (t *Timer) SetTraceId(traceId string) *Timer {
	return syncx.WithLockReturnValue(&t.mu, func() *Timer {
		t.traceId = traceId
		return t
	})
}

// GetTraceId 获取计时器的跟踪Id
func (t *Timer) GetTraceId() string {
	return syncx.WithRLockReturnValue(&t.mu, func() string {
		return t.traceId
	})
}

// GetStartTime 获取计时器开始时间
func (t *Timer) GetStartTime() time.Time {
	return syncx.WithRLockReturnValue(&t.mu, func() time.Time {
		return t.startTime
	})
}

// GetEndTime 获取计时器的结束时间
func (t *Timer) GetEndTime() time.Time {
	return syncx.WithRLockReturnValue(&t.mu, func() time.Time {
		return t.endTime
	})
}

// GetDuration 获取计时器的实际持续时间
func (t *Timer) GetDuration() time.Duration {
	return syncx.WithRLockReturnValue(&t.mu, func() time.Duration {
		return t.duration
	})
}

// GetPauseDuration 获取计时器的暂停持续时间
func (t *Timer) GetPauseDuration() time.Duration {
	return syncx.WithRLockReturnValue(&t.mu, func() time.Duration {
		return t.pauseDuration
	})
}

// GetPaused 获取计时器的暂停状态
func (t *Timer) GetPaused() bool {
	return syncx.WithRLockReturnValue(&t.mu, func() bool {
		return t.paused
	})
}

// PrintLog 打印计时器的日志
func (t *Timer) PrintLog() {
	syncx.WithRLock(&t.mu, func() {
		fmt.Printf("Trace ID: %s, Duration Run Time: %v\n", t.traceId, t.GetDuration())
	})
}

// 确保 Timer 实现了 TimerInterface 接口
var _ TimerInterface = (*Timer)(nil)
