/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-24 00:30:55
 * @FilePath: \go-toolbox\pkg\breaker\metrics.go
 * @Description: 通用指标收集器
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package breaker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// MetricsCollector 通用指标收集器
type MetricsCollector struct {
	// 执行统计
	executionCount    map[string]*int64 // 执行总次数
	successCount      map[string]*int64 // 成功次数
	failureCount      map[string]*int64 // 失败次数
	runningCount      map[string]*int64 // 当前运行中的数量
	totalDuration     map[string]*int64 // 总执行时间(毫秒)
	lastExecutionTime map[string]*int64 // 最后执行时间戳
	// 全局统计
	totalExecutions int64
	totalSuccess    int64
	totalFailure    int64
	activeCount     int64

	// 性能统计
	avgExecutionTime map[string]float64
	maxExecutionTime map[string]int64
	minExecutionTime map[string]int64

	mu sync.RWMutex
}

// GetExecutionCount 获取指定名称的执行次数
func (mc *MetricsCollector) GetExecutionCount(name string) int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	if mc.executionCount[name] == nil {
		return 0
	}
	return atomic.LoadInt64(mc.executionCount[name])
}

// GetSuccessCount 获取指定名称的成功次数
func (mc *MetricsCollector) GetSuccessCount(name string) int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	if mc.successCount[name] == nil {
		return 0
	}
	return atomic.LoadInt64(mc.successCount[name])
}

// GetFailureCount 获取指定名称的失败次数
func (mc *MetricsCollector) GetFailureCount(name string) int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	if mc.failureCount[name] == nil {
		return 0
	}
	return atomic.LoadInt64(mc.failureCount[name])
}

// GetRunningCount 获取指定名称的当前运行数
func (mc *MetricsCollector) GetRunningCount(name string) int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	if mc.runningCount[name] == nil {
		return 0
	}
	return atomic.LoadInt64(mc.runningCount[name])
}

// GetAvgExecutionTime 获取指定名称的平均执行时间
func (mc *MetricsCollector) GetAvgExecutionTime(name string) float64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.avgExecutionTime[name]
}

// GetMaxExecutionTime 获取指定名称的最大执行时间
func (mc *MetricsCollector) GetMaxExecutionTime(name string) int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.maxExecutionTime[name]
}

// GetMinExecutionTime 获取指定名称的最小执行时间
func (mc *MetricsCollector) GetMinExecutionTime(name string) int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.minExecutionTime[name]
}

// GetLastExecutionTime 获取指定名称的最后执行时间
func (mc *MetricsCollector) GetLastExecutionTime(name string) int64 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	if mc.lastExecutionTime[name] == nil {
		return 0
	}
	return atomic.LoadInt64(mc.lastExecutionTime[name])
}

// GetTotalExecutions 获取全局执行次数
func (mc *MetricsCollector) GetTotalExecutions() int64 {
	return atomic.LoadInt64(&mc.totalExecutions)
}

// GetTotalSuccess 获取全局成功次数
func (mc *MetricsCollector) GetTotalSuccess() int64 {
	return atomic.LoadInt64(&mc.totalSuccess)
}

// GetTotalFailure 获取全局失败次数
func (mc *MetricsCollector) GetTotalFailure() int64 {
	return atomic.LoadInt64(&mc.totalFailure)
}

// GetActiveCount 获取全局活跃数
func (mc *MetricsCollector) GetActiveCount() int64 {
	return atomic.LoadInt64(&mc.activeCount)
}

// NewMetricsCollector 创建指标收集器
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		executionCount:    make(map[string]*int64),
		successCount:      make(map[string]*int64),
		failureCount:      make(map[string]*int64),
		runningCount:      make(map[string]*int64),
		totalDuration:     make(map[string]*int64),
		lastExecutionTime: make(map[string]*int64),
		avgExecutionTime:  make(map[string]float64),
		maxExecutionTime:  make(map[string]int64),
		minExecutionTime:  make(map[string]int64),
	}
}

// RecordStart 记录开始执行
func (mc *MetricsCollector) RecordStart(name string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// 初始化计数器
	if mc.executionCount[name] == nil {
		var count int64 = 0
		mc.executionCount[name] = &count
	}
	if mc.runningCount[name] == nil {
		var count int64 = 0
		mc.runningCount[name] = &count
	}

	atomic.AddInt64(mc.executionCount[name], 1)
	atomic.AddInt64(mc.runningCount[name], 1)
	atomic.AddInt64(&mc.totalExecutions, 1)
	atomic.AddInt64(&mc.activeCount, 1)
}

// RecordSuccess 记录执行成功
func (mc *MetricsCollector) RecordSuccess(name string, duration time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.successCount[name] == nil {
		var count int64 = 0
		mc.successCount[name] = &count
	}
	if mc.totalDuration[name] == nil {
		var dur int64 = 0
		mc.totalDuration[name] = &dur
	}
	if mc.lastExecutionTime[name] == nil {
		var t int64 = 0
		mc.lastExecutionTime[name] = &t
	}
	if mc.runningCount[name] == nil {
		var count int64 = 0
		mc.runningCount[name] = &count
	}

	atomic.AddInt64(mc.successCount[name], 1)
	atomic.AddInt64(mc.runningCount[name], -1)
	atomic.AddInt64(mc.totalDuration[name], duration.Milliseconds())
	atomic.StoreInt64(mc.lastExecutionTime[name], time.Now().Unix())
	atomic.AddInt64(&mc.totalSuccess, 1)
	atomic.AddInt64(&mc.activeCount, -1)

	// 更新执行时间统计
	mc.updateExecutionTimeStats(name, duration.Milliseconds())
}

// RecordFailure 记录执行失败
func (mc *MetricsCollector) RecordFailure(name string, duration time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.failureCount[name] == nil {
		var count int64 = 0
		mc.failureCount[name] = &count
	}
	if mc.runningCount[name] == nil {
		var count int64 = 0
		mc.runningCount[name] = &count
	}

	atomic.AddInt64(mc.failureCount[name], 1)
	atomic.AddInt64(mc.runningCount[name], -1)
	atomic.AddInt64(&mc.totalFailure, 1)
	atomic.AddInt64(&mc.activeCount, -1)

	// 失败也记录执行时间
	mc.updateExecutionTimeStats(name, duration.Milliseconds())
}

// RecordRateLimited 记录被限流
func (mc *MetricsCollector) RecordRateLimited(name string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.failureCount[name] == nil {
		var count int64 = 0
		mc.failureCount[name] = &count
	}

	// 被限流视为执行失败的一种
	atomic.AddInt64(mc.failureCount[name], 1)
}

// updateExecutionTimeStats 更新执行时间统计
func (mc *MetricsCollector) updateExecutionTimeStats(name string, durationMs int64) {
	// 确保字段已初始化
	if mc.executionCount[name] == nil || mc.totalDuration[name] == nil {
		return
	}

	// 更新平均值
	execCount := atomic.LoadInt64(mc.executionCount[name])
	if execCount > 0 {
		totalDur := atomic.LoadInt64(mc.totalDuration[name])
		mc.avgExecutionTime[name] = float64(totalDur) / float64(execCount)
	}

	// 更新最大值
	if durationMs > mc.maxExecutionTime[name] {
		mc.maxExecutionTime[name] = durationMs
	}

	// 更新最小值
	if mc.minExecutionTime[name] == 0 || durationMs < mc.minExecutionTime[name] {
		mc.minExecutionTime[name] = durationMs
	}
}

// GetMetrics 获取单个指标
func (mc *MetricsCollector) GetMetrics(name string) *Metrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	return &Metrics{
		Name:              name,
		ExecutionCount:    mc.getInt64Value(mc.executionCount[name]),
		SuccessCount:      mc.getInt64Value(mc.successCount[name]),
		FailureCount:      mc.getInt64Value(mc.failureCount[name]),
		RunningCount:      mc.getInt64Value(mc.runningCount[name]),
		AvgExecutionTime:  mc.avgExecutionTime[name],
		MaxExecutionTime:  mc.maxExecutionTime[name],
		MinExecutionTime:  mc.minExecutionTime[name],
		LastExecutionTime: mc.getInt64Value(mc.lastExecutionTime[name]),
		SuccessRate:       mc.calculateSuccessRate(name),
	}
}

// GetAllMetrics 获取所有指标
func (mc *MetricsCollector) GetAllMetrics() map[string]*Metrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metrics := make(map[string]*Metrics)
	for name := range mc.executionCount {
		metrics[name] = mc.GetMetrics(name)
	}
	return metrics
}

// GetGlobalMetrics 获取全局统计
func (mc *MetricsCollector) GetGlobalMetrics() *GlobalMetrics {
	return &GlobalMetrics{
		TotalExecutions: atomic.LoadInt64(&mc.totalExecutions),
		TotalSuccess:    atomic.LoadInt64(&mc.totalSuccess),
		TotalFailure:    atomic.LoadInt64(&mc.totalFailure),
		ActiveCount:     atomic.LoadInt64(&mc.activeCount),
		SuccessRate:     mc.calculateGlobalSuccessRate(),
	}
}

// calculateSuccessRate 计算成功率
func (mc *MetricsCollector) calculateSuccessRate(name string) float64 {
	execCount := mc.getInt64Value(mc.executionCount[name])
	if execCount == 0 {
		return 0
	}
	successCount := mc.getInt64Value(mc.successCount[name])
	return float64(successCount) / float64(execCount) * 100
}

// calculateGlobalSuccessRate 计算全局成功率
func (mc *MetricsCollector) calculateGlobalSuccessRate() float64 {
	total := atomic.LoadInt64(&mc.totalExecutions)
	if total == 0 {
		return 0
	}
	success := atomic.LoadInt64(&mc.totalSuccess)
	return float64(success) / float64(total) * 100
}

// getInt64Value 安全获取int64值
func (mc *MetricsCollector) getInt64Value(ptr *int64) int64 {
	if ptr == nil {
		return 0
	}
	return atomic.LoadInt64(ptr)
}

// Reset 重置指标
func (mc *MetricsCollector) Reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.executionCount = make(map[string]*int64)
	mc.successCount = make(map[string]*int64)
	mc.failureCount = make(map[string]*int64)
	mc.runningCount = make(map[string]*int64)
	mc.totalDuration = make(map[string]*int64)
	mc.lastExecutionTime = make(map[string]*int64)
	mc.avgExecutionTime = make(map[string]float64)
	mc.maxExecutionTime = make(map[string]int64)
	mc.minExecutionTime = make(map[string]int64)

	atomic.StoreInt64(&mc.totalExecutions, 0)
	atomic.StoreInt64(&mc.totalSuccess, 0)
	atomic.StoreInt64(&mc.totalFailure, 0)
	atomic.StoreInt64(&mc.activeCount, 0)
}

// Metrics 通用指标
type Metrics struct {
	Name              string  `json:"name"`
	ExecutionCount    int64   `json:"execution_count"`
	SuccessCount      int64   `json:"success_count"`
	FailureCount      int64   `json:"failure_count"`
	RunningCount      int64   `json:"running_count"`
	AvgExecutionTime  float64 `json:"avg_execution_time_ms"`
	MaxExecutionTime  int64   `json:"max_execution_time_ms"`
	MinExecutionTime  int64   `json:"min_execution_time_ms"`
	LastExecutionTime int64   `json:"last_execution_time"`
	SuccessRate       float64 `json:"success_rate"`
}

// GlobalMetrics 全局指标
type GlobalMetrics struct {
	TotalExecutions int64   `json:"total_executions"`
	TotalSuccess    int64   `json:"total_success"`
	TotalFailure    int64   `json:"total_failure"`
	ActiveCount     int64   `json:"active_count"`
	SuccessRate     float64 `json:"success_rate"`
}

// MetricsSnapshot 指标快照,用于导出
type MetricsSnapshot struct {
	GlobalMetrics *GlobalMetrics      `json:"global_metrics"`
	Metrics       map[string]*Metrics `json:"metrics"`
	Timestamp     int64               `json:"timestamp"`
}

// GetSnapshot 获取指标快照
func (mc *MetricsCollector) GetSnapshot() *MetricsSnapshot {
	return &MetricsSnapshot{
		GlobalMetrics: mc.GetGlobalMetrics(),
		Metrics:       mc.GetAllMetrics(),
		Timestamp:     time.Now().Unix(),
	}
}

// PrometheusExporter Prometheus格式导出器
type PrometheusExporter struct {
	collector *MetricsCollector
	namespace string
	service   string
}

// NewPrometheusExporter 创建Prometheus导出器
func NewPrometheusExporter(collector *MetricsCollector, namespace, service string) *PrometheusExporter {
	return &PrometheusExporter{
		collector: collector,
		namespace: namespace,
		service:   service,
	}
}

// Export 导出Prometheus格式指标
func (pe *PrometheusExporter) Export() string {
	var output string
	prefix := fmt.Sprintf("%s_%s_", pe.namespace, pe.service)

	// 导出全局指标
	global := pe.collector.GetGlobalMetrics()
	output += fmt.Sprintf("# HELP %stotal_executions Total number of executions\n", prefix)
	output += fmt.Sprintf("# TYPE %stotal_executions counter\n", prefix)
	output += fmt.Sprintf("%stotal_executions %d\n", prefix, global.TotalExecutions)

	output += fmt.Sprintf("# HELP %stotal_success Total number of successful executions\n", prefix)
	output += fmt.Sprintf("# TYPE %stotal_success counter\n", prefix)
	output += fmt.Sprintf("%stotal_success %d\n", prefix, global.TotalSuccess)

	output += fmt.Sprintf("# HELP %stotal_failure Total number of failed executions\n", prefix)
	output += fmt.Sprintf("# TYPE %stotal_failure counter\n", prefix)
	output += fmt.Sprintf("%stotal_failure %d\n", prefix, global.TotalFailure)

	output += fmt.Sprintf("# HELP %sactive_count Number of currently active executions\n", prefix)
	output += fmt.Sprintf("# TYPE %sactive_count gauge\n", prefix)
	output += fmt.Sprintf("%sactive_count %d\n", prefix, global.ActiveCount)

	// 导出所有指标
	allMetrics := pe.collector.GetAllMetrics()
	for name, metrics := range allMetrics {
		output += fmt.Sprintf("%sexecution_count{name=\"%s\"} %d\n", prefix, name, metrics.ExecutionCount)
		output += fmt.Sprintf("%ssuccess_count{name=\"%s\"} %d\n", prefix, name, metrics.SuccessCount)
		output += fmt.Sprintf("%sfailure_count{name=\"%s\"} %d\n", prefix, name, metrics.FailureCount)
		output += fmt.Sprintf("%savg_execution_time_ms{name=\"%s\"} %.2f\n", prefix, name, metrics.AvgExecutionTime)
		output += fmt.Sprintf("%ssuccess_rate{name=\"%s\"} %.2f\n", prefix, name, metrics.SuccessRate)
	}

	return output
}
