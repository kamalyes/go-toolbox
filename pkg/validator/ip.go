/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-02-19 10:25:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-02-19 10:25:55
 * @FilePath: \go-toolbox\pkg\validator\ip.go
 * @Description:
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

	// 解析客户端 IP
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	// 检查每个规则
	for _, cidr := range cidrList {
		// 精确匹配
		if ip == cidr {
			return true
		}
		// CIDR 格式匹配
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(clientIP) {
			return true
		}
	}
	return false
}
