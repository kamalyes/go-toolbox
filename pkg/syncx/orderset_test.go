/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-11 13:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:59:04
 * @FilePath: \go-toolbox\pkg\syncx\orderset_test.go
 * @Description: ordered set 有序集合单元测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedSetBasic(t *testing.T) {
	set := NewOrderedSet[int]()

	// Add & Contains 测试
	set.Add(1)
	set.Add(2)
	set.Add(3)
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.True(t, set.Contains(3))
	assert.False(t, set.Contains(4))

	// Len 测试
	assert.Equal(t, 3, set.Len())

	// Remove 测试
	set.Remove(2)
	assert.False(t, set.Contains(2))
	assert.Equal(t, 2, set.Len())

	// Clear 测试
	set.Clear()
	assert.Equal(t, 0, set.Len())
	assert.False(t, set.Contains(1))
}

func TestOrderedSetConcurrent(t *testing.T) {
	set := NewOrderedSet[int]()
	var wg sync.WaitGroup
	concurrency := 100
	operations := 1000

	// 并发添加
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				set.Add(j)
			}
		}(i)
	}
	wg.Wait()

	// 因为元素唯一，最终元素数量应为 operations
	assert.Equal(t, operations, set.Len())

	// 并发删除
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				set.Remove(j)
			}
		}(i)
	}
	wg.Wait()

	// 删除完毕，集合应为空
	assert.Equal(t, 0, set.Len())
}

func TestOrderedSetConcurrentMixed(t *testing.T) {
	set := NewOrderedSet[string]()
	var wg sync.WaitGroup
	concurrency := 50
	operations := 500

	// 并发混合操作：添加、删除、Contains
	wg.Add(concurrency * 3)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				set.Add("item-" + strconv.Itoa(j))
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				set.Remove("item-" + strconv.Itoa(j))
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				set.Contains("item-" + strconv.Itoa(j))
			}
		}(i)
	}
	wg.Wait()

	// 最终元素数量不确定，确保不会崩溃且内部状态一致
	elems := set.Elements()
	uniqueMap := make(map[string]struct{})
	for _, e := range elems {
		uniqueMap[e] = struct{}{}
	}
	// 元素切片长度和唯一元素数应一致，保证无重复
	assert.Equal(t, len(elems), len(uniqueMap))

	// 元素数量应小于等于 operations
	assert.LessOrEqual(t, len(elems), operations)
}

func TestOrderedSetString(t *testing.T) {
	set := NewOrderedSet[int]()
	set.Add(10)
	set.Add(20)
	set.Add(30)
	str := set.String()
	expected := "OrderedSet{10, 20, 30}"
	assert.Equal(t, expected, str)
}

func TestNewOrderedSetFromSlice(t *testing.T) {
	items := []int{3, 1, 2, 3, 2, 4}
	set := NewOrderedSetFromSlice(items)

	// 元素顺序应为首次出现顺序，且无重复
	expected := []int{3, 1, 2, 4}
	assert.Equal(t, expected, set.Elements())
	assert.Equal(t, len(expected), set.Len())
}

func TestElementsReturnsCopy(t *testing.T) {
	set := NewOrderedSet[int]()
	set.Add(1)
	set.Add(2)
	set.Add(3)

	elems := set.Elements()
	assert.Equal(t, []int{1, 2, 3}, elems)

	// 修改 elems 不影响集合内部状态
	elems[0] = 100
	assert.Equal(t, 1, set.Elements()[0])
}

func TestRemoveAndContainsOnEmptySet(t *testing.T) {
	set := NewOrderedSet[string]()

	// 删除不存在元素不报错
	set.Remove("nonexistent")

	// Contains 返回 false
	assert.False(t, set.Contains("nonexistent"))
}

func TestAddDuplicateElements(t *testing.T) {
	set := NewOrderedSet[int]()
	set.Add(5)
	set.Add(5)
	set.Add(5)

	// 只添加一次，长度为1
	assert.Equal(t, 1, set.Len())
	assert.True(t, set.Contains(5))
}
