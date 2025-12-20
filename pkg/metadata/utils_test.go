/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 17:00:00
 * @FilePath: \go-toolbox\pkg\metadata\utils_test.go
 * @Description: 元数据工具函数测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package metadata

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTLSVersionString(t *testing.T) {
	tests := []struct {
		name     string
		version  uint16
		expected string
	}{
		{"TLS 1.3", tls.VersionTLS13, "TLS 1.3"},
		{"TLS 1.2", tls.VersionTLS12, "TLS 1.2"},
		{"TLS 1.1", tls.VersionTLS11, "TLS 1.1"},
		{"TLS 1.0", tls.VersionTLS10, "TLS 1.0"},
		{"SSL 3.0", 0x0300, "SSL 3.0"},
		{"Unknown", 0x0000, ""},
		{"Invalid", 0xFFFF, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTLSVersionString(tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		name       string
		acceptLang string
		wantLang   string
		wantRegion string
		wantFull   string
	}{
		{
			name:       "Chinese Simplified",
			acceptLang: "zh-CN,zh;q=0.9,en;q=0.8",
			wantLang:   "zh",
			wantRegion: "CN",
			wantFull:   "zh-CN",
		},
		{
			name:       "English US",
			acceptLang: "en-US,en;q=0.9",
			wantLang:   "en",
			wantRegion: "US",
			wantFull:   "en-US",
		},
		{
			name:       "French France",
			acceptLang: "fr-FR",
			wantLang:   "fr",
			wantRegion: "FR",
			wantFull:   "fr-FR",
		},
		{
			name:       "Language Only",
			acceptLang: "en",
			wantLang:   "en",
			wantRegion: "",
			wantFull:   "en",
		},
		{
			name:       "Language Only with Quality",
			acceptLang: "ja;q=0.9",
			wantLang:   "ja",
			wantRegion: "",
			wantFull:   "ja",
		},
		{
			name:       "Empty String",
			acceptLang: "",
			wantLang:   "",
			wantRegion: "",
			wantFull:   "",
		},
		{
			name:       "Whitespace",
			acceptLang: "  ",
			wantLang:   "",
			wantRegion: "",
			wantFull:   "",
		},
		{
			name:       "With Spaces",
			acceptLang: " en-US , zh-CN",
			wantLang:   "en",
			wantRegion: "US",
			wantFull:   "en-US",
		},
		{
			name:       "Complex with Multiple Languages",
			acceptLang: "de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7",
			wantLang:   "de",
			wantRegion: "DE",
			wantFull:   "de-DE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang, region, full := ParseAcceptLanguage(tt.acceptLang)
			assert.Equal(t, tt.wantLang, lang)
			assert.Equal(t, tt.wantRegion, region)
			assert.Equal(t, tt.wantFull, full)
		})
	}
}

func TestGetRemoteIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		expected   string
	}{
		{
			name:       "IPv4 with Port",
			remoteAddr: "192.168.1.1:12345",
			expected:   "192.168.1.1",
		},
		{
			name:       "IPv6 with Port",
			remoteAddr: "[2001:db8::1]:8080",
			expected:   "2001:db8::1",
		},
		{
			name:       "IPv6 Full with Port",
			remoteAddr: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:443",
			expected:   "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		},
		{
			name:       "IPv4 Only",
			remoteAddr: "192.168.1.1",
			expected:   "192.168.1.1",
		},
		{
			name:       "Localhost with Port",
			remoteAddr: "127.0.0.1:8080",
			expected:   "127.0.0.1",
		},
		{
			name:       "IPv6 Localhost with Port",
			remoteAddr: "[::1]:8080",
			expected:   "::1",
		},
		{
			name:       "Empty String",
			remoteAddr: "",
			expected:   "",
		},
		{
			name:       "Only Port",
			remoteAddr: ":8080",
			expected:   ":8080",
		},
		{
			name:       "No Port",
			remoteAddr: "example.com",
			expected:   "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRemoteIP(tt.remoteAddr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRemotePort(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		expected   string
	}{
		{
			name:       "IPv4 with Port",
			remoteAddr: "192.168.1.1:12345",
			expected:   "12345",
		},
		{
			name:       "IPv6 with Port",
			remoteAddr: "[2001:db8::1]:8080",
			expected:   "8080",
		},
		{
			name:       "IPv6 Full with Port",
			remoteAddr: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:443",
			expected:   "443",
		},
		{
			name:       "Localhost with Port",
			remoteAddr: "127.0.0.1:8080",
			expected:   "8080",
		},
		{
			name:       "IPv6 Localhost with Port",
			remoteAddr: "[::1]:8080",
			expected:   "8080",
		},
		{
			name:       "HTTPS Port",
			remoteAddr: "example.com:443",
			expected:   "443",
		},
		{
			name:       "HTTP Port",
			remoteAddr: "example.com:80",
			expected:   "80",
		},
		{
			name:       "No Port",
			remoteAddr: "192.168.1.1",
			expected:   "",
		},
		{
			name:       "IPv6 No Port",
			remoteAddr: "[2001:db8::1]",
			expected:   "",
		},
		{
			name:       "Empty String",
			remoteAddr: "",
			expected:   "",
		},
		{
			name:       "Only Port",
			remoteAddr: ":8080",
			expected:   "",
		},
		{
			name:       "Colon at End",
			remoteAddr: "192.168.1.1:",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRemotePort(tt.remoteAddr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
