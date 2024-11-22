/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:57:11
 * @FilePath: \go-toolbox\tests\stringx_remove_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/stretchr/testify/assert"
)

func TestRemoveAll(t *testing.T) {
	result := stringx.RemoveAll("aa-bb-cc-dd", "-")
	assert.Equal(t, "aabbccdd", result)
}

func TestRemoveAny(t *testing.T) {
	result := stringx.RemoveAny("aa-bb-cc-dd", []string{"-", "b"})
	assert.Equal(t, "aaccdd", result)
}

func TestRemoveAllLineBreaks(t *testing.T) {
	result := stringx.RemoveAllLineBreaks("Hello\r\nWorld")
	assert.Equal(t, "HelloWorld", result)
}

func TestRemovePrefix(t *testing.T) {
	result := stringx.RemovePrefix("hello", "he")
	assert.Equal(t, "llo", result)
}

func TestRemovePrefixIgnoreCase(t *testing.T) {
	result := stringx.RemovePrefixIgnoreCase("hELLo", "he")
	assert.Equal(t, "LLo", result)

	result = stringx.RemovePrefixIgnoreCase("HeLLo", "he")
	assert.Equal(t, "LLo", result)

	result = stringx.RemovePrefixIgnoreCase("heLlo", "he")
	assert.Equal(t, "Llo", result)
}

func TestRemoveSuffix(t *testing.T) {
	result := stringx.RemoveSuffix("hello", "lo")
	assert.Equal(t, "hel", result)
}

func TestRemoveSuffixIgnoreCase(t *testing.T) {
	result := stringx.RemoveSuffixIgnoreCase("helLO", "lo")
	assert.Equal(t, "hel", result)

	result = stringx.RemovePrefixIgnoreCase("HeLlo", "he")
	assert.Equal(t, "Llo", result)

	result = stringx.RemovePrefixIgnoreCase("heLlo", "he")
	assert.Equal(t, "Llo", result)
}
