/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-11 00:55:15
 * @FilePath: \go-toolbox\tests\sign_rsa_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"hash"
	"os"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/sign"
	"github.com/stretchr/testify/assert"
)

// 辅助函数：生成 RSA 密钥对并进行有效性断言
func generateAndAssertRsaKeyPair(size sign.RsaKeySize, t *testing.T) *sign.RsaKeyPair {
	keyPair, err := sign.GenerateRsaKeyPair(size)
	assert.NoError(t, err, "生成 RSA 密钥对时发生错误")
	assert.NotNil(t, keyPair.PrivateKey, "私钥应不为 nil")
	assert.NotNil(t, keyPair.PublicKey, "公钥应不为 nil")
	return keyPair
}

// 辅助函数：导出 RSA 密钥到 PEM 格式并进行有效性断言
func exportAndAssertRsaKeysToPEM(keyPair *sign.RsaKeyPair, t *testing.T) {
	privPEM, err := sign.ExportRsaPrivateKeyToPEM(keyPair.PrivateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, privPEM, "导出的私钥 PEM 不应为空")

	pubPEM, err := sign.ExportRsaPublicKeyToPEM(keyPair.PublicKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, pubPEM, "导出的公钥 PEM 不应为空")
}

// 辅助函数：测试 RSA 加密和解密功能，使用不同的哈希函数
func testRsaCryptoEncryptDecryptWithHashFunc(keyPair *sign.RsaKeyPair, hashFuncs []func() hash.Hash, t *testing.T) {
	for _, hashFunc := range hashFuncs {
		rsaCrypto := sign.NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, hashFunc)

		originalText := []byte("Hello, RSA!") // 原始文本
		salt := []byte("salt")                // 盐值
		encryptedText, err := rsaCrypto.EncryptSalt(originalText, salt)
		assert.NoError(t, err, "加密时发生错误")
		assert.NotNil(t, encryptedText, "加密后的文本不应为 nil")

		decryptedText, err := rsaCrypto.Decrypt(encryptedText)
		assert.NoError(t, err, "解密时发生错误")
		assert.Equal(t, originalText, decryptedText[len(salt):], "解密后的文本应与原文本匹配")
	}
}

// 辅助函数：测试 RSA Base64 解密功能，使用不同的哈希函数
func testRsaCryptoDecryptBase64WithHashFunc(keyPair *sign.RsaKeyPair, hashFuncs []func() hash.Hash, t *testing.T) {
	for _, hashFunc := range hashFuncs {
		rsaCrypto := sign.NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, hashFunc)

		originalText := []byte("Hello, RSA!") // 原始文本
		salt := []byte("salt")                // 盐值
		encryptedText, err := rsaCrypto.EncryptSalt(originalText, salt)
		assert.NoError(t, err)

		// 将加密文本转换为 Base64
		encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedText)

		// 测试 Base64 解密
		decryptedText, err := rsaCrypto.DecryptBase64(encryptedBase64)
		assert.NoError(t, err, "Base64 解密时发生错误")
		assert.Equal(t, originalText, decryptedText[len(salt):], "解密后的文本应与原文本匹配")
	}
}

// 测试 EncryptRandSalt 函数
func TestEncryptRandSalt(t *testing.T) {
	keyPair := generateAndAssertRsaKeyPair(sign.RsaKeySize2048, t)

	hashFuncs := []func() hash.Hash{
		sha256.New,
		sha512.New,
		sha1.New,
	}
	for _, hashFunc := range hashFuncs {
		// 生成 RSA 密钥对
		rsaCrypto := sign.NewRsaCryptoFromKeys(keyPair.PrivateKey, keyPair.PublicKey, hashFunc)

		input := []byte("Hello, World!")

		// 测试默认盐长度
		encryptedData, salt, err := rsaCrypto.EncryptRandSalt(input)
		assert.NoError(t, err, "加密失败")
		assert.Equal(t, 16, len(salt), "默认盐长度应为 16")
		assert.NotEmpty(t, encryptedData, "加密数据应不为空")

		// 解密测试
		decryptedData, err := rsaCrypto.Decrypt(encryptedData)
		assert.NoError(t, err, "解密失败")
		assert.Equal(t, input, decryptedData[16:], "解密后的数据应与原始数据匹配")

		// 测试自定义盐长度
		customSaltLength := 32
		encryptedData, salt, err = rsaCrypto.EncryptRandSalt(input, customSaltLength)
		assert.NoError(t, err, "加密失败")
		assert.Equal(t, customSaltLength, len(salt), "自定义盐长度应为 %d", customSaltLength)
		assert.NotEmpty(t, encryptedData, "加密数据应不为空")

		// 解密测试
		decryptedData, err = rsaCrypto.Decrypt(encryptedData)
		assert.NoError(t, err, "解密失败")
		assert.Equal(t, input, decryptedData[customSaltLength:], "解密后的数据应与原始数据匹配")
	}

}

// 测试生成 RSA 密钥对
func TestGenerateRsaKeyPair(t *testing.T) {
	keySizes := []sign.RsaKeySize{sign.RsaKeySize512, sign.RsaKeySize1024, sign.RsaKeySize2048, sign.RsaKeySize4096}
	for _, size := range keySizes {
		generateAndAssertRsaKeyPair(size, t)
	}
}

// 测试导出 RSA 私钥和公钥为 PEM 格式
func TestExportRsaKeysToPEM(t *testing.T) {
	keyPair := generateAndAssertRsaKeyPair(sign.RsaKeySize2048, t)
	exportAndAssertRsaKeysToPEM(keyPair, t)
}

// 测试 RSA 加解密功能，使用不同的哈希函数
func TestRsaCryptoEncryptDecryptWithHashFunc(t *testing.T) {
	keyPair := generateAndAssertRsaKeyPair(sign.RsaKeySize2048, t)

	hashFuncs := []func() hash.Hash{
		sha1.New,
		sha256.New224,
		sha256.New,
		sha512.New384,
		sha512.New,
	}

	testRsaCryptoEncryptDecryptWithHashFunc(keyPair, hashFuncs, t)
}

// 测试 RSA Base64 解密功能，使用不同的哈希函数
func TestRsaCryptoDecryptBase64WithHashFunc(t *testing.T) {
	keyPair := generateAndAssertRsaKeyPair(sign.RsaKeySize2048, t)

	hashFuncs := []func() hash.Hash{
		sha256.New,
		sha512.New,
		sha1.New,
	}

	testRsaCryptoDecryptBase64WithHashFunc(keyPair, hashFuncs, t)
}

// 测试从私钥文件创建 RSA 加解密器
func TestNewRsaCryptoFromPrivateFile(t *testing.T) {
	keyPair := generateAndAssertRsaKeyPair(sign.RsaKeySize2048, t)

	privPEM, err := sign.ExportRsaPrivateKeyToPEM(keyPair.PrivateKey)
	assert.NoError(t, err)

	// 将私钥写入临时文件
	tempFile, err := os.CreateTemp("", "private_key.pem")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(privPEM))
	assert.NoError(t, err)
	tempFile.Close()

	// 从文件中创建 RSA 加解密器
	rsaCrypto, err := sign.NewRsaCryptoFromPrivateFile(tempFile.Name(), sha256.New)
	assert.NoError(t, err)
	assert.NotNil(t, rsaCrypto.GetPrivateKey(), "私钥应不为 nil")
	assert.NotNil(t, rsaCrypto.GetPublicKey(), "公钥应不为 nil")
}

// 测试从 PEM 格式公钥创建 RSA 加解密器
func TestNewRsaCryptoFromPublicPEM(t *testing.T) {
	keyPair := generateAndAssertRsaKeyPair(sign.RsaKeySize2048, t)

	pubPEM, err := sign.ExportRsaPublicKeyToPEM(keyPair.PublicKey)
	assert.NoError(t, err)

	// 输出 PEM 格式公钥，便于调试
	t.Logf("公钥 PEM:\n%s", pubPEM)

	// 确保 PEM 格式正确
	if err := isValidPEM(string(pubPEM)); err != nil {
		t.Fatalf("公钥 PEM 格式不正确: %v", err) // 输出 PEM 内容以便调试
	}

	rsaCrypto, err := sign.NewRsaCryptoFromPublicPEM([]byte(pubPEM), sha256.New)
	assert.NoError(t, err)
	assert.NotNil(t, rsaCrypto.GetPublicKey(), "公钥应不为 nil")
}

// 使用标准库检查 PEM 格式的有效性
func isValidPEM(pemData string) error {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return errors.New("无法解析 PEM 块")
	}
	if block.Type != "PUBLIC KEY" {
		return errors.New("PEM 块类型不正确")
	}
	return nil
}

// 测试解析 PEM 格式公钥
func TestParsePublicKey(t *testing.T) {
	keyPair := generateAndAssertRsaKeyPair(sign.RsaKeySize2048, t)

	pubPEM, err := sign.ExportRsaPublicKeyToPEM(keyPair.PublicKey)
	assert.NoError(t, err)

	// 输出 PEM 格式公钥，便于调试
	t.Logf("公钥 PEM:\n%s", pubPEM)

	publicKey, err := sign.ParsePublicKey([]byte(pubPEM))
	if err != nil {
		t.Fatalf("公钥解析错误: %v", err) // 输出详细错误信息
	}
	assert.NotNil(t, publicKey, "解析后的公钥应不为 nil")

	// 校验生成的公钥和解析的公钥是否一致
	// 重新编码解析后的公钥为 PEM 格式
	reEncodedPEM, err := sign.ExportRsaPublicKeyToPEM(publicKey)
	assert.NoError(t, err)

	// 比较原始 PEM 和重新编码的 PEM 是否一致
	assert.Equal(t, pubPEM, reEncodedPEM, "原始公钥 PEM 和重新编码的公钥 PEM 应该一致")
}

// 测试解析 PEM 格式私钥
func TestParsePrivateKey(t *testing.T) {
	keyPair := generateAndAssertRsaKeyPair(sign.RsaKeySize2048, t)

	privPEM, err := sign.ExportRsaPrivateKeyToPEM(keyPair.PrivateKey)
	assert.NoError(t, err)

	privateKey, err := sign.ParsePrivateKey([]byte(privPEM))
	assert.NoError(t, err)
	assert.NotNil(t, privateKey, "解析后的私钥应不为 nil")
	// 校验生成的私钥和解析的私钥是否一致
	// 重新编码解析后的私钥为 PEM 格式
	reEncodedPEM, err := sign.ExportRsaPrivateKeyToPEM(privateKey)
	assert.NoError(t, err)

	// 比较原始 PEM 和重新编码的 PEM 是否一致
	assert.Equal(t, privPEM, reEncodedPEM, "原始私钥 PEM 和重新编码的私钥 PEM 应该一致")
}
