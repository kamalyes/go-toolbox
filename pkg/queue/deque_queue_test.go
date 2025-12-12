/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-21 17:20:15
 * @FilePath: \go-toolbox\pkg\queue\deque_queue_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDequePushPop(t *testing.T) {
	q := NewDeque()

	// 测试 PushBack 和 PopFront
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	assert.Equal(t, 1, q.PopFront(), "Expected the first element to be 1")
	assert.Equal(t, 2, q.PopFront(), "Expected the second element to be 2")
	assert.Equal(t, 3, q.PopFront(), "Expected the third element to be 3")

	// 测试 PopFront 在空队列上
	assert.Panics(t, func() { q.PopFront() }, "Expected panic on PopFront from empty deque")
}

func TestDequePushFront(t *testing.T) {
	q := NewDeque()

	q.PushFront(1)
	q.PushFront(2)
	q.PushFront(3)

	assert.Equal(t, 1, q.PopBack(), "Expected the last element to be 1")
	assert.Equal(t, 2, q.PopBack(), "Expected the second last element to be 2")
	assert.Equal(t, 3, q.PopBack(), "Expected the third last element to be 3")

	// 测试 PopBack 在空队列上
	assert.Panics(t, func() { q.PopBack() }, "Expected panic on PopBack from empty deque")
}

func TestDequeFrontBack(t *testing.T) {
	q := NewDeque()

	q.PushBack(1)
	q.PushBack(2)

	front, err := q.Front()
	assert.NoError(t, err, "Expected no error when accessing Front")
	assert.Equal(t, 1, front, "Expected Front to be 1")

	back, err := q.Back()
	assert.NoError(t, err, "Expected no error when accessing Back")
	assert.Equal(t, 2, back, "Expected Back to be 2")

	// 测试 Front 和 Back 在空队列上
	q.PopFront()
	q.PopBack()

	_, err = q.Front()
	assert.Error(t, err, "Expected error when accessing Front from empty deque")

	_, err = q.Back()
	assert.Error(t, err, "Expected error when accessing Back from empty deque")
}

func TestDequeInsert(t *testing.T) {
	q := NewDeque()

	q.PushBack(1)
	q.PushBack(2)

	q.Insert(1, 3) // Insert 3 at index 1 (between 1 and 2)

	assert.Equal(t, 1, q.At(0), "Expected At(0) to be 1")
	assert.Equal(t, 3, q.At(1), "Expected At(1) to be 3")
	assert.Equal(t, 2, q.At(2), "Expected At(2) to be 2")

	// 测试无效索引
	assert.Panics(t, func() { q.Insert(-1, 4) }, "Expected panic on invalid index")
	assert.Panics(t, func() { q.Insert(4, 4) }, "Expected panic on invalid index")
}

func TestDequeClear(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)

	q.Clear()
	assert.Equal(t, 0, q.Len(), "Expected length to be 0 after clear")
	assert.Equal(t, 16, q.Cap(), "Expected capacity to remain the same after clear")
}

func TestDequeRotate(t *testing.T) {
	q := NewDeque()

	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)
	q.PushBack(4)

	q.Rotate(2) // 向前旋转 2 次
	assert.Equal(t, 3, q.PopFront(), "Expected the first element to be 3 after rotation")
	assert.Equal(t, 4, q.PopFront(), "Expected the second element to be 4 after rotation")
	assert.Equal(t, 1, q.PopFront(), "Expected the third element to be 1 after rotation")
	assert.Equal(t, 2, q.PopFront(), "Expected the fourth element to be 2 after rotation")

	// 测试负旋转
	q.PushBack(5)
	q.PushBack(6)
	q.Rotate(-1) // 向后旋转 1 次
	assert.Equal(t, 6, q.PopFront(), "Expected the first element to be 6 after backward rotation")
	assert.Equal(t, 5, q.PopFront(), "Expected the second element to be 5 after backward rotation")
}

func TestDequeIndex(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	index := q.Index(func(item interface{}) bool {
		return item == 2
	})
	assert.Equal(t, 1, index, "Expected index of item 2 to be 1")

	index = q.RIndex(func(item interface{}) bool {
		return item == 2
	})
	assert.Equal(t, 1, index, "Expected reverse index of item 2 to be 1")

	index = q.Index(func(item interface{}) bool {
		return item == 4
	})
	assert.Equal(t, -1, index, "Expected index of non-existent item to be -1")
}

func TestDequeGrow(t *testing.T) {
	q := NewDeque()

	// 测试 Grow 方法
	q.Grow(20) // 增长容量以容纳 20 个项目
	assert.GreaterOrEqual(t, q.Cap(), 20, "Expected capacity to be at least 20 after Grow")
}

func TestDequeGrowNegative(t *testing.T) {
	q := NewDeque()
	assert.Panics(t, func() { q.Grow(-1) }, "Expected panic when trying to grow with a negative number")
}

func TestDequeRotateNegativeTooLarge(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)
	q.Rotate(-5) // 旋转一个比长度大的负数
	assert.Equal(t, 2, q.PopFront(), "Expected the first element to be 3 after negative rotation")
	assert.Equal(t, 3, q.PopFront(), "Expected the second element to be 1 after negative rotation")
	assert.Equal(t, 1, q.PopFront(), "Expected the third element to be 2 after negative rotation")
}

func TestDequeInsertAtHead(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.Insert(0, 3) // 在索引 0 插入 3

	assert.Equal(t, 3, q.At(0), "Expected At(0) to be 3")
	assert.Equal(t, 1, q.At(1), "Expected At(1) to be 1")
	assert.Equal(t, 2, q.At(2), "Expected At(2) to be 2")
}

func TestDequeInsertAtTail(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.Insert(2, 3) // 在索引 2 插入 3

	assert.Equal(t, 1, q.At(0), "Expected At(0) to be 1")
	assert.Equal(t, 2, q.At(1), "Expected At(1) to be 2")
	assert.Equal(t, 3, q.At(2), "Expected At(2) to be 3")
}

func TestDequeIter(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	count := 0
	iter := q.Iter()
	iter(func(item interface{}) bool {
		assert.Equal(t, count+1, item, "Expected item to be in iter sequential order")
		count++
		return true
	})
	assert.Equal(t, 3, count, "Expected to iterate over 3 items")
}

func TestDequeRIter(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	count := 3
	rIter := q.RIter() // 获取反向迭代器
	rIter(func(item interface{}) bool {
		assert.Equal(t, count, item, "Expected item to be in iter reverse sequential order")
		count--
		return true
	})
	assert.Equal(t, 0, count, "Expected to iterate over 3 items in reverse order")
}

func TestDequePopFrontEmpty(t *testing.T) {
	q := NewDeque()
	assert.Panics(t, func() { q.PopFront() }, "Expected panic when popping from an empty deque")
}

func TestDequePopBackEmpty(t *testing.T) {
	q := NewDeque()
	assert.Panics(t, func() { q.PopBack() }, "Expected panic when popping from an empty deque")
}

func TestDequeIterPopFront(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	// 使用 IterPopFront 迭代并移除元素
	iter := q.IterPopFront()
	count := 0
	iter(func(item interface{}) bool {
		assert.Equal(t, count+1, item, "Expected item to be in sequential order")
		count++
		return true
	})

	assert.Equal(t, 3, count, "Expected to iterate over 3 items")
	assert.Equal(t, 0, q.Len(), "Expected queue length to be 0 after IterPopFront")
}

func TestDequeIterPopBack(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	// 使用 IterPopBack 迭代并移除元素
	iter := q.IterPopBack()
	count := 3
	iter(func(item interface{}) bool {
		assert.Equal(t, count, item, "Expected item to be in reverse sequential order")
		count--
		return true
	})

	assert.Equal(t, 0, count, "Expected to iterate over 3 items in reverse order")
	assert.Equal(t, 0, q.Len(), "Expected queue length to be 0 after IterPopBack")
}

func TestDequeIterPopFrontPartial(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	// 使用 IterPopFront 迭代并移除部分元素
	iter := q.IterPopFront()
	count := 0
	iter(func(item interface{}) bool {
		assert.Equal(t, count+1, item, "Expected item to be in sequential order")
		count++
		return count != 2 // 停止迭代
	})

	assert.Equal(t, 2, count, "Expected to iterate over 2 items")
	assert.Equal(t, 1, q.Len(), "Expected queue length to be 1 after partial IterPopFront")
	assert.Equal(t, 3, q.PopFront(), "Expected remaining item to be 3")
}

func TestDequeIterPopBackPartial(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)

	// 使用 IterPopBack 迭代并移除部分元素
	iter := q.IterPopBack()
	count := 3
	iter(func(item interface{}) bool {
		assert.Equal(t, count, item, "Expected item to be in reverse sequential order")
		count--
		return count != 1 // 停止迭代
	})

	assert.Equal(t, 1, q.Len(), "Expected queue length to be 1 after partial IterPopBack")
	assert.Equal(t, 1, q.PopBack(), "Expected remaining item to be 1")
}
