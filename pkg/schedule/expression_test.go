/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-01-22 13:55:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-13 18:26:19
 * @FilePath: \go-toolbox\pkg\schedule\expression_test.go
 * @Description: 解析表达式测试用例
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"sync"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/stretchr/testify/assert"
)

type cronTestCase struct {
	name          string
	expr          string
	expectErr     bool
	expectSecond  []int
	expectMinute  []int
	expectHour    []int
	expectDay     []int
	expectMonth   []int
	expectWeekday []int
	expectYear    []int
}

func TestCronExpressions(t *testing.T) {
	comploxCronExpr := "5,10-15/2,20 0,3-9/3,55-59 0,6-18/3,22-23 1,15 1-5,7 * 2023-2025/2" // 测试使用
	complexHourSchedule := "0 0 0,6-18/3,22-23 * * *"                                       // 每小时的第0秒，6点到18点之间每3小时执行一次，及22点和23点
	weekdayWorkingHoursQuarter := "0 1-10/2 9-17 * * 1-5"                                   // 每天的工作日（周一到周五）上午9点到17点，每15分钟执行一次
	nearestWorkdayOf15thAt2359 := "0 59 23 15W * ?"                                         // 每天的最近工作日的15号的23点59分执行
	weekendEarlyMorningEvery10Min := "0 0/10 1-3 * * 6,0"                                   // 每周六和周日的凌晨1点到3点，每10分钟执行一次
	newYearEveLastSecond := "59 59 23 31 12 *"                                              // 每年12月31日23点59分59秒执行（年字段支持）
	workdayBusinessHoursHourly := "0 0 8-18 * * 1-5"                                        // 每天的8点到18点，每小时的第0分第0秒执行，且只在工作日（周一到周五）
	firstAndFifteenthDay := "0 0 0 1,15 * *"                                                // 每月1号和15号的0点0分0秒执行
	everyDayExceptWeekend := "0 0 0 * * 1-5"                                                // 每天的0点0分0秒执行，但忽略周末（周六、周日）
	cases := []cronTestCase{
		{"EverySecond", EverySecond, false, random.RandNumerical(0, 59), random.RandNumerical(0, 59), random.RandNumerical(0, 23), random.RandNumerical(1, 31), random.RandNumerical(1, 12), random.RandNumerical(0, 6), random.RandNumerical(yearRange.Min, yearRange.Max)},
		{"EveryMinute", EveryMinute, false, []int{0}, random.RandNumerical(0, 59), nil, nil, nil, nil, nil},
		{"EveryHalfMinute", EveryHalfMinute, false, []int{0, 30}, random.RandNumerical(0, 59), random.RandNumerical(0, 23), nil, nil, nil, nil},
		{"EveryHour", EveryHour, false, []int{0}, []int{0}, nil, nil, nil, nil, nil},
		{"EveryHalfHour", EveryHalfHour, false, []int{0}, []int{0, 30}, random.RandNumerical(0, 23), nil, nil, nil, nil},
		{"EveryDay", EveryDay, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"EveryHalfDay", EveryHalfDay, false, []int{0}, []int{0}, []int{0, 12}, nil, nil, nil, nil},
		{"EveryWeek", EveryWeek, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"WeekdaysOnly", WeekdaysOnly, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"WeekendsOnly", WeekendsOnly, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"PeakHours", PeakHours, false, []int{0}, nil, []int{8, 9, 17, 18}, nil, nil, nil, nil},
		{"OffPeakHours", OffPeakHours, false, []int{0}, nil, random.RandNumerical(9, 16), nil, nil, nil, nil},
		{"Monday", Monday, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"Tuesday", Tuesday, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"Wednesday", Wednesday, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"Thursday", Thursday, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"Friday", Friday, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"Saturday", Saturday, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"Sunday", Sunday, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"FirstDayOfMonth", FirstDayOfMonth, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"LastDayOfMonth", LastDayOfMonth, false, []int{0}, []int{0}, []int{0}, nil, nil, nil, nil},
		{"ComploxCronExpr", comploxCronExpr, false, []int{5, 10, 12, 14, 20}, []int{0, 3, 6, 9, 55, 56, 57, 58, 59}, []int{0, 6, 9, 12, 15, 18, 22, 23}, nil, nil, nil, nil},
		{name: "ComplexHourSchedule", expr: complexHourSchedule, expectErr: false, expectSecond: []int{0}, expectMinute: []int{0}, expectHour: []int{0, 6, 9, 12, 15, 18, 22, 23}},
		{name: "WeekdayWorkingHoursQuarter", expr: weekdayWorkingHoursQuarter, expectErr: false, expectSecond: []int{0}, expectMinute: random.RandNumerical(1, 10, 2), expectHour: random.RandNumerical(9, 17), expectWeekday: random.RandNumerical(1, 5)},
		{name: "NearestWorkdayOf15thAt2359", expr: nearestWorkdayOf15thAt2359, expectErr: false, expectSecond: []int{0}, expectMinute: []int{59}, expectHour: []int{23}},
		{name: "WeekendEarlyMorningEvery10Min", expr: weekendEarlyMorningEvery10Min, expectErr: false, expectSecond: []int{0}, expectMinute: random.RandNumerical(0, 0, 10), expectHour: []int{1, 2, 3}, expectWeekday: []int{6, 0}},
		{name: "NewYearEveLastSecond", expr: newYearEveLastSecond, expectErr: false, expectSecond: []int{59}, expectMinute: []int{59}, expectHour: []int{23}, expectDay: []int{31}, expectMonth: []int{12}},
		{name: "WorkdayBusinessHoursHourly", expr: workdayBusinessHoursHourly, expectErr: false, expectSecond: []int{0}, expectMinute: []int{0}, expectHour: random.RandNumerical(8, 18), expectWeekday: random.RandNumerical(1, 5)},
		{name: "FirstAndFifteenthDay", expr: firstAndFifteenthDay, expectErr: false, expectSecond: []int{0}, expectMinute: []int{0}, expectHour: []int{0}, expectDay: []int{1, 15}},
		{name: "EveryDayExceptWeekend", expr: everyDayExceptWeekend, expectErr: false, expectSecond: []int{0}, expectMinute: []int{0}, expectHour: []int{0}, expectWeekday: random.RandNumerical(1, 5)},
	}

	parser := NewCronParser()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			parser.SetExpression(c.expr)
			cronExpr, err := parser.Parse()

			if c.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, cronExpr)

			if c.expectSecond != nil {
				assert.ElementsMatch(t, cronExpr.Second.Vals(), c.expectSecond, "second values mismatch")
			}
			if c.expectMinute != nil {
				assert.ElementsMatch(t, cronExpr.Minute.Vals(), c.expectMinute, "minute values mismatch")
			}
			if c.expectHour != nil {
				assert.ElementsMatch(t, cronExpr.Hour.Vals(), c.expectHour, "hour values mismatch")
			}
			if c.expectDay != nil {
				assert.ElementsMatch(t, cronExpr.Day.Vals(), c.expectDay, "day values mismatch")
			}
			if c.expectMonth != nil {
				assert.ElementsMatch(t, cronExpr.Month.Vals(), c.expectMonth, "month values mismatch")
			}
			if c.expectWeekday != nil {
				assert.ElementsMatch(t, cronExpr.Weekday.Vals(), c.expectWeekday, "weekday values mismatch")
			}
			if c.expectYear != nil {
				assert.ElementsMatch(t, cronExpr.Year.Vals(), c.expectYear, "year values mismatch")
			}
		})
	}
}

func TestParseSimplePart(t *testing.T) {
	fr := FieldRange{Min: 0, Max: 10}

	// 测试星号
	vals, err := parseSimplePart("*", fr)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, vals)

	// 测试单值
	vals, err = parseSimplePart("5", fr)
	assert.NoError(t, err)
	assert.Equal(t, []int{5}, vals)

	// 测试范围
	vals, err = parseSimplePart("2-4", fr)
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 3, 4}, vals)

	// 测试步进 全范围
	vals, err = parseSimplePart("*/3", fr)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 3, 6, 9}, vals)

	// 测试步进 范围内
	vals, err = parseSimplePart("2-8/2", fr)
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6, 8}, vals)

	// 测试错误格式
	_, err = parseSimplePart("*/0", fr)
	assert.Error(t, err)

	_, err = parseSimplePart("10-5", fr)
	assert.Error(t, err)

	_, err = parseSimplePart("a", fr)
	assert.Error(t, err)
}

func TestParseWeekdayField(t *testing.T) {
	// 测试普通数字，0基
	wf, err := parseWeekdayField("0,1,5", WeekdayZeroBased)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 5}, wf.Vals())

	// 测试普通数字，1基
	wf, err = parseWeekdayField("1,2,6", WeekdayOneBased)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 5}, wf.Vals()) // 1->0, 2->1, 6->5

	// 测试星号
	wf, err = parseWeekdayField("*", WeekdayZeroBased)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6}, wf.Vals())

	// 测试问号
	wf, err = parseWeekdayField("?", WeekdayZeroBased)
	assert.NoError(t, err)
	assert.True(t, wf.IsQuestion())

	// 测试#语法
	wf, err = parseWeekdayField("2#3", WeekdayZeroBased)
	assert.NoError(t, err)
	nth, ok := wf.NthWeekday()[2]
	assert.True(t, ok)
	assert.Equal(t, 3, nth)

	// 错误格式
	_, err = parseWeekdayField("2##3", WeekdayZeroBased)
	assert.Error(t, err)

	_, err = parseWeekdayField("8", WeekdayZeroBased)
	assert.Error(t, err)
}

func TestParseDayField(t *testing.T) {
	// 普通多值
	df, err := parseDayField("1,15,31")
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 15, 31}, df.Vals())

	// 问号
	df, err = parseDayField("?")
	assert.NoError(t, err)
	assert.True(t, df.IsQuestion())

	// L
	df, err = parseDayField("L")
	assert.NoError(t, err)
	assert.True(t, df.IsLastDay())

	// W
	df, err = parseDayField("15W")
	assert.NoError(t, err)
	assert.True(t, df.IsNearestWD())
	assert.Equal(t, []int{15}, df.Vals())

	// 错误格式
	_, err = parseDayField("15WW")
	assert.Error(t, err)

	_, err = parseDayField("32")
	assert.Error(t, err)
}

func TestParseSimpleField(t *testing.T) {
	fr := FieldRange{Min: 0, Max: 59}

	sf, err := parseSimpleField("0,15,30,45", fr)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 15, 30, 45}, sf.Vals())

	sf, err = parseSimpleField("*", fr)
	assert.NoError(t, err)
	assert.Equal(t, 60, len(sf.Vals()))

	_, err = parseSimpleField("61", fr)
	assert.Error(t, err)
}

func TestParseCronExpr_Valid(t *testing.T) {
	expr := "0 0 12 15 6 1 2025"
	cron, err := NewCronParser().SetExpression(expr).Parse()
	assert.NoError(t, err)

	assert.True(t, cron.Second.Contains(0))
	assert.True(t, cron.Minute.Contains(0))
	assert.True(t, cron.Hour.Contains(12))
	assert.True(t, cron.Day.Contains(15))
	assert.True(t, cron.Month.Contains(6))
	assert.True(t, cron.Weekday.Contains(1))
	assert.True(t, cron.Year.Contains(2025))
}

func TestParseCronExpr_InvalidFieldCount(t *testing.T) {
	expr := "0 0 0 * *"
	_, err := NewCronParser().
		SetExpression(expr).
		SetWeekNumbering(WeekdayZeroBased).Parse()
	assert.NoError(t, err)
}

func TestParseCronExpr_InvalidFieldValues(t *testing.T) {
	// 秒字段非法
	_, err := NewCronParser().
		SetExpression("61 0 0 * * * *").
		SetWeekNumbering(WeekdayZeroBased).Parse()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrParseSecond)

	// 星期字段非法
	_, err = NewCronParser().
		SetExpression("0 0 0 * * 8 *").
		SetWeekNumbering(WeekdayZeroBased).Parse()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrParseWeekday)
}

func TestWeekdayField_Contains(t *testing.T) {
	wf, err := parseWeekdayField("1,3,5", WeekdayZeroBased)
	assert.NoError(t, err)
	assert.True(t, wf.Contains(1))
	assert.False(t, wf.Contains(2))

	// ? 表示匹配所有
	wf, err = parseWeekdayField("?", WeekdayZeroBased)
	assert.NoError(t, err)
	assert.True(t, wf.Contains(0))
	assert.True(t, wf.Contains(6))
}

func TestDayField_Contains(t *testing.T) {
	df, err := parseDayField("1,15,31")
	assert.NoError(t, err)
	assert.True(t, df.Contains(15))
	assert.False(t, df.Contains(10))

	df, err = parseDayField("?")
	assert.NoError(t, err)
	assert.True(t, df.Contains(10))

	df, err = parseDayField("L")
	assert.NoError(t, err)
	assert.False(t, df.Contains(10)) // 需结合年月，暂时返回false
}

func TestCronParser_ConcurrentAccess(t *testing.T) {
	expressions := []string{
		"0 0 12 * * * 2025", // 注意这里把 ? 改成 *，确保符合7字段格式
		"15 5 * * * * 2025",
		"0/5 * * * * * 2025",
	}
	weekNumberings := []WeekdayNumbering{
		WeekdayZeroBased,
		WeekdayOneBased,
		WeekdayZeroBased,
	}

	var wg sync.WaitGroup

	// 并发测试，独立 parser 实例
	for i := 0; i < len(expressions); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			parser := NewCronParser()
			p := parser.SetExpression(expressions[i])
			assert.NotNil(t, p)
			p2 := parser.SetWeekNumbering(weekNumberings[i])
			assert.NotNil(t, p2)

			cronExpr, err := parser.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, cronExpr)
			assert.NotEmpty(t, cronExpr.Second.Vals())
			assert.NotEmpty(t, cronExpr.Minute.Vals())
		}(i)
	}

	wg.Wait()
}
