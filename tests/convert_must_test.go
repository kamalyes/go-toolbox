/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-05-27 18:05:57
 * @FilePath: \go-toolbox\tests\convert_must_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/stretchr/testify/assert"
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
		result := convert.MustString(test.input)
		assert.Equal(t, test.expected, result, "convert.MustString(%v) = %s; want %s", test.input, result, test.expected)
	}
}

func TestMustIntT_ErrorCases(t *testing.T) {
	// 定义一些不支持的类型进行测试
	unsupportedValues := []any{
		"string",         // 字符串
		true,             // 布尔值
		[]int{1, 2, 3},   // 切片
		map[string]int{}, // 映射
		nil,              // nil 值
	}

	for _, val := range unsupportedValues {
		result, err := convert.MustIntT[int](val, nil) // 传递 nil 作为默认的取整模式
		assert.Error(t, err, "Expected an error for input %v, but got result %d", val, result)
	}
}

func TestMustIntT_SuccessCases(t *testing.T) {
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
		result, err := convert.MustIntT[int](test.input, nil) // 传递 nil 作为默认的取整模式
		assert.NoError(t, err, "Unexpected error for input %v: %v", test.input, err)
		assert.Equal(t, test.expected, result, "MustIntT(%v) = %d; want %d", test.input, result, test.expected)
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
		mode := convert.RoundUp // 设置取整模式为向上取整
		result, err := convert.MustIntT[int](test.input, &mode)
		assert.NoError(t, err, "Unexpected error for input %v: %v", test.input, err)
		assert.Equal(t, test.expected, result, "MustIntT(%v) = %d; want %d", test.input, result, test.expected)
	}
}

func TestMustIntT_RoundDown(t *testing.T) {
	tests := []struct {
		input    any
		expected int
	}{
		{float32(3.7), 3}, // 测试浮点数向下取整
		{float64(4.9), 4}, // 测试浮点数向下取整
	}

	for _, test := range tests {
		mode := convert.RoundDown // 设置取整模式为向下取整
		result, err := convert.MustIntT[int](test.input, &mode)
		assert.NoError(t, err, "Unexpected error for input %v: %v", test.input, err)
		assert.Equal(t, test.expected, result, "MustIntT(%v) = %d; want %d", test.input, result, test.expected)
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
		result := convert.MustBool(test.input)
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
			result = convert.NumberSliceToStringSlice(v)
		case []int:
			result = convert.NumberSliceToStringSlice(v)
		case []float64:
			result = convert.NumberSliceToStringSlice(v)
		case []int64:
			result = convert.NumberSliceToStringSlice(v)
		case []uint:
			result = convert.NumberSliceToStringSlice(v)
		case []float32:
			result = convert.NumberSliceToStringSlice(v)
		default:
			t.Fatalf("unsupported type: %T", v)
		}

		assert.Equal(t, test.expected, result, "Expected %v, got %v", test.expected, result)
	}
}

func TestStringSliceToFloat64Slice_RoundUp(t *testing.T) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	mode := convert.RoundUp
	expected := []float64{2.0, 3.0, 4.0, 4.0, 6.0}

	result, err := convert.StringSliceToNumberSlice[float64](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToFloat64Slice_RoundDown(t *testing.T) {
	input := []string{"1.5", "2.3", "3.7", "4.0", "5.9"}
	mode := convert.RoundDown
	expected := []float64{1.0, 2.0, 3.0, 4.0, 5.0}

	result, err := convert.StringSliceToNumberSlice[float64](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToIntSlice(t *testing.T) {
	input := []string{"1", "2", "3"}
	expected := []int{1, 2, 3}

	result, err := convert.StringSliceToNumberSlice[int](input, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToInt64Slice(t *testing.T) {
	input := []string{"1000", "2000", "3000"}
	expected := []int64{1000, 2000, 3000}

	result, err := convert.StringSliceToNumberSlice[int64](input, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToNumberSlice_InvalidInput(t *testing.T) {
	input := []string{"a", "b", "c"}
	mode := convert.RoundDown

	result, err := convert.StringSliceToNumberSlice[int](input, &mode)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestStringSliceToNumberSlice_EmptySlice(t *testing.T) {
	input := []string{}
	expected := []int{}

	result, err := convert.StringSliceToNumberSlice[int](input, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToNumberSlice_NilSlice(t *testing.T) {
	var input []string
	expected := []int{}

	result, err := convert.StringSliceToNumberSlice[int](input, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToIntSlice_NegativeNumbers_RoundUp(t *testing.T) {
	input := []string{"-1", "-2", "-3"}
	mode := convert.RoundUp
	expected := []int{-1, -2, -3}

	result, err := convert.StringSliceToNumberSlice[int](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToIntSlice_NegativeNumbers_RoundDown(t *testing.T) {
	input := []string{"-1.5", "-2.3", "-3.7"}
	mode := convert.RoundDown
	expected := []int{-2, -3, -4}

	result, err := convert.StringSliceToNumberSlice[int](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToInt64Slice_LargeNumbers(t *testing.T) {
	input := []string{"9223372036854775807", "9223372036854775806"}
	_, err := convert.StringSliceToNumberSlice[int64](input, nil)
	assert.NoError(t, err)
}

func TestStringSliceToFloat64Slice_SmallNumbers(t *testing.T) {
	input := []string{"0.0001", "0.0002", "0.0003"}
	mode := convert.RoundDown
	expected := []float64{0.0, 0.0, 0.0}

	result, err := convert.StringSliceToNumberSlice[float64](input, &mode)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestStringSliceToFloat64Slice_ComplexInput(t *testing.T) {
	input := []string{"1", "2.5", "3.14", "4.0", "5.999"}
	mode := convert.RoundUp
	expected := []float64{1.0, 3.0, 4.0, 4.0, 6.0}

	result, err := convert.StringSliceToNumberSlice[float64](input, &mode)
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
		err := convert.ParseFloat(test.input, &result)

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
		err := convert.ParseFloat(test.input, &result)

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
		mode     convert.RoundMode
		expected float64
		hasError bool
	}{
		{"3.14", convert.RoundNone, 3.14, false},
		{"3.14", convert.RoundNearest, 3.0, false},
		{"3.6", convert.RoundNearest, 4.0, false},
		{"3.1", convert.RoundUp, 4.0, false},
		{"3.9", convert.RoundDown, 3.0, false},
		{"-3.14", convert.RoundNone, -3.14, false},
		{"-3.6", convert.RoundNearest, -4.0, false},
		{"invalid", convert.RoundNone, 0, true},          // 无效字符串
		{float64(5.5), convert.RoundNone, 5.5, false},    // 测试 float64 类型
		{float32(5.5), convert.RoundNone, 5.5, false},    // 测试 float32 类型
		{float32(5.5), convert.RoundNearest, 6.0, false}, // 测试 float32 类型与取整
	}

	for _, test := range tests {
		result, err := convert.MustFloatT[float64](test.input, test.mode)
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
		mode     convert.RoundMode
		expected []float64
		hasError bool
	}{
		// 测试 float64 类型
		{[]string{"3.14", "2.71"}, convert.RoundNone, []float64{3.14, 2.71}, false},
		{[]string{"3.14", "2.71"}, convert.RoundNearest, []float64{3.0, 3.0}, false},
		{[]string{"3.6", "3.1"}, convert.RoundUp, []float64{4.0, 4.0}, false},
		{[]string{"3.9", "3.2"}, convert.RoundDown, []float64{3.0, 3.0}, false},
		{nil, convert.RoundNone, []float64{}, false},                               // 处理 nil 切片
		{[]string{"invalid"}, convert.RoundNone, nil, true},                        // 包含无效字符串
		{[]string{"1.5", "2.5"}, convert.RoundNearest, []float64{2.0, 3.0}, false}, // 更新四舍五入预期
	}

	for _, test := range tests {
		// 测试 float64
		result, err := convert.StringSliceToFloatSlice[float64](test.input, test.mode)
		if test.hasError {
			assert.Error(t, err, "输入: %v, 预期错误: %v", test.input, test.hasError)
		} else {
			assert.NoError(t, err, "输入: %v, 实际错误: %v", test.input, err)
			assert.Equal(t, test.expected, result, "输入: %v, 预期: %v, 实际: %v", test.input, test.expected, result)
		}

		// 测试 float32
		result32, err32 := convert.StringSliceToFloatSlice[float32](test.input, test.mode)
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

func TestStringSliceToInterfaceSlice_TableDriven(t *testing.T) {
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
			result := convert.StringSliceToInterfaceSlice(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToNumberSlice_TableDriven(t *testing.T) {
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
			result, err := convert.ToNumberSlice[int](tc.input, ",")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestMustToNumberSlice_TableDriven(t *testing.T) {
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
					convert.MustToNumberSlice[int](tc.input, ",")
				})
			} else {
				assert.NotPanics(t, func() {
					got := convert.MustToNumberSlice[int](tc.input, ",")
					assert.Equal(t, tc.expected, got)
				})
			}
		})
	}
}
