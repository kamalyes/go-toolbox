/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-09 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 13:17:27
 * @FilePath: \go-toolbox\pkg\types\reflect.go
 * @Description: 反射工具函数
 */
package types

import (
	"reflect"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ProtoMessageType protobuf 消息类型
var ProtoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

// IsProtoMessageType 判断类型是否实现 proto.Message 接口
func IsProtoMessageType(t reflect.Type) bool {
	return t != nil && t.Implements(ProtoMessageType)
}

// IsExportedField 判断结构体字段是否可按 JSON 规则参与导出处理
func IsExportedField(field reflect.StructField) bool {
	return field.PkgPath == "" || field.Anonymous
}

// ExtractJSONKey 从结构体字段标签中提取 JSON 键名
func ExtractJSONKey(fieldType reflect.StructField) string {
	jsonTag := fieldType.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return ""
	}
	if idx := strings.Index(jsonTag, ","); idx >= 0 {
		jsonTag = jsonTag[:idx]
	}
	if jsonTag == "" {
		return fieldType.Name
	}
	return jsonTag
}

// JSONFieldName 获取结构体字段的 JSON 字段名；没有显式名称时返回 Go 字段名
func JSONFieldName(fieldType reflect.StructField) string {
	if name := ExtractJSONKey(fieldType); name != "" {
		return name
	}
	return fieldType.Name
}

// HasJSONTagOption 判断结构体字段的 json tag 是否包含指定选项
func HasJSONTagOption(fieldType reflect.StructField, options ...string) bool {
	jsonTag := fieldType.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" || len(options) == 0 {
		return false
	}

	idx := strings.Index(jsonTag, ",")
	if idx < 0 || idx == len(jsonTag)-1 {
		return false
	}

	for _, tagOption := range strings.Split(jsonTag[idx+1:], ",") {
		for _, option := range options {
			if tagOption == option {
				return true
			}
		}
	}
	return false
}

// EnsureStructDefaults 确保结构体字段默认值
func EnsureStructDefaults(v reflect.Value) {
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			fieldType := t.Field(i)
			if IsProtoMessageType(fieldType.Type) {
				field.Set(reflect.New(fieldType.Type.Elem()))
			} else if fieldType.Type.Elem().Kind() == reflect.Struct {
				field.Set(reflect.New(fieldType.Type.Elem()))
			}
		}
	}
}

// NewProtoMessage 创建 protobuf 消息实例
func NewProtoMessage[T proto.Message]() T {
	var zero T
	t := reflect.TypeOf(zero).Elem()
	return reflect.New(t).Interface().(T)
}

// IsNil 判断是否为 nil
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

// IsCEmpty 判断是否为空
func IsCEmpty[T comparable](v T) bool {
	var zero T
	return v == zero
}

// DerefValue 解引用指针，返回解引用后的值和是否成功
func DerefValue(value interface{}) (interface{}, bool) {
	if value == nil {
		return nil, false
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, false
		}
		return rv.Elem().Interface(), true
	}
	return value, true
}

// UnwrapPBValue 解包 protobuf wrapper 类型，返回底层值
// 支持 wrapperspb 全部 9 种类型：String, Bool, Int32, Int64, UInt32, UInt64, Float, Double, Bytes
// 如果 wrapper 为 nil 或 typed nil，返回 nil
// 如果不是 wrapper 类型，返回原始值
func UnwrapPBValue(iface interface{}) interface{} {
	switch v := iface.(type) {
	case *wrapperspb.StringValue:
		if !IsNil(v) {
			return v.Value
		}
		return nil
	case *wrapperspb.BoolValue:
		if !IsNil(v) {
			return v.Value
		}
		return nil
	case *wrapperspb.Int32Value:
		if !IsNil(v) {
			return v.Value
		}
		return nil
	case *wrapperspb.Int64Value:
		if !IsNil(v) {
			return v.Value
		}
		return nil
	case *wrapperspb.UInt32Value:
		if !IsNil(v) {
			return v.Value
		}
		return nil
	case *wrapperspb.UInt64Value:
		if !IsNil(v) {
			return v.Value
		}
		return nil
	case *wrapperspb.FloatValue:
		if !IsNil(v) {
			return v.Value
		}
		return nil
	case *wrapperspb.DoubleValue:
		if !IsNil(v) {
			return v.Value
		}
		return nil
	case *wrapperspb.BytesValue:
		if !IsNil(v) {
			return v.Value
		}
		return nil
	default:
		return iface
	}
}

// ResolveModelKey 解析 Model 结构体字段的 key
// 优先使用 gorm column tag，其次 json tag，最后使用字段名
func ResolveModelKey(fieldType reflect.StructField) string {
	if tag, ok := fieldType.Tag.Lookup("gorm"); ok {
		if col := ExtractGormColumn(tag); col != "" {
			return col
		}
	}
	if jsonTag := fieldType.Tag.Get("json"); jsonTag == "-" {
		return "-"
	}
	if name := ExtractJSONKey(fieldType); name != "" {
		return name
	}
	return fieldType.Name
}

// ExtractGormColumn 从 gorm tag 中提取 column 名
func ExtractGormColumn(tag string) string {
	for _, part := range strings.Split(tag, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return ""
}

// ResolvePBKey 解析 PB 结构体字段的 key
// 优先使用 protobuf tag 的 name，其次 json tag，最后使用字段名
func ResolvePBKey(fieldType reflect.StructField) string {
	if tag, ok := fieldType.Tag.Lookup("protobuf"); ok {
		if name := ExtractPBTagValue(tag, "name"); name != "" {
			return name
		}
	}
	if jsonTag := fieldType.Tag.Get("json"); jsonTag == "-" {
		return "-"
	}
	if name := ExtractJSONKey(fieldType); name != "" {
		return name
	}
	return fieldType.Name
}

// ExtractPBTagValue 从 protobuf tag 中提取指定键的值
func ExtractPBTagValue(tag string, key string) string {
	for _, part := range strings.Split(tag, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, key+"=") {
			return strings.TrimPrefix(part, key+"=")
		}
	}
	return ""
}

// UnwrapModelValue 解引用指针类型，返回底层值
// 如果不是指针或指针为 nil，返回原始值
func UnwrapModelValue(iface interface{}) interface{} {
	v := reflect.ValueOf(iface)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		return v.Elem().Interface()
	}
	return iface
}
