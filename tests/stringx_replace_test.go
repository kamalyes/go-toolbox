/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 15:17:25
 * @FilePath: \go-toolbox\tests\stringx_replace_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
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
	tests := []struct {
		input       string
		startIndex  int
		endIndex    int
		replacedStr string
		expected    string
	}{
		{"hello world", 6, 11, "*", "hello *****"},
		{"hello world", 0, 5, "*", "***** world"},
		{"hello world", 0, 0, "*", "hello world"}, // 替换长度为0
		{"hello world", 5, 5, "*", "hello world"}, // 替换长度为0
		{"hello world", 5, 6, "*", "hello*world"}, // 替换单个字符
		{"hello", 1, 4, "X", "hXXXo"},             // 替换多个字符
		{"", 0, 1, "*", ""},                       // 空字符串
		{"test", -1, 2, "#", "##st"},              // startIndex < 0
		{"test", 1, 10, "#", "t###"},              // endIndex 超出范围
		{"test", 3, 1, "#", "test"},               // startIndex > endIndex
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("input: %s, start: %d, end: %d", test.input, test.startIndex, test.endIndex), func(t *testing.T) {
			result := stringx.ReplaceWithIndex(test.input, test.startIndex, test.endIndex, test.replacedStr)
			assert.Equal(t, test.expected, result)
		})
	}
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
