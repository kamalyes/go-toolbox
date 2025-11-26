/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-13 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-26 21:08:32
 * @FilePath: \go-toolbox\pkg\safe\safe_access.go
 * @Description: 安全访问装饰器 - 类似JavaScript的可选链操作符
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"fmt"
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/types"
	"reflect"
	"strings"
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

	// 检查字段是否可导出（公开）
	if !field.CanInterface() {
		return &SafeAccess{valid: false}
	}

	// 如果是指针类型且为nil
	if field.Kind() == reflect.Ptr && field.IsNil() {
		return &SafeAccess{valid: false}
	}

	// 如果是指针类型且不为nil，解引用获取实际值
	fieldValue := field.Interface()
	if field.Kind() == reflect.Ptr && !field.IsNil() {
		fieldValue = field.Elem().Interface()
	}

	return &SafeAccess{
		value: fieldValue,
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

	// 使用 convert.MustBool 进行强大的类型转换
	result := convert.MustBool(s.value)

	if !result && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return result
}

// Int 获取整数值，支持强大的类型自动转换
func (s *SafeAccess) Int(defaultValue ...int) int {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	// 使用 convert.MustIntT 进行类型转换
	result, err := convert.MustIntT[int](s.value, nil)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return result
}

// Int64 获取int64值，支持强大的类型自动转换
func (s *SafeAccess) Int64(defaultValue ...int64) int64 {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	result, err := convert.MustIntT[int64](s.value, nil)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return result
}

// Int32 获取int32值
func (s *SafeAccess) Int32(defaultValue ...int32) int32 {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	result, err := convert.MustIntT[int32](s.value, nil)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return result
}

// Uint 获取uint值，支持强大的类型自动转换
func (s *SafeAccess) Uint(defaultValue ...uint) uint {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	result, err := convert.MustIntT[uint](s.value, nil)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return result
}

// Uint64 获取uint64值
func (s *SafeAccess) Uint64(defaultValue ...uint64) uint64 {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	result, err := convert.MustIntT[uint64](s.value, nil)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return result
}

// Float32 获取float32值，支持强大的类型自动转换
func (s *SafeAccess) Float32(defaultValue ...float32) float32 {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	result, err := convert.MustFloatT[float32](s.value, convert.RoundNone)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return result
}

// Float64 获取float64值，支持强大的类型自动转换
func (s *SafeAccess) Float64(defaultValue ...float64) float64 {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	result, err := convert.MustFloatT[float64](s.value, convert.RoundNone)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return result
}

// String 获取字符串值
func (s *SafeAccess) String(defaultValue ...string) string {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}

	// 使用 convert.MustString 进行强大的类型转换
	result := convert.MustString(s.value)

	// 如果转换结果为空且提供了默认值，使用默认值
	if result == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return result
}

// StringOr 获取字符串值，如果无效或为空则返回默认值
func (s *SafeAccess) StringOr(defaultValue string) string {
	if !s.valid {
		return defaultValue
	}

	// 使用 convert.MustString 进行强大的类型转换
	result := convert.MustString(s.value)
	if result == "" {
		return defaultValue
	}
	return result
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
		// 使用 convert.MustString 进行类型转换
		return convert.MustString(v)
	}
	return ""
}

// SafeGetBool 安全获取map中的布尔值
func SafeGetBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		// 使用 convert.MustBool 进行类型转换
		return convert.MustBool(v)
	}
	return false
}

// SafeGetStringSlice 安全获取map中的字符串切片
func SafeGetStringSlice(m map[string]interface{}, key string) []string {
	if v, ok := m[key]; ok {
		// 直接使用[]string类型
		if slice, ok := v.([]string); ok {
			return slice
		}
		// 处理[]interface{}类型
		if slice, ok := v.([]interface{}); ok {
			result := make([]string, 0, len(slice))
			for _, item := range slice {
				// 使用convert.MustString进行类型转换
				result = append(result, convert.MustString(item))
			}
			return result
		}
	}
	return nil
}

// splitFieldPath 分割字段路径字符串
// 例如: "Config.Database.Host" => ["Config", "Database", "Host"]
func splitFieldPath(path string) []string {
	if path == "" {
		return []string{}
	}
	return strings.Split(path, ".")
}

// 通用便捷方法

// At 通用字段路径访问方法，支持链式调用和默认值
func (s *SafeAccess) At(fieldPath string, defaultValue ...interface{}) *SafeAccess {
	// 空路径视为无效
	if fieldPath == "" {
		if len(defaultValue) > 0 {
			return Safe(defaultValue[0])
		}
		return &SafeAccess{valid: false}
	}

	// 支持路径访问，如 "Config.Database.Host"
	fields := splitFieldPath(fieldPath)
	current := s

	for _, field := range fields {
		current = current.Field(field)
		if !current.valid {
			break
		}
	}

	if !current.valid && len(defaultValue) > 0 {
		return Safe(defaultValue[0])
	}

	return current
}

// BoolAt 获取布尔字段的便捷方法
func (s *SafeAccess) BoolAt(fieldPath string, defaultValue ...bool) bool {
	return s.At(fieldPath).Bool(defaultValue...)
}

// StringAt 获取字符串字段的便捷方法
func (s *SafeAccess) StringAt(fieldPath string, defaultValue ...string) string {
	return s.At(fieldPath).String(defaultValue...)
}

// StringOrAt 获取字符串字段，支持空值默认值
func (s *SafeAccess) StringOrAt(fieldPath string, defaultValue string) string {
	return s.At(fieldPath).StringOr(defaultValue)
}

// IntAt 获取整数字段的便捷方法
func (s *SafeAccess) IntAt(fieldPath string, defaultValue ...int) int {
	return s.At(fieldPath).Int(defaultValue...)
}

// DurationAt 获取时间间隔字段的便捷方法
func (s *SafeAccess) DurationAt(fieldPath string, defaultValue ...time.Duration) time.Duration {
	return s.At(fieldPath).Duration(defaultValue...)
}

// ValueAt 获取任意类型字段的原始值
func (s *SafeAccess) ValueAt(fieldPath string, defaultValue ...interface{}) interface{} {
	result := s.At(fieldPath, defaultValue...)
	return result.Value()
}

// GetIntValue 从SafeAccess中获取int值,支持int/int64/float64类型转换
func (sa *SafeAccess) GetIntValue(defaultValue int) int {
	if !sa.valid {
		return defaultValue
	}

	// 使用 convert.MustIntT 进行类型转换
	result, err := convert.MustIntT[int](sa.value, nil)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetInt64Value 从SafeAccess中获取int64值,支持int/int64/float64类型转换
func (sa *SafeAccess) GetInt64Value(defaultValue int64) int64 {
	if !sa.valid {
		return defaultValue
	}

	// 使用 convert.MustIntT 进行类型转换
	result, err := convert.MustIntT[int64](sa.value, nil)
	if err != nil {
		return defaultValue
	}
	return result
}

// ==================== 泛型转换方法 ====================

// As 泛型转换方法 - 支持所有数值类型
func As[T types.Numerical](s *SafeAccess, defaultValue ...T) T {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		var zero T
		return zero
	}

	result, err := convert.MustIntT[T](s.value, nil)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		var zero T
		return zero
	}
	return result
}

// AsFloat 泛型浮点数转换方法
func AsFloat[T types.Float](s *SafeAccess, mode convert.RoundMode, defaultValue ...T) T {
	if !s.valid {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		var zero T
		return zero
	}

	result, err := convert.MustFloatT[T](s.value, mode)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		var zero T
		return zero
	}
	return result
}

// AsString 智能字符串转换
func (s *SafeAccess) AsString(timeLayout ...string) string {
	if !s.valid {
		return ""
	}
	return convert.MustString(s.value, timeLayout...)
}

// AsBool 智能布尔转换
func (s *SafeAccess) AsBool() bool {
	if !s.valid {
		return false
	}
	return convert.MustBool(s.value)
}

// AsJSON 转换为JSON字符串
func (s *SafeAccess) AsJSON(indent bool) (string, error) {
	if !s.valid {
		return "", fmt.Errorf("invalid value")
	}
	var data []byte
	var err error
	if indent {
		data, err = convert.MustJSONIndent(s.value)
	} else {
		data, err = convert.MustJSON(s.value)
	}
	return string(data), err
}

// ==================== 切片操作 ====================

// AsSlice 泛型切片转换
func AsSlice[T types.Numerical](s *SafeAccess) ([]T, error) {
	if !s.valid {
		return nil, fmt.Errorf("invalid value")
	}

	// 如果已经是对应类型的切片
	if slice, ok := s.value.([]T); ok {
		return slice, nil
	}

	// 如果是字符串切片
	if strSlice, ok := s.value.([]string); ok {
		return convert.StringSliceToNumberSlice[T](strSlice, nil)
	}

	// 如果是interface{}切片
	if interfaceSlice, ok := s.value.([]interface{}); ok {
		result := make([]T, 0, len(interfaceSlice))
		for _, v := range interfaceSlice {
			num, err := convert.MustIntT[T](v, nil)
			if err != nil {
				return nil, err
			}
			result = append(result, num)
		}
		return result, nil
	}

	return nil, fmt.Errorf("cannot convert %T to []%T", s.value, *new(T))
}

// AsFloatSlice 泛型浮点数切片转换
func AsFloatSlice[T types.Float](s *SafeAccess, mode convert.RoundMode) ([]T, error) {
	if !s.valid {
		return nil, fmt.Errorf("invalid value")
	}

	// 如果已经是对应类型的切片
	if slice, ok := s.value.([]T); ok {
		return slice, nil
	}

	// 如果是字符串切片
	if strSlice, ok := s.value.([]string); ok {
		return convert.StringSliceToFloatSlice[T](strSlice, mode)
	}

	// 如果是interface{}切片
	if interfaceSlice, ok := s.value.([]interface{}); ok {
		result := make([]T, 0, len(interfaceSlice))
		for _, v := range interfaceSlice {
			num, err := convert.MustFloatT[T](v, mode)
			if err != nil {
				return nil, err
			}
			result = append(result, num)
		}
		return result, nil
	}

	return nil, fmt.Errorf("cannot convert %T to []%T", s.value, *new(T))
}

// AsStringSlice 字符串切片转换
func (s *SafeAccess) AsStringSlice() []string {
	if !s.valid {
		return nil
	}

	// 如果已经是字符串切片
	if slice, ok := s.value.([]string); ok {
		return slice
	}

	// 如果是interface{}切片
	if interfaceSlice, ok := s.value.([]interface{}); ok {
		result := make([]string, 0, len(interfaceSlice))
		for _, v := range interfaceSlice {
			result = append(result, convert.MustString(v))
		}
		return result
	}

	return nil
}

// ==================== 链式操作增强 ====================

// Map 映射转换（支持泛型）
func Map[T any, R any](s *SafeAccess, fn func(T) R) *SafeAccess {
	if !s.valid {
		return s
	}
	if val, ok := s.value.(T); ok {
		return Safe(fn(val))
	}
	return &SafeAccess{valid: false}
}

// FlatMap 扁平化映射
func (s *SafeAccess) FlatMap(fn func(interface{}) *SafeAccess) *SafeAccess {
	if !s.valid {
		return s
	}
	return fn(s.value)
}

// OrDefault 提供默认值
func OrDefault[T any](s *SafeAccess, defaultValue T) T {
	if !s.valid {
		return defaultValue
	}
	if val, ok := s.value.(T); ok {
		return val
	}
	return defaultValue
}

// Must 强制获取值，无效时panic
func Must[T any](s *SafeAccess) T {
	if !s.valid {
		panic("SafeAccess: value is invalid")
	}
	if val, ok := s.value.(T); ok {
		return val
	}
	var zero T
	panic(fmt.Sprintf("SafeAccess: cannot convert %T to %T", s.value, zero))
}

// ==================== 条件操作 ====================

// When 条件执行
func (s *SafeAccess) When(predicate func(interface{}) bool, fn func(interface{}) interface{}) *SafeAccess {
	if !s.valid {
		return s
	}
	if predicate(s.value) {
		return Safe(fn(s.value))
	}
	return s
}

// Unless 条件排除
func (s *SafeAccess) Unless(predicate func(interface{}) bool, fn func(interface{}) interface{}) *SafeAccess {
	if !s.valid {
		return s
	}
	if !predicate(s.value) {
		return Safe(fn(s.value))
	}
	return s
}

// IsEmpty 检查是否为空值
func (s *SafeAccess) IsEmpty() bool {
	if !s.valid {
		return true
	}

	switch v := s.value.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}

// IsNonEmpty 检查是否非空
func (s *SafeAccess) IsNonEmpty() bool {
	return !s.IsEmpty()
}

// ==================== 类型检查 ====================

// IsType 检查值是否为指定类型
func IsType[T any](s *SafeAccess) bool {
	if !s.valid {
		return false
	}
	_, ok := s.value.(T)
	return ok
}

// IsNumber 检查是否为数值类型
func (s *SafeAccess) IsNumber() bool {
	if !s.valid {
		return false
	}
	switch s.value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}

// IsString 检查是否为字符串
func (s *SafeAccess) IsString() bool {
	if !s.valid {
		return false
	}
	_, ok := s.value.(string)
	return ok
}

// IsBool 检查是否为布尔值
func (s *SafeAccess) IsBool() bool {
	if !s.valid {
		return false
	}
	_, ok := s.value.(bool)
	return ok
}

// IsSlice 检查是否为切片
func (s *SafeAccess) IsSlice() bool {
	if !s.valid {
		return false
	}
	val := reflect.ValueOf(s.value)
	return val.Kind() == reflect.Slice
}

// IsMap 检查是否为map
func (s *SafeAccess) IsMap() bool {
	if !s.valid {
		return false
	}
	val := reflect.ValueOf(s.value)
	return val.Kind() == reflect.Map
}

// ==================== 集合操作 ====================

// Len 获取长度（字符串、切片、map等）
func (s *SafeAccess) Len() int {
	if !s.valid {
		return 0
	}

	switch v := s.value.(type) {
	case string:
		return len(v)
	case []interface{}:
		return len(v)
	case []string:
		return len(v)
	case map[string]interface{}:
		return len(v)
	default:
		val := reflect.ValueOf(s.value)
		switch val.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
			return val.Len()
		default:
			return 0
		}
	}
}

// Keys 获取map的所有键
func (s *SafeAccess) Keys() []string {
	if !s.valid {
		return nil
	}

	if m, ok := s.value.(map[string]interface{}); ok {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		return keys
	}

	return nil
}

// Values 获取map的所有值
func (s *SafeAccess) Values() []interface{} {
	if !s.valid {
		return nil
	}

	if m, ok := s.value.(map[string]interface{}); ok {
		values := make([]interface{}, 0, len(m))
		for _, v := range m {
			values = append(values, v)
		}
		return values
	}

	return nil
}

// Contains 检查切片或map是否包含指定值
func (s *SafeAccess) Contains(target interface{}) bool {
	if !s.valid {
		return false
	}

	// 检查map
	if m, ok := s.value.(map[string]interface{}); ok {
		if key, ok := target.(string); ok {
			_, exists := m[key]
			return exists
		}
		return false
	}

	// 检查切片
	val := reflect.ValueOf(s.value)
	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(val.Index(i).Interface(), target) {
				return true
			}
		}
	}

	return false
}
