/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-25 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\schema.go
 * @Description: JSON Schema 验证 - 高性能嵌套验证
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// JSONSchema JSON Schema 定义
type JSONSchema struct {
	Type                 string                `json:"type,omitempty"`                 // string, number, integer, boolean, array, object, null
	Properties           map[string]JSONSchema `json:"properties,omitempty"`           // 对象属性定义
	Required             []string              `json:"required,omitempty"`             // 必需字段
	Items                *JSONSchema           `json:"items,omitempty"`                // 数组元素 schema
	Enum                 []interface{}         `json:"enum,omitempty"`                 // 枚举值
	Minimum              *float64              `json:"minimum,omitempty"`              // 最小值
	Maximum              *float64              `json:"maximum,omitempty"`              // 最大值
	MinLength            *int                  `json:"minLength,omitempty"`            // 最小长度
	MaxLength            *int                  `json:"maxLength,omitempty"`            // 最大长度
	Pattern              string                `json:"pattern,omitempty"`              // 正则模式
	MinItems             *int                  `json:"minItems,omitempty"`             // 数组最小元素数
	MaxItems             *int                  `json:"maxItems,omitempty"`             // 数组最大元素数
	UniqueItems          bool                  `json:"uniqueItems,omitempty"`          // 数组元素唯一性
	AdditionalProperties interface{}           `json:"additionalProperties,omitempty"` // 额外属性
}

// ValidateJSONSchema 验证 JSON 数据是否符合 Schema
// data: JSON 字符串或已解析的数据
// schema: JSON Schema 定义（可以是 JSONSchema 结构体或 JSON 字符串）
//
// 返回 CompareResult，包含验证结果和详细信息
func ValidateJSONSchema(data interface{}, schema interface{}) CompareResult {
	// 解析 schema
	var schemaObj JSONSchema
	switch s := schema.(type) {
	case JSONSchema:
		schemaObj = s
	case string:
		if err := json.Unmarshal([]byte(s), &schemaObj); err != nil {
			return CompareResult{
				Success: false,
				Message: fmt.Sprintf("Schema 解析失败: %v", err),
			}
		}
	case map[string]interface{}:
		// 将 map 转换为 JSONSchema
		schemaBytes, _ := json.Marshal(s)
		if err := json.Unmarshal(schemaBytes, &schemaObj); err != nil {
			return CompareResult{
				Success: false,
				Message: fmt.Sprintf("Schema 转换失败: %v", err),
			}
		}
	default:
		return CompareResult{
			Success: false,
			Message: "不支持的 Schema 类型",
		}
	}

	// 解析数据
	var dataObj interface{}
	switch d := data.(type) {
	case string:
		if err := json.Unmarshal([]byte(d), &dataObj); err != nil {
			return CompareResult{
				Success: false,
				Message: fmt.Sprintf("JSON 数据解析失败: %v", err),
				Actual:  d,
			}
		}
	default:
		dataObj = d
	}

	// 验证数据
	if err := validateValue(dataObj, schemaObj, ""); err != nil {
		return CompareResult{
			Success: false,
			Message: err.Error(),
			Actual:  fmt.Sprintf("%v", dataObj),
			Expect:  "符合 Schema 定义",
		}
	}

	return CompareResult{
		Success: true,
		Message: "JSON Schema 验证通过",
		Actual:  fmt.Sprintf("%v", dataObj),
		Expect:  "符合 Schema 定义",
	}
}

// validateValue 递归验证值
func validateValue(value interface{}, schema JSONSchema, path string) error {
	// 类型验证
	if schema.Type != "" {
		if err := validateType(value, schema.Type, path); err != nil {
			return err
		}
	}

	// 枚举验证
	if len(schema.Enum) > 0 {
		found := false
		for _, enumVal := range schema.Enum {
			if fmt.Sprintf("%v", value) == fmt.Sprintf("%v", enumVal) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s: 值必须是枚举之一 %v，实际: %v", path, schema.Enum, value)
		}
	}

	// 根据类型进行特定验证
	switch schema.Type {
	case "string":
		if err := validateString(value, schema, path); err != nil {
			return err
		}
	case "number", "integer":
		if err := validateNumber(value, schema, path); err != nil {
			return err
		}
	case "array":
		if err := validateArray(value, schema, path); err != nil {
			return err
		}
	case "object":
		if err := validateObject(value, schema, path); err != nil {
			return err
		}
	}

	return nil
}

// validateType 验证类型 - 使用公共辅助函数
func validateType(value interface{}, expectedType string, path string) error {
	if value == nil {
		if expectedType == "null" {
			return nil
		}
		return fmt.Errorf("%s: 期望类型 %s，实际类型 null", path, expectedType)
	}

	// 使用公共函数获取类型
	kind := GetReflectKind(value)

	switch expectedType {
	case "null":
		return fmt.Errorf("%s: 期望 null，实际非 null", path)

	case "boolean":
		if kind != reflect.Bool {
			return fmt.Errorf("%s: 期望类型 boolean，实际类型 %v", path, kind)
		}

	case "string":
		if kind != reflect.String {
			return fmt.Errorf("%s: 期望类型 string，实际类型 %v", path, kind)
		}

	case "integer":
		if IsIntegerKind(kind) {
			return nil
		}
		// 允许 float64 但必须是整数
		if kind == reflect.Float32 || kind == reflect.Float64 {
			if f, ok := ToFloat64(value); ok && IsWholeNumber(f) {
				return nil
			}
			return fmt.Errorf("%s: 期望整数，实际: %v", path, value)
		}
		return fmt.Errorf("%s: 期望类型 integer，实际类型 %v", path, kind)

	case "number":
		if !IsNumericKind(kind) {
			return fmt.Errorf("%s: 期望类型 number，实际类型 %v", path, kind)
		}

	case "array":
		if kind != reflect.Slice && kind != reflect.Array {
			return fmt.Errorf("%s: 期望类型 array，实际类型 %v", path, kind)
		}

	case "object":
		if kind != reflect.Map {
			return fmt.Errorf("%s: 期望类型 object，实际类型 %v", path, kind)
		}
	}

	return nil
} // validateString 验证字符串 - 使用公共辅助函数
func validateString(value interface{}, schema JSONSchema, path string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s: 不是字符串类型", path)
	}

	strLen := len(str) // 只计算一次长度

	// 长度验证
	if schema.MinLength != nil && strLen < *schema.MinLength {
		return fmt.Errorf("%s: 字符串长度不能小于 %d，实际: %d", path, *schema.MinLength, strLen)
	}
	if schema.MaxLength != nil && strLen > *schema.MaxLength {
		return fmt.Errorf("%s: 字符串长度不能大于 %d，实际: %d", path, *schema.MaxLength, strLen)
	}

	// 正则验证 - 使用公共缓存函数
	if schema.Pattern != "" {
		re, err := GetCompiledRegex(schema.Pattern)
		if err != nil {
			return fmt.Errorf("%s: 正则表达式无效: %v", path, err)
		}
		if !re.MatchString(str) {
			return fmt.Errorf("%s: 字符串不匹配模式 %s", path, schema.Pattern)
		}
	}

	return nil
}

// validateNumber 验证数字 - 使用公共辅助函数
func validateNumber(value interface{}, schema JSONSchema, path string) error {
	num, ok := ToFloat64(value)
	if !ok {
		return fmt.Errorf("%s: 不是数字类型", path)
	}

	// 范围验证
	if schema.Minimum != nil && num < *schema.Minimum {
		return fmt.Errorf("%s: 数值不能小于 %v，实际: %v", path, *schema.Minimum, num)
	}
	if schema.Maximum != nil && num > *schema.Maximum {
		return fmt.Errorf("%s: 数值不能大于 %v，实际: %v", path, *schema.Maximum, num)
	}

	return nil
}

// validateArray 验证数组 - 优化性能
func validateArray(value interface{}, schema JSONSchema, path string) error {
	arr, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("%s: 不是数组类型", path)
	}

	arrLen := len(arr) // 只计算一次长度

	// 长度验证
	if schema.MinItems != nil && arrLen < *schema.MinItems {
		return fmt.Errorf("%s: 数组元素数量不能小于 %d，实际: %d", path, *schema.MinItems, arrLen)
	}
	if schema.MaxItems != nil && arrLen > *schema.MaxItems {
		return fmt.Errorf("%s: 数组元素数量不能大于 %d，实际: %d", path, *schema.MaxItems, arrLen)
	}

	// 唯一性验证 - 使用 map 提高性能
	if schema.UniqueItems && arrLen > 0 {
		seen := make(map[string]struct{}, arrLen)
		for _, item := range arr {
			key := fmt.Sprintf("%v", item)
			if _, exists := seen[key]; exists {
				return fmt.Errorf("%s: 数组元素必须唯一，发现重复: %v", path, item)
			}
			seen[key] = struct{}{}
		}
	}

	// 元素验证 - 快速路径：如果没有 Items 定义，跳过验证
	if schema.Items != nil {
		// 预先计算路径前缀，避免重复拼接
		pathPrefix := path
		if pathPrefix != "" {
			pathPrefix += "["
		} else {
			pathPrefix = "["
		}

		for i, item := range arr {
			itemPath := fmt.Sprintf("%s%d]", pathPrefix, i)
			if err := validateValue(item, *schema.Items, itemPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateObject 验证对象 - 优化嵌套性能
func validateObject(value interface{}, schema JSONSchema, path string) error {
	obj, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%s: 不是对象类型", path)
	}

	// 必需字段验证 - 快速检查
	for _, required := range schema.Required {
		if _, exists := obj[required]; !exists {
			return fmt.Errorf("%s: 缺少必需字段 '%s'", path, required)
		}
	}

	// 预先计算路径分隔符
	var pathSep string
	if path == "" {
		pathSep = ""
	} else {
		pathSep = "."
	}

	// 属性验证 - 优化嵌套对象的路径拼接
	hasAdditionalProps := schema.AdditionalProperties != nil
	hasProperties := len(schema.Properties) > 0

	for key, val := range obj {
		propSchema, hasPropSchema := schema.Properties[key]

		if hasPropSchema {
			// 有定义的属性，按 schema 验证
			propPath := path + pathSep + key
			if err := validateValue(val, propSchema, propPath); err != nil {
				return err
			}
		} else if hasAdditionalProps {
			// 没有定义的属性，检查 additionalProperties
			switch ap := schema.AdditionalProperties.(type) {
			case bool:
				if !ap {
					return fmt.Errorf("%s: 不允许额外的属性 '%s'", path, key)
				}
			case map[string]interface{}:
				// additionalProperties 是一个 schema
				var apSchema JSONSchema
				apBytes, _ := json.Marshal(ap)
				if err := json.Unmarshal(apBytes, &apSchema); err == nil {
					propPath := path + pathSep + key
					if err := validateValue(val, apSchema, propPath); err != nil {
						return err
					}
				}
			case JSONSchema:
				// 直接是 JSONSchema 类型
				propPath := path + pathSep + key
				if err := validateValue(val, ap, propPath); err != nil {
					return err
				}
			}
		} else if hasProperties {
			// 有 properties 定义但当前属性不在其中，且没有 additionalProperties
			// 默认允许（宽松模式）
			continue
		}
	}

	return nil
}

// ValidateStructWithSchema 验证结构体是否符合 JSON Schema
//
// 将结构体转换为 JSON 后进行 Schema 验证
func ValidateStructWithSchema(structData interface{}, schema interface{}) CompareResult {
	// 将结构体转换为 JSON
	jsonBytes, err := json.Marshal(structData)
	if err != nil {
		return CompareResult{
			Success: false,
			Message: fmt.Sprintf("结构体序列化失败: %v", err),
		}
	}

	// 解析为 map 以便验证
	var dataMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &dataMap); err != nil {
		return CompareResult{
			Success: false,
			Message: fmt.Sprintf("JSON 解析失败: %v", err),
		}
	}

	return ValidateJSONSchema(dataMap, schema)
}

// SchemaBuilder JSON Schema 构建器（链式调用）
type SchemaBuilder struct {
	schema JSONSchema
}

// NewSchemaBuilder 创建 Schema 构建器
func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		schema: JSONSchema{
			Properties: make(map[string]JSONSchema),
		},
	}
}

// Type 设置类型
func (b *SchemaBuilder) Type(t string) *SchemaBuilder {
	b.schema.Type = t
	return b
}

// Required 设置必需字段
func (b *SchemaBuilder) Required(fields ...string) *SchemaBuilder {
	b.schema.Required = append(b.schema.Required, fields...)
	return b
}

// Property 添加属性
func (b *SchemaBuilder) Property(name string, propSchema JSONSchema) *SchemaBuilder {
	if b.schema.Properties == nil {
		b.schema.Properties = make(map[string]JSONSchema)
	}
	b.schema.Properties[name] = propSchema
	return b
}

// StringProperty 添加字符串属性
func (b *SchemaBuilder) StringProperty(name string, minLen, maxLen int) *SchemaBuilder {
	schema := JSONSchema{Type: "string"}
	if minLen > 0 {
		schema.MinLength = &minLen
	}
	if maxLen > 0 {
		schema.MaxLength = &maxLen
	}
	return b.Property(name, schema)
}

// NumberProperty 添加数字属性
func (b *SchemaBuilder) NumberProperty(name string, min, max *float64) *SchemaBuilder {
	schema := JSONSchema{
		Type:    "number",
		Minimum: min,
		Maximum: max,
	}
	return b.Property(name, schema)
}

// ArrayProperty 添加数组属性
func (b *SchemaBuilder) ArrayProperty(name string, items JSONSchema) *SchemaBuilder {
	schema := JSONSchema{
		Type:  "array",
		Items: &items,
	}
	return b.Property(name, schema)
}

// Enum 设置枚举值
func (b *SchemaBuilder) Enum(values ...interface{}) *SchemaBuilder {
	b.schema.Enum = values
	return b
}

// Build 构建 Schema
func (b *SchemaBuilder) Build() JSONSchema {
	return b.schema
}

// BuildJSON 构建 Schema JSON 字符串
func (b *SchemaBuilder) BuildJSON() string {
	data, _ := json.MarshalIndent(b.schema, "", "  ")
	return string(data)
}

// QuickSchema 快速创建简单的对象 Schema（语法糖）
//
// 示例:
//
//	schema := QuickSchema(map[string]string{
//	    "name": "string",
//	    "age": "number",
//	    "email": "string",
//	}, "name", "email")  // name 和 email 是必需的
func QuickSchema(properties map[string]string, required ...string) JSONSchema {
	schema := JSONSchema{
		Type:       "object",
		Properties: make(map[string]JSONSchema),
		Required:   required,
	}

	for name, typeStr := range properties {
		schema.Properties[name] = JSONSchema{Type: typeStr}
	}

	return schema
}

// FormatSchemaError 格式化 Schema 验证错误信息
func FormatSchemaError(result CompareResult) string {
	if result.Success {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("JSON Schema 验证失败:\n")
	sb.WriteString(result.Message)

	return sb.String()
}
