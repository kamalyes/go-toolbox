/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 19:55:56
 * @FilePath: \go-toolbox\tests\validator_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"errors"
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
		{"", true},
		{"Hello", false},
		{nil, true},
		{0, true},
		{1, false},
		{[]int{}, true},
		{[]int{1, 2}, false},
		{map[string]int{}, true},
		{map[string]int{"key": 1}, false},
		{struct{}{}, true},
		{TestStruct{}, true},
		{TestStruct{Name: "Test"}, false},
	}

	for _, test := range tests {
		t.Run(func() string {
			if test.value == nil {
				return "nil"
			}
			return reflect.TypeOf(test.value).String()
		}(), func(t *testing.T) {
			// 直接调用 IsEmptyValue
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

func TestVerify(t *testing.T) {
	tests := []struct {
		input    TestStruct
		expected error
	}{
		{TestStruct{Name: "John", Age: 30, Email: "john@example.com"}, nil},
		{TestStruct{Name: "", Age: 30, Email: "john@example.com"}, errors.New("Name值不能为空")},
		{TestStruct{Name: "John", Age: -1, Email: "john@example.com"}, errors.New("Age长度或值不在合法范围,ge=0")},
		{TestStruct{Name: "John", Age: 30, Email: "invalid-email"}, errors.New("Email格式校验不通过")},
	}

	roleMap := validator.Rules{
		"Name":  {"notEmpty"},
		"Age":   {"ge=0"},
		"Email": {"regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"},
	}

	for _, test := range tests {
		t.Run(test.input.Name, func(t *testing.T) {
			err := validator.Verify(test.input, roleMap)
			if test.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expected.Error())
			}
		})
	}
}
