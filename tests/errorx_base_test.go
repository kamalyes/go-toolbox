/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-13 11:27:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-08 15:55:55
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
		{"an error occurred", errors.New("original error"), "an error occurred: original error"},
		{"another error", nil, ""},
		{"wrapped error", errors.New("something went wrong"), "wrapped error: something went wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := errorx.WrapError(tt.message, tt.err)

			if tt.expected == "" {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.EqualError(t, got, tt.expected)
			}
		})
	}
}

func TestConcurrentErrorCreation(t *testing.T) {
	errorx.ResetErrorMap()
	var wg sync.WaitGroup
	const numGoroutines = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errType := errorx.ErrorType(i)
			errorx.RegisterError(errType, "resource not found")
			errorx.NewError(errType)
		}(i)
	}

	wg.Wait()
	count := len(errorx.GetErrorMap())
	assert.Equal(t, numGoroutines, count, "错误计数不正确")
}

func TestConcurrentErrorRegistration(t *testing.T) {
	errorx.ResetErrorMap()
	var wg sync.WaitGroup
	const numGoroutines = 50
	const errType = errorx.ErrorType(1)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errorx.RegisterError(errType, "resource not found")
		}()
	}

	wg.Wait()
	assert.Equal(t, 1, len(errorx.GetErrorMap()), "错误映射不应包含重复的错误类型")
}

func TestNewErrorUnknownType(t *testing.T) {
	errorx.ResetErrorMap()
	unknownError := errorx.NewError(errorx.ErrorType(999))
	assert.EqualError(t, unknownError, "unknown error", "应返回未知错误消息")
}

func TestConcurrentErrorRetrieval(t *testing.T) {
	errorx.ResetErrorMap()
	const numGoroutines = 100
	var wg sync.WaitGroup

	errorx.RegisterError(1, "resource not found")

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := errorx.NewError(1)
			assert.EqualError(t, err, "resource not found", "应返回正确的错误消息")
		}()
	}

	wg.Wait()
}

func TestResetErrorMap(t *testing.T) {
	errorx.RegisterError(1, "resource not found")
	errorx.ResetErrorMap()
	assert.Empty(t, errorx.GetErrorMap(), "错误映射应为空")
}

func TestRegisterDifferentMessages(t *testing.T) {
	errorx.ResetErrorMap()
	errorx.RegisterError(1, "first error")
	errorx.RegisterError(1, "second error") // Should not register again

	assert.Equal(t, 1, len(errorx.GetErrorMap()), "错误映射应仅包含一个错误类型")
}

func TestErrorMessageFormatting(t *testing.T) {
	errorx.ResetErrorMap()
	errorx.RegisterError(1, "error occurred with code %d")
	err := errorx.NewError(1, 404)
	assert.EqualError(t, err, "error occurred with code 404", "错误消息格式化不正确")
}

func TestConcurrentResetErrorMap(t *testing.T) {
	errorx.RegisterError(1, "resource not found")
	var wg sync.WaitGroup
	const numGoroutines = 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errorx.ResetErrorMap()
		}()
	}

	wg.Wait()
	assert.Empty(t, errorx.GetErrorMap(), "错误映射应为空")
}

func TestPrintErrorMap(t *testing.T) {
	errorx.ResetErrorMap()
	errorx.RegisterError(1, "resource not found")
	errorx.RegisterError(2, "another error")
	errorx.PrintErrorMap() // 确保不会引发错误
}
