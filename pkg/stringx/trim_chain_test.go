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

func TestTrimProtocolChain(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"http protocol":             {input: "http://example.com", expected: "example.com"},
		"https protocol":            {input: "https://example.com/path", expected: "example.com/path"},
		"ftp protocol":              {input: "ftp://ftp.example.com", expected: "ftp.example.com"},
		"ws protocol":               {input: "ws://example.com:8080", expected: "example.com:8080"},
		"wss protocol":              {input: "wss://example.com/socket", expected: "example.com/socket"},
		"no protocol":               {input: "example.com", expected: "example.com"},
		"empty string":              {input: "", expected: ""},
		"http with trailing spaces": {input: "http://example.com  ", expected: "example.com"},
		"no protocol with spaces":   {input: "  example.com  ", expected: "example.com"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := New(tc.input).TrimProtocolChain().Value()
			assert.Equal(t, tc.expected, result)
		})
	}
}
