/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 02:29:53
 * @FilePath: \go-toolbox\tests\stringx_trim_test.go
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

func TestAllTrimFunctions(t *testing.T) {
	t.Run("TestTrim", TestTrim)
	t.Run("TestTrimStart", TestTrimStart)
	t.Run("TestTrimEnd", TestTrimEnd)
	t.Run("TestCleanEmpty", TestCleanEmpty)
}

// TestTrim tests the Trim function
func TestTrim(t *testing.T) {
	// Test empty string
	assert.Equal(t, "", stringx.Trim(""))

	// Test string with leading and trailing spaces
	assert.Equal(t, "hello", stringx.Trim("  hello  "))

	// Test string with leading spaces
	assert.Equal(t, "world", stringx.Trim("   world "))

	// Test string with trailing spaces
	assert.Equal(t, "hi", stringx.Trim("hi    "))
}

// TestTrimStart tests the TrimStart function
func TestTrimStart(t *testing.T) {
	// Test empty string
	assert.Equal(t, "", stringx.TrimStart(""))

	// Test string with leading and trailing spaces
	assert.Equal(t, "hello  ", stringx.TrimStart("  hello  "))

	// Test string with leading spaces
	assert.Equal(t, "world ", stringx.TrimStart("   world "))

	// Test string with trailing spaces
	assert.Equal(t, "hi    ", stringx.TrimStart("hi    "))
}

// TestTrimEnd tests the TrimEnd function
func TestTrimEnd(t *testing.T) {
	// Test empty string
	assert.Equal(t, "", stringx.TrimEnd(""))

	// Test string with leading and trailing spaces
	assert.Equal(t, "  hello", stringx.TrimEnd("  hello  "))

	// Test string with leading spaces
	assert.Equal(t, "   world", stringx.TrimEnd("   world "))

	// Test string with trailing spaces
	assert.Equal(t, "hi", stringx.TrimEnd("hi    "))
}

// TestCleanEmpty tests the CleanEmpty function
func TestCleanEmpty(t *testing.T) {
	// Test empty string
	assert.Equal(t, "", stringx.CleanEmpty(""))

	// Test string with spaces
	assert.Equal(t, "helloworld", stringx.CleanEmpty(" hello world "))

	// Test string without spaces
	assert.Equal(t, "hello", stringx.CleanEmpty("hello"))
}
