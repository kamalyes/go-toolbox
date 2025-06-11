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

import "net"

func IsIPAllowed(ip string, cidrList []string) bool {
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	for _, cidr := range cidrList {
		if ip == cidr {
			return true
		}
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
