/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 00:15:15
 * @FilePath: \go-toolbox\tests\osx_file_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
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
	files, err := osx.GetDirFiles(tempDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 检查返回的文件数量
	if len(files) != 3 {
		t.Fatalf("Expected 3 files, got %d", len(files))
	}
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
	if err := createTestImage(validFilename); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	t.Cleanup(func() { os.Remove(validFilename) }) // 清理测试文件

	// 创建无效图像
	invalidImageFilename := "invalid_image.png"
	file, err := os.Create(invalidImageFilename)
	if err != nil {
		t.Fatalf("Failed to create invalid image file: %v", err)
	}
	file.WriteString("This is not an image.")
	file.Close()
	t.Cleanup(func() { os.Remove(invalidImageFilename) }) // 清理测试文件

	for filename, tc := range testCases {
		t.Run(filename, func(t *testing.T) {
			exists := osx.CheckImageExists(filename) == nil
			if exists != tc.shouldExist {
				t.Errorf("Expected existence of %s to be %v, but got %v", filename, tc.shouldExist, exists)
			}
		})
	}
}

// TestSaveImage 测试 SaveImage 函数
func TestSaveImage(t *testing.T) {
	// 创建一个临时图像文件
	tempImageFilename := "temp_image.png"
	if err := createTestImage(tempImageFilename); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
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
			if err != nil {
				t.Fatalf("Failed to read temp image data: %v", err)
			}

			err = osx.SaveImage(filename, imgData, tc.quality)
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error status %v for %s, but got %v", tc.expectErr, filename, err != nil)
			}
			t.Cleanup(func() { os.Remove(filename) }) // 清理测试文件
			if err == nil && osx.CheckImageExists(filename) != nil {
				t.Errorf("Expected %s to exist after saving, but it does not", filename)
			}
		})
	}
}

func TestWriteContentToFile(t *testing.T) {
	// 定义测试用的文件路径和内容
	testFilePath := "test_output.txt"
	testContent := "Hello, World!"

	// 调用 WriteContentToFile 函数
	err := osx.WriteContentToFile(testFilePath, testContent)

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
		result := osx.FileNameWithoutExt(testCase.input)
		assert.Equal(t, testCase.expected, result, fmt.Sprintf("FileNameWithoutExt(%q) 应该返回 %q", testCase.input, testCase.expected))
	}
}

// TestRemoveIfExist 测试 RemoveIfExist 函数
func TestRemoveIfExist(t *testing.T) {
	// 创建一个临时文件用于测试
	tempFile := filepath.Join(osx.MkdirTemp(), "tempfile")
	defer os.Remove(tempFile) // 清理：如果测试失败，将在测试完成后删除文件
	// 但注意，如果 RemoveIfExist 成功，这个文件将在测试中被删除

	// 创建文件
	err := os.WriteFile(tempFile, []byte("test"), 0644)
	assert.NoError(t, err, "ioutil.WriteFile 应该成功创建文件")

	// 调用 RemoveIfExist
	err = osx.RemoveIfExist(tempFile)
	assert.NoError(t, err, "RemoveIfExist 应该成功删除文件")

	// 检查文件是否已删除
	assert.False(t, osx.FileExists(tempFile), "文件应该已被删除")
}

// TestCreateIfNotExist 测试 CreateIfNotExist 函数
func TestCreateIfNotExist(t *testing.T) {
	// 创建一个临时文件用于测试
	tempFile := filepath.Join(osx.MkdirTemp(), "tempfile")
	defer os.Remove(tempFile) // 测试完成后删除临时文件

	// 调用 CreateIfNotExist
	file, err := osx.CreateIfNotExist(tempFile)
	assert.NoError(t, err, "CreateIfNotExist 应该成功创建文件")
	assert.NotNil(t, file, "返回的文件句柄不应该为空")

	// 检查文件是否存在
	assert.True(t, osx.FileExists(tempFile), "文件应该存在")

	// 关闭文件句柄
	file.Close()

	// 尝试再次创建同一个文件，应该返回错误
	_, err = osx.CreateIfNotExist(tempFile)
	assert.Error(t, err, "CreateIfNotExist 应该返回错误，因为文件已经存在")
}

func TestHash(t *testing.T) {
	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "tempFile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// 写入内容到临时文件
	_, err = tempFile.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	// 计算文件的MD5哈希值
	hash, err := osx.ComputeHashes(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// 验证哈希值是否正确
	expectedMd5Hash := "65a8e27d8879283831b664bd8b7f0ad4"
	expectedSha1Hash := "0a0a9f2a6772942557ab5355d76af442f8f65e01"
	expectedSha224Hash := "72a23dfa411ba6fde01dbfabf3b00a709c93ebf273dc29e2d8b261ff"
	expectedSha256Hash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
	expectedSha384Hash := "5485cc9b3365b4305dfb4e8337e0a598a574f8242bf17289e0dd6c20a3cd44a089de16ab4ab308f63e44b1170eb5f515"
	expectedSha512Hash := "374d794a95cdcfd8b35993185fef9ba368f160d8daf432d08ba9f1ed1e5abe6cc69291e0fa2fe0006a52570ef18c19def4e617c33ce52ef0a6e5fbe318cb0387"
	assert.Equal(t, expectedMd5Hash, hash.MD5, fmt.Sprintf("Expected md5 %s, got %s", expectedMd5Hash, hash.MD5))
	assert.Equal(t, expectedSha1Hash, hash.SHA1, fmt.Sprintf("Expected sha1 %s, got %s", expectedSha1Hash, hash.SHA1))
	assert.Equal(t, expectedSha224Hash, hash.SHA224, fmt.Sprintf("Expected sha224 %s, got %s", expectedSha224Hash, hash.SHA224))
	assert.Equal(t, expectedSha256Hash, hash.SHA256, fmt.Sprintf("Expected sha256 %s, got %s", expectedSha256Hash, hash.SHA256))
	assert.Equal(t, expectedSha384Hash, hash.SHA384, fmt.Sprintf("Expected sha256 %s, got %s", expectedSha384Hash, hash.SHA384))
	assert.Equal(t, expectedSha512Hash, hash.SHA512, fmt.Sprintf("Expected sha512 %s, got %s", expectedSha512Hash, hash.SHA512))
}
