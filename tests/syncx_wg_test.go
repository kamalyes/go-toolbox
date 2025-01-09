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

// TestNormalGoroutines 测试正常的 goroutine
func TestNormalGoroutines(t *testing.T) {
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

// TestPanicGoroutine 测试带有 panic 的 goroutine
func TestPanicGoroutine(t *testing.T) {
	var wg syncx.WaitGroup

	wg.GoTry(func() {
		panic("This is a panic in goroutine 3")
	})

	wg.GoTry(func() {
		fmt.Println("Goroutine 4 is running")
	})

	// 等待所有 goroutine 完成
	err := wg.Wait()
	assert.Error(t, err, "Expected an error from Wait() due to panic")
	assert.EqualError(t, err, "发生了未知错误: This is a panic in goroutine 3", "Expected error message to match")
}

// TestMultiplePanicGoroutines 测试多个 panic 的处理
func TestMultiplePanicGoroutines(t *testing.T) {
	var wg syncx.WaitGroup

	wg.GoTry(func() {
		panic("Panic in goroutine 1")
	})

	wg.GoTry(func() {
		panic("Panic in goroutine 2")
	})

	// 等待所有 goroutine 完成
	err := wg.Wait()
	assert.Error(t, err, "Expected an error from Wait() due to panic")
	assert.Contains(t, err.Error(), "发生了未知错误:", "Expected error message to contain '发生了未知错误:'")
}

// TestMaxConcurrency 测试最大并发限制
func TestMaxConcurrency(t *testing.T) {
	maxConcurrency := uint(2)
	wg := syncx.NewWaitGroup(maxConcurrency)

	// 启动超过最大并发数量的 goroutine
	for i := 0; i < 5; i++ {
		wg.Go(func() {
			fmt.Println("TestMaxConcurrency")
		})
	}

	// 等待所有 goroutine 完成
	wg.Wait()
}
