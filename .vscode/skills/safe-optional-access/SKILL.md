---
name: safe-optional-access
description: 安全访问与空值处理工具，提供可选链式字段访问、指针零值默认、时间哈希。当需要安全访问嵌套结构体字段、处理nil指针默认值、或需要时间桶哈希时使用。
---

# safe - 安全访问与空值处理

提供链式安全字段访问、nil指针零值默认值、时间桶哈希，避免空指针panic。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/safe"
```

链式安全访问：
```go
name := safe.Safe(user).Field("Profile").Field("Name").String("unknown")
age := safe.Safe(user).Field("Age").Int(0)
```

指针默认值：
```go
val := safe.Ptr[string, string](nilPtr, "default")
```

## 完整API索引

### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Safe` | `func(v interface{}) *SafeAccess` | 创建SafeAccess链式访问器 |
| `Ptr[T,R]` | `func(ptr *T, defaultValue R) R` | nil指针返回默认值 |
| `StringPtr` | `func(s string) *string` | 字符串指针工具 |
| `IntPtr` | `func(i int) *int` | 整数指针工具 |
| `BoolPtr` | `func(b bool) *bool` | 布尔指针工具 |
| `SlicePtr[T]` | `func(s []T) *[]T` | 切片指针工具 |
| `TimeToTimestampPB` | `func(t time.Time) int64` | 时间转PB时间戳 |
| `NewTemporalHasher` | `func(opts ...TemporalHasherOption) *TemporalHasher` | 创建时间桶哈希器 |
| `WithWindow` | `func(window time.Duration) TemporalHasherOption` | 设置时间窗口 |
| `WithLength` | `func(length int) TemporalHasherOption` | 设置哈希长度 |
| `WithSeparator` | `func(sep string) TemporalHasherOption` | 设置分隔符 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `SafeAccess` | 安全访问链式类型 |
| `TemporalHasher` | 时间桶哈希器类型 |
| `TemporalHasherOption` | 时间哈希器配置选项类型 |
| `SourceFunc` | 来源函数类型 |
| `ProcessFunc` | 处理函数类型 |
| `ContextKey` | 上下文键类型 |
| `MetadataAdapter` | 元数据适配器类型 |
| `Marshaler` | 序列化接口类型 |

### SafeAccess 方法

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Field` | `func(fieldName string) *SafeAccess` | 链式访问结构体字段 |
| `Bool` | `func(defaults ...bool) bool` | 安全取bool值，带默认值 |
| `Int` | `func(defaults ...int) int` | 安全取int值，带默认值 |
| `String` | `func(defaults ...string) string` | 安全取string值，带默认值 |
| `Float64` | `func(defaults ...float64) float64` | 安全取float64值，带默认值 |
| `Value` | `func() interface{}` | 安全取原始值（无默认值，返回零值） |

## 注意事项

- `Safe` 接收 `interface{}`，链式调用中遇到nil字段不会panic，而是返回零值
- `Field` 使用反射，性能敏感场景请避免在热路径中使用
- `Ptr` 的泛型参数 `T` 为指针基类型，`R` 为返回值类型