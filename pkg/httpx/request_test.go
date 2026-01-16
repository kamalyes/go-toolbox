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

	assert.Equal(t, "value", req.Header().Get("X-Custom-Header"))
}

func TestRequestSend(t *testing.T) {
	// 创建一个测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, ContentTypeApplicationJSON, r.Header.Get(HeaderContentType))
		w.Write([]byte("Success"))
	}))
	defer server.Close()

	client := &http.Client{}
	req := NewRequest(context.Background(), client, "POST", server.URL).
		SetHeader(HeaderContentType, ContentTypeApplicationJSON).
		SetBody(map[string]string{"name": "test"})

	resp, err := req.Send()
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Response.Body)
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
		// multipart/form-data 会包含 boundary，所以使用 Contains 检查
		contentType := r.Header.Get(HeaderContentType)
		assert.Contains(t, contentType, "multipart/form-data")
		w.Write([]byte("File Uploaded"))
	}))
	defer server.Close()

	client := &http.Client{}
	req := NewRequest(context.Background(), client, "POST", server.URL).
		SetBodyMultipart("file", "test.txt", []byte("File content"))

	resp, err := req.Send()
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Response.Body)
	assert.NoError(t, err)
	assert.Equal(t, "File Uploaded", string(body))
}

// TestRequestSetQueries 测试批量设置查询参数
func TestRequestSetQueries(t *testing.T) {
	req := setupRequest("GET", testURL).
		SetQueries(map[string]string{
			"page": "1",
			"size": "10",
			"sort": "desc",
		})

	assert.Equal(t, "1", req.queryValues.Get("page"))
	assert.Equal(t, "10", req.queryValues.Get("size"))
	assert.Equal(t, "desc", req.queryValues.Get("sort"))
}

// TestRequestAddQueries 测试批量添加查询参数
func TestRequestAddQueries(t *testing.T) {
	req := setupRequest("GET", testURL).
		AddQuery("key", "value1").
		AddQueries(map[string]string{
			"key": "value2",
			"foo": "bar",
		})

	// AddQueries 会添加，不会覆盖
	values := req.queryValues["key"]
	assert.Len(t, values, 2)
	assert.Contains(t, values, "value1")
	assert.Contains(t, values, "value2")
	assert.Equal(t, "bar", req.queryValues.Get("foo"))
}

// TestRequestSetHeaders 测试批量设置请求头
func TestRequestSetHeaders(t *testing.T) {
	req := setupRequest("GET", testURL).
		SetHeaders(map[string]string{
			"X-Custom-1": "value1",
			"X-Custom-2": "value2",
		})

	assert.Equal(t, "value1", req.headers.Get("X-Custom-1"))
	assert.Equal(t, "value2", req.headers.Get("X-Custom-2"))
}

// TestRequestAddHeaders 测试批量添加请求头
func TestRequestAddHeaders(t *testing.T) {
	req := setupRequest("GET", testURL).
		AddHeader("X-Test", "first").
		AddHeaders(map[string]string{
			"X-Test":  "second",
			"X-Other": "value",
		})

	values := req.headers["X-Test"]
	assert.Len(t, values, 2)
	assert.Equal(t, "value", req.headers.Get("X-Other"))
}

// TestRequestSetUserAgent 测试设置 User-Agent
func TestRequestSetUserAgent(t *testing.T) {
	req := setupRequest("GET", testURL).
		SetUserAgent("MyApp/1.0")

	assert.Equal(t, "MyApp/1.0", req.headers.Get("User-Agent"))
}

// TestRequestSetAuthorization 测试设置 Authorization
func TestRequestSetAuthorization(t *testing.T) {
	req := setupRequest("GET", testURL).
		SetAuthorization("token123")

	assert.Equal(t, "token123", req.headers.Get("Authorization"))
}

// TestRequestSetBearerToken 测试设置 Bearer Token
func TestRequestSetBearerToken(t *testing.T) {
	req := setupRequest("GET", testURL).
		SetBearerToken("mytoken")

	assert.Equal(t, "mytoken", req.headers.Get("Authorization"))
}

// TestRequestSetContentType 测试设置 Content-Type
func TestRequestSetContentType(t *testing.T) {
	req := setupRequest("POST", testURL).
		SetContentType("application/xml")

	assert.Equal(t, "application/xml", req.headers.Get("Content-Type"))
}

// TestRequestSetAccept 测试设置 Accept
func TestRequestSetAccept(t *testing.T) {
	req := setupRequest("GET", testURL).
		SetAccept("application/json")

	assert.Equal(t, "application/json", req.headers.Get("Accept"))
}

// TestRequestSetBodyJSON 测试设置 JSON 请求体
func TestRequestSetBodyJSON(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"age":  30,
	}

	req := setupRequest("POST", testURL).
		SetBodyJSON(data)

	assert.Equal(t, ContentTypeApplicationJSON, req.headers.Get("Content-Type"))
	assert.NotNil(t, req.body)
	assert.NotNil(t, req.bodyEncodeFunc)
}

// TestRequestSetBodyString 测试设置字符串请求体
func TestRequestSetBodyString(t *testing.T) {
	req := setupRequest("POST", testURL).
		SetBodyString("raw string body")

	assert.NotNil(t, req.bodyBytes)
}

// TestRequestSetQueryValues 测试直接设置 url.Values
func TestRequestSetQueryValues(t *testing.T) {
	values := url.Values{
		"key1": []string{"value1"},
		"key2": []string{"value2a", "value2b"},
	}

	req := setupRequest("GET", testURL).
		SetQueryValues(values)

	assert.Equal(t, "value1", req.queryValues.Get("key1"))
	assert.Len(t, req.queryValues["key2"], 2)
}

// TestRequestGetFullURL 测试获取完整 URL
func TestRequestGetFullURL(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		queries  map[string]string
		expected string
	}{
		{
			name:     "无查询参数",
			endpoint: "http://example.com/api",
			queries:  nil,
			expected: "http://example.com/api",
		},
		{
			name:     "有查询参数",
			endpoint: "http://example.com/api",
			queries:  map[string]string{"page": "1", "size": "10"},
			expected: "http://example.com/api?page=1&size=10",
		},
		{
			name:     "已有查询参数",
			endpoint: "http://example.com/api?existing=true",
			queries:  map[string]string{"page": "1"},
			expected: "http://example.com/api?existing=true&page=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := setupRequest("GET", tt.endpoint)
			if tt.queries != nil {
				req.SetQueries(tt.queries)
			}
			fullURL := req.GetFullURL()
			assert.Contains(t, fullURL, tt.endpoint)
		})
	}
}

// TestRequestClone 测试克隆请求
func TestRequestClone(t *testing.T) {
	original := setupRequest("POST", testURL).
		SetHeader("X-Custom", "value").
		SetQuery("page", "1").
		SetBody(map[string]string{"name": "test"})

	cloned := original.Clone()

	// 验证克隆的基本属性
	assert.Equal(t, original.method, cloned.method)
	assert.Equal(t, original.endpoint, cloned.endpoint)
	assert.Equal(t, original.headers.Get("X-Custom"), cloned.headers.Get("X-Custom"))
	assert.Equal(t, original.queryValues.Get("page"), cloned.queryValues.Get("page"))

	// 修改克隆不应影响原始请求
	cloned.SetHeader("X-Custom", "new-value")
	assert.Equal(t, "value", original.headers.Get("X-Custom"))
	assert.Equal(t, "new-value", cloned.headers.Get("X-Custom"))
}

// TestRequestSetBodyMultipartWithFields 测试设置多字段 multipart
func TestRequestSetBodyMultipartWithFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		assert.NoError(t, err)

		// 检查普通字段
		assert.Equal(t, "test-user", r.FormValue("username"))
		assert.Equal(t, "user@test.com", r.FormValue("email"))

		w.Write([]byte("OK"))
	}))
	defer server.Close()

	fields := map[string]string{
		"username": "test-user",
		"email":    "user@test.com",
	}

	files := map[string]FileField{
		"avatar": {
			FileName: "avatar.png",
			Content:  []byte("fake image content"),
		},
	}

	client := &http.Client{}
	req := NewRequest(context.Background(), client, "POST", server.URL).
		SetBodyMultipartWithFields(fields, files)

	resp, err := req.Send()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestRequestMustSend 测试 MustSend 方法
func TestRequestMustSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Success"))
	}))
	defer server.Close()

	client := &http.Client{}
	req := NewRequest(context.Background(), client, "GET", server.URL)

	// 正常情况不应 panic
	assert.NotPanics(t, func() {
		resp := req.MustSend()
		assert.NotNil(t, resp)
	})
}

// TestRequestMustSendPanic 测试 MustSend 在错误时 panic
func TestRequestMustSendPanic(t *testing.T) {
	client := &http.Client{}
	req := NewRequest(context.Background(), client, "GET", "http://invalid-url-that-does-not-exist-12345.com")

	// 应该 panic
	assert.Panics(t, func() {
		req.MustSend()
	})
}

// TestRequestErrorHandling 测试错误处理
func TestRequestErrorHandling(t *testing.T) {
	req := setupRequest("POST", testURL).
		SetBodyMultipart("file", "test.txt", []byte("content"))

	// 模拟错误：设置一个会失败的条件
	req.err = assert.AnError

	resp, err := req.Send()
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Nil(t, resp.Response)
}

// TestChainedMethods 测试链式调用
func TestChainedMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "MyApp/1.0", r.Header.Get("User-Agent"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := &http.Client{}
	resp, err := NewRequest(context.Background(), client, "POST", server.URL).
		SetUserAgent("MyApp/1.0").
		SetContentType("application/json").
		SetAccept("application/json").
		SetQuery("page", "1").
		SetBodyJSON(map[string]string{"test": "data"}).
		Send()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
