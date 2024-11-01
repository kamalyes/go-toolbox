/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 09:16:55
 * @FilePath: \go-toolbox\pkg\random\model.go
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
