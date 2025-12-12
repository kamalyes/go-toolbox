/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:28:39
 * @FilePath: \go-toolbox\pkg\stringx\contains_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsChain(t *testing.T) {
	result := New("Hello, World!").ContainsChain("World")
	assert.True(t, result)
}

func TestContainsIgnoreCaseChain(t *testing.T) {
	result := New("Hello, World!").ContainsIgnoreCaseChain("world")
	assert.True(t, result)
}

func TestContainsAnyChain(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := New("I like banana").ContainsAnyChain(strList)
	assert.True(t, result)
}

func TestContainsAnyIgnoreCaseChain(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := New("I like BANANA").ContainsAnyIgnoreCaseChain(strList)
	assert.True(t, result)
}

func TestContainsAllChain(t *testing.T) {
	strList := []string{"Hello", "World"}
	result := New("Hello, World!").ContainsAllChain(strList)
	assert.True(t, result)
}

func TestContainsBlankChain(t *testing.T) {
	result := New("Hello World").ContainsBlankChain()
	assert.True(t, result)

	result = New("HelloWorld").ContainsBlankChain()
	assert.False(t, result)
}

func TestGetContainsStrChain(t *testing.T) {
	strList := []string{"apple", "banana", "orange"}
	result := New("I like banana").GetContainsStrChain(strList)
	assert.Equal(t, "banana", result)

	result = New("I like grapes").GetContainsStrChain(strList)
	assert.Equal(t, "", result)
}
