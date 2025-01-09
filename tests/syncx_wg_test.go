/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-08 13:06:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-09 15:55:23
 * @FilePath: \go-toolbox\tests\syncx_wg_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

// TestWaitGroupNormalGoroutines 测试正常的 goroutine
func TestWaitGroupNormalGoroutines(t *testing.T) {
	var wg syncx.WaitGroup

	wg.Go(func() {
		fmt.Println("Goroutine 1 is running")
	})

	wg.Go(func() {
		fmt.Println("Goroutine 2 is running")
	})

	// 等待所有 goroutine 完成
	err := wg.Wait()
	assert.NoError(t, err, "Expected no error from Wait()")
}

// TestWaitGroupPanicGoroutine 测试带有 panic 的 goroutine
func TestWaitGroupPanicGoroutine(t *testing.T) {
	wg := syncx.NewWaitGroup(true)

	wg.Go(func() {
		panic("This is a panic in goroutine 3")
	})

	wg.Go(func() {
		fmt.Println("Goroutine 4 is running")
	})

	// 等待所有 goroutine 完成
	err := wg.Wait()
	assert.Error(t, err, "Expected an error from Wait() due to panic")
	assert.EqualError(t, err, "发生了未知错误: This is a panic in goroutine 3", "Expected error message to match")
}

// TestWaitGroupMultiplePanicGoroutines 测试多个 panic 的处理
func TestWaitGroupMultiplePanicGoroutines(t *testing.T) {
	wg := syncx.NewWaitGroup(true)

	wg.Go(func() {
		panic("Panic in goroutine 1")
	})

	wg.Go(func() {
		panic("Panic in goroutine 2")
	})

	// 等待所有 goroutine 完成
	err := wg.Wait()
	assert.Error(t, err, "Expected an error from Wait() due to panic")
	assert.Contains(t, err.Error(), "发生了未知错误:", "Expected error message to contain '发生了未知错误:'")
}

// TestWaitGroupMaxConcurrency 测试最大并发限制
func TestWaitGroupMaxConcurrency(t *testing.T) {
	wg := syncx.NewWaitGroup(false, 2)

	// 启动超过最大并发数量的 goroutine
	for i := 0; i < 5; i++ {
		wg.Go(func() {
			fmt.Println("TestWaitGroupMaxConcurrency")
		})
	}

	// 等待所有 goroutine 完成
	wg.Wait()
}

// TestWaitGroupGetError 测试 GetError 方法
func TestWaitGroupGetError(t *testing.T) {
	wg := syncx.NewWaitGroup(false, 2)

	// Simulate an error
	wg.Go(func() {
		wg.SetError(fmt.Errorf("test error"))
	})

	// Wait for goroutines to finish
	wg.Wait()

	// Check if GetError returns the expected error
	err := wg.GetError()
	assert.Error(t, err, "Expected an error from GetError()")
	assert.EqualError(t, err, "test error", "Expected error message to match")
}

// TestWaitGroupGetChannelSize 测试 GetChannelSize 方法
func TestWaitGroupGetChannelSize(t *testing.T) {
	wg := syncx.NewWaitGroup(false, 3)

	// Check initial channel size
	assert.Equal(t, 0, wg.GetChannelSize(), "Expected initial channel size to be 0")

	// Start a goroutine to fill the channel
	wg.Go(func() {
		fmt.Println("TestWaitGroupGetChannelSize")
	})

	// Check channel size after starting a goroutine
	assert.Equal(t, 1, wg.GetChannelSize(), "Expected channel size to be 1 after starting one goroutine")

	// Release the channel position
	wg.Wait() // Wait for the goroutine to finish

	// Check final channel size
	assert.Equal(t, 0, wg.GetChannelSize(), "Expected channel size to be 0 after goroutine completes")
}
