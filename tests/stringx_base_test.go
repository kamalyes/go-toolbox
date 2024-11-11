/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 01:58:19
 * @FilePath: \go-toolbox\tests\stringx_base_test.go
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

func TestConvertCase(t *testing.T) {
	tests := []struct {
		input    string
		style    stringx.CharacterStyle
		expected string
	}{
		// 测试蛇形命名法
		{"HelloWorld", stringx.SnakeCharacterStyle, "hello_world"},
		{"helloWorld", stringx.SnakeCharacterStyle, "hello_world"},
		{"Hello_World", stringx.SnakeCharacterStyle, "hello_world"},
		{" Hello World", stringx.SnakeCharacterStyle, "hello_world"},
		{"Hello World", stringx.SnakeCharacterStyle, "hello_world"},
		{" ", stringx.SnakeCharacterStyle, ""}, // 空格测试
		{"", stringx.SnakeCharacterStyle, ""},  // 空字符串测试

		// 测试每个单词首字母大写的风格
		{"hello_world", stringx.StudlyCharacterStyle, "HelloWorld"},
		{"helloWorld", stringx.StudlyCharacterStyle, "HelloWorld"},
		{"hello world", stringx.StudlyCharacterStyle, "HelloWorld"},
		{" Hello World", stringx.StudlyCharacterStyle, "HelloWorld"},
		{"Hello_World", stringx.StudlyCharacterStyle, "HelloWorld"},
		{" ", stringx.StudlyCharacterStyle, ""}, // 空格测试
		{"", stringx.StudlyCharacterStyle, ""},  // 空字符串测试

		// 测试驼峰命名法
		{"hello_world", stringx.CamelCharacterStyle, "helloWorld"},
		{"HelloWorld", stringx.CamelCharacterStyle, "helloWorld"},
		{"hello world", stringx.CamelCharacterStyle, "helloWorld"},
		{" Hello World", stringx.CamelCharacterStyle, "helloWorld"},
		{"Hello_World", stringx.CamelCharacterStyle, "helloWorld"},
		{" ", stringx.CamelCharacterStyle, ""}, // 空格测试
		{"", stringx.CamelCharacterStyle, ""},  // 空字符串测试

		// 测试无效的 CharacterStyle
		{"HelloWorld", stringx.CharacterStyle(999), "HelloWorld"}, // 无效的 caseType 应返回原字符串

	}

	for _, test := range tests {
		result := stringx.ConvertCharacterStyle(test.input, test.style)
		assert.Equal(t, test.expected, result, fmt.Sprintf("ConvertCase(%q, %v) = %q; want %q", test.input, test.style, result, test.expected))
	}
}
