# stringx 详细示例

## 1. 链式调用 StringX

```go
// 创建并链式操作（New 返回 *StringX，方法修改内部 value）
result := stringx.New("  Hello World  ").
    TrimChain().
    ToLowerChain().
    ReplaceAllChain("world", "go")

fmt.Println(result.String()) // "hello go"

// 多步链式
result := stringx.New("user@example.com").
    SubBeforeChain("@", false).
    ToUpperChain()
fmt.Println(result.String()) // "USER"

// 获取值
val := stringx.New("hello").Value() // "hello"
```

## 2. 子串提取

```go
// SubBefore - 提取分隔符前的内容（第三个参数 isLastSeparator）
user := stringx.SubBefore("user@host", "@", false)     // "user"
userLast := stringx.SubBefore("a.b.c", ".", true)       // "a.b"（取最后一个分隔符）

// SubAfter - 提取分隔符后的内容
domain := stringx.SubAfter("user@host", "@", false)      // "host"
ext := stringx.SubAfter("file.tar.gz", ".tar", false)    // ".gz"

// SubBetween - 提取两标记之间
name := stringx.SubBetween("[hello]", "[", "]")   // "hello"

// SubBetweenAll - 提取所有匹配
parts := stringx.SubBetweenAll("a<b>c<d>e", "<", ">") // ["b", "d"]

// SubString - 按位置截取
s := stringx.SubString("hello world", 0, 5) // "hello"
```

## 3. 脱敏隐藏

```go
// 手机号脱敏（[startInclude, endExclude) 区间）
phone := stringx.Hide("13812345678", 3, 7)  // "138****5678"

// 邮箱脱敏
email := stringx.Hide("user@example.com", 2, 5) // "us***@example.com"

// 身份证脱敏
id := stringx.Hide("110101199001011234", 6, 14) // "110101********1234"
```

## 4. 填充对齐

```go
// Pad - 使用 * 填充，通过 *Paddler 指定位置（默认 Middle）
padded := stringx.Pad("abc", 10)  // 默认中间填充
padded := stringx.Pad("abc", 10, &stringx.Paddler{Position: stringx.Left})  // "*******abc"
padded := stringx.Pad("abc", 10, &stringx.Paddler{Position: stringx.Right}) // "abc*******"
padded := stringx.Pad("abc", 10, &stringx.Paddler{Position: stringx.Middle}) // 中间填充

// FillBefore / FillAfter - 使用指定字符填充
padded := stringx.FillBefore("abc", "0", 10)  // "0000000abc"
padded := stringx.FillAfter("abc", "-", 10)  // "abc-------"

// CJK对齐
padded := stringx.FillAfter("中文", " ", 10) // "中文      "
// DisplayWidth考虑CJK宽度
width := stringx.DisplayWidth("中文") // 4
```

## 5. 前后缀判断

```go
// 基本前后缀
ok := stringx.StartWith("hello", "he")              // true
ok := stringx.EndWith("hello", "lo")                // true

// 忽略大小写
ok := stringx.StartWithIgnoreCase("Hello", "he")    // true
ok := stringx.EndWithIgnoreCase("Hello", "LO")     // true

// 任一匹配（接收 []string 切片）
ok := stringx.StartWithAny("hello", []string{"he", "hi"})     // true
ok := stringx.EndWithAny("hello", []string{"lo", "go"})       // true
ok := stringx.StartWithAnyIgnoreCase("Hello", []string{"HE"}) // true
```

## 6. 替换操作

```go
// 替换前n个
s := stringx.Replace("hello hello", "hello", "hi", 1) // "hi hello"

// 替换所有
s := stringx.ReplaceAll("hello hello", "hello", "hi") // "hi hi"

// 按索引范围替换 [startIndex, endIndex)
s := stringx.ReplaceWithIndex("hello world", 0, 5, "hi") // "hi world"

// 按正则表达式替换（regex 字符串 + 替换函数）
s := stringx.ReplaceWithMatcher("hello 123 world", `\d+`, func(match string) string {
    return "***"
}) // "hello *** world"

// 替换特殊字符为指定 rune
s := stringx.ReplaceSpecialChars("hello@world#test", '_') // "hello_world_test"
```

## 7. 修剪操作

```go
s := stringx.TrimProtocol("https://example.com") // "example.com"
s := stringx.TrimAll("aaabbbccc", "abc")           // ""
s := stringx.TrimPrefix("prefix_value", "prefix_") // "value"
s := stringx.TrimPrefixIgnoreCase("PREFIX_value", "prefix_") // "value"
s := stringx.TrimAllLineBreaks("hello\nworld\r\n") // "helloworld"
s := stringx.TrimSymbols("hello, world!") // "helloworld"
s := stringx.TrimAny("hello world test", []string{"hello ", " test"}) // "world"
```

## 8. 命名风格转换

```go
// 各种命名法转换
pascal := stringx.ToPascalCase("user_name")   // "UserName"
camel := stringx.ToCamelCase("UserName")      // "userName"
snake := stringx.ToSnakeCase("UserName")       // "user_name"
kebab := stringx.ToKebabCase("UserName")       // "user-name"

// 通过 CharacterStyle 枚举转换
s := stringx.ConvertCharacterStyle("user_name", stringx.PascalCharacterStyle) // "UserName"

// 获取所有命名变体
variants := stringx.NormalizeFieldName("user_name")
// ["user_name", "UserName", "userName", "user-name"]
```

## 9. 包含判断

```go
ok := stringx.Contains("hello world", "world")           // true
ok := stringx.ContainsIgnoreCase("Hello World", "WORLD") // true
ok := stringx.ContainsAny("hello", []string{"he", "hi"}) // true
ok := stringx.ContainsAll("hello world", []string{"hello", "world"}) // true
ok := stringx.ContainsBlank("hello world") // true
found := stringx.GetContainsStr("hello", []string{"hi", "he"}) // "he"
```

## 10. 分割操作

```go
parts := stringx.Split("a,b,c", ",")              // ["a", "b", "c"]
parts := stringx.SplitLimit("a,b,c", ",", 2)      // ["a", "b,c"]
parts := stringx.SplitTrim("a , b , c", ",")       // ["a", "b", "c"]
parts := stringx.SplitByLen("abcdef", 2)           // ["ab", "cd", "ef"]
parts := stringx.Cut("abcdef", 3)                  // ["ab", "cd", "ef"]

// 泛型分割转换
nums := stringx.SplitAfterMapping("1,2,3", ",", func(s string) (int, error) {
    return strconv.Atoi(s)
}) // [1, 2, 3]

// 去重
unique := stringx.UniqueStringSlice([]string{"a", "b", "a", "", "c"}) // ["a", "b", "c"]
```

## 11. 索引查找

```go
idx := stringx.IndexOf("hello world", "world")      // 6
idx := stringx.IndexOfIgnoreCase("Hello World", "world") // 6
idx := stringx.LastIndexOf("hello hello", "hello")  // 6
idx := stringx.IndexOfByRange("hello world", "world", 3, 11) // 6
idx := stringx.OrdinalIndexOf("a-b-c-d", "-", 2)    // 3（第2次出现）
```

## 12. 工具函数

```go
length := stringx.Length("中文")        // 2 (rune长度)
width := stringx.DisplayWidth("中文")   // 4 (CJK宽度为2)
reversed := stringx.Reverse("hello")   // "olleh"
joined := stringx.Coalesce("a", "b", "c") // "abc"
md5 := stringx.CalculateMD5Hash("hello") // "5d41402abc4b2a76b9719d911017c592"
slug := stringx.SanitizeSlug("Hello World!") // "hello-world"
dir := stringx.NormalizeSQLDirection("asc", "desc") // "ASC"
```

## 13. 快速格式化

```go
// 快速整数转字符串（0-9999 使用缓存）
s := stringx.FastItoa(42)    // "42"
s := stringx.FastFloat(3.14, 2) // "3.14"

// 追加整数到 buffer
buf := make([]byte, 0, 16)
buf = stringx.FastAppendInt(buf, 2026) // []byte("2026")

// 格式化时间到 buffer
t := time.Now()
buf = stringx.FastFormatTime(buf, t)         // "2026/1/5 10:30:45 "
buf = stringx.FastFormatTimeISO(buf, t)      // "2026-01-05 10:30:45"
buf = stringx.FastFormatTimeCompact(buf, t)  // "20260105103045"
```

## 14. 域名处理

```go
prefix := stringx.ExtractDomainPrefix("www.example.com", "example.com") // "www"
prefix := stringx.ExtractDomainPrefix("example.com", "example.com")     // "@"
ok := stringx.IsSubdomain("www.example.com", "example.com")              // true
prefix, primary := stringx.SplitDomain("www.api.example.com", "example.com") // "www.api", "example.com"
```

## 15. 格式化与截断

```go
// map 格式化
s := stringx.Format("{a} and {b}", map[string]interface{}{"a": "aValue", "b": "bValue"}) // "aValue and bValue"

// 有序格式化
s := stringx.IndexedFormat("this is {0} for {1}", []interface{}{"a", "b"}) // "this is a for b"

// 截断
s := stringx.Truncate("hello world", 5)               // "hello"
s := stringx.TruncateAppendEllipsis("hello world", 5)  // "hello..."
s := stringx.TruncateMessage("very long log message", 10) // "very long ..."

// 前后缀补充
s := stringx.AddPrefixIfNot("value", "prefix_")  // "prefix_value"
s := stringx.AddPrefixIfNot("prefix_value", "prefix_") // "prefix_value"
s := stringx.AddSuffixIfNot("hello", ".txt")     // "hello.txt"
```

## 16. JSON 引用

```go
// 按 JSON 规则转义字符串
quoted := stringx.QuoteJSONBytes(`hello "world"`)
// 结果: []byte(`"hello \"world\""`)
```

## 17. 字段解析

```go
// 解析数字字段（带范围校验）
v, err := stringx.ParseFieldInt("5", 0, 10)  // 5, nil
v, err := stringx.ParseFieldInt("15", 0, 10) // 0, error（超出范围）

// 支持通配符
v, err := stringx.ParseFieldIntOrWildcard("*", "*", 0, 10) // -1, nil（通配符返回 -1）
v, err := stringx.ParseFieldIntOrWildcard("5", "*", 0, 10) // 5, nil
```
