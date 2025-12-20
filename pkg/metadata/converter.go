/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 15:55:00
 * @FilePath: \go-toolbox\pkg\metadata\converter.go
 * @Description: 元数据转换工具
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package metadata

import (
	"reflect"
	"strings"
)

// ToMap 将 RequestMetadata 转换为 map[string]interface{}
func (m *RequestMetadata) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		// 基础请求信息
		"user_agent":     m.UserAgent,
		"client_ip":      m.ClientIP,
		"remote_addr":    m.RemoteAddr,
		"request_uri":    m.RequestURI,
		"query_string":   m.QueryString,
		"request_method": m.RequestMethod,
		"request_host":   m.RequestHost,

		// 来源信息
		"origin":  m.Origin,
		"referer": m.Referer,

		// 代理和转发信息
		"x_forwarded_for":   m.XForwardedFor,
		"x_real_ip":         m.XRealIP,
		"x_forwarded_proto": m.XForwardedProto,
		"x_forwarded_host":  m.XForwardedHost,
		"x_forwarded_port":  m.XForwardedPort,

		// 客户端偏好信息
		"accept_language": m.AcceptLanguage,
		"accept_encoding": m.AcceptEncoding,
		"accept":          m.Accept,

		// WebSocket 协议信息
		"sec_websocket_key":        m.SecWebSocketKey,
		"sec_websocket_version":    m.SecWebSocketVersion,
		"sec_websocket_protocol":   m.SecWebSocketProtocol,
		"sec_websocket_extensions": m.SecWebSocketExtensions,

		// 连接信息
		"connection": m.Connection,
		"upgrade":    m.Upgrade,

		// CDN 和安全信息
		"cf_ray":           m.CFRay,
		"cf_connecting_ip": m.CFConnectingIP,
		"cf_ipcountry":     m.CFIPCountry,
		"x_request_id":     m.XRequestID,
		"x_correlation_id": m.XCorrelationID,

		// 缓存和条件请求
		"cache_control":     m.CacheControl,
		"if_none_match":     m.IfNoneMatch,
		"if_modified_since": m.IfModifiedSince,

		// 客户端提示信息
		"sec_ch_ua":          m.SecCHUA,
		"sec_ch_ua_mobile":   m.SecCHUAMobile,
		"sec_ch_ua_platform": m.SecCHUAPlatform,
		"dnt":                m.DNT,

		// 协议信息
		"protocol": m.Protocol,

		// User-Agent 解析结果
		"browser":         m.Browser,
		"browser_version": m.BrowserVersion,
		"os":              m.OS,
		"os_version":      m.OSVersion,
		"device":          m.Device,
		"device_type":     m.DeviceType,
		"device_vendor":   m.DeviceVendor,
		"is_bot":          m.IsBot,
		"bot_name":        m.BotName,
		"is_mobile":       m.IsMobile,
		"is_tablet":       m.IsTablet,
	}

	// 只在 TLS 存在时添加 TLS 相关字段
	if m.TLSVersion > 0 {
		result["tls_version"] = m.TLSVersion
		result["tls_cipher_suite"] = m.TLSCipherSuite
		result["tls_server_name"] = m.TLSServerName
	}

	return result
}

// FromMap 从 map[string]interface{} 创建 RequestMetadata（使用反射自动填充）
func FromMap(data map[string]interface{}) *RequestMetadata {
	metadata := &RequestMetadata{}

	val := reflect.ValueOf(metadata).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 从 json tag 获取 key
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		// 解析 json tag（去除 omitempty 等选项）
		key := strings.Split(jsonTag, ",")[0]
		if key == "" || key == "-" {
			continue
		}

		// 从 map 中获取值
		mapValue, ok := data[key]
		if !ok {
			continue
		}

		// 根据字段类型设置值
		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if str, ok := mapValue.(string); ok {
				field.SetString(str)
			}
		case reflect.Uint16:
			switch v := mapValue.(type) {
			case uint16:
				field.SetUint(uint64(v))
			case float64:
				field.SetUint(uint64(v))
			case int:
				field.SetUint(uint64(v))
			case int64:
				field.SetUint(uint64(v))
			}
		case reflect.Bool:
			if b, ok := mapValue.(bool); ok {
				field.SetBool(b)
			}
		}
	}

	return metadata
}
