/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 16:55:01
 * @FilePath: \go-toolbox\pkg\stringx\index.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"strings"
)

type searchOptions struct {
	start int
	end   int
}

// WithStart 用于配置起始位置
func WithStart(start int) func(*searchOptions) {
	return func(so *searchOptions) {
		so.start = start
	}
}

// WithEnd 用于配置结束位置
func WithEnd(end int) func(*searchOptions) {
	return func(so *searchOptions) {
		so.end = end
	}
}

// SafeIndexOfByRange 在指定范围内查找指定字符，避免下标溢出
func SafeIndexOfByRange(str string, subStr string, options ...func(*searchOptions)) (index int) {
	if subStr == str || len(subStr) == len(str) {
		return 0
	}

	index = -1
	defaultOptions := searchOptions{
		start: 0,
		end:   len(str),
	}

	for _, option := range options {
		option(&defaultOptions)
	}

	if defaultOptions.start < 0 {
		defaultOptions.start = 0
	}

	if defaultOptions.end < 0 {
		defaultOptions.end = 0
	} else if defaultOptions.end > len(str) {
		defaultOptions.end = len(str)
	}

	if defaultOptions.start > len(str) {
		return index
	}

	if defaultOptions.start >= defaultOptions.end {
		return index
	}

	index = strings.Index(str[defaultOptions.start:defaultOptions.end], subStr)
	if index != -1 {
		index += defaultOptions.start
	}
	return
}

// IndexOf 返回字符在原始字符串的下标
func IndexOf(str string, subStr string) int {
	return strings.Index(str, subStr)
}

// IndexOfByRange 指定范围内查找指定字符
func IndexOfByRange(str string, subStr string, start int, end int) int {
	return SafeIndexOfByRange(str, subStr, WithStart(start), WithEnd(end))
}

// IndexOfByRangeStart 指定范围内查找指定字符
func IndexOfByRangeStart(str string, subStr string, start int) int {
	return SafeIndexOfByRange(str, subStr, WithStart(start))
}

// IndexOfIgnoreCase 返回字符在原始字符串的下标(大小写不敏感)
func IndexOfIgnoreCase(str string, subStr string) int {
	return IndexOf(strings.ToLower(str), strings.ToLower(subStr))
}

// IndexOfIgnoreCaseByRange 从指定下标开始，返回在字符串中的下标 (大小不敏感)
func IndexOfIgnoreCaseByRange(str string, subStr string, start int) int {
	return IndexOfByRangeStart(strings.ToLower(str), strings.ToLower(subStr), start)
}

// LastIndexOf 返回最后出现指定字符串的下标
func LastIndexOf(str string, subStr string) int {
	return strings.LastIndex(str, subStr)
}

// LastIndexOfIgnoreCase 返回最后出现指定字符串的下标（大小写不敏感）
func LastIndexOfIgnoreCase(str string, subStr string) int {
	return LastIndexOf(strings.ToLower(str), strings.ToLower(subStr))
}

// LastIndexOfByRangeStart 从指定下标开始，返回最后出现指定字符串的下标
func LastIndexOfByRangeStart(str string, subStr string, start int) int {
	return SafeIndexOfByRange(str, subStr, WithStart(start))
}

// OrdinalIndexOf 返回字符串 subStr 在字符串 str 中第 ordinal 次出现的位置。
// 如果 str="" 或 subStr=" 或 ordinal≥0 则返回-1
func OrdinalIndexOf(str string, subStr string, ordinal int, start ...int) int {
	if subStr == str || len(subStr) == len(str) {
		return 0
	}
	findIndex := 0
	ordinalIndex := 0
	for i := 0; i < len(str); i++ {
		idx := SafeIndexOfByRange(str, subStr, WithStart(findIndex), WithEnd(len(str)-findIndex))
		switch idx {
		case -1:
			return -1
		}
		findIndex = idx + 1
		ordinalIndex += 1
		if ordinalIndex == ordinal {
			break
		}
	}
	return findIndex - 1
}
