/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-28 11:54:16
 * @FilePath: \go-toolbox\uuid\uuid.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package uuid

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kamalyes/go-toolbox/convert"
)

const (
	NumSource = "0123456789"
	HexSource = "ABCDEF0123456789"
	StrSource = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// RandomStr 随机一个字符串
func RandomStr(length int) string {
	var sb strings.Builder
	if length > 0 {
		for i := 0; i < length; i++ {
			sb.WriteString(string(StrSource[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(StrSource))]))
		}
	}
	return sb.String()
}

// RandomNum 随机一个数字字符串
func RandomNum(length int) string {
	var sb strings.Builder
	if length > 0 {
		for i := 0; i < length; i++ {
			sb.WriteString(string(NumSource[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(NumSource))]))
		}
	}
	return sb.String()
}

// RandomHex 随机一个hex字符串
func RandomHex(bytesLen int) string {
	var sb strings.Builder
	for i := 0; i < bytesLen<<1; i++ {
		sb.WriteString(string(HexSource[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(HexSource))]))
	}
	return sb.String()
}

// SameSubStr 创建多少个相同 子字符串的字符串
func SameSubStr(subStr string, repeat int) string {
	var sb strings.Builder
	for i := 0; i < repeat; i++ {
		sb.WriteString(subStr)
	}
	return sb.String()
}

// UUID 生成uuid
func UUID() string {
	id := uuid.New()
	return strings.ReplaceAll(id.String(), "-", "")
}

// UniqueID 根据指定字段生成 uuid
func UniqueID(fields ...interface{}) string {
	if len(fields) == 0 {
		return UUID()
	}
	var buf strings.Builder
	for i := range fields {
		field := fields[i]
		buf.WriteString(convert.AsString(field))
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
