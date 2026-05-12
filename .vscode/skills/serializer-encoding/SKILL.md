---
name: serializer-encoding
description: 序列化编码工具，提供多格式（JSON/Gob/Msgpack/Protobuf）与多压缩（Gzip/Zlib/Zstd）泛型序列化、proto-aware JSON、protobuf JSON、结构化 JSON 错误与性能基准。当需要对结构体/包含 proto.Message 的对象进行序列化/反序列化、选择压缩算法、或对比 JSON/protojson 性能时使用。
---

# serializer - 序列化编码

提供多格式多压缩的泛型序列化/反序列化、proto-aware JSON 编解码与性能基准测试。

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
jsonStr := serializer.ToJSON(obj)
obj := serializer.FromJSON[MyStruct](jsonStr)
```

包含 protobuf 字段的结构体：

```go
type Payload struct {
	Name *wrapperspb.StringValue `json:"name"`
	Age  *wrapperspb.Int32Value  `json:"age"`
}

data, err := serializer.JSONMarshal(&Payload{
	Name: wrapperspb.String("alice"),
	Age:  wrapperspb.Int32(18),
})
// 输出: {"name":"alice","age":18}
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
| `JSONMarshal[T]` | `func(value T) ([]byte, error)` | 标准 JSON / proto-aware JSON 序列化 |
| `JSONUnmarshal[T]` | `func(data []byte, target *T) error` | 标准 JSON / proto-aware JSON 反序列化 |
| `ToJSON[T]` | `func(v T) string` | 快捷JSON序列化，失败返回空字符串 |
| `FromJSON[T]` | `func(jsonStr string) T` | 快捷JSON反序列化，失败返回零值 |
| `ProtoJSONMarshal` | `func(m interface{}) (string, error)` | 序列化 protobuf 消息为 JSON 字符串 |
| `ProtoJSONUnmarshal` | `func(a, b interface{}) error` | 自动识别 proto.Message 与 JSON 字符串/[]byte，兼容任意参数顺序 |

#### JSON错误函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewJSONNilTargetError` | `func() error` | 创建 JSON 目标为空错误 |
| `IsJSONNilTargetError` | `func(err error) bool` | 判断 JSON 目标为空错误 |
| `NewJSONExpectedObjectError` | `func() error` | 创建期望 JSON 对象错误 |
| `IsJSONExpectedObjectError` | `func(err error) bool` | 判断期望 JSON 对象错误 |
| `NewJSONExpectedArrayError` | `func() error` | 创建期望 JSON 数组错误 |
| `IsJSONExpectedArrayError` | `func(err error) bool` | 判断期望 JSON 数组错误 |
| `NewJSONFieldError` | `func(name string, err error) error` | 包装字段级 JSON 错误 |
| `NewJSONItemError` | `func(index int, err error) error` | 包装数组/切片元素级 JSON 错误 |
| `NewJSONKeyError` | `func(key string, err error) error` | 包装 map 键级 JSON 错误 |

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

## proto-aware JSON 说明

- `JSONMarshal` / `JSONUnmarshal` 会自动检测结构体、切片、数组、map 中是否包含 `proto.Message`。
- 普通 Go 字段走标准 `encoding/json`，protobuf 字段走 `protojson`，因此 wrapper、Timestamp、Duration、FieldMask、Any、Struct、DescriptorProto 都使用 protobuf 官方 JSON 形态。
- 类型判断和字段元信息使用缓存；对象和数组使用字节扫描快路径，减少 `map[string]json.RawMessage` / `[]json.RawMessage` 中间分配。
- 反序列化时 `omitempty` / `omitzero` 只影响 marshal，不会跳过目标字段写入。

## 性能基准

运行 serializer JSON 对照 benchmark：

```bash
go test ./pkg/serializer -run ^$ -bench "BenchmarkJSON(Marshal|Unmarshal)(ProtoStruct|GeneratedProtoPayload|GeneratedProtoPayloadTraditional)$" -benchmem
```

当前真实 generated protobuf 场景包含 `Any`、`Struct`、`Timestamp`、`Duration`、`FieldMask`、`DescriptorProto`、map 和 repeated 字段，并提供传统 `RawMessage + protojson` 对照基准。

## 注意事项

- `NewGob` 要求结构体字段全部导出，否则反序列化会丢失字段
- `NewFast` 和 `NewUltraCompact` 使用msgpack编码，需确保字段类型兼容
- 压缩序列化器在小数据量时可能比非压缩更大
- 包含 protobuf 字段的 JSON 场景优先使用 `JSONMarshal` / `JSONUnmarshal`，不要直接用 `encoding/json`
