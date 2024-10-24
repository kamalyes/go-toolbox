/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-24 11:25:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-24 11:25:16
 * @FilePath: \go-toolbox\zipx\gzip_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package zipx

import (
	"reflect"
	"testing"
	"time"

	"github.com/kamalyes/go-toolbox/json"
	"github.com/kamalyes/go-toolbox/random"
)

func TestAllGzipFunctions(t *testing.T) {
	t.Run("TestGzipCompressJSON", TestGzipCompressJSON)
	t.Run("TestGzipDecompressJSON", TestGzipDecompressJSON)
	t.Run("TestCompressDecompress", TestCompressDecompress)
}

// 定义测试模型
type TestModel struct {
	Name      string         `json:"name"`
	Age       int            `json:"age"`
	Salary    float64        `json:"salary"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	Tags      []string       `json:"tags"`
	Settings  map[string]int `json:"settings"`
}

// helper function to compress a TestModel instance
func compressModel(model TestModel) ([]byte, error) {
	modelJSON, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	return GzipCompress(modelJSON)
}

// helper function to decompress data into a TestModel instance
func decompressModel(compressedData []byte) (TestModel, error) {
	var decompressedModel TestModel
	decompressedJSON, _, err := GzipDecompress(compressedData)
	if err != nil {
		return decompressedModel, err
	}
	err = json.Unmarshal(decompressedJSON, &decompressedModel)
	return decompressedModel, err
}

// TestGzipCompressJSON 测试 GzipCompressJSON 函数
func TestGzipCompressJSON(t *testing.T) {
	modelJSON, err := random.GenerateRandomModel(&TestModel{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var model TestModel
	err = json.Unmarshal([]byte(modelJSON), &model)
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
		modelJSON, err := random.GenerateRandomModel(&TestModel{})
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}

		var model TestModel
		err = json.Unmarshal([]byte(modelJSON), &model)
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
	modelJSON, err := random.GenerateRandomModel(&TestModel{})
	if err != nil {
		b.Fatalf("Expected no error, got %v", err)
	}

	var model TestModel
	err = json.Unmarshal([]byte(modelJSON), &model)
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
