/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:10:00
 * @FilePath: \go-toolbox\pkg\idgen\idgen_test.go
 * @Description: ID 生成器测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"github.com/kamalyes/go-toolbox/pkg/assert"
	"regexp"
	"sync"
	"testing"
)

// TestDefaultIDGenerator 测试默认 Hex 生成器
func TestDefaultIDGenerator(t *testing.T) {
	gen := NewDefaultIDGenerator()

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(traceID, "TraceID 不应为空")
		assert.Equal(32, len(traceID), "TraceID 应为 32 字符")
		assert.True(regexp.MustCompile("^[0-9a-f]{32}$").MatchString(traceID), "TraceID 应为 hex 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(spanID, "SpanID 不应为空")
		assert.Equal(16, len(spanID), "SpanID 应为 16 字符")
		assert.True(regexp.MustCompile("^[0-9a-f]{16}$").MatchString(spanID), "SpanID 应为 hex 格式")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(requestID, "RequestID 不应为空")
		assert.True(regexp.MustCompile(`^\d+-\d+$`).MatchString(requestID), "RequestID 应为 timestamp-counter 格式")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(correlationID, "CorrelationID 不应为空")
		assert.Equal(36, len(correlationID), "CorrelationID 应为 36 字符 (UUID 格式)")
		assert.True(regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`).MatchString(correlationID), "CorrelationID 应为 UUID 格式")
	})

	t.Run("Uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 1000; i++ {
			id := gen.GenerateTraceID()
			assert.False(ids[id], "生成的 ID 应唯一")
			ids[id] = true
		}
	})
}

// TestUUIDGenerator 测试 UUID 生成器
func TestUUIDGenerator(t *testing.T) {
	gen := NewUUIDGenerator()

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(traceID, "TraceID 不应为空")
		assert.Equal(36, len(traceID), "TraceID 应为 36 字符")
		assert.True(regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(traceID), "TraceID 应为 UUID v4 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(spanID, "SpanID 不应为空")
		assert.Equal(16, len(spanID), "SpanID 应为 16 字符")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(requestID, "RequestID 不应为空")
		assert.True(regexp.MustCompile(`^[0-9a-f]{8}-\d+$`).MatchString(requestID), "RequestID 应包含 UUID 前缀和计数器")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(correlationID, "CorrelationID 不应为空")
		assert.Equal(36, len(correlationID), "CorrelationID 应为 36 字符")
	})
}

// TestNanoIDGenerator 测试 NanoID 生成器
func TestNanoIDGenerator(t *testing.T) {
	gen := NewNanoIDGenerator()

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(traceID, "TraceID 不应为空")
		assert.Equal(21, len(traceID), "TraceID 应为 21 字符")
		assert.True(regexp.MustCompile(`^[0-9A-Za-z_-]{21}$`).MatchString(traceID), "TraceID 应为 NanoID 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(spanID, "SpanID 不应为空")
		assert.Equal(16, len(spanID), "SpanID 应为 16 字符")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(requestID, "RequestID 不应为空")
		assert.True(regexp.MustCompile(`^[0-9A-Za-z_-]{10}-\d+$`).MatchString(requestID), "RequestID 应包含 NanoID 前缀和计数器")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(correlationID, "CorrelationID 不应为空")
		assert.Equal(21, len(correlationID), "CorrelationID 应为 21 字符")
	})
}

// TestSnowflakeGenerator 测试 Snowflake 生成器
func TestSnowflakeGenerator(t *testing.T) {
	gen := NewSnowflakeGenerator(1, 1)

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(traceID, "TraceID 不应为空")
		assert.True(regexp.MustCompile(`^\d+$`).MatchString(traceID), "TraceID 应为纯数字")
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

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(correlationID, "CorrelationID 不应为空")
		assert.True(regexp.MustCompile(`^\d+$`).MatchString(correlationID), "CorrelationID 应为纯数字")
	})

	t.Run("Monotonic", func(t *testing.T) {
		var lastID int64
		for i := 0; i < 100; i++ {
			id := gen.generateSnowflake()
			assert.True(id > lastID, "Snowflake ID 应单调递增")
			lastID = id
		}
	})
}

// TestULIDGenerator 测试 ULID 生成器
func TestULIDGenerator(t *testing.T) {
	gen := NewULIDGenerator()

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(traceID, "TraceID 不应为空")
		assert.Equal(26, len(traceID), "TraceID 应为 26 字符")
		assert.True(regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`).MatchString(traceID), "TraceID 应为 ULID 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(spanID, "SpanID 不应为空")
		assert.Equal(16, len(spanID), "SpanID 应为 16 字符")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(requestID, "RequestID 不应为空")
		assert.True(regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{10}-\d+$`).MatchString(requestID), "RequestID 应包含 ULID 前缀和计数器")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(correlationID, "CorrelationID 不应为空")
		assert.Equal(26, len(correlationID), "CorrelationID 应为 26 字符")
	})
}

// TestFactory 测试工厂函数
func TestFactory(t *testing.T) {
	t.Run("NewIDGenerator with GeneratorType", func(t *testing.T) {
		gen := NewIDGenerator(GeneratorTypeUUID)
		assert.NotNil(gen, "生成器不应为 nil")
		id := gen.GenerateTraceID()
		assert.True(regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(id), "应生成 UUID v4")
	})

	t.Run("NewIDGenerator with string", func(t *testing.T) {
		gen := NewIDGenerator("nanoid")
		assert.NotNil(gen, "生成器不应为 nil")
		id := gen.GenerateTraceID()
		assert.Equal(21, len(id), "应生成 NanoID")
	})

	t.Run("NewIDGenerator with invalid type", func(t *testing.T) {
		gen := NewIDGenerator(12345)
		assert.NotNil(gen, "生成器不应为 nil")
		id := gen.GenerateTraceID()
		assert.Equal(32, len(id), "应回退到默认生成器")
	})

	t.Run("NewIDGeneratorFromString (deprecated)", func(t *testing.T) {
		gen := NewIDGeneratorFromString("ulid")
		assert.NotNil(gen, "生成器不应为 nil")
		id := gen.GenerateTraceID()
		assert.Equal(26, len(id), "应生成 ULID")
	})
}

// BenchmarkDefaultIDGenerator 基准测试 - 默认生成器
func BenchmarkDefaultIDGenerator(b *testing.B) {
	gen := NewDefaultIDGenerator()
	b.Run("GenerateTraceID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateTraceID()
		}
	})
}

// BenchmarkUUIDGenerator 基准测试 - UUID 生成器
func BenchmarkUUIDGenerator(b *testing.B) {
	gen := NewUUIDGenerator()
	b.Run("GenerateTraceID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateTraceID()
		}
	})
}

// BenchmarkNanoIDGenerator 基准测试 - NanoID 生成器
func BenchmarkNanoIDGenerator(b *testing.B) {
	gen := NewNanoIDGenerator()
	b.Run("GenerateTraceID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateTraceID()
		}
	})
}

// BenchmarkSnowflakeGenerator 基准测试 - Snowflake 生成器
func BenchmarkSnowflakeGenerator(b *testing.B) {
	gen := NewSnowflakeGenerator(1, 1)
	b.Run("GenerateTraceID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateTraceID()
		}
	})
}

// BenchmarkULIDGenerator 基准测试 - ULID 生成器
func BenchmarkULIDGenerator(b *testing.B) {
	gen := NewULIDGenerator()
	b.Run("GenerateTraceID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateTraceID()
		}
	})
}

// TestConcurrentGeneration 测试并发生成
func TestConcurrentGeneration(t *testing.T) {
	generators := []struct {
		name string
		gen  IDGenerator
	}{
		{"Default", NewDefaultIDGenerator()},
		{"UUID", NewUUIDGenerator()},
		{"NanoID", NewNanoIDGenerator()},
		{"Snowflake", NewSnowflakeGenerator(1, 1)},
		{"ULID", NewULIDGenerator()},
	}

	for _, tt := range generators {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			ids := make(map[string]bool)
			mu := sync.Mutex{}

			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					id := tt.gen.GenerateTraceID()
					mu.Lock()
					assert.False(ids[id], "并发生成的 ID 应唯一")
					ids[id] = true
					mu.Unlock()
				}()
			}

			wg.Wait()
			assert.Equal(100, len(ids), "应生成 100 个唯一 ID")
		})
	}
}
