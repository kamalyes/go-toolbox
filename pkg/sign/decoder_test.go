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
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

type decodeTestPayload struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func TestDecodeJSON(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `{"id":"job_001","status":"success"}`
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	decoder := NewEncryptedDecoder(WithAESKey(key))
	decoded, err := DecodeJSON[decodeTestPayload](decoder, []byte(ciphertext))
	require.NoError(t, err)
	require.Equal(t, []byte(ciphertext), decoded.Ciphertext)
	require.Equal(t, []byte(plainText), decoded.Plaintext)
	require.Equal(t, "job_001", decoded.Payload.ID)
	require.Equal(t, "success", decoded.Payload.Status)
}

func TestDecodeJSONTo(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `{"id":"job_002","status":"running"}`
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	var payload decodeTestPayload
	decoder := NewEncryptedDecoder(WithAESKey(key))
	decodedPlainText, err := decoder.DecodeJSONTo([]byte(ciphertext), &payload)
	require.NoError(t, err)
	require.Equal(t, []byte(plainText), decodedPlainText)
	require.Equal(t, "job_002", payload.ID)
	require.Equal(t, "running", payload.Status)
}

func TestDecodeProtoJSON(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `{"id":"job_001","status":"success"}`
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	decoder := NewEncryptedDecoder(WithAESKey(key))
	decoded, err := DecodeProtoJSON(decoder, []byte(ciphertext), func() *structpb.Struct {
		return &structpb.Struct{}
	})
	require.NoError(t, err)
	require.Equal(t, "job_001", decoded.Payload.GetFields()["id"].GetStringValue())
	require.Equal(t, "success", decoded.Payload.GetFields()["status"].GetStringValue())
}

func TestDecodeProtoJSONTo(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `{"id":"job_004","status":"proto"}`
	ciphertext, err := AesEncrypt(plainText, key)
	require.NoError(t, err)

	payload := &structpb.Struct{}
	decoder := NewEncryptedDecoder(WithAESKey(key))
	decodedPlainText, err := decoder.DecodeProtoJSONTo([]byte(ciphertext), payload)
	require.NoError(t, err)
	require.Equal(t, []byte(plainText), decodedPlainText)
	require.Equal(t, "job_004", payload.GetFields()["id"].GetStringValue())
	require.Equal(t, "proto", payload.GetFields()["status"].GetStringValue())
}

func TestEncryptedDecoderRejectsMissingKey(t *testing.T) {
	_, err := NewEncryptedDecoder().Decrypt([]byte("ciphertext"))
	require.ErrorIs(t, err, ErrMissingAESKey)
}
