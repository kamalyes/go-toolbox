/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-05 16:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-05 11:15:15
 * @FilePath: \go-toolbox\pkg\types\slice_test.go
 * @Description: 切片操作测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestContains 测试 Contains 方法
func TestContains(t *testing.T) {
	t.Run("字符串切片", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		assert.True(t, Contains(slice, "banana"))
		assert.False(t, Contains(slice, "orange"))
	})

	t.Run("整数切片", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		assert.True(t, Contains(slice, 3))
		assert.False(t, Contains(slice, 10))
	})

	t.Run("空切片", func(t *testing.T) {
		var slice []string
		assert.False(t, Contains(slice, "test"))
	})
}

// TestContainsAny 测试 ContainsAny 方法
func TestContainsAny(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	t.Run("包含一个", func(t *testing.T) {
		assert.True(t, ContainsAny(slice, "banana", "orange"))
	})

	t.Run("包含多个", func(t *testing.T) {
		assert.True(t, ContainsAny(slice, "apple", "banana"))
	})

	t.Run("都不包含", func(t *testing.T) {
		assert.False(t, ContainsAny(slice, "orange", "grape"))
	})

	t.Run("空参数", func(t *testing.T) {
		assert.False(t, ContainsAny(slice))
	})
}

// TestContainsAll 测试 ContainsAll 方法
func TestContainsAll(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	t.Run("包含所有", func(t *testing.T) {
		assert.True(t, ContainsAll(slice, "apple", "banana"))
	})

	t.Run("缺少一个", func(t *testing.T) {
		assert.False(t, ContainsAll(slice, "apple", "orange"))
	})

	t.Run("空参数", func(t *testing.T) {
		assert.True(t, ContainsAll(slice))
	})
}

// TestIndexOf 测试 IndexOf 方法
func TestIndexOf(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	assert.Equal(t, 0, IndexOf(slice, "apple"))
	assert.Equal(t, 1, IndexOf(slice, "banana"))
	assert.Equal(t, 2, IndexOf(slice, "cherry"))
	assert.Equal(t, -1, IndexOf(slice, "orange"))
}

// TestFilter 测试 Filter 方法
func TestFilter(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6}

	evens := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})

	assert.Equal(t, []int{2, 4, 6}, evens)
}

// TestMapTR 测试 MapTR 方法
func TestMapTR(t *testing.T) {
	numbers := []int{1, 2, 3}

	doubled := MapTR(numbers, func(n int) int {
		return n * 2
	})

	assert.Equal(t, []int{2, 4, 6}, doubled)

	// 测试类型转换
	strings := MapTR(numbers, func(n int) string {
		return fmt.Sprintf("%d", n)
	})

	assert.Equal(t, []string{"1", "2", "3"}, strings)
}

// TestUnique 测试 Unique 方法
func TestUnique(t *testing.T) {
	t.Run("有重复", func(t *testing.T) {
		slice := []int{1, 2, 2, 3, 3, 3, 4}
		unique := Unique(slice)
		assert.Equal(t, []int{1, 2, 3, 4}, unique)
	})

	t.Run("无重复", func(t *testing.T) {
		slice := []int{1, 2, 3, 4}
		unique := Unique(slice)
		assert.Equal(t, []int{1, 2, 3, 4}, unique)
	})

	t.Run("空切片", func(t *testing.T) {
		var slice []int
		unique := Unique(slice)
		assert.Empty(t, unique)
	})
}

// TestReverse 测试 Reverse 方法
func TestReverse(t *testing.T) {
	t.Run("正常反转", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		reversed := Reverse(slice)
		assert.Equal(t, []int{5, 4, 3, 2, 1}, reversed)
		// 确保原切片未修改
		assert.Equal(t, []int{1, 2, 3, 4, 5}, slice)
	})

	t.Run("单元素", func(t *testing.T) {
		slice := []int{1}
		reversed := Reverse(slice)
		assert.Equal(t, []int{1}, reversed)
	})

	t.Run("空切片", func(t *testing.T) {
		var slice []int
		reversed := Reverse(slice)
		assert.Empty(t, reversed)
	})
}

// TestChunk 测试 Chunk 方法
func TestChunk(t *testing.T) {
	t.Run("正常分块", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7}
		chunks := Chunk(slice, 3)
		assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7}}, chunks)
	})

	t.Run("大小等于长度", func(t *testing.T) {
		slice := []int{1, 2, 3}
		chunks := Chunk(slice, 3)
		assert.Equal(t, [][]int{{1, 2, 3}}, chunks)
	})

	t.Run("大小大于长度", func(t *testing.T) {
		slice := []int{1, 2, 3}
		chunks := Chunk(slice, 5)
		assert.Equal(t, [][]int{{1, 2, 3}}, chunks)
	})

	t.Run("大小为1", func(t *testing.T) {
		slice := []int{1, 2, 3}
		chunks := Chunk(slice, 1)
		assert.Equal(t, [][]int{{1}, {2}, {3}}, chunks)
	})

	t.Run("大小为0或负数", func(t *testing.T) {
		slice := []int{1, 2, 3}
		assert.Nil(t, Chunk(slice, 0))
		assert.Nil(t, Chunk(slice, -1))
	})

	t.Run("空切片", func(t *testing.T) {
		var slice []int
		chunks := Chunk(slice, 2)
		assert.Empty(t, chunks)
	})
}

// BenchmarkContains 性能测试
func BenchmarkContains(b *testing.B) {
	slice := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Contains(slice, "e")
	}
}

// BenchmarkUnique 性能测试
func BenchmarkUnique(b *testing.B) {
	slice := []int{1, 2, 3, 2, 4, 3, 5, 1, 6, 4, 7, 5}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unique(slice)
	}
}
