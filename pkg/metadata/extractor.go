/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 16:50:00
 * @FilePath: \go-toolbox\pkg\metadata\extractor.go
 * @Description: HTTP 请求元数据提取器
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package metadata

import (
	"net/http"

	"github.com/kamalyes/go-toolbox/pkg/netx"
	"github.com/kamalyes/go-toolbox/pkg/useragent"
)

// ExtractRequestMetadata 从 HTTP 请求中提取所有元数据
func ExtractRequestMetadata(r *http.Request) *RequestMetadata {
	userAgent := r.Header.Get("User-Agent")

	// 构建相对路径的 RequestURI
	requestURI := r.URL.Path
	if r.URL.RawQuery != "" {
		requestURI += "?" + r.URL.RawQuery
	}

	metadata := &RequestMetadata{
		// 基础请求信息
		UserAgent:     userAgent,
		ClientIP:      netx.GetClientIP(r),
		RemoteAddr:    r.RemoteAddr,
		RequestURI:    requestURI,
		QueryString:   r.URL.RawQuery,
		RequestMethod: r.Method,
		RequestHost:   r.Host,

		// 来源信息
		Origin:  r.Header.Get("Origin"),
		Referer: r.Header.Get("Referer"),

		// 代理和转发信息
		XForwardedFor:   r.Header.Get("X-Forwarded-For"),
		XRealIP:         r.Header.Get("X-Real-IP"),
		XForwardedProto: r.Header.Get("X-Forwarded-Proto"),
		XForwardedHost:  r.Header.Get("X-Forwarded-Host"),
		XForwardedPort:  r.Header.Get("X-Forwarded-Port"),

		// 客户端偏好信息
		AcceptLanguage: r.Header.Get("Accept-Language"),
		AcceptEncoding: r.Header.Get("Accept-Encoding"),
		Accept:         r.Header.Get("Accept"),

		// WebSocket 协议信息
		SecWebSocketKey:        r.Header.Get("Sec-WebSocket-Key"),
		SecWebSocketVersion:    r.Header.Get("Sec-WebSocket-Version"),
		SecWebSocketProtocol:   r.Header.Get("Sec-WebSocket-Protocol"),
		SecWebSocketExtensions: r.Header.Get("Sec-WebSocket-Extensions"),

		// 连接信息
		Connection: r.Header.Get("Connection"),
		Upgrade:    r.Header.Get("Upgrade"),

		// CDN 和安全信息
		CFRay:          r.Header.Get("CF-Ray"),
		CFConnectingIP: r.Header.Get("CF-Connecting-IP"),
		CFIPCountry:    r.Header.Get("CF-IPCountry"),
		XRequestID:     r.Header.Get("X-Request-ID"),
		XCorrelationID: r.Header.Get("X-Correlation-ID"),

		// 缓存和条件请求
		CacheControl:    r.Header.Get("Cache-Control"),
		IfNoneMatch:     r.Header.Get("If-None-Match"),
		IfModifiedSince: r.Header.Get("If-Modified-Since"),

		// 客户端提示信息
		SecCHUA:         r.Header.Get("Sec-CH-UA"),
		SecCHUAMobile:   r.Header.Get("Sec-CH-UA-Mobile"),
		SecCHUAPlatform: r.Header.Get("Sec-CH-UA-Platform"),
		DNT:             r.Header.Get("DNT"),
	}

	// TLS 信息
	if r.TLS != nil {
		metadata.Protocol = "https"
		metadata.TLSVersion = r.TLS.Version
		metadata.TLSCipherSuite = r.TLS.CipherSuite
		metadata.TLSServerName = r.TLS.ServerName
	} else {
		metadata.Protocol = "http"
	}

	// 解析 User-Agent
	if userAgent != "" {
		parsed := useragent.Parse(userAgent)
		metadata.Browser = parsed.Browser
		metadata.BrowserVersion = parsed.BrowserVersion
		metadata.OS = parsed.OS
		metadata.OSVersion = parsed.OSVersion
		metadata.Device = parsed.Device
		metadata.DeviceType = parsed.DeviceType
		metadata.DeviceVendor = parsed.DeviceVendor
		metadata.IsBot = parsed.IsBot
		metadata.BotName = parsed.BotName
		metadata.IsMobile = parsed.IsMobile
		metadata.IsTablet = parsed.IsTablet
	}

	return metadata
}
