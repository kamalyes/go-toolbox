/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-09 01:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 21:28:15
 * @FilePath: \go-toolbox\pkg\convert\errors.go
 * @Description: 转换错误包装与泛型 JSON 转换
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"encoding/json"
	"fmt"
)

// ConvertError 用于包装转换过程中产生的错误
type ConvertError struct {
	Op  string // 操作名称
	Err error  // 原始错误
}

// Error 实现 error 接口，返回错误信息
func (e *ConvertError) Error() string {
	return fmt.Sprintf("convert error during %s: %v", e.Op, e.Err)
}

// UnmarshalJSON 泛型 JSON 反序列化（复用现有 json 包功能）
func UnmarshalJSON[T any](jsonData []byte) (*T, error) {
	var result T
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, &ConvertError{"unmarshaling JSON", err}
	}
	return &result, nil
}

// MarshalJSON 泛型 JSON 序列化（复用现有 json 包功能）
func MarshalJSON[T any](data T) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, &ConvertError{"marshalling to JSON", err}
	}
	return jsonData, nil
}
