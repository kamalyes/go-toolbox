/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-11 16:26:56
 * @FilePath: \go-toolbox\tests\mathx_ternary_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
	"testing"
	"time"

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

// 测试 IfDoWithError 函数
func TestIfDoWithError(t *testing.T) {
	doFuncSuccess := func() (string, error) {
		return "ok", nil
	}
	doFuncFail := func() (string, error) {
		return "", errors.New("fail")
	}

	val, err := mathx.IfDoWithError(true, doFuncSuccess, "default")
	assert.NoError(t, err)
	assert.Equal(t, "ok", val)

	val, err = mathx.IfDoWithError(true, doFuncFail, "default")
	assert.Error(t, err)
	assert.Equal(t, "", val)

	val, err = mathx.IfDoWithError(false, doFuncFail, "default")
	assert.NoError(t, err)
	assert.Equal(t, "default", val)
}

// 测试 IfDoAsync 函数
func TestIfDoAsync(t *testing.T) {
	doFunc := func() int {
		time.Sleep(10 * time.Millisecond)
		return 42
	}

	ch := mathx.IfDoAsync(true, doFunc, 0)
	val := <-ch
	assert.Equal(t, 42, val)

	ch = mathx.IfDoAsync(false, doFunc, 99)
	val = <-ch
	assert.Equal(t, 99, val)
}

// TestIfDoAsyncWithTimeout 测试异步执行带超时的 IfDoAsyncWithTimeout 函数
func TestIfDoAsyncWithTimeout(t *testing.T) {
	assert := assert.New(t)

	slowFunc := func() int {
		time.Sleep(50 * time.Millisecond)
		return 42
	}

	// 执行时间小于超时，返回正常结果
	ch1 := mathx.IfDoAsyncWithTimeout(true, slowFunc, 0, 100)
	res1 := <-ch1
	assert.Equal(42, res1, "未超时应返回正常结果")

	// 执行时间大于超时，返回零值
	ch2 := mathx.IfDoAsyncWithTimeout(true, slowFunc, 0, 10)
	res2 := <-ch2
	assert.Equal(0, res2, "超时应返回类型零值")

	// 条件为 false，直接返回默认值
	ch3 := mathx.IfDoAsyncWithTimeout(false, slowFunc, 99, 100)
	res3 := <-ch3
	assert.Equal(99, res3, "条件为 false 应返回默认值")
}

// TestIfElseAndIfChain 测试多条件链判断 IfElse 和 IfChain
func TestIfElseAndIfChain(t *testing.T) {
	assert := assert.New(t)

	conds := []bool{false, true, false}
	values := []string{"a", "b", "c"}
	defVal := "default"

	res := mathx.IfElse(conds, values, defVal)
	assert.Equal("b", res, "IfElse 应返回第一个为 true 的对应值")

	pairs := []mathx.ConditionValue[int]{
		{Cond: false, Value: 1},
		{Cond: true, Value: 2},
		{Cond: false, Value: 3},
	}
	res2 := mathx.IfChain(pairs, 999)
	assert.Equal(2, res2, "IfChain 应返回第一个为 true 的对应值")
}

// TestIfDoWithErrorAsync 测试异步带错误返回的函数
func TestIfDoWithErrorAsync(t *testing.T) {
	assert := assert.New(t)

	// 模拟成功函数
	successFunc := func() (int, error) {
		return 100, nil
	}

	// 模拟失败函数
	failFunc := func() (int, error) {
		return 0, errors.New("fail error")
	}

	// 条件为 true，成功执行
	ch1 := mathx.IfDoWithErrorAsync(true, successFunc, 999)
	res1 := <-ch1
	assert.NoError(res1.Err, "成功执行时错误应为 nil")
	assert.Equal(100, res1.Result, "应返回成功结果")

	// 条件为 true，执行失败
	ch2 := mathx.IfDoWithErrorAsync(true, failFunc, 999)
	res2 := <-ch2
	assert.Error(res2.Err, "执行失败应返回错误")
	assert.Equal(0, res2.Result, "失败时结果应为函数返回值")

	// 条件为 false，返回默认值且无错误
	ch3 := mathx.IfDoWithErrorAsync(false, successFunc, 999)
	res3 := <-ch3
	assert.NoError(res3.Err, "条件为 false 时错误应为 nil")
	assert.Equal(999, res3.Result, "条件为 false 应返回默认值")
}
