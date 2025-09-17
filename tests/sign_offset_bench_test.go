/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-16 18:55:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-17 09:58:26
 * @FilePath: \go-toolbox\tests\sign_offset_bench_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/sign"
)

func BenchmarkProtonOffsetCipherEncrypt(b *testing.B) {
	psc := sign.NewProtonOffsetCipher()
	data := []byte("Hello, World!")

	for i := 0; i < b.N; i++ {
		_, err := psc.Encrypt(data)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}
	}
}

func BenchmarkProtonOffsetCipherDecrypt(b *testing.B) {
	psc := sign.NewProtonOffsetCipher()
	data := []byte("Hello, World!")
	encrypted, _ := psc.Encrypt(data)

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := psc.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}

func BenchmarkProtonOffsetCipherEncryptLong(b *testing.B) {
	psc := sign.NewProtonOffsetCipher()
	data := make([]byte, 1024*1024) // 1 MB 数据
	for i := range data {
		data[i] = 'A' // 填充数据
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := psc.Encrypt(data)
		if err != nil {
			b.Fatalf("Encryption failed: %v", err)
		}
	}
}

func BenchmarkProtonOffsetCipherDecryptLong(b *testing.B) {
	psc := sign.NewProtonOffsetCipher()
	data := make([]byte, 1024*1024) // 1 MB 数据
	for i := range data {
		data[i] = 'A' // 填充数据
	}
	encrypted, _ := psc.Encrypt(data)

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := psc.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decryption failed: %v", err)
		}
	}
}
