---
name: idgen-identifier
description: ID生成器，提供多种唯一ID生成策略（UUID/NanoID/Snowflake/ShortFlake/ULID）。当需要生成分布式唯一ID、短ID、或有序ID时使用。
---

# idgen - ID生成器

提供多种唯一ID生成策略，支持UUID、NanoID、Snowflake、ShortFlake和ULID。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/idgen"
```

基本使用：
```go
gen := idgen.NewDefaultIDGenerator()
id := gen.Generate() // 生成ID
```

指定类型：
```go
uuid := idgen.NewUUIDGenerator().Generate()
snow := idgen.NewSnowflakeGenerator(1, 1).Generate()
```

## 完整API索引

### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewIDGenerator` | `func(generatorType GeneratorType) IDGenerator` | 按类型创建ID生成器 |
| `NewIDGeneratorFromString` | `func(type string) (IDGenerator, error)` | 按字符串创建ID生成器（Deprecated） |
| `NewDefaultIDGenerator` | `func() *DefaultIDGenerator` | 创建默认ID生成器（UUID） |
| `NewUUIDGenerator` | `func() *UUIDGenerator` | 创建UUID生成器 |
| `NewNanoIDGenerator` | `func() *NanoIDGenerator` | 创建NanoID生成器 |
| `NewSnowflakeGenerator` | `func(workerID, datacenter int64) *SnowflakeGenerator` | 创建Snowflake生成器 |
| `NewShortFlakeGenerator` | `func(nodeID int64) *ShortFlakeGenerator` | 创建ShortFlake生成器 |
| `NewShortFlakeBase62Generator` | `func(nodeID int64) *ShortFlakeBase62Generator` | 创建Base62短ID生成器 |
| `NewULIDGenerator` | `func() *ULIDGenerator` | 创建ULID生成器 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `IDGenerator` | ID生成器接口，含 `Generate()` 方法 |
| `GeneratorType` | 生成器类型枚举 |
| `DefaultIDGenerator` | 默认ID生成器类型 |
| `UUIDGenerator` | UUID生成器类型 |
| `NanoIDGenerator` | NanoID生成器类型 |
| `SnowflakeGenerator` | Snowflake生成器类型 |
| `ShortFlakeGenerator` | ShortFlake生成器类型 |
| `ShortFlakeBase62Generator` | Base62短ID生成器类型 |
| `ULIDGenerator` | ULID生成器类型 |

## 注意事项

- `NewIDGeneratorFromString` 已废弃，建议使用 `NewIDGenerator`
- Snowflake需要配置唯一的 workerID 和 datacenterID
- ULID基于时间戳，适合需要时间排序的场景