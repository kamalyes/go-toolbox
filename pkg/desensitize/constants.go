/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 22:19:05
 * @FilePath: \go-toolbox\pkg\desensitize\constants.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package desensitize

const (
	CustomExtension DesensitizeType = iota + 1 // 自定义扩展
	ChineseName                                // 中文名称
	IDCard                                     // 身份证号
	PhoneNumber                                // 手机号码
	MobilePhone                                // 移动电话号码
	Address                                    // 地址
	Email                                      // 邮箱
	Password                                   // 密码
	CarLicense                                 // 车牌号：油车、电车
	BankCard                                   // 银行卡号
	IPV4                                       // ipv4
	IPV6                                       // ipv6
)
