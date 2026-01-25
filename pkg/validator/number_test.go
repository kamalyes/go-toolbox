/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\number_test.go
 * @Description: 数值验证测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareNumbers(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		actual   int
		expect   int
		op       CompareOperator
		wantPass bool
	}{
		{"相等-通过", 100, 100, OpEqual, true},
		{"相等-失败", 100, 200, OpEqual, false},
		{"不相等-通过", 100, 200, OpNotEqual, true},
		{"不相等-失败", 100, 100, OpNotEqual, false},
		{"大于-通过", 200, 100, OpGreaterThan, true},
		{"大于-失败", 100, 200, OpGreaterThan, false},
		{"大于等于-通过-大于", 200, 100, OpGreaterThanOrEqual, true},
		{"大于等于-通过-等于", 100, 100, OpGreaterThanOrEqual, true},
		{"大于等于-失败", 100, 200, OpGreaterThanOrEqual, false},
		{"小于-通过", 100, 200, OpLessThan, true},
		{"小于-失败", 200, 100, OpLessThan, false},
		{"小于等于-通过-小于", 100, 200, OpLessThanOrEqual, true},
		{"小于等于-通过-等于", 100, 100, OpLessThanOrEqual, true},
		{"小于等于-失败", 200, 100, OpLessThanOrEqual, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareNumbers(tt.actual, tt.expect, tt.op)
			a.Equal(tt.wantPass, result.Success, "message: %s", result.Message)
		})
	}
}

func TestCompareNumbersSymbolOperators(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		actual   int
		expect   int
		op       CompareOperator
		wantPass bool
	}{
		{"符号相等-通过", 100, 100, OpSymbolEqual, true},
		{"符号相等-失败", 100, 200, OpSymbolEqual, false},
		{"符号不相等-通过", 100, 200, OpSymbolNotEqual, true},
		{"符号不相等-失败", 100, 100, OpSymbolNotEqual, false},
		{"符号大于-通过", 200, 100, OpSymbolGreaterThan, true},
		{"符号大于-失败", 100, 200, OpSymbolGreaterThan, false},
		{"符号大于等于-通过-大于", 200, 100, OpSymbolGreaterThanOrEqual, true},
		{"符号大于等于-通过-等于", 100, 100, OpSymbolGreaterThanOrEqual, true},
		{"符号大于等于-失败", 100, 200, OpSymbolGreaterThanOrEqual, false},
		{"符号小于-通过", 100, 200, OpSymbolLessThan, true},
		{"符号小于-失败", 200, 100, OpSymbolLessThan, false},
		{"符号小于等于-通过-小于", 100, 200, OpSymbolLessThanOrEqual, true},
		{"符号小于等于-通过-等于", 100, 100, OpSymbolLessThanOrEqual, true},
		{"符号小于等于-失败", 200, 100, OpSymbolLessThanOrEqual, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareNumbers(tt.actual, tt.expect, tt.op)
			a.Equal(tt.wantPass, result.Success, "message: %s", result.Message)
		})
	}
}

func TestCompareNumbersFloat(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		actual   float64
		expect   float64
		op       CompareOperator
		wantPass bool
	}{
		{"浮点数相等", 3.14, 3.14, OpEqual, true},
		{"浮点数大于", 3.14, 2.71, OpGreaterThan, true},
		{"浮点数小于", 2.71, 3.14, OpLessThan, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareNumbers(tt.actual, tt.expect, tt.op)
			a.Equal(tt.wantPass, result.Success)
		})
	}
}

func TestCompareNumbersInvalidOperator(t *testing.T) {
	a := assert.New(t)
	result := CompareNumbers(100, 200, OpContains)
	a.False(result.Success, "CompareNumbers() should fail with string operator")
	a.NotEmpty(result.Message, "CompareNumbers() should have error message for invalid operator")
}

// Benchmark tests
func BenchmarkCompareNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareNumbers(100, 200, OpLessThan)
	}
}
