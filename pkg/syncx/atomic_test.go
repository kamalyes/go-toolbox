/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-05 15:08:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:56:39
 * @FilePath: \go-toolbox\pkg\syncx\atomic_test.go
 * @Description: atomic 原子操作单元测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicBool(t *testing.T) {
	b := NewBool(false)

	// 测试初始值
	assert.Equal(t, false, b.Load(), "expected false")

	// 测试存储值
	b.Store(true)
	assert.Equal(t, true, b.Load(), "expected true")

	// 测试切换值
	old := b.Toggle()
	assert.Equal(t, true, old, "expected true")
	assert.Equal(t, false, b.Load(), "expected false")
}

func TestAtomicInt32(t *testing.T) {
	i32 := NewInt32(10)

	// 测试初始值
	assert.Equal(t, int32(10), i32.Load(), "expected 10")

	// 测试增加
	i32.Add(5)
	assert.Equal(t, int32(15), i32.Load(), "expected 15")

	// 测试减少
	i32.Sub(3)
	assert.Equal(t, int32(12), i32.Load(), "expected 12")

	// 测试交换
	old := i32.Swap(20)
	assert.Equal(t, int32(12), old, "expected 12")
	assert.Equal(t, int32(20), i32.Load(), "expected 20")

	// 测试比较交换
	assert.True(t, i32.CAS(20, 30), "expected CAS to succeed")
	assert.Equal(t, int32(30), i32.Load(), "expected 30")
}

func TestAtomicUint32(t *testing.T) {
	u32 := NewUint32(10)

	// 测试初始值
	assert.Equal(t, uint32(10), u32.Load(), "expected 10")

	// 测试增加
	u32.Add(5)
	assert.Equal(t, uint32(15), u32.Load(), "expected 15")

	// 测试减少
	u32.Sub(3)
	assert.Equal(t, uint32(12), u32.Load(), "expected 12")

	// 测试交换
	old := u32.Swap(uint32(20))
	assert.Equal(t, uint32(12), old, "expected 12")
	assert.Equal(t, uint32(20), u32.Load(), "expected 20")

	// 测试比较交换
	assert.True(t, u32.CAS(20, 30), "expected CAS to succeed")
	assert.Equal(t, uint32(30), u32.Load(), "expected 30")
}

func TestAtomicInt64(t *testing.T) {
	i64 := NewInt64(10)

	// 测试初始值
	assert.Equal(t, int64(10), i64.Load(), "expected 10")

	// 测试增加
	i64.Add(5)
	assert.Equal(t, int64(15), i64.Load(), "expected 15")

	// 测试减少
	i64.Sub(3)
	assert.Equal(t, int64(12), i64.Load(), "expected 12")

	// 测试交换
	old := i64.Swap(20)
	assert.Equal(t, int64(12), old, "expected 12")
	assert.Equal(t, int64(20), i64.Load(), "expected 20")

	// 测试比较交换
	assert.True(t, i64.CAS(20, 30), "expected CAS to succeed")
	assert.Equal(t, int64(30), i64.Load(), "expected 30")
}

func TestAtomicUint64(t *testing.T) {
	u64 := NewUint64(10)

	// 测试初始值
	assert.Equal(t, uint64(10), u64.Load(), "expected 10")

	// 测试增加
	u64.Add(5)
	assert.Equal(t, uint64(15), u64.Load(), "expected 15")

	// 测试减少
	u64.Sub(3)
	assert.Equal(t, uint64(12), u64.Load(), "expected 12")

	// 测试交换
	old := u64.Swap(20)
	assert.Equal(t, uint64(12), old, "expected 12")
	assert.Equal(t, uint64(20), u64.Load(), "expected 20")

	// 测试比较交换
	assert.True(t, u64.CAS(20, 30), "expected CAS to succeed")
	assert.Equal(t, uint64(30), u64.Load(), "expected 30")
}

func TestAtomicUintptr(t *testing.T) {
	ptr := NewUintptr(10)

	// 测试初始值
	assert.Equal(t, uintptr(10), ptr.Load(), "expected 10")

	// 测试交换
	old := ptr.Swap(20)
	assert.Equal(t, uintptr(10), old, "expected 10")
	assert.Equal(t, uintptr(20), ptr.Load(), "expected 20")

	// 测试比较交换
	assert.True(t, ptr.CAS(20, 30), "expected CAS to succeed")
	assert.Equal(t, uintptr(30), ptr.Load(), "expected 30")
}

func TestConcurrentAtomicBool(t *testing.T) {
	b := NewBool(false)
	var wg sync.WaitGroup

	// 启动多个 goroutine 进行并发写入和读取
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				b.Store(true)
			} else {
				b.Store(false)
			}
		}(i)
	}

	wg.Wait()
	// 读取最终值
	finalValue := b.Load()
	assert.True(t, finalValue == true || finalValue == false, "final value should be either true or false")
}

func TestConcurrentAtomicInt32(t *testing.T) {
	i32 := NewInt32(0)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				i32.Add(1)
			} else {
				i32.Sub(1)
			}
		}(i)
	}

	wg.Wait()
	assert.Equal(t, int32(0), i32.Load(), "expected final value to be 0")
}

func TestConcurrentAtomicUint32(t *testing.T) {
	u32 := NewUint32(0)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				u32.Add(1)
			} else {
				u32.Sub(1)
			}
		}(i)
	}

	wg.Wait()
	assert.Equal(t, uint32(0), u32.Load(), "expected final value to be 0")
}

func TestConcurrentAtomicInt64(t *testing.T) {
	i64 := NewInt64(0)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				i64.Add(1)
			} else {
				i64.Sub(1)
			}
		}(i)
	}

	wg.Wait()
	assert.Equal(t, int64(0), i64.Load(), "expected final value to be 0")
}

func TestConcurrentAtomicUint64(t *testing.T) {
	u64 := NewUint64(0)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				u64.Add(1)
			} else {
				u64.Sub(1)
			}
		}(i)
	}

	wg.Wait()
	assert.Equal(t, uint64(0), u64.Load(), "expected final value to be 0")
}

func TestAtomicValue_Int(t *testing.T) {
	av := NewAtomicValue(123)
	assert.Equal(t, 123, av.Load())

	av.Store(456)
	assert.Equal(t, 456, av.Load())

	old := av.Swap(789)
	assert.Equal(t, 456, old)
	assert.Equal(t, 789, av.Load())

	ok := av.CompareAndSwap(789, 1000)
	assert.True(t, ok)
	assert.Equal(t, 1000, av.Load())

	ok2 := av.CompareAndSwap(789, 2000)
	assert.False(t, ok2)
	assert.Equal(t, 1000, av.Load())
}

func TestAtomicValue_String(t *testing.T) {
	av := NewAtomicValue("foo")
	assert.Equal(t, "foo", av.Load())
	av.Store("bar")
	assert.Equal(t, "bar", av.Load())
	old := av.Swap("baz")
	assert.Equal(t, "bar", old)
	assert.Equal(t, "baz", av.Load())
	ok := av.CompareAndSwap("baz", "qux")
	assert.True(t, ok)
	assert.Equal(t, "qux", av.Load())
}
