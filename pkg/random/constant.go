/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 10:50:01
 * @FilePath: \go-toolbox\pkg\random\constant.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package random

const (
	// CAPITAL 包含大写字母
	CAPITAL RandType = iota + 1 // 自定义扩展
	// LOWERCASE 包含小写字母
	LOWERCASE
	// SPECIAL 包含特殊字符
	SPECIAL
	// NUMBER 包含数字
	NUMBER
)

const (
	// 自定义数字区间
	DEC_BYTES = "0123456789"
	// 自定义hex区间
	HEX_BYTES = "ABCDEF0123456789"
	// 自定义字符区间
	ALPHA_BYTES  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	LETTER_BYTES = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 自定义十进制

	LETTER_IDX_BITS = 6                      // 6 bits to represent a letter index
	LETTER_IDX_MASK = 1<<LETTER_IDX_BITS - 1 // All 1-bits, as many as letterIdxBits
	LETTER_IDX_MAX  = 31 / LETTER_IDX_BITS   // # of letter indices fitting in 31 bits
)
