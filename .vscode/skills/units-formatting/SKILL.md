---
name: units-formatting
description: 数据单位格式化工具，提供十进制/二进制单位转换与人类可读格式化。当需要将字节数格式化为KB/MB/GB、或解析人类可读的单位字符串时使用。
---

# units - 数据单位格式化

提供十进制与二进制单位转换、人类可读格式化与解析。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/units"
```

格式化：
```go
s := units.HumanSize(1024 * 1024)      // "1 MiB"
s := units.BytesSize(1024 * 1024)       // "1MiB"
```

解析：
```go
b, err := units.ParseBytes("1GiB")      // 1073741824
b, err := units.ParseSizeDecimal("1GB") // 1000000000
```

## 完整API索引

### 常量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `KB` | int64 | 1000字节（十进制KB） |
| `MB` | int64 | 1000000字节（十进制MB） |
| `GB` | int64 | 1000000000字节（十进制GB） |
| `TB` | int64 | 1000000000000字节（十进制TB） |
| `PB` | int64 | 1000000000000000字节（十进制PB） |
| `KiB` | int64 | 1024字节（二进制KiB） |
| `MiB` | int64 | 1048576字节（二进制MiB） |
| `GiB` | int64 | 1073741824字节（二进制GiB） |
| `TiB` | int64 | 1099511627776字节（二进制TiB） |
| `PiB` | int64 | 1125899906842624字节（二进制PiB） |

### 变量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `DecimalMap` | map | 十进制单位映射 |
| `BinaryMap` | map | 二进制单位映射 |
| `DecimalAbbrs` | []string | 十进制单位缩写 |
| `BinaryAbbrs` | []string | 二进制单位缩写 |

### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `CustomSize` | `func(format string, size float64, base float64, unitAbbrs []string) string` | 自定义单位格式化 |
| `HumanSizeWithPrecision` | `func(size float64, precision int) string` | 人类可读大小（指定精度） |
| `HumanSize` | `func(size float64) string` | 人类可读大小 |
| `BytesSize` | `func(size float64) string` | 字节数格式化 |
| `FormatBytes` | `func(bytes int64) string` | 格式化字节数 |
| `ParseBytes` | `func(size string) (int64, error)` | 解析字节字符串 |
| `ParseSizeDecimal` | `func(size string) (int64, error)` | 解析十进制大小字符串 |
| `ParseSizeBinary` | `func(size string) (int64, error)` | 解析二进制大小字符串 |

## 注意事项

- 十进制单位（KB/MB/GB）基于1000，二进制单位（KiB/MiB/GiB）基于1024
- `ParseBytes` 自动识别十进制/二进制单位后缀
- `HumanSize` 使用二进制单位（MiB/GiB等）