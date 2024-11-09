/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 00:15:15
 * @FilePath: \go-toolbox\tests\osx_path_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/stretchr/testify/assert"
)

// TestMkdirIfNotExist 测试 MkdirIfNotExist 函数
func TestMkdirIfNotExist(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir := osx.MkdirTemp()
	defer os.RemoveAll(tempDir) // 测试完成后删除临时目录

	// 创建一个子目录路径
	subDir := filepath.Join(tempDir, "subdir")

	// 调用 MkdirIfNotExist
	err := osx.MkdirIfNotExist(subDir)
	assert.NoError(t, err, "MkdirIfNotExist 应该成功创建目录")

	// 检查目录是否存在
	info, err := os.Stat(subDir)
	assert.NoError(t, err, "os.Stat 应该成功获取目录信息")
	assert.True(t, info.IsDir(), "创建的应该是一个目录")
}

func TestDirHasContent(t *testing.T) {
	// 创建一个临时目录
	tempDir := osx.MkdirTemp()
	defer os.RemoveAll(tempDir)

	// 创建一个空文件
	emptyFile := filepath.Join(tempDir, "empty.txt")
	os.WriteFile(emptyFile, []byte(""), 0644)

	// 创建一个非空文件
	nonEmptyFile := filepath.Join(tempDir, "nonempty.txt")
	os.WriteFile(nonEmptyFile, []byte("Hello, World!"), 0644)

	// 测试空目录（有空文件）
	fileExists, files, _ := osx.DirHasContent(tempDir)
	assert.Equal(t, false, !fileExists, fmt.Sprintf("Expected directory no non-empty files :%#v", files))

	// 删除空文件，添加非空文件
	os.Remove(emptyFile)

	// 测试非空目录
	fileExists2, files2, _ := osx.DirHasContent(tempDir)
	assert.Equal(t, false, !fileExists2, fmt.Sprintf("Expected directory to have no non-empty files :%#v", files2))
}

func TestCopy(t *testing.T) {
	// 创建一个临时文件作为源文件
	srcFile, err := ioutil.TempFile("", "srcFile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())
	srcFile.Write([]byte("Hello, World!"))
	srcFile.Close()

	// 创建一个临时目录作为目标目录
	tempDir := osx.MkdirTemp()
	defer os.RemoveAll(tempDir)
	destFile := filepath.Join(tempDir, "destFile.txt")

	// 执行复制操作
	err = osx.Copy(srcFile.Name(), destFile)
	if err != nil {
		t.Fatal(err)
	}

	// 验证目标文件内容
	destContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(destContent) != "Hello, World!" {
		t.Error("File content does not match")
	}
}
