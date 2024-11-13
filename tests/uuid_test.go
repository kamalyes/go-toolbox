/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 18:09:25
 * @FilePath: \go-toolbox\tests\uuid_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/uuid"
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

	result := uuid.SameSubStr(subStr, repeat)
	assert.Equal(t, expected, result)
}

func TestUUID(t *testing.T) {
	uuid := uuid.UUID()
	assert.NotEmpty(t, uuid)
}

func TestUniqueID(t *testing.T) {
	id := uuid.UniqueID("hello", "world")
	assert.NotEmpty(t, id)
	id = uuid.UniqueID()
	assert.NotEmpty(t, id)
}

func TestMd5(t *testing.T) {
	src := "hello"
	expected := "5d41402abc4b2a76b9719d911017c592"

	result := uuid.Md5(src)
	assert.Equal(t, expected, result)
}
