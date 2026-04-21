# stringx 详细示例

## 1. 链式调用 StringX

```go
// 创建并链式操作
result := stringx.New("  Hello World  ").
    TrimChain().
    ToLowerChain().
    ReplaceChain("world", "go")

fmt.Println(result.String()) // "hello go"

// 多步链式
result := stringx.New("user@example.com").
    SubBeforeChain("@").
    ToUpperChain()
fmt.Println(result.String()) // "USER"
```

## 2. 子串提取

```go
// SubBefore - 提取分隔符前的内容
user := stringx.SubBefore("user@host", "@")     // "user"
path := stringx.SubBefore("/a/b/c", "/")          // ""

// SubAfter - 提取分隔符后的内容
domain := stringx.SubAfter("user@host", "@")      // "host"
ext := stringx.SubAfter("file.tar.gz", ".tar")    // ".gz"

// SubBetween - 提取两标记之间
name := stringx.SubBetween("[hello]", "[", "]")   // "hello"

// SubBetweenAll - 提取所有匹配
parts := stringx.SubBetweenAll("a<b>c<d>e", "<", ">") // ["b", "d"]

// SubString - 按位置截取
s := stringx.SubString("hello world", 0, 5) // "hello"
```

## 3. 脱敏隐藏

```go
// 手机号脱敏
phone := stringx.Hide("13812345678", 3, 7)  // "138****5678"

// 邮箱脱敏
email := stringx.Hide("user@example.com", 2, 5) // "us***@example.com"

// 身份证脱敏
id := stringx.Hide("110101199001011234", 6, 14) // "110101********1234"
```

## 4. 填充对齐

```go
// 左填充（右对齐）
padded := stringx.Pad("abc", 10, stringx.PadLeft, ' ')  // "       abc"

// 右填充（左对齐）
padded := stringx.Pad("abc", 10, stringx.PadRight, '-')  // "abc-------"

// CJK对齐
padded := stringx.Pad("中文", 10, stringx.PadRight, ' ') // "中文      "
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
ok := stringx.EndWithIgnoreCase("Hello", "LO")       // true

// 任一匹配
ok := stringx.StartWithAny("hello", "he", "hi")     // true
ok := stringx.EndWithAny("hello", "lo", "go")        // true
```

## 6. 替换操作

```go
// 替换前n个
s := stringx.Replace("hello hello", "hello", "hi", 1) // "hi hello"

// 替换所有
s := stringx.ReplaceAll("hello hello", "hello", "hi") // "hi hi"

// 按索引范围替换
s := stringx.ReplaceWithIndex("hello world", 0, 5, "hi") // "hi world"

// 按匹配函数替换
s := stringx.ReplaceWithMatcher("hello 123 world", func(s string) bool {
    for _, r := range s { return r >= '0' && r <= '9' }
    return false
}, "***") // "hello *** world"

// 替换特殊字符
s := stringx.ReplaceSpecialChars("hello@world#test", "_") // "hello_world_test"
```

## 7. 修剪操作

```go
s := stringx.TrimProtocol("https://example.com") // "example.com"
s := stringx.TrimAll("aaabbbccc", "abc")         // ""
s := stringx.TrimPrefix("prefix_value", "prefix_") // "value"
s := stringx.TrimPrefixIgnoreCase("PREFIX_value", "prefix_") // "value"
s := stringx.TrimAllLineBreaks("hello\nworld\r\n") // "helloworld"
s := stringx.TrimSymbols("hello, world!") // "hello world"
```

## 8. 工具函数

```go
length := stringx.Length("中文")        // 2 (rune长度)
width := stringx.DisplayWidth("中文")   // 4 (CJK宽度为2)
idx := stringx.IndexOf("hello world", "world") // 6

field := stringx.NormalizeFieldName("user_name") // "UserName" 或类似规范化
```