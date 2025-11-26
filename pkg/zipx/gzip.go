/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-08 12:05:55
 * @FilePath: \go-toolbox\pkg\zipx\gzip.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package zipx

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
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
