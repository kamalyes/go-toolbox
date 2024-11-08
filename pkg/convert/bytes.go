/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 01:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 01:05:55
 * @FilePath: \go-toolbox\pkg\convert\bytes.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"unsafe"
)

// BytesToHex 将字节数组转换为十六进制字符串
func BytesToHex(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

// BytesBCC 计算字节数组的 BCC（块校验字符）
func BytesToBCC(data []byte) byte {
	var bcc byte
	for _, b := range data {
		bcc ^= b
	}
	return bcc
}

// ByteToBinStr 将单个字节转换为二进制字符串
func ByteToBinStr(b byte) string {
	return fmt.Sprintf("%08b", b)
}

// BytesToBinStr 将字节数组转换为二进制字符串
func BytesToBinStr(bs []byte) string {
	var buf bytes.Buffer
	buf.Grow(len(bs) * 8) // 预分配内存，减少内存开销
	for _, v := range bs {
		buf.WriteString(ByteToBinStr(v))
	}
	return buf.String()
}

// BytesToBinStrWithSplit 将字节数组转换为二进制字符串，并添加分隔符
func BytesToBinStrWithSplit(bs []byte, split string) string {
	if len(bs) == 0 {
		return ""
	}
	var buf bytes.Buffer
	buf.Grow(len(bs)*(8+len(split)) - len(split)) // 预分配内存，考虑分隔符的长度
	for i, v := range bs {
		if i > 0 {
			buf.WriteString(split)
		}
		buf.WriteString(ByteToBinStr(v))
	}
	return buf.String()
}

// SliceByteToString 将字节切片转换为字符串
func SliceByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToSliceByte 将字符串转换为字节切片
func StringToSliceByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
