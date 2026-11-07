/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 11:55:16
 * @FilePath: \go-toolbox\pkg\uuid\uuid.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package uuid

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/kamalyes/go-toolbox/pkg/convert"
)

// SameSubStr 创建多少个相同 子字符串的字符串
func SameSubStr(subStr string, repeat int) string {
	var sb strings.Builder
	for i := 0; i < repeat; i++ {
		sb.WriteString(subStr)
	}
	return sb.String()
}

// UUID 生成uuid，输出 RFC 4122 v4 的 32 位小写十六进制，不含连字符
func UUID() string {
	var b [16]byte
	// crypto/rand.Read 自 Go 1.24 起不再返回错误，熵源不可用时直接 panic
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // 版本位固定为 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 位固定为 RFC 4122
	return hex.EncodeToString(b[:])
}

// UniqueID 根据指定字段生成 uuid
func UniqueID(fields ...interface{}) string {
	if len(fields) == 0 {
		return UUID()
	}
	var buf strings.Builder
	for i := range fields {
		field := fields[i]
		buf.WriteString(convert.MustString(field))
	}
	s := strings.TrimSpace(buf.String())
	if s == "" {
		return UUID()
	}
	return Md5(s)
}

// Md5 md5加密
func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
