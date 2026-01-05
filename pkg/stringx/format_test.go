/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:53:15
 * @FilePath: \go-toolbox\pkg\stringx\format_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillBefore(t *testing.T) {
	result := FillBefore("hello", ".", 10)
	assert.Equal(t, ".....hello", result)
}

func TestFillAfter(t *testing.T) {
	result := FillAfter("hello", ".", 10)
	assert.Equal(t, "hello.....", result)
}

func TestFormat(t *testing.T) {
	params := map[string]interface{}{
		"a": "aValue",
		"b": "bValue",
	}
	result := Format("{a} and {b}", params)
	assert.Equal(t, "aValue and bValue", result)
}

func TestIndexedFormat(t *testing.T) {
	result := IndexedFormat("this is {0} for {1}", []interface{}{"a", "b"})
	assert.Equal(t, "this is a for b", result)
}

func TestTruncateAppendEllipsis(t *testing.T) {
	tests := []struct {
		input    string
		maxChars int
		expected string
	}{
		{"这是一个测试字符串199665889@#￥￥", 10, "这是一个测试字符串1..."},
		{"这是一个测试字符串12356789@#￥￥", 50, "这是一个测试字符串12356789@#￥￥"},
		{"", 10, ""},
	}

	for _, test := range tests {
		result := TruncateAppendEllipsis(test.input, test.maxChars)
		if result != test.expected {
			t.Errorf("TruncateAppendEllipsis(%q, %d) = %q; want %q", test.input, test.maxChars, result, test.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	result := Truncate("This is another long string", 10)
	assert.Equal(t, "This is an", result)
}

func TestAddPrefixIfNot(t *testing.T) {
	result := AddPrefixIfNot("world", "hello ")
	assert.Equal(t, "hello world", result)
}

func TestAddSuffixIfNot(t *testing.T) {
	result := AddSuffixIfNot("hello", " world")
	assert.Equal(t, "hello world", result)
}

func TestSanitizeSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "基本转换 - 空格和大写",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "去除特殊字符",
			input:    "Hello World!",
			expected: "hello-world",
		},
		{
			name:     "多个特殊字符",
			input:    "My--Project__123",
			expected: "my-project-123",
		},
		{
			name:     "首尾连字符",
			input:    "  -test-  ",
			expected: "test",
		},
		{
			name:     "连续连字符",
			input:    "hello---world",
			expected: "hello-world",
		},
		{
			name:     "混合大小写和特殊字符",
			input:    "Show-Dev-Platform_Name!@#",
			expected: "show-dev-platform-name",
		},
		{
			name:     "只有字母和数字",
			input:    "abc123XYZ",
			expected: "abc123xyz",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "只有特殊字符",
			input:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "下划线转连字符",
			input:    "hello_world_test",
			expected: "hello-world-test",
		},
		{
			name:     "中文和特殊字符混合",
			input:    "游戏-Game-平台",
			expected: "game",
		},
		{
			name:     "实际场景 - show页面名称",
			input:    "show-dev-Platform Name 123",
			expected: "show-dev-platform-name-123",
		},
		{
			name:     "实际场景 - game页面名称",
			input:    "game-prod-MyGame!@#",
			expected: "game-prod-mygame",
		},
		{
			name:     "多个空格",
			input:    "hello    world",
			expected: "hello-world",
		},
		{
			name:     "Tab和换行符",
			input:    "hello\tworld\ntest",
			expected: "hello-world-test",
		},
		{
			name:     "复杂混合 - 多种分隔符",
			input:    "Hello___World---Test___123",
			expected: "hello-world-test-123",
		},
		{
			name:     "URL中的特殊字符",
			input:    "http://example.com/path?query=1",
			expected: "httpexamplecompathquery1",
		},
		{
			name:     "邮箱格式",
			input:    "user@example.com",
			expected: "userexamplecom",
		},
		{
			name:     "大量连续特殊字符",
			input:    "test!@#$%^&*()test",
			expected: "testtest",
		},
		{
			name:     "Unicode表情符号",
			input:    "Hello 🎮 World 🚀",
			expected: "hello-world",
		},
		{
			name:     "混合中英文数字符号",
			input:    "项目Project_2024@#Version-1.0",
			expected: "project-2024version-10",
		},
		{
			name:     "超长连字符序列",
			input:    "test----------test",
			expected: "test-test",
		},
		{
			name:     "开头结尾都是特殊字符",
			input:    "!!!test!!!",
			expected: "test",
		},
		{
			name:     "只有连字符和空格",
			input:    "--- --- ---",
			expected: "",
		},
		{
			name:     "路径分隔符",
			input:    "path/to/some/file.txt",
			expected: "pathtosomefiletxt",
		},
		{
			name:     "Windows路径",
			input:    "C:\\Users\\Admin\\Documents",
			expected: "cusersadmindocuments",
		},
		{
			name:     "SQL注入尝试",
			input:    "'; DROP TABLE users; --",
			expected: "drop-table-users",
		},
		{
			name:     "HTML标签",
			input:    "<script>alert('test')</script>",
			expected: "scriptalerttestscript",
		},
		{
			name:     "多语言混合",
			input:    "English中文日本語한국어",
			expected: "english",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeSlug(tt.input)
			assert.Equal(t, tt.expected, result, "输入: %q", tt.input)
		})
	}
}

func TestSanitizeSlugEdgeCases(t *testing.T) {
	// 测试极端情况
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "极长字符串",
			input:    strings.Repeat("Hello-World-", 100),
			expected: strings.Repeat("hello-world-", 99) + "hello-world",
		},
		{
			name:     "大量Unicode字符",
			input:    "测试🎮🚀💻⚡️🔥✨🌟",
			expected: "",
		},
		{
			name:     "零宽字符",
			input:    "test\u200B\u200C\u200Dtest",
			expected: "testtest",
		},
		{
			name:     "重音字符",
			input:    "café résumé naïve",
			expected: "caf-rsum-nave",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeSlug(tt.input)
			assert.Equal(t, tt.expected, result, "输入: %q", tt.input)
		})
	}
}

func TestSanitizeSlugPerformance(t *testing.T) {
	// 性能验证测试
	inputs := []string{
		"Simple Test",
		"Complex___Test---With!!!Many@@@Special###Characters",
		strings.Repeat("test-", 1000),
		"中文English日本語한국어Mixed",
	}

	for _, input := range inputs {
		result := SanitizeSlug(input)
		// 验证结果不为nil且是有效字符串
		assert.NotNil(t, result)
		// 验证结果中没有非法字符
		for _, ch := range result {
			assert.True(t, (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-',
				"结果包含非法字符: %c in %q", ch, result)
		}
		// 验证没有连续连字符
		assert.NotContains(t, result, "--", "结果包含连续连字符")
		// 验证首尾没有连字符
		if len(result) > 0 {
			assert.NotEqual(t, '-', rune(result[0]), "结果开头有连字符")
			assert.NotEqual(t, '-', rune(result[len(result)-1]), "结果结尾有连字符")
		}
	}
}

func TestSanitizeSlugChain(t *testing.T) {
	result := New("Hello World!").SanitizeSlugChain().Value()
	assert.Equal(t, "hello-world", result)

	result2 := New("My--Project__123").SanitizeSlugChain().Value()
	assert.Equal(t, "my-project-123", result2)
}

// BenchmarkSanitizeSlug 性能测试
func BenchmarkSanitizeSlug(b *testing.B) {
	testCases := []string{
		"Hello World",
		"My--Project__123",
		"show-dev-Platform Name 123",
		"!@#$%^&*()",
		"abc123XYZ",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			SanitizeSlug(tc)
		}
	}
}

// BenchmarkSanitizeSlugShort 短字符串性能测试
func BenchmarkSanitizeSlugShort(b *testing.B) {
	input := "Hello World"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SanitizeSlug(input)
	}
}

// BenchmarkSanitizeSlugMedium 中等字符串性能测试
func BenchmarkSanitizeSlugMedium(b *testing.B) {
	input := "show-dev-Platform_Name_With___Multiple___Separators"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SanitizeSlug(input)
	}
}

// BenchmarkSanitizeSlugLong 长字符串性能测试
func BenchmarkSanitizeSlugLong(b *testing.B) {
	longString := "This is a very long string with MANY special characters !@#$%^&*() and spaces that needs to be sanitized into a proper slug format"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SanitizeSlug(longString)
	}
}

// BenchmarkSanitizeSlugSpecialChars 大量特殊字符性能测试
func BenchmarkSanitizeSlugSpecialChars(b *testing.B) {
	input := "!!!Hello@@@World###Test$$$123%%%END^^^"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SanitizeSlug(input)
	}
}

// BenchmarkSanitizeSlugUnicode Unicode字符性能测试
func BenchmarkSanitizeSlugUnicode(b *testing.B) {
	input := "Hello世界🎮Test测试💻End"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SanitizeSlug(input)
	}
}

// BenchmarkSanitizeSlugWorstCase 最坏情况：连续分隔符
func BenchmarkSanitizeSlugWorstCase(b *testing.B) {
	input := "test___---___---___test___---___---___end"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SanitizeSlug(input)
	}
}
