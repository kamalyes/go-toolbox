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
)

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
			// proto 消息指针的 Elem 本身就是 struct，与下面的分支完全等价，
			// 因此不需要为了识别 proto 而在此引入 protobuf 依赖
			// [EN] A proto message pointer's Elem is already a struct, so the struct
			// branch below covers it identically; no protobuf import needed here
			if fieldType.Type.Elem().Kind() == reflect.Struct {
				field.Set(reflect.New(fieldType.Type.Elem()))
			}
		}
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
