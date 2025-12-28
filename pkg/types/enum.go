/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-28 00:00:00 00:00:00
 * @FilePath: \go-toolbox\pkg\types\enum.go
 * @Description: 枚举类型验证器
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package types

import (
	"fmt"
	"sort"
)

// EnumValidator 枚举验证器,用于验证枚举类型的有效性
// T 必须是可比较的类型 (comparable),通常是 string、int 等
//
// 使用示例:
//
//	type UserRole string
//	const (
//	    RoleAdmin    UserRole = "admin"
//	    RoleUser     UserRole = "user"
//	    RoleGuest    UserRole = "guest"
//	)
//
//	validator := NewEnumValidator(RoleAdmin, RoleUser, RoleGuest)
//	if validator.IsValid(RoleAdmin) {
//	    fmt.Println("Valid role")
//	}
type EnumValidator[T comparable] struct {
	validValues map[T]struct{} // 使用 struct{} 减少内存占用
}

// NewEnumValidator 创建一个新的枚举验证器
// 参数 values 为所有有效的枚举值
//
// 示例:
//
//	validator := NewEnumValidator("pending", "approved", "rejected")
func NewEnumValidator[T comparable](values ...T) *EnumValidator[T] {
	validator := &EnumValidator[T]{
		validValues: make(map[T]struct{}, len(values)),
	}
	for _, v := range values {
		validator.validValues[v] = struct{}{}
	}
	return validator
}

// IsValid 检查给定值是否为有效的枚举值
//
// 示例:
//
//	if validator.IsValid("pending") {
//	    fmt.Println("Valid status")
//	}
func (v *EnumValidator[T]) IsValid(value T) bool {
	_, exists := v.validValues[value]
	return exists
}

// MustBeValid 验证给定值,如果无效则返回错误
//
// 示例:
//
//	if err := validator.MustBeValid(status); err != nil {
//	    return err
//	}
func (v *EnumValidator[T]) MustBeValid(value T) error {
	if !v.IsValid(value) {
		return fmt.Errorf("invalid enum value: %v, expected one of: %v", value, v.GetValidValues())
	}
	return nil
}

// GetValidValues 获取所有有效的枚举值列表
// 返回的列表顺序不固定
//
// 示例:
//
//	validValues := validator.GetValidValues()
//	fmt.Printf("Valid values: %v\n", validValues)
func (v *EnumValidator[T]) GetValidValues() []T {
	values := make([]T, 0, len(v.validValues))
	for value := range v.validValues {
		values = append(values, value)
	}
	return values
}

// GetValidValuesString 获取所有有效枚举值的字符串表示
// 对于 string 类型的枚举特别有用,返回排序后的列表
//
// 示例:
//
//	fmt.Printf("Valid roles: %v\n", validator.GetValidValuesString())
func (v *EnumValidator[T]) GetValidValuesString() []string {
	values := v.GetValidValues()
	strValues := make([]string, len(values))
	for i, val := range values {
		strValues[i] = fmt.Sprintf("%v", val)
	}
	sort.Strings(strValues)
	return strValues
}

// Count 返回有效枚举值的数量
func (v *EnumValidator[T]) Count() int {
	return len(v.validValues)
}

// Contains 检查验证器是否包含指定的枚举值
// 功能与 IsValid 相同,提供更语义化的方法名
func (v *EnumValidator[T]) Contains(value T) bool {
	return v.IsValid(value)
}

// Add 添加新的有效枚举值
// 如果值已存在,则不会重复添加
//
// 示例:
//
//	validator.Add("new_status")
func (v *EnumValidator[T]) Add(values ...T) {
	for _, value := range values {
		v.validValues[value] = struct{}{}
	}
}

// Remove 移除指定的枚举值
// 如果值不存在,不会产生错误
//
// 示例:
//
//	validator.Remove("deprecated_status")
func (v *EnumValidator[T]) Remove(values ...T) {
	for _, value := range values {
		delete(v.validValues, value)
	}
}

// Clear 清空所有有效枚举值
func (v *EnumValidator[T]) Clear() {
	v.validValues = make(map[T]struct{})
}

// Clone 创建验证器的副本
func (v *EnumValidator[T]) Clone() *EnumValidator[T] {
	clone := &EnumValidator[T]{
		validValues: make(map[T]struct{}, len(v.validValues)),
	}
	for value := range v.validValues {
		clone.validValues[value] = struct{}{}
	}
	return clone
}
