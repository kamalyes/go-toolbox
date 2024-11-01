/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-24 11:25:16
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
)

// GzipCompress 压缩任意数据并返回压缩后的数据
func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	// 将数据写入 gzip writer
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GzipDecompress 解压缩任意数据并返回解压缩后的数据和压缩数据的大小
func GzipDecompress(compressedData []byte) ([]byte, int, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, 0, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return nil, 0, err
	}

	return buf.Bytes(), len(compressedData), nil
}
