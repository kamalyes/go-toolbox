/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-21 03:56:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-21 03:57:15
 * @FilePath: \go-toolbox\tests\osx_env_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"os"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/stretchr/testify/assert"
)

// TestGetenv 测试 Getenv 函数的功能
func TestGetenv(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("TEST_STRING", "hello")
	os.Setenv("TEST_INT", "123")
	os.Setenv("TEST_FLOAT", "123.45")
	os.Setenv("TEST_BOOL", "true")
	os.Setenv("TEST_EMPTY_STRING", "")
	os.Setenv("TEST_INVALID_INT", "not_a_number")
	os.Setenv("TEST_INVALID_FLOAT", "not_a_float")
	os.Setenv("TEST_INVALID_BOOL", "not_a_bool")
	os.Setenv("TEST_NEGATIVE_INT", "-10")
	os.Setenv("TEST_ZERO_FLOAT", "0.0")
	os.Setenv("TEST_UINT", "456")
	os.Setenv("TEST_UINT64", "789")
	os.Setenv("TEST_FLOAT32", "123.456")
	os.Setenv("TEST_FLOAT64", "654.321")

	// 创建一个断言对象
	assert := assert.New(t)

	// 测试字符串类型
	assert.Equal("hello", osx.Getenv("TEST_STRING", "default"), "应该返回环境变量的值")

	// 测试整数类型
	assert.Equal(123, osx.Getenv("TEST_INT", 0), "应该返回解析后的整数值")

	// 测试无符号整数类型
	assert.Equal(uint(456), osx.Getenv("TEST_UINT", uint(0)), "应该返回解析后的无符号整数值")

	// 测试浮点数类型
	assert.Equal(123.45, osx.Getenv("TEST_FLOAT", 0.0), "应该返回解析后的浮点数值")

	// 测试浮点数类型 float32
	assert.Equal(float32(123.456), osx.Getenv("TEST_FLOAT32", float32(0)), "应该返回解析后的 float32 值")

	// 测试浮点数类型 float64
	assert.Equal(654.321, osx.Getenv("TEST_FLOAT64", 0.0), "应该返回解析后的 float64 值")

	// 测试布尔类型
	assert.Equal(true, osx.Getenv("TEST_BOOL", false), "应该返回解析后的布尔值")

	// 测试空字符串，应该返回默认值
	assert.Equal("default", osx.Getenv("TEST_EMPTY_STRING", "default"), "应该返回空字符串")

	// 测试不存在的环境变量，应该返回默认值
	assert.Equal("default", osx.Getenv("NON_EXISTENT", "default"), "应该返回默认值")

	// 测试解析错误的整数值
	assert.Equal(0, osx.Getenv("TEST_INVALID_INT", 0), "应该返回默认值，因为解析失败")

	// 测试解析错误的浮点数值
	assert.Equal(0.0, osx.Getenv("TEST_INVALID_FLOAT", 0.0), "应该返回默认值，因为解析失败")

	// 测试解析错误的布尔值
	assert.Equal(false, osx.Getenv("TEST_INVALID_BOOL", false), "应该返回默认值，因为解析失败")

	// 测试负数整数
	assert.Equal(-10, osx.Getenv("TEST_NEGATIVE_INT", 0), "应该返回解析后的负整数值")

	// 测试零浮点数
	assert.Equal(0.0, osx.Getenv("TEST_ZERO_FLOAT", 1.0), "应该返回解析后的零浮点数值")

	// 测试无符号整数类型 uint64
	assert.Equal(uint64(789), osx.Getenv("TEST_UINT64", uint64(0)), "应该返回解析后的 uint64 值")

	// 清理测试环境变量
	os.Unsetenv("TEST_STRING")
	os.Unsetenv("TEST_INT")
	os.Unsetenv("TEST_FLOAT")
	os.Unsetenv("TEST_BOOL")
	os.Unsetenv("TEST_EMPTY_STRING")
	os.Unsetenv("TEST_INVALID_INT")
	os.Unsetenv("TEST_INVALID_FLOAT")
	os.Unsetenv("TEST_INVALID_BOOL")
	os.Unsetenv("TEST_NEGATIVE_INT")
	os.Unsetenv("TEST_ZERO_FLOAT")
	os.Unsetenv("TEST_UINT")
	os.Unsetenv("TEST_UINT64")
	os.Unsetenv("TEST_FLOAT32")
	os.Unsetenv("TEST_FLOAT64")
}
