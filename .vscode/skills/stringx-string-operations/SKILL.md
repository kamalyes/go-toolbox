---
name: stringx-string-operations
description: 字符串操作链式工具，提供链式调用、子串提取、隐藏脱敏、填充对齐、显示宽度计算、前后缀判断、替换等。当需要对字符串做链式变换、截取前后子串、脱敏或格式化对齐时使用。
---

# stringx - 字符串操作链式工具

提供链式字符串变换、子串提取、脱敏隐藏、填充对齐、显示宽度计算与前后缀匹配。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/stringx"
```

链式调用：
```go
result := stringx.New("hello world").ToUpperChain().ReplaceChain("WORLD", "GO")
```

子串提取与隐藏：
```go
before := stringx.SubBefore("user@host", "@")
hidden := stringx.Hide("13812345678", 3, 7)
```

## 完整API索引

### 函数

#### 构造函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `New` | `func(value string) StringX` | 创建StringX链式对象 |

#### 大小写转换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ToLower` | `func(s string) string` | 转小写 |
| `ToUpper` | `func(s string) string` | 转大写 |
| `ToTitle` | `func(s string) string` | 转标题格式 |

#### 修剪操作

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Trim` | `func(s string) string` | 去除两端空白 |
| `TrimStart` | `func(s string) string` | 去除前导空白 |
| `TrimEnd` | `func(s string) string` | 去除尾部空白 |
| `CleanEmpty` | `func(s string) string` | 清理空白字符 |
| `TrimProtocol` | `func(s string) string` | 去除协议前缀 |
| `TrimAll` | `func(s, cutset string) string` | 去除所有指定字符 |
| `TrimAny` | `func(s, cutset string) string` | 去除任意指定字符 |
| `TrimAllLineBreaks` | `func(s string) string` | 去除所有换行符 |
| `TrimPrefix` | `func(s, prefix string) string` | 去除指定前缀 |
| `TrimPrefixIgnoreCase` | `func(s, prefix string) string` | 去除指定前缀（忽略大小写） |
| `TrimSuffix` | `func(s, suffix string) string` | 去除指定后缀 |
| `TrimSuffixIgnoreCase` | `func(s, suffix string) string` | 去除指定后缀（忽略大小写） |
| `TrimSymbols` | `func(s string) string` | 去除符号字符 |

#### 子串提取

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `SubBefore` | `func(s, sep string) string` | 提取分隔符前的子串 |
| `SubAfter` | `func(s, sep string) string` | 提取分隔符后的子串 |
| `SubBetween` | `func(s, before, after string) string` | 提取两标记之间的子串 |
| `SubBetweenAll` | `func(s, before, after string) []string` | 提取所有两标记之间的子串 |
| `SubString` | `func(s string, start, length int) string` | 按位置截取子串 |

#### 前后缀判断

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `StartWith` | `func(s, prefix string) bool` | 判断前缀 |
| `StartWithIgnoreCase` | `func(s, prefix string) bool` | 判断前缀（忽略大小写） |
| `StartWithAny` | `func(s string, prefixes ...string) bool` | 判断是否以任一前缀开头 |
| `EndWith` | `func(s, suffix string) bool` | 判断后缀 |
| `EndWithIgnoreCase` | `func(s, suffix string) bool` | 判断后缀（忽略大小写） |
| `EndWithAny` | `func(s string, suffixes ...string) bool` | 判断是否以任一后缀结尾 |

#### 替换

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Replace` | `func(s, old, new string, n int) string` | 替换前n个匹配 |
| `ReplaceAll` | `func(s, old, new string) string` | 替换所有匹配 |
| `ReplaceWithIndex` | `func(s string, start, end int, new string) string` | 按索引范围替换 |
| `ReplaceWithMatcher` | `func(s string, matcher func(string) bool, new string) string` | 按匹配函数替换 |
| `ReplaceSpecialChars` | `func(s, new string) string` | 替换特殊字符 |

#### 其他操作

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Hide` | `func(s string, start, end int) string` | 隐藏中间字符（脱敏） |
| `Pad` | `func(s string, length int, pos PadPosition, char rune) string` | 按位置填充对齐 |
| `RuneWidth` | `func(s string) int` | 计算rune宽度 |
| `DisplayWidth` | `func(s string) int` | 计算显示宽度（支持CJK） |
| `Length` | `func(s string) int` | 计算rune长度 |
| `IndexOf` | `func(s, substr string) int` | 查找子串位置 |
| `NormalizeFieldName` | `func(s string) string` | 规范化字段名 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `StringX` | 链式字符串操作类型 |
| `PadPosition` | 填充位置常量类型 |
| `Paddler` | 填充器接口 |

### 常量/变量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `PadLeft` | PadPosition | 左填充 |
| `PadRight` | PadPosition | 右填充 |

### StringX 链式方法

`StringX` 支持以下链式方法（每个方法返回新的 `StringX`）：

`ToLowerChain`, `ToUpperChain`, `ToTitleChain`, `TrimChain`, `TrimStartChain`, `TrimEndChain`, `CleanEmptyChain`, `TrimProtocolChain`, `TrimAllChain`, `TrimAnyChain`, `TrimAllLineBreaksChain`, `TrimPrefixChain`, `TrimPrefixIgnoreCaseChain`, `TrimSuffixChain`, `TrimSuffixIgnoreCaseChain`, `TrimSymbolsChain`, `ReplaceChain`, `ReplaceAllChain`, `HideChain`, `PadChain`, `SubBeforeChain`, `SubAfterChain`, `SubBetweenChain`, `String()`

## 常用示例

详细用法参阅 → [reference.md](reference.md)

## 注意事项

- 链式方法返回新 `StringX`，不修改原对象，可安全复用
- `DisplayWidth` 对CJK字符计宽为2，ASCII计1，用于终端对齐
- `Hide` 的 start/end 参数为rune索引，非字节偏移
- `SubBetweenAll` 返回所有匹配，若无匹配返回空切片