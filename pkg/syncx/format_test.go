/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 13:05:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-13 17:50:20
 * @FilePath: \go-toolbox\pkg\syncx\format_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildContentExtra(t *testing.T) {
	// 测试空数据
	result := BuildContentExtra(map[string]interface{}{})
	assert.Equal(t, "{}", result)

	// 测试非空数据
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	result = BuildContentExtra(data)
	assert.JSONEq(t, `{"key1":"value1","key2":123,"key3":true}`, result)
}

func TestGetStringFromData(t *testing.T) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	assert.Equal(t, "value1", GetStringFromData(data, "key1"))
	assert.Equal(t, "", GetStringFromData(data, "key2")) // 非字符串
	assert.Equal(t, "", GetStringFromData(data, "key3")) // 不存在的键
	assert.Equal(t, "", GetStringFromData(nil, "key1"))  // nil 数据
}

func TestGetBoolFromData(t *testing.T) {
	data := map[string]interface{}{
		"key1": true,
		"key2": "not a bool",
	}

	assert.Equal(t, true, GetBoolFromData(data, "key1"))
	assert.Equal(t, false, GetBoolFromData(data, "key2")) // 非布尔值
	assert.Equal(t, false, GetBoolFromData(data, "key3")) // 不存在的键
	assert.Equal(t, false, GetBoolFromData(nil, "key1"))  // nil 数据
}

func TestGetInt64FromData(t *testing.T) {
	data := map[string]interface{}{
		"key1": int64(12345),
		"key2": 6789,
		"key3": 12.34,
	}

	assert.Equal(t, int64(12345), GetInt64FromData(data, "key1"))
	assert.Equal(t, int64(6789), GetInt64FromData(data, "key2"))
	assert.Equal(t, int64(12), GetInt64FromData(data, "key3")) // 浮点数转换
	assert.Equal(t, int64(0), GetInt64FromData(data, "key4"))  // 不存在的键
	assert.Equal(t, int64(0), GetInt64FromData(nil, "key1"))   // nil 数据
}

func TestParseContentExtraToMap(t *testing.T) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
		"key4": []string{"item1", "item2"},
	}

	result := ParseContentExtraToMap(data)

	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "123", result["key2"])                // 整数转换为字符串
	assert.Equal(t, "true", result["key3"])               // 布尔值转换为字符串
	assert.JSONEq(t, `["item1","item2"]`, result["key4"]) // 非字符串值转换为JSON字符串
}
