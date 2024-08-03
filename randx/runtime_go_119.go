//go:build go1.19
// +build go1.19

/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 21:55:06
 * @FilePath: \go-toolbox\randx\runtime_go_119.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package randx

import (
	_ "unsafe"
)

//go:linkname FastRand64 runtime.fastrand64
func FastRand64() uint64

//go:linkname FastRandu runtime.fastrandu
func FastRandu() uint
