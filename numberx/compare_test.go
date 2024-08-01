/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 20:22:39
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
		intList  []int
		expected int
		err      error
	}{
		"Empty_list": {
			intList:  []int{},
			expected: 0,
			err:      errors.New("intList is empty"),
		},
		"Single_element_list": {
			intList:  []int{5},
			expected: 5,
			err:      nil,
		},
		"Positive_numbers_list": {
			intList:  []int{10, 5, 8, 3, 12, 6},
			expected: 3,
			err:      nil,
		},
		"Negative_numbers_list": {
			intList:  []int{-10, -5, -8, -3, -12, -6},
			expected: -12,
			err:      nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := Min(tc.intList)
			assert.Equal(t, tc.expected, res)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestNumberx_Max(t *testing.T) {
	tests := map[string]struct {
		intList  []int
		expected int
		err      error
	}{
		"Empty_list": {
			intList:  []int{},
			expected: 0,
			err:      errors.New("intList is empty"),
		},
		"Single_element_list": {
			intList:  []int{5},
			expected: 5,
			err:      nil,
		},
		"Positive_numbers_list": {
			intList:  []int{10, 5, 8, 3, 12, 6},
			expected: 12,
			err:      nil,
		},
		"Negative_numbers_list": {
			intList:  []int{-10, -5, -8, -3, -12, -6},
			expected: -3,
			err:      nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := Max(tc.intList)
			assert.Equal(t, tc.expected, res)
			assert.Equal(t, tc.err, err)
		})
	}
}
