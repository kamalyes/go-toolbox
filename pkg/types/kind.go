/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-13 13:08:33
 * @FilePath: \go-toolbox\pkg\types\kind.go
 * @Description: 类型判断辅助函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package types

import "reflect"

// GetReflectKind 获取反射类型的Kind
func GetReflectKind(value interface{}) reflect.Kind {
	if value == nil {
		return reflect.Invalid
	}
	return reflect.ValueOf(value).Kind()
}

// IsNumericKind 判断是否为数值类型
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

// IsIntegerKind 判断是否为整数类型
func IsIntegerKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

// IsFloatKind 判断是否为浮点数类型
func IsFloatKind(kind reflect.Kind) bool {
	return kind == reflect.Float32 || kind == reflect.Float64
}

// ToFloat64OK 尝试将值转换为 float64 类型，返回转换结果和是否成功
func ToFloat64OK(value interface{}) (float64, bool) {
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

// IsWholeNumber 判断是否为整数
func IsWholeNumber(f float64) bool {
	return f == float64(int64(f))
}

// IsFuncType 判断是否为函数类型
func IsFuncType[T any]() bool {
	var zero T
	tp := reflect.TypeOf(zero)
	if tp == nil {
		return false
	}
	return tp.Kind() == reflect.Func
}
