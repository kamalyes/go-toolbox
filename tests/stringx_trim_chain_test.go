/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-22 10:07:57
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:07:57
 * @FilePath: \go-toolbox\tests\stringx_trim_chain_test.go
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

func TestTrimChain(t *testing.T) {
	result := stringx.New("  Hello, World!  ").TrimChain().Value()
	assert.Equal(t, "Hello, World!", result)

	result = stringx.New("   ").TrimChain().Value()
	assert.Equal(t, "", result) // All spaces
}

func TestTrimStartChain(t *testing.T) {
	result := stringx.New("  Hello, World!").TrimStartChain().Value()
	assert.Equal(t, "Hello, World!", result)

	result = stringx.New("   ").TrimStartChain().Value()
	assert.Equal(t, "", result) // All spaces
}

func TestTrimEndChain(t *testing.T) {
	result := stringx.New("Hello, World!  ").TrimEndChain().Value()
	assert.Equal(t, "Hello, World!", result)

	result = stringx.New("   ").TrimEndChain().Value()
	assert.Equal(t, "", result) // All spaces
}

func TestCleanEmptyChain(t *testing.T) {
	result := stringx.New("H e llo, W o rld!").CleanEmptyChain().Value()
	assert.Equal(t, "Hello,World!", result)

	result = stringx.New("   ").CleanEmptyChain().Value()
	assert.Equal(t, "", result) // All spaces
}
