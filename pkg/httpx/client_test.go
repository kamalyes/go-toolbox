/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-28 18:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-29 18:07:36
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
)

// mapsEqual 检查两个映射是否相等，支持忽略长度比较
func mapsEqual(t *testing.T, m1, m2 map[string][]string, ignoreLen bool) {
	if !ignoreLen {
		assert.Equal(t, m1, m2)
	}
	for key, aValue := range m1 {
		bValue := m2[key]
		assert.ElementsMatch(t, aValue, bValue)
	}
}

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
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)
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
