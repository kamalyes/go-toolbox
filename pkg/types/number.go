/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-13 15:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 15:55:18
 * @FilePath: \go-toolbox\internal\types\norm.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package types

// Uint
type Uint interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Int
type Int interface {
	int | int8 | int16 | int32 | int64
}

// Float
type Float interface {
	float32 | float64
}

// Numerical 是一个接口，表示一系列数值类型，包括有符号和无符号的整数以及浮点数。
type Numerical interface {
	Int | Uint | Float
}
