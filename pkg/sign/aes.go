/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-23 17:37:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-23 17:58:07
 * @FilePath: \go-toolbox\pkg\sign\aes.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/sha3"
)

// GenerateByteKey 生成一个指定字节的密钥
func GenerateByteKey(password string, length int) []byte {
	hash := sha3.New256()
	hash.Write([]byte(password))
	return hash.Sum(nil)[:length]
}

// AesEncrypt 加密函数
func AesEncrypt(plainText string, key []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key cannot be empty")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7 填充
	plainTextBytes := []byte(plainText)
	plainTextBytes = pkcs7Padding(plainTextBytes, block.BlockSize())

	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AesDecryptRaw 解密函数（直接接收原始字节，不经过 base64 解码）
// 适用于上层已通过 JSON/protojson 对 bytes 字段做 base64 解码的场景（如 grpc-gateway 的 body 绑定）
// cipherBytes 格式：iv(16字节) + aes-cbc-pkcs7-ciphertext
func AesDecryptRaw(cipherBytes []byte, key []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key cannot be empty")
	}
	if len(cipherBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := cipherBytes[:aes.BlockSize]
	encrypted := cipherBytes[aes.BlockSize:]

	if len(encrypted)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	plainTextBytes := make([]byte, len(encrypted))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainTextBytes, encrypted)

	plainTextBytes, err = pkcs7Unpadding(plainTextBytes)
	if err != nil {
		return "", err
	}

	return string(plainTextBytes), nil
}

// AesDecrypt 解密函数
func AesDecrypt(cipherText string, key []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key cannot be empty")
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherTextBytes, cipherTextBytes)

	// 去掉 PKCS7 填充
	plainTextBytes, err := pkcs7Unpadding(cipherTextBytes)
	if err != nil {
		return "", err
	}

	return string(plainTextBytes), nil
}

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// 去掉 PKCS7 填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("data is empty")
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("invalid padding")
	}
	return data[:(length - unpadding)], nil
}

// AesEncryptWithIV 使用自定义 IV 的 AES-CBC-PKCS7 加密
// [EN] AES-CBC-PKCS7 encryption with custom IV
func AesEncryptWithIV(plainText string, key, iv []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key cannot be empty")
	}
	if len(iv) != aes.BlockSize {
		return "", errors.New("IV length must equal block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7 填充
	plainTextBytes := []byte(plainText)
	plainTextBytes = pkcs7Padding(plainTextBytes, block.BlockSize())

	cipherText := make([]byte, len(plainTextBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AesDecryptWithIV 使用自定义 IV 的 AES-CBC-PKCS7 解密
// [EN] AES-CBC-PKCS7 decryption with custom IV
func AesDecryptWithIV(cipherText string, key, iv []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key cannot be empty")
	}
	if len(iv) != aes.BlockSize {
		return "", errors.New("IV length must equal block size")
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	plainTextBytes := make([]byte, len(cipherTextBytes))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainTextBytes, cipherTextBytes)

	// 去掉 PKCS7 填充
	plainTextBytes, err = pkcs7Unpadding(plainTextBytes)
	if err != nil {
		return "", err
	}

	return string(plainTextBytes), nil
}
