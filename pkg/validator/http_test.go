/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\http_test.go
 * @Description: HTTP 验证测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStatusCode(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		actual   int
		expect   int
		op       CompareOperator
		wantPass bool
	}{
		{"200 OK 相等", 200, 200, OpEqual, true},
		{"404 Not Found 相等", 404, 404, OpEqual, true},
		{"不匹配", 200, 404, OpEqual, false},
		{"500 错误", 500, 500, OpEqual, true},
		{"大于", 404, 200, OpGreaterThan, true},
		{"小于", 200, 404, OpLessThan, true},
		{"大于等于-相等", 200, 200, OpGreaterThanOrEqual, true},
		{"大于等于-大于", 404, 200, OpGreaterThanOrEqual, true},
		{"小于等于-相等", 200, 200, OpLessThanOrEqual, true},
		{"小于等于-小于", 200, 404, OpLessThanOrEqual, true},
		{"不等于-通过", 200, 404, OpNotEqual, true},
		{"不等于-失败", 200, 200, OpNotEqual, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateStatusCode(tt.actual, tt.expect, tt.op)
			a.Equal(tt.wantPass, result.Success, "message: %s", result.Message)
		})
	}
}

func TestValidateStatusCodeRange(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		actual   int
		min      int
		max      int
		wantPass bool
	}{
		{"2xx成功", 200, 200, 299, true},
		{"2xx边界-下限", 200, 200, 299, true},
		{"2xx边界-上限", 299, 200, 299, true},
		{"4xx失败", 404, 200, 299, false},
		{"5xx失败", 500, 200, 299, false},
		{"低于范围", 199, 200, 299, false},
		{"高于范围", 300, 200, 299, false},
		{"单个状态码", 200, 200, 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateStatusCodeRange(tt.actual, tt.min, tt.max)
			a.Equal(tt.wantPass, result.Success, "message: %s", result.Message)
			if !result.Success {
				a.NotEmpty(result.Message, "ValidateStatusCodeRange() should have error message when failed")
			}
		})
	}
}

func TestValidateHeader(t *testing.T) {
	a := assert.New(t)
	headers := map[string]string{"Content-Type": "application/json", "X-Test": "abc"}
	result := ValidateHeader(headers, "Content-Type", "json", OpContains)
	a.True(result.Success)
	result = ValidateHeader(headers, "X-Test", "abc", OpEqual)
	a.True(result.Success)
	result = ValidateHeader(headers, "X-Test", "xyz", OpEqual)
	a.False(result.Success)
	result = ValidateHeader(headers, "Missing", "x", OpEqual)
	a.False(result.Success)
}

func TestValidateContentType(t *testing.T) {
	a := assert.New(t)
	headers := map[string]string{"Content-Type": "application/json; charset=utf-8"}
	result := ValidateContentType(headers, "json")
	a.True(result.Success)
	result = ValidateContentType(headers, "xml")
	a.False(result.Success)
}

// Benchmark tests
func BenchmarkValidateStatusCodeRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateStatusCodeRange(200, 200, 299)
	}
}
