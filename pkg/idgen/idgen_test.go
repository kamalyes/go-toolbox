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
	"regexp"
	"strconv"
	"sync"
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
