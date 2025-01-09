/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:11:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 20:20:55
 * @FilePath: \go-toolbox\tests\context_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"context"
	"strconv"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/contextx"
)

// 基准测试：并发设置值
func BenchmarkConcurrentSet(b *testing.B) {
	parentCtx := context.Background()
	customCtx := contextx.NewContext(parentCtx, nil)

	// 使用 RunParallel 进行并发基准测试
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 为每个 goroutine 生成唯一的键和值
			key := "testKey" + strconv.Itoa(b.N) // 使用 b.N 作为唯一标识
			value := "testValue" + strconv.Itoa(b.N)
			customCtx.Set(key, value) // 并发设置值
		}
	})
}

// 基准测试：并发获取值
func BenchmarkConcurrentGetValue(b *testing.B) {
	parentCtx := context.Background()
	customCtx := contextx.NewContext(parentCtx, nil)

	// 预先填充一些值到上下文中
	const numValues = 100
	for i := 0; i < numValues; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		customCtx.Set(key, value) // 设置初始值
	}

	// 使用 RunParallel 进行并发基准测试
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 随机选择一个键
			key := "key" + strconv.Itoa(int(b.N%numValues))             // 随机选择一个键
			expectedValue := "value" + strconv.Itoa(int(b.N%numValues)) // 预期的值
			actualValue := customCtx.Value(key)                         // 并发获取值

			// 断言获取的值是否正确
			if actualValue != expectedValue {
				b.Errorf("expected %s, got %s", expectedValue, actualValue)
			}
		}
	})
}

// 基准测试：并发删除键
func BenchmarkConcurrentDeleteKey(b *testing.B) {
	parentCtx := context.Background()
	customCtx := contextx.NewContext(parentCtx, nil)

	// 预先填充一些值到上下文中
	const numValues = 100
	for i := 0; i < numValues; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		customCtx.Set(key, value) // 设置初始值
	}

	// 使用 RunParallel 进行并发基准测试
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "key" + strconv.Itoa(int(b.N%numValues)) // 随机选择一个键
			customCtx.Remove(key)                           // 并发删除键
		}
	})
}
