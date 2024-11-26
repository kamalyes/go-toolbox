/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-13 11:27:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-30 16:48:27
 * @FilePath: \go-toolbox\tests\errorx_base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
	"sync"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	tests := []struct {
		message  string
		err      error
		expected string
	}{
		{"an error occurred", errors.New("original error"), "an error occurred: original error"}, // 普通错误
		{"another error", nil, ""}, // nil 错误
		{"wrapped error", errors.New("something went wrong"), "wrapped error: something went wrong"}, // 包装错误
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := errorx.WrapError(tt.message, tt.err)

			if tt.expected == "" {
				assert.Nil(t, got) // 如果预期是 nil，断言返回值为 nil
			} else {
				assert.NotNil(t, got)                  // 断言返回值不为 nil
				assert.EqualError(t, got, tt.expected) // 断言返回的错误信息与预期相等
			}
		})
	}
}

func TestConcurrentErrorCreation(t *testing.T) {
	// 清除之前的错误计数
	errorx.ResetErrorMap()
	// 使用 WaitGroup 来等待多个 goroutine 完成
	var wg sync.WaitGroup
	const numGoroutines = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) { // 将 i 作为参数传递到 goroutine 中
			defer wg.Done()
			errType := errorx.ErrorType(i)
			// 注册错误类型
			errorx.RegisterError(errType, "resource not found")
			errorx.NewError(errType)
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 获取当前错误数量
	count := len(errorx.GetErrorMap())
	assert.Equal(t, numGoroutines, count, "错误计数不正确")
}
