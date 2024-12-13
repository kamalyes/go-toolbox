/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 13:06:30
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-13 13:08:20
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
