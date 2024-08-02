/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-02 14:39:19
 * @FilePath: \go-toolbox\numberx\compare_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package numberx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareAllFunctions(t *testing.T) {
	t.Run("TestNumberx_Min", TestNumberx_Min)
	t.Run("TestNumberx_Max", TestNumberx_Max)
}

func TestNumberx_Min(t *testing.T) {
	tests := map[string]struct {
		intList  []interface{}
		expected interface{}
		err      error
	}{
		"Empty_list": {
			intList:  []interface{}(nil),
			expected: nil,
			err:      errors.New("list is empty"),
		},
		"Single_element_list": {
			intList:  []interface{}{5},
			expected: 5,
			err:      nil,
		},
		"Positive_numbers_list": {
			intList:  []interface{}{10, 5, 8, 3, 12, 6},
			expected: 3,
			err:      nil,
		},
		"Negative_numbers_list": {
			intList:  []interface{}{-10, -5, -8, -3, -12, -6},
			expected: -12,
			err:      nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := MinMax(tc.intList, func(a, b interface{}) interface{} {
				val1 := a.(int)
				val2 := b.(int)
				if val1 < val2 {
					return val1
				}
				return val2
			})
			assert.Equal(t, tc.expected, res)
			assert.Equal(t, tc.err, err)
		})
	}

	// 测试不同类型的列表
	floatList := []interface{}{7.5, 3.2, 8.7, 2.1}
	floatMax, _ := MinMax(floatList, func(a, b interface{}) interface{} {
		if a.(float64) < b.(float64) {
			return a
		}
		return b
	})
	assert.Equal(t, float64(2.1), floatMax)

	boolList := []interface{}{true, false, true}
	boolMin, _ := MinMax(boolList, func(a, b interface{}) interface{} {
		if a.(bool) {
			return a
		}
		return b
	})
	assert.Equal(t, true, boolMin)

	runeList := []interface{}{'a', 'b', 'c'}
	runeMin, _ := MinMax(runeList, func(a, b interface{}) interface{} {
		if a.(rune) < b.(rune) {
			return a
		}
		return b
	})
	assert.Equal(t, 'a', runeMin)
}

func TestNumberx_Max(t *testing.T) {
	tests := map[string]struct {
		intList  []interface{}
		expected interface{}
		err      error
	}{
		"Empty_list": {
			intList:  []interface{}(nil),
			expected: nil,
			err:      errors.New("list is empty"),
		},
		"Single_element_list": {
			intList:  []interface{}{5},
			expected: 5,
			err:      nil,
		},
		"Positive_numbers_list": {
			intList:  []interface{}{10, 5, 8, 3, 12, 6},
			expected: 12,
			err:      nil,
		},
		"Negative_numbers_list": {
			intList:  []interface{}{-10, -5, -8, -3, -12, -6},
			expected: -3,
			err:      nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := MinMax(tc.intList, func(a, b interface{}) interface{} {
				val1 := a.(int)
				val2 := b.(int)
				if val1 > val2 {
					return val1
				}
				return val2
			})
			assert.Equal(t, tc.expected, res)
			assert.Equal(t, tc.err, err)
		})
	}

	// 测试不同类型的列表
	floatList := []interface{}{7.5, 3.2, 8.7, 2.1}
	floatMax, _ := MinMax(floatList, func(a, b interface{}) interface{} {
		if a.(float64) > b.(float64) {
			return a
		}
		return b
	})
	assert.Equal(t, float64(8.7), floatMax)

	boolList := []interface{}{true, false, true}
	boolMax, _ := MinMax(boolList, func(a, b interface{}) interface{} {
		if a.(bool) {
			return a
		}
		return b
	})
	assert.Equal(t, true, boolMax)

	runeList := []interface{}{'a', 'b', 'c'}
	runeMax, _ := MinMax(runeList, func(a, b interface{}) interface{} {
		if a.(rune) > b.(rune) {
			return a
		}
		return b
	})
	assert.Equal(t, 'c', runeMax)
}
