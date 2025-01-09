/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 13:06:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-09 11:55:10
 * @FilePath: \go-toolbox\tests\syncx_wg_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

// 测试 WaitGroupWithMutex
func TestWaitGroupWithMutex(t *testing.T) {
	w := &syncx.WaitGroupWithMutex{}
	var count int64 // Use int64 for atomic operations

	const numGoroutines = 5

	w.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer w.Done()
			// Simulate work
			atomic.AddInt64(&count, int64(i)) // Atomic increment
			time.Sleep(100 * time.Millisecond)
		}(i)
	}

	w.Wait()
	assert.Equal(t, int64(10), count) // Compare with int64
}
