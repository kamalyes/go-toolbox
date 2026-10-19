/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-18 22:11:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-18 23:01:58
 * @FilePath: \go-toolbox\pkg\moment\moment.go
 * @Description: 链式时间值类型，按单位做加减与取边界
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package moment

import (
	"time"
)

// ============================================================================
// Moment 链式时间值类型
//
// 借 dayjs 的用法，但不照搬它的实现
//
// dayjs 的两个招牌是不可变与链式，而 Go 的 time.Time 本来就不可变、本来就能链，
// 所以真正缺的只是「按单位加减」与「取单位边界」这两组操作
//
// 做成值类型而不是指针：Moment 就是一个 time.Time 的壳（24 字节），整条链上不产生
// 任何堆分配dayjs 每次操作都要 new 一个对象，这一点 Go 这边能做得更好
// 单位用类型化常量而不是字符串：dayjs 要在运行期解析 'day' / 'days' / 'd'，
// 而这里写错单位编译期就报错，也没有运行期的字符串比较
//
// 与 dayjs 的语义差异，都是有意选择：
//
//   - 月末钳位：1 月 31 日加一个月得到 2 月 28 日（夏季 29 日），与 dayjs 一致
//     Go 自带的 AddDate 在这里会溢出归一，得到 3 月 3 日钳位更合直觉
//   - 周起始取周一（ISO 8601，也是中文习惯），dayjs 默认周日
//     需要别的起点用 StartOfWeekOn
//   - EndOf 返回该单位的最后一纳秒，dayjs 到毫秒为止——Go 的精度更高，不必自降
//   - 加减天与周走 AddDate 而不是乘 24 小时：有夏令时的时区里一天不总是 24 小时
//
// [EN] A chainable time value type in the spirit of dayjs, without copying its
// implementation. dayjs's headline features are immutability and chaining, both
// of which Go's time.Time already has, so what is actually missing is arithmetic
// by unit and unit boundaries.
//
// It is a value type rather than a pointer: a Moment is a 24-byte shell around a
// time.Time, so a whole chain allocates nothing on the heap, where dayjs
// allocates an object per operation. Units are typed constants rather than
// strings: dayjs parses 'day' / 'days' / 'd' at run time, while a wrong unit
// here fails to compile and costs no run-time string comparison.
//
// The deliberate differences from dayjs: month arithmetic clamps to the end of
// the target month, matching dayjs where Go's own AddDate would overflow into
// the next month; the week starts on Monday per ISO 8601 rather than dayjs's
// Sunday, with StartOfWeekOn for other choices; EndOf returns the unit's last
// nanosecond rather than stopping at milliseconds; and day and week arithmetic
// goes through AddDate instead of multiplying by 24 hours, because a day is not
// always 24 hours in a zone that observes daylight saving
// ============================================================================

// Unit 时间单位
type Unit uint8

// 支持的时间单位
const (
	UnitNanosecond Unit = iota
	UnitMicrosecond
	UnitMillisecond
	UnitSecond
	UnitMinute
	UnitHour
	UnitDay
	UnitWeek
	UnitMonth
	UnitQuarter
	UnitYear
)

// String 返回单位名称
func (u Unit) String() string {
	switch u {
	case UnitNanosecond:
		return "nanosecond"
	case UnitMicrosecond:
		return "microsecond"
	case UnitMillisecond:
		return "millisecond"
	case UnitSecond:
		return "second"
	case UnitMinute:
		return "minute"
	case UnitHour:
		return "hour"
	case UnitDay:
		return "day"
	case UnitWeek:
		return "week"
	case UnitMonth:
		return "month"
	case UnitQuarter:
		return "quarter"
	case UnitYear:
		return "year"
	default:
		return "unknown"
	}
}

// Moment 包装 time.Time 的链式时间值
//
// 值类型，可直接复制；零值等价于 time.Time 的零值
type Moment struct {
	t time.Time
}

// Now 返回当前时刻
func Now() Moment { return Moment{t: time.Now()} }

// From 由 time.Time 构造
func From(t time.Time) Moment { return Moment{t: t} }

// FromUnix 由 Unix 秒与纳秒构造
func FromUnix(sec, nsec int64) Moment { return Moment{t: time.Unix(sec, nsec)} }

// Time 取出底层的 time.Time
func (m Moment) Time() time.Time { return m.t }

// IsZero 是否为零值时刻
func (m Moment) IsZero() bool { return m.t.IsZero() }

// Format 按 layout 格式化
func (m Moment) Format(layout string) string { return m.t.Format(layout) }

// String 按 DefaultTimeFormat 格式化，便于打印
func (m Moment) String() string { return m.t.Format(DefaultTimeFormat) }

// In 换到指定时区，时刻本身不变
func (m Moment) In(loc *time.Location) Moment { return Moment{t: m.t.In(loc)} }

// Unix 返回 Unix 秒
func (m Moment) Unix() int64 { return m.t.Unix() }

// UnixMilli 返回 Unix 毫秒
func (m Moment) UnixMilli() int64 { return m.t.UnixMilli() }

// 日期与时间分量
func (m Moment) Year() int             { return m.t.Year() }
func (m Moment) Month() time.Month     { return m.t.Month() }
func (m Moment) Day() int              { return m.t.Day() }
func (m Moment) Hour() int             { return m.t.Hour() }
func (m Moment) Minute() int           { return m.t.Minute() }
func (m Moment) Second() int           { return m.t.Second() }
func (m Moment) Nanosecond() int       { return m.t.Nanosecond() }
func (m Moment) Weekday() time.Weekday { return m.t.Weekday() }

// Quarter 返回季度，1 到 4
func (m Moment) Quarter() int { return (int(m.t.Month())-1)/3 + 1 }

// 比较
func (m Moment) Before(other Moment) bool { return m.t.Before(other.t) }
func (m Moment) After(other Moment) bool  { return m.t.After(other.t) }
func (m Moment) Equal(other Moment) bool  { return m.t.Equal(other.t) }

// Add 按单位加 n，n 可为负
//
// 天与周走 AddDate 而不是乘 24 小时：有夏令时的时区里一天不总是 24 小时，
// 乘 24 小时会在切换日把时刻挪错一小时
// 月、季、年按月份推进并把日钳到目标月末，见 addMonths
//
// [EN] Adds n of the given unit; n may be negative. Days and weeks go through
// AddDate rather than multiplying by 24 hours, because a day is not always 24
// hours where daylight saving applies and multiplying shifts the clock by an
// hour on the switching day. Months, quarters and years advance by month and
// clamp the day to the end of the target month
func (m Moment) Add(n int, u Unit) Moment {
	switch u {
	case UnitNanosecond:
		return Moment{t: m.t.Add(time.Duration(n))}
	case UnitMicrosecond:
		return Moment{t: m.t.Add(time.Duration(n) * time.Microsecond)}
	case UnitMillisecond:
		return Moment{t: m.t.Add(time.Duration(n) * time.Millisecond)}
	case UnitSecond:
		return Moment{t: m.t.Add(time.Duration(n) * time.Second)}
	case UnitMinute:
		return Moment{t: m.t.Add(time.Duration(n) * time.Minute)}
	case UnitHour:
		return Moment{t: m.t.Add(time.Duration(n) * time.Hour)}
	case UnitDay:
		return Moment{t: m.t.AddDate(0, 0, n)}
	case UnitWeek:
		return Moment{t: m.t.AddDate(0, 0, 7*n)}
	case UnitMonth:
		return Moment{t: addMonths(m.t, n)}
	case UnitQuarter:
		return Moment{t: addMonths(m.t, 3*n)}
	case UnitYear:
		return Moment{t: addMonths(m.t, 12*n)}
	default:
		return m
	}
}

// Sub 按单位减 n
func (m Moment) Sub(n int, u Unit) Moment { return m.Add(-n, u) }

// addMonths 加 n 个月，并把日钳到目标月的最后一天
//
// Go 自带的 AddDate 在这里会溢出归一：1 月 31 日加一个月得到 3 月 3 日
// 钳位（得到 2 月 28 日）更合直觉，dayjs 也是这个语义
//
// [EN] Adds n months, clamping the day to the last day of the target month.
// Go's own AddDate overflows instead, taking 31 January plus one month to 3
// March; clamping to 28 February is the more intuitive result and matches dayjs
func addMonths(t time.Time, n int) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	// 把「年 + 月」折成单一的月序号再加，避免月份越界时手工处理进位
	total := year*12 + int(month) - 1 + n
	newYear := total / 12
	newMonth := total % 12
	if newMonth < 0 {
		newMonth += 12
		newYear--
	}
	target := time.Month(newMonth + 1)

	if dim := DaysInMonth(newYear, target); day > dim {
		day = dim
	}
	return time.Date(newYear, target, day, hour, min, sec, t.Nanosecond(), t.Location())
}

// StartOf 取该单位的起始时刻
//
// 周的起点取周一（ISO 8601，也是中文习惯），dayjs 默认周日，需要别的起点用
// StartOfWeekOn
//
// [EN] The start of the given unit. The week starts on Monday per ISO 8601,
// where dayjs defaults to Sunday; use StartOfWeekOn for another start
func (m Moment) StartOf(u Unit) Moment {
	t := m.t
	loc := t.Location()
	year, month, day := t.Date()
	hour, min, sec := t.Clock()

	switch u {
	case UnitNanosecond:
		return m
	case UnitMicrosecond:
		return Moment{t: time.Date(year, month, day, hour, min, sec, t.Nanosecond()/1e3*1e3, loc)}
	case UnitMillisecond:
		return Moment{t: time.Date(year, month, day, hour, min, sec, t.Nanosecond()/1e6*1e6, loc)}
	case UnitSecond:
		return Moment{t: time.Date(year, month, day, hour, min, sec, 0, loc)}
	case UnitMinute:
		return Moment{t: time.Date(year, month, day, hour, min, 0, 0, loc)}
	case UnitHour:
		return Moment{t: time.Date(year, month, day, hour, 0, 0, 0, loc)}
	case UnitDay:
		return Moment{t: time.Date(year, month, day, 0, 0, 0, 0, loc)}
	case UnitWeek:
		return m.StartOfWeekOn(time.Monday)
	case UnitMonth:
		return Moment{t: time.Date(year, month, 1, 0, 0, 0, 0, loc)}
	case UnitQuarter:
		firstMonth := time.Month((m.Quarter()-1)*3 + 1)
		return Moment{t: time.Date(year, firstMonth, 1, 0, 0, 0, 0, loc)}
	case UnitYear:
		return Moment{t: time.Date(year, time.January, 1, 0, 0, 0, 0, loc)}
	default:
		return m
	}
}

// StartOfWeekOn 取本周起始时刻，周的起点由 start 指定
//
// [EN] The start of the week, with the week's first day given by start
func (m Moment) StartOfWeekOn(start time.Weekday) Moment {
	day := m.StartOf(UnitDay)
	// 距本周起点往回几天，取模保证落在 0 到 6
	back := (int(day.t.Weekday()) - int(start) + 7) % 7
	return Moment{t: day.t.AddDate(0, 0, -back)}
}

// EndOf 取该单位的结束时刻，即最后一纳秒
//
// dayjs 到毫秒为止，而 Go 的精度到纳秒，这里不自降精度
// 用「下一单位的起点减一纳秒」算，天然处理月末、闰年与夏令时
//
// [EN] The end of the unit, its last nanosecond. dayjs stops at milliseconds
// while Go has nanosecond precision, which is not given up here. Computed as the
// next unit's start minus one nanosecond, which handles month ends, leap years
// and daylight saving for free
func (m Moment) EndOf(u Unit) Moment {
	if u == UnitNanosecond {
		return m
	}
	next := m.StartOf(u).Add(1, u)
	// 季度的 Add 已是三个月，无需另算
	return Moment{t: next.t.Add(-time.Nanosecond)}
}

// EndOfWeekOn 取本周结束时刻，周的起点由 start 指定
func (m Moment) EndOfWeekOn(start time.Weekday) Moment {
	s := m.StartOfWeekOn(start)
	return Moment{t: s.t.AddDate(0, 0, 7).Add(-time.Nanosecond)}
}

// Diff 返回 m 减 other 的差值，以该单位计，向零取整
//
// 月、季、年按日历计算而不是按固定天数相除：按天数除会把二月与七月当成一样长
// 做法是先比月份差，再看日与时分秒是否已经走满一个月
//
// [EN] The difference between m and other in the given unit, truncated toward
// zero. Months, quarters and years are counted by calendar rather than by
// dividing a fixed number of days, which would treat February and July as
// equally long: the month difference is taken first, then adjusted by whether
// the day and clock time have completed a whole month
func (m Moment) Diff(other Moment, u Unit) int {
	switch u {
	case UnitMonth, UnitQuarter, UnitYear:
		months := monthsBetween(other.t, m.t)
		switch u {
		case UnitQuarter:
			return months / 3
		case UnitYear:
			return months / 12
		default:
			return months
		}
	}

	d := m.t.Sub(other.t)
	switch u {
	case UnitNanosecond:
		return int(d)
	case UnitMicrosecond:
		return int(d / time.Microsecond)
	case UnitMillisecond:
		return int(d / time.Millisecond)
	case UnitSecond:
		return int(d / time.Second)
	case UnitMinute:
		return int(d / time.Minute)
	case UnitHour:
		return int(d / time.Hour)
	case UnitDay, UnitWeek:
		// 按日历天数算，避开夏令时让一天不足 24 小时的情况
		days := calendarDaysBetween(other.t, m.t)
		if u == UnitWeek {
			return days / 7
		}
		return days
	default:
		return 0
	}
}

// monthsBetween 返回 from 到 to 之间走满的月数，可为负
//
// 反向统一折成正向再取负，以保证 Diff(a,b) == -Diff(b,a)
//
// 不这么做的话月末钳位会让两个方向不对称：1 月 31 日加一个月钳到 2 月 28 日，
// 算正向是满一个月；而 2 月 28 日减一个月落在 1 月 28 日，已经越过 1 月 31 日，
// 算反向就不满一个月于是同一对时刻正向得 1、反向得 0，而不是 -1
// 「差值」理应反对称，所以只定义一个方向，另一个方向取负
//
// [EN] The backward direction is folded into the forward one and negated so that
// Diff(a,b) == -Diff(b,a). Without this, clamping makes the two directions
// disagree: 31 January plus a month clamps to 28 February, a full month forward,
// while 28 February minus a month lands on 28 January, already past 31 January
// and so not a full month — the same pair would give 1 forward and 0 backward
// rather than -1. A difference ought to be antisymmetric, so only one direction
// is defined and the other negates it
func monthsBetween(from, to time.Time) int {
	if to.Before(from) {
		return -monthsBetween(to, from)
	}
	months := (to.Year()-from.Year())*12 + int(to.Month()) - int(from.Month())
	// 月份差是按月份号算的，还要看日与时刻是否真的走满
	if months > 0 && addMonths(from, months).After(to) {
		months--
	}
	return months
}

// calendarDaysBetween 返回 from 到 to 之间走满的日历天数，可为负
//
// 两端都先归到当日零点再相减，再按 24 小时折算——归零之后差值必然是整天的
// 倍数加上夏令时的偏移，四舍五入即可得到整天数
func calendarDaysBetween(from, to time.Time) int {
	a := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	b := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())
	hours := b.Sub(a).Hours()
	if hours >= 0 {
		return int((hours + 12) / 24)
	}
	return -int((-hours + 12) / 24)
}
