---
name: validator-validation
description: 校验器工具，提供JSON Schema校验、结构体校验、JSON字节扫描、protobuf wrapper 解包、字符串/数值比较、IP黑白名单、路径匹配、正则缓存、空值检测。当需要校验JSON数据/结构体、扫描 JSON 字段、做IP访问控制、或匹配路径与正则时使用。
---

# validator - 校验器

提供JSON Schema校验、结构体校验、JSON 字节扫描、protobuf wrapper 解包、IP黑白名单、路径匹配、正则缓存与空值检测。

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
| `UnwrapProtobufWrapper` | `func(value interface{}) (interface{}, bool)` | 解包 protobuf wrapper，返回底层值 |
| `IsEmptyAfterDeref` | `func(value interface{}) (interface{}, bool)` | 解引用并判断过滤条件是否为空，支持 protobuf wrapper |

> `IsNil`、`IsCEmpty`、`IsFuncType`、`DerefValue` 已保留为兼容入口，新代码优先使用 `types` 包同名能力。

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

#### JSON字节扫描

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `IsJSONNull` | `func(data []byte) bool` | 判断去空白后是否为 JSON null，忽略大小写 |
| `SkipJSONSpaces` | `func(data []byte, i int) int` | 跳过 JSON 空白字符并返回下一个位置 |
| `ScanJSONString` | `func(data []byte, start int) (int, error)` | 扫描 JSON 字符串并返回结束后一位 |
| `ScanJSONValueEnd` | `func(data []byte, start int) (int, error)` | 扫描任意 JSON 值并返回结束后一位 |

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

#### 辅助函数（迁移说明）

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `types.Ptr` | `func(v T) *T` | 通用指针工具，替代旧的 `StringPtr`/`IntPtr`/`Float64Ptr`/`BoolPtr` |
| `types.GetReflectKind` | `func(value interface{}) reflect.Kind` | 获取反射 Kind |
| `types.IsNumericKind` | `func(kind reflect.Kind) bool` | 判断数值 Kind |
| `types.ToFloat64OK` | `func(value interface{}) (float64, bool)` | 数值转 float64 |

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
- JSON scanner 只负责定位 JSON 值边界，不替代完整 schema 校验
- 通用类型判断和反射工具已下沉到 `types`，validator 仅保留校验语义和兼容入口