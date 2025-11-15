# Matcher - 通用规则匹配引擎

## 概述

`matcher` 包提供了一个优雅的规则匹配引擎，用于替代复杂的 if-for 嵌套逻辑。

## 核心特性

- ✅ **声明式规则定义** - 使用链式API定义规则
- ✅ **优先级支持** - 自动按优先级排序匹配
- ✅ **类型安全** - 泛型支持，编译时类型检查
- ✅ **高性能** - 避免不必要的计算
- ✅ **易扩展** - 支持自定义条件和规则

## 快速开始

### 基本用法

```go
import "github.com/kamalyes/go-toolbox/pkg/matcher"

// 定义结果类型
type Action struct {
    Name string
    Code int
}

// 创建匹配器
m := matcher.NewMatcher[*Action]()

// 添加规则
m.AddRule(
    matcher.NewChainRule(&Action{Name: "admin", Code: 1}).
        When(matcher.MatchString("role", "admin")).
        WithPriority(100),
)

m.AddRule(
    matcher.NewChainRule(&Action{Name: "user", Code: 2}).
        When(matcher.MatchString("role", "user")).
        WithPriority(50),
)

// 执行匹配
ctx := matcher.NewContext().Set("role", "admin")
if result, ok := m.Match(ctx); ok {
    fmt.Printf("Matched: %s\n", result.Name)
}
```

### 复杂条件组合

```go
// 路由匹配 + 方法匹配 + IP白名单
m.AddRule(
    matcher.NewChainRule(myAction).
        When(matcher.MatchPattern("path", "/api/*")).
        When(matcher.MatchMethodIn([]string{"GET", "POST"})).
        When(matcher.MatchStringNotIn("ip", blacklist)).
        WithPriority(100),
)

// 多条件OR
m.AddRule(
    matcher.NewChainRule(myAction).
        When(matcher.MatchAny(
            matcher.MatchPrefix("path", "/public"),
            matcher.MatchString("auth", "none"),
        )).
        WithPriority(50),
)

// 多条件AND
m.AddRule(
    matcher.NewChainRule(myAction).
        When(matcher.MatchAll(
            matcher.MatchPattern("path", "/admin/*"),
            matcher.MatchString("role", "admin"),
            matcher.MatchBool("verified", true),
        )).
        WithPriority(200),
)
```

### 自定义条件

```go
// 自定义匹配逻辑
customCondition := func(ctx *matcher.Context) bool {
    userID := ctx.GetString("user_id")
    return len(userID) > 0 && userID[0] == 'V'
}

m.AddRule(
    matcher.NewChainRule(vipAction).
        When(customCondition).
        WithPriority(150),
)
```

## API 文档

### Context 上下文

```go
ctx := matcher.NewContext()
ctx.Set("key", "value")           // 设置值
val, ok := ctx.Get("key")         // 获取值
str := ctx.GetString("key")       // 获取字符串
slice := ctx.GetStringSlice("key") // 获取字符串切片
bool := ctx.GetBool("key")        // 获取布尔值
```

### Matcher 匹配器

```go
m := matcher.NewMatcher[T]()      // 创建匹配器
m.AddRule(rule)                   // 添加单个规则
m.AddRules(rule1, rule2, ...)     // 批量添加规则
result, ok := m.Match(ctx)        // 匹配第一个
results := m.MatchAll(ctx)        // 匹配所有
```

### ChainRule 链式规则

```go
rule := matcher.NewChainRule(result).
    When(condition1).                 // 添加条件
    When(condition2).                 // 链式添加
    WithPriority(100)                 // 设置优先级
```

### 内置条件函数

| 函数 | 说明 |
|-----|------|
| `MatchString(key, value)` | 字符串精确匹配 |
| `MatchStringIn(key, list)` | 字符串在列表中 |
| `MatchStringNotIn(key, list)` | 字符串不在列表中 |
| `MatchPattern(key, pattern)` | 路径模式匹配 |
| `MatchPrefix(key, prefix)` | 前缀匹配 |
| `MatchSuffix(key, suffix)` | 后缀匹配 |
| `MatchContains(key, substring)` | 包含匹配 |
| `MatchBool(key, expected)` | 布尔值匹配 |
| `MatchAny(conditions...)` | 任意条件满足 |
| `MatchAll(conditions...)` | 所有条件满足 |
| `MatchNot(condition)` | 条件取反 |
| `MatchMethodIn(methods)` | HTTP方法匹配 |
| `MatchWildcard(key, pattern)` | 通配符匹配 |

## 使用场景

1. **限流规则匹配** - 替代复杂的路由/IP/用户规则嵌套
2. **权限验证** - 基于角色/资源/操作的权限规则
3. **路由分发** - 根据请求特征分发到不同处理器
4. **配置选择** - 根据环境/场景选择不同配置
5. **策略模式** - 实现灵活的策略选择逻辑

## 性能优化

- 规则按优先级自动排序，高优先级规则先匹配
- 短路求值，匹配成功立即返回
- 避免不必要的类型断言和反射
- 支持批量规则操作

## 最佳实践

1. **优先级设计** - 白名单 > 黑名单 > 特定规则 > 通用规则
2. **条件顺序** - 快速失败的条件放在前面
3. **复用规则** - 相同规则可以复用，避免重复定义
4. **上下文预填充** - 一次性设置所有需要的上下文数据
