/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:00:00
 * @FilePath: \go-toolbox\pkg\idgen\uuid.go
 * @Description: UUID v4 生成器（高性能版本）
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
)

// UUIDGenerator UUID v4 生成器（高性能版本）
type UUIDGenerator struct {
	counter uint64
}

// NewUUIDGenerator 创建 UUID 生成器
func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

// GenerateTraceID 生成跟踪ID（完整 UUID v4，36字符）
// 格式: 8-4-4-4-12 hex（含版本位和变体位）
// 示例: "550e8400-e29b-41d4-a716-446655440000"
func (g *UUIDGenerator) GenerateTraceID() string {
	return g.generateUUID()
}

// GenerateSpanID 生成跨度ID（UUID 前16字符，去除连字符）
// 格式: UUID 前8字节hex（16字符）
// 示例: "550e8400e29b41d4"
// 与 TraceID 的区别: 截取 UUID 时间戳部分，更短，同一 Trace 内唯一
func (g *UUIDGenerator) GenerateSpanID() string {
	uuid := g.generateUUID()
	return uuid[:16]
}

// GenerateRequestID 生成请求ID（UUID前缀+计数器后缀）
// 格式: UUID前8字符-递增计数器
// 示例: "550e8400-1"
// 与 TraceID 的区别: 带计数器后缀，可按请求顺序排序
func (g *UUIDGenerator) GenerateRequestID() string {
	counter := atomic.AddUint64(&g.counter, 1)
	uuid := g.generateUUID()

	var sb strings.Builder
	sb.Grow(20)
	sb.WriteString(uuid[:8])
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatUint(counter, 10))

	return sb.String()
}

// GenerateCorrelationID 生成关联ID（完整 UUID v4，36字符）
// 格式: 与 TraceID 相同的 UUID v4 格式
// 示例: "a1b2c3d4-e5f6-41d4-a716-446655440000"
// 与 TraceID 的区别: 独立生成，不与 TraceID 共享，用于跨系统关联
func (g *UUIDGenerator) GenerateCorrelationID() string {
	return g.generateUUID()
}

// generateUUID 零分配 UUID 生成（stack buffer）
func (g *UUIDGenerator) generateUUID() string {
	var buf [36]byte
	var randomBytes [16]byte
	rand.Read(randomBytes[:])

	randomBytes[6] = (randomBytes[6] & 0x0f) | 0x40
	randomBytes[8] = (randomBytes[8] & 0x3f) | 0x80

	const hexDigits = "0123456789abcdef"
	for i := 0; i < 4; i++ {
		buf[i*2] = hexDigits[randomBytes[i]>>4]
		buf[i*2+1] = hexDigits[randomBytes[i]&0xf]
	}
	buf[8] = '-'
	for i := 4; i < 6; i++ {
		buf[9+(i-4)*2] = hexDigits[randomBytes[i]>>4]
		buf[10+(i-4)*2] = hexDigits[randomBytes[i]&0xf]
	}
	buf[13] = '-'
	for i := 6; i < 8; i++ {
		buf[14+(i-6)*2] = hexDigits[randomBytes[i]>>4]
		buf[15+(i-6)*2] = hexDigits[randomBytes[i]&0xf]
	}
	buf[18] = '-'
	for i := 8; i < 10; i++ {
		buf[19+(i-8)*2] = hexDigits[randomBytes[i]>>4]
		buf[20+(i-8)*2] = hexDigits[randomBytes[i]&0xf]
	}
	buf[23] = '-'
	for i := 10; i < 16; i++ {
		buf[24+(i-10)*2] = hexDigits[randomBytes[i]>>4]
		buf[25+(i-10)*2] = hexDigits[randomBytes[i]&0xf]
	}

	return string(buf[:])
}

// FormatTraceID 将标准 UUID 格式化为无连字符的 TraceID（32字符）
// 用于兼容 OpenTelemetry 的 TraceID 格式
func FormatTraceID(uuid string) string {
	return strings.ReplaceAll(uuid, "-", "")
}

// FormatSpanID 从 UUID 中提取 SpanID（16字符）
func FormatSpanID(uuid string) string {
	clean := strings.ReplaceAll(uuid, "-", "")
	if len(clean) >= 16 {
		return clean[:16]
	}
	return fmt.Sprintf("%-16s", clean)
}
