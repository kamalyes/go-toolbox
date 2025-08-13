/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 11:25:16
 * @FilePath: \go-toolbox\tests\aes_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/go-toolbox/pkg/sign"
)

func BenchmarkAesEncryptDecrypt(b *testing.B) {
	var password = "example1235678"
	var byteKey = sign.GenerateByteKey(password, 32)

	// 生成随机字符串作为测试输入
	plainText := random.FRandString(4096) // 4 KB

	b.Run("EncryptDecrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			encryptedText, err := sign.AesEncrypt(plainText, byteKey)
			if err != nil {
				b.Fatalf("Encryption failed: %v", err)
			}

			decryptedText, err := sign.AesDecrypt(encryptedText, byteKey)
			if err != nil {
				b.Fatalf("Decryption failed: %v", err)
			}

			if decryptedText != plainText {
				b.Errorf("Decrypted text does not match original. Got: %s, Want: %s", decryptedText, plainText)
			}
		}
	})
}
