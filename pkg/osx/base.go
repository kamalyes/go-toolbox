/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 20:09:14
 * @FilePath: \go-toolbox\pkg\osx\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package osx

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/go-toolbox/pkg/stringx"
)

// GetHostName 获取主机名，如果失败则返回错误
func GetHostName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("无法获取主机名: %v", err)
	}
	return hostname, nil
}

// 获取主机名函数
func SafeGetHostName() string {
	output, err := GetHostName()
	if err != nil || output == "" {
		// 如果获取主机名失败或返回空字符串，则生成随机字符串
		return stringx.ReplaceSpecialChars(random.FRandAlphaString(8), 'x')
	}
	return stringx.ReplaceSpecialChars(output, 'x')
}

// HashUnixMicroCipherText
func HashUnixMicroCipherText() string {
	var (
		nowUnixMicro = time.Now().UnixMicro()
		hostName     = SafeGetHostName()
		randStr      = random.RandString(10, 4)
		plainText    = fmt.Sprintf("%s%s%d", hostName, randStr, nowUnixMicro)
		cipherText   = stringx.CalculateMD5Hash(plainText)
	)
	return cipherText
}

// getNetworkInterfaces 返回所有网络接口的 IP 地址
func getNetworkInterfaces() ([]net.IP, error) {
	var ips []net.IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("error getting network interfaces: %w", err)
	}

	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			ips = append(ips, ipNet.IP)
		}
	}
	return ips, nil
}

// separateIPs 将 IP 地址列表分为内部和外部 IP 地址
func separateIPs(ips []net.IP) (internalIPs, externalIPs []string) {
	for _, ip := range ips {
		if ip.IsPrivate() {
			internalIPs = append(internalIPs, ip.String())
		} else {
			externalIPs = append(externalIPs, ip.String())
		}
	}
	return internalIPs, externalIPs
}

// GetLocalInterfaceIeIp 返回本地网卡对应的外部和内部 IP 地址
func GetLocalInterfaceIeIp() (string, string, error) {
	ips, err := getNetworkInterfaces()
	if err != nil {
		return "", "", err
	}

	internalIPs, externalIPs := separateIPs(ips)

	var internalIP, externalIP string
	if len(internalIPs) > 0 {
		internalIP = internalIPs[0]
	}
	if len(externalIPs) > 0 {
		externalIP = externalIPs[0]
	}

	return externalIP, internalIP, nil
}

// GetLocalInterfaceIps 查询本机网卡所有IP
func GetLocalInterfaceIps() ([]net.IP, error) {
	ips, err := getNetworkInterfaces()
	if err != nil {
		return nil, err
	}
	return ips, nil
}

// GetClientPublicIP 从HTTP请求中获取客户端的公网IP地址。
// 如果无法确定公网IP，则返回空字符串和错误（如果适用）。
func GetClientPublicIP(r *http.Request) (string, error) {
	// 检查HTTP头部中的IP地址
	headers := []string{
		r.Header.Get("X-Forwarded-For"),
		r.Header.Get("X-Real-Ip"),
	}

	for _, header := range headers {
		// 分割并修剪空格，获取第一个IP地址
		ips := strings.Fields(header)
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			// 检查IP地址是否有效且不是本地地址
			if ip != "" && !stringx.HasLocalIP(ip) {
				return ip, nil
			}
		}
	}

	// 回退到使用请求的远程地址
	remoteAddr, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		// 这里我们返回错误，以便调用者可以处理它
		return "", fmt.Errorf("failed to parse remote address: %w", err)
	}

	// 检查远程地址是否不是本地地址
	if !stringx.HasLocalIP(remoteAddr) {
		return remoteAddr, nil
	}

	// 如果没有找到有效的公网IP，返回错误
	return "", fmt.Errorf("no valid public IP address found")
}

// GetConNetPublicIp 联网获取本机公网 IP
func GetConNetPublicIp(urls ...string) (string, error) {
	url := "http://myexternalip.com/raw"
	if len(urls) > 0 {
		url = urls[0]
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
