/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-16 11:55:28
 * @FilePath: \go-toolbox\pkg\validator\ip_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package validator

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试常量 - IP 地址
const (
	testIPv4_1         = "192.168.1.100"
	testIPv4_2         = "192.168.1.50"
	testIPv4_3         = "192.168.1.1"
	testIPv4_4         = "10.0.0.1"
	testIPv4_5         = "10.0.0.5"
	testIPv4_6         = "172.16.0.1"
	testIPv4_Localhost = "127.0.0.1"
	testIPv4_Zero      = "0.0.0.0"
	testIPv6_1         = "2001:db8::1"
	testIPv6_2         = "2001:db8::abcd"
)

// 测试常量 - CIDR 范围
const (
	testCIDR_192      = "192.168.1.0/24"
	testCIDR_10       = "10.0.0.0/8"
	testCIDR_172      = "172.16.0.0/12"
	testCIDR_192Block = "192.168.0.0/16"
	testCIDR_127      = "127.0.0.0/8"
	testCIDR_IPv6     = "2001:db8::/64"
)

// 测试常量 - 通配符
const (
	testWildcard = "*"
)

// 公共测试用例数据
var (
	// 异常IP地址测试用例
	invalidIPCases = []struct {
		name string
		ip   string
	}{
		{"Empty IP", ""},
		{"Invalid IP format", "invalid-ip"},
		{"Invalid IP with numbers", "256.256.256.256"},
		{"Incomplete IP", "192.168.1"},
		{"IP with port", "192.168.1.1:8080"},
	}

	// IPv6测试用例
	ipv6Cases = []struct {
		name    string
		ip      string
		pattern string
		want    bool
	}{
		{"IPv6 exact match", testIPv6_1, testIPv6_1, true},
		{"IPv6 exact no match", testIPv6_1, "2001:db8::2", false},
		{"IPv6 CIDR match", testIPv6_2, testCIDR_IPv6, true},
		{"IPv6 CIDR no match", "2001:db9::1", testCIDR_IPv6, false},
	}
)

func TestIsIPAllowed(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		cidrList []string
		want     bool
	}{
		// ========== 通配符测试 ==========
		{"Wildcard allows any IP", testIPv4_1, []string{testWildcard}, true},
		{"Wildcard in list", testIPv4_4, []string{testCIDR_192, testWildcard, testCIDR_172}, true},
		{"Wildcard with specific IPs", "8.8.8.8", []string{testWildcard}, true},
		{"Multiple wildcards", "1.2.3.4", []string{testWildcard, testWildcard}, true},

		// ========== 空列表测试 ==========
		{"Empty list allows all", testIPv4_1, []string{}, true},
		{"Nil list allows all", testIPv4_4, nil, true},
		{"Empty list with any IP", testIPv4_Zero, []string{}, true},
		{"Empty list with IPv6", testIPv6_1, []string{}, true},

		// ========== 精确匹配测试 ==========
		{"Exact IP match", testIPv4_1, []string{testIPv4_1}, true},
		{"Exact IP no match", testIPv4_1, []string{"192.168.1.101"}, false},
		{"Multiple exact IPs one match", testIPv4_5, []string{testIPv4_3, testIPv4_5, testIPv4_6}, true},
		{"Multiple exact IPs none match", "10.0.0.6", []string{testIPv4_3, testIPv4_5, testIPv4_6}, false},

		// ========== CIDR 格式测试 ==========
		{"IP in CIDR", testIPv4_2, []string{testCIDR_192}, true},
		{"IP not in CIDR", "192.168.2.1", []string{testCIDR_192}, false},
		{"Multiple CIDRs one match", testIPv4_5, []string{testCIDR_192, testCIDR_10}, true},
		{"Multiple CIDRs none match", testIPv4_6, []string{testCIDR_192, testCIDR_10}, false},
		{"CIDR /32 exact match", testIPv4_1, []string{testIPv4_1 + "/32"}, true},
		{"CIDR /32 no match", "192.168.1.101", []string{testIPv4_1 + "/32"}, false},
		{"Large CIDR /8", "10.255.255.255", []string{testCIDR_10}, true},
		{"CIDR boundary test lower", "192.168.1.0", []string{testCIDR_192}, true},
		{"CIDR boundary test upper", "192.168.1.255", []string{testCIDR_192}, true},
		{"Outside CIDR boundary", "192.168.0.255", []string{testCIDR_192}, false},

		// ========== IPv6 测试 ==========
		{"IPv6 exact match", testIPv6_1, []string{testIPv6_1}, true},
		{"IPv6 exact no match", testIPv6_1, []string{"2001:db8::2"}, false},
		{"IPv6 CIDR match", testIPv6_2, []string{testCIDR_IPv6}, true},
		{"IPv6 CIDR no match", "2001:db9::1", []string{testCIDR_IPv6}, false},
		{"IPv6 loopback", "::1", []string{"::1"}, true},
		{"IPv6 any", "::", []string{"::"}, true},
		{"IPv6 with wildcard", testIPv6_1, []string{testWildcard}, true},

		// ========== 混合测试 ==========
		{"Mixed IPv4 and IPv6", testIPv4_3, []string{testCIDR_IPv6, testCIDR_192}, true},
		{"Mixed exact and CIDR", testIPv4_1, []string{testIPv4_4, testCIDR_192, testIPv4_6}, true},
		{"Wildcard with other rules", "1.2.3.4", []string{testCIDR_192, testWildcard}, true},

		// ========== 异常和边界情况 ==========
		{"Malformed CIDR", testIPv4_3, []string{"192.168.1.0/99"}, false},
		{"CIDR list contains empty string", testIPv4_2, []string{""}, false},
		{"CIDR list contains invalid CIDR", testIPv4_2, []string{"invalid-cidr"}, false},
		{"IP equals CIDR but CIDR invalid", testIPv4_2, []string{testIPv4_2 + "/33"}, false},

		// ========== 特殊 IP 地址 ==========
		{"Localhost IPv4", testIPv4_Localhost, []string{testIPv4_Localhost}, true},
		{"Localhost in CIDR", testIPv4_Localhost, []string{testCIDR_127}, true},
		{"Broadcast IP", "255.255.255.255", []string{"255.255.255.255"}, true},
		{"Zero IP", testIPv4_Zero, []string{testIPv4_Zero}, true},
		{"Private IP Class A", "10.1.2.3", []string{testCIDR_10}, true},
		{"Private IP Class B", "172.16.5.6", []string{testCIDR_172}, true},
		{"Private IP Class C", "192.168.100.1", []string{testCIDR_192Block}, true},

		// ========== 优先级测试 ==========
		{"Exact match before CIDR", testIPv4_4, []string{testIPv4_4, testCIDR_10}, true},
		{"Wildcard takes precedence", testIPv4_3, []string{testWildcard, "invalid-rule"}, true},

		// ========== 性能相关边界测试 ==========
		{"Many rules one match", testIPv4_3, []string{
			testCIDR_10, testCIDR_172, testCIDR_192Block, testIPv4_3,
		}, true},
		{"Many rules no match", "8.8.8.8", []string{
			testCIDR_10, testCIDR_172, testCIDR_192Block, testIPv4_Localhost,
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIPAllowed(tt.ip, tt.cidrList)
			assert.Equal(t, tt.want, got, "IP: %s, Rules: %v", tt.ip, tt.cidrList)
		})
	}

	// 测试公共的异常IP用例
	t.Run("Invalid IPs", func(t *testing.T) {
		for _, tc := range invalidIPCases {
			t.Run(tc.name, func(t *testing.T) {
				got := IsIPAllowed(tc.ip, []string{testCIDR_192})
				assert.False(t, got, "Invalid IP %s should not be allowed", tc.ip)
			})
		}
	})
}

func TestIsIPBlocked(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		blacklist []string
		want      bool
	}{
		// ========== 空列表测试 ==========
		{"Empty list does not block", testIPv4_1, []string{}, false},
		{"Nil list does not block", testIPv4_4, nil, false},

		// ========== 精确匹配测试 ==========
		{"Exact IP in blacklist", testIPv4_1, []string{testIPv4_1}, true},
		{"Exact IP not in blacklist", testIPv4_1, []string{"192.168.1.101"}, false},
		{"Multiple exact IPs one blocked", testIPv4_5, []string{testIPv4_3, testIPv4_5, testIPv4_6}, true},
		{"Multiple exact IPs none blocked", "10.0.0.6", []string{testIPv4_3, testIPv4_5, testIPv4_6}, false},

		// ========== CIDR 格式测试 ==========
		{"IP in CIDR blacklist", testIPv4_2, []string{testCIDR_192}, true},
		{"IP not in CIDR blacklist", "192.168.2.1", []string{testCIDR_192}, false},
		{"Multiple CIDRs one block", testIPv4_5, []string{testCIDR_192, testCIDR_10}, true},
		{"Multiple CIDRs none block", testIPv4_6, []string{testCIDR_192, testCIDR_10}, false},

		// ========== 特殊 IP 地址 ==========
		{"Localhost IPv4 in blacklist", testIPv4_Localhost, []string{testIPv4_Localhost}, true},
		{"Localhost in CIDR blacklist", testIPv4_Localhost, []string{testCIDR_127}, true},
		{"Broadcast IP in blacklist", "255.255.255.255", []string{"255.255.255.255"}, true},
		{"Zero IP in blacklist", testIPv4_Zero, []string{testIPv4_Zero}, true},
		{"Private IP Class A in blacklist", "10.1.2.3", []string{testCIDR_10}, true},
		{"Private IP Class B in blacklist", "172.16.5.6", []string{testCIDR_172}, true},
		{"Private IP Class C in blacklist", "192.168.100.1", []string{testCIDR_192Block}, true},

		// ========== 边界异常 ==========
		{"CIDR list contains empty string", testIPv4_2, []string{""}, false},
		{"CIDR list contains invalid CIDR", testIPv4_2, []string{"invalid-cidr"}, false},
		{"Malformed CIDR", testIPv4_3, []string{"192.168.1.0/99"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIPBlocked(tt.ip, tt.blacklist)
			assert.Equal(t, tt.want, got, "IP: %s, Rules: %v", tt.ip, tt.blacklist)
		})
	}

	// 测试公共的异常IP用例
	t.Run("Invalid IPs", func(t *testing.T) {
		for _, tc := range invalidIPCases {
			t.Run(tc.name, func(t *testing.T) {
				got := IsIPBlocked(tc.ip, []string{testCIDR_192})
				assert.False(t, got, "Invalid IP %s should return false", tc.ip)
			})
		}
	})

	// 测试IPv6用例
	t.Run("IPv6 Tests", func(t *testing.T) {
		for _, tc := range ipv6Cases {
			t.Run(tc.name, func(t *testing.T) {
				got := IsIPBlocked(tc.ip, []string{tc.pattern})
				assert.Equal(t, tc.want, got, "IPv6: %s, Pattern: %s", tc.ip, tc.pattern)
			})
		}
	})
}

func TestMatchIPPattern(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		pattern string
		want    bool
	}{
		// ========== 通配符测试 ==========
		{"Wildcard match", testIPv4_1, testWildcard, true},
		{"Exact match", testIPv4_1, testIPv4_1, true},
		{"No match", testIPv4_1, "192.168.1.101", false},

		// ========== CIDR 格式测试 ==========
		{"IP in CIDR", testIPv4_2, testCIDR_192, true},
		{"IP not in CIDR", "192.168.2.1", testCIDR_192, false},

		// ========== 异常和边界情况 ==========
		{"Empty IP", "", testCIDR_192, false},
		{"Invalid IP format", "invalid-ip", testCIDR_192, false},
		{"Invalid IP with numbers", "256.256.256.256", testCIDR_192, false},
		{"Malformed CIDR", testIPv4_3, "192.168.1.0/99", false},
		{"Incomplete IP", "192.168.1", testCIDR_192, false},
		{"IP with port", testIPv4_3 + ":8080", testCIDR_192, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchIPPattern(tt.ip, tt.pattern)
			assert.Equal(t, tt.want, got, "IP: %s, Pattern: %s", tt.ip, tt.pattern)
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"Private Class A", "10.1.2.3", true},
		{"Private Class B", "172.16.5.6", true},
		{"Private Class C", "192.168.100.1", true},
		{"Localhost", "127.0.0.1", true},
		{"Link-local", "169.254.0.1", true},
		{"IPv6 loopback", "::1", true},
		{"IPv6 unique local", "fc00::1", true},
		{"Not a private IP", "8.8.8.8", false},
		{"Invalid IP format", "invalid-ip", false},
		{"Invalid IP with numbers", "256.256.256.256", false},
		{"Empty IP", "", false},
		{"IP with port", "10.1.2.3:8080", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPrivateIP(tt.ip)
			assert.Equal(t, tt.want, got, "IP: %s", tt.ip)
		})
	}
}

func TestMatchIPInList(t *testing.T) {
	tests := []struct {
		name   string
		ip     string
		ipList []string
		want   bool
	}{
		// ========== 分隔符列表测试 ==========
		{"Semicolon separated IPs - match first", "192.168.0.10", []string{"192.168.0.10;192.168.0.1"}, true},
		{"Semicolon separated IPs - match second", "192.168.0.1", []string{"192.168.0.10;192.168.0.1"}, true},
		{"Comma separated IPs - match first", "192.168.0.1", []string{"192.168.0.1,192.168.0.56"}, true},
		{"Comma separated IPs - match second", "192.168.0.56", []string{"192.168.0.1,192.168.0.56"}, true},
		{"Semicolon separated IPs - no match", "192.168.0.99", []string{"192.168.0.10;192.168.0.1"}, false},
		{"Comma separated IPs - no match", "192.168.0.99", []string{"192.168.0.1,192.168.0.56"}, false},
		{"Semicolon separated multiple IPs with spaces", "10.0.0.5", []string{"10.0.0.1 ; 10.0.0.5 ; 10.0.0.10"}, true},
		{"Comma separated multiple IPs with spaces", "172.16.0.8", []string{"172.16.0.5, 172.16.0.8, 172.16.0.15"}, true},
		{"Semicolon separated mixed with CIDR", "192.168.1.50", []string{"192.168.1.0/24;10.0.0.0/8"}, true},
		{"Comma separated mixed with wildcard", "192.168.2.100", []string{"192.168.1.*,192.168.2.*"}, true},
		{"Semicolon separated mixed with IP range", "172.16.0.5", []string{"172.16.0.1-172.16.0.10;192.168.1.1-192.168.1.10"}, true},
		{"Comma separated with empty items", "192.168.0.10", []string{"192.168.0.10,,"}, true},
		{"Semicolon separated with empty items", "192.168.0.1", []string{";192.168.0.1;"}, true},
		{"Semicolon separated no match in range", "172.16.0.20", []string{"172.16.0.1-172.16.0.10;192.168.1.1-192.168.1.10"}, false},
		{"Comma separated with three IPs - match middle", "192.168.0.50", []string{"192.168.0.10,192.168.0.50,192.168.0.100"}, true},

		// ========== 换行符分隔测试 ==========
		{"Newline separated IPs - match first", "192.168.0.15", []string{"192.168.0.15\n192.168.0.1"}, true},
		{"Newline separated IPs - match second", "192.168.0.1", []string{"192.168.0.15\n192.168.0.1"}, true},
		{"Newline separated IPs - no match", "192.168.0.99", []string{"192.168.0.15\n192.168.0.1"}, false},
		{"Multiple newlines with empty lines", "10.0.0.5", []string{"10.0.0.1\n\n10.0.0.5\n\n10.0.0.10"}, true},
		{"CRLF newline separated", "172.16.0.8", []string{"172.16.0.5\r\n172.16.0.8\r\n172.16.0.15"}, true},

		// ========== 空格分隔测试 ==========
		{"Space separated IPs - match first", "192.168.0.15", []string{"192.168.0.15 192.168.0.17"}, true},
		{"Space separated IPs - match second", "192.168.0.17", []string{"192.168.0.15 192.168.0.17"}, true},
		{"Space separated IPs - no match", "192.168.0.99", []string{"192.168.0.15 192.168.0.17"}, false},
		{"Multiple spaces between IPs", "10.0.0.5", []string{"10.0.0.1  10.0.0.5  10.0.0.10"}, true},
		{"Tab separated IPs", "172.16.0.8", []string{"172.16.0.5\t172.16.0.8\t172.16.0.15"}, true},

		// ========== 混合分隔符测试 ==========
		{"Mixed semicolon and newline", "192.168.1.50", []string{"192.168.1.50;10.0.0.1\n172.16.0.1"}, true},
		{"Mixed comma and space", "192.168.2.100", []string{"192.168.1.100,192.168.2.100 192.168.3.100"}, true},
		{"Mixed all delimiters", "10.0.0.5", []string{"192.168.0.1;10.0.0.5\n172.16.0.1,192.168.0.2 192.168.0.3"}, true},
		{"Newline with CIDR", "192.168.1.50", []string{"192.168.1.0/24\n10.0.0.0/8"}, true},
		{"Space with wildcard", "192.168.2.100", []string{"192.168.1.* 192.168.2.*"}, true},
		{"Newline with IP range", "172.16.0.5", []string{"172.16.0.1-172.16.0.10\n192.168.1.1-192.168.1.10"}, true},

		// ========== 精确匹配测试 ==========
		{"Exact match single", testIPv4_1, []string{testIPv4_1}, true},
		{"Exact match multiple", testIPv4_2, []string{testIPv4_1, testIPv4_2, testIPv4_3}, true},
		{"Exact no match", testIPv4_1, []string{testIPv4_2, testIPv4_3}, false},
		{"Empty list", testIPv4_1, []string{}, false},
		{"Nil list", testIPv4_1, nil, false},

		// ========== CIDR 匹配测试 ==========
		{"CIDR match /24", testIPv4_2, []string{testCIDR_192}, true},
		{"CIDR no match /24", "192.168.2.1", []string{testCIDR_192}, false},
		{"CIDR match /8", "10.255.255.1", []string{testCIDR_10}, true},
		{"CIDR match /12", "172.31.255.1", []string{testCIDR_172}, true},
		{"CIDR boundary lower", "192.168.1.0", []string{testCIDR_192}, true},
		{"CIDR boundary upper", "192.168.1.255", []string{testCIDR_192}, true},
		{"Outside CIDR", "192.168.0.255", []string{testCIDR_192}, false},

		// ========== IP 范围匹配测试 ==========
		{"Range match start", "172.16.0.1", []string{"172.16.0.1-172.16.0.10"}, true},
		{"Range match end", "172.16.0.10", []string{"172.16.0.1-172.16.0.10"}, true},
		{"Range match middle", "172.16.0.5", []string{"172.16.0.1-172.16.0.10"}, true},
		{"Range no match below", "172.16.0.0", []string{"172.16.0.1-172.16.0.10"}, false},
		{"Range no match above", "172.16.0.11", []string{"172.16.0.1-172.16.0.10"}, false},
		{"Range with spaces", "172.16.0.5", []string{" 172.16.0.1 - 172.16.0.10 "}, true},
		{"Range invalid format", "172.16.0.5", []string{"172.16.0.1-"}, false},
		{"Range malformed", "172.16.0.5", []string{"invalid-range"}, false},

		// ========== 通配符匹配测试 ==========
		{"Wildcard last octet", "192.168.2.100", []string{"192.168.2.*"}, true},
		{"Wildcard last octet no match", "192.168.3.100", []string{"192.168.2.*"}, false},
		{"Wildcard two octets", "192.168.50.100", []string{"192.168.*.*"}, true},
		{"Wildcard all octets", "1.2.3.4", []string{"*.*.*.*"}, true},
		{"Wildcard middle octet", "192.100.1.50", []string{"192.*.1.50"}, true},
		{"Wildcard middle no match", "192.100.2.50", []string{"192.*.1.50"}, false},
		{"Wildcard first octet", "50.168.1.1", []string{"*.168.1.1"}, true},

		// ========== 混合规则测试 ==========
		{"Mixed exact and CIDR", testIPv4_1, []string{testIPv4_2, testCIDR_192}, true},
		{"Mixed CIDR and range", "172.16.0.5", []string{testCIDR_192, "172.16.0.1-172.16.0.10"}, true},
		{"Mixed all types", "192.168.2.50", []string{testIPv4_1, testCIDR_10, "172.16.0.1-172.16.0.10", "192.168.2.*"}, true},
		{"Mixed none match", "8.8.8.8", []string{testIPv4_1, testCIDR_192, "172.16.0.1-172.16.0.10", "10.0.0.*"}, false},

		// ========== 边界测试 ==========
		{"List with empty string", testIPv4_1, []string{""}, false},
		{"List with invalid entries", testIPv4_1, []string{"invalid", "bad-cidr"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchIPInList(tt.ip, tt.ipList)
			assert.Equal(t, tt.want, got, "IP: %s, List: %v", tt.ip, tt.ipList)
		})
	}

	// 测试公共的异常IP用例
	t.Run("Invalid IPs", func(t *testing.T) {
		for _, tc := range invalidIPCases {
			t.Run(tc.name, func(t *testing.T) {
				got := MatchIPInList(tc.ip, []string{testCIDR_192})
				assert.False(t, got, "Invalid IP %s should not match", tc.ip)
			})
		}
	})

	// 测试IPv6用例
	t.Run("IPv6 Tests", func(t *testing.T) {
		tests := []struct {
			name   string
			ip     string
			ipList []string
			want   bool
		}{
			{"IPv6 exact match", testIPv6_1, []string{testIPv6_1}, true},
			{"IPv6 CIDR match", testIPv6_2, []string{testCIDR_IPv6}, true},
			{"IPv6 no match", "2001:db9::1", []string{testIPv6_1, testCIDR_IPv6}, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := MatchIPInList(tt.ip, tt.ipList)
				assert.Equal(t, tt.want, got, "IPv6: %s, List: %v", tt.ip, tt.ipList)
			})
		}
	})
}

func TestIsIPInRange(t *testing.T) {
	tests := []struct {
		name  string
		ip    string
		start string
		end   string
		want  bool
	}{
		// ========== 正常范围测试 ==========
		{"In range start", "172.16.0.1", "172.16.0.1", "172.16.0.10", true},
		{"In range end", "172.16.0.10", "172.16.0.1", "172.16.0.10", true},
		{"In range middle", "172.16.0.5", "172.16.0.1", "172.16.0.10", true},
		{"Below range", "172.16.0.0", "172.16.0.1", "172.16.0.10", false},
		{"Above range", "172.16.0.11", "172.16.0.1", "172.16.0.10", false},

		// ========== 单IP范围测试 ==========
		{"Single IP range match", "192.168.1.1", "192.168.1.1", "192.168.1.1", true},
		{"Single IP range no match", "192.168.1.2", "192.168.1.1", "192.168.1.1", false},

		// ========== 大范围测试 ==========
		{"Large range start", "10.0.0.0", "10.0.0.0", "10.255.255.255", true},
		{"Large range end", "10.255.255.255", "10.0.0.0", "10.255.255.255", true},
		{"Large range middle", "10.128.50.100", "10.0.0.0", "10.255.255.255", true},
		{"Outside large range", "11.0.0.1", "10.0.0.0", "10.255.255.255", false},

		// ========== 跨网段测试 ==========
		{"Cross subnet", "192.168.50.1", "192.168.1.1", "192.168.100.1", true},
		{"Cross class C", "192.168.255.255", "192.168.0.0", "192.169.0.0", true},

		// ========== 边界值测试 ==========
		{"All zeros", "0.0.0.0", "0.0.0.0", "255.255.255.255", true},
		{"All 255s", "255.255.255.255", "0.0.0.0", "255.255.255.255", true},
		{"Broadcast in range", "192.168.1.255", "192.168.1.0", "192.168.1.255", true},

		// ========== IPv6 测试 (应该返回 false) ==========
		{"IPv6 not supported", "2001:db8::1", "2001:db8::1", "2001:db8::10", false},

		// ========== 异常情况 ==========
		{"Invalid IP", "invalid", "172.16.0.1", "172.16.0.10", false},
		{"Invalid start", "172.16.0.5", "invalid", "172.16.0.10", false},
		{"Invalid end", "172.16.0.5", "172.16.0.1", "invalid", false},
		{"Empty IP", "", "172.16.0.1", "172.16.0.10", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := ParseIPHelper(tt.ip)
			start := ParseIPHelper(tt.start)
			end := ParseIPHelper(tt.end)
			got := IsIPInRange(ip, start, end)
			assert.Equal(t, tt.want, got, "IP: %s, Range: %s-%s", tt.ip, tt.start, tt.end)
		})
	}
}

func TestMatchIPWithWildcard(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		pattern string
		want    bool
	}{
		// ========== 单段通配符测试 ==========
		{"Last octet wildcard match", "192.168.1.100", "192.168.1.*", true},
		{"Last octet wildcard no match", "192.168.2.100", "192.168.1.*", false},
		{"First octet wildcard match", "50.168.1.1", "*.168.1.1", true},
		{"First octet wildcard no match", "50.168.1.2", "*.168.1.1", false},
		{"Middle octet wildcard match", "192.50.1.1", "192.*.1.1", true},
		{"Middle octet wildcard no match", "192.50.2.1", "192.*.1.2", false},

		// ========== 多段通配符测试 ==========
		{"Two octets wildcard", "192.168.50.100", "192.168.*.*", true},
		{"Three octets wildcard", "192.50.100.200", "192.*.*.*", true},
		{"All octets wildcard", "1.2.3.4", "*.*.*.*", true},
		{"First two octets wildcard", "50.100.1.1", "*.*.1.1", true},

		// ========== 无通配符测试 ==========
		{"No wildcard exact match", "192.168.1.1", "192.168.1.1", true},
		{"No wildcard no match", "192.168.1.1", "192.168.1.2", false},

		// ========== 边界值测试 ==========
		{"Wildcard with zero", "0.0.0.0", "0.0.0.*", true},
		{"Wildcard with 255", "255.255.255.255", "255.255.255.*", true},
		{"All wildcards", "1.2.3.4", "*.*.*.*", true},

		// ========== 异常情况测试 ==========
		{"Different length IP", "192.168.1", "192.168.1.*", false},
		{"Different length pattern", "192.168.1.1", "192.168.*", false},
		{"Empty IP", "", "192.168.1.*", false},
		{"Empty pattern", "192.168.1.1", "", false},
		{"IPv6 format", "2001:db8::1", "2001:db8::*", false}, // 因为分隔符不是'.'
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchIPWithWildcard(tt.ip, tt.pattern)
			assert.Equal(t, tt.want, got, "IP: %s, Pattern: %s", tt.ip, tt.pattern)
		})
	}
}

func TestValidateIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		// ========== 有效 IPv4 ==========
		{"Valid IPv4", "192.168.1.1", false},
		{"Valid IPv4 with zeros", "10.0.0.1", false},
		{"Valid IPv4 max values", "255.255.255.255", false},
		{"Valid IPv4 min values", "0.0.0.0", false},
		{"Valid IPv4 localhost", "127.0.0.1", false},

		// ========== 有效 IPv6 ==========
		{"Valid IPv6", "2001:db8::1", false},
		{"Valid IPv6 loopback", "::1", false},
		{"Valid IPv6 any", "::", false},
		{"Valid IPv6 full", "2001:0db8:0000:0000:0000:0000:0000:0001", false},

		// ========== 无效格式 ==========
		{"Invalid format", "invalid-ip", true},
		{"Invalid numbers", "256.256.256.256", true},
		{"Incomplete IP", "192.168.1", true},
		{"Too many octets", "192.168.1.1.1", true},
		{"Empty string", "", true},
		{"With port", "192.168.1.1:8080", true},
		{"With CIDR", "192.168.1.0/24", true},
		{"Letters in IP", "192.168.1.abc", true},
		{"Negative numbers", "192.-168.1.1", true},
	}

	ipBase := &IPBase{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ipBase.ValidateIP(tt.ip)
			if tt.wantErr {
				assert.Error(t, err, "IP: %s", tt.ip)
			} else {
				assert.NoError(t, err, "IP: %s", tt.ip)
			}
		})
	}
}

// ParseIPHelper 辅助函数,用于测试 IsIPInRange
func ParseIPHelper(ip string) net.IP {
	return net.ParseIP(ip)
}
