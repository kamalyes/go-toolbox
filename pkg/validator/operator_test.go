/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\operator_test.go
 * @Description: 操作符和结果类型测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareResult(t *testing.T) {
	a := assert.New(t)

	result := CompareResult{
		Success: true,
		Message: "测试消息",
		Actual:  "实际值",
		Expect:  "期望值",
	}

	a.True(result.Success, "CompareResult.Success should be true")
	a.Equal("测试消息", result.Message)
	a.Equal("实际值", result.Actual)
	a.Equal("期望值", result.Expect)
}

func TestCompareOperators(t *testing.T) {
	a := assert.New(t)

	// 测试操作符常量定义
	a.Equal(CompareOperator("eq"), OpEqual)
	a.Equal(CompareOperator("ne"), OpNotEqual)
	a.Equal(CompareOperator("gt"), OpGreaterThan)
	a.Equal(CompareOperator("gte"), OpGreaterThanOrEqual)
	a.Equal(CompareOperator("lt"), OpLessThan)
	a.Equal(CompareOperator("lte"), OpLessThanOrEqual)
	a.Equal(CompareOperator("contains"), OpContains)
	a.Equal(CompareOperator("not_contains"), OpNotContains)
	a.Equal(CompareOperator("has_prefix"), OpHasPrefix)
	a.Equal(CompareOperator("has_suffix"), OpHasSuffix)
	a.Equal(CompareOperator("regex"), OpRegex)
	a.Equal(CompareOperator("empty"), OpEmpty)
	a.Equal(CompareOperator("not_empty"), OpNotEmpty)

	// 测试符号别名
	a.Equal(CompareOperator("="), OpSymbolEqual)
	a.Equal(CompareOperator("!="), OpSymbolNotEqual)
	a.Equal(CompareOperator(">"), OpSymbolGreaterThan)
	a.Equal(CompareOperator(">="), OpSymbolGreaterThanOrEqual)
	a.Equal(CompareOperator("<"), OpSymbolLessThan)
	a.Equal(CompareOperator("<="), OpSymbolLessThanOrEqual)
}
