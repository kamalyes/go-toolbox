/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 23:10:59
 * @FilePath: \go-toolbox\pkg\stringx\format_chain_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillBeforeChain(t *testing.T) {
	result := New("Hello").FillBeforeChain("*", 10).Value()
	assert.Equal(t, "*****Hello", result)
}

func TestFillAfterChain(t *testing.T) {
	result := New("Hello").FillAfterChain("*", 10).Value()
	assert.Equal(t, "Hello*****", result)
}

func TestFormatChain(t *testing.T) {
	params := map[string]interface{}{"a": "aValue", "b": "bValue"}
	result := New("{a} and {b}").FormatChain(params).Value()
	assert.Equal(t, "aValue and bValue", result)
}

func TestIndexedFormatChain(t *testing.T) {
	params := []interface{}{"a", "b"}
	result := New("this is {0} for {1}").IndexedFormatChain(params).Value()
	assert.Equal(t, "this is a for b", result)
}

func TestTruncateAppendEllipsisChain(t *testing.T) {
	result := New("This is a long string.").TruncateAppendEllipsisChain(5).Value()
	assert.Equal(t, "This ...", result)
}

func TestTruncateChain(t *testing.T) {
	result := New("This is a long string.").TruncateChain(10).Value()
	assert.Equal(t, "This is a ", result)
}

func TestAddPrefixIfNotChain(t *testing.T) {
	result := New("World").AddPrefixIfNotChain("Hello ").Value()
	assert.Equal(t, "Hello World", result)

	result = New("Hello World").AddPrefixIfNotChain("Hello ").Value()
	assert.Equal(t, "Hello World", result)
}

func TestAddSuffixIfNotChain(t *testing.T) {
	result := New("Hello").AddSuffixIfNotChain(" World").Value()
	assert.Equal(t, "Hello World", result)

	result = New("Hello World").AddSuffixIfNotChain(" World").Value()
	assert.Equal(t, "Hello World", result)
}
