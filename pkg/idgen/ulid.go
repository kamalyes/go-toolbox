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

// GenerateTraceID 生成跟踪ID（完整 ULID，26字符）
// 格式: 10字符时间戳 + 16字符随机部分（Crockford Base32）
// 示例: "01ARZ3NDEKTSV4RRFFQ69G5FAV"
// 时间排序: 字典序 = 时间序，天然适合追踪
func (g *ULIDGenerator) GenerateTraceID() string {
	return g.generateULID()
}

// GenerateSpanID 生成跨度ID（ULID 后16字符，即随机部分）
// 格式: 截取 ULID 的随机部分
// 示例: "TSV4RRFFQ69G5FAV"
// 与 TraceID 的区别: 去掉时间戳前缀，更短，同一 Trace 内唯一
func (g *ULIDGenerator) GenerateSpanID() string {
	return g.generateULID()[10:]
}

// GenerateRequestID 生成请求ID（ULID前缀+计数器后缀）
// 格式: ULID前10字符(时间戳部分)-递增计数器
// 示例: "01ARZ3NDEK-1"
// 与 TraceID 的区别: 仅保留时间戳部分+计数器，可排序
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

// GenerateCorrelationID 生成关联ID（双 ULID 拼接，52字符）
// 格式: ULID-ULID
// 示例: "01ARZ3NDEKTSV4RRFFQ69G5FAV-01ARZ3ND0LWM4PQXX7HVK5FAVZ"
// 与 TraceID 的区别: 双段拼接，增强唯一性，适合跨系统关联
func (g *ULIDGenerator) GenerateCorrelationID() string {
	ulid1 := g.generateULID()
	ulid2 := g.generateULID()

	var sb strings.Builder
	sb.Grow(53)
	sb.WriteString(ulid1)
	sb.WriteByte('-')
	sb.WriteString(ulid2)

	return sb.String()
}

// generateULID 零分配版本（stack buffer + 快速 Base32 编码）
func (g *ULIDGenerator) generateULID() string {
	timestamp := uint64(time.Now().UnixMilli())

	var randomBytes [10]byte
	rand.Read(randomBytes[:])

	const encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	var result [26]byte

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
