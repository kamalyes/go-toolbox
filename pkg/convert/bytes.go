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
	"strings"
	"unsafe"
)

// byteBinTable 预生成的字节二进制字符串查找表
// 索引 0..255 对应 "00000000".."11111111"
// 避免每次调用 fmt.Sprintf("%08b", b) 解析格式字符串
var byteBinTable [256]string

// hexUpperTable 大写十六进制编码查找表
// 每个字节对应 2 个大写 hex 字符，避免 hex.EncodeToString + strings.ToUpper 两次分配
var hexUpperTable [256][2]byte

func init() {
	// 初始化二进制查找表
	const digits = "01"
	for i := 0; i < 256; i++ {
		bin := []byte("00000000")
		for j := 7; j >= 0; j-- {
			if i&(1<<uint(7-j)) != 0 {
				bin[j] = digits[1]
			} else {
				bin[j] = digits[0]
			}
		}
		byteBinTable[i] = string(bin)
	}

	// 初始化大写 hex 查找表
	const hexChars = "0123456789ABCDEF"
	for i := 0; i < 256; i++ {
		hexUpperTable[i] = [2]byte{hexChars[i>>4], hexChars[i&0x0F]}
	}
}

// BytesToHex 将字节数组转换为大写十六进制字符串
func BytesToHex(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	// 直接构造 []byte 转 string，避免中间字符串到 ToUpper 的拷贝
	buf := make([]byte, len(data)*2)
	for i, b := range data {
		buf[i*2] = hexUpperTable[b][0]
		buf[i*2+1] = hexUpperTable[b][1]
	}
	return string(buf)
}

// BytesToBCC 计算字节数组的 BCC（块校验字符）
func BytesToBCC(data []byte) byte {
	var bcc byte
	for _, b := range data {
		bcc ^= b
	}
	return bcc
}

// ByteToBinStr 将单个字节转换为 8 位二进制字符串
func ByteToBinStr(b byte) string {
	return byteBinTable[b]
}

// BytesToBinStr 将字节数组转换为二进制字符串
func BytesToBinStr(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow(len(bs) * 8)
	for _, v := range bs {
		sb.WriteString(byteBinTable[v])
	}
	return sb.String()
}

// BytesToBinStrWithSplit 将字节数组转换为二进制字符串，并添加分隔符
func BytesToBinStrWithSplit(bs []byte, split string) string {
	if len(bs) == 0 {
		return ""
	}
	var sb strings.Builder
	// 预估容量：每个字节 8 字符 + 分隔符长度（最后一个不加分隔符）
	sb.Grow(len(bs)*(8+len(split)) - len(split))
	for i, v := range bs {
		if i > 0 {
			sb.WriteString(split)
		}
		sb.WriteString(byteBinTable[v])
	}
	return sb.String()
}

// SliceByteToString 将字节切片转换为字符串（零拷贝）
//
// 注意：返回的 string 与 b 共享底层数组，修改 b 会影响返回值
// 仅适用于只读场景或临时字符串
func SliceByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
