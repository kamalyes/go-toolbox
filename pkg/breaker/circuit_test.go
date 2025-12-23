/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-23 23:50:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-23 23:50:00
 * @FilePath: \go-toolbox\pkg\breaker\circuit_test.go
 * @Description: 熔断器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package breaker

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCircuit(t *testing.T) {
	config := Config{
		MaxFailures:       3,
		ResetTimeout:      time.Second,
		HalfOpenSuccesses: 2,
	}

	cb := New("test-circuit", config)

	assert.NotNil(t, cb)
	assert.Equal(t, "test-circuit", cb.name)
	assert.Equal(t, int32(3), cb.maxFailures)
	assert.Equal(t, time.Second, cb.resetTimeout)
	assert.Equal(t, int32(2), cb.halfOpenSuccesses)
	assert.Equal(t, StateClosed, cb.GetState())
}

func TestNewCircuitWithDefaults(t *testing.T) {
	cb := New("test", Config{})

	assert.Equal(t, int32(5), cb.maxFailures)
	assert.Equal(t, 30*time.Second, cb.resetTimeout)
	assert.Equal(t, int32(2), cb.halfOpenSuccesses)
}

func TestCircuitStateClosed(t *testing.T) {
	cb := New("test", Config{MaxFailures: 3})

	assert.True(t, cb.AllowRequest())
	assert.Equal(t, StateClosed, cb.GetState())
}

func TestCircuitStateTransitionToOpen(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:  3,
		ResetTimeout: time.Second,
	})

	// 记录失败次数达到阈值
	cb.RecordFailure()
	assert.Equal(t, StateClosed, cb.GetState())

	cb.RecordFailure()
	assert.Equal(t, StateClosed, cb.GetState())

	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())
}

func TestCircuitStateOpen(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:  2,
		ResetTimeout: 100 * time.Millisecond,
	})

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()

	assert.Equal(t, StateOpen, cb.GetState())
	assert.False(t, cb.AllowRequest())
}

func TestCircuitStateTransitionToHalfOpen(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:  2,
		ResetTimeout: 50 * time.Millisecond,
	})

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())

	// 等待重置超时
	time.Sleep(60 * time.Millisecond)

	// 应该允许请求并进入半开状态
	assert.True(t, cb.AllowRequest())
	assert.Equal(t, StateHalfOpen, cb.GetState())
}

func TestCircuitStateHalfOpenToClosed(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:       2,
		ResetTimeout:      50 * time.Millisecond,
		HalfOpenSuccesses: 2,
	})

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())

	// 进入半开状态
	time.Sleep(60 * time.Millisecond)
	cb.AllowRequest()
	assert.Equal(t, StateHalfOpen, cb.GetState())

	// 记录成功
	cb.RecordSuccess()
	assert.Equal(t, StateHalfOpen, cb.GetState())

	cb.RecordSuccess()
	assert.Equal(t, StateClosed, cb.GetState())
}

func TestCircuitStateHalfOpenToOpen(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:  2,
		ResetTimeout: 50 * time.Millisecond,
	})

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()

	// 进入半开状态
	time.Sleep(60 * time.Millisecond)
	cb.AllowRequest()
	assert.Equal(t, StateHalfOpen, cb.GetState())

	// 半开状态下失败
	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())
}

func TestCircuitExecuteSuccess(t *testing.T) {
	cb := New("test", Config{MaxFailures: 3})

	err := cb.Execute(func() error {
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, StateClosed, cb.GetState())
}

func TestCircuitExecuteFailure(t *testing.T) {
	cb := New("test", Config{MaxFailures: 2})

	testErr := errors.New("test error")
	err := cb.Execute(func() error {
		return testErr
	})

	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

func TestCircuitExecuteWhenOpen(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:  1,
		ResetTimeout: time.Second,
	})

	// 触发熔断
	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())

	err := cb.Execute(func() error {
		return nil
	})

	assert.Error(t, err)
	assert.Equal(t, ErrOpen, err)
}

func TestCircuitRecordSuccessInClosedState(t *testing.T) {
	cb := New("test", Config{MaxFailures: 3})

	// 先记录一些失败
	cb.RecordFailure()
	cb.RecordFailure()
	assert.Equal(t, int32(2), cb.failures)

	// 记录成功应该重置失败计数
	cb.RecordSuccess()
	assert.Equal(t, int32(0), cb.failures)
	assert.Equal(t, StateClosed, cb.GetState())
}

func TestCircuitRecordSuccessInHalfOpenState(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:       2,
		ResetTimeout:      50 * time.Millisecond,
		HalfOpenSuccesses: 3,
	})

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()

	// 进入半开状态
	time.Sleep(60 * time.Millisecond)
	cb.AllowRequest()
	assert.Equal(t, StateHalfOpen, cb.GetState())

	// 记录成功
	cb.RecordSuccess()
	assert.Equal(t, int32(1), cb.successes)
	assert.Equal(t, StateHalfOpen, cb.GetState())

	cb.RecordSuccess()
	assert.Equal(t, int32(2), cb.successes)
	assert.Equal(t, StateHalfOpen, cb.GetState())

	cb.RecordSuccess()
	assert.Equal(t, StateClosed, cb.GetState())
	assert.Equal(t, int32(0), cb.successes)
	assert.Equal(t, int32(0), cb.failures)
}

func TestCircuitRecordFailureInClosedState(t *testing.T) {
	cb := New("test", Config{MaxFailures: 3})

	cb.RecordFailure()
	assert.Equal(t, int32(1), cb.failures)
	assert.Equal(t, StateClosed, cb.GetState())

	cb.RecordFailure()
	assert.Equal(t, int32(2), cb.failures)
	assert.Equal(t, StateClosed, cb.GetState())

	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())
}

func TestCircuitRecordFailureInHalfOpenState(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:       2,
		ResetTimeout:      50 * time.Millisecond,
		HalfOpenSuccesses: 2,
	})

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()

	// 进入半开状态
	time.Sleep(60 * time.Millisecond)
	cb.AllowRequest()
	assert.Equal(t, StateHalfOpen, cb.GetState())

	// 记录成功
	cb.RecordSuccess()
	assert.Equal(t, int32(1), cb.successes)

	// 半开状态下失败应该立即回到开启状态
	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())
	assert.Equal(t, int32(0), cb.successes)
}

func TestCircuitOnStateChange(t *testing.T) {
	var stateChanges []struct {
		from State
		to   State
	}

	cb := New("test", Config{
		MaxFailures:  2,
		ResetTimeout: 50 * time.Millisecond,
		OnStateChange: func(from, to State) {
			stateChanges = append(stateChanges, struct {
				from State
				to   State
			}{from, to})
		},
	})

	// 触发熔断
	cb.RecordFailure()
	cb.RecordFailure()

	assert.Len(t, stateChanges, 1)
	assert.Equal(t, StateClosed, stateChanges[0].from)
	assert.Equal(t, StateOpen, stateChanges[0].to)

	// 进入半开状态
	time.Sleep(60 * time.Millisecond)
	cb.AllowRequest()

	assert.Len(t, stateChanges, 2)
	assert.Equal(t, StateOpen, stateChanges[1].from)
	assert.Equal(t, StateHalfOpen, stateChanges[1].to)

	// 恢复到关闭状态
	cb.RecordSuccess()
	cb.RecordSuccess()

	assert.Len(t, stateChanges, 3)
	assert.Equal(t, StateHalfOpen, stateChanges[2].from)
	assert.Equal(t, StateClosed, stateChanges[2].to)
}

func TestCircuitGetStats(t *testing.T) {
	cb := New("test-circuit", Config{MaxFailures: 3})

	cb.RecordFailure()
	cb.RecordFailure()

	stats := cb.GetStats()

	assert.Equal(t, "test-circuit", stats["name"])
	assert.Equal(t, "closed", stats["state"])
	assert.Equal(t, int32(2), stats["failures"])
}

func TestCircuitStats(t *testing.T) {
	cb := New("test-circuit", Config{MaxFailures: 3})

	cb.RecordFailure()

	stats := cb.Stats()

	assert.Equal(t, "test-circuit", stats.Name)
	assert.Equal(t, "closed", stats.State)
	assert.Equal(t, int32(1), stats.Failures)
}

func TestStateString(t *testing.T) {
	assert.Equal(t, "closed", StateClosed.String())
	assert.Equal(t, "open", StateOpen.String())
	assert.Equal(t, "half-open", StateHalfOpen.String())
	assert.Equal(t, "unknown", State(999).String())
}

func TestCircuitConcurrentAccess(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:  10,
		ResetTimeout: time.Second,
	})

	done := make(chan bool)
	workers := 10

	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				cb.Execute(func() error {
					if j%2 == 0 {
						return nil
					}
					return errors.New("error")
				})
			}
			done <- true
		}()
	}

	for i := 0; i < workers; i++ {
		<-done
	}

	// 验证状态一致性
	state := cb.GetState()
	assert.True(t, state == StateClosed || state == StateOpen || state == StateHalfOpen)
}

func TestCircuitAllowRequestConcurrency(t *testing.T) {
	cb := New("test", Config{
		MaxFailures:  5,
		ResetTimeout: 100 * time.Millisecond,
	})

	done := make(chan bool)
	workers := 20

	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < 50; j++ {
				cb.AllowRequest()
				time.Sleep(time.Millisecond)
			}
			done <- true
		}()
	}

	for i := 0; i < workers; i++ {
		<-done
	}

	assert.NotNil(t, cb.GetState())
}
