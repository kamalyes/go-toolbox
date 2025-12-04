/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-04 18:15:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 18:29:53
 * @FilePath: \go-toolbox\pkg\zipx\gzip_test.go
 * @Description: Gzip 压缩解压缩测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package zipx

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGzipCompress 测试基本 Gzip 压缩功能
func TestGzipCompress(t *testing.T) {
	testData := []byte("Hello, World! This is a test string for gzip compression.")

	compressed, err := GzipCompress(testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressed, "Gzip压缩后数据为空")

	// 验证压缩后数据长度应该小于原始数据（对于重复内容）
	largeData := bytes.Repeat([]byte("B"), 1000)
	compressedLarge, err := GzipCompress(largeData)
	assert.NoError(t, err)

	assert.Less(t, len(compressedLarge), len(largeData), "Gzip压缩效果不佳: 原始=%d, 压缩后=%d", len(largeData), len(compressedLarge))

	t.Logf("Gzip压缩测试成功: 原始=%d字节, 压缩后=%d字节, 压缩率=%.2f%%",
		len(largeData), len(compressedLarge), float64(len(compressedLarge))/float64(len(largeData))*100)
}

// TestGzipDecompress 测试基本 Gzip 解压缩功能
func TestGzipDecompress(t *testing.T) {
	testData := []byte("Hello, World! This is a test string for gzip compression and decompression.")

	// 先压缩
	compressed, err := GzipCompress(testData)
	assert.NoError(t, err)

	// 再解压缩
	decompressed, err := GzipDecompress(compressed)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, bytes.Equal(testData, decompressed), "Gzip解压缩后数据不一致:\n原始: %s\n解压后: %s", testData, decompressed)

	t.Logf("Gzip解压缩测试成功: 原始=%d字节, 解压后=%d字节", len(testData), len(decompressed))
}

// TestGzipCompressObject_SingleMessage 测试单个消息对象 Gzip 压缩
func TestGzipCompressObjectSingleMessage(t *testing.T) {
	// 创建测试消息（截断时间以移除单调时钟）
	testMsg := TestMessage{
		ID:        "gzip_msg_123456",
		Content:   "这是一条用于Gzip压缩测试的消息，包含中文内容！",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "gzip_user_789",
		Type:      1,
	}

	// 压缩对象
	compressed, err := GzipCompressObject(testMsg)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressed, "Gzip压缩后数据为空")

	// 解压缩对象
	decompressed, err := GzipDecompressObject[TestMessage](compressed)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, reflect.DeepEqual(testMsg, decompressed), "Gzip对象压缩解压缩后数据不一致:\n原始: %+v\n解压后: %+v", testMsg, decompressed)

	t.Logf("Gzip对象压缩测试成功: 压缩后大小=%d字节", len(compressed))
}

// TestGzipCompressObjectWithSize 测试带原始大小返回的泛型压缩函数
func TestGzipCompressObjectWithSize(t *testing.T) {
	// 创建测试消息（截断时间以移除单调时钟）
	testMsg := TestMessage{
		ID:        "msg_123456",
		Content:   "这是一条测试消息，包含中文内容！",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "user_789",
		Type:      1,
	}

	// 使用新函数进行压缩，同时获取原始大小
	compressed, originalSize, err := GzipCompressObjectWithSize(testMsg)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressed, "压缩后数据为空")
	assert.Greater(t, originalSize, 0, "原始数据大小应该大于0")

	// 验证原始大小的正确性
	jsonData, _ := json.Marshal(testMsg)
	assert.Equal(t, len(jsonData), originalSize, "返回的原始大小应该与手动序列化的大小一致")

	// 解压缩验证
	decompressed, err := GzipDecompressObject[TestMessage](compressed)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(testMsg, decompressed), "对象压缩解压缩后数据不一致")

	t.Logf("带大小的Gzip压缩测试成功: 原始大小=%d字节, 压缩后大小=%d字节, 压缩率=%.2f%%",
		originalSize, len(compressed), float64(len(compressed))/float64(originalSize)*100)
}

// TestGzipCompressObject_MessageSlice 测试消息数组 Gzip 压缩
func TestGzipCompressObjectMessageSlice(t *testing.T) {
	// 创建多条测试消息（截断时间以移除单调时钟）
	baseTime := time.Now().Truncate(time.Second)
	messages := []TestMessage{
		{
			ID:        "gzip_msg_001",
			Content:   "Gzip第一条消息",
			Timestamp: baseTime,
			UserID:    "gzip_user_001",
			Type:      1,
		},
		{
			ID:        "gzip_msg_002",
			Content:   "Gzip第二条消息，内容更长一些，用于测试Gzip压缩效果",
			Timestamp: baseTime.Add(1 * time.Minute),
			UserID:    "gzip_user_002",
			Type:      2,
		},
		{
			ID:        "gzip_msg_003",
			Content:   "Gzip第三条消息",
			Timestamp: baseTime.Add(2 * time.Minute),
			UserID:    "gzip_user_003",
			Type:      1,
		},
	}

	// 压缩消息数组
	compressed, err := GzipCompressObject(messages)
	assert.NoError(t, err)

	// 解压缩消息数组
	decompressed, err := GzipDecompressObject[[]TestMessage](compressed)
	assert.NoError(t, err)

	// 验证数组长度
	assert.Equal(t, len(messages), len(decompressed), "Gzip消息数组长度不一致: 原始=%d, 解压后=%d", len(messages), len(decompressed))

	// 验证每条消息
	for i, original := range messages {
		assert.True(t, reflect.DeepEqual(original, decompressed[i]), "Gzip第%d条消息不一致:\n原始: %+v\n解压后: %+v", i, original, decompressed[i])
	}

	t.Logf("Gzip消息数组压缩测试成功: %d条消息, 压缩后大小=%d字节", len(messages), len(compressed))
}

// TestMultiGZipCompress 测试多重 Gzip 压缩
func TestMultiGZipCompress(t *testing.T) {
	testData := []byte("Hello, World! This is a test for multiple gzip compression rounds.")

	// 测试3轮压缩
	compressed, err := MultiGZipCompress(testData, 3)
	assert.NoError(t, err)

	// 测试3轮解压缩
	decompressed, err := MultiGZipDecompress(compressed, 3)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, bytes.Equal(testData, decompressed), "Gzip多重压缩解压后数据不一致")

	t.Logf("Gzip多重压缩测试成功: 原始=%d字节, 3轮压缩后=%d字节", len(testData), len(compressed))
}

// TestGzipConcurrency 测试 Gzip 并发安全性
func TestGzipConcurrency(t *testing.T) {
	const numGoroutines = 50
	const numOperations = 30

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperations)

	// 创建测试数据（截断时间以移除单调时钟）
	baseTime := time.Now().Truncate(time.Second)
	messages := make([]TestMessage, 5)
	for i := 0; i < 5; i++ {
		messages[i] = TestMessage{
			ID:        "gzip_concurrent_msg_" + string(rune('A'+i)),
			Content:   "Gzip并发测试消息内容",
			Timestamp: baseTime.Add(time.Duration(i) * time.Second),
			UserID:    "gzip_concurrent_user",
			Type:      rand.Intn(3) + 1,
		}
	}

	// 启动多个协程进行并发压缩解压缩
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				// 压缩
				compressed, err := GzipCompressObject(messages)
				if err != nil {
					errors <- err
					return
				}

				// 解压缩
				decompressed, err := GzipDecompressObject[[]TestMessage](compressed)
				if err != nil {
					errors <- err
					return
				}

				// 验证
				assert.Equal(t, len(messages), len(decompressed), "Gzip并发测试数据长度不一致")
			}
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()
	close(errors)

	// 检查错误
	errorCount := 0
	for err := range errors {
		assert.Error(t, err, "Gzip并发测试出错: %v", err)
		errorCount++
	}

	if errorCount == 0 {
		t.Logf("Gzip并发测试成功: %d个协程, 每个执行%d次操作", numGoroutines, numOperations)
	}
}

// TestGzipVsZlibComparison 对比 Gzip 和 Zlib 压缩效果
func TestGzipVsZlibComparison(t *testing.T) {
	// 创建大量重复数据进行压缩对比
	largeData := bytes.Repeat([]byte("重复数据用于测试压缩效果对比"), 200)

	// Gzip 压缩
	gzipStart := time.Now()
	gzipCompressed, err := GzipCompress(largeData)
	assert.NoError(t, err)
	gzipTime := time.Since(gzipStart)

	// Zlib 压缩
	zlibStart := time.Now()
	zlibCompressed, err := ZlibCompress(largeData)
	assert.NoError(t, err)
	zlibTime := time.Since(zlibStart)

	// 计算压缩率
	gzipRatio := float64(len(gzipCompressed)) / float64(len(largeData)) * 100
	zlibRatio := float64(len(zlibCompressed)) / float64(len(largeData)) * 100

	t.Logf("压缩对比测试结果:")
	t.Logf("原始数据: %d字节", len(largeData))
	t.Logf("Gzip: 压缩后=%d字节, 压缩率=%.2f%%, 耗时=%v", len(gzipCompressed), gzipRatio, gzipTime)
	t.Logf("Zlib: 压缩后=%d字节, 压缩率=%.2f%%, 耗时=%v", len(zlibCompressed), zlibRatio, zlibTime)

	// 验证解压缩
	gzipDecompressed, err := GzipDecompress(gzipCompressed)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(largeData, gzipDecompressed), "Gzip解压缩数据不一致")

	zlibDecompressed, err := ZlibDecompress(zlibCompressed)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(largeData, zlibDecompressed), "Zlib解压缩数据不一致")
}

// TestGzipLargeObjectArray 测试大型对象数组 Gzip 压缩
func TestGzipLargeObjectArray(t *testing.T) {
	// 创建大量复杂对象（截断时间以移除单调时钟）
	const objectCount = 500
	baseTime := time.Now().Truncate(time.Second)
	complexObjects := make([]TestComplexObject, objectCount)

	for i := 0; i < objectCount; i++ {
		complexObjects[i] = TestComplexObject{
			Name: "Gzip测试用户" + string(rune(i)),
			Age:  20 + i%50,
			Tags: []string{"gzip", "test", "large", "array"},
			Metadata: map[string]interface{}{
				"index":   i,
				"type":    "gzip_test",
				"enabled": i%2 == 0,
			},
			Messages: []TestMessage{
				{
					ID:        "gzip_large_msg_" + string(rune(i)),
					Content:   "Gzip大型数组测试消息",
					Timestamp: baseTime.Add(time.Duration(i) * time.Second),
					UserID:    "gzip_large_user",
					Type:      i%4 + 1,
				},
			},
		}
	}

	start := time.Now()

	// 压缩大型对象数组
	compressed, err := GzipCompressObject(complexObjects)
	assert.NoError(t, err)

	compressTime := time.Since(start)

	start = time.Now()

	// 解压缩大型对象数组
	decompressed, err := GzipDecompressObject[[]TestComplexObject](compressed)
	assert.NoError(t, err)

	decompressTime := time.Since(start)

	// 验证数据
	assert.Equal(t, objectCount, len(decompressed), "Gzip大型数组解压缩后数量不一致: 期望=%d, 实际=%d", objectCount, len(decompressed))

	// 验证部分数据
	assert.Equal(t, complexObjects[0].Name, decompressed[0].Name, "Gzip大型数组数据验证失败")

	t.Logf("Gzip大型对象数组测试成功: %d个对象, 压缩耗时=%v, 解压耗时=%v, 压缩后大小=%d字节",
		objectCount, compressTime, decompressTime, len(compressed))
}

// BenchmarkGzipCompress 基准测试 - Gzip 原始压缩
func BenchmarkGzipCompress(b *testing.B) {
	testData := []byte("Hello, World! This is a benchmark test for gzip compression performance.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GzipCompress(testData)
		assert.NoError(b, err)
	}
}

// BenchmarkGzipCompressObject 基准测试 - Gzip 对象压缩
func BenchmarkGzipCompressObject(b *testing.B) {
	testMsg := TestMessage{
		ID:        "gzip_bench_msg_123",
		Content:   "这是一条用于Gzip性能测试的消息",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "gzip_bench_user",
		Type:      1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GzipCompressObject(testMsg)
		assert.NoError(b, err)
	}
}

// BenchmarkGzipVsZlibPerformance 性能对比基准测试
func BenchmarkGzipVsZlibPerformance(b *testing.B) {
	// 准备测试数据（截断时间以移除单调时钟）
	baseTime := time.Now().Truncate(time.Second)
	messages := make([]TestMessage, 50)
	for i := 0; i < 50; i++ {
		messages[i] = TestMessage{
			ID:        "perf_msg_" + string(rune(i)),
			Content:   "性能对比测试消息内容，用于比较Gzip和Zlib的性能差异",
			Timestamp: baseTime.Add(time.Duration(i) * time.Second),
			UserID:    "perf_user",
			Type:      i%3 + 1,
		}
	}

	b.Run("Gzip", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := GzipCompressObject(messages)
			assert.NoError(b, err)
		}
	})

	b.Run("Zlib", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := ZlibCompressObject(messages)
			assert.NoError(b, err)
		}
	})
}

// TestGzipErrorHandling 测试Gzip错误处理
func TestGzipErrorHandling(t *testing.T) {
	// 测试解压缩空数据
	_, err := GzipDecompress([]byte{})
	assert.Error(t, err, "Expected error for empty compressed data")

	// 测试解压缩无效数据
	_, err = GzipDecompress([]byte{0x00, 0x01, 0x02})
	assert.Error(t, err, "Expected error for invalid compressed data")

	t.Logf("Gzip错误处理测试成功")
}

// TestGzipEmptyData 测试Gzip空数据处理
func TestGzipEmptyData(t *testing.T) {
	// 压缩空数据
	compressedData, err := GzipCompress([]byte{})
	assert.NoError(t, err, "Compression error for empty data")
	assert.NotZero(t, len(compressedData), "Compressed data for empty input is empty")

	// 解压缩压缩后的空数据
	decompressedData, err := GzipDecompress(compressedData)
	assert.NoError(t, err, "Decompression error for empty compressed data")
	assert.Empty(t, decompressedData, "Decompressed data for empty input should be empty")

	t.Logf("Gzip空数据处理测试成功")
}

// TestGzipMegabyteData 测试Gzip大数据处理
func TestGzipMegabyteData(t *testing.T) {
	// 创建一个大数据(1MB)
	originalData := bytes.Repeat([]byte("A"), 1<<20)

	// 压缩数据
	start := time.Now()
	compressedData, err := GzipCompress(originalData)
	assert.NoError(t, err, "Compression error")
	assert.NotZero(t, len(compressedData), "Compressed data is empty")
	compressTime := time.Since(start)

	// 解压缩数据
	start = time.Now()
	decompressedData, err := GzipDecompress(compressedData)
	assert.NoError(t, err, "Decompression error")
	assert.True(t, bytes.Equal(originalData, decompressedData), "Decompressed data does not match original data")
	decompressTime := time.Since(start)

	// 计算压缩率
	compressionRatio := float64(len(compressedData)) / float64(len(originalData)) * 100

	t.Logf("Gzip大数据测试成功: 原始=%d字节, 压缩后=%d字节, 压缩率=%.2f%%, 压缩耗时=%v, 解压耗时=%v",
		len(originalData), len(compressedData), compressionRatio, compressTime, decompressTime)
}

// TestMultiGzipCompressObject 测试泛型多次压缩
func TestMultiGzipCompressObject(t *testing.T) {
	// 创建测试消息
	testMsg := TestMessage{
		ID:        "multi_gzip_msg_123",
		Content:   "这是一条用于多次Gzip压缩测试的消息",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "multi_gzip_user",
		Type:      1,
	}

	// 泛型3轮压缩
	compressed, err := MultiGZipCompressObject(testMsg, 3)
	assert.NoError(t, err)

	// 泛型3轮解压缩
	decompressed, err := MultiGZipDecompressObject[TestMessage](compressed, 3)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, reflect.DeepEqual(testMsg, decompressed), "Gzip多次泛型压缩解压缩后数据不一致")

	t.Logf("Gzip泛型多次压缩测试成功: 压缩后大小=%d字节", len(compressed))
}
