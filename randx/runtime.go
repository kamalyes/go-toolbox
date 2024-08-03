/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-03 21:32:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 21:53:30
 * @FilePath: \go-toolbox\randx\runtime.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package randx

import (
	_ "unsafe"
)

// FastRand 随机数
//
//go:linkname FastRand runtime.fastrand
func FastRand() uint32

// FastRandn 等同于 FastRandn() % n, 但更快
// See https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
//
//go:linkname FastRandn runtime.fastrandn
func FastRandn(n uint32) uint32
