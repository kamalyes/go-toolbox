---
name: types-reflection
description: 类型与反射工具，提供 proto.Message 类型判断、JSON tag 解析、导出字段判断、nil/零值/函数类型判断、指针创建、类型兼容性检查、数值 Kind 判断。当需要下沉通用类型判断、处理反射字段、判断 protobuf 类型、或做严格类型兼容检查时使用。
---

# types - 类型与反射工具

提供通用类型判断、反射辅助、JSON tag 解析、protobuf 类型判断和类型兼容性检查。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/types"
```

JSON 字段解析：

```go
field, _ := reflect.TypeOf(User{}).FieldByName("Name")
name := types.JSONFieldName(field)
omit := types.HasJSONTagOption(field, "omitempty", "omitzero")
```

protobuf 类型判断：

```go
t := reflect.TypeOf(wrapperspb.String("x"))
ok := types.IsProtoMessageType(t)
```

## 完整API索引

### protobuf 与结构体字段

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ProtoMessageType` | `reflect.Type` | protobuf 消息接口类型 |
| `IsProtoMessageType` | `func(t reflect.Type) bool` | 判断类型是否实现 `proto.Message` |
| `IsExportedField` | `func(field reflect.StructField) bool` | 判断字段是否可导出处理，匿名字段也允许 |
| `ExtractJSONKey` | `func(field reflect.StructField) string` | 从 json tag 中提取字段名 |
| `JSONFieldName` | `func(field reflect.StructField) string` | 获取 JSON 字段名，无显式名称时返回 Go 字段名 |
| `HasJSONTagOption` | `func(field reflect.StructField, options ...string) bool` | 判断 json tag 是否包含指定选项 |
| `EnsureStructDefaults` | `func(v reflect.Value)` | 初始化结构体中的 protobuf/结构体指针字段 |
| `NewProtoMessage[T]` | `func() T` | 创建新的 protobuf 消息实例 |

### nil、零值和指针

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `IsNil` | `func(x interface{}) bool` | 判断 interface 或其底层值是否为 nil |
| `IsCEmpty[T]` | `func(v T) bool` | 判断可比较类型是否为零值 |
| `IsFuncType[T]` | `func() bool` | 判断泛型类型是否为函数 |
| `DerefValue` | `func(value interface{}) (interface{}, bool)` | 解引用指针，返回底层值 |
| `Ptr[T]` | `func(v T) *T` | 创建任意类型指针 |

### Kind 与数值类型

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GetReflectKind` | `func(value interface{}) reflect.Kind` | 获取反射 Kind，nil 返回 Invalid |
| `IsNumericKind` | `func(kind reflect.Kind) bool` | 判断是否为数值 Kind |
| `IsIntegerKind` | `func(kind reflect.Kind) bool` | 判断是否为整数 Kind |
| `IsFloatKind` | `func(kind reflect.Kind) bool` | 判断是否为浮点 Kind |
| `ToFloat64OK` | `func(value interface{}) (float64, bool)` | 尝试将数值转为 float64 |
| `IsWholeNumber` | `func(f float64) bool` | 判断浮点数是否为整数值 |

### 类型兼容性

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `CheckTypeCompatibility` | `func(srcType, dstType reflect.Type) error` | 严格检查源类型是否可赋给目标类型 |

## 使用建议

- 通用类型判断、反射 Kind 判断和 JSON tag 解析优先放在 `types`。
- `validator` 只保留校验语义；`mathx` 只保留条件表达式语义；`convert` 只保留类型转换语义。
- 处理包含 protobuf 字段的 JSON serializer 时，使用 `IsProtoMessageType`、`JSONFieldName`、`HasJSONTagOption` 复用类型判断。
