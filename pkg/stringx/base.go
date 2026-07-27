/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-27 23:55:26
 * @FilePath: \go-toolbox\pkg\stringx\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// StringX 是一个结构体，用于封装字符串值并提供操作方法。
type StringX struct {
	value string
}

// New 创建一个新的 StringX 实例。
func New(value string) *StringX {
	return &StringX{value: value}
}

// Value 返回当前字符串的值
func (s *StringX) Value() string {
	return s.value
}

// ToLower 将字符串转换为小写
// 直接委托标准库 strings.ToLower：标准库已做极致优化（ASCII 快速路径 + SIMD），
func ToLower(str string) string {
	return strings.ToLower(str)
}

// ToLowerChain 将字符串转换为小写（链式调用）
func (s *StringX) ToLowerChain() *StringX {
	s.value = ToLower(s.value)
	return s
}

// ToUpper 将字符串转换为大写
// 直接委托标准库 strings.ToUpper：标准库已做极致优化（ASCII 快速路径 + SIMD）
func ToUpper(str string) string {
	return strings.ToUpper(str)
}

// ToUpperChain 将字符串转换为大写（链式调用）
func (s *StringX) ToUpperChain() *StringX {
	s.value = ToUpper(s.value)
	return s
}

// ToTitle 将字符串转换为标题格式（每个单词首字母大写）
func ToTitle(str string) string {
	// 预分配结果切片，避免多次内存分配
	result := make([]rune, 0, len(str))
	space := true // 用于标记是否在单词的开头

	for _, r := range str {
		if space {
			result = append(result, unicode.ToUpper(r)) // 首字母大写
			space = false
		} else {
			result = append(result, unicode.ToLower(r)) // 其他字符小写
		}

		if unicode.IsSpace(r) {
			space = true // 遇到空格，标记为下一个单词的开头
		}
	}

	return string(result)
}

// ToTitleChain 将字符串转换为标题格式（链式调用）
func (s *StringX) ToTitleChain() *StringX {
	s.value = ToTitle(s.value)
	return s
}

// Length 计算字符（rune）长度
// 使用 utf8.RuneCountInString 零分配遍历，替代 []rune(str) 的堆分配
func Length(str string) int {
	return utf8.RuneCountInString(str)
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
	if value == "" || position < 0 {
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
	if str == "" || searchStr == "" || len(searchStr) > len(str) {
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
		return ""
	}
	// 预计算总长度，一次分配
	total := 0
	for _, v := range s {
		total += len(v)
	}
	var b strings.Builder
	b.Grow(total)
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
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
	KebabCharacterStyle                        // 表示短横线命名法（例如：hello-world）
	PascalCharacterStyle                       // 表示帕斯卡命名法（例如：HelloWorld，同 StudlyCharacterStyle）
)

// ConvertCharacterStyle 根据指定的 CharacterStyle 将字符串转换为相应的格式
func ConvertCharacterStyle(input string, caseType CharacterStyle) string {
	trimmedStr := strings.TrimSpace(input)
	if trimmedStr == "" {
		return trimmedStr // 如果输入为空，直接返回
	}

	converters := map[CharacterStyle]func(string) string{
		SnakeCharacterStyle:  ToSnakeCase,
		StudlyCharacterStyle: ToPascalCase, // 使用统一的 ToPascalCase
		CamelCharacterStyle:  ToCamelCase,
		KebabCharacterStyle:  ToKebabCase,
		PascalCharacterStyle: ToPascalCase,
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

// ToPascalCase 将字符串转换为帕斯卡命名法（首字母大写的驼峰，例如：UserName）
// 支持输入格式：camelCase, snake_case, kebab-case, 普通字符串（带空格）
// 增强版：合并了原 toStudlyCase 的逻辑，统一处理所有分隔符
func ToPascalCase(s string) string {
	if s == "" {
		return s
	}

	// 如果包含下划线、连字符或空格，先分割处理
	if strings.ContainsAny(s, "_- ") {
		// 统一替换为下划线再分割
		normalized := strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), " ", "_")
		parts := strings.Split(normalized, "_")
		var builder strings.Builder
		for _, part := range parts {
			if len(part) > 0 {
				// 首字母大写，其余小写
				builder.WriteRune(unicode.ToUpper(rune(part[0])))
				if len(part) > 1 {
					builder.WriteString(strings.ToLower(part[1:]))
				}
			}
		}
		return builder.String()
	}

	// 如果已经是驼峰式，只需首字母大写
	if len(s) > 0 && unicode.IsLower(rune(s[0])) {
		return string(unicode.ToUpper(rune(s[0]))) + s[1:]
	}

	return s
}

// ToSnakeCase 将字符串转换为蛇形命名法（例如：user_name）
// 支持输入格式：camelCase, PascalCase, kebab-case
func ToSnakeCase(s string) string {
	// 先去除首尾空格
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// 预分配：最坏情况每个字符前加一个下划线
	result := make([]byte, 0, len(s)*2)
	lastIsUnderscore := false

	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		// 将连字符和空格都转换为下划线
		if r == '-' || r == ' ' {
			// 避免连续的下划线
			if len(result) > 0 && !lastIsUnderscore {
				result = append(result, '_')
				lastIsUnderscore = true
			}
			continue
		}

		if unicode.IsUpper(r) {
			// 如果不是第一个字符，且前一个字符不是大写，添加下划线
			if i > 0 {
				prevRune := rune(s[i-1])
				if !unicode.IsUpper(prevRune) && prevRune != '_' && prevRune != ' ' && prevRune != '-' {
					if len(result) > 0 && !lastIsUnderscore {
						result = append(result, '_')
						lastIsUnderscore = true
					}
				}
			}
			result = append(result, byte(unicode.ToLower(r)))
			lastIsUnderscore = false
		} else {
			result = append(result, s[i])
			lastIsUnderscore = false
		}
	}

	return string(result)
}

// ToCamelCase 将字符串转换为驼峰命名法（首字母小写，例如：userName）
// 支持输入格式：PascalCase, snake_case, kebab-case
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if len(pascal) > 0 {
		return string(unicode.ToLower(rune(pascal[0]))) + pascal[1:]
	}
	return pascal
}

// ToKebabCase 将字符串转换为短横线命名法（例如：user-name）
// 支持输入格式：camelCase, PascalCase, snake_case
func ToKebabCase(s string) string {
	return strings.ReplaceAll(ToSnakeCase(s), "_", "-")
}

// NormalizeFieldName 规范化字段名，返回所有可能的命名风格变体
// 支持的输入格式：camelCase, PascalCase, snake_case, kebab-case
// 返回顺序：原始名称、PascalCase、camelCase、snake_case、kebab-case
func NormalizeFieldName(fieldName string) []string {
	if fieldName == "" {
		return []string{}
	}

	variants := make([]string, 0, 5)

	// 原始名称
	variants = append(variants, fieldName)

	// 转换为 PascalCase (首字母大写的驼峰)
	pascalCase := ToPascalCase(fieldName)
	if pascalCase != fieldName {
		variants = append(variants, pascalCase)
	}

	// 转换为 camelCase (首字母小写的驼峰)
	camelCase := ToCamelCase(fieldName)
	if camelCase != fieldName && camelCase != pascalCase {
		variants = append(variants, camelCase)
	}

	// 转换为 snake_case
	snakeCase := ToSnakeCase(fieldName)
	if snakeCase != fieldName {
		variants = append(variants, snakeCase)
	}

	// 转换为 kebab-case
	kebabCase := ToKebabCase(fieldName)
	if kebabCase != fieldName && kebabCase != snakeCase {
		variants = append(variants, kebabCase)
	}

	return variants
}

func ToInt(s string) (int, error) {
	return strconv.Atoi(Trim(s))
}

// FindKeysByValue 函数
func FindKeysByValue(data map[string]string, searchValue string) []string {
	var result []string
	for key, value := range data {
		// 将值以逗号分隔成切片
		values := strings.Split(value, ",")
		// 检查切片中是否有完全匹配的值
		for _, v := range values {
			if strings.TrimSpace(v) == searchValue {
				result = append(result, key)
				break
			}
		}
	}
	return result
}

// ToSliceByte 将字符串转换为字节切片
func ToSliceByte(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// TruncateMessage 截断消息内容用于日志显示
func TruncateMessage(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}

// NormalizeSQLDirection 规范化排序方向，仅允许 "ASC" 或 "DESC"（不区分大小写）
// direction: 输入的排序方向
// defaultDirection: 默认方向（当输入无效时）
// 返回规范化后的排序方向（大写）
func NormalizeSQLDirection(direction, defaultDirection string) string {
	// 转换为大写进行比较
	upperDirection := ToUpper(direction)

	// 检查是否为有效的排序方向
	if upperDirection == "ASC" || upperDirection == "DESC" {
		return upperDirection
	}

	// 如果无效，返回默认方向
	return ToUpper(defaultDirection)
}
