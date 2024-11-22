/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 10:39:01
 * @FilePath: \go-toolbox\tests\stringx_sub_test.go
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

func TestAllSubFunctions(t *testing.T) {
	t.Run("TestSubBefore", TestSubBefore)
	t.Run("TestSubAfter", TestSubAfter)
	t.Run("TestSubBetween", TestSubBetween)
	t.Run("TestSubBetweenAll", TestSubBetweenAll)
}

func TestSubBefore(t *testing.T) {
	result := stringx.SubBefore("abcdef", "d", false)
	assert.Equal(t, "abc", result)

	result = stringx.SubBefore("abcdef", "d", true)
	assert.Equal(t, "abc", result)

	result = stringx.SubBefore("abcdef", "x", false)
	assert.Equal(t, "abcdef", result)

	result = stringx.SubBefore("abcdef", "a", false)
	assert.Equal(t, "", result)
}

func TestSubAfter(t *testing.T) {
	result := stringx.SubAfter("abcdef", "d", true)
	assert.Equal(t, "ef", result)

	result = stringx.SubAfter("abcdef", "x", false)
	assert.Equal(t, "", result)

	result = stringx.SubAfter("abcdef", "f", false)
	assert.Equal(t, "", result)
}

func TestSubBetween(t *testing.T) {
	result := stringx.SubBetween("abc123def456ghi", "abc", "def")
	assert.Equal(t, "123", result)

	result = stringx.SubBetween("abc123def456ghi", "def", "ghi")
	assert.Equal(t, "456", result)

	result = stringx.SubBetween("abc123def456ghi", "abc", "xyz")
	assert.Equal(t, "", result)
}

func TestSubBetweenAll(t *testing.T) {
	result := stringx.SubBetweenAll("a1b2c3d4e5f6", "b", "e")
	assert.ElementsMatch(t, []string{"2c3d4"}, result)

	result = stringx.SubBetweenAll("a1b2c3d4e5f6", "x", "y")
	assert.ElementsMatch(t, []string{}, result)
}
