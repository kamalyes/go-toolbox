/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 10:50:50
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

func TestMustInt(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int
	}{
		{"123", 123},
		{123, 123},
		{nil, 0},
		{true, 1},
		{false, 0},
		{3.14, 3},
	}

	for _, test := range tests {
		result, _ := convert.MustInt(test.input)
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