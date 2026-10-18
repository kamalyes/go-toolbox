/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-08 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 23:04:47
 * @FilePath: \go-toolbox\pkg\random\rand_bench_test.go
 * @Description: 随机数生成基准测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package random

import (
	"math"
	"math/rand"
	"sync"
	"testing"
)

// 性能测试函数
func BenchmarkFRandBytesJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := FRandBytesJSON(1024) // 测试生成1024字节的随机字节字符串
		if err != nil {
			b.Error(err) // 如果有错误，记录
		}
	}
}

// BenchmarkGenerateRandModel 性能测试 GenerateRandModel 函数
func BenchmarkGenerateRandModel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		model := &TestModel{}
		_, _, err := GenerateRandModel(model)
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
	}
}

func BenchmarkRandBytesParallel(b *testing.B) {
	b.Run("FRandBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandBytes(20)
			}
		})
	})
	b.Run("FRandAlphaBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandAlphaBytes(20)
			}
		})
	})
	b.Run("FRandHexBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandHexBytes(20)
			}
		})
	})
	b.Run("FRandDecBytes", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandDecBytes(20)
			}
		})
	})
}

func BenchmarkRandInt(b *testing.B) {
	b.Run("RandInt", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = RandInt(0, i)
		}
	})
	b.Run("FRandInt", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = FRandInt(0, i)
		}
	})
	b.Run("FRandUint32", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = FRandUint32(0, uint32(i))
		}
	})
	b.Run("FastIntn", func(b *testing.B) {
		for i := 1; i < b.N; i++ {
			_ = FastIntn(i)
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
				_ = FRandInt(0, math.MaxInt32)
			}
		})
	})
	b.Run("FRandUint32", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FRandUint32(0, math.MaxInt32)
			}
		})
	})
	b.Run("FastIntn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = FastIntn(math.MaxInt32)
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

// ============================================================================
// 随机字符串与生成器构造
//
// 原先 RandInt / RandFloat / RandString 以及 duration 里一共十六处都走 NewRand()，
// 而 NewRand 每次都要构造一整套生成器：rand.NewSource 分配 607 项 int64 表
// （约 4.9KB）并跑完播种预热，一次 7.4 微秒RandString 还要在循环里取 n 个数，
// UUID 调五次 RandString，于是一个 UUID 要三十多微秒
// 这组基准守住改用共享快速随机源后的结果
//
// [EN] Sixteen call sites used to go through NewRand(), which constructs a whole
// generator per call: rand.NewSource allocates a 607-entry int64 table of about
// 4.9KB and runs the seeding warm-up, measured at 7.4 microseconds. RandString
// drew n numbers in a loop and UUID calls RandString five times, so one UUID
// cost over thirty microseconds. These benchmarks guard the shared fast source
// ============================================================================

var benchStrSink string

func BenchmarkRandString(b *testing.B) {
	for _, c := range []struct {
		name string
		n    int
		mode RandType
	}{
		{"10位/特殊", 10, SPECIAL},
		{"16位/数字", 16, NUMBER},
		{"32位/大小写数字", 32, CAPITAL | LOWERCASE | NUMBER},
		{"64位/全集", 64, CAPITAL | LOWERCASE | SPECIAL | NUMBER},
	} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchStrSink = RandString(c.n, c.mode)
			}
		})
	}
}

func BenchmarkRandStringParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			benchStrSink = RandString(16, CAPITAL|LOWERCASE|NUMBER)
		}
	})
}

func BenchmarkUUIDString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchStrSink = UUID()
	}
}

// BenchmarkNewRandSeeded NewRand 带显式种子时仍要完整构造
//
// 这是 rand.NewSource 播种 607 项表的固有代价，不能优化掉——给了种子就必须
// 可复现，TestNewRand 钉着 NewRand(1).Int63() 的精确值
// 留着这个基准是为了说明：热路径不该调它，而不是它该变快
//
// [EN] NewRand with an explicit seed still constructs fully. That is the
// inherent cost of seeding rand.NewSource's 607-entry table and cannot be
// optimised away: an explicit seed must stay reproducible, and TestNewRand pins
// NewRand(1).Int63() exactly. The benchmark exists to show hot paths should not
// call it, not that it should get faster
func BenchmarkNewRandSeeded(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewRand(1)
	}
}
