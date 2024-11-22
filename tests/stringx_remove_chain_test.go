/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:41:02
 * @FilePath: \go-toolbox\tests\stringx_remove_chain_test.go
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

func TestRemoveAllChain(t *testing.T) {
	result := stringx.New("aa-bb-cc-dd").RemoveAllChain("-").Value()
	assert.Equal(t, "aabbccdd", result)
}

func TestRemoveAnyChain(t *testing.T) {
	strsToRemove := []string{"a", "b"}
	result := stringx.New("aa-bb-cc-dd").RemoveAnyChain(strsToRemove).Value()
	assert.Equal(t, "--cc-dd", result)
}

func TestRemoveAllLineBreaksChain(t *testing.T) {
	result := stringx.New("Hello\nWorld\r\n").RemoveAllLineBreaksChain().Value()
	assert.Equal(t, "HelloWorld", result)
}

func TestRemovePrefixChain(t *testing.T) {
	result := stringx.New("HelloWorld").RemovePrefixChain("Hello").Value()
	assert.Equal(t, "World", result)

	result = stringx.New("World").RemovePrefixChain("Hello").Value()
	assert.Equal(t, "World", result)
}

func TestRemovePrefixIgnoreCaseChain(t *testing.T) {
	result := stringx.New("HelloWorld").RemovePrefixIgnoreCaseChain("hello").Value()
	assert.Equal(t, "World", result)

	result = stringx.New("World").RemovePrefixIgnoreCaseChain("hello").Value()
	assert.Equal(t, "World", result)
}

func TestRemoveSuffixChain(t *testing.T) {
	result := stringx.New("HelloWorld").RemoveSuffixChain("World").Value()
	assert.Equal(t, "Hello", result)

	result = stringx.New("Hello").RemoveSuffixChain("World").Value()
	assert.Equal(t, "Hello", result)
}

func TestRemoveSuffixIgnoreCaseChain(t *testing.T) {
	result := stringx.New("HelloWorld").RemoveSuffixIgnoreCaseChain("world").Value()
	assert.Equal(t, "Hello", result)

	result = stringx.New("Hello").RemoveSuffixIgnoreCaseChain("world").Value()
	assert.Equal(t, "Hello", result)
}
