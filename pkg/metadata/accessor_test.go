/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 17:05:58
 * @FilePath: \go-toolbox\pkg\metadata\accessor_test.go
 * @Description: 元数据访问器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHeader(t *testing.T) {
	m := &RequestMetadata{
		UserAgent:              "Mozilla/5.0",
		ClientIP:               "192.168.1.1",
		RemoteAddr:             "192.168.1.1:12345",
		RequestURI:             "/api/test",
		QueryString:            "id=123",
		RequestMethod:          "GET",
		RequestHost:            "example.com",
		Origin:                 "https://example.com",
		Referer:                "https://example.com/page",
		XForwardedFor:          "203.0.113.1",
		XRealIP:                "203.0.113.1",
		XForwardedProto:        "https",
		XForwardedHost:         "example.com",
		XForwardedPort:         "443",
		AcceptLanguage:         "zh-CN",
		AcceptEncoding:         "gzip",
		Accept:                 "application/json",
		SecWebSocketKey:        "test-key",
		SecWebSocketVersion:    "13",
		SecWebSocketProtocol:   "chat",
		SecWebSocketExtensions: "permessage-deflate",
		Connection:             "Upgrade",
		Upgrade:                "websocket",
		CFRay:                  "ray-123",
		CFConnectingIP:         "203.0.113.100",
		CFIPCountry:            "US",
		XRequestID:             "req-123",
		XCorrelationID:         "corr-123",
		CacheControl:           "no-cache",
		IfNoneMatch:            "etag-123",
		IfModifiedSince:        "Mon, 01 Jan 2024 00:00:00 GMT",
		SecCHUA:                `"Chrome"; v="120"`,
		SecCHUAMobile:          "?0",
		SecCHUAPlatform:        `"macOS"`,
		DNT:                    "1",
		Protocol:               "https",
		TLSServerName:          "example.com",
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"UserAgent", "User-Agent", "Mozilla/5.0"},
		{"UserAgentUnderscore", "user_agent", "Mozilla/5.0"},
		{"ClientIP", "Client-IP", "192.168.1.1"},
		{"RemoteAddr", "Remote-Addr", "192.168.1.1:12345"},
		{"RequestURI", "Request-URI", "/api/test"},
		{"QueryString", "Query-String", "id=123"},
		{"RequestMethod", "Request-Method", "GET"},
		{"RequestHost", "Request-Host", "example.com"},
		{"Origin", "Origin", "https://example.com"},
		{"Referer", "Referer", "https://example.com/page"},
		{"XForwardedFor", "X-Forwarded-For", "203.0.113.1"},
		{"XRealIP", "X-Real-IP", "203.0.113.1"},
		{"XForwardedProto", "X-Forwarded-Proto", "https"},
		{"XForwardedHost", "X-Forwarded-Host", "example.com"},
		{"XForwardedPort", "X-Forwarded-Port", "443"},
		{"AcceptLanguage", "Accept-Language", "zh-CN"},
		{"AcceptEncoding", "Accept-Encoding", "gzip"},
		{"Accept", "Accept", "application/json"},
		{"SecWebSocketKey", "Sec-WebSocket-Key", "test-key"},
		{"SecWebSocketVersion", "Sec-WebSocket-Version", "13"},
		{"SecWebSocketProtocol", "Sec-WebSocket-Protocol", "chat"},
		{"SecWebSocketExtensions", "Sec-WebSocket-Extensions", "permessage-deflate"},
		{"Connection", "Connection", "Upgrade"},
		{"Upgrade", "Upgrade", "websocket"},
		{"CFRay", "CF-Ray", "ray-123"},
		{"CFConnectingIP", "CF-Connecting-IP", "203.0.113.100"},
		{"CFIPCountry", "CF-IPCountry", "US"},
		{"XRequestID", "X-Request-ID", "req-123"},
		{"XCorrelationID", "X-Correlation-ID", "corr-123"},
		{"CacheControl", "Cache-Control", "no-cache"},
		{"IfNoneMatch", "If-None-Match", "etag-123"},
		{"IfModifiedSince", "If-Modified-Since", "Mon, 01 Jan 2024 00:00:00 GMT"},
		{"SecCHUA", "Sec-CH-UA", `"Chrome"; v="120"`},
		{"SecCHUAMobile", "Sec-CH-UA-Mobile", "?0"},
		{"SecCHUAPlatform", "Sec-CH-UA-Platform", `"macOS"`},
		{"DNT", "DNT", "1"},
		{"Protocol", "Protocol", "https"},
		{"TLSServerName", "TLS-Server-Name", "example.com"},
		{"UnknownKey", "Unknown-Key", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.GetHeader(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetHeader(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		check func(*testing.T, *RequestMetadata)
	}{
		{
			name:  "UserAgent",
			key:   "User-Agent",
			value: "Mozilla/5.0",
			check: func(t *testing.T, m *RequestMetadata) {
				assert.Equal(t, "Mozilla/5.0", m.UserAgent)
			},
		},
		{
			name:  "ClientIP",
			key:   "Client-IP",
			value: "192.168.1.1",
			check: func(t *testing.T, m *RequestMetadata) {
				assert.Equal(t, "192.168.1.1", m.ClientIP)
			},
		},
		{
			name:  "RemoteAddr",
			key:   "Remote-Addr",
			value: "192.168.1.1:12345",
			check: func(t *testing.T, m *RequestMetadata) {
				assert.Equal(t, "192.168.1.1:12345", m.RemoteAddr)
			},
		},
		{
			name:  "RequestURI",
			key:   "Request-URI",
			value: "/api/test",
			check: func(t *testing.T, m *RequestMetadata) {
				assert.Equal(t, "/api/test", m.RequestURI)
			},
		},
		{
			name:  "Origin",
			key:   "Origin",
			value: "https://example.com",
			check: func(t *testing.T, m *RequestMetadata) {
				assert.Equal(t, "https://example.com", m.Origin)
			},
		},
		{
			name:  "XForwardedFor",
			key:   "X-Forwarded-For",
			value: "203.0.113.1",
			check: func(t *testing.T, m *RequestMetadata) {
				assert.Equal(t, "203.0.113.1", m.XForwardedFor)
			},
		},
		{
			name:  "CFRay",
			key:   "CF-Ray",
			value: "ray-123",
			check: func(t *testing.T, m *RequestMetadata) {
				assert.Equal(t, "ray-123", m.CFRay)
			},
		},
		{
			name:  "SecWebSocketKey",
			key:   "Sec-WebSocket-Key",
			value: "test-key",
			check: func(t *testing.T, m *RequestMetadata) {
				assert.Equal(t, "test-key", m.SecWebSocketKey)
			},
		},
		{
			name:  "UnknownKey",
			key:   "Unknown-Key",
			value: "value",
			check: func(t *testing.T, m *RequestMetadata) {
				// 未知的键不应该改变任何字段
				assert.Empty(t, m.UserAgent)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newM := &RequestMetadata{}
			newM.SetHeader(tt.key, tt.value)
			tt.check(t, newM)
		})
	}
}

func TestSetHeaderAllFields(t *testing.T) {
	// 测试所有字段都可以被设置
	m := &RequestMetadata{}

	headers := map[string]string{
		"User-Agent":               "Mozilla/5.0",
		"Client-IP":                "192.168.1.1",
		"Remote-Addr":              "192.168.1.1:12345",
		"Request-URI":              "/api/test",
		"Query-String":             "id=123",
		"Request-Method":           "GET",
		"Request-Host":             "example.com",
		"Origin":                   "https://example.com",
		"Referer":                  "https://example.com/page",
		"X-Forwarded-For":          "203.0.113.1",
		"X-Real-IP":                "203.0.113.1",
		"X-Forwarded-Proto":        "https",
		"X-Forwarded-Host":         "example.com",
		"X-Forwarded-Port":         "443",
		"Accept-Language":          "zh-CN",
		"Accept-Encoding":          "gzip",
		"Accept":                   "application/json",
		"Sec-WebSocket-Key":        "test-key",
		"Sec-WebSocket-Version":    "13",
		"Sec-WebSocket-Protocol":   "chat",
		"Sec-WebSocket-Extensions": "permessage-deflate",
		"Connection":               "Upgrade",
		"Upgrade":                  "websocket",
		"CF-Ray":                   "ray-123",
		"CF-Connecting-IP":         "203.0.113.100",
		"CF-IPCountry":             "US",
		"X-Request-ID":             "req-123",
		"X-Correlation-ID":         "corr-123",
		"Cache-Control":            "no-cache",
		"If-None-Match":            "etag-123",
		"If-Modified-Since":        "Mon, 01 Jan 2024 00:00:00 GMT",
		"Sec-CH-UA":                `"Chrome"; v="120"`,
		"Sec-CH-UA-Mobile":         "?0",
		"Sec-CH-UA-Platform":       `"macOS"`,
		"DNT":                      "1",
		"Protocol":                 "https",
		"TLS-Server-Name":          "example.com",
	}

	// 设置所有头
	for key, value := range headers {
		m.SetHeader(key, value)
	}

	// 验证所有头都被正确设置
	assert.Equal(t, "Mozilla/5.0", m.UserAgent)
	assert.Equal(t, "192.168.1.1", m.ClientIP)
	assert.Equal(t, "192.168.1.1:12345", m.RemoteAddr)
	assert.Equal(t, "/api/test", m.RequestURI)
	assert.Equal(t, "id=123", m.QueryString)
	assert.Equal(t, "GET", m.RequestMethod)
	assert.Equal(t, "example.com", m.RequestHost)
	assert.Equal(t, "https://example.com", m.Origin)
	assert.Equal(t, "https://example.com/page", m.Referer)
	assert.Equal(t, "203.0.113.1", m.XForwardedFor)
	assert.Equal(t, "203.0.113.1", m.XRealIP)
	assert.Equal(t, "https", m.XForwardedProto)
	assert.Equal(t, "example.com", m.XForwardedHost)
	assert.Equal(t, "443", m.XForwardedPort)
	assert.Equal(t, "zh-CN", m.AcceptLanguage)
	assert.Equal(t, "gzip", m.AcceptEncoding)
	assert.Equal(t, "application/json", m.Accept)
	assert.Equal(t, "test-key", m.SecWebSocketKey)
	assert.Equal(t, "13", m.SecWebSocketVersion)
	assert.Equal(t, "chat", m.SecWebSocketProtocol)
	assert.Equal(t, "permessage-deflate", m.SecWebSocketExtensions)
	assert.Equal(t, "Upgrade", m.Connection)
	assert.Equal(t, "websocket", m.Upgrade)
	assert.Equal(t, "ray-123", m.CFRay)
	assert.Equal(t, "203.0.113.100", m.CFConnectingIP)
	assert.Equal(t, "US", m.CFIPCountry)
	assert.Equal(t, "req-123", m.XRequestID)
	assert.Equal(t, "corr-123", m.XCorrelationID)
	assert.Equal(t, "no-cache", m.CacheControl)
	assert.Equal(t, "etag-123", m.IfNoneMatch)
	assert.Equal(t, "Mon, 01 Jan 2024 00:00:00 GMT", m.IfModifiedSince)
	assert.Equal(t, `"Chrome"; v="120"`, m.SecCHUA)
	assert.Equal(t, "?0", m.SecCHUAMobile)
	assert.Equal(t, `"macOS"`, m.SecCHUAPlatform)
	assert.Equal(t, "1", m.DNT)
	assert.Equal(t, "https", m.Protocol)
	assert.Equal(t, "example.com", m.TLSServerName)
}

func TestGetHeaderEmptyMetadata(t *testing.T) {
	m := &RequestMetadata{}

	// 空元数据应该返回空字符串
	assert.Empty(t, m.GetHeader("User-Agent"))
	assert.Empty(t, m.GetHeader("Origin"))
	assert.Empty(t, m.GetHeader("Unknown-Key"))
}
