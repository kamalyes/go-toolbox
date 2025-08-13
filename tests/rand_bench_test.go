/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 13:54:40
 * @FilePath: \go-toolbox\tests\rand_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"math"
	"math/rand"
	"sync"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/random"
)

// 性能测试函数
func BenchmarkFRandBytesJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := random.FRandBytesJSON(1024) // 测试生成1024字节的随机字节字符串
		if err != nil {
			b.Error(err) // 如果有错误，记录
		}
	}
}

// BenchmarkGenerateRandModel 性能测试 GenerateRandModel 函数
func BenchmarkGenerateRandModel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := &TestModel{}
		_, _, err := random.GenerateRandModel(model)
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
	}
}

func BenchmarkRandBytesParallel(b *testing.B) {
	b.Run("FRandBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandBytes(20)
			}
		})
	})
	b.Run("FRandAlphaBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandAlphaBytes(20)
			}
		})
	})
	b.Run("FRandHexBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandHexBytes(20)
			}
		})
	})
	b.Run("FRandDecBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandDecBytes(20)
			}
		})
	})
}

func BenchmarkRandInt(b *testing.B) {
	b.Run("RandInt", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = random.RandInt(0, i)
		}
	})
	b.Run("FRandInt", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = random.FRandInt(0, i)
		}
	})
	b.Run("FRandUint32", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = random.FRandUint32(0, uint32(i))
		}
	})
	b.Run("FastIntn", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = random.FastIntn(i)
		}
	})
	b.Run("std.rand.Intn", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = rand.Intn(i)
		}
	})
}

func BenchmarkRandInt32Parallel(b *testing.B) {
	b.Run("FRandInt", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandInt(0, math.MaxInt32)
			}
		})
	})
	b.Run("FRandUint32", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FRandUint32(0, math.MaxInt32)
			}
		})
	})
	b.Run("FastIntn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = random.FastIntn(math.MaxInt32)
			}
		})
	})
	var mu sync.Mutex
	b.Run("std.rand.Intn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mu.Lock()
				_ = rand.Intn(math.MaxInt32)
				mu.Unlock()
			}
		})
	})
}
