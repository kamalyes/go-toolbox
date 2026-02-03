/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-23 17:37:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:15:05
 * @FilePath: \go-toolbox\pkg\sign\aes_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func TestAesEncryptDecrypt(t *testing.T) {
	var password = "example1235678"
	var byteKey = GenerateByteKey(password, 32)

	tamperedCiphertexts := map[string]string{
		"tampered ciphertext": "tampered",
	}

	testCases := []struct {
		name             string
		plainText        string
		expectEncryptErr bool
		expectDecryptErr bool
		key              []byte
	}{
		{"normal", "Hello, World!", false, false, byteKey},
		{"empty string", "", false, false, byteKey},
		{"long string", "A long string that exceeds the typical block size to test the AES encryption and decryption functionality.", false, false, byteKey},
		{"special characters", "Special characters: !@#$%^&*()_+[]{}|;':\",.<>?/`~", false, false, byteKey},
		{"unicode", "Unicode test: 你好，世界！", false, false, byteKey},
		{"tampered ciphertext", "Hello, World!", false, true, byteKey},
		{"non-string input", string([]byte{0x00, 0x01}), false, false, byteKey},
		{"empty key", "Hello, World!", true, true, []byte{}},
	}

	for _, tc := range testCases {
		encryptedText, err := AesEncrypt(tc.plainText, tc.key)

		if tc.expectEncryptErr {
			assert.Error(t, err, "%s: Expected error for encryption, got none", tc.name)
			continue // Skip decryption step
		}

		assert.NoError(t, err, "%s: Encryption failed: %v", tc.name, err)

		// Check for tampered ciphertext
		if tamperedText, exists := tamperedCiphertexts[tc.name]; exists {
			encryptedText = tamperedText // Use tampered ciphertext
		}

		decryptedText, err := AesDecrypt(encryptedText, tc.key)

		if tc.expectDecryptErr {
			assert.Error(t, err, "%s: Expected error for decryption, got none", tc.name)
		} else {
			assert.NoError(t, err, "%s: Decryption failed: %v", tc.name, err)
			assert.Equal(t, tc.plainText, decryptedText, "%s: Decrypted text does not match original. Got: %s, Want: %s", tc.name, decryptedText, tc.plainText)
		}

		// 在一条日志中打印密文、密钥和解密后的值
		t.Logf("Test Case: %s | Encrypted Text (Base64): %s | Key (Base64): %s | Decrypted Text: %s",
			tc.name,
			base64.StdEncoding.EncodeToString([]byte(encryptedText)),
			base64.StdEncoding.EncodeToString(tc.key),
			decryptedText)
	}
}

func TestAes(t *testing.T) {
	password := "mysecretpassword"
	keyLength := 16 // AES-128
	key := GenerateByteKey(password, keyLength)

	tests := []interface{}{
		"Hello, World!",
		12345,
		3.14159265359,
		true,
		[]byte{1, 2, 3, 4, 5},
		[]int{1, 2, 3, 4, 5},
		[]float64{1.1, 2.2, 3.3, 4.4, 5.5},
		"中文测试",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	}

	for _, test := range tests {
		fmt.Println("======================================")
		fmt.Printf("Original Text: %v\n", test)

		// 转换为字符串
		originalText := fmt.Sprintf("%v", test)
		fmt.Printf("Original Text (String): %v\n", originalText)

		// 加密
		encryptedText, err := AesEncrypt(originalText, key)
		assert.NoError(t, err, "AesEncrypt error: %v", err)
		if err == nil {
			fmt.Printf("Encrypted Text: %v\n", encryptedText)

			// 解密
			decryptedText, err := AesDecrypt(encryptedText, key)
			assert.NoError(t, err, "AesDecrypt error: %v", err)
			assert.Equal(t, originalText, decryptedText, "Decrypted text does not match the original text")
		}

		fmt.Println("======================================")
	}
}

func BenchmarkAesEncryptDecrypt(b *testing.B) {
	var password = "example1235678"
	var byteKey = GenerateByteKey(password, 32)

	// 生成随机字符串作为测试输入
	plainText := generateRandomString(4096) // 4 KB

	b.Run("EncryptDecrypt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			encryptedText, err := AesEncrypt(plainText, byteKey)
			if err != nil {
				b.Fatalf("Encryption failed: %v", err)
			}

			decryptedText, err := AesDecrypt(encryptedText, byteKey)
			if err != nil {
				b.Fatalf("Decryption failed: %v", err)
			}

			assert.Equal(b, plainText, decryptedText, "Decrypted text does not match original. Got: %s, Want: %s", decryptedText, plainText)
		}
	})
}

// TestAesEncryptDecryptWithIV 测试自定义 IV 的加密解密
func TestAesEncryptDecryptWithIV(t *testing.T) {
	password := "example1235678"
	key := GenerateByteKey(password, 32)
	iv := []byte("gtfrdbhytfredrji") // 16 字节 IV

	testCases := []struct {
		name             string
		plainText        string
		key              []byte
		iv               []byte
		expectEncryptErr bool
		expectDecryptErr bool
	}{
		{"normal", "Hello, World!", key, iv, false, false},
		{"empty string", "", key, iv, false, false},
		{"long string", "A long string that exceeds the typical block size to test the AES encryption and decryption functionality.", key, iv, false, false},
		{"special characters", "Special characters: !@#$%^&*()_+[]{}|;':\",.<>?/`~", key, iv, false, false},
		{"unicode", "Unicode test: 你好，世界！", key, iv, false, false},
		{"json data", `{"key":"test/upload.txt"}`, key, iv, false, false},
		{"empty key", "Hello, World!", []byte{}, iv, true, true},
		{"invalid IV length", "Hello, World!", key, []byte("short"), true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encryptedText, err := AesEncryptWithIV(tc.plainText, tc.key, tc.iv)

			if tc.expectEncryptErr {
				assert.Error(t, err, "Expected error for encryption")
				return
			}

			assert.NoError(t, err, "Encryption failed: %v", err)
			assert.NotEmpty(t, encryptedText, "Encrypted text should not be empty")

			decryptedText, err := AesDecryptWithIV(encryptedText, tc.key, tc.iv)

			if tc.expectDecryptErr {
				assert.Error(t, err, "Expected error for decryption")
				return
			}

			assert.NoError(t, err, "Decryption failed: %v", err)
			assert.Equal(t, tc.plainText, decryptedText, "Decrypted text does not match original")

			t.Logf("Encrypted (Base64): %s | Decrypted: %s", encryptedText, decryptedText)
		})
	}
}

// TestAesWithIVCompatibility 测试与 Cloudflare Worker 的兼容性
func TestAesWithIVCompatibility(t *testing.T) {
	// 使用与 Cloudflare Worker 相同的配置
	key := []byte("mkjnhbgvfrquedhsgdbchgyutrfdhsij") // 32 字节
	iv := []byte("gtfrdbhytfredrji")                  // 16 字节

	// 测试 JSON 数据（模拟 R2 上传签名）
	jsonData := `{"key":"test/upload.txt"}`

	encrypted, err := AesEncryptWithIV(jsonData, key, iv)
	assert.NoError(t, err, "Encryption failed")
	assert.NotEmpty(t, encrypted, "Encrypted text should not be empty")

	decrypted, err := AesDecryptWithIV(encrypted, key, iv)
	assert.NoError(t, err, "Decryption failed")
	assert.Equal(t, jsonData, decrypted, "Decrypted text does not match original")

	t.Logf("Original: %s", jsonData)
	t.Logf("Encrypted: %s", encrypted)
	t.Logf("Decrypted: %s", decrypted)
}

// BenchmarkAesEncryptDecryptWithIV 性能测试
func BenchmarkAesEncryptDecryptWithIV(b *testing.B) {
	password := "example1235678"
	key := GenerateByteKey(password, 32)
	iv := []byte("gtfrdbhytfredrji")
	plainText := generateRandomString(4096) // 4 KB

	b.Run("EncryptDecryptWithIV", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			encryptedText, err := AesEncryptWithIV(plainText, key, iv)
			if err != nil {
				b.Fatalf("Encryption failed: %v", err)
			}

			decryptedText, err := AesDecryptWithIV(encryptedText, key, iv)
			if err != nil {
				b.Fatalf("Decryption failed: %v", err)
			}

			if plainText != decryptedText {
				b.Fatalf("Decrypted text does not match original")
			}
		}
	})
}
