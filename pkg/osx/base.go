/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-16 19:51:55
 * @FilePath: \go-toolbox\pkg\osx\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/go-toolbox/pkg/stringx"
)

func PanicGetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		// 如果获取主机名时出现错误，则返回一个带有错误信息的错误
		panic(fmt.Errorf("无法获取主机名: %v", err))
	}
	return hostname
}

// 获取主机名函数
func SafeGetHostName() string {
	output, err := os.Hostname()
	if err != nil {
		output = random.FRandAlphaString(8)
	}
	return stringx.ReplaceSpecialChars(output, 'x')
}

// HashUnixMicroCipherText
func HashUnixMicroCipherText() string {
	var (
		nowUnixMicro = time.Now().UnixMicro()
		hostName     = PanicGetHostName()
		randStr      = random.RandString(10, 4)
		plainText    = fmt.Sprintf("%s%s%d", hostName, randStr, nowUnixMicro)
		cipherText   = stringx.CalculateMD5Hash(plainText)
	)
	return cipherText
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

// GetCurrentPath 获取当前工作目录的路径
func GetCurrentPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		// 返回一个格式化的错误，而不是忽略它
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}
	return dir, nil
}
