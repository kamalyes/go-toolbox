/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-19 10:25:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-19 10:25:55
 * @FilePath: \go-toolbox\pkg\validator\base.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package validator

import (
	"reflect"
	"strings"
	"unicode"
)

// isEmptyValue checks if a reflect.Value is empty.
func IsEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0.0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil() || IsEmptyValue(v.Elem())
	case reflect.Struct:
		isEmpty := true // 默认假设结构体是空的
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Interface && v.Field(i).IsNil() {
				continue // 如果字段是interface{}且为nil，继续检查其他字段
			}
			if !IsEmptyValue(v.Field(i)) {
				isEmpty = false // 只要有一个字段不为空，整体就不为空
				break
			}
		}
		return isEmpty
	default:
		return false
	}
}

// IsCEmpty 判断元素是否为类型零值。
// Params：
//   - v: 需要判断的元素，类型为 T。
//
// Returns:
//   - 返回布尔值，true 表示 v 是类型的零值，false 表示非零值。
func IsCEmpty[T comparable](v T) bool {
	var zero T
	return v == zero
}

// HasEmpty checks if any element in the slice is empty.
func HasEmpty(elems []interface{}) (bool, int) {
	if len(elems) == 0 {
		return true, 0
	}

	emptyCount := 0
	for _, elem := range elems {
		if IsEmptyValue(reflect.ValueOf(elem)) {
			emptyCount++
		}
	}

	return emptyCount > 0, emptyCount
}

// IsAllEmpty checks if all elements in the slice are empty.
func IsAllEmpty(elems []interface{}) bool {
	for _, elem := range elems {
		if !IsEmptyValue(reflect.ValueOf(elem)) {
			return false
		}
	}
	return true
}

// IsUndefined checks if a string is "undefined" (case insensitive).
func IsUndefined(str string) bool {
	return strings.EqualFold(strings.TrimSpace(str), "undefined")
}

// ContainsChinese checks if a string contains any Chinese characters.
func ContainsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

// EmptyToDefault returns defaultStr if str is empty; otherwise, returns str.
func EmptyToDefault(str string, defaultStr string) string {
	if IsEmptyValue(reflect.ValueOf(str)) {
		return defaultStr
	}
	return str
}

// IsNil 判断传入的接口值是否为 nil
// 先判断接口本身是否为 nil，若不是 nil，则通过反射检查其底层值是否为 nil
// 适用于指针、切片、映射、通道、函数和接口类型的 nil 判断
func IsNil(x interface{}) bool {
	if x == nil {
		return true
	}
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

// IsFuncType 判断T是否为函数类型，利用反射
func IsFuncType[T any]() bool {
	var zero T
	tp := reflect.TypeOf(zero)
	if tp == nil {
		return false
	}
	return tp.Kind() == reflect.Func
}
