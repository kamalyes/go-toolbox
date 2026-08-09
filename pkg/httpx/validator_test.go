/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-10 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-10 00:00:00
 * @FilePath: \go-toolbox\pkg\httpx\validator_test.go
 * @Description: HTTP 方法校验测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package httpx

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsValidMethod 测试 HTTP 方法校验函数
func TestIsValidMethod(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   bool
	}{
		// 标准方法
		{"GET", http.MethodGet, true},
		{"POST", http.MethodPost, true},
		{"PUT", http.MethodPut, true},
		{"DELETE", http.MethodDelete, true},
		{"PATCH", http.MethodPatch, true},
		{"HEAD", http.MethodHead, true},
		{"OPTIONS", http.MethodOptions, true},
		{"CONNECT", http.MethodConnect, true},
		{"TRACE", http.MethodTrace, true},

		// 大小写不敏感
		{"小写 get", "get", true},
		{"小写 post", "post", true},
		{"混合大小写 Get", "Get", true},
		{"混合大小写 PoSt", "PoSt", true},

		// 无效方法
		{"空字符串", "", false},
		{"未知方法", "INVALID", false},
		{"未知方法小写", "invalid", false},
		{"FOO", "FOO", false},
		{"包含空格", "GET ", false},
		{"包含数字", "GET1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidMethod(tt.method)
			assert.Equal(t, tt.want, got, "IsValidMethod(%q) 返回值不符合预期", tt.method)
		})
	}
}

// BenchmarkIsValidMethod 性能基准测试
func BenchmarkIsValidMethod(b *testing.B) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		"INVALID",
		"get",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, m := range methods {
			IsValidMethod(m)
		}
	}
}
