/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-24 11:25:16
 * @FilePath: \go-toolbox\json\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package json

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

// AppendKeysToJSON 将键值对追加到 JSON 字符串中
func AppendKeysToJSON(originalJSON string, pairs *KeyValuePairs) (string, error) {
	var originalMap map[string]interface{}
	if err := Unmarshal([]byte(originalJSON), &originalMap); err != nil {
		return "", err
	}

	for key, value := range pairs.pairs {
		originalMap[key] = value
	}

	updatedJSON, err := Marshal(originalMap)
	if err != nil {
		return "", err
	}

	return string(updatedJSON), nil
}
