/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-16 08:57:23
 * @FilePath: \go-toolbox\stringx\replace_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllReplaceFunctions(t *testing.T) {
	t.Run("TestReplace", TestReplace)
	t.Run("TestReplaceAll", TestReplaceAll)
	t.Run("TestReplaceWithIndex", TestReplaceWithIndex)
	t.Run("TestPad", TestPad)
	t.Run("TestReplaceIgnoreCase", TestReplaceIgnoreCase)
	t.Run("TestEndWithIgnoreCase", TestEndWithIgnoreCase)
	t.Run("TestReplaceWithMatcher", TestReplaceWithMatcher)
	t.Run("TestHide", TestHide)
	t.Run("TestReplaceSpecialChars", TestReplaceSpecialChars)

}

func TestReplace(t *testing.T) {
	result := Replace("hello, world", "hello", "hi", 1)
	assert.Equal(t, "hi, world", result)
}

func TestReplaceAll(t *testing.T) {
	result := ReplaceAll("hello, hello, world", "hello", "hi")
	assert.Equal(t, "hi, hi, world", result)
}

func TestReplaceWithIndex(t *testing.T) {
	result := ReplaceWithIndex("abcdefghij", 2, 6, "****")
	assert.Equal(t, "ab****ghij", result)
}

func TestPad(t *testing.T) {
	resultDefault := Pad("hello", 10)
	assert.Equal(t, "hell*****o", resultDefault)

	result := Pad("hello", 10, &Paddler{Position: Middle})
	assert.Equal(t, "hell*****o", result)

	result = Pad("world", 8, &Paddler{Position: Left})
	assert.Equal(t, "***world", result)

	result = Pad("ok", 5, &Paddler{Position: Right})
	assert.Equal(t, "ok***", result)
}

func TestReplaceIgnoreCase(t *testing.T) {
	result := ReplaceIgnoreCase("Hello, World", "hello", "hi", 1)
	assert.Equal(t, "hi, world", result)
}

func TestReplaceWithMatcher(t *testing.T) {
	result := ReplaceWithMatcher("hello 123 world 456", `\d+`, func(s string) string {
		return "xxx"
	})
	assert.Equal(t, "hello xxx world xxx", result)
}

func TestHide(t *testing.T) {
	result := Hide("password12345", 8, 10)
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
		output := ReplaceSpecialChars(test.input, 'X')
		if output != test.expected {
			t.Errorf("ReplaceSpecialChars(%q) = %q; expected %q", test.input, output, test.expected)
		}
	}
}
