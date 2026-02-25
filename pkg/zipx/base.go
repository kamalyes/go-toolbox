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
