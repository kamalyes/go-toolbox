/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 17:50:03
 * @FilePath: \go-toolbox\numberx\compare.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package numberx

import "errors"

// MinMaxFunc 是一个函数类型，用于计算最小值或最大值
type MinMaxFunc func(a, b int) int

// Min 返回整数中的最小值
func Min(intList []int) (int, error) {
	return minMax(intList, func(a, b int) int {
		if a < b {
			return a
		}
		return b
	})
}

// Max 返回整数中的最大值
func Max(intList []int) (int, error) {
	return minMax(intList, func(a, b int) int {
		if a > b {
			return a
		}
		return b
	})
}

// minMax 是一个通用的函数，用于计算最小值或最大值
func minMax(intList []int, f MinMaxFunc) (int, error) {
	if len(intList) == 0 {
		return 0, errors.New("intList is empty")
	}
	result := intList[0]
	for _, v := range intList[1:] {
		result = f(result, v)
	}
	return result, nil
}
