/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 17:26:07
 * @FilePath: \go-toolbox\stringx\trim.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import "strings"

// Trim 除去字符串头尾部的空白
func Trim(str string) string {
	if IsEmpty(str) {
		return str
	}
	return strings.TrimSpace(str)
}

// TrimStart 除去字符串头部的空白
func TrimStart(str string) string {
	return strings.TrimLeftFunc(str, func(r rune) bool {
		if ' ' == r {
			return true
		}
		return false
	})
}

// TrimEnd 除去字符串尾部的空白
func TrimEnd(str string) string {
	return strings.TrimRightFunc(str, func(r rune) bool {
		if ' ' == r {
			return true
		}
		return false
	})
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
