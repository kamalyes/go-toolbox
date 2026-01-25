/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\http.go
 * @Description: HTTP 相关验证（状态码、Header 等）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import "fmt"

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
