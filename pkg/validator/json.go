/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\json.go
 * @Description: JSON 验证和 JSONPath 支持
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/kamalyes/go-jsonpath"
)

// ValidateJSON 验证JSON结构（基础版，仅验证格式）
func ValidateJSON(data []byte) error {
	var v interface{}
	return json.Unmarshal(data, &v)
}

// IsJSONNull 判断字节数据去除空白后是否为 JSON null
func IsJSONNull(data []byte) bool {
	return bytes.EqualFold(bytes.TrimSpace(data), []byte("null"))
}

// SkipJSONSpaces 跳过 JSON 字节流中的空白字符，并返回下一个非空白位置
func SkipJSONSpaces(data []byte, i int) int {
	for i < len(data) {
		switch data[i] {
		case ' ', '\n', '\r', '\t':
			i++
		default:
			return i
		}
	}
	return i
}

// ScanJSONString 扫描 JSON 字符串，并返回字符串结束后一位的位置
func ScanJSONString(data []byte, start int) (int, error) {
	if start >= len(data) || data[start] != '"' {
		return 0, fmt.Errorf("expected JSON string")
	}
	for i := start + 1; i < len(data); i++ {
		switch data[i] {
		case '\\':
			i++
		case '"':
			return i + 1, nil
		}
	}
	return 0, fmt.Errorf("unterminated JSON string")
}

// ScanJSONValueEnd 扫描任意 JSON 值，并返回值结束后一位的位置
func ScanJSONValueEnd(data []byte, start int) (int, error) {
	if start >= len(data) {
		return 0, fmt.Errorf("expected JSON value")
	}

	switch data[start] {
	case '"':
		return ScanJSONString(data, start)
	case '{', '[':
		return scanJSONCompositeEnd(data, start)
	default:
		return scanJSONScalarEnd(data, start)
	}
}

// scanJSONCompositeEnd 扫描 JSON 对象或数组，并校验括号匹配
func scanJSONCompositeEnd(data []byte, start int) (int, error) {
	stack := make([]byte, 0, 4)
	for i := start; i < len(data); i++ {
		next, done, updated, err := scanJSONCompositeToken(data, i, stack)
		if err != nil || done {
			return next, err
		}
		stack = updated
		i = next
	}
	return 0, fmt.Errorf("unterminated JSON composite value")
}

// scanJSONCompositeToken 处理对象或数组中的单个 token
func scanJSONCompositeToken(data []byte, i int, stack []byte) (next int, done bool, updated []byte, err error) {
	switch data[i] {
	case '"':
		return scanJSONCompositeString(data, i, stack)
	case '{':
		return i, false, append(stack, '}'), nil
	case '[':
		return i, false, append(stack, ']'), nil
	case '}', ']':
		return scanJSONCompositeClose(data[i], i, stack)
	default:
		return i, false, stack, nil
	}
}

// scanJSONCompositeString 跳过复合 JSON 值中的字符串内容
func scanJSONCompositeString(data []byte, i int, stack []byte) (next int, done bool, updated []byte, err error) {
	end, err := ScanJSONString(data, i)
	if err != nil {
		return 0, false, stack, err
	}
	return end - 1, false, stack, nil
}

// scanJSONCompositeClose 处理复合 JSON 值中的闭合括号
func scanJSONCompositeClose(token byte, i int, stack []byte) (next int, done bool, updated []byte, err error) {
	last := len(stack) - 1
	if last < 0 || stack[last] != token {
		return 0, false, stack, fmt.Errorf("mismatched JSON delimiter")
	}
	stack = stack[:last]
	if len(stack) == 0 {
		return i + 1, true, stack, nil
	}
	return i, false, stack, nil
}

// scanJSONScalarEnd 扫描 JSON 标量值，并返回值结束后一位的位置
func scanJSONScalarEnd(data []byte, start int) (int, error) {
	for i := start; i < len(data); i++ {
		switch data[i] {
		case ',', '}', ']':
			return i, nil
		case ' ', '\n', '\r', '\t':
			return skipJSONScalarSpaces(data, i), nil
		}
	}
	return len(data), nil
}

// skipJSONScalarSpaces 跳过标量值末尾空白，并确保空白后是合法分隔符
func skipJSONScalarSpaces(data []byte, i int) int {
	end := i
	for i < len(data) {
		switch data[i] {
		case ' ', '\n', '\r', '\t':
			i++
		case ',', '}', ']':
			return end
		default:
			return i
		}
	}
	return end
}

// ValidateJSONWithData 验证JSON并返回解析后的数据
func ValidateJSONWithData(body []byte) (any, error) {
	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("响应不是有效的JSON: %w", err)
	}
	return data, nil
}

// ValidateJSONField 验证JSON字段值
func ValidateJSONField(body []byte, field string, expected any) CompareResult {
	result := CompareResult{
		Expect: fmt.Sprintf("%v", expected),
	}

	// 解析JSON
	data, err := ValidateJSONWithData(body)
	if err != nil {
		result.Message = err.Error()
		return result
	}

	// 检查是否为对象
	dataMap, ok := data.(map[string]any)
	if !ok {
		result.Message = "JSON根节点不是对象"
		return result
	}

	// 检查字段是否存在
	actual, ok := dataMap[field]
	if !ok {
		result.Message = fmt.Sprintf("字段不存在: %s", field)
		return result
	}

	// 比较字段值
	result.Actual = fmt.Sprintf("%v", actual)
	result.Success = actual == expected
	if !result.Success {
		result.Message = fmt.Sprintf("字段值不匹配: %s, 期望: %v, 实际: %v", field, expected, actual)
	}

	return result
}

// ValidateJSONFields 批量验证JSON字段
func ValidateJSONFields(body []byte, rules map[string]any) []CompareResult {
	results := make([]CompareResult, 0, len(rules))
	for field, expected := range rules {
		result := ValidateJSONField(body, field, expected)
		results = append(results, result)
	}
	return results
}

// ValidateJSONPath 验证JSONPath查询结果
func ValidateJSONPath(body []byte, jsonPath string, expected any, op CompareOperator) CompareResult {
	result := CompareResult{
		Expect: fmt.Sprintf("%v", expected),
	}

	// 检查参数
	if jsonPath == "" {
		result.Message = "JSONPath 不能为空"
		return result
	}

	// 解析JSON
	data, err := ValidateJSONWithData(body)
	if err != nil {
		result.Message = err.Error()
		return result
	}

	// 使用 jsonpath 库查询
	value, err := jsonpath.JsonPathLookup(data, jsonPath)
	if err != nil {
		result.Message = fmt.Sprintf("JSON路径查询失败: %v", err)
		result.Actual = "查询失败"
		return result
	}

	// 转换为字符串进行比较
	actualStr := fmt.Sprintf("%v", value)
	expectStr := fmt.Sprintf("%v", expected)
	result.Actual = actualStr

	// 如果没有指定操作符，默认使用等于
	if op == "" {
		op = OpEqual
	}

	// 使用 CompareStrings 进行比较
	compareResult := CompareStrings(actualStr, expectStr, op)
	result.Success = compareResult.Success
	if !compareResult.Success {
		result.Message = compareResult.Message
	} else {
		result.Message = "JSONPath 验证通过"
	}

	return result
}

// ValidateJSONPathExists 验证JSONPath路径是否存在（不验证值）
func ValidateJSONPathExists(body []byte, jsonPath string) CompareResult {
	result := CompareResult{
		Expect: jsonPath,
	}

	// 解析JSON
	data, err := ValidateJSONWithData(body)
	if err != nil {
		result.Message = err.Error()
		return result
	}

	// 使用 jsonpath 库查询
	value, err := jsonpath.JsonPathLookup(data, jsonPath)
	if err != nil {
		result.Message = fmt.Sprintf("JSON路径不存在: %v", err)
		result.Actual = "路径不存在"
		return result
	}

	result.Success = true
	result.Actual = fmt.Sprintf("%v", value)
	result.Message = "JSONPath 路径存在"
	return result
}
