/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:15:15
 * @FilePath: \go-toolbox\pkg\sign\rsa_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义常量以避免重复的字符串
const (
	errSaltGeneration = "盐生成失败: %v"
	errEncryption     = "加密失败: %v"
	errDecryption     = "解密失败: %v"
)

// 生成 RSA 密钥对的辅助函数
func generateRsaKeyPair(b *testing.B, keySize RsaKeySize) *RsaKeyPair {
	keyPair, err := GenerateRsaKeyPair(keySize)
	if err != nil {
		b.Fatal(err)
	}
	return keyPair
}

func BenchmarkRsaEncryption(b *testing.B) {
	// 1. 生成 RSA 密钥对
	keyPair := generateRsaKeyPair(b, RsaKeySize2048)

	// 2. 创建 RSA 加解密器
	rsaCrypto := NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, sha256.New)

	// 3. 原始数据
	data := []byte("Hello, RSA performance testing!")

	// 4. 生成随机盐
	salt := make([]byte, 16) // 16字节的盐
	if _, err := rand.Read(salt); err != nil {
		b.Fatalf(errSaltGeneration, err)
	}

	// 5. 执行基准测试
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := rsaCrypto.EncryptSalt(data, salt)
		assert.NoError(b, err, errEncryption, err)
	}
}

func BenchmarkRsaDecryption(b *testing.B) {
	// 1. 生成 RSA 密钥对
	keyPair := generateRsaKeyPair(b, RsaKeySize2048)

	// 2. 创建 RSA 加解密器
	rsaCrypto := NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, sha256.New)

	// 3. 原始数据
	data := []byte("Hello, RSA performance testing!")

	// 4. 生成随机盐
	salt := make([]byte, 16) // 16字节的盐
	if _, err := rand.Read(salt); err != nil {
		b.Fatalf(errSaltGeneration, err)
	}

	// 5. 先加密以获取加密数据
	encrypted, err := rsaCrypto.EncryptSalt(data, salt)
	assert.NoError(b, err, errEncryption, err)

	// 6. 执行基准测试
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := rsaCrypto.Decrypt(encrypted)
		assert.NoError(b, err, errDecryption, err)
	}
}
