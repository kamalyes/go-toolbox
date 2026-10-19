/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-18 22:38:05
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-18 23:02:11
 * @FilePath: \go-toolbox\pkg\moment\moment_test.go
 * @Description: 链式时间值类型的测试，重点覆盖与 dayjs 有意不同的几处语义
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package moment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 带夏令时的时区，用来验证按天加减不会被 24 小时假设坑到
func dstZone(t *testing.T) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skipf("当前环境没有时区数据库，跳过夏令时相关断言: %v", err)
	}
	return loc
}

// TestMomentAddMonthClamps 月末钳位：1 月 31 日加一个月应得 2 月末，而不是溢出到 3 月
//
// Go 自带的 AddDate 在这里会归一到 3 月 3 日（或闰年 3 月 2 日）
// dayjs 钳到 2 月末，本类型取 dayjs 的语义
func TestMomentAddMonthClamps(t *testing.T) {
	cases := []struct {
		name          string
		start         time.Time
		n             int
		wantY         int
		wantM         time.Month
		wantD         int
		goAddDateGets string
	}{
		{"1月31日+1月 平年", time.Date(2026, 1, 31, 10, 30, 0, 0, time.UTC), 1, 2026, time.February, 28, "2026-03-03"},
		{"1月31日+1月 闰年", time.Date(2028, 1, 31, 10, 30, 0, 0, time.UTC), 1, 2028, time.February, 29, "2028-03-02"},
		{"3月31日+1月", time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC), 1, 2026, time.April, 30, "2026-05-01"},
		{"5月31日-1月", time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC), -1, 2026, time.April, 30, "2026-05-01"},
		{"12月31日+1月 跨年", time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC), 1, 2027, time.January, 31, ""},
		{"1月31日-1月 跨年", time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC), -1, 2025, time.December, 31, ""},
		{"2月29日+12月 到平年", time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC), 12, 2029, time.February, 28, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := From(c.start).Add(c.n, UnitMonth)
			assert.Equal(t, c.wantY, got.Year(), "年")
			assert.Equal(t, c.wantM, got.Month(), "月")
			assert.Equal(t, c.wantD, got.Day(), "日")
			// 时分秒必须保持不变
			assert.Equal(t, c.start.Hour(), got.Hour(), "小时应保持不变")
			assert.Equal(t, c.start.Minute(), got.Minute(), "分钟应保持不变")

			if c.goAddDateGets != "" {
				// 说明一下：Go 自带的 AddDate 在同样输入下会溢出到下个月
				assert.Equal(t, c.goAddDateGets, c.start.AddDate(0, c.n, 0).Format(time.DateOnly),
					"这条断言只是记录 AddDate 的溢出行为，若 Go 改了语义这里会提醒")
			}
		})
	}
}

// TestMomentAddYearQuarterClamps 年与季同样要钳位
func TestMomentAddYearQuarterClamps(t *testing.T) {
	leapDay := From(time.Date(2028, 2, 29, 12, 0, 0, 0, time.UTC))
	assert.Equal(t, "2029-02-28", leapDay.Add(1, UnitYear).Format(time.DateOnly), "闰日加一年应钳到 2 月 28")
	assert.Equal(t, "2032-02-29", leapDay.Add(4, UnitYear).Format(time.DateOnly), "再逢闰年应回到 2 月 29")

	may31 := From(time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, "2026-08-31", may31.Add(1, UnitQuarter).Format(time.DateOnly))
	assert.Equal(t, "2026-11-30", may31.Add(2, UnitQuarter).Format(time.DateOnly), "11 月只有 30 天，应钳位")
}

// TestMomentAddDayAcrossDST 按天加减必须走日历，不能按 24 小时算
//
// 美东 2026-03-08 是夏令时开始日，那天只有 23 小时
// 用「加 24 小时」会落在 3 月 9 日的 01:30 而不是 00:30
func TestMomentAddDayAcrossDST(t *testing.T) {
	loc := dstZone(t)
	// 夏令时开始前一天的 00:30
	start := time.Date(2026, 3, 7, 0, 30, 0, 0, loc)

	plus2 := From(start).Add(2, UnitDay)
	assert.Equal(t, "2026-03-09", plus2.Format(time.DateOnly))
	assert.Equal(t, 0, plus2.Hour(), "跨夏令时后钟面小时应保持 0")
	assert.Equal(t, 30, plus2.Minute(), "钟面分钟应保持 30")

	// 对照：按 24 小时乘会偏一小时
	naive := start.Add(2 * 24 * time.Hour)
	assert.Equal(t, 1, naive.Hour(),
		"这条断言只是记录乘 24 小时的偏差，若时区数据变动这里会提醒")

	// 按周同理
	plusWeek := From(start).Add(1, UnitWeek)
	assert.Equal(t, "2026-03-14", plusWeek.Format(time.DateOnly))
	assert.Equal(t, 0, plusWeek.Hour(), "加一周后钟面小时应保持 0")
}

// TestMomentStartOfEndOf 各单位的起止边界
func TestMomentStartOfEndOf(t *testing.T) {
	// 2026-10-18 是周日
	base := From(time.Date(2026, 10, 18, 15, 37, 29, 123456789, time.UTC))
	assert.Equal(t, time.Sunday, base.Weekday(), "基准日应为周日，否则本测试的周边界预期不成立")

	cases := []struct {
		unit             Unit
		wantStart, wantE string
	}{
		{UnitSecond, "2026-10-18 15:37:29.000000000", "2026-10-18 15:37:29.999999999"},
		{UnitMinute, "2026-10-18 15:37:00.000000000", "2026-10-18 15:37:59.999999999"},
		{UnitHour, "2026-10-18 15:00:00.000000000", "2026-10-18 15:59:59.999999999"},
		{UnitDay, "2026-10-18 00:00:00.000000000", "2026-10-18 23:59:59.999999999"},
		// 周一为起点，所以 10-18（周日）所在周是 10-12 到 10-18
		{UnitWeek, "2026-10-12 00:00:00.000000000", "2026-10-18 23:59:59.999999999"},
		{UnitMonth, "2026-10-01 00:00:00.000000000", "2026-10-31 23:59:59.999999999"},
		{UnitQuarter, "2026-10-01 00:00:00.000000000", "2026-12-31 23:59:59.999999999"},
		{UnitYear, "2026-01-01 00:00:00.000000000", "2026-12-31 23:59:59.999999999"},
	}
	const layout = "2006-01-02 15:04:05.000000000"
	for _, c := range cases {
		t.Run(c.unit.String(), func(t *testing.T) {
			assert.Equal(t, c.wantStart, base.StartOf(c.unit).Format(layout), "StartOf")
			assert.Equal(t, c.wantE, base.EndOf(c.unit).Format(layout), "EndOf")
		})
	}
}

// TestMomentStartOfWeekOn 周起点可指定，默认周一而非 dayjs 的周日
func TestMomentStartOfWeekOn(t *testing.T) {
	// 2026-10-18 周日
	sunday := From(time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC))
	assert.Equal(t, "2026-10-12", sunday.StartOf(UnitWeek).Format(time.DateOnly), "默认周一为起点")
	assert.Equal(t, "2026-10-18", sunday.StartOfWeekOn(time.Sunday).Format(time.DateOnly), "指定周日为起点")
	assert.Equal(t, "2026-10-24", sunday.EndOfWeekOn(time.Sunday).Format(time.DateOnly), "周日起点的周末")

	// 周一当天，两种起点下结果不同
	monday := From(time.Date(2026, 10, 19, 12, 0, 0, 0, time.UTC))
	assert.Equal(t, time.Monday, monday.Weekday())
	assert.Equal(t, "2026-10-19", monday.StartOfWeekOn(time.Monday).Format(time.DateOnly))
	assert.Equal(t, "2026-10-18", monday.StartOfWeekOn(time.Sunday).Format(time.DateOnly))
}

// TestMomentQuarter 季度编号
func TestMomentQuarter(t *testing.T) {
	for month, wantQ := range map[time.Month]int{
		time.January: 1, time.March: 1, time.April: 2, time.June: 2,
		time.July: 3, time.September: 3, time.October: 4, time.December: 4,
	} {
		got := From(time.Date(2026, month, 15, 0, 0, 0, 0, time.UTC)).Quarter()
		assert.Equal(t, wantQ, got, "%s 应属第 %d 季度", month, wantQ)
	}
}

// TestMomentDiff 差值按单位计算，月与年按日历而非固定天数
func TestMomentDiff(t *testing.T) {
	a := From(time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC))

	// 按日历：1月31日到2月28日还差一天才满一个月
	assert.Equal(t, 0, From(time.Date(2026, 2, 27, 0, 0, 0, 0, time.UTC)).Diff(a, UnitMonth),
		"未走满一个月应为 0")
	assert.Equal(t, 1, From(time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)).Diff(a, UnitMonth),
		"钳位后 2 月 28 日即满一个月")
	assert.Equal(t, 12, From(time.Date(2027, 1, 31, 0, 0, 0, 0, time.UTC)).Diff(a, UnitMonth))
	assert.Equal(t, 1, From(time.Date(2027, 1, 31, 0, 0, 0, 0, time.UTC)).Diff(a, UnitYear))
	assert.Equal(t, 4, From(time.Date(2027, 1, 31, 0, 0, 0, 0, time.UTC)).Diff(a, UnitQuarter))

	// 负方向
	assert.Equal(t, -1, a.Diff(From(time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)), UnitMonth))

	// 时间类单位
	t0 := From(time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC))
	t1 := From(time.Date(2026, 10, 20, 6, 30, 0, 0, time.UTC))
	assert.Equal(t, 2, t1.Diff(t0, UnitDay))
	assert.Equal(t, 54, t1.Diff(t0, UnitHour))
	assert.Equal(t, 54*60+30, t1.Diff(t0, UnitMinute))
	assert.Equal(t, 0, t1.Diff(t0, UnitWeek), "不足一周应为 0")
	assert.Equal(t, 1, From(time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)).Diff(t0, UnitWeek))
}

// TestMomentDiffDayAcrossDST 跨夏令时的天数差按日历算
//
// 美东 3 月 8 日只有 23 小时，按小时相除会把两天算成 1 天
func TestMomentDiffDayAcrossDST(t *testing.T) {
	loc := dstZone(t)
	from := From(time.Date(2026, 3, 7, 12, 0, 0, 0, loc))
	to := From(time.Date(2026, 3, 9, 12, 0, 0, 0, loc))
	assert.Equal(t, 2, to.Diff(from, UnitDay), "跨夏令时仍应是 2 天")

	// 对照：按小时数相除会少一天
	hours := to.Time().Sub(from.Time()).Hours()
	assert.Equal(t, 47.0, hours, "这两天实际只隔 47 小时，记录于此以说明为何不能按 24 小时相除")
}

// TestMomentPreservesLocation 所有链式操作都要保留时区
func TestMomentPreservesLocation(t *testing.T) {
	loc := time.FixedZone("CST", 8*3600)
	m := From(time.Date(2026, 10, 18, 15, 30, 0, 0, loc))

	for name, got := range map[string]Moment{
		"Add/day":       m.Add(1, UnitDay),
		"Add/month":     m.Add(1, UnitMonth),
		"Sub/year":      m.Sub(1, UnitYear),
		"StartOf/month": m.StartOf(UnitMonth),
		"EndOf/year":    m.EndOf(UnitYear),
		"StartOf/week":  m.StartOf(UnitWeek),
	} {
		assert.Equal(t, loc.String(), got.Time().Location().String(), "%s 应保留时区", name)
	}

	// In 换时区后时刻不变
	utc := m.In(time.UTC)
	assert.True(t, utc.Equal(m), "换时区不应改变时刻")
	assert.Equal(t, 7, utc.Hour(), "+0800 的 15:30 在 UTC 应是 07:30")
}

// TestMomentChainIsAllocationFree 链式调用不应产生堆分配
//
// Moment 是值类型，整条链都在栈上这是它相对 dayjs 的实质优势——dayjs 每次
// 操作都要 new 一个对象一旦有人把它改成指针接收者或在链中引入接口，这里会红
//
// [EN] A chain must not allocate. Moment is a value type so the whole chain
// stays on the stack, which is the substantive advantage over dayjs, where every
// operation allocates an object. This goes red if someone switches to pointer
// receivers or introduces an interface into the chain
func TestMomentChainIsAllocationFree(t *testing.T) {
	base := time.Date(2026, 10, 18, 15, 30, 0, 0, time.UTC)
	var sink int

	allocs := testing.AllocsPerRun(200, func() {
		sink = From(base).
			Add(7, UnitDay).
			Add(1, UnitMonth).
			StartOf(UnitMonth).
			EndOf(UnitDay).
			Sub(2, UnitHour).
			Day()
	})
	assert.Zero(t, allocs, "链式调用应零分配，实测每次 %.1f 次", allocs)
	assert.NotZero(t, sink)

	assert.Zero(t, testing.AllocsPerRun(200, func() {
		sink = From(base).Diff(From(base).Add(3, UnitMonth), UnitDay)
	}), "Diff 应零分配")
}

// TestMomentZeroValue 零值可用，不应 panic
func TestMomentZeroValue(t *testing.T) {
	var m Moment
	assert.True(t, m.IsZero())
	assert.NotPanics(t, func() {
		_ = m.Add(1, UnitDay)
		_ = m.StartOf(UnitYear)
		_ = m.EndOf(UnitMonth)
		_ = m.Quarter()
		_ = m.String()
	})
}

// TestUnitString 单位名称齐全，便于错误信息与日志
func TestUnitString(t *testing.T) {
	for u := UnitNanosecond; u <= UnitYear; u++ {
		assert.NotEqual(t, "unknown", u.String(), "单位 %d 缺少名称", u)
	}
	assert.Equal(t, "unknown", Unit(200).String())
}

// BenchmarkMomentChain 链式调用的开销
func BenchmarkMomentChain(b *testing.B) {
	base := time.Date(2026, 10, 18, 15, 30, 0, 0, time.UTC)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = From(base).Add(7, UnitDay).StartOf(UnitMonth).EndOf(UnitDay).Day()
	}
}

func BenchmarkMomentAddMonth(b *testing.B) {
	base := time.Date(2026, 1, 31, 15, 30, 0, 0, time.UTC)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = From(base).Add(1, UnitMonth)
	}
}

func BenchmarkMomentDiffMonth(b *testing.B) {
	a := From(time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC))
	c := From(time.Date(2027, 3, 15, 0, 0, 0, 0, time.UTC))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Diff(a, UnitMonth)
	}
}

// TestMomentDiffAntisymmetric 差值必须反对称：Diff(a,b) == -Diff(b,a)
//
// 月末钳位很容易破坏这条性质，用成对时刻全量扫一遍守住
//
// [EN] Diff must be antisymmetric. Month-end clamping breaks this easily, so a
// sweep over pairs of instants pins it
func TestMomentDiffAntisymmetric(t *testing.T) {
	var times []Moment
	for _, d := range []struct {
		y    int
		m    time.Month
		day  int
		hour int
	}{
		{2026, time.January, 1, 0}, {2026, time.January, 31, 23},
		{2026, time.February, 1, 0}, {2026, time.February, 28, 12},
		{2028, time.February, 29, 6}, {2026, time.March, 1, 0},
		{2026, time.March, 31, 18}, {2026, time.April, 30, 9},
		{2026, time.December, 31, 23}, {2027, time.January, 1, 0},
		{2025, time.June, 15, 7}, {2029, time.November, 30, 15},
	} {
		times = append(times, From(time.Date(d.y, d.m, d.day, d.hour, 17, 3, 0, time.UTC)))
	}

	units := []Unit{UnitSecond, UnitMinute, UnitHour, UnitDay, UnitWeek, UnitMonth, UnitQuarter, UnitYear}
	for _, u := range units {
		for i, a := range times {
			for j, b := range times {
				ab, ba := a.Diff(b, u), b.Diff(a, u)
				assert.Equal(t, ab, -ba,
					"单位 %s 下第 %d 与第 %d 个时刻不反对称: %d 对 %d (%s / %s)",
					u, i, j, ab, ba, a.Format(time.DateOnly), b.Format(time.DateOnly))
			}
		}
	}
}
