/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 10:43:17
 * @FilePath: \go-toolbox\stringx\start_end_with_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

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
	assert.True(t, StartWith("hello", "he"))
	assert.False(t, StartWith("hello", "lo"))
	assert.True(t, StartWith("", ""))
}

func TestStartWithIgnoreCase(t *testing.T) {
	assert.True(t, StartWithIgnoreCase("Hello", "he"))
	assert.False(t, StartWithIgnoreCase("Hello", "lo"))
}

func TestStartWithAny(t *testing.T) {
	assert.True(t, StartWithAny("hello", []string{"he", "lo", "hi"}))
	assert.False(t, StartWithAny("hello", []string{"lo", "hi"}))
	assert.False(t, StartWithAny("", []string{"lo", "hi"}))
}

func TestStartWithAnyIgnoreCase(t *testing.T) {
	assert.True(t, StartWithAnyIgnoreCase("Hello", []string{"he", "Lo", "HI"}))
	assert.False(t, StartWithAnyIgnoreCase("Hello", []string{"lo", "hi"}))
}

func TestEndWith(t *testing.T) {
	assert.True(t, EndWith("hello", "lo"))
	assert.False(t, EndWith("hello", "he"))
	assert.True(t, EndWith("", ""))
}

func TestEndWithIgnoreCase(t *testing.T) {
	assert.True(t, EndWithIgnoreCase("Hello", "LO"))
	assert.False(t, EndWithIgnoreCase("Hello", "HE"))
}

func TestEndWithAny(t *testing.T) {
	assert.True(t, EndWithAny("hello", []string{"lo", "hi", "la"}))
	assert.False(t, EndWithAny("hello", []string{"he", "hi"}))
	assert.False(t, EndWithAny("", []string{"he", "hi"}))
}

func TestEndWithAnyIgnoreCase(t *testing.T) {
	assert.True(t, EndWithAnyIgnoreCase("hello", []string{"LO", "hi", "la"}))
	assert.False(t, EndWithAnyIgnoreCase("hello", []string{"HE", "hi"}))
}
