/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 11:50:50
 * @FilePath: \go-toolbox\pkg\json\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package json

import "strings"

// KeyValuePairs 是一个用于存储键值对的结构体
type KeyValuePairs struct {
	pairs map[string]interface{}
}

// NewKeyValuePairs 创建一个新的 KeyValuePairs 实例
func NewKeyValuePairs() *KeyValuePairs {
	return &KeyValuePairs{pairs: make(map[string]interface{})}
}

// Add 向 KeyValuePairs 中添加一个键值对
func (kv *KeyValuePairs) Add(key string, value interface{}) *KeyValuePairs {
	kv.pairs[key] = value
	return kv
}

// AppendKeysToJSONMarshal 将键值对追加到 JSON 中
func AppendKeysToJSONMarshal(originalJSON string, pairs *KeyValuePairs) ([]byte, error) {
	var originalMap map[string]interface{}
	if err := Unmarshal([]byte(originalJSON), &originalMap); err != nil {
		return []byte{}, err
	}

	for key, value := range pairs.pairs {
		originalMap[key] = value
	}

	return Marshal(originalMap)
}

// ReplaceKeys 替换 JSON 中的键（key）中的指定字符串为目标字符串
func ReplaceKeys(data interface{}, oldStr, newStr string) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		replacedMap := make(map[string]interface{})
		for k, value := range v {
			// 替换键中的字符串
			newKey := strings.ReplaceAll(k, oldStr, newStr)
			replacedValue, err := ReplaceKeys(value, oldStr, newStr) // 递归处理值
			if err != nil {
				return nil, err
			}
			replacedMap[newKey] = replacedValue
		}
		return replacedMap, nil
	case []interface{}:
		for i, value := range v {
			replacedValue, err := ReplaceKeys(value, oldStr, newStr) // 递归处理值
			if err != nil {
				return nil, err
			}
			v[i] = replacedValue
		}
		return v, nil
	default:
		return data, nil
	}
}
