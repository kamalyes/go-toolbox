/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-03-09 17:20:11
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-03-09 17:20:11
 * @FilePath: \go-toolbox\pkg\safe\temporal_hasher_bench_test.go
 * @Description: 临时哈希生成器性能测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"testing"
	"time"
)

var (
	benchTemporalHashResult string
)

func BenchmarkTemporalHasher_Hash(b *testing.B) {
	hasher := NewTemporalHasher(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchTemporalHashResult = hasher.Hash("user123", "device456", "iOS")
	}
}

func BenchmarkTemporalHasher_HashMap(b *testing.B) {
	hasher := NewTemporalHasher(nil)
	kvMap := map[string]string{
		"userId":   "user123",
		"deviceId": "device456",
		"platform": "iOS",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchTemporalHashResult = hasher.HashMap(kvMap)
	}
}

func BenchmarkTemporalHasher_HashAt(b *testing.B) {
	hasher := NewTemporalHasher(nil)
	t := time.Now()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchTemporalHashResult = hasher.HashAt(t, "user123", "device456", "iOS")
	}
}

func BenchmarkTemporalHasher_DifferentPartsCount(b *testing.B) {
	hasher := NewTemporalHasher(nil)

	tests := []struct {
		name  string
		parts []string
	}{
		{"2个参数", []string{"user123", "device456"}},
		{"3个参数", []string{"user123", "device456", "iOS"}},
		{"5个参数", []string{"user123", "device456", "iOS", "Mozilla/5.0", "1.0.0"}},
		{"10个参数", []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchTemporalHashResult = hasher.Hash(tt.parts...)
			}
		})
	}
}

func BenchmarkTemporalHasher_DifferentHashLength(b *testing.B) {
	tests := []struct {
		name   string
		length int
	}{
		{"8字符", 8},
		{"12字符", 12},
		{"16字符", 16},
		{"32字符", 32},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			hasher := NewTemporalHasher(
				WithWindow(5*time.Minute),
				WithLength(tt.length),
			)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchTemporalHashResult = hasher.Hash("user123", "device456", "iOS")
			}
		})
	}
}

func BenchmarkTemporalHasher_DifferentMapSize(b *testing.B) {
	tests := []struct {
		name  string
		kvMap map[string]string
	}{
		{
			"2个键值对",
			map[string]string{
				"userId":   "user123",
				"deviceId": "device456",
			},
		},
		{
			"5个键值对",
			map[string]string{
				"userId":     "user123",
				"deviceId":   "device456",
				"platform":   "iOS",
				"User-Agent": "Mozilla/5.0",
				"version":    "1.0.0",
			},
		},
		{
			"10个键值对",
			map[string]string{
				"userId":     "user123",
				"deviceId":   "device456",
				"platform":   "iOS",
				"User-Agent": "Mozilla/5.0",
				"version":    "1.0.0",
				"language":   "zh-CN",
				"timezone":   "Asia/Shanghai",
				"network":    "WiFi",
				"carrier":    "China Mobile",
				"model":      "iPhone 13",
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			hasher := NewTemporalHasher(nil)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchTemporalHashResult = hasher.HashMap(tt.kvMap)
			}
		})
	}
}

func BenchmarkTemporalHasher_IsExpired(b *testing.B) {
	hasher := NewTemporalHasher(nil)
	hash := hasher.Hash("user123", "device456", "iOS")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = hasher.IsExpired(hash, "user123", "device456", "iOS")
	}
}

func BenchmarkTemporalHasher_Parallel(b *testing.B) {
	hasher := NewTemporalHasher(nil)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchTemporalHashResult = hasher.Hash("user123", "device456", "iOS")
		}
	})
}

func BenchmarkTemporalHasher_ParallelWithMap(b *testing.B) {
	hasher := NewTemporalHasher(nil)
	kvMap := map[string]string{
		"userId":   "user123",
		"deviceId": "device456",
		"platform": "iOS",
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			benchTemporalHashResult = hasher.HashMap(kvMap)
		}
	})
}
