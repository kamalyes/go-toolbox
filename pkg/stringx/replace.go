/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 15:15:18
 * @FilePath: \go-toolbox\pkg\stringx\replace.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"regexp"
	"strings"
	"sync"
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
	// 如果 startIndex 小于0，设置为0
	if startIndex < 0 {
		startIndex = 0
	}

	// 如果 endIndex 超出字符串长度，设置为字符串长度
	if endIndex > len(str) {
		endIndex = len(str)
	}

	// 如果 startIndex 大于 endIndex，返回原始字符串
	if startIndex > endIndex {
		return str
	}

	runes := []rune(str)
	replaceCount := endIndex - startIndex

	return string(runes[:startIndex]) + RepeatByLength(replacedStr, replaceCount) + string(runes[endIndex:])
}

// regexCache 已编译正则的缓存，按模式串缓存
//
// 模式是调用方在运行期传进来的，没法在包级预编译，但同一个模式通常会被反复使用
// （循环里脱敏、按同一规则批量替换），而 regexp.MustCompile 要走完整的编译流程，
// 比一次匹配贵得多
//
// 用 sync.Map 而不是 map 加锁：这里是读多写极少的典型形态，模式集合在程序跑起来
// 之后基本固定下来
//
// 注意缓存不设上限如果模式是由不可信输入拼出来的，缓存会随之无限增长；这种用法
// 本身就不该把外部输入当正则（既有正则注入风险，也有灾难性回溯风险），调用方应当
// 只传自己代码里写死的模式
//
// [EN] A cache of compiled regexps keyed by the pattern string. The pattern
// arrives at run time so it cannot be compiled at package level, yet the same
// pattern is typically used over and over — masking inside a loop, batch
// replacement by one rule — and MustCompile runs a full compilation, far dearer
// than a single match. sync.Map rather than a mutex-guarded map because this is
// the classic read-mostly shape, with the set of patterns settling once the
// program is running. The cache is unbounded: if patterns are built from
// untrusted input it will grow without limit, but that usage is wrong to begin
// with — treating external input as a regexp invites both injection and
// catastrophic backtracking — and callers should pass only patterns written in
// their own code
var regexCache sync.Map // map[string]*regexp.Regexp

// cachedRegex 取已编译的正则，未缓存则编译并存入
//
// 与 regexp.MustCompile 一样在模式非法时 panic，保持 ReplaceWithMatcher 的既有行为
func cachedRegex(pattern string) *regexp.Regexp {
	if v, ok := regexCache.Load(pattern); ok {
		return v.(*regexp.Regexp)
	}
	re := regexp.MustCompile(pattern)
	// 并发下可能有多个 goroutine 同时编译同一个模式，存入时谁先谁后无所谓：
	// 编译结果是等价的，*Regexp 本身也可并发使用
	actual, _ := regexCache.LoadOrStore(pattern, re)
	return actual.(*regexp.Regexp)
}

// ReplaceWithMatcher 通过正则表达式替换字符串
//
// 正则会按模式串缓存，同一模式只编译一次，原因见 regexCache
func ReplaceWithMatcher(str string, regex string, replaceFun func(string) string) string {
	return cachedRegex(regex).ReplaceAllStringFunc(str, replaceFun)
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
