/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-28 18:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-30 15:26:55
 * @FilePath: \go-toolbox\pkg\httpx\client.go
 * @Description: HTTP 客户端实现
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Client 是一个封装 http.Client 的结构体
type Client struct {
	ctx    context.Context
	client *http.Client
}

// NewHttpClient 创建一个使用自定义 http.Client 的 Client 实例
func NewHttpClient(client *http.Client) *Client {
	return newClient(context.Background(), client)
}

// NewClientWithContext 创建一个使用自定义 http.Client 的 Client 实例
func NewClientWithContext(client *http.Client, ctx context.Context) *Client {
	return newClient(ctx, client)
}

// NewDefaultHttpClient 创建一个使用默认 http.Client 的 Client 实例
func NewDefaultHttpClient() *Client {
	return NewDefaultHttpClientWithContext(context.Background())
}

// NewDefaultHttpClientWithContext 创建一个使用默认 http.Client 和自定义上下文的 Client 实例
func NewDefaultHttpClientWithContext(ctx context.Context) *Client {
	return newClient(ctx, http.DefaultClient)
}

// NewCustomDefaultClient 创建一个使用自定义的 http.Client 的 Client 实例
func NewCustomDefaultClient() *Client {
	return NewCustomDefaultClientWithContext(context.Background())
}

// NewCustomDefaultClientWithContext 创建一个使用自定义的 http.Client 和自定义上下文的 Client 实例
func NewCustomDefaultClientWithContext(ctx context.Context) *Client {
	customClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   1000,
			MaxConnsPerHost:       1000,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: time.Second,
		},
	}
	return newClient(ctx, customClient)
}

// newClient 是一个辅助函数，用于初始化 Client
func newClient(ctx context.Context, client *http.Client) *Client {
	return &Client{
		ctx:    ctx,
		client: client,
	}
}

// NewRequest 创建一个新的 Request 实例
func (c *Client) NewRequest(method, endpoint string) *Request {
	return &Request{
		ctx:         c.ctx,
		client:      c.client,
		method:      method,
		endpoint:    endpoint,
		queryValues: make(url.Values),
		headers:     make(http.Header),
	}
}

// Get 创建一个 GET 请求
func (c *Client) Get(endpoint string) *Request {
	return c.NewRequest(http.MethodGet, endpoint)
}

// Post 创建一个 POST 请求
func (c *Client) Post(endpoint string) *Request {
	return c.NewRequest(http.MethodPost, endpoint)
}

// Put 创建一个 PUT 请求
func (c *Client) Put(endpoint string) *Request {
	return c.NewRequest(http.MethodPut, endpoint)
}

// Delete 创建一个 DELETE 请求
func (c *Client) Delete(endpoint string) *Request {
	return c.NewRequest(http.MethodDelete, endpoint)
}

// Patch 创建一个 PATCH 请求
func (c *Client) Patch(endpoint string) *Request {
	return c.NewRequest(http.MethodPatch, endpoint)
}

// Head 创建一个 HEAD 请求
func (c *Client) Head(endpoint string) *Request {
	return c.NewRequest(http.MethodHead, endpoint)
}

// Options 创建一个 OPTIONS 请求
func (c *Client) Options(endpoint string) *Request {
	return c.NewRequest(http.MethodOptions, endpoint)
}

// Connect 创建一个 CONNECT 请求
func (c *Client) Connect(endpoint string) *Request {
	return c.NewRequest(http.MethodConnect, endpoint)
}

// Trace 创建一个 TRACE 请求
func (c *Client) Trace(endpoint string) *Request {
	return c.NewRequest(http.MethodTrace, endpoint)
}
