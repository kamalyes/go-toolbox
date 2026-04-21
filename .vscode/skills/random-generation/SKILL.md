---
name: random-generation
description: 随机生成工具，提供随机数/字符串/字节/时长、数值范围、邮箱/手机/姓名/身份证等假数据、域名关键词构建、姓氏管理。当需要生成随机测试数据、随机ID片段、或构造仿真用户信息时使用。
---

# random - 随机生成

提供随机基本类型生成与假数据生成（邮箱/手机/姓名/身份证），用于测试与仿真。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/random"
```

基本随机：
```go
n := random.RandInt(1, 100)
s := random.RandString(16)
h := random.RandHex(8)
```

假数据：
```go
email := random.RandomEmail()
phone := random.RandomPhone()
```

## 完整API索引

### 函数

#### 基本随机

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewRand` | `func(seed ...int64) *rand.Rand` | 创建随机源 |
| `RandInt` | `func(min, max int) int` | 随机整数[min,max] |
| `RandString` | `func(length int) string` | 随机字母数字字符串 |
| `RandHex` | `func(length int) string` | 随机十六进制字符串 |
| `RandNumber` | `func(length int) string` | 随机数字字符串 |
| `RandBytes` | `func(length int) []byte` | 随机字节切片 |

#### 快速随机（非crypto安全）

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `FRandInt` | `func(min, max int) int` | 快速随机整数 |
| `FRandUint32` | `func() uint32` | 快速随机uint32 |
| `FastRand64` | `func() uint64` | 快速随机uint64 |
| `FastRand` | `func() int` | 快速随机int |
| `FastRandn` | `func(n int) int` | 快速随机[0,n) |
| `FRandBytesJSON[T]` | `func(n int) T` | 快速随机字节的JSON反序列化 |

#### 字符串切片

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `RandStringSlice` | `func(length, count int) []string` | 随机字符串切片 |

#### 注册与模型生成

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `RegisterGenerator` | `func(name string, fn RandGeneratorFunc)` | 注册自定义随机生成器 |
| `GetGenerator` | `func(name string) (RandGeneratorFunc, bool)` | 获取已注册生成器 |
| `UnregisterGenerator` | `func(name string)` | 注销生成器 |
| `ListRegisteredGenerators` | `func() []string` | 列出已注册生成器 |
| `ClearAllGenerators` | `func()` | 清除所有生成器 |
| `GenerateRandModel[T]` | `func() T` | 使用注册生成器生成随机模型 |

#### 假数据

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `RandomEmail` | `func() string` | 随机邮箱 |
| `RandomPhone` | `func() string` | 随机手机号 |
| `RandomName` | `func() string` | 随机姓名 |
| `RandomIDCard` | `func() string` | 随机身份证号 |

#### 域名与姓氏

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewDomainKeywordBuilder` | `func(baseKeyword string) *DomainKeywordBuilder` | 创建域名关键词构建器 |
| `NewSurnameManager` | `func(data ...string) *SurnameManager` | 创建姓氏管理器 |

#### 时间与数值

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `RandDuration` | `func(min, max time.Duration) time.Duration` | 随机时间段 |
| `RandTimeBetween` | `func(start, end time.Time) time.Time` | 随机时间（两时间之间） |
| `GenerateAvailablePort` | `func() (int, error)` | 生成可用端口 |
| `RandNumericalStep[T]` | `func(start, end, step T) T` | 随机数（起止+步长） |
| `RandNumerical[T]` | `func(min, max T) T` | 随机浮点数（范围内） |
| `UUID` | `func() string` | 生成UUID |

### 类型

| 导出名称 | 说明 |
|---|---|
| `RandType` | 随机类型常量（CAPITAL/LOWERCASE/SPECIAL/NUMBER） |
| `RandGeneratorFunc` | 随机生成器函数类型 |
| `RandGeneratorRegistry` | 随机生成器注册表类型 |
| `DomainKeywordBuilder` | 域名关键词构建器类型 |
| `SurnameInfo` | 姓氏信息类型 |
| `SurnameManager` | 姓氏管理器类型 |

### 常量/变量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `DEC_BYTES` | string | 数字字符集 |
| `HEX_BYTES` | string | 十六进制字符集 |
| `ALPHA_BYTES` | string | 字母字符集 |
| `LETTER_BYTES` | string | 字母数字字符集 |
| `CAPITAL` | RandType | 大写字母类型 |
| `LOWERCASE` | RandType | 小写字母类型 |
| `SPECIAL` | RandType | 特殊字符类型 |
| `NUMBER` | RandType | 数字类型 |

## 注意事项

- `FastRand` 系列性能优于 `RandString` 但非crypto安全
- `RandomPhone` / `RandomIDCard` 仅用于测试，不保证校验位合法
- `RandNumericalStep` 的步长参数必须大于0