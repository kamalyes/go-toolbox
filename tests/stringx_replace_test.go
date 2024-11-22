/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-16 08:57:23
 * @FilePath: \go-toolbox\tests\stringx_replace_test.go
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

func TestReplace(t *testing.T) {
	result := stringx.Replace("hello, world", "hello", "hi", 1)
	assert.Equal(t, "hi, world", result)
}

func TestReplaceAll(t *testing.T) {
	result := stringx.ReplaceAll("hello, hello, world", "hello", "hi")
	assert.Equal(t, "hi, hi, world", result)
}

func TestReplaceWithIndex(t *testing.T) {
	result := stringx.ReplaceWithIndex("abcdefghij", 2, 6, "****")
	assert.Equal(t, "ab****ghij", result)
}

func TestPad(t *testing.T) {
	resultDefault := stringx.Pad("hello", 10)
	assert.Equal(t, "hell*****o", resultDefault)

	result := stringx.Pad("hello", 10, &stringx.Paddler{Position: stringx.Middle})
	assert.Equal(t, "hell*****o", result)

	result = stringx.Pad("world", 8, &stringx.Paddler{Position: stringx.Left})
	assert.Equal(t, "***world", result)

	result = stringx.Pad("ok", 5, &stringx.Paddler{Position: stringx.Right})
	assert.Equal(t, "ok***", result)
}

func TestReplaceWithMatcher(t *testing.T) {
	result := stringx.ReplaceWithMatcher("hello 123 world 456", `\d+`, func(s string) string {
		return "xxx"
	})
	assert.Equal(t, "hello xxx world xxx", result)
}

func TestHide(t *testing.T) {
	result := stringx.Hide("password12345", 8, 10)
	assert.Equal(t, "password**345", result)
}

func TestReplaceSpecialChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, World!", "HelloXXWorldX"},
		{"Go is fun.", "GoXisXfunX"},
		{"Special #chars#", "SpecialXXcharsX"},
		{"1234-5678", "1234X5678"},
		{"NoSpecialChars", "NoSpecialChars"},
		{"", ""},
		{"!@#$%^&*()", "XXXXXXXXXX"},
	}

	for _, test := range tests {
		output := stringx.ReplaceSpecialChars(test.input, 'X')
		if output != test.expected {
			t.Errorf("ReplaceSpecialChars(%q) = %q; expected %q", test.input, output, test.expected)
		}
	}
}
