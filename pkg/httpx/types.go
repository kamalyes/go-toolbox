/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-28 18:55:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-28 18:55:55
 * @FilePath: \go-toolbox\pkg\httpx\types.go
 * @Description: HTTP 相关常量和类型定义
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package httpx

import "io"

// HTTP 请求和响应相关的常量定义
const (
	// HTTP Headers
	HeaderContentType   = "Content-Type"  // HTTP 请求和响应头中的 Content-Type
	HeaderAccept        = "Accept"        // HTTP 请求头中的 Accept
	HeaderAuthorization = "Authorization" // HTTP 请求头中的 Authorization
	HeaderUserAgent     = "User-Agent"    // HTTP 请求头中的 User-Agent
	HeaderCookie        = "Cookie"        // HTTP 请求头中的 Cookie

	// 常见的 Content-Type
	ContentTypeTextPlain                    = "text/plain"                        // 纯文本格式
	ContentTypeTextPlainCharacterUTF8       = "text/plain; charset=utf-8"         // UTF-8 编码的纯文本
	ContentTypeApplicationJSON              = "application/json"                  // JSON 格式
	ContentTypeApplicationXNDJSON           = "application/x-ndjson"              // NDJSON 格式
	ContentTypeApplicationXSnappyFramed     = "application/x-snappy-framed"       // Snappy 压缩格式的 JSON
	ContentTypeApplicationJSONCharacterUTF8 = "application/json; charset=utf-8"   // UTF-8 编码的 JSON
	ContentTypeApplicationXML               = "application/xml"                   // XML 格式
	ContentTypeApplicationXMLCharacterUTF8  = "application/xml; charset=utf-8"    // UTF-8 编码的 XML
	ContentTypeTextHTML                     = "text/html"                         // HTML 文本格式
	ContentTypeTextHTMLCharacterUTF8        = "text/html; charset=utf-8"          // UTF-8 编码的 HTML 文本
	ContentTypeTextXML                      = "text/xml"                          // XML 文本格式
	ContentTypeTextXMLCharacterUTF8         = "text/xml; charset=utf-8"           // UTF-8 编码的 XML 文本
	ContentTypeTextCSV                      = "text/csv"                          // CSV 文本格式
	ContentTypeTextCSVCharacterUTF8         = "text/csv; charset=utf-8"           // UTF-8 编码的 CSV 文本
	ContentTypeImageJPEG                    = "image/jpeg"                        // JPEG 图片格式
	ContentTypeImagePNG                     = "image/png"                         // PNG 图片格式
	ContentTypeImageGIF                     = "image/gif"                         // GIF 图片格式
	ContentTypeImageWebP                    = "image/webp"                        // WebP 图片格式
	ContentTypeApplicationOctetStream       = "application/octet-stream"          // 二进制流格式
	ContentTypeMultipartFormData            = "multipart/form-data"               // 表单数据格式
	ContentTypeWWWFormURLEncoded            = "application/x-www-form-urlencoded" // URL 编码的表单数据
)

// BodyEncodeFunc 定义请求体编码函数类型
// 该函数接收任意类型的请求体并返回 io.Reader 和可能的错误
type BodyEncodeFunc func(body any) (io.Reader, error)
