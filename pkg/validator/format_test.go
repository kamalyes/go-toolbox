/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-25 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\format_test.go
 * @Description: 格式验证函数测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"
)

func TestValidateString(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		expect   string
		op       CompareOperator
		wantPass bool
	}{
		{"Equal - Pass", "hello", "hello", OpEqual, true},
		{"Equal - Fail", "hello", "world", OpEqual, false},
		{"NotEqual - Pass", "hello", "world", OpNotEqual, true},
		{"Contains - Pass", "hello world", "world", OpContains, true},
		{"Contains - Fail", "hello", "world", OpContains, false},
		{"HasPrefix - Pass", "hello world", "hello", OpHasPrefix, true},
		{"HasSuffix - Pass", "hello world", "world", OpHasSuffix, true},
		{"Empty - Pass", "   ", "", OpEmpty, true},
		{"Empty - Fail", "hello", "", OpEmpty, false},
		{"NotEmpty - Pass", "hello", "", OpNotEmpty, true},
		{"Regex - Pass", "hello123", "^hello[0-9]+$", OpRegex, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateString(tt.actual, tt.expect, tt.op)
			if result.Success != tt.wantPass {
				t.Errorf("ValidateString() = %v, want %v, message: %s", result.Success, tt.wantPass, result.Message)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		wantPass bool
	}{
		{"Valid email", "user@example.com", true},
		{"Valid email with name", "Test User <user@example.com>", true},
		{"Invalid - no @", "userexample.com", false},
		{"Invalid - no domain", "user@", false},
		{"Invalid - no TLD", "user@example", false},
		{"Empty", "", false},
		{"Valid - subdomain", "user@mail.example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEmail(tt.email)
			if result.Success != tt.wantPass {
				t.Errorf("ValidateEmail(%s) = %v, want %v, message: %s", tt.email, result.Success, tt.wantPass, result.Message)
			}
		})
	}
}

func TestValidateIPAddress(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		wantPass bool
	}{
		{"Valid IPv4", "192.168.1.1", true},
		{"Valid IPv4 - localhost", "127.0.0.1", true},
		{"Valid IPv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"Valid IPv6 - short", "2001:db8::1", true},
		{"Valid IPv6 - localhost", "::1", true},
		{"Invalid - malformed", "192.168.1", false},
		{"Invalid - out of range", "256.256.256.256", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateIP(tt.ip)
			if result.Success != tt.wantPass {
				t.Errorf("ValidateIP(%s) = %v, want %v, message: %s", tt.ip, result.Success, tt.wantPass, result.Message)
			}
		})
	}
}

func TestValidateProtocol(t *testing.T) {
	tests := []struct {
		name              string
		url               string
		allowedProtocols  []string
		wantPass          bool
	}{
		{"HTTP - default", "http://example.com", nil, true},
		{"HTTPS - default", "https://example.com", nil, true},
		{"WS - default", "ws://example.com", nil, true},
		{"WSS - default", "wss://example.com", nil, true},
		{"FTP - default", "ftp://example.com", nil, true},
		{"FTPS - default", "ftps://example.com", nil, true},
		{"HTTP only - pass", "http://example.com", []string{"http"}, true},
		{"HTTP only - fail", "https://example.com", []string{"http"}, false},
		{"WS/WSS - pass", "wss://example.com/ws", []string{"ws", "wss"}, true},
		{"Invalid - no scheme", "example.com", nil, false},
		{"Invalid - no host", "http://", nil, false},
		{"Empty", "", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateProtocol(tt.url, tt.allowedProtocols...)
			if result.Success != tt.wantPass {
				t.Errorf("ValidateProtocol(%s, %v) = %v, want %v, message: %s", tt.url, tt.allowedProtocols, result.Success, tt.wantPass, result.Message)
			}
		})
	}
}

func TestValidateHTTP(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantPass bool
	}{
		{"Valid HTTP", "http://example.com", true},
		{"Valid HTTPS", "https://example.com", true},
		{"Valid with path", "https://example.com/path/to/resource", true},
		{"Valid with query", "https://example.com?key=value", true},
		{"Invalid - WS protocol", "ws://example.com", false},
		{"Invalid - FTP protocol", "ftp://example.com", false},
		{"Invalid - no scheme", "example.com", false},
		{"Invalid - no host", "http://", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateHTTP(tt.url)
			if result.Success != tt.wantPass {
				t.Errorf("ValidateHTTP(%s) = %v, want %v, message: %s", tt.url, result.Success, tt.wantPass, result.Message)
			}
		})
	}
}

func TestValidateWebSocket(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantPass bool
	}{
		{"Valid WS", "ws://example.com", true},
		{"Valid WSS", "wss://example.com", true},
		{"Valid with path", "wss://example.com/ws/chat", true},
		{"Invalid - HTTP protocol", "http://example.com", false},
		{"Invalid - no scheme", "example.com", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateWebSocket(tt.url)
			if result.Success != tt.wantPass {
				t.Errorf("ValidateWebSocket(%s) = %v, want %v, message: %s", tt.url, result.Success, tt.wantPass, result.Message)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name     string
		uuid     string
		wantPass bool
	}{
		{"Valid UUID v4", "550e8400-e29b-41d4-a716-446655440000", true},
		{"Valid UUID v1", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", true},
		{"Invalid - no hyphens", "550e8400e29b41d4a716446655440000", false},
		{"Invalid - wrong format", "550e8400-e29b-41d4-a716", false},
		{"Invalid - wrong chars", "550e8400-e29b-41d4-a716-44665544000g", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUUID(tt.uuid)
			if result.Success != tt.wantPass {
				t.Errorf("ValidateUUID(%s) = %v, want %v, message: %s", tt.uuid, result.Success, tt.wantPass, result.Message)
			}
		})
	}
}

func TestValidateBase64(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		wantPass bool
	}{
		{"Valid Standard Base64", "SGVsbG8gV29ybGQ=", true},
		{"Valid URL-safe Base64", "SGVsbG8gV29ybGQ", true},
		{"Valid without padding", "SGVsbG8gV29ybGQ", true},
		{"Invalid - wrong chars", "Hello@World!", false},
		{"Empty", "", false},
		{"Valid complex", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateBase64(tt.str)
			if result.Success != tt.wantPass {
				t.Errorf("ValidateBase64(%s) = %v, want %v, message: %s", tt.str, result.Success, tt.wantPass, result.Message)
			}
		})
	}
}

// 基准测试
func BenchmarkValidateEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateEmail("user@example.com")
	}
}

func BenchmarkValidateIPAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateIP("192.168.1.1")
	}
}

func BenchmarkValidateProtocol(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateProtocol("https://example.com/path")
	}
}

func BenchmarkValidateHTTP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateHTTP("https://example.com/path")
	}
}

func BenchmarkValidateWebSocket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateWebSocket("wss://example.com/ws")
	}
}

func BenchmarkValidateUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateUUID("550e8400-e29b-41d4-a716-446655440000")
	}
}

func BenchmarkValidateBase64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateBase64("SGVsbG8gV29ybGQ=")
	}
}
