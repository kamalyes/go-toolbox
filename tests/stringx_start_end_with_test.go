/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 10:55:17
 * @FilePath: \go-toolbox\tests\start_end_with_test.go
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

func TestAllStartEndWithFunctions(t *testing.T) {
	t.Run("TestStartWith", TestStartWith)
	t.Run("TestStartWithIgnoreCase", TestStartWithIgnoreCase)
	t.Run("TestStartWithAny", TestStartWithAny)
	t.Run("TestStartWithAnyIgnoreCase", TestStartWithAnyIgnoreCase)
	t.Run("TestEndWith", TestEndWith)
	t.Run("TestEndWithIgnoreCase", TestEndWithIgnoreCase)
	t.Run("TestEndWithAny", TestEndWithAny)
	t.Run("TestEndWithAnyIgnoreCase", TestEndWithAnyIgnoreCase)

}

func TestStartWith(t *testing.T) {
	assert.True(t, stringx.StartWith("hello", "he"))
	assert.False(t, stringx.StartWith("hello", "lo"))
	assert.True(t, stringx.StartWith("", ""))
}

func TestStartWithIgnoreCase(t *testing.T) {
	assert.True(t, stringx.StartWithIgnoreCase("Hello", "he"))
	assert.False(t, stringx.StartWithIgnoreCase("Hello", "lo"))
}

func TestStartWithAny(t *testing.T) {
	assert.True(t, stringx.StartWithAny("hello", []string{"he", "lo", "hi"}))
	assert.False(t, stringx.StartWithAny("hello", []string{"lo", "hi"}))
	assert.False(t, stringx.StartWithAny("", []string{"lo", "hi"}))
}

func TestStartWithAnyIgnoreCase(t *testing.T) {
	assert.True(t, stringx.StartWithAnyIgnoreCase("Hello", []string{"he", "Lo", "HI"}))
	assert.False(t, stringx.StartWithAnyIgnoreCase("Hello", []string{"lo", "hi"}))
}

func TestEndWith(t *testing.T) {
	assert.True(t, stringx.EndWith("hello", "lo"))
	assert.False(t, stringx.EndWith("hello", "he"))
	assert.True(t, stringx.EndWith("", ""))
}

func TestEndWithIgnoreCase(t *testing.T) {
	assert.True(t, stringx.EndWithIgnoreCase("Hello", "LO"))
	assert.False(t, stringx.EndWithIgnoreCase("Hello", "HE"))
}

func TestEndWithAny(t *testing.T) {
	assert.True(t, stringx.EndWithAny("hello", []string{"lo", "hi", "la"}))
	assert.False(t, stringx.EndWithAny("hello", []string{"he", "hi"}))
	assert.False(t, stringx.EndWithAny("", []string{"he", "hi"}))
}

func TestEndWithAnyIgnoreCase(t *testing.T) {
	assert.True(t, stringx.EndWithAnyIgnoreCase("hello", []string{"LO", "hi", "la"}))
	assert.False(t, stringx.EndWithAnyIgnoreCase("hello", []string{"HE", "hi"}))
}
