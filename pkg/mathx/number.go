/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-13 23:20:55
 * @FilePath: \go-toolbox\pkg\mathx\number.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package mathx

import (
	"bytes"
	"fmt"
	"math"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/types"
)

// Decimals 转换为包含小数点后指定位数的字符串
func Decimals[T types.Numerical](num T, digit int) string {
	// 计算除数
	divisor := T(math.Pow10(digit))
	// 将整数转换为浮点数，然后除以除数
	flt := float64(num) / float64(divisor)
	// 格式化为字符串，保留小数点后指定位数
	result := fmt.Sprintf("%.*f", digit, flt)
	return result
}

// AtLeast 返回 x 和 lower 中的较大值。
// 参数:
// x - 要比较的第一个数值
// lower - 要比较的第二个数值（下限）
// 返回值:
// 返回 x 和 lower 中的较大值。
func AtLeast[T types.Numerical](x, lower T) T {
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
func AtMost[T types.Numerical](x, upper T) T {
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
func Between[T types.Numerical](x, lower, upper T) T {
	if x < lower {
		return lower
	}
	if x > upper {
		return upper
	}
	return x
}

// LongestCommonPrefix 返回两个字符串的最长公共前缀的长度。
// 参数:
// a - 第一个字符串
// b - 第二个字符串
// 返回值:
// 返回两个字符串的最长公共前缀的长度。
func LongestCommonPrefix(a, b string) int {
	// 计算两个字符串的最小长度
	maxLength := AtMost(len(a), len(b))

	// 遍历两个字符串，比较字符
	for i := 0; i < maxLength; i++ {
		if a[i] != b[i] {
			return i // 返回公共前缀的长度
		}
	}
	return maxLength // 如果完全相同，返回最小长度
}

// CountPathSegments 计算路径中指定前缀的参数数量，默认为 ":" 和 "*"。
func CountPathSegments(path string, prefixes ...string) int {
	// 如果没有提供前缀，则使用默认前缀
	if len(prefixes) == 0 {
		prefixes = []string{":", "*"}
	}

	count := 0
	// 遍历所有前缀并计算出现的次数
	for _, prefix := range prefixes {
		count += bytes.Count(convert.StringToSliceByte(path), []byte(prefix))
	}
	return count
}

// ZeroValue 返回类型 T 的零值
func ZeroValue[T any]() T {
	var t T
	return t
}

// EqualSlices 比较两个切片是否相等，支持任意类型
func EqualSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
