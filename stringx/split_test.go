/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 10:49:11
 * @FilePath: \go-toolbox\stringx\split_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"testing"

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
	result := Split("one,two,three,four", ",")
	assert.Equal(t, []string{"one", "two", "three", "four"}, result)
}

func TestSplitLimit(t *testing.T) {
	result := SplitLimit("one,two,three,four", ",", 2)
	assert.Equal(t, []string{"one", "two,three,four"}, result)
}

func TestSplitTrim(t *testing.T) {
	result := SplitTrim(" one , two , three , four ", ",")
	assert.Equal(t, []string{"one", "two", "three", "four"}, result)
}

func TestSplitTrimLimit(t *testing.T) {
	result := SplitTrimLimit(" one , two , three , four ", ",", 2)
	assert.Equal(t, []string{"one", "two , three , four"}, result)
}

func TestSplitByLen(t *testing.T) {
	result := SplitByLen("abcdefghij", 3)
	assert.Equal(t, []string{"abc", "def", "ghi", "j"}, result)
}

func TestCut(t *testing.T) {
	result := Cut("abcdefghij", 3)
	assert.Equal(t, []string{"abcd", "efg", "hij"}, result)

	result = Cut("abcdefghij", 4)
	assert.Equal(t, []string{"abc", "def", "gh", "ij"}, result)
}
