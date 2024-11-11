/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-12 22:29:55
 * @FilePath: \go-toolbox\tests\convert_must_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
)

func TestMustString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{"hello", "hello"},
		{[]byte("world"), "world"},
		{nil, ""},
		{true, "true"},
		{42, "42"},
		{3.14, "3.14"},
		{time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), "2024-01-01T12:00:00Z"},
	}

	for _, test := range tests {
		result := convert.MustString(test.input)
		if result != test.expected {
			t.Errorf("convert.MustString(%v) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestMustIntT_ErrorCases(t *testing.T) {
	// 定义一些不支持的类型进行测试
	unsupportedValues := []any{
		"string",         // 字符串
		3.14,             // 浮点数
		true,             // 布尔值
		[]int{1, 2, 3},   // 切片
		map[string]int{}, // 映射
		nil,              // nil 值
	}

	for _, val := range unsupportedValues {
		result, err := convert.MustIntT[int](val)
		if err == nil {
			t.Errorf("Expected an error for input %v, but got result %d", val, result)
		}
	}
}

func TestMustInt(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int
	}{
		{int(10), 10},
		{int8(20), 20},
		{int16(30), 30},
		{int32(40), 40},
		{int64(50), 50},
		{uint(60), 60},
		{uint8(70), 70},
		{uint16(80), 80},
		{uint32(90), 90},
		{uint64(100), 100},
	}

	for _, test := range tests {
		result, _ := convert.MustIntT[int](test.input)
		if result != test.expected {
			t.Errorf("MustInt(%v) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestMustBool(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{"1", true},
		{"true", true},
		{"false", false},
		{0, false},
		{1, true},
		{nil, false},
		{true, true},
		{false, false},
	}

	for _, test := range tests {
		result := convert.MustBool(test.input)
		if result != test.expected {
			t.Errorf("MustBool(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}
