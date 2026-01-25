/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\regex.go
 * @Description: 正则表达式验证
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"fmt"
	"regexp"
)

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
