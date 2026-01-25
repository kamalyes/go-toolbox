/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\json_test.go
 * @Description: JSON 验证测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestValidateJSONPath(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		body     []byte
		path     string
		expected any
		op       CompareOperator
		wantPass bool
	}{
		{
			name:     "简单路径-字符串",
			body:     []byte(`{"name":"test"}`),
			path:     "$.name",
			expected: "test",
			op:       OpEqual,
			wantPass: true,
		},
		{
			name:     "简单路径-数字",
			body:     []byte(`{"value":123}`),
			path:     "$.value",
			expected: 123,
			op:       OpEqual,
			wantPass: true,
		},
		{
			name:     "嵌套路径",
			body:     []byte(`{"user":{"name":"test"}}`),
			path:     "$.user.name",
			expected: "test",
			op:       OpEqual,
			wantPass: true,
		},
		{
			name:     "数组索引",
			body:     []byte(`{"items":["a","b","c"]}`),
			path:     "$.items[0]",
			expected: "a",
			op:       OpEqual,
			wantPass: true,
		},
		{
			name:     "路径不存在",
			body:     []byte(`{"name":"test"}`),
			path:     "$.missing",
			expected: "x",
			op:       OpEqual,
			wantPass: false,
		},
		{
			name:     "值不匹配",
			body:     []byte(`{"name":"test"}`),
			path:     "$.name",
			expected: "wrong",
			op:       OpEqual,
			wantPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONPath(tt.body, tt.path, tt.expected, tt.op)
			a.Equal(tt.wantPass, result.Success, "message: %s", result.Message)
		})
	}
}

func TestValidateJSONPathExists(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		body     []byte
		path     string
		wantPass bool
	}{
		{
			name:     "路径存在",
			body:     []byte(`{"name":"test"}`),
			path:     "$.name",
			wantPass: true,
		},
		{
			name:     "嵌套路径存在",
			body:     []byte(`{"user":{"name":"test"}}`),
			path:     "$.user.name",
			wantPass: true,
		},
		{
			name:     "路径不存在",
			body:     []byte(`{"name":"test"}`),
			path:     "$.missing",
			wantPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONPathExists(tt.body, tt.path)
			a.Equal(tt.wantPass, result.Success, "message: %s", result.Message)
		})
	}
}

// Benchmark tests
func BenchmarkValidateJSON(b *testing.B) {
	data := []byte(`{"name":"test","value":123,"nested":{"key":"value"}}`)
	for i := 0; i < b.N; i++ {
		ValidateJSON(data)
	}
}

func BenchmarkValidateJSONPath(b *testing.B) {
	data := []byte(`{"user":{"name":"test","age":30}}`)
	for i := 0; i < b.N; i++ {
		ValidateJSONPath(data, "$.user.name", "test", OpEqual)
	}
}
