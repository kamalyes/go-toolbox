/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-02 16:17:50
 * @FilePath: \go-toolbox\desensitize\desensitize_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package desensitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateAllFunctions(t *testing.T) {
	t.Run("TestSensitiveData", TestSensitiveData)
	t.Run("TestChinesName", TestChinesName)
	t.Run("TestPhoneNumber", TestPhoneNumber)
	t.Run("TestCarLicense", TestCarLicense)
	t.Run("TestBankCard", TestBankCard)
	t.Run("TestIPv4", TestIPv4)
	t.Run("TestIPv6", TestIPv6)
}

func TestSensitiveData(t *testing.T) {
	assert.Equal(t, "1***", SensitiveData("1234", 0, 0))
	assert.Equal(t, "1***", SensitiveData("1234", 2, 2))
	assert.Equal(t, "12*4", SensitiveData("1234", 2, 3))
	assert.Equal(t, "1***", SensitiveData("1234", 0, 5))
	assert.Equal(t, "", SensitiveData("", 0, 0))
}

func TestChinesName(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected string
		option   DesensitizeOptions
	}{
		"ChineseNameNoDesensitizeOptions": {
			input:    "石晓浩",
			expected: "石**",
		},
		"ChineseNameNULLDesensitizeOptions": {
			input:    "石晓浩",
			expected: "石**",
			option:   DesensitizeOptions{},
		},
		"ChineseNameErrIndex": {
			input:    "阿哈利",
			expected: "阿哈利",
			option:   DesensitizeOptions{ChineseNameStartIndex: 2, ChineseNameEndIndex: 1},
		},
		"ChineseNameIndex": {
			input:    "阿哈利晓浩",
			expected: "阿哈利*浩",
			option:   DesensitizeOptions{ChineseNameStartIndex: 3, ChineseNameEndIndex: 4},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, Desensitize(tc.input, ChineseName, tc.option))
		})
	}
}

func TestPhoneNumber(t *testing.T) {
	assert.Equal(t, "181****8789", phoneNumber("18175698789", 3, 7))
	assert.Equal(t, "181*****789", phoneNumber("1817789", 3, 7))
}

func TestCarLicense(t *testing.T) {
	assert.Equal(t, "浙A1****B", carLicense("浙A12345B"))
	assert.Equal(t, "浙A1****Z", carLicense("浙A12345Z"))
	assert.Equal(t, "", carLicense(""))
}

func TestBankCard(t *testing.T) {
	assert.Equal(t, "1234 **** **** *** 6789", bankCard("123456789", 16))
	assert.Equal(t, "1234 **** **** *** 6789", bankCard("1234 5678 9", 16))
	assert.Equal(t, "1234 **** **** **** 6789", bankCard("123456789", 19))
	assert.Equal(t, "1234 **** **** **** 6789", bankCard("1234 5678 9", 19))
	assert.Equal(t, "", bankCard("", 19))
}

func TestIPv4(t *testing.T) {
	assert.Equal(t, "127.*.*.*", ipv4("127.0.0.1"))
	assert.Equal(t, "", ipv4(""))
}

func TestIPv6(t *testing.T) {
	assert.Equal(t, "2001:*:*:*:*:*:*:*", ipv6("2001:0db8:86a3:08d3:1319:8a2e:0370:7344"))
	assert.Equal(t, "", ipv6(""))
}
