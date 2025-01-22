/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-22 09:52:35
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
