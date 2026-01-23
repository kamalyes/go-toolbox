/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 10:50:01
 * @FilePath: \go-toolbox\pkg\validator\compare_test.go
 * @Description: 比较和验证功能测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
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

func TestValidateJSON(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"有效JSON对象", []byte(`{"name":"test","value":123}`), false},
		{"有效JSON数组", []byte(`[1,2,3]`), false},
		{"有效JSON字符串", []byte(`"hello"`), false},
		{"有效JSON数字", []byte(`123`), false},
		{"有效JSON布尔", []byte(`true`), false},
		{"有效JSON null", []byte(`null`), false},
		{"无效JSON", []byte(`{name:"test"}`), true},
		{"空字符串", []byte(``), true},
		{"不完整JSON", []byte(`{"name":"test"`), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJSON(tt.data)
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
			}
		})
	}
}

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

func TestCompareResult(t *testing.T) {
	a := assert.New(t)

	result := CompareResult{
		Success: true,
		Message: "测试消息",
		Actual:  "实际值",
		Expect:  "期望值",
	}

	a.True(result.Success, "CompareResult.Success should be true")
	a.Equal("测试消息", result.Message)
	a.Equal("实际值", result.Actual)
	a.Equal("期望值", result.Expect)
}

func TestValidateJSONWithData(t *testing.T) {
	a := assert.New(t)
	data := []byte(`{"name":"test","value":123}`)
	parsed, err := ValidateJSONWithData(data)
	a.NoError(err)
	m, ok := parsed.(map[string]any)
	a.True(ok)
	a.Equal("test", m["name"])
	a.Equal(float64(123), m["value"])

	invalid := []byte(`{name:"test"}`)
	_, err = ValidateJSONWithData(invalid)
	a.Error(err)
}

func TestValidateJSONField(t *testing.T) {
	a := assert.New(t)
	data := []byte(`{"name":"test","value":123}`)
	result := ValidateJSONField(data, "name", "test")
	a.True(result.Success)
	result = ValidateJSONField(data, "value", float64(123))
	a.True(result.Success)
	result = ValidateJSONField(data, "missing", "x")
	a.False(result.Success)
	result = ValidateJSONField(data, "name", "x")
	a.False(result.Success)
}

func TestValidateJSONFields(t *testing.T) {
	a := assert.New(t)
	data := []byte(`{"name":"test","value":123}`)
	rules := map[string]any{"name": "test", "value": float64(123)}
	results := ValidateJSONFields(data, rules)
	a.Len(results, 2)
	for _, r := range results {
		a.True(r.Success)
	}
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

func TestValidateRegex(t *testing.T) {
	a := assert.New(t)
	body := []byte("abc123")
	result := ValidateRegex(body, "^[a-z]+[0-9]+$")
	a.True(result.Success)
	result = ValidateRegex(body, "^[0-9]+[a-z]+$")
	a.False(result.Success)
	result = ValidateRegex(body, "")
	a.True(result.Success)
	result = ValidateRegex(body, "[")
	a.False(result.Success)
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

func BenchmarkCompareNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareNumbers(100, 200, OpLessThan)
	}
}

func BenchmarkValidateJSON(b *testing.B) {
	data := []byte(`{"name":"test","value":123,"nested":{"key":"value"}}`)
	for i := 0; i < b.N; i++ {
		ValidateJSON(data)
	}
}

func BenchmarkValidateStatusCodeRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateStatusCodeRange(200, 200, 299)
	}
}
