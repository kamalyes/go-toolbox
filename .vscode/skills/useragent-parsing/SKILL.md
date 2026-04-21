---
name: useragent-parsing
description: User-Agent解析工具，提供浏览器/OS/设备信息提取、Bot检测、移动端判断。当需要解析HTTP User-Agent字符串、判断客户端类型、或检测爬虫时使用。
---

# useragent - User-Agent解析

提供User-Agent字符串解析，提取浏览器、操作系统、设备信息与Bot检测。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/useragent"
```

解析User-Agent：
```go
result := useragent.Parse(r.Header.Get("User-Agent"))
fmt.Println(result.Browser, result.OS, result.DeviceType)
```

## 完整API索引

### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Parse` | `func(ua string) *ParsedUserAgent` | 解析User-Agent字符串 |
| `ParseVersion` | `func(ver string) VersionNo` | 解析版本号 |
| `New` | `func() *UserAgent` | 创建UserAgent解析器 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `ParsedUserAgent` | 解析结果类型，包含以下字段：Raw, Browser, BrowserVersion, OS, OSVersion, Device, DeviceType, DeviceVendor, IsBot, BotName, IsMobile, IsTablet 等 |
| `VersionNo` | 版本号类型，包含 Major, Minor, Patch 字段 |
| `RgType` | 正则类型 |
| `DeviceType` | 设备类型枚举 |
| `UserAgent` | UserAgent解析器类型 |

### 常量/变量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `AllBrowsers` | []string | 所有已知浏览器列表 |
| `PopularBrowsers` | []string | 常见浏览器列表 |
| `UnpopularBrowsers` | []string | 不常见浏览器列表 |
| `AllOS` | []string | 所有已知操作系统列表 |
| `PopularOS` | []string | 常见操作系统列表 |
| `UnpopularOS` | []string | 不常见操作系统列表 |
| `StabilizeUserAgents` | []string | 稳定化UA列表 |

(`AllBrowsers`/`PopularBrowsers`/`UnpopularBrowsers` 包含 Chrome、Firefox、Safari、Edge 等;
`AllOS`/`PopularOS`/`UnpopularOS` 包含 Windows、macOS、Linux、Android、iOS 等;
另有各种浏览器/OS/设备名称常量)

## 注意事项

- `Parse` 对未知UA返回空字段而非nil
- `IsBot` 通过预定义Bot列表判断，自定义Bot需扩展
- `ParseVersion` 将 "1.2.3" 解析为 `{Major:1, Minor:2, Patch:3}`