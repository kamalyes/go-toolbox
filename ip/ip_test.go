/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-02 10:13:03
 * @FilePath: \go-toolbox\ip\ip_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package ip

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPFunctions(t *testing.T) {
	t.Run("TestGetClientIP", TestGetClientIP)
	t.Run("TestGetLocalIp", TestGetLocalIp)
	t.Run("TestGetPublicIP", TestGetPublicIP)
	t.Run("TestIsIpAddress", TestIsIpAddress)
	t.Run("TestIsIPv4", TestIsIPv4)
	t.Run("TestIsIPv6", TestIsIPv6)
}

func TestGetClientIP(t *testing.T) {
	// 创建一个模拟的http.Request对象
	req, err := http.NewRequest("GET", "https://pkg.go.dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	GetClientIP(req)
}

func TestGetLocalIp(t *testing.T) {
	// 测试GetLocalIp函数
	ips, err := GetLocalIp()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(ips), 0)
}

func TestGetPublicIP(t *testing.T) {
	// 测试GetPublicIP函数
	publicIP, err := GetPublicIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, publicIP)
}

func TestIsIpAddress(t *testing.T) {
	// 测试IsIpAddress函数
	validIP := "192.168.1.1"
	invalidIP := "invalidip"
	assert.True(t, IsIpAddress(validIP))
	assert.False(t, IsIpAddress(invalidIP))
}

func TestIsIPv4(t *testing.T) {
	// 测试IsIPv4函数
	validIPv4 := "192.168.1.1"
	invalidIPv4 := "invalidipv4"
	assert.True(t, IsIPv4(validIPv4))
	assert.False(t, IsIPv4(invalidIPv4))
}

func TestIsIPv6(t *testing.T) {
	// 测试IsIPv6函数
	validIPv6 := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	invalidIPv6 := "invalidipv6"
	assert.True(t, IsIPv6(validIPv6))
	assert.False(t, IsIPv6(invalidIPv6))
}
