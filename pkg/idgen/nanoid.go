/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:00:00
 * @FilePath: \go-toolbox\pkg\idgen\nanoid.go
 * @Description: NanoID 生成器（高性能版本）
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"crypto/rand"
	"strconv"
	"strings"
	"sync/atomic"
)

// NanoIDGenerator NanoID 生成器（高性能版本，预分配字母表）
type NanoIDGenerator struct {
	counter  uint64
	alphabet []byte // 使用 []byte 提升性能
	size     int
}

// NewNanoIDGenerator 创建 NanoID 生成器
func NewNanoIDGenerator() *NanoIDGenerator {
	return &NanoIDGenerator{
		alphabet: []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_-"),
		size:     21,
	}
}

// GenerateTraceID 生成跟踪ID（完整 NanoID，21字符）
// 格式: URL 安全的 21 字符随机字符串
// 示例: "V1StGXR8_Z5jdHi6B-myT"
func (g *NanoIDGenerator) GenerateTraceID() string {
	return g.generateNanoID()
}

// GenerateSpanID 生成跨度ID（NanoID 前16字符）
// 格式: 截取 NanoID 前16字符
// 示例: "V1StGXR8_Z5jdHi6"
// 与 TraceID 的区别: 更短，同一 Trace 内唯一即可
func (g *NanoIDGenerator) GenerateSpanID() string {
	return g.generateNanoID()[:16]
}

// GenerateRequestID 生成请求ID（NanoID前缀+计数器后缀）
// 格式: NanoID前10字符-递增计数器
// 示例: "V1StGXR8_Z-1"
// 与 TraceID 的区别: 带计数器后缀，可按请求顺序排序
func (g *NanoIDGenerator) GenerateRequestID() string {
	counter := atomic.AddUint64(&g.counter, 1)
	nanoID := g.generateNanoID()

	var sb strings.Builder
	sb.Grow(22)
	sb.WriteString(nanoID[:10])
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatUint(counter, 10))

	return sb.String()
}

// GenerateCorrelationID 生成关联ID（完整 NanoID，21字符）
// 格式: 与 TraceID 相同长度的 NanoID
// 示例: "K9mPxR2vN4qL7wBjT5cYs"
// 与 TraceID 的区别: 独立生成，用于跨系统关联
func (g *NanoIDGenerator) GenerateCorrelationID() string {
	return g.generateNanoID()
}

// generateNanoID 零分配版本（使用 stack buffer）
func (g *NanoIDGenerator) generateNanoID() string {
	var randomBytes [21]byte
	rand.Read(randomBytes[:])

	var id [21]byte
	alphabetLen := len(g.alphabet)
	for i := 0; i < g.size; i++ {
		id[i] = g.alphabet[int(randomBytes[i])%alphabetLen]
	}

	return string(id[:])
}
