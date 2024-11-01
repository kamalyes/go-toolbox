//go:build go1.20
// +build go1.20

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 21:55:06
 * @FilePath: \go-toolbox\pkg\convert\runtime_go_119.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package convert

import (
	"unsafe"
)

// S2B converts string to byte slice without a memory allocation.
// Ref: https://github.com/golang/go/issues/53003#issuecomment-1140276077
func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// B2S converts byte slice to string without a memory allocation.
// Slower: unsafe.String(unsafe.SliceData(b), len(b))
// strings.Clone(): unsafe.String(&b[0], len(b))
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
