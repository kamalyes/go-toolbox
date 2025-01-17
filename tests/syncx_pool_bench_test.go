/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 13:32:15
 * @FilePath: \go-toolbox\tests\syncx_pool_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// BenchmarkLimitedPool 测试 LimitedPool 的性能
func BenchmarkLimitedPool(b *testing.B) {
	const numAllocations = 100000
	pool := syncx.NewLimitedPool(32, 1024) // 初始化池，最小32字节，最大1024字节

	b.ResetTimer() // 重置计时器，确保不包括初始化时间

	for i := 0; i < b.N; i++ {
		buf := pool.Get(64) // 从池中获取一个大小为64字节的切片
		if buf == nil {
			b.Fatal("Get returned nil buffer")
		}
		pool.Put(buf) // 将切片放回池中
	}
}
