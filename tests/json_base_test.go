/*
 * @Author: kamalyes 501893067@qq.com
 * @Date:2024-10-24 10:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-03 22:53:55
 * @FilePath: \go-toolbox\tests\json_base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"reflect"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/json"
)

// 辅助函数，创建测试用例
func newTestCase(originalJSON string, pairs *json.KeyValuePairs, expected map[string]interface{}) struct {
	originalJSON string
	pairs        *json.KeyValuePairs
	expected     map[string]interface{}
} {
	return struct {
		originalJSON string
		pairs        *json.KeyValuePairs
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
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("解析 JSON 失败: %v", err)
	}
	return result
}

// 测试将键值对追加到 JSON 字符串的功能
func TestAppendKeysToJSON(t *testing.T) {
	tests := []struct {
		originalJSON string
		pairs        *json.KeyValuePairs
		expected     map[string]interface{}
	}{
		newTestCase(
			`{"name": "Alice", "age": 30}`,
			json.NewKeyValuePairs().Add("city", "New York").Add("country", "USA"),
			map[string]interface{}{
				"name":    "Alice",
				"age":     float64(30),
				"city":    "New York",
				"country": "USA",
			},
		),
		newTestCase(
			`{}`,
			json.NewKeyValuePairs().Add("key1", "value1"),
			map[string]interface{}{
				"key1": "value1",
			},
		),
		newTestCase(
			`{"existing": "value"}`,
			json.NewKeyValuePairs().Add("newKey", "newValue"),
			map[string]interface{}{
				"existing": "value",
				"newKey":   "newValue",
			},
		),
		newTestCase(
			`{"a": 1}`,
			json.NewKeyValuePairs().Add("b", 2).Add("c", 3),
			map[string]interface{}{
				"a": float64(1),
				"b": float64(2),
				"c": float64(3),
			},
		),
	}

	// 遍历测试用例
	for _, tt := range tests {
		updatedJSON, err := json.AppendKeysToJSON(tt.originalJSON, tt.pairs)
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
		if err := json.Unmarshal([]byte(test.input), &data); err != nil {
			t.Fatalf("Failed to unmarshal input: %v", err)
		}

		replacedData, err := json.ReplaceKeys(data, "_", "-")
		if err != nil {
			t.Fatalf("Error during replacement: %v", err)
		}

		var expectedData map[string]interface{}
		if err := json.Unmarshal([]byte(test.expected), &expectedData); err != nil {
			t.Fatalf("Failed to unmarshal expected output: %v", err)
		}

		if !reflect.DeepEqual(replacedData, expectedData) {
			t.Errorf("Expected %+v, but got %+v", expectedData, replacedData)
		}
	}
}
