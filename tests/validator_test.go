/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-11 17:17:15
 * @FilePath: \go-toolbox\tests\validator_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"reflect"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/validator"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"notEmpty"`
	Age   int    `validate:"ge=0"`
	Email string `validate:"regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
}

func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{"", true},                                // 空字符串
		{"Hello", false},                          // 非空字符串
		{nil, true},                               // nil 值
		{0, true},                                 // 整数 0
		{1, false},                                // 非零整数
		{[]int{}, true},                           // 空切片
		{[]int{1, 2}, false},                      // 非空切片
		{map[string]int{}, true},                  // 空映射
		{map[string]int{"key": 1}, false},         // 非空映射
		{struct{}{}, true},                        // 空结构体
		{TestStruct{}, true},                      // 自定义结构体，所有字段零值
		{TestStruct{Name: "Test"}, false},         // 自定义结构体，非零字段
		{TestStruct{Name: "", Age: 0}, true},      // 自定义结构体，所有字段零值
		{TestStruct{Name: "Test", Age: 1}, false}, // 自定义结构体，至少一个非零字段
		{struct{ A int }{1}, false},               // 非空结构体
		{struct{ A interface{} }{nil}, true},      // 包含 nil 的结构体
		{make(chan int), false},                   // 非空通道
	}

	for _, test := range tests {
		t.Run(func() string {
			if test.value == nil {
				return "nil"
			}
			return reflect.TypeOf(test.value).String()
		}(), func(t *testing.T) {
			result := validator.IsEmptyValue(reflect.ValueOf(test.value))
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestHasEmpty(t *testing.T) {
	tests := []struct {
		elems    []interface{}
		expected bool
		count    int
	}{
		{[]interface{}{"", "data", nil}, true, 2},
		{[]interface{}{"data1", "data2"}, false, 0},
		{[]interface{}{0, 1, 2}, true, 1},
		{[]interface{}{0, "", nil}, true, 3},
	}

	for _, test := range tests {
		t.Run("HasEmpty", func(t *testing.T) {
			result, count := validator.HasEmpty(test.elems)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.count, count)
		})
	}
}

func TestIsAllEmpty(t *testing.T) {
	tests := []struct {
		elems    []interface{}
		expected bool
	}{
		{[]interface{}{"", nil}, true},
		{[]interface{}{"data", nil}, false},
		{[]interface{}{0, 0}, true},
		{[]interface{}{1, 0}, false},
	}

	for _, test := range tests {
		t.Run("IsAllEmpty", func(t *testing.T) {
			result := validator.IsAllEmpty(test.elems)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsUndefined(t *testing.T) {
	tests := []struct {
		str      string
		expected bool
	}{
		{"undefined", true},
		{"Undefined", true},
		{"UNDEFINED", true},
		{"defined", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			result := validator.IsUndefined(test.str)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestContainsChinese(t *testing.T) {
	tests := []struct {
		str      string
		expected bool
	}{
		{"Hello 你好", true},
		{"Hello World", false},
		{"", false},
		{"123", false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			result := validator.ContainsChinese(test.str)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestEmptyToDefault(t *testing.T) {
	tests := []struct {
		str        string
		defaultStr string
		expected   string
	}{
		{"", "default", "default"},
		{"value", "default", "value"},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			result := validator.EmptyToDefault(test.str, test.defaultStr)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsNil(t *testing.T) {
	// 测试 nil interface
	var nilInterface interface{}
	assert.True(t, validator.IsNil(nilInterface), "Expected nil interface to return true")

	// 测试 nil map
	var nilMap map[string]int
	assert.True(t, validator.IsNil(nilMap), "Expected nil map to return true")

	// 测试空 map
	emptyMap := make(map[string]int)
	assert.False(t, validator.IsNil(emptyMap), "Expected empty map to return false")

	// 测试非 nil map
	nonNilMap := map[string]int{"key": 1}
	assert.False(t, validator.IsNil(nonNilMap), "Expected non-nil map to return false")

	// 测试指向 nil 的 map
	var ptrToNilMap *map[string]int
	assert.True(t, validator.IsNil(ptrToNilMap), "Expected pointer to nil map to return true")

	// 测试指向空 map 的指针
	ptrToEmptyMap := &emptyMap
	assert.False(t, validator.IsNil(ptrToEmptyMap), "Expected pointer to empty map to return false")

	// 测试非 nil 指针
	num := 42
	ptrToNum := &num
	assert.False(t, validator.IsNil(ptrToNum), "Expected pointer to non-nil value to return false")

	// 测试 nil 切片
	var nilSlice []int
	assert.True(t, validator.IsNil(nilSlice), "Expected nil slice to return true")

	// 测试空切片
	emptySlice := []int{}
	assert.False(t, validator.IsNil(emptySlice), "Expected empty slice to return false")

	// 测试非 nil 切片
	nonNilSlice := []int{1, 2, 3}
	assert.False(t, validator.IsNil(nonNilSlice), "Expected non-nil slice to return false")

	// 测试指向 nil 切片的指针
	var ptrToNilSlice *[]int
	assert.True(t, validator.IsNil(ptrToNilSlice), "Expected pointer to nil slice to return true")

	// 测试指向空切片的指针
	ptrToEmptySlice := &emptySlice
	assert.False(t, validator.IsNil(ptrToEmptySlice), "Expected pointer to empty slice to return false")

	// 测试 nil 通道
	var nilChan chan int
	assert.True(t, validator.IsNil(nilChan), "Expected nil channel to return true")

	// 测试空通道
	emptyChan := make(chan int)
	assert.False(t, validator.IsNil(emptyChan), "Expected empty channel to return false")

	// 测试指向 nil 通道的指针
	var ptrToNilChan *chan int
	assert.True(t, validator.IsNil(ptrToNilChan), "Expected pointer to nil channel to return true")

	// 测试指向非 nil 通道的指针
	nonNilChan := make(chan int, 1)
	assert.False(t, validator.IsNil(nonNilChan), "Expected non-nil channel to return false")

	// 测试 nil 接口
	var nilInterfaceValue interface{}
	assert.True(t, validator.IsNil(nilInterfaceValue), "Expected nil interface value to return true")

	// 测试指向非 nil 接口的指针
	var nonNilInterfaceValue interface{} = 42
	ptrToNonNilInterface := &nonNilInterfaceValue
	assert.False(t, validator.IsNil(ptrToNonNilInterface), "Expected pointer to non-nil interface to return false")
}

func TestIsIPAllowed(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		cidrList []string
		want     bool
	}{
		// 正常情况
		{"Exact IP match", "192.168.1.100", []string{"192.168.1.100"}, true},
		{"IP in CIDR", "192.168.1.50", []string{"192.168.1.0/24"}, true},
		{"IP not in CIDR", "192.168.2.1", []string{"192.168.1.0/24"}, false},
		{"Multiple CIDRs one match", "10.0.0.5", []string{"192.168.1.0/24", "10.0.0.0/8"}, true},
		{"Multiple CIDRs none match", "172.16.0.1", []string{"192.168.1.0/24", "10.0.0.0/8"}, false},

		// IPv6
		{"IPv6 exact match", "2001:db8::1", []string{"2001:db8::1"}, true},
		{"IPv6 CIDR match", "2001:db8::abcd", []string{"2001:db8::/64"}, true},
		{"IPv6 CIDR no match", "2001:db9::1", []string{"2001:db8::/64"}, false},

		// 异常和边界情况
		{"Empty IP", "", []string{"192.168.1.0/24"}, false},
		{"Invalid IP format", "invalid-ip", []string{"192.168.1.0/24"}, false},
		{"Empty CIDR list", "192.168.1.50", []string{}, false},
		{"Nil CIDR list", "192.168.1.50", nil, false},
		{"CIDR list contains empty string", "192.168.1.50", []string{""}, false},
		{"CIDR list contains invalid CIDR", "192.168.1.50", []string{"invalid-cidr"}, false},
		{"IP equals CIDR but CIDR invalid", "192.168.1.50", []string{"192.168.1.50/33"}, false}, // 33不是合法掩码
		{"IP equals CIDR string but IP invalid", "999.999.999.999", []string{"999.999.999.999"}, false},

		// IP equals CIDR string exact match优先
		{"IP equals CIDR string exact match", "10.0.0.1", []string{"10.0.0.1", "10.0.0.0/8"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.IsIPAllowed(tt.ip, tt.cidrList)
			assert.Equal(t, tt.want, got)
		})
	}
}
