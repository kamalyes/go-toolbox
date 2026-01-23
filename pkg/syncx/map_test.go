/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-15 11:55:15
 * @FilePath: \go-toolbox\pkg\syncx\map_test.go
 * @Description: map 映射单元测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	m := NewMap[string, int]()

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
	size := m.Size()
	assert.Equal(t, 2, size, "Expected size 2 after deletion")

	// 测试 Keys
	keys := m.Keys()
	assert.Equal(t, 2, len(keys), "Expected 2 keys")
	assert.Contains(t, keys, "key3", "Expected keys to contain 'key3'")
	assert.Contains(t, keys, "key4", "Expected keys to contain 'key4'")

	// 测试 Values
	values := m.Values()
	assert.Equal(t, 2, len(values), "Expected 2 values")

	// 对返回的值进行排序
	sort.Ints(values)

	// 检查值是否为预期的内容
	expectedValues := []int{3, 4}
	sort.Ints(expectedValues)
	assert.Equal(t, expectedValues, values, "expected values to contain 3 and 4")

	// 测试 Clear
	m.Clear()
	size = m.Size()
	assert.Equal(t, 0, size, "Expected size 0 after clear")
}

func TestMapSwap(t *testing.T) {
	m := NewMap[string, int]()

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

func TestMapClear(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("key1", 1)
	m.Store("key2", 2)

	// 清空 Map
	m.Clear()

	// 验证 Map 为空
	size := m.Size()
	assert.Equal(t, 0, size, "Expected size 0 after clear")

	// 验证 Load 方法返回值
	_, ok := m.Load("key1")
	assert.False(t, ok, "expected key1 to be deleted after clear")
	_, ok = m.Load("key2")
	assert.False(t, ok, "expected key2 to be deleted after clear")
}

func TestMapLoadAndDelete(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("key1", 1)

	// 测试 LoadAndDelete
	val, ok := m.LoadAndDelete("key1")
	assert.True(t, ok, "expected key1 to exist")
	assert.Equal(t, 1, val, "expected value to be 1")

	// 再次尝试加载已删除的键
	_, ok = m.Load("key1")
	assert.False(t, ok, "expected key1 to be deleted")
}

func TestMapEquals(t *testing.T) {
	m := NewMap[string, int]()
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

func TestMapSize_Concurrent(t *testing.T) {
	m := NewMap[string, int]()

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
	size := m.Size()
	assert.Equal(t, 100, size, "Expected size 100")
}

func TestMapKeysAndValues(t *testing.T) {
	m := NewMap[string, int]()
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

	CopyMeta(src, dst)

	assert.Equal(t, "value1", dst["key1"], "expected dst['key1'] = 'value1'")
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
		result := MetaStringToMap(test.meta)
		assert.True(t, mapsEqual(result, test.expected),
			"MetaStringToMap(%q) = %v; want %v", test.meta, result, test.expected)
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
		result := MetaMapToString(test.meta)
		assert.Equal(t, test.expected, result,
			"MetaMapToString(%v) should equal expected", test.meta)
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

func TestMapClone(t *testing.T) {
	// 创建一个新的 Map 实例并添加一些键值对
	originalMap := NewMap[string, int]()
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

// TestMapFilter 测试 Filter 方法
func TestMapFilter(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)
	m.Store("d", 4)

	// 过滤出值大于 2 的元素
	result := m.Filter(func(k string, v int) bool {
		return v > 2
	})

	sort.Ints(result) // 排序以便比较
	assert.Equal(t, []int{3, 4}, result, "Filter 应返回值大于 2 的元素")

	// 测试空结果
	emptyResult := m.Filter(func(k string, v int) bool {
		return v > 10
	})
	assert.Empty(t, emptyResult, "Filter 应返回空切片")
}

// TestMapFilterKeys 测试 FilterKeys 方法
func TestMapFilterKeys(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("apple", 10)
	m.Store("banana", 20)
	m.Store("cherry", 30)

	// 过滤出键包含 'a' 的元素
	result := m.FilterKeys(func(k string, v int) bool {
		return k[0] == 'a' || k[0] == 'b'
	})

	sort.Strings(result)
	assert.Equal(t, []string{"apple", "banana"}, result, "FilterKeys 应返回以 a 或 b 开头的键")
}

// TestMapFilterMap 测试 FilterMap 方法
func TestMapFilterMap(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	// 过滤出值为偶数的元素
	filteredMap := m.FilterMap(func(k string, v int) bool {
		return v%2 == 0
	})

	assert.Equal(t, 1, filteredMap.Size(), "FilterMap 应返回 1 个元素")
	val, ok := filteredMap.Load("b")
	assert.True(t, ok, "b 应该存在")
	assert.Equal(t, 2, val, "b 的值应为 2")

	_, notExists := filteredMap.Load("a")
	assert.False(t, notExists, "a 不应该存在")
}

// TestMapForEach 测试 ForEach 方法
func TestMapForEach(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("x", 10)
	m.Store("y", 20)
	m.Store("z", 30)

	sum := 0
	m.ForEach(func(k string, v int) {
		sum += v
	})

	assert.Equal(t, 60, sum, "ForEach 应计算所有值的总和")
}

// TestMapUpdate 测试 Update 方法
func TestMapUpdate(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("counter", 10)

	// 更新存在的键
	updated := m.Update("counter", func(v int) int {
		return v + 5
	})
	assert.True(t, updated, "Update 应返回 true")

	val, _ := m.Load("counter")
	assert.Equal(t, 15, val, "counter 应被更新为 15")

	// 更新不存在的键
	notUpdated := m.Update("nonexistent", func(v int) int {
		return v + 1
	})
	assert.False(t, notUpdated, "Update 不存在的键应返回 false")
}

// TestMapGetOrStore 测试 GetOrStore 方法
func TestMapGetOrStore(t *testing.T) {
	m := NewMap[string, string]()
	m.Store("key1", "value1")

	// 获取已存在的键
	val := m.GetOrStore("key1", "default")
	assert.Equal(t, "value1", val, "GetOrStore 应返回已存在的值")

	// 获取不存在的键（使用默认值）
	val2 := m.GetOrStore("key2", "default")
	assert.Equal(t, "default", val2, "GetOrStore 应存储并返回默认值")

	// 验证默认值已被存储
	val3, ok := m.Load("key2")
	assert.True(t, ok, "key2 应该存在")
	assert.Equal(t, "default", val3, "key2 的值应为 default")
}

// TestMapGetOrCompute 测试 GetOrCompute 方法
func TestMapGetOrCompute(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("existing", 100)

	computeCallCount := 0
	compute := func() int {
		computeCallCount++
		return 42
	}

	// 获取已存在的键（不应调用 compute）
	val := m.GetOrCompute("existing", compute)
	assert.Equal(t, 100, val, "GetOrCompute 应返回已存在的值")
	assert.Equal(t, 0, computeCallCount, "compute 不应被调用")

	// 获取不存在的键（应调用 compute）
	val2 := m.GetOrCompute("new", compute)
	assert.Equal(t, 42, val2, "GetOrCompute 应返回计算的值")
	assert.Equal(t, 1, computeCallCount, "compute 应被调用一次")

	// 验证计算的值已被存储
	val3, ok := m.Load("new")
	assert.True(t, ok, "new 应该存在")
	assert.Equal(t, 42, val3, "new 的值应为 42")
}

// TestMapDeleteIf 测试 DeleteIf 方法
func TestMapDeleteIf(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)
	m.Store("d", 4)

	// 删除值为偶数的元素
	deletedCount := m.DeleteIf(func(k string, v int) bool {
		return v%2 == 0
	})

	assert.Equal(t, 2, deletedCount, "DeleteIf 应删除 2 个元素")
	assert.Equal(t, 2, m.Size(), "Map 应剩余 2 个元素")

	// 验证奇数元素仍存在
	_, ok1 := m.Load("a")
	assert.True(t, ok1, "a 应该存在")
	_, ok3 := m.Load("c")
	assert.True(t, ok3, "c 应该存在")

	// 验证偶数元素已删除
	_, ok2 := m.Load("b")
	assert.False(t, ok2, "b 不应该存在")
	_, ok4 := m.Load("d")
	assert.False(t, ok4, "d 不应该存在")
}

// TestMapAny 测试 Any 方法
func TestMapAny(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	// 存在满足条件的元素
	result := m.Any(func(k string, v int) bool {
		return v > 2
	})
	assert.True(t, result, "Any 应返回 true（存在值大于 2 的元素）")

	// 不存在满足条件的元素
	result2 := m.Any(func(k string, v int) bool {
		return v > 10
	})
	assert.False(t, result2, "Any 应返回 false（不存在值大于 10 的元素）")

	// 空 Map
	emptyMap := NewMap[string, int]()
	result3 := emptyMap.Any(func(k string, v int) bool {
		return true
	})
	assert.False(t, result3, "空 Map 的 Any 应返回 false")
}

// TestMapAll 测试 All 方法
func TestMapAll(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("a", 2)
	m.Store("b", 4)
	m.Store("c", 6)

	// 所有元素都满足条件
	result := m.All(func(k string, v int) bool {
		return v%2 == 0
	})
	assert.True(t, result, "All 应返回 true（所有值都是偶数）")

	// 存在不满足条件的元素
	m.Store("d", 3)
	result2 := m.All(func(k string, v int) bool {
		return v%2 == 0
	})
	assert.False(t, result2, "All 应返回 false（存在奇数）")

	// 空 Map
	emptyMap := NewMap[string, int]()
	result3 := emptyMap.All(func(k string, v int) bool {
		return false
	})
	assert.True(t, result3, "空 Map 的 All 应返回 true")
}

// TestMapCount 测试 Count 方法
func TestMapCount(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)
	m.Store("d", 4)

	// 计算值大于 2 的元素数量
	count := m.Count(func(k string, v int) bool {
		return v > 2
	})
	assert.Equal(t, 2, count, "Count 应返回 2（值大于 2 的元素有 2 个）")

	// 计算所有元素
	totalCount := m.Count(func(k string, v int) bool {
		return true
	})
	assert.Equal(t, 4, totalCount, "Count 应返回 4（总共 4 个元素）")

	// 不满足条件的元素
	zeroCount := m.Count(func(k string, v int) bool {
		return v > 10
	})
	assert.Equal(t, 0, zeroCount, "Count 应返回 0（没有值大于 10 的元素）")
}

// TestMapIsEmpty 测试 IsEmpty 方法
func TestMapIsEmpty(t *testing.T) {
	m := NewMap[string, int]()

	// 空 Map
	assert.True(t, m.IsEmpty(), "新创建的 Map 应该为空")

	// 添加元素后
	m.Store("key", 1)
	assert.False(t, m.IsEmpty(), "添加元素后 Map 不应为空")

	// 删除元素后
	m.Delete("key")
	assert.True(t, m.IsEmpty(), "删除所有元素后 Map 应该为空")
}

// TestMapToMap 测试 ToMap 方法
func TestMapToMap(t *testing.T) {
	m := NewMap[string, int]()
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	regularMap := m.ToMap()

	assert.Equal(t, 3, len(regularMap), "ToMap 应返回包含 3 个元素的 map")
	assert.Equal(t, 1, regularMap["a"], "a 的值应为 1")
	assert.Equal(t, 2, regularMap["b"], "b 的值应为 2")
	assert.Equal(t, 3, regularMap["c"], "c 的值应为 3")

	// 修改普通 map 不应影响 syncx.Map
	regularMap["d"] = 4
	_, exists := m.Load("d")
	assert.False(t, exists, "syncx.Map 不应包含 d")
}

// TestFromMap 测试 FromMap 方法
func TestFromMap(t *testing.T) {
	regularMap := map[string]int{
		"x": 10,
		"y": 20,
		"z": 30,
	}

	m := FromMap(regularMap)

	assert.Equal(t, 3, m.Size(), "FromMap 应创建包含 3 个元素的 syncx.Map")

	val1, ok1 := m.Load("x")
	assert.True(t, ok1, "x 应该存在")
	assert.Equal(t, 10, val1, "x 的值应为 10")

	val2, ok2 := m.Load("y")
	assert.True(t, ok2, "y 应该存在")
	assert.Equal(t, 20, val2, "y 的值应为 20")

	val3, ok3 := m.Load("z")
	assert.True(t, ok3, "z 应该存在")
	assert.Equal(t, 30, val3, "z 的值应为 30")
}

// BenchmarkMapFilter 基准测试：Filter 方法
func BenchmarkMapFilter(b *testing.B) {
	m := NewMap[int, int]()
	for i := 0; i < 1000; i++ {
		m.Store(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Filter(func(k, v int) bool {
			return v%2 == 0
		})
	}
}

// BenchmarkMapForEach 基准测试：ForEach 方法
func BenchmarkMapForEach(b *testing.B) {
	m := NewMap[int, int]()
	for i := 0; i < 1000; i++ {
		m.Store(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		m.ForEach(func(k, v int) {
			sum += v
		})
	}
}

// BenchmarkMapToMap 基准测试：ToMap 方法
func BenchmarkMapToMap(b *testing.B) {
	m := NewMap[int, int]()
	for i := 0; i < 1000; i++ {
		m.Store(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.ToMap()
	}
}
