/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-21 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 21:00:00
 * @FilePath: \go-toolbox\pkg\idgen\ulid.go
 * @Description: ULID 生成器（时间排序友好，高性能版本）
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package idgen

import (
	"crypto/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// ULIDGenerator ULID 生成器（时间排序友好，高性能版本）
type ULIDGenerator struct {
	counter uint64
}

// NewULIDGenerator 创建 ULID 生成器
func NewULIDGenerator() *ULIDGenerator {
	return &ULIDGenerator{}
}

// GenerateTraceID 生成跟踪ID
func (g *ULIDGenerator) GenerateTraceID() string {
	return g.generateULID()
}

// GenerateSpanID 生成跨度ID
func (g *ULIDGenerator) GenerateSpanID() string {
	return g.generateULID()[:16]
}

// GenerateRequestID 生成请求ID
func (g *ULIDGenerator) GenerateRequestID() string {
	counter := atomic.AddUint64(&g.counter, 1)
	ulid := g.generateULID()

	var sb strings.Builder
	sb.Grow(22)
	sb.WriteString(ulid[:10])
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatUint(counter, 10))

	return sb.String()
}

// GenerateCorrelationID 生成关联ID
func (g *ULIDGenerator) GenerateCorrelationID() string {
	return g.generateULID()
}

// generateULID 零分配版本（stack buffer + 快速 Base32 编码）
func (g *ULIDGenerator) generateULID() string {
	timestamp := uint64(time.Now().UnixMilli())

	var randomBytes [10]byte
	rand.Read(randomBytes[:])

	const encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	var result [26]byte

	// 时间戳编码（10字符）- 快速位运算
	result[0] = encoding[(timestamp>>45)&0x1F]
	result[1] = encoding[(timestamp>>40)&0x1F]
	result[2] = encoding[(timestamp>>35)&0x1F]
	result[3] = encoding[(timestamp>>30)&0x1F]
	result[4] = encoding[(timestamp>>25)&0x1F]
	result[5] = encoding[(timestamp>>20)&0x1F]
	result[6] = encoding[(timestamp>>15)&0x1F]
	result[7] = encoding[(timestamp>>10)&0x1F]
	result[8] = encoding[(timestamp>>5)&0x1F]
	result[9] = encoding[timestamp&0x1F]

	// 随机部分编码（16字符）
	result[10] = encoding[(randomBytes[0]>>3)&0x1F]
	result[11] = encoding[((randomBytes[0]<<2)|(randomBytes[1]>>6))&0x1F]
	result[12] = encoding[(randomBytes[1]>>1)&0x1F]
	result[13] = encoding[((randomBytes[1]<<4)|(randomBytes[2]>>4))&0x1F]
	result[14] = encoding[((randomBytes[2]<<1)|(randomBytes[3]>>7))&0x1F]
	result[15] = encoding[(randomBytes[3]>>2)&0x1F]
	result[16] = encoding[((randomBytes[3]<<3)|(randomBytes[4]>>5))&0x1F]
	result[17] = encoding[randomBytes[4]&0x1F]
	result[18] = encoding[(randomBytes[5]>>3)&0x1F]
	result[19] = encoding[((randomBytes[5]<<2)|(randomBytes[6]>>6))&0x1F]
	result[20] = encoding[(randomBytes[6]>>1)&0x1F]
	result[21] = encoding[((randomBytes[6]<<4)|(randomBytes[7]>>4))&0x1F]
	result[22] = encoding[((randomBytes[7]<<1)|(randomBytes[8]>>7))&0x1F]
	result[23] = encoding[(randomBytes[8]>>2)&0x1F]
	result[24] = encoding[((randomBytes[8]<<3)|(randomBytes[9]>>5))&0x1F]
	result[25] = encoding[randomBytes[9]&0x1F]

	return string(result[:])
}
