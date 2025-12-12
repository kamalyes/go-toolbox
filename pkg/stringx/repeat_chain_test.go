/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-22 10:07:57
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:30:15
 * @FilePath: \go-toolbox\pkg\stringx\repeat_chain_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatChain(t *testing.T) {
	result := New("Hello").RepeatChain(3).Value()
	assert.Equal(t, "HelloHelloHello", result)
}

func TestRepeatByLengthChain(t *testing.T) {
	result := New("abc").RepeatByLengthChain(10).Value()
	assert.Equal(t, "abcabcabca", result)

	result = New("abc").RepeatByLengthChain(2).Value()
	assert.Equal(t, "ab", result)

	result = New("abc").RepeatByLengthChain(0).Value()
	assert.Equal(t, "", result)

	result = New("").RepeatByLengthChain(5).Value()
	assert.Equal(t, "", result)
}

func TestRepeatAndJoinChain(t *testing.T) {
	result := New("Hello").RepeatAndJoinChain(", ", 3).Value()
	assert.Equal(t, "Hello, Hello, Hello", result)

	result = New("Hi").RepeatAndJoinChain(" - ", 2).Value()
	assert.Equal(t, "Hi - Hi", result)

	result = New("Test").RepeatAndJoinChain(" ", 0).Value()
	assert.Equal(t, "", result)
}
