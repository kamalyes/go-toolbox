/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 01:55:19
 * @FilePath: \go-toolbox\tests\next_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"net/http"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/next"
	"github.com/stretchr/testify/assert"
)

func TestIPFunctions(t *testing.T) {
	t.Run("TestGetClientIP", TestGetClientIP)
	t.Run("TestGetLocalIp", TestGetLocalIp)
	t.Run("TestGetPublicIP", TestGetPublicIP)
	t.Run("TestIsIpAddress", TestIsIpAddress)
	t.Run("TestIsIPv4", TestIsIPv4)
	t.Run("TestIsIPv6", TestIsIPv6)
	t.Run("TestHasLocalIPAddr", TestHasLocalIPAddr)
}

func TestGetClientIP(t *testing.T) {
	// 创建一个模拟的http.Request对象
	req, err := http.NewRequest("GET", "https://pkg.go.dev", nil)
	if err != nil {
		t.Fatal(err)
	}

	next.GetClientIP(req)
}

func TestGetLocalIp(t *testing.T) {
	// 测试GetLocalIp函数
	ips, err := next.GetLocalIp()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(ips), 0)
}

func TestGetPublicIP(t *testing.T) {
	// 测试GetPublicIP函数
	publicIP, err := next.GetPublicIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, publicIP)
}

func TestIsIpAddress(t *testing.T) {
	// 测试IsIpAddress函数
	validIP := "192.168.1.1"
	invalidIP := "invalidip"
	assert.True(t, next.IsIpAddress(validIP))
	assert.False(t, next.IsIpAddress(invalidIP))
}

func TestIsIPv4(t *testing.T) {
	// 测试IsIPv4函数
	validIPv4 := "192.168.1.1"
	invalidIPv4 := "invalidipv4"
	assert.True(t, next.IsIPv4(validIPv4))
	assert.False(t, next.IsIPv4(invalidIPv4))
}

func TestIsIPv6(t *testing.T) {
	// 测试IsIPv6函数
	validIPv6 := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	invalidIPv6 := "invalidipv6"
	assert.True(t, next.IsIPv6(validIPv6))
	assert.False(t, next.IsIPv6(invalidIPv6))
}

func TestHasLocalIPAddr(t *testing.T) {
	tests := []struct {
		ip      string
		isLocal bool
	}{
		{"127.0.0.1", true},        // 回环地址
		{"10.0.0.1", true},         // 私有地址
		{"172.16.0.1", true},       // 私有地址
		{"192.168.1.1", true},      // 私有地址
		{"169.254.1.1", true},      // 链路本地地址
		{"8.8.8.8", false},         // 公共地址
		{"255.255.255.255", false}, // 广播地址
		{"::1", true},              // IPv6 回环地址
		{"fc00::1", true},          // 唯一本地地址
		{"fe80::1", true},          // 链路本地地址
		{"2001:db8::1", false},     // 公共地址
		{"invalid-ip", false},      // 无效地址
	}

	for _, test := range tests {
		result := next.HasLocalIPAddr(test.ip)
		if result != test.isLocal {
			t.Errorf("HasLocalIPAddr(%q) = %v; want %v", test.ip, result, test.isLocal)
		}
	}
}
