/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-24 11:25:16
 * @FilePath: \go-toolbox\zipx\zlib_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package zipx

import (
	"reflect"
	"testing"
)


func TestAllZlibFunctions(t *testing.T) {
	t.Run("TestZlibCompress", TestZlibCompress)
	t.Run("TestZlibDecompress", TestZlibDecompress)

}

// TestZlibCompress 测试 ZlibCompress 函数
func TestZlibCompress(t *testing.T) {
	data := []byte("Test data for zlib compression")

	compressedData, err := ZlibCompress(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(compressedData) == 0 {
		t.Error("Expected compressed data to be non-empty")
	}
}

// TestZlibDecompress 测试 ZlibDecompress 函数
func TestZlibDecompress(t *testing.T) {
	data := []byte("Test data for zlib compression")

	compressedData, err := ZlibCompress(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decompressedData, err := ZlibDecompress(compressedData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(decompressedData, data) {
		t.Errorf("Expected decompressed data to be %s, got %s", data, decompressedData)
	}
}
