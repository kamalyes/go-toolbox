/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 15:26:51
 * @FilePath: \go-toolbox\tests\stringx_split_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	result := stringx.Split("one,two,three,four", ",")
	assert.Equal(t, []string{"one", "two", "three", "four"}, result)
}

func TestSplitLimit(t *testing.T) {
	result := stringx.SplitLimit("one,two,three,four", ",", 2)
	assert.Equal(t, []string{"one", "two,three,four"}, result)
}

func TestSplitTrim(t *testing.T) {
	result := stringx.SplitTrim(" one , two , three , four ", ",")
	assert.Equal(t, []string{"one", "two", "three", "four"}, result)
}

func TestSplitTrimLimit(t *testing.T) {
	result := stringx.SplitTrimLimit(" one , two , three , four ", ",", 2)
	assert.Equal(t, []string{"one", "two , three , four"}, result)
}

func TestSplitByLen_Cut(t *testing.T) {
	tests := []struct {
		name     string
		function func(string, int) []string
		input    string
		param    int
		expected []string
	}{
		{
			name:     "SplitByLen - normal case",
			function: stringx.SplitByLen,
			input:    "HelloWorld",
			param:    3,
			expected: []string{"Hel", "loW", "orl", "d"},
		},
		{
			name:     "SplitByLen - length greater than string",
			function: stringx.SplitByLen,
			input:    "Hi",
			param:    5,
			expected: []string{"Hi"},
		},
		{
			name:     "SplitByLen - empty string",
			function: stringx.SplitByLen,
			input:    "",
			param:    3,
			expected: []string{},
		},
		{
			name:     "SplitByLen - zero length",
			function: stringx.SplitByLen,
			input:    "Test",
			param:    0,
			expected: []string{},
		},
		{
			name:     "Cut - normal case",
			function: stringx.Cut,
			input:    "HelloWorld",
			param:    3,
			expected: []string{"Hell", "oWo", "rld"},
		},
		{
			name:     "Cut - n greater than string length",
			function: stringx.Cut,
			input:    "Hi",
			param:    5,
			expected: []string{"H", "i", "", "", ""},
		},
		{
			name:     "Cut - empty string",
			function: stringx.Cut,
			input:    "",
			param:    3,
			expected: []string{},
		},
		{
			name:     "Cut - zero parts",
			function: stringx.Cut,
			input:    "Test",
			param:    0,
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.function(test.input, test.param)
			assert.Equal(t, test.expected, result, "Expected %v but got %v", test.expected, result)
		})
	}
}

func TestSplitAfterMapping(t *testing.T) {
	tests := []struct {
		input       string
		separator   string
		mapping     func(s string) (int, error)
		expected    []int
		expectPanic bool
	}{
		{
			input:     "1,2,3",
			separator: ",",
			mapping: func(s string) (int, error) {
				var i int
				_, err := fmt.Sscanf(s, "%d", &i)
				return i, err
			},
			expected:    []int{1, 2, 3},
			expectPanic: false,
		},
		{
			input:     "4,5,6",
			separator: ",",
			mapping: func(s string) (int, error) {
				var i int
				_, err := fmt.Sscanf(s, "%d", &i)
				return i, err
			},
			expected:    []int{4, 5, 6},
			expectPanic: false,
		},
		{
			input:     "7,8,a",
			separator: ",",
			mapping: func(s string) (int, error) {
				var i int
				_, err := fmt.Sscanf(s, "%d", &i)
				return i, err
			},
			expected:    []int{7, 8},
			expectPanic: true, // Expect panic due to invalid mapping
		},
	}

	for _, test := range tests {
		if test.expectPanic {
			assert.Panics(t, func() {
				stringx.SplitAfterMapping(test.input, test.separator, test.mapping)
			}, "Expected panic for input %q, but did not panic", test.input)
		} else {
			result := stringx.SplitAfterMapping(test.input, test.separator, test.mapping)
			assert.Equal(t, test.expected, result, "For input %q, expected %v but got %v", test.input, test.expected, result)
		}
	}
}
