/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 17:21:35
 * @FilePath: \go-toolbox\tests\osx_base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/stretchr/testify/assert"
)

func TestAllSysBaseFunctions(t *testing.T) {
	t.Run("TestSafeGetHostName", TestSafeGetHostName)
	t.Run("TestHashUnixMicroCipherText", TestHashUnixMicroCipherText)
}

func TestSafeGetHostName(t *testing.T) {
	actual := osx.SafeGetHostName()
	assert.NotEmpty(t, actual, "HostNames should match")
}

// TestHashUnixMicroCipherText 测试 HashUnixMicroCipherText 函数
func TestHashUnixMicroCipherText(t *testing.T) {
	hash1 := osx.HashUnixMicroCipherText()
	hash2 := osx.HashUnixMicroCipherText()

	// 验证生成的哈希值不为空
	assert.NotEqual(t, hash1, "")
	assert.NotEqual(t, hash2, "")
	assert.Equal(t, len(hash1), 32)
	assert.NotEqual(t, hash1, hash2)
}

func TestGetServerIP(t *testing.T) {
	externalIP, internalIP, err := osx.GetLocalInterfaceIeIp()
	assert.Nil(t, err)
	if externalIP != "" {
		t.Logf("externalIP %s", externalIP)
	}
	if internalIP != "" {
		t.Logf("internalIP %s", internalIP)
	}
}

func TestGetLocalInterfaceIps(t *testing.T) {
	ips, err := osx.GetLocalInterfaceIps()
	assert.Nil(t, err)
	assert.NotEmpty(t, ips, fmt.Sprintf("Expected at least one global unicast IP, got: %v", ips))
	for _, ip := range ips {
		assert.NotEmpty(t, ip, fmt.Sprintf("Invalid IP address: %s", ip))
	}
}

func TestGetClientPublicIP_XForwardedFor(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	testIp := "1.2.3.4"
	req.Header.Set("X-Forwarded-For", testIp)
	ip, err := osx.GetClientPublicIP(req)
	assert.Nil(t, err)
	assert.Equal(t, testIp, ip)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected IP %s, got: %s", testIp, ip))
}

func TestGetClientPublicIP_XRealIp(t *testing.T) {
	testIp := "113.168.80.129"
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-Ip", testIp)
	ip, err := osx.GetClientPublicIP(req)
	assert.Nil(t, err)
	assert.Equal(t, testIp, ip)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected IP %s, got: %s", testIp, ip))
}

func TestGetClientPublicIP_RemoteAddr(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "115.10.11.12:12345"
	spIp := strings.Split(req.RemoteAddr, ":")[0]
	ip, err := osx.GetClientPublicIP(req)
	assert.Nil(t, err)
	assert.Equal(t, spIp, ip)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected IP %s, got: %s", spIp, ip))
}

func TestGetConNetPublicIp(t *testing.T) {
	ip, err := osx.GetConNetPublicIp()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected public IP, got: %s", ip))
}

func TestGetClientPublicIP_NoValidIp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "127.0.0.1") // Localhost IP
	req.Header.Set("X-Real-Ip", "169.254.0.1")     // Link-local IP
	req.RemoteAddr = "192.168.1.1:12345"           // Private IP
	ip, err := osx.GetClientPublicIP(req)
	assert.Nil(t, err)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected public IP, got: %s", ip))
}

// TestGetCallerInfo 测试 GetCallerInfo 函数
func TestGetCallerInfo(t *testing.T) {
	caller := osx.GetCallerInfo(0)
	assert.Equal(t, caller.FuncName, "GetCallerInfo")

	caller = osx.GetCallerInfo(1)
	assert.Equal(t, caller.FuncName, "TestGetCallerInfo")
	assert.Equal(t, caller.Line, 116)

	caller = osx.GetCallerInfo(2)
	assert.Equal(t, caller.FuncName, "tRunner")
}
