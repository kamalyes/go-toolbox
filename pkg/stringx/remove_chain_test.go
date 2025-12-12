/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:30:02
 * @FilePath: \go-toolbox\pkg\stringx\remove_chain_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveAllChain(t *testing.T) {
	result := New("aa-bb-cc-dd").RemoveAllChain("-").Value()
	assert.Equal(t, "aabbccdd", result)
}

func TestRemoveAnyChain(t *testing.T) {
	strsToRemove := []string{"a", "b"}
	result := New("aa-bb-cc-dd").RemoveAnyChain(strsToRemove).Value()
	assert.Equal(t, "--cc-dd", result)
}

func TestRemoveAllLineBreaksChain(t *testing.T) {
	result := New("Hello\nWorld\r\n").RemoveAllLineBreaksChain().Value()
	assert.Equal(t, "HelloWorld", result)
}

func TestRemovePrefixChain(t *testing.T) {
	result := New("HelloWorld").RemovePrefixChain("Hello").Value()
	assert.Equal(t, "World", result)

	result = New("World").RemovePrefixChain("Hello").Value()
	assert.Equal(t, "World", result)
}

func TestRemovePrefixIgnoreCaseChain(t *testing.T) {
	result := New("HelloWorld").RemovePrefixIgnoreCaseChain("hello").Value()
	assert.Equal(t, "World", result)

	result = New("World").RemovePrefixIgnoreCaseChain("hello").Value()
	assert.Equal(t, "World", result)
}

func TestRemoveSuffixChain(t *testing.T) {
	result := New("HelloWorld").RemoveSuffixChain("World").Value()
	assert.Equal(t, "Hello", result)

	result = New("Hello").RemoveSuffixChain("World").Value()
	assert.Equal(t, "Hello", result)
}

func TestRemoveSuffixIgnoreCaseChain(t *testing.T) {
	result := New("HelloWorld").RemoveSuffixIgnoreCaseChain("world").Value()
	assert.Equal(t, "Hello", result)

	result = New("Hello").RemoveSuffixIgnoreCaseChain("world").Value()
	assert.Equal(t, "Hello", result)
}
