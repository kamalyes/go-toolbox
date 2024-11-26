/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-29 19:07:21
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-30 20:23:00
 * @FilePath: \go-toolbox\pkg\httpx\request_test.go
 * @Description: HTTP 请求封装测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义常量
const (
	testURL = "http://example.com" // 定义测试用的 URL
)

func setupRequest(method, url string) *Request {
	client := &http.Client{}
	return NewRequest(context.Background(), client, method, url)
}

func TestNewRequest(t *testing.T) {
	tests := map[string]struct {
		method   string
		expected string
	}{
		"GET Request":  {"GET", testURL},
		"POST Request": {"POST", testURL},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := setupRequest(tt.method, tt.expected)

			assert.NotNil(t, req)
			assert.Equal(t, tt.method, req.method)
			assert.Equal(t, tt.expected, req.endpoint)
		})
	}
}

func TestRequestSetHeader(t *testing.T) {
	req := setupRequest("GET", testURL).
		SetHeader("X-Custom-Header", "value")

	assert.Equal(t, "value", req.GetHeaders().Get("X-Custom-Header"))
}

func TestRequestSend(t *testing.T) {
	// 创建一个测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, ContentTypeApplicationJSON, r.Header.Get("Content-Type"))
		w.Write([]byte("Success"))
	}))
	defer server.Close()

	client := &http.Client{}
	req := NewRequest(context.Background(), client, "POST", server.URL).
		SetHeader("Content-Type", ContentTypeApplicationJSON).
		SetBody(map[string]string{"name": "test"})

	resp, err := req.Send()
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Success", string(body))
}

func TestRequestSetBodyForm(t *testing.T) {
	req := setupRequest("POST", testURL).
		SetBodyForm(url.Values{"key": {"value"}})

	assert.NotNil(t, req.GetBody())
}

func TestRequestSetBodyMultipart(t *testing.T) {
	// 创建一个测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, ContentTypeMultipartFormData, r.Header.Get("Content-Type"))
		w.Write([]byte("File Uploaded"))
	}))
	defer server.Close()

	client := &http.Client{}
	req := NewRequest(context.Background(), client, "POST", server.URL).
		SetBodyMultipart("file", "test.txt", []byte("File content"))

	resp, err := req.Send()
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "File Uploaded", string(body))
}
