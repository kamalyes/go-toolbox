# convert 详细示例

## 1. 泛型类型转换

```go
// MustString - 任意类型转string
s := convert.MustString[int](42)                  // "42"
s := convert.MustString[float64](3.14)            // "3.14"
s := convert.MustString[time.Time](t)             // 使用默认 time.RFC3339 layout
s := convert.MustString[time.Time](t, "2006-01-02") // 自定义 layout
s := convert.MustString[*timestamppb.Timestamp](ts) // 支持 *timestamppb.Timestamp

// MustIntT - 任意类型转整数（返回 error）
n, err := convert.MustIntT[int]("123", nil)       // 123, nil
n, err := convert.MustIntT[int](float64(3.14), nil) // 3, nil（mode=nil 默认 RoundNone）
n, err := convert.MustIntT[int64]("0x10", nil)    // 支持 0x 前缀的字符串

// MustFloatT - 任意类型转浮点（返回 error）
f, err := convert.MustFloatT[float64]("3.14", convert.RoundNone)

// ToFloat64 - 任意类型转 float64（返回 error）
f, err := convert.ToFloat64("3.14")
f, err := convert.ToFloat64(42)

// Float64ToInt - float64 转整数（含范围检查，返回 error）
n, err := convert.Float64ToInt[int](3.14, convert.RoundDown) // 3, nil
n, err := convert.Float64ToInt[uint](-1.0, convert.RoundNone) // err：负值不能转为无符号

// ParseFloat - 字符串解析为浮点（校验 NaN/Inf）
var f float64
err := convert.ParseFloat("3.14", &f)

// MustBool - 任意类型转 bool
b := convert.MustBool[string]("true")  // true
b := convert.MustBool[int](1)          // true
b := convert.MustBool[int](0)           // false

// MustConvertTo - 泛型强制转换（返回 bool 表示是否成功）
val, ok := convert.MustConvertTo[string](123)     // "123", true
val, ok := convert.MustConvertTo[int]("42")       // 42, true
val, ok := convert.MustConvertTo[bool]("true")     // true, true
val, ok := convert.MustConvertTo[float64]("3.14") // 3.14, true
val, ok := convert.MustConvertTo[[]byte]("hello") // []byte("hello"), true
```

## 2. JSON/YAML 编解码

```go
// YAML 转 JSON（返回 error）
yamlData := []byte("name: hello\nage: 30")
jsonBytes, err := convert.YAMLToJSON(yamlData)

// JSON 转 YAML（返回 error）
jsonData := []byte(`{"name":"hello","age":30}`)
yamlBytes, err := convert.JSONToYAML(jsonData)

// 字符串版本（返回 error）
jsonStr, err := convert.YAMLStringToJSON("name: hello\nage: 30")
yamlStr, err := convert.JSONStringToYAML(`{"name":"hello","age":30}`)

// YAML 转 interface{} / map
iface, err := convert.YAMLToInterface(yamlData)
m, err := convert.YAMLToMap(yamlData)

// interface{} / map 转 YAML
yamlBytes, err := convert.InterfaceToYAML(someData)
yamlBytes, err := convert.MapToYAML(map[string]interface{}{"k": "v"})

// 泛型序列化/反序列化（返回指针）
type Config struct {
    Name string `yaml:"name"`
}
cfgPtr, err := convert.UnmarshalYAML[Config](yamlData) // 返回 *Config
data, err := convert.MarshalYAML[Config](*cfgPtr)

cfgPtr, err := convert.UnmarshalJSON[Config](jsonData) // 返回 *Config
data, err := convert.MarshalJSON[Config](*cfgPtr)

// JSON 字节（返回 error，需自行 string() 转字符串）
jsonBytes, err := convert.MustJSON(myStruct)
indentBytes, err := convert.MustJSONIndent(myStruct)
jsonStr := string(jsonBytes)

// 字符串切片与 JSON 数组互转
jsonStr := convert.StringsToJSON([]string{"a", "b", "c"}) // `["a","b","c"]`
strs, err := convert.StringsFromJSON(`["a","b","c"]`)
```

## 3. 字节/十六进制/二进制转换

```go
// 零拷贝 string <-> []byte（unsafe，仅用于只读场景）
b := convert.S2B("hello")        // []byte("hello")
s := convert.B2S([]byte("hello")) // "hello"
s := convert.SliceByteToString([]byte("hello")) // 等价 B2S

// Hex（BytesToHex 输出大写）
hex := convert.BytesToHex([]byte{0xde, 0xad}) // "DEAD"
raw, err := convert.HexToBytes("dead")        // []byte{0xde, 0xad}

// Binary（查找表优化）
bin := convert.ByteToBinStr(0xff)                          // "11111111"
bin := convert.BytesToBinStr([]byte{0xff, 0x00})           // "1111111100000000"
bin := convert.BytesToBinStrWithSplit([]byte{0xff, 0x00}, " ") // "11111111 00000000"

// Decimal（使用 uint64）
dec, err := convert.HexToDec("ff")   // 255, nil
hex := convert.DecToHex(uint64(255))  // "FF"
bin := convert.DecToBin(uint64(255))  // "11111111"
bin, err := convert.HexToBin("ff")    // "11111111", nil

// BCC 校验码
bccStr, err := convert.HexToBCC("3031")              // hex 字符串
bccByte := convert.BytesToBCC([]byte{0x30, 0x31})    // byte
```

## 4. Base64 编解码

```go
// 标准编码/解码（支持 []byte / string 输入）
encoded, err := convert.B64Encode([]byte("hello"))
encoded, err := convert.B64Encode("hello")
decoded, err := convert.B64Decode(encoded)

// URL 安全编码/解码
urlEncoded, err := convert.B64UrlEncode([]byte("https://example.com"))
decoded, err := convert.B64UrlDecode(urlEncoded)

// Base64 字符串转字节
imgBytes, err := convert.B64ToByte(base64Str)
```

## 5. IP 转换

```go
// IPv4 与数值互转
n, err := convert.IP2Long(net.ParseIP("192.168.1.1")) // 3232235777, nil
ip, err := convert.Long2IP(uint(3232235777))         // 192.168.1.1, nil

// IPv6 会返回错误
_, err := convert.IP2Long(net.ParseIP("::1")) // err: invalid ipv4 format
```

## 6. 切片转换

```go
// 数字切片转字符串切片
strs := convert.NumberSliceToStringSlice[int]([]int{1, 2, 3}) // ["1","2","3"]

// 字符串切片转数字切片（返回 error）
nums, err := convert.StringSliceToNumberSlice[int]([]string{"1", "2", "3"}, nil)

// 字符串切片转浮点切片（返回 error）
floats, err := convert.StringSliceToFloatSlice[float64]([]string{"1.1", "2.2"}, convert.RoundNone)

// 任意切片/数组转 []any
ifaceSlice := convert.AnySliceToInterfaceSlice([]string{"a", "b"})
ifaceSlice := convert.AnySliceToInterfaceSlice([3]int{1, 2, 3}) // 也支持数组

// []string 转 []any（兼容包装）
ifaceSlice := convert.StringSliceToInterfaceSlice([]string{"a", "b"})

// []any 转换
strSlice := convert.InterfaceSliceToStringSlice([]any{"a", 1, true})
intSlice := convert.InterfaceSliceToIntSlice([]any{"1", "2", "x"}, nil) // [1,2]，跳过 "x"

// ToNumberSlice：支持字符串按分隔符拆分
nums, err := convert.ToNumberSlice[int]("1,2,3", ",")     // [1,2,3]
nums, err := convert.ToNumberSlice[int]([]string{"1", "2"}, ",") // [1,2]
nums := convert.MustToNumberSlice[int]("1,2,3", ",")     // panic on error
```

## 7. Map 转换

```go
// map[any]any 转 map[string]any（仅保留字符串键）
m := convert.InterfaceMapToStringMap(map[any]any{"k": "v", 1: "ignored"})

// 对象转 map（struct 走 json tag）
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
m := convert.ParseObjectToMap(&User{Name: "alice", Age: 18}) // map[name:alice age:18]
m := convert.ParseObjectToMap(map[string]any{"k": "v"})

// 键值对参数转 map（单参数对象走 ParseObjectToMap）
m := convert.ParseKVPairsToMap("name", "hello", "age", 30) // map[name:hello age:30]
m := convert.ParseKVPairsToMap(&User{Name: "alice"})       // 走对象解析
```

## 8. 键值对转换

```go
// map 转键值对切片
kvs := convert.MapToKVPairs(map[string]any{"user_id": "123", "action": "login"})
// ["user_id","123","action","login"]

kvs := convert.MapStringToKVPairs(map[string]string{"Content-Type": "application/json"})
kvs := convert.MapAnyToKVPairs(map[string]any{"k": "v"}) // 返回 []any

// 合并多个键值对切片
all := convert.MergeKVPairs(
    []any{"user_id", "123"},
    []any{"action", "login"},
)

// 键值对切片转 map
m := convert.KVPairsToMap([]any{"key1", "value1", "key2", 123})

// 便捷构造与追加
kvs := convert.KVPairs("user_id", "123", "action", "login")
kvs = convert.AddKVPair(kvs, "ip", "127.0.0.1")
kvs = convert.AddKVPairs(kvs, map[string]any{"device": "mobile"})

// 合并多个 map 并转为键值对切片
kvs := convert.MergeMapToKVPairs(
    map[string]any{"event": "login"},
    map[string]any{"ip": "127.0.0.1"},
)
```

## 9. 统计格式化

```go
// 时长格式化（秒 -> 人类可读）
s := convert.FormatDuration(45)        // "45s"
s := convert.FormatDuration(90)        // "1m 30s"
s := convert.FormatDuration(3665)      // "1h 1m 5s"
s := convert.FormatDuration(86400)     // "1d"
s := convert.FormatDuration(nil)       // "N/A"

// 数量格式化
s := convert.FormatCount(123)  // "123"
s := convert.FormatCount(nil)  // "0"

// 百分比格式化
s := convert.FormatPercentage(85.567, 1) // "85.6%"
s := convert.FormatPercentage(nil, 1)     // "0%"
s := convert.FormatPercentage(100, 0)   // "100%"
```

## 10. 字段映射转换 Transformer

```go
// 使用 NewTransformer 链式调用
err := convert.NewTransformer().
    SetDst(&dstStruct).
    SetSrc(srcStruct).
    SetOptions(&convert.TransformFieldsOptions{
        StrictTypeCheck: true,
        TimeFormat:      time.RFC3339,
        TransTagName:    "transform",
    }).
    Transform()

// 兼容旧用法
err := convert.TransformFields(&dstStruct, srcStruct, &convert.TransformFieldsOptions{
    TransTagName: "transform",
})

// 选项链式设置
opts := (&convert.TransformFieldsOptions{}).
    SetStrictTypeCheck(true).
    SetTimeFormat(time.DateTime)
```

## 11. AppendValue 高效追加

```go
buf := make([]byte, 0, 64)
buf = convert.AppendValue(buf, 42)         // 复用 stringx.FastAppendInt
buf = convert.AppendValue(buf, "hello")
buf = convert.AppendValue(buf, 3.14)       // 复用 stringx.FastFloat
buf = convert.AppendValue(buf, true)
buf = convert.AppendValue(buf, []byte("bytes"))
buf = convert.AppendValue(buf, error(nil)) // "<nil>"
```
