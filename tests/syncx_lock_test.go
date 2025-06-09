/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 13:06:30
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-09 11:36:21
 * @FilePath: \go-toolbox\tests\syncx_lock_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"sync"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

// MockLocker 是一个用于测试的简单实现 Locker 接口的结构体
type MockLocker struct {
	mu sync.Mutex
}

func (l *MockLocker) Lock() {
	l.mu.Lock()
}

func (l *MockLocker) Unlock() {
	l.mu.Unlock()
}

// TestWithLock 测试 WithLock 函数
func TestWithLock(t *testing.T) {
	counter := 0
	lock := &MockLocker{}

	// 使用 WithLock 执行操作
	syncx.WithLock(lock, func() {
		counter++
	})

	assert.Equal(t, 1, counter, "Expected counter to be 1")
}

// TestWithLockReturnValue 测试 WithLockReturnValue 函数
func TestWithLockReturnValue(t *testing.T) {
	lock := &MockLocker{}

	// 使用 WithLockReturnValue 执行操作
	result := syncx.WithLockReturnValue(lock, func() int {
		return 42 // 返回一个结果
	})

	assert.Equal(t, 42, result, "Expected result to be 42")
}

// TestWithLockConcurrent 测试 WithLock 在并发情况下的表现
func TestWithLockConcurrent(t *testing.T) {
	counter := 0
	lock := &MockLocker{}
	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			syncx.WithLock(lock, func() {
				counter++
			})
		}()
	}

	wg.Wait()

	assert.Equal(t, goroutines, counter, "Expected counter to be %d", goroutines)
}

// MockRWLocker 是一个用于测试的简单实现 Locker 接口的结构体
type MockRWLocker struct {
	mu sync.RWMutex
}

func (l *MockRWLocker) RLock() {
	l.mu.RLock()
}

func (l *MockRWLocker) RUnlock() {
	l.mu.RUnlock()
}

// TestWithLockReturn 测试 WithLockReturn 函数
func TestWithLockReturn(t *testing.T) {
	lock := &MockLocker{}

	// 使用 WithLockReturn 执行操作
	result, err := syncx.WithLockReturn(lock, func() (int, error) {
		return 42, nil // 返回一个结果
	})

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 42, result, "Expected result to be 42")
}

// TestWithLockReturnString 测试 WithLockReturn 函数返回字符串
func TestWithLockReturnString(t *testing.T) {
	lock := &MockLocker{}

	// 使用 WithLockReturn 执行操作
	result, err := syncx.WithLockReturn(lock, func() (string, error) {
		return "Hello, World!", nil // 返回一个字符串
	})

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "Hello, World!", result, "Expected result to be 'Hello, World!'")
}

// TestWithLockReturnBool 测试 WithLockReturn 函数返回布尔值
func TestWithLockReturnBool(t *testing.T) {
	lock := &MockLocker{}

	// 使用 WithLockReturn 执行操作
	result, err := syncx.WithLockReturn(lock, func() (bool, error) {
		return true, nil // 返回一个布尔值
	})

	assert.NoError(t, err, "Expected no error")
	assert.True(t, result, "Expected result to be true")
}

// TestWithLockReturnComplex 测试 WithLockReturn 函数返回复杂类型
func TestWithLockReturnComplex(t *testing.T) {
	lock := &MockLocker{}

	type CustomStruct struct {
		Name  string
		Value int
	}

	// 使用 WithLockReturn 执行操作
	result, err := syncx.WithLockReturn(lock, func() (CustomStruct, error) {
		return CustomStruct{Name: "Test", Value: 100}, nil // 返回一个自定义结构体
	})

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, CustomStruct{Name: "Test", Value: 100}, result, "Expected result to match the custom struct")
}

// TestWithLockReturnConcurrent 测试 WithLockReturn 在并发情况下的表现
func TestWithLockReturnConcurrent(t *testing.T) {
	lock := &MockLocker{}
	const goroutines = 100

	var wg sync.WaitGroup
	results := make([]int, goroutines)

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			result, err := syncx.WithLockReturn(lock, func() (int, error) {
				return index, nil // 返回当前索引
			})
			assert.NoError(t, err, "Expected no error")
			results[index] = result
		}(i)
	}

	wg.Wait()

	// 确保每个索引都有对应的结果
	for i := 0; i < goroutines; i++ {
		assert.Equal(t, i, results[i], "Expected result at index %d to be %d", i, i)
	}
}

// TestWithRLock 测试 WithRLock 函数
func TestWithRLock(t *testing.T) {
	counter := 0
	lock := &MockRWLocker{}

	// 使用 WithRLock 执行操作
	syncx.WithRLock(lock, func() {
		counter++
	})

	assert.Equal(t, 1, counter, "Expected counter to be 1")
}

// TestWithRLockReturn 测试 WithRLockReturn 函数
func TestWithRLockReturn(t *testing.T) {
	lock := &MockRWLocker{}

	// 使用 WithRLockReturn 执行操作
	result, err := syncx.WithRLockReturn(lock, func() (int, error) {
		return 42, nil // 返回一个结果
	})

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 42, result, "Expected result to be 42")
}

// TestWithRLockReturnConcurrent 测试 WithRLockReturn 在并发情况下的表现
func TestWithRLockReturnConcurrent(t *testing.T) {
	lock := &MockRWLocker{}
	const goroutines = 100

	var wg sync.WaitGroup
	results := make([]int, goroutines)

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			result, err := syncx.WithRLockReturn(lock, func() (int, error) {
				return index, nil // 返回当前索引
			})
			assert.NoError(t, err, "Expected no error")
			results[index] = result
		}(i)
	}

	wg.Wait()

	// 确保每个索引都有对应的结果
	for i := 0; i < goroutines; i++ {
		assert.Equal(t, i, results[i], "Expected result at index %d to be %d", i, i)
	}
}

// TestWithRLockReturnValue 测试 WithRLockReturnValue 函数
func TestWithRLockReturnValue(t *testing.T) {
	lock := &MockRWLocker{}

	// 使用 WithRLockReturnValue 执行操作
	result := syncx.WithRLockReturnValue(lock, func() int {
		return 42 // 返回一个结果
	})

	assert.Equal(t, 42, result, "Expected result to be 42")
}

// --------------------------------------------------
// 模拟一个支持 TryLock 的简单互斥锁实现
type TryMutex struct {
	mu sync.Mutex
}

func (m *TryMutex) Lock() {
	m.mu.Lock()
}

func (m *TryMutex) Unlock() {
	m.mu.Unlock()
}

func (m *TryMutex) TryLock() bool {
	// 这里示例用 false，实际不调用
	return false
}

// 用 channel 模拟的简单 TryMutex 实现
type ChanTryMutex struct {
	ch chan struct{}
}

func NewChanTryMutex() *ChanTryMutex {
	m := &ChanTryMutex{ch: make(chan struct{}, 1)}
	m.ch <- struct{}{} // 初始化时可用
	return m
}

func (m *ChanTryMutex) Lock() {
	<-m.ch
}

func (m *ChanTryMutex) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *ChanTryMutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
		return false
	}
}

// --------------------------------------------------
// 模拟一个支持 TryRLock 的读写锁实现
type TryRWMutex struct {
	mu         sync.RWMutex
	tryRLockCh chan struct{}
}

func NewTryRWMutex() *TryRWMutex {
	m := &TryRWMutex{
		tryRLockCh: make(chan struct{}, 1),
	}
	m.tryRLockCh <- struct{}{}
	return m
}

func (m *TryRWMutex) RLock() {
	m.mu.RLock()
}

func (m *TryRWMutex) RUnlock() {
	m.mu.RUnlock()
}

func (m *TryRWMutex) TryRLock() bool {
	select {
	case <-m.tryRLockCh:
		m.mu.RLock()
		return true
	default:
		return false
	}
}

func (m *TryRWMutex) RUnlockTry() {
	m.mu.RUnlock()
	m.releaseTryRLockSignal()
}

// 新增：安全恢复信号，非阻塞放入
func (m *TryRWMutex) releaseTryRLockSignal() {
	select {
	case m.tryRLockCh <- struct{}{}:
	default:
	}
}

// --------------------------------------------------
// 这里模拟 syncx 包里的错误变量
var ErrLockNotAcquired = assert.AnError

// 这里模拟 syncx 包里的辅助函数（简化版）
func WithTryLock(mu interface {
	TryLock() bool
	Unlock()
}, fn func()) error {
	if !mu.TryLock() {
		return ErrLockNotAcquired
	}
	defer mu.Unlock()
	fn()
	return nil
}

func WithTryLockReturn[T any](mu interface {
	TryLock() bool
	Unlock()
}, fn func() (T, error)) (T, error) {
	var zero T
	if !mu.TryLock() {
		return zero, ErrLockNotAcquired
	}
	defer mu.Unlock()
	return fn()
}

func WithTryLockReturnValue[T any](mu interface {
	TryLock() bool
	Unlock()
}, fn func() T) (T, error) {
	var zero T
	if !mu.TryLock() {
		return zero, ErrLockNotAcquired
	}
	defer mu.Unlock()
	return fn(), nil
}

func WithTryRLock(mu interface {
	TryRLock() bool
	RUnlockTry()
}, fn func()) error {
	if !mu.TryRLock() {
		return ErrLockNotAcquired
	}
	defer mu.RUnlockTry()
	fn()
	return nil
}

func WithTryRLockReturn[T any](mu interface {
	TryRLock() bool
	RUnlockTry()
}, fn func() (T, error)) (T, error) {
	var zero T
	if !mu.TryRLock() {
		return zero, ErrLockNotAcquired
	}
	defer mu.RUnlockTry()
	return fn()
}

func WithTryRLockReturnValue[T any](mu interface {
	TryRLock() bool
	RUnlockTry()
}, fn func() T) (T, error) {
	var zero T
	if !mu.TryRLock() {
		return zero, ErrLockNotAcquired
	}
	defer mu.RUnlockTry()
	return fn(), nil
}

// --------------------------------------------------

func TestWithTryLock(t *testing.T) {
	mu := NewChanTryMutex()

	err := WithTryLock(mu, func() {
		// 执行任务
	})
	assert.NoError(t, err)

	mu.Lock()
	err = WithTryLock(mu, func() {
		t.Fatal("不应该执行到这里")
	})
	assert.ErrorIs(t, err, ErrLockNotAcquired)
	mu.Unlock()
}

func TestWithTryLockReturn(t *testing.T) {
	mu := NewChanTryMutex()

	v, err := WithTryLockReturn(mu, func() (int, error) {
		return 42, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 42, v)

	mu.Lock()
	v, err = WithTryLockReturn(mu, func() (int, error) {
		t.Fatal("不应该执行操作")
		return 0, nil
	})
	assert.ErrorIs(t, err, ErrLockNotAcquired)
	assert.Equal(t, 0, v)
	mu.Unlock()
}

func TestWithTryLockReturnValue(t *testing.T) {
	mu := NewChanTryMutex()

	v, err := WithTryLockReturnValue(mu, func() string {
		return "hello"
	})
	assert.NoError(t, err)
	assert.Equal(t, "hello", v)

	mu.Lock()
	v, err = WithTryLockReturnValue(mu, func() string {
		t.Fatal("不应该执行操作")
		return ""
	})
	assert.ErrorIs(t, err, ErrLockNotAcquired)
	assert.Equal(t, "", v)
	mu.Unlock()
}

func TestWithTryRLock(t *testing.T) {
	rwmu := NewTryRWMutex()

	err := WithTryRLock(rwmu, func() {
		// 执行读操作
	})
	assert.NoError(t, err)

	// 安全消耗信号，模拟 TryRLock 失败
	select {
	case <-rwmu.tryRLockCh:
	default:
		t.Fatal("信号通道已空，无法消耗信号")
	}

	err = WithTryRLock(rwmu, func() {
		t.Fatal("不应该执行操作")
	})
	assert.ErrorIs(t, err, ErrLockNotAcquired)

	// 恢复信号，避免影响后续测试
	rwmu.releaseTryRLockSignal()
}

func TestWithTryRLockReturn(t *testing.T) {
	rwmu := NewTryRWMutex()

	v, err := WithTryRLockReturn(rwmu, func() (string, error) {
		return "read", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "read", v)

	// 安全消耗信号，模拟失败
	select {
	case <-rwmu.tryRLockCh:
	default:
		t.Fatal("信号通道已空，无法消耗信号")
	}

	v, err = WithTryRLockReturn(rwmu, func() (string, error) {
		t.Fatal("不应该执行操作")
		return "", nil
	})
	assert.ErrorIs(t, err, ErrLockNotAcquired)
	assert.Equal(t, "", v)

	// 恢复信号
	rwmu.releaseTryRLockSignal()
}

func TestWithTryRLockReturnValue(t *testing.T) {
	rwmu := NewTryRWMutex()

	v, err := WithTryRLockReturnValue(rwmu, func() string {
		return "readValue"
	})
	assert.NoError(t, err)
	assert.Equal(t, "readValue", v)

	// 安全消耗信号，模拟失败
	select {
	case <-rwmu.tryRLockCh:
	default:
		t.Fatal("信号通道已空，无法消耗信号")
	}

	v, err = WithTryRLockReturnValue(rwmu, func() string {
		t.Fatal("不应该执行操作")
		return ""
	})
	assert.ErrorIs(t, err, ErrLockNotAcquired)
	assert.Equal(t, "", v)

	// 恢复信号
	rwmu.releaseTryRLockSignal()
}
