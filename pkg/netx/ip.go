/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-18 17:22:25
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-09-18 17:55:15
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
