/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-16 18:55:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-17 09:51:08
 * @FilePath: \go-toolbox\pkg\sign\offset_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtonOffsetCipher(t *testing.T) {
	// 创建一个新的 ProtonOffsetCipher 实例
	psc := NewProtonOffsetCipher()

	testCases := map[string][]byte{
		"string":       []byte("Hello, World!"),
		"integers":     {1, 2, 3, 4, 5},
		"floats":       {0x3f, 0x80, 0x00, 0x00}, // 1.0 的 IEEE 754 表示
		"long":         []byte("This is a longer test string to check encryption and decryption."),
		"specialChars": []byte("!@#$%^&*()"),
		"booleanTrue":  {1}, // 代表布尔值 true
		"booleanFalse": {0}, // 代表布尔值 false
		"mixed":        []byte("Mix of different types: 12345 & !@#$%"),
	}

	for name, original := range testCases {
		t.Run(name, func(t *testing.T) {
			encrypted, err := psc.Encrypt(original)
			assert.NoError(t, err, "Encryption should not return an error")
			decrypted, err := psc.Decrypt(encrypted)
			assert.NoError(t, err, "Decryption should not return an error")
			assert.Equal(t, original, decrypted, "Decrypted data should equal original")
		})
	}
}

// 随机测试
func TestProtonOffsetCipherMultipleCases(t *testing.T) {
	// 创建一个新的 ProtonOffsetCipher 实例
	psc := NewProtonOffsetCipher()

	// 生成 10,000 个随机测试用例
	for i := 0; i < 10000; i++ {
		// 生成随机长度数据（1 到 1000 字节）
		dataLength := rand.Intn(1000) + 1
		data := make([]byte, dataLength)
		for j := range data {
			data[j] = byte(rand.Intn(1000)) // 随机字节
		}

		// 加密
		encrypted, err := psc.Encrypt(data)
		assert.NoError(t, err, "Encryption should not return an error")

		// 解密
		decrypted, err := psc.Decrypt(encrypted)
		assert.NoError(t, err, "Decryption should not return an error")

		// 验证解密后的数据与原始数据相同
		assert.Equal(t, data, decrypted, "Decrypted data should equal original")
	}
}
