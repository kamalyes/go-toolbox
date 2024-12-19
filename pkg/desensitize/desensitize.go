/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-19 08:15:19
 * @FilePath: \go-toolbox\pkg\desensitize\desensitize.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package desensitize

import (
	"reflect"
	"unicode/utf8"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// Desensitize 数据脱敏
func Desensitize(str string, DesensitizeType DesensitizeType, options ...DesensitizeOptions) string {
	// 如果数据为空，则返回空字符串
	if validator.IsEmptyValue(reflect.ValueOf(str)) {
		return str
	}
	opt := NewDesensitizeOptions()
	if len(options) > 0 {
		opt = options[0]
	}
	var newStr string
	switch DesensitizeType {
	case CustomExtension:
		newStr = SensitiveData(str, opt.CustomExtensionStartIndex, opt.CustomExtensionEndIndex)
	case ChineseName:
		newStr = SensitiveData(str, opt.ChineseNameStartIndex, 1+(stringx.Length(str)-2)/2)
	case IDCard:
		newStr = SensitiveData(str, opt.IdCardStartIndex, stringx.Length(str)-4)
	case PhoneNumber:
		newStr = SensitizePhoneNumber(str, opt.PhoneNumberStartIndex, opt.PhoneNumberEndIndex)
	case MobilePhone:
		newStr = SensitiveData(str, opt.MobilePhoneStartIndex, stringx.Length(str)-4)
	case Address:
		newStr = SensitiveData(str, stringx.Length(str)/3, stringx.Length(str)-3)
	case Email:
		newStr = SensitiveData(str, opt.EmailStartIndex, stringx.IndexOf(str, "@")-2)
	case Password:
		newStr = SensitiveData(str, 0, stringx.Length(str))
	case CarLicense:
		newStr = SensitiveData(str, 3, stringx.Length(str)-2)
	case BankCard:
		newStr = SensitizeBankCard(str, opt.IdCardLength)
	case IPV4:
		newStr = SensitizeIpv4(str)
	case IPV6:
		newStr = SensitizeIpv6(str)
	default:
		newStr = str
	}

	return newStr
}

// 通用脱敏函数
func SensitiveData(str string, start, end int) string {
	// 如果数据为空，则返回空字符串
	if validator.IsEmptyValue(reflect.ValueOf(str)) {
		return str
	}

	// 获取字符数量
	charCount := utf8.RuneCountInString(str)

	// 调整起始位置和结束位置
	if start <= 0 || start >= charCount {
		start = 1
	}

	if end <= 0 || end >= charCount {
		end = charCount
	}

	// 如果起始位置和结束位置相等，则视为需要处理整个文本
	if start == end {
		start = 1
		end = charCount
	}

	// 对指定位置的敏感数据进行隐藏处理
	return stringx.Hide(str, start, end)
}

// 手机号脱敏
func SensitizePhoneNumber(str string, start, end int) string {
	// 空判断
	if validator.IsEmptyValue(reflect.ValueOf(str)) {
		return str
	}
	return SensitiveData(stringx.Pad(str, 11), start, end)
}

// 银行卡号脱敏
func SensitizeBankCard(str string, cardLength int) string {
	// 如果卡号为空，则直接返回原卡号
	if validator.IsEmptyValue(reflect.ValueOf(str)) {
		return str
	}

	// 清理卡号中的空格
	cleanCardNo := stringx.CleanEmpty(str)
	// 使用指定的长度填充卡号
	cleanCardNo = stringx.Pad(cleanCardNo, cardLength)
	// 获取卡号长度
	length := stringx.Length(cleanCardNo)
	// 默认间隔为4
	interval := 4

	// 根据卡号长度和特殊情况调整末尾长度
	endLength := length % interval
	if cardLength == 16 {
		endLength = length % 3
	}
	// 计算中间部分长度
	midLength := length - interval - endLength

	// 生成新的卡号切片，初始包含前4位
	newCardNo := []rune(cleanCardNo[:interval])
	for i := 0; i < midLength; i++ {
		// 每隔4位插入一个空格
		if i%interval == 0 {
			newCardNo = append(newCardNo, ' ')
		}
		// 其余字符用*替换
		newCardNo = append(newCardNo, '*')
	}
	// 在中间部分的最末尾插入一个空格
	newCardNo = append(newCardNo, ' ')
	// 添加卡号的最后4位
	lastFour := cleanCardNo[stringx.Length(cleanCardNo)-4:]
	newCardNo = append(newCardNo, []rune(lastFour)...)

	// 返回格式化后的卡号
	return string(newCardNo)
}

// ipv4 脱敏
func SensitizeIpv4(str string) string {
	if validator.IsEmptyValue(reflect.ValueOf(str)) {
		return str
	}
	newIP := stringx.SubBefore(str, ".", false)
	return newIP + ".*.*.*"
}

// ipv6 脱敏
func SensitizeIpv6(str string) string {
	if validator.IsEmptyValue(reflect.ValueOf(str)) {
		return str
	}
	newIP := stringx.SubBefore(str, ":", false)
	return newIP + ":*:*:*:*:*:*:*"
}
