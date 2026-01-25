/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\number.go
 * @Description: 数值比较验证
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"cmp"
	"fmt"

	"github.com/kamalyes/go-toolbox/pkg/types"
)

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
