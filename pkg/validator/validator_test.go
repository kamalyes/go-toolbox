/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-11 17:17:15
 * @FilePath: \go-toolbox\pkg\validator\validator_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试常量 - IP 地址
const (
	testIPv4_1         = "192.168.1.100"
	testIPv4_2         = "192.168.1.50"
	testIPv4_3         = "192.168.1.1"
	testIPv4_4         = "10.0.0.1"
	testIPv4_5         = "10.0.0.5"
	testIPv4_6         = "172.16.0.1"
	testIPv4_Localhost = "127.0.0.1"
	testIPv4_Zero      = "0.0.0.0"
	testIPv6_1         = "2001:db8::1"
	testIPv6_2         = "2001:db8::abcd"
)

// 测试常量 - CIDR 范围
const (
	testCIDR_192      = "192.168.1.0/24"
	testCIDR_10       = "10.0.0.0/8"
	testCIDR_172      = "172.16.0.0/12"
	testCIDR_192Block = "192.168.0.0/16"
	testCIDR_127      = "127.0.0.0/8"
	testCIDR_IPv6     = "2001:db8::/64"
)

// 测试常量 - 通配符
const (
	testWildcard = "*"
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
			result := IsEmptyValue(reflect.ValueOf(test.value))
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
			result, count := HasEmpty(test.elems)
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
			result := IsAllEmpty(test.elems)
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
			result := IsUndefined(test.str)
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
			result := ContainsChinese(test.str)
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
			result := EmptyToDefault(test.str, test.defaultStr)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsNil(t *testing.T) {
	// 测试 nil interface
	var nilInterface interface{}
	assert.True(t, IsNil(nilInterface), "Expected nil interface to return true")

	// 测试 nil map
	var nilMap map[string]int
	assert.True(t, IsNil(nilMap), "Expected nil map to return true")

	// 测试空 map
	emptyMap := make(map[string]int)
	assert.False(t, IsNil(emptyMap), "Expected empty map to return false")

	// 测试非 nil map
	nonNilMap := map[string]int{"key": 1}
	assert.False(t, IsNil(nonNilMap), "Expected non-nil map to return false")

	// 测试指向 nil 的 map
	var ptrToNilMap *map[string]int
	assert.True(t, IsNil(ptrToNilMap), "Expected pointer to nil map to return true")

	// 测试指向空 map 的指针
	ptrToEmptyMap := &emptyMap
	assert.False(t, IsNil(ptrToEmptyMap), "Expected pointer to empty map to return false")

	// 测试非 nil 指针
	num := 42
	ptrToNum := &num
	assert.False(t, IsNil(ptrToNum), "Expected pointer to non-nil value to return false")

	// 测试 nil 切片
	var nilSlice []int
	assert.True(t, IsNil(nilSlice), "Expected nil slice to return true")

	// 测试空切片
	emptySlice := []int{}
	assert.False(t, IsNil(emptySlice), "Expected empty slice to return false")

	// 测试非 nil 切片
	nonNilSlice := []int{1, 2, 3}
	assert.False(t, IsNil(nonNilSlice), "Expected non-nil slice to return false")

	// 测试指向 nil 切片的指针
	var ptrToNilSlice *[]int
	assert.True(t, IsNil(ptrToNilSlice), "Expected pointer to nil slice to return true")

	// 测试指向空切片的指针
	ptrToEmptySlice := &emptySlice
	assert.False(t, IsNil(ptrToEmptySlice), "Expected pointer to empty slice to return false")

	// 测试 nil 通道
	var nilChan chan int
	assert.True(t, IsNil(nilChan), "Expected nil channel to return true")

	// 测试空通道
	emptyChan := make(chan int)
	assert.False(t, IsNil(emptyChan), "Expected empty channel to return false")

	// 测试指向 nil 通道的指针
	var ptrToNilChan *chan int
	assert.True(t, IsNil(ptrToNilChan), "Expected pointer to nil channel to return true")

	// 测试指向非 nil 通道的指针
	nonNilChan := make(chan int, 1)
	assert.False(t, IsNil(nonNilChan), "Expected non-nil channel to return false")

	// 测试 nil 接口
	var nilInterfaceValue interface{}
	assert.True(t, IsNil(nilInterfaceValue), "Expected nil interface value to return true")

	// 测试指向非 nil 接口的指针
	var nonNilInterfaceValue interface{} = 42
	ptrToNonNilInterface := &nonNilInterfaceValue
	assert.False(t, IsNil(ptrToNonNilInterface), "Expected pointer to non-nil interface to return false")
}

func TestIsIPAllowed(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		cidrList []string
		want     bool
	}{
		// ========== 通配符测试 ==========
		{"Wildcard allows any IP", testIPv4_1, []string{testWildcard}, true},
		{"Wildcard in list", testIPv4_4, []string{testCIDR_192, testWildcard, testCIDR_172}, true},
		{"Wildcard with specific IPs", "8.8.8.8", []string{testWildcard}, true},
		{"Multiple wildcards", "1.2.3.4", []string{testWildcard, testWildcard}, true},

		// ========== 空列表测试 ==========
		{"Empty list allows all", testIPv4_1, []string{}, true},
		{"Nil list allows all", testIPv4_4, nil, true},
		{"Empty list with any IP", testIPv4_Zero, []string{}, true},
		{"Empty list with IPv6", testIPv6_1, []string{}, true},

		// ========== 精确匹配测试 ==========
		{"Exact IP match", testIPv4_1, []string{testIPv4_1}, true},
		{"Exact IP no match", testIPv4_1, []string{"192.168.1.101"}, false},
		{"Multiple exact IPs one match", testIPv4_5, []string{testIPv4_3, testIPv4_5, testIPv4_6}, true},
		{"Multiple exact IPs none match", "10.0.0.6", []string{testIPv4_3, testIPv4_5, testIPv4_6}, false},

		// ========== CIDR 格式测试 ==========
		{"IP in CIDR", testIPv4_2, []string{testCIDR_192}, true},
		{"IP not in CIDR", "192.168.2.1", []string{testCIDR_192}, false},
		{"Multiple CIDRs one match", testIPv4_5, []string{testCIDR_192, testCIDR_10}, true},
		{"Multiple CIDRs none match", testIPv4_6, []string{testCIDR_192, testCIDR_10}, false},
		{"CIDR /32 exact match", testIPv4_1, []string{testIPv4_1 + "/32"}, true},
		{"CIDR /32 no match", "192.168.1.101", []string{testIPv4_1 + "/32"}, false},
		{"Large CIDR /8", "10.255.255.255", []string{testCIDR_10}, true},
		{"CIDR boundary test lower", "192.168.1.0", []string{testCIDR_192}, true},
		{"CIDR boundary test upper", "192.168.1.255", []string{testCIDR_192}, true},
		{"Outside CIDR boundary", "192.168.0.255", []string{testCIDR_192}, false},

		// ========== IPv6 测试 ==========
		{"IPv6 exact match", testIPv6_1, []string{testIPv6_1}, true},
		{"IPv6 exact no match", testIPv6_1, []string{"2001:db8::2"}, false},
		{"IPv6 CIDR match", testIPv6_2, []string{testCIDR_IPv6}, true},
		{"IPv6 CIDR no match", "2001:db9::1", []string{testCIDR_IPv6}, false},
		{"IPv6 loopback", "::1", []string{"::1"}, true},
		{"IPv6 any", "::", []string{"::"}, true},
		{"IPv6 with wildcard", testIPv6_1, []string{testWildcard}, true},

		// ========== 混合测试 ==========
		{"Mixed IPv4 and IPv6", testIPv4_3, []string{testCIDR_IPv6, testCIDR_192}, true},
		{"Mixed exact and CIDR", testIPv4_1, []string{testIPv4_4, testCIDR_192, testIPv4_6}, true},
		{"Wildcard with other rules", "1.2.3.4", []string{testCIDR_192, testWildcard}, true},

		// ========== 异常和边界情况 ==========
		{"Empty IP", "", []string{testCIDR_192}, false},
		{"Invalid IP format", "invalid-ip", []string{testCIDR_192}, false},
		{"Invalid IP with numbers", "256.256.256.256", []string{testCIDR_192}, false},
		{"CIDR list contains empty string", testIPv4_2, []string{""}, false},
		{"CIDR list contains invalid CIDR", testIPv4_2, []string{"invalid-cidr"}, false},
		{"IP equals CIDR but CIDR invalid", testIPv4_2, []string{testIPv4_2 + "/33"}, false}, // 33不是合法掩码
		{"IP equals CIDR string but IP invalid", "999.999.999.999", []string{"999.999.999.999"}, false},
		{"Malformed CIDR", testIPv4_3, []string{"192.168.1.0/99"}, false},
		{"Incomplete IP", "192.168.1", []string{testCIDR_192}, false},
		{"IP with port", testIPv4_3 + ":8080", []string{testCIDR_192}, false},

		// ========== 特殊 IP 地址 ==========
		{"Localhost IPv4", testIPv4_Localhost, []string{testIPv4_Localhost}, true},
		{"Localhost in CIDR", testIPv4_Localhost, []string{testCIDR_127}, true},
		{"Broadcast IP", "255.255.255.255", []string{"255.255.255.255"}, true},
		{"Zero IP", testIPv4_Zero, []string{testIPv4_Zero}, true},
		{"Private IP Class A", "10.1.2.3", []string{testCIDR_10}, true},
		{"Private IP Class B", "172.16.5.6", []string{testCIDR_172}, true},
		{"Private IP Class C", "192.168.100.1", []string{testCIDR_192Block}, true},

		// ========== 优先级测试 ==========
		{"Exact match before CIDR", testIPv4_4, []string{testIPv4_4, testCIDR_10}, true},
		{"Wildcard takes precedence", testIPv4_3, []string{testWildcard, "invalid-rule"}, true},

		// ========== 性能相关边界测试 ==========
		{"Many rules one match", testIPv4_3, []string{
			testCIDR_10, testCIDR_172, testCIDR_192Block, testIPv4_3,
		}, true},
		{"Many rules no match", "8.8.8.8", []string{
			testCIDR_10, testCIDR_172, testCIDR_192Block, testIPv4_Localhost,
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIPAllowed(tt.ip, tt.cidrList)
			assert.Equal(t, tt.want, got, "IP: %s, Rules: %v", tt.ip, tt.cidrList)
		})
	}
}

func TestIsFuncType(t *testing.T) {
	type FuncType func(int) int
	type MyStruct struct{ A int }

	tests := []struct {
		name     string
		typCheck func() bool
		want     bool
	}{
		{"int", func() bool { return IsFuncType[int]() }, false},
		{"string", func() bool { return IsFuncType[string]() }, false},
		{"struct", func() bool { return IsFuncType[MyStruct]() }, false},
		{"pointer", func() bool { return IsFuncType[*MyStruct]() }, false},
		{"slice", func() bool { return IsFuncType[[]int]() }, false},
		{"map", func() bool { return IsFuncType[map[string]int]() }, false},
		{"func type", func() bool { return IsFuncType[FuncType]() }, true},
		{"func literal type", func() bool { return IsFuncType[func(int) int]() }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.typCheck()
			assert.Equal(t, tt.want, got)
		})
	}
}
