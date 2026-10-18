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
	"fmt"
	"math"
	"regexp"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockCounterStore struct {
	data map[string]uint64
	mu   sync.Mutex
}

func (m *mockCounterStore) Increment(key string, delta uint64, initValue uint64) (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	if !ok {
		v = initValue
	}
	v += delta
	m.data[key] = v
	return v, nil
}

// TestDefaultIDGenerator 测试默认 Hex 生成器
func TestDefaultIDGenerator(t *testing.T) {
	gen := NewDefaultIDGenerator()

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "TraceID 不应为空")
		assert.Equal(t, 32, len(traceID), "TraceID 应为 32 字符")
		assert.True(t, regexp.MustCompile("^[0-9a-f]{32}$").MatchString(traceID), "TraceID 应为 hex 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(t, spanID, "SpanID 不应为空")
		assert.Equal(t, 16, len(spanID), "SpanID 应为 16 字符")
		assert.True(t, regexp.MustCompile("^[0-9a-f]{16}$").MatchString(spanID), "SpanID 应为 hex 格式")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(t, requestID, "RequestID 不应为空")
		assert.True(t, regexp.MustCompile(`^\d+-\d+$`).MatchString(requestID), "RequestID 应为 timestamp-counter 格式")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(t, correlationID, "CorrelationID 不应为空")
		assert.Equal(t, 36, len(correlationID), "CorrelationID 应为 36 字符 (UUID 格式)")
		assert.True(t, regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`).MatchString(correlationID), "CorrelationID 应为 UUID 格式")
	})

	t.Run("Uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 1000; i++ {
			id := gen.GenerateTraceID()
			assert.False(t, ids[id], "生成的 ID 应唯一")
			ids[id] = true
		}
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
}

// TestUUIDGenerator 测试 UUID 生成器
func TestUUIDGenerator(t *testing.T) {
	gen := NewUUIDGenerator()

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "TraceID 不应为空")
		assert.Equal(t, 36, len(traceID), "TraceID 应为 36 字符")
		assert.True(t, regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(traceID), "TraceID 应为 UUID v4 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(t, spanID, "SpanID 不应为空")
		assert.Equal(t, 16, len(spanID), "SpanID 应为 16 字符")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(t, requestID, "RequestID 不应为空")
		assert.True(t, regexp.MustCompile(`^[0-9a-f]{8}-\d+$`).MatchString(requestID), "RequestID 应包含 UUID 前缀和计数器")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(t, correlationID, "CorrelationID 不应为空")
		assert.Equal(t, 36, len(correlationID), "CorrelationID 应为 36 字符")
	})

	t.Run("DifferentFormats", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		spanID := gen.GenerateSpanID()
		requestID := gen.GenerateRequestID()
		correlationID := gen.GenerateCorrelationID()

		assert.NotEqual(t, traceID, spanID, "TraceID 和 SpanID 格式应不同")
		assert.NotEqual(t, traceID, requestID, "TraceID 和 RequestID 格式应不同")
		assert.NotEqual(t, traceID, correlationID, "TraceID 和 CorrelationID 格式应不同")
	})
}

// TestNanoIDGenerator 测试 NanoID 生成器
func TestNanoIDGenerator(t *testing.T) {
	gen := NewNanoIDGenerator()

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "TraceID 不应为空")
		assert.Equal(t, 21, len(traceID), "TraceID 应为 21 字符")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z_-]{21}$`).MatchString(traceID), "TraceID 应为 NanoID 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(t, spanID, "SpanID 不应为空")
		assert.Equal(t, 16, len(spanID), "SpanID 应为 16 字符")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(t, requestID, "RequestID 不应为空")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z_-]{10}-\d+$`).MatchString(requestID), "RequestID 应包含 NanoID 前缀和计数器")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(t, correlationID, "CorrelationID 不应为空")
		assert.Equal(t, 21, len(correlationID), "CorrelationID 应为 21 字符")
	})

	t.Run("DifferentFormats", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		spanID := gen.GenerateSpanID()
		requestID := gen.GenerateRequestID()
		correlationID := gen.GenerateCorrelationID()

		assert.NotEqual(t, traceID, spanID, "TraceID 和 SpanID 格式应不同")
		assert.NotEqual(t, traceID, requestID, "TraceID 和 RequestID 格式应不同")
		assert.NotEqual(t, traceID, correlationID, "TraceID 和 CorrelationID 格式应不同")
	})
}

// TestSnowflakeGenerator 测试 Snowflake 生成器
func TestSnowflakeGenerator(t *testing.T) {
	gen := NewSnowflakeGenerator(1, 1)

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "TraceID 不应为空")
		assert.Equal(t, 16, len(traceID), "TraceID 应为 16 字符 hex")
		assert.True(t, regexp.MustCompile(`^[0-9a-f]{16}$`).MatchString(traceID), "TraceID 应为 hex 格式")
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

		assert.NotEqual(t, traceID, spanID, "TraceID(hex 16字符) 和 SpanID(hex 8字符) 格式应不同")
		assert.NotEqual(t, traceID, requestID, "TraceID 和 RequestID 格式应不同")
		assert.NotEqual(t, traceID, correlationID, "TraceID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, spanID, requestID, "SpanID 和 RequestID 格式应不同")
		assert.NotEqual(t, spanID, correlationID, "SpanID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, requestID, correlationID, "RequestID 和 CorrelationID 格式应不同")
	})

	t.Run("Monotonic", func(t *testing.T) {
		var lastID int64
		for i := 0; i < 100; i++ {
			id := gen.Generate()
			assert.True(t, id > lastID, "Snowflake ID 应单调递增")
			lastID = id
		}
	})
}

// TestShortIDGenerator 测试短ID生成器
func TestShortIDGenerator(t *testing.T) {
	gen := NewShortIDGenerator()

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "TraceID 不应为空")
		assert.Equal(t, 10, len(traceID), "TraceID 应为 10 字符")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z]{10}$`).MatchString(traceID), "TraceID 应为 Base62 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(t, spanID, "SpanID 不应为空")
		assert.Equal(t, 8, len(spanID), "SpanID 应为 8 字符")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z]{8}$`).MatchString(spanID), "SpanID 应为 Base62 格式")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(t, requestID, "RequestID 不应为空")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z]{5}-\d+$`).MatchString(requestID), "RequestID 应为 Base62前缀-计数器 格式")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(t, correlationID, "CorrelationID 不应为空")
		assert.Equal(t, 10, len(correlationID), "CorrelationID 应为 10 字符")
		assert.True(t, regexp.MustCompile(`^[0-9A-Za-z]{10}$`).MatchString(correlationID), "CorrelationID 应为 Base62 格式")
	})

	t.Run("DifferentFormats", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		spanID := gen.GenerateSpanID()
		requestID := gen.GenerateRequestID()
		correlationID := gen.GenerateCorrelationID()

		assert.NotEqual(t, traceID, spanID, "TraceID(10字符) 和 SpanID(8字符) 格式应不同")
		assert.NotEqual(t, traceID, requestID, "TraceID 和 RequestID 格式应不同")
		assert.NotEqual(t, traceID, correlationID, "TraceID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, spanID, requestID, "SpanID 和 RequestID 格式应不同")
		assert.NotEqual(t, spanID, correlationID, "SpanID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, requestID, correlationID, "RequestID 和 CorrelationID 格式应不同")
	})

	t.Run("TraceIDTimeSortable", func(t *testing.T) {
		id1 := gen.GenerateTraceID()
		id2 := gen.GenerateTraceID()
		assert.True(t, id2 >= id1, "TraceID 应时间可排序（字典序=时间序）")
	})

	t.Run("Uniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 1000; i++ {
			id := gen.GenerateTraceID()
			assert.False(t, ids[id], "生成的 ID 应唯一")
			ids[id] = true
		}
	})
}

// TestNumericIDGenerator 测试8位纯数字ID生成器
func TestNumericIDGenerator(t *testing.T) {
	gen := NewNumericIDGenerator()

	t.Run("GenerateUserID", func(t *testing.T) {
		userID := gen.GenerateUserID()
		assert.NotEmpty(t, userID, "UserID 不应为空")
		assert.Equal(t, 8, len(userID), "UserID 应为 8 位数字")
		assert.True(t, regexp.MustCompile(`^[1-9]\d{7}$`).MatchString(userID), "UserID 应为 8 位数字（首位非0）")
	})

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "TraceID 不应为空")
		assert.Equal(t, 8, len(traceID), "TraceID 应为 8 位数字")
		assert.True(t, regexp.MustCompile(`^[1-9]\d{7}$`).MatchString(traceID), "TraceID 应为 8 位数字（首位非0）")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(t, spanID, "SpanID 不应为空")
		assert.Equal(t, 8, len(spanID), "SpanID 应为 8 位数字")
		assert.True(t, regexp.MustCompile(`^[1-9]\d{7}$`).MatchString(spanID), "SpanID 应为 8 位数字（首位非0）")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(t, requestID, "RequestID 不应为空")
		assert.Equal(t, 8, len(requestID), "RequestID 应为 8 位数字")
		assert.True(t, regexp.MustCompile(`^[1-9]\d{7}$`).MatchString(requestID), "RequestID 应为 8 位数字（首位非0）")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(t, correlationID, "CorrelationID 不应为空")
		assert.Equal(t, 8, len(correlationID), "CorrelationID 应为 8 位数字")
		assert.True(t, regexp.MustCompile(`^[1-9]\d{7}$`).MatchString(correlationID), "CorrelationID 应为 8 位数字（首位非0）")
	})

	t.Run("DifferentFormats", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		spanID := gen.GenerateSpanID()
		requestID := gen.GenerateRequestID()
		correlationID := gen.GenerateCorrelationID()
		userID := gen.GenerateUserID()

		assert.NotEqual(t, traceID, spanID, "TraceID 和 SpanID 应不同")
		assert.NotEqual(t, traceID, requestID, "TraceID 和 RequestID 应不同")
		assert.NotEqual(t, traceID, correlationID, "TraceID 和 CorrelationID 应不同")
		assert.NotEqual(t, spanID, requestID, "SpanID 和 RequestID 应不同")
		assert.NotEqual(t, spanID, correlationID, "SpanID 和 CorrelationID 应不同")
		assert.NotEqual(t, requestID, correlationID, "RequestID 和 CorrelationID 应不同")
		assert.NotEqual(t, userID, traceID, "UserID 和 TraceID 应不同")
	})

	t.Run("UserIDSequential", func(t *testing.T) {
		freshGen := NewNumericIDGenerator()
		id1 := freshGen.GenerateUserID()
		id2 := freshGen.GenerateUserID()
		id3 := freshGen.GenerateUserID()
		assert.True(t, id2 > id1, "UserID 应严格递增")
		assert.True(t, id3 > id2, "UserID 应严格递增")
	})

	t.Run("UserIDUniqueness", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 10000; i++ {
			id := gen.GenerateUserID()
			assert.False(t, ids[id], "生成的 UserID 应唯一")
			ids[id] = true
		}
	})

	t.Run("AllDigitsInRange", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			id := gen.GenerateUserID()
			num, err := strconv.Atoi(id)
			assert.Nil(t, err, "UserID 应为纯数字")
			assert.True(t, num >= 10000000 && num <= 99999999, "UserID 应在 10000000-99999999 范围内")
		}
	})

	t.Run("NoMutexNoTimeWheel", func(t *testing.T) {
		var wg sync.WaitGroup
		ids := make(chan string, 10000)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					ids <- gen.GenerateUserID()
				}
			}()
		}
		wg.Wait()
		close(ids)

		unique := make(map[string]bool)
		for id := range ids {
			assert.False(t, unique[id], "并发生成的 UserID 应唯一")
			unique[id] = true
		}
		assert.Equal(t, 10000, len(unique), "应生成 10000 个唯一 UserID")
	})

	t.Run("DistributedWorkerID", func(t *testing.T) {
		gen0 := NewNumericIDGeneratorWithWorker(0)
		gen1 := NewNumericIDGeneratorWithWorker(1)
		gen9 := NewNumericIDGeneratorWithWorker(9)

		id0 := gen0.GenerateUserID()
		id1 := gen1.GenerateUserID()
		id9 := gen9.GenerateUserID()

		assert.NotEqual(t, id0, id1, "不同 Worker 的 UserID 应不同")
		assert.NotEqual(t, id1, id9, "不同 Worker 的 UserID 应不同")

		num0, _ := strconv.Atoi(id0)
		num1, _ := strconv.Atoi(id1)
		num9, _ := strconv.Atoi(id9)

		assert.Equal(t, num1-num0, 10000, "Worker 1 与 Worker 0 偏移应为 10000")
		assert.Equal(t, num9-num0, 90000, "Worker 9 与 Worker 0 偏移应为 90000")
	})

	t.Run("WorkerIDNoOverlap", func(t *testing.T) {
		gen0 := NewNumericIDGeneratorWithWorker(0)
		gen1 := NewNumericIDGeneratorWithWorker(1)

		ids0 := make(map[string]bool)
		ids1 := make(map[string]bool)

		for i := 0; i < 9999; i++ {
			ids0[gen0.GenerateUserID()] = true
			ids1[gen1.GenerateUserID()] = true
		}

		for id := range ids0 {
			assert.False(t, ids1[id], "Worker 0 和 Worker 1 不应生成相同 UserID")
		}
	})

	t.Run("PerMachineDailyCapacity", func(t *testing.T) {
		freshGen := NewNumericIDGeneratorWithWorker(0)
		first := freshGen.GenerateUserID()
		numFirst, _ := strconv.Atoi(first)

		for i := 0; i < 9998; i++ {
			freshGen.GenerateUserID()
		}
		last := freshGen.GenerateUserID()
		numLast, _ := strconv.Atoi(last)

		assert.Equal(t, numLast-numFirst, 9999, "每机每天应支持 10000 个 UserID（0-9999）")
	})

	t.Run("WorkerIDModulo", func(t *testing.T) {
		gen := NewNumericIDGeneratorWithWorker(15)
		assert.Equal(t, gen.WorkerID(), uint64(5), "WorkerID 15 %% 10 = 5")
	})

	t.Run("CustomConfig", func(t *testing.T) {
		cfg := NumericIDConfig{
			Epoch:        1704067200,
			Base:         100000000,
			WorkerSpace:  100000,
			MaxWorkers:   5,
			DaySpace:     500000,
			RandomDigits: 9,
			BatchSize:    1000,
		}
		gen := NewNumericIDGeneratorWithConfigAndWorker(cfg, 2)

		userID := gen.GenerateUserID()
		assert.Equal(t, 9, len(userID), "9位配置应生成9位数字")

		spanID := gen.GenerateSpanID()
		assert.Equal(t, 9, len(spanID), "9位配置的 SpanID 应为9位")

		assert.Equal(t, uint64(2), gen.WorkerID(), "WorkerID 应为 2")
		assert.Equal(t, uint64(5), gen.Config().MaxWorkers, "MaxWorkers 应为 5")
	})

	t.Run("CustomConfigWorkerOffset", func(t *testing.T) {
		cfg := NumericIDConfig{
			Epoch:        1704067200,
			Base:         100000000,
			WorkerSpace:  100000,
			MaxWorkers:   5,
			DaySpace:     500000,
			RandomDigits: 9,
			BatchSize:    1000,
		}
		gen0 := NewNumericIDGeneratorWithConfigAndWorker(cfg, 0)
		gen1 := NewNumericIDGeneratorWithConfigAndWorker(cfg, 1)

		id0 := gen0.GenerateUserID()
		id1 := gen1.GenerateUserID()
		num0, _ := strconv.Atoi(id0)
		num1, _ := strconv.Atoi(id1)

		assert.Equal(t, num1-num0, 100000, "自定义配置中 Worker 偏移应为 WorkerSpace=100000")
	})

	t.Run("ConfigValidateFail", func(t *testing.T) {
		badCfg := NumericIDConfig{
			Epoch:       0,
			Base:        10000000,
			WorkerSpace: 10000,
			MaxWorkers:  10,
			DaySpace:    100000,
			BatchSize:   100,
		}
		assert.NotNil(t, badCfg.Validate(), "Epoch=0 应校验失败")

		mismatchCfg := NumericIDConfig{
			Epoch:       1704067200,
			Base:        10000000,
			WorkerSpace: 10000,
			MaxWorkers:  10,
			DaySpace:    99999,
			BatchSize:   100,
		}
		assert.NotNil(t, mismatchCfg.Validate(), "DaySpace != WorkerSpace*MaxWorkers 应校验失败")

		badBatch := NumericIDConfig{
			Epoch:       1704067200,
			Base:        10000000,
			WorkerSpace: 10000,
			MaxWorkers:  10,
			DaySpace:    100000,
			BatchSize:   0,
		}
		assert.NotNil(t, badBatch.Validate(), "BatchSize=0 应校验失败")

		badBatch2 := NumericIDConfig{
			Epoch:       1704067200,
			Base:        10000000,
			WorkerSpace: 10000,
			MaxWorkers:  10,
			DaySpace:    100000,
			BatchSize:   20000,
		}
		assert.NotNil(t, badBatch2.Validate(), "BatchSize > WorkerSpace 应校验失败")
	})

	t.Run("CounterStoreBatchPrefetch", func(t *testing.T) {
		store := &mockCounterStore{data: make(map[string]uint64)}

		cfg := DefaultNumericIDConfig()
		cfg.Store = store
		cfg.BatchSize = 10

		gen := NewNumericIDGeneratorWithConfigAndWorker(cfg, 0)

		ids := make([]uint64, 0, 25)
		for i := 0; i < 25; i++ {
			idStr := gen.GenerateUserID()
			id, _ := strconv.ParseUint(idStr, 10, 64)
			ids = append(ids, id)
		}

		for i := 1; i < len(ids); i++ {
			assert.Equal(t, ids[i]-ids[i-1], uint64(1), fmt.Sprintf("ID应连续递增，ids[%d]=%d ids[%d]=%d", i, ids[i], i-1, ids[i-1]))
		}
	})

	t.Run("CounterStoreRecycle", func(t *testing.T) {
		store := &mockCounterStore{data: make(map[string]uint64)}

		cfg := DefaultNumericIDConfig()
		cfg.Store = store
		cfg.BatchSize = 10

		gen1 := NewNumericIDGeneratorWithConfigAndWorker(cfg, 0)
		id1 := gen1.GenerateUserID()
		_ = gen1.GenerateUserID()
		id3 := gen1.GenerateUserID()
		num3, _ := strconv.Atoi(id3)
		num1, _ := strconv.Atoi(id1)
		assert.Equal(t, num3-num1, 2, "应连续递增")

		gen2 := NewNumericIDGeneratorWithConfigAndWorker(cfg, 0)
		id4 := gen2.GenerateUserID()
		num4, _ := strconv.Atoi(id4)
		assert.True(t, num4 > num3, fmt.Sprintf("重启后应从上次位置继续，num4=%d > num3=%d", num4, num3))
	})

	t.Run("CounterStoreDayKeyIsolation", func(t *testing.T) {
		store := &mockCounterStore{data: make(map[string]uint64)}

		cfg := DefaultNumericIDConfig()
		cfg.Store = store
		cfg.BatchSize = 10

		gen := NewNumericIDGeneratorWithConfigAndWorker(cfg, 0)
		_ = gen.GenerateUserID()

		days := uint64((time.Now().Unix()-cfg.Epoch)/86400) + 1
		key := fmt.Sprintf("numeric:0:%d", days-1)
		_, exists := store.data[key]
		assert.True(t, exists, fmt.Sprintf("应使用按天隔离的 key: %s", key))
	})

	t.Run("CounterStoreAtomicIncrement", func(t *testing.T) {
		store := &mockCounterStore{data: make(map[string]uint64)}

		cfg := DefaultNumericIDConfig()
		cfg.Store = store
		cfg.BatchSize = 5

		gen0 := NewNumericIDGeneratorWithConfigAndWorker(cfg, 0)
		gen1 := NewNumericIDGeneratorWithConfigAndWorker(cfg, 1)

		id0 := gen0.GenerateUserID()
		id1 := gen1.GenerateUserID()
		num0, _ := strconv.Atoi(id0)
		num1, _ := strconv.Atoi(id1)

		assert.True(t, num1-num0 >= 10000, fmt.Sprintf("不同 Worker 的 ID 应有 WorkerSpace 偏移，num0=%d num1=%d", num0, num1))
	})
}

// TestULIDGenerator 测试 ULID 生成器
func TestULIDGenerator(t *testing.T) {
	gen := NewULIDGenerator()

	t.Run("GenerateTraceID", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "TraceID 不应为空")
		assert.Equal(t, 26, len(traceID), "TraceID 应为 26 字符")
		assert.True(t, regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`).MatchString(traceID), "TraceID 应为 ULID 格式")
	})

	t.Run("GenerateSpanID", func(t *testing.T) {
		spanID := gen.GenerateSpanID()
		assert.NotEmpty(t, spanID, "SpanID 不应为空")
		assert.Equal(t, 16, len(spanID), "SpanID 应为 16 字符（ULID 随机部分）")
		assert.True(t, regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{16}$`).MatchString(spanID), "SpanID 应为 ULID 随机部分格式")
	})

	t.Run("GenerateRequestID", func(t *testing.T) {
		requestID := gen.GenerateRequestID()
		assert.NotEmpty(t, requestID, "RequestID 不应为空")
		assert.True(t, regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{10}-\d+$`).MatchString(requestID), "RequestID 应包含 ULID 时间戳前缀和计数器")
	})

	t.Run("GenerateCorrelationID", func(t *testing.T) {
		correlationID := gen.GenerateCorrelationID()
		assert.NotEmpty(t, correlationID, "CorrelationID 不应为空")
		assert.True(t, regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}-[0-9A-HJKMNP-TV-Z]{26}$`).MatchString(correlationID), "CorrelationID 应为 ULID-ULID 格式")
	})

	t.Run("DifferentFormats", func(t *testing.T) {
		traceID := gen.GenerateTraceID()
		spanID := gen.GenerateSpanID()
		requestID := gen.GenerateRequestID()
		correlationID := gen.GenerateCorrelationID()

		assert.NotEqual(t, traceID, spanID, "TraceID(26字符) 和 SpanID(16字符) 格式应不同")
		assert.NotEqual(t, traceID, requestID, "TraceID 和 RequestID 格式应不同")
		assert.NotEqual(t, traceID, correlationID, "TraceID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, spanID, requestID, "SpanID 和 RequestID 格式应不同")
		assert.NotEqual(t, spanID, correlationID, "SpanID 和 CorrelationID 格式应不同")
		assert.NotEqual(t, requestID, correlationID, "RequestID 和 CorrelationID 格式应不同")
	})
}

// TestFactory 测试工厂函数
func TestFactory(t *testing.T) {
	t.Run("NewIDGenerator with GeneratorType", func(t *testing.T) {
		gen := NewIDGenerator(GeneratorTypeUUID)
		assert.NotNil(t, gen, "生成器不应为 nil")
		id := gen.GenerateTraceID()
		assert.True(t, regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(id), "应生成 UUID v4")
	})

	t.Run("NewIDGenerator with string", func(t *testing.T) {
		gen := NewIDGenerator("nanoid")
		assert.NotNil(t, gen, "生成器不应为 nil")
		id := gen.GenerateTraceID()
		assert.Equal(t, 21, len(id), "应生成 NanoID")
	})

	t.Run("NewIDGenerator with invalid type", func(t *testing.T) {
		gen := NewIDGenerator(12345)
		assert.NotNil(t, gen, "生成器不应为 nil")
		id := gen.GenerateTraceID()
		assert.Equal(t, 32, len(id), "应回退到默认生成器")
	})

	t.Run("NewIDGeneratorFromString (deprecated)", func(t *testing.T) {
		gen := NewIDGeneratorFromString("ulid")
		assert.NotNil(t, gen, "生成器不应为 nil")
		id := gen.GenerateTraceID()
		assert.Equal(t, 26, len(id), "应生成 ULID")
	})

	t.Run("SnowflakeUsesDistributedWorkerID", func(t *testing.T) {
		gen := NewIDGenerator("snowflake")
		assert.NotNil(t, gen, "生成器不应为 nil")
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "Snowflake TraceID 不应为空")
		assert.Equal(t, 16, len(traceID), "Snowflake TraceID 应为 16 字符 hex")
	})

	t.Run("ShortFlakeUsesDistributedWorkerID", func(t *testing.T) {
		gen := NewIDGenerator("shortflake")
		assert.NotNil(t, gen, "生成器不应为 nil")
		traceID := gen.GenerateTraceID()
		assert.NotEmpty(t, traceID, "ShortFlake TraceID 不应为空")
	})

	t.Run("NumericUsesDistributedWorkerID", func(t *testing.T) {
		gen := NewIDGenerator("numeric")
		assert.NotNil(t, gen, "生成器不应为 nil")
		userID := gen.(*NumericIDGenerator).GenerateUserID()
		assert.NotEmpty(t, userID, "Numeric UserID 不应为空")
		assert.Equal(t, 8, len(userID), "Numeric UserID 应为 8 位数字")
	})
}

// TestIDType 测试 IDType 枚举
func TestIDType(t *testing.T) {
	assert.Equal(t, string(IDTypeTraceID), "trace_id", "IDTypeTraceID 应为 trace_id")
	assert.Equal(t, string(IDTypeSpanID), "span_id", "IDTypeSpanID 应为 span_id")
	assert.Equal(t, string(IDTypeRequestID), "request_id", "IDTypeRequestID 应为 request_id")
	assert.Equal(t, string(IDTypeCorrelationID), "correlation_id", "IDTypeCorrelationID 应为 correlation_id")
}

// TestIDSpec 测试 IDSpec 规格
func TestIDSpec(t *testing.T) {
	t.Run("GeneratorType Spec", func(t *testing.T) {
		spec := GeneratorTypeDefault.Spec()
		assert.Equal(t, spec.TraceLen, 32, "Default TraceLen 应为 32")
		assert.Equal(t, spec.SpanLen, 16, "Default SpanLen 应为 16")
		assert.True(t, spec.RequestCounter, "Default RequestCounter 应为 true")
		assert.True(t, spec.CorrelationFmt, "Default CorrelationFmt 应为 true")
	})

	t.Run("Unknown GeneratorType", func(t *testing.T) {
		spec := GeneratorType("unknown").Spec()
		assert.Equal(t, spec.TraceLen, 32, "Unknown 应回退到 DefaultSpec")
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

// BenchmarkShortIDGenerator 基准测试 - ShortID 生成器
func BenchmarkShortIDGenerator(b *testing.B) {
	gen := NewShortIDGenerator()
	b.Run("GenerateTraceID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateTraceID()
		}
	})
}

// BenchmarkNumericIDGenerator 基准测试 - NumericID 生成器
func BenchmarkNumericIDGenerator(b *testing.B) {
	gen := NewNumericIDGenerator()
	b.Run("GenerateUserID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateUserID()
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
		{"ShortID", NewShortIDGenerator()},
		{"Numeric", NewNumericIDGenerator()},
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
					assert.False(t, ids[id], "并发生成的 ID 应唯一")
					ids[id] = true
					mu.Unlock()
				}()
			}

			wg.Wait()
			assert.Equal(t, 100, len(ids), "应生成 100 个唯一 ID")
		})
	}
}

// ============================================================================
// 定宽 hex 格式化：与被替换掉的 fmt.Sprintf 实现逐字符对照
// ============================================================================

// 下面四个函数是替换前的 Sprintf 实现原样复刻，仅用于对照
func oldHex16(id int64) string { return fmt.Sprintf("%016x", id) }
func oldHex8(id int64) string  { return fmt.Sprintf("%08x", id&0xFFFFFFFF) }
func oldHex13(id int64) string { return fmt.Sprintf("%013x", id) }
func oldSnowflakeCorr(id int64, counter uint64) string {
	h := fmt.Sprintf("%016x", id)
	var buf [36]byte
	copy(buf[0:8], h[0:8])
	buf[8] = '-'
	copy(buf[9:13], h[8:12])
	buf[13] = '-'
	copy(buf[14:18], "4b10")
	buf[18] = '-'
	copy(buf[19:23], h[12:16])
	buf[23] = '-'
	copy(buf[24:36], fmt.Sprintf("%012x", counter))
	return string(buf[:])
}
func oldShortFlakeCorr(id1, id2 int64) string {
	h1 := fmt.Sprintf("%016x", id1)
	h2 := fmt.Sprintf("%016x", id2)
	h1 = h1[:12] + "4" + h1[13:]
	h2 = "8" + h2[1:]
	var buf [36]byte
	copy(buf[0:8], h1[0:8])
	buf[8] = '-'
	copy(buf[9:13], h1[8:12])
	buf[13] = '-'
	copy(buf[14:18], h1[12:16])
	buf[18] = '-'
	copy(buf[19:23], h2[0:4])
	buf[23] = '-'
	copy(buf[24:36], h2[4:16])
	return string(buf[:])
}

// hexTestIDs 覆盖真实 ID 范围与边界：0、小值、53 位上限、63 位上限
func hexTestIDs() []int64 {
	return []int64{
		0, 1, 0xF, 0x10, 0xFFFFFFFF, 0x100000000,
		1 << 20, 1 << 40, 1 << 52, (1 << 53) - 1,
		1 << 53, 1 << 62, (1 << 62) - 1, 0x7FFFFFFFFFFFFFFF,
		1732184000123456789, 6203154432000000,
	}
}

// TestHexHelpersMatchSprintf 手写定宽 hex 必须与 Sprintf 逐字符一致
//
// hex16 与 hex8 是精确等价而非近似：uint64 最多 16 位 hex、uint32 最多 8 位，
// 不存在值超出宽度的情况而 13 位低于 53 位 ID 的上限，所以那里用的是补零版
//
// [EN] The hand-written fixed-width hex must match Sprintf character for
// character. hex16 and hex8 are exact, not approximate: a uint64 never needs
// more than 16 digits nor a uint32 more than 8. Width 13 is below a 53-bit ID's
// maximum, which is why that path uses the zero-padding variant instead
func TestHexHelpersMatchSprintf(t *testing.T) {
	for _, id := range hexTestIDs() {
		assert.Equal(t, oldHex16(id), hex16(uint64(id)), "hex16(%#x)", id)
		assert.Equal(t, oldHex8(id), hex8(uint32(id)), "hex8(%#x)", id)
		var buf [16]byte
		assert.Equal(t, oldHex13(id), string(appendHexPadded(buf[:0], uint64(id), 13)),
			"appendHexPadded(%#x, 13)", id)
	}
	// 补零版还要在宽度小于值所需位数时照实输出，不截断
	for _, w := range []int{8, 12, 13, 16} {
		for _, v := range []uint64{0, 1, 0xFFFFFFFF, 1 << 52, (1 << 53) - 1, ^uint64(0)} {
			var buf [32]byte
			assert.Equal(t, fmt.Sprintf("%0*x", w, v),
				string(appendHexPadded(buf[:0], v, w)), "width=%d v=%#x", w, v)
		}
	}
}

// TestCorrelationIDLayoutUnchanged 关联 ID 的 8-4-4-4-12 版式不得因改写而变化
func TestCorrelationIDLayoutUnchanged(t *testing.T) {
	ids := hexTestIDs()
	for _, id := range ids {
		for _, counter := range []uint64{1, 0xFF, 1 << 20, (1 << 48) - 1} {
			var h, c [16]byte
			putHex16(&h, uint64(id))
			putHex16(&c, counter)
			var buf [36]byte
			copy(buf[0:8], h[0:8])
			buf[8] = '-'
			copy(buf[9:13], h[8:12])
			buf[13] = '-'
			copy(buf[14:18], "4b10")
			buf[18] = '-'
			copy(buf[19:23], h[12:16])
			buf[23] = '-'
			copy(buf[24:36], c[4:16])
			assert.Equal(t, oldSnowflakeCorr(id, counter), string(buf[:]),
				"snowflake corr(%#x, %d)", id, counter)
		}
	}
	for _, id1 := range ids {
		for _, id2 := range ids {
			var h1, h2 [16]byte
			putHex16(&h1, uint64(id1))
			putHex16(&h2, uint64(id2))
			h1[12] = '4'
			h2[0] = '8'
			var buf [36]byte
			copy(buf[0:8], h1[0:8])
			buf[8] = '-'
			copy(buf[9:13], h1[8:12])
			buf[13] = '-'
			copy(buf[14:18], h1[12:16])
			buf[18] = '-'
			copy(buf[19:23], h2[0:4])
			buf[23] = '-'
			copy(buf[24:36], h2[4:16])
			assert.Equal(t, oldShortFlakeCorr(id1, id2), string(buf[:]),
				"shortflake corr(%#x, %#x)", id1, id2)
		}
	}
}

// TestSnowflakeClockRegression Snowflake 对时钟回拨做钳位，ID 不回头
func TestSnowflakeClockRegression(t *testing.T) {
	g := NewSnowflakeGenerator(1, 1)
	g.Generate()
	g.mu.Lock()
	g.lastTime += 3
	g.mu.Unlock()

	seen := make(map[int64]bool, 256)
	prev := int64(-1)
	for i := 0; i < 200; i++ {
		id := g.Generate()
		assert.False(t, seen[id], "第 %d 个 ID 重复: %d", i, id)
		assert.Greater(t, id, prev, "第 %d 个 ID 不递增", i)
		seen[id] = true
		prev = id
	}
}

// ============================================================================
// 序列号耗尽后的等待：只量 CPU 占用，不断言
// ============================================================================

func processCPUTime() time.Duration {
	var ru syscall.Rusage
	_ = syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

// TestWaitStrategyCPUCost 打印等待期间的 CPU 占用，供改动等待策略时对照
//
// 不做断言：CPU 比例随机器与负载波动，断了只会变成假红
// 改造前两种等法（持锁轮询 time.Now、runtime.Gosched 自旋）都是 100% 占满一个
// 核空转，改成精确睡到毫秒边界后 Snowflake 降到约 34%、ShortFlake 约 2%，
// 而吞吐完全不变——吞吐由每毫秒配额决定，与等法无关
//
// [EN] Logs the CPU burned while waiting, for comparison when the wait strategy
// changes. It asserts nothing: the ratio moves with machine and load, so an
// assertion would only produce false failures
func TestWaitStrategyCPUCost(t *testing.T) {
	cases := []struct {
		name string
		n    int
		gen  func() int64
	}{
		{"Snowflake", 40000, NewSnowflakeGenerator(1, 1).Generate},
		{"ShortFlake", 600, NewShortFlakeGenerator(1).Generate},
	}
	for _, c := range cases {
		c.gen() // 预热，让首个毫秒窗口对齐
		cpu0, wall0 := processCPUTime(), time.Now()
		var sink int64
		for i := 0; i < c.n; i++ {
			sink = c.gen()
		}
		_ = sink
		wall := time.Since(wall0)
		cpu := processCPUTime() - cpu0
		t.Logf("%-11s %6d 个 ID  墙上 %7.2fms  CPU %7.2fms  CPU/墙上 %5.1f%%",
			c.name, c.n, float64(wall.Microseconds())/1000,
			float64(cpu.Microseconds())/1000, 100*float64(cpu)/float64(wall))
	}
}

// ============================================================================
// 基准：定宽 hex 手写与 fmt.Sprintf 的差距
// ============================================================================

var hexSink string

func BenchmarkHexFormatting(b *testing.B) {
	v := uint64(0x18f5a3c7d2e4b109)
	b.Run("16位/Sprintf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hexSink = fmt.Sprintf("%016x", v)
		}
	})
	b.Run("16位/手写", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hexSink = hex16(v)
		}
	})
	b.Run("8位/Sprintf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hexSink = fmt.Sprintf("%08x", uint32(v))
		}
	})
	b.Run("8位/手写", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hexSink = hex8(uint32(v))
		}
	})
	b.Run("13位/Sprintf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hexSink = fmt.Sprintf("%013x", v)
		}
	})
	b.Run("13位/补零", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var buf [16]byte
			hexSink = string(appendHexPadded(buf[:0], v, 13))
		}
	})
}

// ============================================================================
// 参数校验：越界的节点 ID 必须启动即报错，不能静默折回
// ============================================================================

// TestSnowflakeConfigRejectsOutOfRange 越界或负的 ID 必须 panic
//
// 原先是 workerID & 0x1F 静默折回，于是配成 32 的机器和配成 0 的机器变成同一个
// worker，两台机器从此发出完全相同的 ID这个测试守住「宁可启动失败」
func TestSnowflakeConfigRejectsOutOfRange(t *testing.T) {
	cases := []struct {
		name               string
		worker, datacenter int64
	}{
		{"worker 刚越界", SnowflakeMaxWorkerID + 1, 0},
		{"worker 折回会与 0 冲突", 32, 0},
		{"worker 负数", -1, 0},
		{"datacenter 刚越界", 0, SnowflakeMaxDatacenterID + 1},
		{"datacenter 负数", 0, -1},
		{"两者都越界", 999, 999},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Panics(t, func() { NewSnowflakeGenerator(c.worker, c.datacenter) })
			err := SnowflakeConfig{Epoch: 1640995200000, WorkerID: c.worker, DatacenterID: c.datacenter}.Validate()
			assert.Error(t, err, "Validate 应报错")
			assert.Contains(t, err.Error(), "must be in [0,", "错误信息应给出合法区间")
		})
	}
}

// TestSnowflakeConfigAcceptsBoundary 边界值必须可用，且节点位精确落位
func TestSnowflakeConfigAcceptsBoundary(t *testing.T) {
	for _, w := range []int64{0, 1, SnowflakeMaxWorkerID} {
		for _, d := range []int64{0, 1, SnowflakeMaxDatacenterID} {
			g := NewSnowflakeGenerator(w, d)
			id := g.Generate()
			assert.Equal(t, w, (id>>17)&0x1F, "worker 位应为 %d", w)
			assert.Equal(t, d, (id>>12)&0x1F, "datacenter 位应为 %d", d)
		}
	}
}

// TestSnowflakeConfigEpoch 纪元必须为正且不在未来
func TestSnowflakeConfigEpoch(t *testing.T) {
	assert.Error(t, SnowflakeConfig{Epoch: 0}.Validate())
	assert.Error(t, SnowflakeConfig{Epoch: -1}.Validate())
	assert.Error(t, SnowflakeConfig{Epoch: time.Now().Add(time.Hour).UnixMilli()}.Validate())
	assert.NoError(t, DefaultSnowflakeConfig().Validate())
	// Epoch 留空时按默认纪元填充，不应报错
	assert.NotPanics(t, func() { NewSnowflakeGeneratorWithConfig(SnowflakeConfig{WorkerID: 1}) })
}

// TestBoundNodeID 工厂把宿主推导的 ID 收进容量内，且对 MinInt64 不崩
func TestBoundNodeID(t *testing.T) {
	assert.Equal(t, int64(0), boundNodeID(0, 63))
	assert.Equal(t, int64(63), boundNodeID(63, 63))
	assert.Equal(t, int64(0), boundNodeID(64, 63), "64 折回 0，这是容量所限")
	assert.Equal(t, int64(1), boundNodeID(-1, 63), "负数取绝对值")
	assert.Equal(t, int64(0), boundNodeID(math.MinInt64, 63), "MinInt64 取反会溢出，须特判")
	for _, id := range []int64{0, 1, 7, 1023, -5, math.MinInt64, math.MaxInt64} {
		got := boundNodeID(id, SnowflakeMaxWorkerID)
		assert.GreaterOrEqual(t, got, int64(0))
		assert.LessOrEqual(t, got, int64(SnowflakeMaxWorkerID))
	}
}

// TestFactoryBuildsAllTypes 工厂对所有类型都能构造出来，不因参数越界 panic
func TestFactoryBuildsAllTypes(t *testing.T) {
	for _, name := range []string{"uuid", "nanoid", "snowflake", "shortflake", "short", "shortid", "numeric", "ulid", "default", "hex", "logger", "", "unknown"} {
		assert.NotPanics(t, func() {
			g := NewIDGenerator(name)
			assert.NotNil(t, g, "类型 %q", name)
			assert.NotEmpty(t, g.GenerateTraceID(), "类型 %q", name)
		}, "类型 %q", name)
	}
}
