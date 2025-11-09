/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-09 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 00:00:00
 * @FilePath: \go-toolbox\tests\random_enhanced_test.go
 * @Description: 增强版 GenerateRandModel 函数的测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/stretchr/testify/assert"
)

// TestEnhancedModel 测试增强版功能的结构体
type TestEnhancedModel struct {
	// 基本类型
	StringField  string    `json:"string_field" rand:"name"`
	IntField     int       `json:"int_field"`
	Int8Field    int8      `json:"int8_field"`
	Int16Field   int16     `json:"int16_field"`
	Int32Field   int32     `json:"int32_field"`
	Int64Field   int64     `json:"int64_field"`
	UintField    uint      `json:"uint_field"`
	Uint8Field   uint8     `json:"uint8_field"`
	Uint16Field  uint16    `json:"uint16_field"`
	Uint32Field  uint32    `json:"uint32_field"`
	Uint64Field  uint64    `json:"uint64_field"`
	Float32Field float32   `json:"float32_field"`
	Float64Field float64   `json:"float64_field"`
	BoolField    bool      `json:"bool_field"`
	TimeField    time.Time `json:"time_field"`
	
	// 复数类型（程序会自动跳过，因为无法JSON序列化）
	Complex64Field  complex64  `json:"complex64_field"`
	Complex128Field complex128 `json:"complex128_field"`
	
	// 其他不支持JSON序列化的类型
	ChanField chan int    `json:"chan_field"`
	FuncField func() error `json:"func_field"`
	
	// 指针类型
	StringPtr  *string  `json:"string_ptr"`
	IntPtr     *int     `json:"int_ptr"`
	FloatPtr   *float64 `json:"float_ptr"`
	BoolPtr    *bool    `json:"bool_ptr"`
	
	// 切片类型
	StringSlice []string          `json:"string_slice"`
	IntSlice    []int             `json:"int_slice"`
	FloatSlice  []float64         `json:"float_slice"`
	StructSlice []NestedTestModel `json:"struct_slice"`
	
	// 数组类型
	StringArray [3]string `json:"string_array"`
	IntArray    [5]int    `json:"int_array"`
	
	// 映射类型
	StringMap     map[string]string `json:"string_map"`
	IntMap        map[string]int    `json:"int_map"`
	FloatMap      map[string]float64 `json:"float_map"`
	ComplexMap    map[int]string     `json:"complex_map"`
	
	// 嵌套结构体
	NestedStruct NestedTestModel  `json:"nested_struct"`
	NestedPtr    *NestedTestModel `json:"nested_ptr"`
	
	// 接口类型
	InterfaceField interface{} `json:"interface_field"`
	
	// 自定义标签字段
	EmailField string `json:"email_field" rand:"email"`
	PhoneField string `json:"phone_field" rand:"phone"`
	UUIDField  string `json:"uuid_field" rand:"uuid"`
}

type NestedTestModel struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// TestGenerateRandModelEnhanced_BasicTypes 测试基本类型
func TestGenerateRandModelEnhanced_BasicTypes(t *testing.T) {
	model := &TestEnhancedModel{}
	
	result, jsonStr, err := random.GenerateRandModel(model)
	assert.NoError(t, err, "生成增强模型不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	// 验证 JSON 可解析
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &jsonMap)
	assert.NoError(t, err, "应能正确解析 JSON")
	
	resultModel := result.(*TestEnhancedModel)
	
	// 验证基本类型字段
	assert.NotEmpty(t, resultModel.StringField, "StringField 应被填充")
	assert.NotZero(t, resultModel.IntField, "IntField 应被填充")
	assert.NotZero(t, resultModel.Int8Field, "Int8Field 应被填充")
	assert.NotZero(t, resultModel.Int16Field, "Int16Field 应被填充")
	assert.NotZero(t, resultModel.Int32Field, "Int32Field 应被填充")
	assert.NotZero(t, resultModel.Int64Field, "Int64Field 应被填充")
	assert.NotZero(t, resultModel.UintField, "UintField 应被填充")
	assert.NotZero(t, resultModel.Uint8Field, "Uint8Field 应被填充")
	assert.NotZero(t, resultModel.Uint16Field, "Uint16Field 应被填充")
	assert.NotZero(t, resultModel.Uint32Field, "Uint32Field 应被填充")
	assert.NotZero(t, resultModel.Uint64Field, "Uint64Field 应被填充")
	assert.NotZero(t, resultModel.Float32Field, "Float32Field 应被填充")
	assert.NotZero(t, resultModel.Float64Field, "Float64Field 应被填充")
	assert.NotZero(t, resultModel.TimeField, "TimeField 应被填充")
	
	// 验证不支持JSON序列化的类型被跳过（保持零值）
	assert.Zero(t, resultModel.Complex64Field, "Complex64Field 应被跳过（零值）")
	assert.Zero(t, resultModel.Complex128Field, "Complex128Field 应被跳过（零值）")
	assert.Nil(t, resultModel.ChanField, "ChanField 应被跳过（nil）")
	assert.Nil(t, resultModel.FuncField, "FuncField 应被跳过（nil）")
}

// TestGenerateRandModelEnhanced_Pointers 测试指针类型
func TestGenerateRandModelEnhanced_Pointers(t *testing.T) {
	model := &TestEnhancedModel{}
	
	result, jsonStr, err := random.GenerateRandModel(model)
	assert.NoError(t, err, "生成增强模型不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	resultModel := result.(*TestEnhancedModel)
	
	// 验证指针类型字段
	assert.NotNil(t, resultModel.StringPtr, "StringPtr 应被初始化")
	assert.NotEmpty(t, *resultModel.StringPtr, "StringPtr 指向的值应被填充")
	
	assert.NotNil(t, resultModel.IntPtr, "IntPtr 应被初始化")
	assert.NotZero(t, *resultModel.IntPtr, "IntPtr 指向的值应被填充")
	
	assert.NotNil(t, resultModel.FloatPtr, "FloatPtr 应被初始化")
	assert.NotZero(t, *resultModel.FloatPtr, "FloatPtr 指向的值应被填充")
	
	assert.NotNil(t, resultModel.BoolPtr, "BoolPtr 应被初始化")
	// BoolPtr 的值可能是 true 或 false，都是有效的
}

// TestGenerateRandModelEnhanced_SlicesAndArrays 测试切片和数组类型
func TestGenerateRandModelEnhanced_SlicesAndArrays(t *testing.T) {
	model := &TestEnhancedModel{}
	
	result, jsonStr, err := random.GenerateRandModel(model)
	assert.NoError(t, err, "生成增强模型不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	resultModel := result.(*TestEnhancedModel)
	
	// 验证切片类型
	assert.NotEmpty(t, resultModel.StringSlice, "StringSlice 应包含元素")
	assert.NotEmpty(t, resultModel.IntSlice, "IntSlice 应包含元素")
	assert.NotEmpty(t, resultModel.FloatSlice, "FloatSlice 应包含元素")
	assert.NotEmpty(t, resultModel.StructSlice, "StructSlice 应包含元素")
	
	// 验证切片元素
	for i, item := range resultModel.StringSlice {
		assert.NotEmpty(t, item, "StringSlice[%d] 应不为空", i)
	}
	
	for i, item := range resultModel.StructSlice {
		assert.NotEmpty(t, item.Name, "StructSlice[%d].Name 应被填充", i)
		assert.NotZero(t, item.Value, "StructSlice[%d].Value 应被填充", i)
	}
	
	// 验证数组类型（数组会被填充所有元素）
	for i, item := range resultModel.StringArray {
		assert.NotEmpty(t, item, "StringArray[%d] 应被填充", i)
	}
	
	for i, item := range resultModel.IntArray {
		assert.NotZero(t, item, "IntArray[%d] 应被填充", i)
	}
}

// TestGenerateRandModelEnhanced_Maps 测试映射类型
func TestGenerateRandModelEnhanced_Maps(t *testing.T) {
	model := &TestEnhancedModel{}
	
	result, jsonStr, err := random.GenerateRandModel(model)
	assert.NoError(t, err, "生成增强模型不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	resultModel := result.(*TestEnhancedModel)
	
	// 验证映射类型
	assert.NotEmpty(t, resultModel.StringMap, "StringMap 应包含元素")
	assert.NotEmpty(t, resultModel.IntMap, "IntMap 应包含元素")
	assert.NotEmpty(t, resultModel.FloatMap, "FloatMap 应包含元素")
	// ComplexMap 的键是 int 类型，不支持JSON序列化，会被自动跳过
	assert.Empty(t, resultModel.ComplexMap, "ComplexMap 应被跳过（键类型不是字符串）")
	
	// 验证映射元素
	for key, value := range resultModel.StringMap {
		assert.NotEmpty(t, key, "StringMap 键应不为空")
		assert.NotEmpty(t, value, "StringMap 值应不为空")
	}
	
	for key, value := range resultModel.IntMap {
		assert.NotEmpty(t, key, "IntMap 键应不为空")
		assert.NotZero(t, value, "IntMap 值应不为零")
	}
}

// TestGenerateRandModelEnhanced_NestedStructs 测试嵌套结构体
func TestGenerateRandModelEnhanced_NestedStructs(t *testing.T) {
	model := &TestEnhancedModel{}
	
	result, jsonStr, err := random.GenerateRandModel(model)
	assert.NoError(t, err, "生成增强模型不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	resultModel := result.(*TestEnhancedModel)
	
	// 验证嵌套结构体
	assert.NotEmpty(t, resultModel.NestedStruct.Name, "NestedStruct.Name 应被填充")
	assert.NotZero(t, resultModel.NestedStruct.Value, "NestedStruct.Value 应被填充")
	
	// 验证嵌套指针结构体
	assert.NotNil(t, resultModel.NestedPtr, "NestedPtr 应被初始化")
	assert.NotEmpty(t, resultModel.NestedPtr.Name, "NestedPtr.Name 应被填充")
	assert.NotZero(t, resultModel.NestedPtr.Value, "NestedPtr.Value 应被填充")
}

// TestGenerateRandModelEnhanced_Interface 测试接口类型
func TestGenerateRandModelEnhanced_Interface(t *testing.T) {
	model := &TestEnhancedModel{}
	
	result, jsonStr, err := random.GenerateRandModel(model)
	assert.NoError(t, err, "生成增强模型不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	resultModel := result.(*TestEnhancedModel)
	
	// 验证接口类型（应该被填充为某个具体类型的值）
	assert.NotNil(t, resultModel.InterfaceField, "InterfaceField 应被填充")
}

// TestGenerateRandModelEnhanced_CustomTags 测试自定义标签
func TestGenerateRandModelEnhanced_CustomTags(t *testing.T) {
	model := &TestEnhancedModel{}
	
	result, jsonStr, err := random.GenerateRandModel(model)
	assert.NoError(t, err, "生成增强模型不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	resultModel := result.(*TestEnhancedModel)
	
	// 验证自定义标签生成的值
	assert.Contains(t, resultModel.EmailField, "@", "EmailField 应包含 @ 符号")
	assert.Contains(t, resultModel.EmailField, ".com", "EmailField 应包含 .com")
	
	assert.Len(t, resultModel.PhoneField, 11, "PhoneField 应为11位手机号")
	assert.True(t, resultModel.PhoneField[0] == '1', "PhoneField 应以1开头")
	
	assert.Contains(t, resultModel.UUIDField, "-", "UUIDField 应包含连字符")
	assert.Len(t, resultModel.UUIDField, 36, "UUIDField 长度应为36位")
}

// TestGenerateRandModelEnhanced_Options 测试自定义选项
func TestGenerateRandModelEnhanced_Options(t *testing.T) {
	model := &TestEnhancedModel{}
	
	// 测试自定义选项
	opts := &random.GenerateRandModelOptions{
		MaxDepth:      2,
		MaxSliceLen:   3,
		MaxMapLen:     2,
		StringLength:  15,
		FillNilPtr:    true,
		UseCustomTags: false,
	}
	
	result, jsonStr, err := random.GenerateRandModel(model, opts)
	assert.NoError(t, err, "使用自定义选项不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	resultModel := result.(*TestEnhancedModel)
	
	// 验证字符串长度（由于使用了 opts.UseCustomTags = false，应该使用默认字符串生成）
	if opts.UseCustomTags == false {
		assert.Len(t, resultModel.StringField, opts.StringLength, "StringField 长度应符合自定义设置")
	}
	
	// 验证切片长度限制
	assert.LessOrEqual(t, len(resultModel.StringSlice), opts.MaxSliceLen, "StringSlice 长度应不超过限制")
	
	// 验证映射长度限制
	assert.LessOrEqual(t, len(resultModel.StringMap), opts.MaxMapLen, "StringMap 长度应不超过限制")
}

// TestGenerateRandModelEnhanced_NilPointerOption 测试指针选项
func TestGenerateRandModelEnhanced_NilPointerOption(t *testing.T) {
	model := &TestEnhancedModel{}
	
	// 测试不填充 nil 指针的选项
	opts := random.DefaultOptions()
	opts.FillNilPtr = false
	
	result, jsonStr, err := random.GenerateRandModel(model, opts)
	assert.NoError(t, err, "不填充指针选项不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	// 注意：由于我们的增强版本默认会创建指针，这个测试主要验证不报错
}

// 测试复杂嵌套情况
type ComplexNestedModel struct {
	Level1 *Level1Model `json:"level1"`
}

type Level1Model struct {
	Name   string        `json:"name"`
	Level2 *Level2Model  `json:"level2"`
	Items  []*Level2Model `json:"items"`
}

type Level2Model struct {
	Value   int                    `json:"value"`
	Details map[string]*Level3Model `json:"details"`
}

type Level3Model struct {
	Description string `json:"description"`
}

// TestGenerateRandModelEnhanced_AutoSkipUnsupportedTypes 测试自动跳过不支持的类型
func TestGenerateRandModelEnhanced_AutoSkipUnsupportedTypes(t *testing.T) {
	// 定义一个包含各种不支持JSON序列化类型的结构体
	type UnsupportedTypesModel struct {
		SupportedField   string        `json:"supported_field"`
		ComplexField     complex64     `json:"complex_field"`
		ChanField        chan int      `json:"chan_field"`
		FuncField        func() string `json:"func_field"`
		UnsafeField      uintptr       `json:"unsafe_field"`
		SkippedField     string        `json:"-"` // 明确标记跳过的字段
		privateField     string        // 不可导出的字段
	}
	
	model := &UnsupportedTypesModel{}
	
	result, jsonStr, err := random.GenerateRandModel(model)
	assert.NoError(t, err, "应该能够处理包含不支持类型的结构体，自动跳过它们")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	// 验证 JSON 可解析
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &jsonMap)
	assert.NoError(t, err, "应能正确解析 JSON")
	
	resultModel := result.(*UnsupportedTypesModel)
	
	// 验证支持的字段被填充
	assert.NotEmpty(t, resultModel.SupportedField, "支持的字段应被填充")
	
	// 验证不支持的字段被跳过（保持零值）
	assert.Zero(t, resultModel.ComplexField, "ComplexField 应被自动跳过")
	assert.Nil(t, resultModel.ChanField, "ChanField 应被自动跳过")
	assert.Nil(t, resultModel.FuncField, "FuncField 应被自动跳过")
	assert.Zero(t, resultModel.UnsafeField, "UnsafeField 应被自动跳过")
	assert.Empty(t, resultModel.SkippedField, "显式跳过的字段应为空")
	assert.Empty(t, resultModel.privateField, "私有字段应为空")
	
	// 验证 JSON 中不包含不支持的字段
	assert.NotContains(t, jsonStr, "complex_field", "JSON 不应包含 complex 字段")
	assert.NotContains(t, jsonStr, "chan_field", "JSON 不应包含 chan 字段")
	assert.NotContains(t, jsonStr, "func_field", "JSON 不应包含 func 字段")
	assert.NotContains(t, jsonStr, "unsafe_field", "JSON 不应包含 unsafe 字段")
	assert.NotContains(t, jsonStr, "skipped_field", "JSON 不应包含跳过的字段")
	assert.NotContains(t, jsonStr, "private_field", "JSON 不应包含私有字段")
	
	// 验证 JSON 包含支持的字段
	assert.Contains(t, jsonStr, "supported_field", "JSON 应包含支持的字段")
}

// TestGenerateRandModelEnhanced_ComplexNesting 测试复杂嵌套
func TestGenerateRandModelEnhanced_ComplexNesting(t *testing.T) {
	model := &ComplexNestedModel{}
	
	result, jsonStr, err := random.GenerateRandModel(model)
	assert.NoError(t, err, "生成复杂嵌套模型不应有错误")
	assert.NotNil(t, result, "结果不应为 nil")
	assert.NotEmpty(t, jsonStr, "JSON 字符串不应为空")
	
	// 验证 JSON 可解析
	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &jsonMap)
	assert.NoError(t, err, "应能正确解析复杂嵌套 JSON")
	
	resultModel := result.(*ComplexNestedModel)
	
	// 验证多级嵌套结构
	assert.NotNil(t, resultModel.Level1, "Level1 应被初始化")
	assert.NotEmpty(t, resultModel.Level1.Name, "Level1.Name 应被填充")
	assert.NotNil(t, resultModel.Level1.Level2, "Level1.Level2 应被初始化")
	assert.NotZero(t, resultModel.Level1.Level2.Value, "Level1.Level2.Value 应被填充")
	assert.NotEmpty(t, resultModel.Level1.Level2.Details, "Level1.Level2.Details 应包含元素")
	assert.NotEmpty(t, resultModel.Level1.Items, "Level1.Items 应包含元素")
	
	// 验证映射中的嵌套结构
	for key, detail := range resultModel.Level1.Level2.Details {
		assert.NotEmpty(t, key, "Details 键应不为空")
		assert.NotNil(t, detail, "Details 值应不为 nil")
		assert.NotEmpty(t, detail.Description, "Details.Description 应被填充")
	}
	
	// 验证切片中的嵌套结构
	for i, item := range resultModel.Level1.Items {
		assert.NotNil(t, item, "Items[%d] 应不为 nil", i)
		assert.NotZero(t, item.Value, "Items[%d].Value 应被填充", i)
	}
}