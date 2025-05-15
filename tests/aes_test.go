/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-23 17:37:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-05-15 17:35:46
 * @FilePath: \go-toolbox\tests\aes_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/sign"
	"github.com/stretchr/testify/assert"
)

func TestAesEncryptDecrypt(t *testing.T) {
	var password = "example1235678"
	var byteKey = sign.GenerateByteKey(password, 32)

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
		encryptedText, err := sign.AesEncrypt(tc.plainText, tc.key)

		if tc.expectEncryptErr {
			if err == nil {
				t.Errorf("%s: Expected error for encryption, got none", tc.name)
			}
			continue // Skip decryption step
		}

		if err != nil {
			t.Fatalf("%s: Encryption failed: %v", tc.name, err)
		}

		// Check for tampered ciphertext
		if tamperedText, exists := tamperedCiphertexts[tc.name]; exists {
			encryptedText = tamperedText // Use tampered ciphertext
		}

		decryptedText, err := sign.AesDecrypt(encryptedText, tc.key)

		if tc.expectDecryptErr {
			if err == nil {
				t.Errorf("%s: Expected error for decryption, got none", tc.name)
			}
		} else {
			if err != nil {
				t.Fatalf("%s: Decryption failed: %v", tc.name, err)
			}

			if decryptedText != tc.plainText {
				t.Errorf("%s: Decrypted text does not match original. Got: %s, Want: %s", tc.name, decryptedText, tc.plainText)
			}
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
	key := sign.GenerateByteKey(password, keyLength)

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
		encryptedText, err := sign.AesEncrypt(originalText, key)
		if err != nil {
			fmt.Printf("AesEncrypt error: %v\n", err)
		} else {
			fmt.Printf("Encrypted Text: %v\n", encryptedText)
		}

		// 解密
		decryptedText, err := sign.AesDecrypt(encryptedText, key)
		if err != nil {
			fmt.Printf("AesDecrypt error: %v\n", err)
		} else {
			fmt.Printf("Decrypted Text: %v\n", decryptedText)
		}

		// 验证
		assert.Nil(t, err, "AesDecrypt error")
		assert.Equal(t, originalText, decryptedText, "Decrypted text does not match the original text")
		if assert.Equal(t, originalText, decryptedText) {
			fmt.Println("Verification: PASSED")
		} else {
			fmt.Println("Verification: FAILED")
		}
		fmt.Println("======================================")
	}
}
