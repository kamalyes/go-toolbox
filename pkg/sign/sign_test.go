/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-02 10:30:31
 * @FilePath: \go-toolbox\pkg\sign\sign_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"testing"
)

func TestSign(t *testing.T) {
	t.Run("TestHmacSha256Base64", TestHmacSha256Base64)
	t.Run("TestHmacSha256Hex", TestHmacSha256Hex)
	t.Run("TestSHA256", TestSHA256)
}

func TestHmacSha256Base64(t *testing.T) {
	message := "Hello, world!"
	secret := "mysecret"
	expected := "x14OheqQEArT8lK2ruEjRGcpU+djdJINAThNE0iWK7g="

	signature := HmacSha256Base64(message, secret)
	if signature != expected {
		t.Errorf("HmacSha256Base64: Expected: %s, but got: %s", expected, signature)
	}
}

func TestHmacSha256Hex(t *testing.T) {
	message := "Hello, world!"
	secret := "mysecret"
	expected := "c75e0e85ea90100ad3f252b6aee12344672953e76374920d01384d1348962bb8"

	signature := HmacSha256Hex(message, secret)
	if signature != expected {
		t.Errorf("HmacSha256Hex: Expected: %s, but got: %s", expected, signature)
	}

	// Test with empty secret
	signature = HmacSha256Hex(message, "")
	expected = "0d192eb5bc5e4407192197cbf9e1658295fa3ff995b3ff914f3cc7c38d83b10f"
	if signature != expected {
		t.Errorf("HmacSha256Hex with empty Expected: %s, but got: %s", expected, signature)
	}
}

func TestSHA256(t *testing.T) {
	text := "Hello, world!"
	expected := "0c6c7bc085d215537fba23b05171d387534160c7eaa1e5b448d5dea9be498c39"

	signature := SHA256(text)
	if signature != expected {
		t.Errorf("SHA256: Expected: %s, but got: %s", expected, signature)
	}

	// Test with special characters
	text = "Hello, $%^&*"
	expected = "549803836a05c042130331d661677917458a5cedd151fbab04fc6b60ac981031"
	signature = SHA256(text)
	if signature != expected {
		t.Errorf("SHA256 with special characters: Expected: %s, but got: %s", expected, signature)
	}
}
