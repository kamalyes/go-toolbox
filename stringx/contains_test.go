/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 17:05:53
 * @FilePath: \go-toolbox\stringx\contains_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllContainsFunctions(t *testing.T) {
	t.Run("TestContains", TestContains)
	t.Run("TestContainsIgnoreCase", TestContainsIgnoreCase)
	t.Run("TestContainsAny", TestContainsAny)
	t.Run("TestContainsAnyIgnoreCase", TestContainsAnyIgnoreCase)
	t.Run("TestContainsAll", TestContainsAll)
	t.Run("TestContainsBlank", TestContainsBlank)
	t.Run("TestGetContainsStr", TestGetContainsStr)
}

func TestContains(t *testing.T) {
	result := Contains("hello world", "lo")
	assert.True(t, result)
}

func TestContainsIgnoreCase(t *testing.T) {
	result := ContainsIgnoreCase("Hello WoRld", "lLow")
	assert.False(t, result)
}

func TestContainsAny(t *testing.T) {
	searchStrs := []string{"apple", "banana", "orange"}
	result := ContainsAny("I like apples", searchStrs)
	assert.True(t, result)
}

func TestContainsAnyIgnoreCase(t *testing.T) {
	searchStrs := []string{"apple", "banana", "orange"}
	result := ContainsAnyIgnoreCase("I like BaNAnas", searchStrs)
	assert.True(t, result)
}

func TestContainsAll(t *testing.T) {
	searchStrs := []string{"apple", "banana", "orange"}
	result := ContainsAll("I like apples and bananas", searchStrs)
	assert.False(t, result)
}

func TestContainsBlank(t *testing.T) {
	result := ContainsBlank("Hello,  World")
	assert.True(t, result)
}

func TestGetContainsStr(t *testing.T) {
	searchStrs := []string{"apple", "banana", "orange"}
	result := GetContainsStr("I like oranges", searchStrs)
	assert.Equal(t, "orange", result)
}
