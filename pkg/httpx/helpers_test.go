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
