/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 00:15:15
 * @FilePath: \go-toolbox\pkg\osx\path.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// MkdirIfNotExist 如果目录不存在则创建它
func MkdirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// DirHasContent 检查目录是否有内容（即是否有非空文件）
func DirHasContent(dir string) (bool, []string, error) {
	var files []string
	// 读取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, nil, err
	}

	// 遍历目录项
	for _, entry := range entries {
		// 获取文件的完整路径
		path := filepath.Join(dir, entry.Name())

		// 获取文件信息
		info, err := entry.Info()
		if err != nil {
			// 如果无法获取文件信息，则跳过该文件
			continue
		}

		// 检查文件是否为非空文件
		if !info.IsDir() && info.Size() > 0 {
			files = append(files, path)
		}
	}

	return true, files, nil
}

// Copy 复制文件从源路径到目标路径
func Copy(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destDir := filepath.Dir(dest)
	err = MkdirIfNotExist(destDir)
	if err != nil {
		return err
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	destFile.Chmod(os.ModePerm) // 设置文件权限

	_, err = io.Copy(destFile, srcFile)
	return err
}

// MkdirTemp 创建一个临时目录，如果创建失败则程序退出
func MkdirTemp() string {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatalln(err)
	}
	return dir
}
