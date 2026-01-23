/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 10:50:01
 * @FilePath: \go-toolbox\pkg\validator\compare.go
 * @Description: 比较和验证相关功能
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"cmp"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/kamalyes/go-toolbox/pkg/types"
)

// CompareOperator 比较操作符
type CompareOperator string

const (
	OpEqual              CompareOperator = "eq"
	OpNotEqual           CompareOperator = "ne"
	OpGreaterThan        CompareOperator = "gt"
	OpGreaterThanOrEqual CompareOperator = "gte"
	OpLessThan           CompareOperator = "lt"
	OpLessThanOrEqual    CompareOperator = "lte"
	OpContains           CompareOperator = "contains"
	OpNotContains        CompareOperator = "not_contains"
	OpHasPrefix          CompareOperator = "has_prefix"
	OpHasSuffix          CompareOperator = "has_suffix"
	OpRegex              CompareOperator = "regex"
	OpEmpty              CompareOperator = "empty"
	OpNotEmpty           CompareOperator = "not_empty"

	// 符号别名
	OpSymbolEqual              CompareOperator = "="
	OpSymbolNotEqual           CompareOperator = "!="
	OpSymbolGreaterThan        CompareOperator = ">"
	OpSymbolGreaterThanOrEqual CompareOperator = ">="
	OpSymbolLessThan           CompareOperator = "<"
	OpSymbolLessThanOrEqual    CompareOperator = "<="
)

// CompareResult 比较结果
type CompareResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Actual  string `json:"actual"`
	Expect  string `json:"expect"`
}

// CompareStrings 比较两个字符串
func CompareStrings(actual, expect string, op CompareOperator) CompareResult {
	result := CompareResult{
		Actual: actual,
		Expect: expect,
	}

	switch op {
	case OpEqual, OpSymbolEqual:
		result.Success = actual == expect
	case OpNotEqual, OpSymbolNotEqual:
		result.Success = actual != expect
	case OpContains:
		result.Success = strings.Contains(actual, expect)
	case OpNotContains:
		result.Success = !strings.Contains(actual, expect)
	case OpHasPrefix:
		result.Success = strings.HasPrefix(actual, expect)
	case OpHasSuffix:
		result.Success = strings.HasSuffix(actual, expect)
	case OpEmpty:
		result.Success = actual == ""
	case OpNotEmpty:
		result.Success = actual != ""
	case OpRegex:
		matched, err := regexp.MatchString(expect, actual)
		if err != nil {
			result.Message = fmt.Sprintf("正则表达式错误: %v", err)
			return result
		}
		result.Success = matched
	default:
		result.Message = "不支持的操作符"
	}

	if !result.Success && result.Message == "" {
		result.Message = fmt.Sprintf("比较失败: 期望 %s %s, 实际 %s",
			expect, op, actual)
	}

	return result
}

// CompareNumbers 比较两个数值 - 使用 cmp.Compare 支持泛型比较
func CompareNumbers[T types.Numerical](actual, expect T, op CompareOperator) CompareResult {
	result := CompareResult{
		Actual: fmt.Sprintf("%v", actual),
		Expect: fmt.Sprintf("%v", expect),
	}

	// 使用 cmp.Compare 进行比较（返回 -1, 0, 1）
	cmpResult := cmp.Compare(actual, expect)

	switch op {
	case OpEqual, OpSymbolEqual:
		result.Success = cmpResult == 0
	case OpNotEqual, OpSymbolNotEqual:
		result.Success = cmpResult != 0
	case OpGreaterThan, OpSymbolGreaterThan:
		result.Success = cmpResult > 0
	case OpGreaterThanOrEqual, OpSymbolGreaterThanOrEqual:
		result.Success = cmpResult >= 0
	case OpLessThan, OpSymbolLessThan:
		result.Success = cmpResult < 0
	case OpLessThanOrEqual, OpSymbolLessThanOrEqual:
		result.Success = cmpResult <= 0
	default:
		result.Message = "不支持的数值操作符"
	}

	if !result.Success && result.Message == "" {
		result.Message = fmt.Sprintf("数值比较失败: 期望 %v %s %v, 实际 %v",
			expect, op, expect, actual)
	}

	return result
}

// ValidateJSON 验证JSON结构（基础版，仅验证格式）
func ValidateJSON(data []byte) error {
	var v interface{}
	return json.Unmarshal(data, &v)
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

// ValidateStatusCode 验证HTTP状态码 - 支持多种比较操作符
func ValidateStatusCode(statusCode, expected int, op CompareOperator) CompareResult {
	if op == "" {
		op = OpEqual // 默认相等比较
	}
	return CompareNumbers(statusCode, expected, op)
}

// ValidateStatusCodeRange 验证HTTP状态码在范围内
func ValidateStatusCodeRange(actual, min, max int) CompareResult {
	result := CompareResult{
		Actual: fmt.Sprintf("%d", actual),
		Expect: fmt.Sprintf("%d-%d", min, max),
	}

	result.Success = actual >= min && actual <= max
	if !result.Success {
		result.Message = fmt.Sprintf("状态码 %d 不在范围 [%d, %d] 内",
			actual, min, max)
	}

	return result
}

// ValidateContains 验证响应体包含指定字符串
func ValidateContains(body []byte, substring string) CompareResult {
	result := CompareResult{
		Actual: string(body),
		Expect: substring,
	}

	if substring == "" {
		result.Success = true
		return result
	}

	result.Success = strings.Contains(string(body), substring)
	if !result.Success {
		result.Message = fmt.Sprintf("响应不包含: %s", substring)
	}

	return result
}

// ValidateNotContains 验证响应体不包含指定字符串
func ValidateNotContains(body []byte, substring string) CompareResult {
	result := CompareResult{
		Actual: string(body),
		Expect: fmt.Sprintf("不包含: %s", substring),
	}

	if substring == "" {
		result.Success = true
		return result
	}

	result.Success = !strings.Contains(string(body), substring)
	if !result.Success {
		result.Message = fmt.Sprintf("响应包含不应存在的内容: %s", substring)
	}

	return result
}

// ValidateRegex 验证响应体匹配正则表达式
func ValidateRegex(body []byte, pattern string) CompareResult {
	result := CompareResult{
		Actual: string(body),
		Expect: fmt.Sprintf("匹配正则: %s", pattern),
	}

	if pattern == "" {
		result.Success = true
		return result
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		result.Message = fmt.Sprintf("正则表达式编译失败: %v", err)
		return result
	}

	result.Success = regex.Match(body)
	if !result.Success {
		result.Message = fmt.Sprintf("响应不匹配正则: %s", pattern)
	}

	return result
}

// ValidateHeader 验证HTTP头部字段
func ValidateHeader(headers map[string]string, key, expected string, op CompareOperator) CompareResult {
	actual, ok := headers[key]
	if !ok {
		return CompareResult{
			Success: false,
			Message: fmt.Sprintf("Header 不存在: %s", key),
			Expect:  expected,
		}
	}

	if op == "" {
		op = OpEqual
	}

	return CompareStrings(actual, expected, op)
}

// ValidateContentType 验证 Content-Type
func ValidateContentType(headers map[string]string, expected string) CompareResult {
	return ValidateHeader(headers, "Content-Type", expected, OpContains)
}
