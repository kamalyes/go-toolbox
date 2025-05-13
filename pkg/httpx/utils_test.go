/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-28 18:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-30 20:21:57
 * @FilePath: \go-toolbox\pkg\httpx\utils_test.go
 * @Description: HTTP 相关工具测试用例
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUrlPath(t *testing.T) {
	tests := []struct {
		urlString string
		expected  string
	}{
		{"http://example.com/path/to/resource", "/path/to/resource"},
		{"https://example.com/another/path?query=param", "/another/path"},
		{"ftp://example.com/file.txt", "/file.txt"},
		{"http://example.com/", "/"},
		{"invalid-url", "invalid-url"},
	}

	for _, test := range tests {
		t.Run(test.urlString, func(t *testing.T) {
			result := ParseUrlPath(test.urlString)
			assert.Equal(t, test.expected, result)
		})
	}
}
