---
name: metadata-extraction
description: 元数据提取工具，提供HTTP请求元数据提取、语言检测、TLS信息获取。当需要从HTTP请求提取结构化元数据、检测客户端语言、或获取TLS版本信息时使用。
---

# metadata - 元数据提取

提供HTTP请求元数据提取、语言检测与TLS信息获取。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/metadata"
```

提取请求元数据：
```go
md := metadata.ExtractRequestMetadata(r)
```

语言检测：
```go
lang := metadata.ExtractLanguage(r)
lang := metadata.ExtractLanguageWithDefault(r, "en")
```

## 完整API索引

### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ExtractRequestMetadata` | `func(r *http.Request) *RequestMetadata` | 提取HTTP请求元数据 |
| `NewMetadataExtractor` | `func(ctx context.Context, r *http.Request) *MetadataExtractor` | 创建元数据提取器 |
| `NewMetadataExtractorFromRequest` | `func(r *http.Request) *MetadataExtractor` | 从请求创建元数据提取器 |
| `FromContextSource` | `func() SourceFunc` | 上下文来源函数 |
| `FromQuerySource` | `func(key string) SourceFunc` | 查询参数来源函数 |
| `FromHeaderSource` | `func(key string) SourceFunc` | 请求头来源函数 |
| `FromCookieSource` | `func(key string) SourceFunc` | Cookie来源函数 |
| `FromAcceptLanguageSource` | `func() SourceFunc` | Accept-Language来源函数 |
| `FromMap` | `func(data map[string]string) SourceFunc` | 从map创建来源函数 |
| `NewLanguageExtractor` | `func(defaultLang string) *LanguageExtractor` | 创建语言提取器 |
| `ExtractLanguage` | `func(r *http.Request) string` | 提取客户端语言 |
| `ExtractLanguageWithDefault` | `func(r *http.Request, defaultLang string) string` | 提取语言带默认值 |
| `GetTLSVersionString` | `func(state *tls.ConnectionState) string` | 获取TLS版本字符串 |
| `ParseAcceptLanguage` | `func(header string) []string` | 解析Accept-Language头部 |
| `NormalizeLanguage` | `func(lang string) string` | 规范化语言代码 |
| `GetRemoteIP` | `func(r *http.Request) string` | 获取远程IP |
| `GetRemotePort` | `func(r *http.Request) string` | 获取远程端口 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `RequestMetadata` | 请求元数据类型 |
| `MetadataExtractor` | 元数据提取器类型 |
| `LanguageExtractor` | 语言提取器类型 |
| `SourceFunc` | 来源函数类型 |
| `ProcessFunc` | 处理函数类型 |
| `ContextKey` | 上下文键类型 |
| `MetadataAdapter` | 元数据适配器类型 |
| `Marshaler` | 序列化接口类型 |

## 注意事项

- `ExtractRequestMetadata` 会从请求中提取IP、User-Agent、语言等信息
- `ParseAcceptLanguage` 按 q 值权重排序
- `GetTLSVersionString` 在无TLS连接时返回空字符串