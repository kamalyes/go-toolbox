/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 15:01:18
 * @FilePath: \go-toolbox\pkg\syncx\pool_test.go
 * @Description: LimitedPool 池单元测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitedPool(t *testing.T) {
	// 创建一个 LimitedPool，最小大小为 16，最大大小为 64
	pool := NewLimitedPool(16, 64)

	// 测试获取和释放字节切片
	t.Run("TestGetAndPut", func(t *testing.T) {
		// 获取一个大小为 32 的字节切片
		buf := pool.Get(32)
		assert.NotNil(t, buf)
		assert.Equal(t, 32, len(*buf))

		// 将字节切片放回池中
		pool.Put(buf)

		// 再次获取，应该能复用之前的字节切片
		buf2 := pool.Get(32)
		assert.NotNil(t, buf2)
		assert.Equal(t, 32, len(*buf2))
		assert.Equal(t, buf, buf2, "Expected to get the same slice from the pool")
	})

	// 测试边界条件
	t.Run("TestBoundaryConditions", func(t *testing.T) {
		// 测试最小大小
		bufMin := pool.Get(16)
		assert.NotNil(t, bufMin)
		assert.Equal(t, 16, len(*bufMin))
		pool.Put(bufMin)

		// 测试最大大小
		bufMax := pool.Get(64)
		assert.NotNil(t, bufMax)
		assert.Equal(t, 64, len(*bufMax))
		pool.Put(bufMax)

		// 测试超出最大大小的请求
		bufTooLarge := pool.Get(128)
		assert.Nil(t, bufTooLarge, "Expected nil for size > maxSize")
	})

	// 测试不同大小的请求
	t.Run("TestDifferentSizes", func(t *testing.T) {
		sizes := []int{16, 32, 48, 64}
		for _, size := range sizes {
			buf := pool.Get(size)
			if buf == nil {
				continue // 继续下一个测试
			}

			assert.Equal(t, size, len(*buf), "Expected buffer length to be %d", size)
			pool.Put(buf)
		}
	})
	// 测试获取未放回的切片
	t.Run("TestUnreturnedSlice", func(t *testing.T) {
		buf1 := pool.Get(16)
		buf2 := pool.Get(32)

		// buf1 和 buf2 应该是不同的切片
		assert.NotEqual(t, buf1, buf2, "Expected different slices for consecutive Get calls")

		pool.Put(buf1)
		pool.Put(buf2)
	})
}

// Test for generic Pool[T]
func TestGenericPool(t *testing.T) {
	// 创建一个 int 类型的泛型池
	pool := NewPool(func() int { return 42 })

	// Get 应该返回 new() 的值
	v := pool.Get()
	assert.Equal(t, 42, v)

	// Put 一个新值，再次 Get 应该能拿到刚放进去的值
	pool.Put(100)
	v2 := pool.Get()
	assert.Equal(t, 100, v2)

	// Put nil 值（对于 int 类型无意义，但测试泛型池的健壮性）
	// 对于非指针类型，nil 不适用，这里仅做类型安全性验证
	// pool.Put(nil) // 编译期会报错，说明泛型池类型安全

	// 测试 string 类型
	strPool := NewPool(func() string { return "hello" })
	s := strPool.Get()
	assert.Equal(t, "hello", s)
	strPool.Put("world")
	s2 := strPool.Get()
	assert.Equal(t, "world", s2)
}
