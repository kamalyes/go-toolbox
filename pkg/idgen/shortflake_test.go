/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:30:00
 * @FilePath: \go-toolbox\pkg\idgen\shortflake_test.go
 * @Description: ShortFlake 生成器测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"regexp"
	"strconv"
	"sync"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/assert"
)

// TestShortFlakeGenerator 测试短 Snowflake 生成器
func TestShortFlakeGenerator(t *testing.T) {
	gen := NewShortFlakeGenerator(1)

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(traceID, "TraceID 不应为空")
		assert.True(regexp.MustCompile(`^\d+$`).MatchString(traceID), "TraceID 应为纯数字")
		
		// 验证长度（53位最多16位数字）
		id, _ := strconv.ParseInt(traceID, 10, 64)
		assert.True(id > 0, "ID 应为正数")
		assert.True(id < 9007199254740992, "ID 应小于 2^53")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(spanID, "SpanID 不应为空")
		assert.True(regexp.MustCompile(`^\d+$`).MatchString(spanID), "SpanID 应为纯数字")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(requestID, "RequestID 不应为空")
		assert.True(regexp.MustCompile(`^\d+$`).MatchString(requestID), "RequestID 应为纯数字")
	})

	t.Run("Monotonic", func(t *testing.T) {
		var lastID int64
		for i := 0; i < 1000; i++ {
			id := gen.Generate()
			assert.True(id > lastID, "ShortFlake ID 应单调递增")
			lastID = id
		}
	})

	t.Run("Uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 10000; i++ {
			id := gen.GenerateTraceID()
			assert.False(ids[id], "生成的 ID 应唯一")
			ids[id] = true
		}
	})
}

// TestShortFlakeBase62Generator 测试 Base62 编码生成器
func TestShortFlakeBase62Generator(t *testing.T) {
	gen := NewShortFlakeBase62Generator(1)

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(traceID, "TraceID 不应为空")
		assert.True(len(traceID) >= 9 && len(traceID) <= 10, "TraceID 应为 9-10 字符")
		assert.True(regexp.MustCompile(`^[0-9A-Za-z]+$`).MatchString(traceID), "TraceID 应为 Base62 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(spanID, "SpanID 不应为空")
		assert.True(len(spanID) >= 9 && len(spanID) <= 10, "SpanID 应为 9-10 字符")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(requestID, "RequestID 不应为空")
		assert.True(len(requestID) >= 9 && len(requestID) <= 10, "RequestID 应为 9-10 字符")
	})

	t.Run("Uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 10000; i++ {
			id := gen.GenerateTraceID()
			assert.False(ids[id], "生成的 ID 应唯一")
			ids[id] = true
		}
	})
}

// TestShortFlakeConcurrent 测试并发生成
func TestShortFlakeConcurrent(t *testing.T) {
	gen := NewShortFlakeGenerator(1)
	
	var wg sync.WaitGroup
	ids := make(map[int64]bool)
	mu := sync.Mutex{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := gen.Generate()
			mu.Lock()
			assert.False(ids[id], "并发生成的 ID 应唯一")
			ids[id] = true
			mu.Unlock()
		}()
	}

	wg.Wait()
	assert.Equal(100, len(ids), "应生成 100 个唯一 ID")
}

// BenchmarkShortFlakeGenerator 基准测试
func BenchmarkShortFlakeGenerator(b *testing.B) {
	gen := NewShortFlakeGenerator(1)
	b.Run("Generate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.Generate()
		}
	})
}

// BenchmarkShortFlakeBase62Generator 基准测试 - Base62
func BenchmarkShortFlakeBase62Generator(b *testing.B) {
	gen := NewShortFlakeBase62Generator(1)
	b.Run("GenerateTraceID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateTraceID()
		}
	})
}
