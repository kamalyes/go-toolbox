/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 17:00:34
 * @FilePath: \go-toolbox\stringx\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"reflect"
	"strings"

	"github.com/kamalyes/go-toolbox/validator"
)

// Length 计算长度
func Length(str string) int {
	if validator.IsEmptyValue(reflect.ValueOf(str)) {
		return 0
	}
	strRune := []rune(str)
	return len(strRune)
}

// Reverse 反转给定的字符串
func Reverse(str string) string {
	// 将字符串转换为 rune 切片，以处理 Unicode 字符
	runes := []rune(str)
	n := len(runes)

	// 反转 rune 切片
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}

	// 返回反转后的字符串
	return string(runes)
}

// Equals 比较2个字符串（大小写敏感）
func Equals(str1 string, str2 string) bool {
	return str1 == str2
}

// EqualsIgnoreCase 比较2个字符串（大小写不敏感）
func EqualsIgnoreCase(str1 string, str2 string) bool {
	return strings.EqualFold(str1, str2)
}

// InsertSpaces 插入空值
// InsertSpaces 插入空值
func InsertSpaces(str string, interval int) string {
	var buffer strings.Builder
	count := 0
	totalChars := 0

	for _, char := range str {
		buffer.WriteRune(char)
		count++
		totalChars++

		if count == interval {
			if totalChars < len(str) {
				buffer.WriteRune(' ')
				count = 0
			}
		}
	}

	return buffer.String()
}

// EqualsAny 给定字符串是否与提供的中任一字符串相同，相同则返回true，没有相同的返回false;
// 如果参与比对的字符串列表为空，返回false
func EqualsAny(str1 string, str2 []string) bool {
	if len(str2) == 0 {
		return false
	}
	for _, s := range str2 {
		if Equals(str1, s) {
			return true
		}
	}
	return false
}

// EqualsAnyIgnoreCase 给定字符串是否与提供的中任一字符串相同（忽略大小写），相同则返回true，没有相同的返回false;
// 如果参与比对的字符串列表为空，返回false
func EqualsAnyIgnoreCase(str1 string, str2 []string) bool {
	if len(str2) == 0 {
		return false
	}
	for _, s := range str2 {
		if EqualsIgnoreCase(str1, s) {
			return true
		}
	}
	return false
}

// EqualsAt 字符串指定位置的字符是否与给定字符相同
func EqualsAt(value string, position int, subStr string) bool {
	if validator.IsEmptyValue(reflect.ValueOf(value)) || position < 0 {
		return false
	}
	return len(value) > position && Equals(subStr, string(value[position]))
}

// Count 统计指定内容中包含指定字符串的数量
func Count(str string, searchStr string) int {
	hasEmpty, _ := validator.HasEmpty([]interface{}{str, searchStr})
	if hasEmpty || len(searchStr) > len(str) {
		return 0
	}
	return strings.Count(str, searchStr)
}

// CompareIgnoreCase 比较两个字符串，用于排序(大小写不敏感)
//
//	0 相等；<0 小于； >0 大于
func CompareIgnoreCase(str1 string, str2 string) int {
	return strings.Compare(strings.ToLower(str1), strings.ToLower(str2))
}
