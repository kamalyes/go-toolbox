/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 17:21:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:15:01
 * @FilePath: \go-toolbox\pkg\sign\rsa.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// 错误定义
var (
	ErrPrivateKey  = errors.New("私钥文件错误")
	ErrPublicKey   = errors.New("公钥解析错误")
	ErrNotRsaKey   = errors.New("不是RSA密钥")
	ErrEncryptFail = errors.New("加密失败")
	ErrDecryptFail = errors.New("解密失败")
)

// PEM类型常量
const (
	privateKeyType = "PRIVATE KEY"
	publicKeyType  = "PUBLIC KEY"
)

// RsaKeySize 定义了支持的RSA密钥大小
type RsaKeySize int

const (
	RsaKeySize512  RsaKeySize = 512
	RsaKeySize1024 RsaKeySize = 1024
	RsaKeySize2048 RsaKeySize = 2048
	RsaKeySize4096 RsaKeySize = 4096
)

// RsaKeyPair 包含RSA密钥对
type RsaKeyPair struct {
	HashCrypto HashCryptoFunc
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// RsaCrypto 定义了RSA加解密器的接口
type RsaCrypto interface {
	Encrypt(input []byte) ([]byte, error)
	Decrypt(input []byte) ([]byte, error)
	DecryptBase64(input string) ([]byte, error)
	GetPrivateKey() *rsa.PrivateKey
	GetPublicKey() *rsa.PublicKey
}

type HashCryptoFunc int

const (
	HashCryptoFuncMd5 HashCryptoFunc = iota
	HashCryptoFuncSHA1
	HashCryptoFuncSHA224
	HashCryptoFuncSHA384
	HashCryptoFuncSHA256
	HashCryptoFuncSHA512
)

// rsaBase 是RSA加密器和解密器的公共基础结构体
type rsaBase struct {
	keySize RsaKeySize // 密钥大小
}

// rsaCryptoImpl 是RSA加解密器的实现
type rsaCryptoImpl struct {
	rsaBase
	hashCrypto HashCryptoFunc
	privateKey *rsa.PrivateKey // 私钥
	publicKey  *rsa.PublicKey  // 公钥
}

// encodePEM 封装了PEM编码的逻辑
func encodePEM(blockType string, bytes []byte) string {
	block := &pem.Block{
		Type:  blockType,
		Bytes: bytes,
	}
	return string(pem.EncodeToMemory(block))
}

// GenerateRsaKeyPair 生成指定大小的 RSA 密钥对
func GenerateRsaKeyPair(keySize RsaKeySize) (*RsaKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, int(keySize))
	if err != nil {
		return nil, err
	}
	return &RsaKeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// ExportRsaPrivateKeyToPEM 将 RSA 私钥导出为 PEM 格式，支持 PKCS#1 和 PKCS#8
func ExportRsaPrivateKeyToPEM(privateKey *rsa.PrivateKey, pkcs8 bool) (string, error) {
	var privBytes []byte
	var err error
	if pkcs8 {
		privBytes, err = x509.MarshalPKCS8PrivateKey(privateKey)
	} else {
		privBytes = x509.MarshalPKCS1PrivateKey(privateKey)
	}
	if err != nil {
		return "", err
	}
	return encodePEM(privateKeyType, privBytes), nil
}

// ExportRsaPublicKeyToPEM 将 RSA 公钥导出为 PEM 格式
func ExportRsaPublicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return encodePEM(publicKeyType, pubBytes), nil
}

// NewRsaCryptoFromKeys 根据公钥和私钥创建RSA加解密器
func NewRsaCryptoFromKeys(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, keySize RsaKeySize) RsaCrypto {
	return &rsaCryptoImpl{
		rsaBase:    rsaBase{keySize: keySize},
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// Decrypt 实现了解密功能，使用 OAEP
func (r *rsaCryptoImpl) Decrypt(input []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privateKey, input, nil)
}

// DecryptBase64 先将Base64编码的字符串解码，然后解密
func (r *rsaCryptoImpl) DecryptBase64(input string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return r.Decrypt(decoded)
}

// Encrypt 实现了加密功能，使用 OAEP
func (r *rsaCryptoImpl) Encrypt(input []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, r.publicKey, input, nil)
}

// GetPrivateKey 返回 rsaCryptoImpl 中的私钥
func (r *rsaCryptoImpl) GetPrivateKey() *rsa.PrivateKey {
	return r.privateKey
}

// GetPublicKey 返回 rsaCryptoImpl 中的公钥
func (r *rsaCryptoImpl) GetPublicKey() *rsa.PublicKey {
	return r.publicKey
}

// NewRsaCryptoFromPrivateFile 根据私钥文件创建RSA加解密器
func NewRsaCryptoFromPrivateFile(filePath string, keySize RsaKeySize) (RsaCrypto, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPrivateKey, err)
	}

	privateKey, err := parsePrivateKey(content, keySize)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey
	return NewRsaCryptoFromKeys(privateKey, publicKey, keySize), nil
}

// parsePrivateKey 解析PEM格式的私钥并验证大小
func parsePrivateKey(content []byte, keySize RsaKeySize) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, ErrPrivateKey
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	if privateKey.N.BitLen() != int(keySize) {
		return nil, fmt.Errorf("私钥大小不匹配，期望%d位，实际%d位", keySize, privateKey.N.BitLen())
	}

	return privateKey, nil
}

// NewRsaCryptoFromPublicPEM 根据PEM格式的公钥字节数组创建RSA加解密器
func NewRsaCryptoFromPublicPEM(pemKey []byte, keySize RsaKeySize) (RsaCrypto, error) {
	block, _ := pem.Decode(pemKey)
	if block == nil {
		return nil, ErrPublicKey
	}

	publicKey, err := parsePublicKey(block.Bytes, keySize)
	if err != nil {
		return nil, err
	}

	return NewRsaCryptoFromKeys(nil, publicKey, keySize), nil
}

// parsePublicKey 解析PEM格式的公钥并验证大小
func parsePublicKey(blockBytes []byte, keySize RsaKeySize) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(blockBytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, ErrNotRsaKey
	}

	if rsaPub.N.BitLen() != int(keySize) {
		return nil, fmt.Errorf("公钥大小不匹配，期望%d位，实际%d位", keySize, rsaPub.N.BitLen())
	}

	return rsaPub, nil
}
