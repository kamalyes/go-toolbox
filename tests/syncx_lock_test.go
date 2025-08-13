/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 13:06:30
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-06 17:37:31
 * @FilePath: \go-toolbox\tests\syncx_lock_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"sync"
	"sync/atomic"
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

func (l *MockRWLocker) Lock() {
	l.mu.Lock()
}

func (l *MockRWLocker) Unlock() {
	l.mu.Unlock()
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

// TestWithUnlockThenLock 测试 WithUnlockThenLock 函数
func TestWithUnlockThenLock(t *testing.T) {
	mockLocker := &MockLocker{}
	var counter int

	// 在锁定的情况下增加计数器
	mockLocker.Lock()
	counter = 0

	syncx.WithUnlockThenLock(mockLocker, func() {
		counter++
	})

	assert.Equal(t, 1, counter, "计数器应该增加到 1")
}

// TestWithRUnlockThenLock 测试 WithRUnlockThenLock 函数
func TestWithRUnlockThenLock(t *testing.T) {
	mockRLocker := &MockRWLocker{}
	var counter int

	// 在锁定的情况下增加计数器
	mockRLocker.Lock()
	counter = 0

	syncx.WithRUnlockThenLock(mockRLocker, func() {
		counter++
	})
	assert.Equal(t, 1, counter, "计数器应该增加到 1")
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

// SimpleTryMutex 是个支持 TryLock 的简单互斥锁实现
type SimpleTryMutex struct {
	state int32
}

func (m *SimpleTryMutex) Lock() {
	for !m.TryLock() {
		// 自旋等待
	}
}

func (m *SimpleTryMutex) Unlock() {
	if atomic.LoadInt32(&m.state) == 0 {
		panic("unlock of unlocked SimpleTryMutex")
	}
	atomic.StoreInt32(&m.state, 0)
}

func (m *SimpleTryMutex) TryLock() bool {
	return atomic.CompareAndSwapInt32(&m.state, 0, 1)
}

// SimpleTryRWMutex 是个支持 TryRLock 的简单读写锁实现
type SimpleTryRWMutex struct {
	state int32
}

const (
	writeBit = 1 << 16
	readMask = writeBit - 1
)

func (m *SimpleTryRWMutex) Lock() {
	for {
		if atomic.CompareAndSwapInt32(&m.state, 0, writeBit) {
			return
		}
	}
}

func (m *SimpleTryRWMutex) Unlock() {
	if atomic.LoadInt32(&m.state) != writeBit {
		panic("unlock of unlocked SimpleTryRWMutex")
	}
	atomic.StoreInt32(&m.state, 0)
}

func (m *SimpleTryRWMutex) TryLock() bool {
	return atomic.CompareAndSwapInt32(&m.state, 0, writeBit)
}

func (m *SimpleTryRWMutex) RLock() {
	for {
		s := atomic.LoadInt32(&m.state)
		if s&writeBit != 0 {
			continue
		}
		if s&readMask == readMask {
			panic("reader count overflow")
		}
		if atomic.CompareAndSwapInt32(&m.state, s, s+1) {
			return
		}
	}
}

func (m *SimpleTryRWMutex) RUnlock() {
	for {
		s := atomic.LoadInt32(&m.state)
		if s&readMask == 0 {
			panic("unlock of unlocked read lock")
		}
		if atomic.CompareAndSwapInt32(&m.state, s, s-1) {
			return
		}
	}
}

func (m *SimpleTryRWMutex) TryRLock() bool {
	for {
		s := atomic.LoadInt32(&m.state)
		if s&writeBit != 0 {
			return false
		}
		if s&readMask == readMask {
			panic("reader count overflow")
		}
		if atomic.CompareAndSwapInt32(&m.state, s, s+1) {
			return true
		}
	}
}

func TestWithTryLock(t *testing.T) {
	lock := &SimpleTryMutex{}

	// 成功获取锁并执行操作
	err := syncx.WithTryLock(lock, func() error {
		return nil
	})
	assert.NoError(t, err)

	// 手动先锁住，测试不能获取锁
	ok := lock.TryLock()
	assert.True(t, ok)

	err = syncx.WithTryLock(lock, func() error {
		t.Fatal("不应该执行")
		return nil
	})
	assert.Equal(t, syncx.ErrLockNotAcquired, err)
	lock.Unlock()
}

func TestWithTryLockReturn(t *testing.T) {
	lock := &SimpleTryMutex{}

	val, err := syncx.WithTryLockReturn(lock, func() (int, error) {
		return 42, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 42, val)

	ok := lock.TryLock()
	assert.True(t, ok)

	val, err = syncx.WithTryLockReturn(lock, func() (int, error) {
		t.Fatal("不应该执行")
		return 0, nil
	})
	assert.Equal(t, syncx.ErrLockNotAcquired, err)
	assert.Equal(t, 0, val)
	lock.Unlock()
}

func TestWithTryLockReturnValue(t *testing.T) {
	lock := &SimpleTryMutex{}

	val, err := syncx.WithTryLockReturnValue(lock, func() string {
		return "hello"
	})
	assert.NoError(t, err)
	assert.Equal(t, "hello", val)

	ok := lock.TryLock()
	assert.True(t, ok)

	val, err = syncx.WithTryLockReturnValue(lock, func() string {
		t.Fatal("不应该执行")
		return "fail"
	})
	assert.Equal(t, syncx.ErrLockNotAcquired, err)
	assert.Equal(t, "", val)
	lock.Unlock()
}

func TestWithTryRLock(t *testing.T) {
	lock := &SimpleTryRWMutex{}

	err := syncx.WithTryRLock(lock, func() error {
		return nil
	})
	assert.NoError(t, err)

	ok := lock.TryLock()
	assert.True(t, ok)

	err = syncx.WithTryRLock(lock, func() error {
		t.Fatal("不应该执行")
		return nil
	})
	assert.Equal(t, syncx.ErrLockNotAcquired, err)
	lock.Unlock()
}

func TestWithTryRLockReturn(t *testing.T) {
	lock := &SimpleTryRWMutex{}

	val, err := syncx.WithTryRLockReturn(lock, func() (string, error) {
		return "read", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "read", val)

	ok := lock.TryLock()
	assert.True(t, ok)

	val, err = syncx.WithTryRLockReturn(lock, func() (string, error) {
		t.Fatal("不应该执行")
		return "", nil
	})
	assert.Equal(t, syncx.ErrLockNotAcquired, err)
	assert.Equal(t, "", val)
	lock.Unlock()
}

func TestWithTryRLockReturnValue(t *testing.T) {
	lock := &SimpleTryRWMutex{}

	val, err := syncx.WithTryRLockReturnValue(lock, func() int {
		return 123
	})
	assert.NoError(t, err)
	assert.Equal(t, 123, val)

	ok := lock.TryLock()
	assert.True(t, ok)

	val, err = syncx.WithTryRLockReturnValue(lock, func() int {
		t.Fatal("不应该执行")
		return 456
	})
	assert.Equal(t, syncx.ErrLockNotAcquired, err)
	assert.Equal(t, 0, val)
	lock.Unlock()
}
