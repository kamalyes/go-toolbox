---
name: validator-validation
description: 校验器工具，提供JSON Schema校验、结构体校验、字符串/数值比较、IP黑白名单、路径匹配、正则缓存、空值检测。当需要校验JSON数据/结构体、做IP访问控制、或匹配路径与正则时使用。
---

# validator - 校验器

提供JSON Schema校验、结构体校验、IP黑白名单、路径匹配、正则缓存与空值检测。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/validator"
```

JSON Schema校验：
```go
err := validator.ValidateJSON(jsonData, schema)
```

IP与路径匹配：
```go
allowed := validator.IsIPAllowed(ip, allowList)
matched := validator.MatchPathInList("/api/v1/users", patterns)
```

## 完整API索引

### 函数

#### 空值检测

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `IsEmptyValue` | `func(v interface{}) bool` | 检测值是否为空 |
| `IsEmptyPointer` | `func(v interface{}) bool` | 检测指针是否为nil |
| `IsEmptyStruct` | `func(v interface{}) bool` | 检测结构体是否为零值 |
| `IsUndefined` | `func(v interface{}) bool` | 检测值是否未定义 |
| `IsNull` | `func(v interface{}) bool` | 检测值是否为null |

#### 字符串比较

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `CompareStrings` | `func(actual, expect, op string) CompareResult` | 字符串比较操作 |
| `ValidateContains` | `func(body, pattern string) bool` | 验证包含子串 |
| `ValidateNotContains` | `func(body, pattern string) bool` | 验证不包含子串 |
| `ValidateRegex` | `func(body, pattern string) bool` | 验证正则匹配 |

#### JSON校验

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ValidateJSON` | `func(body, schema string) error` | 用JSON Schema校验JSON数据 |
| `ValidateJSONWithData` | `func(body, schema string, data interface{}) error` | 用JSON Schema校验JSON带数据 |
| `ValidateJSONField` | `func(body, field, schema string) bool` | 校验JSON单个字段 |
| `ValidateJSONFields` | `func(body string, fields map[string]string) bool` | 校验JSON多个字段 |
| `ValidateJSONPath` | `func(body, path, schema string) error` | 校验JSON指定路径 |
| `ValidateJSONPathExists` | `func(body, path string) bool` | 校验JSON路径是否存在 |
| `ValidateJSONSchema` | `func(schema string) error` | 校验JSON Schema本身合法性 |
| `ValidateStructWithSchema` | `func(v interface{}, schema string) error` | 用Schema校验结构体 |
| `NewSchemaBuilder` | `func() *SchemaBuilder` | 创建Schema构建器 |
| `QuickSchema` | `func(properties, required ...) string` | 快捷创建Schema |
| `FormatSchemaError` | `func(result interface{}) string` | 格式化Schema校验错误 |

#### 数值比较

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `CompareNumbers[T]` | `func(actual, expect T, op CompareOperator) CompareResult` | 泛型数值比较 |

#### HTTP校验

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ValidateStatusCode` | `func(statusCode int, expected int) bool` | 验证HTTP状态码 |
| `ValidateStatusCodeRange` | `func(statusCode, min, max int) bool` | 验证状态码范围 |
| `ValidateHeader` | `func(header, key, value string) bool` | 验证HTTP头部 |
| `ValidateContentType` | `func(header, contentType string) bool` | 验证Content-Type |

#### 路径与IP

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `MatchPathInList` | `func(path string, patterns []string) bool` | 路径模式匹配 |
| `MatchIPInList` | `func(ip string, list []string) bool` | IP匹配列表 |
| `IsIPInRange` | `func(ip, cidr string) bool` | IP是否在CIDR范围内 |
| `MatchIPWithWildcard` | `func(ip, pattern string) bool` | IP通配符匹配 |
| `IsIPAllowed` | `func(ip string, allowList []string) bool` | IP白名单判断 |
| `IsIPBlocked` | `func(ip string, blockList []string) bool` | IP黑名单判断 |
| `MatchIPPattern` | `func(ip, pattern string) bool` | IP模式匹配 |
| `IsPrivateIP` | `func(ip string) bool` | 判断是否私有IP |

#### 正则缓存

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GetCompiledRegex` | `func(pattern string) *regexp.Regexp` | 获取编译缓存的正则 |
| `ClearRegexCache` | `func()` | 清除正则缓存 |

#### 辅助函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `StringPtr` | `func(s string) *string` | 字符串指针工具 |
| `IntPtr` | `func(i int) *int` | 整数指针工具 |
| `Float64Ptr` | `func(f float64) *float64` | 浮点指针工具 |
| `BoolPtr` | `func(b bool) *bool` | 布尔指针工具 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `CompareOperator` | 比较操作符类型 |
| `CompareResult` | 比较结果类型，含 `Match` 和 `Message` 字段 |
| `IPBase` | IP基础类型 |
| `JSONSchema` | JSON Schema类型 |
| `SchemaBuilder` | Schema构建器类型 |

## 注意事项

- `GetCompiledRegex` 内部有缓存，相同正则只编译一次
- `ValidateJSONSchema` 仅校验schema格式，不校验数据
- `IsIPAllowed` 同时支持IPv4和IPv6格式