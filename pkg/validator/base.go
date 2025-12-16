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
	"time"
	"unicode"

	"google.golang.org/protobuf/types/known/timestamppb"
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
		return IsEmptyPointer(v)
	case reflect.Struct:
		return IsEmptyStruct(v)
	default:
		return false
	}
}

// IsEmptyPointer 检查指针或接口类型是否为空
func IsEmptyPointer(v reflect.Value) bool {
	if v.IsNil() {
		return true
	}

	// 检查特殊类型
	if isEmpty, ok := CheckEmptyTimePointer(v); ok {
		return isEmpty
	}

	// 递归检查指针指向的值
	return IsEmptyValue(v.Elem())
}

// IsEmptyStruct 检查结构体是否为空
func IsEmptyStruct(v reflect.Value) bool {
	// 检查特殊的时间类型
	if isEmpty, ok := CheckEmptyTimeStruct(v); ok {
		return isEmpty
	}

	// 通用结构体检查 - 所有字段都为空
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		// 跳过 nil 接口字段
		if field.Kind() == reflect.Interface && field.IsNil() {
			continue
		}
		// 只要有一个字段不为空，整体就不为空
		if !IsEmptyValue(field) {
			return false
		}
	}
	return true
}

// CheckEmptyTimePointer 检查指针类型的时间是否为空
func CheckEmptyTimePointer(v reflect.Value) (isEmpty bool, handled bool) {
	// 检查 *time.Time
	if t, ok := v.Interface().(*time.Time); ok {
		return t == nil || IsTimeEmpty(t), true
	}

	// 检查 *timestamppb.Timestamp
	if ts, ok := v.Interface().(*timestamppb.Timestamp); ok {
		return ts == nil || ts.GetSeconds() <= 0, true
	}

	return false, false
}

// CheckEmptyTimeStruct 检查结构体类型的时间是否为空
func CheckEmptyTimeStruct(v reflect.Value) (isEmpty bool, handled bool) {
	// 检查 time.Time
	if t, ok := v.Interface().(time.Time); ok {
		return IsTimeEmpty(&t), true
	}

	// 检查 protobuf Timestamp（避免复制锁）
	typeName := v.Type().String()
	if typeName == "timestamppb.Timestamp" {
		return IsProtobufTimestampEmpty(v), true
	}

	return false, false
}

// IsTimeEmpty 检查 time.Time 是否为空
// 空的定义：零值、Unix 零点或之前的时间
func IsTimeEmpty(t *time.Time) bool {
	if t == nil {
		return true
	}
	return t.IsZero() || t.Unix() <= 0
}

// IsProtobufTimestampEmpty 检查 protobuf Timestamp 是否为空（使用反射避免复制锁）
func IsProtobufTimestampEmpty(v reflect.Value) bool {
	// 直接通过反射获取 Seconds 字段的值（避免调用方法可能的指针接收者问题）
	secondsField := v.FieldByName("Seconds")
	if !secondsField.IsValid() {
		return true
	}

	// 检查 Seconds 是否 <= 0
	return secondsField.Int() <= 0
}

// IsCEmpty 判断元素是否为类型零值
// Params：
//   - v: 需要判断的元素，类型为 T
//
// Returns:
//   - 返回布尔值，true 表示 v 是类型的零值，false 表示非零值
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

// IsSafeFieldName 检查字段名是否安全(仅包含字母、数字、下划线、点号)
// 用于防止 SQL 注入等安全问题
func IsSafeFieldName(field string) bool {
	if field == "" {
		return false
	}
	for _, ch := range field {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '_' || ch == '.') {
			return false
		}
	}
	return true
}

// IsAllowedField 检查字段是否在白名单中
// 如果提供了白名单，检查字段是否在白名单中；否则验证字段名是否安全
// 参数:
//   - field: 要检查的字段名
//   - allowedFields: 可选的白名单切片（可变参数，传入一个[]string切片）
//
// 返回:
//   - true: 字段允许使用
//   - false: 字段不允许使用
func IsAllowedField(field string, allowedFields ...[]string) bool {
	// 如果提供了白名单，检查字段是否在白名单中
	if len(allowedFields) > 0 && len(allowedFields[0]) > 0 {
		for _, allowedField := range allowedFields[0] {
			if field == allowedField {
				return true
			}
		}
		return false
	}
	// 没有白名单，验证字段名是否安全
	return IsSafeFieldName(field)
}
