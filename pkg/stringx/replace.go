/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:07:57
 * @FilePath: \go-toolbox\pkg\stringx\replace.go
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

// ReplaceChain 替换字符串（链式调用）
func (s *StringX) ReplaceChain(searchStr string, replacement string, replaceCount int) *StringX {
	s.value = Replace(s.value, searchStr, replacement, replaceCount)
	return s
}

// ReplaceAll 替换字符串
func ReplaceAll(source string, searchStr string, replacement string) string {
	return Replace(source, searchStr, replacement, -1)
}

// ReplaceAllChain 替换字符串（链式调用）
func (s *StringX) ReplaceAllChain(searchStr string, replacement string) *StringX {
	s.value = ReplaceAll(s.value, searchStr, replacement)
	return s
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

// ReplaceWithMatcher 通过正则表达式替换字符串
func ReplaceWithMatcher(str string, regex string, replaceFun func(string) string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllStringFunc(str, replaceFun)
}

// ReplaceWithMatcherChain 通过正则表达式替换字符串（链式调用）
func (s *StringX) ReplaceWithMatcherChain(regex string, replaceFun func(string) string) *StringX {
	s.value = ReplaceWithMatcher(s.value, regex, replaceFun)
	return s
}

// Hide 替换指定字符串的指定区间内字符为"*" 俗称：脱敏功能
func Hide(str string, startInclude int, endExclude int) string {
	return ReplaceWithIndex(str, startInclude, endExclude, "*")
}

// HideChain 替换指定字符串的指定区间内字符为"*"（链式调用）
func (s *StringX) HideChain(startInclude int, endExclude int) *StringX {
	s.value = Hide(s.value, startInclude, endExclude)
	return s
}

// ReplaceSpecialChars 去掉特殊符号、转为自定义
func ReplaceSpecialChars(str string, replaceValue rune) string {
	// 定义不同类别的特殊字符
	englishPunctuation := `!"#$%&'()*+,-./:;<=>?@[\\]^_` + "`" + `{|}~`
	chinesePunctuation := `，。！？；：“”‘’《》`
	otherSpecialChars := `【】〔〕…· `

	// 将所有特殊字符组合在一起
	specialChars := englishPunctuation + chinesePunctuation + otherSpecialChars
	// 使用 Map 函数将标点符号和特殊字符替换为自定义
	cleanedStr := strings.Map(func(r rune) rune {
		if strings.ContainsRune(specialChars, r) {
			return replaceValue
		}
		return r // 保留非特殊字符
	}, str)
	return cleanedStr
}

// ReplaceSpecialCharsChain 去掉特殊符号、转为自定义（链式调用）
func (s *StringX) ReplaceSpecialCharsChain(replaceValue rune) *StringX {
	s.value = ReplaceSpecialChars(s.value, replaceValue)
	return s
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
