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

// deepCopy 递归地复制值
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
		dst.Set(reflect.MakeMap(src.Type())) // 创建新的映射
		for _, key := range src.MapKeys() {  // 遍历源映射的键
			value := src.MapIndex(key) // 获取键对应的值

			// 深拷贝 key（对于复杂类型的key很重要）
			newKey := reflect.New(key.Type()).Elem()
			deepCopy(newKey, key)

			// 深拷贝 value
			newValue := reflect.New(value.Type()).Elem()
			deepCopy(newValue, value)

			dst.SetMapIndex(newKey, newValue) // 设置目标映射的值
		}

	case reflect.Slice: // 处理切片类型
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type())) // 如果切片为nil，设置目标为该类型的零值
			return
		}
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
			dst.Set(src)
			return
		}

		for i := 0; i < src.NumField(); i++ { // 遍历源结构体的字段
			srcField := src.Field(i)             // 获取源字段值
			dstField := dst.Field(i)             // 获取目标字段
			fieldType := src.Type().Field(i)     // 获取字段类型信息
			tag := fieldType.Tag.Get("deepcopy") // 获取字段的deepcopy标签

			// 跳过标记为不复制的字段
			if tag == "-" {
				continue
			}

			// 只复制可设置且导出的字段
			if dstField.CanSet() && fieldType.IsExported() {
				deepCopy(dstField, srcField) // 递归复制字段
			}
		}

	case reflect.Array: // 处理数组类型
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
