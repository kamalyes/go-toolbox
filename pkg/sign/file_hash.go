/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-18 15:21:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-18 15:52:36
 * @FilePath: \go-toolbox\pkg\sign\file_hash.go
 * @Description: 文件多哈希计算（从 osx 迁入：osx 仅为此一个函数就要依赖 sign -> x/crypto）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package sign

import (
	"encoding/hex"
	"hash"
	"io"
	"os"
)

// ComputeHashes 计算文件的多种哈希值，优化为只读一次文件
func ComputeHashes(filePath string) (map[HashCryptoFunc]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建哈希实例
	hashes := map[HashCryptoFunc]hash.Hash{}
	for algo, constructor := range SupportHMACCryptoFunc {
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

	results := make(map[HashCryptoFunc]string, len(hashes))
	for algo, h := range hashes {
		results[algo] = hex.EncodeToString(h.Sum(nil))
	}

	return results, nil
}
