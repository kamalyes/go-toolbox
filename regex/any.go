/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-30 17:26:07
 * @FilePath: \go-toolbox\regex\any.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package regex

import (
	"net"
	"regexp"
	"strconv"
	"strings"
)

type AnyRegs struct{}

func NewAnyRegs() *AnyRegs {
	return &AnyRegs{}
}

// MatchIntOrFloat 整数或者小数
func (g *AnyRegs) MatchIntOrFloat(str string) bool {
	compile := regexp.MustCompile(regIntOrFloat)
	return compile.MatchString(str)
}

// MatchNumber 纯数字
func (g *AnyRegs) MatchNumber(str string) bool {
	compile := regexp.MustCompile(regNumber)
	return compile.MatchString(str)
}

// MatchLenNNumber 长度为n的纯数字
func (g *AnyRegs) MatchLenNNumber(str string, n int) bool {
	nu := strconv.Itoa(n)
	reg := strings.Replace(regLenNNumber, "n", nu, 1)
	compile := regexp.MustCompile(reg)
	return compile.MatchString(str)
}

// MatchGeNNumber 长度不小于n位的纯数字
func (g *AnyRegs) MatchGeNNumber(str string, n int) bool {
	nu := strconv.Itoa(n)
	reg := strings.Replace(regGeNNumber, "n", nu, 1)
	compile := regexp.MustCompile(reg)
	return compile.MatchString(str)
}

// MatchMNIntervalNumber 长度m~n位的纯数字
func (g *AnyRegs) MatchMNIntervalNumber(str string, m, n int) bool {
	mu := strconv.Itoa(m)
	nu := strconv.Itoa(n)
	reg := strings.Replace(regMNIntervalNumber, "m", mu, 1)
	reg = strings.Replace(reg, "n", nu, 1)
	compile := regexp.MustCompile(reg)
	return compile.MatchString(str)
}

// MatchStartingWithNonZero 非零开头的纯数字
func (g *AnyRegs) MatchStartingWithNonZero(str string) bool {
	compile := regexp.MustCompile(regStartingWithNonZero)
	return compile.MatchString(str)
}

// MatchNNovelsOfRealNumber 有n位小数的正实数
func (g *AnyRegs) MatchNNovelsOfRealNumber(str string, n int) bool {
	nu := strconv.Itoa(n)
	reg := strings.Replace(regNNovelsOfRealNumber, "n", nu, 1)
	compile := regexp.MustCompile(reg)
	return compile.MatchString(str)
}

// MatchMNNovelsOfRealNumber m~n位小数的正实数
func (g *AnyRegs) MatchMNNovelsOfRealNumber(str string, m, n int) bool {
	mu := strconv.Itoa(m)
	nu := strconv.Itoa(n)
	reg := strings.Replace(regMNNovelsOfRealNumber, "m", mu, 1)
	reg = strings.Replace(reg, "n", nu, 1)
	compile := regexp.MustCompile(reg)
	return compile.MatchString(str)
}

// MatchNanZeroNumber 非零的正整数
func (g *AnyRegs) MatchNanZeroNumber(str string) bool {
	compile := regexp.MustCompile(regNanZeroNumber)
	return compile.MatchString(str)
}

// MatchNanZeroNegNumber 非零的负整数
func (g *AnyRegs) MatchNanZeroNegNumber(str string) bool {
	compile := regexp.MustCompile(regNanZeroNegNumber)
	return compile.MatchString(str)
}

// MatchNLeCharacter 长度为n的字符，特殊字符除外
func (g *AnyRegs) MatchNLeCharacter(str string, n int) bool {
	nu := strconv.Itoa(n)
	reg := strings.Replace(regNLeCharacter, "n", nu, 1)
	compile := regexp.MustCompile(reg)
	return compile.MatchString(str)
}

// MatchEnCharacter 纯英文字符串,大小写不敏感
func (g *AnyRegs) MatchEnCharacter(str string) bool {
	compile := regexp.MustCompile(regEnCharacter)
	return compile.MatchString(str)
}

// MatchEnCharacterDotUnderLine 检查字符串是否为只包含字母、数字、点号和下划线
func MatchEnCharacterDotUnderLine(str string) bool {
	compile := regexp.MustCompile(regEnCharacterDotUnderLine)
	return compile.MatchString(str)
}

// MatchUpEnCharacter 纯大写英文字符串
func (g *AnyRegs) MatchUpEnCharacter(str string) bool {
	compile := regexp.MustCompile(regUpEnCharacter)
	return compile.MatchString(str)
}

// MatchLowerEnCharacter 纯小写英文字符串
func (g *AnyRegs) MatchLowerEnCharacter(str string) bool {
	compile := regexp.MustCompile(regLowerEnCharacter)
	return compile.MatchString(str)
}

// MatchNumberEnCharacter 数字和26个英文字母组成的字符串,大小写不敏感
func (g *AnyRegs) MatchNumberEnCharacter(str string) bool {
	compile := regexp.MustCompile(regNumberEnCharacter)
	return compile.MatchString(str)
}

// MatchNumberEnUnderscores 数字和26个英文字母组成的字符串,大小写不敏感
func (g *AnyRegs) MatchNumberEnUnderscores(str string) bool {
	compile := regexp.MustCompile(regNumberEnUnderscores)
	return compile.MatchString(str)
}

// MatchPass1 密码1 由数字、26个英文字母或者下划线组成的英文开头的字符串, 长度m~n位
func (g *AnyRegs) MatchPass1(str string, m, n int) bool {
	mu := strconv.Itoa(m)
	nu := strconv.Itoa(n)
	reg := strings.Replace(regPass1, "m", mu, 1)
	reg = strings.Replace(reg, "n", nu, 1)
	compile := regexp.MustCompile(reg)
	return compile.MatchString(str)
}

// MatchPass2 密码2
// 密码长度至少为8个字符。
// 包含至少一个小写字母。
// 包含至少一个大写字母。
// 包含至少一个数字。
// 包含至少一个特殊字符（例如 !@#$%^&*() 等
func (g *AnyRegs) MatchPass2(str string) bool {
	if len(str) < 8 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(str)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(str)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(str)
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(str)

	return hasLower && hasUpper && hasDigit && hasSpecial
}

// MatchIsContainSpecialCharacter 验证是否包含特殊字符串
func (g *AnyRegs) MatchIsContainSpecialCharacter(str string) bool {
	compile := regexp.MustCompile(regIsContainSpecialCharacter)
	return compile.MatchString(str)
}

// MatchChineseCharacter 纯汉字
func (g *AnyRegs) MatchChineseCharacter(str string) bool {
	compile := regexp.MustCompile(regChineseCharacter)
	return compile.MatchString(str)
}

// MatchEmail email
func (g *AnyRegs) MatchEmail(str string) bool {
	compile := regexp.MustCompile(regEmail)
	return compile.MatchString(str)
}

// MatchChinesePhoneNumber 大陆手机号
func (g *AnyRegs) MatchChinesePhoneNumber(str string) bool {
	compile := regexp.MustCompile(regChinesePhoneNumber)
	return compile.MatchString(str)
}

// MatchChineseIDCardNumber 验证大陆身份证号
func (g *AnyRegs) MatchChineseIDCardNumber(id string) bool {
	compile := regexp.MustCompile(regChineseIDCardNumber)
	if !compile.MatchString(id) {
		return false
	}
	switch len(id) {
	case 15:
		id = id[:6] + "19" + id[6:]
		return id == id+calculateChecksum(id)
	case 18:
		// Validate the checksum of 18-digit ID card
		return calculateChecksum(id[:17]) == string(id[17])
	}

	return false
}

// calculateChecksum calculates the checksum for the given 17-digit ID card number.
func calculateChecksum(id string) string {
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkMap := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

	sum := 0
	for i, char := range id {
		num, _ := strconv.Atoi(string(char))
		sum += num * weights[i]
	}

	return string(checkMap[sum%11])
}

// MatchContainChineseCharacter 大陆手机号
func (g *AnyRegs) MatchContainChineseCharacter(str string) bool {
	compile := regexp.MustCompile(regContainChineseCharacter)
	return compile.MatchString(str)
}

// MatchDoubleByte 匹配双字节字符(包括汉字在内)
func (g *AnyRegs) MatchDoubleByte(input string) bool {
	re := regexp.MustCompile(regDoubleByte)
	return re.MatchString(input)
}

// MatchEmptyLine 匹配零个或多个空白字符（包括空格、制表符、换页符等）
func (g *AnyRegs) MatchEmptyLine(input string) bool {
	re := regexp.MustCompile(regEmptyLine)
	return re.MatchString(input)
}

// MatchIPv4 ipv4
func (g *AnyRegs) MatchIPv4(input string) bool {
	re := regexp.MustCompile(regIPv4)
	return re.MatchString(input)
}

// MatchIPv6 ipv6
func (g *AnyRegs) MatchIPv6(input string) bool {
	return net.ParseIP(input) != nil && net.ParseIP(input).To4() == nil
}

// RemoveSymbols 正则匹配去掉所有符号
func RemoveSymbols(s string) string {
	reg := regexp.MustCompile(`[^\w]+`)
	return reg.ReplaceAllString(s, "")
}
