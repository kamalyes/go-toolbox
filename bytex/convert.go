/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 18:31:24
 * @FilePath: \go-toolbox\bytex\convert.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package bytex

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// BytesToHex 字节数组转hex
// []byte{0x55, 0xAA} 被转成 55AA
func BytesToHex(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

// HexToBytes 将hex 字符串转成 byte数组
// AABBCC 转成字节数组 []byte{0xAA, 0xBB, 0xCC}
func HexToBytes(hexStr string) []byte {
	decodeString, _ := hex.DecodeString(hexStr)
	return decodeString
}

// HexBCC 计算BCC校验码
func HexBCC(hexStr string) string {
	hexToBytes := HexToBytes(hexStr)
	length := len(hexToBytes)
	if length < 1 {
		return ""
	}
	bcc := hexToBytes[0]
	if length > 1 {
		for i := 1; i < length; i++ {
			bcc = bcc ^ hexToBytes[i]
		}
	}
	return BytesToHex([]byte{bcc & 0xFF})
}

// BytesBCC 计算 BCC
func BytesBCC(bytes []byte) byte {
	bcc := bytes[0]
	if len(bytes) > 1 {
		for i := 1; i < len(bytes); i++ {
			bcc = bcc ^ bytes[i]
		}
	}
	return bcc & 0xFF
}

// DecToHex 十进进制转16进制
func DecToHex(n uint64) string {
	s := strconv.FormatUint(n, 16)
	s = strings.ToUpper(s)
	length := len(s)
	if length%2 == 1 {
		s = "0" + s
	}
	return s
}

// HexToDec 十六进制转10进制
func HexToDec(h string) uint64 {
	n, err := strconv.ParseUint(h, 16, 64)
	if err != nil {
		return 0
	}
	return n
}

// DecToBin 十进制转二进制
func DecToBin(n uint64) string {
	s := strconv.FormatUint(n, 2)
	length := len(s)
	mod := length % 8
	if mod != 0 {
		prefixNum := 8 - mod
		var sb strings.Builder
		for i := 0; i < prefixNum; i++ {
			sb.WriteString("0")
		}
		s = sb.String() + s
	}
	return s
}

// HexToBin 16进制转二进制
func HexToBin(h string) string {
	n, err := strconv.ParseUint(h, 16, 64)
	if err != nil {
		return ""
	}
	return DecToBin(n)
}

// ByteToBinStr 将byte 以8个bit位的形式展示
func ByteToBinStr(b byte) string {
	return fmt.Sprintf("%08b", b)
}

// BytesToBinStr 将byte数组转成8个bit位一组的字符串
func BytesToBinStr(bs []byte) string {
	if len(bs) <= 0 {
		return ""
	}
	buf := bytes.NewBuffer([]byte{})
	for _, v := range bs {
		buf.WriteString(fmt.Sprintf("%08b", v))
	}
	return buf.String()
}

// BytesToBinStrWithSplit 将byte数组转8个bit一组的字符串并且带分割符
func BytesToBinStrWithSplit(bs []byte, split string) string {
	length := len(bs)
	if length <= 0 {
		return ""
	}
	buf := bytes.NewBuffer([]byte{})
	for i := 0; i < length-1; i++ {
		v := bs[i]
		buf.WriteString(fmt.Sprintf("%08b", v))
		buf.WriteString(split)
	}
	buf.WriteString(fmt.Sprintf("%08b", bs[length-1]))
	return buf.String()
}

// HexSuffixZero hex 后补位
func HexSuffixZero(hex string, byteSize int) string {
	data1 := HexToBytes(hex)
	data2 := make([]byte, byteSize)
	copy(data2, data1)
	return BytesToHex(data2)
}

func HexPrefixZero(hex string, byteSize int) string {
	data1 := HexToBytes(hex)
	data2 := append(make([]byte, byteSize-len(data1)), data1...)
	return BytesToHex(data2)
}

// GBKSuffixZero GBK 编码按字节右补0
func GBKSuffixZero(gbkStr string, byteSize int) string {
	data1, _ := io.ReadAll(transform.NewReader(bytes.NewReader([]byte(gbkStr)), simplifiedchinese.GBK.NewEncoder()))
	data2 := make([]byte, byteSize)
	copy(data2, data1)
	return BytesToHex(data2)
}

// GBKSuffixSpace 编码按字节右补空格
func GBKSuffixSpace(chinese string, byteSize int) (hex string) {
	data1, _ := io.ReadAll(transform.NewReader(bytes.NewReader([]byte(chinese)), simplifiedchinese.GBK.NewEncoder()))
	data2 := make([]byte, byteSize)
	copy(data2, data1)
	for i := len(data1); i < len(data2); i++ {
		data2[i] = 0x20
	}
	return string(data2)
}

// HexReverse 字节颠倒
func HexReverse(hex string) string {
	toBytes := HexToBytes(hex)
	length := len(toBytes)
	if length <= 1 {
		return hex
	}
	for i := range toBytes {
		a := toBytes[i]
		toBytes[i] = toBytes[length-1-i]
		toBytes[length-1-i] = a
		if i == length/2 {
			break
		}
	}
	return BytesToHex(toBytes)
}
