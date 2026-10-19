/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-25 00:00:00
 * @FilePath: \go-toolbox\pkg\zipx\base.go
 * @Description: 压缩结果通用结构定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package zipx

import "fmt"

// CompressResult 压缩结果，包含压缩数据和统计信息
type CompressResult struct {
	Data           []byte  // 压缩后的数据
	OriginalSize   int     // 原始数据大小（字节）
	CompressedSize int     // 压缩后数据大小（字节）
	Ratio          float64 // 压缩率（压缩后/压缩前，越小压缩效果越好）
}

// String 返回压缩结果的字符串表示
func (r *CompressResult) String() string {
	return fmt.Sprintf("CompressResult{OriginalSize: %d bytes, CompressedSize: %d bytes, Ratio: %.2f%%}",
		r.OriginalSize, r.CompressedSize, r.Ratio*100)
}

// SavedBytes 返回节省的字节数
func (r *CompressResult) SavedBytes() int {
	return r.OriginalSize - r.CompressedSize
}

// SavedPercent 返回节省的百分比
func (r *CompressResult) SavedPercent() float64 {
	if r.OriginalSize == 0 {
		return 0
	}
	return (1 - r.Ratio) * 100
}

// newCompressResult 创建压缩结果
func newCompressResult(original, compressed []byte) *CompressResult {
	originalSize := len(original)
	compressedSize := len(compressed)

	var ratio float64
	if originalSize > 0 {
		ratio = float64(compressedSize) / float64(originalSize)
	}

	return &CompressResult{
		Data:           compressed,
		OriginalSize:   originalSize,
		CompressedSize: compressedSize,
		Ratio:          ratio,
	}
}

// ============================================================================
// 压缩级别
//
// 为什么要把级别暴露出来：小载荷的开销八成在 Reset，不在压缩
//
// 池化的 writer 每次复用前必须 Reset，而 flate 在级别 2 到 9 上的 Reset 要把
// hashHead 与 hashPrev 两张 32768 项的表清零，合计 256KB 的 memset实测
// （Apple M4）约 9.6 微秒，且这笔开销与载荷大小无关：
//
//	级别   Reset+写入+关闭(72字节)   仅 Reset
//	1      2141 ns                  10 ns
//	6      12069 ns                 9595 ns   ← 占 80%
//	9      12077 ns                 9629 ns
//
// 压 1.5MB 时它摊薄到 0.5%，压几十到几百字节时却占了八成级别 1 走的是
// deflateFast，表小得多，Reset 只要几纳秒
//
// 小载荷换到级别 1 的体积代价很小（实测）：
//
//	载荷              原始    级别1   级别6   级别1 多出
//	JSON 单对象       99      110     110     0.0%
//	日志行            91      111     111     0.0%
//	HTTP 头           247     235     232     1.3%
//	JSON 数组 1KB     726     215     214     0.5%
//
// 注意前两行：gzip 把数据压大了头尾固定开销约 18 字节，几十字节的载荷无论用
// 哪个级别都会变大——这种数据不该压
//
// 默认级别保持不变（DefaultCompression），因为输出字节是可观察的，不能静默改动
// 需要取舍的调用方用 WithLevel 系列显式指定
//
// [EN] Why the level is exposed: for small payloads four fifths of the cost is
// Reset, not compression. A pooled writer must be Reset before reuse, and
// flate's Reset at levels 2 through 9 zeroes the 32768-entry hashHead and
// hashPrev tables, a 256KB memset measured at about 9.6 microseconds — a cost
// independent of payload size. It amortises to 0.5% on 1.5MB but takes four
// fifths of the time on a few hundred bytes. Level 1 uses deflateFast with far
// smaller tables and resets in nanoseconds, and measured on realistic small
// payloads it costs between 0% and 1.3% more bytes. Note that gzip makes the
// smallest payloads larger: the header and footer are about 18 bytes, so tens of
// bytes grow at any level and should not be compressed at all. The default level
// is unchanged because output bytes are observable and must not shift silently;
// callers who want the trade-off ask for it with the WithLevel functions
// ============================================================================

// 压缩级别的取值范围，与 compress/flate 一致
const (
	// MinCompressLevel 最小级别：flate.HuffmanOnly（只做 Huffman 编码，不找匹配）
	MinCompressLevel = -2
	// MaxCompressLevel 最大级别：flate.BestCompression
	MaxCompressLevel = 9
	// levelPoolCount 池数组长度，下标为 level - MinCompressLevel
	levelPoolCount = MaxCompressLevel - MinCompressLevel + 1
)

// compressLevelIndex 把级别换算成池下标，并校验取值
//
// writer 的级别在构造时固定，Reset 不能改级别，所以每个级别各备一个池
func compressLevelIndex(level int) (int, error) {
	if level < MinCompressLevel || level > MaxCompressLevel {
		return 0, fmt.Errorf("zipx: 压缩级别 %d 超出范围 [%d, %d]", level, MinCompressLevel, MaxCompressLevel)
	}
	return level - MinCompressLevel, nil
}
