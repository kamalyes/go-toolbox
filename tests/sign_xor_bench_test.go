/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-17 10:06:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-17 10:06:55
 * @FilePath: \go-toolbox\tests\sign_xor_bench_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/sign"
)

func BenchmarkXORCipherEncrypt(b *testing.B) {
	key := byte(0xAA)
	xor := sign.NewXORCipher(key)
	data := []byte("Hello, World!")

	for i := 0; i < b.N; i++ {
		_, err := xor.Encrypt(data)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}
	}
}

func BenchmarkXORCipherDecrypt(b *testing.B) {
	key := byte(0xAA)
	xor := sign.NewXORCipher(key)
	data := []byte("Hello, World!")
	encrypted, _ := xor.Encrypt(data)

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := xor.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}

func BenchmarkXORCipherEncryptLong(b *testing.B) {
	key := byte(0xAA)
	xor := sign.NewXORCipher(key)
	data := make([]byte, 1024*1024) // 1 MB 数据
	for i := range data {
		data[i] = 'A' // 填充数据
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := xor.Encrypt(data)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}
	}
}

func BenchmarkXORCipherDecryptLong(b *testing.B) {
	key := byte(0xAA)
	xor := sign.NewXORCipher(key)
	data := make([]byte, 1024*1024) // 1 MB 数据
	for i := range data {
		data[i] = 'A' // 填充数据
	}
	encrypted, _ := xor.Encrypt(data)

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := xor.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}
