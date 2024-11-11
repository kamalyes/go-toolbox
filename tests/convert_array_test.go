/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 15:56:07
 * @FilePath: \go-toolbox\tests\array_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/array"
	"github.com/stretchr/testify/assert"
)

func TestArrayAllFunctions(t *testing.T) {
	t.Run("TestInterfaceArrayDiffSet", TestInterfaceArrayDiffSet)
	t.Run("TestInterfaceArrayUnion", TestInterfaceArrayUnion)
	t.Run("TestIsInterfaceArrayExistElement", TestIsInterfaceArrayExistElement)
	t.Run("TestIsExistRepeatInInterfaceArray", TestIsExistRepeatInInterfaceArray)
	t.Run("TestRemoveEmptyInterfaceInArray", TestRemoveEmptyInterfaceInArray)
	t.Run("TestInt64ToStringWithDecimals", TestInt64ToStringWithDecimals)
	t.Run("TestRemoveDuplicatesInInterfaceSlice", TestRemoveDuplicatesInInterfaceSlice)
	t.Run("TestRemoveZeroInInterfaceSlice", TestRemoveZeroInInterfaceSlice)

}

func TestInterfaceArrayDiffSet(t *testing.T) {
	cases := []struct {
		arr1 interface{}
		arr2 interface{}
		want interface{}
	}{
		{[]string{"a", "b", "c"}, []string{"b", "c", "d"}, []interface{}{"a", "d"}},
		{[]string{}, []string{"b", "c", "d"}, []interface{}{"b", "c", "d"}},
		{[]string{"a", "b", "c"}, []string{}, []interface{}{"a", "b", "c"}},
		{[]int{1, 2, 3}, []int{2, 3, 4}, []interface{}{1, 4}},
		{[]string{"apple", "banana"}, []string{"banana", "cherry"}, []interface{}{"apple", "cherry"}},
	}

	for _, tc := range cases {
		result := array.InterfaceArrayDiffSet(tc.arr1, tc.arr2)

		// 检查结果是否与期望一致（长度和值）
		if reflect.ValueOf(result).Len() != reflect.ValueOf(tc.want).Len() {
			t.Errorf("InterfaceArrayDiffSet(%v, %v) returned result with length %d, want length %d", tc.arr1, tc.arr2, reflect.ValueOf(result).Len(), reflect.ValueOf(tc.want).Len())
		} else {
			// 检查结果是否包含所有期望的元素（忽略顺序）
			for _, v := range tc.want.([]interface{}) {
				found := false
				for i := 0; i < reflect.ValueOf(result).Len(); i++ {
					if reflect.DeepEqual(reflect.ValueOf(result).Index(i).Interface(), v) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("InterfaceArrayDiffSet(%v, %v) does not contain expected value: %v", tc.arr1, tc.arr2, v)
				}
			}
		}
	}
}

// 测试 InterfaceArrayUnion 函数，验证计算多种类型数组并集的功能
func TestInterfaceArrayUnion(t *testing.T) {
	// 定义测试案例，包括两个输入数组和期望的并集结果
	testCases := []struct {
		arr1 interface{}
		arr2 interface{}
		want interface{}
	}{
		{[]int{1, 2, 3}, []int{2, 3, 4}, []interface{}{1, 2, 3, 4}},                                             // 测试两个整数数组的并集
		{[]string{"apple", "banana"}, []string{"banana", "cherry"}, []interface{}{"cherry", "apple", "banana"}}, // 测试两个字符串数组的并集
		{[]int8{1, 2, 3}, []int8{2, 3, 4}, []interface{}{int8(1), int8(2), int8(3), int8(4)}},                   // 测试 int8 类型数组的并集
	}

	// 遍历每个测试案例并执行测试
	for _, tc := range testCases {
		result := array.InterfaceArrayUnion(tc.arr1, tc.arr2)

		// 检查结果是否与期望一致（长度和值）
		if reflect.ValueOf(result).Len() != reflect.ValueOf(tc.want).Len() {
			t.Errorf("InterfaceArrayUnion(%v, %v) returned result with length %d, want length %d", tc.arr1, tc.arr2, reflect.ValueOf(result).Len(), reflect.ValueOf(tc.want).Len())
		} else {
			// 检查结果是否包含所有期望的元素（忽略顺序）
			for _, v := range tc.want.([]interface{}) {
				found := false
				for i := 0; i < reflect.ValueOf(result).Len(); i++ {
					if reflect.DeepEqual(reflect.ValueOf(result).Index(i).Interface(), v) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("InterfaceArrayUnion(%v, %v) does not contain expected value: %v", tc.arr1, tc.arr2, v)
				}
			}
		}
	}
}

func TestIsInterfaceArrayExistElement(t *testing.T) {
	// 测试 IsStrArrayExistArray 函数，检测指定字符串是否存在于字符串数组中
	cases := []struct {
		array   []interface{}
		element interface{}
		want    bool
	}{
		{nil, "nil", false},
		{[]interface{}{"a", "b", "c"}, "b", true},
		{[]interface{}{"apple", "banana", "cherry"}, "orange", false},
		{[]interface{}{1, 2, 3, 4}, 3, true},
		{[]interface{}{true, false, true}, false, true},
	}

	for _, tc := range cases {
		got := array.IsInterfaceArrayExistElement(tc.array, tc.element)
		if got != tc.want {
			t.Errorf("IsStrArrayExistArray(%v, %v) = %v; want %v", tc.array, tc.element, got, tc.want)
		}
	}
}

func TestIsExistRepeatInInterfaceArray(t *testing.T) {
	cases := []struct {
		array    []interface{}
		expected bool
	}{
		{[]interface{}{"a", "b", "c"}, false},
		{[]interface{}{"a", "b", "c", "b"}, true},
		{[]interface{}{"苹果", "香蕉", "橙子"}, false},
		{[]interface{}{"苹果", "香蕉", "苹果", "橙子"}, true},
		{[]interface{}{1, 2, 3, 4, 3}, true},
		{[]interface{}{true, false, true}, true},
	}

	for _, tc := range cases {
		result := array.IsExistRepeatInInterfaceArray(tc.array)
		if result != tc.expected {
			t.Errorf("IsExistRepeatInInterfaceArray(%v) returned %v, expected %v", tc.array, result, tc.expected)
		}
	}
}

func TestInt64ToStringWithDecimals(t *testing.T) {
	cases := []struct {
		num      int64
		digit    int
		expected string
	}{
		{1234, 2, "12.34"},
		{5000, 3, "5.000"},
		{987654321, 5, "9876.54321"},
	}

	for _, tc := range cases {
		result := array.Int64ToStringWithDecimals(tc.num, tc.digit)
		if result != tc.expected {
			t.Errorf("Int64ToStringWithDecimals(%d, %d) returned %s, expected %s", tc.num, tc.digit, result, tc.expected)
		}
	}
}

func TestRemoveEmptyInterfaceInArray(t *testing.T) {
	cases := []struct {
		array    []interface{}
		expected []interface{}
	}{
		{[]interface{}{"a", "", "b", nil, "c"}, []interface{}{"a", "b", nil, "c"}},
		{[]interface{}{1, 0, 2, 3, nil}, []interface{}{1, 2, 3, nil}},
		{[]interface{}{true, false, true, nil}, []interface{}{true, true, nil}},
	}

	for _, tc := range cases {
		result := array.RemoveEmptyInterfaceInArray(tc.array)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("RemoveEmptyInterfaceInArray(%v) returned %v, expected %v", tc.array, result, tc.expected)
		}
	}
}

func TestRemoveDuplicatesInInterfaceSlice(t *testing.T) {
	cases := []struct {
		numbers  []interface{}
		expected []interface{}
	}{
		{[]interface{}{1, 2, 3, 2, 1}, []interface{}{1, 2, 3}},
		{[]interface{}{"a", "b", "c", "b"}, []interface{}{"a", "b", "c"}},
	}

	for _, tc := range cases {
		result := array.RemoveDuplicatesInInterfaceSlice(tc.numbers)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("RemoveDuplicatesInInterfaceSlice(%v) returned %v, expected %v", tc.numbers, result, tc.expected)
		}
	}
}

func TestRemoveZeroInInterfaceSlice(t *testing.T) {
	cases := []struct {
		array    []interface{}
		expected []interface{}
	}{
		{[]interface{}{1, 0, 2, 0, 3}, []interface{}{1, 2, 3}},
		{[]interface{}{0, 0, 0, 0}, []interface{}(nil)},
	}

	for _, tc := range cases {
		result := array.RemoveZeroInInterfaceSlice(tc.array)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("RemoveZeroInInterfaceSlice(%v) returned %v, expected %v", tc.array, result, tc.expected)
		}
	}
}

// TestArrayChunk 测试 ArrayChunk 函数的基本功能
func TestArrayChunk(t *testing.T) {
	tests := []struct {
		input    []int
		size     int
		expected [][]int
	}{
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 3, [][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {9}}},
		{[]int{1, 2, 3}, 1, [][]int{{1}, {2}, {3}}},
		{[]int{1, 2, 3}, 5, [][]int{{1, 2, 3}}},
		{[]int{}, 2, [][]int{}},
		{[]int{1, 2, 3}, 0, nil},
	}

	for _, test := range tests {
		result := array.Chunk(test.input, test.size)
		assert.Equal(t, test.expected, result, fmt.Sprintf("ArrayChunk(%v, %d) = %v; expected %v", test.input, test.size, result, test.expected))
	}
}

// equal 比较两个切片是否相等
func equal(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
