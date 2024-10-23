/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-10-23 17:37:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-23 17:58:07
 * @FilePath: \go-toolbox\sign\aes.go
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
