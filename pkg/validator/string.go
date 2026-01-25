/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\string.go
 * @Description: 字符串比较和验证
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

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
