/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 17:03:40
 * @FilePath: \go-toolbox\stringx\base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllBaseFunctions(t *testing.T) {
	t.Run("TestLength", TestLength)
	t.Run("TestContainsIgnoreCase", TestContainsIgnoreCase)
	t.Run("TestContainsAny", TestContainsAny)
	t.Run("TestContainsAnyIgnoreCase", TestContainsAnyIgnoreCase)
	t.Run("TestContainsAll", TestContainsAll)
	t.Run("TestContainsBlank", TestContainsBlank)
	t.Run("TestGetContainsStr", TestGetContainsStr)
	t.Run("TestReverse", TestReverse)
	t.Run("TestEquals", TestEquals)
	t.Run("TestEqualsIgnoreCase", TestEqualsIgnoreCase)
	t.Run("TestInsertSpaces", TestInsertSpaces)
	t.Run("TestEqualsAny", TestEqualsAny)
	t.Run("TestEqualsAnyIgnoreCase", TestEqualsAnyIgnoreCase)
	t.Run("TestEqualsAt", TestEqualsAt)
	t.Run("TestCount", TestCount)
	t.Run("TestCompareIgnoreCase", TestCompareIgnoreCase)
}

func TestLength(t *testing.T) {
	result := Length("Hello, World!")
	assert.Equal(t, 13, result)
}

func TestReverse(t *testing.T) {
	result := Reverse("Hello, World!")
	assert.Equal(t, "!dlroW ,olleH", result)
}

func TestEquals(t *testing.T) {
	result := Equals("hello", "hello")
	assert.True(t, result)
}

func TestEqualsIgnoreCase(t *testing.T) {
	result := EqualsIgnoreCase("HELLO", "hello")
	assert.True(t, result)
}

func TestInsertSpaces(t *testing.T) {
	result := InsertSpaces("1234567890", 2)
	assert.Equal(t, "12 34 56 78 90", result)
}

func TestEqualsAny(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := EqualsAny("banana", strList)
	assert.True(t, result)
}

func TestEqualsAnyIgnoreCase(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := EqualsAnyIgnoreCase("OrAnGe", strList)
	assert.True(t, result)
}

func TestEqualsAt(t *testing.T) {
	result := EqualsAt("hello", 1, "e")
	assert.True(t, result)
}

func TestCount(t *testing.T) {
	result := Count("banana", "a")
	assert.Equal(t, 3, result)
}

func TestCompareIgnoreCase(t *testing.T) {
	result := CompareIgnoreCase("apple", "BANANA")
	assert.Less(t, result, 0)
}
