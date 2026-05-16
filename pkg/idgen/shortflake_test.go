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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestShortFlakeGenerator 测试短 Snowflake 生成器
func TestShortFlakeGenerator(t *testing.T) {
	gen := NewShortFlakeGenerator(1)

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "TraceID 不应为空")
		assert.Equal(t, 13, len(traceID), "TraceID 应为 13 字符 hex")
		assert.True(t, regexp.MustCompile(`^[0-9a-f]{13}$`).MatchString(traceID), "TraceID 应为 hex 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(t, spanID, "SpanID 不应为空")
		assert.Equal(t, 8, len(spanID), "SpanID 应为 8 字符 hex")
		assert.True(t, regexp.MustCompile(`^[0-9a-f]{8}$`).MatchString(spanID), "SpanID 应为 hex 格式")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(t, requestID, "RequestID 不应为空")
		assert.True(t, regexp.MustCompile(`^\d+-\d+$`).MatchString(requestID), "RequestID 应为 数字-计数器 格式")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(t, correlationID, "CorrelationID 不应为空")
		assert.Equal(t, 36, len(correlationID), "CorrelationID 应为 36 字符 UUID 格式")
		assert.True(t, regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`).MatchString(correlationID), "CorrelationID 应为 UUID 格式")
	})

	t.Run("DifferentFormats", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		spanID := gen.GenerateSpanID()
		requestID := gen.GenerateRequestID()
		correlationID := gen.GenerateCorrelationID()

		assert.NotEqual(t, traceID, spanID, "TraceID(hex 13字符) 和 SpanID(hex 8字符) 格式应不同")
		assert.NotEqual(t, traceID, requestID, "TraceID 和 RequestID 格式应不同")
		assert.NotEqual(t, traceID, correlationID, "TraceID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, spanID, requestID, "SpanID 和 RequestID 格式应不同")
		assert.NotEqual(t, spanID, correlationID, "SpanID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, requestID, correlationID, "RequestID 和 CorrelationID 格式应不同")
	})

	t.Run("Monotonic", func(t *testing.T) {
		var lastID int64
		for i := 0; i < 1000; i++ {
			id := gen.Generate()
			assert.True(t, id > lastID, "ShortFlake ID 应单调递增")
			lastID = id
		}
	})

	t.Run("Uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 10000; i++ {
			id := gen.GenerateTraceID()
			assert.False(t, ids[id], "生成的 ID 应唯一")
			ids[id] = true
		}
	})
}

func TestShortFlakeBase62Generator(t *testing.T) {
	gen := NewShortFlakeBase62Generator(1)

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "TraceID 不应为空")
		assert.True(t, len(traceID) >= 9 && len(traceID) <= 10, "TraceID 应为 9-10 字符")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z]+$`).MatchString(traceID), "TraceID 应为 Base62 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(t, spanID, "SpanID 不应为空")
		assert.True(t, len(spanID) >= 1 && len(spanID) <= 7, "SpanID 应为 1-7 字符 Base62")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z]+$`).MatchString(spanID), "SpanID 应为 Base62 格式")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(t, requestID, "RequestID 不应为空")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z]+-\d+$`).MatchString(requestID), "RequestID 应为 Base62前缀-计数器 格式")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(t, correlationID, "CorrelationID 不应为空")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z]+-[0-9A-Za-z]+$`).MatchString(correlationID), "CorrelationID 应为 Base62-Base62 格式")
	})

	t.Run("DifferentFormats", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		spanID := gen.GenerateSpanID()
		requestID := gen.GenerateRequestID()
		correlationID := gen.GenerateCorrelationID()

		assert.NotEqual(t, traceID, spanID, "TraceID 和 SpanID 格式应不同")
		assert.NotEqual(t, traceID, requestID, "TraceID 和 RequestID 格式应不同")
		assert.NotEqual(t, traceID, correlationID, "TraceID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, spanID, requestID, "SpanID 和 RequestID 格式应不同")
		assert.NotEqual(t, spanID, correlationID, "SpanID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, requestID, correlationID, "RequestID 和 CorrelationID 格式应不同")
	})

	t.Run("Uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 10000; i++ {
			id := gen.GenerateTraceID()
			assert.False(t, ids[id], "生成的 ID 应唯一")
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
			assert.False(t, ids[id], "并发生成的 ID 应唯一")
			ids[id] = true
			mu.Unlock()
		}()
	}

	wg.Wait()
	assert.Equal(t, 100, len(ids), "应生成 100 个唯一 ID")
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
