/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-23 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\regex_test.go
 * @Description: 正则表达式验证测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRegex(t *testing.T) {
	a := assert.New(t)
	
	tests := []struct {
		name     string
		body     []byte
		pattern  string
		wantPass bool
	}{
		{"匹配成功", []byte("abc123"), "^[a-z]+[0-9]+$", true},
		{"匹配失败", []byte("123abc"), "^[a-z]+[0-9]+$", false},
		{"空模式", []byte("abc123"), "", true},
		{"无效正则", []byte("abc123"), "[", false},
		{"邮箱匹配", []byte("test@example.com"), `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, true},
		{"URL匹配", []byte("https://example.com"), `^https?://`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateRegex(tt.body, tt.pattern)
			a.Equal(tt.wantPass, result.Success, "message: %s", result.Message)
		})
	}
}
