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

// TestMetricsCollectorGetExecutionCount 测试获取执行次数
func TestMetricsCollectorGetExecutionCount(t *testing.T) {
	mc := NewMetricsCollector()

	// 不存在的 name 返回 0
	assert.Equal(t, int64(0), mc.GetExecutionCount("not-exist"))

	mc.RecordStart("test-name")
	mc.RecordStart("test-name")
	assert.Equal(t, int64(2), mc.GetExecutionCount("test-name"))
}

// TestMetricsCollectorGetSuccessCount 测试获取成功次数
func TestMetricsCollectorGetSuccessCount(t *testing.T) {
	mc := NewMetricsCollector()

	// 不存在的 name 返回 0
	assert.Equal(t, int64(0), mc.GetSuccessCount("not-exist"))

	mc.RecordStart("test-name")
	mc.RecordSuccess("test-name", 50*time.Millisecond)
	assert.Equal(t, int64(1), mc.GetSuccessCount("test-name"))
}

// TestMetricsCollectorGetFailureCount 测试获取失败次数
func TestMetricsCollectorGetFailureCount(t *testing.T) {
	mc := NewMetricsCollector()

	// 不存在的 name 返回 0
	assert.Equal(t, int64(0), mc.GetFailureCount("not-exist"))

	mc.RecordStart("test-name")
	mc.RecordFailure("test-name", 50*time.Millisecond)
	assert.Equal(t, int64(1), mc.GetFailureCount("test-name"))
}

// TestMetricsCollectorGetRunningCount 测试获取运行中数量
func TestMetricsCollectorGetRunningCount(t *testing.T) {
	mc := NewMetricsCollector()

	// 不存在的 name 返回 0
	assert.Equal(t, int64(0), mc.GetRunningCount("not-exist"))

	mc.RecordStart("test-name")
	assert.Equal(t, int64(1), mc.GetRunningCount("test-name"))
}

// TestMetricsCollectorGetAvgExecutionTime 测试获取平均执行时间
func TestMetricsCollectorGetAvgExecutionTime(t *testing.T) {
	mc := NewMetricsCollector()

	// 不存在的 name 返回 0
	assert.Equal(t, float64(0), mc.GetAvgExecutionTime("not-exist"))

	mc.RecordStart("test-name")
	mc.RecordSuccess("test-name", 100*time.Millisecond)
	assert.Equal(t, float64(100), mc.GetAvgExecutionTime("test-name"))
}

// TestMetricsCollectorGetMaxExecutionTime 测试获取最大执行时间
func TestMetricsCollectorGetMaxExecutionTime(t *testing.T) {
	mc := NewMetricsCollector()

	// 不存在的 name 返回 0
	assert.Equal(t, int64(0), mc.GetMaxExecutionTime("not-exist"))

	mc.RecordStart("test-name")
	mc.RecordSuccess("test-name", 100*time.Millisecond)
	assert.Equal(t, int64(100), mc.GetMaxExecutionTime("test-name"))
}

// TestMetricsCollectorGetMinExecutionTime 测试获取最小执行时间
func TestMetricsCollectorGetMinExecutionTime(t *testing.T) {
	mc := NewMetricsCollector()

	// 不存在的 name 返回 0
	assert.Equal(t, int64(0), mc.GetMinExecutionTime("not-exist"))

	mc.RecordStart("test-name")
	mc.RecordSuccess("test-name", 100*time.Millisecond)
	assert.Equal(t, int64(100), mc.GetMinExecutionTime("test-name"))
}

// TestMetricsCollectorGetLastExecutionTime 测试获取最后执行时间
func TestMetricsCollectorGetLastExecutionTime(t *testing.T) {
	mc := NewMetricsCollector()

	// 不存在的 name 返回 0
	assert.Equal(t, int64(0), mc.GetLastExecutionTime("not-exist"))

	mc.RecordStart("test-name")
	mc.RecordSuccess("test-name", 100*time.Millisecond)
	assert.NotEqual(t, int64(0), mc.GetLastExecutionTime("test-name"))
}

// TestMetricsCollectorRecordSuccessWithoutStart 测试未调用 RecordStart 直接 RecordSuccess
func TestMetricsCollectorRecordSuccessWithoutStart(t *testing.T) {
	mc := NewMetricsCollector()

	// 未调用 RecordStart 直接 RecordSuccess，触发各字段的初始化分支
	mc.RecordSuccess("no-start", 50*time.Millisecond)
	assert.Equal(t, int64(1), mc.GetSuccessCount("no-start"))
}

// TestMetricsCollectorRecordFailureWithoutStart 测试未调用 RecordStart 直接 RecordFailure
func TestMetricsCollectorRecordFailureWithoutStart(t *testing.T) {
	mc := NewMetricsCollector()

	// 未调用 RecordStart 直接 RecordFailure，触发各字段的初始化分支
	mc.RecordFailure("no-start", 50*time.Millisecond)
	assert.Equal(t, int64(1), mc.GetFailureCount("no-start"))
}

// TestMetricsCollectorRecordRateLimited 测试记录被限流
func TestMetricsCollectorRecordRateLimited(t *testing.T) {
	mc := NewMetricsCollector()

	// 记录限流，应计入失败次数
	mc.RecordRateLimited("rate-limited")
	assert.Equal(t, int64(1), mc.GetFailureCount("rate-limited"))

	// 多次记录
	mc.RecordRateLimited("rate-limited")
	assert.Equal(t, int64(2), mc.GetFailureCount("rate-limited"))
}

// TestMetricsCollectorSuccessRateZero 测试成功率为 0 的情况
func TestMetricsCollectorSuccessRateZero(t *testing.T) {
	mc := NewMetricsCollector()

	// 不存在的 name 的成功率应为 0
	metrics := mc.GetMetrics("not-exist")
	assert.NotNil(t, metrics)
	assert.Equal(t, float64(0), metrics.SuccessRate)
}

// TestMetricsCollectorGlobalSuccessRateZero 测试全局成功率为 0 的情况
func TestMetricsCollectorGlobalSuccessRateZero(t *testing.T) {
	mc := NewMetricsCollector()

	// 没有任何执行记录时全局成功率应为 0
	global := mc.GetGlobalMetrics()
	assert.NotNil(t, global)
	assert.Equal(t, float64(0), global.SuccessRate)
}

// TestPrometheusExporter 测试 Prometheus 导出器
func TestPrometheusExporter(t *testing.T) {
	mc := NewMetricsCollector()

	// 记录一些数据
	mc.RecordStart("api-call")
	mc.RecordSuccess("api-call", 100*time.Millisecond)
	mc.RecordStart("api-call")
	mc.RecordFailure("api-call", 50*time.Millisecond)

	exporter := NewPrometheusExporter(mc, "myapp", "breaker")
	assert.NotNil(t, exporter)

	output := exporter.Export()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "myapp_breaker_total_executions")
	assert.Contains(t, output, "myapp_breaker_total_success")
	assert.Contains(t, output, "myapp_breaker_total_failure")
	assert.Contains(t, output, "myapp_breaker_active_count")
	assert.Contains(t, output, "execution_count")
	assert.Contains(t, output, "success_count")
	assert.Contains(t, output, "failure_count")
}

// TestPrometheusExporterEmpty 测试 Prometheus 导出器在没有数据时的输出
func TestPrometheusExporterEmpty(t *testing.T) {
	mc := NewMetricsCollector()
	exporter := NewPrometheusExporter(mc, "ns", "svc")

	output := exporter.Export()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "ns_svc_total_executions 0")
	assert.Contains(t, output, "ns_svc_total_success 0")
	assert.Contains(t, output, "ns_svc_total_failure 0")
	assert.Contains(t, output, "ns_svc_active_count 0")
}

// TestMetricsCollectorGlobalMetrics 测试全局指标聚合
func TestMetricsCollectorGlobalMetrics(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordStart("a")
	mc.RecordSuccess("a", 100*time.Millisecond)
	mc.RecordStart("b")
	mc.RecordFailure("b", 50*time.Millisecond)

	global := mc.GetGlobalMetrics()
	assert.Equal(t, int64(2), global.TotalExecutions)
	assert.Equal(t, int64(1), global.TotalSuccess)
	assert.Equal(t, int64(1), global.TotalFailure)
	assert.Equal(t, int64(0), global.ActiveCount)
	assert.InDelta(t, 50.0, global.SuccessRate, 0.1)
}
