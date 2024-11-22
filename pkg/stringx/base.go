/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 13:02:50
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

// StringX 是一个结构体，用于封装字符串值并提供操作方法。
type StringX struct {
	value string
}

// New 创建一个新的 StringX 实例。
func New(value string) *StringX {
	return &StringX{value: value}
}

// Get 返回当前字符串的值
func (s *StringX) Value() string {
	return s.value
}

// ToLower 将字符串转换为小写
func ToLower(str string) string {
	// 预分配结果切片，避免多次内存分配
	result := make([]rune, len(str))
	resultIndex := 0 // 结果切片的索引

	for _, r := range str {
		result[resultIndex] = unicode.ToLower(r)
		resultIndex++
	}

	return string(result[:resultIndex]) // 返回有效部分
}

// ToLowerChain 将字符串转换为小写（链式调用）
func (s *StringX) ToLowerChain() *StringX {
	s.value = ToLower(s.value)
	return s
}

// ToUpper 将字符串转换为大写
func ToUpper(str string) string {
	// 预分配结果切片，避免多次内存分配
	result := make([]rune, len(str))
	resultIndex := 0

	for _, r := range str {
		result[resultIndex] = unicode.ToUpper(r)
		resultIndex++
	}

	return string(result[:resultIndex])
}

// ToUpperChain 将字符串转换为大写（链式调用）
func (s *StringX) ToUpperChain() *StringX {
	s.value = ToUpper(s.value)
	return s
}

// ToTitle 将字符串转换为标题格式（每个单词首字母大写）
func ToTitle(str string) string {
	// 预分配结果切片，避免多次内存分配
	result := make([]rune, len(str))
	resultIndex := 0
	space := true // 用于标记是否在单词的开头

	for _, r := range str {
		if space {
			result[resultIndex] = unicode.ToUpper(r) // 首字母大写
			space = false
		} else {
			result[resultIndex] = unicode.ToLower(r) // 其他字符小写
		}
		resultIndex++

		if unicode.IsSpace(r) {
			space = true // 遇到空格，标记为下一个单词的开头
		}
	}

	return string(result[:resultIndex]) // 返回有效部分
}

// ToTitleChain 将字符串转换为标题格式（链式调用）
func (s *StringX) ToTitleChain() *StringX {
	s.value = ToTitle(s.value)
	return s
}

// Length 计算长度
func Length(str string) int {
	strRune := []rune(str)
	return len(strRune)
}

// LengthChain 计算长度（链式调用）
func (s *StringX) LengthChain() int {
	return Length(s.value)
}

// Reverse 反转给定的字符串
func Reverse(str string) string {
	runes := []rune(str)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	return string(runes)
}

// ReverseChain 反转字符串（链式调用）
func (s *StringX) ReverseChain() *StringX {
	s.value = Reverse(s.value)
	return s
}

// Equals 比较2个字符串（大小写敏感）
func Equals(str1 string, str2 string) bool {
	return str1 == str2
}

// EqualsChain 比较2个字符串（链式调用）
func (s *StringX) EqualsChain(str2 string) bool {
	return Equals(s.value, str2)
}

// EqualsIgnoreCase 比较2个字符串（大小写不敏感）
func EqualsIgnoreCase(str1 string, str2 string) bool {
	return strings.EqualFold(str1, str2)
}

// EqualsIgnoreCaseChain 比较2个字符串（链式调用，忽略大小写）
func (s *StringX) EqualsIgnoreCaseChain(str2 string) bool {
	return EqualsIgnoreCase(s.value, str2)
}

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

// InsertSpacesChain 插入空值（链式调用）
func (s *StringX) InsertSpacesChain(interval int) *StringX {
	s.value = InsertSpaces(s.value, interval)
	return s
}

// EqualsAny 给定字符串是否与提供的中任一字符串相同
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

// EqualsAnyChain 给定字符串是否与提供的中任一字符串相同（链式调用）
func (s *StringX) EqualsAnyChain(strs []string) bool {
	return EqualsAny(s.value, strs)
}

// EqualsAnyIgnoreCase 给定字符串是否与提供的中任一字符串相同（忽略大小写）
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

// EqualsAnyIgnoreCaseChain 给定字符串是否与提供的中任一字符串相同（链式调用，忽略大小写）
func (s *StringX) EqualsAnyIgnoreCaseChain(strs []string) bool {
	return EqualsAnyIgnoreCase(s.value, strs)
}

// EqualsAt 字符串指定位置的字符是否与给定字符相同
func EqualsAt(value string, position int, subStr string) bool {
	if validator.IsEmptyValue(reflect.ValueOf(value)) || position < 0 {
		return false
	}
	return len(value) > position && Equals(subStr, string(value[position]))
}

// EqualsAtChain 字符串指定位置的字符是否与给定字符相同（链式调用）
func (s *StringX) EqualsAtChain(position int, subStr string) bool {
	return EqualsAt(s.value, position, subStr)
}

// Count 统计指定内容中包含指定字符串的数量
func Count(str string, searchStr string) int {
	hasEmpty, _ := validator.HasEmpty([]interface{}{str, searchStr})
	if hasEmpty || len(searchStr) > len(str) {
		return 0
	}
	return strings.Count(str, searchStr)
}

// CountChain 统计指定内容中包含指定字符串的数量（链式调用）
func (s *StringX) CountChain(searchStr string) int {
	return Count(s.value, searchStr)
}

// CompareIgnoreCase 比较两个字符串，用于排序(大小写不敏感)
func CompareIgnoreCase(str1 string, str2 string) int {
	return strings.Compare(strings.ToLower(str1), strings.ToLower(str2))
}

// CompareIgnoreCaseChain 比较两个字符串（链式调用，用于排序，忽略大小写）
func (s *StringX) CompareIgnoreCaseChain(str2 string) int {
	return CompareIgnoreCase(s.value, str2)
}

// ExtractValue 从给定的字符串中提取指定的键的值
func ExtractValue(extra string, key string, searchPrefix string) string {
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

// CoalesceChain 高性能字符串拼接（链式调用）
func (s *StringX) CoalesceChain(strs ...string) *StringX {
	s.value = Coalesce(append([]string{s.value}, strs...)...)
	return s
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

// ConvertCharacterStyleChain 根据指定的 CharacterStyle 将字符串转换为相应的格式（链式调用）
func (s *StringX) ConvertCharacterStyleChain(caseType CharacterStyle) *StringX {
	s.value = ConvertCharacterStyle(s.value, caseType)
	return s
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
