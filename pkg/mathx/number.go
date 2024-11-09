/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-10 00:18:15
 * @FilePath: \go-toolbox\pkg\mathx\number.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package mathx

import "errors"

// MaxInt 返回a和b中较大的一个。
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// MinInt 返回a和b中较小的一个。
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MinMaxFunc 是用于计算最小值或最大值的函数类型，接收两个interface{}类型的参数，返回一个interface{}类型的结果。
type MinMaxFunc func(a, b interface{}) interface{}

// MinMax 是一个通用的函数，用于计算列表中元素的最小值或最大值。
// 它接收一个interface{}类型的切片和一个MinMaxFunc类型的函数，
// 根据提供的函数决定是计算最小值还是最大值。
// 如果列表为空，则返回错误。
func MinMax(list []interface{}, f MinMaxFunc) (interface{}, error) {
	// 检查列表是否为空
	if len(list) == 0 {
		return nil, errors.New("列表为空") // 优化注释，使其更简洁明了
	}

	// 初始化结果为列表的第一个元素
	result := list[0]

	// 遍历列表中的其余元素，使用提供的函数更新结果
	for _, v := range list[1:] {
		result = f(result, v)
	}

	// 返回最终的结果和nil错误（表示无错误）
	return result, nil
}

// Numerical 是一个接口，表示一系列数值类型，包括有符号和无符号的整数以及浮点数。
type Numerical interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// AtLeast 返回 x 和 lower 中的较大值。
// 参数:
// x - 要比较的第一个数值
// lower - 要比较的第二个数值（下限）
// 返回值:
// 返回 x 和 lower 中的较大值。
func AtLeast[T Numerical](x, lower T) T {
	if x < lower {
		return lower
	}
	return x
}

// AtMost 返回 x 和 upper 中的较小值。
// 参数:
// x - 要比较的第一个数值
// upper - 要比较的第二个数值（上限）
// 返回值:
// 返回 x 和 upper 中的较小值。
func AtMost[T Numerical](x, upper T) T {
	if x > upper {
		return upper
	}
	return x
}

// Between 将 x 的值限制在 [lower, upper] 范围内。
// 如果 x 小于 lower，则返回 lower；
// 如果 x 大于 upper，则返回 upper；
// 否则，返回 x 本身。
// 参数:
// x - 要限制的数值
// lower - 范围的下限
// upper - 范围的上限
// 返回值:
// 返回 x 被限制在 [lower, upper] 范围内的值。
func Between[T Numerical](x, lower, upper T) T {
	if x < lower {
		return lower
	}
	if x > upper {
		return upper
	}
	return x
}
