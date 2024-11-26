/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-28 18:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-30 20:22:37
 * @FilePath: \go-toolbox\pkg\httpx\response_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseDecodeRespBodyJSON(t *testing.T) {
	// 准备一个 JSON 响应
	expected := map[string]string{"name": "陈明勇"}
	body, err := json.Marshal(expected)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader(body))
	req.Header.Set(HeaderContentType, ContentTypeApplicationJSON)

	rr := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderContentType, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}).ServeHTTP(rr, req)

	// 修正这里的字段名
	resp := &Response{Response: rr.Result()}

	var result map[string]string
	err = resp.DecodeRespBody(&result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestResponseDecodeRespBodyXML(t *testing.T) {
	type ResponseData struct {
		Name string `xml:"name"`
	}
	// 准备一个 XML 响应
	expected := ResponseData{Name: "陈明勇"}
	body, err := xml.Marshal(expected)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader(body))
	req.Header.Set("Content-Type", ContentTypeTextXML)

	rr := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ContentTypeTextXML)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}).ServeHTTP(rr, req)

	resp := &Response{Response: rr.Result()} // 修正这里的字段名

	var result ResponseData
	err = resp.DecodeRespBody(&result) // 确保这里是指向结构体的指针
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestResponseDecodeRespBodyTextPlain(t *testing.T) {
	// 准备一个纯文本响应
	expected := "hello world"
	body := []byte(expected)

	req := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader(body))
	req.Header.Set(HeaderContentType, ContentTypeTextPlain)

	rr := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderContentType, ContentTypeTextPlain)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}).ServeHTTP(rr, req)

	resp := &Response{Response: rr.Result()} // 修正这里的字段名

	var result string
	err := resp.DecodeRespBody(&result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestResponseDecodeRespBodyInvalidContentType(t *testing.T) {
	// 准备一个无效的 Content-Type 响应
	body := []byte("invalid content")
	req := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader(body))
	req.Header.Set(HeaderContentType, "application/unknown")

	rr := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}).ServeHTTP(rr, req)

	resp := &Response{Response: rr.Result()} // 修正这里的字段名

	var result string
	err := resp.DecodeRespBody(&result)
	assert.Error(t, err)
}

func TestResponseDecodeRespBodyTextPlainError(t *testing.T) {
	// 准备一个纯文本响应
	expected := "hello world"
	body := []byte(expected)

	req := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader(body))
	req.Header.Set(HeaderContentType, ContentTypeTextPlain)

	rr := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderContentType, ContentTypeTextPlain)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}).ServeHTTP(rr, req)

	resp := &Response{Response: rr.Result()} // 修正这里的字段名

	var result int // 错误的目标类型
	err := resp.DecodeRespBody(&result)
	assert.Error(t, err)
	assert.Equal(t, "expected dst to be *string, but got *int", err.Error())
}
