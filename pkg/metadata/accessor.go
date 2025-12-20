/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 15:55:00
 * @FilePath: \go-toolbox\pkg\metadata\accessor.go
 * @Description: 元数据访问器
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package metadata

import "strings"

// GetHeader 从 RequestMetadata 中获取指定的头信息
func (m *RequestMetadata) GetHeader(key string) string {
	key = strings.ToLower(strings.ReplaceAll(key, "-", "_"))

	switch key {
	case "user_agent":
		return m.UserAgent
	case "client_ip":
		return m.ClientIP
	case "remote_addr":
		return m.RemoteAddr
	case "request_uri":
		return m.RequestURI
	case "query_string":
		return m.QueryString
	case "request_method":
		return m.RequestMethod
	case "request_host":
		return m.RequestHost
	case "origin":
		return m.Origin
	case "referer":
		return m.Referer
	case "x_forwarded_for":
		return m.XForwardedFor
	case "x_real_ip":
		return m.XRealIP
	case "x_forwarded_proto":
		return m.XForwardedProto
	case "x_forwarded_host":
		return m.XForwardedHost
	case "x_forwarded_port":
		return m.XForwardedPort
	case "accept_language":
		return m.AcceptLanguage
	case "accept_encoding":
		return m.AcceptEncoding
	case "accept":
		return m.Accept
	case "sec_websocket_key":
		return m.SecWebSocketKey
	case "sec_websocket_version":
		return m.SecWebSocketVersion
	case "sec_websocket_protocol":
		return m.SecWebSocketProtocol
	case "sec_websocket_extensions":
		return m.SecWebSocketExtensions
	case "connection":
		return m.Connection
	case "upgrade":
		return m.Upgrade
	case "cf_ray":
		return m.CFRay
	case "cf_connecting_ip":
		return m.CFConnectingIP
	case "cf_ipcountry":
		return m.CFIPCountry
	case "x_request_id":
		return m.XRequestID
	case "x_correlation_id":
		return m.XCorrelationID
	case "cache_control":
		return m.CacheControl
	case "if_none_match":
		return m.IfNoneMatch
	case "if_modified_since":
		return m.IfModifiedSince
	case "sec_ch_ua":
		return m.SecCHUA
	case "sec_ch_ua_mobile":
		return m.SecCHUAMobile
	case "sec_ch_ua_platform":
		return m.SecCHUAPlatform
	case "dnt":
		return m.DNT
	case "protocol":
		return m.Protocol
	case "tls_server_name":
		return m.TLSServerName
	default:
		return ""
	}
}

// SetHeader 设置 RequestMetadata 中指定的头信息
func (m *RequestMetadata) SetHeader(key, value string) {
	key = strings.ToLower(strings.ReplaceAll(key, "-", "_"))

	switch key {
	case "user_agent":
		m.UserAgent = value
	case "client_ip":
		m.ClientIP = value
	case "remote_addr":
		m.RemoteAddr = value
	case "request_uri":
		m.RequestURI = value
	case "query_string":
		m.QueryString = value
	case "request_method":
		m.RequestMethod = value
	case "request_host":
		m.RequestHost = value
	case "origin":
		m.Origin = value
	case "referer":
		m.Referer = value
	case "x_forwarded_for":
		m.XForwardedFor = value
	case "x_real_ip":
		m.XRealIP = value
	case "x_forwarded_proto":
		m.XForwardedProto = value
	case "x_forwarded_host":
		m.XForwardedHost = value
	case "x_forwarded_port":
		m.XForwardedPort = value
	case "accept_language":
		m.AcceptLanguage = value
	case "accept_encoding":
		m.AcceptEncoding = value
	case "accept":
		m.Accept = value
	case "sec_websocket_key":
		m.SecWebSocketKey = value
	case "sec_websocket_version":
		m.SecWebSocketVersion = value
	case "sec_websocket_protocol":
		m.SecWebSocketProtocol = value
	case "sec_websocket_extensions":
		m.SecWebSocketExtensions = value
	case "connection":
		m.Connection = value
	case "upgrade":
		m.Upgrade = value
	case "cf_ray":
		m.CFRay = value
	case "cf_connecting_ip":
		m.CFConnectingIP = value
	case "cf_ipcountry":
		m.CFIPCountry = value
	case "x_request_id":
		m.XRequestID = value
	case "x_correlation_id":
		m.XCorrelationID = value
	case "cache_control":
		m.CacheControl = value
	case "if_none_match":
		m.IfNoneMatch = value
	case "if_modified_since":
		m.IfModifiedSince = value
	case "sec_ch_ua":
		m.SecCHUA = value
	case "sec_ch_ua_mobile":
		m.SecCHUAMobile = value
	case "sec_ch_ua_platform":
		m.SecCHUAPlatform = value
	case "dnt":
		m.DNT = value
	case "protocol":
		m.Protocol = value
	case "tls_server_name":
		m.TLSServerName = value
	}
}
