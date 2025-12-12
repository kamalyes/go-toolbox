/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:30:15
 * @FilePath: \go-toolbox\pkg\stringx\remove_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveAll(t *testing.T) {
	result := RemoveAll("aa-bb-cc-dd", "-")
	assert.Equal(t, "aabbccdd", result)
}

func TestRemoveAny(t *testing.T) {
	result := RemoveAny("aa-bb-cc-dd", []string{"-", "b"})
	assert.Equal(t, "aaccdd", result)
}

func TestRemoveAllLineBreaks(t *testing.T) {
	result := RemoveAllLineBreaks("Hello\r\nWorld")
	assert.Equal(t, "HelloWorld", result)
}

func TestRemovePrefix(t *testing.T) {
	result := RemovePrefix("hello", "he")
	assert.Equal(t, "llo", result)
}

func TestRemovePrefixIgnoreCase(t *testing.T) {
	result := RemovePrefixIgnoreCase("hELLo", "he")
	assert.Equal(t, "LLo", result)

	result = RemovePrefixIgnoreCase("HeLLo", "he")
	assert.Equal(t, "LLo", result)

	result = RemovePrefixIgnoreCase("heLlo", "he")
	assert.Equal(t, "Llo", result)
}

func TestRemoveSuffix(t *testing.T) {
	result := RemoveSuffix("hello", "lo")
	assert.Equal(t, "hel", result)
}

func TestRemoveSuffixIgnoreCase(t *testing.T) {
	result := RemoveSuffixIgnoreCase("helLO", "lo")
	assert.Equal(t, "hel", result)

	result = RemovePrefixIgnoreCase("HeLlo", "he")
	assert.Equal(t, "Llo", result)

	result = RemovePrefixIgnoreCase("heLlo", "he")
	assert.Equal(t, "Llo", result)
}

// TestRemoveSymbols 测试 RemoveSymbols
func TestRemoveSymbols(t *testing.T) {
	input := "Hello, World! 123"
	result := RemoveSymbols(input)
	assert.Equal(t, "HelloWorld123", result, "Expected cleaned string")
}

// TestRemoveSymbolsChain 测试 RemoveSymbols 的有效情况
func TestRemoveSymbolsChain(t *testing.T) {
	input := New("Hello, World! 123")
	result := input.RemoveSymbolsChain().Value()
	assert.Equal(t, "HelloWorld123", result, "Expected cleaned string")
}
