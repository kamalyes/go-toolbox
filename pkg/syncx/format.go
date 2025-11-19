/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-12-13 13:05:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-13 17:50:20
 * @FilePath: \go-toolbox\pkg\syncx\format.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import "encoding/json"

// BuildContentExtra 构建额外内容JSON
func BuildContentExtra(data map[string]interface{}) string {
	if len(data) == 0 {
		return "{}"
	}

	if dataBytes, err := json.Marshal(data); err == nil {
		return string(dataBytes)
	}

	return "{}"
}

// GetStringFromData 从 Data map 中获取字符串值
func GetStringFromData(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}
	if val, ok := data[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

// GetBoolFromData 从 Data map 中获取布尔值
func GetBoolFromData(data map[string]interface{}, key string) bool {
	if data == nil {
		return false
	}
	if val, ok := data[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}

func GetInt64FromData(data map[string]interface{}, key string) int64 {
	if data == nil {
		return 0
	}
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		}
	}
	return 0
}

// ParseContentExtraToMap 将message.Data解析为map[string]string格式
// 用于protobuf的ContentExtra字段
func ParseContentExtraToMap(data map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for key, value := range data {
		if strVal, ok := value.(string); ok {
			result[key] = strVal
		} else {
			// 将非字符串值转换为JSON字符串
			if jsonBytes, err := json.Marshal(value); err == nil {
				result[key] = string(jsonBytes)
			}
		}
	}

	return result
}
