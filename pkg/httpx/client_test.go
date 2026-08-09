/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-28 18:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-05-27 15:05:08
 * @FilePath: \go-toolbox\pkg\httpx\client_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package httpx

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 验证 http.Client 的 Transport 配置
func validateTransport(t *testing.T, client *http.Client) {
	assert.NotNil(t, client)
	transport, ok := client.Transport.(*http.Transport)
	assert.True(t, ok, "Expected client.Transport to be of type *http.Transport")

	// 验证 Transport 的配置
	assert.Equal(t, 0, transport.MaxIdleConns)
	assert.Equal(t, 1000, transport.MaxIdleConnsPerHost)
	assert.Equal(t, 1000, transport.MaxConnsPerHost)
	assert.Equal(t, 60*time.Second, transport.IdleConnTimeout)
	assert.Equal(t, 10*time.Second, transport.TLSHandshakeTimeout)
	assert.Equal(t, time.Second, transport.ExpectContinueTimeout)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.False(t, transport.TLSClientConfig.InsecureSkipVerify)
}

// TestNewHttpClient 测试使用自定义 HTTP 客户端创建
func TestNewHttpClient(t *testing.T) {
	client := NewHttpClient(http.DefaultClient)
	assert.NotNil(t, client)
	assert.Equal(t, http.DefaultClient, client.client)
}

// TestNewClientWithContext 测试使用自定义 HTTP 客户端和上下文创建
func TestNewClientWithContext(t *testing.T) {
	ctx := context.Background()
	client := NewClientWithContext(http.DefaultClient, ctx)
	assert.NotNil(t, client)
	assert.Equal(t, ctx, client.ctx)
	assert.Equal(t, http.DefaultClient, client.client)
}

// TestNewDefaultHttpClient 测试创建默认客户端
func TestNewDefaultHttpClient(t *testing.T) {
	client := NewDefaultHttpClient()
	assert.NotNil(t, client)
	assert.Equal(t, http.DefaultClient, client.client)
}

// TestNewDefaultHttpClientWithContext 测试创建默认客户端和自定义上下文
func TestNewDefaultHttpClientWithContext(t *testing.T) {
	ctx := context.Background()
	client := NewDefaultHttpClientWithContext(ctx)
	assert.NotNil(t, client)
	assert.Equal(t, ctx, client.ctx)
	assert.Equal(t, http.DefaultClient, client.client)
}

// TestNewCustomDefaultClient 测试创建自定义默认客户端
func TestNewCustomDefaultClient(t *testing.T) {
	client := NewCustomDefaultClient()
	assert.NotNil(t, client)
	validateTransport(t, client.client)
}

// TestNewCustomDefaultClientWithContext 测试创建自定义默认客户端和上下文
func TestNewCustomDefaultClientWithContext(t *testing.T) {
	ctx := context.Background()
	client := NewCustomDefaultClientWithContext(ctx)
	assert.NotNil(t, client)
	assert.Equal(t, ctx, client.ctx)
	validateTransport(t, client.client)
}

// TestClient_Request 测试客户端请求方法
func TestClient_Request(t *testing.T) {
	client := NewDefaultHttpClient()
	url := "http://localhost:8080"
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
		http.MethodConnect,
		http.MethodTrace,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			var r *Request
			switch method {
			case http.MethodGet:
				r = client.Get(url)
			case http.MethodPost:
				r = client.Post(url)
			case http.MethodPut:
				r = client.Put(url)
			case http.MethodDelete:
				r = client.Delete(url)
			case http.MethodPatch:
				r = client.Patch(url)
			case http.MethodHead:
				r = client.Head(url)
			case http.MethodOptions:
				r = client.Options(url)
			case http.MethodConnect:
				r = client.Connect(url)
			case http.MethodTrace:
				r = client.Trace(url)
			}

			compareRequest(t, r, method, client)
		})
	}
}

// compareRequest 比较请求的各个字段
func compareRequest(t *testing.T, r *Request, method string, client *Client) {
	assert.NotNil(t, r)
	assert.Equal(t, "http://localhost:8080", r.GetURL())
	assert.Equal(t, method, r.method)
	assert.Equal(t, client.client, r.client)
}

// TestDefaultClientConfig 测试默认配置函数返回值
func TestDefaultClientConfig(t *testing.T) {
	cfg := defaultClientConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, 30*time.Second, cfg.timeout)
	assert.Equal(t, 0, cfg.maxIdleConns)
	assert.Equal(t, 1000, cfg.maxIdleConnsPerHost)
	assert.Equal(t, 1000, cfg.maxConnsPerHost)
	assert.Equal(t, 60*time.Second, cfg.idleConnTimeout)
	assert.Equal(t, 10*time.Second, cfg.tlsHandshakeTimeout)
	assert.False(t, cfg.insecureSkipVerify)
	assert.NotNil(t, cfg.ctx)
}

// TestNewClient 测试使用函数式选项创建客户端
func TestNewClient(t *testing.T) {
	// 无选项时使用默认配置
	client := NewClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.client)
	validateTransport(t, client.client)
	// 默认上下文应为 context.Background()
	assert.NotNil(t, client.ctx)
}

// TestNewClientWithOptions 测试使用所有选项创建客户端
func TestNewClientWithOptions(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{ k string }{"k"}, "v")
	client := NewClient(
		WithTimeout(10*time.Second),
		WithMaxIdleConns(50),
		WithMaxIdleConnsPerHost(20),
		WithMaxConnsPerHost(30),
		WithIdleConnTimeout(15*time.Second),
		WithTLSHandshakeTimeout(5*time.Second),
		WithInsecureSkipVerify(true),
		WithContext(ctx),
	)
	assert.NotNil(t, client)
	assert.Equal(t, ctx, client.ctx)

	transport, ok := client.client.Transport.(*http.Transport)
	require.True(t, ok, "Transport 应为 *http.Transport 类型")
	assert.Equal(t, 10*time.Second, client.client.Timeout)
	assert.Equal(t, 50, transport.MaxIdleConns)
	assert.Equal(t, 20, transport.MaxIdleConnsPerHost)
	assert.Equal(t, 30, transport.MaxConnsPerHost)
	assert.Equal(t, 15*time.Second, transport.IdleConnTimeout)
	assert.Equal(t, 5*time.Second, transport.TLSHandshakeTimeout)
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
}

// TestClientOptionFunctions 测试每个选项函数是否能正确修改配置
func TestClientOptionFunctions(t *testing.T) {
	cfg := defaultClientConfig()

	WithTimeout(7 * time.Second)(cfg)
	assert.Equal(t, 7*time.Second, cfg.timeout)

	WithMaxIdleConns(11)(cfg)
	assert.Equal(t, 11, cfg.maxIdleConns)

	WithMaxIdleConnsPerHost(22)(cfg)
	assert.Equal(t, 22, cfg.maxIdleConnsPerHost)

	WithMaxConnsPerHost(33)(cfg)
	assert.Equal(t, 33, cfg.maxConnsPerHost)

	WithIdleConnTimeout(13 * time.Second)(cfg)
	assert.Equal(t, 13*time.Second, cfg.idleConnTimeout)

	WithTLSHandshakeTimeout(17 * time.Second)(cfg)
	assert.Equal(t, 17*time.Second, cfg.tlsHandshakeTimeout)

	WithInsecureSkipVerify(true)(cfg)
	assert.True(t, cfg.insecureSkipVerify)

	ctx := context.Background()
	WithContext(ctx)(cfg)
	assert.Equal(t, ctx, cfg.ctx)
}

// TestNewClientRequest 通过 Client.NewRequest 创建请求
func TestNewClientRequest(t *testing.T) {
	client := NewDefaultHttpClient()
	req := client.NewRequest(http.MethodGet, "http://localhost:9999")
	assert.NotNil(t, req)
	assert.Equal(t, http.MethodGet, req.Method())
	assert.Equal(t, "http://localhost:9999", req.URL())
}
