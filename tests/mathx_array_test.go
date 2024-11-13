/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 15:55:06
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 13:20:15
 * @FilePath: \go-toolbox\tests\mathx_array_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/stretchr/testify/assert"
)

// TestArrayMinMax 测试 ArrayMinMax 函数
func TestArrayMinMax(t *testing.T) {
	tests := []struct {
		name      string
		list      []int
		wantMin   int
		wantMax   int
		expectErr bool
	}{
		{"MinMaxPositive", []int{1, 2, 3, 4, 5}, 1, 5, false},
		{"MinMaxNegative", []int{-1, -2, -3, -4, -5}, -5, -1, false},
		{"MinMaxMixed", []int{-1, 2, -3, 4, -5}, -5, 4, false},
		{"MinMaxEmpty", []int{}, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var min, max int
			var err error

			if len(tt.list) == 0 {
				err = errors.New("list is empty")
			} else {
				min = tt.list[0]
				max = tt.list[0]

				for _, v := range tt.list {
					min = mathx.AtMost(min, v)
					max = mathx.AtLeast(max, v)
				}
			}

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMin, min)
				assert.Equal(t, tt.wantMax, max)
			}
		})
	}
}

// TestArrayUnion 测试 ArrayUnion 函数
func TestArrayUnion(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{"UnionPositive", []int{1, 2, 3}, []int{2, 3, 4}, []int{1, 2, 3, 4}},
		{"UnionEmptyA", []int{}, []int{2, 3}, []int{2, 3}},
		{"UnionEmptyB", []int{1, 2, 3}, []int{}, []int{1, 2, 3}},
		{"UnionNoCommon", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.ArrayUnion(tt.a, tt.b)
			assert.ElementsMatch(t, tt.want, result)
		})
	}
}

// 测试 ArrayFisherYates 函数
func TestArrayFisherYates(t *testing.T) {

	tests := [][]int{
		{1, 2, 3, 4, 5},
		{10, 20, 30, 40, 50},
	}

	for _, input := range tests {
		original := make([]int, len(input))
		copy(original, input) // 复制输入以便后续比较

		mathx.ArrayFisherYates(input) // 调用洗牌函数

		// 使用 assert 检查数组是否被打乱
		assert.NotEqual(t, original, input, "ArrayFisherYates did not shuffle the array: original %v, got %v", original, input)
	}
}

// TestArrayContains 测试 ArrayContains 函数
func TestArrayContains(t *testing.T) {
	tests := []struct {
		name     string
		array    []int
		element  int
		expected bool
	}{
		{"ContainsTrue", []int{1, 2, 3}, 2, true},
		{"ContainsFalse", []int{1, 2, 3}, 4, false},
		{"ContainsEmpty", []int{}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.ArrayContains(tt.array, tt.element)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestArrayHasDuplicates 测试 ArrayHasDuplicates 函数
func TestArrayHasDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		array    []int
		expected bool
	}{
		{"HasDuplicatesTrue", []int{1, 2, 2, 3}, true},
		{"HasDuplicatesFalse", []int{1, 2, 3}, false},
		{"HasDuplicatesEmpty", []int{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.ArrayHasDuplicates(tt.array)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestArrayRemoveEmpty 测试 ArrayRemoveEmpty 函数
func TestArrayRemoveEmpty(t *testing.T) {
	tests := []struct {
		name     string
		array    []interface{}
		expected []interface{}
	}{
		{"RemoveEmpty", []interface{}{1, "", nil, 2}, []interface{}{1, 2}},
		{"RemoveAllEmpty", []interface{}{nil, "", nil}, []interface{}{}},
		{"RemoveNoEmpty", []interface{}{1, 2, 3}, []interface{}{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.ArrayRemoveEmpty(tt.array)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

// TestArrayRemoveDuplicates 测试 ArrayRemoveDuplicates 函数
func TestArrayRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		array    []int
		expected []int
	}{
		{"RemoveDuplicates", []int{1, 2, 2, 3}, []int{1, 2, 3}},
		{"RemoveNoDuplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"RemoveAllDuplicates", []int{1, 1, 1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.ArrayRemoveDuplicates(tt.array)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

// TestArrayRemoveZero 测试 ArrayRemoveZero 函数
func TestArrayRemoveZero(t *testing.T) {
	tests := []struct {
		name     string
		array    []int
		expected []int
	}{
		{"RemoveZeros", []int{0, 1, 2, 0, 3}, []int{1, 2, 3}},
		{"NoZeros", []int{1, 2, 3}, []int{1, 2, 3}},
		{"AllZeros", []int{0, 0, 0}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.ArrayRemoveZero(tt.array)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

// TestArrayChunk 测试 ArrayChunk 函数
func TestArrayChunk(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		size     int
		expected [][]int
	}{
		{"ChunkSize2", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{"ChunkSize3", []int{1, 2, 3, 4, 5}, 3, [][]int{{1, 2, 3}, {4, 5}}},
		{"ChunkSize0", []int{1, 2, 3}, 0, nil},
		{"ChunkSizeNegative", []int{1, 2, 3}, -1, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.ArrayChunk(tt.slice, tt.size)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestArrayDiffSetStrings(t *testing.T) {
	cases := []struct {
		arr1 []string
		arr2 []string
		want []interface{}
	}{
		{[]string{"a", "b", "c"}, []string{"b", "c", "d"}, []interface{}{"a", "d"}},
		{[]string{}, []string{"b", "c", "d"}, []interface{}{"b", "c", "d"}},
		{[]string{"a", "b", "c"}, []string{}, []interface{}{"a", "b", "c"}},
		{[]string{"apple", "banana"}, []string{"banana", "cherry"}, []interface{}{"apple", "cherry"}},
	}

	for _, tc := range cases {
		result := mathx.ArrayDiffSetSorted(tc.arr1, tc.arr2)
		assert.ElementsMatch(t, tc.want, result, "ArrayDiffSet(%v, %v) = %v; want %v", tc.arr1, tc.arr2, result, tc.want)
	}
}

func TestArrayDiffSetInts(t *testing.T) {
	cases := []struct {
		arr1 []int
		arr2 []int
		want []interface{}
	}{
		{[]int{1, 2, 3}, []int{2, 3, 4}, []interface{}{1, 4}},
		{[]int{}, []int{2, 3, 4}, []interface{}{2, 3, 4}},
		{[]int{1, 2, 3}, []int{}, []interface{}{1, 2, 3}},
		{[]int{1, 2, 3, 3}, []int{2, 3, 4}, []interface{}{1, 4}},
	}

	for _, tc := range cases {
		result := mathx.ArrayDiffSetSorted(tc.arr1, tc.arr2)
		assert.ElementsMatch(t, tc.want, result, "ArrayDiffSet(%v, %v) = %v; want %v", tc.arr1, tc.arr2, result, tc.want)
	}
}

func TestArrayDiffSetFloats(t *testing.T) {
	cases := []struct {
		arr1 []float64
		arr2 []float64
		want []interface{}
	}{
		{[]float64{1.1, 2.2, 3.3}, []float64{2.2, 3.3, 4.4}, []interface{}{1.1, 4.4}},
		{[]float64{}, []float64{2.2, 3.3}, []interface{}{2.2, 3.3}},
		{[]float64{1.1, 2.2, 3.3}, []float64{}, []interface{}{1.1, 2.2, 3.3}},
	}

	for _, tc := range cases {
		result := mathx.ArrayDiffSetSorted(tc.arr1, tc.arr2)
		assert.ElementsMatch(t, tc.want, result, "ArrayDiffSet(%v, %v) = %v; want %v", tc.arr1, tc.arr2, result, tc.want)
	}
}
