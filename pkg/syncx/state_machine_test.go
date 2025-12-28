/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 00:00:00 09:09:21
 * @FilePath: \go-toolbox\pkg\syncx\state_machine_test.go
 * @Description: 状态机测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package syncx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义测试用的状态类型
type ConnectionState string

const (
	StateDisconnected ConnectionState = "disconnected"
	StateConnecting   ConnectionState = "connecting"
	StateConnected    ConnectionState = "connected"
	StateReconnecting ConnectionState = "reconnecting"
)

func TestNewStateMachine(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	assert.NotNil(t, sm)
	assert.Equal(t, StateDisconnected, sm.CurrentState())
}

func TestStateMachine_AllowTransition(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)

	sm.AllowTransition(StateDisconnected, StateConnecting)
	sm.AllowTransition(StateConnecting, StateConnected)
	sm.AllowTransition(StateConnecting, StateDisconnected)

	assert.True(t, sm.CanTransitionTo(StateConnecting))
	assert.False(t, sm.CanTransitionTo(StateConnected))
}

func TestStateMachine_AllowTransitions(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)

	sm.AllowTransitions(StateDisconnected, StateConnecting, StateConnected)

	assert.True(t, sm.CanTransitionTo(StateConnecting))
	assert.True(t, sm.CanTransitionTo(StateConnected))
	assert.False(t, sm.CanTransitionTo(StateReconnecting))
}

func TestStateMachine_DisallowTransition(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)

	sm.AllowTransitions(StateDisconnected, StateConnecting, StateConnected)
	assert.True(t, sm.CanTransitionTo(StateConnecting))

	sm.DisallowTransition(StateDisconnected, StateConnecting)
	assert.False(t, sm.CanTransitionTo(StateConnecting))
	assert.True(t, sm.CanTransitionTo(StateConnected))
}

func TestStateMachine_TransitionTo(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnecting)
	sm.AllowTransition(StateConnecting, StateConnected)

	t.Run("valid transition", func(t *testing.T) {
		err := sm.TransitionTo(StateConnecting)
		assert.NoError(t, err)
		assert.Equal(t, StateConnecting, sm.CurrentState())
	})

	t.Run("invalid transition", func(t *testing.T) {
		err := sm.TransitionTo(StateDisconnected)
		assert.Error(t, err)
		assert.Equal(t, StateConnecting, sm.CurrentState()) // 状态不变
	})

	t.Run("valid second transition", func(t *testing.T) {
		err := sm.TransitionTo(StateConnected)
		assert.NoError(t, err)
		assert.Equal(t, StateConnected, sm.CurrentState())
	})

	t.Run("same state transition", func(t *testing.T) {
		err := sm.TransitionTo(StateConnected)
		assert.NoError(t, err) // 转换到相同状态不报错
		assert.Equal(t, StateConnected, sm.CurrentState())
	})
}

func TestStateMachine_MustTransitionTo(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnecting)

	t.Run("valid transition", func(t *testing.T) {
		assert.NotPanics(t, func() {
			sm.MustTransitionTo(StateConnecting)
		})
		assert.Equal(t, StateConnecting, sm.CurrentState())
	})

	t.Run("invalid transition panics", func(t *testing.T) {
		assert.Panics(t, func() {
			sm.MustTransitionTo(StateDisconnected)
		})
	})
}

func TestStateMachine_OnTransition(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnecting)
	sm.AllowTransition(StateConnecting, StateConnected)

	var transitions []struct {
		from, to ConnectionState
	}

	sm.OnTransition(func(from, to ConnectionState) {
		transitions = append(transitions, struct{ from, to ConnectionState }{from, to})
	})

	_ = sm.TransitionTo(StateConnecting)
	_ = sm.TransitionTo(StateConnected)

	assert.Equal(t, 2, len(transitions))
	assert.Equal(t, StateDisconnected, transitions[0].from)
	assert.Equal(t, StateConnecting, transitions[0].to)
	assert.Equal(t, StateConnecting, transitions[1].from)
	assert.Equal(t, StateConnected, transitions[1].to)
}

func TestStateMachine_OnEnter(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnected)

	var enteredFrom ConnectionState
	sm.OnEnter(StateConnected, func(from ConnectionState) {
		enteredFrom = from
	})

	_ = sm.TransitionTo(StateConnected)

	assert.Equal(t, StateDisconnected, enteredFrom)
}

func TestStateMachine_OnExit(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnected)

	var exitedTo ConnectionState
	sm.OnExit(StateDisconnected, func(to ConnectionState) {
		exitedTo = to
	})

	_ = sm.TransitionTo(StateConnected)

	assert.Equal(t, StateConnected, exitedTo)
}

func TestStateMachine_GetAllowedTransitions(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransitions(StateDisconnected, StateConnecting, StateConnected)

	allowed := sm.GetAllowedTransitions()
	assert.Equal(t, 2, len(allowed))

	// 转换到新状态后,允许的转换应该改变
	sm.AllowTransition(StateConnecting, StateConnected)
	_ = sm.TransitionTo(StateConnecting)
	allowed = sm.GetAllowedTransitions()
	assert.Equal(t, 1, len(allowed))
	assert.Equal(t, StateConnected, allowed[0])
}

func TestStateMachine_Reset(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnected)

	_ = sm.TransitionTo(StateConnected)
	assert.Equal(t, StateConnected, sm.CurrentState())

	sm.Reset(StateDisconnected)
	assert.Equal(t, StateDisconnected, sm.CurrentState())
}

func TestStateMachine_ClearCallbacks(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnected)

	callbackCalled := false
	sm.OnTransition(func(from, to ConnectionState) {
		callbackCalled = true
	})

	_ = sm.TransitionTo(StateConnected)
	assert.True(t, callbackCalled)

	callbackCalled = false
	sm.Reset(StateDisconnected)
	sm.ClearCallbacks()

	_ = sm.TransitionTo(StateConnected)
	assert.False(t, callbackCalled) // 回调已被清除
}

func TestStateMachine_WithAllowAnyTransition(t *testing.T) {
	sm := NewStateMachine(StateDisconnected, WithAllowAnyTransition[ConnectionState]())

	// 不需要配置转换规则,任意转换都允许
	assert.True(t, sm.CanTransitionTo(StateConnecting))
	assert.True(t, sm.CanTransitionTo(StateConnected))
	assert.True(t, sm.CanTransitionTo(StateReconnecting))

	err := sm.TransitionTo(StateConnected)
	assert.NoError(t, err)
	assert.Equal(t, StateConnected, sm.CurrentState())
}

func TestStateMachine_ConcurrentAccess(t *testing.T) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnecting)
	sm.AllowTransition(StateConnecting, StateConnected)
	sm.AllowTransition(StateConnected, StateDisconnected)

	var wg sync.WaitGroup
	numGoroutines := 100

	// 并发读取当前状态
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sm.CurrentState()
		}()
	}

	// 并发检查转换
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sm.CanTransitionTo(StateConnecting)
		}()
	}

	wg.Wait()
}

// 基准测试
func BenchmarkStateMachine_CurrentState(b *testing.B) {
	sm := NewStateMachine(StateDisconnected)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sm.CurrentState()
	}
}

func BenchmarkStateMachine_CanTransitionTo(b *testing.B) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnecting)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sm.CanTransitionTo(StateConnecting)
	}
}

func BenchmarkStateMachine_TransitionTo(b *testing.B) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnecting)
	sm.AllowTransition(StateConnecting, StateDisconnected)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			_ = sm.TransitionTo(StateConnecting)
		} else {
			_ = sm.TransitionTo(StateDisconnected)
		}
	}
}

func BenchmarkStateMachine_TransitionWithCallback(b *testing.B) {
	sm := NewStateMachine(StateDisconnected)
	sm.AllowTransition(StateDisconnected, StateConnecting)
	sm.AllowTransition(StateConnecting, StateDisconnected)

	sm.OnTransition(func(from, to ConnectionState) {
		// 简单的回调
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			_ = sm.TransitionTo(StateConnecting)
		} else {
			_ = sm.TransitionTo(StateDisconnected)
		}
	}
}

// 示例测试
func ExampleStateMachine() {
	// 创建连接状态机
	sm := NewStateMachine(StateDisconnected)

	// 配置允许的状态转换
	sm.AllowTransition(StateDisconnected, StateConnecting)
	sm.AllowTransition(StateConnecting, StateConnected)
	sm.AllowTransition(StateConnecting, StateDisconnected)
	sm.AllowTransition(StateConnected, StateDisconnected)

	// 注册状态转换回调
	sm.OnTransition(func(from, to ConnectionState) {
		// log.Printf("State: %s -> %s", from, to)
	})

	// 执行状态转换
	_ = sm.TransitionTo(StateConnecting)
	_ = sm.TransitionTo(StateConnected)

	// Output:
}
