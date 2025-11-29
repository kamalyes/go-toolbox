/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-09 01:00:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 01:00:55
 * @FilePath: \go-toolbox\pkg\convert\ip.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package convert

import (
	"errors"
	"math"
	"net"
)

// IP2Long 把net.IP转为数值
func IP2Long(ip net.IP) (uint, error) {
	b := ip.To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}

	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

// Long2IP 把数值转为net.IP
func Long2IP(i uint) (net.IP, error) {
	// 使用uint32(math.MaxUint32)避免32位架构上的常量溢出
	if uint64(i) > uint64(uint32(math.MaxUint32)) {
		return nil, errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip, nil
}
