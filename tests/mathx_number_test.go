/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-12 15:10:05
 * @FilePath: \go-toolbox\tests\mathx_number_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/stretchr/testify/assert"
)

// Decimals 测试
func TestDecimals(t *testing.T) {
	tests := []struct {
		name  string
		num   float64
		digit int
		want  string
	}{
		{"PositiveInteger", 12345, 2, "123.45"},
		{"PositiveFloat", 12345.6789, 3, "12.346"},
		{"NegativeInteger", -12345, 2, "-123.45"},
		{"NegativeFloat", -12345.6789, 4, "-1.2346"},
		{"Zero", 0, 3, "0.000"},
		{"SmallNumber", 0.12345, 4, "0.0000"},
		{"LargeNumber", 123456789, 5, "1234.56789"},
		{"NegativeSmallNumber", -0.12345, 4, "-0.0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.Decimals(tt.num, tt.digit)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMathxAtLeast(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, lower, expected int
	}{
		{5, 3, 3},
		{2, 10, 2},
		{-1, -5, -5},
	}
	for _, tt := range tests {
		result := mathx.AtLeast(tt.x, tt.lower)
		assert.Equal(tt.expected, result)
	}
}

func TestMathxAtMost(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, upper, expected int
	}{
		{8, 6, 8},
		{4, 2, 4},
		{-3, -7, -3},
	}
	for _, tt := range tests {
		result := mathx.AtMost(tt.x, tt.upper)
		assert.Equal(tt.expected, result)
	}
}

func TestMathxBetween(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, lower, upper, expected int
	}{
		{4, 1, 10, 4},
		{0, -5, 5, 0},
		{12, 1, 10, 10},
		{-8, -5, -3, -5},
	}
	for _, tt := range tests {
		result := mathx.Between(tt.x, tt.lower, tt.upper)
		assert.Equal(tt.expected, result)
	}
}

func TestAtLeastFloat64(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, lower, expected float64
	}{
		{5.5, 3.3, 3.3},
		{2.2, 10.0, 2.2},
		{-1.1, -5.5, -5.5},
	}
	for _, tt := range tests {
		result := mathx.AtLeast(tt.x, tt.lower)
		assert.Equal(tt.expected, result)
	}
}

func TestBetweenFloat64(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		x, lower, upper, expected float64
	}{
		{4.4, 1.1, 10.0, 4.4},
		{0.0, -5.5, 5.5, 0.0},
		{12.3, 1.0, 10.0, 10.0},
		{-8.8, -5.5, -3.3, -5.5},
	}
	for _, tt := range tests {
		result := mathx.Between(tt.x, tt.lower, tt.upper)
		assert.Equal(tt.expected, result)
	}
}

func TestZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		expected any
	}{
		{"int", 0},
		{"int8", int8(0)},
		{"int16", int16(0)},
		{"int32", int32(0)},
		{"int64", int64(0)},
		{"uint", uint(0)},
		{"uint8", uint8(0)},
		{"uint16", uint16(0)},
		{"uint32", uint32(0)},
		{"uint64", uint64(0)},
		{"float32", float32(0.0)},
		{"float64", float64(0.0)},
		{"string", ""},
		{"bool", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result any
			switch tt.expected.(type) {
			case int:
				result = mathx.ZeroValue[int]()
			case int8:
				result = mathx.ZeroValue[int8]()
			case int16:
				result = mathx.ZeroValue[int16]()
			case int32:
				result = mathx.ZeroValue[int32]()
			case int64:
				result = mathx.ZeroValue[int64]()
			case uint:
				result = mathx.ZeroValue[uint]()
			case uint8:
				result = mathx.ZeroValue[uint8]()
			case uint16:
				result = mathx.ZeroValue[uint16]()
			case uint32:
				result = mathx.ZeroValue[uint32]()
			case uint64:
				result = mathx.ZeroValue[uint64]()
			case float32:
				result = mathx.ZeroValue[float32]()
			case float64:
				result = mathx.ZeroValue[float64]()
			case string:
				result = mathx.ZeroValue[string]()
			case bool:
				result = mathx.ZeroValue[bool]()
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEqualIntSlices(t *testing.T) {
	tests := []struct {
		a        []int
		b        []int
		expected bool
	}{
		{[]int{1, 2, 3}, []int{1, 2, 3}, true},
		{[]int{1, 2, 3}, []int{3, 2, 1}, false},
		{[]int{1, 2, 3}, []int{1, 2}, false},
		{[]int{}, []int{}, true},
		{nil, nil, true},
		{[]int{1, 2, 3}, nil, false},
		{nil, []int{1, 2, 3}, false},
	}

	for _, test := range tests {
		result := mathx.EqualSlices(test.a, test.b)
		assert.Equal(t, test.expected, result, "EqualSlices(%v, %v) = %v; expected %v", test.a, test.b, result, test.expected)
	}
}

func TestEqualFloatSlices(t *testing.T) {
	tests := []struct {
		a        []float64
		b        []float64
		expected bool
	}{
		{[]float64{1.1, 2.2, 3.3}, []float64{1.1, 2.2, 3.3}, true},
		{[]float64{1.1, 2.2, 3.3}, []float64{3.3, 2.2, 1.1}, false},
		{[]float64{1.1, 2.2, 3.3}, []float64{1.1, 2.2}, false},
		{[]float64{}, []float64{}, true},
		{nil, nil, true},
		{[]float64{1.1, 2.2, 3.3}, nil, false},
		{nil, []float64{1.1, 2.2, 3.3}, false},
	}

	for _, test := range tests {
		result := mathx.EqualSlices(test.a, test.b)
		assert.Equal(t, test.expected, result, "EqualSlices(%v, %v) = %v; expected %v", test.a, test.b, result, test.expected)
	}
}

func TestEqualStringSlices(t *testing.T) {
	tests := []struct {
		a        []string
		b        []string
		expected bool
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, true},
		{[]string{"a", "b", "c"}, []string{"c", "b", "a"}, false},
		{[]string{"a", "b", "c"}, []string{"a", "b"}, false},
		{[]string{}, []string{}, true},
		{nil, nil, true},
		{[]string{"a", "b"}, nil, false},
		{nil, []string{"a", "b"}, false},
	}

	for _, test := range tests {
		result := mathx.EqualSlices(test.a, test.b)
		assert.Equal(t, test.expected, result, "EqualSlices(%v, %v) = %v; expected %v", test.a, test.b, result, test.expected)
	}
}

func TestEqualBoolSlices(t *testing.T) {
	tests := []struct {
		a        []bool
		b        []bool
		expected bool
	}{
		{[]bool{true, false, true}, []bool{true, false, true}, true},
		{[]bool{true, false, true}, []bool{false, true, true}, false},
		{[]bool{true, false}, []bool{true, false, true}, false},
		{[]bool{}, []bool{}, true},
		{nil, nil, true},
		{[]bool{true, false}, nil, false},
		{nil, []bool{true, false}, false},
	}

	for _, test := range tests {
		result := mathx.EqualSlices(test.a, test.b)
		assert.Equal(t, test.expected, result, "EqualSlices(%v, %v) = %v; expected %v", test.a, test.b, result, test.expected)
	}
}

func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected int
	}{
		{"flower", "flow", 4},         // 正常情况
		{"dog", "racecar", 0},         // 无公共前缀
		{"", "", 0},                   // 空字符串
		{"", "abc", 0},                // 一个空字符串
		{"abc", "", 0},                // 一个空字符串
		{"abcde", "abcfg", 3},         // 部分公共前缀
		{"prefix", "prefixsuffix", 6}, // 完全相同的前缀
		{"abcdef", "abcxyz", 3},       // 部分公共前缀
		{"same", "same", 4},           // 完全相同的字符串
		{"abc", "abcd", 3},            // 另一个边界情况
	}

	for _, tt := range tests {
		t.Run(tt.a+"_"+tt.b, func(t *testing.T) {
			got := mathx.LongestCommonPrefix(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("LongestCommonPrefix(%q, %q) = %d; want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestParseIntOrName(t *testing.T) {
	names := map[string]uint{
		"jan": 1,
		"feb": 2,
		"mar": 3,
		"dec": 12,
	}

	tests := []struct {
		input    string
		expected uint
		wantErr  bool
	}{
		{"jan", 1, false},
		{"FEB", 2, false},
		{"10", 10, false},
		{"0", 0, false},
		{"abc", 0, true},
	}

	for _, tt := range tests {
		got, err := mathx.ParseIntOrName(tt.input, names)
		if tt.wantErr {
			assert.Error(t, err, "ParseIntOrName(%q) should return error", tt.input)
		} else {
			assert.NoError(t, err, "ParseIntOrName(%q) unexpected error", tt.input)
			assert.Equal(t, tt.expected, got, "ParseIntOrName(%q) got wrong result", tt.input)
		}
	}
}

// 自定义结构体用于测试
type Person struct {
	Name string
	Age  int
}

func TestSafeGetIndexWithErr(t *testing.T) {
	// 字符串切片
	strSlice := []string{"apple", "banana", "cherry"}
	valStr, err := mathx.SafeGetIndexWithErr(strSlice, 1)
	assert.NoError(t, err)
	assert.Equal(t, "banana", valStr)

	// 整型切片
	intSlice := []int{10, 20, 30}
	valInt, err := mathx.SafeGetIndexWithErr(intSlice, 2)
	assert.NoError(t, err)
	assert.Equal(t, 30, valInt)

	// 结构体切片
	personSlice := []Person{
		{"Alice", 30},
		{"Bob", 25},
	}
	valPerson, err := mathx.SafeGetIndexWithErr(personSlice, 0)
	assert.NoError(t, err)
	assert.Equal(t, Person{"Alice", 30}, valPerson)

	// 指针切片
	ptrSlice := []*Person{
		{"Charlie", 40},
		nil,
	}
	valPtr, err := mathx.SafeGetIndexWithErr(ptrSlice, 1)
	assert.NoError(t, err)
	assert.Nil(t, valPtr) // 索引1是nil指针

	// 索引越界测试
	_, err = mathx.SafeGetIndexWithErr(strSlice, 5)
	assert.Error(t, err)

	_, err = mathx.SafeGetIndexWithErr(intSlice, -1)
	assert.Error(t, err)
}

func TestSafeGetIndexOrDefault(t *testing.T) {
	// 字符串切片，索引合法
	strSlice := []string{"apple", "banana", "cherry"}
	valStr := mathx.SafeGetIndexOrDefault(strSlice, 1, "default")
	assert.Equal(t, "banana", valStr)

	// 字符串切片，索引越界
	valStr = mathx.SafeGetIndexOrDefault(strSlice, 5, "default")
	assert.Equal(t, "default", valStr)

	// 整型切片，索引合法
	intSlice := []int{10, 20, 30}
	valInt := mathx.SafeGetIndexOrDefault(intSlice, 2, -1)
	assert.Equal(t, 30, valInt)

	// 整型切片，索引越界
	valInt = mathx.SafeGetIndexOrDefault(intSlice, -1, -1)
	assert.Equal(t, -1, valInt)

	// 结构体切片，索引合法
	personSlice := []Person{
		{"Alice", 30},
		{"Bob", 25},
	}
	valPerson := mathx.SafeGetIndexOrDefault(personSlice, 0, Person{"Default", 0})
	assert.Equal(t, Person{"Alice", 30}, valPerson)

	// 结构体切片，索引越界
	valPerson = mathx.SafeGetIndexOrDefault(personSlice, 5, Person{"Default", 0})
	assert.Equal(t, Person{"Default", 0}, valPerson)

	// 指针切片，索引合法且元素为nil指针
	ptrSlice := []*Person{
		{"Charlie", 40},
		nil,
	}
	valPtr := mathx.SafeGetIndexOrDefault(ptrSlice, 1, nil)
	assert.Nil(t, valPtr)

	// 指针切片，索引越界
	valPtr = mathx.SafeGetIndexOrDefault(ptrSlice, 10, nil)
	assert.Nil(t, valPtr)
}

func TestSafeGetIndexOrDefaultNoSpace(t *testing.T) {
	tests := []struct {
		slice      []string
		index      int
		defaultVal string
		want       string
	}{
		{[]string{"a b c", " d e f ", "g h i"}, 0, "default", "abc"},
		{[]string{"a b c", " d e f ", "g h i"}, 1, "default", "def"},
		{[]string{"a b c", " d e f ", "g h i"}, 2, "default", "ghi"},
		{[]string{"a b c", " d e f ", "g h i"}, 3, "default", "default"}, // 越界返回默认值
		{[]string{}, 0, "default", "default"},                            // 空切片返回默认值
		{[]string{" no space "}, 0, "default", "nospace"},
	}

	for _, tt := range tests {
		got := mathx.SafeGetIndexOrDefaultNoSpace(tt.slice, tt.index, tt.defaultVal)
		if got != tt.want {
			t.Errorf("SafeGetIndexOrDefaultNoSpace(%v, %d, %q) = %q; want %q",
				tt.slice, tt.index, tt.defaultVal, got, tt.want)
		}
	}
}
