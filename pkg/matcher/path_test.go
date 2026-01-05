/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-05 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-05 00:00:00
 * @FilePath: \go-toolbox\pkg\matcher\path_test.go
 * @Description: 路径匹配器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathMatcher(t *testing.T) {
	tests := []struct {
		name      string
		matchType PathMatcherType
		pattern   string
		path      string
		expected  bool
	}{
		{"ExactMatch", PathMatchExact, "/api/v1/resource", "/api/v1/resource", true},
		{"ExactMatchFail", PathMatchExact, "/api/v1/resource", "/api/v1/other", false},
		{"PrefixMatch", PathMatchPrefix, "/api", "/api/v1/resource", true},
		{"PrefixMatchFail", PathMatchPrefix, "/api", "/v1/resource", false},
		{"SuffixMatch", PathMatchSuffix, "resource", "/api/v1/resource", true},
		{"SuffixMatchFail", PathMatchSuffix, "resource", "/api/v1/", false},
		{"GlobMatch", PathMatchGlob, "/api/*", "/api/v1/resource", true},
		{"GlobMatchFail", PathMatchGlob, "/api/*", "/v1/resource", false},
		{"RegexMatch", PathMatchRegex, `^/api/\d+/resource$`, "/api/123/resource", true},
		{"RegexMatchFail", PathMatchRegex, `^/api/\d+/resource$`, "/api/resource", false},
		{"ContainsMatch", PathMatchContains, "v1", "/api/v1/resource", true},
		{"ContainsMatchFail", PathMatchContains, "v1", "/api/resource", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm, err := NewPathMatcher(tt.matchType, tt.pattern)
			assert.NoError(t, err, "创建路径匹配器失败")

			result := pm.Match(tt.path)
			assert.Equal(t, tt.expected, result, "路径匹配结果不符合预期")
		})
	}
}

func TestNormalizePath(t *testing.T) {
	assert.Equal(t, "/api/v1/resource", NormalizePath("//api/v1/resource"), "路径标准化失败")
	assert.Equal(t, "/api/v1/resource", NormalizePath("/api/v1/resource/"), "路径标准化失败")
	assert.Equal(t, "/", NormalizePath("/"), "路径标准化失败")
	assert.Equal(t, "/", NormalizePath(""), "路径标准化失败")
}

func TestExtractPathSegments(t *testing.T) {
	assert.Equal(t, []string{"api", "v1", "resource"}, ExtractPathSegments("/api/v1/resource"), "路径段提取失败")
	assert.Equal(t, []string{}, ExtractPathSegments("/"), "路径段提取失败")
	assert.Equal(t, []string{}, ExtractPathSegments(""), "路径段提取失败")
}

func TestMatchPathWithMethod(t *testing.T) {
	allowedMethods := []string{"GET", "POST"}
	assert.True(t, MatchPathWithMethod("/api/v1/resource", "GET", "/api/v1/resource", allowedMethods), "应该允许GET方法")
	assert.False(t, MatchPathWithMethod("/api/v1/resource", "DELETE", "/api/v1/resource", allowedMethods), "不应该允许DELETE方法")
}

// TestMatchMethod 测试 HTTP 方法匹配
func TestMatchMethod(t *testing.T) {
	tests := []struct {
		name     string
		methods  []string
		method   string
		expected bool
	}{
		{
			name:     "空列表匹配所有",
			methods:  []string{},
			method:   "GET",
			expected: true,
		},
		{
			name:     "nil列表匹配所有",
			methods:  nil,
			method:   "POST",
			expected: true,
		},
		{
			name:     "精确匹配",
			methods:  []string{"GET", "POST"},
			method:   "GET",
			expected: true,
		},
		{
			name:     "大小写不敏感",
			methods:  []string{"GET", "POST"},
			method:   "get",
			expected: true,
		},
		{
			name:     "不匹配",
			methods:  []string{"GET", "POST"},
			method:   "DELETE",
			expected: false,
		},
		{
			name:     "单个方法匹配",
			methods:  []string{"POST"},
			method:   "POST",
			expected: true,
		},
		{
			name:     "多个方法不匹配",
			methods:  []string{"GET", "POST", "PUT"},
			method:   "PATCH",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchMethod(tt.methods, tt.method)
			assert.Equal(t, tt.expected, result)
		})
	}
}
