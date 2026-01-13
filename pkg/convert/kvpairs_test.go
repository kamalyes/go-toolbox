/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-13 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-13 10:36:15
 * @FilePath: \go-toolbox\pkg\convert\kvpairs_test.go
 * @Description: 键值对转换工具测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToKVPairs(t *testing.T) {
	// 正常转换
	m := map[string]interface{}{"key1": "value1", "key2": 123, "key3": true}
	result := MapToKVPairs(m)
	assert.Len(t, result, 6)
	assert.Equal(t, m, KVPairsToMap(result))

	// 空map
	assert.Nil(t, MapToKVPairs(map[string]interface{}{}))

	// nil map
	assert.Nil(t, MapToKVPairs(nil))

	// 单个键值对
	single := map[string]interface{}{"single": "value"}
	result = MapToKVPairs(single)
	assert.Len(t, result, 2)
	assert.Equal(t, single, KVPairsToMap(result))
}

func TestMapStringToKVPairs(t *testing.T) {
	m := map[string]string{"name": "test", "email": "test@example.com"}
	result := MapStringToKVPairs(m)
	assert.Len(t, result, 4)
	assert.Nil(t, MapStringToKVPairs(map[string]string{}))
}

func TestMapAnyToKVPairs(t *testing.T) {
	m := map[string]any{"key1": "value1", "key2": 123}
	assert.Len(t, MapAnyToKVPairs(m), 4)
	assert.Nil(t, MapAnyToKVPairs(map[string]any{}))
}

func TestMergeKVPairs(t *testing.T) {
	kv1 := []interface{}{"key1", "value1"}
	kv2 := []interface{}{"key2", "value2"}
	kv3 := []interface{}{"key3", "value3"}

	assert.Len(t, MergeKVPairs(kv1, kv2), 4)
	assert.Len(t, MergeKVPairs(kv1, kv2, kv3), 6)
	assert.Nil(t, MergeKVPairs())
	assert.Equal(t, kv1, MergeKVPairs(kv1))
}

func TestKVPairs(t *testing.T) {
	result := KVPairs("key1", "value1", "key2", 123)
	assert.Len(t, result, 4)
	assert.Equal(t, []interface{}{"key1", "value1", "key2", 123}, result)
	assert.Len(t, KVPairs(), 0)
}

func TestKVPairsToMap(t *testing.T) {
	// 正常转换
	kvs := []interface{}{"key1", "value1", "key2", 123}
	expected := map[string]interface{}{"key1": "value1", "key2": 123}
	assert.Equal(t, expected, KVPairsToMap(kvs))

	// 边界情况
	assert.Nil(t, KVPairsToMap([]interface{}{}))
	assert.Nil(t, KVPairsToMap(nil))

	// 奇数个元素
	assert.Equal(t, map[string]interface{}{"key1": "value1"},
		KVPairsToMap([]interface{}{"key1", "value1", "key2"}))

	// 非字符串键被跳过
	assert.Equal(t, map[string]interface{}{"key1": "value1"},
		KVPairsToMap([]interface{}{"key1", "value1", 123, "value2"}))
}

func TestAddKVPair(t *testing.T) {
	result := AddKVPair([]interface{}{"key1", "value1"}, "key2", "value2")
	assert.Len(t, result, 4)
	assert.Equal(t, "key2", result[2])
	assert.Equal(t, "value2", result[3])

	result = AddKVPair([]interface{}{}, "key1", "value1")
	assert.Len(t, result, 2)
}

func TestAddKVPairs(t *testing.T) {
	initial := []interface{}{"key1", "value1"}
	result := AddKVPairs(initial, map[string]interface{}{"key2": "value2", "key3": "value3"})
	assert.Len(t, result, 6)

	result = AddKVPairs(initial, map[string]interface{}{})
	assert.Len(t, result, 2)
}

func TestRoundTrip(t *testing.T) {
	original := map[string]interface{}{"key1": "value1", "key2": 123, "key3": true}
	assert.Equal(t, original, KVPairsToMap(MapToKVPairs(original)))
}

// ============= 基准测试 =============

func BenchmarkMapToKVPairs(b *testing.B) {
	m := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapToKVPairs(m)
	}
}

func BenchmarkMapStringToKVPairs(b *testing.B) {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapStringToKVPairs(m)
	}
}

func BenchmarkKVPairsToMap(b *testing.B) {
	kvs := []interface{}{"key1", "value1", "key2", "value2", "key3", "value3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = KVPairsToMap(kvs)
	}
}

func BenchmarkMergeKVPairs(b *testing.B) {
	kv1 := []interface{}{"key1", "value1", "key2", "value2"}
	kv2 := []interface{}{"key3", "value3", "key4", "value4"}
	kv3 := []interface{}{"key5", "value5", "key6", "value6"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MergeKVPairs(kv1, kv2, kv3)
	}
}

func BenchmarkAddKVPair(b *testing.B) {
	initial := []interface{}{"key1", "value1", "key2", "value2"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AddKVPair(initial, "key3", "value3")
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	m := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kvs := MapToKVPairs(m)
		_ = KVPairsToMap(kvs)
	}
}

func TestMergeMapToKVPairs(t *testing.T) {
	// 合并两个 map
	map1 := map[string]interface{}{"event": "login", "user_id": "123"}
	map2 := map[string]interface{}{"ip": "127.0.0.1", "device": "mobile"}
	result := MergeMapToKVPairs(map1, map2)
	assert.Len(t, result, 8)

	// 验证所有键都存在
	resultMap := KVPairsToMap(result)
	assert.Equal(t, "login", resultMap["event"])
	assert.Equal(t, "123", resultMap["user_id"])
	assert.Equal(t, "127.0.0.1", resultMap["ip"])
	assert.Equal(t, "mobile", resultMap["device"])

	// 空 maps
	assert.Nil(t, MergeMapToKVPairs())

	// 单个 map
	single := MergeMapToKVPairs(map1)
	assert.Len(t, single, 4)
	assert.Equal(t, map1, KVPairsToMap(single))

	// 合并三个 map
	map3 := map[string]interface{}{"status": "success"}
	result = MergeMapToKVPairs(map1, map2, map3)
	assert.Len(t, result, 10)
	resultMap = KVPairsToMap(result)
	assert.Equal(t, "success", resultMap["status"])

	// 包含空 map
	result = MergeMapToKVPairs(map1, map[string]interface{}{}, map3)
	assert.Len(t, result, 6)
}

func BenchmarkMergeMapToKVPairs(b *testing.B) {
	map1 := map[string]interface{}{"key1": "value1", "key2": "value2"}
	map2 := map[string]interface{}{"key3": "value3", "key4": "value4"}
	map3 := map[string]interface{}{"key5": "value5", "key6": "value6"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MergeMapToKVPairs(map1, map2, map3)
	}
}
