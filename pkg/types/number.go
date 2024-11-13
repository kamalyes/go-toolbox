/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-13 15:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 23:20:55
 * @FilePath: \go-toolbox\pkg\types\number.go
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

// MinMaxFunc 是用于计算最小值或最大值的函数类型，接收两个interface{}类型的参数，返回一个interface{}类型的结果。
type MinMaxFunc[T any] func(a, b T) T
