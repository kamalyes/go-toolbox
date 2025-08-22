/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-21 18:57:15
 * @FilePath: \go-toolbox\tests\deque_queue_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/queue"
)

func BenchmarkDequePushBack(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
}

func BenchmarkDequePushFront(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < b.N; i++ {
		q.PushFront(i)
	}
}

func BenchmarkDequePopFront(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	for i := 0; i < b.N; i++ {
		q.PopFront()
	}
}

func BenchmarkDequePopBack(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	for i := 0; i < b.N; i++ {
		q.PopBack()
	}
}

func BenchmarkDequeInsert(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	for i := 0; i < b.N; i++ {
		q.Insert(i/2, i) // 在中间插入元素
	}
}

func BenchmarkDequeRotate(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < 1000; i++ { // 初始化 1000 个元素
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	for i := 0; i < b.N; i++ {
		q.Rotate(100) // 每次旋转 100 次
	}
}

func BenchmarkDequeIter(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < 1000; i++ {
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	iter := q.Iter()
	for i := 0; i < b.N; i++ {
		iter(func(item interface{}) bool {
			_ = item // 仅访问元素
			return true
		})
	}
}

func BenchmarkDequeRIter(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < 1000; i++ {
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	rIter := q.RIter()
	for i := 0; i < b.N; i++ {
		rIter(func(item interface{}) bool {
			_ = item // 仅访问元素
			return true
		})
	}
}

func BenchmarkDequeIterPopFront(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < 1000; i++ {
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	for i := 0; i < b.N; i++ {
		iter := q.IterPopFront()
		iter(func(item interface{}) bool {
			return true // 仅迭代
		})
	}
}

func BenchmarkDequeIterPopBack(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < 1000; i++ {
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	for i := 0; i < b.N; i++ {
		iter := q.IterPopBack()
		iter(func(item interface{}) bool {
			return true // 仅迭代
		})
	}
}

func BenchmarkDequeIterPopFrontPartial(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < 1000; i++ {
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	for i := 0; i < b.N; i++ {
		iter := q.IterPopFront()
		count := 0
		iter(func(item interface{}) bool {
			count++
			return count != 2 // 停止迭代
		})
	}
}

func BenchmarkDequeIterPopBackPartial(b *testing.B) {
	q := queue.NewDeque()
	for i := 0; i < 1000; i++ {
		q.PushBack(i)
	}
	b.ResetTimer() // 重置计时器，排除初始化的时间
	for i := 0; i < b.N; i++ {
		iter := q.IterPopBack()
		count := 3
		iter(func(item interface{}) bool {
			count--
			return count != 1 // 停止迭代
		})
	}
}
