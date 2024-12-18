/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 13:06:30
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-18 18:31:09
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

func (l *MockRWLocker) Lock() {
	l.mu.Lock()
}

func (l *MockRWLocker) Unlock() {
	l.mu.Unlock()
}

func (l *MockRWLocker) RLock() {
	l.mu.RLock()
}

func (l *MockRWLocker) RUnlock() {
	l.mu.RUnlock()
}

// TestWithRLock 测试 WithLock 函数在 RWMutex 下的表现
func TestWithRLock(t *testing.T) {
	counter := 0
	lock := &MockRWLocker{}
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
