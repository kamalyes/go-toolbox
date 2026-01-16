/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-16 13:35:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-16 21:15:00
 * @FilePath: \go-toolbox\pkg\httpx\params_test.go
 * @Description: 参数构建器测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParams(t *testing.T) {
	builder := NewParams()

	assert.NotNil(t, builder, "builder 不应为 nil")
	assert.NotNil(t, builder.params, "builder.params 不应为 nil")
	assert.Equal(t, 0, len(builder.params), "初始 params 长度应为 0")
}

func TestNewParamsWithBase(t *testing.T) {
	base := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	builder := NewParamsWithBase(base)

	assert.NotNil(t, builder, "builder 不应为 nil")
	assert.Equal(t, 2, len(builder.params), "基于base创建的params长度应为2")
	assert.Equal(t, "value1", builder.params["key1"], "key1 取值错误")
	assert.Equal(t, "value2", builder.params["key2"], "key2 取值错误")

	// 确保是深拷贝，修改原 map 不影响 builder
	base["key1"] = "modified"
	assert.Equal(t, "value1", builder.params["key1"], "深拷贝验证失败：原map修改影响了builder")
}

func TestParamsBuilder_Set(t *testing.T) {
	params := NewParams().
		Set("key1", "value1").
		Set("key2", "value2").
		Build()

	assert.Equal(t, 2, len(params), "参数数量应为 2")
	assert.Equal(t, "value1", params["key1"])
	assert.Equal(t, "value2", params["key2"])
}

func TestParamsBuilder_Add(t *testing.T) {
	params := NewParams().
		Add("key1", "value1").
		Add("key2", "value2").
		Build()

	assert.Equal(t, 2, len(params), "参数数量应为 2")
	assert.Equal(t, "value1", params["key1"])
	assert.Equal(t, "value2", params["key2"])
}

func TestParamsBuilder_Delete(t *testing.T) {
	builder := NewParams().
		Set("key1", "value1").
		Set("key2", "value2")

	assert.Equal(t, 2, builder.Len())

	builder.Delete("key1")
	assert.Equal(t, 1, builder.Len())
	assert.False(t, builder.Has("key1"))
	assert.True(t, builder.Has("key2"))
}

func TestParamsBuilder_Get(t *testing.T) {
	builder := NewParams().
		Set("key1", "value1")

	value, ok := builder.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", value)

	value, ok = builder.Get("nonexistent")
	assert.False(t, ok)
	assert.Equal(t, "", value)
}

func TestParamsBuilder_Has(t *testing.T) {
	builder := NewParams().
		Set("key1", "value1")

	assert.True(t, builder.Has("key1"))
	assert.False(t, builder.Has("key2"))
}

func TestParamsBuilder_Clear(t *testing.T) {
	builder := NewParams().
		Set("key1", "value1").
		Set("key2", "value2")

	assert.Equal(t, 2, builder.Len())

	builder.Clear()
	assert.Equal(t, 0, builder.Len())
	assert.False(t, builder.Has("key1"))
	assert.False(t, builder.Has("key2"))
}

func TestParamsBuilder_Len(t *testing.T) {
	builder := NewParams()
	assert.Equal(t, 0, builder.Len())

	builder.Set("key1", "value1")
	assert.Equal(t, 1, builder.Len())

	builder.Set("key2", "value2")
	assert.Equal(t, 2, builder.Len())

	builder.Delete("key1")
	assert.Equal(t, 1, builder.Len())
}

func TestParamsBuilder_Clone(t *testing.T) {
	original := NewParams().
		Set("key1", "value1").
		Set("key2", "value2")

	cloned := original.Clone()

	assert.Equal(t, original.Len(), cloned.Len())
	assert.Equal(t, original.Build(), cloned.Build())

	// 修改克隆不应影响原始
	cloned.Set("key3", "value3")
	assert.Equal(t, 2, original.Len())
	assert.Equal(t, 3, cloned.Len())
}

func TestParamsBuilder_ToSlice(t *testing.T) {
	builder := NewParams().
		Set("key1", "value1").
		Set("key2", "value2")

	slice := builder.ToSlice()

	assert.Equal(t, 4, len(slice), "切片长度应为 4")
	assert.Contains(t, slice, "key1")
	assert.Contains(t, slice, "value1")
	assert.Contains(t, slice, "key2")
	assert.Contains(t, slice, "value2")
}

func TestParamsBuilder_Keys(t *testing.T) {
	builder := NewParams().
		Set("key1", "value1").
		Set("key2", "value2")

	keys := builder.Keys()

	assert.Equal(t, 2, len(keys))
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
}

func TestParamsBuilder_Values(t *testing.T) {
	builder := NewParams().
		Set("key1", "value1").
		Set("key2", "value2")

	values := builder.Values()

	assert.Equal(t, 2, len(values))
	assert.Contains(t, values, "value1")
	assert.Contains(t, values, "value2")
}

func TestParamsBuilder_Merge(t *testing.T) {
	builder1 := NewParams().
		Set("key1", "value1").
		Set("key2", "value2")

	builder2 := NewParams().
		Set("key2", "updated_value2").
		Set("key3", "value3")

	builder1.Merge(builder2)

	assert.Equal(t, 3, builder1.Len())
	assert.Equal(t, "value1", builder1.params["key1"])
	assert.Equal(t, "updated_value2", builder1.params["key2"], "应该覆盖原有值")
	assert.Equal(t, "value3", builder1.params["key3"])

	// 测试 nil 合并
	builder1.Merge(nil)
	assert.Equal(t, 3, builder1.Len(), "合并 nil 不应改变参数")
}

func TestParamsBuilder_SetAny(t *testing.T) {
	builder := NewParams().
		SetAny("int", 123).
		SetAny("float", 45.67).
		SetAny("bool", true).
		SetAny("string", "test")

	assert.Equal(t, "123", builder.params["int"])
	assert.Equal(t, "45.67", builder.params["float"])
	assert.Equal(t, "true", builder.params["bool"])
	assert.Equal(t, "test", builder.params["string"])
}

func TestParamsBuilder_SetAnyIf(t *testing.T) {
	builder := NewParams().
		SetAnyIf(true, "key1", 123).
		SetAnyIf(false, "key2", 456)

	assert.True(t, builder.Has("key1"))
	assert.Equal(t, "123", builder.params["key1"])
	assert.False(t, builder.Has("key2"))
}

func TestParamsBuilder_SetIf(t *testing.T) {
	builder := NewParams().
		SetIf(true, "key1", "value1").
		SetIf(false, "key2", "value2")

	assert.True(t, builder.Has("key1"))
	assert.Equal(t, "value1", builder.params["key1"])
	assert.False(t, builder.Has("key2"))
}

func TestParamsBuilder_SetNotEmpty(t *testing.T) {
	builder := NewParams().
		SetNotEmpty("key1", "value1").
		SetNotEmpty("key2", "")

	assert.True(t, builder.Has("key1"))
	assert.Equal(t, "value1", builder.params["key1"])
	assert.False(t, builder.Has("key2"))
}

func TestParamsBuilder_SetMultiple(t *testing.T) {
	params := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	builder := NewParams().SetMultiple(params)

	assert.Equal(t, 3, builder.Len())
	assert.Equal(t, "value1", builder.params["key1"])
	assert.Equal(t, "value2", builder.params["key2"])
	assert.Equal(t, "value3", builder.params["key3"])
}

func TestParamsBuilder_ChainedCalls(t *testing.T) {
	params := NewParams().
		Set("key1", "value1").
		SetIf(true, "key2", "value2").
		SetNotEmpty("key3", "value3").
		SetAny("key4", 123).
		Build()

	assert.Equal(t, 4, len(params))
	assert.Equal(t, "value1", params["key1"])
	assert.Equal(t, "value2", params["key2"])
	assert.Equal(t, "value3", params["key3"])
	assert.Equal(t, "123", params["key4"])
}

func TestParamsBuilder_Overwrite(t *testing.T) {
	// 测试重复设置同一个 key 会覆盖
	params := NewParams().
		Set("key", "value1").
		Set("key", "value2").
		Build()

	assert.Equal(t, 1, len(params), "重复设置key后长度仍为1")
	assert.Equal(t, "value2", params["key"], "重复设置应覆盖为最新值")
}

func TestParamsBuilder_EmptyBuild(t *testing.T) {
	params := NewParams().Build()

	assert.NotNil(t, params, "空Build不应返回nil")
	assert.Equal(t, 0, len(params), "空Build返回的map长度应为0")
}

// 性能测试
func BenchmarkParamsBuilder_Simple(b *testing.B) {
	b.ResetTimer() // 重置计时器，排除初始化耗时
	for i := 0; i < b.N; i++ {
		_ = NewParams().
			Set("key1", "value1").
			Set("key2", "value2").
			Set("key3", "value3").
			Build()
	}
}

func BenchmarkParamsBuilder_WithConditions(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewParams().
			Set("base", "value").
			SetIf(true, "key1", "value1").
			SetIf(false, "key2", "value2").
			SetNotEmpty("key3", "value3").
			SetNotEmpty("key4", "").
			Build()
	}
}

func BenchmarkParamsBuilder_Complex(b *testing.B) {
	additional := map[string]string{
		"extra1": "value1",
		"extra2": "value2",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewParams().
			Set("id", "123").
			SetIf(true, "type", "premium").
			SetNotEmpty("name", "test").
			SetMultiple(additional).
			Build()
	}
}

// 对比传统 map 方式
func BenchmarkTraditionalMap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		params := make(map[string]string)
		params["key1"] = "value1"
		params["key2"] = "value2"
		params["key3"] = "value3"
		_ = params
	}
}
