---
name: matcher-routing
description: 路由匹配工具，提供字符串匹配（精确/包含/前缀/后缀/正则/通配符）、路径匹配、链式规则匹配。当需要做URL路由匹配、请求路径分发、或多条件组合匹配时使用。
---

# matcher - 路由匹配

提供字符串匹配、路径匹配与链式规则匹配。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/matcher"
```

字符串匹配：
```go
m := matcher.NewMatcher[string]()
m.AddRule(matcher.NewChainRule("admin", matcher.MatchPrefix("admin")))
result, ok := m.Match("/admin/dashboard")
```

路径匹配：
```go
pm := matcher.NewPathMatcher(matcher.PathMatcherTypeGlob, "/api/v1/*")
matched := pm.Match("/api/v1/users")
```

## 完整API索引

### 函数

#### 匹配器构建

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewMatcher[T]` | `func() *Matcher[T]` | 创建泛型匹配器 |
| `NewChainRule[T]` | `func(result T) *ChainRule[T]` | 创建链式规则 |

#### 字符串匹配

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `MatchString` | `func(pattern, target string) bool` | 精确字符串匹配 |
| `MatchStringIn` | `func(target string, patterns []string) bool` | 字符串在列表中匹配 |
| `MatchStringNotIn` | `func(target string, patterns []string) bool` | 字符串不在列表中 |
| `MatchStringInCaseInsensitive` | `func(target string, patterns []string) bool` | 忽略大小写在列表中匹配 |
| `MatchPattern` | `func(pattern, target string) bool` | 正则模式匹配 |
| `MatchPrefix` | `func(prefix string) func(string) bool` | 前缀匹配 |
| `MatchSuffix` | `func(suffix string) func(string) bool` | 后缀匹配 |
| `MatchContains` | `func(substr string) func(string) bool` | 包含匹配 |
| `MatchBool` | `func(cond bool) func(string) bool` | 布尔条件匹配 |
| `MatchAny` | `func(matchers ...func(string) bool) func(string) bool` | 任一匹配 |
| `MatchAll` | `func(matchers ...func(string) bool) func(string) bool` | 全部匹配 |
| `MatchNot` | `func(matcher func(string) bool) func(string) bool` | 取反匹配 |
| `MatchMethodIn` | `func(method string, methods []string) bool` | HTTP方法匹配 |
| `MatchWildcard` | `func(pattern, target string) bool` | 通配符匹配 |

#### 路径匹配

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewPathMatcher` | `func(matchType PathMatcherType, pattern string) PathMatcher` | 创建路径匹配器 |
| `MatchPathGlob` | `func(pattern, path string) bool` | 路径glob匹配 |
| `MatchPathWithMethod` | `func(method, path string, patterns map[string]string) bool` | 带HTTP方法的路径匹配 |
| `MatchMethod` | `func(method string, allowed string) bool` | HTTP方法匹配 |
| `NormalizePath` | `func(path string) string` | 规范化路径 |
| `ExtractPathSegments` | `func(path string) []string` | 提取路径段 |
| `NewPathMatcherBuilder` | `func() *PathMatcherBuilder` | 创建路径匹配器构建器 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `Rule[T]` | 规则接口类型 |
| `Matcher[T]` | 泛型匹配器类型 |
| `MatcherStats` | 匹配器统计类型 |
| `MatchMiddleware[T]` | 匹配中间件类型 |
| `ChainRule[T]` | 链式规则类型 |
| `PathMatcherType` | 路径匹配器类型枚举 |
| `PathMatcher` | 路径匹配器类型 |
| `PathMatcherBuilder` | 路径匹配器构建器类型 |

## 注意事项

- `MatchPattern` 使用正则表达式，高频场景建议预编译
- `MatchWildcard` 支持 `*` 通配符，不支持 `**`
- `NewPathMatcherBuilder` 支持链式配置多个路径规则