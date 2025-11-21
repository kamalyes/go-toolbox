/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 16:10:53
 * @FilePath: \go-toolbox\tests\netx_ip_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/netx"
	"github.com/stretchr/testify/assert"
)

func TestGetLocalInterfaceIPAndExternalIP(t *testing.T) {
	externalIP, internalIP, err := netx.GetLocalInterfaceIPAndExternalIP()
	assert.Nil(t, err)
	assert.NotEmpty(t, externalIP)
	assert.NotEmpty(t, internalIP)
	t.Logf("externalIP %s", externalIP)
	t.Logf("internalIP %s", internalIP)
}

func TestGetLocalInterfaceIPs(t *testing.T) {
	ips, err := netx.GetLocalInterfaceIPs()
	assert.Nil(t, err)
	assert.NotEmpty(t, ips, fmt.Sprintf("Expected at least one global unicast IP, got: %v", ips))
	for _, ip := range ips {
		assert.NotEmpty(t, ip, fmt.Sprintf("Invalid IP address: %s", ip))
	}
}

func TestGetConNetPublicIP(t *testing.T) {
	ip, err := netx.GetConNetPublicIP()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip, fmt.Sprintf("Expected public IP, got: %s", ip))
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		expectedIP string
	}{
		{
			name: "X-Forwarded-For header",
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.1, 10.0.0.1",
			},
			remoteAddr: "127.0.0.1:8080",
			expectedIP: "192.168.1.1",
		},
		{
			name: "X-Real-IP header",
			headers: map[string]string{
				"X-Real-IP": "203.0.113.1",
			},
			remoteAddr: "127.0.0.1:8080",
			expectedIP: "203.0.113.1",
		},
		{
			name:       "RemoteAddr fallback",
			headers:    map[string]string{},
			remoteAddr: "192.0.2.1:8080",
			expectedIP: "192.0.2.1",
		},
		{
			name:       "No IP headers",
			headers:    map[string]string{},
			remoteAddr: "127.0.0.1:8080",
			expectedIP: "127.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			req.RemoteAddr = tt.remoteAddr

			ip := netx.GetClientIP(req)
			assert.Equal(t, tt.expectedIP, ip)
		})
	}
}
