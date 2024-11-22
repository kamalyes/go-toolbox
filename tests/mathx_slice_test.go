/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 15:55:06
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 11:09:17
 * @FilePath: \go-toolbox\tests\mathx_slice_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/types"
	"github.com/stretchr/testify/assert"
)

// TestSliceMinMax 测试 SliceMinMax 函数
func TestSliceMinMax(t *testing.T) {
	tests := []struct {
		name      string
		list      []int
		f         types.MinMaxFunc[int]
		expected  int
		expectErr bool
	}{
		{"Empty list", []int{}, mathx.MinFunc[int], 0, true},
		{"Single element", []int{5}, mathx.MinFunc[int], 5, false},
		{"Multiple elements - Min", []int{3, 1, 4, 1, 5, 9}, mathx.MinFunc[int], 1, false},
		{"Multiple elements - Max", []int{3, 1, 4, 1, 5, 9}, mathx.MaxFunc[int], 9, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := mathx.SliceMinMax(test.list, test.f)
			if test.expectErr {
				assert.Error(t, err) // 断言期望错误
			} else {
				assert.NoError(t, err)                 // 断言没有错误
				assert.Equal(t, test.expected, result) // 断言结果相等
			}
		})
	}
}

func TestSliceAtMostAtLeast(t *testing.T) {
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

// TestSliceUnion 测试 SliceUnion 函数
func TestSliceUnion(t *testing.T) {
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
			result := mathx.SliceUnion(tt.a, tt.b)
			assert.ElementsMatch(t, tt.want, result)
		})
	}
}

func TestSliceEqual(t *testing.T) {
	// 测试整数切片
	intSlice1 := []int{1, 2, 3, 4, 5}
	intSlice2 := []int{1, 2, 3, 4, 5}
	intSlice3 := []int{1, 2, 3, 4, 6}

	assert.True(t, mathx.SliceEqual(intSlice1, intSlice2), "Expected slices to be equal")
	assert.False(t, mathx.SliceEqual(intSlice1, intSlice3), "Expected slices to be different")

	// 测试字符串切片
	strSlice1 := []string{"a", "b", "c"}
	strSlice2 := []string{"a", "b", "c"}
	strSlice3 := []string{"a", "b", "d"}

	assert.True(t, mathx.SliceEqual(strSlice1, strSlice2), "Expected slices to be equal")
	assert.False(t, mathx.SliceEqual(strSlice1, strSlice3), "Expected slices to be different")

	// 测试自定义结构体切片
	type Person struct {
		Name string
		Age  int
	}

	personSlice1 := []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	personSlice2 := []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	personSlice3 := []Person{{Name: "Charlie", Age: 35}}

	assert.True(t, mathx.SliceEqual(personSlice1, personSlice2), "Expected slices to be equal")
	assert.False(t, mathx.SliceEqual(personSlice1, personSlice3), "Expected slices to be different")
}

// TestSliceFisherYates 测试 Fisher-Yates 洗牌算法
func TestSliceFisherYates(t *testing.T) {
	tests := [][]int{
		{1, 2, 3, 4, 5},
		{10, 20, 30, 40, 50},
	}

	for _, original := range tests {
		// 进行多次洗牌测试
		shuffledCount := 0
		for i := 0; i < 100; i++ {
			// 复制原始数组以便每次测试都有相同的输入
			testSlice := make([]int, len(original))
			copy(testSlice, original)

			// 调用洗牌函数，设置最大重试次数为 100
			err := mathx.SliceFisherYates(testSlice, 100)
			if err != nil {
				t.Errorf("Error during shuffling: %v", err)
				continue // 继续进行下一个测试
			}

			// 检查洗牌后的数组是否与原数组相同
			if !mathx.SliceEqual(original, testSlice) {
				shuffledCount++
			}
		}

		// 断言至少有一次洗牌结果与原数组不同
		assert.Greater(t, shuffledCount, 0, "SliceFisherYates did not shuffle the slice: original %v", original)
	}
}

func TestSliceQuickSortInPlace(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{3, 6, 8, 10, 1, 2, 1}, []int{1, 1, 2, 3, 6, 8, 10}},
		{[]int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{[]int{1}, []int{1}},
		{[]int{}, []int{}},               // 测试空切片
		{[]int{2, 2, 2}, []int{2, 2, 2}}, // 测试所有元素相同
	}

	for _, test := range tests {
		// 调用快速排序
		mathx.InsertionSort(test.input)

		// 使用 assert 验证排序结果
		assert.Equal(t, test.expected, test.input, "InsertionSort(%v) = %v; expected %v", test.input, test.input, test.expected)
	}
}

// 测试快速排序的正确性
func TestQuickSort(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{3, 2, 1}, []int{1, 2, 3}},
		{[]int{5, 3, 8, 4, 2}, []int{2, 3, 4, 5, 8}},
		{[]int{1, 1, 1}, []int{1, 1, 1}},
		{[]int{}, []int{}}, // 测试空数组
	}

	for _, test := range tests {
		// 复制输入数组以避免修改原始数据
		arr := make([]int, len(test.input))
		copy(arr, test.input)

		mathx.QuickSort(arr, 0, len(arr)-1)

		// 使用 assert 进行验证
		assert.Equal(t, test.expected, arr, "对于输入 %v，期望 %v，但得到 %v", test.input, test.expected, arr)
	}
}

// 测试冒泡排序的正确性
func TestBubbleSort(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int
	}{
		{[]int{3, 2, 1}, []int{1, 2, 3}},
		{[]int{5, 3, 8, 4, 2}, []int{2, 3, 4, 5, 8}},
		{[]int{1, 1, 1}, []int{1, 1, 1}},
		{[]int{}, []int{}}, // 测试空数组
	}

	for _, test := range tests {
		// 复制输入数组以避免修改原始数据
		arr := make([]int, len(test.input))
		copy(arr, test.input)

		mathx.BubbleSort(arr)

		// 使用 assert 进行验证
		assert.Equal(t, test.expected, arr, "对于输入 %v，期望 %v，但得到 %v", test.input, test.expected, arr)
	}
}

// TestSliceContains 测试 SliceContains 函数
func TestSliceContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		element  int
		expected bool
	}{
		{"ContainsTrue", []int{1, 2, 3}, 2, true},
		{"ContainsFalse", []int{1, 2, 3}, 4, false},
		{"ContainsEmpty", []int{}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.SliceContains(tt.slice, tt.element)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSliceHasDuplicates 测试 SliceHasDuplicates 函数
func TestSliceHasDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected bool
	}{
		{"HasDuplicatesTrue", []int{1, 2, 2, 3}, true},
		{"HasDuplicatesFalse", []int{1, 2, 3}, false},
		{"HasDuplicatesEmpty", []int{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.SliceHasDuplicates(tt.slice)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSliceRemoveEmpty 测试 SliceRemoveEmpty 函数
func TestSliceRemoveEmpty(t *testing.T) {
	tests := []struct {
		name     string
		slice    []interface{}
		expected []interface{}
	}{
		{"RemoveEmpty", []interface{}{1, "", nil, 2}, []interface{}{1, 2}},
		{"RemoveAllEmpty", []interface{}{nil, "", nil}, []interface{}{}},
		{"RemoveNoEmpty", []interface{}{1, 2, 3}, []interface{}{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.SliceRemoveEmpty(tt.slice)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

// TestSliceRemoveDuplicates 测试 SliceRemoveDuplicates 函数
func TestSliceRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{"RemoveDuplicates", []int{1, 2, 2, 3}, []int{1, 2, 3}},
		{"RemoveNoDuplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"RemoveAllDuplicates", []int{1, 1, 1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.SliceRemoveDuplicates(tt.slice)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

// TestSliceRemoveZero 测试 SliceRemoveZero 函数
func TestSliceRemoveZero(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{"RemoveZeros", []int{0, 1, 2, 0, 3}, []int{1, 2, 3}},
		{"NoZeros", []int{1, 2, 3}, []int{1, 2, 3}},
		{"AllZeros", []int{0, 0, 0}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mathx.SliceRemoveZero(tt.slice)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

// TestSliceChunk 测试 SliceChunk 函数
func TestSliceChunk(t *testing.T) {
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
			result := mathx.SliceChunk(tt.slice, tt.size)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSliceDiffSetStrings(t *testing.T) {
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
		result := mathx.SliceDiffSetSorted(tc.arr1, tc.arr2)
		assert.ElementsMatch(t, tc.want, result, "SliceDiffSet(%v, %v) = %v; want %v", tc.arr1, tc.arr2, result, tc.want)
	}
}

func TestSliceDiffSetInts(t *testing.T) {
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
		result := mathx.SliceDiffSetSorted(tc.arr1, tc.arr2)
		assert.ElementsMatch(t, tc.want, result, "SliceDiffSet(%v, %v) = %v; want %v", tc.arr1, tc.arr2, result, tc.want)
	}
}

func TestSliceDiffSetFloats(t *testing.T) {
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
		result := mathx.SliceDiffSetSorted(tc.arr1, tc.arr2)
		assert.ElementsMatch(t, tc.want, result, "SliceDiffSet(%v, %v) = %v; want %v", tc.arr1, tc.arr2, result, tc.want)
	}
}

// TestSliceUniq 测试 SliceUniq 函数
func TestSliceUniq(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Empty slice", []int{}, []int{}},
		{"No duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"With duplicates", []int{1, 2, 2, 3, 1}, []int{1, 2, 3}},
		{"All duplicates", []int{1, 1, 1}, []int{1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := mathx.SliceUniq(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestSliceDiff 测试 SliceDiff 函数
func TestSliceDiff(t *testing.T) {
	tests := []struct {
		name      string
		list1     []int
		list2     []int
		expected1 []int
		expected2 []int
	}{
		{"No difference", []int{1, 2}, []int{1, 2}, []int{}, []int{}},
		{"Some difference", []int{1, 2, 3}, []int{2, 3, 4}, []int{1}, []int{4}},
		{"All different", []int{1, 2}, []int{3, 4}, []int{1, 2}, []int{3, 4}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result1, result2 := mathx.SliceDiff(test.list1, test.list2)
			assert.Equal(t, test.expected1, result1)
			assert.Equal(t, test.expected2, result2)
		})
	}
}

// TestSliceWithout 测试 SliceWithout 函数
func TestSliceWithout(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		exclude  []int
		expected []int
	}{
		{"No exclude", []int{1, 2, 3}, []int{}, []int{1, 2, 3}},
		{"Some excluded", []int{1, 2, 3}, []int{2}, []int{1, 3}},
		{"All excluded", []int{1, 2, 3}, []int{1, 2, 3}, []int{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := mathx.SliceWithout(test.input, test.exclude...)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestSliceIntersect 测试 SliceIntersect 函数
func TestSliceIntersect(t *testing.T) {
	tests := []struct {
		name     string
		list1    []int
		list2    []int
		expected []int
	}{
		{"No intersection", []int{1, 2}, []int{3, 4}, []int{}},
		{"Some intersection", []int{1, 2, 3}, []int{2, 3, 4}, []int{2, 3}},
		{"All intersecting", []int{1, 2}, []int{1, 2}, []int{1, 2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := mathx.SliceIntersect(test.list1, test.list2)
			assert.Equal(t, test.expected, result)
		})
	}
}
