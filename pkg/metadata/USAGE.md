# Metadata 使用指南

HTTP 请求元数据提取器，用于提取和管理 HTTP 请求的完整元数据信息。

## 快速开始

### 基础用法

```go
import (
    "net/http"
    "github.com/kamalyes/go-toolbox/pkg/metadata"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // 提取请求元数据
    meta := metadata.ExtractRequestMetadata(r)
    
    // 访问元数据
    fmt.Println("浏览器:", meta.Browser)
    fmt.Println("操作系统:", meta.OS)
    fmt.Println("设备类型:", meta.DeviceType)
    fmt.Println("客户端IP:", meta.ClientIP)
}
```

## 主要功能

### 1. 提取请求元数据

```go
meta := metadata.ExtractRequestMetadata(r)

// User-Agent 信息
meta.Browser           // Chrome, Firefox, Safari 等
meta.BrowserVersion    // 浏览器版本
meta.OS                // Windows, macOS, iOS, Android 等
meta.OSVersion         // 操作系统版本
meta.DeviceType        // mobile, tablet, desktop, bot
meta.DeviceVendor      // Apple, Samsung, Huawei 等
meta.IsBot             // 是否为爬虫
meta.IsMobile          // 是否为移动设备
meta.IsTablet          // 是否为平板

// 请求基础信息
meta.ClientIP          // 客户端 IP
meta.RequestMethod     // GET, POST 等
meta.RequestURI        // 请求路径
meta.RequestHost       // 请求主机

// 代理和转发信息
meta.XForwardedFor     // X-Forwarded-For
meta.XRealIP           // X-Real-IP

// TLS 信息
meta.Protocol          // http 或 https
meta.TLSVersion        // TLS 版本
meta.TLSServerName     // TLS 服务器名称
```

### 2. 转换为 Map

```go
// 转换为 map[string]interface{}
dataMap := meta.ToMap()

// 序列化为 JSON
jsonData, _ := json.Marshal(dataMap)
```

### 3. 从 Map 恢复

```go
// 从 map 恢复元数据
meta := metadata.FromMap(dataMap)
```

### 4. 访问器方法

```go
// 获取头信息
userAgent := meta.GetHeader("User-Agent")
origin := meta.GetHeader("Origin")

// 设置头信息
meta.SetHeader("Custom-Header", "value")
```

### 5. 工具函数

```go
// TLS 版本转字符串
tlsStr := metadata.GetTLSVersionString(meta.TLSVersion)  // "TLS 1.3"

// 解析 Accept-Language
lang, region, full := metadata.ParseAcceptLanguage(meta.AcceptLanguage)
// "zh", "CN", "zh-CN"

// 提取 IP 和端口
ip := metadata.GetRemoteIP(meta.RemoteAddr)       // "192.168.1.1"
port := metadata.GetRemotePort(meta.RemoteAddr)   // "12345"
```

## 应用场景

### 请求日志记录

```go
func LogRequest(r *http.Request) {
    meta := metadata.ExtractRequestMetadata(r)
    log.Printf("请求: %s %s | 浏览器: %s | 设备: %s | IP: %s",
        meta.RequestMethod, meta.RequestURI,
        meta.Browser, meta.DeviceType, meta.ClientIP)
}
```

### 安全审计

```go
func AuditRequest(r *http.Request) {
    meta := metadata.ExtractRequestMetadata(r)
    
    // 检测可疑请求
    if meta.IsBot {
        log.Printf("爬虫访问: %s from %s", meta.BotName, meta.ClientIP)
    }
    
    // 记录地理位置
    if meta.CFIPCountry != "" {
        log.Printf("来自国家: %s", meta.CFIPCountry)
    }
}
```

### 设备适配

```go
func AdaptiveResponse(r *http.Request) string {
    meta := metadata.ExtractRequestMetadata(r)
    
    if meta.IsMobile {
        return "mobile.html"
    } else if meta.IsTablet {
        return "tablet.html"
    }
    return "desktop.html"
}
```

## 测试覆盖率

- **覆盖率**: 97.9%
- **测试用例**: 完整覆盖各种设备、浏览器、爬虫场景
