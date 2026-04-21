---
name: serializer-encoding
description: 序列化编码工具，提供多格式（JSON/Gob/Msgpack/Protobuf）与多压缩（Gzip/Zlib/Zstd）泛型序列化。当需要对结构体进行序列化/反序列化、选择压缩算法、或基准测试序列化性能时使用。
---

# serializer - 序列化编码

提供多格式多压缩的泛型序列化/反序列化与性能基准测试。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/serializer"
```

基本序列化：
```go
s := serializer.New[MyStruct]()
data, err := s.Serialize(obj)
obj, err := s.Deserialize(data)
```

JSON快捷方式：
```go
data, err := serializer.ToJSON[MyStruct](obj)
obj, err := serializer.FromJSON[MyStruct](data)
```

## 完整API索引

### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `New[T]` | `func() *Serializer[T]` | 创建默认序列化器 |
| `NewJSON[T]` | `func() *Serializer[T]` | 创建JSON序列化器 |
| `NewGob[T]` | `func() *Serializer[T]` | 创建Gob序列化器 |
| `NewCompact[T]` | `func() *Serializer[T]` | 创建紧凑序列化器 |
| `NewZlibCompact[T]` | `func() *Serializer[T]` | 创建Zlib紧凑序列化器 |
| `NewFast[T]` | `func() *Serializer[T]` | 创建快速序列化器 |
| `NewUltraCompact[T]` | `func() *Serializer[T]` | 创建超紧凑序列化器 |
| `ToJSON[T]` | `func(v T) ([]byte, error)` | 快捷JSON序列化 |
| `FromJSON[T]` | `func(data []byte) (T, error)` | 快捷JSON反序列化 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `Serializer[T]` | 泛型序列化器类型 |
| `SerializeType` | 序列化格式枚举 |
| `CompressionType` | 压缩类型枚举 |
| `Stats` | 序列化统计类型 |
| `BenchmarkResult` | 基准测试结果类型 |

### 常量/变量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `TypeJSON` | SerializeType | JSON序列化格式 |
| `TypeGob` | SerializeType | Gob序列化格式 |
| `TypeMsgpack` | SerializeType | Msgpack序列化格式 |
| `TypeProtobuf` | SerializeType | Protobuf序列化格式 |
| `CompressionNone` | CompressionType | 无压缩 |
| `CompressionGzip` | CompressionType | Gzip压缩 |
| `CompressionZlib` | CompressionType | Zlib压缩 |
| `CompressionZstd` | CompressionType | Zstd压缩 |

## 注意事项

- `NewGob` 要求结构体字段全部导出，否则反序列化会丢失字段
- `NewFast` 和 `NewUltraCompact` 使用msgpack编码，需确保字段类型兼容
- 压缩序列化器在小数据量时可能比非压缩更大