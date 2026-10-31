---
name: cron-expression
description: Cron 表达式解析与下次触发时间计算，支持 5 字段标准式与 6 字段带秒式、@daily/@every 等描述符、预置别名表、位掩码字段解析与自定义字段组合；当需要解析 cron 表达式或算下一次执行时刻时使用
---

# cron - Cron 表达式解析

把 cron 表达式解析成 `CronSchedule`，再用 `Next(t)` 算下一次触发时刻

## 这个包只解析，不跑任务

**它没有调度循环**：没有 `Start`/`Stop`/`AddJob`，不起 goroutine，不管重入与错过补偿
它给你的是「下一次该在什么时候执行」，循环要自己写：

```go
schedule, err := cron.ParseCronStandard("0 2 * * *")
if err != nil {
    return err
}
for {
    next := schedule.Next(time.Now())
    select {
    case <-time.After(time.Until(next)):
        doJob()
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

好处是零调度器状态、容易测试；代价是漏执行补偿、并发互斥、panic 恢复都得自己做

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/cron"
```

```go
// 5 字段标准式：分 时 日 月 周
s, _ := cron.ParseCronStandard("0 2 * * *")        // 每天 02:00

// 6 字段带秒式：秒 分 时 日 月 周
s, _ = cron.ParseCronWithSeconds("0 0 2 * * *")

// 描述符
s, _ = cron.ParseCronStandard("@daily")
s, _ = cron.ParseCronStandard("@every 1h30m")      // 固定间隔

// 自定义字段组合
p := cron.NewCronParser(cron.CronSecond | cron.CronMinute | cron.CronHour |
                        cron.CronDom | cron.CronMonth | cron.CronDow)
s, _ = p.Parse("*/5 * * * * *")

next := s.Next(time.Now())
```

## 完整API索引

### 解析入口

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ParseCronStandard` | `func(spec string) (CronSchedule, error)` | 5 字段（分 时 日 月 周），也认描述符 |
| `ParseCronWithSeconds` | `func(spec string) (CronSchedule, error)` | 6 字段（秒 分 时 日 月 周） |
| `NewCronParser` | `func(options CronParseOption) *CronParser` | 自定义字段组合的解析器 |
| `Parse` | `func(spec string) (CronSchedule, error)` | `*CronParser` 的方法 |

### CronParseOption 位标志

`CronSecond`、`CronSecondOptional`、`CronMinute`、`CronHour`、`CronDom`、`CronMonth`、
`CronDow`、`CronDowOptional`、`CronDescriptor`

按位或组合，决定表达式有哪些字段。预置组合：`CronStandardParser`、`CronSecondParser`

### 描述符与别名

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GetExpression` | `func(alias string) (string, bool)` | 按别名取表达式，第二个返回值表示是否命中 |
| `ExpressionAliases` | `map` | 别名到表达式的映射表 |

支持的描述符：`@yearly`、`@annually`、`@monthly`、`@weekly`、`@daily`、`@midnight`、
`@hourly`，以及 `@every <duration>`

预置表达式常量：`EverySecond`、`Every10Seconds`、`Every30Seconds`、`EveryMinute`、
`Every5Minutes`、`Every15Minutes`、`Every30Minutes`、`EveryHour`、`Every2Hours`、
`Every3Hours`、`Every6Hours`、`Every12Hours`、`EveryDay`、`EveryWeek`、`EveryMonth`、`EveryYear`

### 调度类型

| 导出名称 | 说明 |
|---|---|
| `CronSchedule` | 调度接口，核心方法 `Next(t time.Time) time.Time` |
| `CronSpecSchedule` | 按字段位掩码匹配的调度（普通 cron 表达式） |
| `CronEverySchedule` | 固定间隔调度（`@every`） |

**`CronSpecSchedule` 的构造与链式设置**（用于不走文本解析、直接按位构造）：

| 导出名称 | 签名 |
|---|---|
| `NewZeroCronSpecSchedule` | `func(loc *time.Location) *CronSpecSchedule` |
| `WithSecond` / `WithMinute` / `WithHour` / `WithDom` / `WithMonth` / `WithDow` | `func(v uint64) *CronSpecSchedule` |

每个字段是一个 `uint64` 位掩码，第 n 位为 1 表示该取值命中

### 字段解析辅助

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ParseExprToBits[T]` | `func(expr string, bounds types.Bounds[T], starBit uint64) (uint64, error)` | 整个表达式转位掩码 |
| `ParseFieldToBits[T]` | `func(field string, bounds types.Bounds[T], starBit uint64) (uint64, error)` | 单字段转位掩码 |
| `ParseFieldWithSpecialChars` | `func(fieldStr string, bounds cronBounds, starBit uint64, isDOM, isDOW bool) (...)` | 处理 `*` `?` `L` 等特殊字符 |
| `ParseIntOrName[T]` | `func(expr string, names map[string]T) (T, error)` | 数字或名称（如 `MON`、`JAN`）转值 |

泛型参数受 `types.Numerical` 约束，边界用 `types.Bounds[T]` 描述

## 注意事项

- **⚠️ 日(`Dom`) 与 周(`Dow`) 同时限定时本包用「与」，这与传统 cron 不同。**
  Vixie cron 与 robfig/cron 在两个字段都被限定时取**或**，本包取**与**
  实测 `ParseCronStandard("0 0 1 * MON")` 从 2026-05-20 起算，
  下一次是 **2026-06-01**（1 号且恰为周一），而不是 2026-05-25（下一个周一）：

  | `Dom` | `Dow` | 本包行为 |
  |---|---|---|
  | `*` | `*` | 每天 |
  | `*` | 限定 | 只看 `Dow` |
  | 限定 | `*` | 只看 `Dom` |
  | 限定 | 限定 | **两者都要满足（与）** |

  从 crontab 或 robfig 迁移表达式时，凡是两个字段都填了具体值的，触发频率会显著变低，
  必须逐条用 `Next` 对照过一遍
  判定逻辑在 `CronSpecSchedule.matchDayOfMonthAndWeek`，末尾是 `domMatch && dowMatch`
- **时区取自 `CronSpecSchedule` 持有的 `*time.Location`**，不是调用 `Next` 时传入时间的时区
  用 `NewZeroCronSpecSchedule(loc)` 显式指定；跨夏令时的「每天 02:30」在缺失的那一小时上
  行为取决于该时区的规则，关键任务不要挑 DST 切换窗口
- `Next` 是**纯函数**：同一个 `schedule` 和同一个 `t` 永远得到同样结果，不记录状态
  所以它不知道「上次有没有执行过」，漏执行补偿要调用方自己判断
- `@every 1h30m` 走 `CronEverySchedule`，是从**上一次 `Next` 的入参**起算固定间隔，
  与按墙上时钟对齐的 cron 字段语义不同（`0 */2 * * *` 会对齐到偶数小时整点，`@every 2h` 不会）
- 本包**不依赖任何三方 cron 库**（如 robfig/cron），位掩码与解析都是自带实现；
  行为与 robfig 高度相似但不保证逐字节一致，迁移时请对关键表达式做一次 `Next` 对照
