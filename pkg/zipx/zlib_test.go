/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-04 18:15:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 18:29:50
 * @FilePath: \go-toolbox\pkg\zipx\zlib_test.go
 * @Description: Zlib 压缩解压缩测试
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

// TestMessage 测试用的消息结构体
type TestMessage struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Type      int       `json:"type"`
}

// TestComplexObject 复杂测试对象
type TestComplexObject struct {
	Name     string                 `json:"name"`
	Age      int                    `json:"age"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Messages []TestMessage          `json:"messages"`
}

// TestZlibCompress 测试基本压缩功能
func TestZlibCompress(t *testing.T) {
	testData := []byte("Hello, World! This is a test string for compression.")

	compressed, err := ZlibCompress(testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressed, "压缩后数据为空")

	// 验证压缩后数据长度应该小于原始数据（对于重复内容）
	largeData := bytes.Repeat([]byte("A"), 1000)
	compressedLarge, err := ZlibCompress(largeData)
	assert.NoError(t, err)

	assert.Less(t, len(compressedLarge), len(largeData), "压缩效果不佳: 原始=%d, 压缩后=%d", len(largeData), len(compressedLarge))
	t.Logf("压缩测试成功: 原始=%d字节, 压缩后=%d字节, 压缩率=%.2f%%",
		len(largeData), len(compressedLarge), float64(len(compressedLarge))/float64(len(largeData))*100)
}

// TestZlibDecompress 测试基本解压缩功能
func TestZlibDecompress(t *testing.T) {
	testData := []byte("Hello, World! This is a test string for compression and decompression.")

	// 先压缩
	compressed, err := ZlibCompress(testData)
	assert.NoError(t, err)

	// 再解压缩
	decompressed, err := ZlibDecompress(compressed)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, bytes.Equal(testData, decompressed), "解压缩后数据不一致:\n原始: %s\n解压后: %s", testData, decompressed)
	t.Logf("解压缩测试成功: 原始=%d字节, 解压后=%d字节", len(testData), len(decompressed))
}

// TestZlibCompressObjectSingleMessage 测试单个消息对象压缩
func TestZlibCompressObjectSingleMessage(t *testing.T) {
	// 创建测试消息（截断时间以移除单调时钟）
	testMsg := TestMessage{
		ID:        "msg_123456",
		Content:   "这是一条测试消息，包含中文内容！",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "user_789",
		Type:      1,
	}

	// 压缩对象
	compressed, err := ZlibCompressObject(testMsg)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressed, "压缩后数据为空")

	// 解压缩对象
	decompressed, err := ZlibDecompressObject[TestMessage](compressed)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, reflect.DeepEqual(testMsg, decompressed), "对象压缩解压缩后数据不一致:\n原始: %+v\n解压后: %+v", testMsg, decompressed)
	t.Logf("对象压缩测试成功: 压缩后大小=%d字节", len(compressed))
}

// TestZlibCompressObjectWithSize 测试带原始大小返回的泛型压缩函数
func TestZlibCompressObjectWithSize(t *testing.T) {
	// 创建测试消息（截断时间以移除单调时钟）
	testMsg := TestMessage{
		ID:        "msg_123456",
		Content:   "这是一条测试消息，包含中文内容！",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "user_789",
		Type:      1,
	}

	// 使用新函数进行压缩，同时获取原始大小
	compressed, originalSize, err := ZlibCompressObjectWithSize(testMsg)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressed, "压缩后数据为空")
	assert.Greater(t, originalSize, 0, "原始数据大小应该大于0")

	// 验证原始大小的正确性
	jsonData, _ := json.Marshal(testMsg)
	assert.Equal(t, len(jsonData), originalSize, "返回的原始大小应该与手动序列化的大小一致")

	// 解压缩验证
	decompressed, err := ZlibDecompressObject[TestMessage](compressed)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(testMsg, decompressed), "对象压缩解压缩后数据不一致")

	t.Logf("带大小的压缩测试成功: 原始大小=%d字节, 压缩后大小=%d字节, 压缩率=%.2f%%",
		originalSize, len(compressed), float64(len(compressed))/float64(originalSize)*100)
}

// TestZlibCompressObjectMessageSlice 测试消息数组压缩
func TestZlibCompressObjectMessageSlice(t *testing.T) {
	// 创建多条测试消息（截断时间以移除单调时钟）
	baseTime := time.Now().Truncate(time.Second)
	messages := []TestMessage{
		{
			ID:        "msg_001",
			Content:   "第一条消息",
			Timestamp: baseTime,
			UserID:    "user_001",
			Type:      1,
		},
		{
			ID:        "msg_002",
			Content:   "第二条消息，内容更长一些，用于测试压缩效果",
			Timestamp: baseTime.Add(1 * time.Minute),
			UserID:    "user_002",
			Type:      2,
		},
		{
			ID:        "msg_003",
			Content:   "第三条消息",
			Timestamp: baseTime.Add(2 * time.Minute),
			UserID:    "user_003",
			Type:      1,
		},
	}

	// 压缩消息数组
	compressed, err := ZlibCompressObject(messages)
	assert.NoError(t, err)

	// 解压缩消息数组
	decompressed, err := ZlibDecompressObject[[]TestMessage](compressed)
	assert.NoError(t, err)

	// 验证数组长度
	assert.Equal(t, len(messages), len(decompressed), "消息数组长度不一致: 原始=%d, 解压后=%d", len(messages), len(decompressed))

	// 验证每条消息
	for i, original := range messages {
		assert.True(t, reflect.DeepEqual(original, decompressed[i]), "第%d条消息不一致:\n原始: %+v\n解压后: %+v", i, original, decompressed[i])
	}

	t.Logf("消息数组压缩测试成功: %d条消息, 压缩后大小=%d字节", len(messages), len(compressed))
}

// TestZlibCompressObjectComplexObject 测试复杂对象压缩
func TestZlibCompressObjectComplexObject(t *testing.T) {
	// 创建复杂测试对象
	complexObj := TestComplexObject{
		Name: "测试用户",
		Age:  25,
		Tags: []string{"开发者", "Go语言", "微服务"},
		Metadata: map[string]interface{}{
			"role":        "admin",
			"permissions": []string{"read", "write", "delete"},
			"config": map[string]interface{}{
				"theme":    "dark",
				"language": "zh-CN",
			},
		},
		Messages: []TestMessage{
			{
				ID:        "msg_complex_001",
				Content:   "复杂对象中的消息1",
				Timestamp: time.Now(),
				UserID:    "complex_user",
				Type:      1,
			},
			{
				ID:        "msg_complex_002",
				Content:   "复杂对象中的消息2",
				Timestamp: time.Now().Add(30 * time.Second),
				UserID:    "complex_user",
				Type:      2,
			},
		},
	}

	// 压缩复杂对象
	compressed, err := ZlibCompressObject(complexObj)
	assert.NoError(t, err)

	// 解压缩复杂对象
	decompressed, err := ZlibDecompressObject[TestComplexObject](compressed)
	assert.NoError(t, err)

	// 验证基本字段
	assert.Equal(t, complexObj.Name, decompressed.Name)
	assert.Equal(t, complexObj.Age, decompressed.Age)

	// 验证标签
	assert.True(t, reflect.DeepEqual(complexObj.Tags, decompressed.Tags), "标签不一致")

	// 验证消息
	assert.Equal(t, len(complexObj.Messages), len(decompressed.Messages), "消息数量不一致")

	t.Logf("复杂对象压缩测试成功: 压缩后大小=%d字节", len(compressed))
}

// TestMultiZlibCompress 测试多重压缩
func TestMultiZlibCompress(t *testing.T) {
	testData := []byte("Hello, World! This is a test for multiple compression rounds.")

	// 测试2轮压缩
	compressed, err := MultiZlibCompress(testData, 2)
	assert.NoError(t, err)

	// 测试2轮解压缩
	decompressed, err := MultiZlibDecompress(compressed, 2)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, bytes.Equal(testData, decompressed), "多重压缩解压后数据不一致")
	t.Logf("多重压缩测试成功: 原始=%d字节, 2轮压缩后=%d字节", len(testData), len(compressed))
}

// TestZlibConcurrency 测试并发安全性
func TestZlibConcurrency(t *testing.T) {
	const numGoroutines = 100
	const numOperations = 50

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperations)

	// 创建测试数据
	messages := make([]TestMessage, 10)
	for i := 0; i < 10; i++ {
		messages[i] = TestMessage{
			ID:        "msg_" + string(rune('A'+i)),
			Content:   "并发测试消息内容",
			Timestamp: time.Now(),
			UserID:    "concurrent_user",
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
				compressed, err := ZlibCompressObject(messages)
				if err != nil {
					errors <- err
					return
				}

				// 解压缩
				decompressed, err := ZlibDecompressObject[[]TestMessage](compressed)
				if err != nil {
					errors <- err
					return
				}

				// 验证
				if len(decompressed) != len(messages) {
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
		assert.Error(t, err, "并发测试出错: %v", err)
		errorCount++
	}

	if errorCount == 0 {
		t.Logf("并发测试成功: %d个协程, 每个执行%d次操作", numGoroutines, numOperations)
	}
}

// TestZlibEmptyData 测试空数据处理
func TestZlibEmptyData(t *testing.T) {
	// 测试空字节数组
	emptyData := []byte{}
	compressed, err := ZlibCompress(emptyData)
	assert.NoError(t, err)

	decompressed, err := ZlibDecompress(compressed)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(decompressed), "解压缩后空数据长度不为0: %d", len(decompressed))

	// 测试空对象数组
	emptyMessages := []TestMessage{}
	compressedObj, err := ZlibCompressObject(emptyMessages)
	assert.NoError(t, err)

	decompressedObj, err := ZlibDecompressObject[[]TestMessage](compressedObj)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(decompressedObj), "解压缩后空对象数组长度不为0: %d", len(decompressedObj))

	t.Log("空数据处理测试成功")
}

// TestZlibLargeData 测试大数据压缩
func TestZlibLargeData(t *testing.T) {
	// 创建大量消息
	const messageCount = 1000
	messages := make([]TestMessage, messageCount)

	for i := 0; i < messageCount; i++ {
		messages[i] = TestMessage{
			ID:        "large_msg_" + string(rune(i)),
			Content:   "这是一条用于测试大数据压缩的消息，内容比较长，重复文本可以提高压缩效率。",
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
			UserID:    "large_test_user",
			Type:      i%3 + 1,
		}
	}

	start := time.Now()

	// 压缩大数据
	compressed, err := ZlibCompressObject(messages)
	assert.NoError(t, err)

	compressTime := time.Since(start)

	start = time.Now()

	// 解压缩大数据
	decompressed, err := ZlibDecompressObject[[]TestMessage](compressed)
	assert.NoError(t, err)

	decompressTime := time.Since(start)

	// 验证数据
	assert.Equal(t, messageCount, len(decompressed), "大数据解压缩后数量不一致: 期望=%d, 实际=%d", messageCount, len(decompressed))

	t.Logf("大数据测试成功: %d条消息, 压缩耗时=%v, 解压耗时=%v, 压缩后大小=%d字节",
		messageCount, compressTime, decompressTime, len(compressed))
}

// BenchmarkZlibCompress 基准测试 - 原始压缩
func BenchmarkZlibCompress(b *testing.B) {
	testData := []byte("Hello, World! This is a benchmark test for compression performance.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ZlibCompress(testData)
		assert.NoError(b, err)
	}
}

// BenchmarkZlibCompressObject 基准测试 - 对象压缩
func BenchmarkZlibCompressObject(b *testing.B) {
	testMsg := TestMessage{
		ID:        "bench_msg_123",
		Content:   "这是一条用于性能测试的消息",
		Timestamp: time.Now(),
		UserID:    "bench_user",
		Type:      1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ZlibCompressObject(testMsg)
		assert.NoError(b, err)
	}
}

// BenchmarkZlibCompressObjectSlice 基准测试 - 数组压缩
func BenchmarkZlibCompressObjectSlice(b *testing.B) {
	messages := make([]TestMessage, 100)
	for i := 0; i < 100; i++ {
		messages[i] = TestMessage{
			ID:        "bench_slice_msg_" + string(rune(i)),
			Content:   "性能测试消息内容",
			Timestamp: time.Now(),
			UserID:    "bench_slice_user",
			Type:      i%3 + 1,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ZlibCompressObject(messages)
		assert.NoError(b, err)
	}
}

// TestZlibErrorHandling 测试Zlib错误处理
func TestZlibErrorHandling(t *testing.T) {
	// 测试解压缩空数据
	_, err := ZlibDecompress([]byte{})
	assert.Error(t, err, "Expected error for empty compressed data")

	// 测试解压缩无效数据
	_, err = ZlibDecompress([]byte{0x00, 0x01, 0x02})
	assert.Error(t, err, "Expected error for invalid compressed data")

	t.Logf("Zlib错误处理测试成功")
}

// TestZlibMegabyteData 测试Zlib大数据处理
func TestZlibMegabyteData(t *testing.T) {
	// 创建一个大数据(1MB)
	originalData := bytes.Repeat([]byte("B"), 1<<20)

	// 压缩数据
	start := time.Now()
	compressedData, err := ZlibCompress(originalData)
	assert.NoError(t, err, "Compression error")
	assert.NotZero(t, len(compressedData), "Compressed data is empty")
	compressTime := time.Since(start)

	// 解压缩数据
	start = time.Now()
	decompressedData, err := ZlibDecompress(compressedData)
	assert.NoError(t, err, "Decompression error")
	assert.True(t, bytes.Equal(originalData, decompressedData), "Decompressed data does not match original data")
	decompressTime := time.Since(start)

	// 计算压缩率
	compressionRatio := float64(len(compressedData)) / float64(len(originalData)) * 100

	t.Logf("Zlib大数据测试成功: 原始=%d字节, 压缩后=%d字节, 压缩率=%.2f%%, 压缩耗时=%v, 解压耗时=%v",
		len(originalData), len(compressedData), compressionRatio, compressTime, decompressTime)
}

// TestMultiZlibCompressObject 测试泛型多次压缩
func TestMultiZlibCompressObject(t *testing.T) {
	// 创建测试消息
	testMsg := TestMessage{
		ID:        "multi_zlib_msg_123",
		Content:   "这是一条用于多次Zlib压缩测试的消息",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "multi_zlib_user",
		Type:      1,
	}

	// 泛型2轮压缩
	compressed, err := MultiZlibCompressObject(testMsg, 2)
	assert.NoError(t, err)

	// 泛型2轮解压缩
	decompressed, err := MultiZlibDecompressObject[TestMessage](compressed, 2)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, reflect.DeepEqual(testMsg, decompressed), "Zlib多次泛型压缩解压缩后数据不一致")

	t.Logf("Zlib泛型多次压缩测试成功: 压缩后大小=%d字节", len(compressed))
}
