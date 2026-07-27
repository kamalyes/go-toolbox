/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-05 15:27:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-05 15:27:15
 * @FilePath: \go-toolbox\pkg\syncx\clone.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"errors"
	"fmt"
	"reflect"
)

// isPrimitiveKind 判断是否为基本类型（不可变，无需深拷贝）
// string/bool/数值类型直接赋值即可，不存在引用共享问题
func isPrimitiveKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	}
	return false
}

// deepCopy 递归地复制值
//
// 性能优化要点：
//  1. Struct 先整体值拷贝（dst.Set(src)，memcpy 级别），只对引用类型字段递归
//     原实现逐字段递归反射，N 个字段 = N 次递归；优化后 1 次 Set + 引用类型字段递归
//  2. Map 的 key 若为基本类型（如 string）直接复用，跳过递归
//  3. Slice 元素若为基本类型，用 reflect.Copy 一次性拷贝，跳过逐元素递归
func deepCopy(dst, src reflect.Value) {
	if !src.IsValid() {
		return // 如果源值无效，直接返回
	}

	switch src.Kind() {
	case reflect.Interface: // 处理接口类型
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type())) // 如果接口为nil，设置目标为该类型的零值
			return
		}
		// 快速路径：基本类型的 interface{}（如 map[string]interface{} 的 string/int 值）
		// 直接赋值，跳过 reflect.New + 递归反射开销
		if isPrimitiveKind(src.Elem().Kind()) {
			dst.Set(src)
			return
		}
		value := src.Elem()                          // 获取接口内部的值
		newValue := reflect.New(value.Type()).Elem() // 创建一个新的值
		deepCopy(newValue, value)                    // 递归复制
		dst.Set(newValue)                            // 设置目标值

	case reflect.Ptr: // 处理指针类型
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type())) // 如果指针为nil，设置目标为该类型的零值
			return
		}
		newPtr := reflect.New(src.Elem().Type()) // 创建一个新的指针
		dst.Set(newPtr)                          // 设置目标为新指针
		deepCopy(newPtr.Elem(), src.Elem())      // 递归复制指针指向的值

	case reflect.Map: // 处理映射类型
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type())) // 如果映射为nil，设置目标为该类型的零值
			return
		}
		dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len())) // 预分配容量，减少扩容
		for _, key := range src.MapKeys() {                     // 遍历源映射的键
			value := src.MapIndex(key) // 获取键对应的值

			// key 快速路径：基本类型直接复用，跳过递归（map[string]X 的 key 是 string 最常见）
			var newKey reflect.Value
			if isPrimitiveKind(key.Kind()) {
				newKey = key
			} else {
				newKey = reflect.New(key.Type()).Elem()
				deepCopy(newKey, key)
			}

			// value 快速路径：基本类型直接 Set，跳过递归
			var newValue reflect.Value
			if isPrimitiveKind(value.Kind()) {
				newValue = value
			} else {
				newValue = reflect.New(value.Type()).Elem()
				deepCopy(newValue, value)
			}

			dst.SetMapIndex(newKey, newValue) // 设置目标映射的值
		}

	case reflect.Slice: // 处理切片类型
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type())) // 如果切片为nil，设置目标为该类型的零值
			return
		}
		// 快速路径：元素为基本类型（如 []string/[]int），用 Copy 一次性拷贝，跳过逐元素递归
		if isPrimitiveKind(src.Type().Elem().Kind()) {
			dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
			reflect.Copy(dst, src)
			return
		}
		// 通用路径：逐元素深拷贝
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap())) // 创建新的切片
		for i := 0; i < src.Len(); i++ {                             // 遍历源切片
			deepCopy(dst.Index(i), src.Index(i)) // 递归复制每个元素
		}

	case reflect.Struct: // 处理结构体类型
		// 特殊处理：如果结构体没有任何导出字段，直接赋值
		// 这包括 time.Time, time.Duration 等标准库类型
		hasExportedField := false
		for i := 0; i < src.NumField(); i++ {
			if src.Type().Field(i).IsExported() {
				hasExportedField = true
				break
			}
		}
		if !hasExportedField {
			dst.Set(src) // 无导出字段，直接值拷贝
			return
		}

		// ⚡ 性能优化：先整体值拷贝（memcpy 级别），值类型字段（string/int/bool/time.Time 等）一次性完成
		// 原实现逐字段递归反射调用，N 个字段 = N 次递归 + N 次 SetMapIndex
		// 优化后：1 次 Set + 仅引用类型字段递归
		dst.Set(src)

		// 只对引用类型字段和含导出字段的 struct 字段递归深拷贝
		// 值类型字段（string/int/bool/Array/Chan/Func 等）已在 dst.Set(src) 中完成
		for i := 0; i < src.NumField(); i++ {
			fieldType := src.Type().Field(i)     // 获取字段类型信息
			tag := fieldType.Tag.Get("deepcopy") // 获取字段的deepcopy标签

			// 跳过标记为不复制的字段
			if tag == "-" {
				continue
			}

			// 只处理可设置且导出的字段
			dstField := dst.Field(i)
			if !dstField.CanSet() || !fieldType.IsExported() {
				continue
			}

			srcField := src.Field(i)
			kind := fieldType.Type.Kind()

			switch kind {
			case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
				// 引用类型：先清零（消除 Set 留下的共享引用），再深拷贝
				// nil 引用已被 dst.Set(src) 正确拷贝为 nil，跳过
				if !srcField.IsNil() {
					dstField.Set(reflect.Zero(dstField.Type()))
					deepCopy(dstField, srcField)
				}
			case reflect.Struct:
				// 含导出字段的 struct 字段：递归处理其引用类型子字段
				// dst.Set(src) 已值拷贝，但内部引用类型子字段可能共享，需要递归
				// （无导出字段的 struct 如 time.Time 已被 Set 完成，递归内部会直接 Set 返回）
				deepCopy(dstField, srcField)
			}
			// 值类型字段（string/int/bool/Array 等）已在 dst.Set(src) 中完成，跳过
		}

	case reflect.Array: // 处理数组类型
		// 快速路径：元素为基本类型，直接 Copy
		if isPrimitiveKind(src.Type().Elem().Kind()) {
			reflect.Copy(dst, src)
			return
		}
		for i := 0; i < src.Len(); i++ { // 遍历源数组
			deepCopy(dst.Index(i), src.Index(i)) // 递归复制每个元素
		}

	case reflect.Chan, reflect.Func: // 处理通道和函数类型
		dst.Set(reflect.Zero(dst.Type())) // 设置目标为该类型的零值

	default: // 处理基本类型
		dst.Set(src) // 直接设置目标值
	}
}

// DeepCopy 复制源值到目标值
//
// @params dst: 目标值的指针，表示要将源值复制到的位置。必须是一个指向某种类型的指针。
// @params src: 源值的指针，表示要复制的原始数据。也必须是一个指向某种类型的指针。
//
// @return:
//
//	如果成功，返回 nil；如果源值为 nil，返回一个错误。
func DeepCopy(dst, src interface{}) error {
	dstVal := reflect.ValueOf(dst) // 获取目标的反射值
	srcVal := reflect.ValueOf(src) // 获取源的反射值

	// 检查目标和源是否都是指针
	if dstVal.Kind() != reflect.Ptr || srcVal.Kind() != reflect.Ptr {
		panic("DeepCopy: both dst and src must be pointers") // 如果不是指针，抛出异常
	}

	// 检查源是否为nil
	if srcVal.IsNil() {
		return errors.New("DeepCopy: src is nil") // 如果源为nil，返回错误
	}

	// 如果目标为nil，则为目标分配新内存
	if dstVal.IsNil() {
		dstVal.Set(reflect.New(srcVal.Elem().Type())) // 为目标分配新内存
	}

	// 检查类型不匹配
	if dstVal.Type() != srcVal.Type() {
		panic(fmt.Sprintf("DeepCopy: type mismatch: %s != %s", dstVal.Type(), srcVal.Type())) // 抛出异常
	}

	// 执行深度复制
	deepCopy(dstVal.Elem(), srcVal.Elem())
	return nil // 返回nil表示成功
}
