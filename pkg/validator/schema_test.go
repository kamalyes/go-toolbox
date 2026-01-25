/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-25 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\schema_test.go
 * @Description: JSON Schema 验证测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJSONSchemaBasicTypes(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		schema  JSONSchema
		wantOK  bool
		message string
	}{
		{
			name:   "字符串类型-通过",
			data:   `"hello"`,
			schema: JSONSchema{Type: "string"},
			wantOK: true,
		},
		{
			name:   "字符串类型-失败",
			data:   `123`,
			schema: JSONSchema{Type: "string"},
			wantOK: false,
		},
		{
			name:   "数字类型-通过",
			data:   `123`,
			schema: JSONSchema{Type: "number"},
			wantOK: true,
		},
		{
			name:   "整数类型-通过",
			data:   `42`,
			schema: JSONSchema{Type: "integer"},
			wantOK: true,
		},
		{
			name:   "整数类型-失败(小数)",
			data:   `42.5`,
			schema: JSONSchema{Type: "integer"},
			wantOK: false,
		},
		{
			name:   "布尔类型-通过",
			data:   `true`,
			schema: JSONSchema{Type: "boolean"},
			wantOK: true,
		},
		{
			name:   "null类型-通过",
			data:   `null`,
			schema: JSONSchema{Type: "null"},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONSchema(tt.data, tt.schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestValidateJSONSchemaObject(t *testing.T) {
	schema := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"name": {Type: "string"},
			"age":  {Type: "integer"},
		},
		Required: []string{"name"},
	}

	tests := []struct {
		name   string
		data   string
		wantOK bool
	}{
		{
			name:   "完整对象-通过",
			data:   `{"name": "张三", "age": 25}`,
			wantOK: true,
		},
		{
			name:   "只有必需字段-通过",
			data:   `{"name": "张三"}`,
			wantOK: true,
		},
		{
			name:   "缺少必需字段-失败",
			data:   `{"age": 25}`,
			wantOK: false,
		},
		{
			name:   "类型错误-失败",
			data:   `{"name": "张三", "age": "25"}`,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONSchema(tt.data, schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestValidateJSONSchemaArray(t *testing.T) {
	minItems := 1
	maxItems := 5

	schema := JSONSchema{
		Type:     "array",
		MinItems: &minItems,
		MaxItems: &maxItems,
		Items: &JSONSchema{
			Type: "string",
		},
	}

	tests := []struct {
		name   string
		data   string
		wantOK bool
	}{
		{
			name:   "正常数组-通过",
			data:   `["a", "b", "c"]`,
			wantOK: true,
		},
		{
			name:   "空数组-失败(minItems)",
			data:   `[]`,
			wantOK: false,
		},
		{
			name:   "超长数组-失败(maxItems)",
			data:   `["a", "b", "c", "d", "e", "f"]`,
			wantOK: false,
		},
		{
			name:   "元素类型错误-失败",
			data:   `["a", 123, "c"]`,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONSchema(tt.data, schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestValidateJSONSchemaStringConstraints(t *testing.T) {
	minLen := 3
	maxLen := 10

	schema := JSONSchema{
		Type:      "string",
		MinLength: &minLen,
		MaxLength: &maxLen,
		Pattern:   `^[a-z]+$`,
	}

	tests := []struct {
		name   string
		data   string
		wantOK bool
	}{
		{
			name:   "符合所有约束-通过",
			data:   `"hello"`,
			wantOK: true,
		},
		{
			name:   "太短-失败",
			data:   `"ab"`,
			wantOK: false,
		},
		{
			name:   "太长-失败",
			data:   `"verylongstring"`,
			wantOK: false,
		},
		{
			name:   "不匹配正则-失败",
			data:   `"Hello"`,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONSchema(tt.data, schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestValidateJSONSchemaNumberConstraints(t *testing.T) {
	min := 0.0
	max := 100.0

	schema := JSONSchema{
		Type:    "number",
		Minimum: &min,
		Maximum: &max,
	}

	tests := []struct {
		name   string
		data   string
		wantOK bool
	}{
		{
			name:   "范围内-通过",
			data:   `50`,
			wantOK: true,
		},
		{
			name:   "边界值-最小-通过",
			data:   `0`,
			wantOK: true,
		},
		{
			name:   "边界值-最大-通过",
			data:   `100`,
			wantOK: true,
		},
		{
			name:   "小于最小值-失败",
			data:   `-1`,
			wantOK: false,
		},
		{
			name:   "大于最大值-失败",
			data:   `101`,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONSchema(tt.data, schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestValidateJSONSchemaEnum(t *testing.T) {
	schema := JSONSchema{
		Type: "string",
		Enum: []interface{}{"red", "green", "blue"},
	}

	tests := []struct {
		name   string
		data   string
		wantOK bool
	}{
		{
			name:   "枚举值-red-通过",
			data:   `"red"`,
			wantOK: true,
		},
		{
			name:   "枚举值-blue-通过",
			data:   `"blue"`,
			wantOK: true,
		},
		{
			name:   "非枚举值-失败",
			data:   `"yellow"`,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONSchema(tt.data, schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestValidateJSONSchemaUniqueItems(t *testing.T) {
	schema := JSONSchema{
		Type:        "array",
		UniqueItems: true,
		Items: &JSONSchema{
			Type: "string",
		},
	}

	tests := []struct {
		name   string
		data   string
		wantOK bool
	}{
		{
			name:   "唯一元素-通过",
			data:   `["a", "b", "c"]`,
			wantOK: true,
		},
		{
			name:   "重复元素-失败",
			data:   `["a", "b", "a"]`,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONSchema(tt.data, schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestValidateJSONSchemaNestedObject(t *testing.T) {
	schema := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"user": {
				Type: "object",
				Properties: map[string]JSONSchema{
					"name": {Type: "string"},
					"age":  {Type: "integer"},
				},
				Required: []string{"name"},
			},
		},
		Required: []string{"user"},
	}

	tests := []struct {
		name   string
		data   string
		wantOK bool
	}{
		{
			name:   "嵌套对象-通过",
			data:   `{"user": {"name": "张三", "age": 25}}`,
			wantOK: true,
		},
		{
			name:   "嵌套对象缺少必需字段-失败",
			data:   `{"user": {"age": 25}}`,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONSchema(tt.data, schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestValidateStructWithSchema(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	schema := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"name":  {Type: "string"},
			"age":   {Type: "integer"},
			"email": {Type: "string", Pattern: `^[\w\.-]+@[\w\.-]+\.\w+$`},
		},
		Required: []string{"name", "email"},
	}

	tests := []struct {
		name   string
		user   User
		wantOK bool
	}{
		{
			name: "有效用户-通过",
			user: User{
				Name:  "张三",
				Age:   25,
				Email: "zhangsan@example.com",
			},
			wantOK: true,
		},
		{
			name: "邮箱格式错误-失败",
			user: User{
				Name:  "李四",
				Age:   30,
				Email: "invalid-email",
			},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateStructWithSchema(tt.user, schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestSchemaBuilder(t *testing.T) {
	// 使用构建器创建 Schema
	schema := NewSchemaBuilder().
		Type("object").
		StringProperty("name", 1, 50).
		NumberProperty("age", floatPtr(0), floatPtr(150)).
		Required("name").
		Build()

	// 验证数据
	data := `{"name": "测试", "age": 25}`
	result := ValidateJSONSchema(data, schema)
	assert.True(t, result.Success, "Schema构建器生成的Schema应该有效")
}

func TestQuickSchema(t *testing.T) {
	// 快速创建简单 Schema
	schema := QuickSchema(map[string]string{
		"name":  "string",
		"age":   "number",
		"email": "string",
	}, "name", "email")

	tests := []struct {
		name   string
		data   string
		wantOK bool
	}{
		{
			name:   "完整数据-通过",
			data:   `{"name": "张三", "age": 25, "email": "test@example.com"}`,
			wantOK: true,
		},
		{
			name:   "缺少必需字段-失败",
			data:   `{"name": "张三"}`,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONSchema(tt.data, schema)
			assert.Equal(t, tt.wantOK, result.Success, "验证结果不符: %s", result.Message)
		})
	}
}

func TestValidateJSONSchemaAdditionalProperties(t *testing.T) {
	// 不允许额外属性
	schemaNoAdditional := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"name": {Type: "string"},
		},
		AdditionalProperties: false,
	}

	result := ValidateJSONSchema(`{"name": "test", "age": 25}`, schemaNoAdditional)
	assert.False(t, result.Success, "不应允许额外属性")

	// 允许额外属性
	schemaAllowAdditional := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"name": {Type: "string"},
		},
		AdditionalProperties: true,
	}

	result = ValidateJSONSchema(`{"name": "test", "age": 25}`, schemaAllowAdditional)
	assert.True(t, result.Success, "应允许额外属性")
}

func TestFormatSchemaError(t *testing.T) {
	result := CompareResult{
		Success: false,
		Message: "验证失败: 缺少必需字段 'name'",
	}

	errorMsg := FormatSchemaError(result)
	assert.Contains(t, errorMsg, "JSON Schema 验证失败")
	assert.Contains(t, errorMsg, "缺少必需字段")
}

func TestValidateJSONSchemaComplexExample(t *testing.T) {
	// 复杂的电商订单 Schema
	schema := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"orderId": {
				Type:    "string",
				Pattern: `^ORD-\d{8}$`,
			},
			"customer": {
				Type: "object",
				Properties: map[string]JSONSchema{
					"name":  {Type: "string", MinLength: intPtr(1)},
					"email": {Type: "string", Pattern: `^[\w\.-]+@[\w\.-]+\.\w+$`},
				},
				Required: []string{"name", "email"},
			},
			"items": {
				Type:     "array",
				MinItems: intPtr(1),
				Items: &JSONSchema{
					Type: "object",
					Properties: map[string]JSONSchema{
						"productId": {Type: "string"},
						"quantity":  {Type: "integer", Minimum: floatPtr(1)},
						"price":     {Type: "number", Minimum: floatPtr(0)},
					},
					Required: []string{"productId", "quantity", "price"},
				},
			},
			"status": {
				Type: "string",
				Enum: []interface{}{"pending", "paid", "shipped", "delivered"},
			},
		},
		Required: []string{"orderId", "customer", "items", "status"},
	}

	validOrder := `{
		"orderId": "ORD-20260125",
		"customer": {
			"name": "张三",
			"email": "zhangsan@example.com"
		},
		"items": [
			{
				"productId": "PROD-001",
				"quantity": 2,
				"price": 99.99
			}
		],
		"status": "pending"
	}`

	result := ValidateJSONSchema(validOrder, schema)
	assert.True(t, result.Success, "有效订单应该通过验证: %s", result.Message)

	invalidOrder := `{
		"orderId": "INVALID",
		"customer": {
			"name": "李四"
		},
		"items": [],
		"status": "unknown"
	}`

	result = ValidateJSONSchema(invalidOrder, schema)
	assert.False(t, result.Success, "无效订单应该验证失败")
}

// TestDeepNestedValidation 深度嵌套验证测试
func TestDeepNestedValidation(t *testing.T) {
	// 5层嵌套的复杂 schema
	schema := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"level1": {
				Type: "object",
				Properties: map[string]JSONSchema{
					"level2": {
						Type: "object",
						Properties: map[string]JSONSchema{
							"level3": {
								Type: "object",
								Properties: map[string]JSONSchema{
									"level4": {
										Type: "object",
										Properties: map[string]JSONSchema{
											"level5": {
												Type: "array",
												Items: &JSONSchema{
													Type: "object",
													Properties: map[string]JSONSchema{
														"id":   {Type: "integer"},
														"name": {Type: "string"},
													},
													Required: []string{"id", "name"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	validData := `{
		"level1": {
			"level2": {
				"level3": {
					"level4": {
						"level5": [
							{"id": 1, "name": "item1"},
							{"id": 2, "name": "item2"}
						]
					}
				}
			}
		}
	}`

	result := ValidateJSONSchema(validData, schema)
	assert.True(t, result.Success, "深度嵌套数据验证失败: %s", result.Message)

	// 缺少必需字段
	invalidData := `{
		"level1": {
			"level2": {
				"level3": {
					"level4": {
						"level5": [
							{"id": 1}
						]
					}
				}
			}
		}
	}`

	result = ValidateJSONSchema(invalidData, schema)
	assert.False(t, result.Success, "应该检测到缺少必需字段")
	assert.Contains(t, result.Message, "name", "错误信息应包含缺失的字段名")
}

// TestLargeArrayValidation 大数组验证测试
func TestLargeArrayValidation(t *testing.T) {
	schema := JSONSchema{
		Type: "array",
		Items: &JSONSchema{
			Type: "object",
			Properties: map[string]JSONSchema{
				"id":    {Type: "integer"},
				"email": {Type: "string", Pattern: `^[\w\.-]+@[\w\.-]+\.\w+$`},
			},
			Required: []string{"id", "email"},
		},
	}

	// 构造100个元素的数组
	data := `[`
	for i := 0; i < 100; i++ {
		if i > 0 {
			data += ","
		}
		data += `{"id":` + string(rune('0'+i%10)) + `,"email":"user` + string(rune('0'+i%10)) + `@test.com"}`
	}
	data += `]`

	result := ValidateJSONSchema(data, schema)
	assert.True(t, result.Success, "大数组验证失败: %s", result.Message)
}

// BenchmarkJSONSchemaSimple 简单对象基准测试
func BenchmarkJSONSchemaSimple(b *testing.B) {
	schema := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"name": {Type: "string"},
			"age":  {Type: "integer"},
		},
		Required: []string{"name"},
	}

	data := `{"name": "张三", "age": 25}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateJSONSchema(data, schema)
	}
}

// BenchmarkJSONSchemaNested 嵌套对象基准测试
func BenchmarkJSONSchemaNested(b *testing.B) {
	schema := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"user": {
				Type: "object",
				Properties: map[string]JSONSchema{
					"profile": {
						Type: "object",
						Properties: map[string]JSONSchema{
							"name": {Type: "string"},
							"age":  {Type: "integer"},
						},
					},
				},
			},
		},
	}

	data := `{"user": {"profile": {"name": "张三", "age": 25}}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateJSONSchema(data, schema)
	}
}

// BenchmarkJSONSchemaArray 数组基准测试
func BenchmarkJSONSchemaArray(b *testing.B) {
	schema := JSONSchema{
		Type: "array",
		Items: &JSONSchema{
			Type: "object",
			Properties: map[string]JSONSchema{
				"id":   {Type: "integer"},
				"name": {Type: "string"},
			},
		},
	}

	data := `[{"id": 1, "name": "item1"}, {"id": 2, "name": "item2"}, {"id": 3, "name": "item3"}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateJSONSchema(data, schema)
	}
}

// BenchmarkJSONSchemaDeepNested 深度嵌套基准测试
func BenchmarkJSONSchemaDeepNested(b *testing.B) {
	schema := JSONSchema{
		Type: "object",
		Properties: map[string]JSONSchema{
			"level1": {
				Type: "object",
				Properties: map[string]JSONSchema{
					"level2": {
						Type: "object",
						Properties: map[string]JSONSchema{
							"level3": {
								Type: "object",
								Properties: map[string]JSONSchema{
									"level4": {
										Type: "array",
										Items: &JSONSchema{
											Type: "object",
											Properties: map[string]JSONSchema{
												"id":   {Type: "integer"},
												"name": {Type: "string"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	data := `{"level1": {"level2": {"level3": {"level4": [{"id": 1, "name": "item1"}]}}}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateJSONSchema(data, schema)
	}
}

// 辅助函数 - 使用公共辅助函数
func intPtr(i int) *int {
	return IntPtr(i)
}

func floatPtr(f float64) *float64 {
	return Float64Ptr(f)
}
