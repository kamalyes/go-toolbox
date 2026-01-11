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
	"bytes"
	"fmt"
	"net"
	"strings"
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

// MatchIPInList 检查 IP 是否匹配列表中的任一规则
// 支持多种格式:
// 1. 分隔符列表: 支持分号、逗号、换行符、空格分隔
//   - "192.168.0.10;192.168.0.1"
//   - "192.168.0.1,192.168.0.56"
//   - "192.168.0.15\n192.168.0.1"
//   - "192.168.0.15 192.168.0.17"
//
// 2. 精确匹配: "192.168.1.100"
// 3. CIDR: "192.168.1.0/24"
// 4. IP范围: "172.16.0.1-172.16.0.10"
// 5. 通配符: "192.168.2.*" 或 "192.168.*.*"
func MatchIPInList(ip string, ipList []string) bool {
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	for _, pattern := range ipList {
		// 1. 分隔符列表匹配 (支持分号、逗号、换行符、空格、Tab)
		// 注意：如果包含 '-' 可能是 IP 范围，需要特殊处理
		hasSeparator := strings.ContainsAny(pattern, ";,\n\t") ||
			(strings.Contains(pattern, " ") && !strings.Contains(pattern, "-"))

		if hasSeparator {
			// 统一处理多种分隔符
			pattern = strings.ReplaceAll(pattern, "\r\n", "\n")
			pattern = strings.ReplaceAll(pattern, "\t", "\n")
			pattern = strings.ReplaceAll(pattern, ";", "\n")
			pattern = strings.ReplaceAll(pattern, ",", "\n")

			// 按换行符分割
			lines := strings.Split(pattern, "\n")
			var subPatterns []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				// 如果行内还有空格，继续按空格分割
				if strings.Contains(line, " ") {
					parts := strings.Fields(line)
					subPatterns = append(subPatterns, parts...)
				} else {
					subPatterns = append(subPatterns, line)
				}
			}

			// 递归检查每个子模式
			for _, subPattern := range subPatterns {
				if subPattern == "" {
					continue
				}
				if MatchIPInList(ip, []string{subPattern}) {
					return true
				}
			}
			continue
		}
		// 2. 精确匹配
		if ip == pattern {
			return true
		}

		// 3. CIDR 格式匹配
		_, ipNet, err := net.ParseCIDR(pattern)
		if err == nil && ipNet.Contains(clientIP) {
			return true
		}

		// 4. IP范围匹配 (例如: "172.16.0.1-172.16.0.10")
		if strings.Contains(pattern, "-") {
			parts := strings.Split(pattern, "-")
			if len(parts) == 2 {
				startIP := net.ParseIP(strings.TrimSpace(parts[0]))
				endIP := net.ParseIP(strings.TrimSpace(parts[1]))
				if startIP != nil && endIP != nil {
					if IsIPInRange(clientIP, startIP, endIP) {
						return true
					}
				}
			}
		}

		// 5. 通配符匹配 (例如: "192.168.2.*" 或 "192.168.*.*")
		if strings.Contains(pattern, "*") {
			if MatchIPWithWildcard(ip, pattern) {
				return true
			}
		}
	}
	return false
}

// IsIPInRange 检查 IP 是否在指定范围内
func IsIPInRange(ip, start, end net.IP) bool {
	// 将IPv4地址转换为4字节表示
	ipv4 := ip.To4()
	startv4 := start.To4()
	endv4 := end.To4()

	// 如果任何一个不是IPv4,则不匹配
	if ipv4 == nil || startv4 == nil || endv4 == nil {
		return false
	}

	// 比较字节数组
	return bytes.Compare(ipv4, startv4) >= 0 && bytes.Compare(ipv4, endv4) <= 0
}

// MatchIPWithWildcard 使用通配符匹配 IP
func MatchIPWithWildcard(ip, pattern string) bool {
	// 分割IP和模式为段
	ipParts := strings.Split(ip, ".")
	patternParts := strings.Split(pattern, ".")

	// 长度必须相同
	if len(ipParts) != len(patternParts) {
		return false
	}

	// 逐段比较
	for i := 0; i < len(ipParts); i++ {
		if patternParts[i] != "*" && ipParts[i] != patternParts[i] {
			return false
		}
	}

	return true
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

	return MatchIPInList(ip, cidrList)
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

	return MatchIPInList(ip, blacklist)
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

	return MatchIPInList(ip, []string{pattern})
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

	return MatchIPInList(ip, privateIPBlocks)
}
