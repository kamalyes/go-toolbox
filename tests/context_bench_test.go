/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:11:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-21 11:57:55
 * @FilePath: \go-toolbox\tests\context_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"context"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/contextx"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

func BenchmarkNewContext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		pool := syncx.NewLimitedPool(32, 1024)
		_ = contextx.NewContext(ctx, pool)
	}
}

// 基准测试：并发设置字符串值
func BenchmarkWithStringValue(b *testing.B) {
	ctx := context.Background()
	pool := syncx.NewLimitedPool(32, 1024)
	c := contextx.NewContext(ctx, pool)

	b.ResetTimer() // 重置计时器，以排除设置上下文的开销

	for i := 0; i < b.N; i++ {
		key := i // 使用整数作为键
		value := "test string value"
		if err := c.WithValue(key, value); err != nil {
			b.Fatalf("failed to set value: %v", err)
		}
	}
}

// 基准测试：并发设置整数值
func BenchmarkWithIntValue(b *testing.B) {
	ctx := context.Background()
	pool := syncx.NewLimitedPool(32, 1024)
	c := contextx.NewContext(ctx, pool)

	b.ResetTimer() // 重置计时器，以排除设置上下文的开销

	for i := 0; i < b.N; i++ {
		key := i    // 使用整数作为键
		value := 42 // 整数值
		if err := c.WithValue(key, value); err != nil {
			b.Fatalf("failed to set value: %v", err)
		}
	}
}

func BenchmarkWithStructValue(b *testing.B) {
	ctx := context.Background()
	pool := syncx.NewLimitedPool(32, 1024)
	c := contextx.NewContext(ctx, pool)

	b.ResetTimer() // 重置计时器，以排除设置上下文的开销

	for i := 0; i < b.N; i++ {
		key := i                                  // 使用整数作为键
		value := TestStruct{Name: "test", Age: i} // 结构体值
		if err := c.WithValue(key, value); err != nil {
			b.Fatalf("failed to set value: %v", err)
		}
	}
}

// 基准测试：并发设置切片值
func BenchmarkWithSliceValue(b *testing.B) {
	ctx := context.Background()
	pool := syncx.NewLimitedPool(32, 1024)
	c := contextx.NewContext(ctx, pool)

	b.ResetTimer() // 重置计时器，以排除设置上下文的开销

	for i := 0; i < b.N; i++ {
		key := i                            // 使用整数作为键
		value := []byte("test slice value") // 切片值
		if err := c.WithValue(key, value); err != nil {
			b.Fatalf("failed to set value: %v", err)
		}
	}
}

// 基准测试：并发设置空接口值
func BenchmarkWithInterfaceValue(b *testing.B) {
	ctx := context.Background()
	pool := syncx.NewLimitedPool(32, 1024)
	c := contextx.NewContext(ctx, pool)

	b.ResetTimer() // 重置计时器，以排除设置上下文的开销

	for i := 0; i < b.N; i++ {
		key := i                                     // 使用整数作为键
		value := interface{}("test interface value") // 空接口值
		if err := c.WithValue(key, value); err != nil {
			b.Fatalf("failed to set value: %v", err)
		}
	}
}

// 基准测试：合并上下文
func BenchmarkMergeContext(b *testing.B) {
	ctx1 := context.Background()
	ctx2 := context.Background()
	pool := syncx.NewLimitedPool(32, 1024)
	c1 := contextx.NewContext(ctx1, pool)
	c2 := contextx.NewContext(ctx2, pool)

	// 预先设置一些值
	for i := 0; i < 1000; i++ {
		c1.WithValue(i, []byte("value from ctx1"))
		c2.WithValue(i, []byte("value from ctx2"))
	}

	b.ResetTimer() // 重置计时器，以排除设置上下文的开销

	for i := 0; i < b.N; i++ {
		_ = contextx.MergeContext(c1, c2)
	}
}

// 基准测试：并发获取值
func BenchmarkValue(b *testing.B) {
	ctx := context.Background()
	pool := syncx.NewLimitedPool(32, 1024)
	c := contextx.NewContext(ctx, pool)

	// 预先设置一些值
	for i := 0; i < 1000; i++ {
		key := i
		value := []byte("test value")
		c.WithValue(key, value)
	}

	b.ResetTimer() // 重置计时器，以排除设置上下文的开销

	for i := 0; i < b.N; i++ {
		key := i % 1000 // 访问先前设置的值
		_ = c.Value(key)
	}
}

// 基准测试：并发删除键
func BenchmarkRemove(b *testing.B) {
	ctx := context.Background()
	pool := syncx.NewLimitedPool(32, 1024)
	c := contextx.NewContext(ctx, pool)

	// 预先设置一些值
	for i := 0; i < 1000; i++ {
		key := i
		value := []byte("test value")
		c.WithValue(key, value)
	}

	b.ResetTimer() // 重置计时器，以排除设置上下文的开销

	for i := 0; i < b.N; i++ {
		key := i % 1000 // 删除先前设置的值
		c.Remove(key)
	}
}
