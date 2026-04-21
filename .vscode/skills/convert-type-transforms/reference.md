# convert 详细示例

## 1. 泛型类型转换

```go
// MustString - 任意类型转string
s := convert.MustString[int](42)           // "42"
s := convert.MustString[float64](3.14)      // "3.14"
s := convert.MustString[time.Time](t)       // 使用默认layout
s := convert.MustString[time.Time](t, "2006-01-02") // 自定义layout

// MustIntT - 任意类型转整数
n := convert.MustIntT[string]("123")       // 123
n := convert.MustIntT[float64](3.14)        // 3

// MustBool - 任意类型转bool
b := convert.MustBool[string]("true")       // true
b := convert.MustBool[int](1)               // true
b := convert.MustBool[int](0)               // false

// MustConvertTo - 泛型强制转换
val := convert.MustConvertTo[int]("42")     // 42
```

## 2. JSON/YAML编解码

```go
// YAML转JSON
yamlData := []byte("name: hello\nage: 30")
jsonBytes := convert.YAMLToJSON(yamlData)

// JSON转YAML
jsonData := []byte(`{"name":"hello","age":30}`)
yamlBytes := convert.JSONToYAML(jsonData)

// 字符串版
jsonStr := convert.YAMLStringToJSON("name: hello\nage: 30")
yamlStr := convert.JSONStringToYAML(`{"name":"hello","age":30}`)

// 泛型序列化/反序列化
type Config struct {
    Name string `yaml:"name"`
}
cfg, err := convert.UnmarshalYAML[Config](yamlData)
data, err := convert.MarshalYAML[Config](cfg)
cfg, err := convert.UnmarshalJSON[Config](jsonData)
data, err := convert.MarshalJSON[Config](cfg)

// JSON缩进
s := convert.MustJSONIndent(myStruct)
s := convert.MustJSON(myStruct)
```

## 3. 字节/十六进制/二进制转换

```go
// Hex
hex := convert.BytesToHex([]byte{0xde, 0xad})  // "dead"
raw := convert.HexToBytes("dead")                // []byte{0xde, 0xad}

// Binary
bin := convert.ByteToBinStr(0xff)                // "11111111"
bin := convert.BytesToBinStr([]byte{0xff, 0x00}) // "1111111100000000"
bin := convert.BytesToBinStrWithSplit([]byte{0xff, 0x00}, " ") // "11111111 00000000"

// Decimal conversions
dec := convert.HexToDec("ff")   // 255
hex := convert.DecToHex(255)    // "ff"
bin := convert.DecToBin(255)     // "11111111"
bin := convert.HexToBin("ff")    // "11111111"

// BCC checksum
bcc := convert.HexToBCC("3031")  // BCC码字符串
bcc := convert.BytesToBCC([]byte{0x30, 0x31}) // byte
```

## 4. 切片转换

```go
// 数字切片转字符串切片
strs := convert.NumberSliceToStringSlice[int]([]int{1, 2, 3}) // ["1","2","3"]

// 字符串切片转数字切片
nums := convert.StringSliceToNumberSlice[int]([]string{"1","2","3"}, convert.RoundModeDown)

// 浮点切片
floats := convert.StringSliceToFloatSlice[float64]([]string{"1.1","2.2"}, convert.RoundModeDown)

// interface切片互转
ifaceSlice := convert.AnySliceToInterfaceSlice([]string{"a","b"})
strSlice := convert.InterfaceSliceToStringSlice(ifaceSlice)
intSlice := convert.InterfaceSliceToIntSlice(ifaceSlice, convert.RoundModeDown)
```

## 5. 字段映射转换 Transformer

```go
// 使用Transformer进行结构体字段映射
transformer := convert.NewTransformer()
err := convert.TransformFields(&dstStruct, srcStruct, convert.TransformFieldsOptions{
    // 配置选项
})

// KV解析
m := convert.ParseKVPairsToMap("name", "hello", "age", "30") // map[name:hello age:30]

// 对象转map
m := convert.ParseObjectToMap(myStruct)
```