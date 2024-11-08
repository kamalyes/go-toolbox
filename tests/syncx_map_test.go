/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 12:28:17
 * @FilePath: \go-toolbox\tests\syncx_map_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	m := syncx.NewMap[string, int]()

	// 测试 Store 和 Load
	m.Store("key1", 1)
	val, ok := m.Load("key1")
	assert.True(t, ok, "expected key1 to exist")
	assert.Equal(t, 1, val, "expected value to be 1")

	// 测试 LoadOrStore
	val, ok = m.LoadOrStore("key1", 2)
	assert.True(t, ok, "expected key1 to exist")
	assert.Equal(t, 1, val, "expected value to be 1")

	val, ok = m.LoadOrStore("key2", 2)
	assert.False(t, ok, "expected key2 to not exist")
	assert.Equal(t, 2, val, "expected value to be 2")

	// 测试 CompareAndSwap
	ok = m.CompareAndSwap("key1", 1, 3)
	assert.True(t, ok, "expected CompareAndSwap to succeed")
	val, ok = m.Load("key1")
	assert.True(t, ok, "expected key1 to exist")
	assert.Equal(t, 3, val, "expected value to be 3")

	// 测试 CompareAndDelete
	ok = m.CompareAndDelete("key1", 3)
	assert.True(t, ok, "expected CompareAndDelete to succeed")
	_, ok = m.Load("key1")
	assert.False(t, ok, "expected key1 to be deleted")

	// 测试 Equals 方法
	m.Store("key2", 5)
	isEqual := func(existing int) bool {
		return existing == 5
	}
	assert.True(t, m.Equals("key2", 5, isEqual), "expected key2 to be equal to 5")

	// 测试 Range
	m.Store("key3", 3)
	m.Store("key4", 4)
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true // 继续迭代
	})
	assert.Equal(t, 3, count, "expected 3 items in the map")

	// 测试 Delete
	m.Delete("key2")
	_, ok = m.Load("key2")
	assert.False(t, ok, "expected key2 to be deleted")
}

func TestCopyMetaWithExistingKeys(t *testing.T) {
	src := map[string]string{
		"key1": "value1",
	}
	dst := map[string]string{
		"key1": "old_value",
	}

	syncx.CopyMeta(src, dst)

	if dst["key1"] != "value1" {
		t.Errorf("expected dst['key1'] = 'value1', got '%s'", dst["key1"])
	}
}

func TestMetaStringToMap(t *testing.T) {
	tests := []struct {
		meta     string
		expected map[string]string
	}{
		{"key1=value1&key2=value2", map[string]string{"key1": "value1", "key2": "value2"}},
		{"key2=value2&key1=value1", map[string]string{"key1": "value1", "key2": "value2"}}, // 顺序不同
		{"", map[string]string{}},                                        // 空字符串
		{"invalid_query_string", map[string]string{}},                    // 无法解析的字符串
		{"key1=value1&key1=value2", map[string]string{"key1": "value1"}}, // 重复键
	}

	for _, test := range tests {
		result := syncx.MetaStringToMap(test.meta)
		if !mapsEqual(result, test.expected) {
			t.Errorf("MetaStringToMap(%q) = %v; want %v", test.meta, result, test.expected)
		}
	}
}

func TestMetaMapToString(t *testing.T) {
	tests := []struct {
		meta     map[string]string
		expected string
	}{
		{map[string]string{"key1": "value1", "key2": "value2"}, "key1=value1&key2=value2"},
		{map[string]string{"key2": "value2", "key1": "value1"}, "key1=value1&key2=value2"}, // 顺序不同
		{map[string]string{}, ""}, // 空映射
		{map[string]string{"key1": "value with spaces"}, "key1=value+with+spaces"},                                       // 处理空格
		{map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}, "key1=value1&key2=value2&key3=value3"}, // 多个键值对
	}

	for _, test := range tests {
		result := syncx.MetaMapToString(test.meta)
		if result != test.expected {
			t.Errorf("MetaMapToString(%v) = %q; want %q", test.meta, result, test.expected)
		}
	}
}

// 辅助函数：比较两个映射是否相等
func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if bValue, exists := b[key]; !exists || value != bValue {
			return false
		}
	}
	return true
}
