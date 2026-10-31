---
name: protox-protobuf
description: protobuf 互操作工具，提供 wrapperspb 包装类型双向转换、timestamppb 与 time.Time 互转（带默认时长）、proto.Message 类型判断与泛型实例创建、wrapper 解包，是 go-toolbox 中唯一 import protobuf 的通用包
---

# protox - protobuf 互操作

提供 wrapperspb 包装类型的双向转换、timestamppb 与 `time.Time` 互转、`proto.Message` 类型判断与泛型实例创建

## 这个包为什么存在

**它是 go-toolbox 里唯一直接 import `google.golang.org/protobuf` 的通用工具包。**

这些函数原先分散在 `safe`、`types`、`mathx` 三个包里，而那三个包被几乎所有业务代码引用
结果是：只想用 `safe.New(x).Field("Name").String("")` 做一次链式安全访问的人，
也要编译整个 protobuf 运行时
把它们集中到 protox 之后，那三个包的依赖图里都没有 protobuf 了（可用
`go list -deps ./pkg/safe | grep protobuf` 验证，应为空）

**所有搬迁都只换包名，签名没有变。**

| 原位置 | 搬来的符号 |
|---|---|
| `safe` | `TimeToTimestampPB`、`StringToPB`、`BoolToPB`、`Int32ToPB`、`Int64ToPB`、`DoubleToPB`、`PtrToTime`、`PtrToString`、`PtrToBool`、`PtrToInt32`、`PtrToInt64`、`PtrToDouble` |
| `types` | `ProtoMessageType`、`IsProtoMessageType`、`NewProtoMessage[T]`、`UnwrapPBValue` |
| `mathx` | `IfProtoTimeOr`、`IfProtoTimeOrPtr`、`IfTimeToProto` |

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/protox"
```

Go 值转 PB 包装类型（nil 入参返回 nil，不会 panic）：
```go
pbName := protox.StringToPB(&name)        // *string        -> *wrapperspb.StringValue
pbAge  := protox.Int32ToPB(&age)          // *int32         -> *wrapperspb.Int32Value
pbTime := protox.TimeToTimestampPB(&now)  // *time.Time     -> *timestamppb.Timestamp
```

PB 包装类型转回 Go 值：
```go
name := protox.PtrToString(pbName)        // *wrapperspb.StringValue -> *string
t    := protox.PtrToTime(pbTime)          // *timestamppb.Timestamp  -> *time.Time
```

时间转换带默认时长（PB 时间为 nil 时用 `duration` 兜底）：
```go
expire   := protox.IfProtoTimeOr(req.ExpireAt, 24*time.Hour)      // 返回 time.Time
expirePtr := protox.IfProtoTimeOrPtr(req.ExpireAt, 24*time.Hour)  // 返回 *time.Time
pbExpire := protox.IfTimeToProto(&expire, 24*time.Hour)           // 返回 *timestamppb.Timestamp
```

泛型创建与类型判断：
```go
msg := protox.NewProtoMessage[*pb.User]()          // 约束 T proto.Message
ok  := protox.IsProtoMessageType(reflect.TypeOf(v))
raw := protox.UnwrapPBValue(anyWrapper)            // 解包 wrapperspb，非 wrapper 原样返回
```

## 完整API索引

### Go 值 → PB 包装类型

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `StringToPB` | `func(src *string) *wrapperspb.StringValue` | nil 返回 nil |
| `BoolToPB` | `func(src *bool) *wrapperspb.BoolValue` | nil 返回 nil |
| `Int32ToPB` | `func(src *int32) *wrapperspb.Int32Value` | nil 返回 nil |
| `Int64ToPB` | `func(src *int64) *wrapperspb.Int64Value` | nil 返回 nil |
| `DoubleToPB` | `func(src *float64) *wrapperspb.DoubleValue` | nil 返回 nil |
| `TimeToTimestampPB` | `func(src *time.Time) *timestamppb.Timestamp` | nil 返回 nil |

### PB 包装类型 → Go 值

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `PtrToString` | `func(src *wrapperspb.StringValue) *string` | nil 返回 nil |
| `PtrToBool` | `func(src *wrapperspb.BoolValue) *bool` | nil 返回 nil |
| `PtrToInt32` | `func(src *wrapperspb.Int32Value) *int32` | nil 返回 nil |
| `PtrToInt64` | `func(src *wrapperspb.Int64Value) *int64` | nil 返回 nil |
| `PtrToDouble` | `func(src *wrapperspb.DoubleValue) *float64` | nil 返回 nil |
| `PtrToTime` | `func(src *timestamppb.Timestamp) *time.Time` | nil 返回 nil |

### 带默认时长的时间转换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `IfProtoTimeOr` | `func(protoTime *timestamppb.Timestamp, duration time.Duration) time.Time` | PB 时间为 nil 时按 duration 兜底 |
| `IfProtoTimeOrPtr` | `func(protoTime *timestamppb.Timestamp, duration time.Duration) *time.Time` | 同上，返回指针 |
| `IfTimeToProto` | `func(t *time.Time, duration time.Duration) *timestamppb.Timestamp` | t 为 nil 时按 duration 兜底 |

### 类型判断与创建

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ProtoMessageType` | `reflect.Type`（变量） | `*proto.Message` 的 Elem，用于接口实现判断 |
| `IsProtoMessageType` | `func(t reflect.Type) bool` | 判断类型是否实现 `proto.Message` |
| `NewProtoMessage[T]` | `func[T proto.Message]() T` | 创建 T 的新实例 |
| `UnwrapPBValue` | `func(iface interface{}) interface{}` | 解包 wrapperspb 全部 9 种；非 wrapper 返回原值 |

## 注意事项

- **所有转换函数的 nil 入参都返回 nil 或零值，不会 panic**，调用处不必先判空
- 新代码不要再从 `safe`、`types`、`mathx` 找这些函数——那边已经没有了，编译会直接报
  undefined；对照表见上文「这个包为什么存在」
- `UnwrapPBValue` 覆盖 wrapperspb 的全部 9 种包装类型；传入非 wrapper 时原样返回，
  所以可以无条件调用
- `types` 里保留的 protobuf 相关函数只做 **tag 文本解析**（`ExtractPBTagValue`、
  `ResolvePBKey`），它们只认字符串、不 import protobuf，不要搬到本包
- 引入本包就等于引入 `google.golang.org/protobuf`。如果一个包只是顺带转一两个
  wrapper，优先考虑把转换放在调用方的边界层，而不是让底层包依赖它
