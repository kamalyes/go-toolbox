/*
 * @Author: kamalyes 501893067@qq.com
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-20 09:00:00
 * @FilePath: \go-toolbox\pkg\syncx\timewheel_bench_test.go
 * @Description: 分片时间轮性能基准测试
 *   对比维度：
 *     - HashedWheelTimer Schedule/Refresh/Cancel 吞吐
 *     - time.AfterFunc 基线（标准库单次定时器）
 *     - 全量扫描基线（模拟 checkHeartbeat O(N) 扫描）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// HashedWheelTimer 基准测试
// ============================================================================

// BenchmarkTimerSchedule 测试无 key 调度吞吐
func BenchmarkTimerSchedule(b *testing.B) {
	timer := NewHashedWheelTimer()
	defer timer.Stop()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timer.Schedule(time.Second, func() {})
	}
}

// BenchmarkTimerScheduleWithKey 测试带 key 调度吞吐
func BenchmarkTimerScheduleWithKey(b *testing.B) {
	timer := NewHashedWheelTimer()
	defer timer.Stop()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timer.ScheduleWithKey(fmt.Sprintf("key-%d", i), time.Second, func() {})
	}
}

// BenchmarkTimerRefresh 测试刷新吞吐（心跳热路径）
func BenchmarkTimerRefresh(b *testing.B) {
	timer := NewHashedWheelTimer()
	defer timer.Stop()

	// 预热一个 key
	timer.ScheduleWithKey("bench-key", time.Second, func() {})

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		timer.Refresh("bench-key", time.Second, func() {})
	}
}

// BenchmarkTimerCancel 测试取消吞吐
func BenchmarkTimerCancel(b *testing.B) {
	timer := NewHashedWheelTimer()
	defer timer.Stop()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		task := timer.Schedule(time.Second, func() {})
		timer.Cancel(task)
	}
}

// BenchmarkTimerConcurrentSchedule 并发调度吞吐（多 goroutine）
func BenchmarkTimerConcurrentSchedule(b *testing.B) {
	timer := NewHashedWheelTimer()
	defer timer.Stop()

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			timer.Schedule(time.Second, func() {})
		}
	})
}

// BenchmarkTimerConcurrentRefresh 并发刷新吞吐（多 goroutine 同 key）
func BenchmarkTimerConcurrentRefresh(b *testing.B) {
	timer := NewHashedWheelTimer()
	defer timer.Stop()

	timer.ScheduleWithKey("bench-key", time.Second, func() {})

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			timer.Refresh("bench-key", time.Second, func() {})
		}
	})
}

// ============================================================================
// time.AfterFunc 基线对比
// ============================================================================

// BenchmarkTimeAfterFuncSchedule time.AfterFunc 调度基线
func BenchmarkTimeAfterFuncSchedule(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		t := time.AfterFunc(time.Second, func() {})
		t.Stop()
	}
}

// BenchmarkTimeAfterFuncConcurrent time.AfterFunc 并发调度基线
func BenchmarkTimeAfterFuncConcurrent(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			t := time.AfterFunc(time.Second, func() {})
			t.Stop()
		}
	})
}

// ============================================================================
// 全量扫描基线（模拟 checkHeartbeat O(N) 扫描）
// ============================================================================

// scanClient 模拟客户端心跳数据
type scanClient struct {
	id         string
	lastActive int64 // 纳秒时间戳
}

// scanRegistry 模拟旧版 checkHeartbeat 的全量扫描注册表
type scanRegistry struct {
	mu      sync.RWMutex
	clients map[string]*scanClient
}

func newScanRegistry(n int) *scanRegistry {
	r := &scanRegistry{clients: make(map[string]*scanClient, n)}
	now := time.Now().UnixNano()
	for i := 0; i < n; i++ {
		c := &scanClient{
			id:         fmt.Sprintf("client-%d", i),
			lastActive: now,
		}
		r.clients[c.id] = c
	}
	return r
}

// scanTimeouts 模拟 checkHeartbeat：遍历所有客户端检查超时
func (r *scanRegistry) scanTimeouts(timeout time.Duration) int {
	cutoff := time.Now().Add(-timeout).UnixNano()
	count := 0
	r.mu.RLock()
	for _, c := range r.clients {
		if c.lastActive < cutoff {
			count++
		}
	}
	r.mu.RUnlock()
	return count
}

// refreshScan 模拟旧版心跳更新（map 写操作）
func (r *scanRegistry) refreshScan(id string) {
	now := time.Now().UnixNano()
	r.mu.Lock()
	if c, ok := r.clients[id]; ok {
		c.lastActive = now
	}
	r.mu.Unlock()
}

// BenchmarkScanHeartbeat_10k 全量扫描 1 万连接
func BenchmarkScanHeartbeat_10k(b *testing.B) {
	r := newScanRegistry(10000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.scanTimeouts(90 * time.Second)
	}
}

// BenchmarkScanHeartbeat_100k 全量扫描 10 万连接
func BenchmarkScanHeartbeat_100k(b *testing.B) {
	r := newScanRegistry(100000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.scanTimeouts(90 * time.Second)
	}
}

// BenchmarkScanRefresh_10k 旧版心跳刷新（map 写锁）1 万连接
func BenchmarkScanRefresh_10k(b *testing.B) {
	r := newScanRegistry(10000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.refreshScan(fmt.Sprintf("client-%d", i%10000))
	}
}

// ============================================================================
// 综合对比：100 万连接场景模拟
// ============================================================================

// BenchmarkTimerMillionConnections 时间轮管理 100 万连接调度+刷新
func BenchmarkTimerMillionConnections(b *testing.B) {
	timer := NewHashedWheelTimer()
	defer timer.Stop()

	// 预填充 100 万任务
	const N = 1000000
	b.StopTimer()
	for i := 0; i < N; i++ {
		timer.ScheduleWithKey(fmt.Sprintf("conn-%d", i), 90*time.Second, func() {})
	}
	b.StartTimer()

	b.ResetTimer()
	b.ReportAllocs()
	// 测试刷新吞吐（心跳热路径）
	for i := 0; i < b.N; i++ {
		timer.Refresh(fmt.Sprintf("conn-%d", i%N), 90*time.Second, func() {})
	}
}

// BenchmarkScanMillionConnections 全量扫描 100 万连接
func BenchmarkScanMillionConnections(b *testing.B) {
	r := newScanRegistry(1000000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.scanTimeouts(90 * time.Second)
	}
}

// ============================================================================
// 辅助：确保 benchmark 中回调不干扰计时
// ============================================================================

// noopCounter 用于 benchmark 回调计数（避免编译器优化掉）
var benchCounter atomic.Int64
