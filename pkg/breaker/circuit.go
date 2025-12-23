/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-23 23:50:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-24 00:35:10
 * @FilePath: \go-toolbox\pkg\breaker\circuit.go
 * @Description: 熔断器实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package breaker

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
)

// State 熔断器状态
type State int32

const (
	StateClosed   State = iota // 关闭(正常)
	StateOpen                  // 开启(熔断)
	StateHalfOpen              // 半开(尝试恢复)
)

// String 返回状态字符串
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Circuit 熔断器
type Circuit struct {
	name              string
	maxFailures       int32         // 最大失败次数
	resetTimeout      time.Duration // 重置超时时间
	halfOpenSuccesses int32         // 半开状态需要的成功次数
	state             int32         // 当前状态
	failures          int32         // 失败计数
	successes         int32         // 成功计数
	lastFailureTime   int64         // 最后失败时间
	mu                sync.RWMutex
	onStateChange     func(from, to State)
}

// Config 熔断器配置
type Config struct {
	MaxFailures       int32
	ResetTimeout      time.Duration
	HalfOpenSuccesses int32
	OnStateChange     func(from, to State)
}

// New 创建熔断器
func New(name string, config Config) *Circuit {
	config.MaxFailures = mathx.IF(config.MaxFailures == 0, 5, config.MaxFailures)
	config.ResetTimeout = mathx.IF(config.ResetTimeout == 0, 30*time.Second, config.ResetTimeout)
	config.HalfOpenSuccesses = mathx.IF(config.HalfOpenSuccesses == 0, 2, config.HalfOpenSuccesses)

	return &Circuit{
		name:              name,
		maxFailures:       config.MaxFailures,
		resetTimeout:      config.ResetTimeout,
		halfOpenSuccesses: config.HalfOpenSuccesses,
		state:             int32(StateClosed),
		onStateChange:     config.OnStateChange,
	}
}

// Execute 执行带熔断保护的操作
func (c *Circuit) Execute(fn func() error) error {
	if !c.AllowRequest() {
		return ErrOpen
	}

	err := fn()
	if err != nil {
		c.RecordFailure()
		return err
	}

	c.RecordSuccess()
	return nil
}

// AllowRequest 是否允许请求
func (c *Circuit) AllowRequest() bool {
	state := State(atomic.LoadInt32(&c.state))

	switch state {
	case StateClosed:
		return true
	case StateOpen:
		// 检查是否应该进入半开状态
		lastFailure := atomic.LoadInt64(&c.lastFailureTime)
		if time.Since(time.Unix(0, lastFailure)) > c.resetTimeout {
			c.setState(StateHalfOpen)
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// RecordSuccess 记录成功
func (c *Circuit) RecordSuccess() {
	state := State(atomic.LoadInt32(&c.state))

	switch state {
	case StateClosed:
		atomic.StoreInt32(&c.failures, 0)
	case StateHalfOpen:
		successes := atomic.AddInt32(&c.successes, 1)
		if successes >= c.halfOpenSuccesses {
			c.setState(StateClosed)
			atomic.StoreInt32(&c.successes, 0)
			atomic.StoreInt32(&c.failures, 0)
		}
	}
}

// RecordFailure 记录失败
func (c *Circuit) RecordFailure() {
	atomic.StoreInt64(&c.lastFailureTime, time.Now().UnixNano())
	state := State(atomic.LoadInt32(&c.state))

	switch state {
	case StateClosed:
		failures := atomic.AddInt32(&c.failures, 1)
		if failures >= c.maxFailures {
			c.setState(StateOpen)
		}
	case StateHalfOpen:
		c.setState(StateOpen)
		atomic.StoreInt32(&c.successes, 0)
	}
}

// setState 设置状态
func (c *Circuit) setState(newState State) {
	c.mu.Lock()
	defer c.mu.Unlock()

	oldState := State(atomic.LoadInt32(&c.state))
	if oldState == newState {
		return
	}

	atomic.StoreInt32(&c.state, int32(newState))

	if c.onStateChange != nil {
		c.onStateChange(oldState, newState)
	}
}

// GetState 获取当前状态
func (c *Circuit) GetState() State {
	return State(atomic.LoadInt32(&c.state))
}

// GetStats 获取统计信息
func (c *Circuit) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"name":     c.name,
		"state":    c.GetState().String(),
		"failures": atomic.LoadInt32(&c.failures),
	}
}

// Stats 统计信息
func (c *Circuit) Stats() CircuitStats {
	return CircuitStats{
		Name:     c.name,
		State:    c.GetState().String(),
		Failures: atomic.LoadInt32(&c.failures),
	}
}

// CircuitStats 熔断器统计
type CircuitStats struct {
	Name     string
	State    string
	Failures int32
}
