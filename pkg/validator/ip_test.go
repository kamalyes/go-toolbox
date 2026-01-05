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
		{"Empty IP", "", []string{testCIDR_192}, false},
		{"Invalid IP format", "invalid-ip", []string{testCIDR_192}, false},
		{"Invalid IP with numbers", "256.256.256.256", []string{testCIDR_192}, false},
		{"CIDR list contains empty string", testIPv4_2, []string{""}, false},
		{"CIDR list contains invalid CIDR", testIPv4_2, []string{"invalid-cidr"}, false},
		{"IP equals CIDR but CIDR invalid", testIPv4_2, []string{testIPv4_2 + "/33"}, false}, // 33不是合法掩码
		{"IP equals CIDR string but IP invalid", "999.999.999.999", []string{"999.999.999.999"}, false},
		{"Malformed CIDR", testIPv4_3, []string{"192.168.1.0/99"}, false},
		{"Incomplete IP", "192.168.1", []string{testCIDR_192}, false},
		{"IP with port", testIPv4_3 + ":8080", []string{testCIDR_192}, false},

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

		// ========== IPv6 测试 ==========
		{"IPv6 exact block", testIPv6_1, []string{testIPv6_1}, true},
		{"IPv6 exact no block", testIPv6_1, []string{"2001:db8::2"}, false},
		{"IPv6 CIDR block", testIPv6_2, []string{testCIDR_IPv6}, true},
		{"IPv6 CIDR no block", "2001:db9::1", []string{testCIDR_IPv6}, false},

		// ========== 异常和边界情况 ==========
		{"Empty IP", "", []string{testCIDR_192}, false},
		{"Invalid IP format", "invalid-ip", []string{testCIDR_192}, false},
		{"Invalid IP with numbers", "256.256.256.256", []string{testCIDR_192}, false},
		{"CIDR list contains empty string", testIPv4_2, []string{""}, false},
		{"CIDR list contains invalid CIDR", testIPv4_2, []string{"invalid-cidr"}, false},
		{"IP equals CIDR but CIDR invalid", testIPv4_2, []string{testIPv4_2 + "/33"}, false},
		{"IP equals CIDR string but IP invalid", "999.999.999.999", []string{"999.999.999.999"}, false},
		{"Malformed CIDR", testIPv4_3, []string{"192.168.1.0/99"}, false},
		{"Incomplete IP", "192.168.1", []string{testCIDR_192}, false},
		{"IP with port", testIPv4_3 + ":8080", []string{testCIDR_192}, false},

		// ========== 特殊 IP 地址 ==========
		{"Localhost IPv4 in blacklist", testIPv4_Localhost, []string{testIPv4_Localhost}, true},
		{"Localhost in CIDR blacklist", testIPv4_Localhost, []string{testCIDR_127}, true},
		{"Broadcast IP in blacklist", "255.255.255.255", []string{"255.255.255.255"}, true},
		{"Zero IP in blacklist", testIPv4_Zero, []string{testIPv4_Zero}, true},
		{"Private IP Class A in blacklist", "10.1.2.3", []string{testCIDR_10}, true},
		{"Private IP Class B in blacklist", "172.16.5.6", []string{testCIDR_172}, true},
		{"Private IP Class C in blacklist", "192.168.100.1", []string{testCIDR_192Block}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIPBlocked(tt.ip, tt.blacklist)
			assert.Equal(t, tt.want, got, "IP: %s, Rules: %v", tt.ip, tt.blacklist)
		})
	}
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
