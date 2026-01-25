/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-25 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-25 00:00:00
 * @FilePath: \go-toolbox\pkg\validator\format.go
 * @Description: 格式验证函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

// UUID 正则表达式
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ValidateString 验证字符串（支持各种操作符）
func ValidateString(actual, expect string, op CompareOperator) CompareResult {
	result := CompareResult{
		Actual: actual,
		Expect: expect,
	}

	switch op {
	case OpEqual, OpSymbolEqual:
		result.Success = actual == expect
		if !result.Success {
			result.Message = fmt.Sprintf("字符串不相等: 期望 '%s', 实际 '%s'", expect, actual)
		}

	case OpNotEqual, OpSymbolNotEqual:
		result.Success = actual != expect
		if !result.Success {
			result.Message = fmt.Sprintf("字符串应该不相等: 都是 '%s'", actual)
		}

	case OpContains:
		result.Success = strings.Contains(actual, expect)
		if !result.Success {
			result.Message = fmt.Sprintf("字符串不包含: '%s' 中未找到 '%s'", actual, expect)
		}

	case OpNotContains:
		result.Success = !strings.Contains(actual, expect)
		if !result.Success {
			result.Message = fmt.Sprintf("字符串不应包含: '%s' 中找到了 '%s'", actual, expect)
		}

	case OpHasPrefix:
		result.Success = strings.HasPrefix(actual, expect)
		if !result.Success {
			result.Message = fmt.Sprintf("字符串前缀不匹配: '%s' 不以 '%s' 开头", actual, expect)
		}

	case OpHasSuffix:
		result.Success = strings.HasSuffix(actual, expect)
		if !result.Success {
			result.Message = fmt.Sprintf("字符串后缀不匹配: '%s' 不以 '%s' 结尾", actual, expect)
		}

	case OpEmpty:
		result.Success = len(strings.TrimSpace(actual)) == 0
		result.Expect = "empty string"
		if !result.Success {
			result.Message = fmt.Sprintf("字符串不为空: '%s'", actual)
		}

	case OpNotEmpty:
		result.Success = len(strings.TrimSpace(actual)) > 0
		result.Expect = "non-empty string"
		if !result.Success {
			result.Message = "字符串为空"
		}

	case OpRegex:
		re, err := regexp.Compile(expect)
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("正则表达式编译失败: %s", err.Error())
		} else {
			result.Success = re.MatchString(actual)
			if !result.Success {
				result.Message = fmt.Sprintf("正则表达式不匹配: '%s' 不匹配模式 '%s'", actual, expect)
			}
		}

	default:
		result.Success = false
		result.Message = fmt.Sprintf("不支持的字符串操作符: %s", op)
	}

	if result.Success && result.Message == "" {
		result.Message = "字符串验证通过"
	}

	return result
}

// ValidateEmail 验证 Email 格式
func ValidateEmail(email string) CompareResult {
	result := CompareResult{
		Actual: email,
		Expect: "valid email format",
	}

	email = strings.TrimSpace(email)
	if email == "" {
		result.Success = false
		result.Message = "Email 地址为空"
		return result
	}

	// 使用标准库的 mail.ParseAddress 验证
	addr, err := mail.ParseAddress(email)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("无效的 Email 格式: %s", err.Error())
		return result
	}

	// 验证域名部分
	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		result.Success = false
		result.Message = "Email 格式错误: 缺少 @ 符号"
		return result
	}

	domain := parts[1]
	if domain == "" || !strings.Contains(domain, ".") {
		result.Success = false
		result.Message = "Email 格式错误: 无效的域名"
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("Email 格式验证通过: %s", email)
	return result
}

// ValidateIP 验证 IP 地址（支持 IPv4 和 IPv6）
func ValidateIP(ipStr string) CompareResult {
	result := CompareResult{
		Actual: ipStr,
		Expect: "valid IP address (IPv4 or IPv6)",
	}

	ipStr = strings.TrimSpace(ipStr)
	if ipStr == "" {
		result.Success = false
		result.Message = "IP 地址为空"
		return result
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		result.Success = false
		result.Message = fmt.Sprintf("无效的 IP 地址: %s", ipStr)
		return result
	}

	// 判断 IP 类型
	ipType := "IPv4"
	if ip.To4() == nil {
		ipType = "IPv6"
	}

	result.Success = true
	result.Message = fmt.Sprintf("%s 地址验证通过: %s", ipType, ipStr)
	return result
}

// ValidateProtocol 验证 URL 协议（支持多种协议：http, https, ws, wss, ftp, ftps 等）
// allowedProtocols: 允许的协议列表，为空则允许所有常见协议
func ValidateProtocol(urlStr string, allowedProtocols ...string) CompareResult {
	result := CompareResult{
		Actual: urlStr,
	}

	// 设置默认允许的协议
	if len(allowedProtocols) == 0 {
		allowedProtocols = []string{"http", "https", "ws", "wss", "ftp", "ftps"}
	}
	result.Expect = fmt.Sprintf("valid URL with protocol: %v", allowedProtocols)

	urlStr = strings.TrimSpace(urlStr)
	if urlStr == "" {
		result.Success = false
		result.Message = "URL 为空"
		return result
	}

	// 解析 URL
	u, err := url.Parse(urlStr)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("无效的 URL: %s", err.Error())
		return result
	}

	// 检查协议是否为空
	if u.Scheme == "" {
		result.Success = false
		result.Message = "URL 缺少协议"
		return result
	}

	// 检查协议是否在允许列表中
	schemeValid := false
	for _, allowed := range allowedProtocols {
		if strings.EqualFold(u.Scheme, allowed) {
			schemeValid = true
			break
		}
	}

	if !schemeValid {
		result.Success = false
		result.Message = fmt.Sprintf("不支持的协议: %s (允许: %v)", u.Scheme, allowedProtocols)
		return result
	}

	// 检查主机名（某些协议可能不需要）
	if u.Host == "" && (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "ws" || u.Scheme == "wss" || u.Scheme == "ftp" || u.Scheme == "ftps") {
		result.Success = false
		result.Message = "URL 缺少主机名"
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("URL 协议验证通过: %s (%s)", u.Scheme, urlStr)
	return result
}

// ValidateHTTP 验证 HTTP/HTTPS URL（便捷函数）
func ValidateHTTP(urlStr string) CompareResult {
	return ValidateProtocol(urlStr, "http", "https")
}

// ValidateWebSocket 验证 WebSocket URL（便捷函数）
func ValidateWebSocket(urlStr string) CompareResult {
	return ValidateProtocol(urlStr, "ws", "wss")
}

// ValidateUUID 验证 UUID 格式（支持 UUID v1-v5）
func ValidateUUID(uuidStr string) CompareResult {
	result := CompareResult{
		Actual: uuidStr,
		Expect: "valid UUID format",
	}

	uuidStr = strings.TrimSpace(uuidStr)
	if uuidStr == "" {
		result.Success = false
		result.Message = "UUID 为空"
		return result
	}

	// 使用正则表达式验证 UUID 格式
	if !uuidRegex.MatchString(uuidStr) {
		result.Success = false
		result.Message = fmt.Sprintf("无效的 UUID 格式: %s (应为 xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)", uuidStr)
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("UUID 格式验证通过: %s", uuidStr)
	return result
}

// ValidateBase64 验证 Base64 编码
func ValidateBase64(str string) CompareResult {
	result := CompareResult{
		Actual: str,
		Expect: "valid Base64 encoding",
	}

	str = strings.TrimSpace(str)
	if str == "" {
		result.Success = false
		result.Message = "Base64 字符串为空"
		return result
	}

	// 尝试解码 Base64
	_, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		// 尝试 URL 安全的 Base64
		_, err = base64.URLEncoding.DecodeString(str)
		if err != nil {
			// 尝试无填充的 Base64
			_, err = base64.RawStdEncoding.DecodeString(str)
			if err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("无效的 Base64 编码: %s", err.Error())
				return result
			}
		}
	}

	result.Success = true
	result.Message = "Base64 编码验证通过"
	return result
}
