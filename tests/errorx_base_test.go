/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-13 11:27:59
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 11:29:10
 * @FilePath: \go-toolbox\tests\errorx_base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"github.com/stretchr/testify/assert"
)

func TestWrapError(t *testing.T) {
	tests := []struct {
		message  string
		err      error
		expected string
	}{
		{"an error occurred", errors.New("original error"), "an error occurred: original error"}, // 普通错误
		{"another error", nil, ""}, // nil 错误
		{"wrapped error", errors.New("something went wrong"), "wrapped error: something went wrong"}, // 包装错误
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			got := errorx.WrapError(tt.message, tt.err)

			if tt.expected == "" {
				assert.Nil(t, got) // 如果预期是 nil，断言返回值为 nil
			} else {
				assert.NotNil(t, got)                  // 断言返回值不为 nil
				assert.EqualError(t, got, tt.expected) // 断言返回的错误信息与预期相等
			}
		})
	}
}
