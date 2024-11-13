/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 23:32:32
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
		{5, 3, 5},
		{2, 10, 10},
		{-1, -5, -1},
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
		{8, 6, 6},
		{4, 2, 2},
		{-3, -7, -7},
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
		{5.5, 3.3, 5.5},
		{2.2, 10.0, 10.0},
		{-1.1, -5.5, -1.1},
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
