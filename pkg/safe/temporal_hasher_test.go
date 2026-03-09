/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-03-09 17:18:05
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-03-09 22:15:33
 * @FilePath: \go-toolbox\pkg\safe\temporal_hasher_test.go
 * @Description: 临时哈希生成器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemporalHasher_Hash(t *testing.T) {
	hasher := NewTemporalHasher(nil)

	// 相同条件生成相同哈希
	hash1 := hasher.Hash("user123", "device456", "iOS")
	hash2 := hasher.Hash("user123", "device456", "iOS")
	assert.Equal(t, hash1, hash2)
	assert.Equal(t, 12, len(hash1))

	// 不同参数生成不同哈希
	hash3 := hasher.Hash("user456", "device456", "iOS")
	assert.NotEqual(t, hash1, hash3)

	// 参数顺序不影响结果（内部会排序）
	hash4 := hasher.Hash("iOS", "device456", "user123")
	assert.Equal(t, hash1, hash4)
}

func TestTemporalHasher_HashMap(t *testing.T) {
	hasher := NewTemporalHasher(nil)

	kvMap1 := map[string]string{
		"userId":   "user123",
		"deviceId": "device456",
		"platform": "iOS",
	}

	kvMap2 := map[string]string{
		"platform": "iOS",
		"userId":   "user123",
		"deviceId": "device456",
	}

	// map 顺序不同，但内容相同
	hash1 := hasher.HashMap(kvMap1)
	hash2 := hasher.HashMap(kvMap2)
	assert.Equal(t, hash1, hash2)

	// 不同内容生成不同哈希
	kvMap3 := map[string]string{
		"userId":   "user456",
		"deviceId": "device456",
		"platform": "iOS",
	}
	hash3 := hasher.HashMap(kvMap3)
	assert.NotEqual(t, hash1, hash3)
}

func TestTemporalHasher_TimeWindow(t *testing.T) {
	tests := []struct {
		name   string
		window time.Duration
	}{
		{"30秒窗口", 30 * time.Second},
		{"1分钟窗口", 1 * time.Minute},
		{"5分钟窗口", 5 * time.Minute},
		{"1小时窗口", 1 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewTemporalHasher(
				WithWindow(tt.window),
				WithLength(12),
			)

			// 使用窗口对齐的时间点（从 Unix 纪元开始的整数倍窗口）
			baseTime := time.Date(2025, 3, 9, 0, 0, 0, 0, time.UTC)

			// 窗口内第一个时间点
			hash1 := hasher.HashAt(baseTime, "user123", "device456")

			// 窗口内中间时间点（窗口的一半）
			midTime := baseTime.Add(tt.window / 2)
			hash2 := hasher.HashAt(midTime, "user123", "device456")
			assert.Equal(t, hash1, hash2, "窗口内中间时间点应该生成相同哈希")

			// 窗口边界前1秒
			beforeBoundary := baseTime.Add(tt.window - time.Second)
			hash3 := hasher.HashAt(beforeBoundary, "user123", "device456")
			assert.Equal(t, hash1, hash3, "窗口边界前应该生成相同哈希")

			// 窗口边界（下一个窗口开始）
			nextWindow := baseTime.Add(tt.window)
			hash4 := hasher.HashAt(nextWindow, "user123", "device456")
			assert.NotEqual(t, hash1, hash4, "下一个窗口应该生成不同哈希")

			// 下一个窗口内的时间点
			nextWindowMid := nextWindow.Add(tt.window / 2)
			hash5 := hasher.HashAt(nextWindowMid, "user123", "device456")
			assert.Equal(t, hash4, hash5, "同一窗口内应该生成相同哈希")
			assert.NotEqual(t, hash1, hash5, "不同窗口应该生成不同哈希")
		})
	}
}

func TestTemporalHasher_CustomConfig(t *testing.T) {
	// 自定义配置：10 分钟窗口，16 字符长度
	hasher := NewTemporalHasher(
		WithWindow(10*time.Minute),
		WithLength(16),
		WithSeparator(":"),
	)

	hash := hasher.Hash("user123", "device456")
	assert.Equal(t, 16, len(hash))
	assert.Equal(t, 10*time.Minute, hasher.Window())
	assert.Equal(t, 16, hasher.Length())
}

func TestTemporalHasher_EmptyParts(t *testing.T) {
	hasher := NewTemporalHasher(nil)

	// 无参数
	hash1 := hasher.Hash()
	hash2 := hasher.Hash()
	assert.Equal(t, hash1, hash2)

	// 有参数
	hash3 := hasher.Hash("user123")
	assert.NotEqual(t, hash1, hash3)
}

func TestTemporalHasher_WebSocketClientID(t *testing.T) {
	// 模拟 WebSocket ClientID 生成场景
	hasher := NewTemporalHasher(
		WithWindow(5*time.Minute),
		WithLength(12),
	)

	// 场景1：同一用户在不同设备登录
	device1Hash := hasher.Hash("user123", "device1")
	device2Hash := hasher.Hash("user123", "device2")
	device3Hash := hasher.Hash("user123", "device3")

	assert.NotEqual(t, device1Hash, device2Hash)
	assert.NotEqual(t, device1Hash, device3Hash)
	assert.NotEqual(t, device2Hash, device3Hash)

	// 场景2：同一设备登录不同用户
	user1Hash := hasher.Hash("user1", "device123")
	user2Hash := hasher.Hash("user2", "device123")
	user3Hash := hasher.Hash("user3", "device123")

	assert.NotEqual(t, user1Hash, user2Hash)
	assert.NotEqual(t, user1Hash, user3Hash)
	assert.NotEqual(t, user2Hash, user3Hash)

	// 场景3：同一用户同一设备短期内重连（复用 ClientID）
	t1 := time.Date(2025, 3, 9, 0, 0, 0, 0, time.UTC)
	hash1 := hasher.HashAt(t1, "user123", "device456")

	t2 := t1.Add(2 * time.Minute) // 2 分钟后重连
	hash2 := hasher.HashAt(t2, "user123", "device456")
	assert.Equal(t, hash1, hash2) // 复用相同的 ClientID

	// 场景4：超过时间窗口后重连（生成新 ClientID）
	t3 := t1.Add(6 * time.Minute) // 6 分钟后重连
	hash3 := hasher.HashAt(t3, "user123", "device456")
	assert.NotEqual(t, hash1, hash3) // 生成新的 ClientID
}

func TestTemporalHasher_WithHeaders(t *testing.T) {
	hasher := NewTemporalHasher()

	// 使用 map 模拟 headers
	headers1 := map[string]string{
		"userId":     "user123",
		"deviceId":   "device456",
		"User-Agent": "Mozilla/5.0",
		"Platform":   "iOS",
	}

	headers2 := map[string]string{
		"Platform":   "iOS",
		"userId":     "user123",
		"User-Agent": "Mozilla/5.0",
		"deviceId":   "device456",
	}

	// headers 顺序不同，但内容相同
	hash1 := hasher.HashMap(headers1)
	hash2 := hasher.HashMap(headers2)
	assert.Equal(t, hash1, hash2)

	// 不同 headers 生成不同哈希
	headers3 := map[string]string{
		"userId":     "user123",
		"deviceId":   "device456",
		"User-Agent": "Chrome",
		"Platform":   "Android",
	}
	hash3 := hasher.HashMap(headers3)
	assert.NotEqual(t, hash1, hash3)
}

func TestTemporalHasher_MultipleTimeWindows(t *testing.T) {
	tests := []struct {
		name        string
		window      time.Duration
		windowCount int
	}{
		{"5个30秒窗口", 30 * time.Second, 5},
		{"10个1分钟窗口", 1 * time.Minute, 10},
		{"5个5分钟窗口", 5 * time.Minute, 5},
		{"5个1小时窗口", 1 * time.Hour, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewTemporalHasher(
				WithWindow(tt.window),
				WithLength(12),
			)

			// 使用窗口对齐的时间（从 Unix 纪元开始的整数倍）
			baseTime := time.Date(2025, 3, 9, 0, 0, 0, 0, time.UTC)
			parts := []string{"user123", "device456"}

			// 生成多个时间窗口的哈希
			hashes := make([]string, tt.windowCount)
			for i := range tt.windowCount {
				t := baseTime.Add(time.Duration(i) * tt.window)
				hashes[i] = hasher.HashAt(t, parts...)
			}

			// 验证不同时间窗口生成不同哈希
			for i := range tt.windowCount {
				for j := i + 1; j < tt.windowCount; j++ {
					assert.NotEqual(t, hashes[i], hashes[j],
						"窗口%d和窗口%d应该生成不同哈希", i, j)
				}
			}

			// 验证每个窗口内的一致性
			for i := range tt.windowCount {
				windowStart := baseTime.Add(time.Duration(i) * tt.window)
				windowMid := windowStart.Add(tt.window / 2)
				windowEnd := windowStart.Add(tt.window - time.Second)

				hashStart := hasher.HashAt(windowStart, parts...)
				hashMid := hasher.HashAt(windowMid, parts...)
				hashEnd := hasher.HashAt(windowEnd, parts...)

				assert.Equal(t, hashStart, hashMid, "窗口%d内中间时间应该一致", i)
				assert.Equal(t, hashStart, hashEnd, "窗口%d内结束时间应该一致", i)
			}
		})
	}
}
