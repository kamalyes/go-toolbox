/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-22 10:07:57
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-22 10:07:57
 * @FilePath: \go-toolbox\tests\stringx_sub_chain_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/stretchr/testify/assert"
)

func TestSubBeforeChain(t *testing.T) {
	result := stringx.New("Hello, World!").SubBeforeChain(",", false).Value()
	assert.Equal(t, "Hello", result)

	result = stringx.New("Hello, World!").SubBeforeChain(",", true).Value()
	assert.Equal(t, "Hello", result)

	result = stringx.New("Hello World!").SubBeforeChain(",", false).Value()
	assert.Equal(t, "Hello World!", result) // Separator not found
}

func TestSubAfterChain(t *testing.T) {
	result := stringx.New("Hello, World!").SubAfterChain(",", false).Value()
	assert.Equal(t, " World!", result)

	result = stringx.New("Hello, World!").SubAfterChain(",", true).Value()
	assert.Equal(t, " World!", result)

	result = stringx.New("Hello World!").SubAfterChain(",", false).Value()
	assert.Equal(t, "", result) // Separator not found
}

func TestSubBetweenChain(t *testing.T) {
	result := stringx.New("Hello [World]!").SubBetweenChain("[", "]").Value()
	assert.Equal(t, "World", result)

	result = stringx.New("Hello World!").SubBetweenChain("[", "]").Value()
	assert.Equal(t, "", result) // Delimiters not found
}
