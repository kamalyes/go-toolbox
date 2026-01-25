/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-25 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\helpers.go
 * @Description: 验证器通用辅助函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"reflect"
	"regexp"
	"sync"
)

// 正则缓存 - 避免重复编译
var (
	regexCache   = make(map[string]*regexp.Regexp)
	regexCacheMu sync.RWMutex
)

// GetCompiledRegex 获取编译的正则（带缓存）- 公共函数供其他模块使用
func GetCompiledRegex(pattern string) (*regexp.Regexp, error) {
	regexCacheMu.RLock()
	re, exists := regexCache[pattern]
	regexCacheMu.RUnlock()

	if exists {
		return re, nil
	}

	regexCacheMu.Lock()
	defer regexCacheMu.Unlock()

	// 双重检查
	if re, exists := regexCache[pattern]; exists {
		return re, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	regexCache[pattern] = re
	return re, nil
}

// ClearRegexCache 清空正则缓存（用于测试或内存释放）
func ClearRegexCache() {
	regexCacheMu.Lock()
	defer regexCacheMu.Unlock()
	regexCache = make(map[string]*regexp.Regexp)
}

// GetReflectKind 获取值的 reflect.Kind（带 nil 处理）
func GetReflectKind(value interface{}) reflect.Kind {
	if value == nil {
		return reflect.Invalid
	}
	return reflect.ValueOf(value).Kind()
}

// IsNumericKind 判断是否为数值类型的 Kind
func IsNumericKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// IsIntegerKind 判断是否为整数类型的 Kind
func IsIntegerKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

// ToFloat64 将数值类型转换为 float64
func ToFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

// IsWholeNumber 判断浮点数是否为整数（没有小数部分）
func IsWholeNumber(f float64) bool {
	return f == float64(int64(f))
}

// StringPtr 创建字符串指针（辅助函数）
func StringPtr(s string) *string {
	return &s
}

// IntPtr 创建整数指针（辅助函数）
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr 创建 float64 指针（辅助函数）
func Float64Ptr(f float64) *float64 {
	return &f
}

// BoolPtr 创建布尔指针（辅助函数）
func BoolPtr(b bool) *bool {
	return &b
}
