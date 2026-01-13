/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 18:26:12
 * @FilePath: \go-toolbox\pkg\zipx\zlib.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package zipx

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"sync"
)

const (
	// ZlibPrefix 是用于标识 zlib 压缩数据的前缀
	ZlibPrefix = "ZLIB:"
	// ZlibPrefixLen 是 zlib 前缀的长度
	ZlibPrefixLen = len(ZlibPrefix)
)

// 创建一个 sync.Pool 来复用 bytes.Buffer
var (
	zlibBuffer = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer) // 创建新的 bytes.Buffer
		},
	}
	// 添加读取缓冲区池用于解压缩优化
	zlibReadBuffer = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

// ZlibCompress 压缩数据（修复版本）
func ZlibCompress(data []byte) ([]byte, error) {
	// 从池中获取一个缓冲区
	buf := zlibBuffer.Get().(*bytes.Buffer)
	defer zlibBuffer.Put(buf) // 使用后放回池中

	// 清空缓冲区
	buf.Reset()

	// 创建新的 zlib.Writer（不要重用，因为 zlib.Writer 没有有效的 Reset 方法）
	writer := zlib.NewWriter(buf)
	defer writer.Close()

	if _, err := writer.Write(data); err != nil {
		return nil, err // 写入数据时出错
	}

	if err := writer.Close(); err != nil {
		return nil, err // 关闭 writer 时出错
	}

	// 必须返回副本!不能返回buf.Bytes(),因为buf会被放回Pool重用
	// 直接返回buf.Bytes()会导致并发竞争,其他goroutine可能修改同一buffer
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// ZlibDecompress 解压缩数据（优化版本，使用对象池）
func ZlibDecompress(compressedData []byte) ([]byte, error) {
	// 创建一个新的读取器
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err // 创建读取器时出错
	}
	defer reader.Close() // 使用完后关闭读取器

	// 从对象池获取缓冲区以减少分配
	buf := zlibReadBuffer.Get().(*bytes.Buffer)
	buf.Reset()
	defer zlibReadBuffer.Put(buf)

	if _, err := io.Copy(buf, reader); err != nil {
		return nil, err // 复制数据时出错
	}

	// 创建副本以避免对象池重用时的数据污染
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil // 返回解压后的字节切片副本
}

// MultiZlibCompress 支持多次压缩
func MultiZlibCompress(data []byte, times int) ([]byte, error) {
	var err error
	compressedData := data

	for i := 0; i < times; i++ {
		compressedData, err = ZlibCompress(compressedData) // 进行多次压缩
		if err != nil {
			return nil, err // 压缩时出错
		}
	}
	return compressedData, nil // 返回最终的压缩数据
}

// MultiZlibDecompress 支持多次解压缩
func MultiZlibDecompress(compressedData []byte, times int) ([]byte, error) {
	var err error
	decompressedData := compressedData

	for i := 0; i < times; i++ {
		decompressedData, err = ZlibDecompress(decompressedData) // 进行多次解压缩
		if err != nil {
			return nil, err // 解压缩时出错
		}
	}
	return decompressedData, nil // 返回最终的解压数据
}

// ZlibCompressObject 泛型压缩函数，支持任意类型自动JSON序列化
func ZlibCompressObject[T any](obj T) ([]byte, error) {
	// 序列化对象为JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	// 压缩JSON数据
	return ZlibCompress(data)
}

// ZlibCompressObjectWithSize 泛型压缩函数，返回压缩后的数据和原始JSON数据大小
func ZlibCompressObjectWithSize[T any](obj T) ([]byte, int, error) {
	// 序列化对象为JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, 0, err
	}
	// 压缩JSON数据
	compressedData, err := ZlibCompress(data)
	if err != nil {
		return nil, 0, err
	}
	return compressedData, len(data), nil
}

// ZlibDecompressObject 泛型解压缩函数，支持自动JSON反序列化
func ZlibDecompressObject[T any](compressedData []byte) (T, error) {
	var result T

	// 解压缩数据
	data, err := ZlibDecompress(compressedData)
	if err != nil {
		return result, err
	}

	// 反序列化JSON
	err = json.Unmarshal(data, &result)
	return result, err
}

// MultiZlibCompressObject 泛型多次压缩函数，支持任意类型自动JSON序列化
func MultiZlibCompressObject[T any](obj T, times int) ([]byte, error) {
	// 序列化对象为JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	// 多次压缩JSON数据
	return MultiZlibCompress(data, times)
}

// MultiZlibDecompressObject 泛型多次解压缩函数，支持自动JSON反序列化
func MultiZlibDecompressObject[T any](compressedData []byte, times int) (T, error) {
	var result T

	// 多次解压缩数据
	data, err := MultiZlibDecompress(compressedData, times)
	if err != nil {
		return result, err
	}

	// 反序列化JSON
	err = json.Unmarshal(data, &result)
	return result, err
}

// ZlibCompressWithPrefix 压缩数据并添加 ZLIB: 前缀
// 返回带前缀的压缩数据，适用于需要明确标识压缩格式的场景
func ZlibCompressWithPrefix(data []byte) ([]byte, error) {
	compressed, err := ZlibCompress(data)
	if err != nil {
		return nil, err
	}
	result := make([]byte, ZlibPrefixLen+len(compressed))
	copy(result, []byte(ZlibPrefix))
	copy(result[ZlibPrefixLen:], compressed)
	return result, nil
}

// ZlibDecompressWithPrefix 解压缩带 ZLIB: 前缀的数据
// 如果数据带有前缀，自动去除后解压；否则直接返回原数据
func ZlibDecompressWithPrefix(data []byte) ([]byte, error) {
	if len(data) > ZlibPrefixLen && string(data[:ZlibPrefixLen]) == ZlibPrefix {
		return ZlibDecompress(data[ZlibPrefixLen:])
	}
	// 如果没有前缀，直接返回原数据（假设未压缩）
	return data, nil
}

// IsZlibCompressed 检查数据是否带有 ZLIB 压缩前缀
func IsZlibCompressed(data []byte) bool {
	return len(data) > ZlibPrefixLen && string(data[:ZlibPrefixLen]) == ZlibPrefix
}
