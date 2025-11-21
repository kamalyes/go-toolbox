/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:30:00
 * @FilePath: \go-toolbox\pkg\idgen\shortflake.go
 * @Description: ShortFlake 短 ID 生成器（高性能、紧凑型）
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"fmt"
	"sync"
	"time"
)

// ShortFlakeGenerator 短 Snowflake 生成器
// 生成 53 位整数（JavaScript Number 安全范围）
// 时间戳(41位) + 机器ID(6位) + 序列号(6位) = 53位
// 最大值: 9007199254740991 (约 9PB，16位数字)
type ShortFlakeGenerator struct {
	epoch    int64  // 自定义纪元（毫秒）
	nodeID   int64  // 节点ID (0-63)
	sequence uint64 // 序列号
	lastTime int64  // 上次生成时间
	mu       sync.Mutex
}

// NewShortFlakeGenerator 创建短 Snowflake 生成器
// nodeID: 0-63 (支持64个节点)
func NewShortFlakeGenerator(nodeID int64) *ShortFlakeGenerator {
	return &ShortFlakeGenerator{
		epoch:    1640995200000, // 2022-01-01 00:00:00
		nodeID:   nodeID & 0x3F, // 6位，最大63
		sequence: 0,
	}
}

// GenerateTraceID 生成跟踪ID
func (g *ShortFlakeGenerator) GenerateTraceID() string {
	return fmt.Sprintf("%d", g.generate())
}

// GenerateSpanID 生成跨度ID
func (g *ShortFlakeGenerator) GenerateSpanID() string {
	return fmt.Sprintf("%d", g.generate())
}

// GenerateRequestID 生成请求ID
func (g *ShortFlakeGenerator) GenerateRequestID() string {
	return fmt.Sprintf("%d", g.generate())
}

// GenerateCorrelationID 生成关联ID
func (g *ShortFlakeGenerator) GenerateCorrelationID() string {
	return fmt.Sprintf("%d", g.generate())
}

// Generate 生成 53 位整数 ID
func (g *ShortFlakeGenerator) Generate() int64 {
	return g.generate()
}

// generate 内部生成方法
func (g *ShortFlakeGenerator) generate() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < g.lastTime {
		now = g.lastTime
	}

	if now == g.lastTime {
		// 同一毫秒内，序列号递增
		g.sequence = (g.sequence + 1) & 0x3F // 6位，最大63
		if g.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for now <= g.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	// 时间戳(41位) + 节点ID(6位) + 序列号(6位)
	timestamp := (now - g.epoch) << 12
	node := g.nodeID << 6

	return timestamp | node | int64(g.sequence)
}

// ShortFlakeBase62Generator Base62 编码的短 ID 生成器
// 将 53 位数字编码为 9-10 字符的字符串
type ShortFlakeBase62Generator struct {
	*ShortFlakeGenerator
}

// NewShortFlakeBase62Generator 创建 Base62 编码生成器
func NewShortFlakeBase62Generator(nodeID int64) *ShortFlakeBase62Generator {
	return &ShortFlakeBase62Generator{
		ShortFlakeGenerator: NewShortFlakeGenerator(nodeID),
	}
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateTraceID 生成 Base62 编码的 TraceID
func (g *ShortFlakeBase62Generator) GenerateTraceID() string {
	return g.encodeBase62(g.generate())
}

// GenerateSpanID 生成 Base62 编码的 SpanID
func (g *ShortFlakeBase62Generator) GenerateSpanID() string {
	return g.encodeBase62(g.generate())
}

// GenerateRequestID 生成 Base62 编码的 RequestID
func (g *ShortFlakeBase62Generator) GenerateRequestID() string {
	return g.encodeBase62(g.generate())
}

// GenerateCorrelationID 生成 Base62 编码的 CorrelationID
func (g *ShortFlakeBase62Generator) GenerateCorrelationID() string {
	return g.encodeBase62(g.generate())
}

// encodeBase62 将整数编码为 Base62 字符串（9-10字符）
func (g *ShortFlakeBase62Generator) encodeBase62(num int64) string {
	if num == 0 {
		return "0"
	}

	var result [11]byte // 最多11字符
	pos := len(result) - 1

	for num > 0 {
		result[pos] = base62Chars[num%62]
		num /= 62
		pos--
	}

	return string(result[pos+1:])
}
