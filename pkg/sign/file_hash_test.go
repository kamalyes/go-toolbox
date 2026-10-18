/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-18 15:36:52
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-18 15:52:36
 * @FilePath: \go-toolbox\pkg\sign\file_hash_test.go
 * @Description: 文件多哈希计算测试（随 ComputeHashes 从 osx 迁入）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package sign

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "tempFile.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// 写入测试内容
	_, err = tempFile.Write([]byte("Hello, World!"))
	assert.NoError(t, err)
	tempFile.Close()

	// 计算哈希，返回map
	hashMap, err := ComputeHashes(tempFile.Name())
	assert.NoError(t, err)

	// 验证哈希值是否正确
	expected := map[HashCryptoFunc]string{
		AlgorithmMD5:    "65a8e27d8879283831b664bd8b7f0ad4",
		AlgorithmSHA1:   "0a0a9f2a6772942557ab5355d76af442f8f65e01",
		AlgorithmSHA224: "72a23dfa411ba6fde01dbfabf3b00a709c93ebf273dc29e2d8b261ff",
		AlgorithmSHA256: "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f",
		AlgorithmSHA384: "5485cc9b3365b4305dfb4e8337e0a598a574f8242bf17289e0dd6c20a3cd44a089de16ab4ab308f63e44b1170eb5f515",
		AlgorithmSHA512: "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387",
	}

	// 遍历期望值逐个断言
	for algo, expectedHash := range expected {
		actualHash, ok := hashMap[algo]
		assert.True(t, ok, fmt.Sprintf("应该包含算法 %s 的哈希结果", algo))
		assert.Equal(t, expectedHash, actualHash, fmt.Sprintf("算法 %s 哈希不匹配", algo))
	}
}
