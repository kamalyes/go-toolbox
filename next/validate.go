/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-17 16:05:19
 * @FilePath: \go-toolbox\next\validate.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package next

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/thinkeridea/go-extend/exnet"
)

const (
	// GET_PUBLIC_IP_URL 查询公网IP的URL
	GET_PUBLIC_IP_URL = "http://myexternalip.com/raw"
)

var (
	v *validator.Validate
)

// GetClientIP
/**
 *  @Description: 获取用户真实IP
 *  @param r
 *  @return ip
 */
func GetClientIP(r *http.Request) (ip string) {
	ip = exnet.ClientPublicIP(r)
	if ip == "" {
		ip = exnet.ClientIP(r)
	}
	return
}

// GetLocalIp
/**
 *  @Description: 查询本机内网IP
 *  @return ips ip列表
 *  @return err 错误
 */
func GetLocalIp() (ips []string, err error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("get ip interfaces error:", err)
		return
	}

	for _, i := range netInterfaces {
		address, errRet := i.Addrs()
		if errRet != nil {
			continue
		}

		for _, addr := range address {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				if ip.IsGlobalUnicast() {
					ips = append(ips, ip.String())
				}
			}
		}
	}
	return
}

// GetPublicIP
/**
 *  @Description: 获取本机公网ip
 *  @return ip
 *  @return err
 */
func GetPublicIP() (ip string, err error) {
	resp, errHttp := http.Get(GET_PUBLIC_IP_URL)
	if errHttp != nil {
		return ip, errHttp
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	// 读取标准输出
	body, errRead := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 || errRead != nil {
		err = errors.New("请求失败")
		return ip, err
	}
	ip = string(body)
	return ip, nil
}

// IsIpAddress
/**
 * @Description: 判断是否为合法的ip地址
 * @param ip ip地址
 * @return ok 是否合法
 */
func IsIpAddress(ip string) (ok bool) {
	if v == nil {
		v = validator.New()
	}
	errs := v.Var(ip, "ip")
	if errs != nil {
		return false
	} else {
		return true
	}
}

// IsIPv4
/**
 * @Description: 判断是否为合法的ipv4地址
 * @param ip ip地址
 * @return ok 是否合法
 */
func IsIPv4(ip string) (ok bool) {
	if v == nil {
		v = validator.New()
	}
	errs := v.Var(ip, "ipv4")
	if errs != nil {
		return false
	} else {
		return true
	}
}

// IsIPv6
/**
 * @Description: 判断是否为合法的ipv6地址
 * @param ip ip地址
 * @return ok 是否合法
 */
func IsIPv6(ip string) (ok bool) {
	if v == nil {
		v = validator.New()
	}
	errs := v.Var(ip, "ipv6")
	if errs != nil {
		return false
	} else {
		return true
	}
}

// HasLocalIPAddr 检测 IP 地址字符串是否是内网地址
func HasLocalIPAddr(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && HasLocalIP(parsedIP)
}

// HasLocalIP 检测 IP 地址是否是内网地址
func HasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true // 回环地址是内网地址
	}

	if ip4 := ip.To4(); ip4 != nil {
		// 检查 IPv4 内网地址
		return isPrivateIPv4(ip4)
	}

	// 检查 IPv6 内网地址
	return isPrivateIPv6(ip)
}

// isPrivateIPv4 检查 IPv4 地址是否是内网地址
func isPrivateIPv4(ip4 net.IP) bool {
	switch {
	case ip4[0] == 10: // 10.0.0.0/8
		return true
	case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31: // 172.16.0.0/12
		return true
	case ip4[0] == 169 && ip4[1] == 254: // 169.254.0.0/16
		return true
	case ip4[0] == 192 && ip4[1] == 168: // 192.168.0.0/16
		return true
	default:
		return false
	}
}

// isPrivateIPv6 检查 IPv6 地址是否是内网地址
func isPrivateIPv6(ip net.IP) bool {
	return (ip[0] == 0xFE && (ip[1]&0xC0) == 0x80) || (ip[0] == 0xFC || ip[0] == 0xFD)
}
