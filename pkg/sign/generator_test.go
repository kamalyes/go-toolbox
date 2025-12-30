/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-30 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-30 15:15:56
 * @FilePath: \go-stress\go-toolbox\pkg\sign\generator_test.go
 * @Description: 签名生成器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sign

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGenerator(t *testing.T) {
	// 测试默认配置
	gen := NewGenerator(nil)
	assert.NotNil(t, gen)
	assert.Equal(t, "X-Sign", gen.config.HeaderName)
	assert.Equal(t, "X-Timestamp", gen.config.TimestampHeader)
	assert.Equal(t, AlgorithmSHA256, gen.config.Algorithm)

	// 测试自定义配置
	config := &GeneratorConfig{
		Enabled:         true,
		HeaderName:      "Custom-Sign",
		TimestampHeader: "Custom-Timestamp",
		SecretKey:       "test-secret",
		Algorithm:       AlgorithmSHA256,
	}
	gen = NewGenerator(config)
	assert.Equal(t, "Custom-Sign", gen.config.HeaderName)
	assert.Equal(t, "Custom-Timestamp", gen.config.TimestampHeader)
}

func TestGeneratorGenerateHeaders(t *testing.T) {
	config := &GeneratorConfig{
		Enabled:      true,
		SecretKey:    "test-secret-key",
		Algorithm:    AlgorithmSHA256,
		IncludeBody:  true,
		IncludeQuery: true,
	}

	gen := NewGenerator(config)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	queryParams := url.Values{
		"page": []string{"1"},
		"size": []string{"10"},
	}

	result := gen.GenerateHeaders("POST", "/api/users", `{"name":"test"}`, headers, queryParams)

	// 验证生成的 headers
	assert.NotEmpty(t, result["X-Sign"])
	assert.NotEmpty(t, result["X-Timestamp"])
	assert.NotEmpty(t, result["X-Nonce"])
	assert.Equal(t, "application/json", result["Content-Type"])
}

func TestGeneratorDisabled(t *testing.T) {
	config := &GeneratorConfig{
		Enabled: false,
	}

	gen := NewGenerator(config)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	result := gen.GenerateHeaders("GET", "/api/test", "", headers, nil)

	// 未启用时，应该返回原始 headers
	assert.Equal(t, headers, result)
	assert.Empty(t, result["X-Sign"])
}

func TestGeneratorWithExtraHeaders(t *testing.T) {
	config := &GeneratorConfig{
		Enabled:   true,
		SecretKey: "test-secret",
		Extra: map[string]string{
			"X-App-ID":  "test-app",
			"X-Version": "1.0.0",
		},
	}

	gen := NewGenerator(config)

	result := gen.GenerateHeaders("GET", "/api/test", "", map[string]string{}, nil)

	assert.Equal(t, "test-app", result["X-App-ID"])
	assert.Equal(t, "1.0.0", result["X-Version"])
	assert.NotEmpty(t, result["X-Nonce"], "应该包含 nonce header")
}

func TestGeneratorCustomFormat(t *testing.T) {
	config := &GeneratorConfig{
		Enabled:   true,
		SecretKey: "test-secret",
		Format:    "{method}\n{path}\n{timestamp}\n{nonce}",
	}

	gen := NewGenerator(config)

	result := gen.GenerateHeaders("POST", "/api/users", "", map[string]string{}, nil)

	assert.NotEmpty(t, result["X-Sign"])
	assert.NotEmpty(t, result["X-Timestamp"])
	assert.NotEmpty(t, result["X-Nonce"], "应该包含 nonce header")
}

func TestGeneratorIncludeHeaders(t *testing.T) {
	config := &GeneratorConfig{
		Enabled:   true,
		SecretKey: "test-secret",
		IncludeHeaders: []string{
			"Content-Type",
			"Authorization",
		},
	}

	gen := NewGenerator(config)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token123",
	}

	result := gen.GenerateHeaders("POST", "/api/users", "", headers, nil)

	assert.NotEmpty(t, result["X-Sign"])
}

func TestGeneratorVerify(t *testing.T) {
	config := &GeneratorConfig{
		Enabled:      true,
		SecretKey:    "test-secret",
		Algorithm:    AlgorithmSHA256,
		IncludeBody:  true,
		IncludeQuery: true,
	}

	gen := NewGenerator(config)

	method := "POST"
	path := "/api/users"
	body := `{"name":"test"}`
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	queryParams := url.Values{
		"page": []string{"1"},
	}

	// 生成签名
	result := gen.GenerateHeaders(method, path, body, headers, queryParams)
	signature := result["X-Sign"]
	timestamp := result["X-Timestamp"]

	// 验证签名
	isValid := gen.Verify(signature, method, path, body, headers, queryParams, timestamp)
	assert.True(t, isValid)

	// 验证错误的签名
	isValid = gen.Verify("wrong-signature", method, path, body, headers, queryParams, timestamp)
	assert.False(t, isValid)
}

func TestSortedQueryString(t *testing.T) {
	gen := NewGenerator(&GeneratorConfig{})

	tests := []struct {
		name     string
		params   url.Values
		expected string
	}{
		{
			name:     "空参数",
			params:   url.Values{},
			expected: "",
		},
		{
			name: "单个参数",
			params: url.Values{
				"key": []string{"value"},
			},
			expected: "key=value",
		},
		{
			name: "多个参数按字母序",
			params: url.Values{
				"c": []string{"3"},
				"a": []string{"1"},
				"b": []string{"2"},
			},
			expected: "a=1&b=2&c=3",
		},
		{
			name: "多值参数",
			params: url.Values{
				"tags": []string{"go", "test", "sign"},
			},
			expected: "tags=go&tags=test&tags=sign",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.sortedQueryString(tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildHeaderString(t *testing.T) {
	config := &GeneratorConfig{
		IncludeHeaders: []string{"Content-Type", "Authorization", "X-Custom"},
	}
	gen := NewGenerator(config)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
		"X-Custom":      "value",
		"X-Ignore":      "ignored",
	}

	result := gen.buildHeaderString(headers)

	// 应该按字母序排序，并且只包含配置中指定的 headers
	assert.Contains(t, result, "Authorization=Bearer token")
	assert.Contains(t, result, "Content-Type=application/json")
	assert.Contains(t, result, "X-Custom=value")
	assert.NotContains(t, result, "X-Ignore")
}

func BenchmarkGenerateHeaders(b *testing.B) {
	config := &GeneratorConfig{
		Enabled:      true,
		SecretKey:    "test-secret-key",
		Algorithm:    AlgorithmSHA256,
		IncludeBody:  true,
		IncludeQuery: true,
	}

	gen := NewGenerator(config)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	queryParams := url.Values{
		"page": []string{"1"},
		"size": []string{"10"},
	}

	body := `{"name":"test","age":30}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.GenerateHeaders("POST", "/api/users", body, headers, queryParams)
	}
}

// TestRealAPIRequest 测试真实的 API 请求场景
func TestRealAPIRequest(t *testing.T) {
	const secretKey = "test-secret-key-123"

	// 创建服务端签名验证器
	serverConfig := &GeneratorConfig{
		Enabled:      true,
		SecretKey:    secretKey,
		Algorithm:    AlgorithmSHA256,
		IncludeBody:  true,
		IncludeQuery: true,
	}
	serverGen := NewGenerator(serverConfig)

	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 读取请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "读取请求体失败", http.StatusBadRequest)
			return
		}

		// 获取签名相关的 headers
		signature := r.Header.Get("X-Sign")
		timestamp := r.Header.Get("X-Timestamp")
		nonce := r.Header.Get("X-Nonce")

		if signature == "" || timestamp == "" {
			http.Error(w, "缺少签名信息", http.StatusUnauthorized)
			return
		}

		// 验证 nonce 存在
		if nonce == "" {
			http.Error(w, "缺少 nonce", http.StatusUnauthorized)
			return
		}

		// 构建用于验证的 headers
		verifyHeaders := make(map[string]string)
		for k, v := range r.Header {
			if len(v) > 0 {
				verifyHeaders[k] = v[0]
			}
		}

		// 验证签名
		isValid := serverGen.Verify(
			signature,
			r.Method,
			r.URL.Path,
			string(body),
			verifyHeaders,
			r.URL.Query(),
			timestamp,
		)

		if !isValid {
			http.Error(w, "签名验证失败", http.StatusForbidden)
			return
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "签名验证通过",
		})
	}))
	defer server.Close()

	// 客户端生成签名并发送请求
	clientConfig := &GeneratorConfig{
		Enabled:      true,
		SecretKey:    secretKey,
		Algorithm:    AlgorithmSHA256,
		IncludeBody:  true,
		IncludeQuery: true,
	}
	clientGen := NewGenerator(clientConfig)

	// 准备请求数据
	requestBody := `{"username":"test","email":"test@example.com"}`
	queryParams := url.Values{
		"page": []string{"1"},
		"size": []string{"10"},
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 生成签名 headers
	signedHeaders := clientGen.GenerateHeaders("POST", "/api/users", requestBody, headers, queryParams)

	// 构建完整 URL
	fullURL := server.URL + "/api/users?" + queryParams.Encode()

	// 发送请求
	req, err := http.NewRequest("POST", fullURL, strings.NewReader(requestBody))
	assert.NoError(t, err)

	// 设置 headers
	for k, v := range signedHeaders {
		req.Header.Set(k, v)
	}

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 验证响应
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.True(t, result["success"].(bool))
	assert.Equal(t, "签名验证通过", result["message"])
}

// TestRealAPIRequestWithInvalidSignature 测试无效签名
func TestRealAPIRequestWithInvalidSignature(t *testing.T) {
	const secretKey = "test-secret-key-123"

	// 创建服务端验证器
	serverConfig := &GeneratorConfig{
		Enabled:     true,
		SecretKey:   secretKey,
		Algorithm:   AlgorithmSHA256,
		IncludeBody: true,
	}
	serverGen := NewGenerator(serverConfig)

	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		signature := r.Header.Get("X-Sign")
		timestamp := r.Header.Get("X-Timestamp")

		verifyHeaders := make(map[string]string)
		for k, v := range r.Header {
			if len(v) > 0 {
				verifyHeaders[k] = v[0]
			}
		}

		isValid := serverGen.Verify(signature, r.Method, r.URL.Path, string(body), verifyHeaders, r.URL.Query(), timestamp)

		if !isValid {
			http.Error(w, "签名验证失败", http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 发送错误签名的请求
	req, _ := http.NewRequest("POST", server.URL+"/api/test", strings.NewReader(`{"test":"data"}`))
	req.Header.Set("X-Sign", "invalid-signature")
	req.Header.Set("X-Timestamp", "1234567890")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 应该返回 403
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

// TestRealAPIRequestWithQueryParams 测试查询参数排序一致性
func TestRealAPIRequestWithQueryParams(t *testing.T) {
	const secretKey = "test-secret-key-123"

	serverConfig := &GeneratorConfig{
		Enabled:      true,
		SecretKey:    secretKey,
		Algorithm:    AlgorithmSHA256,
		IncludeQuery: true,
	}
	serverGen := NewGenerator(serverConfig)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Sign")
		timestamp := r.Header.Get("X-Timestamp")

		verifyHeaders := make(map[string]string)
		for k, v := range r.Header {
			if len(v) > 0 {
				verifyHeaders[k] = v[0]
			}
		}

		// 服务端验证时会自动排序查询参数
		isValid := serverGen.Verify(signature, r.Method, r.URL.Path, "", verifyHeaders, r.URL.Query(), timestamp)

		if !isValid {
			t.Logf("验证失败 - 方法:%s, 路径:%s", r.Method, r.URL.Path)
			t.Logf("查询参数: %v", r.URL.Query())
			http.Error(w, "签名验证失败", http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	clientConfig := &GeneratorConfig{
		Enabled:      true,
		SecretKey:    secretKey,
		Algorithm:    AlgorithmSHA256,
		IncludeQuery: true,
	}
	clientGen := NewGenerator(clientConfig)

	// 测试不同顺序的查询参数
	tests := []struct {
		name   string
		params url.Values
	}{
		{
			name: "字母序参数",
			params: url.Values{
				"a": []string{"1"},
				"b": []string{"2"},
				"c": []string{"3"},
			},
		},
		{
			name: "逆序参数",
			params: url.Values{
				"z": []string{"26"},
				"y": []string{"25"},
				"x": []string{"24"},
			},
		},
		{
			name: "随机序参数",
			params: url.Values{
				"page":   []string{"1"},
				"size":   []string{"10"},
				"sort":   []string{"desc"},
				"filter": []string{"active"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 客户端生成签名
			signedHeaders := clientGen.GenerateHeaders("GET", "/api/test", "", map[string]string{}, tt.params)

			// 构建 URL（参数可能是乱序的）
			fullURL := server.URL + "/api/test?" + tt.params.Encode()

			req, _ := http.NewRequest("GET", fullURL, nil)
			for k, v := range signedHeaders {
				req.Header.Set(k, v)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// 应该验证成功
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

// TestRealAPIRequestWithMultiValueParams 测试多值查询参数
func TestRealAPIRequestWithMultiValueParams(t *testing.T) {
	const secretKey = "test-secret-key-123"

	serverConfig := &GeneratorConfig{
		Enabled:      true,
		SecretKey:    secretKey,
		Algorithm:    AlgorithmSHA256,
		IncludeQuery: true,
	}
	serverGen := NewGenerator(serverConfig)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get("X-Sign")
		timestamp := r.Header.Get("X-Timestamp")

		verifyHeaders := make(map[string]string)
		for k, v := range r.Header {
			if len(v) > 0 {
				verifyHeaders[k] = v[0]
			}
		}

		isValid := serverGen.Verify(signature, r.Method, r.URL.Path, "", verifyHeaders, r.URL.Query(), timestamp)

		if !isValid {
			http.Error(w, "签名验证失败", http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	clientConfig := &GeneratorConfig{
		Enabled:      true,
		SecretKey:    secretKey,
		Algorithm:    AlgorithmSHA256,
		IncludeQuery: true,
	}
	clientGen := NewGenerator(clientConfig)

	// 多值参数
	queryParams := url.Values{
		"tags": []string{"go", "test", "api"},
		"ids":  []string{"1", "2", "3"},
	}

	signedHeaders := clientGen.GenerateHeaders("GET", "/api/search", "", map[string]string{}, queryParams)

	fullURL := server.URL + "/api/search?" + queryParams.Encode()
	req, _ := http.NewRequest("GET", fullURL, nil)
	for k, v := range signedHeaders {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
