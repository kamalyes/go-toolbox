/*
 * @Author: kamalyes 501893067@qq.com
 * @Date:2024-10-24 10:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-24 11:00:16
 * @FilePath: \go-toolbox\json\base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package json

import (
	"testing"
)

// 辅助函数，创建测试用例
func newTestCase(originalJSON string, pairs *KeyValuePairs, expected map[string]interface{}) struct {
	originalJSON string
	pairs        *KeyValuePairs
	expected     map[string]interface{}
} {
	return struct {
		originalJSON string
		pairs        *KeyValuePairs
		expected     map[string]interface{}
	}{
		originalJSON: originalJSON,
		pairs:        pairs,
		expected:     expected,
	}
}

// 辅助函数，检查 JSON 对象是否与期望值匹配
func checkJSONEquality(t testing.TB, result, expected map[string]interface{}) {
	for key, expectedValue := range expected {
		if result[key] != expectedValue {
			t.Errorf("对于键 %s, 期望 %v, 实际 %v", key, expectedValue, result[key])
		}
	}
}

// 辅助函数，解析 JSON 字符串并返回 map
func parseJSON(t testing.TB, jsonStr string) map[string]interface{} {
	var result map[string]interface{}
	if err := Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("解析 JSON 失败: %v", err)
	}
	return result
}

// 测试将键值对追加到 JSON 字符串的功能
func TestAppendKeysToJSON(t *testing.T) {
	tests := []struct {
		originalJSON string
		pairs        *KeyValuePairs
		expected     map[string]interface{}
	}{
		newTestCase(
			`{"name": "Alice", "age": 30}`,
			NewKeyValuePairs().Add("city", "New York").Add("country", "USA"),
			map[string]interface{}{
				"name":    "Alice",
				"age":     float64(30),
				"city":    "New York",
				"country": "USA",
			},
		),
		newTestCase(
			`{}`,
			NewKeyValuePairs().Add("key1", "value1"),
			map[string]interface{}{
				"key1": "value1",
			},
		),
		newTestCase(
			`{"existing": "value"}`,
			NewKeyValuePairs().Add("newKey", "newValue"),
			map[string]interface{}{
				"existing": "value",
				"newKey":   "newValue",
			},
		),
		newTestCase(
			`{"a": 1}`,
			NewKeyValuePairs().Add("b", 2).Add("c", 3),
			map[string]interface{}{
				"a": float64(1),
				"b": float64(2),
				"c": float64(3),
			},
		),
	}

	// 遍历测试用例
	for _, tt := range tests {
		updatedJSON, err := AppendKeysToJSON(tt.originalJSON, tt.pairs)
		if err != nil {
			t.Fatalf("期望没有错误，实际错误为 %v", err)
		}

		result := parseJSON(t, updatedJSON)
		checkJSONEquality(t, result, tt.expected)
	}
}

// BenchmarkAppendKeysToJSON 基准测试 AppendKeysToJSON 函数
func BenchmarkAppendKeysToJSON(b *testing.B) {
	originalJSON := `{"name": "Alice", "age": 30}`
	pairs := NewKeyValuePairs().
		Add("city", "New York").
		Add("country", "USA").
		Add("occupation", "Engineer").
		Add("hobby", "Photography")

	expected := map[string]interface{}{
		"name":       "Alice",
		"age":        float64(30),
		"city":       "New York",
		"country":    "USA",
		"occupation": "Engineer",
		"hobby":      "Photography",
	}

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		updatedJSON, err := AppendKeysToJSON(originalJSON, pairs)
		if err != nil {
			b.Fatalf("期望没有错误，实际错误为 %v", err)
		}

		result := parseJSON(b, updatedJSON)
		checkJSONEquality(b, result, expected)
	}
}
