/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-17 10:02:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-17 10:35:16
 * @FilePath: \go-toolbox\tests\sign_xor_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"math/rand"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/sign"
	"github.com/stretchr/testify/assert"
)

// 测试 XORCipher 的加密和解密
func TestXORCipher(t *testing.T) {
	key := byte(0xAA)
	cipher := sign.NewXORCipher(key)

	// 测试数据
	data := []byte("Hello, World!")

	// 加密
	encrypted, err := cipher.Encrypt(data)
	assert.NoError(t, err, "Encryption should not return an error")
	assert.NotEqual(t, data, encrypted, "Encrypted data should not be the same as original data")

	// 解密
	decrypted, err := cipher.Decrypt(encrypted)
	assert.NoError(t, err, "Decryption should not return an error")
	assert.Equal(t, data, decrypted, "Decrypted data should match original data")
}

// 随机测试
func TestXORCipherMultipleCases(t *testing.T) {
	key := byte(0xAA)
	cipher := sign.NewXORCipher(key)

	for i := 0; i < 10000; i++ {
		// 生成随机数据
		data := make([]byte, rand.Intn(1000)+1) // 随机长度 1 到 1000 字节
		for j := range data {
			data[j] = byte(rand.Intn(1000)) // 随机字节
		}

		// 加密
		encrypted, err := cipher.Encrypt(data)
		assert.NoError(t, err, "Encryption should not return an error")
		assert.NotEqual(t, data, encrypted, "Encrypted data should not be the same as original data")

		// 解密
		decrypted, err := cipher.Decrypt(encrypted)
		assert.NoError(t, err, "Decryption should not return an error")
		assert.Equal(t, data, decrypted, "Decrypted data should match original data")
	}
}
