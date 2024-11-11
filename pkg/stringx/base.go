/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 13:08:58
 * @FilePath: \go-toolbox\pkg\stringx\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"strings"
	"unicode"

	"github.com/kamalyes/go-toolbox/pkg/validator"
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

// ExtractValue 从给定的字符串中提取指定的键的值
func ExtractValue(extra string, key string, searchPrefix string) string {
	// 构造要查找的键名
	searchKey := key + "="
	if strings.Contains(extra, searchKey) {
		start := strings.Index(extra, searchKey) + len(searchKey)
		end := strings.Index(extra[start:], searchPrefix)
		if end == -1 {
			return extra[start:] // 如果没有分号，返回到字符串末尾
		}
		return extra[start : start+end]
	}
	return "" // 如果没有找到，返回空字符串
}

// CalculateMD5Hash 计算string md5 hash值
func CalculateMD5Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Coalesce 高性能字符串拼接
func Coalesce(s ...string) string {
	if len(s) == 0 {
		return "" // 如果没有输入，则返回空字符串
	}
	var str strings.Builder
	for _, v := range s {
		str.WriteString(v)
	}

	return str.String() // 返回拼接后的字符串
}

type CharacterStyle int

const (
	SnakeCharacterStyle  CharacterStyle = iota // 表示蛇形命名法（例如：hello_world）
	StudlyCharacterStyle                       // 表示每个单词首字母大写的风格（例如：HelloWorld）
	CamelCharacterStyle                        // 表示驼峰命名法（例如：helloWorld）
)

// ConvertCharacterStyle 根据指定的 CharacterStyle 将字符串转换为相应的格式
func ConvertCharacterStyle(input string, caseType CharacterStyle) string {
	trimmedStr := strings.TrimSpace(input)
	if trimmedStr == "" {
		return trimmedStr // 如果输入为空，直接返回
	}

	converters := map[CharacterStyle]func(string) string{
		SnakeCharacterStyle:  toSnakeCase,
		StudlyCharacterStyle: toStudlyCase,
		CamelCharacterStyle:  toCamelCase,
	}

	if converter, exists := converters[caseType]; exists {
		return converter(trimmedStr)
	}
	return trimmedStr // 默认返回原字符串
}

// toSnakeCase 将字符串转换为蛇形命名法
func toSnakeCase(s string) string {
	var builder strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && (s[i-1] != '_' && s[i-1] != ' ') {
				builder.WriteRune('_') // 在大写字母前添加下划线
			}
			builder.WriteRune(unicode.ToLower(r)) // 转为小写
		} else if r == '_' {
			if builder.Len() == 0 || builder.String()[builder.Len()-1] != '_' {
				builder.WriteRune('_') // 添加下划线，但避免重复
			}
		} else if r == ' ' {
			if builder.Len() > 0 && builder.String()[builder.Len()-1] != '_' {
				builder.WriteRune('_') // 将空格转换为下划线
			}
		} else {
			builder.WriteRune(r) // 直接添加小写字母或其他字符
		}
	}
	return builder.String()
}

// toStudlyCase 将字符串转换为每个单词首字母大写的风格
func toStudlyCase(s string) string {
	var builder strings.Builder
	s = strings.ReplaceAll(s, "_", " ") // 将下划线替换为空格
	words := strings.Fields(s)          // 按空格分割单词

	for _, word := range words {
		if len(word) > 0 {
			builder.WriteRune(unicode.ToUpper(rune(word[0]))) // 首字母大写
			builder.WriteString(word[1:])                     // 追加剩余部分
		}
	}
	return builder.String()
}

// toCamelCase 将字符串转换为驼峰命名法
func toCamelCase(s string) string {
	studly := toStudlyCase(s)                                    // 首先转换为 Studly Case
	return string(unicode.ToLower(rune(studly[0]))) + studly[1:] // 将首字母转换为小写
}
