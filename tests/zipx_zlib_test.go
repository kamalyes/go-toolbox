/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 12:46:09
 * @FilePath: \go-toolbox\tests\zipx_zlib_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package tests

import (
	"bytes"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/zipx"
	"github.com/stretchr/testify/assert"
)

func TestAllZlibFunctions(t *testing.T) {
	t.Run("TestZlibCompressDecompress", TestZlibCompressDecompress)
	t.Run("TestZlibDecompressEmpty", TestZlibDecompressEmpty)
	t.Run("TestZlibDecompressInvalidData", TestZlibDecompressInvalidData)
	t.Run("TestMultiZlibCompressDecompress", TestMultiZlibCompressDecompress)
	t.Run("TestZlibCompressLargeData", TestZlibCompressLargeData)
	t.Run("TestZlibCompressEmpty", TestZlibCompressEmpty)
}

// TestZlibCompressDecompress 测试压缩和解压缩的功能
func TestZlibCompressDecompress(t *testing.T) {
	// 测试数据
	originalData := []byte("Hello, World! This is a test for zlib compression.")

	// 压缩数据
	compressedData, err := zipx.ZlibCompress(originalData)
	assert.NoError(t, err, "Compression error")
	assert.NotZero(t, len(compressedData), "Compressed data is empty")

	// 解压缩数据
	decompressedData, err := zipx.ZlibDecompress(compressedData)
	assert.NoError(t, err, "Decompression error")
	assert.True(t, bytes.Equal(originalData, decompressedData), "Decompressed data does not match original data")
}

// TestZlibDecompressEmpty 测试解压缩空数据的情况
func TestZlibDecompressEmpty(t *testing.T) {
	// 测试解压缩空数据
	_, err := zipx.ZlibDecompress([]byte{})
	assert.Error(t, err, "Expected error for empty compressed data")
}

// TestZlibDecompressInvalidData 测试解压缩无效数据的情况
func TestZlibDecompressInvalidData(t *testing.T) {
	// 测试解压缩无效数据
	_, err := zipx.ZlibDecompress([]byte{0x00, 0x01, 0x02})
	assert.Error(t, err, "Expected error for invalid compressed data")
}

// TestMultiZlibCompressDecompress 测试多次压缩和解压缩
func TestMultiZlibCompressDecompress(t *testing.T) {
	originalData := []byte("This is a test for multiple compressions.")

	// 多次压缩
	compressedData, err := zipx.MultiZlibCompress(originalData, 2)
	assert.NoError(t, err, "Compression error")

	// 多次解压缩
	decompressedData, err := zipx.MultiZlibDecompress(compressedData, 2)
	assert.NoError(t, err, "Decompression error")

	// 检查解压缩后的数据是否与原始数据匹配
	assert.True(t, bytes.Equal(originalData, decompressedData), "Decompressed data does not match original data")
	t.Log("Compression and decompression successful!")
}

// TestZlibCompressLargeData 测试大数据的压缩和解压缩
func TestZlibCompressLargeData(t *testing.T) {
	// 创建一个大数据
	originalData := bytes.Repeat([]byte("A"), 1<<20) // 1 MB of 'A'

	// 压缩数据
	compressedData, err := zipx.ZlibCompress(originalData)
	assert.NoError(t, err, "Compression error")
	assert.NotZero(t, len(compressedData), "Compressed data is empty")

	// 解压缩数据
	decompressedData, err := zipx.ZlibDecompress(compressedData)
	assert.NoError(t, err, "Decompression error")

	// 检查解压缩后的数据是否与原始数据相同
	assert.True(t, bytes.Equal(originalData, decompressedData), "Decompressed data does not match original data")
}

// TestZlibCompressEmpty 测试压缩空数据的情况
func TestZlibCompressEmpty(t *testing.T) {
	// 压缩空数据
	compressedData, err := zipx.ZlibCompress([]byte{})
	assert.NoError(t, err, "Compression error for empty data")
	assert.NotZero(t, len(compressedData), "Compressed data for empty input is empty")

	// 解压缩压缩后的空数据
	decompressedData, err := zipx.ZlibDecompress(compressedData)
	assert.NoError(t, err, "Decompression error for empty compressed data")

	// 检查解压缩后的数据是否为空
	assert.Empty(t, decompressedData, "Decompressed data for empty input should be empty")
}
