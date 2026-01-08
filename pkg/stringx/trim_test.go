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
	t.Run("TestTrimProtocol", TestTrimProtocol)
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

// TestTrimProtocol tests the TrimProtocol function
func TestTrimProtocol(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"empty string":              {input: "", expected: ""},
		"http protocol":             {input: "http://example.com", expected: "example.com"},
		"https protocol":            {input: "https://example.com/path", expected: "example.com/path"},
		"ftp protocol":              {input: "ftp://ftp.example.com", expected: "ftp.example.com"},
		"ws protocol":               {input: "ws://example.com:8080", expected: "example.com:8080"},
		"wss protocol":              {input: "wss://example.com/socket", expected: "example.com/socket"},
		"file protocol":             {input: "file://path/to/file", expected: "path/to/file"},
		"no protocol":               {input: "example.com", expected: "example.com"},
		"no protocol with path":     {input: "example.com/path/to/page", expected: "example.com/path/to/page"},
		"custom protocol":           {input: "custom://custom.example.com", expected: "custom.example.com"},
		"http with trailing spaces": {input: "http://example.com  ", expected: "example.com"},
		"https with leading spaces": {input: "https://  example.com", expected: "example.com"},
		"no protocol with spaces":   {input: "  example.com  ", expected: "example.com"},
		"protocol with both spaces": {input: "  https://example.com  ", expected: "example.com"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, TrimProtocol(tc.input))
		})
	}
}
