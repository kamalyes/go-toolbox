/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-09 00:25:15
 * @FilePath: \go-toolbox\tests\location_amap_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/location"
	"github.com/stretchr/testify/assert"
)

// 模拟 Amap API 的响应
const mockResponse = `{
	"status": "1",
	"info": "OK",
	"infocode": "10000",
	"location": "116.481028,39.989643"
}`

// TestGetGPSByIpAmap 测试 GetGPSByIpAmap 函数
func TestGetGPSByIpAmap(t *testing.T) {
	// 创建一个新的 HTTP 服务器，返回模拟响应
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close() // 确保在测试结束时关闭服务器

	// 准备测试参数
	amapKey := "test_key"
	amapSign := "test_sign"
	amapUrl := server.URL // 使用模拟服务器的 URL
	ip := "8.8.8.8"       // 示例 IP 地址

	// 调用被测试的函数
	data, err := location.GetGPSByIpAmap(amapKey, amapSign, amapUrl, ip)
	assert.NoError(t, err, "期望没有错误，但得到 %v", err)

	// 验证响应
	assert.Equal(t, "1", data["status"], "期望状态为 '1'，但得到 %v", data["status"])

	location, ok := data["location"].(string)
	assert.True(t, ok, "期望 location 为字符串")
	assert.Equal(t, "116.481028,39.989643", location, "期望 location 为 '116.481028,39.989643'，但得到 %v", location)
}

// TestGetGPSByIpAmap_Error 测试 GetGPSByIpAmap 函数的错误情况
func TestGetGPSByIpAmap_Error(t *testing.T) {
	// 创建一个新的 HTTP 服务器，返回错误响应
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close() // 确保在测试结束时关闭服务器

	// 准备测试参数
	amapKey := "test_key"
	amapSign := "test_sign"
	amapUrl := server.URL // 使用模拟服务器的 URL
	ip := "8.8.8.8"       // 示例 IP 地址

	// 调用被测试的函数
	data, err := location.GetGPSByIpAmap(amapKey, amapSign, amapUrl, ip)
	assert.Error(t, err, "期望有错误，但没有得到")

	// 验证数据为空
	assert.Nil(t, data, "期望数据为 nil，但得到 %v", data)
}

// TestGetGPSByIpAmap_MissingConfig 测试 GetGPSByIpAmap 函数缺少配置的情况
func TestGetGPSByIpAmap_MissingConfig(t *testing.T) {
	// 调用被测试的函数，缺少参数
	data, err := location.GetGPSByIpAmap("", "", "", "")
	assert.Error(t, err, "期望缺少配置时有错误，但没有得到")

	// 验证数据为空
	assert.Nil(t, data, "期望数据为 nil，但得到 %v", data)
}
