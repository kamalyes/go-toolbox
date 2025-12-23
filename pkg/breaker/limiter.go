/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-23 23:50:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-23 23:50:00
 * @FilePath: \go-toolbox\pkg\breaker\limiter.go
 * @Description: 限流器实现（令牌桶算法）
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package breaker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Limiter 限流器(令牌桶算法)
type Limiter struct {
	rate       int32 // 每秒令牌数
	capacity   int32 // 桶容量
	tokens     int32 // 当前令牌数
	lastUpdate int64 // 最后更新时间(纳秒)
	mu         sync.Mutex
}

// NewLimiter 创建限流器
func NewLimiter(rate, capacity int32) *Limiter {
	return &Limiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastUpdate: time.Now().UnixNano(),
	}
}

// AllowN 是否允许N个请求
func (l *Limiter) AllowN(n int32) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UnixNano()
	l.refill(now)

	if atomic.LoadInt32(&l.tokens) >= n {
		atomic.AddInt32(&l.tokens, -n)
		return true
	}
	return false
}

// Allow 是否允许单个请求
func (l *Limiter) Allow() bool {
	return l.AllowN(1)
}

// Wait 等待直到可以执行
func (l *Limiter) Wait(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()

	for {
		if l.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			continue
		}
	}
}

// refill 补充令牌
func (l *Limiter) refill(now int64) {
	lastUpdate := atomic.LoadInt64(&l.lastUpdate)
	elapsed := time.Duration(now - lastUpdate)

	// 计算应该添加的令牌数
	tokensToAdd := int32(elapsed.Seconds() * float64(l.rate))
	if tokensToAdd > 0 {
		currentTokens := atomic.LoadInt32(&l.tokens)
		newTokens := currentTokens + tokensToAdd
		if newTokens > l.capacity {
			newTokens = l.capacity
		}
		atomic.StoreInt32(&l.tokens, newTokens)
		atomic.StoreInt64(&l.lastUpdate, now)
	}
}

// GetAvailableTokens 获取可用令牌数
func (l *Limiter) GetAvailableTokens() int32 {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill(time.Now().UnixNano())
	return atomic.LoadInt32(&l.tokens)
}

// Stats 获取统计信息
func (l *Limiter) Stats() LimiterStats {
	return LimiterStats{
		Rate:            l.rate,
		Capacity:        l.capacity,
		AvailableTokens: l.GetAvailableTokens(),
	}
}

// LimiterStats 限流器统计
type LimiterStats struct {
	Rate            int32
	Capacity        int32
	AvailableTokens int32
}
