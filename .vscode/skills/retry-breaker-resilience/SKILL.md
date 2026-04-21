---
name: retry-breaker-resilience
description: 重试与熔断工具，提供可配置的重试策略（退避、最大次数、可重试错误判定）和熔断器（三态切换、半开探测）。当需要为不稳定操作添加重试逻辑、或需要熔断保护下游服务时使用。
---

# retry + breaker - 重试与熔断

提供链式配置的重试执行器和三态熔断器，构建弹性调用链路。

## 快速开始

```go
import (
    "github.com/kamalyes/go-toolbox/pkg/retry"
    "github.com/kamalyes/go-toolbox/pkg/breaker"
)
```

重试执行：
```go
err := retry.NewRetry().
    SetAttemptCount(3).
    SetInterval(100 * time.Millisecond).
    Do(func() error { return callAPI() })
```

熔断保护：
```go
b := breaker.New("service-a", breaker.Config{})
err := b.Execute(func() error { return callService() })
```

## 完整API索引

### retry 包

#### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewRetry` | `func() *Retry` | 创建重试执行器 |
| `NewRetryWithCtx` | `func(ctx context.Context) *Retry` | 创建带上下文的重试执行器 |

#### Retry 链式配置方法

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `SetAttemptCount` | `func(n int) *Retry` | 设置最大重试次数 |
| `SetInterval` | `func(d time.Duration) *Retry` | 设置重试间隔 |
| `SetMaxInterval` | `func(d time.Duration) *Retry` | 设置最大重试间隔 |
| `SetBackoffMultiplier` | `func(f float64) *Retry` | 设置退避乘数 |
| `SetJitter` | `func(d time.Duration) *Retry` | 设置抖动时间 |
| `SetJitterPercent` | `func(p float64) *Retry` | 设置抖动百分比 |
| `SetErrCallback` | `func(fn ErrCallbackFunc) *Retry` | 设置错误回调 |
| `SetSuccessCallback` | `func(fn SuccessCallbackFunc) *Retry` | 设置成功回调 |
| `SetConditionFunc` | `func(fn func(error) bool) *Retry` | 设置可重试条件判定 |
| `Do` | `func(fn DoFun) error` | 执行函数并按策略重试 |

#### 类型

| 导出名称 | 说明 |
|---|---|
| `Retry` | 重试执行器类型 |
| `DoFun` | 重试执行函数类型 |
| `ErrCallbackFunc` | 错误回调函数类型 |
| `SuccessCallbackFunc` | 成功回调函数类型 |

### breaker 包

#### 函数

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `New` | `func(name string, config Config) *Circuit` | 创建熔断器 |

#### Circuit 方法

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Execute` | `func(fn func() error) error` | 在熔断保护下执行函数 |
| `AllowRequest` | `func() bool` | 判断是否允许请求通过 |
| `RecordSuccess` | `func()` | 记录成功 |
| `RecordFailure` | `func()` | 记录失败 |
| `State` | `func() State` | 获取当前状态 |

#### 类型

| 导出名称 | 说明 |
|---|---|
| `Circuit` | 熔断器类型 |
| `Config` | 熔断器配置类型 |
| `State` | 熔断器状态类型 |

#### 常量/变量

| 导出名称 | 值/类型 | 说明 |
|---|---|---|
| `StateClosed` | State | 熔断器关闭（正常）状态 |
| `StateOpen` | State | 熔断器打开（熔断）状态 |
| `StateHalfOpen` | State | 熔断器半开（探测）状态 |
| `ErrOpen` | error | 熔断器打开错误 |

## 注意事项

- `retry.Do` 在不可重试错误时立即返回，仅对可重试错误重试
- `breaker.Execute` 在 Open 状态下直接返回 `ErrOpen`，不会调用fn
- 熔断器状态切换依赖 `RecordSuccess/RecordFailure`，务必在业务代码中正确调用