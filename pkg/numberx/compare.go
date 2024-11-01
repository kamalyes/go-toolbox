/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-01 02:25:56
 * @FilePath: \go-toolbox\pkg\numberx\compare.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package numberx

import "errors"

// MinMaxFunc 是用于计算最小值或最大值的函数类型
type MinMaxFunc func(a, b interface{}) interface{}

// MinMax 是一个通用的函数，用于计算最小值或最大值
func MinMax(list []interface{}, f MinMaxFunc) (interface{}, error) {
	if len(list) == 0 {
		return nil, errors.New("list is empty")
	}
	result := list[0]
	for _, v := range list[1:] {
		result = f(result, v)
	}
	return result, nil
}
