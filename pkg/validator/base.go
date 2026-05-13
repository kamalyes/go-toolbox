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

	"github.com/kamalyes/go-toolbox/pkg/types"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
		str := strings.TrimSpace(v.String())
		return str == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0.0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Func:
		// 函数类型：nil 为空
		return v.IsNil()
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

	// 检查 protobuf wrapper 类型
	if isEmpty, ok := CheckEmptyWrapperPointer(v); ok {
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
	if !v.CanInterface() {
		return false, false
	}

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

// isWrapperStringValueEmpty 检查 *wrapperspb.StringValue 是否为空
func isWrapperStringValueEmpty(val *wrapperspb.StringValue) bool {
	str := strings.TrimSpace(val.Value)
	return str == ""
}

// CheckEmptyWrapperPointer 检查 protobuf wrapper 类型是否为空
func CheckEmptyWrapperPointer(v reflect.Value) (isEmpty bool, handled bool) {
	if !v.CanInterface() || !isProtobufWrapperType(v.Type()) {
		return false, false
	}
	// 使用类型断言直接判断
	if val, ok := v.Interface().(*wrapperspb.StringValue); ok {
		if val == nil {
			return true, true
		}
		return isWrapperStringValueEmpty(val), true
	}

	if val, ok := v.Interface().(*wrapperspb.BytesValue); ok {
		return val == nil || len(val.Value) == 0, true
	}

	return v.IsNil(), true
}

// CheckEmptyTimeStruct 检查结构体类型的时间是否为空
func CheckEmptyTimeStruct(v reflect.Value) (isEmpty bool, handled bool) {
	if !v.CanInterface() {
		return false, false
	}

	// 使用类型断言直接判断
	switch val := v.Interface().(type) {
	case time.Time:
		return IsTimeEmpty(&val), true
	case timestamppb.Timestamp:
		return val.GetSeconds() <= 0, true
	default:
		return false, false
	}
}

// IsTimeEmpty 检查 time.Time 是否为空
// 空的定义：零值、Unix 零点或之前的时间
func IsTimeEmpty(t *time.Time) bool {
	if t == nil {
		return true
	}
	return t.IsZero() || t.Unix() <= 0
}

// IsTimeValid 检查时间值是否有效（非nil且非零值）
// 支持 time.Time 和 *time.Time 类型
// 其他类型始终返回 true（视为有效）
// 适用场景: SQL 构建器中时间范围过滤条件的有效性判断
func IsTimeValid(timeVal interface{}) bool {
	if timeVal == nil {
		return false
	}

	// 处理 *time.Time 类型
	if ptr, ok := timeVal.(*time.Time); ok {
		return ptr != nil && !ptr.IsZero() && ptr.After(time.Unix(0, 0))
	}

	// 处理 time.Time 类型
	if t, ok := timeVal.(time.Time); ok {
		return !t.IsZero() && t.After(time.Unix(0, 0))
	}

	// 其他类型认为有效
	return true
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
// Deprecated: 请使用 types.IsCEmpty 替代
func IsCEmpty[T comparable](v T) bool {
	return types.IsCEmpty(v)
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

// IsNull checks if a string is "null" (case insensitive).
func IsNull(str string) bool {
	return strings.EqualFold(strings.TrimSpace(str), "null")
}

// IfNullOrUndefined returns trueVal if str is "null" or "undefined"; otherwise, returns falseVal.
func IfNullOrUndefined(str string) bool {
	return IsNull(str) || IsUndefined(str)
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
// Deprecated: 请使用 types.IsNil 替代
func IsNil(x interface{}) bool {
	return types.IsNil(x)
}

// IsFuncType 判断传入的类型是否为函数类型
// Deprecated: 请使用 types.IsFuncType 替代
func IsFuncType[T any]() bool {
	return types.IsFuncType[T]()
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

// DerefValue 解引用 interface{} 中的指针，返回底层值
// Deprecated: 请使用 types.DerefValue 替代
func DerefValue(value interface{}) (interface{}, bool) {
	return types.DerefValue(value)
}

// protobufWrapperTypes 存储所有 protobuf wrapper 类型的反射类型
var protobufWrapperTypes = map[reflect.Type]bool{
	reflect.TypeOf((*wrapperspb.StringValue)(nil)).Elem(): true,
	reflect.TypeOf((*wrapperspb.Int32Value)(nil)).Elem():  true,
	reflect.TypeOf((*wrapperspb.Int64Value)(nil)).Elem():  true,
	reflect.TypeOf((*wrapperspb.UInt32Value)(nil)).Elem(): true,
	reflect.TypeOf((*wrapperspb.UInt64Value)(nil)).Elem(): true,
	reflect.TypeOf((*wrapperspb.BoolValue)(nil)).Elem():   true,
	reflect.TypeOf((*wrapperspb.FloatValue)(nil)).Elem():  true,
	reflect.TypeOf((*wrapperspb.DoubleValue)(nil)).Elem(): true,
	reflect.TypeOf((*wrapperspb.BytesValue)(nil)).Elem():  true,
}

// isProtobufWrapperType 检查类型是否为 protobuf wrapper 类型
func isProtobufWrapperType(t reflect.Type) bool {
	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return protobufWrapperTypes[t]
}

// invokeGetValue 调用 protobuf wrapper 的 GetValue 方法
func invokeGetValue(v reflect.Value) (interface{}, bool) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, true
		}
		method := v.MethodByName("GetValue")
		if method.IsValid() {
			results := method.Call(nil)
			if len(results) > 0 {
				return results[0].Interface(), true
			}
		}
		return nil, false
	}

	ptr := reflect.New(v.Type())
	ptr.Elem().Set(v)
	method := ptr.MethodByName("GetValue")
	if !method.IsValid() {
		return nil, false
	}
	results := method.Call(nil)
	if len(results) == 0 {
		return nil, false
	}
	return results[0].Interface(), true
}

// UnwrapProtobufWrapper 解包 protobuf 包装器，返回底层值
// 如果值为 nil 或指向 nil 指针，则返回 (nil, true)
// 如果值是指针且非 nil，返回 (dereferencedValue, true)
// 如果值不是指针，返回 (originalValue, true)
func UnwrapProtobufWrapper(value interface{}) (interface{}, bool) {
	if value == nil {
		return nil, false
	}
	rv := reflect.ValueOf(value)
	if !isProtobufWrapperType(rv.Type()) {
		return nil, false
	}
	return invokeGetValue(rv)
}

// isEmptyUnwrappedProtobufWrapper 判断解包后的 protobuf 包装器是否为空
// 支持字符串、字节数组、空值、未定义值、空字符串
// 返回 true 表示为空，false 表示不为空
func isEmptyUnwrappedProtobufWrapper(value interface{}) bool {
	switch val := value.(type) {
	case nil:
		return true
	case string:
		str := strings.TrimSpace(val)
		return str == ""
	case []byte:
		return len(val) == 0
	default:
		return false
	}
}

// tryUnwrapAndCheckEmpty 尝试解包 protobuf wrapper 并检查是否为空
// 如果成功解包，返回 (value, true, isEmpty)
// 如果不是 wrapper 类型，返回 (nil, false, false)
func tryUnwrapAndCheckEmpty(value interface{}) (interface{}, bool, bool) {
	unwrapped, ok := UnwrapProtobufWrapper(value)
	if !ok {
		return nil, false, false
	}
	return unwrapped, true, isEmptyUnwrappedProtobufWrapper(unwrapped)
}

// IsEmptyAfterDeref 判断值是否为空（用于过滤条件场景）
// 支持自动解引用指针：nil 指针为空，非 nil 指针检查其底层值
// 内部复用 IsEmptyValue 进行统一的空值判断
// 返回解引用后的值和是否为空的标志，避免调用方重复调用 DerefValue
//
// 适用场景: SQL 构建器中的 IfNotEmpty 系列方法
func IsEmptyAfterDeref(value interface{}) (interface{}, bool) {
	if unwrapped, handled, empty := tryUnwrapAndCheckEmpty(value); handled {
		if empty {
			return nil, true
		}
		return unwrapped, false
	}

	deref, ok := DerefValue(value)
	if !ok {
		return nil, true
	}

	if unwrapped, handled, empty := tryUnwrapAndCheckEmpty(deref); handled {
		if empty {
			return nil, true
		}
		return unwrapped, false
	}

	// bool 的 false 是有效过滤值（例如 status=false），不应被当作空值跳过。
	if _, isBool := deref.(bool); isBool {
		return deref, false
	}

	if IsEmptyValue(reflect.ValueOf(deref)) {
		return nil, true
	}
	return deref, false
}

// NormalizeFilterValue 归一化过滤条件值。
// 支持 protobuf wrapper 解包，并将切片/数组递归转换为 []interface{}。
func NormalizeFilterValue(value interface{}) interface{} {
	if normalized, ok := UnwrapProtobufWrapper(value); ok {
		return normalized
	}

	if values, ok := value.([]interface{}); ok {
		return NormalizeFilterValueSlice(values)
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return value
	}
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		if rv.Kind() == reflect.Slice && rv.IsNil() {
			return nil
		}
		return normalizeReflectFilterValueSlice(rv)
	}

	return value
}

// NormalizeFilterValueSlice 递归归一化过滤条件值切片。
func NormalizeFilterValueSlice(values []interface{}) []interface{} {
	if values == nil {
		return nil
	}

	normalized := make([]interface{}, len(values))
	for i, value := range values {
		normalized[i] = NormalizeFilterValue(value)
	}
	return normalized
}

// NormalizeFilterValueIfNotEmpty 先按过滤条件规则判断空值，再返回归一化后的值。
func NormalizeFilterValueIfNotEmpty(value interface{}) (interface{}, bool) {
	deref, empty := IsEmptyAfterDeref(value)
	if empty {
		return nil, true
	}
	return NormalizeFilterValue(deref), false
}

// normalizeReflectFilterValueSlice 递归归一化反射值切片。
// 支持 protobuf wrapper 解包，并将切片/数组递归转换为 []interface{}。
func normalizeReflectFilterValueSlice(values reflect.Value) []interface{} {
	normalized := make([]interface{}, values.Len())
	for i := 0; i < values.Len(); i++ {
		normalized[i] = NormalizeFilterValue(values.Index(i).Interface())
	}
	return normalized
}
