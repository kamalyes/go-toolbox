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
	"encoding/json"
	"fmt"

	"github.com/kamalyes/go-jsonpath"
)

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
