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
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestValue(t *testing.T) {
	tests := []struct {
		name       string
		contextKey interface{}
		header     string
		query      string
		expected   string
	}{
		{"Context value", "testKey", "", "", "contextValue"},
		{"Header value", nil, "headerValue", "", "headerValue"},
		{"Query value", nil, "", "queryValue", "queryValue"},
		{"No value", nil, "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: http.Header{},
				URL:    &url.URL{},
			}
			// 设置查询参数
			if tt.query != "" {
				query := url.Values{}
				query.Set("queryName", tt.query)
				req.URL.RawQuery = query.Encode()
			}
			if tt.contextKey != nil {
				req = req.WithContext(context.WithValue(req.Context(), tt.contextKey, "contextValue"))
			}
			req.Header.Set("Header-Name", tt.header)

			assert.Equal(t, tt.expected, GetRequestValue(req, tt.contextKey, "Header-Name", "queryName"))
		})
	}
}

func TestReadRequestBody(t *testing.T) {
	tests := []struct {
		name        string
		body        []byte
		expected    []byte
		expectError bool
	}{
		{"Read non-empty", []byte("Hello"), []byte("Hello"), false},
		{"Read empty", []byte(""), []byte(""), false},
		{"Nil body", nil, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rBody := io.NopCloser(bytes.NewBuffer(tt.body))
			if tt.body == nil {
				rBody = nil
			}

			r := &http.Request{Body: rBody}
			bodyBytes, err := ReadRequestBody(r)

			assert.Equal(t, tt.expectError, err != nil)
			assert.Equal(t, tt.expected, bodyBytes)
		})
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
