/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 23:20:15
 * @FilePath: \go-toolbox\pkg\uuid\uuid_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringUtilFunctions(t *testing.T) {
	t.Run("TestSameSubStr", TestSameSubStr)
	t.Run("TestUUID", TestUUID)
	t.Run("TestUniqueID", TestUniqueID)
	t.Run("TestMd5", TestMd5)
}

func TestSameSubStr(t *testing.T) {
	subStr := "abc"
	repeat := 3
	expected := "abcabcabc"

	result := SameSubStr(subStr, repeat)
	assert.Equal(t, expected, result)
}

func TestUUID(t *testing.T) {
	uuid := UUID()
	assert.NotEmpty(t, uuid)
}

func TestUniqueID(t *testing.T) {
	id := UniqueID("hello", "world")
	assert.NotEmpty(t, id)
	id = UniqueID()
	assert.NotEmpty(t, id)
}

func TestMd5(t *testing.T) {
	src := "hello"
	expected := "5d41402abc4b2a76b9719d911017c592"

	result := Md5(src)
	assert.Equal(t, expected, result)
}
