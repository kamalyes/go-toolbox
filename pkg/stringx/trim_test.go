/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 02:29:53
 * @FilePath: \go-toolbox\pkg\stringx\trim_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

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
	assert.Equal(t, "", Trim(""))

	// Test string with leading and trailing spaces
	assert.Equal(t, "hello", Trim("  hello  "))

	// Test string with leading spaces
	assert.Equal(t, "world", Trim("   world "))

	// Test string with trailing spaces
	assert.Equal(t, "hi", Trim("hi    "))
}

// TestTrimStart tests the TrimStart function
func TestTrimStart(t *testing.T) {
	// Test empty string
	assert.Equal(t, "", TrimStart(""))

	// Test string with leading and trailing spaces
	assert.Equal(t, "hello  ", TrimStart("  hello  "))

	// Test string with leading spaces
	assert.Equal(t, "world ", TrimStart("   world "))

	// Test string with trailing spaces
	assert.Equal(t, "hi    ", TrimStart("hi    "))
}

// TestTrimEnd tests the TrimEnd function
func TestTrimEnd(t *testing.T) {
	// Test empty string
	assert.Equal(t, "", TrimEnd(""))

	// Test string with leading and trailing spaces
	assert.Equal(t, "  hello", TrimEnd("  hello  "))

	// Test string with leading spaces
	assert.Equal(t, "   world", TrimEnd("   world "))

	// Test string with trailing spaces
	assert.Equal(t, "hi", TrimEnd("hi    "))
}

// TestCleanEmpty tests the CleanEmpty function
func TestCleanEmpty(t *testing.T) {
	// Test empty string
	assert.Equal(t, "", CleanEmpty(""))

	// Test string with spaces
	assert.Equal(t, "helloworld", CleanEmpty(" hello world "))

	// Test string without spaces
	assert.Equal(t, "hello", CleanEmpty("hello"))
}
