/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-28 09:52:21
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 22:19:15
 * @FilePath: \go-toolbox\tests\osx_file_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/osx"
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
