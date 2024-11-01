/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 10:49:11
 * @FilePath: \go-toolbox\tests\split_test.go
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

func TestAllSplitFunctions(t *testing.T) {
	t.Run("TestSplit", TestSplit)
	t.Run("TestSplitLimit", TestSplitLimit)
	t.Run("TestSplitTrim", TestSplitTrim)
	t.Run("TestSplitTrimLimit", TestSplitTrimLimit)
	t.Run("TestSplitByLen", TestSplitByLen)
	t.Run("TestEndWithIgnoreCase", TestEndWithIgnoreCase)
	t.Run("TestEndWithAny", TestEndWithAny)
	t.Run("TestCut", TestCut)

}

func TestSplit(t *testing.T) {
	result := stringx.Split("one,two,three,four", ",")
	assert.Equal(t, []string{"one", "two", "three", "four"}, result)
}

func TestSplitLimit(t *testing.T) {
	result := stringx.SplitLimit("one,two,three,four", ",", 2)
	assert.Equal(t, []string{"one", "two,three,four"}, result)
}

func TestSplitTrim(t *testing.T) {
	result := stringx.SplitTrim(" one , two , three , four ", ",")
	assert.Equal(t, []string{"one", "two", "three", "four"}, result)
}

func TestSplitTrimLimit(t *testing.T) {
	result := stringx.SplitTrimLimit(" one , two , three , four ", ",", 2)
	assert.Equal(t, []string{"one", "two , three , four"}, result)
}

func TestSplitByLen(t *testing.T) {
	result := stringx.SplitByLen("abcdefghij", 3)
	assert.Equal(t, []string{"abc", "def", "ghi", "j"}, result)
}

func TestCut(t *testing.T) {
	result := stringx.Cut("abcdefghij", 3)
	assert.Equal(t, []string{"abcd", "efg", "hij"}, result)

	result = stringx.Cut("abcdefghij", 4)
	assert.Equal(t, []string{"abc", "def", "gh", "ij"}, result)
}
