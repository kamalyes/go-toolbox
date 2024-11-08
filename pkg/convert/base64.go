/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 01:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 10:11:15
 * @FilePath: \go-toolbox\pkg\convert\base64.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"bytes"
	"encoding/base64"
	"errors"
	"sync"
)

// 创建一个对象池，用于复用 bytes.Buffer
var baseBufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// encodeBase64 编码 Base64 的通用函数
func encodeBase64(data interface{}, encoding *base64.Encoding) (string, error) {
	// 从对象池中获取一个 bytes.Buffer
	b := baseBufferPool.Get().(*bytes.Buffer)
	b.Reset() // 重置 Buffer 的内容

	encoder := base64.NewEncoder(encoding, b)

	var input []byte
	switch v := data.(type) {
	case []byte:
		input = v
	case string:
		input = []byte(v)
	default:
		return "", errors.New("unsupported type")
	}

	if _, err := encoder.Write(input); err != nil {
		return "", err
	}
	if err := encoder.Close(); err != nil {
		return "", err
	}

	// 将 Buffer 的内容转换为字符串
	result := b.String()

	// 将 Buffer 归还到对象池中
	baseBufferPool.Put(b)

	return result, nil
}

// B64Encode Base64 编码
func B64Encode(data interface{}) (string, error) {
	return encodeBase64(data, base64.StdEncoding)
}

// B64Decode Base64 解码
func B64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// B64UrlEncode Base64 URL 安全编码
func B64UrlEncode(data interface{}) (string, error) {
	return encodeBase64(data, base64.URLEncoding)
}

// B64UrlDecode Base64 URL 安全解码
func B64UrlDecode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

// B64ToByte 将 Base64 字符串解码为字节切片
// 参数：imageBase64 - 要解码的 Base64 字符串
// 返回：解码后的字节切片和可能的错误
func B64ToByte(imageBase64 string) ([]byte, error) {
	// 解码 Base64 字符串
	image, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return nil, err // 返回错误
	}

	return image, nil // 返回解码后的字节
}
