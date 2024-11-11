/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 10:30:31
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:31:12
 * @FilePath: \go-toolbox\tests\sign_rsa_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"encoding/base64"
	"os"
	"strconv"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/sign"
	"github.com/stretchr/testify/assert"
)

// 公共测试数据结构体
type TestRsaData struct {
	OriginalMessage []byte
	EncryptedBase64 string
}

var testRsaData = TestRsaData{
	OriginalMessage: []byte("Hello, RSA!"),
	EncryptedBase64: base64.StdEncoding.EncodeToString([]byte("Hello, RSA with Base64!")),
}

// 辅助函数：生成密钥对并创建 RSA 加解密器
func createRsaCrypto(keySize sign.RsaKeySize) (sign.RsaCrypto, error) {
	keyPair, err := sign.GenerateRsaKeyPair(keySize)
	if err != nil {
		return nil, err
	}
	return sign.NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, keySize), nil
}

func TestRsaCrypto(t *testing.T) {
	// 定义支持的密钥大小
	keySizes := []sign.RsaKeySize{
		sign.RsaKeySize512,
		sign.RsaKeySize1024,
		sign.RsaKeySize2048,
		sign.RsaKeySize4096,
	}

	for _, keySize := range keySizes {
		keySizeStr := strconv.Itoa(int(keySize))
		t.Run("TestRsaCrypto_"+keySizeStr, func(t *testing.T) {
			rsaCrypto, err := createRsaCrypto(keySize)
			assert.NoError(t, err, "生成密钥对失败")

			// 测试用例：正向加密
			encryptedMessage, err := rsaCrypto.Encrypt(testRsaData.OriginalMessage)
			assert.NoError(t, err, "加密失败")

			// 测试用例：逆向解密
			decryptedMessage, err := rsaCrypto.Decrypt(encryptedMessage)
			assert.NoError(t, err, "解密失败")

			// 验证解密后的消息是否与原始消息相同
			assert.Equal(t, testRsaData.OriginalMessage, decryptedMessage, "解密后的消息与原始消息不匹配")
		})
	}
}

func TestRsaCryptoBase64(t *testing.T) {
	// 定义支持的密钥大小
	keySizes := []sign.RsaKeySize{
		sign.RsaKeySize512,
		sign.RsaKeySize1024,
		sign.RsaKeySize2048,
		sign.RsaKeySize4096,
	}

	for _, keySize := range keySizes {
		keySizeStr := strconv.Itoa(int(keySize))
		t.Run("TestRsaCryptoBase64_"+keySizeStr, func(t *testing.T) {
			rsaCrypto, err := createRsaCrypto(keySize)
			assert.NoError(t, err, "生成密钥对失败")

			// 测试用例：正向加密
			encryptedMessage, err := rsaCrypto.Encrypt([]byte(testRsaData.EncryptedBase64))
			assert.NoError(t, err, "加密失败")

			// 将加密后的消息转换为 Base64
			encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedMessage)

			// 测试用例：逆向解密
			decryptedMessage, err := rsaCrypto.DecryptBase64(encryptedBase64)
			assert.NoError(t, err, "Base64解密失败")

			// 验证解密后的消息是否与原始消息相同
			assert.Equal(t, []byte(testRsaData.EncryptedBase64), decryptedMessage, "解密后的消息与原始消息不匹配")
		})
	}
}

func TestNewRsaCryptoFromPrivateFile(t *testing.T) {
	// 生成测试私钥
	keySize := sign.RsaKeySize2048
	privateKey, _ := sign.GenerateRsaKeyPair(keySize)
	pemData, err := sign.ExportRsaPrivateKeyToPEM(privateKey.PrivateKey, false)
	assert.NoError(t, err)

	// 写入临时文件
	filePath := "test_private_key.pem"
	err = os.WriteFile(filePath, []byte(pemData), 0644)
	assert.NoError(t, err)
	defer os.Remove(filePath) // 清理临时文件

	// 测试创建RSA加解密器
	rsaCrypto, err := sign.NewRsaCryptoFromPrivateFile(filePath, keySize)
	assert.NoError(t, err)
	assert.NotNil(t, rsaCrypto)
	assert.Equal(t, privateKey.PrivateKey, rsaCrypto.GetPrivateKey())
}

func TestNewRsaCryptoFromPublicPEM(t *testing.T) {
	// 生成测试密钥对
	keySize := sign.RsaKeySize2048
	privateKey, _ := sign.GenerateRsaKeyPair(keySize)
	pemData, err := sign.ExportRsaPublicKeyToPEM(privateKey.PublicKey)
	assert.NoError(t, err)

	// 测试创建RSA加解密器
	rsaCrypto, err := sign.NewRsaCryptoFromPublicPEM([]byte(pemData), keySize)
	assert.NoError(t, err)
	assert.NotNil(t, rsaCrypto)
}

func TestInvalidPrivateKeyFile(t *testing.T) {
	// 测试无效的私钥文件
	filePath := "invalid_private_key.pem"
	_, err := sign.NewRsaCryptoFromPrivateFile(filePath, sign.RsaKeySize2048)
	assert.Error(t, err)
}

func TestInvalidPublicKeyPEM(t *testing.T) {
	// 测试无效的公钥PEM
	_, err := sign.NewRsaCryptoFromPublicPEM([]byte("invalid public key"), sign.RsaKeySize2048)
	assert.Error(t, err)
}
