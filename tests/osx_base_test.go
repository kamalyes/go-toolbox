/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 02:31:11
 * @FilePath: \go-toolbox\tests\osx_base_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/stretchr/testify/assert"
)

func TestAllSysBaseFunctions(t *testing.T) {
	t.Run("TestSafeGetHostName", TestSafeGetHostName)
	t.Run("TestHashUnixMicroCipherText", TestHashUnixMicroCipherText)
}

func TestSafeGetHostName(t *testing.T) {
	actual := osx.SafeGetHostName()
	assert.NotEmpty(t, actual, "HostNames should match")
}

// TestHashUnixMicroCipherText 测试 HashUnixMicroCipherText 函数
func TestHashUnixMicroCipherText(t *testing.T) {
	hash1 := osx.HashUnixMicroCipherText()
	hash2 := osx.HashUnixMicroCipherText()

	// 验证生成的哈希值不为空
	if hash1 == "" {
		t.Error("HashUnixMicroCipherText 生成的哈希值为空")
	}

	// 验证生成的哈希值长度是否为32（MD5哈希值长度）
	if len(hash1) != 32 {
		t.Errorf("期望哈希值长度为32，但得到的长度为 %d", len(hash1))
	}

	// 由于时间戳和随机字符串的原因，连续两次调用的结果应该不同
	if hash1 == hash2 {
		t.Error("连续两次调用 HashUnixMicroCipherText 生成的哈希值相同，期望不同")
	}
}

// TestGetCurrentPath 测试 GetCurrentPath 函数
func TestGetCurrentPath(t *testing.T) {
	expectedDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get expected working directory: %v", err)
	}

	actualDir, err := osx.GetCurrentPath()
	if err != nil {
		t.Fatalf("GetCurrentPath() returned an error: %v", err)
	}

	// 比较实际路径和预期路径
	if actualDir != expectedDir {
		t.Errorf("Expected %s, but got %s", expectedDir, actualDir)
	}

	// 额外检查路径是否是绝对路径
	if !filepath.IsAbs(actualDir) {
		t.Errorf("Expected an absolute path, but got a relative path: %s", actualDir)
	}
}
