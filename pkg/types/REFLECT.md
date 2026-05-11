# Types 反射工具模块

本模块提供通用的反射工具函数，用于处理 protobuf 消息和结构体的反射操作。

## 功能列表

### ProtoMessageType

```go
var ProtoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
```

protobuf 消息类型的反射类型，用于判断一个类型是否实现了 `proto.Message` 接口。

### ExtractJSONKey

```go
func ExtractJSONKey(fieldType reflect.StructField) string
```

从结构体字段的 tag 中提取 JSON 键名。

**参数：**

- `fieldType` - 结构体字段信息

**返回：**

- JSON 键名，如果字段没有 json tag 或 tag 为 "-" 则返回空字符串

**示例：**

```go
type User struct {
    ID   string `json:"id"`
    Name string `json:"name,omitempty"`
    Age  int    `json:"-"`
    Addr string
}

// ExtractJSONKey 会返回：
// ID 字段 -> "id"
// Name 字段 -> "name"
// Age 字段 -> "" (因为 json:"-")
// Addr 字段 -> "Addr" (没有 json tag，使用字段名)
```

### EnsureStructDefaults

```go
func EnsureStructDefaults(v reflect.Value)
```

确保结构体的 protobuf 指针字段和嵌套结构体指针字段被初始化。

**参数：**

- `v` - 结构体的反射值

**功能：**

- 遍历结构体的所有字段
- 对于实现了 `proto.Message` 接口的 nil 指针字段，初始化为新的 protobuf 消息实例
- 对于指向结构体的 nil 指针字段，初始化为新的结构体实例

**示例：**

```go
type Config struct {
    Name *wrapperspb.StringValue `json:"name"`
    Age  *wrapperspb.Int32Value  `json:"age"`
    Meta *Metadata               `json:"meta"`
}

type Metadata struct {
    Key string `json:"key"`
}

var cfg Config
v := reflect.ValueOf(&cfg).Elem()
EnsureStructDefaults(v)
// 现在 cfg.Name, cfg.Age, cfg.Meta 都不再是 nil
```

### NewProtoMessage

```go
func NewProtoMessage[T proto.Message]() T
```

创建一个新的 protobuf 消息实例。

**类型参数：**

- `T` - protobuf 消息类型，必须实现 `proto.Message` 接口

**返回：**

- 新创建的 protobuf 消息实例

**示例：**

```go
// 创建 StringValue 实例
sv := NewProtoMessage[*wrapperspb.StringValue]()
sv.Value = "hello"

// 创建 Int32Value 实例
iv := NewProtoMessage[*wrapperspb.Int32Value]()
iv.Value = 42
```

## 使用场景

1. **数据库 JSON 序列化**：配合 `ProtoJSON[T]` 使用，自动处理 protobuf 消息的 JSON 序列化
2. **动态结构体处理**：在需要通过反射操作结构体字段时使用
3. **protobuf 消息初始化**：确保 protobuf 指针字段在使用前被正确初始化

## 注意事项

1. `EnsureStructDefaults` 只处理指针类型的字段，非指针字段不受影响
2. `ExtractJSONKey` 遵循 Go 标准库的 json tag 解析规则
3. `NewProtoMessage` 要求类型参数必须是指向 protobuf 消息的指针类型
