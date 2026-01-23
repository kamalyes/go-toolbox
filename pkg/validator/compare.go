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

// ValidateJSON 验证JSON结构
func ValidateJSON(data []byte) error {
	var v interface{}
	return json.Unmarshal(data, &v)
}

// ValidateStatusCode 验证HTTP状态码
func ValidateStatusCode(actual, expect int) CompareResult {
	return CompareNumbers(actual, expect, OpEqual)
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
