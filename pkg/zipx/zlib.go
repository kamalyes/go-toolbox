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

var (
	zlibBuffer = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	zlibReadBuffer = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	// zlibWriterPool 复用 zlib.Writer，避免每次压缩都分配新对象
	zlibWriterPool = sync.Pool{
		New: func() interface{} {
			return zlib.NewWriter(nil)
		},
	}
	// zlibReaderPool 复用 zlib 解压 reader（通过 zlib.Resetter 接口），避免每次解压都分配新对象
	// zlib.Reader 未导出，但 NewReader 返回的 io.ReadCloser 实现了 zlib.Resetter 接口
	zlibReaderPool = sync.Pool{
		New: func() interface{} {
			var b bytes.Buffer
			w := zlib.NewWriter(&b)
			w.Close()
			reader, err := zlib.NewReader(bytes.NewReader(b.Bytes()))
			if err != nil {
				panic("zipx: failed to initialize zlib reader pool: " + err.Error())
			}
			return reader
		},
	}
)

// ZlibCompress 压缩数据
func ZlibCompress(data []byte) ([]byte, error) {
	buf := zlibBuffer.Get().(*bytes.Buffer)
	defer zlibBuffer.Put(buf)
	buf.Reset()

	writer := zlibWriterPool.Get().(*zlib.Writer)
	writer.Reset(buf) // zlib.Writer 支持 Reset，可安全复用
	closed := false
	defer func() {
		if !closed {
			writer.Close()
		}
		zlibWriterPool.Put(writer)
	}()

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}
	closed = true

	// 返回副本，避免 buf 被放回 Pool 后产生数据竞争
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// ZlibCompressWithInfo 压缩数据并返回压缩信息
func ZlibCompressWithInfo(data []byte) (*CompressResult, error) {
	compressed, err := ZlibCompress(data)
	if err != nil {
		return nil, err
	}
	return newCompressResult(data, compressed), nil
}

// ZlibDecompress 解压缩数据（使用对象池优化）
func ZlibDecompress(compressedData []byte) ([]byte, error) {
	reader := zlibReaderPool.Get().(io.ReadCloser)
	defer zlibReaderPool.Put(reader)

	if err := reader.(zlib.Resetter).Reset(bytes.NewReader(compressedData), nil); err != nil {
		return nil, err
	}
	defer reader.Close()

	buf := zlibReadBuffer.Get().(*bytes.Buffer)
	buf.Reset()
	defer zlibReadBuffer.Put(buf)

	if _, err := io.Copy(buf, reader); err != nil {
		return nil, err
	}

	// 返回副本，避免 buf 被放回 Pool 后产生数据竞争
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
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

// MultiZlibCompressWithInfo 支持多次压缩并返回压缩信息
func MultiZlibCompressWithInfo(data []byte, times int) (*CompressResult, error) {
	compressed, err := MultiZlibCompress(data, times)
	if err != nil {
		return nil, err
	}
	return newCompressResult(data, compressed), nil
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

// ZlibCompressObjectWithInfo 泛型压缩函数，支持任意类型自动JSON序列化并返回压缩信息
func ZlibCompressObjectWithInfo[T any](obj T) (*CompressResult, error) {
	compressed, originalSize, err := ZlibCompressObjectWithSize(obj)
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

// MultiZlibCompressObjectWithInfo 泛型多次压缩函数，支持任意类型自动JSON序列化并返回压缩信息
func MultiZlibCompressObjectWithInfo[T any](obj T, times int) (*CompressResult, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	compressed, err := MultiZlibCompress(data, times)
	if err != nil {
		return nil, err
	}
	return newCompressResult(data, compressed), nil
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

// ZlibCompressWithPrefixInfo 压缩数据并添加 ZLIB: 前缀，同时返回压缩信息
func ZlibCompressWithPrefixInfo(data []byte) (*CompressResult, error) {
	compressed, err := ZlibCompressWithPrefix(data)
	if err != nil {
		return nil, err
	}
	return newCompressResult(data, compressed), nil
}

// IsZlibCompressed 检查数据是否带有 ZLIB 压缩前缀
func IsZlibCompressed(data []byte) bool {
	return len(data) > ZlibPrefixLen && string(data[:ZlibPrefixLen]) == ZlibPrefix
}

// ZlibSmartDecompress 智能解压缩函数
// 自动检测数据是否被压缩，如果是则解压，否则直接返回原数据
// 适用于需要兼容压缩/未压缩数据的场景
func ZlibSmartDecompress(data []byte) ([]byte, error) {
	// 1. 检查是否有 ZLIB 前缀
	if IsZlibCompressed(data) {
		return ZlibDecompress(data[ZlibPrefixLen:])
	}

	// 2. 尝试直接解压缩（处理没有前缀但被压缩的数据）
	decompressed, err := ZlibDecompress(data)
	if err == nil {
		return decompressed, nil
	}

	// 3. 解压失败，返回原数据
	return data, nil
}

// ZlibSmartDecompressObject 智能解压缩对象函数
// 自动检测数据是否被压缩，支持自动JSON反序列化
// 适用于需要兼容压缩/未压缩数据的场景
func ZlibSmartDecompressObject[T any](data []byte) (T, error) {
	var result T

	// 智能解压缩
	decompressed, err := ZlibSmartDecompress(data)
	if err != nil {
		return result, err
	}

	// 反序列化
	if err := json.Unmarshal(decompressed, &result); err != nil {
		return result, err
	}

	return result, nil
}
