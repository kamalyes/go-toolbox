---
name: httpx-client
description: HTTP客户端工具，提供可配置的HTTP客户端、请求构建、参数处理、响应解码、Cookie管理。当需要构建HTTP请求、配置连接池/超时/TLS、或解析响应体时使用。
---

# httpx - HTTP客户端

提供可配置的HTTP客户端、请求构建、参数处理与响应解码。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/httpx"
```

创建客户端：
```go
client := httpx.NewClient(httpx.WithTimeout(10 * time.Second))
```

构建请求：
```go
req := httpx.NewRequest(ctx, client, "GET", "https://api.example.com/users")
```

## 完整API索引

### 函数

#### 客户端构建

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewClient` | `func(opts ...ClientOption) *Client` | 创建可配置的HTTP客户端 |
| `NewHttpClient` | `func(client *http.Client) *Client` | 从标准http.Client创建 |
| `WithTimeout` | `func(timeout time.Duration) ClientOption` | 设置超时 |
| `WithMaxIdleConns` | `func(n int) ClientOption` | 设置最大空闲连接数 |
| `WithMaxIdleConnsPerHost` | `func(n int) ClientOption` | 设置每主机最大空闲连接数 |
| `WithMaxConnsPerHost` | `func(n int) ClientOption` | 设置每主机最大连接数 |
| `WithIdleConnTimeout` | `func(timeout time.Duration) ClientOption` | 设置空闲连接超时 |
| `WithTLSHandshakeTimeout` | `func(timeout time.Duration) ClientOption` | 设置TLS握手超时 |
| `WithInsecureSkipVerify` | `func(skip bool) ClientOption` | 设置跳过TLS验证 |
| `WithContext` | `func(ctx context.Context) ClientOption` | 设置上下文 |

#### 请求构建

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewRequest` | `func(ctx, client, method, endpoint) *Request` | 创建HTTP请求 |
| `ReadAndCacheResponseBody` | `func(resp *http.Response) ([]byte, error)` | 读取并缓存响应体 |
| `DecodeRespBody` | `func(resp *http.Response, v interface{}) error` | 解码响应体 |
| `NormalizeBaseURL` | `func(url string) string` | 规范化基础URL |
| `IsValidMethod` | `func(method string) bool` | 校验HTTP方法 |
| `BuildParams` | `func(params map[string]interface{}) string` | 构建查询参数 |
| `WithParam` | `func(params, key, value string) string` | 添加查询参数 |
| `WithParamNotEmpty` | `func(params, key, value string) string` | 添加非空查询参数 |
| `GetRequestValue` | `func(r *http.Request, key string) string` | 从请求头或查询参数获取值 |
| `GetValueFromHeaderOrQuery` | `func(r *http.Request, key string) string` | 从头部或查询参数获取值 |
| `ReadRequestBody` | `func(r *http.Request) ([]byte, error)` | 读取请求体 |
| `NewParams` | `func() *ParamsBuilder` | 创建参数构建器 |
| `NewParamsWithBase` | `func(base map[string]string) *ParamsBuilder` | 从基础map创建参数构建器 |
| `GetCookies` | `func(url string) ([]*http.Cookie, error)` | 获取URL的cookies |

### 类型

| 导出名称 | 说明 |
|---|---|
| `Client` | HTTP客户端类型 |
| `ClientOption` | 客户端配置选项类型 |
| `Response` | 响应封装类型 |
| `Request` | 请求构建类型 |
| `FileField` | 文件上传字段类型 |
| `BodyEncodeFunc` | 请求体编码函数类型 |
| `ParamsBuilder` | 参数构建器类型 |

## 注意事项

- `WithInsecureSkipVerify` 仅用于测试环境，生产环境请勿跳过TLS验证
- `ReadAndCacheResponseBody` 读取后允许响应体被多次读取
- `GetRequestValue` 优先从请求头获取，其次从查询参数获取