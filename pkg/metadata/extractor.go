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
	"context"
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

type SourceFunc func(ctx context.Context, r *http.Request, key string) string
type ProcessFunc func(val string) string

type MetadataExtractor struct {
	ctx     context.Context
	r       *http.Request
	sources []struct {
		fn      SourceFunc
		key     string
		process ProcessFunc
	}
	defVal string
}

// NewMetadataExtractor 创建一个元数据提取器，从 context 和 request 中提取值
func NewMetadataExtractor(ctx context.Context, r *http.Request) *MetadataExtractor {
	return &MetadataExtractor{ctx: ctx, r: r}
}

// NewMetadataExtractorFromRequest 从 http.Request 创建元数据提取器，自动使用 request 的 context
func NewMetadataExtractorFromRequest(r *http.Request) *MetadataExtractor {
	var ctx context.Context
	if r != nil {
		ctx = r.Context()
	}
	return &MetadataExtractor{ctx: ctx, r: r}
}

// addSource 添加一个来源函数，用于从指定来源提取值
func (me *MetadataExtractor) addSource(fn SourceFunc, key string, process ProcessFunc) *MetadataExtractor {
	me.sources = append(me.sources, struct {
		fn      SourceFunc
		key     string
		process ProcessFunc
	}{fn, key, process})
	return me
}

// Default 设置默认值，当所有来源均未提取到值时返回该默认值
func (me *MetadataExtractor) Default(val string) *MetadataExtractor {
	me.defVal = val
	return me
}

// Get 执行提取操作，按添加的来源顺序尝试提取值
func (me *MetadataExtractor) Get() string {
	for _, src := range me.sources {
		v := src.fn(me.ctx, me.r, src.key)
		if v != "" && src.process != nil {
			v = src.process(v)
		}
		if v != "" {
			return v
		}
	}
	return me.defVal
}

// 便捷的链式调用方法

// FromContext 从 context 中提取值
func (me *MetadataExtractor) FromContext(key ContextKey) *MetadataExtractor {
	return me.addSource(FromContextSource, key.String(), nil)
}

// FromQuery 从 URL query 参数中提取值
func (me *MetadataExtractor) FromQuery(key string) *MetadataExtractor {
	return me.addSource(FromQuerySource, key, nil)
}

// FromHeader 从 HTTP header 中提取值
func (me *MetadataExtractor) FromHeader(key string) *MetadataExtractor {
	return me.addSource(FromHeaderSource, key, nil)
}

// FromCookie 从 HTTP cookie 中提取值
func (me *MetadataExtractor) FromCookie(key string) *MetadataExtractor {
	return me.addSource(FromCookieSource, key, nil)
}

// WithProcess 为最后添加的来源设置处理函数
func (me *MetadataExtractor) WithProcess(process ProcessFunc) *MetadataExtractor {
	if len(me.sources) > 0 {
		me.sources[len(me.sources)-1].process = process
	}
	return me
}

// 内置来源函数（独立函数形式，用于自定义场景）

type ContextKey string

func (ck ContextKey) String() string {
	return string(ck)
}

// FromContextSource 从 context 中提取字符串值
func FromContextSource(ctx context.Context, r *http.Request, key string) string {
	if ctx != nil {
		if v := ctx.Value(ContextKey(key)); v != nil {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// FromQuerySource 从 URL query 参数中提取值
func FromQuerySource(ctx context.Context, r *http.Request, key string) string {
	if r != nil {
		return r.URL.Query().Get(key)
	}
	return ""
}

// FromHeaderSource 从 HTTP header 中提取值
func FromHeaderSource(ctx context.Context, r *http.Request, key string) string {
	if r != nil {
		return r.Header.Get(key)
	}
	return ""
}

// FromCookieSource 从 HTTP cookie 中提取值
func FromCookieSource(ctx context.Context, r *http.Request, key string) string {
	if r != nil {
		if c, err := r.Cookie(key); err == nil {
			return c.Value
		}
	}
	return ""
}
