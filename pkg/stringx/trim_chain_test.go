/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-22 10:07:57
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:28:19
 * @FilePath: \go-toolbox\pkg\stringx\trim_chain_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimChain(t *testing.T) {
	result := New("  Hello, World!  ").TrimChain().Value()
	assert.Equal(t, "Hello, World!", result)

	result = New("   ").TrimChain().Value()
	assert.Equal(t, "", result) // All spaces
}

func TestTrimStartChain(t *testing.T) {
	result := New("  Hello, World!").TrimStartChain().Value()
	assert.Equal(t, "Hello, World!", result)

	result = New("   ").TrimStartChain().Value()
	assert.Equal(t, "", result) // All spaces
}

func TestTrimEndChain(t *testing.T) {
	result := New("Hello, World!  ").TrimEndChain().Value()
	assert.Equal(t, "Hello, World!", result)

	result = New("   ").TrimEndChain().Value()
	assert.Equal(t, "", result) // All spaces
}

func TestCleanEmptyChain(t *testing.T) {
	result := New("H e llo, W o rld!").CleanEmptyChain().Value()
	assert.Equal(t, "Hello,World!", result)

	result = New("   ").CleanEmptyChain().Value()
	assert.Equal(t, "", result) // All spaces
}
