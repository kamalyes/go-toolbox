/*
 * @Author: kamalyes 501893067@qq.com
 * @Date:2024-10-24 10:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 23:15:50
 * @FilePath: \go-toolbox\pkg\json\base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package json

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 辅助函数，创建测试用例
func newTestCase(originalJSON string, pairs *KeyValuePairs, expected map[string]interface{}, expectError bool) struct {
	originalJSON string
	pairs        *KeyValuePairs
	expected     map[string]interface{}
	expectError  bool
} {
	return struct {
		originalJSON string
		pairs        *KeyValuePairs
		expected     map[string]interface{}
		expectError  bool
	}{
		originalJSON: originalJSON,
		pairs:        pairs,
		expected:     expected,
		expectError:  expectError,
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
func parseJSON(t testing.TB, jsonStr []byte) map[string]interface{} {
	var result map[string]interface{}
	if err := Unmarshal(jsonStr, &result); err != nil {
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
		expectError  bool // 新增字段，指示是否期望错误
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
			false,
		),
		newTestCase(
			`{}`,
			NewKeyValuePairs().Add("key1", "value1"),
			map[string]interface{}{
				"key1": "value1",
			},
			false,
		),
		newTestCase(
			`{"existing": "value"}`,
			NewKeyValuePairs().Add("newKey", "newValue"),
			map[string]interface{}{
				"existing": "value",
				"newKey":   "newValue",
			},
			false,
		),
		newTestCase(
			`{"a": 1}`,
			NewKeyValuePairs().Add("b", 2).Add("c", 3),
			map[string]interface{}{
				"a": float64(1),
				"b": float64(2),
				"c": float64(3),
			},
			false,
		),
		// 无效的 JSON 字符串
		newTestCase(
			`{"name": "Alice", "age": 30,}`,
			NewKeyValuePairs().Add("city", "New York"),
			nil,  // 预期返回 nil
			true, // 期望错误
		),
		// 传入 nil 的 JSON 字符串
		newTestCase(
			`nil`,
			NewKeyValuePairs().Add("key", "value"),
			nil,  // 预期返回 nil
			true, // 期望错误
		),
	}

	// 遍历测试用例
	for _, tt := range tests {
		updatedJSON, err := AppendKeysToJSONMarshal(tt.originalJSON, tt.pairs)
		if tt.expectError {
			// 如果期望错误，检查是否返回了错误
			if err == nil {
				t.Fatalf("期望错误，但没有返回错误")
			}
			continue // 继续下一个测试用例
		}

		// 如果不期望错误，检查返回的 JSON
		if err != nil {
			t.Fatalf("期望没有错误，实际错误为 %v", err)
		}

		result := parseJSON(t, updatedJSON)
		checkJSONEquality(t, result, tt.expected)
	}
}

func TestReplaceKeysComplex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `{"user_info": {"first_name": "Alice", "last_name": "Smith", "contact_details": {"email_address": "alice@example.com", "phone_number": "123-456-7890"}}}`,
			expected: `{"user-info":{"first-name":"Alice","last-name":"Smith","contact-details":{"email-address":"alice@example.com","phone-number":"123-456-7890"}}}`,
		},
		{
			input:    `{"products": [{"product_id": 1, "product_name": "Widget", "product_details": {"weight_kg": 2.5, "dimensions_cm": {"length": 10, "width": 5, "height": 2}}}, {"product_id": 2, "product_name": "Gadget", "product_details": {"weight_kg": 1.2, "dimensions_cm": {"length": 8, "width": 4, "height": 1}}}]} `,
			expected: `{"products":[{"product-id":1,"product-name":"Widget","product-details":{"weight-kg":2.5,"dimensions-cm":{"length":10,"width":5,"height":2}}},{"product-id":2,"product-name":"Gadget","product-details":{"weight-kg":1.2,"dimensions-cm":{"length":8,"width":4,"height":1}}}]} `,
		},
		{
			input:    `{"metadata": {"created_at": "2023-01-01", "updated_at": "2023-01-02"}, "items": [{"item_id": 1, "item_name": "Item One"}, {"item_id": 2, "item_name": "Item Two"}]}`,
			expected: `{"metadata":{"created-at":"2023-01-01","updated-at":"2023-01-02"},"items":[{"item-id":1,"item-name":"Item One"},{"item-id":2,"item-name":"Item Two"}]}`,
		},
		{
			input:    `{"nested": {"level_one": {"level_two": {"level_three": {"key_one": "value1", "key_two": "value2"}}}}}`,
			expected: `{"nested":{"level-one":{"level-two":{"level-three":{"key-one":"value1","key-two":"value2"}}}}}`,
		},
		{
			input:    `{"array_of_objects": [{"name": "John_Doe", "age": 30}, {"name": "Jane_Doe", "age": 25}], "status": "active"}`,
			expected: `{"array-of-objects":[{"name":"John_Doe","age":30},{"name":"Jane_Doe","age":25}],"status":"active"}`,
		},
	}

	for _, test := range tests {
		var data map[string]interface{}
		if err := Unmarshal([]byte(test.input), &data); err != nil {
			t.Fatalf("Failed to unmarshal input: %v", err)
		}

		replacedData, err := ReplaceKeys(data, "_", "-")
		if err != nil {
			t.Fatalf("Error during replacement: %v", err)
		}

		var expectedData map[string]interface{}
		if err := Unmarshal([]byte(test.expected), &expectedData); err != nil {
			t.Fatalf("Failed to unmarshal expected output: %v", err)
		}

		if !reflect.DeepEqual(replacedData, expectedData) {
			t.Errorf("Expected %+v, but got %+v", expectedData, replacedData)
		}
	}
}

type JsonUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMarshalWithExtraField_Object(t *testing.T) {
	u := JsonUser{Name: "Alice", Age: 30}
	b, err := MarshalWithExtraField(u, "extra", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := Unmarshal(b, &result); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	// 验证原字段
	if result["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %v", result["name"])
	}
	if age, ok := result["age"].(float64); !ok || age != 30 {
		t.Errorf("expected age=30, got %v", result["age"])
	}

	// 验证额外字段
	if result["extra"] != "hello" {
		t.Errorf("expected extra=hello, got %v", result["extra"])
	}
}

func TestMarshalWithExtraField_Map(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	b, err := MarshalWithExtraField(m, "x", 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := Unmarshal(b, &result); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if result["a"] != float64(1) || result["b"] != float64(2) {
		t.Errorf("unexpected map values: %v", result)
	}
	if result["x"] != float64(123) {
		t.Errorf("expected x=123, got %v", result["x"])
	}
}

func TestMarshalWithExtraField_Array_Fail(t *testing.T) {
	arr := []int{1, 2, 3}
	_, err := MarshalWithExtraField(arr, "extra", "fail")
	if err == nil {
		t.Errorf("expected error for array input, got nil")
	}
}

func TestMarshalWithExtraField_String_Fail(t *testing.T) {
	s := "hello"
	_, err := MarshalWithExtraField(s, "extra", "fail")
	if err == nil {
		t.Errorf("expected error for string input, got nil")
	}
}

func TestMarshalWithExtraField_Number_Fail(t *testing.T) {
	n := 42
	_, err := MarshalWithExtraField(n, "extra", "fail")
	if err == nil {
		t.Errorf("expected error for number input, got nil")
	}
}

func TestCompact(t *testing.T) {
	cases := []string{
		`{ "a": 1, "b": 2 }`,
		`hello`,
		``,
		`[{"x":1},{"y":2}]`,
		`123`,
		`{"outer":{"inner":{"key":"value"}}, "arr":[1,2,3]}`,
		`{"a":{"b":{"c":{"d":4}}}}`,
	}

	for _, input := range cases {
		compacted := Compact([]byte(input))

		var expectedObj interface{}
		var actualObj interface{}

		// 尝试解析原始输入
		err1 := Unmarshal([]byte(input), &expectedObj)
		// 尝试解析压缩结果
		err2 := Unmarshal([]byte(compacted), &actualObj)

		if err1 == nil && err2 == nil {
			// 都是合法 JSON，比较解析后结构是否相等
			assert.Equal(t, expectedObj, actualObj)
		} else {
			// 不是 JSON，直接比较字符串
			assert.Equal(t, input, string(compacted))
		}
	}
}
