/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-05 16:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 11:32:25
 * @FilePath: \go-toolbox\pkg\httpx\helpers_test.go
 * @Description: HTTP 辅助方法测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetUserID 测试获取用户ID
func TestGetUserID(t *testing.T) {
	type contextKey string
	const userIDKey contextKey = "user_id"

	tests := []struct {
		name        string
		setupReq    func() *http.Request
		contextKey  interface{}
		headerKey   string
		expectedID  string
		description string
	}{
		{
			name: "从上下文获取",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				ctx := req.Context()
				req = req.WithContext(context.WithValue(ctx, userIDKey, "user123"))
				return req
			},
			contextKey:  userIDKey,
			headerKey:   "X-User-ID",
			expectedID:  "user123",
			description: "应该从上下文优先获取用户ID",
		},
		{
			name: "从请求头获取",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("X-User-ID", "user456")
				return req
			},
			contextKey:  userIDKey,
			headerKey:   "X-User-ID",
			expectedID:  "user456",
			description: "应该从请求头获取用户ID",
		},
		{
			name: "上下文优先于请求头",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				ctx := req.Context()
				req = req.WithContext(context.WithValue(ctx, userIDKey, "context_user"))
				req.Header.Set("X-User-ID", "header_user")
				return req
			},
			contextKey:  userIDKey,
			headerKey:   "X-User-ID",
			expectedID:  "context_user",
			description: "上下文优先级应该高于请求头",
		},
		{
			name: "都不存在返回空",
			setupReq: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			contextKey:  userIDKey,
			headerKey:   "X-User-ID",
			expectedID:  "",
			description: "都不存在时应该返回空字符串",
		},
		{
			name: "上下文值类型错误",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				ctx := req.Context()
				req = req.WithContext(context.WithValue(ctx, userIDKey, 123)) // 不是字符串
				req.Header.Set("X-User-ID", "header_user")
				return req
			},
			contextKey:  userIDKey,
			headerKey:   "X-User-ID",
			expectedID:  "header_user",
			description: "上下文值类型错误时应该降级到请求头",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupReq()
			userID := GetUserID(req, tt.contextKey, tt.headerKey)
			assert.Equal(t, tt.expectedID, userID, tt.description)
		})
	}
}

// TestGetUserID_NilContextKey 测试 nil context key
func TestGetUserID_NilContextKey(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-User-ID", "user789")

	userID := GetUserID(req, nil, "X-User-ID")
	assert.Equal(t, "user789", userID, "nil context key应该从请求头获取")
}

// TestGetUserID_EmptyHeaderKey 测试空 header key
func TestGetUserID_EmptyHeaderKey(t *testing.T) {
	type contextKey string
	const userIDKey contextKey = "user_id"

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := req.Context()
	req = req.WithContext(context.WithValue(ctx, userIDKey, "user999"))

	userID := GetUserID(req, userIDKey, "")
	assert.Equal(t, "user999", userID, "空header key应该只从上下文获取")
}

// BenchmarkGetUserID 性能测试
func BenchmarkGetUserID(b *testing.B) {
	type contextKey string
	const userIDKey contextKey = "user_id"

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := req.Context()
	req = req.WithContext(context.WithValue(ctx, userIDKey, "bench_user"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetUserID(req, userIDKey, "X-User-ID")
	}
}

// TestBuildParams 测试构建请求参数
func TestBuildParams(t *testing.T) {
	tests := []struct {
		name     string
		base     map[string]string
		opts     []func(map[string]string)
		expected map[string]string
	}{
		{
			name: "仅基础参数",
			base: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			opts: nil,
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "基础参数加可选参数",
			base: map[string]string{
				"domain": "example.com",
			},
			opts: []func(map[string]string){
				WithParam(true, "auto_renew", "1"),
				WithParam(false, "private", "1"),
			},
			expected: map[string]string{
				"domain":     "example.com",
				"auto_renew": "1",
			},
		},
		{
			name: "使用WithParamNotEmpty",
			base: map[string]string{
				"domain": "example.com",
			},
			opts: []func(map[string]string){
				WithParamNotEmpty("coupon", "SAVE10"),
				WithParamNotEmpty("empty", ""),
			},
			expected: map[string]string{
				"domain": "example.com",
				"coupon": "SAVE10",
			},
		},
		{
			name: "空基础参数",
			base: map[string]string{},
			opts: []func(map[string]string){
				WithParam(true, "key1", "value1"),
			},
			expected: map[string]string{
				"key1": "value1",
			},
		},
		{
			name: "nil可选参数",
			base: map[string]string{
				"key1": "value1",
			},
			opts: []func(map[string]string){
				nil,
				WithParam(true, "key2", "value2"),
				nil,
			},
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildParams(tt.base, tt.opts...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestWithParam 测试条件添加参数
func TestWithParam(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		key       string
		value     string
		shouldAdd bool
	}{
		{
			name:      "条件为true应该添加",
			condition: true,
			key:       "auto_renew",
			value:     "1",
			shouldAdd: true,
		},
		{
			name:      "条件为false不应该添加",
			condition: false,
			key:       "private",
			value:     "1",
			shouldAdd: false,
		},
		{
			name:      "true且value为空也应该添加",
			condition: true,
			key:       "empty_key",
			value:     "",
			shouldAdd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := make(map[string]string)
			fn := WithParam(tt.condition, tt.key, tt.value)
			fn(params)

			if tt.shouldAdd {
				assert.Contains(t, params, tt.key)
				assert.Equal(t, tt.value, params[tt.key])
			} else {
				assert.NotContains(t, params, tt.key)
			}
		})
	}
}

// TestWithParamNotEmpty 测试非空添加参数
func TestWithParamNotEmpty(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		shouldAdd bool
	}{
		{
			name:      "非空字符串应该添加",
			key:       "coupon",
			value:     "DISCOUNT20",
			shouldAdd: true,
		},
		{
			name:      "空字符串不应该添加",
			key:       "empty_coupon",
			value:     "",
			shouldAdd: false,
		},
		{
			name:      "空格字符串不应该添加",
			key:       "spaces",
			value:     "   ",
			shouldAdd: false,
		},
		{
			name:      "有效字符串应该添加",
			key:       "payment_id",
			value:     "PAY123",
			shouldAdd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := make(map[string]string)
			fn := WithParamNotEmpty(tt.key, tt.value)
			fn(params)

			if tt.shouldAdd {
				assert.Contains(t, params, tt.key)
				assert.Equal(t, tt.value, params[tt.key])
			} else {
				assert.NotContains(t, params, tt.key)
			}
		})
	}
}

// TestBuildParams_Integration 集成测试
func TestBuildParams_Integration(t *testing.T) {
	// 模拟域名注册请求
	params := BuildParams(
		map[string]string{
			"domain": "example.com",
			"years":  "2",
		},
		WithParam(true, "auto_renew", "1"),
		WithParam(false, "private", "1"),
		WithParamNotEmpty("coupon", "SAVE10"),
		WithParamNotEmpty("payment_id", ""),
	)

	expected := map[string]string{
		"domain":     "example.com",
		"years":      "2",
		"auto_renew": "1",
		"coupon":     "SAVE10",
	}

	assert.Equal(t, expected, params)
}

// BenchmarkBuildParams 性能测试
func BenchmarkBuildParams(b *testing.B) {
	base := map[string]string{
		"domain": "example.com",
		"years":  "1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildParams(
			base,
			WithParam(true, "auto_renew", "1"),
			WithParam(true, "private", "1"),
			WithParamNotEmpty("coupon", "SAVE10"),
		)
	}
}

// BenchmarkWithParam 性能测试
func BenchmarkWithParam(b *testing.B) {
	params := make(map[string]string)
	fn := WithParam(true, "key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(params)
	}
}

// BenchmarkWithParamNotEmpty 性能测试
func BenchmarkWithParamNotEmpty(b *testing.B) {
	params := make(map[string]string)
	fn := WithParamNotEmpty("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(params)
	}
}
