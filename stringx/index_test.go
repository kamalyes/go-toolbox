/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 16:39:40
 * @FilePath: \go-toolbox\stringx\index_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllIndexFunctions(t *testing.T) {
	t.Run("IndexOf", TestIndexOf)
	t.Run("IndexOfByRange", TestIndexOfByRange)
	t.Run("IndexOfByRangeStart", TestIndexOfByRangeStart)
	t.Run("IndexOfIgnoreCase", TestIndexOfIgnoreCase)
	t.Run("IndexOfIgnoreCaseByRange", TestIndexOfIgnoreCaseByRange)
	t.Run("LastIndexOf", TestLastIndexOf)
	t.Run("LastIndexOfIgnoreCase", TestLastIndexOfIgnoreCase)
	t.Run("LastIndexOfByRangeStart", TestLastIndexOfByRangeStart)
	t.Run("OrdinalIndexOf", TestOrdinalIndexOf)
}

func TestIndexOf(t *testing.T) {
	result := IndexOf("", "o")
	assert.Equal(t, -1, result)

	result = IndexOf("hello world", "o")
	assert.Equal(t, 4, result)
}

func TestIndexOfByRange(t *testing.T) {
	result := IndexOfByRange("hello world", "o", 10, 5)
	assert.Equal(t, -1, result)

	result = IndexOfByRange("hello world", "o", 5, 10)
	assert.Equal(t, 7, result)
}

func TestIndexOfByRangeStart(t *testing.T) {
	result := IndexOfByRangeStart("hello world", "o", 100)
	assert.Equal(t, -1, result)

	result = IndexOfByRangeStart("hello world", "o", 5)
	assert.Equal(t, 7, result)
}

func TestIndexOfIgnoreCase(t *testing.T) {
	result := IndexOfIgnoreCase("", "")
	assert.Equal(t, 0, result)

	result = IndexOfIgnoreCase("Hello WorLd", "llo")
	assert.Equal(t, 2, result)
}

func TestIndexOfIgnoreCaseByRange(t *testing.T) {
	result := IndexOfIgnoreCaseByRange("", "hello", 5)
	assert.Equal(t, -1, result)

	result = IndexOfIgnoreCaseByRange("Hello123 WorLd", "123", 5)
	assert.Equal(t, 5, result)
}

func TestLastIndexOf(t *testing.T) {
	result := LastIndexOf("hello world", "x")
	assert.Equal(t, -1, result)

	result = LastIndexOf("hello world", "o")
	assert.Equal(t, 7, result)
}

func TestLastIndexOfIgnoreCase(t *testing.T) {
	result := LastIndexOfIgnoreCase("", "l")
	assert.Equal(t, -1, result)

	result = LastIndexOfIgnoreCase("Hello WorLd", "l")
	assert.Equal(t, 9, result)
}

func TestLastIndexOfByRangeStart(t *testing.T) {
	result := LastIndexOfByRangeStart("hello world", "f", 10)
	assert.Equal(t, -1, result)

	result = LastIndexOfByRangeStart("hello world", "o", 4)
	assert.Equal(t, 4, result)
}

func TestOrdinalIndexOf(t *testing.T) {
	result := OrdinalIndexOf("", "a", 1)
	assert.Equal(t, -1, result)

	result = OrdinalIndexOf("aabcbcdd", "b", 2)
	assert.Equal(t, 4, result)

	result = OrdinalIndexOf("ABCDEFGHIJKCLMNOPQRST", "C", 2)
	assert.Equal(t, 11, result)

	result = OrdinalIndexOf("", "", 3)
	assert.Equal(t, 0, result)
}

func TestIndexOfEmptyString(t *testing.T) {
	result := IndexOf("", "hello")
	assert.Equal(t, -1, result)

	result = IndexOf("hello world", "lo")
	assert.Equal(t, 3, result)

	result = IndexOf("", "")
	assert.Equal(t, 0, result)
}
