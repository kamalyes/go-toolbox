/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-22 09:26:11
 * @FilePath: \go-toolbox\tests\mathx_ternary_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/stretchr/testify/assert"
)

// 自定义类型
type IFType struct {
	Value string
}

// 测试 IF 函数的不同类型
func TestIF(t *testing.T) {
	tests := map[string]struct {
		condition bool
		trueVal   interface{}
		falseVal  interface{}
		expected  interface{}
	}{
		"condition true":    {60 > 50, "Hello", "World", "Hello"},
		"condition false":   {15 > 60, "Hello", "World", "World"},
		"string true":       {true, "Hello", "World", "Hello"},
		"string false":      {false, "Hello", "World", "World"},
		"int true":          {true, 10, 20, 10},
		"int false":         {false, 10, 20, 20},
		"bool true":         {true, true, false, true},
		"bool false":        {false, true, false, false},
		"float true":        {true, 3.14, 2.71, 3.14},
		"float false":       {false, 3.14, 2.71, 2.71},
		"custom type true":  {true, IFType{Value: "Hello"}, IFType{Value: "World"}, IFType{Value: "Hello"}},
		"custom type false": {false, IFType{Value: "Hello"}, IFType{Value: "World"}, IFType{Value: "World"}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, mathx.IF(tt.condition, tt.trueVal, tt.falseVal))
		})
	}
}

// 测试 IfDo 函数的不同类型
func TestIfDo(t *testing.T) {
	tests := map[string]struct {
		condition  bool
		do         func() interface{}
		defaultVal interface{}
		expected   interface{}
	}{
		"string true":       {true, func() interface{} { return "Hello" }, "World", "Hello"},
		"string false":      {false, func() interface{} { return "Hello" }, "World", "World"},
		"int true":          {true, func() interface{} { return 100 }, 0, 100},
		"int false":         {false, func() interface{} { return 100 }, 0, 0},
		"bool true":         {true, func() interface{} { return true }, false, true},
		"bool false":        {false, func() interface{} { return true }, false, false},
		"float true":        {true, func() interface{} { return 3.14 }, 2.71, 3.14},
		"float false":       {false, func() interface{} { return 3.14 }, 2.71, 2.71},
		"custom type true":  {true, func() interface{} { return IFType{Value: "Hello"} }, IFType{Value: "World"}, IFType{Value: "Hello"}},
		"custom type false": {false, func() interface{} { return IFType{Value: "Hello"} }, IFType{Value: "World"}, IFType{Value: "World"}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, mathx.IfDo(tt.condition, tt.do, tt.defaultVal))
		})
	}
}

// 测试 IfDoAF 函数的不同类型
func TestIfDoAF(t *testing.T) {
	tests := map[string]struct {
		condition   bool
		do          func() interface{}
		defaultFunc func() interface{}
		expected    interface{}
	}{
		"string true":       {true, func() interface{} { return "Hello" }, func() interface{} { return "World" }, "Hello"},
		"string false":      {false, func() interface{} { return "Hello" }, func() interface{} { return "World" }, "World"},
		"int true":          {true, func() interface{} { return 100 }, func() interface{} { return 0 }, 100},
		"int false":         {false, func() interface{} { return 100 }, func() interface{} { return 0 }, 0},
		"bool true":         {true, func() interface{} { return true }, func() interface{} { return false }, true},
		"bool false":        {false, func() interface{} { return true }, func() interface{} { return false }, false},
		"float true":        {true, func() interface{} { return 3.14 }, func() interface{} { return 2.71 }, 3.14},
		"float false":       {false, func() interface{} { return 3.14 }, func() interface{} { return 2.71 }, 2.71},
		"custom type true":  {true, func() interface{} { return IFType{Value: "Hello"} }, func() interface{} { return IFType{Value: "World"} }, IFType{Value: "Hello"}},
		"custom type false": {false, func() interface{} { return IFType{Value: "Hello"} }, func() interface{} { return IFType{Value: "World"} }, IFType{Value: "World"}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, mathx.IfDoAF(tt.condition, tt.do, tt.defaultFunc))
		})
	}
}
