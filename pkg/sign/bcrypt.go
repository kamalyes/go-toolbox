/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-10 21:51:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-08 13:17:52
 * @FilePath: \go-toolbox\pkg\sign\bcrypt.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package sign

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
