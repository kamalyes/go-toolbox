# Types 反射工具模块

本模块提供通用的反射工具函数，用于处理 protobuf 消息和结构体的反射操作
源码位于 `pkg/types/reflect.go`

## 功能列表

### ProtoMessageType

```go
var ProtoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
```

protobuf 消息类型的反射类型，用于判断一个类型是否实现了 `proto.Message` 接口

### IsProtoMessageType

```go
func IsProtoMessageType(t reflect.Type) bool
```

判断类型是否实现 `proto.Message` 接口`t` 为 nil 时返回 false

```go
t := reflect.TypeOf((*wrapperspb.StringValue)(nil))
types.IsProtoMessageType(t) // true

t2 := reflect.TypeOf(0)
types.IsProtoMessageType(t2) // false
```

### IsExportedField

```go
func IsExportedField(field reflect.StructField) bool
```

判断结构体字段是否可按 JSON 规则参与导出处理`PkgPath` 为空（导出字段）或
匿名字段（嵌入字段）时返回 true

```go
field, _ := reflect.TypeOf(User{}).FieldByName("Name")
types.IsExportedField(field) // true（导出字段）
```

### ExtractJSONKey

```go
func ExtractJSONKey(fieldType reflect.StructField) string
```

从结构体字段的 tag 中提取 JSON 键名

**参数：**

- `fieldType` - 结构体字段信息

**返回：**

- JSON 键名；字段没有 json tag、tag 为 "-" 或空名时返回空字符串

**行为说明：**

- 无 json tag：返回 `""`
- `json:"-"`：返回 `""`
- `json:"name"`：返回 `"name"`
- `json:"name,omitempty"`：去除选项后返回 `"name"`
- `json:",omitempty"`（仅选项无名称）：返回字段名

```go
type User struct {
    ID   string `json:"id"`
    Name string `json:"name,omitempty"`
    Age  int    `json:"-"`
    Addr string
}

// ExtractJSONKey 返回：
// ID 字段   -> "id"
// Name 字段 -> "name"
// Age 字段  -> ""  (因为 json:"-")
// Addr 字段 -> ""  (没有 json tag)
```

### JSONFieldName

```go
func JSONFieldName(fieldType reflect.StructField) string
```

获取结构体字段的 JSON 字段名优先使用 `ExtractJSONKey` 的结果，
没有显式名称时返回 Go 字段名

```go
type User struct {
    ID   string `json:"id"`
    Addr string
}

// JSONFieldName 返回：
// ID 字段   -> "id"
// Addr 字段 -> "Addr" (没有 json tag，回退到字段名)
```

### HasJSONTagOption

```go
func HasJSONTagOption(fieldType reflect.StructField, options ...string) bool
```

判断结构体字段的 json tag 是否包含指定选项（如 `omitempty`、`omitzero`）
无 json tag、tag 为 "-"、无选项或字段未指定任何匹配选项时返回 false

```go
type User struct {
    Name string `json:"name,omitempty"`
    Age  int    `json:"age"`
}

field, _ := reflect.TypeOf(User{}).FieldByName("Name")
types.HasJSONTagOption(field, "omitempty")          // true
types.HasJSONTagOption(field, "omitempty", "string") // true（满足任一即可）
```

### EnsureStructDefaults

```go
func EnsureStructDefaults(v reflect.Value)
```

确保结构体的 protobuf 指针字段和嵌套结构体指针字段被初始化

**参数：**

- `v` - 结构体的反射值（非 Struct 类型时直接返回）

**功能：**

- 遍历结构体的所有字段
- 对于实现了 `proto.Message` 接口的 nil 指针字段，初始化为新的 protobuf 消息实例
- 对于指向结构体的 nil 指针字段，初始化为新的结构体实例

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
types.EnsureStructDefaults(v)
// 现在 cfg.Name, cfg.Age, cfg.Meta 都不再是 nil
```

### NewProtoMessage

```go
func NewProtoMessage[T proto.Message]() T
```

创建一个新的 protobuf 消息实例

**类型参数：**

- `T` - protobuf 消息类型，必须实现 `proto.Message` 接口

**返回：**

- 新创建的 protobuf 消息实例

```go
// 创建 StringValue 实例
sv := types.NewProtoMessage[*wrapperspb.StringValue]()
sv.Value = "hello"

// 创建 Int32Value 实例
iv := types.NewProtoMessage[*wrapperspb.Int32Value]()
iv.Value = 42
```

### UnwrapPBValue

```go
func UnwrapPBValue(iface interface{}) interface{}
```

解包 protobuf wrapper 类型，返回底层值支持 `wrapperspb` 全部 9 种类型：
`String`、`Bool`、`Int32`、`Int64`、`UInt32`、`UInt64`、`Float`、`Double`、`Bytes`

- wrapper 为 nil 或 typed nil 时返回 nil
- 不是 wrapper 类型时返回原始值

```go
types.UnwrapPBValue(wrapperspb.String("hello"))   // "hello"
types.UnwrapPBValue(wrapperspb.Int32(42))          // int32(42)
types.UnwrapPBValue((*wrapperspb.StringValue)(nil)) // nil
types.UnwrapPBValue("plain string")                 // "plain string"
```

### ResolveModelKey

```go
func ResolveModelKey(fieldType reflect.StructField) string
```

解析 Model 结构体字段的 key，按优先级返回：
1. gorm tag 的 `column:` 值
2. json tag 的名称（`json:"-"` 时返回 `"-"`）
3. Go 字段名

```go
type User struct {
    Name string `gorm:"column:user_name" json:"name"`
    Age  int    `json:"age"`
    Raw  string `json:"-"`
}

// ResolveModelKey 返回：
// Name -> "user_name" (gorm column 优先)
// Age  -> "age"
// Raw  -> "-" (json:"-")
```

### ExtractGormColumn

```go
func ExtractGormColumn(tag string) string
```

从 gorm tag 字符串中提取 `column:` 名称支持分号分隔的多段 tag
未找到 `column:` 段时返回空字符串

```go
types.ExtractGormColumn("column:user_name;size:255") // "user_name"
types.ExtractGormColumn("size:255")                   // ""
```

### ResolvePBKey

```go
func ResolvePBKey(fieldType reflect.StructField) string
```

解析 PB 结构体字段的 key，按优先级返回：
1. protobuf tag 的 `name=` 值
2. json tag 的名称（`json:"-"` 时返回 `"-"`）
3. Go 字段名

```go
type User struct {
    Name string `protobuf:"name=user_name,json=name" json:"name"`
}

// ResolvePBKey 返回 "user_name" (protobuf name 优先)
```

### ExtractPBTagValue

```go
func ExtractPBTagValue(tag string, key string) string
```

从 protobuf tag 字符串中提取指定键的值tag 以逗号分隔，键值格式为 `key=value`
未找到时返回空字符串

```go
types.ExtractPBTagValue("name=user_name,json=name,proto3", "name") // "user_name"
types.ExtractPBTagValue("name=user_name,json=name,proto3", "json") // "name"
```

### UnwrapModelValue

```go
func UnwrapModelValue(iface interface{}) interface{}
```

解引用指针类型，返回底层值不是指针或指针为 nil 时返回原始值

```go
n := 42
types.UnwrapModelValue(&n)  // 42
types.UnwrapModelValue(42)  // 42
types.UnwrapModelValue((*int)(nil)) // nil
```

## 使用场景

1. **数据库 JSON 序列化**：配合 `ProtoJSON[T]` 使用，自动处理 protobuf 消息的 JSON 序列化
2. **动态结构体处理**：在需要通过反射操作结构体字段时使用
3. **protobuf 消息初始化**：确保 protobuf 指针字段在使用前被正确初始化
4. **ORM 字段映射**：通过 `ResolveModelKey` / `ResolvePBKey` 统一获取字段名
5. **wrapper 解包**：通过 `UnwrapPBValue` 将 protobuf wrapper 还原为基本类型

## 注意事项

1. `EnsureStructDefaults` 只处理指针类型的字段，非指针字段不受影响
2. `ExtractJSONKey` 无 json tag 时返回空字符串；需要回退到字段名时使用 `JSONFieldName`
3. `NewProtoMessage` 要求类型参数必须是指向 protobuf 消息的指针类型
4. `UnwrapPBValue` 对 typed nil（如 `(*wrapperspb.StringValue)(nil)`）返回 nil
5. `ResolveModelKey` / `ResolvePBKey` 不会受 `ExtractJSONKey` 空字符串结果影响，
   它们会在 json 名称为空时回退到 Go 字段名
