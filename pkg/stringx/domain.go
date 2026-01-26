/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-26 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-26 00:00:00
 * @FilePath: \go-toolbox\pkg\stringx\domain.go
 * @Description: 域名处理工具函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package stringx

import "strings"

// RootDomainPrefix 根域名前缀标识符（用于DNS记录等场景）
const RootDomainPrefix = "@"

// ExtractDomainPrefix 从完整域名中提取前缀
// 示例: ExtractDomainPrefix("www.example.com", "example.com") -> "www"
// 示例: ExtractDomainPrefix("example.com", "example.com") -> "@"
func ExtractDomainPrefix(fullDomain, primaryDomain string) string {
	if fullDomain == "" || fullDomain == primaryDomain {
		return RootDomainPrefix
	}

	if primaryDomain == "" || len(fullDomain) <= len(primaryDomain) {
		return fullDomain
	}

	suffix := "." + primaryDomain
	if !strings.HasSuffix(fullDomain, suffix) {
		return fullDomain
	}

	return strings.TrimSuffix(fullDomain, suffix)
}

// IsSubdomain 判断是否为子域名
func IsSubdomain(subdomain, primaryDomain string) bool {
	return subdomain != primaryDomain && strings.HasSuffix(subdomain, "."+primaryDomain)
}

// SplitDomain 分割域名为前缀和主域名
// 返回: prefix, primaryDomain
// 示例: SplitDomain("www.api.example.com", "example.com") -> "www.api", "example.com"
func SplitDomain(fullDomain, primaryDomain string) (string, string) {
	prefix := ExtractDomainPrefix(fullDomain, primaryDomain)
	if prefix == RootDomainPrefix {
		return "", primaryDomain
	}
	return prefix, primaryDomain
}
