/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-23 09:09:56
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-24 09:22:56
 * @FilePath: \go-toolbox\tests\syncx_defer_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

// TestWithDefer 测试 WithDefer 函数
func TestWithDefer(t *testing.T) {
	var executedOperation bool
	var executedCleanup bool

	// 定义操作和清理函数
	operation := func() {
		executedOperation = true
	}
	cleanup := func() {
		executedCleanup = true
	}

	// 调用 WithDefer 函数
	syncx.WithDefer(operation, cleanup)

	// 使用 assert 来验证操作和清理函数是否被调用
	assert.True(t, executedOperation, "Operation should be executed")
	assert.True(t, executedCleanup, "Cleanup should be executed")
}

// TestWithDeferReturnValue 测试 WithDeferReturnValue 函数
func TestWithDeferReturnValue(t *testing.T) {
	var executedCleanup bool
	operation := func() int {
		return 42 // 返回一个测试值
	}
	cleanup := func() {
		executedCleanup = true
	}

	result := syncx.WithDeferReturnValue(operation, cleanup)

	// 使用 assert 来验证结果和清理函数是否被调用
	assert.Equal(t, 42, result, "Expected result should be 42")
	assert.True(t, executedCleanup, "Cleanup should be executed")
}

// TestWithDeferReturn 测试 WithDeferReturn 函数
func TestWithDeferReturn(t *testing.T) {
	var executedCleanup bool
	operation := func() (int, error) {
		return 42, nil // 返回一个测试值和无错误
	}
	cleanup := func() {
		executedCleanup = true
	}

	result, err := syncx.WithDeferReturn(operation, cleanup)

	// 使用 assert 来验证结果和清理函数是否被调用
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 42, result, "Expected result should be 42")
	assert.True(t, executedCleanup, "Cleanup should be executed")
}

// TestWithDeferReturnWithError 测试 WithDeferReturn 函数的错误处理
func TestWithDeferReturnWithError(t *testing.T) {
	var executedCleanup bool
	operation := func() (int, error) {
		return 0, assert.AnError // 返回一个错误
	}
	cleanup := func() {
		executedCleanup = true
	}

	result, err := syncx.WithDeferReturn(operation, cleanup)

	// 使用 assert 来验证结果和清理函数是否被调用
	assert.Error(t, err, "Expected an error")
	assert.Equal(t, 0, result, "Expected result should be 0")
	assert.True(t, executedCleanup, "Cleanup should be executed")
}
