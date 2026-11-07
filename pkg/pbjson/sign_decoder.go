/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-29 13:15:55
 * @FilePath: \go-toolbox\pkg\pbjson\sign_decoder.go
 * @Description: sign 解密器的 Proto JSON 载荷解码
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package pbjson

import (
	"errors"
	"fmt"

	"github.com/kamalyes/go-toolbox/pkg/sign"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// DecodeProtoJSONTo 解密Proto JSON到目标
// 原为 sign.EncryptedDecoder 的方法，因携带 protobuf 依赖而拆出本模块，
// 改为以解密器为首参的函数，行为不变
func DecodeProtoJSONTo(d *sign.EncryptedDecoder, ciphertext []byte, target proto.Message) ([]byte, error) {
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

// DecodeProtoJSON 解密Proto JSON
func DecodeProtoJSON[T proto.Message](d *sign.EncryptedDecoder, ciphertext []byte, newPayload func() T) (*sign.Decoded[T], error) {
	if newPayload == nil {
		return nil, errors.New("new payload function is required")
	}

	payload := newPayload()
	plainText, err := DecodeProtoJSONTo(d, ciphertext, payload)
	if err != nil {
		return nil, err
	}
	return &sign.Decoded[T]{
		Ciphertext: append([]byte(nil), ciphertext...),
		Plaintext:  plainText,
		Payload:    payload,
	}, nil
}
