# httpx 标准 API 对照表

本文档说明了 httpx 包与 `net/http` 标准库的 API 对照关系。

## Request 标准方法

### 上下文管理

| 标准方法 | Deprecated 方法 | 说明 |
|---------|----------------|------|
| `Context() context.Context` | `GetCtx()` | 返回请求的上下文 |
| `WithContext(ctx) *Request` | - | 返回使用新上下文的请求副本 |

### 请求信息访问

| 标准方法 | Deprecated 方法 | 说明 |
|---------|----------------|------|
| `URL() string` | `GetURL()` | 返回请求的 URL |
| `Method() string` | `GetMethod()` | 返回请求的方法 |
| `Header() http.Header` | `GetHeaders()` | 返回请求头 |
| `Query() url.Values` | `GetQueryValues()` | 返回查询参数 |
| `FullURL() string` | `GetFullURL()` | 返回包含查询参数的完整 URL |
| `Client() *http.Client` | `GetClient()` | 返回 HTTP 客户端 |
| `Error() error` | `GetError()` | 返回错误信息 |

### 内部使用方法（不推荐外部调用）

- `GetBody()` - 返回请求体（内部使用）
- `GetBodyBytes()` - 返回请求体字节流（内部使用）
- `GetBodyEncodeFunc()` - 返回编码函数（内部使用）

### Cookie 管理

| 标准方法 | 说明 |
|---------|------|
| `Cookie(name string) (*http.Cookie, error)` | 获取指定名称的 Cookie |
| `AddCookie(cookie *http.Cookie) *Request` | 添加 Cookie 到请求 |

## Response 标准方法

### 响应体读取

| 标准方法 | Deprecated 方法 | 说明 |
|---------|----------------|------|
| `Body() ([]byte, error)` | `GetBody()` | 读取响应体（标准方法） |
| `Bytes() ([]byte, error)` | - | Body() 的别名 |
| `String() (string, error)` | - | 读取响应体并转换为字符串 |

### 响应体解码

| 标准方法 | Deprecated 方法 | 说明 |
|---------|----------------|------|
| `JSON(dst any) error` | - | 解码 JSON 响应体 |
| `XML(dst any) error` | - | 解码 XML 响应体 |
| `Decode(dst any) error` | `DecodeRespBody(dst)` | 根据 Content-Type 自动解码 |

### 错误处理

| 标准方法 | Deprecated 方法 | 说明 |
|---------|----------------|------|
| `Error() error` | `GetError()` | 返回错误信息 |
| `IsError() bool` | - | 检查是否有错误 |
| `OK() bool` | - | 检查状态码是否为 200 |

### 其他方法

| 标准方法 | Deprecated 方法 | 说明 |
|---------|----------------|------|
| `Cookies() []*http.Cookie` | - | 获取响应的 Cookie（继承自 http.Response） |
| - | `GetCookies()` | 获取响应的 Cookie（废弃，直接使用 Cookies()） |

## Client 标准方法

### 创建客户端

```go
// 标准方法：使用函数选项模式
client := httpx.NewClient(
    httpx.WithTimeout(30 * time.Second),
    httpx.WithMaxIdleConns(100),
    httpx.WithMaxIdleConnsPerHost(10),
)

// 便捷方法
client := httpx.NewDefaultHttpClient(30 * time.Second)
```

### 函数选项

- `WithTimeout(duration)` - 设置超时时间
- `WithMaxIdleConns(n)` - 设置最大空闲连接数
- `WithMaxIdleConnsPerHost(n)` - 设置每个 Host 的最大空闲连接数
- `WithMaxConnsPerHost(n)` - 设置每个 Host 的最大连接数
- `WithIdleConnTimeout(duration)` - 设置空闲连接超时
- `WithTLSHandshakeTimeout(duration)` - 设置 TLS 握手超时
- `WithExpectContinueTimeout(duration)` - 设置 Expect: 100-continue 超时
- `WithResponseHeaderTimeout(duration)` - 设置响应头超时
- `WithInsecureSkipVerify(skip)` - 是否跳过 TLS 证书验证

## 参数构建助手

### BuildParams

```go
params := httpx.BuildParams(
    url.Values{"api_key": []string{apiKey}},
    httpx.WithParam(autoRenew, "auto_renew", "1"),
    httpx.WithParamNotEmpty("domain", domainName),
)
```

### WithParam

```go
// 条件参数：只有当 condition 为 true 时才添加
httpx.WithParam(condition bool, key, value string) func(url.Values)
```

### WithParamNotEmpty

```go
// 非空参数：只有当 value 非空时才添加
httpx.WithParamNotEmpty(key, value string) func(url.Values)
```

## 迁移指南

### Request 迁移

```go
// 旧代码
ctx := req.GetCtx()
url := req.GetURL()
headers := req.GetHeaders()
err := req.GetError()

// 新代码（标准方法）
ctx := req.Context()
url := req.URL()
headers := req.Header()
err := req.Error()
```

### Response 迁移

```go
// 旧代码
body, err := resp.GetBody()
err := resp.GetError()
err := resp.DecodeRespBody(&result)

// 新代码（标准方法）
body, err := resp.Body()        // 或 resp.Bytes()
err := resp.Error()
err := resp.Decode(&result)     // 或 resp.JSON(&result)
```

### 注意事项

1. **Body 方法名冲突**：由于我们添加了 `Body()` 方法，在访问原生 `http.Response.Body` 时需要使用 `resp.Response.Body`
2. **Deprecated 方法**：所有 `GetXxx()` 方法都已废弃，建议使用新的标准方法
3. **向后兼容**：所有废弃方法仍然可用，但会在内部调用新的标准方法
4. **编译器警告**：使用废弃方法时，IDE 会显示 `Deprecated` 提示

## 设计原则

1. **与 net/http 保持一致**：方法名和行为尽可能与标准库保持一致
2. **链式调用支持**：所有 Setter 方法返回 `*Request` 以支持链式调用
3. **错误处理**：统一使用 `Error()` 方法获取错误，而不是多个 `GetError()` 方法
4. **简洁的 API**：提供 `Bytes()`、`String()` 等便捷方法
5. **类型安全**：所有方法都有明确的类型定义
