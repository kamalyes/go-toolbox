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
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// SnowflakeGenerator Snowflake 分布式ID生成器
type SnowflakeGenerator struct {
	epoch      int64
	workerID   int64
	datacenter int64
	sequence   uint64
	lastTime   int64
	counter    uint64
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

// GenerateTraceID 生成跟踪ID（16进制格式，带时间排序特征）
// 格式: hex(snowflake_id)，16字符
// 示例: "018f5a3c7d2e4b10"
// 与 RequestID 的区别: 使用 hex 编码，更适合作为追踪标识
func (g *SnowflakeGenerator) GenerateTraceID() string {
	id := g.Generate()
	return fmt.Sprintf("%016x", id)
}

// GenerateSpanID 生成跨度ID（hex截取后8字符）
// 格式: hex(snowflake_id) 后8字符
// 示例: "7d2e4b10"
// 与 TraceID 的区别: 更短，仅保留低位部分，同 Trace 内唯一
func (g *SnowflakeGenerator) GenerateSpanID() string {
	id := g.Generate()
	return fmt.Sprintf("%08x", id&0xFFFFFFFF)
}

// GenerateRequestID 生成请求ID（纯数字+计数器后缀）
// 格式: snowflake数字-递增计数器
// 示例: "1732184000123456789-1"
// 与 TraceID 的区别: 纯数字+计数器，可排序，方便日志检索
func (g *SnowflakeGenerator) GenerateRequestID() string {
	counter := atomic.AddUint64(&g.counter, 1)
	id := g.Generate()

	var sb strings.Builder
	sb.Grow(28)
	sb.WriteString(strconv.FormatInt(id, 10))
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatUint(counter, 10))

	return sb.String()
}

// GenerateCorrelationID 生成关联ID（UUID 风格的 hex 格式）
// 格式: 8-4-4-4-4 hex（由 snowflake ID 拆分填充）
// 示例: "018f5a3c-7d2e-4b10-a1c3-e5f678901234"
// 与 TraceID 的区别: UUID 格式，含连字符，适合跨系统传递和关联
func (g *SnowflakeGenerator) GenerateCorrelationID() string {
	id := g.Generate()
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

	counter := atomic.AddUint64(&g.counter, 1)
	suffix := fmt.Sprintf("%012x", counter)
	copy(buf[24:36], suffix)

	return string(buf[:])
}

// Generate 生成 Snowflake ID
func (g *SnowflakeGenerator) Generate() int64 {
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
