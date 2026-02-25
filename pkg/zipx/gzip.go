/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-04 18:37:14
 * @FilePath: \go-toolbox\pkg\zipx\gzip.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package zipx

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"sync"
)

const (
	// GzipPrefix 是用于标识 gzip 压缩数据的前缀
	GzipPrefix = "GZIP:"
	// GzipPrefixLen 是 gzip 前缀的长度
	GzipPrefixLen = len(GzipPrefix)
)

var (
	gzipWriter sync.Pool // gzip.Writer 的对象池
	gzipReader sync.Pool // gzip.Reader 的对象池
	gzipBuffer sync.Pool // bytes.Buffer 的对象池
)

func init() {
	// 初始化对象池
	gzipWriter = sync.Pool{New: func() interface{} {
		return gzip.NewWriter(nil) // 创建新的 gzip.Writer
	}}
	gzipReader = sync.Pool{New: func() interface{} {
		return new(gzip.Reader) // 创建新的 gzip.Reader
	}}
	gzipBuffer = sync.Pool{New: func() interface{} {
		return bytes.NewBuffer(nil) // 创建新的 bytes.Buffer
	}}
}

// GzipCompress 使用 gzip 压缩数据
func GzipCompress(data []byte) ([]byte, error) {
	buf := gzipBuffer.Get().(*bytes.Buffer) // 从对象池获取 bytes.Buffer
	buf.Reset()                             // 重置缓冲区以避免数据污染
	defer gzipBuffer.Put(buf)               // 使用完后将缓冲区放回池中

	writer := gzipWriter.Get().(*gzip.Writer) // 从对象池获取 gzip.Writer
	writer.Reset(buf)                         // 将 writer 绑定到缓冲区
	defer func() {
		writer.Close()         // 关闭 writer 以刷新剩余数据
		gzipWriter.Put(writer) // 使用完后将 writer 放回池中
	}()

	if _, err := writer.Write(data); err != nil {
		return nil, err // 写入数据时出错
	}

	if err := writer.Close(); err != nil {
		return nil, err // 关闭 writer 时出错
	}

	// 创建副本以避免对象池重用时的数据污染
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil // 返回压缩后的字节切片副本
}

// GzipCompressWithInfo 使用 gzip 压缩数据并返回压缩信息
func GzipCompressWithInfo(data []byte) (*CompressResult, error) {
	compressed, err := GzipCompress(data)
	if err != nil {
		return nil, err
	}
	return newCompressResult(data, compressed), nil
}

// GzipDecompress 解压缩 gzip 压缩的数据
func GzipDecompress(compressedData []byte) ([]byte, error) {
	buf := bytes.NewBuffer(compressedData) // 创建一个新的缓冲区，读取压缩数据

	reader := gzipReader.Get().(*gzip.Reader) // 从对象池获取 gzip.Reader
	defer gzipReader.Put(reader)              // 使用完后将 reader 放回池中

	if err := reader.Reset(buf); err != nil {
		return nil, err // 重置 reader 时出错
	}
	defer reader.Close() // 使用完后关闭 reader

	data, err := io.ReadAll(reader) // 读取解压后的数据
	if err != nil {
		return nil, err // 读取数据时出错
	}
	return data, nil // 返回解压后的字节切片
}

// MultiGZipCompress 支持多次压缩
func MultiGZipCompress(data []byte, times int) ([]byte, error) {
	var err error
	compressedData := data

	for i := 0; i < times; i++ {
		compressedData, err = GzipCompress(compressedData) // 进行多次压缩
		if err != nil {
			return nil, err // 压缩时出错
		}
	}
	return compressedData, nil // 返回最终的压缩数据
}

// MultiGZipCompressWithInfo 支持多次压缩并返回压缩信息
func MultiGZipCompressWithInfo(data []byte, times int) (*CompressResult, error) {
	compressed, err := MultiGZipCompress(data, times)
	if err != nil {
		return nil, err
	}
	return newCompressResult(data, compressed), nil
}

// MultiGZipDecompress 支持多次解压缩
func MultiGZipDecompress(compressedData []byte, times int) ([]byte, error) {
	var err error
	decompressedData := compressedData

	for i := 0; i < times; i++ {
		decompressedData, err = GzipDecompress(decompressedData) // 进行多次解压缩
		if err != nil {
			return nil, err // 解压缩时出错
		}
	}
	return decompressedData, nil // 返回最终的解压数据
}

// GzipCompressObject 泛型压缩函数，支持任意类型自动JSON序列化
func GzipCompressObject[T any](obj T) ([]byte, error) {
	// 序列化对象为JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	// 压缩JSON数据
	return GzipCompress(data)
}

// GzipCompressObjectWithInfo 泛型压缩函数，支持任意类型自动JSON序列化并返回压缩信息
func GzipCompressObjectWithInfo[T any](obj T) (*CompressResult, error) {
	compressed, originalSize, err := GzipCompressObjectWithSize(obj)
	if err != nil {
		return nil, err
	}
	return &CompressResult{
		Data:           compressed,
		OriginalSize:   originalSize,
		CompressedSize: len(compressed),
		Ratio:          float64(len(compressed)) / float64(originalSize),
	}, nil
}

// GzipCompressObjectWithSize 泛型压缩函数，返回压缩后的数据和原始JSON数据大小
func GzipCompressObjectWithSize[T any](obj T) ([]byte, int, error) {
	// 序列化对象为JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, 0, err
	}
	// 压缩JSON数据
	compressedData, err := GzipCompress(data)
	if err != nil {
		return nil, 0, err
	}
	return compressedData, len(data), nil
}

// GzipDecompressObject 泛型解压缩函数，支持自动JSON反序列化
func GzipDecompressObject[T any](compressedData []byte) (T, error) {
	var result T

	// 解压缩数据
	data, err := GzipDecompress(compressedData)
	if err != nil {
		return result, err
	}

	// 反序列化JSON
	err = json.Unmarshal(data, &result)
	return result, err
}

// MultiGZipCompressObject 泛型多次压缩函数，支持任意类型自动JSON序列化
func MultiGZipCompressObject[T any](obj T, times int) ([]byte, error) {
	// 序列化对象为JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	// 多次压缩JSON数据
	return MultiGZipCompress(data, times)
}

// MultiGZipCompressObjectWithInfo 泛型多次压缩函数，支持任意类型自动JSON序列化并返回压缩信息
func MultiGZipCompressObjectWithInfo[T any](obj T, times int) (*CompressResult, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	compressed, err := MultiGZipCompress(data, times)
	if err != nil {
		return nil, err
	}
	return newCompressResult(data, compressed), nil
}

// MultiGZipDecompressObject 泛型多次解压缩函数，支持自动JSON反序列化
func MultiGZipDecompressObject[T any](compressedData []byte, times int) (T, error) {
	var result T

	// 多次解压缩数据
	data, err := MultiGZipDecompress(compressedData, times)
	if err != nil {
		return result, err
	}

	// 反序列化JSON
	err = json.Unmarshal(data, &result)
	return result, err
}

// GzipCompressWithPrefix 压缩数据并添加 GZIP: 前缀
// 返回带前缀的压缩数据，适用于需要明确标识压缩格式的场景
func GzipCompressWithPrefix(data []byte) ([]byte, error) {
	compressed, err := GzipCompress(data)
	if err != nil {
		return nil, err
	}
	result := make([]byte, GzipPrefixLen+len(compressed))
	copy(result, []byte(GzipPrefix))
	copy(result[GzipPrefixLen:], compressed)
	return result, nil
}

// GzipCompressWithPrefixInfo 压缩数据并添加 GZIP: 前缀，同时返回压缩信息
func GzipCompressWithPrefixInfo(data []byte) (*CompressResult, error) {
	compressed, err := GzipCompressWithPrefix(data)
	if err != nil {
		return nil, err
	}
	return newCompressResult(data, compressed), nil
}

// IsGzipCompressed 检查数据是否带有 GZIP 压缩前缀
func IsGzipCompressed(data []byte) bool {
	return len(data) > GzipPrefixLen && string(data[:GzipPrefixLen]) == GzipPrefix
}

// GzipSmartDecompress 智能解压缩函数
// 自动检测数据是否被压缩，如果是则解压，否则直接返回原数据
// 适用于需要兼容压缩/未压缩数据的场景
func GzipSmartDecompress(data []byte) ([]byte, error) {
	// 1. 检查是否有 GZIP 前缀
	if IsGzipCompressed(data) {
		return GzipDecompress(data[GzipPrefixLen:])
	}

	// 2. 尝试直接解压缩（处理没有前缀但被压缩的数据）
	decompressed, err := GzipDecompress(data)
	if err == nil {
		return decompressed, nil
	}

	// 3. 解压失败，返回原数据
	return data, nil
}

// GzipSmartDecompressObject 智能解压缩对象函数
// 自动检测数据是否被压缩，支持自动JSON反序列化
// 适用于需要兼容压缩/未压缩数据的场景
func GzipSmartDecompressObject[T any](data []byte) (T, error) {
	var result T

	// 智能解压缩
	decompressed, err := GzipSmartDecompress(data)
	if err != nil {
		return result, err
	}

	// 反序列化
	if err := json.Unmarshal(decompressed, &result); err != nil {
		return result, err
	}

	return result, nil
}
