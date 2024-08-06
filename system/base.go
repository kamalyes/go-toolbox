/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-06 15:52:39
 * @FilePath: \go-toolbox\system\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package system

import (
	"net"
	"os"
	"strings"

	"github.com/kamalyes/go-toolbox/random"
)

// 获取主机名函数
func SafeGetHostName() string {
	output, err := os.Hostname()
	if err != nil {
		output = random.FRandAlphaString(8)
	}
	output = strings.ReplaceAll(output, "-", "_")
	return string(output)
}

func ServerIP() (string, string) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", ""
	}
	var internalIP, externalIP string
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipAddr := ipNet.IP.String()
				if strings.HasPrefix(ipAddr, "10.") || strings.HasPrefix(ipAddr, "192.") || strings.HasPrefix(ipAddr, "127.") {
					internalIP = ipAddr
				} else {
					externalIP = ipAddr
				}
			}

		}
	}

	return externalIP, internalIP
}
