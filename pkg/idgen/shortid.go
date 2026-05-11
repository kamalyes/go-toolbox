/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-05-15 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-15 18:10:15
 * @FilePath: \go-toolbox\pkg\idgen\shortid.go
 * @Description: 8~10位短ID生成器（Base62 编码，无锁设计）
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

// ShortIDGenerator 8~10位短ID生成器
// 使用 Base62 编码，时间戳+序列号组合，无互斥锁设计
// TraceID: 10字符 (6字符ms时间戳 + 4字符序列号)
// SpanID: 8字符 (纯随机)
// RequestID: ~10字符 (5字符时间戳-计数器)
// CorrelationID: 10字符 (纯随机)
type ShortIDGenerator struct {
	epoch    int64
	traceSeq uint64
	counter  uint64
}

// NewShortIDGenerator 创建短ID生成器
func NewShortIDGenerator() *ShortIDGenerator {
	return &ShortIDGenerator{
		epoch: 1640995200000,
	}
}

// GenerateTraceID 生成跟踪ID（10字符 Base62，时间可排序）
// 格式: 6字符Base62(毫秒时间戳) + 4字符Base62(原子序列号)
// 示例: "0If09Q4b2x"
// 特点: 前6字符编码时间戳，字典序=时间序，后4字符序列号保证同毫秒唯一(1476万/ms)
func (g *ShortIDGenerator) GenerateTraceID() string {
	ts := uint64(time.Now().UnixMilli() - g.epoch)
	seq := atomic.AddUint64(&g.traceSeq, 1)
	return encodeBase62Pad(ts, 6) + encodeBase62Pad(seq, 4)
}

// GenerateSpanID 生成跨度ID（8字符 Base62，纯随机）
// 格式: 8字符随机Base62
// 示例: "K9mPxR2v"
// 与 TraceID 的区别: 无时间戳前缀，纯随机，更短，同 Trace 内唯一
func (g *ShortIDGenerator) GenerateSpanID() string {
	return randomBase62(8)
}

// GenerateRequestID 生成请求ID（时间戳前缀+计数器，可排序）
// 格式: 5字符Base62(秒级时间戳)-递增计数器
// 示例: "0If09-1"
// 与 TraceID 的区别: 秒级精度+计数器，方便按请求顺序排序
func (g *ShortIDGenerator) GenerateRequestID() string {
	counter := atomic.AddUint64(&g.counter, 1)
	ts := uint64(time.Now().Unix() - g.epoch/1000)

	var sb strings.Builder
	sb.Grow(12)
	sb.WriteString(encodeBase62Pad(ts, 5))
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatUint(counter, 10))

	return sb.String()
}

// GenerateCorrelationID 生成关联ID（10字符 Base62，纯随机）
// 格式: 10字符随机Base62
// 示例: "xY7wN4qL2v"
// 与 TraceID 的区别: 无时间戳前缀，纯随机，适合跨系统关联
func (g *ShortIDGenerator) GenerateCorrelationID() string {
	return randomBase62(10)
}

// encodeBase62Pad 将 uint64 编码为固定宽度的 Base62 字符串
// 不足宽度时左侧填充 '0'
func encodeBase62Pad(num uint64, width int) string {
	var buf [12]byte
	for i := width - 1; i >= 0; i-- {
		buf[i] = base62Chars[num%62]
		num /= 62
	}
	return string(buf[:width])
}

// randomBase62 生成指定长度的随机 Base62 字符串
func randomBase62(length int) string {
	var randomBytes [12]byte
	rand.Read(randomBytes[:length])

	var result [12]byte
	for i := 0; i < length; i++ {
		result[i] = base62Chars[int(randomBytes[i])%62]
	}
	return string(result[:length])
}
