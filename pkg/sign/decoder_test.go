/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 12:17:56
 * @FilePath: \go-toolbox\pkg\sign\decoder_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type decodeBuildCallbackPayload struct {
	BuildID    string            `json:"buildId"`
	ExternalID string            `json:"externalId"`
	Status     string            `json:"status"`
	Message    string            `json:"message"`
	Artifact   decodeArtifact    `json:"artifact"`
	Steps      []decodeBuildStep `json:"steps"`
	Metadata   map[string]string `json:"metadata"`
}

type decodeArtifact struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

type decodeBuildStep struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Duration int64  `json:"duration"`
}

func TestDecodeJSON(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `{"buildId":"build_001","externalId":"easycli_job_001","status":"success","message":"build completed","artifact":{"url":"https://cdn.example.com/app.tar.gz","sha256":"abc123","size":2048},"steps":[{"name":"install","status":"success","duration":12},{"name":"compile","status":"success","duration":31}],"metadata":{"commit":"abcdef","branch":"main"}}`
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	decoder := NewEncryptedDecoder(WithAESKey(key))
	decoded, err := DecodeJSON[decodeBuildCallbackPayload](decoder, []byte(ciphertext))
	require.NoError(t, err)
	require.Equal(t, []byte(ciphertext), decoded.Ciphertext)
	require.Equal(t, []byte(plainText), decoded.Plaintext)
	require.Equal(t, "build_001", decoded.Payload.BuildID)
	require.Equal(t, "easycli_job_001", decoded.Payload.ExternalID)
	require.Equal(t, "success", decoded.Payload.Status)
	require.Equal(t, "https://cdn.example.com/app.tar.gz", decoded.Payload.Artifact.URL)
	require.Equal(t, int64(2048), decoded.Payload.Artifact.Size)
	require.Len(t, decoded.Payload.Steps, 2)
	require.Equal(t, "compile", decoded.Payload.Steps[1].Name)
	require.Equal(t, "abcdef", decoded.Payload.Metadata["commit"])
}

func TestDecodeJSONTo(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `{"buildId":"build_002","externalId":"easycli_job_002","status":"running","message":"build is running","artifact":{"url":"","sha256":"","size":0},"steps":[{"name":"compile","status":"running","duration":8}],"metadata":{"commit":"fedcba","branch":"release"}}`
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	var payload decodeBuildCallbackPayload
	decoder := NewEncryptedDecoder(WithAESKey(key))
	decodedPlainText, err := decoder.DecodeJSONTo([]byte(ciphertext), &payload)
	require.NoError(t, err)
	require.Equal(t, []byte(plainText), decodedPlainText)
	require.Equal(t, "build_002", payload.BuildID)
	require.Equal(t, "running", payload.Status)
	require.Equal(t, "release", payload.Metadata["branch"])
}

func TestDecodeJSONWithAESPassword(t *testing.T) {
	password := "webhook-callback-secret"
	key := GenerateByteKey(password, 32)
	plainText := `{"buildId":"build_003","externalId":"easycli_job_003","status":"failed","message":"compile failed","artifact":{"url":"","sha256":"","size":0},"steps":[{"name":"compile","status":"failed","duration":19}],"metadata":{"commit":"badc0de","error":"syntax"}}`
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	decoder := NewEncryptedDecoder(WithAESPassword(password))
	decoded, err := DecodeJSON[decodeBuildCallbackPayload](decoder, []byte(ciphertext))
	require.NoError(t, err)
	require.Equal(t, "build_003", decoded.Payload.BuildID)
	require.Equal(t, "failed", decoded.Payload.Status)
	require.Equal(t, "syntax", decoded.Payload.Metadata["error"])
}

func TestDecodeProtoJSON(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `"buildId,status,artifact.url,metadata.commit"`
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	decoder := NewEncryptedDecoder(WithAESKey(key))
	decoded, err := DecodeProtoJSON(decoder, []byte(ciphertext), func() *fieldmaskpb.FieldMask {
		return &fieldmaskpb.FieldMask{}
	})
	require.NoError(t, err)
	require.Equal(t, []string{"build_id", "status", "artifact.url", "metadata.commit"}, decoded.Payload.Paths)
}

func TestDecodeProtoJSONTo(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `"buildId,status,steps.name"`
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	payload := &fieldmaskpb.FieldMask{}
	decoder := NewEncryptedDecoder(WithAESKey(key))
	decodedPlainText, err := decoder.DecodeProtoJSONTo([]byte(ciphertext), payload)
	require.NoError(t, err)
	require.Equal(t, []byte(plainText), decodedPlainText)
	require.Equal(t, []string{"build_id", "status", "steps.name"}, payload.Paths)
}

func TestEncryptedDecoderRejectsMissingKey(t *testing.T) {
	_, err := NewEncryptedDecoder().Decrypt([]byte("ciphertext"))
	require.ErrorIs(t, err, ErrMissingAESKey)
}

func TestDecryptWithRawCiphertext(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `{"buildId":"build_raw","status":"success"}`

	// 加密得到 base64 字符串
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	// 模拟 grpc-gateway 的行为：对 bytes 字段做 base64 解码
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	require.NoError(t, err)

	// 使用 WithRawCiphertext 选项，直接传入原始字节
	decoder := NewEncryptedDecoder(WithAESKey(key), WithRawCiphertext())
	decoded, err := DecodeJSON[decodeBuildCallbackPayload](decoder, cipherBytes)
	require.NoError(t, err)
	require.Equal(t, "build_raw", decoded.Payload.BuildID)
	require.Equal(t, "success", decoded.Payload.Status)
}

func TestAesDecryptRaw(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := "hello raw decrypt"

	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	// base64 解码得到原始字节
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	require.NoError(t, err)

	// AesDecryptRaw 直接解密原始字节
	decrypted, err := AesDecryptRaw(cipherBytes, key)
	require.NoError(t, err)
	require.Equal(t, plainText, decrypted)
}

// TestWebhookE2E_SimulatesWorkerEncryptAndGatewayDecrypt 模拟完整 webhook 链路：
// 1. Worker 端加密（SHA3-256 派生密钥 + AES-CBC-PKCS7 + iv+ciphertext + base64）
// 2. grpc-gateway 对 bytes 字段做 base64 解码（得到原始字节）
// 3. Go 端用 WithRawCiphertext 解密
// 验证两端结构完全匹配
func TestWebhookE2E_SimulatesWorkerEncryptAndGatewayDecrypt(t *testing.T) {
	// 模拟 Worker 的 bodyEncryptionSecret（与 Worker normalizeSecret 后一致）
	secret := "webhook-callback-secret"

	// 模拟 Worker 的业务 payload（JSON 明文）
	plainText := `{"domain":"example.com","zone_id":"123456","status":"success","request_id":"req_001"}`

	// Step 1: 模拟 Worker 端加密
	// Worker: deriveWebhookEncryptionKey(secret) = SHA3-256(secret.trim()) 取前 32 字节
	key := GenerateByteKey(secret, 32)

	// Worker: AES-CBC 加密（Web Crypto API 自动 PKCS7 padding）
	// Go 端用 AesEncrypt 模拟（同样是 AES-CBC-PKCS7，单次 padding）
	ciphertextBase64, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	// Worker: buildWebhookRequestBody(ciphertext) = JSON.stringify(ciphertext)
	// 即 HTTP body = "base64..." (带双引号的 JSON string)
	httpBody := `"` + ciphertextBase64 + `"`

	// Step 2: 模拟 grpc-gateway 对 bytes 字段的处理
	// protojson 解析 {"body": "base64..."} 时，对 bytes 字段做 base64 解码
	// 等价于：base64.StdEncoding.DecodeString(ciphertextBase64)
	// 结果是 iv(16) + encrypted 的原始字节
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	require.NoError(t, err)

	// 验证结构：前 16 字节是 IV，剩余是密文（长度是 16 的倍数）
	require.GreaterOrEqual(t, len(cipherBytes), 32, "cipherBytes should be at least iv(16)+1block(16)")
	require.Equal(t, 0, len(cipherBytes)%16, "cipherBytes length must be multiple of 16")

	// Step 3: 模拟 Go 端 webhook_service.go 的解密
	// decoder := sign.NewEncryptedDecoder(sign.WithAESPassword(secret), sign.WithRawCiphertext())
	decoder := NewEncryptedDecoder(WithAESPassword(secret), WithRawCiphertext())

	// DecodeProtoJSON 模拟：先 Decrypt 得到明文，再 protojson.Unmarshal
	plainTextBytes, err := decoder.Decrypt(cipherBytes)
	require.NoError(t, err)

	// 验证明文与原始 payload 完全一致
	require.Equal(t, plainText, string(plainTextBytes))

	// Step 4: 验证 HTTP body 格式正确（带双引号的 base64 字符串）
	// 这是 Worker 发送给 grpc-gateway 的实际 body
	require.True(t, strings.HasPrefix(httpBody, `"`), "HTTP body should be JSON string with quotes")
	require.True(t, strings.HasSuffix(httpBody, `"`), "HTTP body should end with quote")
}

// TestWebhookE2E_NoDoublePadding 验证不会出现双重 PKCS7 padding
// Worker 端 Web Crypto API 的 AES-CBC 自动 padding，Go 端 AesDecrypt/AesDecryptRaw 移除一次 padding
// 如果 Worker 端手动 padding 会导致明文末尾残留 padding 字节
func TestWebhookE2E_NoDoublePadding(t *testing.T) {
	secret := "test-secret"
	key := GenerateByteKey(secret, 32)

	// 明文长度恰好是 16 的倍数 - 14（这样手动 padding 会添加 14 字节，值为 0x0e）
	// 这模拟了之前 bug 的场景：\x0e 出现在明文中
	plainText := strings.Repeat("a", 16*10-14) // 146 字节，PKCS7 需要添加 14 字节 padding

	// 正确方式：单次 PKCS7 padding（Go AesEncrypt 自动处理）
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	require.NoError(t, err)

	decoder := NewEncryptedDecoder(WithAESKey(key), WithRawCiphertext())
	decrypted, err := decoder.Decrypt(cipherBytes)
	require.NoError(t, err)

	// 解密后的明文应该与原始明文完全一致，没有残留 padding 字节
	require.Equal(t, plainText, string(decrypted))
	require.NotContains(t, string(decrypted), "\x0e", "decrypted text should not contain padding byte 0x0e")
}
