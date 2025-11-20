/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-20 11:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-20 21:27:25
 * @FilePath: \go-toolbox\pkg\mathx\ternary_extended_test.go
 * @Description: 扩展三元运算符功能测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestIfNotNil(t *testing.T) {
	// 测试指针不为空的情况
	value := 42
	result := IfNotNil(&value, 0)
	assert.Equal(t, 42, result)

	// 测试指针为空的情况
	var nilPtr *int
	result = IfNotNil(nilPtr, 100)
	assert.Equal(t, 100, result)
}

func TestIfNotEmpty(t *testing.T) {
	// 测试非空字符串
	result := IfNotEmpty("hello", "default")
	assert.Equal(t, "hello", result)

	// 测试空字符串
	result = IfNotEmpty("", "default")
	assert.Equal(t, "default", result)
}

func TestIfNotZero(t *testing.T) {
	// 测试非零值
	result := IfNotZero(42, 100)
	assert.Equal(t, 42, result)

	// 测试零值
	result = IfNotZero(0, 100)
	assert.Equal(t, 100, result)

	// 测试字符串零值
	strResult := IfNotZero("", "default")
	assert.Equal(t, "default", strResult)
}

func TestIfContains(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	// 测试包含的情况
	result := IfContains(slice, "banana", "found", "not found")
	assert.Equal(t, "found", result)

	// 测试不包含的情况
	result = IfContains(slice, "grape", "found", "not found")
	assert.Equal(t, "not found", result)
}

func TestIfAny(t *testing.T) {
	// 测试有条件为真的情况
	conditions := []bool{false, false, true, false}
	result := IfAny(conditions, "success", "failure")
	assert.Equal(t, "success", result)

	// 测试所有条件为假的情况
	conditions = []bool{false, false, false}
	result = IfAny(conditions, "success", "failure")
	assert.Equal(t, "failure", result)
}

func TestIfAll(t *testing.T) {
	// 测试所有条件为真的情况
	conditions := []bool{true, true, true}
	result := IfAll(conditions, "success", "failure")
	assert.Equal(t, "success", result)

	// 测试有条件为假的情况
	conditions = []bool{true, false, true}
	result = IfAll(conditions, "success", "failure")
	assert.Equal(t, "failure", result)

	// 测试空条件数组
	conditions = []bool{}
	result = IfAll(conditions, "success", "failure")
	assert.Equal(t, "failure", result)
}

func TestIfCount(t *testing.T) {
	conditions := []bool{true, false, true, true, false}

	// 测试达到阈值的情况
	result := IfCount(conditions, 3, "enough", "not enough")
	assert.Equal(t, "enough", result)

	// 测试未达到阈值的情况
	result = IfCount(conditions, 4, "enough", "not enough")
	assert.Equal(t, "not enough", result)
}

func TestIfMap(t *testing.T) {
	// 测试条件为真时的映射
	result := IfMap(true, "hello", strings.ToUpper, "default")
	assert.Equal(t, "HELLO", result)

	// 测试条件为假时返回默认值
	result = IfMap(false, "hello", strings.ToUpper, "default")
	assert.Equal(t, "default", result)
}

func TestIfMapElse(t *testing.T) {
	toUpper := func(s string) string { return strings.ToUpper(s) }
	toLower := func(s string) string { return strings.ToLower(s) }

	// 测试条件为真
	result := IfMapElse(true, "Hello", toUpper, toLower)
	assert.Equal(t, "HELLO", result)

	// 测试条件为假
	result = IfMapElse(false, "Hello", toUpper, toLower)
	assert.Equal(t, "hello", result)
}

func TestIfFilter(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6}
	isEven := func(n int) bool { return n%2 == 0 }

	// 测试使用过滤
	result := IfFilter(true, numbers, isEven)
	expected := []int{2, 4, 6}
	assert.Equal(t, expected, result)

	// 测试不使用过滤
	result = IfFilter(false, numbers, isEven)
	assert.Equal(t, numbers, result)
}

func TestIfValidate(t *testing.T) {
	isPositive := func(n int) bool { return n > 0 }

	// 测试验证通过
	result := IfValidate(5, isPositive, "valid", "invalid")
	assert.Equal(t, "valid", result)

	// 测试验证失败
	result = IfValidate(-5, isPositive, "valid", "invalid")
	assert.Equal(t, "invalid", result)
}

func TestIfCast(t *testing.T) {
	var value interface{} = "hello"

	// 测试成功转换
	result := IfCast[string](value, "default")
	assert.Equal(t, "hello", result)

	// 测试转换失败
	result = IfCast[string](123, "default")
	assert.Equal(t, "default", result)
}

func TestIfBetween(t *testing.T) {
	// 测试在范围内
	result := IfBetween(5, 1, 10, 100, 200)
	assert.Equal(t, 100, result)

	// 测试超出范围
	result = IfBetween(15, 1, 10, 100, 200)
	assert.Equal(t, 200, result)

	// 测试边界值
	result = IfBetween(1, 1, 10, 100, 200)
	assert.Equal(t, 100, result)
}

func TestIfSwitch(t *testing.T) {
	cases := map[string]string{
		"red":    "stop",
		"green":  "go",
		"yellow": "caution",
	}

	// 测试匹配的情况
	result := IfSwitch("red", cases, "unknown")
	assert.Equal(t, "stop", result)

	// 测试不匹配的情况
	result = IfSwitch("blue", cases, "unknown")
	assert.Equal(t, "unknown", result)
}

func TestIfTryParse(t *testing.T) {
	parser := func(s string) (int, error) {
		return strconv.Atoi(s)
	}

	// 测试解析成功
	result := IfTryParse("123", parser, 0)
	assert.Equal(t, 123, result)

	// 测试解析失败
	result = IfTryParse("abc", parser, 0)
	assert.Equal(t, 0, result)
}

func TestIfSafeIndex(t *testing.T) {
	slice := []string{"a", "b", "c"}

	// 测试有效索引
	result := IfSafeIndex(slice, 1, "default")
	assert.Equal(t, "b", result)

	// 测试无效索引
	result = IfSafeIndex(slice, 5, "default")
	assert.Equal(t, "default", result)

	// 测试负索引
	result = IfSafeIndex(slice, -1, "default")
	assert.Equal(t, "default", result)
}

func TestIfSafeKey(t *testing.T) {
	m := map[string]int{
		"apple":  1,
		"banana": 2,
		"orange": 3,
	}

	// 测试存在的键
	result := IfSafeKey(m, "banana", 0)
	assert.Equal(t, 2, result)

	// 测试不存在的键
	result = IfSafeKey(m, "grape", 0)
	assert.Equal(t, 0, result)
}

func TestIfMulti(t *testing.T) {
	values := []string{"red", "green", "blue"}

	// 测试匹配的情况
	result := IfMulti("green", values, "primary", "not primary")
	assert.Equal(t, "primary", result)

	// 测试不匹配的情况
	result = IfMulti("purple", values, "primary", "not primary")
	assert.Equal(t, "not primary", result)
}

func TestIfPipeline(t *testing.T) {
	funcs := []func(string) string{
		strings.ToUpper,
		func(s string) string { return s + "!" },
		func(s string) string { return ">>> " + s },
	}

	// 测试条件为真时的管道处理
	result := IfPipeline(true, "hello", funcs, "default")
	expected := ">>> HELLO!"
	assert.Equal(t, expected, result)

	// 测试条件为假时返回默认值
	result = IfPipeline(false, "hello", funcs, "default")
	assert.Equal(t, "default", result)
}

func TestIfLazy(t *testing.T) {
	expensiveTrue := func() string {
		return "expensive true result"
	}
	expensiveFalse := func() string {
		return "expensive false result"
	}

	// 测试条件为真
	result := IfLazy(true, expensiveTrue, expensiveFalse)
	assert.Equal(t, "expensive true result", result)

	// 测试条件为假
	result = IfLazy(false, expensiveTrue, expensiveFalse)
	assert.Equal(t, "expensive false result", result)
}

func TestIfMemoized(t *testing.T) {
	cache := make(map[string]string)
	callCount := 0

	computeFn := func() string {
		callCount++
		return "computed result"
	}

	// 第一次调用
	result := IfMemoized(true, "key1", cache, computeFn, "default")
	assert.Equal(t, "computed result", result)
	assert.Equal(t, 1, callCount)

	// 第二次调用相同key，应该使用缓存
	result = IfMemoized(true, "key1", cache, computeFn, "default")
	assert.Equal(t, "computed result", result)
	assert.Equal(t, 1, callCount)

	// 条件为假时返回默认值
	result = IfMemoized(false, "key2", cache, computeFn, "default")
	assert.Equal(t, "default", result)
	assert.Equal(t, 1, callCount)
}

// 基准测试
func BenchmarkIfNotNil(b *testing.B) {
	value := 42
	for i := 0; i < b.N; i++ {
		IfNotNil(&value, 0)
	}
}

func BenchmarkIfContains(b *testing.B) {
	slice := []string{"apple", "banana", "orange", "grape", "watermelon"}
	for i := 0; i < b.N; i++ {
		IfContains(slice, "banana", "found", "not found")
	}
}

func BenchmarkIfLazy(b *testing.B) {
	trueFn := func() string { return "true result" }
	falseFn := func() string { return "false result" }

	for i := 0; i < b.N; i++ {
		IfLazy(i%2 == 0, trueFn, falseFn)
	}
}
