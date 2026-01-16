/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-16 13:35:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-16 13:16:15
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

	assert.Equal(t, 2, len(params), "Set方法后params长度应为2")
	assert.Equal(t, "value1", params["key1"], "key1 取值错误")
	assert.Equal(t, "value2", params["key2"], "key2 取值错误")
}

func TestParamsBuilder_SetIf(t *testing.T) {
	t.Run("condition true", func(t *testing.T) {
		params := NewParams().
			Set("base", "value").
			SetIf(true, "key1", "value1").
			SetIf(false, "key2", "value2").
			Build()

		assert.Equal(t, 2, len(params), "条件为true时应包含base和key1")
		assert.Equal(t, "value", params["base"], "base 取值错误")
		assert.Equal(t, "value1", params["key1"], "key1 取值错误")

		// 关键修复：检查不存在的key，而非取值
		_, key2Exists := params["key2"]
		assert.False(t, key2Exists, "条件为false时key2不应存在")
	})

	t.Run("condition false", func(t *testing.T) {
		params := NewParams().
			SetIf(false, "key", "value").
			Build()

		assert.Equal(t, 0, len(params), "条件为false时params长度应为0")
	})

	t.Run("dynamic condition", func(t *testing.T) {
		age := 25
		name := ""

		params := NewParams().
			SetIf(age > 18, "age", "adult").
			SetIf(name != "", "name", name).
			Build()

		assert.Equal(t, 1, len(params), "动态条件应只包含age")
		assert.Equal(t, "adult", params["age"], "age 取值错误")

		// 关键修复：检查不存在的key
		_, nameExists := params["name"]
		assert.False(t, nameExists, "name为空时不应存在")
	})
}

func TestParamsBuilder_SetNotEmpty(t *testing.T) {
	t.Run("non-empty string", func(t *testing.T) {
		params := NewParams().
			SetNotEmpty("key1", "value1").
			SetNotEmpty("key2", "").
			SetNotEmpty("key3", "   "). // 空白字符串会被validator.IsEmptyValue判定为空
			Build()

		assert.Equal(t, 1, len(params), "仅非空字符串key1应存在")
		assert.Equal(t, "value1", params["key1"], "key1 取值错误")

		// 验证空字符串和空白字符串的key不存在
		_, key2Exists := params["key2"]
		assert.False(t, key2Exists, "空字符串key2不应存在")

		_, key3Exists := params["key3"]
		assert.False(t, key3Exists, "空白字符串key3不应存在")
	})

	t.Run("empty string not set", func(t *testing.T) {
		params := NewParams().
			SetNotEmpty("key", "").
			Build()

		_, exists := params["key"]
		assert.False(t, exists, "空字符串key不应存在")
	})

	t.Run("mixed empty and non-empty", func(t *testing.T) {
		params := NewParams().
			Set("base", "always").
			SetNotEmpty("optional1", "has_value").
			SetNotEmpty("optional2", "").
			SetNotEmpty("optional3", "another_value").
			Build()

		assert.Equal(t, 3, len(params), "应包含base、optional1、optional3")
		assert.Equal(t, "always", params["base"], "base 取值错误")
		assert.Equal(t, "has_value", params["optional1"], "optional1 取值错误")
		assert.Equal(t, "another_value", params["optional3"], "optional3 取值错误")

		_, optional2Exists := params["optional2"]
		assert.False(t, optional2Exists, "空字符串optional2不应存在")
	})
}

func TestParamsBuilder_SetMultiple(t *testing.T) {
	additional := map[string]string{
		"key3": "value3",
		"key4": "value4",
	}

	params := NewParams().
		Set("key1", "value1").
		Set("key2", "value2").
		SetMultiple(additional).
		Build()

	assert.Equal(t, 4, len(params), "批量设置后应包含4个参数")
	assert.Equal(t, "value1", params["key1"], "key1 取值错误")
	assert.Equal(t, "value2", params["key2"], "key2 取值错误")
	assert.Equal(t, "value3", params["key3"], "key3 取值错误")
	assert.Equal(t, "value4", params["key4"], "key4 取值错误")
}

func TestParamsBuilder_ChainedCalls(t *testing.T) {
	// 测试复杂的链式调用
	isVIP := true
	discount := 0
	coupon := "SAVE20"

	params := NewParams().
		Set("user_id", "12345").
		Set("product", "laptop").
		SetIf(isVIP, "vip_level", "gold").
		SetIf(discount > 0, "discount", "10").
		SetNotEmpty("coupon", coupon).
		SetMultiple(map[string]string{
			"currency": "USD",
			"quantity": "1",
		}).
		Build()

	assert.Equal(t, 6, len(params), "链式调用后应包含6个参数")
	assert.Equal(t, "12345", params["user_id"], "user_id 取值错误")
	assert.Equal(t, "laptop", params["product"], "product 取值错误")
	assert.Equal(t, "gold", params["vip_level"], "vip_level 取值错误")
	assert.Equal(t, "SAVE20", params["coupon"], "coupon 取值错误")
	assert.Equal(t, "USD", params["currency"], "currency 取值错误")
	assert.Equal(t, "1", params["quantity"], "quantity 取值错误")

	// 关键修复：检查discount是否存在
	_, hasDiscount := params["discount"]
	assert.False(t, hasDiscount, "discount=0时不应存在该key")
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
