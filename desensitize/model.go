/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-02 11:30:59
 * @FilePath: \go-toolbox\desensitize\model.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package desensitize

type DesensitizeType int

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

// DesensitizeOptions 脱敏选项
type DesensitizeOptions struct {
	CustomExtensionStartIndex int
	CustomExtensionEndIndex   int
	ChineseNameStartIndex     int
	ChineseNameEndIndex       int
	IdCardStartIndex          int
	IdCardEndIndex            int
	IdCardLength              int
	PhoneNumberStartIndex     int
	PhoneNumberEndIndex       int
	AddressLength             int
}

// NewDesensitizeOptions 创建带有默认值的 DesensitizeOptions
func NewDesensitizeOptions() DesensitizeOptions {
	return DesensitizeOptions{
		CustomExtensionStartIndex: 1,
		CustomExtensionEndIndex:   1,
		ChineseNameStartIndex:     1,
		ChineseNameEndIndex:       1,
		IdCardStartIndex:          7,
		IdCardEndIndex:            13,
		IdCardLength:              19,
		PhoneNumberStartIndex:     3,
		PhoneNumberEndIndex:       6,
		AddressLength:             8,
	}
}
