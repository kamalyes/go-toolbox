/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 10:50:50
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:59:39
 * @FilePath: \go-toolbox\pkg\syncx\set_bench_test.go
 * @Description: set 集合基准测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"testing"
)

func BenchmarkSetAdd(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < b.N; i++ {
		s.Add(i) // 向集合中添加元素
	}
}

func BenchmarkSetHas(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i) // 预先添加 1000 个元素
	}

	b.ResetTimer() // 重置计时器以排除设置时间
	for i := 0; i < b.N; i++ {
		s.Has(i % 1000) // 检查元素是否存在
	}
}

func BenchmarkSetDelete(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i) // 预先添加 1000 个元素
	}

	b.ResetTimer() // 重置计时器以排除设置时间
	for i := 0; i < b.N; i++ {
		s.Delete(i % 1000) // 删除元素
	}
}

func BenchmarkSetAddAll(b *testing.B) {
	s := NewSet[int]()
	b.Run("Add 1000", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.AddAll(0, 1, 2, 3, 4, 5, 6, 7, 8, 9) // 添加多个元素
		}
	})
}

func BenchmarkSetDeleteAll(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i) // 预先添加 1000 个元素
	}

	b.ResetTimer() // 重置计时器以排除设置时间
	b.Run("Delete All", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.DeleteAll(0, 1, 2, 3, 4, 5, 6, 7, 8, 9) // 删除多个元素
		}
	})
}
