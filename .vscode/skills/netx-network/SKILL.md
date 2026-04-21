---
name: netx-network
description: 网络工具包，提供本地IP获取、公网IP获取、客户端IP提取。当需要获取本机/私有/公网IP地址、或从HTTP请求提取客户端IP时使用。
---

# netx - 网络工具

提供本地IP、私有IP、公网IP获取与HTTP客户端IP提取。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/netx"
```

获取IP：
```go
ips := netx.GetLocalInterfaceIPs()
privateIP := netx.GetPrivateIP()
publicIP := netx.GetConNetPublicIP()
```

从请求获取客户端IP：
```go
clientIP := netx.GetClientIP(r)
```

## 完整API索引

### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GetLocalInterfaceIPs` | `func() []string` | 获取所有本地网络接口IP |
| `GetPrivateIP` | `func() string` | 获取私有IP地址 |
| `GetLocalInterfaceIPAndExternalIP` | `func(urls ...string) (string, string)` | 获取本地IP和外部IP |
| `GetConNetPublicIP` | `func(urls ...string) string` | 获取公网IP地址 |
| `GetClientIP` | `func(r *http.Request) string` | 从HTTP请求提取客户端IP |

## 注意事项

- `GetConNetPublicIP` 需要外部网络访问，可传入备用URL
- `GetClientIP` 会依次检查 X-Forwarded-For、X-Real-IP 等头部