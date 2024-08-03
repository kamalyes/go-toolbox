/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 10:39:01
 * @FilePath: \go-toolbox\stringx\sub_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllSubFunctions(t *testing.T) {
	t.Run("TestSubBefore", TestSubBefore)
	t.Run("TestSubAfter", TestSubAfter)
	t.Run("TestSubBetween", TestSubBetween)
	t.Run("TestSubBetweenAll", TestSubBetweenAll)
}

func TestSubBefore(t *testing.T) {
	result := SubBefore("abcdef", "d", false)
	assert.Equal(t, "abc", result)

	result = SubBefore("abcdef", "d", true)
	assert.Equal(t, "abc", result)

	result = SubBefore("abcdef", "x", false)
	assert.Equal(t, "abcdef", result)

	result = SubBefore("abcdef", "a", false)
	assert.Equal(t, "", result)
}

func TestSubAfter(t *testing.T) {
	result := SubAfter("abcdef", "d", true)
	assert.Equal(t, "def", result)

	result = SubAfter("abcdef", "x", false)
	assert.Equal(t, "", result)

	result = SubAfter("abcdef", "f", false)
	assert.Equal(t, "", result)
}

func TestSubBetween(t *testing.T) {
	result := SubBetween("abc123def456ghi", "abc", "def")
	assert.Equal(t, "123", result)

	result = SubBetween("abc123def456ghi", "def", "ghi")
	assert.Equal(t, "456", result)

	result = SubBetween("abc123def456ghi", "abc", "xyz")
	assert.Equal(t, "", result)
}

func TestSubBetweenAll(t *testing.T) {
	result := SubBetweenAll("a1b2c3d4e5f6", "b", "e")
	assert.ElementsMatch(t, []string{"2c3d4"}, result)

	result = SubBetweenAll("a1b2c3d4e5f6", "x", "y")
	assert.ElementsMatch(t, []string{}, result)
}
