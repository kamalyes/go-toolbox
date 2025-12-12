/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-17 10:06:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-17 10:06:55
 * @FilePath: \go-toolbox\pkg\sign\xor_bench_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkXORCipherEncrypt(b *testing.B) {
	key := byte(0xAA)
	xor := NewXORCipher(key)
	data := []byte("Hello, World!")

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := xor.Encrypt(data)
		if err != nil {
			assert.Fail(b, errEncryption, err)
		}
	}
}

func BenchmarkXORCipherDecrypt(b *testing.B) {
	key := byte(0xAA)
	xor := NewXORCipher(key)
	data := []byte("Hello, World!")
	encrypted, err := xor.Encrypt(data)
	if err != nil {
		b.Fatalf(errEncryption, err)
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := xor.Decrypt(encrypted)
		if err != nil {
			assert.Fail(b, errDecryption, err)
		}
	}
}

func BenchmarkXORCipherEncryptLong(b *testing.B) {
	key := byte(0xAA)
	xor := NewXORCipher(key)
	data := make([]byte, 1024*1024) // 1 MB 数据
	for i := range data {
		data[i] = 'A' // 填充数据
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := xor.Encrypt(data)
		if err != nil {
			assert.Fail(b, errEncryption, err)
		}
	}
}

func BenchmarkXORCipherDecryptLong(b *testing.B) {
	key := byte(0xAA)
	xor := NewXORCipher(key)
	data := make([]byte, 1024*1024) // 1 MB 数据
	for i := range data {
		data[i] = 'A' // 填充数据
	}
	encrypted, err := xor.Encrypt(data)
	if err != nil {
		b.Fatalf(errEncryption, err)
	}

	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		_, err := xor.Decrypt(encrypted)
		if err != nil {
			assert.Fail(b, errDecryption, err)
		}
	}
}
