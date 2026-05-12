---
name: convert-type-transforms
description: 类型转换工具包，提供泛型强制转换、JSON/YAML互转、字节与十六进制转换、字段映射转换。当需要将任意类型安全转换为string/int/float/bool、或进行JSON与YAML/Hex/BCC互转、对象字段映射时使用。
---

# convert - 类型转换工具包

提供泛型类型强制转换、JSON/YAML编解码、字节与十六进制互转、字段映射转换等类型变换工具。

> 快速数字/时间格式化能力已迁移到 `stringx`，如 `FastAppendInt`、`FastFormatTime`、`FastItoa`。`convert.AppendValue` 内部会复用 `stringx` 的快速格式化能力。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/convert"
```

泛型类型转换：
```go
s := convert.MustString[int](42)
n := convert.MustIntT[string]("123")
b := convert.MustBool[string]("true")
```

JSON/YAML互转：
```go
jsonBytes := convert.YAMLToJSON(yamlBytes)
yamlBytes := convert.JSONToYAML(jsonBytes)
```

## 完整API索引

### 函数

#### 泛型类型转换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `MustString[T]` | `func(v T, timeLayout...string) string` | 将任意基本类型转为string |
| `MustIntT[T]` | `func(value, mode) T` | 将值转为整数类型T |
| `MustFloatT[T]` | `func(value, mode) T` | 将值转为浮点类型T |
| `ToFloat64` | `func(value) float64` | 将值转为float64 |
| `Float64ToInt[T]` | `func(value, mode) T` | float64转整数类型T |
| `ParseFloat[T]` | `func(v, value*) T` | 解析浮点数 |
| `MustBool[T]` | `func(v T) bool` | 将值转为bool |
| `MustConvertTo[T]` | `func(value) T` | 泛型强制类型转换 |

#### JSON转换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `MustJSONIndent` | `func(v) string` | 将值序列化为缩进JSON字符串 |
| `MustJSON` | `func(v) string` | 将值序列化为JSON字符串 |

#### 切片转换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `StringsToJSON` | `func(s []string) string` | 字符串切片序列化为JSON数组字符串 |
| `StringsFromJSON` | `func(jsonStr string) ([]string, error)` | JSON数组字符串反序列化为字符串切片 |
| `NumberSliceToStringSlice[T]` | `func(numbers []T) []string` | 数字切片转字符串切片 |
| `StringSliceToNumberSlice[T]` | `func(input []string, mode) []T` | 字符串切片转数字切片 |
| `StringSliceToFloatSlice[T]` | `func(input []string, mode) []T` | 字符串切片转浮点切片 |
| `AnySliceToInterfaceSlice` | `func(slice) []interface{}` | 任意切片转interface切片 |
| `StringSliceToInterfaceSlice` | `func(slice) []interface{}` | 字符串切片转interface切片 |
| `InterfaceSliceToStringSlice` | `func(slice) []string` | interface切片转字符串切片 |
| `InterfaceSliceToIntSlice` | `func(slice, mode) []int` | interface切片转整数切片 |
| `ToNumberSlice[T]` | `func(slice) []T` | 切片转数字切片 |
| `MustToNumberSlice[T]` | `func(slice) []T` | 切片强制转数字切片 |

#### Map转换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `InterfaceMapToStringMap` | `func(m) map[string]string` | interface map转string map |
| `ParseObjectToMap` | `func(obj) map[string]interface{}` | 将对象解析为map |
| `ParseKVPairsToMap` | `func(keysAndValues...) map[string]string` | 将键值对解析为map |

#### 辅助函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `AppendValue` | `func(buf, v)` | 将值追加到buffer |

> 新代码中如需直接快速格式化数字/时间，请使用 `stringx.FastAppendInt`、`stringx.FastFormatTime` 等。

#### 字节/十六进制/二进制/十进制转换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `BytesToHex` | `func(b []byte) string` | 字节切片转十六进制字符串 |
| `HexToBytes` | `func(hex string) []byte` | 十六进制字符串转字节切片 |
| `ByteToBinStr` | `func(b byte) string` | 单字节转二进制字符串 |
| `BytesToBinStr` | `func(b []byte) string` | 字节切片转二进制字符串 |
| `BytesToBinStrWithSplit` | `func(b []byte, split string) string` | 字节切片转二进制字符串（带分隔符） |
| `HexToDec` | `func(hex string) int64` | 十六进制转十进制 |
| `DecToHex` | `func(n int64) string` | 十进制转十六进制 |
| `DecToBin` | `func(n int64) string` | 十进制转二进制 |
| `HexToBin` | `func(hex string) string` | 十六进制转二进制 |
| `HexToBCC` | `func(hex string) string` | 十六进制转BCC码 |
| `BytesToBCC` | `func(b []byte) byte` | 字节切片转BCC校验码 |

#### YAML/JSON编解码

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `YAMLToJSON` | `func(yamlBytes) []byte` | YAML字节转JSON字节 |
| `JSONToYAML` | `func(jsonBytes) []byte` | JSON字节转YAML字节 |
| `YAMLStringToJSON` | `func(yamlStr) string` | YAML字符串转JSON字符串 |
| `JSONStringToYAML` | `func(jsonStr) string` | JSON字符串转YAML字符串 |
| `YAMLToInterface` | `func(yamlBytes) interface{}` | YAML字节转interface |
| `YAMLToMap` | `func(yamlBytes) map[string]interface{}` | YAML字节转map |
| `InterfaceToYAML` | `func(v) []byte` | interface转YAML字节 |
| `MapToYAML` | `func(m) []byte` | map转YAML字节 |
| `UnmarshalYAML[T]` | `func(data) (T, error)` | YAML泛型反序列化 |
| `MarshalYAML[T]` | `func(v T) ([]byte, error)` | YAML泛型序列化 |
| `UnmarshalJSON[T]` | `func(data) (T, error)` | JSON泛型反序列化 |
| `MarshalJSON[T]` | `func(v T) ([]byte, error)` | JSON泛型序列化 |

#### 字段映射转换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewTransformer` | `func() *Transformer` | 创建字段转换器 |
| `TransformFields` | `func(dst, src, opts) error` | 字段映射转换 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `RoundMode` | 舍入模式类型 |
| `ConvertError` | 转换错误类型 |
| `TransformFieldsOptions` | 字段映射转换选项 |
| `Transformer` | 字段转换器类型 |

### 常量/变量

（无独立常量/变量导出）

## 常用示例

详细用法参阅 → [reference.md](reference.md)

## 注意事项

- `Must*` 函数在转换失败时panic，如需安全转换请先验证输入
- `YAMLToJSON` 返回 `[]byte`，需用 `string()` 包装才能得到字符串
- `MustConvertTo[T]` 依赖 `cast` 库，仅支持常见基本类型互转
- `TransformFields` 按字段名匹配，可用 `TransformFieldsOptions` 自定义映射
- 字段映射中的严格类型兼容检查已下沉到 `types.CheckTypeCompatibility`