/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 17:26:07
 * @FilePath: \go-toolbox\desensitize\model.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package desensitize

type DesensitizeType int

const (
	UserId      DesensitizeType = 1  // 用户ID
	ChineseName DesensitizeType = 2  // 中文名称
	IDCard      DesensitizeType = 3  // 身份证号
	FixedPhone  DesensitizeType = 4  // 固定电话号码
	MobilePhone DesensitizeType = 5  // 移动电话号码
	Address     DesensitizeType = 6  // 地址
	Email       DesensitizeType = 7  // 邮箱
	Password    DesensitizeType = 8  // 密码
	CarLicense  DesensitizeType = 9  // 车牌号：油车、电车
	BankCard    DesensitizeType = 10 // 银行卡号
	IPV4        DesensitizeType = 11 // ipv4
	IPV6        DesensitizeType = 12 // ipv6
)
