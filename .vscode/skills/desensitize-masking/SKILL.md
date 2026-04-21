---
name: desensitize-masking
description: 数据脱敏与掩码工具，提供多种脱敏类型（手机号/银行卡/IP/PEM等）、数据掩码器、规则注册。当需要对敏感数据进行脱敏处理、掩码显示、或自定义脱敏规则时使用。
---

# desensitize - 数据脱敏与掩码

提供多种内置脱敏类型、数据掩码器与自定义脱敏规则注册。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/desensitize"
```

基本脱敏：
```go
masked := desensitize.Desensitize("13812345678", desensitize.PhoneNumber)
```

敏感数据区间：
```go
result := desensitize.SensitiveData("hello world", 2, 5)
```

## 完整API索引

### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Desensitize` | `func(str string, dtype DesensitizeType, opts ...DesensitizeOptions) string` | 按类型脱敏 |
| `SensitiveData` | `func(str string, start, end int) string` | 按区间脱敏 |
| `SensitizePhoneNumber` | `func(str string) string` | 手机号脱敏 |
| `SensitizeBankCard` | `func(str string) string` | 银行卡号脱敏 |
| `SensitizeIpv4` | `func(str string) string` | IPv4地址脱敏 |
| `SensitizeIpv6` | `func(str string) string` | IPv6地址脱敏 |
| `SensitizePEMKey` | `func(str string) string` | PEM密钥脱敏 |
| `NewDesensitizeOptions` | `func() *DesensitizeOptions` | 创建脱敏选项 |
| `RegisterDesensitizer` | `func(dtype DesensitizeType, desensitizer Desensitizer)` | 注册自定义脱敏器 |
| `OperateByRule` | `func(dtype DesensitizeType, in string) string` | 按规则执行脱敏 |
| `Desensitization` | `func(obj interface{}) interface{}` | 对对象进行脱敏 |
| `NewMasker` | `func(configs ...MaskerConfig) *DataMasker` | 创建数据掩码器 |
| `DefaultMaskerConfig` | `func() MaskerConfig` | 获取默认掩码配置 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `DesensitizeType` | 脱敏类型枚举 |
| `DesensitizeOptions` | 脱敏选项类型 |
| `Desensitizer` | 脱敏器接口 |
| `DefaultDesensitizer` | 默认脱敏器类型 |
| `MaskerConfig` | 掩码配置类型 |
| `DataMasker` | 数据掩码器类型 |

## 注意事项

- `Desensitize` 按类型自动选择脱敏规则，也可通过 `RegisterDesensitizer` 注册自定义规则
- `SensitiveData` 的 start/end 为rune索引
- `Desensitization` 使用标签反射对结构体进行脱敏