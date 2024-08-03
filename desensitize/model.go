/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 22:18:46
 * @FilePath: \go-toolbox\desensitize\model.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package desensitize

type DesensitizeType int

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
	MobilePhoneStartIndex     int
	EmailStartIndex           int
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
		MobilePhoneStartIndex:     3,
		EmailStartIndex:           1,
		AddressLength:             8,
	}
}
