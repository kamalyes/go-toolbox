/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-19 11:19:10
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

func TestSensitiveData(t *testing.T) {
	assert.Equal(t, "1***", desensitize.SensitiveData("1234", 0, 0))
	assert.Equal(t, "1***", desensitize.SensitiveData("1234", 2, 2))
	assert.Equal(t, "12*4", desensitize.SensitiveData("1234", 2, 3))
	assert.Equal(t, "1***", desensitize.SensitiveData("1234", 0, 5))
	assert.Equal(t, "", desensitize.SensitiveData("", 0, 0))
}

func TestDesensitizeAllTypes(t *testing.T) {
	desensitizeOptions := desensitize.NewDesensitizeOptions()
	testCases := map[string]struct {
		input           string
		expected        string
		desensitizeType desensitize.DesensitizeType
		option          desensitize.DesensitizeOptions
	}{
		"TestSensitiveData": {
			input:           "123456789",
			expected:        "123****89",
			desensitizeType: desensitize.CustomExtension,
			option:          desensitize.DesensitizeOptions{CustomExtensionStartIndex: 3, CustomExtensionEndIndex: 7},
		},
		"TestChineseName": {
			input:           "李四",
			expected:        "李*",
			desensitizeType: desensitize.ChineseName,
		},
		"TestIDCard": {
			input:           "123456789012345678",
			expected:        "123456********5678",
			desensitizeType: desensitize.IDCard,
			option:          desensitizeOptions,
		},
		"TestPassWord": {
			input:           "12678",
			expected:        "1****",
			desensitizeType: desensitize.Password,
		},
		"TestPhoneNumber": {
			input:           "18175698789",
			expected:        "181****8789",
			desensitizeType: desensitize.PhoneNumber,
			option:          desensitizeOptions,
		},
		"TestMobilePhone": {
			input:           "13812345678",
			expected:        "138****5678",
			desensitizeType: desensitize.MobilePhone,
			option:          desensitizeOptions,
		},
		"TestAddress": {
			input:           "北京市朝阳区某某街道123号",
			expected:        "北京市朝*******23号",
			desensitizeType: desensitize.Address,
			option:          desensitizeOptions,
		},
		"TestEmail": {
			input:           "example@test.com",
			expected:        "e****le@test.com",
			desensitizeType: desensitize.Email,
			option:          desensitizeOptions,
		},
		"TestCarLicense": {
			input:           "浙A12345B",
			expected:        "浙A1***5B",
			desensitizeType: desensitize.CarLicense,
		},
		"TestBankCard": {
			input:           "1234567890123456",
			expected:        "1234 **** **** **** 3456",
			desensitizeType: desensitize.BankCard,
			option:          desensitizeOptions,
		},
		"TestIPv4": {
			input:           "192.168.1.1",
			expected:        "192.*.*.*",
			desensitizeType: desensitize.IPV4,
		},
		"TestIPv6": {
			input:           "2001:0db8:86a3:08d3:1319:8a2e:0370:7344",
			expected:        "2001:*:*:*:*:*:*:*",
			desensitizeType: desensitize.IPV6,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, desensitize.Desensitize(tc.input, tc.desensitizeType, tc.option))
		})
	}
}

func TestPhoneNumber(t *testing.T) {
	assert.Equal(t, "", desensitize.SensitizePhoneNumber("", 3, 7))
	assert.Equal(t, "181****8789", desensitize.SensitizePhoneNumber("18175698789", 3, 7))
	assert.Equal(t, "181*****789", desensitize.SensitizePhoneNumber("1817789", 3, 7))
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
