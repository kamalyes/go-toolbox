/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 16:05:00
 * @FilePath: \go-toolbox\pkg\metadata\utils.go
 * @Description: 元数据工具函数
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package metadata

import "strings"

// GetTLSVersionString 将 TLS 版本号转换为可读字符串
func GetTLSVersionString(version uint16) string {
	switch version {
	case 0x0304:
		return "TLS 1.3"
	case 0x0303:
		return "TLS 1.2"
	case 0x0302:
		return "TLS 1.1"
	case 0x0301:
		return "TLS 1.0"
	case 0x0300:
		return "SSL 3.0"
	default:
		return ""
	}
}

// ParseAcceptLanguage 解析 Accept-Language 头获取主要语言和地区代码
// 返回: 语言代码(如 "zh"), 地区代码(如 "CN"), 完整标签(如 "zh-CN")
func ParseAcceptLanguage(acceptLang string) (language, region, fullTag string) {
	if acceptLang == "" {
		return "", "", ""
	}

	// Accept-Language 格式: zh-CN,zh;q=0.9,en;q=0.8
	// 提取第一个语言标签（优先级最高）
	endIdx := len(acceptLang)
	for i, c := range acceptLang {
		if c == ',' || c == ';' {
			endIdx = i
			break
		}
	}

	fullTag = strings.TrimSpace(acceptLang[:endIdx])
	if fullTag == "" {
		return "", "", ""
	}

	// 标准化：替换下划线为连字符
	fullTag = strings.ReplaceAll(fullTag, "_", "-")

	// 解析语言和地区代码
	// 格式可能是: "zh-CN", "en-US", "zh", "en"
	parts := strings.Split(fullTag, "-")
	if len(parts) >= 1 {
		language = strings.ToLower(parts[0])
	}
	if len(parts) >= 2 {
		region = strings.ToUpper(parts[1])
		// 重新组合标准化后的 fullTag
		fullTag = language + "-" + region
	} else {
		// 单一语言代码
		fullTag = language
	}

	return language, region, fullTag
}

// NormalizeLanguage 标准化语言代码
// 例如: "zh-cn" -> "zh-CN", "zh_CN" -> "zh-CN", "EN" -> "en"
func NormalizeLanguage(lang string) string {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return ""
	}

	// 替换下划线为连字符
	lang = strings.ReplaceAll(lang, "_", "-")

	// 处理常见的语言代码格式
	parts := strings.Split(lang, "-")
	if len(parts) == 2 {
		// 例如: zh-cn -> zh-CN, EN-us -> en-US
		return strings.ToLower(parts[0]) + "-" + strings.ToUpper(parts[1])
	}

	// 单一语言代码，统一小写
	return strings.ToLower(lang)
}

// GetRemoteIP 从 RemoteAddr 中提取 IP 地址（去除端口）
func GetRemoteIP(remoteAddr string) string {
	if remoteAddr == "" {
		return ""
	}

	// RemoteAddr 格式: IP:Port 或 [IPv6]:Port
	// 使用快速字符串查找，从右向左找第一个冒号
	lastColon := -1
	for i := len(remoteAddr) - 1; i >= 0; i-- {
		if remoteAddr[i] == ':' {
			lastColon = i
			break
		}
		// 如果遇到 ]，说明是 IPv6，停止查找
		if remoteAddr[i] == ']' {
			break
		}
	}

	if lastColon > 0 {
		ip := remoteAddr[:lastColon]
		// 去除 IPv6 的方括号
		if len(ip) > 2 && ip[0] == '[' && ip[len(ip)-1] == ']' {
			return ip[1 : len(ip)-1]
		}
		return ip
	}

	return remoteAddr
}

// GetRemotePort 从 RemoteAddr 中提取端口号
func GetRemotePort(remoteAddr string) string {
	if remoteAddr == "" {
		return ""
	}

	// 从右向左找第一个冒号
	lastColon := -1
	for i := len(remoteAddr) - 1; i >= 0; i-- {
		if remoteAddr[i] == ':' {
			lastColon = i
			break
		}
		if remoteAddr[i] == ']' {
			break
		}
	}

	if lastColon > 0 && lastColon < len(remoteAddr)-1 {
		return remoteAddr[lastColon+1:]
	}

	return ""
}
