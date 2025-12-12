/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:53:15
 * @FilePath: \go-toolbox\pkg\stringx\format_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillBefore(t *testing.T) {
	result := FillBefore("hello", ".", 10)
	assert.Equal(t, ".....hello", result)
}

func TestFillAfter(t *testing.T) {
	result := FillAfter("hello", ".", 10)
	assert.Equal(t, "hello.....", result)
}

func TestFormat(t *testing.T) {
	params := map[string]interface{}{
		"a": "aValue",
		"b": "bValue",
	}
	result := Format("{a} and {b}", params)
	assert.Equal(t, "aValue and bValue", result)
}

func TestIndexedFormat(t *testing.T) {
	result := IndexedFormat("this is {0} for {1}", []interface{}{"a", "b"})
	assert.Equal(t, "this is a for b", result)
}

func TestTruncateAppendEllipsis(t *testing.T) {
	tests := []struct {
		input    string
		maxChars int
		expected string
	}{
		{"这是一个测试字符串199665889@#￥￥", 10, "这是一个测试字符串1..."},
		{"这是一个测试字符串12356789@#￥￥", 50, "这是一个测试字符串12356789@#￥￥"},
		{"", 10, ""},
	}

	for _, test := range tests {
		result := TruncateAppendEllipsis(test.input, test.maxChars)
		if result != test.expected {
			t.Errorf("TruncateAppendEllipsis(%q, %d) = %q; want %q", test.input, test.maxChars, result, test.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	result := Truncate("This is another long string", 10)
	assert.Equal(t, "This is an", result)
}

func TestAddPrefixIfNot(t *testing.T) {
	result := AddPrefixIfNot("world", "hello ")
	assert.Equal(t, "hello world", result)
}

func TestAddSuffixIfNot(t *testing.T) {
	result := AddSuffixIfNot("hello", " world")
	assert.Equal(t, "hello world", result)
}
