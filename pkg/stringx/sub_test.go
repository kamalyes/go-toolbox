/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:53:15
 * @FilePath: \go-toolbox\pkg\stringx\sub_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllSubFunctions(t *testing.T) {
	t.Run("TestSubBefore", TestSubBefore)
	t.Run("TestSubAfter", TestSubAfter)
	t.Run("TestSubBetween", TestSubBetween)
	t.Run("TestSubBetweenAll", TestSubBetweenAll)
}

func TestSubBefore(t *testing.T) {
	result := SubBefore("abcdef", "d", false)
	assert.Equal(t, "abc", result)

	result = SubBefore("abcdef", "d", true)
	assert.Equal(t, "abc", result)

	result = SubBefore("abcdef", "x", false)
	assert.Equal(t, "abcdef", result)

	result = SubBefore("abcdef", "a", false)
	assert.Equal(t, "", result)
}

func TestSubAfter(t *testing.T) {
	result := SubAfter("abcdef", "d", true)
	assert.Equal(t, "ef", result)

	result = SubAfter("abcdef", "x", false)
	assert.Equal(t, "", result)

	result = SubAfter("abcdef", "f", false)
	assert.Equal(t, "", result)
}

func TestSubBetween(t *testing.T) {
	result := SubBetween("abc123def456ghi", "abc", "def")
	assert.Equal(t, "123", result)

	result = SubBetween("abc123def456ghi", "def", "ghi")
	assert.Equal(t, "456", result)

	result = SubBetween("abc123def456ghi", "abc", "xyz")
	assert.Equal(t, "", result)
}

func TestSubBetweenAll(t *testing.T) {
	result := SubBetweenAll("a1b2c3d4e5f6", "b", "e")
	assert.ElementsMatch(t, []string{"2c3d4"}, result)

	result = SubBetweenAll("a1b2c3d4e5f6", "x", "y")
	assert.ElementsMatch(t, []string{}, result)
}

// TestSubString 测试 SubString 函数
func TestSubString(t *testing.T) {
	// 测试用例
	tests := []struct {
		input    string
		start    int
		length   int
		expected string
	}{
		{"Hello, World!", 0, 5, "Hello"},
		{"Hello, World!", 7, 5, "World"},
		{"Hello, World!", 0, 20, "Hello, World!"},
		{"Hello, World!", -1, 5, ""},
		{"Hello, World!", 13, 5, ""},
	}

	for _, test := range tests {
		result := SubString(test.input, test.start, test.length)
		assert.Equal(t, test.expected, result)
	}
}

// TestSubStringChain 测试 SubStringChain 方法
func TestSubStringChain(t *testing.T) {
	s := &StringX{value: "Hello, World!"}

	// 测试链式调用
	result := s.SubStringChain(7, 5)
	assert.Equal(t, "World", result.value)

	// 测试超出范围
	s = &StringX{value: "Hello, World!"}
	result = s.SubStringChain(0, 20)
	assert.Equal(t, "Hello, World!", result.value)

	// 测试负数起始位置
	s = &StringX{value: "Hello, World!"}
	result = s.SubStringChain(-1, 5)
	assert.Equal(t, "", result.value)

	// 测试起始位置超出
	s = &StringX{value: "Hello, World!"}
	result = s.SubStringChain(13, 5)
	assert.Equal(t, "", result.value)
}
