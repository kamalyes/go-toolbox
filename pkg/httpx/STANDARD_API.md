# httpx 标准 API 对照表

本文档说明了 httpx 包与 `net/http` 标准库的 API 对照关系

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
| - | `GetQueryValues()` | 返回查询参数（注：暂无标准方法别名，使用 GetQueryValues()） |
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
| - | `GetCookies() ([]*http.Cookie, error)` | 获取响应的 Cookie（废弃，直接使用 Cookies()） |

## Client 标准方法

### 创建客户端

```go
// 标准方法：使用函数选项模式
client := httpx.NewClient(
    httpx.WithTimeout(30 * time.Second),
    httpx.WithMaxIdleConns(100),
    httpx.WithMaxIdleConnsPerHost(10),
)

// 便捷方法（使用默认 http.Client，无参数）
client := httpx.NewDefaultHttpClient()

// 使用自定义 http.Client
client := httpx.NewHttpClient(myHTTPClient)

// 使用自定义 http.Client 和上下文
client := httpx.NewClientWithContext(myHTTPClient, ctx)

// 使用自定义配置的 http.Client
client := httpx.NewCustomDefaultClient()
```

### 函数选项

- `WithTimeout(duration)` - 设置请求超时时间
- `WithMaxIdleConns(n)` - 设置最大空闲连接数
- `WithMaxIdleConnsPerHost(n)` - 设置每个主机最大空闲连接数
- `WithMaxConnsPerHost(n)` - 设置每个主机最大连接数
- `WithIdleConnTimeout(duration)` - 设置空闲连接超时时间
- `WithTLSHandshakeTimeout(duration)` - 设置 TLS 握手超时时间
- `WithInsecureSkipVerify(skip)` - 是否跳过 TLS 证书验证（生产环境慎用）
- `WithContext(ctx)` - 设置请求上下文

### 客户端请求方法

Client 提供了所有 HTTP 方法的便捷创建函数：

```go
client.Get(endpoint)       // GET 请求
client.Post(endpoint)      // POST 请求
client.Put(endpoint)       // PUT 请求
client.Delete(endpoint)    // DELETE 请求
client.Patch(endpoint)     // PATCH 请求
client.Head(endpoint)      // HEAD 请求
client.Options(endpoint)   // OPTIONS 请求
client.Connect(endpoint)   // CONNECT 请求
client.Trace(endpoint)     // TRACE 请求
client.NewRequest(method, endpoint) // 自定义方法
```

## 参数构建助手

### BuildParams

```go
// 基础参数通过 base 传入，可选参数通过 opts 传入
func BuildParams(base map[string]string, opts ...func(map[string]string)) map[string]string

// 条件参数：只有当 condition 为 true 时才添加
func WithParam(condition bool, key, value string) func(map[string]string)

// 非空参数：只有当 value 非空时才添加
func WithParamNotEmpty(key, value string) func(map[string]string)
```

使用示例：

```go
params := httpx.BuildParams(
    map[string]string{"api_key": apiKey},
    httpx.WithParam(autoRenew, "auto_renew", "1"),
    httpx.WithParamNotEmpty("domain", domainName),
)
```

### ParamsBuilder 参数构建器

支持链式调用的参数构建器，提供更丰富的操作：

```go
builder := httpx.NewParams().
    Set("key1", "value1").
    SetIf(condition, "key2", "value2").
    SetNotEmpty("key3", value3).
    SetAny("count", 42)

// 获取结果
params := builder.Build()       // map[string]string
slice := builder.ToSlice()      // []interface{}{"key1", "value1", ...}（用于日志）
keys := builder.Keys()
values := builder.Values()
len := builder.Len()
has := builder.Has("key1")      // true

// 其他操作
builder.Delete("key1")
builder.Clear()
clone := builder.Clone()
merged := builder.Merge(otherBuilder)

// 从基础参数创建
builder := httpx.NewParamsWithBase(baseMap)

// 批量设置
builder.SetMultiple(map[string]string{"a": "1", "b": "2"})
```

## Request 请求体设置方法

Request 提供了丰富的链式 Setter 方法：

### 查询参数

```go
req.AddQuery(key, value string)         // 添加查询参数
req.SetQuery(key, value string)         // 设置查询参数（覆盖）
req.SetQueries(queries map[string]string) // 批量设置
req.AddQueries(queries map[string]string) // 批量添加
req.SetQueryValues(values url.Values)   // 直接设置 url.Values
```

### 请求头

```go
req.SetHeader(key, value string)        // 设置请求头
req.AddHeader(key, value string)        // 添加请求头
req.SetHeaders(headers map[string]string) // 批量设置
req.AddHeaders(headers map[string]string) // 批量添加
req.SetUserAgent(ua string)             // 设置 User-Agent
req.SetAuthorization(token string)      // 设置 Authorization
req.SetBearerToken(token string)        // 设置 Bearer Token
req.SetContentType(contentType string)  // 设置 Content-Type
req.SetAccept(accept string)            // 设置 Accept
```

### 请求体

```go
req.SetBody(body any)                   // 设置请求体（自动检测 io.Reader 或 JSON 编码）
req.SetBodyJSON(body any)               // 设置 JSON 请求体（自动设置 Content-Type）
req.SetBodyString(body string)          // 设置字符串请求体
req.SetBodyRaw(body []byte)             // 设置原始字节请求体
req.SetBodyForm(data url.Values)        // 设置表单请求体
req.SetBodyMultipart(fieldName, fileName string, content []byte) // 单文件上传
req.SetBodyMultipartWithFields(fields map[string]string, files map[string]FileField) // 多字段+多文件上传
req.SetBodyEncodeFunc(fn BodyEncodeFunc) // 自定义编码函数
```

### 其他方法

```go
req.SetEndpoint(endpoint string)        // 设置请求 URL
req.Clone() *Request                    // 克隆请求
req.Send() (Response, error)            // 执行请求
req.Do(ctx context.Context) ([]byte, error) // 执行请求并返回响应体字节
req.MustSend() Response                 // 执行请求，失败则 panic
```

## HTTP 辅助函数

### 请求值获取

```go
// 从 HTTP 请求中获取指定值（查找顺序：Context → Header → Query）
func GetRequestValue(r *http.Request, contextKey interface{}, headerName, queryName string) string

// 从请求头或查询参数中获取值（查找顺序：Header → Query）
func GetValueFromHeaderOrQuery(r *http.Request, headerName, queryName string) string
```

### 请求体读取

```go
// 读取 HTTP 请求体（支持重复读取，读取后重新设置 Body）
func ReadRequestBody(r *http.Request) ([]byte, error)
```

### Cookie 获取

```go
// 发送 GET 请求获取指定 URL 的 Cookie
func GetCookies(url string) ([]*http.Cookie, error)
```

### URL 工具函数

```go
// 规范化基础 URL，确保包含协议前缀
// 不含 http:// 或 https:// 时自动添加 https://
func NormalizeBaseURL(baseURL string) string
```

示例：

```go
httpx.NormalizeBaseURL("www.example.com")     // "https://www.example.com"
httpx.NormalizeBaseURL("http://example.com")  // "http://example.com"
httpx.NormalizeBaseURL("")                    // ""
```

### 方法校验

```go
// 校验 HTTP 方法是否有效
func IsValidMethod(method string) bool
```

### 响应辅助方法

```go
resp.Close() error                    // 关闭响应体
resp.CheckStatus() error              // 检查状态码是否为 200
resp.LogResponse()                    // 记录响应日志

// 读取并缓存响应体（用于多次访问）
func ReadAndCacheResponseBody(resp *http.Response) (string, error)
```

## 类型与常量

### 核心类型

```go
// 请求体编码函数类型
type BodyEncodeFunc func(body any) (io.Reader, error)

// 文件字段结构（用于 multipart 上传）
type FileField struct {
    FileName string
    Content  []byte
}
```

### HTTP Header 常量

```go
HeaderContentType   = "Content-Type"
HeaderAccept        = "Accept"
HeaderAuthorization = "Authorization"
HeaderUserAgent     = "User-Agent"
HeaderCookie        = "Cookie"
```

### Content-Type 常量

```go
ContentTypeApplicationJSON              = "application/json"
ContentTypeApplicationJSONCharacterUTF8 = "application/json; charset=utf-8"
ContentTypeApplicationXML               = "application/xml"
ContentTypeApplicationXMLCharacterUTF8  = "application/xml; charset=utf-8"
ContentTypeTextPlain                    = "text/plain"
ContentTypeTextPlainCharacterUTF8       = "text/plain; charset=utf-8"
ContentTypeTextHTML                     = "text/html"
ContentTypeTextCSV                      = "text/csv"
ContentTypeApplicationOctetStream       = "application/octet-stream"
ContentTypeMultipartFormData            = "multipart/form-data"
ContentTypeWWWFormURLEncoded            = "application/x-www-form-urlencoded"
// ... 及更多图片格式常量
```

### 错误类型

```go
const (
    ErrUnsupportedContentType  errorx.ErrorType = 30000 + iota // 不支持的内容类型
    ErrInvalidMethod                                           // 无效的请求方法
    ErrBodyEncodeFuncNotSet                                    // 请求体编码函数未设置
    ErrExpectedDestinationType                                 // 期望的目标类型不匹配
    ErrRequestStatusCode                                       // 请求状态码错误
)
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
