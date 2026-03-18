/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-18 16:37:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-18 16:37:15
 * @FilePath: \go-toolbox\pkg\validator\path_test.go
 * @Description: 路径匹配功能测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchPathInList(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{"精确匹配", "/health", []string{"/health", "/metrics"}, true},
		{"前缀匹配", "/api/v1/users", []string{"/api/v1", "/admin"}, true},
		{"不匹配", "/admin/users", []string{"/api", "/health"}, false},
		{"空列表", "/api/users", []string{}, false},
		{"空路径", "", []string{"/api", ""}, true},
		{"多个模式匹配第一个", "/health/check", []string{"/health", "/api"}, true},
		{"多个模式匹配最后一个", "/api/v2/users", []string{"/admin", "/api"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchPathInList(tt.path, tt.patterns)
			assert.Equal(t, tt.want, got)
		})
	}
}
