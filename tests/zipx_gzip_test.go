/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 11:55:55
 * @FilePath: \go-toolbox\tests\zipx_gzip_test.go
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

func TestAllGzipFunctions(t *testing.T) {
	t.Run("TestGzipCompressDecompress", TestGzipCompressDecompress)
	t.Run("TestGzipDecompressEmpty", TestGzipDecompressEmpty)
	t.Run("TestGzipDecompressInvalidData", TestGzipDecompressInvalidData)
	t.Run("TestMultiGzipCompressDecompress", TestMultiGzipCompressDecompress)
	t.Run("TestGzipCompressLargeData", TestGzipCompressLargeData)
	t.Run("TestGzipCompressEmpty", TestGzipCompressEmpty)
}

// TestGzipCompressDecompress 测试压缩和解压缩的功能
func TestGzipCompressDecompress(t *testing.T) {
	// 测试数据
	originalData := []byte("Hello, World! This is a test for gzip compression.")

	// 压缩数据
	compressedData, err := zipx.GzipCompress(originalData)
	assert.NoError(t, err, "Compression error")
	assert.NotZero(t, len(compressedData), "Compressed data is empty")

	// 解压缩数据
	decompressedData, err := zipx.GzipDecompress(compressedData)
	assert.NoError(t, err, "Decompression error")
	assert.True(t, bytes.Equal(originalData, decompressedData), "Decompressed data does not match original data")
}

// TestGzipDecompressEmpty 测试解压缩空数据的情况
func TestGzipDecompressEmpty(t *testing.T) {
	// 测试解压缩空数据
	_, err := zipx.GzipDecompress([]byte{})
	assert.Error(t, err, "Expected error for empty compressed data")
}

// TestGzipDecompressInvalidData 测试解压缩无效数据的情况
func TestGzipDecompressInvalidData(t *testing.T) {
	// 测试解压缩无效数据
	_, err := zipx.GzipDecompress([]byte{0x00, 0x01, 0x02})
	assert.Error(t, err, "Expected error for invalid compressed data")
}

// TestMultiGzipCompressDecompress 测试多次压缩和解压缩
func TestMultiGzipCompressDecompress(t *testing.T) {
	originalData := []byte("This is a test for multiple compressions.")

	// 多次压缩
	compressedData, err := zipx.MultiGZipCompress(originalData, 2)
	assert.NoError(t, err, "Compression error")

	// 多次解压缩
	decompressedData, err := zipx.MultiGZipDecompress(compressedData, 2)
	assert.NoError(t, err, "Decompression error")
	assert.True(t, bytes.Equal(originalData, decompressedData), "Decompressed data does not match original data")

	t.Log("Compression and decompression successful!")
}

// TestGzipCompressLargeData 测试大数据的压缩和解压缩
func TestGzipCompressLargeData(t *testing.T) {
	// 创建一个大数据
	originalData := bytes.Repeat([]byte("A"), 1<<20) // 1 MB of 'A'

	// 压缩数据
	compressedData, err := zipx.GzipCompress(originalData)
	assert.NoError(t, err, "Compression error")
	assert.NotZero(t, len(compressedData), "Compressed data is empty")

	t.Logf("Compressed data length: %d", len(compressedData))

	// 解压缩数据
	decompressedData, err := zipx.GzipDecompress(compressedData)
	assert.NoError(t, err, "Decompression error")
	assert.True(t, bytes.Equal(originalData, decompressedData), "Decompressed data does not match original data")
}

// TestGzipCompressEmpty 测试压缩空数据的情况
func TestGzipCompressEmpty(t *testing.T) {
	// 压缩空数据
	compressedData, err := zipx.GzipCompress([]byte{})
	assert.NoError(t, err, "Compression error for empty data")
	assert.NotZero(t, len(compressedData), "Compressed data for empty input is empty")

	// 解压缩压缩后的空数据
	decompressedData, err := zipx.GzipDecompress(compressedData)
	assert.NoError(t, err, "Decompression error for empty compressed data")
	assert.Empty(t, decompressedData, "Decompressed data for empty input should be empty")
}
