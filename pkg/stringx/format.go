/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:07:57
 * @FilePath: \go-toolbox\pkg\stringx\format.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"fmt"
	"strings"
	"unicode"
)

// FillBefore 将已有字符串填充为规定长度，如果已有字符串超过这个长度则返回这个字符串
func FillBefore(str string, char string, length int) string {
	return Fill(str, char, length, true)
}

// FillBeforeChain 将已有字符串填充为规定长度（链式调用）
func (s *StringX) FillBeforeChain(char string, length int) *StringX {
	s.value = FillBefore(s.value, char, length)
	return s
}

// FillAfter 将已有字符串填充为规定长度，如果已有字符串超过这个长度则返回这个字符串
func FillAfter(str string, char string, length int) string {
	return Fill(str, char, length, false)
}

// FillAfterChain 将已有字符串填充为规定长度（链式调用）
func (s *StringX) FillAfterChain(char string, length int) *StringX {
	s.value = FillAfter(s.value, char, length)
	return s
}

// Fill 将已有字符串填充为规定长度，如果已有字符串超过这个长度则返回这个字符串
func Fill(str string, char string, length int, isPre bool) string {
	if len(str) >= length {
		return str
	}
	fillLength := length - len(str)
	fillStr := strings.Repeat(char, fillLength)
	if isPre {
		return fillStr + str
	}
	return str + fillStr
}

// Format 通过map中的参数 格式化字符串
// map = {a: "aValue", b: "bValue"} format("{a} and {b}", map) ---=》 aValue and bValue
func Format(template string, params map[string]interface{}) string {
	// 遍历map中的键值对
	for key, value := range params {
		// 构造占位符，例如 "{a}"
		placeholder := fmt.Sprintf("{%s}", key)
		// 将占位符替换为对应的值
		template = strings.ReplaceAll(template, placeholder, fmt.Sprintf("%v", value))
	}
	return template
}

// FormatChain 通过map中的参数格式化字符串（链式调用）
func (s *StringX) FormatChain(params map[string]interface{}) *StringX {
	s.value = Format(s.value, params)
	return s
}

// IndexedFormat 有序的格式化文本，使用{number}做为占位符
// 通常使用：format("this is {0} for {1}", "a", "b") =》 this is a for b
func IndexedFormat(template string, params []interface{}) string {
	// 遍历所有参数
	for i, param := range params {
		placeholder := fmt.Sprintf("{%d}", i)
		// 将占位符替换为对应的值
		template = strings.ReplaceAll(template, placeholder, fmt.Sprintf("%v", param))
	}
	return template
}

// IndexedFormatChain 有序的格式化文本（链式调用）
func (s *StringX) IndexedFormatChain(params []interface{}) *StringX {
	s.value = IndexedFormat(s.value, params)
	return s
}

// TruncateAppendEllipsis 截断字符串，使用不超过maxChars字符长度。截断后自动追加省略号(...) 用于存储数据库varchar且编码为UTF-8的字段
func TruncateAppendEllipsis(str string, maxChars int) string {
	// 如果字符串本身就比 maxChars 短，则不需要截断
	if len(str) <= maxChars {
		return str
	}

	// 初始化变量
	runes := []rune(str) // 将字符串转换为 rune 切片
	if len(runes) > maxChars {
		runes = runes[:maxChars] // 截断到 maxChars 个字符
	}

	return string(runes) + "..." // 返回截断字符串加省略号
}

// TruncateAppendEllipsisChain 截断字符串（链式调用）
func (s *StringX) TruncateAppendEllipsisChain(maxBytes int) *StringX {
	s.value = TruncateAppendEllipsis(s.value, maxBytes)
	return s
}

// Truncate 截断字符串，使用不超过maxBytes长度
func Truncate(str string, maxBytes int) string {
	return str[0:maxBytes]
}

// TruncateChain 截断字符串（链式调用）
func (s *StringX) TruncateChain(maxBytes int) *StringX {
	s.value = Truncate(s.value, maxBytes)
	return s
}

// AddPrefixIfNot 如果给定字符串不是以prefix开头的，在开头补充 prefix
func AddPrefixIfNot(str string, prefix string) string {
	if StartWith(str, prefix) {
		return str
	}
	return prefix + str
}

// AddPrefixIfNotChain 如果给定字符串不是以prefix开头的（链式调用）
func (s *StringX) AddPrefixIfNotChain(prefix string) *StringX {
	s.value = AddPrefixIfNot(s.value, prefix)
	return s
}

// AddSuffixIfNot 如果给定字符串不是以suffix结尾的，在尾部补充 suffix
func AddSuffixIfNot(str string, suffix string) string {
	if EndWith(str, suffix) {
		return str
	}
	return str + suffix
}

// AddSuffixIfNotChain 如果给定字符串不是以suffix结尾的，在尾部补充 suffix（链式调用）
func (s *StringX) AddSuffixIfNotChain(prefix string) *StringX {
	s.value = AddSuffixIfNot(s.value, prefix)
	return s
}

// SanitizeSlug 格式化字符串为 URL 友好的 slug 格式
// 规则：小写、去除空格、只保留字母数字和连字符、去除连续连字符、去除首尾连字符
// 适用场景：项目名称、域名、URL slug 等
//
// 示例：
//
//	SanitizeSlug("Hello World!")      // "hello-world"
//	SanitizeSlug("My--Project__123")  // "my-project-123"
//	SanitizeSlug("  -test-  ")        // "test"
func SanitizeSlug(name string) string {
	if name == "" {
		return ""
	}

	// 预分配 Builder，避免多次内存分配
	var builder strings.Builder
	builder.Grow(len(name))

	lastWasHyphen := true // 用于跟踪上一个字符是否为连字符，初始为 true 以去除开头连字符

	for _, r := range name {
		// 将字符转为小写
		r = unicode.ToLower(r)

		// 判断字符类型并处理
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			// 字母或数字，直接添加
			builder.WriteRune(r)
			lastWasHyphen = false
		} else if r == '-' || r == '_' || r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			// 分隔符转换为连字符，但避免连续连字符
			if !lastWasHyphen && builder.Len() > 0 {
				builder.WriteRune('-')
				lastWasHyphen = true
			}
		}
		// 其他特殊字符直接忽略
	}

	result := builder.String()

	// 去除尾部连字符（开头连字符已在循环中处理）
	if len(result) > 0 && result[len(result)-1] == '-' {
		result = result[:len(result)-1]
	}

	return result
}

// SanitizeSlugChain 格式化字符串为 URL 友好的 slug 格式（链式调用）
func (s *StringX) SanitizeSlugChain() *StringX {
	s.value = SanitizeSlug(s.value)
	return s
}
