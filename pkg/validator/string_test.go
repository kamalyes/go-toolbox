/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\string_test.go
 * @Description: 字符串验证测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareStrings(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		actual   string
		expect   string
		op       CompareOperator
		wantPass bool
	}{
		{"相等-通过", "hello", "hello", OpEqual, true},
		{"相等-失败", "hello", "world", OpEqual, false},
		{"不相等-通过", "hello", "world", OpNotEqual, true},
		{"不相等-失败", "hello", "hello", OpNotEqual, false},
		{"包含-通过", "hello world", "world", OpContains, true},
		{"包含-失败", "hello world", "test", OpContains, false},
		{"不包含-通过", "hello world", "test", OpNotContains, true},
		{"不包含-失败", "hello world", "hello", OpNotContains, false},
		{"前缀-通过", "hello world", "hello", OpHasPrefix, true},
		{"前缀-失败", "hello world", "world", OpHasPrefix, false},
		{"后缀-通过", "hello world", "world", OpHasSuffix, true},
		{"后缀-失败", "hello world", "hello", OpHasSuffix, false},
		{"为空-通过", "", "", OpEmpty, true},
		{"为空-失败", "hello", "", OpEmpty, false},
		{"非空-通过", "hello", "", OpNotEmpty, true},
		{"非空-失败", "", "", OpNotEmpty, false},
		{"正则-通过", "abc123", "^[a-z]+[0-9]+$", OpRegex, true},
		{"正则-失败", "123abc", "^[a-z]+[0-9]+$", OpRegex, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareStrings(tt.actual, tt.expect, tt.op)
			a.Equal(tt.wantPass, result.Success, "message: %s", result.Message)
			a.Equal(tt.actual, result.Actual)
		})
	}
}

func TestCompareStringsSymbolOperators(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		actual   string
		expect   string
		op       CompareOperator
		wantPass bool
	}{
		{"符号相等-通过", "hello", "hello", OpSymbolEqual, true},
		{"符号相等-失败", "hello", "world", OpSymbolEqual, false},
		{"符号不相等-通过", "hello", "world", OpSymbolNotEqual, true},
		{"符号不相等-失败", "hello", "hello", OpSymbolNotEqual, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareStrings(tt.actual, tt.expect, tt.op)
			a.Equal(tt.wantPass, result.Success, "message: %s", result.Message)
		})
	}
}

func TestCompareStringsInvalidOperator(t *testing.T) {
	a := assert.New(t)
	result := CompareStrings("hello", "world", "invalid")
	a.False(result.Success, "CompareStrings() should fail with invalid operator")
	a.NotEmpty(result.Message, "CompareStrings() should have error message for invalid operator")
}

func TestCompareStringsInvalidRegex(t *testing.T) {
	a := assert.New(t)
	result := CompareStrings("hello", "[invalid(", OpRegex)
	a.False(result.Success, "CompareStrings() should fail with invalid regex")
	a.NotEmpty(result.Message, "CompareStrings() should have error message for invalid regex")
}

func TestValidateContains(t *testing.T) {
	a := assert.New(t)
	body := []byte("hello world")
	result := ValidateContains(body, "world")
	a.True(result.Success)
	result = ValidateContains(body, "test")
	a.False(result.Success)
	result = ValidateContains(body, "")
	a.True(result.Success)
}

func TestValidateNotContains(t *testing.T) {
	a := assert.New(t)
	body := []byte("hello world")
	result := ValidateNotContains(body, "test")
	a.True(result.Success)
	result = ValidateNotContains(body, "hello")
	a.False(result.Success)
	result = ValidateNotContains(body, "")
	a.True(result.Success)
}

// Benchmark tests
func BenchmarkCompareStringsEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareStrings("hello world", "hello world", OpEqual)
	}
}

func BenchmarkCompareStringsContains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareStrings("hello world", "world", OpContains)
	}
}

func BenchmarkCompareStringsRegex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareStrings("abc123", "^[a-z]+[0-9]+$", OpRegex)
	}
}
