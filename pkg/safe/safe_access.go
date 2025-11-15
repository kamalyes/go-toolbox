/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-13 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-15 16:31:15
 * @FilePath: \engine-im-service\go-toolbox\pkg\safe\safe_access.go
 * @Description: 安全访问装饰器 - 类似JavaScript的可选链操作符
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"reflect"
	"time"
)

type SafeAccess struct {
	value interface{}
	valid bool
}

// Safe 创建安全访问装饰器
func Safe(v interface{}) *SafeAccess {
	return &SafeAccess{
		value: v,
		valid: v != nil,
	}
}

// Field 安全访问字段，支持链式调用
func (s *SafeAccess) Field(fieldName string) *SafeAccess {
	if !s.valid || s.value == nil {
		return &SafeAccess{valid: false}
	}

	// 首先检查是否是 map[string]interface{}
	if m, ok := s.value.(map[string]interface{}); ok {
		if value, exists := m[fieldName]; exists && value != nil {
			return &SafeAccess{
				value: value,
				valid: true,
			}
		}
		return &SafeAccess{valid: false}
	}

	rv := reflect.ValueOf(s.value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return &SafeAccess{valid: false}
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return &SafeAccess{valid: false}
	}

	field := rv.FieldByName(fieldName)
	if !field.IsValid() {
		return &SafeAccess{valid: false}
	}

	// 如果是指针类型且为nil
	if field.Kind() == reflect.Ptr && field.IsNil() {
		return &SafeAccess{valid: false}
	}

	return &SafeAccess{
		value: field.Interface(),
		valid: true,
	}
}

// Bool 获取布尔值，如果无效则返回默认值
func (s *SafeAccess) Bool(defaultValue ...bool) bool {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	if b, ok := s.value.(bool); ok {
		return b
	}
	if bp, ok := s.value.(*bool); ok && bp != nil {
		return *bp
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// Int 获取整数值
func (s *SafeAccess) Int(defaultValue ...int) int {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	if i, ok := s.value.(int); ok {
		return i
	}
	if ip, ok := s.value.(*int); ok && ip != nil {
		return *ip
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// String 获取字符串值
func (s *SafeAccess) String(defaultValue ...string) string {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}

	if s, ok := s.value.(string); ok {
		return s
	}
	if sp, ok := s.value.(*string); ok && sp != nil {
		return *sp
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// StringOr 获取字符串值，如果无效或为空则返回默认值
func (s *SafeAccess) StringOr(defaultValue string) string {
	if !s.valid {
		return defaultValue
	}

	if str, ok := s.value.(string); ok {
		if str == "" {
			return defaultValue
		}
		return str
	}
	if sp, ok := s.value.(*string); ok && sp != nil {
		if *sp == "" {
			return defaultValue
		}
		return *sp
	}

	return defaultValue
}

// Duration 获取时间间隔值
func (s *SafeAccess) Duration(defaultValue ...time.Duration) time.Duration {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	if d, ok := s.value.(time.Duration); ok {
		return d
	}
	if dp, ok := s.value.(*time.Duration); ok && dp != nil {
		return *dp
	}
	if str, ok := s.value.(string); ok {
		if parsed, err := time.ParseDuration(str); err == nil {
			return parsed
		}
	}
	if i, ok := s.value.(int); ok {
		return time.Duration(i)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// Value 获取原始值
func (s *SafeAccess) Value() interface{} {
	if !s.valid {
		return nil
	}
	return s.value
}

// IsValid 检查值是否有效
func (s *SafeAccess) IsValid() bool {
	return s.valid
}

// OrElse 如果当前值无效，返回备用值
func (s *SafeAccess) OrElse(alternative interface{}) *SafeAccess {
	if s.valid {
		return s
	}
	return Safe(alternative)
}

// IfPresent 如果值存在则执行函数
func (s *SafeAccess) IfPresent(fn func(interface{})) *SafeAccess {
	if s.valid && s.value != nil {
		fn(s.value)
	}
	return s
}

// Map 转换值（类似JavaScript的map）
func (s *SafeAccess) Map(fn func(interface{}) interface{}) *SafeAccess {
	if !s.valid {
		return s
	}
	return Safe(fn(s.value))
}

// Filter 过滤值（类似JavaScript的filter）
func (s *SafeAccess) Filter(predicate func(interface{}) bool) *SafeAccess {
	if !s.valid || !predicate(s.value) {
		return &SafeAccess{valid: false}
	}
	return s
}

// SafeGetString 安全获取map中的字符串值
func SafeGetString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return ""
}

// SafeGetBool 安全获取map中的布尔值
func SafeGetBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// SafeGetStringSlice 安全获取map中的字符串切片
func SafeGetStringSlice(m map[string]interface{}, key string) []string {
	if v, ok := m[key]; ok {
		if slice, ok := v.([]interface{}); ok {
			result := make([]string, 0, len(slice))
			for _, item := range slice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return nil
}

// 常用便捷方法

// Enabled 获取 Enabled 字段的布尔值
func (s *SafeAccess) Enabled(defaultValue ...bool) bool {
	return s.Field("Enabled").Bool(defaultValue...)
}

// Host 获取 Host 字段的字符串值
func (s *SafeAccess) Host(defaultValue ...string) string {
	return s.Field("Host").String(defaultValue...)
}

// Port 获取 Port 字段的整数值
func (s *SafeAccess) Port(defaultValue ...int) int {
	return s.Field("Port").Int(defaultValue...)
}

// Timeout 获取 Timeout 字段的时间间隔值
func (s *SafeAccess) Timeout(defaultValue ...time.Duration) time.Duration {
	return s.Field("Timeout").Duration(defaultValue...)
}
