---
name: mathx-conditional-logic
description: 数学工具与三元运算包，提供泛型条件表达式、空值/零值默认值、切片统计函数、克隆。当需要在表达式中做三元判断、处理空指针零值默认值、或计算百分比/均值/最值时使用。
---

# mathx - 数学工具与三元运算

提供泛型条件运算（IF/IfDo系列）、空值零值安全默认值、切片统计、概率与克隆。

> `mathx` 的 nil/零值判断依赖已调整为 `types`，例如 `IfNil`、`IfCEmpty` 通过 `types.IsNil` / `types.IsCEmpty` 判断。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/mathx"
```

三元表达式：
```go
result := mathx.IF(age >= 18, "adult", "minor")
```

安全默认值：
```go
val := mathx.IfNotEmpty(name, "default")
num := mathx.IfNotZero(count, 10)
```

## 完整API索引

### 函数

#### 条件运算

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `IF[T]` | `func(condition bool, trueVal, falseVal T) T` | 三元运算符泛型版本 |
| `IfDo[T]` | `func(condition bool, do func() T, defaultVal T) T` | 条件为true时执行函数返回结果 |
| `IfDoAF[T]` | `func(condition bool, do func() T, defaultVal ...T) T` | IfDo的废弃别名 |
| `IfDoWithError[T]` | `func(condition bool, do func() (T, error), defaultVal T) (T, error)` | 条件执行带错误返回 |
| `IfDoAsync[T]` | `func(condition bool, do func() T, defaultVal ...T) <-chan T` | 条件异步执行 |

#### 安全默认值

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `IfNotEmpty` | `func(str, defaultVal string) string` | 空字符串时返回默认值 |
| `IfNotZero[T]` | `func(val, defaultVal T) T` | 零值时返回默认值 |
| `IfLeZero[T]` | `func(val, defaultVal T) T` | 小于等于零时返回默认值 |
| `IfSafeIndex[T]` | `func(slice []T, index int, defaultVal T) T` | 安全索引访问，越界返回默认值 |
| `IfSafeKey[K,V]` | `func(m map[K]V, key K, defaultVal V) V` | 安全键访问，不存在返回默认值 |
| `DefaultIfNilPtr[T]` | `func(param, defaultValue T) T` | nil指针返回默认值 |
| `IfNil[T]` | `func(val interface{}, trueVal, falseVal T) T` | nil 时返回 trueVal |
| `IfNotNilValue[T]` | `func(val interface{}, trueVal, falseVal T) T` | 非 nil 时返回 trueVal |
| `IfCEmpty[T,R]` | `func(val T, trueVal, falseVal R) R` | 可比较类型零值时返回 trueVal |
| `IfNotCEmpty[T,R]` | `func(val T, trueVal, falseVal R) R` | 可比较类型非零值时返回 trueVal |

#### 统计函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Percentile` | `func(values []float64, p float64) float64` | 计算百分位数 |
| `Percentiles` | `func(values []float64, percentiles ...float64) []float64` | 计算多个百分位数 |
| `Percentage` | `func(part, total float64) float64` | 计算百分比 |
| `FormatPercentage` | `func(part, total float64, precision int) string` | 格式化百分比字符串 |
| `Mean` | `func(values []float64) float64` | 计算均值 |
| `MinSlice` | `func(values []float64) float64` | 计算最小值 |
| `MaxSlice` | `func(values []float64) float64` | 计算最大值 |
| `SliceMinMax[T]` | `func(list []T, f func(T) float64) (min, max float64)` | 泛型切片最值 |

#### 其他

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `SliceFisherYates[T]` | `func(slice []T, maxRetries int) []T` | Fisher-Yates洗牌 |
| `NewProba` | `func() *Proba` | 创建概率器 |
| `NewUnstable` | `func(deviation float64) *Unstable` | 创建不稳定随机器 |
| `Clone[T]` | `func(value T, seen ...interface{}) T` | 深克隆值 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `DoFunc[T]` | 条件执行函数类型 |
| `DoFuncWithError[T]` | 带错误的条件执行函数类型 |
| `Proba` | 概率器类型 |
| `Unstable` | 不稳定随机器类型 |

## 注意事项

- `IF` 即时求值两个分支，如需延迟求值用 `IfDo`
- `IfSafeIndex` 对越界索引返回默认值，不会panic
- `Percentile` 要求切片已排序或内部会先排序
- 通用 nil/零值/函数类型判断优先放在 `types`，`mathx` 只保留条件表达式语义