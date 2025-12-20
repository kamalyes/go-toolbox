/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 15:55:00
 * @FilePath: \go-toolbox\pkg\metadata\types.go
 * @Description: HTTP 请求元数据结构定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package metadata

// RequestMetadata HTTP 请求元数据
type RequestMetadata struct {
	// 基础请求信息
	UserAgent     string `json:"user_agent"`
	ClientIP      string `json:"client_ip"`
	RemoteAddr    string `json:"remote_addr"`
	RequestURI    string `json:"request_uri"`
	QueryString   string `json:"query_string"`
	RequestMethod string `json:"request_method"`
	RequestHost   string `json:"request_host"`

	// 来源信息
	Origin  string `json:"origin"`
	Referer string `json:"referer"`

	// 代理和转发信息
	XForwardedFor   string `json:"x_forwarded_for"`
	XRealIP         string `json:"x_real_ip"`
	XForwardedProto string `json:"x_forwarded_proto"`
	XForwardedHost  string `json:"x_forwarded_host"`
	XForwardedPort  string `json:"x_forwarded_port"`

	// 客户端偏好信息
	AcceptLanguage string `json:"accept_language"`
	AcceptEncoding string `json:"accept_encoding"`
	Accept         string `json:"accept"`

	// WebSocket 协议信息
	SecWebSocketKey        string `json:"sec_websocket_key"`
	SecWebSocketVersion    string `json:"sec_websocket_version"`
	SecWebSocketProtocol   string `json:"sec_websocket_protocol"`
	SecWebSocketExtensions string `json:"sec_websocket_extensions"`

	// 连接信息
	Connection string `json:"connection"`
	Upgrade    string `json:"upgrade"`

	// CDN 和安全信息
	CFRay          string `json:"cf_ray"`
	CFConnectingIP string `json:"cf_connecting_ip"`
	CFIPCountry    string `json:"cf_ipcountry"`
	XRequestID     string `json:"x_request_id"`
	XCorrelationID string `json:"x_correlation_id"`

	// 缓存和条件请求
	CacheControl    string `json:"cache_control"`
	IfNoneMatch     string `json:"if_none_match"`
	IfModifiedSince string `json:"if_modified_since"`

	// 客户端提示信息（Client Hints）
	SecCHUA         string `json:"sec_ch_ua"`
	SecCHUAMobile   string `json:"sec_ch_ua_mobile"`
	SecCHUAPlatform string `json:"sec_ch_ua_platform"`
	DNT             string `json:"dnt"`

	// TLS 信息
	Protocol       string `json:"protocol"`
	TLSVersion     uint16 `json:"tls_version,omitempty"`
	TLSCipherSuite uint16 `json:"tls_cipher_suite,omitempty"`
	TLSServerName  string `json:"tls_server_name,omitempty"`

	// User-Agent 解析结果（来自 useragent.Parse）
	Browser        string `json:"browser,omitempty"`         // 浏览器名称: Chrome, Firefox, Safari
	BrowserVersion string `json:"browser_version,omitempty"` // 浏览器版本号
	OS             string `json:"os,omitempty"`              // 操作系统: Windows, macOS, Android, iOS
	OSVersion      string `json:"os_version,omitempty"`      // 操作系统版本号
	Device         string `json:"device,omitempty"`          // 设备名称
	DeviceType     string `json:"device_type,omitempty"`     // 设备类型: mobile/tablet/desktop/bot
	DeviceVendor   string `json:"device_vendor,omitempty"`   // 设备厂商: Apple, Samsung, Huawei
	IsBot          bool   `json:"is_bot"`                    // 是否为爬虫/机器人
	BotName        string `json:"bot_name,omitempty"`        // 爬虫名称
	IsMobile       bool   `json:"is_mobile"`                 // 是否为移动设备
	IsTablet       bool   `json:"is_tablet"`                 // 是否为平板设备
}
