/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 15:37:35
 * @FilePath: \go-toolbox\tests\syncx_map_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

func BenchmarkMap_Store(b *testing.B) {
	m := syncx.NewMap[int, string]()
	for i := 0; i < b.N; i++ {
		m.Store(i, "value") // 存储键值对
	}
}

func BenchmarkMap_Load(b *testing.B) {
	m := syncx.NewMap[int, string]()
	for i := 0; i < 1000; i++ {
		m.Store(i, "value") // 预先存储 1000 个键值对
	}

	b.ResetTimer() // 重置计时器以排除设置时间
	for i := 0; i < b.N; i++ {
		m.Load(i % 1000) // 加载键值对
	}
}

func BenchmarkMap_Delete(b *testing.B) {
	m := syncx.NewMap[int, string]()
	for i := 0; i < 1000; i++ {
		m.Store(i, "value") // 预先存储 1000 个键值对
	}

	b.ResetTimer() // 重置计时器以排除设置时间
	for i := 0; i < b.N; i++ {
		m.Delete(i % 1000) // 删除键值对
	}
}

func BenchmarkMap_LoadOrStore(b *testing.B) {
	m := syncx.NewMap[int, string]()
	b.ResetTimer() // 重置计时器以排除设置时间
	for i := 0; i < b.N; i++ {
		m.LoadOrStore(i, "value") // 加载或存储键值对
	}
}

func BenchmarkMap_Range(b *testing.B) {
	m := syncx.NewMap[int, string]()
	for i := 0; i < 1000; i++ {
		m.Store(i, "value") // 预先存储 1000 个键值对
	}

	b.ResetTimer() // 重置计时器以排除设置时间
	for i := 0; i < b.N; i++ {
		m.Range(func(k int, v string) bool {
			return true // 遍历所有键值对
		})
	}
}
