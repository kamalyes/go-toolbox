---
name: location-geo
description: 通过高德开放平台 IP 定位接口把 IP 解析为 GPS 坐标，内网地址会匹配到当前出口公网地址；当需要由 IP 反查地理位置时使用
---

# location - IP 转 GPS 坐标

只有一个函数，通过高德（Amap）开放平台的 IP 定位接口把 IP 解析为坐标

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/location"
```

```go
data, err := location.GetGPSByIpAmap(amapKey, amapSign, amapUrl, ip)
if err != nil {
    return err
}
// data 是高德返回的原始 JSON 反序列化结果
```

## 完整API索引

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GetGPSByIpAmap` | `func(amapKey, amapSign, amapUrl, ip string) (map[string]interface{}, error)` | 查询 IP 对应的 GPS 坐标；传入内网地址时，高德会按请求的出口公网地址返回 |

| 常量 | 值 | 说明 |
|---|---|---|
| `TIME_OUT` | `3 * time.Second` | 第三方接口请求超时时间 |

## 注意事项

- **这是一个会发起外部网络请求的包**，不是纯计算。它依赖高德服务可用、
  依赖有效的 `amapKey`，并且会受配额与限流影响。放在热路径上要自己加缓存
- **返回值是 `map[string]interface{}`，不是结构体**。字段名与层级由高德接口决定，
  取值前要判断键是否存在、类型是否符合预期，不要直接断言
  上游改字段不会导致编译失败，只会在运行时取到零值
- 超时固定为 `TIME_OUT`（3 秒），**不接受 `context.Context`**，
  因此无法随调用方取消，也无法按调用点调整超时。在需要可取消的链路里要特别注意
- `amapSign` 用于高德的数字签名校验；只有在控制台开启了签名校验时才需要，
  未开启时传空串
- 测试会真打外部接口，**没有 key 或网络不可达时相关测试无法通过**，
  CI 里应跳过或自备 mock
