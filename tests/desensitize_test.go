/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 01:53:57
 * @FilePath: \go-toolbox\tests\desensitize_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/desensitize"
	"github.com/stretchr/testify/assert"
)

func TestDesensitizeAllFunctions(t *testing.T) {
	t.Run("TestSensitiveData", TestSensitiveData)
	t.Run("TestChinesName", TestChinesName)
	t.Run("TestPhoneNumber", TestPhoneNumber)
	t.Run("TestCarLicense", TestCarLicense)
	t.Run("TestBankCard", TestBankCard)
	t.Run("TestIPv4", TestIPv4)
	t.Run("TestIPv6", TestIPv6)
}

func TestSensitiveData(t *testing.T) {
	assert.Equal(t, "1***", desensitize.SensitiveData("1234", 0, 0))
	assert.Equal(t, "1***", desensitize.SensitiveData("1234", 2, 2))
	assert.Equal(t, "12*4", desensitize.SensitiveData("1234", 2, 3))
	assert.Equal(t, "1***", desensitize.SensitiveData("1234", 0, 5))
	assert.Equal(t, "", desensitize.SensitiveData("", 0, 0))
}

func TestChinesName(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected string
		option   desensitize.DesensitizeOptions
	}{
		"ChineseNameNoDesensitizeOptions": {
			input:    "石晓浩",
			expected: "石**",
		},
		"ChineseNameNULLDesensitizeOptions": {
			input:    "石晓浩",
			expected: "石**",
			option:   desensitize.DesensitizeOptions{},
		},
		"ChineseNameErrIndex": {
			input:    "阿哈利",
			expected: "阿哈利",
			option:   desensitize.DesensitizeOptions{ChineseNameStartIndex: 2, ChineseNameEndIndex: 1},
		},
		"ChineseNameIndex": {
			input:    "阿哈利晓浩",
			expected: "阿哈利*浩",
			option:   desensitize.DesensitizeOptions{ChineseNameStartIndex: 3, ChineseNameEndIndex: 4},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, desensitize.Desensitize(tc.input, desensitize.ChineseName, tc.option))
		})
	}
}

func TestPhoneNumber(t *testing.T) {
	assert.Equal(t, "181****8789", desensitize.SensitizePhoneNumber("18175698789", 3, 7))
	assert.Equal(t, "181*****789", desensitize.SensitizePhoneNumber("1817789", 3, 7))
}

func TestCarLicense(t *testing.T) {
	assert.Equal(t, "浙A1****B", desensitize.SensitizeCarLicense("浙A12345B"))
	assert.Equal(t, "浙A1****Z", desensitize.SensitizeCarLicense("浙A12345Z"))
	assert.Equal(t, "", desensitize.SensitizeCarLicense(""))
}

func TestBankCard(t *testing.T) {
	assert.Equal(t, "1234 **** **** *** 6789", desensitize.SensitizeBankCard("123456789", 16))
	assert.Equal(t, "1234 **** **** *** 6789", desensitize.SensitizeBankCard("1234 5678 9", 16))
	assert.Equal(t, "1234 **** **** **** 6789", desensitize.SensitizeBankCard("123456789", 19))
	assert.Equal(t, "1234 **** **** **** 6789", desensitize.SensitizeBankCard("1234 5678 9", 19))
	assert.Equal(t, "", desensitize.SensitizeBankCard("", 19))
}

func TestIPv4(t *testing.T) {
	assert.Equal(t, "127.*.*.*", desensitize.SensitizeIpv4("127.0.0.1"))
	assert.Equal(t, "", desensitize.SensitizeIpv4(""))
}

func TestIPv6(t *testing.T) {
	assert.Equal(t, "2001:*:*:*:*:*:*:*", desensitize.SensitizeIpv6("2001:0db8:86a3:08d3:1319:8a2e:0370:7344"))
	assert.Equal(t, "", desensitize.SensitizeIpv6(""))
}
