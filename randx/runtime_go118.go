//go:build !go1.19
// +build !go1.19

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 21:44:06
 * @FilePath: \go-toolbox\randx\runtime_go_119.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package randx

import (
	_ "unsafe"
)

func FastRand64() uint64 {
	return (uint64(FastRand()) << 32) | uint64(FastRand())
}

func FastRandu() uint {
	if PtrSize == 8 {
		return uint(FastRand64())
	}
	return uint(FastRand())
}
