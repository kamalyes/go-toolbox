/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 22:10:16
 * @FilePath: \go-toolbox\pkg\validator\operator.go
 * @Description: 比较操作符定义和结果类型
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

// CompareOperator 比较操作符
type CompareOperator string

// String 返回操作符的字符串表示
func (op CompareOperator) String() string {
	return string(op)
}

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
