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
	"strings"
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

// TestZlibCompressWithPrefix 测试带前缀的压缩功能
func TestZlibCompressWithPrefix(t *testing.T) {
	testData := []byte("Hello, World! This is a test for prefix compression.")

	// 压缩并添加前缀
	compressed, err := ZlibCompressWithPrefix(testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, compressed, "压缩后数据为空")

	// 验证前缀存在
	assert.True(t, IsZlibCompressed(compressed), "压缩数据应该带有ZLIB前缀")
	assert.True(t, len(compressed) > ZlibPrefixLen, "压缩数据长度应大于前缀长度")
	assert.Equal(t, ZlibPrefix, string(compressed[:ZlibPrefixLen]), "前缀不匹配")

	t.Logf("带前缀压缩成功: 原始=%d字节, 压缩后=%d字节(含前缀)", len(testData), len(compressed))
}

// TestZlibDecompressWithPrefix 测试带前缀的解压缩功能
func TestZlibDecompressWithPrefix(t *testing.T) {
	testData := []byte("Hello, World! Testing decompression with prefix.")

	// 压缩并添加前缀
	compressed, err := ZlibCompressWithPrefix(testData)
	assert.NoError(t, err)

	// 解压缩（自动识别前缀）
	decompressed, err := ZlibSmartDecompress(compressed)
	assert.NoError(t, err)

	// 验证数据一致性
	assert.True(t, bytes.Equal(testData, decompressed), "解压缩后数据不一致")

	t.Logf("带前缀解压缩成功: 解压后=%d字节", len(decompressed))
}

// TestZlibDecompressWithPrefix_NoPrefix 测试无前缀数据的解压缩
func TestZlibDecompressWithPrefixNoPrefix(t *testing.T) {
	testData := []byte("This data has no compression prefix")

	// 直接解压（没有前缀应该返回原数据）
	decompressed, err := ZlibSmartDecompress(testData)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(testData, decompressed), "无前缀数据应该原样返回")

	t.Logf("无前缀数据处理成功: 原样返回%d字节", len(decompressed))
}

// TestIsZlibCompressed 测试压缩前缀检测
func TestIsZlibCompressed(t *testing.T) {
	// 测试带前缀的数据
	compressedData, _ := ZlibCompressWithPrefix([]byte("test data"))
	assert.True(t, IsZlibCompressed(compressedData), "应该检测到ZLIB前缀")

	// 测试不带前缀的数据
	normalData := []byte("normal data without prefix")
	assert.False(t, IsZlibCompressed(normalData), "不应该检测到ZLIB前缀")

	// 测试空数据
	emptyData := []byte("")
	assert.False(t, IsZlibCompressed(emptyData), "空数据不应该检测到前缀")

	t.Log("前缀检测测试通过")
}

// TestZlibPrefixRoundTrip 测试带前缀的完整往返
func TestZlibPrefixRoundTrip(t *testing.T) {
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
			compressed, err := ZlibCompressWithPrefix(tc.data)
			assert.NoError(t, err)
			assert.True(t, IsZlibCompressed(compressed))

			// 解压
			decompressed, err := ZlibSmartDecompress(compressed)
			assert.NoError(t, err)
			assert.True(t, bytes.Equal(tc.data, decompressed))

			t.Logf("%s: 原始=%d字节, 压缩=%d字节", tc.name, len(tc.data), len(compressed))
		})
	}
}

// TestZlibSmartDecompress 测试智能解压缩功能
func TestZlibSmartDecompress(t *testing.T) {
	testData := []byte("Hello, World! This is a test for smart decompression.")

	t.Run("压缩数据（带前缀）", func(t *testing.T) {
		// 压缩并添加前缀
		compressed, err := ZlibCompressWithPrefix(testData)
		assert.NoError(t, err)

		// 智能解压缩
		decompressed, err := ZlibSmartDecompress(compressed)
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(testData, decompressed), "智能解压缩后数据不一致")

		t.Logf("带前缀压缩数据智能解压成功: %d字节", len(decompressed))
	})

	t.Run("压缩数据（无前缀）", func(t *testing.T) {
		// 压缩但不添加前缀
		compressed, err := ZlibCompress(testData)
		assert.NoError(t, err)

		// 智能解压缩（应该能自动识别并解压）
		decompressed, err := ZlibSmartDecompress(compressed)
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(testData, decompressed), "智能解压缩后数据不一致")

		t.Logf("无前缀压缩数据智能解压成功: %d字节", len(decompressed))
	})

	t.Run("未压缩数据", func(t *testing.T) {
		// 未压缩的原始数据
		decompressed, err := ZlibSmartDecompress(testData)
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(testData, decompressed), "未压缩数据应该原样返回")

		t.Logf("未压缩数据智能处理成功: %d字节", len(decompressed))
	})

	t.Run("空数据", func(t *testing.T) {
		emptyData := []byte{}
		decompressed, err := ZlibSmartDecompress(emptyData)
		assert.NoError(t, err)
		assert.Empty(t, decompressed, "空数据应该返回空")

		t.Log("空数据智能处理成功")
	})
}

// TestZlibSmartDecompressObject 测试智能解压缩对象功能
func TestZlibSmartDecompressObject(t *testing.T) {
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
		compressed, err := ZlibCompressWithPrefix(jsonData)
		assert.NoError(t, err)

		// 智能解压缩对象
		decompressed, err := ZlibSmartDecompressObject[TestMessage](compressed)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testMsg, decompressed), "智能解压缩对象后数据不一致")

		t.Logf("带前缀压缩对象智能解压成功")
	})

	t.Run("压缩对象（无前缀）", func(t *testing.T) {
		// 使用标准压缩（无前缀）
		compressed, err := ZlibCompressObject(testMsg)
		assert.NoError(t, err)

		// 智能解压缩对象
		decompressed, err := ZlibSmartDecompressObject[TestMessage](compressed)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testMsg, decompressed), "智能解压缩对象后数据不一致")

		t.Logf("无前缀压缩对象智能解压成功")
	})

	t.Run("未压缩JSON数据", func(t *testing.T) {
		// 直接序列化，不压缩
		jsonData, err := json.Marshal(testMsg)
		assert.NoError(t, err)

		// 智能解压缩对象（应该能识别未压缩的JSON）
		decompressed, err := ZlibSmartDecompressObject[TestMessage](jsonData)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testMsg, decompressed), "未压缩JSON数据应该能正确解析")

		t.Logf("未压缩JSON数据智能处理成功")
	})
}

// TestZlibSmartDecompressBackwardCompatibility 测试向后兼容性
func TestZlibSmartDecompressBackwardCompatibility(t *testing.T) {
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
		decompressed, err := ZlibSmartDecompressObject[TestMessage](oldData)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(testMsg, decompressed), "应该能读取旧的未压缩数据")

		t.Log("✅ 能够读取旧的未压缩数据")
	})

	t.Run("场景2: 新数据已压缩，旧代码禁用压缩", func(t *testing.T) {
		// 模拟新数据：压缩存储
		compressed, err := ZlibCompressObject(testMsg)
		assert.NoError(t, err)

		// 新代码使用智能解压缩读取
		decompressed, err := ZlibSmartDecompressObject[TestMessage](compressed)
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
		result1, err := ZlibSmartDecompressObject[TestMessage](uncompressedData)
		assert.NoError(t, err)
		assert.Equal(t, messages[0].ID, result1.ID)

		// 第二条压缩
		compressedData, _ := ZlibCompressObject(messages[1])
		result2, err := ZlibSmartDecompressObject[TestMessage](compressedData)
		assert.NoError(t, err)
		assert.Equal(t, messages[1].ID, result2.ID)

		t.Log("✅ 能够同时处理压缩和未压缩的混合数据")
	})
}

// TestZlibSmartDecompressThresholdChange 测试压缩阈值变化场景
func TestZlibSmartDecompressThresholdChange(t *testing.T) {
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
		result, err := ZlibSmartDecompressObject[TestMessage](smallData)
		assert.NoError(t, err)
		assert.Equal(t, smallMsg.ID, result.ID)

		t.Log("✅ 小消息（未压缩）处理成功")
	})

	t.Run("场景2: 阈值1024，大消息压缩", func(t *testing.T) {
		// 大消息压缩
		largeData, _ := json.Marshal(largeMsg)
		t.Logf("大消息大小: %d字节", len(largeData))

		compressed, err := ZlibCompress(largeData)
		assert.NoError(t, err)
		t.Logf("压缩后大小: %d字节", len(compressed))

		// 智能解压缩应该能处理
		result, err := ZlibSmartDecompressObject[TestMessage](compressed)
		assert.NoError(t, err)
		assert.Equal(t, largeMsg.ID, result.ID)

		t.Log("✅ 大消息（已压缩）处理成功")
	})

	t.Run("场景3: 提高阈值到2048，之前压缩的数据仍能读取", func(t *testing.T) {
		// 之前在1024阈值下压缩的数据
		largeData, _ := json.Marshal(largeMsg)
		compressed, _ := ZlibCompress(largeData)

		// 现在提高阈值到2048，智能解压缩仍能处理
		result, err := ZlibSmartDecompressObject[TestMessage](compressed)
		assert.NoError(t, err)
		assert.Equal(t, largeMsg.ID, result.ID)

		t.Log("✅ 提高阈值后，旧的压缩数据仍能读取")
	})

	t.Run("场景4: 降低阈值到512，之前未压缩的数据仍能读取", func(t *testing.T) {
		// 之前在1024阈值下未压缩的小消息
		smallData, _ := json.Marshal(smallMsg)

		// 现在降低阈值到512，智能解压缩仍能处理
		result, err := ZlibSmartDecompressObject[TestMessage](smallData)
		assert.NoError(t, err)
		assert.Equal(t, smallMsg.ID, result.ID)

		t.Log("✅ 降低阈值后，旧的未压缩数据仍能读取")
	})
}

// TestZlibSmartDecompressConcurrency 测试智能解压缩的并发安全性
func TestZlibSmartDecompressConcurrency(t *testing.T) {
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
			compressedData[i], _ = ZlibCompressObject(msg)
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
				_, err := ZlibSmartDecompressObject[TestMessage](data)
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

// TestZlibSmartDecompressEdgeCases 测试智能解压缩的边界情况
func TestZlibSmartDecompressEdgeCases(t *testing.T) {
	t.Run("空数据", func(t *testing.T) {
		result, err := ZlibSmartDecompress([]byte{})
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("只有前缀无数据", func(t *testing.T) {
		prefixOnly := []byte(ZlibPrefix)
		result, err := ZlibSmartDecompress(prefixOnly)
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
		_, err := ZlibSmartDecompressObject[TestMessage](invalidJSON)
		assert.Error(t, err, "无效JSON应该返回错误")
	})

	t.Run("损坏的压缩数据", func(t *testing.T) {
		// 创建看起来像压缩数据但实际损坏的数据
		corruptedData := []byte{0x78, 0x9c, 0x00, 0x00, 0x00, 0xFF, 0xFF}
		result, err := ZlibSmartDecompress(corruptedData)
		// 智能解压缩应该能处理：尝试解压失败后返回原数据
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		t.Log("损坏的压缩数据被当作未压缩数据返回")
	})

	t.Run("超大数据", func(t *testing.T) {
		// 创建10MB数据
		largeData := bytes.Repeat([]byte("X"), 10<<20)
		compressed, err := ZlibCompress(largeData)
		assert.NoError(t, err)

		decompressed, err := ZlibSmartDecompress(compressed)
		assert.NoError(t, err)
		assert.Equal(t, len(largeData), len(decompressed))
		t.Logf("超大数据处理成功: 原始=%dMB, 压缩=%dMB", len(largeData)>>20, len(compressed)>>20)
	})
}

// BenchmarkZlibSmartDecompress 基准测试 - 智能解压缩
func BenchmarkZlibSmartDecompress(b *testing.B) {
	testData := []byte("Hello, World! This is a benchmark test for smart decompression.")

	b.Run("压缩数据", func(b *testing.B) {
		compressed, _ := ZlibCompress(testData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ZlibSmartDecompress(compressed)
		}
	})

	b.Run("未压缩数据", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ZlibSmartDecompress(testData)
		}
	})
}

// BenchmarkZlibSmartDecompressObject 基准测试 - 智能解压缩对象
func BenchmarkZlibSmartDecompressObject(b *testing.B) {
	testMsg := TestMessage{
		ID:        "bench_msg",
		Content:   "性能测试消息",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "bench_user",
		Type:      1,
	}

	b.Run("压缩对象", func(b *testing.B) {
		compressed, _ := ZlibCompressObject(testMsg)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ZlibSmartDecompressObject[TestMessage](compressed)
		}
	})

	b.Run("未压缩JSON", func(b *testing.B) {
		jsonData, _ := json.Marshal(testMsg)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ZlibSmartDecompressObject[TestMessage](jsonData)
		}
	})
}

// TestZlibCompressWithInfo 测试 Zlib 压缩并返回压缩信息
func TestZlibCompressWithInfo(t *testing.T) {
	testData := []byte(strings.Repeat("Test Data ", 100))

	result, err := ZlibCompressWithInfo(testData)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, len(testData), result.OriginalSize)
	assert.Less(t, result.CompressedSize, result.OriginalSize)

	// 验证压缩数据可以解压
	decompressed, err := ZlibDecompress(result.Data)
	assert.NoError(t, err)
	assert.Equal(t, testData, decompressed)
}

// TestMultiZlibCompressWithInfo 测试多次 Zlib 压缩并返回压缩信息
func TestMultiZlibCompressWithInfo(t *testing.T) {
	testData := []byte(strings.Repeat("Multi zlib test ", 50))

	result, err := MultiZlibCompressWithInfo(testData, 2)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, len(testData), result.OriginalSize)

	// 验证多次压缩数据可以解压
	decompressed, err := MultiZlibDecompress(result.Data, 2)
	assert.NoError(t, err)
	assert.Equal(t, testData, decompressed)
}

// TestZlibCompressObjectWithInfo 测试 Zlib 压缩对象并返回压缩信息
func TestZlibCompressObjectWithInfo(t *testing.T) {
	obj := TestMessage{
		ID:        "zlib_test_456",
		Content:   "zlib test content",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "user_zlib",
		Type:      1,
	}

	result, err := ZlibCompressObjectWithInfo(obj)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Greater(t, result.OriginalSize, 0)

	// 验证可以解压并反序列化
	decompressed, err := ZlibDecompressObject[TestMessage](result.Data)
	assert.NoError(t, err)
	assert.Equal(t, obj.ID, decompressed.ID)
	assert.Equal(t, obj.Content, decompressed.Content)
	assert.Equal(t, obj.UserID, decompressed.UserID)
	assert.Equal(t, obj.Type, decompressed.Type)
}

// TestMultiZlibCompressObjectWithInfo 测试多次 Zlib 压缩对象并返回压缩信息
func TestMultiZlibCompressObjectWithInfo(t *testing.T) {
	obj := TestMessage{
		ID:        "multi_zlib_999",
		Content:   "multi zlib content",
		Timestamp: time.Now().Truncate(time.Second),
		UserID:    "user_multi_zlib",
		Type:      3,
	}

	result, err := MultiZlibCompressObjectWithInfo(obj, 3)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Greater(t, result.OriginalSize, 0)

	// 验证可以解压并反序列化
	decompressed, err := MultiZlibDecompressObject[TestMessage](result.Data, 3)
	assert.NoError(t, err)
	assert.Equal(t, obj.ID, decompressed.ID)
	assert.Equal(t, obj.Content, decompressed.Content)
	assert.Equal(t, obj.UserID, decompressed.UserID)
}

// TestZlibCompressWithPrefixInfo 测试带前缀的 Zlib 压缩并返回压缩信息
func TestZlibCompressWithPrefixInfo(t *testing.T) {
	testData := []byte(strings.Repeat("Zlib prefix test ", 50))

	result, err := ZlibCompressWithPrefixInfo(testData)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, len(testData), result.OriginalSize)
	assert.Greater(t, result.CompressedSize, 0)

	// 验证数据包含前缀
	assert.True(t, IsZlibCompressed(result.Data))

	// 验证可以智能解压
	decompressed, err := ZlibSmartDecompress(result.Data)
	assert.NoError(t, err)
	assert.Equal(t, testData, decompressed)
}
