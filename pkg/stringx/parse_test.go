/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-06 17:51:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-06 17:51:22
 * @FilePath: \go-toolbox\pkg\stringx\parse_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFieldIntOrWildcard(t *testing.T) {
	wildcard := "*"

	tests := []struct {
		name      string
		field     string
		wildcard  string
		min       int
		max       int
		wantVal   int
		wantError bool
	}{
		{"通配符*", "*", wildcard, 0, 59, -1, false},
		{"正常数字30", "30", wildcard, 0, 59, 30, false},
		{"边界值0", "0", wildcard, 0, 59, 0, false},
		{"边界值59", "59", wildcard, 0, 59, 59, false},
		{"超出范围60", "60", wildcard, 0, 59, 0, true},
		{"无效数字abc", "abc", wildcard, 0, 59, 0, true},
		{"其他通配符?", "?", "?", 0, 59, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ParseFieldIntOrWildcard(tt.field, tt.wildcard, tt.min, tt.max)
			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, tt.wantVal, val)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVal, val)
			}
		})
	}
}

func TestParseFieldInt(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		min       int
		max       int
		wantVal   int
		wantError bool
	}{
		{"正常数字15", "15", 0, 59, 15, false},
		{"边界值0", "0", 0, 59, 0, false},
		{"边界值59", "59", 0, 59, 59, false},
		{"超出范围60", "60", 0, 59, 0, true},
		{"无效数字abc", "abc", 0, 59, 0, true},
		{"通配符*", "*", 0, 59, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ParseFieldInt(tt.field, tt.min, tt.max)
			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, tt.wantVal, val)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVal, val)
			}
		})
	}
}
