/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-23 23:50:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-23 23:50:00
 * @FilePath: \go-toolbox\pkg\breaker\limiter_test.go
 * @Description: 限流器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package breaker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLimiter(t *testing.T) {
	limiter := NewLimiter(10, 100)

	assert.NotNil(t, limiter)
	assert.Equal(t, int32(10), limiter.rate)
	assert.Equal(t, int32(100), limiter.capacity)
	assert.Equal(t, int32(100), limiter.tokens)
}

func TestLimiterAllow(t *testing.T) {
	limiter := NewLimiter(10, 10)

	// 初始应该有10个令牌
	for i := 0; i < 10; i++ {
		assert.True(t, limiter.Allow(), "request %d should be allowed", i)
	}

	// 令牌耗尽
	assert.False(t, limiter.Allow(), "request should be denied when tokens exhausted")
}

func TestLimiterAllowN(t *testing.T) {
	limiter := NewLimiter(10, 10)

	// 消耗3个令牌
	assert.True(t, limiter.AllowN(3))
	assert.Equal(t, int32(7), limiter.tokens)

	// 消耗5个令牌
	assert.True(t, limiter.AllowN(5))
	assert.Equal(t, int32(2), limiter.tokens)

	// 尝试消耗5个令牌（不够）
	assert.False(t, limiter.AllowN(5))
	assert.Equal(t, int32(2), limiter.tokens)

	// 消耗2个令牌
	assert.True(t, limiter.AllowN(2))
	assert.Equal(t, int32(0), limiter.tokens)
}

func TestLimiterAllowNZero(t *testing.T) {
	limiter := NewLimiter(10, 10)

	// 消耗0个令牌应该总是成功
	assert.True(t, limiter.AllowN(0))
	assert.Equal(t, int32(10), limiter.tokens)
}

func TestLimiterRefill(t *testing.T) {
	limiter := NewLimiter(10, 10)

	// 消耗所有令牌
	for i := 0; i < 10; i++ {
		limiter.Allow()
	}
	assert.Equal(t, int32(0), limiter.tokens)

	// 等待令牌补充（1秒应该补充10个）
	time.Sleep(1100 * time.Millisecond)

	// 应该有新令牌
	assert.True(t, limiter.Allow())
}

func TestLimiterRefillPartial(t *testing.T) {
	limiter := NewLimiter(10, 10)

	// 消耗所有令牌
	for i := 0; i < 10; i++ {
		limiter.Allow()
	}

	// 等待0.5秒（应该补充5个令牌）
	time.Sleep(600 * time.Millisecond)

	// 应该能消耗大约5个令牌
	count := 0
	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			count++
		}
	}
	assert.GreaterOrEqual(t, count, 4, "should have at least 4 tokens")
	assert.LessOrEqual(t, count, 6, "should have at most 6 tokens")
}

func TestLimiterRefillCapacity(t *testing.T) {
	limiter := NewLimiter(10, 5)

	// 等待足够长时间让令牌超过容量
	time.Sleep(2 * time.Second)

	// 最多只能有capacity个令牌
	count := 0
	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			count++
		}
	}
	assert.Equal(t, 5, count, "should not exceed capacity")
}

func TestLimiterWait(t *testing.T) {
	limiter := NewLimiter(10, 1)

	// 消耗令牌
	limiter.Allow()

	ctx := context.Background()
	start := time.Now()

	// 等待直到有令牌
	err := limiter.Wait(ctx)

	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, elapsed, 80*time.Millisecond, "should wait for token refill")
}

func TestLimiterWaitWithCancel(t *testing.T) {
	limiter := NewLimiter(10, 1)

	// 消耗所有令牌
	limiter.Allow()

	ctx, cancel := context.WithCancel(context.Background())

	// 立即取消
	cancel()

	err := limiter.Wait(ctx)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestLimiterWaitWithTimeout(t *testing.T) {
	limiter := NewLimiter(1, 1)

	// 消耗令牌
	limiter.Allow()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := limiter.Wait(ctx)

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestLimiterGetAvailableTokens(t *testing.T) {
	limiter := NewLimiter(10, 10)

	assert.Equal(t, int32(10), limiter.GetAvailableTokens())

	limiter.AllowN(3)
	assert.Equal(t, int32(7), limiter.GetAvailableTokens())

	limiter.AllowN(7)
	assert.Equal(t, int32(0), limiter.GetAvailableTokens())
}

func TestLimiterGetAvailableTokensAfterRefill(t *testing.T) {
	limiter := NewLimiter(10, 10)

	// 消耗所有令牌
	limiter.AllowN(10)
	assert.Equal(t, int32(0), limiter.GetAvailableTokens())

	// 等待补充
	time.Sleep(500 * time.Millisecond)

	tokens := limiter.GetAvailableTokens()
	assert.GreaterOrEqual(t, tokens, int32(3), "should have refilled some tokens")
	assert.LessOrEqual(t, tokens, int32(7), "should not exceed expected refill")
}

func TestLimiterStats(t *testing.T) {
	limiter := NewLimiter(10, 100)

	stats := limiter.Stats()

	assert.Equal(t, int32(10), stats.Rate)
	assert.Equal(t, int32(100), stats.Capacity)
	assert.Equal(t, int32(100), stats.AvailableTokens)

	// 消耗一些令牌
	limiter.AllowN(30)

	stats = limiter.Stats()
	assert.Equal(t, int32(70), stats.AvailableTokens)
}

func TestLimiterConcurrentAllow(t *testing.T) {
	limiter := NewLimiter(100, 100)

	done := make(chan bool)
	workers := 10
	requestsPerWorker := 20

	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < requestsPerWorker; j++ {
				limiter.Allow()
			}
			done <- true
		}()
	}

	for i := 0; i < workers; i++ {
		<-done
	}

	// 验证令牌数不会变成负数
	tokens := limiter.GetAvailableTokens()
	assert.GreaterOrEqual(t, tokens, int32(0), "tokens should not be negative")
}

func TestLimiterConcurrentAllowN(t *testing.T) {
	limiter := NewLimiter(100, 100)

	done := make(chan bool)
	workers := 10

	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				limiter.AllowN(2)
			}
			done <- true
		}()
	}

	for i := 0; i < workers; i++ {
		<-done
	}

	tokens := limiter.GetAvailableTokens()
	assert.GreaterOrEqual(t, tokens, int32(0))
}

func TestLimiterConcurrentWait(t *testing.T) {
	limiter := NewLimiter(10, 5)

	done := make(chan bool)
	workers := 5

	for i := 0; i < workers; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			limiter.Wait(ctx)
			done <- true
		}()
	}

	for i := 0; i < workers; i++ {
		<-done
	}

	assert.True(t, true, "all workers completed")
}

func TestLimiterHighRate(t *testing.T) {
	limiter := NewLimiter(1000, 1000)

	// 初始应该能快速消耗所有令牌
	count := 0
	for i := 0; i < 1000; i++ {
		if limiter.Allow() {
			count++
		}
	}
	assert.Equal(t, 1000, count)

	// 令牌耗尽
	assert.False(t, limiter.Allow())

	// 等待1秒应该补充1000个
	time.Sleep(1100 * time.Millisecond)

	count = 0
	for i := 0; i < 1000; i++ {
		if limiter.Allow() {
			count++
		}
	}
	assert.GreaterOrEqual(t, count, 900, "should refill close to 1000 tokens")
}

func TestLimiterLowRate(t *testing.T) {
	limiter := NewLimiter(1, 5)

	// 消耗所有令牌
	for i := 0; i < 5; i++ {
		limiter.Allow()
	}

	// 等待2秒应该补充2个令牌
	time.Sleep(2100 * time.Millisecond)

	count := 0
	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			count++
		}
	}
	assert.GreaterOrEqual(t, count, 1, "should have at least 1 token")
	assert.LessOrEqual(t, count, 3, "should have at most 3 tokens")
}

func TestLimiterZeroRate(t *testing.T) {
	limiter := NewLimiter(0, 10)

	// 消耗所有初始令牌
	for i := 0; i < 10; i++ {
		limiter.Allow()
	}

	// 等待一段时间
	time.Sleep(500 * time.Millisecond)

	// 不应该有新令牌（因为rate为0）
	assert.False(t, limiter.Allow())
}
