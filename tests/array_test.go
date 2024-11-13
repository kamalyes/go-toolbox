/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 15:55:06
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 13:08:58
 * @FilePath: \go-toolbox\tests\array_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/array"
)

// 测试用例
func TestMinMax(t *testing.T) {
	minFunc := func(a, b interface{}) interface{} {
		if a.(int) < b.(int) {
			return a
		}
		return b
	}

	maxFunc := func(a, b interface{}) interface{} {
		if a.(int) > b.(int) {
			return a
		}
		return b
	}

	tests := []struct {
		name        string
		list        []interface{}
		f           array.MinMaxFunc
		expected    interface{}
		expectError bool
	}{
		{"Min with positive integers", []interface{}{3, 1, 4, 2}, minFunc, 1, false},
		{"Max with positive integers", []interface{}{3, 1, 4, 2}, maxFunc, 4, false},
		{"Min with negative integers", []interface{}{-1, -3, -2}, minFunc, -3, false},
		{"Max with negative integers", []interface{}{-1, -3, -2}, maxFunc, -1, false},
		{"Empty list", []interface{}{}, minFunc, nil, true},
		{"Single element", []interface{}{5}, minFunc, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := array.MinMax(tt.list, tt.f)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError && got != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}
