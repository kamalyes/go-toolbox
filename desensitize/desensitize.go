/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 17:26:07
 * @FilePath: \go-toolbox\desensitize\desensitize.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package desensitize

import "github.com/kamalyes/go-toolbox/stringx"

// Desensitize 数据脱敏
func Desensitize(str string, DesensitizeType DesensitizeType) string {
	if stringx.IsEmpty(str) {
		return ""
	}
	var newStr string
	switch DesensitizeType {
	case UserId:
		newStr = userId(str)
	case ChineseName:
		newStr = chineseName(str)
	case IDCard:
		newStr = idCard(str, 1, 2)
	case FixedPhone:
		newStr = fixedPhone(str)
	case MobilePhone:
		newStr = mobilePhone(str)
	case Address:
		newStr = address(str, 8)
	case Email:
		newStr = email(str)
	case Password:
		newStr = password(str)
	case CarLicense:
		newStr = carLicense(str)
	case BankCard:
		newStr = bankCard(str)
	case IPV4:
		newStr = ipv4(str)
	case IPV6:
		newStr = ipv6(str)
	default:
		newStr = str
	}

	return newStr
}

// 用户ID脱敏
func userId(id string) string {
	return "0"
}

// 中文名称脱敏
func chineseName(name string) string {
	// 只展示第一个字符
	if stringx.IsEmpty(name) {
		return ""
	}
	return stringx.Hide(name, 1, stringx.Length(name))
}

// 身份证脱敏
func idCard(cardNo string, front int, end int) string {
	if stringx.IsEmpty(cardNo) {
		return ""
	}
	if front+end > stringx.Length(cardNo) {
		return ""
	}
	if front < 0 || end < 0 {
		return ""
	}

	return stringx.Hide(cardNo, front, stringx.Length(cardNo)-end)
}

// 固定电话脱敏
func fixedPhone(phone string) string {
	if stringx.IsEmpty(phone) {
		return ""
	}

	return stringx.Hide(phone, 4, stringx.Length(phone)-2)
}

// 移动电话脱敏
func mobilePhone(phone string) string {
	if stringx.IsEmpty(phone) {
		return ""
	}

	return stringx.Hide(phone, 3, stringx.Length(phone)-4)
}

// 地址脱敏,后 s 位进行脱敏
func address(addr string, sensitiveSize int) string {
	if stringx.IsEmpty(addr) {
		return ""
	}
	length := stringx.Length(addr)
	return stringx.Hide(addr, length-sensitiveSize, length)
}

// 邮箱脱敏
func email(email string) string {
	if stringx.IsEmpty(email) {
		return ""
	}
	index := stringx.IndexOf(email, "@")
	if index <= 1 {
		return email
	}
	return stringx.Hide(email, 1, index)
}

// 密码脱敏
func password(pass string) string {
	if stringx.IsEmpty(pass) {
		return ""
	}
	return stringx.RepeatByLength("*", stringx.Length(pass))
}

// 车牌号脱敏
func carLicense(carNo string) string {
	if stringx.IsEmpty(carNo) {
		return ""
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
func bankCard(cardNo string) string {
	if stringx.IsEmpty(cardNo) {
		return cardNo
	}

	cleanCarNo := stringx.CleanEmpty(cardNo)

	if stringx.Length(cleanCarNo) < 9 {
		return cardNo
	}
	length := stringx.Length(cleanCarNo)
	endLength := length % 4
	if endLength == 0 {
		endLength = 4
	}
	midLength := length - 4 - endLength

	newCardNo := []rune(cleanCarNo[:4])
	for i := 0; i < midLength; i++ {
		if i%4 == 0 {
			newCardNo = append(newCardNo, ' ')
		}
		newCardNo = append(newCardNo, '*')
	}
	newCardNo = append(newCardNo, ' ')
	newCardNo = append(newCardNo, []rune(cleanCarNo[length-endLength:length])...)

	return string(newCardNo)
}

// ipv4 脱敏
func ipv4(ip string) string {
	newIP := stringx.SubBefore(ip, ".", false)
	return newIP + ".*.*.*"
}

// ipv6 脱敏
func ipv6(ip string) string {
	newIP := stringx.SubBefore(ip, ":", false)
	return newIP + ":*:*:*:*:*:*:*"
}
