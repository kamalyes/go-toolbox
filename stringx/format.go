/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 18:33:45
 * @FilePath: \go-toolbox\stringx\format.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseStr 解析字符串
func ParseStr(source interface{}) string {
	switch v := source.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.Itoa(int(v))
	case int16:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		return time.Time.Format(v, "2006-01-02 15:04:05")
	case bool:
		return strconv.FormatBool(v)
	default:
		{
			b, _ := json.Marshal(v)
			return string(b)
		}
	}
}

// FillBefore 将已有字符串填充为规定长度，如果已有字符串超过这个长度则返回这个字符串
func FillBefore(str string, char string, length int) string {
	return Fill(str, char, length, true)
}

// FillAfter 将已有字符串填充为规定长度，如果已有字符串超过这个长度则返回这个字符串
func FillAfter(str string, char string, length int) string {
	return Fill(str, char, length, false)
}

// Fill 将已有字符串填充为规定长度，如果已有字符串超过这个长度则返回这个字符串
func Fill(str string, char string, length int, isPre bool) string {
	if len(str) >= length {
		return str
	}
	fillLength := length - len(str)
	fillStr := strings.Repeat(char, fillLength)
	if isPre {
		return fillStr + str
	}
	return str + fillStr
}

// Format 通过map中的参数 格式化字符串
// map = {a: "aValue", b: "bValue"} format("{a} and {b}", map) ---=》 aValue and bValue
func Format(template string, params map[string]interface{}) string {
	// 遍历map中的键值对
	for key, value := range params {
		// 构造占位符，例如 "{a}"
		placeholder := fmt.Sprintf("{%s}", key)
		// 将占位符替换为对应的值
		template = strings.ReplaceAll(template, placeholder, fmt.Sprintf("%v", value))
	}
	return template
}

// IndexedFormat 有序的格式化文本，使用{number}做为占位符
// 通常使用：format("this is {0} for {1}", "a", "b") =》 this is a for b
func IndexedFormat(template string, params []interface{}) string {

	// 遍历所有参数
	for i, param := range params {
		placeholder := fmt.Sprintf("{%d}", i)
		// 将占位符替换为对应的值
		template = ReplaceAll(template, placeholder, fmt.Sprintf("%v", param))
	}
	return template
}

// TruncateAppendEllipsis 截断字符串，使用不超过maxBytes长度。截断后自动追加省略号(...) 用于存储数据库varchar且编码为UTF-8的字段
func TruncateAppendEllipsis(str string, maxBytes int) string {
	// 如果字符串本身就比 maxBytes 短，则不需要截断
	if len(str) <= maxBytes {
		return str
	}
	s := str[0 : maxBytes-3]

	return string(append([]byte(s), '.', '.', '.'))

}

// Truncate 截断字符串，使用不超过maxBytes长度
func Truncate(str string, maxBytes int) string {
	return str[0:maxBytes]
}

// AddPrefixIfNot 如果给定字符串不是以prefix开头的，在开头补充 prefix
func AddPrefixIfNot(str string, prefix string) string {
	if StartWith(str, prefix) {
		return str
	}
	return prefix + str
}

// AddSuffixIfNot 如果给定字符串不是以suffix结尾的，在尾部补充 suffix
func AddSuffixIfNot(str string, suffix string) string {
	if EndWith(str, suffix) {
		return str
	}
	return str + suffix
}
