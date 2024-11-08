/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 01:15:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 12:22:59
 * @FilePath: \go-toolbox\pkg\convert\radix.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// HexToBytes 将十六进制字符串转换为字节数组。
// 如果十六进制字符串的长度为奇数或转换失败，则返回错误。
func HexToBytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		return nil, errors.New("十六进制字符串长度必须为偶数")
	}
	// 预分配字节切片以避免重新分配
	bytes := make([]byte, len(hexStr)/2)
	_, err := hex.Decode(bytes, []byte(hexStr))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// HexToDec 将十六进制字符串转换为十进制数。
// 如果转换失败，则返回错误。
func HexToDec(h string) (uint64, error) {
	return strconv.ParseUint(h, 16, 64)
}

// HexToBin 将十六进制字符串转换为二进制字符串。
// 如果转换失败，则返回错误。
func HexToBin(h string) (string, error) {
	n, err := HexToDec(h)
	if err != nil {
		return "", err
	}
	return DecToBin(n), nil
}

// HexToBCC 计算十六进制字符串的 BCC
func HexToBCC(hexStr string) (string, error) {
	bytes, err := HexToBytes(hexStr)
	if err != nil {
		return "", err
	}
	bcc := BytesToBCC(bytes)
	return hex.EncodeToString([]byte{bcc}), nil
}

// DecToHex 十进制转为十六进制字符串
func DecToHex(n uint64) string {
	s := strconv.FormatUint(n, 16)
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return strings.ToUpper(s)
}

// DecToBin 将十进制数转换为二进制字符串，并补齐到8位。
func DecToBin(n uint64) string {
	return fmt.Sprintf("%08b", n)
}
