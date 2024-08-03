//go:build !go1.20
// +build !go1.20

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 21:51:05
 * @FilePath: \go-toolbox\convert\convert_go119.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"unsafe"
)

// S2B StringToBytes converts string to byte slice without a memory allocation.
// Ref: gin
func S2B(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// B2S BytesToString
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
