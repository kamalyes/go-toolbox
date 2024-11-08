/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 20:06:12
 * @FilePath: \go-toolbox\pkg\stringx\regexp.go
 * @Description: 字符串正则表达式工具包
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package stringx

import (
	"net"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// AnyRegs 包含各种正则表达式的结构体
type AnyRegs struct {
	regIntOrFloat                *regexp.Regexp // 整数或小数
	regNumber                    *regexp.Regexp // 纯数字
	regLenNNumber                *regexp.Regexp // 长度为n的纯数字
	regGeNNumber                 *regexp.Regexp // 长度不小于n的纯数字
	regMNIntervalNumber          *regexp.Regexp // 长度在m到n之间的纯数字
	regStartingWithNonZero       *regexp.Regexp // 非零开头的纯数字
	regNNovelsOfRealNumber       *regexp.Regexp // n位小数的正实数
	regMNNovelsOfRealNumber      *regexp.Regexp // m到n位小数的正实数
	regNanZeroNumber             *regexp.Regexp // 非零的正整数
	regNanZeroNegNumber          *regexp.Regexp // 非零的负整数
	regNLeCharacter              *regexp.Regexp // 长度为n的字符
	regEnCharacter               *regexp.Regexp // 纯英文字符串
	regUpEnCharacter             *regexp.Regexp // 纯大写英文字符串
	regLowerEnCharacter          *regexp.Regexp // 纯小写英文字符串
	regEnCharacterDotUnderLine   *regexp.Regexp // 英文、数字、点和下划线
	regNumberEnCharacter         *regexp.Regexp // 数字和英文字符组成的字符串
	regNumberEnUnderscores       *regexp.Regexp // 数字、英文字符或下划线组成的字符串
	regPass1                     *regexp.Regexp // 密码1规则
	regIsContainSpecialCharacter *regexp.Regexp // 是否包含特殊字符
	regEmail                     *regexp.Regexp // email
	regChinesePhoneNumber        *regexp.Regexp // 大陆手机号
	regChineseIDCardNumber       *regexp.Regexp // 大陆身份证号
	regContainChineseCharacter   *regexp.Regexp // 包含中文字符
	regDoubleByte                *regexp.Regexp // 双字节字符
	regEmptyLine                 *regexp.Regexp // 空行
	regIPv4                      *regexp.Regexp // IPv4
	regTime                      *regexp.Regexp // 时间格式
	regHex                       *regexp.Regexp // 十六进制
}

// 正则表达式常量
const (
	regIntOrFloat                = `^[0-9]+\.{0,1}[0-9]{0,2}$`                    // 整数或小数
	regNumber                    = `^[0-9]*$`                                     // 纯数字
	regLenNNumber                = `^\d{n}$`                                      // 长度为n的纯数字
	regGeNNumber                 = `^\d{n,}$`                                     // 长度不小于n的纯数字
	regMNIntervalNumber          = `^\d{m,n}$`                                    // 长度在m到n之间的纯数字
	regStartingWithNonZero       = `^([1-9][0-9]*)`                               // 非零开头的纯数字
	regNNovelsOfRealNumber       = `^[0-9]+(.[0-9]{n})?$`                         // n位小数的正实数
	regMNNovelsOfRealNumber      = `^[0-9]+(.[0-9]{m,n})?$`                       // m到n位小数的正实数
	regNanZeroNumber             = `^\+?[1-9][0-9]*$`                             // 非零的正整数
	regNanZeroNegNumber          = `^\-[1-9][0-9]*$`                              // 非零的负整数
	regNLeCharacter              = `^.{n}$`                                       // 长度为n的字符
	regEnCharacter               = `^[A-Za-z]+$`                                  // 纯英文字符串
	regUpEnCharacter             = `^[A-Z]+$`                                     // 纯大写英文字符串
	regLowerEnCharacter          = `^[a-z]+$`                                     // 纯小写英文字符串
	regEnCharacterDotUnderLine   = `^[a-zA-Z0-9._]+$`                             // 英文、数字、点和下划线
	regNumberEnCharacter         = `^[A-Za-z0-9]+$`                               // 数字和英文字符组成的字符串
	regNumberEnUnderscores       = `^\w+$`                                        // 数字、英文字符或下划线组成的字符串
	regPass1                     = `^[a-zA-Z]\w{m,n}$`                            // 密码1规则
	regIsContainSpecialCharacter = `[!@#\$%\^&\*\(\)_\+\[\]{}|;':",./<>?]`        // 是否包含特殊字符
	regEmail                     = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` // email
	regChinesePhoneNumber        = `^1[3-9]\d{9}$`                                // 大陆手机号
	regChineseIDCardNumber       = `^\d{15}$|^\d{17}(\d|X|x)$`                    // 大陆身份证号
	regContainChineseCharacter   = `[\p{Han}]`                                    // 包含中文字符
	regDoubleByte                = `[^\x00-\xff]`                                 // 双字节字符
	regEmptyLine                 = `^\s*`                                         // 空行
	regIPv4                      = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	regTime                      = `(\d{4}[-/\.]\d{1,2}[-/\.]\d{1,2})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)` // 时间格式
	hexRegex                     = `^[0-9a-fA-F]+`                                                                                                                     // 十六进制
)

// NewAnyRegs 创建一个新的 AnyRegs 实例，并编译所有正则表达式
func NewAnyRegs() *AnyRegs {
	return &AnyRegs{
		regIntOrFloat:                regexp.MustCompile(regIntOrFloat),
		regNumber:                    regexp.MustCompile(regNumber),
		regLenNNumber:                regexp.MustCompile(regLenNNumber),
		regGeNNumber:                 regexp.MustCompile(regGeNNumber),
		regMNIntervalNumber:          regexp.MustCompile(regMNIntervalNumber),
		regStartingWithNonZero:       regexp.MustCompile(regStartingWithNonZero),
		regNNovelsOfRealNumber:       regexp.MustCompile(regNNovelsOfRealNumber),
		regMNNovelsOfRealNumber:      regexp.MustCompile(regMNNovelsOfRealNumber),
		regNanZeroNumber:             regexp.MustCompile(regNanZeroNumber),
		regNanZeroNegNumber:          regexp.MustCompile(regNanZeroNegNumber),
		regNLeCharacter:              regexp.MustCompile(regNLeCharacter),
		regEnCharacter:               regexp.MustCompile(regEnCharacter),
		regUpEnCharacter:             regexp.MustCompile(regUpEnCharacter),
		regLowerEnCharacter:          regexp.MustCompile(regLowerEnCharacter),
		regEnCharacterDotUnderLine:   regexp.MustCompile(regEnCharacterDotUnderLine),
		regNumberEnCharacter:         regexp.MustCompile(regNumberEnCharacter),
		regNumberEnUnderscores:       regexp.MustCompile(regNumberEnUnderscores),
		regPass1:                     regexp.MustCompile(regPass1),
		regIsContainSpecialCharacter: regexp.MustCompile(regIsContainSpecialCharacter),
		regEmail:                     regexp.MustCompile(regEmail),
		regChinesePhoneNumber:        regexp.MustCompile(regChinesePhoneNumber),
		regChineseIDCardNumber:       regexp.MustCompile(regChineseIDCardNumber),
		regContainChineseCharacter:   regexp.MustCompile(regContainChineseCharacter),
		regDoubleByte:                regexp.MustCompile(regDoubleByte),
		regEmptyLine:                 regexp.MustCompile(regEmptyLine),
		regIPv4:                      regexp.MustCompile(regIPv4),
		regTime:                      regexp.MustCompile(regTime),
		regHex:                       regexp.MustCompile(hexRegex),
	}
}

// match 方法封装，减少重复代码
func (g *AnyRegs) match(compiled *regexp.Regexp, str string) bool {
	result := compiled.MatchString(str)
	// log.Printf("AnyRegs Matching string: '%s' against regex: '%s', result: %v", str, compiled.String(), result) // 添加日志打印
	return result
}

// MatchIntOrFloat 检查字符串是否为整数或小数
func (g *AnyRegs) MatchIntOrFloat(str string) bool {
	return g.match(g.regIntOrFloat, str)
}

// MatchNumber 检查字符串是否为纯数字
func (g *AnyRegs) MatchNumber(str string) bool {
	return g.match(g.regNumber, str)
}

// MatchLenNNumber 检查字符串是否为长度为n的纯数字
func (g *AnyRegs) MatchLenNNumber(str string, n int) bool {
	reg := strings.Replace(g.regLenNNumber.String(), "n", strconv.Itoa(n), 1)
	return g.match(regexp.MustCompile(reg), str)
}

// MatchGeNNumber 检查字符串是否为长度不小于n的纯数字
func (g *AnyRegs) MatchGeNNumber(str string, n int) bool {
	reg := strings.Replace(g.regGeNNumber.String(), "n", strconv.Itoa(n), 1)
	return g.match(regexp.MustCompile(reg), str)
}

// MatchMNIntervalNumber 检查字符串是否为长度在m到n之间的纯数字
func (g *AnyRegs) MatchMNIntervalNumber(str string, m, n int) bool {
	mu := strconv.Itoa(m)
	nu := strconv.Itoa(n)
	reg := strings.Replace(g.regMNIntervalNumber.String(), "m", mu, 1)
	reg = strings.Replace(reg, "n", nu, 1)
	return g.match(regexp.MustCompile(reg), str)
}

// MatchStartingWithNonZero 检查字符串是否为非零开头的纯数字
func (g *AnyRegs) MatchStartingWithNonZero(str string) bool {
	return g.match(g.regStartingWithNonZero, str)
}

// MatchNNovelsOfRealNumber 检查字符串是否为有n位小数的正实数
func (g *AnyRegs) MatchNNovelsOfRealNumber(str string, n int) bool {
	reg := strings.Replace(g.regNNovelsOfRealNumber.String(), "n", strconv.Itoa(n), 1)
	return g.match(regexp.MustCompile(reg), str)
}

// MatchMNNovelsOfRealNumber 检查字符串是否为m到n位小数的正实数
func (g *AnyRegs) MatchMNNovelsOfRealNumber(str string, m, n int) bool {
	mu := strconv.Itoa(m)
	nu := strconv.Itoa(n)
	reg := strings.Replace(g.regMNNovelsOfRealNumber.String(), "m", mu, 1)
	reg = strings.Replace(reg, "n", nu, 1)
	return g.match(regexp.MustCompile(reg), str)
}

// MatchNanZeroNumber 检查字符串是否为非零的正整数
func (g *AnyRegs) MatchNanZeroNumber(str string) bool {
	return g.match(g.regNanZeroNumber, str)
}

// MatchNanZeroNegNumber 检查字符串是否为非零的负整数
func (g *AnyRegs) MatchNanZeroNegNumber(str string) bool {
	return g.match(g.regNanZeroNegNumber, str)
}

// MatchNLeCharacter 检查字符串是否为长度为n的字符
func (g *AnyRegs) MatchNLeCharacter(str string, n int) bool {
	reg := strings.Replace(g.regNLeCharacter.String(), "n", strconv.Itoa(n), 1)
	return g.match(regexp.MustCompile(reg), str)
}

// MatchEnCharacter 检查字符串是否为纯英文字符串（大小写不敏感）
func (g *AnyRegs) MatchEnCharacter(str string) bool {
	return g.match(g.regEnCharacter, str)
}

// MatchEnCharacterDotUnderLine 检查字符串是否只包含字母、数字、点号和下划线
func MatchEnCharacterDotUnderLine(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._]+`).MatchString(str)
}

// MatchUpEnCharacter 检查字符串是否为纯大写英文字符串
func (g *AnyRegs) MatchUpEnCharacter(str string) bool {
	return g.match(g.regUpEnCharacter, str)
}

// MatchLowerEnCharacter 检查字符串是否为纯小写英文字符串
func (g *AnyRegs) MatchLowerEnCharacter(str string) bool {
	return g.match(g.regLowerEnCharacter, str)
}

// MatchNumberEnCharacter 检查字符串是否由数字和26个英文字母组成（大小写不敏感）
func (g *AnyRegs) MatchNumberEnCharacter(str string) bool {
	return g.match(g.regNumberEnCharacter, str)
}

// MatchNumberEnUnderscores 检查字符串是否由数字和26个英文字母或下划线组成（大小写不敏感）
func (g *AnyRegs) MatchNumberEnUnderscores(str string) bool {
	return g.match(g.regNumberEnUnderscores, str)
}

// MatchPass1 检查密码1是否符合规则：由数字、26个英文字母或下划线组成的英文开头的字符串，长度在m到n位之间
func (g *AnyRegs) MatchPass1(str string, m, n int) bool {
	mu := strconv.Itoa(m)
	nu := strconv.Itoa(n)
	reg := strings.Replace(g.regPass1.String(), "m", mu, 1)
	reg = strings.Replace(reg, "n", nu, 1)
	return g.match(regexp.MustCompile(reg), str)
}

// MatchPass2 检查密码2是否符合规则：
// 密码长度至少为8个字符。
// 包含至少一个小写字母。
// 包含至少一个大写字母。
// 包含至少一个数字。
// 包含至少一个特殊字符（例如 !@#$%^&*() 等）
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

// MatchIsContainSpecialCharacter 检查字符串是否包含特殊字符
func (g *AnyRegs) MatchIsContainSpecialCharacter(str string) bool {
	return g.match(g.regIsContainSpecialCharacter, str)
}

// IsChineseCharacter 检查字符串是否为纯汉字
func IsChineseCharacter(str string) (isContains bool, count int) {
	for _, v := range str {
		if !unicode.Is(unicode.Han, v) {
			count++
		}
	}

	if count == 0 {
		isContains = true
	}

	return
}

// MatchEmail 检查字符串是否为有效的email
func (g *AnyRegs) MatchEmail(str string) bool {
	return g.match(g.regEmail, str)
}

// MatchChinesePhoneNumber 检查字符串是否为有效的大陆手机号
func (g *AnyRegs) MatchChinesePhoneNumber(str string) bool {
	return g.match(g.regChinesePhoneNumber, str)
}

// MatchChineseIDCardNumber 检查字符串是否为有效的大陆身份证号
func (g *AnyRegs) MatchChineseIDCardNumber(id string) bool {
	if !g.match(g.regChineseIDCardNumber, id) {
		return false
	}
	switch len(id) {
	case 15:
		id = id[:6] + "19" + id[6:] // 将15位身份证号转换为18位
		return id == id+calculateChecksum(id)
	case 18:
		// 验证18位身份证号的校验和
		return calculateChecksum(id[:17]) == string(id[17])
	}

	return false
}

// calculateChecksum 计算给定17位身份证号的校验和
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

// MatchContainChineseCharacter 检查字符串是否包含中文字符
func (g *AnyRegs) MatchContainChineseCharacter(str string) bool {
	return g.match(g.regContainChineseCharacter, str)
}

// MatchDoubleByte 检查字符串是否包含双字节字符（包括汉字）
func (g *AnyRegs) MatchDoubleByte(input string) bool {
	return g.match(g.regDoubleByte, input)
}

// MatchEmptyLine 检查字符串是否为空行
func (g *AnyRegs) MatchEmptyLine(input string) bool {
	return g.match(g.regEmptyLine, input)
}

// MatchIPv4 检查字符串是否为有效的IPv4地址
func (g *AnyRegs) MatchIPv4(input string) bool {
	return g.match(g.regIPv4, input)
}

// MatchIPv6 检查字符串是否为有效的IPv6地址
func (g *AnyRegs) MatchIPv6(input string) bool {
	return net.ParseIP(input) != nil && net.ParseIP(input).To4() == nil
}

// RemoveSymbols 使用正则表达式去掉字符串中的所有符号
func RemoveSymbols(s string) string {
	reg := regexp.MustCompile(`[^\w]+`)
	return reg.ReplaceAllString(s, "")
}

// MatchTime 检查字符串是否符合时间格式
func (g *AnyRegs) MatchTime(input string) bool {
	return g.match(g.regTime, input)
}

// IsHex 检查字符串是否为有效的十六进制数
func (g *AnyRegs) IsHex(input string) bool {
	return g.match(g.regHex, input)
}

// IsTrueString 检查字符串是否表示为 true
func IsTrueString(s string) bool {
	return strings.EqualFold(s, "true") || strings.EqualFold(s, "1") || strings.EqualFold(s, "yes")
}

// HasLocalIP 方法检查给定的IP地址是否是本地地址
func HasLocalIP(ip string) bool {
	// 检查硬编码的本地地址和主机名
	localIPsAndHostnames := []string{"localhost", "127.0.0.1", "::1"}
	for _, local := range localIPsAndHostnames {
		if ip == local {
			return true
		}
	}

	// 解析IP地址，以便我们可以检查IPv4地址的范围
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		// 如果无法解析IP地址，则不是本地地址
		return false
	}

	// 检查IPv4的私有地址范围
	// 注意：这里没有包含192.168.0.1的特定检查，因为它已经在上面的字符串列表中。
	// 如果你需要更广泛的私有地址检查，可以使用以下范围：
	// 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
	privateIPv4Blocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range privateIPv4Blocks {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			// 这通常不会发生，除非CIDR字符串格式错误
			continue
		}
		if ipnet.Contains(parsedIP) {
			return true
		}
	}

	// 如果都不是，则返回false
	return false
}

// IsLinkLocal 检查给定的 IP 是否为链路本地地址
func IsLinkLocal(ip net.IP) bool {
	// 检查 IPv4 链路本地地址
	if ipv4 := ip.To4(); ipv4 != nil {
		return ipv4[0] == 169 && ipv4[1] == 254
	}

	// 检查 IPv6 链路本地地址
	if ipv6 := ip.To16(); ipv6 != nil {
		return ipv6[0] == 0xfe && (ipv6[1]&0xc0) == 0x80 // fe80::/10
	}

	return false
}

// IsUniqueLocalAddress 检查给定的 IP 是否为唯一本地地址（ULA）
func IsUniqueLocalAddress(ip net.IP) bool {
	if ipv6 := ip.To16(); ipv6 != nil {
		return ipv6[0] == 0xFC || ipv6[0] == 0xFD // fc00::/7
	}
	return false
}

// IsGlobalUnicast 检查给定的 IP 地址是否是全球单播地址
func IsGlobalUnicast(ip string) bool {
	// 解析IP地址，以便我们可以检查IPv4地址的范围
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		// 如果无法解析IP地址，则不是本地地址
		return false
	}

	if ip4 := parsedIP.To4(); ip4 != nil {
		// 检查是否是私有地址
		if HasLocalIP(ip) || ip4[0] == 127 || (ip4[0] == 169 && ip4[1] == 254) {
			return false
		}
		return true // 其他情况为全球单播地址
	}

	if ip6 := parsedIP.To16(); ip6 != nil {
		return !parsedIP.IsUnspecified() && !parsedIP.IsLoopback() && !IsLinkLocal(parsedIP) && !parsedIP.IsPrivate() && !isDocumentationAddress(ip6)
	}

	return false
}

// isDocumentationAddress 检查是否是文档专用地址
func isDocumentationAddress(ip net.IP) bool {
	return ip[0] == 0x20 && ip[1] == 0x01 && ip[2] == 0x0d && ip[3] == 0xb8
}
