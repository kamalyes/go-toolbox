/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-13 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 13:55:15
 * @FilePath: \go-toolbox\pkg\safe\nil_panic_detector_test.go
 * @Description: Nil Panic检测工具，帮助发现项目中可能的nil指针访问
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package safe

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilPanicDetector(t *testing.T) {
	// 创建一个临时目录
	tempDir := os.TempDir() + "/nil_panic_detector_test"
	err := os.MkdirAll(tempDir, os.ModePerm)
	assert.NoError(t, err, "创建临时目录失败")

	// 创建一个测试文件，包含可能导致 nil panic 的代码
	testFilePath := tempDir + "/test.go"
	testFileContent := `
package main

func main() {
	var ptr *int
	_ = *ptr // 这将导致 nil panic
	x := map[string]string{"key": "value"}
	_ = x["nonexistent"] // 这里没有 nil 检查
}
`
	err = os.WriteFile(testFilePath, []byte(testFileContent), 0644)
	assert.NoError(t, err, "写入测试文件失败")

	// 创建 NilPanicDetector 实例并扫描目录
	detector := NewNilPanicDetector()
	err = detector.ScanDirectory(tempDir)
	assert.NoError(t, err, "扫描目录时出错")

	// 获取检测到的问题
	issues := detector.GetIssues()

	// 断言检测到的问题数量
	assert.Greater(t, len(issues), 0, "应该检测到至少一个问题")

	// 断言问题的类型和描述
	for _, issue := range issues {
		switch issue.Type {
		case "PointerDereference":
			assert.Equal(t, "指针解引用没有nil检查,可能导致panic", issue.Description)
		case "MapAccess":
			assert.Equal(t, "Map访问需要检查ok值", issue.Description)
		}
	}

	// 清理临时目录
	err = os.RemoveAll(tempDir)
	assert.NoError(t, err, "清理临时目录失败")
}
