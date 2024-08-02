/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-03 00:33:31
 * @FilePath: \go-toolbox\desensitize\desensitize.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package desensitize

import (
	"unicode/utf8"

	"github.com/kamalyes/go-toolbox/stringx"
)

// Desensitize 数据脱敏
func Desensitize(str string, DesensitizeType DesensitizeType, options ...DesensitizeOptions) string {
	// 如果数据为空，则返回空字符串
	if stringx.IsEmpty(str) {
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
		newStr = SensitiveData(str, opt.ChineseNameStartIndex, opt.ChineseNameEndIndex)
	case IDCard:
		newStr = SensitiveData(str, opt.IdCardStartIndex, opt.IdCardEndIndex)
	case PhoneNumber:
		newStr = phoneNumber(str, opt.PhoneNumberStartIndex, opt.PhoneNumberEndIndex)
	case MobilePhone:
		newStr = SensitiveData(str, 3, stringx.Length(str)-4)
	case Address:
		newStr = SensitiveData(str, stringx.Length(str)-opt.AddressLength, stringx.Length(str))
	case Email:
		newStr = SensitiveData(str, 1, stringx.IndexOf(str, "@"))
	case Password:
		newStr = SensitiveData(str, 1, stringx.Length(str))
	case CarLicense:
		newStr = carLicense(str)
	case BankCard:
		newStr = bankCard(str, opt.IdCardLength)
	case IPV4:
		newStr = ipv4(str)
	case IPV6:
		newStr = ipv6(str)
	default:
		newStr = str
	}

	return newStr
}

// 通用脱敏函数
func SensitiveData(data string, start, end int) string {
	// 如果数据为空，则返回空字符串
	if stringx.IsEmpty(data) {
		return data
	}

	// 获取字符数量
	charCount := utf8.RuneCountInString(data)

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
	return stringx.Hide(data, start, end)
}

// 手机号脱敏
func phoneNumber(phone string, start, end int) string {
	// 空判断
	if stringx.IsEmpty(phone) {
		return phone
	}
	phone = stringx.Pad(phone, 11)
	return SensitiveData(phone, start, end)
}

// 车牌号脱敏
func carLicense(carNo string) string {
	// 空判断
	if stringx.IsEmpty(carNo) {
		return carNo
	}

	newCarNo := carNo
	// 普通车牌
	if stringx.Length(carNo) == 7 {
		newCarNo = stringx.Hide(carNo, 3, 6)
	} else if stringx.Length(carNo) == 8 { // 新能源
		newCarNo = stringx.Hide(carNo, 3, 7)
	}
	return newCarNo
}

// 银行卡号脱敏
func bankCard(cardNo string, cardLength int) string {
	// 如果卡号为空，则直接返回原卡号
	if stringx.IsEmpty(cardNo) {
		return cardNo
	}

	// 清理卡号中的空格
	cleanCardNo := stringx.CleanEmpty(cardNo)
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
	lastFour := cleanCardNo[len(cleanCardNo)-4:]
	newCardNo = append(newCardNo, []rune(lastFour)...)

	// 返回格式化后的卡号
	return string(newCardNo)
}

// ipv4 脱敏
func ipv4(ip string) string {
	// 空判断
	if stringx.IsEmpty(ip) {
		return ip
	}
	newIP := stringx.SubBefore(ip, ".", false)
	return newIP + ".*.*.*"
}

// ipv6 脱敏
func ipv6(ip string) string {
	// 空判断
	if stringx.IsEmpty(ip) {
		return ip
	}
	newIP := stringx.SubBefore(ip, ":", false)
	return newIP + ":*:*:*:*:*:*:*"
}
