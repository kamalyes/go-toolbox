/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-28 18:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-30 20:17:08
 * @FilePath: \go-toolbox\pkg\httpx\request.go
 * @Description: HTTP 请求封装
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/kamalyes/go-toolbox/pkg/errorx"
)

// Request 结构体用于封装 HTTP 请求的相关信息
type Request struct {
	ctx            context.Context
	client         *http.Client   // HTTP 客户端
	endpoint       string         // 请求的 URL
	method         string         // 请求方法（GET, POST, etc.）
	headers        http.Header    // 请求头
	queryValues    url.Values     // 查询参数
	body           any            // 请求体，可以是任意类型
	bodyBytes      io.Reader      // 请求体的字节流
	bodyEncodeFunc BodyEncodeFunc // 自定义的请求体编码函数
	err            error
}

// NewRequest 创建一个新的 HTTP 请求
func NewRequest(ctx context.Context, client *http.Client, method, endpoint string) *Request {
	return &Request{
		ctx:         ctx,
		client:      client,
		method:      method,
		endpoint:    endpoint,
		headers:     make(http.Header),
		queryValues: make(url.Values),
	}
}

// Getter 方法

// GetCtx 返回请求的上下文
func (r *Request) GetCtx() context.Context {
	return r.ctx
}

// GetClient 返回 HTTP 客户端
func (r *Request) GetClient() *http.Client {
	return r.client
}

// GetURL 返回请求的 URL
func (r *Request) GetURL() string {
	return r.endpoint
}

// GetMethod 返回请求的方法
func (r *Request) GetMethod() string {
	return r.method
}

// GetHeaders 返回请求头
func (r *Request) GetHeaders() http.Header {
	return r.headers
}

// GetQueryValues 返回查询参数
func (r *Request) GetQueryValues() url.Values {
	return r.queryValues
}

// GetBody 返回请求体
func (r *Request) GetBody() any {
	return r.body
}

// GetBodyBytes 返回请求体的字节流
func (r *Request) GetBodyBytes() io.Reader {
	return r.bodyBytes
}

// GetBodyEncodeFunc 返回自定义的请求体编码函数
func (r *Request) GetBodyEncodeFunc() BodyEncodeFunc {
	return r.bodyEncodeFunc
}

// GetError 返回错误信息
func (r *Request) GetError() error {
	return r.err
}

// Setter 方法

// AddQuery 添加查询参数
func (r *Request) AddQuery(key, value string) *Request {
	r.queryValues.Add(key, value) // 向查询参数中添加键值对
	return r
}

// SetQuery 设置查询参数
func (r *Request) SetQuery(key, value string) *Request {
	r.queryValues.Set(key, value) // 设置查询参数的值，如果已存在则覆盖
	return r
}

// SetHeader 设置请求头
func (r *Request) SetHeader(key, value string) *Request {
	r.headers.Set(key, value) // 设置请求头的键值对
	return r
}

// AddHeader 添加请求头
func (r *Request) AddHeader(key, value string) *Request {
	r.headers.Add(key, value) // 向请求头中添加键值对
	return r
}

// SetBodyEncodeFunc 设置请求体编码函数
func (r *Request) SetBodyEncodeFunc(fn BodyEncodeFunc) *Request {
	r.bodyEncodeFunc = fn // 设置自定义的请求体编码函数
	return r
}

// SetBody 设置请求体
func (r *Request) SetBody(body any) *Request {
	r.body = body
	if bodyReader, ok := body.(io.Reader); ok {
		r.bodyBytes = bodyReader // 如果请求体实现了 io.Reader 接口，则直接赋值
	} else {
		r.bodyEncodeFunc = encodeJSON // 设置默认编码函数
	}
	return r
}

// SetBodyForm 设置请求体为表单数据
func (r *Request) SetBodyForm(data url.Values) *Request {
	r.body = data
	r.bodyEncodeFunc = func(body any) (io.Reader, error) {
		return strings.NewReader(body.(url.Values).Encode()), nil
	}
	r.SetHeader("Content-Type", ContentTypeWWWFormURLEncoded) // 设置 Content-Type 为表单
	return r
}

// SetBodyMultipart 设置请求体为 multipart/form-data
func (r *Request) SetBodyMultipart(fieldName, fileName string, fileContent []byte) *Request {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件字段
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return r // 处理错误
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return r // 处理错误
	}

	// 关闭 writer，完成 multipart 数据
	err = writer.Close()
	if err != nil {
		return r // 处理错误
	}

	r.bodyBytes = &buf
	r.SetHeader("Content-Type", ContentTypeMultipartFormData)
	return r
}

// Send 执行 HTTP 请求
func (r *Request) Send() (Response, error) {
	if !IsValidMethod(r.method) {
		return Response{}, errorx.NewError(ErrInvalidMethod, r.method)
	}

	body, err := r.encodeBody() // 编码请求体
	if err != nil {
		return Response{}, err // 如果编码出错，记录错误并返回
	}

	req, err := http.NewRequestWithContext(r.ctx, r.method, r.endpoint, body)
	if err != nil {
		return Response{}, err // 如果请求创建失败，记录错误并返回
	}
	req.Header = r.headers // 设置请求头
	// 设置查询参数

	req.URL.RawQuery = r.queryValues.Encode()

	// 执行请求

	resp, err := r.client.Do(req)
	if err != nil {
		return Response{}, err // 如果请求执行出错，记录错误并返回
	} // 将原始 HTTP 响应赋值给 Response 结构体

	return Response{Response: resp}, nil
}

// encodeBody 编码请求体
func (r *Request) encodeBody() (io.Reader, error) {
	if r.bodyBytes != nil {
		return r.bodyBytes, nil // 如果 bodyBytes 已经被设置，直接返回
	}
	if r.body != nil {
		if r.bodyEncodeFunc == nil {
			return nil, errorx.NewError(ErrBodyEncodeFuncNotSet) // 检查 bodyEncodeFunc
		}
		body, err := r.bodyEncodeFunc(r.body) // 使用自定义编码函数编码请求体
		if err != nil {
			return nil, err // 如果编码出错，返回错误
		}
		r.bodyBytes = body // 将编码后的字节流赋值给 bodyBytes
		return body, nil
	}
	return nil, nil // 如果没有请求体，返回 nil
}

// JSON 编码请求体的默认函数
func encodeJSON(body any) (io.Reader, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBytes), nil
}
