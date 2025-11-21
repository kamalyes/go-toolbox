/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:00:00
 * @FilePath: \go-toolbox\pkg\idgen\default.go
 * @Description: 默认 Hex ID 生成器（高性能版本）
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"crypto/rand"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// DefaultIDGenerator 默认ID生成器（高性能 Hex 编码）
type DefaultIDGenerator struct {
	counter uint64
	mu      sync.Mutex
}

// NewDefaultIDGenerator 创建默认ID生成器
func NewDefaultIDGenerator() *DefaultIDGenerator {
	return &DefaultIDGenerator{}
}

// GenerateTraceID 生成跟踪ID（零分配优化）
func (g *DefaultIDGenerator) GenerateTraceID() string {
	// 使用 stack buffer 避免堆分配
	var buf [32]byte
	timestamp := uint64(time.Now().UnixNano())

	// 快速 hex 编码时间戳
	const hexDigits = "0123456789abcdef"
	for i := 15; i >= 0; i-- {
		buf[i] = hexDigits[timestamp&0xf]
		timestamp >>= 4
	}

	// 随机部分
	var randomBytes [8]byte
	rand.Read(randomBytes[:])
	for i := 0; i < 16; i++ {
		v := randomBytes[i/2]
		if i%2 == 0 {
			v >>= 4
		}
		buf[16+i] = hexDigits[v&0xf]
	}

	return string(buf[:])
}

// GenerateSpanID 生成跨度ID（零分配优化）
func (g *DefaultIDGenerator) GenerateSpanID() string {
	var buf [16]byte
	var randomBytes [8]byte
	rand.Read(randomBytes[:])

	const hexDigits = "0123456789abcdef"
	for i := 0; i < 16; i++ {
		v := randomBytes[i/2]
		if i%2 == 0 {
			v >>= 4
		}
		buf[i] = hexDigits[v&0xf]
	}

	return string(buf[:])
}

// GenerateRequestID 生成请求ID（使用 strings.Builder 优化）
func (g *DefaultIDGenerator) GenerateRequestID() string {
	counter := atomic.AddUint64(&g.counter, 1)
	timestamp := time.Now().Unix()

	// 预分配容量避免扩容
	var sb strings.Builder
	sb.Grow(32)
	sb.WriteString(strconv.FormatInt(timestamp, 10))
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatUint(counter, 10))

	return sb.String()
}

// GenerateCorrelationID 生成关联ID（零分配优化）
func (g *DefaultIDGenerator) GenerateCorrelationID() string {
	var buf [36]byte
	var randomBytes [16]byte
	rand.Read(randomBytes[:])

	// 设置版本和变体位
	randomBytes[6] = (randomBytes[6] & 0x0f) | 0x40
	randomBytes[8] = (randomBytes[8] & 0x3f) | 0x80

	// 快速 hex 编码（UUID 格式）
	const hexDigits = "0123456789abcdef"
	encodeHex := func(dst []byte, src []byte) {
		for i := 0; i < len(src); i++ {
			dst[i*2] = hexDigits[src[i]>>4]
			dst[i*2+1] = hexDigits[src[i]&0xf]
		}
	}

	encodeHex(buf[0:8], randomBytes[0:4])
	buf[8] = '-'
	encodeHex(buf[9:13], randomBytes[4:6])
	buf[13] = '-'
	encodeHex(buf[14:18], randomBytes[6:8])
	buf[18] = '-'
	encodeHex(buf[19:23], randomBytes[8:10])
	buf[23] = '-'
	encodeHex(buf[24:36], randomBytes[10:16])

	return string(buf[:])
}
