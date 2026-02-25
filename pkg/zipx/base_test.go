/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-25 00:00:00
 * @FilePath: \go-toolbox\pkg\zipx\base_test.go
 * @Description: 压缩结果结构体测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package zipx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCompressResult 测试 CompressResult 结构体基本功能
func TestCompressResult(t *testing.T) {
	original := []byte(strings.Repeat("test data ", 100))
	compressed := []byte("compressed")

	result := newCompressResult(original, compressed)

	assert.Equal(t, len(original), result.OriginalSize)
	assert.Equal(t, len(compressed), result.CompressedSize)

	expectedRatio := float64(len(compressed)) / float64(len(original))
	assert.Equal(t, expectedRatio, result.Ratio)
}

// TestCompressResult_SavedBytes 测试 SavedBytes 方法
func TestCompressResult_SavedBytes(t *testing.T) {
	original := []byte(strings.Repeat("test ", 200))
	compressed := []byte("small")

	result := newCompressResult(original, compressed)
	saved := result.SavedBytes()

	expected := len(original) - len(compressed)
	assert.Equal(t, expected, saved)
}

// TestCompressResult_SavedPercent 测试 SavedPercent 方法
func TestCompressResult_SavedPercent(t *testing.T) {
	original := []byte(strings.Repeat("a", 1000))
	compressed := []byte(strings.Repeat("b", 100))

	result := newCompressResult(original, compressed)
	percent := result.SavedPercent()

	assert.Equal(t, 90.0, percent)
}

// TestCompressResult_String 测试 String 方法
func TestCompressResult_String(t *testing.T) {
	original := []byte("original data")
	compressed := []byte("comp")

	result := newCompressResult(original, compressed)
	str := result.String()

	assert.Contains(t, str, "OriginalSize")
	assert.Contains(t, str, "CompressedSize")
	assert.Contains(t, str, "Ratio")
}

// TestCompressResult_EmptyData 测试空数据场景
func TestCompressResult_EmptyData(t *testing.T) {
	t.Run("原始数据为空", func(t *testing.T) {
		original := []byte{}
		compressed := []byte("compressed")

		result := newCompressResult(original, compressed)

		assert.Equal(t, 0, result.OriginalSize)
		assert.Equal(t, len(compressed), result.CompressedSize)
		assert.Equal(t, 0.0, result.Ratio)
		assert.Equal(t, -len(compressed), result.SavedBytes())
	})

	t.Run("压缩数据为空", func(t *testing.T) {
		original := []byte("original data")
		compressed := []byte{}

		result := newCompressResult(original, compressed)

		assert.Equal(t, len(original), result.OriginalSize)
		assert.Equal(t, 0, result.CompressedSize)
		assert.Equal(t, 0.0, result.Ratio)
		assert.Equal(t, len(original), result.SavedBytes())
		assert.Equal(t, 100.0, result.SavedPercent())
	})

	t.Run("两者都为空", func(t *testing.T) {
		original := []byte{}
		compressed := []byte{}

		result := newCompressResult(original, compressed)

		assert.Equal(t, 0, result.OriginalSize)
		assert.Equal(t, 0, result.CompressedSize)
		assert.Equal(t, 0.0, result.Ratio)
		assert.Equal(t, 0, result.SavedBytes())
		assert.Equal(t, 0.0, result.SavedPercent())
	})
}

// TestCompressResult_NoCompression 测试无压缩效果场景
func TestCompressResult_NoCompression(t *testing.T) {
	t.Run("压缩后大小相同", func(t *testing.T) {
		data := []byte("same size")
		result := newCompressResult(data, data)

		assert.Equal(t, 1.0, result.Ratio)
		assert.Equal(t, 0, result.SavedBytes())
		assert.Equal(t, 0.0, result.SavedPercent())
	})

	t.Run("压缩后反而变大", func(t *testing.T) {
		original := []byte("small")
		compressed := []byte("much larger compressed data")

		result := newCompressResult(original, compressed)

		assert.Greater(t, result.Ratio, 1.0)
		assert.Less(t, result.SavedBytes(), 0)
		assert.Less(t, result.SavedPercent(), 0.0)
	})
}

// TestCompressResult_HighCompression 测试高压缩率场景
func TestCompressResult_HighCompression(t *testing.T) {
	original := []byte(strings.Repeat("A", 10000))
	compressed := []byte("tiny")

	result := newCompressResult(original, compressed)

	assert.Less(t, result.Ratio, 0.01) // 压缩率小于 1%
	assert.Greater(t, result.SavedPercent(), 99.0)
	assert.Greater(t, result.SavedBytes(), 9990)
}

// TestCompressResult_LargeData 测试大数据场景
func TestCompressResult_LargeData(t *testing.T) {
	// 模拟 10MB 数据压缩到 1MB
	originalSize := 10 * 1024 * 1024
	compressedSize := 1 * 1024 * 1024

	original := make([]byte, originalSize)
	compressed := make([]byte, compressedSize)

	result := newCompressResult(original, compressed)

	assert.Equal(t, originalSize, result.OriginalSize)
	assert.Equal(t, compressedSize, result.CompressedSize)
	assert.InDelta(t, 0.1, result.Ratio, 0.001)
	assert.Equal(t, 90.0, result.SavedPercent())
}

// TestCompressResult_EdgeCases 测试边界情况
func TestCompressResult_EdgeCases(t *testing.T) {
	t.Run("单字节数据", func(t *testing.T) {
		original := []byte("A")
		compressed := []byte("B")

		result := newCompressResult(original, compressed)

		assert.Equal(t, 1, result.OriginalSize)
		assert.Equal(t, 1, result.CompressedSize)
		assert.Equal(t, 1.0, result.Ratio)
	})

	t.Run("极小压缩改善", func(t *testing.T) {
		original := []byte(strings.Repeat("x", 1000))
		compressed := []byte(strings.Repeat("y", 999))

		result := newCompressResult(original, compressed)

		assert.InDelta(t, 0.999, result.Ratio, 0.001)
		assert.Equal(t, 1, result.SavedBytes())
		assert.InDelta(t, 0.1, result.SavedPercent(), 0.01)
	})
}

// TestCompressResult_StringFormat 测试字符串格式化
func TestCompressResult_StringFormat(t *testing.T) {
	original := []byte(strings.Repeat("test", 250))
	compressed := []byte(strings.Repeat("x", 100))

	result := newCompressResult(original, compressed)
	str := result.String()

	// 验证包含所有关键信息
	assert.Contains(t, str, "1000")  // OriginalSize
	assert.Contains(t, str, "100")   // CompressedSize
	assert.Contains(t, str, "10.00") // Ratio percentage
	assert.Contains(t, str, "bytes") // 单位
	assert.Contains(t, str, "Ratio") // 字段名
}
