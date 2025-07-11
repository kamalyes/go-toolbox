/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-28 09:52:21
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 10:53:21
 * @FilePath: \go-toolbox\pkg\osx\file.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kamalyes/go-toolbox/pkg/sign"
)

// GetDirFiles 获取指定目录及其子目录下的所有文件路径
// 参数：dir - 目录路径
// 返回：所有文件的路径切片和可能的错误
func GetDirFiles(dir string) ([]string, error) {
	// 读取目录内容
	dirList, err := os.ReadDir(dir)
	if err != nil {
		return nil, err // 返回错误
	}

	// 创建文件路径切片
	var filesRet []string

	// 遍历目录中的每个文件
	for _, file := range dirList {
		// 如果是目录，则递归调用 GetDirFiles
		if file.IsDir() {
			files, err := GetDirFiles(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err // 返回错误
			}

			filesRet = append(filesRet, files...) // 添加子目录中的文件
		} else {
			// 添加文件的完整路径
			filesRet = append(filesRet, filepath.Join(dir, file.Name()))
		}
	}

	return filesRet, nil // 返回文件路径切片
}

// CheckImageExists 检查图像文件是否存在并且有效
func CheckImageExists(filename string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("expected file %s to exist, but it does not", filename)
	}

	// 尝试打开文件并解码以验证其有效性
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open image file %s: %v", filename, err)
	}
	defer file.Close()

	_, _, err = image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image file %s: %v", filename, err)
	}

	return nil
}

// WriteContentToFile 将内容追加写入指定文件
func WriteContentToFile(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// SaveImage 将字节数据保存为指定文件名的图片
func SaveImage(filename string, imgData []byte, quality int) error {
	// 创建一个新的文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取字节数据并解码为图片
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return err
	}

	// 根据文件扩展名选择编码方式
	ext := strings.ToLower(getFileExtension(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	case ".png":
		return png.Encode(file, img)
	case ".bmp":
		return encodeBMP(file, img)
	default:
		return fmt.Errorf("unsupported file format: %s", filename)
	}
}

// encodeBMP 将图像编码为 BMP 格式
func encodeBMP(file *os.File, img image.Image) error {
	// BMP 编码需要使用第三方库，或手动实现
	return fmt.Errorf("BMP encoding is not implemented")
}

// getFileExtension 返回文件的扩展名
func getFileExtension(filename string) string {
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		return filename[idx:]
	}
	return ""
}

// CreateIfNotExist 如果文件不存在则创建它
func CreateIfNotExist(file string) (*os.File, error) {
	if FileExists(file) {
		return nil, fmt.Errorf("%s already exists", file)
	}
	return os.Create(file)
}

// RemoveIfExist 如果文件存在则删除它
func RemoveIfExist(filename string) error {
	if FileExists(filename) {
		return os.Remove(filename)
	}
	return nil
}

// FileExists 检查文件是否存在
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// FileNameWithoutExt 返回不带扩展名的文件名
func FileNameWithoutExt(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}

// ComputeHashes 计算文件的多种哈希值，优化为只读一次文件
func ComputeHashes(filePath string) (map[sign.HashCryptoFunc]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建哈希实例
	hashes := map[sign.HashCryptoFunc]hash.Hash{}
	for algo, constructor := range sign.SupportHMACCryptoFunc {
		hashes[algo] = constructor()
	}

	// 构造 MultiWriter 写入多个哈希器
	writers := make([]io.Writer, 0, len(hashes))
	for _, h := range hashes {
		writers = append(writers, h)
	}
	multiWriter := io.MultiWriter(writers...)

	// 读取文件一次，写入所有哈希器
	if _, err := io.Copy(multiWriter, file); err != nil {
		return nil, err
	}

	results := make(map[sign.HashCryptoFunc]string, len(hashes))
	for algo, h := range hashes {
		results[algo] = hex.EncodeToString(h.Sum(nil))
	}

	return results, nil
}
