/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-26 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-26 13:25:15
 * @FilePath: \go-toolbox\pkg\stringx\domain_test.go
 * @Description: 域名处理工具函数测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDomainPrefix(t *testing.T) {
	tests := []struct {
		name          string
		fullDomain    string
		primaryDomain string
		want          string
	}{
		{"单级子域名", "www.example.com", "example.com", "www"},
		{"多级子域名", "api.staging.example.com", "example.com", "api.staging"},
		{"根域名本身", "example.com", "example.com", RootDomainPrefix},
		{"空完整域名", "", "example.com", RootDomainPrefix},
		{"空主域名", "www.example.com", "", "www.example.com"},
		{"不匹配的域名", "www.other.com", "example.com", "www.other.com"},
		{"主域名比完整域名长", "ex.com", "example.com", "ex.com"},
		{"多级复杂子域名", "h5.abc.ca.example.com", "example.com", "h5.abc.ca"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractDomainPrefix(tt.fullDomain, tt.primaryDomain)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsSubdomain(t *testing.T) {
	tests := []struct {
		name          string
		subdomain     string
		primaryDomain string
		want          bool
	}{
		{"有效子域名", "www.example.com", "example.com", true},
		{"多级子域名", "api.www.example.com", "example.com", true},
		{"相同域名", "example.com", "example.com", false},
		{"不相关域名", "www.other.com", "example.com", false},
		{"空参数", "", "example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSubdomain(tt.subdomain, tt.primaryDomain)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSplitDomain(t *testing.T) {
	tests := []struct {
		name          string
		fullDomain    string
		primaryDomain string
		wantPrefix    string
		wantPrimary   string
	}{
		{"子域名", "www.example.com", "example.com", "www", "example.com"},
		{"多级子域名", "api.staging.example.com", "example.com", "api.staging", "example.com"},
		{"根域名", "example.com", "example.com", "", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrefix, gotPrimary := SplitDomain(tt.fullDomain, tt.primaryDomain)
			assert.Equal(t, tt.wantPrefix, gotPrefix)
			assert.Equal(t, tt.wantPrimary, gotPrimary)
		})
	}
}

// Benchmark
func BenchmarkExtractDomainPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractDomainPrefix("api.staging.example.com", "example.com")
	}
}
