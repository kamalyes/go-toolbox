/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-20 16:58:01
 * @FilePath: \go-toolbox\pkg\metadata\extractor_test.go
 * @Description: HTTP 请求元数据提取器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package metadata

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractRequestMetadataChrome(t *testing.T) {
	// 模拟 Chrome 浏览器请求
	req := httptest.NewRequest("GET", "/api/test?id=123", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("X-Forwarded-For", "203.0.113.1")

	metadata := ExtractRequestMetadata(req)

	// 验证基础信息
	assert.Equal(t, "GET", metadata.RequestMethod)
	assert.Equal(t, "/api/test?id=123", metadata.RequestURI)
	assert.Equal(t, "id=123", metadata.QueryString)
	assert.Equal(t, "https://example.com", metadata.Origin)
	assert.Equal(t, "zh-CN,zh;q=0.9", metadata.AcceptLanguage)
	assert.Equal(t, "203.0.113.1", metadata.XForwardedFor)

	// 验证 User-Agent 解析结果
	assert.Equal(t, "Chrome", metadata.Browser)
	assert.Equal(t, "120", metadata.BrowserVersion)
	assert.Equal(t, "Windows", metadata.OS)
	assert.Equal(t, "10", metadata.OSVersion)
	assert.Equal(t, "desktop", metadata.DeviceType)
	assert.False(t, metadata.IsMobile)
	assert.False(t, metadata.IsTablet)
	assert.False(t, metadata.IsBot)
}

func TestExtractRequestMetadataiPhone(t *testing.T) {
	// 模拟 iPhone Safari 请求
	req := httptest.NewRequest("POST", "/api/users", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")
	req.Header.Set("X-Real-IP", "198.51.100.50")

	metadata := ExtractRequestMetadata(req)

	// 验证移动设备信息
	assert.Equal(t, "Safari", metadata.Browser)
	assert.Equal(t, "16", metadata.BrowserVersion)
	assert.Equal(t, "iOS", metadata.OS)
	assert.Equal(t, "16.6", metadata.OSVersion)
	assert.Equal(t, "mobile", metadata.DeviceType)
	assert.Equal(t, "Apple", metadata.DeviceVendor)
	assert.Equal(t, "iphone", metadata.Device)
	assert.True(t, metadata.IsMobile)
	assert.False(t, metadata.IsTablet)
}

func TestExtractRequestMetadataAndroid(t *testing.T) {
	// 模拟 Android Chrome 请求
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 13; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Mobile Safari/537.36")

	metadata := ExtractRequestMetadata(req)

	// 验证 Android 设备
	assert.Equal(t, "Chrome", metadata.Browser)
	assert.Equal(t, "119", metadata.BrowserVersion)
	assert.Equal(t, "Android", metadata.OS)
	assert.Equal(t, "13", metadata.OSVersion)
	assert.Equal(t, "mobile", metadata.DeviceType)
	assert.Equal(t, "Samsung", metadata.DeviceVendor)
	assert.True(t, metadata.IsMobile)
	assert.False(t, metadata.IsTablet)
}

func TestExtractRequestMetadataIPad(t *testing.T) {
	// 模拟 iPad Safari 请求
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPad; CPU OS 15_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Mobile/15E148 Safari/604.1")

	metadata := ExtractRequestMetadata(req)

	// 验证平板设备
	assert.Equal(t, "Safari", metadata.Browser)
	assert.Equal(t, "iOS", metadata.OS)
	assert.Equal(t, "15.7", metadata.OSVersion)
	assert.Equal(t, "tablet", metadata.DeviceType)
	assert.Equal(t, "Apple", metadata.DeviceVendor)
	assert.Equal(t, "ipad", metadata.Device)
	assert.False(t, metadata.IsMobile)
	assert.True(t, metadata.IsTablet)
}

func TestExtractRequestMetadataBot(t *testing.T) {
	// 模拟 Google Bot 请求
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	metadata := ExtractRequestMetadata(req)

	// 验证爬虫检测
	assert.True(t, metadata.IsBot)
	assert.Equal(t, "Googlebot", metadata.BotName)
	assert.Equal(t, "bot", metadata.DeviceType)
	assert.False(t, metadata.IsMobile)
	assert.False(t, metadata.IsTablet)
}

func TestExtractRequestMetadataWebSocket(t *testing.T) {
	// 模拟 WebSocket 升级请求
	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Protocol", "chat")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")

	metadata := ExtractRequestMetadata(req)

	// 验证 WebSocket 信息
	assert.Equal(t, "Upgrade", metadata.Connection)
	assert.Equal(t, "websocket", metadata.Upgrade)
	assert.Equal(t, "dGhlIHNhbXBsZSBub25jZQ==", metadata.SecWebSocketKey)
	assert.Equal(t, "13", metadata.SecWebSocketVersion)
	assert.Equal(t, "chat", metadata.SecWebSocketProtocol)

	// 验证 User-Agent 解析
	assert.Equal(t, "Chrome", metadata.Browser)
	assert.Equal(t, "Windows", metadata.OS)
}

func TestExtractRequestMetadataTLS(t *testing.T) {
	// 创建带 TLS 的请求
	req := httptest.NewRequest("GET", "https://example.com/", nil)
	req.TLS = &tls.ConnectionState{
		Version:     tls.VersionTLS13,
		CipherSuite: tls.TLS_AES_128_GCM_SHA256,
		ServerName:  "example.com",
	}

	metadata := ExtractRequestMetadata(req)

	// 验证 TLS 信息
	assert.Equal(t, "https", metadata.Protocol)
	assert.Equal(t, uint16(tls.VersionTLS13), metadata.TLSVersion)
	assert.Equal(t, uint16(tls.TLS_AES_128_GCM_SHA256), metadata.TLSCipherSuite)
	assert.Equal(t, "example.com", metadata.TLSServerName)
}

func TestExtractRequestMetadataCDN(t *testing.T) {
	// 模拟通过 Cloudflare CDN 的请求
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("CF-Ray", "7f1234567890-LAX")
	req.Header.Set("CF-Connecting-IP", "203.0.113.100")
	req.Header.Set("CF-IPCountry", "US")
	req.Header.Set("X-Request-ID", "req-abc123")

	metadata := ExtractRequestMetadata(req)

	// 验证 CDN 信息
	assert.Equal(t, "7f1234567890-LAX", metadata.CFRay)
	assert.Equal(t, "203.0.113.100", metadata.CFConnectingIP)
	assert.Equal(t, "US", metadata.CFIPCountry)
	assert.Equal(t, "req-abc123", metadata.XRequestID)
}

func TestExtractRequestMetadataEmptyUserAgent(t *testing.T) {
	// 没有 User-Agent 的请求
	req := httptest.NewRequest("GET", "/", nil)

	metadata := ExtractRequestMetadata(req)

	// User-Agent 相关字段应该为空
	assert.Empty(t, metadata.UserAgent)
	assert.Empty(t, metadata.Browser)
	assert.Empty(t, metadata.BrowserVersion)
	assert.Empty(t, metadata.OS)
	assert.Empty(t, metadata.OSVersion)
	assert.Empty(t, metadata.DeviceType)
	assert.False(t, metadata.IsBot)
	assert.False(t, metadata.IsMobile)
	assert.False(t, metadata.IsTablet)
}

func TestExtractRequestMetadataCompleteRequest(t *testing.T) {
	// 模拟完整的请求（包含所有可能的头信息）
	req := httptest.NewRequest("POST", "https://api.example.com/v1/resource?filter=active", nil)
	req.Host = "api.example.com"
	req.RemoteAddr = "203.0.113.100:12345"

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Origin", "https://app.example.com")
	req.Header.Set("Referer", "https://app.example.com/dashboard")
	req.Header.Set("X-Forwarded-For", "198.51.100.1, 198.51.100.2")
	req.Header.Set("X-Real-IP", "198.51.100.1")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "api.example.com")
	req.Header.Set("X-Forwarded-Port", "443")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("Sec-CH-UA", `"Chrome"; v="120"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)

	req.TLS = &tls.ConnectionState{
		Version:    tls.VersionTLS12,
		ServerName: "api.example.com",
	}

	metadata := ExtractRequestMetadata(req)

	// 验证所有字段
	assert.Equal(t, "POST", metadata.RequestMethod)
	assert.Equal(t, "/v1/resource?filter=active", metadata.RequestURI)
	assert.Equal(t, "filter=active", metadata.QueryString)
	assert.Equal(t, "api.example.com", metadata.RequestHost)
	assert.Equal(t, "203.0.113.100:12345", metadata.RemoteAddr)

	assert.Equal(t, "https://app.example.com", metadata.Origin)
	assert.Equal(t, "https://app.example.com/dashboard", metadata.Referer)

	assert.Equal(t, "198.51.100.1, 198.51.100.2", metadata.XForwardedFor)
	assert.Equal(t, "198.51.100.1", metadata.XRealIP)
	assert.Equal(t, "https", metadata.XForwardedProto)
	assert.Equal(t, "api.example.com", metadata.XForwardedHost)
	assert.Equal(t, "443", metadata.XForwardedPort)

	assert.Equal(t, "en-US,en;q=0.9", metadata.AcceptLanguage)
	assert.Equal(t, "gzip, deflate, br", metadata.AcceptEncoding)
	assert.Equal(t, "application/json", metadata.Accept)

	assert.Equal(t, "no-cache", metadata.CacheControl)
	assert.Equal(t, "1", metadata.DNT)
	assert.Equal(t, `"Chrome"; v="120"`, metadata.SecCHUA)
	assert.Equal(t, "?0", metadata.SecCHUAMobile)
	assert.Equal(t, `"macOS"`, metadata.SecCHUAPlatform)

	assert.Equal(t, "https", metadata.Protocol)
	assert.Equal(t, uint16(tls.VersionTLS12), metadata.TLSVersion)
	assert.Equal(t, "api.example.com", metadata.TLSServerName)

	// User-Agent 解析
	assert.Equal(t, "Chrome", metadata.Browser)
	assert.Equal(t, "120", metadata.BrowserVersion)
	assert.Equal(t, "macOS", metadata.OS)
	assert.Equal(t, "10.15", metadata.OSVersion)
	assert.Equal(t, "desktop", metadata.DeviceType)
	assert.Equal(t, "Apple", metadata.DeviceVendor)
	assert.False(t, metadata.IsMobile)
	assert.False(t, metadata.IsTablet)
	assert.False(t, metadata.IsBot)
}

// 定义常量作为上下文键和其他测试用的字符串

const (
	TestQueryKeyFoo                 = "foo"
	TestHeaderKeyCustom             = "X-Custom-Header"
	TestCookieKeySession            = "session"
	TestContextKey       ContextKey = "contextKey"
	TestContextValue                = "contextValue"
	TestDefaultValue                = "defaultValue"
)

func TestMetadataExtractor(t *testing.T) {
	// 创建一个上下文
	ctx := context.Background()

	// 创建一个 HTTP 请求
	req := httptest.NewRequest("GET", "http://example.com?"+TestQueryKeyFoo+"=bar", nil)
	req.Header.Set(TestHeaderKeyCustom, "headerValue")
	req.AddCookie(&http.Cookie{Name: TestCookieKeySession, Value: "cookieValue"})

	// 创建 MetadataExtractor 实例
	extractor := NewMetadataExtractor(ctx, req)

	// 测试从查询参数中提取值
	extractor.FromQuery(TestQueryKeyFoo)
	assert.Equal(t, "bar", extractor.Get(), "应该从查询参数中提取到值")

	// 测试从 HTTP header 中提取值
	extractor = NewMetadataExtractor(ctx, req) // 重新创建实例
	extractor.FromHeader(TestHeaderKeyCustom)
	assert.Equal(t, "headerValue", extractor.Get(), "应该从 HTTP header 中提取到值")

	// 测试从 HTTP cookie 中提取值
	extractor = NewMetadataExtractor(ctx, req) // 重新创建实例
	extractor.FromCookie(TestCookieKeySession)
	assert.Equal(t, "cookieValue", extractor.Get(), "应该从 HTTP cookie 中提取到值")

	// 测试从 context 中提取值
	ctx = context.WithValue(ctx, TestContextKey, TestContextValue)
	extractor = NewMetadataExtractor(ctx, req) // 重新创建实例
	extractor.FromContext(TestContextKey)
	assert.Equal(t, TestContextValue, extractor.Get(), "应该从 context 中提取到值")

	// 测试默认值
	extractor = NewMetadataExtractor(ctx, req) // 重新创建实例
	extractor.Default(TestDefaultValue)
	assert.Equal(t, TestDefaultValue, extractor.Get(), "所有来源均未提取到值时应该返回默认值")
}
