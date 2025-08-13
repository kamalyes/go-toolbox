/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-15 11:55:15
 * @FilePath: \go-toolbox\tests\syncx_map_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"sort"
	"strconv"
	"testing"
	"time"

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

	// 测试 Size
	if size := m.Size(); size != 2 {
		t.Errorf("Expected size 2, got %d", size)
	}

	// 测试 Keys
	keys := m.Keys()
	if len(keys) != 2 || (keys[0] != "key3" && keys[1] != "key3") {
		t.Errorf("Expected keys to contain 'key3' and 'key4', got %v", keys)
	}

	// 测试 Values
	values := m.Values()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}

	// 对返回的值进行排序
	sort.Ints(values)

	// 检查值是否为预期的内容
	expectedValues := []int{3, 4}
	sort.Ints(expectedValues)
	assert.Equal(t, expectedValues, values, "expected values to contain 3 and 4, got %v", values)

	// 测试 Clear
	m.Clear()
	if size := m.Size(); size != 0 {
		t.Errorf("Expected size 0 after clear, got %d", size)
	}
}

func TestMap_Swap(t *testing.T) {
	m := syncx.NewMap[string, int]()

	// 测试 Swap 时键不存在
	pre, ok := m.Swap("key1", 10)
	assert.Equal(t, 0, pre, "expected pre to be 0 for non-existing key")
	assert.False(t, ok, "expected ok to be false for non-existing key")

	// 存储一个值
	m.Store("key1", 5)

	// 测试 Swap 时键存在
	pre, ok = m.Swap("key1", 10)
	assert.Equal(t, 5, pre, "expected pre to be 5 for existing key")
	assert.True(t, ok, "expected ok to be true for existing key")

	// 确认值已被替换
	val, ok := m.Load("key1")
	assert.True(t, ok, "expected key1 to exist")
	assert.Equal(t, 10, val, "expected value to be 10 after swap")
}

func TestMap_Clear(t *testing.T) {
	m := syncx.NewMap[string, int]()
	m.Store("key1", 1)
	m.Store("key2", 2)

	// 清空 Map
	m.Clear()

	// 验证 Map 为空
	if size := m.Size(); size != 0 {
		t.Errorf("Expected size 0 after clear, got %d", size)
	}

	// 验证 Load 方法返回值
	_, ok := m.Load("key1")
	assert.False(t, ok, "expected key1 to be deleted after clear")
	_, ok = m.Load("key2")
	assert.False(t, ok, "expected key2 to be deleted after clear")
}

func TestMap_LoadAndDelete(t *testing.T) {
	m := syncx.NewMap[string, int]()
	m.Store("key1", 1)

	// 测试 LoadAndDelete
	val, ok := m.LoadAndDelete("key1")
	assert.True(t, ok, "expected key1 to exist")
	assert.Equal(t, 1, val, "expected value to be 1")

	// 再次尝试加载已删除的键
	_, ok = m.Load("key1")
	assert.False(t, ok, "expected key1 to be deleted")
}

func TestMap_Equals(t *testing.T) {
	m := syncx.NewMap[string, int]()
	m.Store("key1", 1)

	// 测试存在的键
	isEqual := func(existing int) bool {
		return existing == 1
	}
	assert.True(t, m.Equals("key1", 1, isEqual), "expected key1 to be equal to 1")

	// 测试不存在的键
	assert.False(t, m.Equals("key2", 1, isEqual), "expected key2 to not exist")

	// 测试不同的比较函数
	isEqualDifferent := func(existing int) bool {
		return existing == 2
	}
	assert.False(t, m.Equals("key1", 1, isEqualDifferent), "expected key1 to not be equal to 2")
}

func TestMap_Size_Concurrent(t *testing.T) {
	m := syncx.NewMap[string, int]()

	// 启动多个 goroutine 来并发存储值
	for i := 0; i < 100; i++ {
		go func(i int) {
			m.Store("key"+strconv.Itoa(i), i) // 使用 strconv.Itoa 将整数转换为字符串
		}(i)
	}

	// 等待所有 goroutine 完成（可以使用 sync.WaitGroup 更好地管理）
	// 这里简单使用 Sleep 来确保所有操作完成
	time.Sleep(1 * time.Second)

	// 验证 Size
	if size := m.Size(); size != 100 {
		t.Errorf("Expected size 100, got %d", size)
	}
}

func TestMap_KeysAndValues(t *testing.T) {
	m := syncx.NewMap[string, int]()
	m.Store("key1", 1)
	m.Store("key2", 2)
	m.Store("key3", 3)

	// 测试 Keys
	keys := m.Keys()
	assert.ElementsMatch(t, []string{"key1", "key2", "key3"}, keys, "expected keys to match")

	// 测试 Values
	values := m.Values()
	assert.ElementsMatch(t, []int{1, 2, 3}, values, "expected values to match")
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

func TestMap_Clone(t *testing.T) {
	// 创建一个新的 Map 实例并添加一些键值对
	originalMap := syncx.NewMap[string, int]()
	originalMap.Store("key1", 1)
	originalMap.Store("key2", 2)
	originalMap.Store("key3", 3)

	// 克隆原始 Map
	clonedMap := originalMap.Clone()

	// 验证克隆后的 Map 是否与原始 Map 相同
	clonedMap.Range(func(key string, value int) bool {
		originalValue, ok := originalMap.Load(key)
		assert.True(t, ok, "Key %s should exist in the original map", key)
		assert.Equal(t, originalValue, value, "Value for key %s should match", key)
		return true
	})

	// 验证克隆后的 Map 是否是独立的
	clonedMap.Store("key4", 4) // 在克隆的 Map 中添加新键
	_, originalExists := originalMap.Load("key4")
	assert.False(t, originalExists, "Original map should not contain key4 after cloning")

	// 验证原始 Map 的值未被改变
	originalMap.Range(func(key string, value int) bool {
		if key == "key4" {
			assert.Fail(t, "Original map should not contain key4")
		}
		return true
	})
}
