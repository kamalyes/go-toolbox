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
	"fmt"
	"math"
	"math/rand"
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

// daysFromCivil 是非平凡算术，用 time.Time 逐日推进做差分对照
//
// 常量另经 Python 的 datetime（外推格里历）独立核算，避免手算错一位还自圆其说
func TestDaysFromCivil(t *testing.T) {
	start := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2100, 12, 31, 0, 0, 0, 0, time.UTC)
	base, n := daysFromCivil(1900, 1, 1), 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		y, mo, dd := d.Date()
		require.Equal(t, base+n, daysFromCivil(y, int(mo), dd),
			"daysFromCivil(%d,%d,%d) 与 time.Time 的逐日推进不一致", y, mo, dd)
		n++
	}

	for _, tc := range []struct{ y, m, d, want int }{
		{1970, 1, 1, 0}, {1970, 1, 2, 1}, {1969, 12, 31, -1},
		{2000, 2, 29, 11016}, {2024, 2, 29, 19782},
		{1900, 3, 1, -25508}, // 1900 不是闰年
		{1, 1, 1, -719162},   // 零值 time.Time
	} {
		assert.Equal(t, tc.want, daysFromCivil(tc.y, tc.m, tc.d),
			"daysFromCivil(%d,%d,%d)", tc.y, tc.m, tc.d)
	}
}

// 溢出必须显式 panic，而不是静默返回一个看起来正常的日期
func TestCycleOverflowPanics(t *testing.T) {
	anchor := utc(2026, time.January, 31, 10)

	cases := []struct {
		name string
		call func()
		want string
	}{
		// 未加检查时返回 2025-12-28，比锚点还早
		{"月周期 MaxInt", func() { MustNewCycle(anchor, 1, UnitMonth).NthStart(math.MaxInt) },
			"超出可表示"},
		// 未加检查时 12*MinInt 回绕成 0，返回锚点本身
		{"年周期 MinInt", func() { MustNewCycle(anchor, 1, UnitYear).NthStart(math.MinInt) },
			"超出可表示"},
		// i*step 本身就溢出
		{"step=13 乘法溢出", func() { MustNewCycle(anchor, 13, UnitMonth).NthStart(math.MaxInt / 2) },
			"溢出"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				require.NotNil(t, r, "应当 panic")
				assert.Contains(t, r.(string), tc.want)
			}()
			tc.call()
		})
	}
}

// 回归闸门：相距很远时估算仍须贴近真值
//
// 原先 estimateIndex 用 t.Sub(anchor)，而 Duration 在约 292 年处饱和，
// 估算被钉住后校正循环要爬几十万步——实测日周期相距 2000 年时单次 IndexAt 要 21.7 毫秒，
// 零值 time.Time 作锚点要 14.5 毫秒
// 这里不断言耗时（会受机器影响而不稳），直接断言估算与真值的差距，
// 那才是当初坏掉的性质
func TestCycleEstimateStaysTightWhenFarApart(t *testing.T) {
	for _, u := range []Unit{UnitDay, UnitWeek, UnitMonth, UnitQuarter, UnitYear} {
		for _, step := range []int{1, 3} {
			c := MustNewCycle(utc(2026, time.January, 31, 10), step, u)
			for _, probe := range []Moment{
				From(time.Time{}), // 公元 1 年
				From(c.Anchor().t.AddDate(-2000, 0, 0)),
				From(c.Anchor().t.AddDate(-500, 0, 0)),
				From(c.Anchor().t.AddDate(500, 0, 0)),
				From(c.Anchor().t.AddDate(2000, 0, 0)),
			} {
				est, real := c.estimateIndex(probe), c.IndexAt(probe)
				diff := est - real
				if diff < 0 {
					diff = -diff
				}
				assert.LessOrEqual(t, diff, 1,
					"单位=%v step=%d t=%s 估算=%d 真值=%d 偏差过大",
					u, step, probe.Format("2006-01-02"), est, real)
			}
		}
	}
}

// 春令时空洞：锚点时刻在那一天不存在，Go 会把它归一化
//
// 这不是缺陷而是墙上时钟运算的固有结果，但行为反直觉，必须钉住：
// 美西 2026-03-08 的 02:00 到 03:00 不存在，02:30 的锚点在那天被归一化为 01:30 而非 03:30，
// 也就是说当地时刻会往回跳一小时，次日恢复
// 绝对时间仍严格单调递增，期号与区间都不受影响
func TestCycleAnchorInDSTGap(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	c := MustNewCycle(From(time.Date(2026, time.March, 7, 2, 30, 0, 0, loc)), 1, UnitDay)

	assert.Equal(t, []string{
		"03-07 02:30 PST",
		"03-08 01:30 PST", // 02:30 不存在，被归一化到 01:30
		"03-09 02:30 PDT",
	}, []string{
		c.NthStart(0).Format("01-02 15:04 MST"),
		c.NthStart(1).Format("01-02 15:04 MST"),
		c.NthStart(2).Format("01-02 15:04 MST"),
	}, "空洞当天当地时刻往回跳一小时，次日恢复")

	// 绝对时间仍严格单调，期号仍正确
	for i := 0; i <= 2; i++ {
		assert.True(t, c.NthStart(i+1).t.After(c.NthStart(i).t), "第%d期后仍严格递增", i)
		assert.Equal(t, i, c.IndexAt(c.NthStart(i)), "第%d期期号", i)
	}
}

// 性质测试：拿不变式撞实现，覆盖手写用例想不到的组合
func TestCycleInvariants(t *testing.T) {
	zones := []string{"UTC", "Asia/Shanghai", "America/Los_Angeles", "Australia/Lord_Howe"}
	units := []Unit{UnitDay, UnitWeek, UnitMonth, UnitQuarter, UnitYear}
	steps := []int{1, 2, 3, 7, 13}
	anchors := []struct{ y, mo, d, h, mi int }{
		{2026, 1, 31, 10, 0}, {2026, 1, 31, 0, 0}, {2026, 1, 31, 23, 59},
		{2024, 2, 29, 12, 0}, {2026, 3, 8, 2, 30}, {2026, 3, 8, 3, 30},
		{2026, 11, 1, 1, 30}, {2026, 12, 31, 23, 0}, {2026, 2, 28, 0, 0},
		{2026, 4, 30, 12, 0}, {2026, 5, 31, 12, 0},
	}
	rng := rand.New(rand.NewSource(20261031))

	for _, zn := range zones {
		loc, err := time.LoadLocation(zn)
		require.NoError(t, err)
		for _, as := range anchors {
			anchor := From(time.Date(as.y, time.Month(as.mo), as.d, as.h, as.mi, 0, 0, loc))
			for _, u := range units {
				for _, step := range steps {
					c := MustNewCycle(anchor, step, u)
					tag := fmt.Sprintf("[%s %v step=%d 锚=%s]", zn, u, step, anchor.Format("2006-01-02 15:04"))

					// NthStart 对 i 严格单调递增——IndexAt 的终止性依赖这条
					prev := c.NthStart(-12)
					for i := -11; i <= 12; i++ {
						cur := c.NthStart(i)
						require.True(t, cur.t.After(prev.t), "%s 第%d期未晚于前一期", tag, i)
						prev = cur
					}

					for i := -6; i <= 6; i++ {
						require.Equal(t, c.NthStart(i+1).t.Add(-time.Nanosecond), c.NthEnd(i).t,
							"%s 第%d期末应为下期起点前一纳秒", tag, i)
						require.Equal(t, i, c.IndexAt(c.NthStart(i)), "%s 第%d期起点往返", tag, i)
						require.Equal(t, i, c.IndexAt(c.NthEnd(i)), "%s 第%d期末往返", tag, i)
					}

					// 任取时刻都应满足 NthStart(i) <= t < NthStart(i+1)，且 Next 严格晚于 t
					for k := 0; k < 30; k++ {
						off := time.Duration(rng.Int63n(int64(4000*24*time.Hour))) - time.Duration(2000*24*time.Hour)
						probe := From(anchor.t.Add(off))
						i := c.IndexAt(probe)
						require.False(t, c.NthStart(i).t.After(probe.t), "%s t=%s 期起点晚于 t", tag, probe.t)
						require.True(t, c.NthStart(i+1).t.After(probe.t), "%s t=%s 下期起点未晚于 t", tag, probe.t)
						require.True(t, c.Next(probe).t.After(probe.t), "%s Next 未严格晚于入参", tag)
					}
				}
			}
		}
	}
}

// 校正必须对任意起始估算都收敛到同一答案，两个方向都要走到
//
// 当前估算是真值上界，正常路径只会向下校正；这里故意喂偏大与偏小的值，
// 把向上那支也走到，免得它成为无人验证的死代码
func TestCycleCorrectIndexBothDirections(t *testing.T) {
	for _, u := range []Unit{UnitDay, UnitWeek, UnitMonth, UnitQuarter, UnitYear} {
		for _, step := range []int{1, 3} {
			c := MustNewCycle(utc(2026, time.January, 31, 10), step, u)
			for _, probe := range []Moment{
				utc(2026, time.March, 15, 7),
				utc(2026, time.January, 31, 10),
				utc(2025, time.June, 2, 23),
				utc(2031, time.December, 9, 0),
			} {
				want := c.IndexAt(probe)
				for _, skew := range []int{-9, -3, -1, 0, 1, 3, 9} {
					assert.Equal(t, want, c.correctIndex(want+skew, probe),
						"单位=%v step=%d t=%s 估算偏移=%d 未收敛到 %d",
						u, step, probe.Format("2006-01-02"), skew, want)
				}
			}
		}
	}
}

// daysFromCivil 在年 <= 0 时走另一条取整分支，用格里历 400 年周期性验证
//
// 格里历每 400 年恰好重复 146097 天，这是个与实现无关的数学恒等式，
// 比我自己手算负年份的常量可靠
func TestDaysFromCivilNegativeYears(t *testing.T) {
	const cycleDays = 146097
	for _, tc := range []struct{ y, m, d int }{
		{2026, 1, 31}, {2000, 2, 29}, {1, 1, 1}, {1, 3, 1},
		{0, 1, 1}, {0, 3, 1}, {400, 2, 29},
	} {
		got := daysFromCivil(tc.y, tc.m, tc.d)
		back := daysFromCivil(tc.y-400, tc.m, tc.d)
		assert.Equal(t, got-cycleDays, back,
			"%d-%d-%d 往前推 400 年应恰好少 %d 天", tc.y, tc.m, tc.d, cycleDays)

		fwd := daysFromCivil(tc.y+400, tc.m, tc.d)
		assert.Equal(t, got+cycleDays, fwd,
			"%d-%d-%d 往后推 400 年应恰好多 %d 天", tc.y, tc.m, tc.d, cycleDays)
	}

	// 年 0 与年 -1 的相邻日仍须连续（跨过取整分支的边界）
	assert.Equal(t, daysFromCivil(0, 1, 1), daysFromCivil(-1, 12, 31)+1,
		"0001-01-01 应紧接 0000-12-31")
	assert.Equal(t, daysFromCivil(1, 1, 1), daysFromCivil(0, 12, 31)+1,
		"年 1 的 1 月 1 日应紧接年 0 的 12 月 31 日")
}
