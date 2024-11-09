/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:55:15
 * @FilePath: \go-toolbox\pkg\sign\hash.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"
	"io"
)

// Hasher 接口定义了一个计算哈希值的方法。
type Hasher interface {
	Hash(io.Reader) (string, error)
}

// GenericHasher 结构体用于通用哈希计算。
type GenericHasher struct {
	h hash.Hash
}

// 定义算法名称常量
type HashCryptoFunc string

const (
	AlgorithmMD5    HashCryptoFunc = "MD5"
	AlgorithmSHA1   HashCryptoFunc = "SHA1"
	AlgorithmSHA224 HashCryptoFunc = "SHA224"
	AlgorithmSHA256 HashCryptoFunc = "SHA256"
	AlgorithmSHA384 HashCryptoFunc = "SHA384"
	AlgorithmSHA512 HashCryptoFunc = "SHA512"
)

// SupportHashCryptoFunc 支持的哈希算法映射
var SupportHashCryptoFunc = map[HashCryptoFunc]func() hash.Hash{
	AlgorithmMD5:    md5.New,
	AlgorithmSHA1:   sha1.New,
	AlgorithmSHA224: sha256.New224,
	AlgorithmSHA256: sha256.New,
	AlgorithmSHA384: sha512.New384,
	AlgorithmSHA512: sha512.New,
}

// ErrUnsupportedAlgorithm 表示不支持的算法错误。
var ErrUnsupportedAlgorithm = errors.New("不支持的哈希算法")

// NewHasher 创建一个新的 GenericHasher 实例。
func NewHasher(algorithm HashCryptoFunc) (Hasher, error) {
	if hasherFunc, exists := SupportHashCryptoFunc[algorithm]; exists {
		return &GenericHasher{h: hasherFunc()}, nil
	}
	return nil, ErrUnsupportedAlgorithm
}

// Hash 计算哈希值并返回十六进制字符串。
func (g *GenericHasher) Hash(r io.Reader) (string, error) {
	if _, err := io.Copy(g.h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(g.h.Sum(nil)), nil
}
