/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-31 22:01:36
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-31 22:19:05
 * @FilePath: \go-toolbox\pkg\moment\cycle_test.go
 * @Description: 订阅周期测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package moment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func utc(y int, mo time.Month, d, h int) Moment {
	return From(time.Date(y, mo, d, h, 0, 0, 0, time.UTC))
}

func dates(ms []Moment) []string {
	out := make([]string, len(ms))
	for i, m := range ms {
		out[i] = m.Format("2006-01-02")
	}
	return out
}

// 这是本类型存在的理由：1 月 31 日订阅，扣款日不能漂移到 28 号
func TestCycleNoDriftFromMonthEnd(t *testing.T) {
	c := MustNewCycle(utc(2026, time.January, 31, 10), 1, UnitMonth)

	assert.Equal(t,
		[]string{"2026-01-31", "2026-02-28", "2026-03-31", "2026-04-30", "2026-05-31", "2026-06-30"},
		dates(c.Schedule(0, 6)),
		"每期必须从锚点算，长月要回到 31 号")

	// 对照：链式累加会漂移，这正是不能那样写的原因
	drift := utc(2026, time.January, 31, 10)
	for i := 0; i < 2; i++ {
		drift = drift.Add(1, UnitMonth)
	}
	assert.Equal(t, "2026-03-28", drift.Format("2006-01-02"), "链式累加漂移到 28 号")
	assert.Equal(t, "2026-03-31", c.NthStart(2).Format("2006-01-02"), "锚定不漂移")
}

func TestCycleQuarterAndYear(t *testing.T) {
	q := MustNewCycle(utc(2026, time.January, 31, 0), 1, UnitQuarter)
	assert.Equal(t,
		[]string{"2026-01-31", "2026-04-30", "2026-07-31", "2026-10-31", "2027-01-31"},
		dates(q.Schedule(0, 5)), "季度跨 3 个月，4 月只有 30 天")

	y := MustNewCycle(utc(2024, time.February, 29, 0), 1, UnitYear)
	assert.Equal(t,
		[]string{"2024-02-29", "2025-02-28", "2026-02-28", "2027-02-28", "2028-02-29"},
		dates(y.Schedule(0, 5)), "闰日订阅在平年钳到 28 号，闰年回到 29 号")
}

func TestCycleMultiStep(t *testing.T) {
	c := MustNewCycle(utc(2026, time.January, 31, 0), 2, UnitMonth)
	assert.Equal(t,
		[]string{"2026-01-31", "2026-03-31", "2026-05-31", "2026-07-31"},
		dates(c.Schedule(0, 4)), "每两个月一期")
}

func TestCycleIndexAtBoundaries(t *testing.T) {
	c := MustNewCycle(utc(2026, time.January, 31, 10), 1, UnitMonth)

	// 区间左闭右开
	assert.Equal(t, 0, c.IndexAt(c.NthStart(0)), "恰好等于期起点归入该期")
	assert.Equal(t, 0, c.IndexAt(c.NthEnd(0)), "期末最后一纳秒仍属该期")
	assert.Equal(t, 1, c.IndexAt(c.NthStart(1)), "下一期起点归入下一期")

	assert.Equal(t, 1, c.IndexAt(utc(2026, time.March, 15, 0)), "3-15 落在 2-28 起的那期")
	assert.Equal(t, 2, c.IndexAt(utc(2026, time.March, 31, 10)), "3-31 恰为第 2 期起点")

	// 锚点之前应为负数，不能因为整除截断而错归到第 0 期
	assert.Equal(t, -1, c.IndexAt(utc(2026, time.January, 15, 0)), "锚点之前是第 -1 期")
	assert.Equal(t, -1, c.IndexAt(utc(2025, time.December, 31, 10)), "恰为第 -1 期起点")
	assert.Equal(t, -2, c.IndexAt(utc(2025, time.December, 1, 0)), "再往前是第 -2 期")
}

func TestCycleNextIsStrictlyAfter(t *testing.T) {
	c := MustNewCycle(utc(2026, time.January, 31, 10), 1, UnitMonth)

	// 扣款语义：当天已扣过，下一次应在下个周期
	assert.Equal(t, "2026-02-28", c.Next(c.NthStart(0)).Format("2006-01-02"),
		"after 恰为期起点时返回再下一期")
	assert.Equal(t, "2026-03-31", c.Next(utc(2026, time.March, 1, 0)).Format("2006-01-02"))
	assert.Equal(t, "2026-02-28", c.Next(utc(2026, time.February, 10, 0)).Format("2006-01-02"))
}

func TestCyclePeriodAt(t *testing.T) {
	c := MustNewCycle(utc(2026, time.January, 31, 10), 1, UnitMonth)
	i, start, end := c.PeriodAt(utc(2026, time.March, 15, 0))

	assert.Equal(t, 1, i)
	assert.Equal(t, "2026-02-28", start.Format("2006-01-02"))
	assert.Equal(t, "2026-03-31", end.Format("2006-01-02"), "期末是下期起点前一纳秒，同一天")
	assert.Equal(t, start.t.Add(-time.Nanosecond), c.NthEnd(0).t, "相邻两期首尾相接不重叠")
}

func TestCycleDayAndWeek(t *testing.T) {
	d := MustNewCycle(utc(2026, time.January, 30, 0), 10, UnitDay)
	assert.Equal(t,
		[]string{"2026-01-30", "2026-02-09", "2026-02-19", "2026-03-01"},
		dates(d.Schedule(0, 4)), "每 10 天一期")

	w := MustNewCycle(utc(2026, time.January, 1, 0), 2, UnitWeek)
	assert.Equal(t,
		[]string{"2026-01-01", "2026-01-15", "2026-01-29"},
		dates(w.Schedule(0, 3)), "每两周一期")
	assert.Equal(t, 1, w.IndexAt(utc(2026, time.January, 20, 0)))
}

// 跨夏令时：按日推进必须按墙上时钟，不能按 24 小时
func TestCycleAcrossDST(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	// 2026-03-08 是美西夏令时开始日（当天只有 23 小时）
	start := From(time.Date(2026, time.March, 6, 12, 0, 0, 0, loc))
	c := MustNewCycle(start, 1, UnitDay)

	for i, want := range []string{"2026-03-06", "2026-03-07", "2026-03-08", "2026-03-09"} {
		got := c.NthStart(i)
		assert.Equal(t, want, got.Format("2006-01-02"), "第 %d 期日期", i)
		assert.Equal(t, 12, got.Hour(), "跨夏令时后仍是当地 12 点，第 %d 期", i)
	}
	assert.Equal(t, 2, c.IndexAt(From(time.Date(2026, time.March, 8, 18, 0, 0, 0, loc))),
		"夏令时切换当天的期号仍然正确")
}

func TestCycleTimezonePreserved(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)
	c := MustNewCycle(From(time.Date(2026, time.January, 31, 9, 30, 0, 0, loc)), 1, UnitMonth)

	n := c.NthStart(1)
	assert.Equal(t, loc.String(), n.Time().Location().String(), "时区随锚点保留")
	assert.Equal(t, 9, n.Hour())
	assert.Equal(t, 30, n.Minute())
}

func TestNewCycleValidation(t *testing.T) {
	start := utc(2026, time.January, 1, 0)

	for _, step := range []int{0, -1} {
		_, err := NewCycle(start, step, UnitMonth)
		assert.ErrorIs(t, err, ErrInvalidCycleStep, "step=%d 应被拒绝", step)
	}

	// 比日更小的单位对订阅没有意义，必须显式报错而不是静默接受
	for _, u := range []Unit{UnitNanosecond, UnitMicrosecond, UnitMillisecond, UnitSecond, UnitMinute, UnitHour} {
		_, err := NewCycle(start, 1, u)
		assert.ErrorIs(t, err, ErrInvalidCycleUnit, "单位 %v 应被拒绝", u)
	}

	for _, u := range []Unit{UnitDay, UnitWeek, UnitMonth, UnitQuarter, UnitYear} {
		_, err := NewCycle(start, 1, u)
		assert.NoError(t, err, "单位 %v 应被接受", u)
	}

	assert.Panics(t, func() { MustNewCycle(start, 0, UnitMonth) })
}

func TestCycleAccessorsAndSchedule(t *testing.T) {
	start := utc(2026, time.January, 31, 10)
	c := MustNewCycle(start, 3, UnitMonth)

	assert.True(t, c.Anchor().Equal(start))
	step, unit := c.Every()
	assert.Equal(t, 3, step)
	assert.Equal(t, UnitMonth, unit)

	assert.Nil(t, c.Schedule(0, 0), "count 非正返回空")
	assert.Nil(t, c.Schedule(0, -1))
	assert.Len(t, c.Schedule(5, 2), 2, "可从任意期号开始")
	assert.Equal(t, c.NthStart(5).t, c.Schedule(5, 2)[0].t)
}

func TestFloorDiv(t *testing.T) {
	// Go 的 / 向零截断，期号需要向下取整
	assert.Equal(t, 0, floorDiv(2, 3))
	assert.Equal(t, -1, floorDiv(-1, 3), "向下取整而非截断为 0")
	assert.Equal(t, -1, floorDiv(-3, 3))
	assert.Equal(t, -2, floorDiv(-4, 3))
	assert.Equal(t, 1, floorDiv(3, 3))
}

// IndexAt 与 NthStart 必须互为逆：任取一期，其内任一点都应回到该期
func TestCycleIndexRoundTrip(t *testing.T) {
	for _, u := range []Unit{UnitDay, UnitWeek, UnitMonth, UnitQuarter, UnitYear} {
		for _, step := range []int{1, 2, 7} {
			c := MustNewCycle(utc(2026, time.January, 31, 10), step, u)
			for i := -8; i <= 8; i++ {
				s, e := c.NthStart(i), c.NthEnd(i)
				assert.Equal(t, i, c.IndexAt(s), "单位=%v step=%d 期=%d 起点", u, step, i)
				assert.Equal(t, i, c.IndexAt(e), "单位=%v step=%d 期=%d 终点", u, step, i)
				assert.True(t, e.t.After(s.t), "单位=%v step=%d 期=%d 终点应晚于起点", u, step, i)
			}
		}
	}
}

// 零值 Cycle 是编程错误，但 panic 要说清缘由，而不是 integer divide by zero
func TestZeroValueCyclePanicsClearly(t *testing.T) {
	var c Cycle

	for name, call := range map[string]func(){
		"IndexAt":  func() { c.IndexAt(utc(2026, time.March, 1, 0)) },
		"NthStart": func() { c.NthStart(1) },
		"NthEnd":   func() { c.NthEnd(0) },
		"PeriodAt": func() { c.PeriodAt(utc(2026, time.March, 1, 0)) },
		"Next":     func() { c.Next(utc(2026, time.March, 1, 0)) },
		"Schedule": func() { c.Schedule(0, 2) },
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				r := recover()
				require.NotNil(t, r, "零值 Cycle 调用 %s 应当 panic", name)
				assert.Contains(t, r.(string), "Cycle 未初始化",
					"panic 信息要指出是未初始化，而不是 integer divide by zero")
			}()
			call()
		})
	}
}

// 在有夏令时的时区做往返，覆盖估算值偏小时的向上校正
//
// estimateIndex 对日/周用 24 小时估算，而夏令时当天只有 23 小时（或 25 小时），
// 估算会偏小一期，必须靠 IndexAt 的第二个循环补回来
func TestCycleIndexRoundTripAcrossDST(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	// 2026-03-08 夏令时开始（少一小时），2026-11-01 夏令时结束（多一小时）
	for _, anchor := range []time.Time{
		time.Date(2026, time.March, 6, 12, 0, 0, 0, loc),
		time.Date(2026, time.October, 30, 12, 0, 0, 0, loc),
	} {
		for _, u := range []Unit{UnitDay, UnitWeek} {
			for _, step := range []int{1, 3} {
				c := MustNewCycle(From(anchor), step, u)
				for i := 0; i <= 6; i++ {
					s, e := c.NthStart(i), c.NthEnd(i)
					assert.Equal(t, i, c.IndexAt(s),
						"锚点=%s 单位=%v step=%d 期=%d 起点", anchor.Format("01-02"), u, step, i)
					assert.Equal(t, i, c.IndexAt(e),
						"锚点=%s 单位=%v step=%d 期=%d 终点", anchor.Format("01-02"), u, step, i)
					// 墙上时钟时刻必须保持不变，而不是被 24 小时算法推偏一小时
					assert.Equal(t, 12, s.Hour(),
						"锚点=%s 单位=%v step=%d 期=%d 应仍是当地 12 点", anchor.Format("01-02"), u, step, i)
				}
			}
		}
	}
}
