/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 13:15:55
 * @FilePath: \go-toolbox\pkg\sign\decoder.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ErrMissingAESKey 缺少AES密钥
var (
	ErrMissingAESKey     = errors.New("aes key is required")
	ErrMissingCiphertext = errors.New("ciphertext is required")
)

// EncryptedDecoder 解密器
type EncryptedDecoder struct {
	aesKey        []byte
	rawCiphertext bool
}

// DecodeOption 解密器选项
type DecodeOption func(*EncryptedDecoder)

type Decoded[T any] struct {
	Ciphertext []byte
	Plaintext  []byte
	Payload    T
}

// NewEncryptedDecoder 创建解密器
func NewEncryptedDecoder(opts ...DecodeOption) *EncryptedDecoder {
	d := &EncryptedDecoder{}
	for _, opt := range opts {
		if opt != nil {
			opt(d)
		}
	}
	return d
}

// WithAESKey 设置AES密钥
func WithAESKey(key []byte) DecodeOption {
	return func(d *EncryptedDecoder) {
		d.aesKey = append([]byte(nil), key...)
	}
}

// WithAESPassword 设置AES密码
func WithAESPassword(password string) DecodeOption {
	return func(d *EncryptedDecoder) {
		d.aesKey = GenerateByteKey(password, 32)
	}
}

// WithRawCiphertext 设置密文为原始字节模式
// 启用后 Decrypt 直接接收 iv+ciphertext 原始字节，不再做 base64 解码
// 适用于 grpc-gateway 等 JSON marshaler 已对 bytes 字段做 base64 解码的场景
func WithRawCiphertext() DecodeOption {
	return func(d *EncryptedDecoder) {
		d.rawCiphertext = true
	}
}

// Decrypt 解密
func (d *EncryptedDecoder) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) == 0 {
		return nil, ErrMissingCiphertext
	}
	if d == nil || len(d.aesKey) == 0 {
		return nil, ErrMissingAESKey
	}

	if d.rawCiphertext {
		plainText, err := AesDecryptRaw(ciphertext, d.aesKey)
		if err != nil {
			return nil, fmt.Errorf("decrypt ciphertext: %w", err)
		}
		return []byte(plainText), nil
	}

	plainText, err := AesDecrypt(strings.TrimSpace(string(ciphertext)), d.aesKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt ciphertext: %w", err)
	}
	return []byte(plainText), nil
}

// DecodeJSONTo 解密JSON到目标
func (d *EncryptedDecoder) DecodeJSONTo(ciphertext []byte, target any) ([]byte, error) {
	plainText, err := d.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, errors.New("json target is required")
	}
	if err := json.Unmarshal(plainText, target); err != nil {
		return nil, fmt.Errorf("decode json payload: %w", err)
	}
	return plainText, nil
}

// DecodeProtoJSONTo 解密Proto JSON到目标
func (d *EncryptedDecoder) DecodeProtoJSONTo(ciphertext []byte, target proto.Message) ([]byte, error) {
	plainText, err := d.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, errors.New("proto target is required")
	}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(plainText, target); err != nil {
		return nil, fmt.Errorf("decode proto json payload: %w", err)
	}
	return plainText, nil
}

// DecodeJSON 解密JSON
func DecodeJSON[T any](d *EncryptedDecoder, ciphertext []byte) (*Decoded[T], error) {
	var payload T
	plainText, err := d.DecodeJSONTo(ciphertext, &payload)
	if err != nil {
		return nil, err
	}
	return &Decoded[T]{
		Ciphertext: append([]byte(nil), ciphertext...),
		Plaintext:  plainText,
		Payload:    payload,
	}, nil
}

// DecodeProtoJSON 解密Proto JSON
func DecodeProtoJSON[T proto.Message](d *EncryptedDecoder, ciphertext []byte, newPayload func() T) (*Decoded[T], error) {
	if newPayload == nil {
		return nil, errors.New("new payload function is required")
	}

	payload := newPayload()
	plainText, err := d.DecodeProtoJSONTo(ciphertext, payload)
	if err != nil {
		return nil, err
	}
	return &Decoded[T]{
		Ciphertext: append([]byte(nil), ciphertext...),
		Plaintext:  plainText,
		Payload:    payload,
	}, nil
}
