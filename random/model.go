/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 22:22:51
 * @FilePath: \go-toolbox\random\model.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package random

type RandType int

// RandOptions
type RandOptions struct {
}

// NewRandOptions 创建带有默认值的 RandOptions
func NewRandOptions() RandOptions {
	return RandOptions{}
}
