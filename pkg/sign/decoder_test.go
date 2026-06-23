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
