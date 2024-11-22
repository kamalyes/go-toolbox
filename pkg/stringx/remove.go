/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:07:57
 * @FilePath: \go-toolbox\pkg\stringx\remove.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"regexp"
	"strings"

	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// RemoveAll 移除字符串中所有给定字符串;removeAll("aa-bb-cc-dd", "-") =》 aabbccdd
func RemoveAll(str string, strToRemove string) string {
	if str == "" || strToRemove == "" {
		return str
	}
	return ReplaceAll(str, strToRemove, "")
}

// RemoveAllChain 移除字符串中所有给定字符串（链式调用）
func (s *StringX) RemoveAllChain(strToRemove string) *StringX {
	s.value = RemoveAll(s.value, strToRemove)
	return s
}

// RemoveSymbols 使用正则表达式去掉字符串中的所有符号
func RemoveSymbols(str string) string {
	reg := regexp.MustCompile(`[^\w]+`)
	return reg.ReplaceAllString(str, "")
}

// RemoveSymbolsChain 使用正则表达式去掉字符串中的所有符号（链式调用）
func (s *StringX) RemoveSymbolsChain() *StringX {
	s.value = RemoveSymbols(s.value)
	return s
}

// RemoveAny 移除字符串中所有给定字符串，当某个字符串出现多次，则全部移除
func RemoveAny(str string, strsToRemove []string) string {
	var result = str
	hasEmpty, _ := validator.HasEmpty([]interface{}{str})
	if !hasEmpty {
		for _, s := range strsToRemove {
			result = RemoveAll(result, s)
		}
	}
	return result
}

// RemoveAnyChain 移除字符串中所有给定字符串（链式调用）
func (s *StringX) RemoveAnyChain(strsToRemove []string) *StringX {
	s.value = RemoveAny(s.value, strsToRemove)
	return s
}

// RemoveAllLineBreaks 去除所有换行符，包括：\r \n
func RemoveAllLineBreaks(str string) string {
	return RemoveAny(str, []string{"\r", "\n"})
}

// RemoveAllLineBreaksChain 去除所有换行符（链式调用）
func (s *StringX) RemoveAllLineBreaksChain() *StringX {
	s.value = RemoveAllLineBreaks(s.value)
	return s
}

// RemovePrefix 去掉指定前缀
func RemovePrefix(str string, prefix string) string {
	hasEmpty, _ := validator.HasEmpty([]interface{}{str, prefix})
	if hasEmpty {
		return str
	}
	if strings.HasPrefix(str, prefix) {
		return str[len(prefix):]
	}
	return str
}

// RemovePrefixChain 去掉指定前缀（链式调用）
func (s *StringX) RemovePrefixChain(prefix string) *StringX {
	s.value = RemovePrefix(s.value, prefix)
	return s
}

// RemovePrefixIgnoreCase 忽略大小写去掉指定前缀
func RemovePrefixIgnoreCase(str string, prefix string) string {
	if strings.HasPrefix(strings.ToLower(str), strings.ToLower(prefix)) {
		return str[len(prefix):]
	}
	return str
}

// RemovePrefixIgnoreCaseChain 忽略大小写去掉指定前缀（链式调用）
func (s *StringX) RemovePrefixIgnoreCaseChain(prefix string) *StringX {
	s.value = RemovePrefixIgnoreCase(s.value, prefix)
	return s
}

// RemoveSuffix 去掉指定后缀
func RemoveSuffix(str string, suffix string) string {
	hasEmpty, _ := validator.HasEmpty([]interface{}{str, suffix})
	if hasEmpty {
		return str
	}
	if strings.HasSuffix(str, suffix) {
		return str[:len(str)-len(suffix)]
	}
	return str
}

// RemoveSuffixChain 去掉指定后缀（链式调用）
func (s *StringX) RemoveSuffixChain(suffix string) *StringX {
	s.value = RemoveSuffix(s.value, suffix)
	return s
}

// RemoveSuffixIgnoreCase 去掉指定后缀(忽略大小写)
func RemoveSuffixIgnoreCase(str string, suffix string) string {
	lowerStr := strings.ToLower(str)
	lowerSuffix := strings.ToLower(suffix)
	if strings.HasSuffix(lowerStr, lowerSuffix) {
		return str[:len(str)-len(suffix)]
	}
	return str
}

// RemoveSuffixIgnoreCaseChain 去掉指定后缀（链式调用，忽略大小写）
func (s *StringX) RemoveSuffixIgnoreCaseChain(suffix string) *StringX {
	s.value = RemoveSuffixIgnoreCase(s.value, suffix)
	return s
}
