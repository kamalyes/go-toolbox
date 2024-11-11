/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-23 17:37:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:24:40
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
	"hash"
	"io"
)

// Hasher 接口定义了一个计算哈希值的方法。
type Hasher interface {
	Hash(io.Reader) (string, error)
}

// MD5Hasher 结构体实现了 Hasher 接口，用于计算 MD5 哈希值。
type MD5Hasher struct {
	h hash.Hash
}

func NewMD5Hasher() Hasher {
	return &MD5Hasher{h: md5.New()}
}

func (h *MD5Hasher) Hash(r io.Reader) (string, error) {
	return hashData(h.h, r)
}

// SHA1Hasher 结构体实现了 Hasher 接口，用于计算 SHA-1 哈希值。
type SHA1Hasher struct {
	h hash.Hash
}

func NewSHA1Hasher() Hasher {
	return &SHA1Hasher{h: sha1.New()}
}

func (h *SHA1Hasher) Hash(r io.Reader) (string, error) {
	return hashData(h.h, r)
}

// SHA224Hasher 结构体实现了 Hasher 接口，用于计算 SHA-224 哈希值。
type SHA224Hasher struct {
	h hash.Hash
}

func NewSHA224Hasher() Hasher {
	return &SHA224Hasher{h: sha256.New224()}
}

func (h *SHA224Hasher) Hash(r io.Reader) (string, error) {
	return hashData(h.h, r)
}

// SHA256Hasher 结构体实现了 Hasher 接口，用于计算 SHA-256 哈希值。
type SHA256Hasher struct {
	h hash.Hash
}

func NewSHA256Hasher() Hasher {
	return &SHA256Hasher{h: sha256.New()}
}

func (h *SHA256Hasher) Hash(r io.Reader) (string, error) {
	return hashData(h.h, r)
}

// SHA384Hasher 结构体实现了 Hasher 接口，用于计算 SHA-384 哈希值。
type SHA384Hasher struct {
	h hash.Hash
}

func NewSHA384Hasher() Hasher {
	return &SHA384Hasher{h: sha512.New384()}
}

func (h *SHA384Hasher) Hash(r io.Reader) (string, error) {
	return hashData(h.h, r)
}

// SHA512Hasher 结构体实现了 Hasher 接口，用于计算 SHA-512 哈希值。
type SHA512Hasher struct {
	h hash.Hash
}

func NewSHA512Hasher() Hasher {
	return &SHA512Hasher{h: sha512.New()}
}

func (h *SHA512Hasher) Hash(r io.Reader) (string, error) {
	return hashData(h.h, r)
}

// hashData 是一个辅助函数，它接受一个已经初始化的哈希实例和一个读取器，并返回计算出的哈希值的十六进制字符串。
func hashData(h hash.Hash, r io.Reader) (string, error) {
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
