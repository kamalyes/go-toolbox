/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-03 13:52:20
 * @FilePath: \go-toolbox\tests\zipx_gzip_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"reflect"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/json"
	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/go-toolbox/pkg/zipx"
)

func TestAllGzipFunctions(t *testing.T) {
	t.Run("TestGzipCompressJSON", TestGzipCompressJSON)
	t.Run("TestGzipDecompressJSON", TestGzipDecompressJSON)
	t.Run("TestCompressDecompress", TestCompressDecompress)
}

// helper function to compress a TestModel instance
func compressModel(model TestModel) ([]byte, error) {
	modelJSON, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	return zipx.GzipCompress(modelJSON)
}

// helper function to decompress data into a TestModel instance
func decompressModel(compressedData []byte) (TestModel, error) {
	var decompressedModel TestModel
	decompressedJSON, _, err := zipx.GzipDecompress(compressedData)
	if err != nil {
		return decompressedModel, err
	}
	err = json.Unmarshal(decompressedJSON, &decompressedModel)
	return decompressedModel, err
}

// TestGzipCompressJSON 测试 GzipCompressJSON 函数
func TestGzipCompressJSON(t *testing.T) {
	_, jsonData, err := random.GenerateRandomModel(&TestModel{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var model TestModel
	err = json.Unmarshal([]byte(jsonData), &model)
	if err != nil {
		t.Fatalf("Failed to unmarshal model JSON: %v", err)
	}

	compressedData, err := compressModel(model)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Logf("Compressed size: %d bytes", len(compressedData))

	if len(compressedData) == 0 {
		t.Error("Expected compressed data to be non-empty")
	}
}

// TestGzipDecompressJSON 测试 GzipDecompressJSON 函数
func TestGzipDecompressJSON(t *testing.T) {
	model := TestModel{Name: "Alice", Age: 30}

	compressedData, err := compressModel(model)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decompressedModel, err := decompressModel(compressedData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(decompressedModel, model) {
		t.Errorf("Expected decompressed model to be %+v, got %+v", model, decompressedModel)
	}
}

// TestCompressDecompress 测试压缩和解压缩的完整流程
func TestCompressDecompress(t *testing.T) {
	model := TestModel{Name: "Alice", Age: 30}

	compressedData, err := compressModel(model)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decompressedModel, err := decompressModel(compressedData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(decompressedModel, model) {
		t.Errorf("Expected decompressed model to be %+v, got %+v", model, decompressedModel)
	}
}

// BenchmarkGzipCompressJSON 性能测试 GzipCompressJSON 函数
func BenchmarkGzipCompressJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, jsonData, err := random.GenerateRandomModel(&TestModel{})
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}

		var model TestModel
		err = json.Unmarshal([]byte(jsonData), &model)
		if err != nil {
			b.Fatalf("Failed to unmarshal model JSON: %v", err)
		}

		_, err = compressModel(model)
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
	}
}

// BenchmarkGzipDecompressJSON 性能测试 GzipDecompressJSON 函数
func BenchmarkGzipDecompressJSON(b *testing.B) {
	_, jsonData, err := random.GenerateRandomModel(&TestModel{})
	if err != nil {
		b.Fatalf("Expected no error, got %v", err)
	}

	var model TestModel
	err = json.Unmarshal([]byte(jsonData), &model)
	if err != nil {
		b.Fatalf("Failed to unmarshal model JSON: %v", err)
	}

	compressedData, err := compressModel(model)
	if err != nil {
		b.Fatalf("Expected no error, got %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := decompressModel(compressedData)
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
	}
}
