/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-16 18:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-17 09:55:15
 * @FilePath: \go-toolbox\pkg\sign\xor.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"errors"
)

// XORCipher 结构体定义
type XORCipher struct {
	Key byte // 加密密钥
}

// NewXORCipher 创建一个新的 XORCipher 实例
func NewXORCipher(key byte) *XORCipher {
	return &XORCipher{Key: key}
}

// Encrypt 加密函数
func (xc *XORCipher) Encrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("数据不能为空")
	}

	encrypted := make([]byte, len(data))

	for i, char := range data {
		// 使用 XOR 运算进行加密
		encrypted[i] = char ^ xc.Key
	}

	return encrypted, nil
}

// Decrypt 解密函数
func (xc *XORCipher) Decrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("加密数据不能为空")
	}

	decrypted := make([]byte, len(data))

	for i, char := range data {
		// 使用 XOR 运算进行解密
		decrypted[i] = char ^ xc.Key
	}

	return decrypted, nil
}
