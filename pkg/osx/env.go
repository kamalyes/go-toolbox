/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-21 03:50:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-21 03:55:26
 * @FilePath: \go-toolbox\pkg\osx\env.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"os"
	"strconv"
)

// Getenv 函数用于获取指定环境变量的值
// 如果环境变量不存在或无法解析，则返回提供的默认值
// T 是一个类型参数，可以是任意类型
func Getenv[T any](key string, defaultValue T) T {
	// 获取环境变量的值
	value := os.Getenv(key)

	// 如果环境变量不存在，直接返回默认值
	if value == "" {
		return defaultValue
	}

	// 使用类型断言判断 defaultValue 的类型
	switch any(defaultValue).(type) {
	case string:
		// 如果是字符串类型，直接返回环境变量的值
		return any(value).(T)
	case int:
		// 如果是整数类型，尝试将字符串解析为整数
		if intValue, err := strconv.Atoi(value); err == nil {
			return any(intValue).(T) // 返回解析后的整数值
		}
	case int32:
		// 如果是 int32 类型，尝试将字符串解析为 int32
		if intValue, err := strconv.ParseInt(value, 10, 32); err == nil {
			return any(int32(intValue)).(T)
		}
	case int64:
		// 如果是 int64 类型，尝试将字符串解析为 int64
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return any(int64(intValue)).(T)
		}
	case uint:
		// 如果是 uint 类型，尝试将字符串解析为 uint
		if uintValue, err := strconv.ParseUint(value, 10, 0); err == nil {
			return any(uint(uintValue)).(T)
		}
	case uint32:
		// 如果是 uint32 类型，尝试将字符串解析为 uint32
		if uintValue, err := strconv.ParseUint(value, 10, 32); err == nil {
			return any(uint32(uintValue)).(T)
		}
	case uint64:
		// 如果是 uint64 类型，尝试将字符串解析为 uint64
		if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			return any(uint64(uintValue)).(T)
		}
	case float32:
		// 如果是 float32 类型，尝试将字符串解析为 float32
		if floatValue, err := strconv.ParseFloat(value, 32); err == nil {
			return any(float32(floatValue)).(T)
		}
	case float64:
		// 如果是浮点数类型，尝试将字符串解析为浮点数
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return any(floatValue).(T) // 返回解析后的浮点数值
		}
	case bool:
		// 如果是布尔类型，尝试将字符串解析为布尔值
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return any(boolValue).(T) // 返回解析后的布尔值
		}
	default:
		// 如果类型不受支持，返回默认值
		return defaultValue
	}

	// 如果解析过程中发生错误，返回默认值
	return defaultValue
}
