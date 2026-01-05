/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-19 10:25:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-19 10:25:55
 * @FilePath: \go-toolbox\pkg\validator\ip.go
 * @Description: IP地址验证器
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"fmt"
	"net"
)

// IPBase IP验证器基类
type IPBase struct{}

// ValidateIP 验证IP地址是否有效
func (b *IPBase) ValidateIP(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("无效的IP地址: %s", ip)
	}
	return nil
}

// checkIPInCIDR 检查 IP 是否在 CIDR 列表中
func checkIPInCIDR(ip string, cidrList []string) bool {
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	for _, cidr := range cidrList {
		// 精确匹配
		if ip == cidr {
			return true
		}
		// CIDR 格式匹配
		_, ipNet, err := net.ParseCIDR(cidr)
		if err == nil && ipNet.Contains(clientIP) {
			return true
		}
	}
	return false
}

// IsIPAllowed 检查 IP 是否在允许列表中
// 支持：
// - 空列表：返回 true（允许所有）
// - 通配符 "*"：返回 true（允许所有）
// - 精确匹配：如 "192.168.1.100"
// - CIDR 格式：如 "192.168.1.0/24"
// - IPv6 支持
func IsIPAllowed(ip string, cidrList []string) bool {
	// 空列表，允许所有
	if len(cidrList) == 0 {
		return true
	}

	// 检查通配符
	for _, cidr := range cidrList {
		if cidr == "*" {
			return true
		}
	}

	return checkIPInCIDR(ip, cidrList)
}

// IsIPBlocked 检查 IP 是否在黑名单中
// 支持：
// - 空列表：返回 false（不阻止）
// - 通配符 "*"：返回 true（阻止所有）
// - 精确匹配：如 "192.168.1.100"
// - CIDR 格式：如 "192.168.1.0/24"
// - IPv6 支持
func IsIPBlocked(ip string, blacklist []string) bool {
	// 空列表，不阻止
	if len(blacklist) == 0 {
		return false
	}

	return checkIPInCIDR(ip, blacklist)
}

// MatchIPPattern 匹配 IP 模式（支持通配符 * 和 CIDR）
// 返回是否匹配
func MatchIPPattern(ip, pattern string) bool {
	// 通配符
	if pattern == "*" {
		return true
	}

	// 精确匹配
	if ip == pattern {
		return true
	}

	return checkIPInCIDR(ip, []string{pattern})
}

// IsPrivateIP 检查是否是私有IP地址
func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// IPv4 私有地址段
	privateIPBlocks := []string{
		"10.0.0.0/8",     // Class A
		"172.16.0.0/12",  // Class B
		"192.168.0.0/16", // Class C
		"127.0.0.0/8",    // Loopback
		"169.254.0.0/16", // Link-local
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 unique local
		"fe80::/10",      // IPv6 link-local
	}

	return checkIPInCIDR(ip, privateIPBlocks)
}
