/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:30:33
 * @FilePath: \go-toolbox\tests\sign_ras_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
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

// 测试 RSA 密钥对生成性能
func BenchmarkRsaKeyPairGeneration(b *testing.B) {
	keySize := sign.RsaKeySize4096

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		generateRsaKeyPair(b, keySize)
	}
}

// 测试 RSA 加解密性能
func BenchmarkRsaEncrypt(b *testing.B) {
	keySize := sign.RsaKeySize4096
	keyPair := generateRsaKeyPair(b, keySize)

	// 测试数据
	data := []byte("This is a test message for RSA encryption performance testing.")

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		crypto := sign.NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, keySize)
		_, err := crypto.Encrypt(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRsaDecrypt(b *testing.B) {
	keySize := sign.RsaKeySize4096
	keyPair := generateRsaKeyPair(b, keySize)

	// 测试数据
	data := []byte("This is a test message for RSA encryption performance testing.")
	crypto := sign.NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, keySize)

	// 先加密数据以便后续解密测试
	encryptedData, err := crypto.Encrypt(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := crypto.Decrypt(encryptedData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
