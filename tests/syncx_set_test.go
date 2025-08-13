/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 15:05:55
 * @FilePath: \go-toolbox\tests\syncx_set_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := syncx.NewSet[string]()

	// 测试 Add 和 Has
	s.Add("apple")
	assert.True(t, s.Has("apple"), "expected set to contain 'apple'")
	assert.False(t, s.Has("banana"), "expected set to not contain 'banana'")

	// 测试 Size
	assert.Equal(t, 1, s.Size(), "expected size to be 1")

	// 测试 AddAll
	s.AddAll("banana", "cherry")
	assert.True(t, s.Has("banana"), "expected set to contain 'banana'")
	assert.True(t, s.Has("cherry"), "expected set to contain 'cherry'")
	assert.Equal(t, 3, s.Size(), "expected size to be 3")

	// 测试 HasAll
	existing, all := s.HasAll("apple", "banana", "grape")
	assert.ElementsMatch(t, existing, []string{"apple", "banana"}, "expected existing elements to match")
	assert.False(t, all, "expected not all elements to exist")

	existing, all = s.HasAll("apple", "banana")
	assert.ElementsMatch(t, existing, []string{"apple", "banana"}, "expected existing elements to match")
	assert.True(t, all, "expected all elements to exist")

	// 测试 Delete
	s.Delete("banana")
	assert.False(t, s.Has("banana"), "expected set to not contain 'banana' after deletion")
	assert.Equal(t, 2, s.Size(), "expected size to be 2 after deletion")

	// 测试 DeleteAll
	s.DeleteAll("apple", "cherry")
	assert.False(t, s.Has("apple"), "expected set to not contain 'apple' after deletion")
	assert.False(t, s.Has("cherry"), "expected set to not contain 'cherry' after deletion")
	assert.Equal(t, 0, s.Size(), "expected size to be 0 after deleting all elements")

	// 测试 Clear
	s.Add("date")
	s.Clear()
	assert.True(t, s.IsEmpty(), "expected set to be empty after clear")
}

func TestSet_Elements(t *testing.T) {
	s := syncx.NewSet[int]()
	s.AddAll(1, 2, 3)

	elements := s.Elements()
	assert.ElementsMatch(t, elements, []int{1, 2, 3}, "expected elements to match")
}

func TestSet_IsEmpty(t *testing.T) {
	s := syncx.NewSet[string]()
	assert.True(t, s.IsEmpty(), "expected new set to be empty")

	s.Add("item")
	assert.False(t, s.IsEmpty(), "expected set to not be empty after adding an item")

	s.Clear()
	assert.True(t, s.IsEmpty(), "expected set to be empty after clearing")
}
