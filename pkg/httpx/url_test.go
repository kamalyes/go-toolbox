/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-16 21:50:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-16 21:50:00
 * @FilePath: \go-toolbox\pkg\httpx\url_test.go
 * @Description: URL 处理工具函数测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "域名不带协议",
			input:    "www.example.com",
			expected: "https://www.example.com",
		},
		{
			name:     "域名不带 www 和协议",
			input:    "example.com",
			expected: "https://example.com",
		},
		{
			name:     "域名带路径不带协议",
			input:    "www.example.com/api",
			expected: "https://www.example.com/api",
		},
		{
			name:     "域名带路径和尾部斜杠不带协议",
			input:    "www.example.com/api/",
			expected: "https://www.example.com/api/",
		},
		{
			name:     "已包含 http 协议",
			input:    "http://www.example.com",
			expected: "http://www.example.com",
		},
		{
			name:     "已包含 https 协议",
			input:    "https://www.example.com",
			expected: "https://www.example.com",
		},
		{
			name:     "已包含 HTTP 协议(大写)",
			input:    "HTTP://www.example.com",
			expected: "HTTP://www.example.com",
		},
		{
			name:     "已包含 HTTPS 协议(大写)",
			input:    "HTTPS://www.example.com",
			expected: "HTTPS://www.example.com",
		},
		{
			name:     "混合大小写协议 Http",
			input:    "Http://www.example.com",
			expected: "Http://www.example.com",
		},
		{
			name:     "混合大小写协议 Https",
			input:    "Https://www.example.com",
			expected: "Https://www.example.com",
		},
		{
			name:     "IP 地址不带协议",
			input:    "192.168.1.1",
			expected: "https://192.168.1.1",
		},
		{
			name:     "IP 地址带端口不带协议",
			input:    "192.168.1.1:8080",
			expected: "https://192.168.1.1:8080",
		},
		{
			name:     "localhost 不带协议",
			input:    "localhost:3000",
			expected: "https://localhost:3000",
		},
		{
			name:     "已包含协议和端口",
			input:    "http://localhost:3000",
			expected: "http://localhost:3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeBaseURL(tt.input)
			assert.Equal(t, tt.expected, result, "NormalizeBaseURL(%q) 返回值不符合预期", tt.input)
		})
	}
}

// BenchmarkNormalizeBaseURL 性能基准测试
func BenchmarkNormalizeBaseURL(t *testing.B) {
	testCases := []string{
		"www.example.com",
		"https://www.example.com",
		"example.com/api/v1",
		"",
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		for _, tc := range testCases {
			NormalizeBaseURL(tc)
		}
	}
}
