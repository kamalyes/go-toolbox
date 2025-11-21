/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:00:00
 * @FilePath: \go-toolbox\pkg\idgen\snowflake.go
 * @Description: Snowflake 分布式ID生成器（高性能版本）
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"fmt"
	"sync"
	"time"
)

// SnowflakeGenerator Snowflake 分布式ID生成器
type SnowflakeGenerator struct {
	epoch      int64
	workerID   int64
	datacenter int64
	sequence   uint64
	lastTime   int64
	mu         sync.Mutex
}

// NewSnowflakeGenerator 创建 Snowflake 生成器
func NewSnowflakeGenerator(workerID, datacenter int64) *SnowflakeGenerator {
	return &SnowflakeGenerator{
		epoch:      1640995200000, // 2022-01-01 00:00:00
		workerID:   workerID & 0x1F,
		datacenter: datacenter & 0x1F,
		sequence:   0,
	}
}

// GenerateTraceID 生成跟踪ID
func (g *SnowflakeGenerator) GenerateTraceID() string {
	return fmt.Sprintf("%d", g.generateSnowflake())
}

// GenerateSpanID 生成跨度ID
func (g *SnowflakeGenerator) GenerateSpanID() string {
	return fmt.Sprintf("%d", g.generateSnowflake())
}

// GenerateRequestID 生成请求ID
func (g *SnowflakeGenerator) GenerateRequestID() string {
	return fmt.Sprintf("%d", g.generateSnowflake())
}

// GenerateCorrelationID 生成关联ID
func (g *SnowflakeGenerator) GenerateCorrelationID() string {
	return fmt.Sprintf("%d", g.generateSnowflake())
}

// generateSnowflake 生成 Snowflake ID
func (g *SnowflakeGenerator) generateSnowflake() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < g.lastTime {
		now = g.lastTime
	}

	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & 0xFFF
		if g.sequence == 0 {
			for now <= g.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	timestamp := (now - g.epoch) << 22
	worker := g.workerID << 17
	dc := g.datacenter << 12

	return timestamp | worker | dc | int64(g.sequence)
}
