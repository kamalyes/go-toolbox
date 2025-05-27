/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-28 18:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-05-27 15:05:08
 * @FilePath: \go-toolbox\pkg\httpx\response.go
 * @Description: HTTP 响应封装
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"

	"github.com/kamalyes/go-toolbox/pkg/errorx"
)

// Response 结构体用于封装 HTTP 响应
type Response struct {
	*http.Response       // 原始 HTTP 响应
	err            error // 处理过程中可能出现的错误
}

// IsError 检查响应是否有错误
func (r *Response) IsError() bool {
	return r.err != nil
}

// GetError 返回错误信息
func (r *Response) GetError() error {
	return r.err
}

// Close 关闭 HTTP 响应体
func (r *Response) Close() error {
	if r.Response != nil && r.Response.Body != nil {
		return r.Response.Body.Close() // 关闭响应体
	}
	return nil
}

// CheckStatus 检查响应状态码
func (r *Response) CheckStatus() error {
	if r.IsError() {
		return r.GetError() // 如果有错误，直接返回
	}
	if r.StatusCode != http.StatusOK {
		return errorx.NewError(ErrRequestStatusCode, r.Status) // 检查状态码
	}
	return nil
}

// LogResponse 日志记录响应信息
func (r *Response) LogResponse() {
	if r != nil && r.Response != nil {
		log.Printf("Request to %s returned status %d", r.Request.URL, r.StatusCode) // 记录请求的 URL 和状态码
	}
}

// ReadAndCacheResponseBody 读取并缓存响应体
func ReadAndCacheResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close() // 确保关闭响应体
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf) // 同时读取和缓存内容

	// 读取内容并忽略返回的内容
	if _, err := io.ReadAll(tee); err != nil {
		return "", err // 读取出错，返回错误
	}
	return buf.String(), nil // 返回缓存的内容
}

// DecodeRespBody 解码响应体到目标结构体
func (r *Response) DecodeRespBody(dst any) error {
	if r.IsError() {
		return r.GetError() // 如果有错误，直接返回
	}
	return DecodeRespBody(r, dst) // 调用 DecodeRespBody 解码响应体
}

// GetBody 读取响应体
func (r *Response) GetBody() ([]byte, error) {
	if r.IsError() {
		return nil, r.GetError()
	}
	defer r.Close() // 确保响应体在使用后关闭
	body, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetCookies 获取响应中的 Cookie
func (r *Response) GetCookies() ([]*http.Cookie, error) {
	return r.Cookies(), nil
}

// DecodeRespBody 根据响应的 Content-Type 解码响应体
func DecodeRespBody(resp *Response, dst any) error {
	contentType := resp.Header.Get(HeaderContentType) // 获取响应的 Content-Type
	switch contentType {
	case ContentTypeApplicationJSON, ContentTypeApplicationJSONCharacterUTF8:
		// 如果是 JSON 格式，使用 JSON 解码器解码
		return json.NewDecoder(resp.Body).Decode(dst)
	case ContentTypeApplicationXML, ContentTypeApplicationXMLCharacterUTF8, ContentTypeTextXML, ContentTypeTextXMLCharacterUTF8:
		// 如果是 XML 格式，使用 XML 解码器解码
		return xml.NewDecoder(resp.Body).Decode(dst)
	case ContentTypeTextPlain, ContentTypeTextPlainCharacterUTF8:
		// 如果是纯文本格式，读取响应体并赋值给目标字符串
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err // 读取出错，返回错误
		}

		// 确保目标是字符串指针类型
		strPtr, ok := dst.(*string)
		if !ok {
			return errorx.NewError(ErrExpectedDestinationType, dst) // 类型不匹配，返回错误
		}
		*strPtr = string(bytes) // 将读取的字节转换为字符串并赋值
		return nil
	default:
		// 不支持的 Content-Type，返回错误
		return errorx.NewError(ErrUnsupportedContentType, contentType)
	}
}
