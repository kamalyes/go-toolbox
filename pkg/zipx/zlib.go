/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-29 10:31:27
 * @FilePath: \engine-im-service\go-toolbox\pkg\zipx\zlib.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package zipx

import (
	"bytes"
	"compress/zlib"
	"io"
	"sync"
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
