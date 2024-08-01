/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 18:21:13
 * @FilePath: \go-toolbox\random\rand_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllRandomFunctions(t *testing.T) {
	t.Run("TestRandInt", TestRandInt)
	t.Run("TestRandFloat", TestRandFloat)
	t.Run("TestRandString", TestRandString)
	t.Run("TestRandomStr", TestRandomStr)
	t.Run("TestRandomNum", TestRandomNum)
	t.Run("TestRandomHex", TestRandomHex)
}

func TestRandInt(t *testing.T) {
	min := 10
	max := 20

	result := RandInt(min, max)

	assert.GreaterOrEqual(t, result, min, "Expected result to be greater than or equal to min")
	assert.LessOrEqual(t, result, max, "Expected result to be less than or equal to max")
}

func TestRandFloat(t *testing.T) {
	min := 10.5
	max := 20.5

	result := RandFloat(min, max)

	assert.GreaterOrEqual(t, result, min, "Expected result to be greater than or equal to min")
	assert.LessOrEqual(t, result, max, "Expected result to be less than or equal to max")
}

func TestRandString(t *testing.T) {
	str := RandString(10, CAPITAL|LOWERCASE|SPECIAL|NUMBER)

	assert.Len(t, str, 10, "Expected string length to be 10")
}

func TestRandomStr(t *testing.T) {
	length := 8
	str := RandomStr(length)

	assert.Len(t, str, length, "Expected string length to be 8")
}

func TestRandomNum(t *testing.T) {
	length := 6
	num := RandomNum(length)

	assert.Len(t, num, length, "Expected number length to be 6")
}

func TestRandomHex(t *testing.T) {
	bytesLen := 4
	hex := RandomHex(bytesLen)

	assert.Len(t, hex, bytesLen*2, "Expected hex length to be twice the bytes length")
}
