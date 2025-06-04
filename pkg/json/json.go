//go:build !jsoniter && !go_json && !(sonic && avx && (linux || windows || darwin) && amd64)

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-02 18:15:36
 * @FilePath: \go-middleware\pkg\json\json.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package json

import "encoding/json"

var (
	// Marshal is exported by go-toolbox/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by go-toolbox/json package.
	Unmarshal = json.Unmarshal
	// MarshalIndent is exported by go-toolbox/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by go-toolbox/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by go-toolbox/json package.
	NewEncoder = json.NewEncoder
)

func MarshalWithExtraField(v any, extraKey string, extraValue any) ([]byte, error) {
	// 先序列化成 JSON 字节
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// 反序列化成 map
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	// 添加额外字段
	m[extraKey] = extraValue

	// 重新序列化
	return json.Marshal(m)
}


// Compact 将输入的 JSON 字节数组解析并重新编码为紧凑格式的 JSON 字符串。
// 如果输入不是合法 JSON，或者编码失败，则返回原始字符串。
// 该函数适用于需要将 JSON 数据以最小化格式存储或传输的场景。
//
// 参数:
//   data - 输入的 JSON 数据，字节切片格式
//
// 返回值:
//   string - 紧凑格式的 JSON 字符串，或原始字符串（当输入不是合法 JSON 时）
func Compact(data []byte) string {
    var obj interface{}
    // 尝试将输入 JSON 反序列化为通用接口
    if err := json.Unmarshal(data, &obj); err != nil {
        // 反序列化失败，说明不是合法 JSON，直接返回原始字符串
        return string(data)
    }
    // 重新编码为紧凑的 JSON 字符串（无多余空格和换行）
    b, err := json.Marshal(obj)
    if err != nil {
        // 编码失败，返回原始字符串
        return string(data)
    }
    return string(b)
}