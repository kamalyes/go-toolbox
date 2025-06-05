/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-09 17:15:16
 * @FilePath: \go-toolbox\pkg\sign\hash.go
 * @Description: 通用 HMAC 签名器实现，支持多种哈希算法，基于 crypto/hmac 包
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package sign

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
)

// Signer 接口定义了签名和验证方法
type Signer interface {
	// Algorithm 返回签名算法名称
	Algorithm() HashCryptoFunc
	// Sign 计算 data 的签名，使用给定的 key，返回签名字节切片或错误
	Sign(data, key []byte) ([]byte, error)
	// Verify 验证 data 的签名是否正确，使用给定的 key 和签名 signature，返回验证结果和错误
	Verify(data, key, signature []byte) (bool, error)
}

// GenericHMACSigner 是一个基于 HMAC 的通用签名器
type GenericHMACSigner struct {
	algorithm HashCryptoFunc   // 算法名称，如 "HS256"
	hashFunc  func() hash.Hash // 哈希函数构造器，如 sha256.New
}

// NewGenericHMACSigner 创建一个新的 GenericHMACSigner 实例
// algorithm 参数是算法名称，hashFunc 是对应的哈希函数构造器
func NewGenericHMACSigner(algorithm HashCryptoFunc, hashFunc func() hash.Hash) *GenericHMACSigner {
	return &GenericHMACSigner{
		algorithm: algorithm,
		hashFunc:  hashFunc,
	}
}

// Algorithm 返回签名算法名称
func (s *GenericHMACSigner) Algorithm() HashCryptoFunc {
	return s.algorithm
}

// Sign 使用 HMAC 算法计算 data 的签名，key 是密钥
func (s *GenericHMACSigner) Sign(data, key []byte) ([]byte, error) {
	// 创建一个新的 HMAC 哈希器
	mac := hmac.New(s.hashFunc, key)
	// 写入数据到哈希器
	_, err := mac.Write(data)
	if err != nil {
		return nil, err
	}
	// 计算并返回签名结果
	return mac.Sum(nil), nil
}

// Verify 验证 data 的签名 signature 是否正确，key 是密钥
func (s *GenericHMACSigner) Verify(data, key, signature []byte) (bool, error) {
	// 计算期望的签名
	expected, err := s.Sign(data, key)
	if err != nil {
		return false, err
	}
	// 使用 hmac.Equal 做安全比较，防止时序攻击
	return hmac.Equal(signature, expected), nil
}

// ErrUnsupportedAlgorithmHMAC 表示不支持的 HMAC 算法错误
var ErrUnsupportedAlgorithmHMAC = errors.New("不支持的 HMAC 算法")

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

// SupportHMACCryptoFunc 支持的 HMAC 哈希算法映射，key 是算法名称，value 是哈希函数构造器
var SupportHMACCryptoFunc = map[HashCryptoFunc]func() hash.Hash{
	AlgorithmMD5:    md5.New,
	AlgorithmSHA1:   sha1.New,
	AlgorithmSHA224: sha256.New224,
	AlgorithmSHA256: sha256.New,
	AlgorithmSHA384: sha512.New384,
	AlgorithmSHA512: sha512.New,
}

// NewHMACSigner 根据算法名称创建对应的通用 HMAC 签名器
func NewHMACSigner(algorithm HashCryptoFunc) (Signer, error) {
	if hashFunc, exists := SupportHMACCryptoFunc[algorithm]; exists {
		return NewGenericHMACSigner(algorithm, hashFunc), nil
	}
	return nil, ErrUnsupportedAlgorithmHMAC
}
