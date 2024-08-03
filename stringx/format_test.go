/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 23:10:59
 * @FilePath: \go-toolbox\stringx\format_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllFormatFunctions(t *testing.T) {
	t.Run("TestFillBefore", TestFillBefore)
	t.Run("TestFillAfter", TestFillAfter)
	t.Run("TestFormat", TestFormat)
	t.Run("TestIndexedFormat", TestIndexedFormat)
	t.Run("TestTruncateAppendEllipsis", TestTruncateAppendEllipsis)
	t.Run("TestTruncate", TestTruncate)
	t.Run("TestAddPrefixIfNot", TestAddPrefixIfNot)
	t.Run("TestAddSuffixIfNot", TestAddSuffixIfNot)
}

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
	result := TruncateAppendEllipsis("This is a very long string", 10)
	assert.Equal(t, "This is...", result)
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
