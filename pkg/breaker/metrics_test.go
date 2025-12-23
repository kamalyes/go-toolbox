/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-23 23:40:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-23 23:55:00
 * @FilePath: \go-toolbox\pkg\breaker\metrics_test.go
 * @Description: 指标收集器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package breaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewMetricsCollector 测试创建指标收集器
func TestNewMetricsCollector(t *testing.T) {
	mc := NewMetricsCollector()

	assert.NotNil(t, mc)
	assert.NotNil(t, mc.GetExecutionCount)
	assert.NotNil(t, mc.GetSuccessCount)
	assert.NotNil(t, mc.GetFailureCount)
}

// TestMetricsCollectorRecordStart 测试记录任务开始
func TestMetricsCollectorRecordStart(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")

	assert.Equal(t, int64(1), mc.GetTotalExecutions())
	assert.Equal(t, int64(1), mc.GetActiveCount())
}

// TestMetricsCollectorRecordSuccess 测试记录任务成功
func TestMetricsCollectorRecordSuccess(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 100*time.Millisecond)

	assert.Equal(t, int64(1), mc.GetTotalSuccess())
	assert.Equal(t, int64(0), mc.GetActiveCount())
}

// TestMetricsCollectorRecordFailure 测试记录任务失败
func TestMetricsCollectorRecordFailure(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")
	mc.RecordFailure("test-", 50*time.Millisecond)

	assert.Equal(t, int64(1), mc.GetTotalFailure())
	assert.Equal(t, int64(0), mc.GetActiveCount())
}

// TestMetricsCollectorMultiples 测试多个任务
func TestMetricsCollectorMultiples(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("1")
	mc.RecordSuccess("1", 100*time.Millisecond)

	mc.RecordStart("2")
	mc.RecordFailure("2", 50*time.Millisecond)

	assert.Equal(t, int64(2), mc.GetTotalExecutions())
	assert.Equal(t, int64(1), mc.GetTotalSuccess())
	assert.Equal(t, int64(1), mc.GetTotalFailure())
}

// TestMetricsCollectorGetMetrics 测试获取任务指标
func TestMetricsCollectorGetMetrics(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 100*time.Millisecond)

	metrics := mc.GetMetrics("test-")

	assert.NotNil(t, metrics)
	assert.Equal(t, "test-", metrics.Name)
	assert.Equal(t, int64(1), metrics.ExecutionCount)
	assert.Equal(t, int64(1), metrics.SuccessCount)
	assert.Equal(t, int64(0), metrics.FailureCount)
}

// TestMetricsCollectorGetAllMetrics 测试获取所有指标
func TestMetricsCollectorGetAllMetrics(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("1")
	mc.RecordSuccess("1", 100*time.Millisecond)

	mc.RecordStart("2")
	mc.RecordSuccess("2", 200*time.Millisecond)

	allMetrics := mc.GetAllMetrics()

	assert.NotNil(t, allMetrics)
	assert.Len(t, allMetrics, 2)
}

// TestMetricsCollectorCalculateAvgDuration 测试计算平均执行时间
func TestMetricsCollectorCalculateAvgDuration(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 100*time.Millisecond)

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 200*time.Millisecond)

	metrics := mc.GetMetrics("test-")

	assert.Equal(t, 150.0, metrics.AvgExecutionTime)
}

// TestMetricsCollectorTrackMaxDuration 测试跟踪最大执行时间
func TestMetricsCollectorTrackMaxDuration(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 100*time.Millisecond)

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 200*time.Millisecond)

	metrics := mc.GetMetrics("test-")

	assert.Equal(t, int64(200), metrics.MaxExecutionTime)
}

// TestMetricsCollectorTrackMinDuration 测试跟踪最小执行时间
func TestMetricsCollectorTrackMinDuration(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 200*time.Millisecond)

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 100*time.Millisecond)

	metrics := mc.GetMetrics("test-")

	assert.Equal(t, int64(100), metrics.MinExecutionTime)
}

// TestMetricsCollectorReset 测试重置指标
func TestMetricsCollectorReset(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 100*time.Millisecond)

	mc.Reset()

	assert.Equal(t, int64(0), mc.GetTotalExecutions())
	assert.Equal(t, int64(0), mc.GetTotalSuccess())
	assert.Equal(t, int64(0), mc.GetTotalFailure())
}

// TestMetricsCollectorConcurrent 测试并发记录
func TestMetricsCollectorConcurrent(t *testing.T) {
	mc := NewMetricsCollector()
	done := make(chan bool, 100)

	for i := 0; i < 100; i++ {
		go func(idx int) {
			Name := "test-"
			mc.RecordStart(Name)
			time.Sleep(time.Millisecond)
			if idx%2 == 0 {
				mc.RecordSuccess(Name, 10*time.Millisecond)
			} else {
				mc.RecordFailure(Name, 10*time.Millisecond)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}

	assert.Equal(t, int64(100), mc.totalExecutions)
	assert.Equal(t, int64(50), mc.totalSuccess)
	assert.Equal(t, int64(50), mc.totalFailure)
}

// TestMetricsCollectorSuccessRate 测试成功率计算
func TestMetricsCollectorSuccessRate(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 100*time.Millisecond)

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 100*time.Millisecond)

	mc.RecordStart("test-")
	mc.RecordFailure("test-", 100*time.Millisecond)

	metrics := mc.GetMetrics("test-")

	// 3次执行，2次成功，成功率应该是 66.67%
	successRate := float64(metrics.SuccessCount) / float64(metrics.ExecutionCount) * 100
	assert.InDelta(t, 66.67, successRate, 0.1)
}

// TestMetricsCollectorSnapshot 测试指标快照
func TestMetricsCollectorSnapshot(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("test-")
	mc.RecordSuccess("test-", 100*time.Millisecond)

	snapshot := mc.GetSnapshot()

	assert.NotNil(t, snapshot)
	assert.Equal(t, int64(1), snapshot.GlobalMetrics.TotalExecutions)
	assert.Equal(t, int64(1), snapshot.GlobalMetrics.TotalSuccess)
}
