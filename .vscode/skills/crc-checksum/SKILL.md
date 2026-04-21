---
name: crc-checksum
description: CRC校验计算工具，提供可配置的CRC算法工厂、缓存工厂、多种CRC参数配置。当需要计算CRC校验和、使用标准或自定义CRC算法、或创建CRC计算器工厂时使用。
---

# crc - CRC校验计算

提供可配置的CRC算法、计算器接口、工厂与缓存工厂。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/crc"
```

创建CRC计算器：
```go
calc := crc.New(crc.Config{
    Width:  32,
    Poly:   0x04C11DB7,
    Init:   0xFFFFFFFF,
    XorOut: 0xFFFFFFFF,
    RefIn:  true,
    RefOut: true,
})
result := calc.Compute(data)
```

## 完整API索引

### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `New` | `func(cfg Config) Calculator` | 创建CRC计算器 |
| `NewFactory` | `func(cfg Config) Factory` | 创建CRC计算器工厂 |
| `NewCachedFactory` | `func(cfg Config) Factory` | 创建带缓存的CRC计算器工厂 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `Config` | CRC配置类型，含 Width, Poly, Init, XorOut, RefIn, RefOut 字段 |
| `Calculator` | CRC计算器接口，含 `Compute(data []byte) uint64` 和 `Reset()` 方法 |
| `Factory` | CRC计算器工厂接口 |

### Config 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| `Width` | uint8 | CRC位宽（8/16/32/64） |
| `Poly` | uint64 | 生成多项式 |
| `Init` | uint64 | 初始值 |
| `XorOut` | uint64 | 输出异或值 |
| `RefIn` | bool | 输入反转 |
| `RefOut` | bool | 输出反转 |

## 注意事项

- `NewFactory` 每次创建新的计算器，适合单次计算
- `NewCachedFactory` 复用计算器，适合高频重复计算场景
- `Config` 的 `Width` 支持 8/16/32/64 位CRC