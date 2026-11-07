/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 13:17:52
 * @FilePath: \go-toolbox\pkg\sign\bcryptx\bcrypt.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

// Package bcryptx 提供 bcrypt 口令哈希
//
// 本包自 pkg/sign 拆出并独立成 module：bcrypt 只存在于 golang.org/x/crypto，
// 标准库无等价实现，而 go-toolbox 核心不携带任何三方依赖
package bcryptx

import "golang.org/x/crypto/bcrypt"

// GenerateFromPassword 使用 bcrypt 对数据进行哈希
// cost: 可选的哈希成本因子，不传则使用默认值 bcrypt.DefaultCost(10)
func GenerateFromPassword(data []byte, cost ...int) ([]byte, error) {
	c := bcrypt.DefaultCost
	if len(cost) > 0 && cost[0] > 0 {
		c = cost[0]
	}
	return bcrypt.GenerateFromPassword(data, c)
}

// CompareHashAndPassword 校验数据与 bcrypt 哈希是否匹配
// 匹配返回 nil，不匹配返回 bcrypt.ErrMismatchedHashAndPassword
func CompareHashAndPassword(hashed, data []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, data)
}
