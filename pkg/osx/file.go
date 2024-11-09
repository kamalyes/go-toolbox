/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-28 09:52:21
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 00:29:06
 * @FilePath: \go-toolbox\pkg\osx\file.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// HashResult 是一个结构体，用于存储不同哈希算法的结果。
type HashResult struct {
	MD5    string
	SHA1   string
	SHA224 string
	SHA256 string
	SHA384 string
	SHA512 string
}

// hashFunc 是一个类型别名，表示计算哈希值的函数类型。
type hashFunc func(io.Reader) (string, error)

// ComputeHashes 计算并返回文件的所有哈希值。
func ComputeHashes(filePath string) (*HashResult, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建哈希函数映射
	hashers := map[string]hashFunc{
		"MD5":    md5Hash,
		"SHA1":   sha1Hash,
		"SHA224": sha224HashSimplified,
		"SHA384": sha384HashSimplified,
		"SHA256": sha256Hash,
		"SHA512": sha512Hash,
	}

	// 计算哈希值
	hashes := &HashResult{}
	for name, harsherFunc := range hashers {
		// 重置文件指针到文件开始
		_, err := file.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		hash, err := harsherFunc(file)
		if err != nil {
			return nil, err
		}
		// 使用字段名来设置哈希值
		switch name {
		case "MD5":
			hashes.MD5 = hash
		case "SHA1":
			hashes.SHA1 = hash
		case "SHA224":
			hashes.SHA224 = hash
		case "SHA384":
			hashes.SHA384 = hash
		case "SHA256":
			hashes.SHA256 = hash
		case "SHA512":
			hashes.SHA512 = hash
		}
	}

	return hashes, nil
}

// md5Hash 计算 MD5 哈希值。
func md5Hash(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// sha1Hash 计算 SHA-1 哈希值。
func sha1Hash(r io.Reader) (string, error) {
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// sha224HashSimplified 计算 SHA-224 哈希值（简化版，不带额外参数）。
func sha224HashSimplified(r io.Reader) (string, error) {
	h := sha256.New224()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil // SHA-224 是 SHA-256 截断的前 28 字节
}

// sha384HashSimplified 计算 SHA-384 哈希值（简化版，不带额外参数）。
func sha384HashSimplified(r io.Reader) (string, error) {
	h := sha512.New384()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil // SHA-384 是 SHA-512 截断的前 48 字节（但通常我们会取前64个字符，即32字节的十六进制表示）
}

// 注意：上面的 sha384HashSimplified 函数实际上返回的是 SHA-384 的完整十六进制字符串，
// 它比 48 字节的原始二进制表示要长，因为每个字节被编码为两个十六进制字符。
// 如果您只需要前48字节的十六进制表示（即96个字符），您应该截取返回的字符串：
// return hex.EncodeToString(h.Sum(nil))[:96], nil
// 但通常我们保留完整的哈希值字符串。

// sha256Hash 计算 SHA-256 哈希值。
func sha256Hash(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// sha512Hash 计算 SHA-512 哈希值。
func sha512Hash(r io.Reader) (string, error) {
	h := sha512.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
