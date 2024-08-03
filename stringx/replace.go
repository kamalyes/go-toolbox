/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 10:54:57
 * @FilePath: \go-toolbox\stringx\replace.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Replace 替换字符串
func Replace(source string, searchStr string, replacement string, replaceCount int) string {
	return strings.Replace(source, searchStr, replacement, replaceCount)
}

// ReplaceAll 替换字符串
func ReplaceAll(source string, searchStr string, replacement string) string {
	return Replace(source, searchStr, replacement, -1)
}

// ReplaceWithIndex 按照指定区间替换字符串
// startIndex 开始位置（包含）
// endIndex 结束位置（不包含）
func ReplaceWithIndex(str string, startIndex int, endIndex int, replacedStr string) string {
	if startIndex < 0 || endIndex > len(str) || startIndex > endIndex {
		return str // 如果索引无效，返回原始字符串
	}
	runes := []rune(str)
	strLen := Length(str)

	// 检查并调整 startIndex 和 endIndex 的有效性
	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex > strLen {
		endIndex = strLen
	}
	if startIndex >= endIndex {
		return str
	}

	replaceCount := endIndex - startIndex

	return string(runes[:startIndex]) + RepeatByLength(replacedStr, replaceCount) + string(runes[endIndex:])
}

type PadPosition int

const (
	Left PadPosition = iota
	Right
	Middle
)

type Paddler struct {
	Position PadPosition
}

func SetPadPosition(p *Paddler, position PadPosition) {
	p.Position = position
}

// Pad 输入的字符长度<minLength时自动补位*
func Pad(input string, minLength int, paddler ...*Paddler) string {
	pad := Paddler{Position: Middle}
	charCount := utf8.RuneCountInString(input)

	if charCount >= minLength {
		return input
	}

	padLen := minLength - charCount
	leftPadLen := padLen / 2
	if leftPadLen < 4 {
		leftPadLen = 4
	}

	if len(paddler) > 0 {
		pad = *paddler[0]
		switch pad.Position {
		case Left:
			return strings.Repeat("*", padLen) + input
		case Right:
			return input + strings.Repeat("*", padLen)
		case Middle:
			return input[:leftPadLen] + strings.Repeat("*", padLen) + input[leftPadLen:]
		}
	}

	// 默认为中间填充
	return input[:leftPadLen] + strings.Repeat("*", padLen) + input[leftPadLen:]
}

// ReplaceIgnoreCase 替换字符串
func ReplaceIgnoreCase(source string, searchStr string, replacement string, replaceCount int) string {
	return Replace(strings.ToLower(source), strings.ToLower(searchStr), replacement, replaceCount)
}

// ReplaceAllIgnoreCase 替换字符串
func ReplaceAllIgnoreCase(source string, searchStr string, replacement string) string {
	return Replace(strings.ToLower(source), strings.ToLower(searchStr), replacement, -1)
}

// ReplaceWithMatcher 通过正则表达式替换字符串
func ReplaceWithMatcher(str string, regex string, replaceFun func(string) string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllStringFunc(str, replaceFun)
}

// Hide 替换指定字符串的指定区间内字符为"*" 俗称：脱敏功能
func Hide(str string, startInclude int, endExclude int) string {
	return ReplaceWithIndex(str, startInclude, endExclude, "*")
}
