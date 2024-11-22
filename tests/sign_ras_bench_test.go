/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:55:15
 * @FilePath: \go-toolbox\tests\sign_ras_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/sign"
)

// 生成 RSA 密钥对的辅助函数
func generateRsaKeyPair(b *testing.B, keySize sign.RsaKeySize) *sign.RsaKeyPair {
	keyPair, err := sign.GenerateRsaKeyPair(keySize)
	if err != nil {
		b.Fatal()
	}
	return keyPair
}

func BenchmarkRsaEncryption(b *testing.B) {
	// 1. 生成 RSA 密钥对
	keyPair := generateRsaKeyPair(b, sign.RsaKeySize2048)

	// 2. 创建 RSA 加解密器
	rsaCrypto := sign.NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, sha256.New)

	// 3. 原始数据
	data := []byte("Hello, RSA performance testing!")

	// 4. 生成随机盐
	salt := make([]byte, 16) // 16字节的盐
	if _, err := rand.Read(salt); err != nil {
		b.Fatalf("盐生成失败: %v", err)
	}

	// 5. 执行基准测试
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := rsaCrypto.EncryptSalt(data, salt)
		if err != nil {
			b.Fatalf("加密失败: %v", err)
		}
	}
}

func BenchmarkRsaDecryption(b *testing.B) {
	// 1. 生成 RSA 密钥对
	keyPair := generateRsaKeyPair(b, sign.RsaKeySize2048)

	// 2. 创建 RSA 加解密器
	rsaCrypto := sign.NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, sha256.New)

	// 3. 原始数据
	data := []byte("Hello, RSA performance testing!")

	// 4. 生成随机盐
	salt := make([]byte, 16) // 16字节的盐
	if _, err := rand.Read(salt); err != nil {
		b.Fatalf("盐生成失败: %v", err)
	}

	// 5. 先加密以获取加密数据
	encrypted, err := rsaCrypto.EncryptSalt(data, salt)
	if err != nil {
		b.Fatalf("加密失败: %v", err)
	}

	// 6. 执行基准测试
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := rsaCrypto.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("解密失败: %v", err)
		}
	}
}
