---
name: moment-datetime
description: 日期时间工具，提供 dayjs 风格的 Moment 链式值类型（StartOf/EndOf/Add/Sub/Diff）、紧凑格式化与解析、柔性日期解析、时区换算、工作日与月末计算、年龄计算与计时器；当需要日期截断、区间计算、多格式日期解析或耗时统计时使用
---

# moment - 日期时间

两套 API：**`Moment` 链式值类型**（新代码优先用这个）和一批以 `time.Time` 为入参的独立函数

## Moment 链式值类型

参照 dayjs 的手感，但是值类型：每次调用返回新的 `Moment`，不改原值，可安全并发读

```go
import "github.com/kamalyes/go-toolbox/pkg/moment"
```

```go
m := moment.Now()
monthStart := m.StartOf(moment.UnitMonth)              // 本月第一天 00:00:00
monthEnd   := m.EndOf(moment.UnitMonth)                // 本月最后一天的最后一纳秒
nextWeek   := m.Add(1, moment.UnitWeek)
days       := m.Diff(other, moment.UnitDay)

// 链式，零分配
label := moment.Now().StartOf(moment.UnitDay).Add(-1, moment.UnitDay).Format("2006-01-02")
```

### 构造

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Now` | `func() Moment` | 当前时间 |
| `From` | `func(t time.Time) Moment` | 包装已有 `time.Time` |
| `FromUnix` | `func(sec, nsec int64) Moment` | 由 Unix 秒 + 纳秒构造 |

### 方法

| 方法 | 签名 | 说明 |
|---|---|---|
| `StartOf` | `func(u Unit) Moment` | 截断到该单位起点 |
| `EndOf` | `func(u Unit) Moment` | 该单位的**最后一纳秒** |
| `StartOfWeekOn` | `func(start time.Weekday) Moment` | 自定周起点的周起点 |
| `EndOfWeekOn` | `func(start time.Weekday) Moment` | 自定周起点的周结束 |
| `Add` / `Sub` | `func(n int, u Unit) Moment` | 按单位加减；`Sub` 等于 `Add(-n, u)` |
| `Diff` | `func(other Moment, u Unit) int` | 相差多少个该单位 |
| `Before` / `After` / `Equal` | `func(other Moment) bool` | 比较 |
| `Format` | `func(layout string) string` | 按 Go layout 格式化 |
| `Time` | `func() time.Time` | 取出底层 `time.Time` |
| `In` | `func(loc *time.Location) Moment` | 换时区 |
| `Unix` / `UnixMilli` | `func() int64` | 时间戳 |
| `Year` `Month` `Day` `Hour` `Minute` `Second` `Nanosecond` `Weekday` `Quarter` | — | 分量取值 |
| `IsZero` | `func() bool` | 是否零值 |
| `String` | `func() string` | 默认字符串表示 |

### Unit 枚举

`UnitNanosecond`、`UnitMicrosecond`、`UnitMillisecond`、`UnitSecond`、`UnitMinute`、
`UnitHour`、`UnitDay`、`UnitWeek`、`UnitMonth`、`UnitQuarter`、`UnitYear`

## 四处刻意与 dayjs 不同，改之前先读这里

1. **月末加减按钳位，不按溢出。** 1 月 31 日加一个月得到 **2 月 28 日**，
   而 Go 自带的 `AddDate` 会溢出归一成 3 月 3 日。钳位更合直觉，dayjs 也是这个语义
2. **`StartOf(UnitWeek)` 的周起点是周一**（ISO 8601），dayjs 默认是周日
   需要周日起点请用 `StartOfWeekOn(time.Sunday)`，不要去改默认值
3. **`EndOf` 返回最后一纳秒**，dayjs 只到毫秒。Go 的精度是纳秒，这里不自降精度
   实现是「下一单位的起点减一纳秒」，月末、闰年、夏令时都自然正确
4. **天与周的加减走 `AddDate`，不是乘 24 小时。** 一天并不总是 24 小时——
   夏令时切换当天是 23 或 25 小时。`AddDate` 按墙上时钟走，乘 `24*time.Hour` 会错一小时

另外 **`Diff` 的反对称性是显式保证的**：月末钳位会破坏
`Diff(a,b) == -Diff(b,a)`，所以实现只定义一个方向，另一个方向取负

## 独立函数（以 time.Time 为入参）

### 柔性与紧凑解析

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `ParseFlexibleDate` | `func(dateStr string) (time.Time, error)` | 多种常见日期写法自动识别 |
| `ParseCompact` | `func(value, format string) (time.Time, error)` | 解析紧凑格式 |
| `StrtoTime` | `func(s string, format ...string) (time.Time, error)` | 字符串转时间，可指定格式 |
| `SafeParseTimeToUnixNano` | `func(timeStr string) int64` | 解析失败返回 0 而不报错 |
| `ConvertStringToTimestamp` | `func(dateString, layout, timeZone string) (int64, error)` | 按时区解析为时间戳 |
| `ParseWithLocation` | `func(timezone, timeStr, defaultDateFormat string) int64` | 按时区解析 |

`ParseFlexibleDate` 先试最常见的那个 layout，不中再用字符类掩码过滤候选，
避免对每个候选 layout 都调一次 `time.Parse`。紧凑日期（如 `20260131`）上这条快路径效果最明显

### 紧凑格式化

| 导出名称 | 说明 |
|---|---|
| `FormatCompact` / `FormatCompactDate` / `FormatCompactDateHour` / `FormatCompactDateTime` / `FormatCompactDateTimeSec` / `FormatCompactDateTimeMilli` | 对给定时间输出紧凑串 |
| `NowCompactDate` / `NowCompactDateHour` / `NowCompactDateTime` / `NowCompactDateTimeSec` / `NowCompactDateTimeMilli` | 对当前时间输出紧凑串 |
| `FormatWithLocation` | `func(timezone string, ts int64, defaultDateFormat string) string` |
| `CharToCode` | `func(layout string) string` | layout 字符转换 |

### 分量与区间

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `Year` `Month` `Day` `Hour` `Minute` `Second` `YearDay` | `func(t ...time.Time) int` | 不传参数即取当前时间 |
| `NowTime` | `func(t []time.Time) time.Time` | 切片为空取当前时间 |
| `Today` / `Yesterday` | `func() (int64, int64)` | 起止时间戳 |
| `LastNDays` / `LastNWeeks` / `LastNMonths` / `LastNYears` | `func(num int) (int64, int64)` | 近 N 个单位的起止时间戳 |
| `CalculateStartAndEndTime` | `func(year int, month time.Month, day int, duration time.Duration) (int64, int64)` | 指定日期加时长的起止 |

### 日历计算

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `DaysInMonth` | `func(year int, month time.Month) int` | 该月天数（含闰年） |
| `LastDayOfMonth` | `func(year int, month time.Month) time.Time` | 月末那天 |
| `DayOfYear` | `func(year, n int) time.Time` | 一年中的第 n 天 |
| `NextWeekday` | `func(year int, month time.Month, day int, target time.Weekday) time.Time` | 下一个指定星期几 |
| `LastWeekdayOfMonth` | `func(year int, month time.Month, target time.Weekday) time.Time` | 该月最后一个指定星期几 |
| `NextWorkDay` | `func(year int, month time.Month, day int) time.Time` | 下一个工作日 |

### 时区与年龄

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `GetServerTimezone` | `func() string` | 服务器时区 |
| `GetTimeOffset` | `func(timezone string, ts int64) int` | 时区偏移秒数 |
| `GetCurrentTimeInfo` | `func() (string, int, time.Time)` | 时区名、偏移、当前时间 |
| `CalculateAge` | `func(birthday string, currentTime time.Time) (int, error)` | 按生日算年龄 |
| `CalculateBirthDate` | `func(age int) time.Time` | 按年龄反推生日 |

### 耗时统计

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `TrackTime` | `func(pre time.Time) time.Duration` | 距 `pre` 过去多久 |
| `HumanDuration` | `func(start, end time.Time) string` | 人类可读的时长 |
| `CalculateTimeDifference` | `func(duration time.Duration) TimeDifference` | 拆成天/时/分/秒 |
| `CurrentNanosecond` / `CurrentMicrosecond` / `CurrentMillisecond` | `func() int64` | 当前时间戳 |

### Timer

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewTimer` | `func() *Timer` | 创建计时器 |
| `NewTimerWithTraceId` | `func(traceId string) *Timer` | 带 traceId 的计时器 |

**Timer 方法**：`Run`、`Pause`、`Resume`、`Finish`、`SetTimeOut`、`SetTraceId`、
`SetLogPrefix`、`GetStartTime`、`GetEndTime`、`GetDuration`、`GetPauseDuration`、
`GetPaused`、`GetTimeOut`、`GetTraceId`；`Set*` 返回 `*Timer` 支持链式

### 类型

| 导出名称 | 说明 |
|---|---|
| `Moment` | 链式日期时间值类型（值语义，方法不改原值） |
| `Unit` | 时间单位枚举（`uint8`） |
| `Timer` | 计时器，支持暂停与恢复 |
| `TimerInterface` | 计时器接口 |
| `TimeDifference` | 时长拆解结果 |

## 注意事项

- **`Moment` 是值类型**，所有方法都返回新值，原值不变；可以安全地在多个 goroutine 里读同一个
  `Moment`，但它不是「可变的时间对象」，`m.Add(1, UnitDay)` 不会改变 `m`
- **跨夏令时的日期加减不要自己乘 `24*time.Hour`**，用 `Add(n, UnitDay)` 或 `AddDate`
- `EndOf` 得到的是**最后一纳秒**，用它做右边界时区间是闭区间；
  如果下游期望半开区间 `[start, next)`，请直接用 `StartOf(u).Add(1, u)` 而不是 `EndOf`
- `Year`/`Month`/`Day` 这类**变参函数不传参数就取当前时间**，
  这让「忘记传参」不会编译报错而是静默取当前时间，调用处注意
- `SafeParseTimeToUnixNano` 解析失败返回 **0**，与「1970-01-01」无法区分，
  需要区分失败请用 `StrtoTime` 并检查 error
- `ParseFlexibleDate` 的候选 layout 表与字符类掩码是配套的：
  加新格式时两处都要改，否则新格式会被掩码提前过滤掉而永远匹配不上
  包内保留了一个朴素实现 `parseFlexibleDateReference` 专门用于差分测试，不要删
