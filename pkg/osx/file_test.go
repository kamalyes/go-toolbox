/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-12 22:17:11
 * @FilePath: \go-toolbox\pkg\osx\file_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/sign"
	"github.com/stretchr/testify/assert"
)

// TestGetDirFiles 测试 GetDirFiles 函数
func TestGetDirFiles(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir := t.TempDir()
	testFile1 := filepath.Join(tempDir, "file1.txt")
	testFile2 := filepath.Join(tempDir, "file2.txt")
	os.WriteFile(testFile1, []byte("test content 1"), 0644)
	os.WriteFile(testFile2, []byte("test content 2"), 0644)

	// 创建一个子目录
	subDir := filepath.Join(tempDir, "subdir")
	os.Mkdir(subDir, 0755)
	subTestFile := filepath.Join(subDir, "file3.txt")
	os.WriteFile(subTestFile, []byte("test content 3"), 0644)

	// 获取目录中的文件
	files, err := GetDirFiles(tempDir)
	assert.NoError(t, err)

	// 检查返回的文件数量
	assert.Equal(t, 3, len(files))
}

// createTestImage 创建一个简单的测试图像并保存到指定文件
func createTestImage(filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// 填充图像为红色
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// TestCheckImageExists 测试 CheckImageExists 函数
func TestCheckImageExists(t *testing.T) {
	// 使用 map 组织测试用例
	testCases := map[string]struct {
		shouldExist bool
	}{
		"valid_image.png":        {true},
		"non_existent_image.png": {false},
		"invalid_image.png":      {false},
	}

	// 创建有效图像
	validFilename := "valid_image.png"
	err := createTestImage(validFilename)
	assert.NoError(t, err)
	t.Cleanup(func() { os.Remove(validFilename) }) // 清理测试文件

	// 创建无效图像
	invalidImageFilename := "invalid_image.png"
	file, err := os.Create(invalidImageFilename)
	assert.NoError(t, err)
	file.WriteString("This is not an image.")
	file.Close()
	t.Cleanup(func() { os.Remove(invalidImageFilename) }) // 清理测试文件

	for filename, tc := range testCases {
		t.Run(filename, func(t *testing.T) {
			exists := CheckImageExists(filename) == nil
			assert.Equal(t, tc.shouldExist, exists, fmt.Sprintf("Expected existence of %s to be %v", filename, tc.shouldExist))
		})
	}
}

// TestSaveImage 测试 SaveImage 函数
func TestSaveImage(t *testing.T) {
	// 创建一个临时图像文件
	tempImageFilename := "temp_image.png"
	err := createTestImage(tempImageFilename)
	assert.NoError(t, err)
	t.Cleanup(func() { os.Remove(tempImageFilename) }) // 清理临时图像文件

	// 使用 map 组织测试用例
	testCases := map[string]struct {
		format    string
		quality   int
		expectErr bool
	}{
		"test_output.jpg": {"jpeg", 80, false},
		"test_output.png": {"png", 0, false},
		"test_output.bmp": {"bmp", 0, true}, // 预期失败
	}

	for filename, tc := range testCases {
		t.Run(filename, func(t *testing.T) {
			imgData, err := os.ReadFile(tempImageFilename) // 读取临时图像文件数据
			assert.NoError(t, err)

			err = SaveImage(filename, imgData, tc.quality)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NoError(t, CheckImageExists(filename))
			}
			t.Cleanup(func() { os.Remove(filename) }) // 清理测试文件
		})
	}
}

func TestWriteContentToFile(t *testing.T) {
	// 定义测试用的文件路径和内容
	testFilePath := "test_output.txt"
	testContent := "Hello, World!"

	// 调用 WriteContentToFile 函数
	err := WriteContentToFile(testFilePath, testContent)

	// 断言没有错误发生
	assert.NoError(t, err)

	// 读取文件内容以验证写入是否成功
	content, err := os.ReadFile(testFilePath)
	assert.NoError(t, err)

	// 断言内容与预期一致
	assert.Equal(t, testContent, string(content))

	// 清理测试生成的文件
	err = os.Remove(testFilePath)
	assert.NoError(t, err)
}

// TestFileNameWithoutExt 测试 FileNameWithoutExt 函数
func TestFileNameWithoutExt(t *testing.T) {
	// 测试文件名和期望结果
	testCases := []struct {
		input    string
		expected string
	}{
		{"example.txt", "example"},
		{"archive.tar.gz", "archive.tar"},
		{"noext", "noext"},
		{".hiddenfile", ""}, // 注意：这个行为可能根据需求有所不同
	}

	for _, testCase := range testCases {
		result := FileNameWithoutExt(testCase.input)
		assert.Equal(t, testCase.expected, result, fmt.Sprintf("FileNameWithoutExt(%q) 应该返回 %q", testCase.input, testCase.expected))
	}
}

// TestRemoveIfExist 测试 RemoveIfExist 函数
func TestRemoveIfExist(t *testing.T) {
	// 创建一个临时文件用于测试
	tempFile := filepath.Join(MkdirTemp(), "tempfile")
	defer os.Remove(tempFile) // 清理：如果测试失败，将在测试完成后删除文件
	// 但注意，如果 RemoveIfExist 成功，这个文件将在测试中被删除

	// 创建文件
	err := os.WriteFile(tempFile, []byte("test"), 0644)
	assert.NoError(t, err, "ioutil.WriteFile 应该成功创建文件")

	// 调用 RemoveIfExist
	err = RemoveIfExist(tempFile)
	assert.NoError(t, err, "RemoveIfExist 应该成功删除文件")

	// 检查文件是否已删除
	assert.False(t, FileExists(tempFile), "文件应该已被删除")
}

// TestCreateIfNotExist 测试 CreateIfNotExist 函数
func TestCreateIfNotExist(t *testing.T) {
	// 创建一个临时文件用于测试
	tempFile := filepath.Join(MkdirTemp(), "tempfile")
	defer os.Remove(tempFile) // 测试完成后删除临时文件

	// 调用 CreateIfNotExist
	file, err := CreateIfNotExist(tempFile)
	assert.NoError(t, err, "CreateIfNotExist 应该成功创建文件")
	assert.NotNil(t, file, "返回的文件句柄不应该为空")

	// 检查文件是否存在
	assert.True(t, FileExists(tempFile), "文件应该存在")

	// 关闭文件句柄
	file.Close()

	// 尝试再次创建同一个文件，应该返回错误
	_, err = CreateIfNotExist(tempFile)
	assert.Error(t, err, "CreateIfNotExist 应该返回错误，因为文件已经存在")
}

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
	expected := map[sign.HashCryptoFunc]string{
		sign.AlgorithmMD5:    "65a8e27d8879283831b664bd8b7f0ad4",
		sign.AlgorithmSHA1:   "0a0a9f2a6772942557ab5355d76af442f8f65e01",
		sign.AlgorithmSHA224: "72a23dfa411ba6fde01dbfabf3b00a709c93ebf273dc29e2d8b261ff",
		sign.AlgorithmSHA256: "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f",
		sign.AlgorithmSHA384: "5485cc9b3365b4305dfb4e8337e0a598a574f8242bf17289e0dd6c20a3cd44a089de16ab4ab308f63e44b1170eb5f515",
		sign.AlgorithmSHA512: "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387",
	}

	// 遍历期望值逐个断言
	for algo, expectedHash := range expected {
		actualHash, ok := hashMap[algo]
		assert.True(t, ok, fmt.Sprintf("应该包含算法 %s 的哈希结果", algo))
		assert.Equal(t, expectedHash, actualHash, fmt.Sprintf("算法 %s 哈希不匹配", algo))
	}
}

// TestFindFiles 测试 FindFiles 函数
func TestFindFiles(t *testing.T) {
	// 创建临时目录结构用于测试
	tempDir := t.TempDir()

	// 创建目录结构
	//  tempDir/
	//  ├── file1.pb.go
	//  ├── file2.pb.go
	//  ├── other.txt
	//  └── subdir/
	//      ├── file3.pb.go
	//      └── nested/
	//          └── file5.pb.go

	pbFile1 := filepath.Join(tempDir, "file1.pb.go")
	pbFile2 := filepath.Join(tempDir, "file2.pb.go")
	otherFile := filepath.Join(tempDir, "other.txt")
	subDir := filepath.Join(tempDir, "subdir")
	pbFile3 := filepath.Join(subDir, "file3.pb.go")
	nestedDir := filepath.Join(subDir, "nested")
	pbFile4 := filepath.Join(nestedDir, "file5.pb.go")

	// 创建文件
	for _, f := range []string{pbFile1, pbFile2, otherFile} {
		err := os.WriteFile(f, []byte("test"), 0644)
		assert.NoError(t, err, "Failed to create test file")
	}

	// 创建子目录和文件
	err := os.MkdirAll(nestedDir, 0755)
	assert.NoError(t, err, "Failed to create test directories")
	for _, f := range []string{pbFile3, pbFile4} {
		err := os.WriteFile(f, []byte("test"), 0644)
		assert.NoError(t, err, "Failed to create test file")
	}

	testCases := []struct {
		name      string
		pattern   string
		expectLen int
		desc      string
	}{
		{
			name:      "simple_glob",
			pattern:   filepath.Join(tempDir, "*.pb.go"),
			expectLen: 2,
			desc:      "Should find only .pb.go files in root directory",
		},
		{
			name:      "recursive_glob",
			pattern:   filepath.Join(tempDir, "**", "*.pb.go"),
			expectLen: 4,
			desc:      "Should find all .pb.go files recursively",
		},
		{
			name:      "non_matching_pattern",
			pattern:   filepath.Join(tempDir, "*.xyz"),
			expectLen: 0,
			desc:      "Should return empty slice for non-matching pattern",
		},
		{
			name:      "txt_files",
			pattern:   filepath.Join(tempDir, "*.txt"),
			expectLen: 1,
			desc:      "Should find .txt files",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			files, err := FindFiles(tc.pattern)
			assert.NoError(t, err, "FindFiles should not return error")
			assert.Equal(t, tc.expectLen, len(files), tc.desc)

			// 验证返回的文件都存在
			for _, f := range files {
				assert.True(t, FileExists(f), "Returned file should exist: "+f)
			}
		})
	}
}

// TestFindFilesRecursive 测试 findFilesRecursive 内部函数（间接测试）
func TestFindFilesRecursive(t *testing.T) {
	tempDir := t.TempDir()

	// 创建深层目录结构
	//  tempDir/
	//  ├── level1/
	//  │   ├── file.txt
	//  │   └── level2/
	//  │       ├── file.txt
	//  │       └── level3/
	//  │           └── file.txt

	level1 := filepath.Join(tempDir, "level1")
	level2 := filepath.Join(level1, "level2")
	level3 := filepath.Join(level2, "level3")

	err := os.MkdirAll(level3, 0755)
	assert.NoError(t, err, "Failed to create test directories")

	testFiles := []string{
		filepath.Join(level1, "file.txt"),
		filepath.Join(level2, "file.txt"),
		filepath.Join(level3, "file.txt"),
	}

	for _, f := range testFiles {
		err := os.WriteFile(f, []byte("test"), 0644)
		assert.NoError(t, err, "Failed to create test file")
	}

	// 使用递归模式查找所有 file.txt
	pattern := filepath.Join(level1, "**", "file.txt")
	files, err := FindFiles(pattern)

	assert.NoError(t, err, "FindFiles should not return error for recursive pattern")
	assert.Equal(t, 3, len(files), "Should find 3 files in nested directories")

	// 验证找到的文件都在正确的位置
	for _, f := range files {
		assert.True(t, FileExists(f), "Returned file should exist: "+f)
		assert.Contains(t, f, "file.txt", "Returned file should be file.txt")
	}
}

// TestFindFilesEdgeCases 测试 FindFiles 的边界情况
func TestFindFilesEdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		pattern   string
		shouldErr bool
		desc      string
	}{
		{
			name:      "empty_pattern",
			pattern:   "",
			shouldErr: false,
			desc:      "Empty pattern should not error (glob behavior)",
		},
		{
			name:      "invalid_recursive_pattern",
			pattern:   "path/with/**/multiple/**/stars/file.txt",
			shouldErr: true,
			desc:      "Pattern with multiple ** should error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FindFiles(tc.pattern)
			if tc.shouldErr {
				assert.Error(t, err, tc.desc)
			} else {
				// 不应该 panic，可能返回错误或空结果
				// 主要验证函数不会崩溃
			}
		})
	}
}
