/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 12:17:56
 * @FilePath: \go-toolbox\pkg\pbjson\sign_decoder_test.go
 * @Description: sign 解密器 Proto JSON 载荷解码测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package pbjson

import (
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/sign"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestDecodeProtoJSON(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `"buildId,status,artifact.url,metadata.commit"`
	ciphertext, err := sign.AesEncrypt(plainText, key)
	require.NoError(t, err)

	decoder := sign.NewEncryptedDecoder(sign.WithAESKey(key))
	decoded, err := DecodeProtoJSON(decoder, []byte(ciphertext), func() *fieldmaskpb.FieldMask {
		return &fieldmaskpb.FieldMask{}
	})
	require.NoError(t, err)
	require.Equal(t, []string{"build_id", "status", "artifact.url", "metadata.commit"}, decoded.Payload.Paths)
}
func TestDecodeProtoJSONTo(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plainText := `"buildId,status,steps.name"`
	ciphertext, err := sign.AesEncrypt(plainText, key)
	require.NoError(t, err)

	payload := &fieldmaskpb.FieldMask{}
	decoder := sign.NewEncryptedDecoder(sign.WithAESKey(key))
	decodedPlainText, err := DecodeProtoJSONTo(decoder, []byte(ciphertext), payload)
	require.NoError(t, err)
	require.Equal(t, []byte(plainText), decodedPlainText)
	require.Equal(t, []string{"build_id", "status", "steps.name"}, payload.Paths)
}
