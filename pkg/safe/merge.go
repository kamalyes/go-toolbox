/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-05 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-05 00:00:00
 * @FilePath: \go-toolbox\pkg\safe\merge.go
 * @Description: 泛型配置合并工具，支持递归合并结构体字段
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import "reflect"

// MergeWithDefaults 合并配置,用默认配置填充nil或零值字段
// 支持泛型,会递归合并所有嵌套结构体、指针、切片、Map等
func MergeWithDefaults[T any](st *T, defaultSts ...*T) *T {
	if st == nil {
		if len(defaultSts) > 0 && defaultSts[0] != nil {
			return defaultSts[0]
		}
		return nil
	}

	if len(defaultSts) == 0 {
		return st
	}

	result := reflect.New(reflect.TypeOf(st).Elem()).Elem()
	result.Set(reflect.ValueOf(st).Elem())

	for _, defaultConfig := range defaultSts {
		if defaultConfig != nil {
			mergeStruct(result, reflect.ValueOf(defaultConfig).Elem())
		}
	}

	return result.Addr().Interface().(*T)
}

// mergeStruct 递归合并结构体字段
func mergeStruct(target, source reflect.Value) {
	for i := 0; i < target.NumField(); i++ {
		field := target.Field(i)
		sourceField := source.Field(i)

		if field.CanSet() {
			mergeField(field, sourceField)
		}
	}
}

// mergeField 合并单个字段
func mergeField(field, defaultField reflect.Value) {
	if !defaultField.IsValid() {
		return
	}

	switch field.Kind() {
	case reflect.Ptr:
		mergePtr(field, defaultField)

	case reflect.Struct:
		mergeStruct(field, defaultField)

	case reflect.Slice:
		mergeSlice(field, defaultField)

	case reflect.Map:
		mergeMap(field, defaultField)

	case reflect.String:
		mergeString(field, defaultField)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		mergeInt(field, defaultField)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		mergeUint(field, defaultField)

	case reflect.Float32, reflect.Float64:
		mergeFloat(field, defaultField)

	case reflect.Bool:
		mergeBool(field, defaultField)
	}
}

// mergePtr 合并指针类型
func mergePtr(field, defaultField reflect.Value) {
	if field.IsNil() {
		if !defaultField.IsNil() {
			field.Set(defaultField)
		}
	} else if !defaultField.IsNil() {
		mergeField(field.Elem(), defaultField.Elem())
	}
}

// mergeSlice 合并切片类型
func mergeSlice(field, defaultField reflect.Value) {
	if field.IsNil() || field.Len() == 0 {
		if !defaultField.IsNil() && defaultField.Len() > 0 {
			field.Set(defaultField)
		}
	}
}

// mergeString 合并字符串类型
func mergeString(field, defaultField reflect.Value) {
	if field.String() == "" && defaultField.String() != "" {
		field.SetString(defaultField.String())
	}
}

// mergeInt 合并整数类型
func mergeInt(field, defaultField reflect.Value) {
	if field.Int() == 0 && defaultField.Int() != 0 {
		field.SetInt(defaultField.Int())
	}
}

// mergeUint 合并无符号整数类型
func mergeUint(field, defaultField reflect.Value) {
	if field.Uint() == 0 && defaultField.Uint() != 0 {
		field.SetUint(defaultField.Uint())
	}
}

// mergeFloat 合并浮点数类型
func mergeFloat(field, defaultField reflect.Value) {
	if field.Float() == 0 && defaultField.Float() != 0 {
		field.SetFloat(defaultField.Float())
	}
}

// mergeBool 合并布尔类型
func mergeBool(field, defaultField reflect.Value) {
	if !field.Bool() && defaultField.Bool() {
		field.SetBool(defaultField.Bool())
	}
}

// mergeMap 合并map键值对
func mergeMap(target, source reflect.Value) {
	// 如果 target 为 nil,直接使用 source
	if target.IsNil() {
		if !source.IsNil() {
			target.Set(source)
		}
		return
	}

	// target 不为 nil,source 也不为 nil,合并键值对
	if !source.IsNil() {
		for _, key := range source.MapKeys() {
			if !target.MapIndex(key).IsValid() {
				target.SetMapIndex(key, source.MapIndex(key))
			}
		}
	}
}
