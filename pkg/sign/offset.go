/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-16 18:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-17 09:55:15
 * @FilePath: \go-toolbox\pkg\sign\offset.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"errors"
)

// ProtonOffsetCipher 结构体定义
type ProtonOffsetCipher struct {
	P       int // 质数
	C       int // 偏移量
	M       int // 模数
	inverse int // 模逆
}

// NewProtonOffsetCipher 创建一个新的 ProtonOffsetCipher 实例，使用默认偏移量
func NewProtonOffsetCipher() *ProtonOffsetCipher {
	return NewProtonOffsetCipherWithPCM(7, 3, 256)
}

// NewProtonOffsetCipherWithPCM 创建一个新的 ProtonOffsetCipher 实例，使用指定的偏移量
func NewProtonOffsetCipherWithPCM(p, c, m int) *ProtonOffsetCipher {
	inverse := modInverse(p, m) // 预计算模逆
	return &ProtonOffsetCipher{P: p, C: c, M: m, inverse: inverse}
}

// Encrypt 加密函数
func (cc *ProtonOffsetCipher) Encrypt(data []byte) ([]byte, error) {
	if cc.P <= 0 || cc.M <= 0 {
		return nil, errors.New("质数和模数必须大于零")
	}

	encrypted := make([]byte, len(data))

	for i, char := range data {
		// 使用质数和偏移量进行复杂运算
		// 加密公式: E(char) = (char * P + C) mod M
		// 其中:
		// char: 原始字符的整数值
		// P: 质数（用于扩展字符值）
		// C: 偏移量（用于增加随机性）
		// M: 模数（用于限制结果范围）
		encrypted[i] = byte((int(char)*cc.P + cc.C) % cc.M)
	}

	return encrypted, nil
}

// Decrypt 解密函数
func (cc *ProtonOffsetCipher) Decrypt(data []byte) ([]byte, error) {
	if cc.P <= 0 || cc.M <= 0 {
		return nil, errors.New("质数和模数必须大于零")
	}

	if len(data) == 0 {
		return nil, errors.New("加密数据不能为空")
	}

	decrypted := make([]byte, len(data))

	for i, char := range data {
		// 使用预计算的模逆进行解密
		// 解密公式: D(char) = ((char - C + M) * inverse) mod M
		// 其中:
		// char: 加密字符的整数值
		// C: 偏移量（与加密时相同）
		// M: 模数（用于限制结果范围）
		// inverse: P 在模 M 下的模逆（用于恢复原始字符值）
		decrypted[i] = byte(((int(char) - cc.C + cc.M) * cc.inverse) % cc.M)
	}

	return decrypted, nil
}

// modInverse 计算模逆
func modInverse(a, m int) int {
	m0, y, x := m, 0, 1
	if m == 1 {
		return 0
	}
	for a > 1 {
		q := a / m
		t := m
		m = a % m
		a = t
		t = y
		y = x - q*y
		x = t
	}
	if x < 0 {
		x += m0
	}
	return x
}
