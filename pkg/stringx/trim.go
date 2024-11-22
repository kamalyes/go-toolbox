/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:37:52
 * @FilePath: \go-toolbox\pkg\stringx\trim.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"reflect"
	"strings"

	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// Trim 除去字符串头尾部的空白
func Trim(str string) string {
	if validator.IsEmptyValue(reflect.ValueOf(str)) {
		return str
	}
	return strings.TrimSpace(str)
}

// TrimChain 除去字符串头尾部的空白（链式调用）
func (s *StringX) TrimChain() *StringX {
	s.value = Trim(s.value)
	return s
}

// TrimStart 除去字符串头部的空白
func TrimStart(str string) string {
	return strings.TrimLeftFunc(str, func(r rune) bool {
		return r == ' '
	})
}

// TrimStartChain 除去字符串头部的空白（链式调用）
func (s *StringX) TrimStartChain() *StringX {
	s.value = TrimStart(s.value)
	return s
}

// TrimEnd 除去字符串尾部的空白
func TrimEnd(str string) string {
	return strings.TrimRightFunc(str, func(r rune) bool {
		return r == ' '
	})
}

// TrimEndChain 除去字符串尾部的空白（链式调用）
func (s *StringX) TrimEndChain() *StringX {
	s.value = TrimEnd(s.value)
	return s
}

// CleanEmpty 清除空白串
func CleanEmpty(str string) string {
	strRune := []rune(str)
	var newRune []rune
	for _, r := range strRune {
		if r != ' ' {
			newRune = append(newRune, r)
		}
	}
	return string(newRune)
}

// CleanEmptyChain 清除空白串（链式调用）
func (s *StringX) CleanEmptyChain() *StringX {
	s.value = CleanEmpty(s.value)
	return s
}
