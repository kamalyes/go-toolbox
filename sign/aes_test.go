/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-23 17:37:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-23 17:50:50
 * @FilePath: \go-toolbox\sign\aes.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"encoding/base64"
	"testing"
)

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

		decryptedText, err := AesDecrypt(encryptedText, tc.key)

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
