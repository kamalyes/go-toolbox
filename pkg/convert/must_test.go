/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\convert\must_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"encoding/json"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stretchr/testify/assert"
)

const (
	unexpectedErrorForInputMsg = "Unexpected error for input %v: %v"
	mustIntTResultMsg          = "MustIntT(%v) = %d; want %d"
)

func TestMustString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{"hello", "hello"},
		{[]byte("world"), "world"},
		{nil, ""},
		{true, "true"},
		{42, "42"},
		{3.14, "3.14"},
		{time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), "2024-01-01T12:00:00Z"},
	}

	for _, test := range tests {
		result := MustString(test.input)
		assert.Equal(t, test.expected, result, "MustString(%v) = %s; want %s", test.input, result, test.expected)
	}
}

func TestMustIntTErrorCases(t *testing.T) {
	// 定义一些不支持的类型进行测试
	unsupportedValues := []any{
		"string",         // 字符串
		true,             // 布尔值
		[]int{1, 2, 3},   // 切片
		map[string]int{}, // 映射
		nil,              // nil 值
	}

	for _, val := range unsupportedValues {
		result, err := MustIntT[int](val, nil) // 传递 nil 作为默认的取整模式
		assert.Error(t, err, "Expected an error for input %v, but got result %d", val, result)
	}
}

func TestMustIntTSuccessCases(t *testing.T) {
	tests := []struct {
		input    any
		expected int
	}{
		{int(10), 10},
		{int8(20), 20},
		{int16(30), 30},
		{int32(40), 40},
		{int64(50), 50},
		{uint(60), 60},
		{uint8(70), 70},
		{uint16(80), 80},
		{uint32(90), 90},
		{uint64(100), 100},
		{float32(3.7), 3}, // 测试浮点数向下取整
		{float64(4.9), 4}, // 测试浮点数向下取整
	}

	for _, test := range tests {
		result, err := MustIntT[int](test.input, nil) // 传递 nil 作为默认的取整模式
		assert.NoError(t, err, unexpectedErrorForInputMsg, test.input, err)
		assert.Equal(t, test.expected, result, mustIntTResultMsg, test.input, result, test.expected)
	}
}

func TestMustIntTConvertConvertRoundUp(t *testing.T) {
	tests := []struct {
		input    any
		expected int
	}{
		{float32(3.2), 4}, // 测试浮点数向上取整
		{float64(4.8), 5}, // 测试浮点数向上取整
	}

	for _, test := range tests {
		mode := RoundUp // 设置取整模式为向上取整
		result, err := MustIntT[int](test.input, &mode)
		assert.NoError(t, err, unexpectedErrorForInputMsg, test.input, err)
		assert.Equal(t, test.expected, result, mustIntTResultMsg, test.input, result, test.expected)
	}
}

func TestMustIntTRoundDown(t *testing.T) {
	tests := []struct {
		input    any
		expected int
	}{
		{float32(3.7), 3}, // 测试浮点数向下取整
		{float64(4.9), 4}, // 测试浮点数向下取整
	}

	for _, test := range tests {
		mode := RoundDown // 设置取整模式为向下取整
		result, err := MustIntT[int](test.input, &mode)
		assert.NoError(t, err, unexpectedErrorForInputMsg, test.input, err)
		assert.Equal(t, test.expected, result, mustIntTResultMsg, test.input, result, test.expected)
	}
}

func TestMustBool(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{"1", true},
		{"true", true},
		{"false", false},
		{0, false},
		{1, true},
		{nil, false},
		{true, true},
		{false, false},
	}

	for _, test := range tests {
		result := MustBool(test.input)
		assert.Equal(t, test.expected, result, "MustBool(%v) = %v; want %v", test.input, result, test.expected)
	}
}

func TestNumberSliceToStringSlice(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []string
	}{
		{
			input:    []uint64{1, 2, 3, 4, 5},
			expected: []string{"1", "2", "3", "4", "5"},
		},
		{
			input:    []int{10, 20, 30},
			expected: []string{"10", "20", "30"},
		},
		{
			input:    []float64{1.1, 2.2, 3.3},
			expected: []string{"1.1", "2.2", "3.3"},
		},
		{
			input:    []int64{-1, -2, -3},
			expected: []string{"-1", "-2", "-3"},
		},
		{
			input:    []uint{100, 200, 300},
			expected: []string{"100", "200", "300"},
		},
		{
			input:    []float32{1.5, 2.5, 3.5},
			expected: []string{"1.5", "2.5", "3.5"},
		},
	}

	for _, test := range tests {
		var result []string
		switch v := test.input.(type) {
		case []uint64:
			result = NumberSliceToStringSlice(v)
		case []int:
			result = NumberSliceToStringSlice(v)
		case []float64:
			result = NumberSliceToStringSlice(v)
		case []int64:
			result = NumberSliceToStringSlice(v)
		case []uint:
			result = NumberSliceToStringSlice(v)
		case []float32:
			result = NumberSliceToStringSlice(v)
		default:
			t.Fatalf("unsupported type: %T", v)
		}

		assert.Equal(t, test.expected, result, "Expected %v, got %v", test.expected, result)
	}
}

func TestStringSliceToFloat64SliceRoundUp(t *testing.T) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	mode := RoundUp
	expected := []float64{2.0, 3.0, 4.0, 4.0, 6.0}

	result, err := StringSliceToNumberSlice[float64](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToFloat64SliceRoundDown(t *testing.T) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	mode := RoundDown
	expected := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

	result, err := StringSliceToNumberSlice[float64](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToIntSlice(t *testing.T) {
	input := []string{"1", "2", "3"}
	expected := []int{1, 2, 3}

	result, err := StringSliceToNumberSlice[int](input, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToInt64Slice(t *testing.T) {
	input := []string{"1000", "2000", "3000"}
	expected := []int64{1000, 2000, 3000}

	result, err := StringSliceToNumberSlice[int64](input, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToNumberSliceInvalidInput(t *testing.T) {
	input := []string{"a", "b", "c"}
	mode := RoundDown

	result, err := StringSliceToNumberSlice[int](input, &mode)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestStringSliceToNumberSliceEmptySlice(t *testing.T) {
	input := []string{}
	expected := []int{}

	result, err := StringSliceToNumberSlice[int](input, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToNumberSliceNilSlice(t *testing.T) {
	var input []string
	expected := []int{}

	result, err := StringSliceToNumberSlice[int](input, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToIntSliceNegativeNumbersRoundUp(t *testing.T) {
	input := []string{"-1", "-2", "-3"}
	mode := RoundUp
	expected := []int{-1, -2, -3}

	result, err := StringSliceToNumberSlice[int](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToIntSliceNegativeNumbersRoundDown(t *testing.T) {
	input := []string{"-1.5", "-2.3", "-3.7"}
	mode := RoundDown
	expected := []int{-2, -3, -4}

	result, err := StringSliceToNumberSlice[int](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToInt64SliceLargeNumbers(t *testing.T) {
	input := []string{"9223372036854775807", "9223372036854775806"}
	_, err := StringSliceToNumberSlice[int64](input, nil)
	assert.NoError(t, err)
}

func TestStringSliceToFloat64SliceSmallNumbers(t *testing.T) {
	input := []string{"0.0001", "0.0002", "0.0003"}
	mode := RoundDown
	expected := []float64{0.0, 0.0, 0.0}

	result, err := StringSliceToNumberSlice[float64](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToFloat64SliceComplexInput(t *testing.T) {
	input := []string{"1", "2.5", "3.14", "4.0", "5.999"}
	mode := RoundUp
	expected := []float64{1.0, 3.0, 4.0, 4.0, 6.0}

	result, err := StringSliceToNumberSlice[float64](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"1.0", 1.0, false},
		{"-1.0", -1.0, false},
		{"0.0", 0.0, false},
		{"3.14", 3.14, false},
		{"-3.14", -3.14, false},
		{"1.5e2", 150.0, false},    // 科学计数法
		{"-1.5e-2", -0.015, false}, // 科学计数法
		{"abc", 0, true},           // 无法解析的字符串
		{"", 0, true},              // 空字符串
		{"NaN", 0, true},           // 非数值
		{"Infinity", 0, true},      // 无穷大
		{"-Infinity", 0, true},     // 负无穷大
	}

	for _, test := range tests {
		var result float64
		err := ParseFloat(test.input, &result)

		if test.hasError {
			assert.Error(t, err, "expected an error for input %q", test.input)
		} else {
			assert.NoError(t, err, "unexpected error for input %q", test.input)
			assert.Equal(t, test.expected, result, "expected %v for input %q, got %v", test.expected, test.input, result)
		}
	}
}

func TestParseFloatFloat32(t *testing.T) {
	tests := []struct {
		input    string
		expected float32
		hasError bool
	}{
		{"1.0", 1.0, false},
		{"-1.0", -1.0, false},
		{"0.0", 0.0, false},
		{"3.14", 3.14, false},
		{"-3.14", -3.14, false},
		{"1.5e2", 150.0, false},    // 科学计数法
		{"-1.5e-2", -0.015, false}, // 科学计数法
		{"abc", 0, true},           // 无法解析的字符串
		{"", 0, true},              // 空字符串
		{"NaN", 0, true},           // 非数值
		{"Infinity", 0, true},      // 无穷大
		{"-Infinity", 0, true},     // 负无穷大
	}

	for _, test := range tests {
		var result float32
		err := ParseFloat(test.input, &result)

		if test.hasError {
			assert.Error(t, err, "expected an error for input %q", test.input)
		} else {
			assert.NoError(t, err, "unexpected error for input %q", test.input)
			assert.Equal(t, test.expected, result, "expected %v for input %q, got %v", test.expected, test.input, result)
		}
	}
}

func TestMustFloatT(t *testing.T) {
	tests := []struct {
		input    any // 修改为 any 以支持不同类型
		mode     RoundMode
		expected float64
		hasError bool
	}{
		{"3.14", RoundNone, 3.14, false},
		{"3.14", RoundNearest, 3.0, false},
		{"3.6", RoundNearest, 4.0, false},
		{"3.1", RoundUp, 4.0, false},
		{"3.9", RoundDown, 3.0, false},
		{"-3.14", RoundNone, -3.14, false},
		{"-3.6", RoundNearest, -4.0, false},
		{"invalid", RoundNone, 0, true},          // 无效字符串
		{float64(5.5), RoundNone, 5.5, false},    // 测试 float64 类型
		{float32(5.5), RoundNone, 5.5, false},    // 测试 float32 类型
		{float32(5.5), RoundNearest, 6.0, false}, // 测试 float32 类型与取整
	}

	for _, test := range tests {
		result, err := MustFloatT[float64](test.input, test.mode)
		if test.hasError {
			assert.Error(t, err, "输入: %s, 预期错误: %v", test.input, test.hasError)
		} else {
			assert.NoError(t, err, "输入: %s, 实际错误: %v", test.input, err)
			assert.Equal(t, test.expected, result, "输入: %s, 预期: %v, 实际: %v", test.input, test.expected, result)
		}
	}
}

func TestStringSliceToFloatSlice(t *testing.T) {
	tests := []struct {
		input    []string
		mode     RoundMode
		expected []float64
		hasError bool
	}{
		// 测试 float64 类型
		{[]string{"3.14", "2.71"}, RoundNone, []float64{3.14, 2.71}, false},
		{[]string{"3.14", "2.71"}, RoundNearest, []float64{3.0, 3.0}, false},
		{[]string{"3.6", "3.1"}, RoundUp, []float64{4.0, 4.0}, false},
		{[]string{"3.9", "3.2"}, RoundDown, []float64{3.0, 3.0}, false},
		{nil, RoundNone, []float64{}, false},                               // 处理 nil 切片
		{[]string{"invalid"}, RoundNone, nil, true},                        // 包含无效字符串
		{[]string{"1.5", "2.5"}, RoundNearest, []float64{2.0, 3.0}, false}, // 更新四舍五入预期
	}

	for _, test := range tests {
		// 测试 float64
		result, err := StringSliceToFloatSlice[float64](test.input, test.mode)
		if test.hasError {
			assert.Error(t, err, "输入: %v, 预期错误: %v", test.input, test.hasError)
		} else {
			assert.NoError(t, err, "输入: %v, 实际错误: %v", test.input, err)
			assert.Equal(t, test.expected, result, "输入: %v, 预期: %v, 实际: %v", test.input, test.expected, result)
		}

		// 测试 float32
		result32, err32 := StringSliceToFloatSlice[float32](test.input, test.mode)
		if test.hasError {
			assert.Error(t, err32, "输入: %v, 预期错误: %v", test.input, test.hasError)
		} else {
			assert.NoError(t, err32, "输入: %v, 实际错误: %v", test.input, err32)

			// 将期望的 float64 转换为 float32
			expected32 := make([]float32, len(test.expected))
			for i, v := range test.expected {
				expected32[i] = float32(v) // 转换为 float32
			}
			assert.Equal(t, expected32, result32, "输入: %v, 预期: %v, 实际: %v", test.input, expected32, result32)
		}
	}
}

func TestStringSliceToInterfaceSliceTableDriven(t *testing.T) {
	tests := map[string]struct {
		input    []string
		expected []interface{}
	}{
		"normal": {
			input:    []string{"a", "b", "c"},
			expected: []interface{}{"a", "b", "c"},
		},
		"empty": {
			input:    []string{},
			expected: []interface{}{},
		},
		"empty string": {
			input:    []string{"", "x", "y"},
			expected: []interface{}{"", "x", "y"},
		},
		"spaces": {
			input:    []string{" ", "  ", "abc"},
			expected: []interface{}{" ", "  ", "abc"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := StringSliceToInterfaceSlice(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToNumberSliceTableDriven(t *testing.T) {
	tests := map[string]struct {
		input    interface{}
		expected []int
		wantErr  bool
	}{
		"string numbers": {
			input:    "1,2,3,4",
			expected: []int{1, 2, 3, 4},
			wantErr:  false,
		},
		"string slice numbers": {
			input:    []string{"10", "20", "30"},
			expected: []int{10, 20, 30},
			wantErr:  false,
		},
		"empty string": {
			input:    "",
			expected: []int{},
			wantErr:  false,
		},
		"string with spaces": {
			input:    " 1 , 2 , 3 ",
			expected: []int{1, 2, 3},
			wantErr:  false,
		},
		"negative numbers": {
			input:    "-1,-2,-3",
			expected: []int{-1, -2, -3},
			wantErr:  false,
		},
		"invalid number": {
			input:   "1,2,abc",
			wantErr: true,
		},
		"wrong type": {
			input:   123,
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ToNumberSlice[int](tc.input, ",")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestMustToNumberSliceTableDriven(t *testing.T) {
	tests := map[string]struct {
		input       string
		expected    []int
		expectPanic bool
	}{
		"normal": {
			input:    "5,6,7",
			expected: []int{5, 6, 7},
		},
		"with spaces": {
			input:    " 8 , 9 , 10 ",
			expected: []int{8, 9, 10},
		},
		"panic case": {
			input:       "1,2,xyz",
			expectPanic: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectPanic {
				assert.Panics(t, func() {
					MustToNumberSlice[int](tc.input, ",")
				})
			} else {
				assert.NotPanics(t, func() {
					got := MustToNumberSlice[int](tc.input, ",")
					assert.Equal(t, tc.expected, got)
				})
			}
		})
	}
}

// TestMustStringWithProtobufTimestamp 测试 protobuf Timestamp 转换
func TestMustStringWithProtobufTimestamp(t *testing.T) {
	// 创建一个 protobuf Timestamp (2024-01-01 12:00:00 UTC)
	pbTime := timestamppb.New(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))

	// 默认格式 (RFC3339)
	result := MustString(pbTime)
	assert.Equal(t, "2024-01-01T12:00:00Z", result)

	// 自定义格式
	customLayout := "2006-01-02 15:04:05"
	result2 := MustString(pbTime, customLayout)
	assert.Equal(t, "2024-01-01 12:00:00", result2)

	// 测试 nil protobuf Timestamp
	var nilPbTime *timestamppb.Timestamp
	result3 := MustString(nilPbTime)
	assert.Equal(t, "", result3) // nil 转换为 ""，而不是 "null"
}

// TestMustStringWithCustomTimeLayout 测试自定义时间格式
func TestMustStringWithCustomTimeLayout(t *testing.T) {
	testTime := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
	customLayout := "2006/01/02 15:04"

	result := MustString(testTime, customLayout)
	assert.Equal(t, "2024/03/15 14:30", result)
}

// TestConvertStructToString 测试结构体转换为字符串
func TestConvertStructToString(t *testing.T) {
	type CustomStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	cs := CustomStruct{Name: "test", Value: 42}
	result := MustString(cs)

	var decoded CustomStruct
	err := json.Unmarshal([]byte(result), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, cs, decoded)
}

// TestConvertPtrToString 测试指针转换为字符串
func TestConvertPtrToString(t *testing.T) {
	// 测试 *int
	intVal := 42
	result := MustString(&intVal)
	assert.Equal(t, "42", result)

	// 测试 *string
	strVal := "hello"
	result2 := MustString(&strVal)
	assert.Equal(t, "hello", result2)

	// 测试 *time.Time
	testTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result3 := MustString(&testTime)
	assert.Equal(t, "2024-01-01T00:00:00Z", result3)
}

// TestMustJSONIndent 测试 JSON 格式化
func TestMustJSONIndent(t *testing.T) {
	type TestStruct struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Active bool   `json:"active"`
	}

	data := TestStruct{
		Name:   "John",
		Age:    30,
		Active: true,
	}

	result, err := MustJSONIndent(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// 验证是否包含换行符和空格
	assert.Contains(t, string(result), "\n")
	assert.Contains(t, string(result), "  ")

	// 验证 JSON 是否有效
	var decoded TestStruct
	err = json.Unmarshal(result, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, data, decoded)
}

// TestMustJSON 测试 JSON 生成
func TestMustJSON(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name  string
		input interface{}
	}{
		{"simple struct", TestStruct{Name: "Alice", Age: 25}},
		{"map", map[string]interface{}{"key": "value", "number": 123}},
		{"slice", []int{1, 2, 3, 4, 5}},
		{"string", "test string"},
		{"number", 42},
		{"bool", true},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MustJSON(tt.input)
			assert.NoError(t, err)
			assert.NotEmpty(t, result)

			// 验证 JSON 是否有效
			var decoded interface{}
			err = json.Unmarshal(result, &decoded)
			assert.NoError(t, err)
		})
	}
}

// TestAnySliceToInterfaceSlice 测试任意切片转换为接口切片
func TestAnySliceToInterfaceSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []interface{}
	}{
		{"int slice", []int{1, 2, 3}, []interface{}{1, 2, 3}},
		{"string slice", []string{"a", "b", "c"}, []interface{}{"a", "b", "c"}},
		{"int array", [3]int{1, 2, 3}, []interface{}{1, 2, 3}},
		{"empty slice", []int{}, []interface{}{}},
		{"nil input", nil, []interface{}{}},
		{"non-slice type", 42, []interface{}{}},
		{"bool slice", []bool{true, false, true}, []interface{}{true, false, true}},
		{"float64 slice", []float64{1.1, 2.2, 3.3}, []interface{}{1.1, 2.2, 3.3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AnySliceToInterfaceSlice(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFloat64ToInt 测试浮点数转换为整数
func TestFloat64ToInt(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		mode      RoundMode
		expected  int
		expectErr bool
	}{
		{"round none", 3.7, RoundNone, 3, false},
		{"round up", 3.2, RoundUp, 4, false},
		{"round down", 3.9, RoundDown, 3, false},
		{"round nearest low", 3.4, RoundNearest, 3, false},
		{"round nearest high", 3.6, RoundNearest, 4, false},
		{"negative round down", -3.9, RoundDown, -4, false},
		{"negative round up", -3.2, RoundUp, -3, false},
		{"zero", 0.0, RoundNone, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Float64ToInt[int](tt.value, tt.mode)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestMustIntTWithRoundNearest 测试 MustIntT 函数的四舍五入
func TestMustIntTWithRoundNearest(t *testing.T) {
	tests := []struct {
		input    float64
		expected int
	}{
		{2.4, 2},
		{2.5, 3},
		{2.6, 3},
		{-2.4, -2},
		{-2.5, -3},
		{-2.6, -3},
	}

	mode := RoundNearest
	for _, test := range tests {
		result, err := MustIntT[int](test.input, &mode)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result, "Input: %v", test.input)
	}
}

// TestNormalizeToStringSlice 测试规范化为字符串切片
func TestNormalizeToStringSlice(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		separator string
		expected  []string
		expectErr bool
	}{
		{"string split", "a,b,c", ",", []string{"a", "b", "c"}, false},
		{"empty string", "", ",", []string{}, false},
		{"string slice", []string{"x", "y", "z"}, ",", []string{"x", "y", "z"}, false},
		{"unsupported type", 123, ",", nil, true},
		{"string with different separator", "a;b;c", ";", []string{"a", "b", "c"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeToStringSlice(tt.input, tt.separator)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
