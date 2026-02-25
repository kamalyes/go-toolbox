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
	"strings"
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

// TestGzipCompressWithPrefix 测试带前缀的压缩功能
func TestGzipCompressWithPrefix(t *testing.T) {
	testData := []byte("Hello, World! This is a test for prefix compression.")

	// 压缩并添加前缀
	compressed, err := GzipCompressWithPrefix(testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressed, "压缩后数据为空")

	// 验证前缀存在
	assert.True(t, IsGzipCompressed(compressed), "压缩数据应该带有GZIP前缀")
	assert.True(t, len(compressed) > GzipPrefixLen, "压缩数据长度应大于前缀长度")
	assert.Equal(t, GzipPrefix, string(compressed[:GzipPrefixLen]), "前缀不匹配")

	t.Logf("带前缀压缩成功: 原始=%d字节, 压缩后=%d字节(含前缀)", len(testData), len(compressed))
}

// TestGzipDecompressWithPrefix 测试带前缀的解压缩功能
func TestGzipDecompressWithPrefix(t *testing.T) {
	testData := []byte("Hello, World! Testing decompression with prefix.")

	// 压缩并添加前缀
	compressed, err := GzipCompressWithPrefix(testData)
	assert.NoError(t, err)

	// 解压缩（自动识别前缀）
	decompressed, err := GzipSmartDecompress(compressed)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, bytes.Equal(testData, decompressed), "解压缩后数据不一致")

	t.Logf("带前缀解压缩成功: 解压后=%d字节", len(decompressed))
}

// TestGzipDecompressWithPrefix_NoPrefix 测试无前缀数据的解压缩
func TestGzipDecompressWithPrefixNoPrefix(t *testing.T) {
	testData := []byte("This data has no compression prefix")

	// 直接解压（没有前缀应该返回原数据）
	decompressed, err := GzipSmartDecompress(testData)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(testData, decompressed), "无前缀数据应该原样返回")

	t.Logf("无前缀数据处理成功: 原样返回%d字节", len(decompressed))
}

// TestIsGzipCompressed 测试压缩前缀检测
func TestIsGzipCompressed(t *testing.T) {
	// 测试带前缀的数据
	compressedData, _ := GzipCompressWithPrefix([]byte("test data"))
	assert.True(t, IsGzipCompressed(compressedData), "应该检测到GZIP前缀")

	// 测试不带前缀的数据
	normalData := []byte("normal data without prefix")
	assert.False(t, IsGzipCompressed(normalData), "不应该检测到GZIP前缀")

	// 测试空数据
	emptyData := []byte("")
	assert.False(t, IsGzipCompressed(emptyData), "空数据不应该检测到前缀")

	t.Log("前缀检测测试通过")
}

// TestGzipPrefixRoundTrip 测试带前缀的完整往返
func TestGzipPrefixRoundTrip(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{"短文本", []byte("short text")},
		{"长文本", bytes.Repeat([]byte("long text data "), 100)},
		{"二进制数据", []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}},
		{"中文内容", []byte("这是一段中文测试内容，用于验证压缩和解压缩功能")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 压缩
			compressed, err := GzipCompressWithPrefix(tc.data)
			assert.NoError(t, err)
			assert.True(t, IsGzipCompressed(compressed))

			// 解压
			decompressed, err := GzipSmartDecompress(compressed)
			assert.NoError(t, err)
			assert.True(t, bytes.Equal(tc.data, decompressed))

			t.Logf("%s: 原始=%d字节, 压缩=%d字节", tc.name, len(tc.data), len(compressed))
		})
	}
}

// TestGzipSmartDecompress 测试智能解压缩功能
func TestGzipSmartDecompress(t *testing.T) {
	testData := []byte("Hello, World! This is a test for smart decompression.")

	t.Run("压缩数据（带前缀）", func(t *testing.T) {
		// 压缩并添加前缀
		compressed, err := GzipCompressWithPrefix(testData)
		assert.NoError(t, err)

		// 智能解压缩
		decompressed, err := GzipSmartDecompress(compressed)
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(testData, decompressed), "智能解压缩后数据不一致")

		t.Logf("带前缀压缩数据智能解压成功: %d字节", len(decompressed))
	})

	t.Run("压缩数据（无前缀）", func(t *testing.T) {
		// 压缩但不添加前缀
		compressed, err := GzipCompress(testData)
		assert.NoError(t, err)

		// 智能解压缩（应该能自动识别并解压）
		decompressed, err := GzipSmartDecompress(compressed)
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(testData, decompressed), "智能解压缩后数据不一致")

		t.Logf("无前缀压缩数据智能解压成功: %d字节", len(decompressed))
	})

	t.Run("未压缩数据", func(t *testing.T) {
		// 未压缩的原始数据
		decompressed, err := GzipSmartDecompress(testData)
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(testData, decompressed), "未压缩数据应该原样返回")

		t.Logf("未压缩数据智能处理成功: %d字节", len(decompressed))
	})

	t.Run("空数据", func(t *testing.T) {
		emptyData := []byte{}
		decompressed, err := GzipSmartDecompress(emptyData)
		assert.NoError(t, err)
		assert.Empty(t, decompressed, "空数据应该返回空")

		t.Log("空数据智能处理成功")
	})
}

// TestGzipSmartDecompressObject 测试智能解压缩对象功能
func TestGzipSmartDecompressObject(t *testing.T) {
	// 创建测试消息
	testMsg := TestMessage{
		ID:        "smart_msg_123",
		Content:   "这是一条用于智能解压缩测试的消息",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "smart_user_789",
		Type:      1,
	}

	t.Run("压缩对象（带前缀）", func(t *testing.T) {
		// 序列化并压缩
		jsonData, _ := json.Marshal(testMsg)
		compressed, err := GzipCompressWithPrefix(jsonData)
		assert.NoError(t, err)

		// 智能解压缩对象
		decompressed, err := GzipSmartDecompressObject[TestMessage](compressed)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testMsg, decompressed), "智能解压缩对象后数据不一致")

		t.Logf("带前缀压缩对象智能解压成功")
	})

	t.Run("压缩对象（无前缀）", func(t *testing.T) {
		// 使用标准压缩（无前缀）
		compressed, err := GzipCompressObject(testMsg)
		assert.NoError(t, err)

		// 智能解压缩对象
		decompressed, err := GzipSmartDecompressObject[TestMessage](compressed)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testMsg, decompressed), "智能解压缩对象后数据不一致")

		t.Logf("无前缀压缩对象智能解压成功")
	})

	t.Run("未压缩JSON数据", func(t *testing.T) {
		// 直接序列化，不压缩
		jsonData, err := json.Marshal(testMsg)
		assert.NoError(t, err)

		// 智能解压缩对象（应该能识别未压缩的JSON）
		decompressed, err := GzipSmartDecompressObject[TestMessage](jsonData)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testMsg, decompressed), "未压缩JSON数据应该能正确解析")

		t.Logf("未压缩JSON数据智能处理成功")
	})
}

// TestGzipSmartDecompressBackwardCompatibility 测试向后兼容性
func TestGzipSmartDecompressBackwardCompatibility(t *testing.T) {
	testMsg := TestMessage{
		ID:        "compat_msg_123",
		Content:   "向后兼容性测试消息",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "compat_user",
		Type:      1,
	}

	t.Run("场景1: 旧数据未压缩，新代码启用压缩", func(t *testing.T) {
		// 模拟旧数据：直接存储JSON（未压缩）
		oldData, err := json.Marshal(testMsg)
		assert.NoError(t, err)

		// 新代码使用智能解压缩读取
		decompressed, err := GzipSmartDecompressObject[TestMessage](oldData)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testMsg, decompressed), "应该能读取旧的未压缩数据")

		t.Log("✅ 能够读取旧的未压缩数据")
	})

	t.Run("场景2: 新数据已压缩，旧代码禁用压缩", func(t *testing.T) {
		// 模拟新数据：压缩存储
		compressed, err := GzipCompressObject(testMsg)
		assert.NoError(t, err)

		// 新代码使用智能解压缩读取
		decompressed, err := GzipSmartDecompressObject[TestMessage](compressed)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testMsg, decompressed), "应该能读取压缩数据")

		t.Log("✅ 能够读取压缩数据")
	})

	t.Run("场景3: 混合数据（部分压缩，部分未压缩）", func(t *testing.T) {
		// 创建多条消息
		messages := []TestMessage{
			testMsg,
			{
				ID:        "compat_msg_456",
				Content:   "第二条消息",
				Timestamp: time.Now().Truncate(time.Second),
				UserID:    "compat_user_2",
				Type:      2,
			},
		}

		// 第一条未压缩
		uncompressedData, _ := json.Marshal(messages[0])
		result1, err := GzipSmartDecompressObject[TestMessage](uncompressedData)
		assert.NoError(t, err)
		assert.Equal(t, messages[0].ID, result1.ID)

		// 第二条压缩
		compressedData, _ := GzipCompressObject(messages[1])
		result2, err := GzipSmartDecompressObject[TestMessage](compressedData)
		assert.NoError(t, err)
		assert.Equal(t, messages[1].ID, result2.ID)

		t.Log("✅ 能够同时处理压缩和未压缩的混合数据")
	})
}

// TestGzipSmartDecompressThresholdChange 测试压缩阈值变化场景
func TestGzipSmartDecompressThresholdChange(t *testing.T) {
	// 创建不同大小的测试数据
	smallMsg := TestMessage{
		ID:        "small_msg",
		Content:   "小消息",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "user_1",
		Type:      1,
	}

	largeMsg := TestMessage{
		ID:        "large_msg",
		Content:   string(bytes.Repeat([]byte("这是一条很长的消息内容，用于测试压缩阈值变化场景。"), 50)),
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "user_2",
		Type:      2,
	}

	t.Run("场景1: 阈值1024，小消息不压缩", func(t *testing.T) {
		// 小消息不压缩（假设小于1024字节）
		smallData, _ := json.Marshal(smallMsg)
		t.Logf("小消息大小: %d字节", len(smallData))

		// 智能解压缩应该能处理
		result, err := GzipSmartDecompressObject[TestMessage](smallData)
		assert.NoError(t, err)
		assert.Equal(t, smallMsg.ID, result.ID)

		t.Log("✅ 小消息（未压缩）处理成功")
	})

	t.Run("场景2: 阈值1024，大消息压缩", func(t *testing.T) {
		// 大消息压缩
		largeData, _ := json.Marshal(largeMsg)
		t.Logf("大消息大小: %d字节", len(largeData))

		compressed, err := GzipCompress(largeData)
		assert.NoError(t, err)
		t.Logf("压缩后大小: %d字节", len(compressed))

		// 智能解压缩应该能处理
		result, err := GzipSmartDecompressObject[TestMessage](compressed)
		assert.NoError(t, err)
		assert.Equal(t, largeMsg.ID, result.ID)

		t.Log("✅ 大消息（已压缩）处理成功")
	})

	t.Run("场景3: 提高阈值到2048，之前压缩的数据仍能读取", func(t *testing.T) {
		// 之前在1024阈值下压缩的数据
		largeData, _ := json.Marshal(largeMsg)
		compressed, _ := GzipCompress(largeData)

		// 现在提高阈值到2048，智能解压缩仍能处理
		result, err := GzipSmartDecompressObject[TestMessage](compressed)
		assert.NoError(t, err)
		assert.Equal(t, largeMsg.ID, result.ID)

		t.Log("✅ 提高阈值后，旧的压缩数据仍能读取")
	})

	t.Run("场景4: 降低阈值到512，之前未压缩的数据仍能读取", func(t *testing.T) {
		// 之前在1024阈值下未压缩的小消息
		smallData, _ := json.Marshal(smallMsg)

		// 现在降低阈值到512，智能解压缩仍能处理
		result, err := GzipSmartDecompressObject[TestMessage](smallData)
		assert.NoError(t, err)
		assert.Equal(t, smallMsg.ID, result.ID)

		t.Log("✅ 降低阈值后，旧的未压缩数据仍能读取")
	})
}

// TestGzipSmartDecompressConcurrency 测试智能解压缩的并发安全性
func TestGzipSmartDecompressConcurrency(t *testing.T) {
	const numGoroutines = 50
	const numOperations = 30

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperations)

	// 准备测试数据（混合压缩和未压缩）
	baseTime := time.Now().Truncate(time.Second)
	compressedData := make([][]byte, 5)
	uncompressedData := make([][]byte, 5)

	for i := 0; i < 5; i++ {
		msg := TestMessage{
			ID:        "concurrent_msg_" + string(rune('A'+i)),
			Content:   "并发测试消息",
			Timestamp: baseTime.Add(time.Duration(i) * time.Second),
			UserID:    "concurrent_user",
			Type:      i%3 + 1,
		}

		// 一半压缩，一半不压缩
		if i%2 == 0 {
			compressedData[i], _ = GzipCompressObject(msg)
		} else {
			uncompressedData[i], _ = json.Marshal(msg)
		}
	}

	// 启动多个协程进行并发智能解压缩
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				idx := (goroutineID + j) % 5

				// 随机选择压缩或未压缩数据
				var data []byte
				if idx%2 == 0 {
					data = compressedData[idx]
				} else {
					data = uncompressedData[idx]
				}

				// 智能解压缩
				_, err := GzipSmartDecompressObject[TestMessage](data)
				if err != nil {
					errors <- err
					return
				}
			}
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()
	close(errors)

	// 检查错误
	errorCount := 0
	for err := range errors {
		t.Errorf("并发测试出错: %v", err)
		errorCount++
	}

	if errorCount == 0 {
		t.Logf("✅ 智能解压缩并发测试成功: %d个协程, 每个执行%d次操作", numGoroutines, numOperations)
	}
}

// TestGzipSmartDecompressEdgeCases 测试智能解压缩的边界情况
func TestGzipSmartDecompressEdgeCases(t *testing.T) {
	t.Run("空数据", func(t *testing.T) {
		result, err := GzipSmartDecompress([]byte{})
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("只有前缀无数据", func(t *testing.T) {
		prefixOnly := []byte(GzipPrefix)
		result, err := GzipSmartDecompress(prefixOnly)
		// 应该尝试解压空数据，可能返回错误或空结果
		if err == nil {
			t.Log("前缀无数据情况处理成功")
		} else {
			t.Logf("前缀无数据返回错误（预期行为）: %v", err)
		}
		_ = result
	})

	t.Run("无效的JSON数据", func(t *testing.T) {
		invalidJSON := []byte("{invalid json")
		_, err := GzipSmartDecompressObject[TestMessage](invalidJSON)
		assert.Error(t, err, "无效JSON应该返回错误")
	})

	t.Run("损坏的压缩数据", func(t *testing.T) {
		// 创建看起来像压缩数据但实际损坏的数据
		corruptedData := []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF}
		result, err := GzipSmartDecompress(corruptedData)
		// 智能解压缩应该能处理：尝试解压失败后返回原数据
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		t.Log("损坏的压缩数据被当作未压缩数据返回")
	})

	t.Run("超大数据", func(t *testing.T) {
		// 创建10MB数据
		largeData := bytes.Repeat([]byte("X"), 10<<20)
		compressed, err := GzipCompress(largeData)
		assert.NoError(t, err)

		decompressed, err := GzipSmartDecompress(compressed)
		assert.NoError(t, err)
		assert.Equal(t, len(largeData), len(decompressed))
		t.Logf("超大数据处理成功: 原始=%dMB, 压缩=%dMB", len(largeData)>>20, len(compressed)>>20)
	})
}

// BenchmarkGzipSmartDecompress 基准测试 - 智能解压缩
func BenchmarkGzipSmartDecompress(b *testing.B) {
	testData := []byte("Hello, World! This is a benchmark test for smart decompression.")

	b.Run("压缩数据", func(b *testing.B) {
		compressed, _ := GzipCompress(testData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GzipSmartDecompress(compressed)
		}
	})

	b.Run("未压缩数据", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GzipSmartDecompress(testData)
		}
	})
}

// BenchmarkGzipSmartDecompressObject 基准测试 - 智能解压缩对象
func BenchmarkGzipSmartDecompressObject(b *testing.B) {
	testMsg := TestMessage{
		ID:        "bench_msg",
		Content:   "性能测试消息",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "bench_user",
		Type:      1,
	}

	b.Run("压缩对象", func(b *testing.B) {
		compressed, _ := GzipCompressObject(testMsg)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GzipSmartDecompressObject[TestMessage](compressed)
		}
	})

	b.Run("未压缩JSON", func(b *testing.B) {
		jsonData, _ := json.Marshal(testMsg)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = GzipSmartDecompressObject[TestMessage](jsonData)
		}
	})
}

// TestGzipCompressWithInfo 测试 Gzip 压缩并返回压缩信息
func TestGzipCompressWithInfo(t *testing.T) {
	testData := []byte(strings.Repeat("Hello World! ", 100))

	result, err := GzipCompressWithInfo(testData)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, len(testData), result.OriginalSize)
	assert.Less(t, result.CompressedSize, result.OriginalSize)
	assert.Less(t, result.Ratio, 1.0)

	// 验证压缩数据可以解压
	decompressed, err := GzipDecompress(result.Data)
	assert.NoError(t, err)
	assert.Equal(t, testData, decompressed)
}

// TestMultiGZipCompressWithInfo 测试多次 Gzip 压缩并返回压缩信息
func TestMultiGZipCompressWithInfo(t *testing.T) {
	testData := []byte(strings.Repeat("Multi compress test ", 50))

	result, err := MultiGZipCompressWithInfo(testData, 2)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, len(testData), result.OriginalSize)

	// 验证多次压缩数据可以解压
	decompressed, err := MultiGZipDecompress(result.Data, 2)
	assert.NoError(t, err)
	assert.Equal(t, testData, decompressed)
}

// TestGzipCompressObjectWithInfo 测试 Gzip 压缩对象并返回压缩信息
func TestGzipCompressObjectWithInfo(t *testing.T) {
	obj := TestMessage{
		ID:        "test_123",
		Content:   "test content",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "user_456",
		Type:      1,
	}

	result, err := GzipCompressObjectWithInfo(obj)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Greater(t, result.OriginalSize, 0)
	assert.Greater(t, result.CompressedSize, 0)

	// 验证可以解压并反序列化
	decompressed, err := GzipDecompressObject[TestMessage](result.Data)
	assert.NoError(t, err)
	assert.Equal(t, obj.ID, decompressed.ID)
	assert.Equal(t, obj.Content, decompressed.Content)
	assert.Equal(t, obj.UserID, decompressed.UserID)
	assert.Equal(t, obj.Type, decompressed.Type)
}

// TestMultiGZipCompressObjectWithInfo 测试多次 Gzip 压缩对象并返回压缩信息
func TestMultiGZipCompressObjectWithInfo(t *testing.T) {
	obj := TestMessage{
		ID:        "multi_gzip_789",
		Content:   "multi compress content",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "user_multi",
		Type:      2,
	}

	result, err := MultiGZipCompressObjectWithInfo(obj, 3)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Greater(t, result.OriginalSize, 0)

	// 验证可以解压并反序列化
	decompressed, err := MultiGZipDecompressObject[TestMessage](result.Data, 3)
	assert.NoError(t, err)
	assert.Equal(t, obj.ID, decompressed.ID)
	assert.Equal(t, obj.Content, decompressed.Content)
	assert.Equal(t, obj.UserID, decompressed.UserID)
}

// TestGzipCompressWithPrefixInfo 测试带前缀的 Gzip 压缩并返回压缩信息
func TestGzipCompressWithPrefixInfo(t *testing.T) {
	testData := []byte(strings.Repeat("Prefix test data ", 50))

	result, err := GzipCompressWithPrefixInfo(testData)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, len(testData), result.OriginalSize)
	assert.Greater(t, result.CompressedSize, 0)

	// 验证数据包含前缀
	assert.True(t, IsGzipCompressed(result.Data))

	// 验证可以智能解压
	decompressed, err := GzipSmartDecompress(result.Data)
	assert.NoError(t, err)
	assert.Equal(t, testData, decompressed)
}
