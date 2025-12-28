/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-13 11:27:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 23:12:10
 * @FilePath: \go-toolbox\pkg\errorx\base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package errorx

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		err      []error
		expected string
	}{
		{"with error", "an error occurred", []error{errors.New("original error")}, "an error occurred: original error"},
		{"without error", "another error", []error{}, "another error"},
		{"with nil error", "nil error", []error{nil}, "nil error"},
		{"with wrapped error", "wrapped error", []error{errors.New("something went wrong")}, "wrapped error: something went wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapError(tt.message, tt.err...)

			assert.NotNil(t, got)
			assert.EqualError(t, got, tt.expected)
		})
	}
}

func TestConcurrentErrorCreation(t *testing.T) {
	ResetErrorMap()
	var wg sync.WaitGroup
	const numGoroutines = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errType := ErrorType(i)
			RegisterError(errType, "resource not found")
			NewError(errType)
		}(i)
	}

	wg.Wait()
	count := len(GetErrorMap())
	assert.Equal(t, numGoroutines, count, "错误计数不正确")
}

func TestConcurrentErrorRegistration(t *testing.T) {
	ResetErrorMap()
	var wg sync.WaitGroup
	const numGoroutines = 50
	const errType = ErrorType(1)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			RegisterError(errType, "resource not found")
		}()
	}

	wg.Wait()
	assert.Equal(t, 1, len(GetErrorMap()), "错误映射不应包含重复的错误类型")
}

func TestNewErrorUnknownType(t *testing.T) {
	ResetErrorMap()
	unknownError := NewError(ErrorType(999))
	assert.EqualError(t, unknownError, "unknown error", "应返回未知错误消息")
}

func TestConcurrentErrorRetrieval(t *testing.T) {
	ResetErrorMap()
	const numGoroutines = 100
	var wg sync.WaitGroup

	RegisterError(1, "resource not found")

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := NewError(1)
			assert.EqualError(t, err, "resource not found", "应返回正确的错误消息")
		}()
	}

	wg.Wait()
}

func TestResetErrorMap(t *testing.T) {
	RegisterError(1, "resource not found")
	ResetErrorMap()
	assert.Empty(t, GetErrorMap(), "错误映射应为空")
}

func TestRegisterDifferentMessages(t *testing.T) {
	ResetErrorMap()
	RegisterError(1, "first error")
	RegisterError(1, "second error") // Should not register again

	assert.Equal(t, 1, len(GetErrorMap()), "错误映射应仅包含一个错误类型")
}

func TestErrorMessageFormatting(t *testing.T) {
	ResetErrorMap()
	RegisterError(1, "error occurred with code %d")
	err := NewError(1, 404)
	assert.EqualError(t, err, "error occurred with code 404", "错误消息格式化不正确")
}

func TestConcurrentResetErrorMap(t *testing.T) {
	RegisterError(1, "resource not found")
	var wg sync.WaitGroup
	const numGoroutines = 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ResetErrorMap()
		}()
	}

	wg.Wait()
	assert.Empty(t, GetErrorMap(), "错误映射应为空")
}

func TestPrintErrorMap(t *testing.T) {
	ResetErrorMap()
	RegisterError(1, "resource not found")
	RegisterError(2, "another error")
	PrintErrorMap() // 确保不会引发错误
}
