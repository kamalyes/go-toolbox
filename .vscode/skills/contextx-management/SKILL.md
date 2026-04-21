---
name: contextx-management
description: 上下文管理工具，提供可取消/超时/值的上下文创建、合并、分离超时、装饰器链、泛型取值。当需要创建带超时/值的上下文、合并多个上下文、或安全提取泛型值时使用。
---

# contextx - 上下文管理

提供上下文创建、合并、分离超时、装饰器链与泛型取值。

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/contextx"
```

创建上下文：
```go
ctx := contextx.NewContext()
ctx := contextx.NewContextWithTimeout(5 * time.Second)
ctx := contextx.NewContextWithValue("key", "value")
```

泛型取值：
```go
val := contextx.MustGet[string](ctx, "key")
val := contextx.GetOrDefault[string](ctx, "key", "default")
```

## 完整API索引

### 函数

#### 上下文创建

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewContext` | `func() *Context` | 创建新上下文 |
| `NewContextWithTimeout` | `func(timeout time.Duration) *Context` | 创建带超时的上下文 |
| `NewContextWithValue` | `func(key, val interface{}) *Context` | 创建带值的上下文 |
| `IsContext` | `func(ctx interface{}) bool` | 判断是否为Context类型 |

#### 上下文操作

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `WithValue` | `func(ctx context.Context, key, val interface{}) context.Context` | 向上下文添加值 |
| `MergeContext` | `func(ctxs ...context.Context) context.Context` | 合并多个上下文 |
| `WithTimeout` | `func(timeout time.Duration, fn func()) error` | 带超时执行函数 |
| `WithTimeoutValue[T]` | `func(timeout time.Duration, fn func() T) T` | 带超时执行并返回值 |
| `OrBackground` | `func(ctx context.Context) context.Context` | 空上下文返回Background |
| `OrWithoutCancel` | `func(ctx context.Context) context.Context` | 去除取消功能的上下文 |

#### 分离超时

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `WithTimeoutFrom` | `func(parent context.Context, timeout time.Duration, fn func()) error` | 从父上下文继承超时 |
| `WithTimeoutOrBackground` | `func(parent context.Context, timeout time.Duration, fn func()) error` | 带超时或Background |
| `NewDetachedTimeout` | `func(parent context.Context, timeout time.Duration) context.Context` | 创建分离超时上下文 |
| `WithDetachedTimeout` | `func(parent context.Context, timeout time.Duration, fn func()) error` | 使用分离超时执行 |

#### 装饰器

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `WithTimeoutDecorators` | `func(timeout time.Duration, decorators ...) context.Context` | 带超时的装饰器链 |
| `WithDeadlineDecorators` | `func(deadline time.Time, decorators ...) context.Context` | 带截止时间的装饰器链 |

#### 泛型取值

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `MustGet[T]` | `func(ctx context.Context, key interface{}) T` | 必须获取泛型值，不存在则panic |
| `MustGetWithMessage[T]` | `func(ctx context.Context, key interface{}, message string) T` | 必须获取泛型值，不存在则panic带消息 |
| `GetOrDefault[T]` | `func(ctx context.Context, key interface{}, defaultVal T) T` | 获取泛型值，不存在返回默认值 |
| `Get[T]` | `func(c context.Context, key interface{}) (T, bool)` | 获取泛型值 |
| `GetValue[T]` | `func(ctx context.Context, key interface{}) (T, bool)` | 获取泛型值 |

### 类型

| 导出名称 | 说明 |
|---|---|
| `Context` | 增强上下文类型 |
| `ContextKey` | 上下文键类型 |
| `SourceFunc` | 来源函数类型 |
| `ProcessFunc` | 处理函数类型 |
| `MetadataAdapter` | 元数据适配器类型 |
| `Marshaler` | 序列化接口类型 |

### Context 方法

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `WithParent` | `func(parent context.Context) *Context` | 设置父上下文 |
| `WithPool` | `func(pool interface{}) *Context` | 设置连接池 |
| `WithCancel` | `func() (*Context, context.CancelFunc)` | 创建可取消上下文 |
| `WithTimeout` | `func(timeout time.Duration) (*Context, context.CancelFunc)` | 创建带超时的可取消上下文 |
| `WithDeadline` | `func(deadline time.Time) (*Context, context.CancelFunc)` | 创建带截止时间的可取消上下文 |

## 注意事项

- `MustGet` 在值不存在时panic，如需安全获取使用 `GetOrDefault`
- `MergeContext` 合并后任一上下文取消则合并上下文取消
- `NewDetachedTimeout` 创建的上下文不受父上下文取消影响