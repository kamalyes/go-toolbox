/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-28 00:00:00
 * @FilePath: \go-toolbox\pkg\convert\object_test.go
 * @Description: 对象解析和转换工具测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseObjectToMap 测试对象解析为 map
func TestParseObjectToMap(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email,omitempty"`
	}

	tests := []struct {
		name     string
		input    interface{}
		expected map[string]interface{}
	}{
		{"nil input", nil, nil},
		{"map[string]interface{}", map[string]interface{}{"key1": "value1", "key2": 123}, map[string]interface{}{"key1": "value1", "key2": 123}},
		{"struct", User{Name: "Alice", Age: 30, Email: "alice@example.com"}, map[string]interface{}{"name": "Alice", "age": 30, "email": "alice@example.com"}},
		{"struct pointer", &User{Name: "Bob", Age: 25, Email: "bob@example.com"}, map[string]interface{}{"name": "Bob", "age": 25, "email": "bob@example.com"}},
		{"non-struct type", "string", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseObjectToMap(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}
			assert.Equal(t, len(tt.expected), len(result))
			for k, v := range tt.expected {
				assert.Equal(t, v, result[k])
			}
		})
	}
}

// TestParseKVPairsToMap 测试键值对解析为 map
func TestParseKVPairsToMap(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name     string
		input    []interface{}
		expected map[string]interface{}
	}{
		{"empty input", []interface{}{}, nil},
		{"key-value pairs", []interface{}{"key1", "value1", "key2", 123}, map[string]interface{}{"key1": "value1", "key2": 123}},
		{"odd number of arguments", []interface{}{"key1", "value1", "key2"}, map[string]interface{}{"key1": "value1", "key2": ""}},
		{"single struct object", []interface{}{User{Name: "Alice", Age: 30}}, map[string]interface{}{"name": "Alice", "age": 30}},
		{"single map object", []interface{}{map[string]interface{}{"key1": "value1", "key2": 123}}, map[string]interface{}{"key1": "value1", "key2": 123}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseKVPairsToMap(tt.input...)
			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}
			assert.Equal(t, len(tt.expected), len(result))
			for k, v := range tt.expected {
				assert.Equal(t, v, result[k])
			}
		})
	}
}

// customStringer 自定义 Stringer 类型（用于测试）
type customStringer struct {
	value string
}

func (c customStringer) String() string {
	return c.value
}

// TestAppendValue 测试值追加到缓冲区
func TestAppendValue(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, "<nil>"},
		{"string", "hello", "hello"},
		{"bytes", []byte("world"), "world"},
		{"int", 123, "123"},
		{"int8", int8(12), "12"},
		{"int16", int16(1234), "1234"},
		{"int32", int32(789), "789"},
		{"int64", int64(456), "456"},
		{"uint", uint(111), "111"},
		{"uint8", uint8(22), "22"},
		{"uint16", uint16(333), "333"},
		{"uint32", uint32(333), "333"},
		{"uint64", uint64(222), "222"},
		{"uintptr", uintptr(0x1234), "0x1234"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"float32", float32(2.71), "2.71"},
		{"float64", 3.14, "3.14"},
		{"complex64", complex64(1 + 2i), "(1+2i)"},
		{"complex128", complex128(3 + 4i), "(3+4i)"},
		{"stringer", customStringer{"custom"}, "custom"},
		{"error", errors.New("test error"), "test error"},
		{"other", struct{ Name string }{"test"}, "{test}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 0, 64)
			result := AppendValue(buf, tt.input)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// BenchmarkParseObjectToMap 性能测试
func BenchmarkParseObjectToMap(b *testing.B) {
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	user := User{Name: "Alice", Age: 30, Email: "alice@example.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ParseObjectToMap(user)
	}
}

// BenchmarkParseKVPairsToMap 性能测试
func BenchmarkParseKVPairsToMap(b *testing.B) {
	kvs := []interface{}{"key1", "value1", "key2", 123, "key3", true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ParseKVPairsToMap(kvs...)
	}
}

// BenchmarkAppendValue 性能测试
func BenchmarkAppendValue(b *testing.B) {
	buf := make([]byte, 0, 1024)
	values := []any{"string", 123, 3.14, true, errors.New("error")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		for _, v := range values {
			buf = AppendValue(buf, v)
		}
	}
}
