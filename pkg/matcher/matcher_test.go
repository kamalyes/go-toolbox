/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-15 02:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-15 02:15:15
 * @FilePath: \go-toolbox\pkg\matcher\matcher_test.go
 * @Description: 匹配器测试 - 验证并发安全性和性能
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package matcher

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/contextx"
	"github.com/stretchr/testify/assert"
)

// TestBasicMatch 基本匹配测试
func TestBasicMatch(t *testing.T) {
	type Action struct {
		Name string
		Code int
	}

	m := NewMatcher[*Action]()

	// 添加规则
	m.AddRule(
		NewChainRule(&Action{Name: "admin", Code: 1}).
			When(MatchString("role", "admin")).
			WithPriority(100),
	)

	m.AddRule(
		NewChainRule(&Action{Name: "user", Code: 2}).
			When(MatchString("role", "user")).
			WithPriority(50),
	)

	// 测试匹配
	ctx := contextx.NewContext().WithValue("role", "admin")
	result, ok := m.Match(ctx)

	assert.True(t, ok, "Expected match but got none")
	assert.Equal(t, "admin", result.Name)
}

// TestConcurrentMatch 并发匹配测试
func TestConcurrentMatch(t *testing.T) {
	type Result struct {
		Value int
	}

	m := NewMatcher[*Result]()

	// 添加多个规则
	for i := 0; i < 100; i++ {
		priority := i
		value := i
		m.AddRule(
			NewChainRule(&Result{Value: value}).
				When(MatchString("key", fmt.Sprintf("value-%d", i))).
				WithPriority(priority),
		)
	}

	// 并发测试
	const goroutines = 1000
	const iterations = 100

	var wg sync.WaitGroup
	var successCount atomic.Int64
	var failCount atomic.Int64

	start := time.Now()

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				ctx := contextx.NewContext().WithValue("key", fmt.Sprintf("value-%d", i%100))
				result, ok := m.Match(ctx)

				if ok {
					assert.Equal(t, i%100, result.Value)
					successCount.Add(1)
				} else {
					failCount.Add(1)
				}
			}
		}(g)
	}

	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("Concurrent test completed:")
	t.Logf("  Goroutines: %d", goroutines)
	t.Logf("  Iterations: %d", iterations)
	t.Logf("  Total ops: %d", goroutines*iterations)
	t.Logf("  Success: %d", successCount.Load())
	t.Logf("  Failed: %d", failCount.Load())
	t.Logf("  Duration: %v", elapsed)
	t.Logf("  Ops/sec: %.2f", float64(goroutines*iterations)/elapsed.Seconds())

	assert.Equal(t, int64(goroutines*iterations), successCount.Load())
}

// TestConcurrentAddAndMatch 并发添加和匹配测试
func TestConcurrentAddAndMatch(t *testing.T) {
	type Result struct {
		ID int
	}

	m := NewMatcher[*Result]()

	var wg sync.WaitGroup
	const addGoroutines = 100
	const matchGoroutines = 500

	// 并发添加规则
	for i := 0; i < addGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			m.AddRule(
				NewChainRule(&Result{ID: id}).
					When(MatchString("id", fmt.Sprintf("id-%d", id))).
					WithPriority(id),
			)
		}(i)
	}

	// 等待添加完成
	wg.Wait()

	// 并发匹配
	var matchSuccess atomic.Int64
	for i := 0; i < matchGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx := contextx.NewContext().WithValue("id", fmt.Sprintf("id-%d", id%addGoroutines))
			if _, ok := m.Match(ctx); ok {
				matchSuccess.Add(1)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Concurrent add and match:")
	t.Logf("  Rules added: %d", addGoroutines)
	t.Logf("  Match attempts: %d", matchGoroutines)
	t.Logf("  Match success: %d", matchSuccess.Load())

	assert.Equal(t, int64(matchGoroutines), matchSuccess.Load())
}

// TestMatchWithCache 缓存测试
func TestMatchWithCache(t *testing.T) {
	type Result struct {
		Value string
	}

	m := NewMatcher[*Result]().EnableCache(time.Second)

	m.AddRule(
		NewChainRule(&Result{Value: "cached"}).
			When(MatchString("key", "value")).
			WithPriority(100),
	)

	ctx := contextx.NewContext().WithValue("key", "value")

	// 第一次匹配（缓存未命中）
	result1, ok1 := m.Match(ctx)
	assert.True(t, ok1, "First match failed")

	stats1 := m.Stats()
	assert.Equal(t, int64(0), stats1["cache_hits"], "First match should not hit cache")

	// 第二次匹配（缓存命中）
	result2, ok2 := m.Match(ctx)
	assert.True(t, ok2, "Second match failed")

	stats2 := m.Stats()
	assert.Equal(t, int64(1), stats2["cache_hits"])

	assert.Equal(t, result1.Value, result2.Value, "Cached result differs from original")

	t.Logf("Cache test stats: %+v", stats2)
}

// TestConcurrentCacheAccess 并发缓存访问测试
func TestConcurrentCacheAccess(t *testing.T) {
	type Result struct {
		Count int
	}

	var callCount atomic.Int64

	m := NewMatcher[*Result]().EnableCache(time.Second)

	// 添加规则，计数调用次数
	m.AddRule(
		NewChainRule(&Result{Count: int(callCount.Add(1))}).
			When(func(ctx *contextx.Context) bool {
				callCount.Add(1)
				return contextx.Get[string](ctx, "key") == "test"
			}).
			WithPriority(100),
	)

	const goroutines = 1000
	var wg sync.WaitGroup
	ctx := contextx.NewContext().WithValue("key", "test")

	start := time.Now()

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Match(ctx)
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	stats := m.Stats()
	t.Logf("Concurrent cache test:")
	t.Logf("  Goroutines: %d", goroutines)
	t.Logf("  Duration: %v", elapsed)
	t.Logf("  Cache hits: %d", stats["cache_hits"])
	t.Logf("  Cache misses: %d", stats["cache_misses"])
	t.Logf("  Call count: %d", callCount.Load())

	// 应该只有第一次调用匹配函数，其他都是缓存命中
	assert.LessOrEqual(t, callCount.Load(), int64(10), "Cache not working properly")
}

// TestMiddleware 中间件测试
func TestMiddleware(t *testing.T) {
	type Result struct {
		Value string
	}

	m := NewMatcher[*Result]()

	var middlewareCalls atomic.Int64

	// 添加中间件
	m.Use(func(ctx *contextx.Context, next func() (*Result, bool)) (*Result, bool) {
		middlewareCalls.Add(1)
		ctx.WithMetadata("middleware", "called")
		return next()
	})

	m.AddRule(
		NewChainRule(&Result{Value: "test"}).
			When(MatchString("key", "value")).
			WithPriority(100),
	)

	ctx := contextx.NewContext().WithValue("key", "value")
	result, ok := m.Match(ctx)

	assert.True(t, ok, "Match failed")
	assert.Equal(t, "test", result.Value)
	assert.Equal(t, int64(1), middlewareCalls.Load())
	assert.Equal(t, "called", ctx.GetMetadata("middleware"))
}

// TestContextTimeout 上下文超时测试
func TestContextTimeout(t *testing.T) {
	type Result struct {
		Value string
	}

	m := NewMatcher[*Result]()

	m.AddRule(
		NewChainRule(&Result{Value: "test"}).
			When(func(ctx *contextx.Context) bool {
				// 模拟慢速匹配
				time.Sleep(100 * time.Millisecond)
				return true
			}).
			WithPriority(100),
	)

	ctx := contextx.NewContext().
		WithValue("key", "value").
		WithTimeout(10 * time.Millisecond)

	time.Sleep(20 * time.Millisecond)

	_, ok := m.Match(ctx)
	assert.False(t, ok, "Match should fail due to timeout")
}

// TestPriority 优先级测试
func TestPriority(t *testing.T) {
	type Result struct {
		Priority int
	}

	m := NewMatcher[*Result]()

	// 添加不同优先级的规则
	m.AddRule(
		NewChainRule(&Result{Priority: 10}).
			When(func(*contextx.Context) bool { return true }).
			WithPriority(10),
	)

	m.AddRule(
		NewChainRule(&Result{Priority: 100}).
			When(func(*contextx.Context) bool { return true }).
			WithPriority(100),
	)

	m.AddRule(
		NewChainRule(&Result{Priority: 50}).
			When(func(*contextx.Context) bool { return true }).
			WithPriority(50),
	)

	ctx := contextx.NewContext()
	result, ok := m.Match(ctx)

	assert.True(t, ok, "Match failed")
	// 应该匹配到优先级最高的
	assert.Equal(t, 100, result.Priority)
}

// BenchmarkMatch 基准测试
func BenchmarkMatch(b *testing.B) {
	type Result struct {
		Value int
	}

	m := NewMatcher[*Result]()

	// 添加100个规则
	for i := 0; i < 100; i++ {
		value := i
		m.AddRule(
			NewChainRule(&Result{Value: value}).
				When(MatchString("key", fmt.Sprintf("value-%d", i))).
				WithPriority(i),
		)
	}

	ctx := contextx.NewContext().WithValue("key", "value-50")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Match(ctx)
		}
	})
}

// BenchmarkMatchWithCache 带缓存的基准测试
func BenchmarkMatchWithCache(b *testing.B) {
	type Result struct {
		Value int
	}

	m := NewMatcher[*Result]().EnableCache(time.Minute)

	for i := 0; i < 100; i++ {
		value := i
		m.AddRule(
			NewChainRule(&Result{Value: value}).
				When(MatchString("key", fmt.Sprintf("value-%d", i))).
				WithPriority(i),
		)
	}

	ctx := contextx.NewContext().WithValue("key", "value-50")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Match(ctx)
		}
	})
}

// BenchmarkConcurrentAddAndMatch 并发添加和匹配基准测试
func BenchmarkConcurrentAddAndMatch(b *testing.B) {
	type Result struct {
		ID int
	}

	b.RunParallel(func(pb *testing.PB) {
		m := NewMatcher[*Result]()

		// 预填充一些规则
		for i := 0; i < 50; i++ {
			id := i
			m.AddRule(
				NewChainRule(&Result{ID: id}).
					When(MatchString("id", fmt.Sprintf("id-%d", id))).
					WithPriority(id),
			)
		}

		i := 0
		for pb.Next() {
			// 交替添加规则和匹配
			if i%2 == 0 {
				m.AddRule(
					NewChainRule(&Result{ID: i}).
						When(MatchString("id", fmt.Sprintf("id-%d", i))).
						WithPriority(i),
				)
			} else {
				ctx := contextx.NewContext().WithValue("id", fmt.Sprintf("id-%d", i%50))
				m.Match(ctx)
			}
			i++
		}
	})
}

// TestStressTest 压力测试
func TestStressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	type Result struct {
		ID int
	}

	m := NewMatcher[*Result]().EnableCache(time.Second)

	const numRules = 1000
	const numGoroutines = 100
	const numIterations = 1000

	// 添加大量规则
	for i := 0; i < numRules; i++ {
		id := i
		m.AddRule(
			NewChainRule(&Result{ID: id}).
				When(MatchString("key", fmt.Sprintf("value-%d", id))).
				WithPriority(id),
		)
	}

	var totalOps atomic.Int64
	var successOps atomic.Int64
	var wg sync.WaitGroup

	start := time.Now()

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()

			for i := 0; i < numIterations; i++ {
				ctx := contextx.NewContext().WithValue("key", fmt.Sprintf("value-%d", i%numRules))
				if _, ok := m.Match(ctx); ok {
					successOps.Add(1)
				}
				totalOps.Add(1)
			}
		}(g)
	}

	wg.Wait()
	elapsed := time.Since(start)

	stats := m.Stats()

	t.Logf("Stress test results:")
	t.Logf("  Rules: %d", numRules)
	t.Logf("  Goroutines: %d", numGoroutines)
	t.Logf("  Iterations per goroutine: %d", numIterations)
	t.Logf("  Total operations: %d", totalOps.Load())
	t.Logf("  Successful matches: %d", successOps.Load())
	t.Logf("  Duration: %v", elapsed)
	t.Logf("  Ops/sec: %.2f", float64(totalOps.Load())/elapsed.Seconds())
	t.Logf("  Stats: %+v", stats)

	assert.Equal(t, totalOps.Load(), successOps.Load(), "All operations should succeed")
}
