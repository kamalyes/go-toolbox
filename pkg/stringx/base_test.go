/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:21:29
 * @FilePath: \go-toolbox\pkg\stringx\base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestLength(t *testing.T) {
	result := Length("Hello, World!")
	assert.Equal(t, 13, result)
}

func TestReverse(t *testing.T) {
	result := Reverse("Hello, World!")
	assert.Equal(t, "!dlroW ,olleH", result)
}

func TestEquals(t *testing.T) {
	result := Equals("hello", "hello")
	assert.True(t, result)
}

func TestEqualsIgnoreCase(t *testing.T) {
	result := EqualsIgnoreCase("HELLO", "hello")
	assert.True(t, result)
}

func TestInsertSpaces(t *testing.T) {
	result := InsertSpaces("1234567890", 2)
	assert.Equal(t, "12 34 56 78 90", result)
}

func TestEqualsAny(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := EqualsAny("banana", strList)
	assert.True(t, result)
}

func TestEqualsAnyIgnoreCase(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := EqualsAnyIgnoreCase("OrAnGe", strList)
	assert.True(t, result)
}

func TestEqualsAt(t *testing.T) {
	result := EqualsAt("hello", 1, "e")
	assert.True(t, result)
}

func TestCount(t *testing.T) {
	result := Count("banana", "a")
	assert.Equal(t, 3, result)
}

func TestCompareIgnoreCase(t *testing.T) {
	result := CompareIgnoreCase("apple", "BANANA")
	assert.Less(t, result, 0)
}

func TestConvertCase(t *testing.T) {
	tests := []struct {
		input    string
		style    CharacterStyle
		expected string
	}{
		// 测试蛇形命名法
		{"HelloWorld", SnakeCharacterStyle, "hello_world"},
		{"helloWorld", SnakeCharacterStyle, "hello_world"},
		{"Hello_World", SnakeCharacterStyle, "hello_world"},
		{" Hello World", SnakeCharacterStyle, "hello_world"},
		{"Hello World", SnakeCharacterStyle, "hello_world"},
		{" ", SnakeCharacterStyle, ""}, // 空格测试
		{"", SnakeCharacterStyle, ""},  // 空字符串测试

		// 测试每个单词首字母大写的风格
		{"hello_world", StudlyCharacterStyle, "HelloWorld"},
		{"helloWorld", StudlyCharacterStyle, "HelloWorld"},
		{"hello world", StudlyCharacterStyle, "HelloWorld"},
		{" Hello World", StudlyCharacterStyle, "HelloWorld"},
		{"Hello_World", StudlyCharacterStyle, "HelloWorld"},
		{" ", StudlyCharacterStyle, ""}, // 空格测试
		{"", StudlyCharacterStyle, ""},  // 空字符串测试

		// 测试驼峰命名法
		{"hello_world", CamelCharacterStyle, "helloWorld"},
		{"HelloWorld", CamelCharacterStyle, "helloWorld"},
		{"hello world", CamelCharacterStyle, "helloWorld"},
		{" Hello World", CamelCharacterStyle, "helloWorld"},
		{"Hello_World", CamelCharacterStyle, "helloWorld"},
		{" ", CamelCharacterStyle, ""}, // 空格测试
		{"", CamelCharacterStyle, ""},  // 空字符串测试

		// 测试无效的 CharacterStyle
		{"HelloWorld", CharacterStyle(999), "HelloWorld"}, // 无效的 caseType 应返回原字符串

	}

	for _, test := range tests {
		result := ConvertCharacterStyle(test.input, test.style)
		assert.Equal(t, test.expected, result, fmt.Sprintf("ConvertCase(%q, %v) = %q; want %q", test.input, test.style, result, test.expected))
	}
}

func TestToInt(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		input    string // 输入字符串
		expected int    // 预期输出
		err      bool   // 预期错误
	}{
		{" 123 ", 123, false}, // 测试正常的正整数
		{"-456", -456, false}, // 测试负整数
		{"0", 0, false},       // 测试零
		{"abc", 0, true},      // 测试非数字输入
		{"", 0, true},         // 测试空字符串
		{"   ", 0, true},      // 测试仅有空白的字符串
	}

	// 遍历每个测试用例
	for _, test := range tests {
		result, err := ToInt(test.input) // 调用 ToInt 函数

		// 使用断言检查错误
		if test.err {
			assert.Error(t, err, "输入: %q 期望错误: %v", test.input, test.err)
		} else {
			assert.NoError(t, err, "输入: %q 不应返回错误", test.input)
		}

		// 使用断言检查返回值
		assert.Equal(t, test.expected, result, "输入: %q 返回值不匹配", test.input)
	}
}

// 测试 FindKeysByValue 函数
func TestFindKeysByValue(t *testing.T) {
	data := map[string]string{
		"key1": "value1, value2, value3",
		"key2": "value4, value5",
		"key3": "value1, value6",
	}

	tests := []struct {
		searchValue string
		expected    []string
	}{
		{"value1", []string{"key1", "key3"}},
		{"value4", []string{"key2"}},
		{"value5", []string{"key2"}},
		{"value6", []string{"key3"}},
		{"value7", []string{}}, // 没有匹配的值
	}

	for _, test := range tests {
		t.Run(test.searchValue, func(t *testing.T) {
			result := FindKeysByValue(data, test.searchValue)
			assert.ElementsMatch(t, test.expected, result, "对于搜索值 %s，期望 %v 但实际得到 %v", test.searchValue, test.expected, result)
		})
	}
}

func TestStringToSliceByte(t *testing.T) {
	s := "hello world"
	b := ToSliceByte(s)

	if string(b) != s {
		t.Errorf("expected '%s', got '%s'", s, b)
	}
}

func TestTruncateMessage(t *testing.T) {
	// 测试用例
	tests := []struct {
		content  string
		maxLen   int
		expected string
	}{
		{"Hello, World!", 5, "Hello..."},
		{"Go is awesome!", 20, "Go is awesome!"},
		{"", 5, ""},
		{"Very long message that exceeds the limit", 10, "Very long ..."},
		{"Short", 10, "Short"},
	}

	for _, test := range tests {
		t.Run(test.content, func(t *testing.T) {
			result := TruncateMessage(test.content, test.maxLen)
			assert.Equal(t, test.expected, result)
		})
	}
}

// 测试 UniqueStringSlice 函数
func TestUniqueStringSlice(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"apple", "banana", "apple", "orange", "banana", ""},
			expected: []string{"apple", "banana", "orange"},
		},
		{
			input:    []string{"", "", "", ""},
			expected: []string{},
		},
		{
			input:    []string{"apple", "banana", "orange"},
			expected: []string{"apple", "banana", "orange"},
		},
		{
			input:    []string{"apple", "apple", "apple"},
			expected: []string{"apple"},
		},
		{
			input:    []string{"banana", "banana", "", "apple", "orange", "banana"},
			expected: []string{"banana", "apple", "orange"},
		},
		{
			input:    []string{},
			expected: []string{},
		},
	}

	for _, test := range tests {
		result := UniqueStringSlice(test.input)
		assert.ElementsMatch(t, test.expected, result, "Input: %v", test.input)
	}
}

// TestNormalizeDirectionEdgeCases 测试 normalizeDirection 的边界情况
func TestNormalizeDirectionEdgeCases(t *testing.T) {
	t.Run("空字符串", func(t *testing.T) {
		result := NormalizeSQLDirection("", "ASC")
		assert.Equal(t, "ASC", result)
	})

	t.Run("ASC小写", func(t *testing.T) {
		result := NormalizeSQLDirection("asc", "DESC")
		assert.Equal(t, "ASC", result)
	})

	t.Run("DESC小写", func(t *testing.T) {
		result := NormalizeSQLDirection("desc", "ASC")
		assert.Equal(t, "DESC", result)
	})

	t.Run("无效方向返回默认值", func(t *testing.T) {
		result := NormalizeSQLDirection("invalid", "DESC")
		assert.Equal(t, "DESC", result)
	})

	t.Run("空字符串使用DESC默认值", func(t *testing.T) {
		result := NormalizeSQLDirection("", "DESC")
		assert.Equal(t, "DESC", result)
	})
}

// legacyToPascalCase 单遍改写之前的实现，仅供差分比对
func legacyToPascalCase(s string) string {
	if s == "" {
		return s
	}
	if strings.ContainsAny(s, "_- ") {
		normalized := strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), " ", "_")
		parts := strings.Split(normalized, "_")
		var builder strings.Builder
		for _, part := range parts {
			if len(part) > 0 {
				builder.WriteRune(unicode.ToUpper(rune(part[0])))
				if len(part) > 1 {
					builder.WriteString(strings.ToLower(part[1:]))
				}
			}
		}
		return builder.String()
	}
	if unicode.IsLower(rune(s[0])) {
		return string(unicode.ToUpper(rune(s[0]))) + s[1:]
	}
	return s
}

// TestToPascalCase_MatchesLegacy 单遍实现与原 ReplaceAll+Split 实现结果必须一致
//
// 分隔符处理最容易走形：连续分隔符、首尾分隔符，以及按字节取首字符这件事
// 任一处偏差都会让按变体查找的调用方少一条命中路径，而那是静默失败
func TestToPascalCase_MatchesLegacy(t *testing.T) {
	inputs := []string{
		"", "a", "A", "_", "-", " ", "__", "--", "  ", "_-", "- ",
		"a_b", "a-b", "a b", "A_B", "a__b", "a--b", "a  b", "a_-b", "a -b",
		"_a", "a_", "_a_", "-a-", " a ", "__a__",
		"user_name", "user-name", "user name", "USER_NAME", "User-Name",
		"userName", "UserName", "URLPath", "HTTP_server", "id_1", "x_1_y",
		"中文_字段", "_中文", "中文 name", "ab_CD_ef", "a_B_c",
	}
	for _, in := range inputs {
		assert.Equal(t, legacyToPascalCase(in), ToPascalCase(in),
			"输入 %q 的 PascalCase 结果与原实现不一致", in)
	}
}

// legacyNormalizeFieldName 去重修正之前的实现，仅供差分比对
func legacyNormalizeFieldName(fieldName string) []string {
	if fieldName == "" {
		return []string{}
	}
	variants := make([]string, 0, 5)
	variants = append(variants, fieldName)
	pascalCase := ToPascalCase(fieldName)
	if pascalCase != fieldName {
		variants = append(variants, pascalCase)
	}
	camelCase := ToCamelCase(fieldName)
	if camelCase != fieldName && camelCase != pascalCase {
		variants = append(variants, camelCase)
	}
	snakeCase := ToSnakeCase(fieldName)
	if snakeCase != fieldName {
		variants = append(variants, snakeCase)
	}
	kebabCase := ToKebabCase(fieldName)
	if kebabCase != fieldName && kebabCase != snakeCase {
		variants = append(variants, kebabCase)
	}
	return variants
}

// TestNormalizeFieldName_MatchesLegacyAfterDedup 新实现等于原实现去重后的结果
//
// 少一个变体就是少一条命中路径，而那种失败是静默的，所以逐个输入比对
func TestNormalizeFieldName_MatchesLegacyAfterDedup(t *testing.T) {
	dedup := func(in []string) []string {
		out := make([]string, 0, len(in))
		for _, v := range in {
			dup := false
			for _, u := range out {
				if u == v {
					dup = true
					break
				}
			}
			if !dup {
				out = append(out, v)
			}
		}
		return out
	}

	inputs := []string{
		"", "a", "A", "_", "-", " ", "__", "ab", "AB", "aB", "Ab",
		"userName", "UserName", "user_name", "user-name", "USER_NAME",
		"user name", "User Name", "userNAME", "URLPath", "HTTPServer",
		"id", "ID", "iD", "x1", "X1", "a_b-c d", "trailing_", "_leading",
		"中文字段", "mixed中文Name", "a__b", "a--b", "a  b",
	}
	for _, in := range inputs {
		assert.Equal(t, dedup(legacyNormalizeFieldName(in)), NormalizeFieldName(in),
			"输入 %q 的变体列表与原实现（去重后）不一致", in)
	}

	// "X1" 原实现本有重复项，现应被去掉
	assert.Equal(t, []string{"X1", "x1"}, NormalizeFieldName("X1"))
	assert.Len(t, legacyNormalizeFieldName("X1"), 3, "原实现确实多一个重复项")
}
