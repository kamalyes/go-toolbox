/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-18 17:22:25
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-21 16:09:51
 * @FilePath: \go-toolbox\pkg\netx\ip.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package netx

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// GetLocalInterfaceIPs 查询本机网卡所有IP
func GetLocalInterfaceIPs() ([]net.IP, error) {
	var localIPs []net.IP
	interfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("error getting network interfaces: %w", err)
	}

	for _, address := range interfaceAddresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			localIPs = append(localIPs, ipNet.IP)
		}
	}
	return localIPs, nil
}

// GetPrivateIP 获取私有 IP
func GetPrivateIP() (string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range netInterfaces {
		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addresses {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.IsPrivate() {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("未找到私有 IP")
}

// GetLocalInterfaceIPAndExternalIP 返回本地网卡对应的外部和内部 IP 地址
func GetLocalInterfaceIPAndExternalIP(urls ...string) (privateIP string, publicIP string, err error) {
	if privateIP, err = GetPrivateIP(); err != nil {
		return
	}
	publicIP, err = GetConNetPublicIP(urls...)
	return privateIP, publicIP, err
}

// GetConNetPublicIP 联网获取本机公网 IP
func GetConNetPublicIP(urls ...string) (string, error) {
	externalIPServiceURL := "http://myexternalip.com/raw"
	if len(urls) > 0 {
		externalIPServiceURL = urls[0]
	}

	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}

	response, err := httpClient.Get(externalIPServiceURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetClientIP 从 HTTP 请求中提取客户端 IP
// 支持 IPv4 和 IPv6（包括 X-Forwarded-For / X-Real-IP 中的裸 IPv6 地址）
func GetClientIP(r *http.Request) string {
	// 1. 尝试从 X-Forwarded-For 头获取
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For 可能包含多个 IP，取第一个
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return NormalizeIP(strings.TrimSpace(ips[0]))
		}
	}

	// 2. 尝试从 X-Real-IP 头获取
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return NormalizeIP(strings.TrimSpace(xri))
	}

	// 3. 从 RemoteAddr 获取（IPv6 格式为 [::1]:port，SplitHostPort 可正确处理）
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return NormalizeIP(r.RemoteAddr)
	}
	return ip
}

// JoinHostPort 将 host 和 port 拼接为网络地址字符串
// 自动处理 IPv6 地址的方括号包裹：
//   - JoinHostPort("127.0.0.1", 8080) → "127.0.0.1:8080"
//   - JoinHostPort("::1", 8080)       → "[::1]:8080"
//   - JoinHostPort("", 8080)          → ":8080"
func JoinHostPort(host string, port int) string {
	return net.JoinHostPort(host, fmt.Sprintf("%d", port))
}

// NormalizeIP 规范化 IP 地址字符串
// 处理以下 IPv6 场景：
//   - 去除方括号：[::1] → ::1
//   - 带端口无方括号：::1:8080 → ::1（尝试识别末尾端口）
//   - 纯 IPv6 地址：2001:db8::1 → 2001:db8::1（不变）
//   - IPv4 地址：192.168.1.1 → 192.168.1.1（不变）
func NormalizeIP(addr string) string {
	if addr == "" {
		return addr
	}

	// 去除方括号（来自 URL 格式的 IPv6，如 [::1]:8080）
	addr = strings.TrimPrefix(addr, "[")
	if idx := strings.LastIndex(addr, "]"); idx >= 0 {
		addr = addr[:idx]
	}

	// 尝试解析为 net.IP，成功说明是纯 IP 地址（已去除端口）
	if ip := net.ParseIP(addr); ip != nil {
		return addr
	}

	// 尝试 SplitHostPort（处理 [::1]:8080 或 192.168.1.1:8080）
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		// 去除可能残留的方括号
		host = strings.TrimPrefix(host, "[")
		host = strings.TrimSuffix(host, "]")
		return host
	}

	// SplitHostPort 失败且不含方括号，可能是裸 IPv6 地址
	// 尝试从右向左找到最后一个冒号，尝试剥离端口部分
	if strings.Contains(addr, ":") && !strings.Contains(addr, ".") {
		lastColon := strings.LastIndex(addr, ":")
		candidate := addr[:lastColon]
		if net.ParseIP(candidate) != nil {
			return candidate
		}
	}

	// 无法解析，原样返回
	return addr
}

// NormalizeListenAddr 规范化 net.Listen 地址，确保 IPv6 地址被方括号包裹
// net.Listen 要求 IPv6 地址格式为 [::1]:port，但配置可能传入裸 IPv6:port
// 此函数自动检测并补齐方括号：
//   - "::1:5081"            → "[::1]:5081"
//   - "2001:db8::1:5081"    → "[2001:db8::1]:5081"
//   - "[::1]:5081"          → "[::1]:5081"（已规范，不变）
//   - "127.0.0.1:5081"      → "127.0.0.1:5081"（IPv4，不变）
//   - ":5081"               → ":5081"（仅端口，不变）
func NormalizeListenAddr(addr string) string {
	// 已经有方括号，无需处理
	if strings.Contains(addr, "[") {
		return addr
	}

	// 不含冒号，不是 IPv6，无需处理
	if !strings.Contains(addr, ":") {
		return addr
	}

	// 尝试直接解析，如果已经是合法的 host:port 则无需处理
	// （IPv4:port 或 :port 会在此成功）
	if _, err := net.ResolveTCPAddr("tcp", addr); err == nil {
		return addr
	}

	// 含冒号但无法解析为 TCP 地址，可能是裸 IPv6 地址
	// 尝试从右向左找到最后一个冒号，剥离端口后尝试解析为 IP
	lastColon := strings.LastIndex(addr, ":")
	if lastColon <= 0 {
		return addr
	}

	hostPart := addr[:lastColon]
	portPart := addr[lastColon+1:]

	if net.ParseIP(hostPart) != nil {
		return "[" + hostPart + "]:" + portPart
	}

	// 无法识别为 IPv6，原样返回
	return addr
}
